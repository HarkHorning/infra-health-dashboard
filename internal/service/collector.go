package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/HarkHorning/infra-health-dashboard/internal/azure"
	"github.com/HarkHorning/infra-health-dashboard/internal/models"
	"github.com/HarkHorning/infra-health-dashboard/internal/repository"
)

// CollectorConfig holds configuration for the metrics collector
type CollectorConfig struct {
	// PollingInterval is how often to collect metrics (default: 5 minutes)
	PollingInterval time.Duration
	// Enabled controls whether the collector is active
	Enabled bool
}

// DefaultCollectorConfig returns sensible default configuration
func DefaultCollectorConfig() CollectorConfig {
	return CollectorConfig{
		PollingInterval: 5 * time.Minute,
		Enabled:         true,
	}
}

// CollectorConfigFromEnv loads collector configuration from environment variables
func CollectorConfigFromEnv() CollectorConfig {
	cfg := DefaultCollectorConfig()

	if interval := os.Getenv("POLLING_INTERVAL"); interval != "" {
		if seconds, err := strconv.Atoi(interval); err == nil {
			cfg.PollingInterval = time.Duration(seconds) * time.Second
		}
	}

	if enabled := os.Getenv("COLLECTOR_ENABLED"); enabled != "" {
		cfg.Enabled = enabled == "true" || enabled == "1"
	}

	return cfg
}

// Collector handles periodic metric collection from Azure
type Collector struct {
	azureClient  *azure.Client
	resourceRepo *repository.ResourceRepository
	metricRepo   *repository.MetricRepository
	config       CollectorConfig

	stopCh chan struct{}
	wg     sync.WaitGroup
	mu     sync.Mutex
	running bool
}

// NewCollector creates a new metrics collector
func NewCollector(
	azureClient *azure.Client,
	resourceRepo *repository.ResourceRepository,
	metricRepo *repository.MetricRepository,
	config CollectorConfig,
) *Collector {
	return &Collector{
		azureClient:  azureClient,
		resourceRepo: resourceRepo,
		metricRepo:   metricRepo,
		config:       config,
		stopCh:       make(chan struct{}),
	}
}

// Start begins the periodic metric collection
func (c *Collector) Start() error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("collector is already running")
	}
	c.running = true
	c.mu.Unlock()

	if !c.config.Enabled {
		log.Println("Collector is disabled, not starting")
		return nil
	}

	log.Printf("Starting metrics collector with polling interval: %v", c.config.PollingInterval)

	c.wg.Add(1)
	go c.collectLoop()

	return nil
}

// Stop gracefully stops the collector
func (c *Collector) Stop() {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return
	}
	c.running = false
	c.mu.Unlock()

	close(c.stopCh)
	c.wg.Wait()
	log.Println("Metrics collector stopped")
}

// collectLoop runs the periodic collection
func (c *Collector) collectLoop() {
	defer c.wg.Done()

	// Collect immediately on start
	c.collectAll()

	ticker := time.NewTicker(c.config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.collectAll()
		case <-c.stopCh:
			return
		}
	}
}

// collectAll collects metrics for all registered Azure resources
func (c *Collector) collectAll() {
	fmt.Println(">>> collectAll() called")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("Starting metrics collection cycle...")

	// Get all Azure VM resources from our database
	resources, err := c.resourceRepo.ListByType("azure_vm")
	if err != nil {
		log.Printf("Error listing resources: %v", err)
		fmt.Printf(">>> Error listing resources: %v\n", err)
		return
	}

	fmt.Printf(">>> Found %d resources in database\n", len(resources))

	if len(resources) == 0 {
		log.Println("No Azure VM resources registered for monitoring")
		return
	}

	log.Printf("Collecting metrics for %d resources", len(resources))

	for _, resource := range resources {
		if err := c.collectForResource(ctx, resource); err != nil {
			log.Printf("Error collecting metrics for resource %s: %v", resource.Name, err)
			continue
		}
	}

	log.Println("Metrics collection cycle completed")
}

// collectForResource collects metrics for a single resource
func (c *Collector) collectForResource(ctx context.Context, resource models.Resource) error {
	// The resource_id should be the full Azure resource URI
	// e.g., /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Compute/virtualMachines/{name}
	resourceURI := resource.ResourceID

	fmt.Printf(">>> Collecting metrics for: %s (URI: %s)\n", resource.Name, resourceURI)
	log.Printf("Collecting metrics for: %s", resource.Name)

	// Fetch VM metrics from Azure
	metricValues, err := c.azureClient.FetchVMMetrics(ctx, resourceURI)
	if err != nil {
		fmt.Printf(">>> Error fetching metrics: %v\n", err)
		return fmt.Errorf("failed to fetch metrics: %w", err)
	}

	fmt.Printf(">>> Got %d metric values from Azure\n", len(metricValues))

	if len(metricValues) == 0 {
		log.Printf("No metrics returned for resource: %s", resource.Name)
		return nil
	}

	// Convert to metric requests and store
	var metricRequests []models.CreateMetricRequest
	for _, mv := range metricValues {
		// Use normalized metric name
		normalizedName := azure.GetNormalizedMetricName(mv.MetricName)
		unit := azure.GetMetricUnit(mv.MetricName)

		metricRequests = append(metricRequests, models.CreateMetricRequest{
			MetricName:  normalizedName,
			MetricValue: mv.Value,
			Unit:        unit,
			RecordedAt:  mv.Timestamp.Format(time.RFC3339),
		})
	}

	// Store metrics in database
	if err := c.metricRepo.InsertBatch(resource.ID, metricRequests); err != nil {
		return fmt.Errorf("failed to store metrics: %w", err)
	}

	log.Printf("Stored %d metrics for resource: %s", len(metricRequests), resource.Name)
	return nil
}

