package repositories

import (
	"context"
	"stock-tracker/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

// SubscriptionRepository defines the interface for subscription data access
type SubscriptionRepository interface {
	// Create stores a new subscription in the repository
	Create(ctx context.Context, subscription *entities.Subscription) error

	// GetByID retrieves a subscription by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Subscription, error)

	// GetByUserID retrieves all subscriptions for a given user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Subscription, error)

	// GetActiveByUserID retrieves the active subscription for a given user ID
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*entities.Subscription, error)

	// Update modifies an existing subscription in the repository
	Update(ctx context.Context, subscription *entities.Subscription) error

	// GetExpiring retrieves all subscriptions that will expire within the given duration
	GetExpiring(ctx context.Context, within time.Duration) ([]*entities.Subscription, error)

	// Delete removes a subscription from the repository
	Delete(ctx context.Context, id uuid.UUID) error

	// GetSubscriptionCount returns the total number of subscriptions
	GetSubscriptionCount(ctx context.Context) (int, error)
}
