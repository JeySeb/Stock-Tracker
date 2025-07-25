package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/model"
	"stock-tracker/internal/domain/repositories"
	"stock-tracker/internal/domain/authentication/validation"
)

// MockStockRepository implements repositories.StockRepository for testing
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) Create(ctx context.Context, stock *model.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Stock, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Stock), args.Error(1)
}

func (m *MockStockRepository) Update(ctx context.Context, stock *model.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStockRepository) GetByTicker(ctx context.Context, ticker string) ([]*model.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).([]*model.Stock), args.Error(1)
}

func (m *MockStockRepository) GetLatestByTicker(ctx context.Context, ticker string) (*model.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*model.Stock), args.Error(1)
}

func (m *MockStockRepository) GetAll(ctx context.Context, filters valueObjects.StockFilters) ([]*model.Stock, *valueObjects.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*model.Stock), args.Get(1).(*valueObjects.Pagination), args.Error(2)
}

func (m *MockStockRepository) GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*model.Stock, error) {
	args := m.Called(ctx, since)
	return args.Get(0).(map[string][]*model.Stock), args.Error(1)
}

func (m *MockStockRepository) BulkCreate(ctx context.Context, stocks []*model.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

func (m *MockStockRepository) BulkUpdate(ctx context.Context, stocks []*model.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

func (m *MockStockRepository) GetTopMoversByTarget(ctx context.Context, limit int) ([]*model.Stock, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*model.Stock), args.Error(1)
}

func (m *MockStockRepository) GetUniqueTickersCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStockRepository) GetBrokerageStats(ctx context.Context) ([]repositories.BrokerageStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repositories.BrokerageStats), args.Error(1)
}

// MockBrokerRepository implements repositories.BrokerRepository for testing
type MockBrokerRepository struct {
	mock.Mock
}

func (m *MockBrokerRepository) Create(ctx context.Context, broker *model.Broker) error {
	args := m.Called(ctx, broker)
	return args.Error(0)
}

func (m *MockBrokerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Broker, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Broker), args.Error(1)
}

func (m *MockBrokerRepository) GetByName(ctx context.Context, name string) (*model.Broker, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*model.Broker), args.Error(1)
}

func (m *MockBrokerRepository) GetAll(ctx context.Context) ([]*model.Broker, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Broker), args.Error(1)
}

func (m *MockBrokerRepository) Update(ctx context.Context, broker *model.Broker) error {
	args := m.Called(ctx, broker)
	return args.Error(0)
}

func (m *MockBrokerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBrokerRepository) UpsertByName(ctx context.Context, broker *model.Broker) error {
	args := m.Called(ctx, broker)
	return args.Error(0)
}
