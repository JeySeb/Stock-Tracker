package repositories

import (
	"context"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
)

// MarketDataRepository defines the interface for market data persistence
type MarketDataRepository interface {
	// Save market data
	SaveMarketData(ctx context.Context, marketData *model.MarketData) error
	
	// Save ingestion log
	SaveIngestionLog(ctx context.Context, log *model.MarketDataIngestionLog) error
	
	// Update ingestion log
	UpdateIngestionLog(ctx context.Context, log *model.MarketDataIngestionLog) error
	
	// Get unique tickers from stocks table
	GetUniqueTickers(ctx context.Context) ([]string, error)
	
	// Check if market data exists for ticker and timestamp
	ExistsMarketData(ctx context.Context, ticker string, dataSource model.DataSource, dataTimestamp time.Time) (bool, error)
	
	// Get latest market data for a ticker
	GetLatestMarketData(ctx context.Context, ticker string) (*model.MarketData, error)
	
	// Get market data stats
	GetMarketDataStats(ctx context.Context, days int) (map[string]interface{}, error)
} 