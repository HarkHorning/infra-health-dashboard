package azure

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

// Client wraps Azure SDK clients for metrics and resource operations
type Client struct {
	credential     *azidentity.DefaultAzureCredential
	metricsClient  *armmonitor.MetricsClient
	resourceClient *armresources.Client
	vmClient       *armcompute.VirtualMachinesClient
	subscriptionID string
}

// Config holds Azure configuration
type Config struct {
	SubscriptionID string
	// ClientID, ClientSecret, TenantID are read from environment by DefaultAzureCredential
}

// ConfigFromEnv loads Azure configuration from environment variables
func ConfigFromEnv() Config {
	return Config{
		SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
	}
}

// NewClient creates a new Azure client with DefaultAzureCredential
// DefaultAzureCredential will try multiple authentication methods in order:
// 1. Environment variables (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)
// 2. Managed Identity (when running in Azure)
// 3. Azure CLI credentials
// 4. Azure Developer CLI credentials
func NewClient(cfg Config) (*Client, error) {
	if cfg.SubscriptionID == "" {
		return nil, fmt.Errorf("AZURE_SUBSCRIPTION_ID environment variable is required")
	}

	// Create credential using DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credential: %w", err)
	}

	// Create metrics client
	metricsClient, err := armmonitor.NewMetricsClient(cfg.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	// Create resource client for listing resources
	resourceClient, err := armresources.NewClient(cfg.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource client: %w", err)
	}

	// Create VM client for listing VMs directly
	vmClient, err := armcompute.NewVirtualMachinesClient(cfg.SubscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create VM client: %w", err)
	}

	return &Client{
		credential:     cred,
		metricsClient:  metricsClient,
		resourceClient: resourceClient,
		vmClient:       vmClient,
		subscriptionID: cfg.SubscriptionID,
	}, nil
}

// NewClientFromEnv creates a new Azure client using environment variables
func NewClientFromEnv() (*Client, error) {
	cfg := ConfigFromEnv()
	return NewClient(cfg)
}

// GetMetricsClient returns the underlying metrics client
func (c *Client) GetMetricsClient() *armmonitor.MetricsClient {
	return c.metricsClient
}

// GetResourceClient returns the underlying resource client
func (c *Client) GetResourceClient() *armresources.Client {
	return c.resourceClient
}

// GetSubscriptionID returns the subscription ID
func (c *Client) GetSubscriptionID() string {
	return c.subscriptionID
}

// ListResources lists all resources in the subscription
func (c *Client) ListResources(ctx context.Context) ([]*armresources.GenericResourceExpanded, error) {
	var resources []*armresources.GenericResourceExpanded

	pager := c.resourceClient.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list resources: %w", err)
		}
		resources = append(resources, page.Value...)
	}

	return resources, nil
}

// ListResourcesByType lists resources of a specific type (e.g., "Microsoft.Compute/virtualMachines")
func (c *Client) ListResourcesByType(ctx context.Context, resourceType string) ([]*armresources.GenericResourceExpanded, error) {
	var resources []*armresources.GenericResourceExpanded

	filter := fmt.Sprintf("resourceType eq '%s'", resourceType)
	fmt.Printf("DEBUG: Listing resources with filter: %s\n", filter)
	fmt.Printf("DEBUG: Subscription ID: %s\n", c.subscriptionID)

	pager := c.resourceClient.NewListPager(&armresources.ClientListOptions{
		Filter: &filter,
	})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list resources by type: %w", err)
		}
		fmt.Printf("DEBUG: Got page with %d resources\n", len(page.Value))
		resources = append(resources, page.Value...)
	}

	fmt.Printf("DEBUG: Total resources found: %d\n", len(resources))
	return resources, nil
}

// VMInfo holds basic VM information
type VMInfo struct {
	ID       string
	Name     string
	Location string
}

// ListVMs lists all VMs in the subscription using the Compute client directly
func (c *Client) ListVMs(ctx context.Context) ([]VMInfo, error) {
	var vms []VMInfo

	fmt.Printf("DEBUG: Listing VMs using Compute client for subscription: %s\n", c.subscriptionID)

	pager := c.vmClient.NewListAllPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list VMs: %w", err)
		}
		fmt.Printf("DEBUG: Got page with %d VMs\n", len(page.Value))

		for _, vm := range page.Value {
			if vm.ID == nil || vm.Name == nil {
				continue
			}
			location := ""
			if vm.Location != nil {
				location = *vm.Location
			}
			vms = append(vms, VMInfo{
				ID:       *vm.ID,
				Name:     *vm.Name,
				Location: location,
			})
			fmt.Printf("DEBUG: Found VM: %s (ID: %s)\n", *vm.Name, *vm.ID)
		}
	}

	fmt.Printf("DEBUG: Total VMs found: %d\n", len(vms))
	return vms, nil
}
