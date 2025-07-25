package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/model"
)

// MockStockAPIClient implements clients.StockAPIClient for testing
type MockStockAPIClient struct {
	mock.Mock
}

func (m *MockStockAPIClient) FetchAllStocks(ctx context.Context) ([]*model.Stock, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Stock), args.Error(1)
}

func (m *MockStockAPIClient) FetchPage(ctx context.Context, nextPage string) ([]*model.Stock, string, error) {
	args := m.Called(ctx, nextPage)
	return args.Get(0).([]*model.Stock), args.String(1), args.Error(2)
}