// CollectOnce performs a single collection cycle (useful for testing or manual triggers)
func (c *Collector) CollectOnce() {
	c.collectAll()
}

// CollectOnceWithResult performs a collection cycle and returns debug info
func (c *Collector) CollectOnceWithResult() map[string]interface{} {
	result := make(map[string]interface{})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Get all Azure VM resources from our database
	resources, err := c.resourceRepo.ListByType("azure_vm")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to list resources: %v", err)
		return result
	}

	result["resources_found"] = len(resources)

	if len(resources) == 0 {
		result["message"] = "No Azure VM resources registered for monitoring"
		return result
	}

	var metricsCollected []map[string]interface{}

	for _, resource := range resources {
		resourceResult := map[string]interface{}{
			"resource_name": resource.Name,
			"resource_id":   resource.ResourceID,
		}

		// Fetch VM metrics from Azure
		metricValues, err := c.azureClient.FetchVMMetrics(ctx, resource.ResourceID)
		if err != nil {
			resourceResult["error"] = err.Error()
			metricsCollected = append(metricsCollected, resourceResult)
			continue
		}

		resourceResult["metrics_from_azure"] = len(metricValues)

		if len(metricValues) == 0 {
			resourceResult["message"] = "No metrics returned from Azure"
			metricsCollected = append(metricsCollected, resourceResult)
			continue
		}

		// Convert to metric requests and store
		var metricRequests []models.CreateMetricRequest
		for _, mv := range metricValues {
			normalizedName := azure.GetNormalizedMetricName(mv.MetricName)
			unit := azure.GetMetricUnit(mv.MetricName)

			metricRequests = append(metricRequests, models.CreateMetricRequest{
				MetricName:  normalizedName,
				MetricValue: mv.Value,
				Unit:        unit,
				RecordedAt:  mv.Timestamp.Format(time.RFC3339),
			})
		}

		// Store metrics in database
		if err := c.metricRepo.InsertBatch(resource.ID, metricRequests); err != nil {
			resourceResult["store_error"] = err.Error()
		} else {
			resourceResult["metrics_stored"] = len(metricRequests)
		}

		metricsCollected = append(metricsCollected, resourceResult)
	}

	result["collection_results"] = metricsCollected
	return result
}

// DiscoverAndRegisterVMs discovers Azure VMs and registers them in the database
func (c *Collector) DiscoverAndRegisterVMs(ctx context.Context) ([]models.Resource, error) {
	fmt.Println("=== DISCOVER START ===")
	fmt.Printf("Subscription ID: %s\n", c.azureClient.GetSubscriptionID())
	log.Println("Discovering Azure VMs using Compute client...")

	// List all VMs in the subscription using the Compute client directly
	vms, err := c.azureClient.ListVMs(ctx)
	if err != nil {
		fmt.Printf("=== ERROR: %v ===\n", err)
		return nil, fmt.Errorf("failed to list VMs: %w", err)
	}

	fmt.Printf("=== FOUND %d VMs ===\n", len(vms))
	log.Printf("Found %d VMs from Azure", len(vms))

	var registered []models.Resource

	for _, vm := range vms {
		// Check if already registered
		existing, _ := c.resourceRepo.GetByResourceID(vm.ID)
		if existing != nil {
			log.Printf("VM already registered: %s", vm.Name)
			registered = append(registered, *existing)
			continue
		}

		// Create resource request
		req := models.CreateResourceRequest{
			ResourceID:   vm.ID,
			ResourceType: "azure_vm",
			Name:         vm.Name,
			Region:       vm.Location,
		}

		// Register in database
		resource, err := c.resourceRepo.Create(req)
		if err != nil {
			log.Printf("Failed to register VM %s: %v", vm.Name, err)
			continue
		}

		log.Printf("Registered VM: %s", vm.Name)
		registered = append(registered, *resource)
	}

	log.Printf("Discovered and registered %d VMs", len(registered))
	return registered, nil
}

// IsRunning returns whether the collector is currently running
func (c *Collector) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

// GetSubscriptionID returns the Azure subscription ID being used
func (c *Collector) GetSubscriptionID() string {
	return c.azureClient.GetSubscriptionID()
}
