package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"stock-tracker/internal/domain/repositories"
	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/infrastructure/clients"
	"stock-tracker/internal/infrastructure/config"
	"stock-tracker/internal/infrastructure/database"
	"stock-tracker/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	//Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	//Initialize the logger
	logger := logger.NewSimpleLogger()
	logger.Info("Starting stock ingestion system")

	//Initialize the database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	// Initialize repositories
	stockRepo := database.NewStockRepository(db.GetPool(), logger)
	brokerRepo := repositories.NewBrokerRepository(db.GetPool())

	// Initialize external clients
	stockAPIClient := clients.NewStockAPIClient(cfg.StockAPIURL, cfg.StockAPIKey, logger)

	// Initialize the use case
	stockIngestionUseCase := usecases.NewStockIngestionUseCase(stockRepo, brokerRepo, stockAPIClient, logger)
	// Initialize the cron job (cron scheduler)
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// First run immediately
	ctx := context.Background()
	if err := stockIngestionUseCase.IngestStocks(ctx); err != nil {
		logger.Error("Initial ingestion failed", "error", err)
	}

	// Add the cron job to schedule the ingestion every hour
	_, err = c.AddFunc("0 * * * *", func() {
		ctx := context.Background()
		if err := stockIngestionUseCase.IngestStocks(ctx); err != nil {
			logger.Error("Ingestion job failed", "error", err)
		}
	})
	if err != nil {
		log.Fatal("Failed to schedule ingestion job", "error", err)
	}

	// Start the cron scheduler
	c.Start()

	// Wait for a signal to stop the application
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop the cron scheduler
	logger.Info("Stopping stock ingestion system")
	c.Stop()

	// Close the database connection
	if err := db.Close(); err != nil {
		logger.Error("Failed to close database connection", "error", err)
	}

	logger.Info("Stock ingestion system stopped")
}
