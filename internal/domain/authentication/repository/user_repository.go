package repositories

import (
	"context"
	"stock-tracker/internal/domain/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User verification
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Statistics
	GetUserCount(ctx context.Context) (int, error)
	GetUsersByTier(ctx context.Context, tier model.UserTier) ([]*model.User, error)
}
