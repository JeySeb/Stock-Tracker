package repositories

import (
	"context"
	stockModel "stock-tracker/internal/domain/stocks/model"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
	"time"

	"github.com/google/uuid"
)

type StockRepository interface {
	//CRUD operations
	Create(ctx context.Context, stock *stockModel.Stock) error
	GetByID(ctx context.Context, id uuid.UUID) (*stockModel.Stock, error)
	Update(ctx context.Context, stock *stockModel.Stock) error
	Delete(ctx context.Context, id uuid.UUID) error

	//Query operations
	GetByTicker(ctx context.Context, ticker string) ([]*stockModel.Stock, error)
	GetLatestByTicker(ctx context.Context, ticker string) (*stockModel.Stock, error)
	GetAll(ctx context.Context, filters stockValidation.StockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error)
	GetAllWithEnhancedFilters(ctx context.Context, filters stockValidation.EnhancedStockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error)
	GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*stockModel.Stock, error)

	//Batch operations
	BulkCreate(ctx context.Context, stocks []*stockModel.Stock) error
	BulkUpdate(ctx context.Context, stocks []*stockModel.Stock) error

	//Analytics queries
	GetTopMoversByTarget(ctx context.Context, limit int) ([]*stockModel.Stock, error)
	GetUniqueTickersCount(ctx context.Context) (int, error)
	GetBrokerageStats(ctx context.Context) ([]BrokerageStats, error)
}

type BrokerageStats struct {
	Brokerage string  `json:"brokerage"`
	Count     int     `json:"count"`
	AvgScore  float64 `json:"avg_score"`
}
