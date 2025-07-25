package repositories

import (
	"context"
	"stock-tracker/internal/domain/model"
	"stock-tracker/internal/domain/authentication/validation"
	"time"

	"github.com/google/uuid"
)

type StockRepository interface {
	//CRUD operations
	Create(ctx context.Context, stock *model.Stock) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Stock, error)
	Update(ctx context.Context, stock *model.Stock) error
	Delete(ctx context.Context, id uuid.UUID) error

	//Query operations
	GetByTicker(ctx context.Context, ticker string) ([]*model.Stock, error)
	GetLatestByTicker(ctx context.Context, ticker string) (*model.Stock, error)
	GetAll(ctx context.Context, filters validation.StockFilters) ([]*model.Stock, *validation.Pagination, error)
	GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*model.Stock, error)

	//Batch operations
	BulkCreate(ctx context.Context, stocks []*model.Stock) error
	BulkUpdate(ctx context.Context, stocks []*model.Stock) error

	//Analytics queries
	GetTopMoversByTarget(ctx context.Context, limit int) ([]*model.Stock, error)
	GetUniqueTickersCount(ctx context.Context) (int, error)
	GetBrokerageStats(ctx context.Context) ([]BrokerageStats, error)
}

type BrokerageStats struct {
	Brokerage string  `json:"brokerage"`
	Count     int     `json:"count"`
	AvgScore  float64 `json:"avg_score"`
}
