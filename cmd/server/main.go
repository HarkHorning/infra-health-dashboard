package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HarkHorning/infra-health-dashboard/internal/api"
	"github.com/HarkHorning/infra-health-dashboard/internal/azure"
	"github.com/HarkHorning/infra-health-dashboard/internal/repository"
	"github.com/HarkHorning/infra-health-dashboard/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (won't override existing env vars)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	} else {
		log.Println("Loaded configuration from .env file")
	}

	// Connect to database
	db, err := repository.NewDBFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Run migrations (creates tables if they don't exist)
	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Initialize repositories
	resourceRepo := repository.NewResourceRepository(db)
	metricRepo := repository.NewMetricRepository(db)

	// Initialize Azure client (optional - only if Azure credentials are configured)
	var collector *service.Collector
	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		log.Printf("Azure client not initialized (metrics collection disabled): %v", err)
	} else {
		log.Println("Azure client initialized successfully")

		// Initialize and start the metrics collector
		collectorConfig := service.CollectorConfigFromEnv()
		collector = service.NewCollector(azureClient, resourceRepo, metricRepo, collectorConfig)

		if err := collector.Start(); err != nil {
			log.Printf("Failed to start collector: %v", err)
		}
	}

	// Setup router with all routes (pass collector for discovery endpoint)
	router := api.SetupRouter(db, collector)

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down...")

		// Stop the collector if running
		if collector != nil {
			collector.Stop()
		}

		// Give ongoing requests time to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		<-ctx.Done()

		os.Exit(0)
	}()

	// Start server
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
