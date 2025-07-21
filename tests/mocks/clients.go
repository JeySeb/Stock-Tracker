package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/entities"
)

// MockStockAPIClient implements clients.StockAPIClient for testing
type MockStockAPIClient struct {
	mock.Mock
}

func (m *MockStockAPIClient) FetchAllStocks(ctx context.Context) ([]*entities.Stock, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Stock), args.Error(1)
}

func (m *MockStockAPIClient) FetchPage(ctx context.Context, nextPage string) ([]*entities.Stock, string, error) {
	args := m.Called(ctx, nextPage)
	return args.Get(0).([]*entities.Stock), args.String(1), args.Error(2)
}
