package usecases

import (
	"context"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
)

type StockUseCase interface {
	GetStocks(ctx context.Context, filters stockValidation.StockFilters) (interface{}, *stockValidation.Pagination, error)
	GetStocksByTicker(ctx context.Context, ticker string) (interface{}, error)
	GetStats(ctx context.Context) (interface{}, error)
}
