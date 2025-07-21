package usecases

import (
	"context"
	"fmt"
	"time"

	"stock-tracker/internal/domain/repositories"
	"stock-tracker/internal/domain/valueObjects"
	"stock-tracker/pkg/logger"
)

type StockQueryUseCase struct {
	stockRepo  repositories.StockRepository
	brokerRepo repositories.BrokerRepository
	logger     logger.Logger
}

func NewStockQueryUseCase(
	stockRepo repositories.StockRepository,
	brokerRepo repositories.BrokerRepository,
	logger logger.Logger,
) StockUseCase {
	return &StockQueryUseCase{
		stockRepo:  stockRepo,
		brokerRepo: brokerRepo,
		logger:     logger,
	}
}

// GetStocks returns stocks with pagination
func (uc *StockQueryUseCase) GetStocks(ctx context.Context, filters valueObjects.StockFilters) (interface{}, *valueObjects.Pagination, error) {
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
	filters := valueObjects.StockFilters{}
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
