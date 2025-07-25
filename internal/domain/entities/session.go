package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Session represents a user's authenticated session
type Session struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token" validate:"required"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SessionDuration defines how long a session remains valid
const SessionDuration = 7 * 24 * time.Hour // 7 days

// NewSession creates a new session instance with validation
func NewSession(userID uuid.UUID, refreshToken, userAgent, ipAddress string) (*Session, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	now := time.Now()
	return &Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    now.Add(SessionDuration),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Extend prolongs the session duration from the current time
func (s *Session) Extend() {
	s.ExpiresAt = time.Now().Add(SessionDuration)
	s.UpdatedAt = time.Now()
}

// Invalidate immediately expires the session
func (s *Session) Invalidate() {
	s.ExpiresAt = time.Now()
	s.UpdatedAt = time.Now()
}

// GetUserID returns the associated user ID
func (s *Session) GetUserID() uuid.UUID {
	return s.UserID
}
