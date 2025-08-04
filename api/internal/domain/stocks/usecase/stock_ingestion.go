package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	stockModel "stock-tracker/internal/domain/stocks/model"
	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
	"stock-tracker/internal/infrastructure/clients"
	"stock-tracker/pkg/logger"
)

type StockIngestionUseCase struct {
	stockRepo   stockRepos.StockRepository
	brokerRepo  stockRepos.BrokerRepository
	apiClient   clients.StockAPIClient
	logger      logger.Logger
	batchSize   int
	workerCount int
}

func NewStockIngestionUseCase(
	stockRepo stockRepos.StockRepository,
	brokerRepo stockRepos.BrokerRepository,
	apiClient clients.StockAPIClient,
	logger logger.Logger,
) *StockIngestionUseCase {
	return &StockIngestionUseCase{
		stockRepo:   stockRepo,
		brokerRepo:  brokerRepo,
		apiClient:   apiClient,
		logger:      logger,
		batchSize:   100,
		workerCount: 5,
	}
}

func (uc *StockIngestionUseCase) IngestStocks(ctx context.Context) error {

	batchID := uuid.New().String()
	startTime := time.Now()

	uc.logger.Info("Starting stock ingestion batch", "batchID", batchID, "startTime", startTime)

	stocks, err := uc.apiClient.FetchAllStocks(ctx)
	if err != nil {
		uc.logger.Error("Failed to fetch stocks from API", "error", err)
		return fmt.Errorf("failed to fetch stocks: %w", err)
	}

	if len(stocks) == 0 {
		uc.logger.Info("No stocks found in API", "batchID", batchID)
		return nil
	}

	uc.logger.Info("Fetched stocks from API", "batchID", batchID, "count", len(stocks))

	// Enrich the stocks with the Brokes IDs

	if err := uc.enrichWithBrokerInfo(ctx, stocks); err != nil {
		uc.logger.Error("Failed to enrich stocks with brokers", "error", err)
		return fmt.Errorf("failed to enrich stocks with brokers: %w", err)
	}

	uc.logger.Info("Enriched stocks with brokers", "batchID", batchID, "count", len(stocks))

	eg, ctx := errgroup.WithContext(ctx)

	//Process Stocks in batches using worker pool
	if err := uc.processStocksInBatches(ctx, eg, stocks); err != nil {
		uc.logger.Error("Failed to process stocks in batches", "error", err)
		return fmt.Errorf("failed to process stocks in batches: %w", err)
	}

	if err := eg.Wait(); err != nil {
		uc.logger.Error("Error during stock ingestion", "error", err)
		return fmt.Errorf("error during stock ingestion: %w", err)
	}
	duration := time.Since(startTime)
	uc.logger.Info("Stock ingestion batch completed, for a total of", len(stocks), "stocks", "batchID", batchID, "duration", duration)
	return nil
}

func (uc *StockIngestionUseCase) enrichWithBrokerInfo(ctx context.Context, stocks []*stockModel.Stock) error {
	// Get all existing brokers
	brokers, err := uc.brokerRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get brokers: %w", err)
	}

	brokerMap := make(map[string]*stockModel.Broker)
	for _, broker := range brokers {
		brokerMap[broker.Name] = broker
	}

	// Create missing brokers and assign IDs
	var newBrokers []*stockModel.Broker
	for _, stock := range stocks {
		if broker, exists := brokerMap[stock.Brokerage]; exists {
			stock.BrokerID = broker.ID
		} else {
			// Create new broker with default credibility score
			newBroker := stockModel.NewBroker(stock.Brokerage, 0.60)
			newBrokers = append(newBrokers, newBroker)
			brokerMap[stock.Brokerage] = newBroker
			stock.BrokerID = newBroker.ID
		}
	}

	// Save new brokers
	for _, broker := range newBrokers {
		if err := uc.brokerRepo.Create(ctx, broker); err != nil {
			uc.logger.Warn("Failed to create broker", "name", broker.Name, "error", err)
		}
	}

	return nil
}

func (uc *StockIngestionUseCase) processStocksInBatches(ctx context.Context, eg *errgroup.Group, stocks []*stockModel.Stock) error {
	batches := uc.createBatches(stocks)

	for i, batch := range batches {
		batchNum := i
		batch := batch // capture loop variable

		eg.Go(func() error {
			return uc.processBatch(ctx, batch, batchNum)
		})
	}

	return nil
}

