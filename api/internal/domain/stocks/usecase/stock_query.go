package usecases

import (
	"context"
	"fmt"
	"time"

	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
	"stock-tracker/pkg/logger"
)

type StockQueryUseCase struct {
	stockRepo  stockRepos.StockRepository
	brokerRepo stockRepos.BrokerRepository
	logger     logger.Logger
}

func NewStockQueryUseCase(
	stockRepo stockRepos.StockRepository,
	brokerRepo stockRepos.BrokerRepository,
	logger logger.Logger,
) StockUseCase {
	return &StockQueryUseCase{
		stockRepo:  stockRepo,
		brokerRepo: brokerRepo,
		logger:     logger,
	}
}

// GetStocks returns stocks with pagination
func (uc *StockQueryUseCase) GetStocks(ctx context.Context, filters stockValidation.StockFilters) (interface{}, *stockValidation.Pagination, error) {
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
func (uc *StockQueryUseCase) GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error) {
	uc.logger.Info("Getting stocks by ticker", "ticker", ticker)

	stocks, err := uc.stockRepo.GetByTicker(ctx, ticker)
	if err != nil {
		uc.logger.Error("Failed to get stocks by ticker", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to retrieve stocks for ticker %s: %w", ticker, err)
	}

	uc.logger.Info("Successfully retrieved stocks by ticker", "ticker", ticker, "count", len(stocks))
	return stocks, nil
}

// GetStats returns basic statistics about the stock data
func (uc *StockQueryUseCase) GetStats(ctx context.Context) (interface{}, error) {
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

// GetStocksWithEnhancedFilters returns stocks with enhanced filtering capabilities
func (uc *StockQueryUseCase) GetStocksWithEnhancedFilters(ctx context.Context, filters stockValidation.EnhancedStockFilters) (interface{}, *stockValidation.Pagination, error) {
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
