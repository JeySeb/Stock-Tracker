package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	stockModel "stock-tracker/internal/domain/stocks/model"
)

// MockStockAPIClient implements clients.StockAPIClient for testing
type MockStockAPIClient struct {
	mock.Mock
}

func (m *MockStockAPIClient) FetchAllStocks(ctx context.Context) ([]*stockModel.Stock, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*stockModel.Stock), args.Error(1)
}

func (m *MockStockAPIClient) FetchPage(ctx context.Context, nextPage string) ([]*stockModel.Stock, string, error) {
	args := m.Called(ctx, nextPage)
	return args.Get(0).([]*stockModel.Stock), args.String(1), args.Error(2)
}
