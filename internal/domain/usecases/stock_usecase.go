package usecases

import (
	"context"
	"stock-tracker/internal/domain/valueObjects"
)

type StockUseCase interface {
	GetStocks(ctx context.Context, filters valueObjects.StockFilters) (interface{}, *valueObjects.Pagination, error)
	GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error)
	GetStats(ctx context.Context) (interface{}, error)
}
