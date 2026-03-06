package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HarkHorning/infra-health-dashboard/internal/models"
	"github.com/HarkHorning/infra-health-dashboard/internal/repository"
	"github.com/HarkHorning/infra-health-dashboard/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler holds the dependencies for HTTP handlers
type Handler struct {
	resourceRepo *repository.ResourceRepository
	metricRepo   *repository.MetricRepository
	collector    *service.Collector
}

// NewHandler creates a new Handler with the given repositories
func NewHandler(resourceRepo *repository.ResourceRepository, metricRepo *repository.MetricRepository, collector *service.Collector) *Handler {
	return &Handler{
		resourceRepo: resourceRepo,
		metricRepo:   metricRepo,
		collector:    collector,
	}
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ListResources handles GET /api/v1/resources
func (h *Handler) ListResources(c *gin.Context) {
	resourceType := c.Query("type")

	var resources []models.Resource
	var err error

	if resourceType != "" {
		resources, err = h.resourceRepo.ListByType(resourceType)
	} else {
		resources, err = h.resourceRepo.List()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if resources == nil {
		resources = []models.Resource{}
	}

	c.JSON(http.StatusOK, resources)
}

// GetResource handles GET /api/v1/resources/:id
func (h *Handler) GetResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	resource, err := h.resourceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resource)
}

// CreateResource handles POST /api/v1/resources
func (h *Handler) CreateResource(c *gin.Context) {
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource, err := h.resourceRepo.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resource)
}

// UpdateResource handles PUT /api/v1/resources/:id
func (h *Handler) UpdateResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resource, err := h.resourceRepo.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resource)
}

// DeleteResource handles DELETE /api/v1/resources/:id
func (h *Handler) DeleteResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	if err := h.resourceRepo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "resource deleted"})
}

// GetResourceMetrics handles GET /api/v1/resources/:id/metrics
func (h *Handler) GetResourceMetrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	// Parse query parameters
	query := models.MetricQuery{
		MetricName: c.Query("metric_name"),
	}

	if from := c.Query("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' time format (use RFC3339)"})
			return
		}
		query.From = t
	}

	if to := c.Query("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' time format (use RFC3339)"})
			return
		}
		query.To = t
	}

	metrics, err := h.metricRepo.GetByResourceID(id, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if metrics == nil {
		metrics = []models.Metric{}
	}

	c.JSON(http.StatusOK, metrics)
}

// CreateResourceMetric handles POST /api/v1/resources/:id/metrics
func (h *Handler) CreateResourceMetric(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	// Verify resource exists
	_, err = h.resourceRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return
	}

	var req models.CreateMetricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metric, err := h.metricRepo.Insert(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, metric)
}

// GetResourceLatestMetrics handles GET /api/v1/resources/:id/metrics/latest
func (h *Handler) GetResourceLatestMetrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	metrics, err := h.metricRepo.GetLatestByResourceID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if metrics == nil {
		metrics = []models.Metric{}
	}

	c.JSON(http.StatusOK, metrics)
}

// GetMetricsSummary handles GET /api/v1/metrics/summary
func (h *Handler) GetMetricsSummary(c *gin.Context) {
	resources, err := h.resourceRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type ResourceSummary struct {
		Resource      *models.Resource `json:"resource"`
		LatestMetrics []models.Metric  `json:"latest_metrics"`
	}

	var summaries []ResourceSummary
	for _, res := range resources {
		metrics, err := h.metricRepo.GetLatestByResourceID(res.ID)
		if err != nil {
			continue
		}
		summaries = append(summaries, ResourceSummary{
			Resource:      &res,
			LatestMetrics: metrics,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total_resources": len(resources),
		"summaries":       summaries,
	})
}

// DiscoverResources handles POST /api/v1/discover
// Discovers Azure VMs and registers them in the database
func (h *Handler) DiscoverResources(c *gin.Context) {
	fmt.Println(">>> HANDLER: DiscoverResources called")

	if h.collector == nil {
		fmt.Println(">>> HANDLER: collector is nil!")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Azure collector not initialized - check Azure credentials",
		})
		return
	}

	fmt.Println(">>> HANDLER: calling DiscoverAndRegisterVMs...")
	resources, err := h.collector.DiscoverAndRegisterVMs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to discover resources",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "discovery completed",
		"resources_found":      len(resources),
		"registered_resources": resources,
	})
}

// DebugAzure handles GET /api/v1/debug/azure
// Shows Azure connection info for debugging
func (h *Handler) DebugAzure(c *gin.Context) {
	if h.collector == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":           "Azure collector not initialized",
			"possible_causes": []string{
				"AZURE_CLIENT_ID not set or invalid",
				"AZURE_CLIENT_SECRET not set or invalid",
				"AZURE_TENANT_ID not set or invalid",
				"AZURE_SUBSCRIPTION_ID not set or invalid",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "Azure client initialized",
		"subscription_id": h.collector.GetSubscriptionID(),
		"message":         "Try POST /api/v1/discover to find VMs",
	})
}

// TriggerCollection handles POST /api/v1/collect
// Manually triggers a metrics collection cycle
func (h *Handler) TriggerCollection(c *gin.Context) {
	if h.collector == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Collector not initialized",
		})
		return
	}

	// Collect and return result directly
	result := h.collector.CollectOnceWithResult()

	c.JSON(http.StatusOK, result)
}
