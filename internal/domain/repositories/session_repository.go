package repositories

import (
	"context"
	"stock-tracker/internal/domain/entities"

	"github.com/google/uuid"
)

// SessionRepository defines the interface for session persistence operations
type SessionRepository interface {
	// Create stores a new session in the database
	Create(ctx context.Context, session *entities.Session) error

	// GetByRefreshToken retrieves a session by its refresh token
	GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.Session, error)

	// DeleteByUserID removes all sessions for a given user ID
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired removes all expired sessions from the database
	DeleteExpired(ctx context.Context) error

	// DeleteByRefreshToken removes a specific session by its refresh token
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error

	// GetByUserID retrieves all active sessions for a given user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Session, error)
}
