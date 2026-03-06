package api

import (
	"github.com/HarkHorning/infra-health-dashboard/internal/repository"
	"github.com/HarkHorning/infra-health-dashboard/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SetupRouter creates and configures the Gin router with all routes
// collector is optional - if nil, discovery endpoints will be disabled
func SetupRouter(db *sqlx.DB, collector *service.Collector) *gin.Engine {
	// Create repositories
	resourceRepo := repository.NewResourceRepository(db)
	metricRepo := repository.NewMetricRepository(db)

	// Create handler with dependencies
	h := NewHandler(resourceRepo, metricRepo, collector)

	// Create router with default middleware (logger and recovery)
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(Logger())
	router.Use(CORS())

	// Health check
	router.GET("/health", h.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Resource routes
		resources := v1.Group("/resources")
		{
			resources.GET("", h.ListResources)
			resources.POST("", h.CreateResource)
			resources.GET("/:id", h.GetResource)
			resources.PUT("/:id", h.UpdateResource)
			resources.DELETE("/:id", h.DeleteResource)

			// Resource metrics routes
			resources.GET("/:id/metrics", h.GetResourceMetrics)
			resources.POST("/:id/metrics", h.CreateResourceMetric)
			resources.GET("/:id/metrics/latest", h.GetResourceLatestMetrics)
		}

		// Metrics summary
		v1.GET("/metrics/summary", h.GetMetricsSummary)

		// Discovery endpoint (only if collector is available)
		if collector != nil {
			v1.POST("/discover", h.DiscoverResources)
			v1.GET("/debug/azure", h.DebugAzure)
			v1.POST("/collect", h.TriggerCollection)
		}
	}

	return router
}
