package repositories

import (
	"context"
	authModel "stock-tracker/internal/domain/authentication/model"
	"stock-tracker/internal/domain/shared/enums"

	"github.com/google/uuid"
)

type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *authModel.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*authModel.User, error)
	GetByEmail(ctx context.Context, email string) (*authModel.User, error)
	Update(ctx context.Context, user *authModel.User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// User verification
	VerifyUser(ctx context.Context, userID uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error

	// Statistics
	GetUserCount(ctx context.Context) (int, error)
	GetUsersByTier(ctx context.Context, tier enums.UserTier) ([]*authModel.User, error)
}
