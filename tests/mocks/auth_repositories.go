package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/entities"
)

// MockUserRepository implements repositories.UserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetUsersByTier(ctx context.Context, tier entities.UserTier) ([]*entities.User, error) {
	args := m.Called(ctx, tier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

// MockSessionRepository implements repositories.SessionRepository for testing
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *entities.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.Session, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Session), args.Error(1)
}

func (m *MockSessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Session), args.Error(1)
}

// MockSubscriptionRepository implements repositories.SubscriptionRepository for testing
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, subscription *entities.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*entities.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(ctx context.Context, subscription *entities.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetExpiring(ctx context.Context, within time.Duration) ([]*entities.Subscription, error) {
	args := m.Called(ctx, within)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetSubscriptionCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}
