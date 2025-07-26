package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	authModel "stock-tracker/internal/domain/authentication/model"
	authValidation "stock-tracker/internal/domain/authentication/validation"
	stockModel "stock-tracker/internal/domain/stocks/model"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
	subscriptionModel "stock-tracker/internal/domain/subscription/model"
	subscriptionValidation "stock-tracker/internal/domain/subscription/validation"
	stockRepos "stock-tracker/internal/domain/stocks/repositories"
)

// MockStockRepository implements repositories.StockRepository for testing
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) Create(ctx context.Context, stock *stockModel.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) GetByID(ctx context.Context, id uuid.UUID) (*stockModel.Stock, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) Update(ctx context.Context, stock *stockModel.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStockRepository) GetByTicker(ctx context.Context, ticker string) ([]*stockModel.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).([]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetLatestByTicker(ctx context.Context, ticker string) (*stockModel.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetAll(ctx context.Context, filters stockValidation.StockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*stockModel.Stock), args.Get(1).(*stockValidation.Pagination), args.Error(2)
}

func (m *MockStockRepository) GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*stockModel.Stock, error) {
	args := m.Called(ctx, since)
	return args.Get(0).(map[string][]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) BulkCreate(ctx context.Context, stocks []*stockModel.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

func (m *MockStockRepository) BulkUpdate(ctx context.Context, stocks []*stockModel.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

func (m *MockStockRepository) GetTopMoversByTarget(ctx context.Context, limit int) ([]*stockModel.Stock, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*stockModel.Stock), args.Error(1)
}

func (m *MockStockRepository) GetUniqueTickersCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockStockRepository) GetBrokerageStats(ctx context.Context) ([]stockRepos.BrokerageStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]stockRepos.BrokerageStats), args.Error(1)
}

// MockBrokerRepository implements stockRepos.BrokerRepository for testing
type MockBrokerRepository struct {
	mock.Mock
}

func (m *MockBrokerRepository) Create(ctx context.Context, broker *stockModel.Broker) error {
	args := m.Called(ctx, broker)
	return args.Error(0)
}

func (m *MockBrokerRepository) GetByID(ctx context.Context, id uuid.UUID) (*stockModel.Broker, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*stockModel.Broker), args.Error(1)
}

func (m *MockBrokerRepository) GetByName(ctx context.Context, name string) (*stockModel.Broker, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*stockModel.Broker), args.Error(1)
}

func (m *MockBrokerRepository) GetAll(ctx context.Context) ([]*stockModel.Broker, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*stockModel.Broker), args.Error(1)
}

func (m *MockBrokerRepository) Update(ctx context.Context, broker *stockModel.Broker) error {
	args := m.Called(ctx, broker)
	return args.Error(0)
}

func (m *MockBrokerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBrokerRepository) UpsertByName(ctx context.Context, broker *stockModel.Broker) error {
	args := m.Called(ctx, broker)
	return args.Error(0)
}
