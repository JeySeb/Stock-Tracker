package repositories

import (
	"context"
	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/valueObjects"
	"time"

	"github.com/google/uuid"
)

type StockRepository interface {
	//CRUD operations
	Create(ctx context.Context, stock *entities.Stock) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Stock, error)
	Update(ctx context.Context, stock *entities.Stock) error
	Delete(ctx context.Context, id uuid.UUID) error

	//Query operations
	GetByTicker(ctx context.Context, ticker string) ([]*entities.Stock, error)
	GetLatestByTicker(ctx context.Context, ticker string) (*entities.Stock, error)
	GetAll(ctx context.Context, filters valueObjects.StockFilters) ([]*entities.Stock, *valueObjects.Pagination, error)
	GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*entities.Stock, error)

	//Batch operations
	BulkCreate(ctx context.Context, stocks []*entities.Stock) error
	BulkUpdate(ctx context.Context, stocks []*entities.Stock) error

	//Analytics queries
	GetTopMoversByTarget(ctx context.Context, limit int) ([]*entities.Stock, error)
	GetUniqueTickersCount(ctx context.Context) (int, error)
	GetBrokerageStats(ctx context.Context) ([]BrokerageStats, error)
}

type BrokerageStats struct {
	Brokerage string  `json:"brokerage"`
	Count     int     `json:"count"`
	AvgScore  float64 `json:"avg_score"`
}
