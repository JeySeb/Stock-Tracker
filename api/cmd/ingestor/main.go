package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	marketDataUsecases "stock-tracker/internal/domain/marketdata/usecase"
	stockUsecases "stock-tracker/internal/domain/stocks/usecase"
	"stock-tracker/internal/infrastructure/clients"
	"stock-tracker/internal/infrastructure/config"
	"stock-tracker/internal/infrastructure/database"
	"stock-tracker/internal/infrastructure/external"
	"stock-tracker/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	// CLI flags to control ingestion modes
	onlyMarketData := flag.Bool("only-market-data", false, "Run only market data ingestion")
	onlyStocks := flag.Bool("only-stocks", false, "Run only stock ingestion")
	runOnce := flag.Bool("run-once", false, "Run the selected ingestion(s) once and exit (no scheduler)")
	flag.Parse()

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

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}()

	// Initialize repositories
	stockRepo := database.NewStockRepository(db.GetPool(), logger)
	brokerRepo := database.NewBrokerRepository(db.GetPool())
	marketDataRepo := database.NewMarketDataRepository(db.GetPool(), logger)

	// Initialize external clients
	stockAPIClient := clients.NewStockAPIClient(cfg.StockAPIURL, cfg.StockAPIKey, logger)
	yahooClient := external.NewYahooFinanceClient(logger)

	// Initialize the use cases
	stockIngestionUseCase := stockUsecases.NewStockIngestionUseCase(stockRepo, brokerRepo, stockAPIClient, logger)
	marketDataIngestionUseCase := marketDataUsecases.NewMarketDataIngestionUseCase(marketDataRepo, yahooClient, logger)
	// Determine which ingestions to run
	if *onlyMarketData && *onlyStocks {
		log.Fatal("Cannot use --only-market-data and --only-stocks together")
	}
	doStocks := !*onlyMarketData
	doMarket := !*onlyStocks

	// Initialize the cron job (cron scheduler)
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))))

	// First run immediately based on flags
	ctx := context.Background()
	if doStocks {
		if err := stockIngestionUseCase.IngestStocks(ctx); err != nil {
			logger.Error("Initial stock ingestion failed", "error", err)
		}
	}

	// First market data ingestion
	if doMarket {
		if err := marketDataIngestionUseCase.IngestMarketData(ctx); err != nil {
			logger.Error("Initial market data ingestion failed", "error", err)
		}
	}

	// If run-once, skip scheduler and exit
	if *runOnce {
		logger.Info("Run-once mode enabled; skipping scheduler and exiting")
		return
	}

	// Add the cron job to schedule the stock ingestion every hour
	if doStocks {
		_, err = c.AddFunc("0 * * * *", func() {
			ctx := context.Background()
			if err := stockIngestionUseCase.IngestStocks(ctx); err != nil {
				logger.Error("Stock ingestion job failed", "error", err)
			}
		})
		if err != nil {
			log.Fatal("Failed to schedule stock ingestion job", "error", err)
		}
	}

	// Add the cron job to schedule the market data ingestion every hour (at 15 minutes past)
	if doMarket {
		_, err = c.AddFunc("15 * * * *", func() {
			ctx := context.Background()
			if err := marketDataIngestionUseCase.IngestMarketData(ctx); err != nil {
				logger.Error("Market data ingestion job failed", "error", err)
			}
		})
		if err != nil {
			log.Fatal("Failed to schedule market data ingestion job", "error", err)
		}
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
