package repositories

import (
	"context"
	"stock-tracker/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User verification
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Statistics
	GetUserCount(ctx context.Context) (int, error)
	GetUsersByTier(ctx context.Context, tier entities.UserTier) ([]*entities.User, error)
}
