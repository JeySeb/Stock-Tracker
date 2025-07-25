package usecases

import (
	"context"
	"stock-tracker/internal/domain/authentication/validation"
)

type StockUseCase interface {
	GetStocks(ctx context.Context, filters validation.StockFilters) (interface{}, *validation.Pagination, error)
	GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error)
	GetStats(ctx context.Context) (interface{}, error)
}