// GetStats returns basic statistics about the stock data
func (uc *StockIngestionUseCase) GetStats(ctx context.Context) (interface{}, error) {
	uc.logger.Info("Getting stock statistics")

	// Get a count by querying with an empty filter
	filters := stockValidation.StockFilters{}
	filters.SetDefaults()
	filters.Limit = 1 // We only need the count, not the actual data

	_, pagination, err := uc.stockRepo.GetAll(ctx, filters)
	if err != nil {
		uc.logger.Error("Failed to get total stocks count", "error", err)
		return nil, fmt.Errorf("failed to retrieve statistics: %w", err)
	}

	stats := map[string]interface{}{
		"total_stocks": pagination.TotalItems,
		"last_updated": time.Now(),
	}

	uc.logger.Info("Successfully retrieved statistics", "total_stocks", pagination.TotalItems)
	return stats, nil
}

// GetStocks returns stocks with pagination
func (uc *StockIngestionUseCase) GetStocks(ctx context.Context, filters stockValidation.StockFilters) (interface{}, *stockValidation.Pagination, error) {
	uc.logger.Info("Getting stocks with filters", "filters", filters)

	stocks, pagination, err := uc.stockRepo.GetAll(ctx, filters)
	if err != nil {
		uc.logger.Error("Failed to get stocks from repository", "error", err)
		return nil, nil, fmt.Errorf("failed to retrieve stocks: %w", err)
	}

	uc.logger.Info("Successfully retrieved stocks", "count", len(stocks), "total", pagination.TotalItems)
	return stocks, pagination, nil
}

// GetStocksByTicker returns stocks for a specific ticker
func (uc *StockIngestionUseCase) GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error) {
	uc.logger.Info("Getting stocks by ticker", "ticker", ticker)

	stocks, err := uc.stockRepo.GetByTicker(ctx, ticker)
	if err != nil {
		uc.logger.Error("Failed to get stocks by ticker", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to retrieve stocks for ticker %s: %w", ticker, err)
	}

	uc.logger.Info("Successfully retrieved stocks by ticker", "ticker", ticker, "count", len(stocks))
	return stocks, nil
}

// GetStocksWithEnhancedFilters returns stocks with enhanced filtering capabilities
func (uc *StockIngestionUseCase) GetStocksWithEnhancedFilters(ctx context.Context, filters stockValidation.EnhancedStockFilters) (interface{}, *stockValidation.Pagination, error) {
	uc.logger.Info("Getting stocks with enhanced filters", "filters", filters)

	// Validate filters
	if err := filters.Validate(); err != nil {
		uc.logger.Error("Invalid enhanced filters", "error", err)
		return nil, nil, fmt.Errorf("invalid filters: %w", err)
	}

	stocks, pagination, err := uc.stockRepo.GetAllWithEnhancedFilters(ctx, filters)
	if err != nil {
		uc.logger.Error("Failed to get stocks from repository with enhanced filters", "error", err)
		return nil, nil, fmt.Errorf("failed to retrieve stocks: %w", err)
	}

	uc.logger.Info("Successfully retrieved stocks with enhanced filters", "count", len(stocks), "total", pagination.TotalItems)
	return stocks, pagination, nil
}

func (uc *StockIngestionUseCase) createBatches(stocks []*stockModel.Stock) [][]*stockModel.Stock {
	var batches [][]*stockModel.Stock

	for i := 0; i < len(stocks); i += uc.batchSize {
		end := i + uc.batchSize
		if end > len(stocks) {
			end = len(stocks)
		}
		batches = append(batches, stocks[i:end])
	}

	return batches
}

func (uc *StockIngestionUseCase) processBatch(ctx context.Context, batch []*stockModel.Stock, batchNum int) error {
	uc.logger.Info("Processing batch", "batch_num", batchNum, "size", len(batch))

	startTime := time.Now()
	err := uc.stockRepo.BulkCreate(ctx, batch)
	duration := time.Since(startTime)

	if err != nil {
		uc.logger.Error("Failed to process batch", "batch_num", batchNum, "error", err)
		return fmt.Errorf("failed to process batch %d: %w", batchNum, err)
	}

	uc.logger.Info("Completed batch", "batch_num", batchNum, "duration", duration)
	return nil
}
