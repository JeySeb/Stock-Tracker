package database

import (
	"context"
	"stock-tracker/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewSessionRepository creates a new SessionRepositoryImpl instance
func NewSessionRepository(db *pgxpool.Pool) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{db: db}
}

// Create stores a new session in the database
func (r *SessionRepositoryImpl) Create(ctx context.Context, session *entities.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
		session.CreatedAt,
	)
	return err
}

// GetByRefreshToken retrieves a session by its refresh token
func (r *SessionRepositoryImpl) GetByRefreshToken(ctx context.Context, refreshToken string) (*entities.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at
		FROM sessions 
		WHERE refresh_token = $1
	`
	session := &entities.Session{}
	err := r.db.QueryRow(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// DeleteByUserID removes all sessions for a given user ID
func (r *SessionRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// DeleteExpired removes all expired sessions from the database
func (r *SessionRepositoryImpl) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, query)
	return err
}

// DeleteByRefreshToken removes a specific session by its refresh token
func (r *SessionRepositoryImpl) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM sessions WHERE refresh_token = $1`
	_, err := r.db.Exec(ctx, query, refreshToken)
	return err
}

// GetByUserID retrieves all active sessions for a given user ID
func (r *SessionRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, user_agent, ip_address, expires_at, created_at
		FROM sessions 
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*entities.Session
	for rows.Next() {
		session := &entities.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RefreshToken,
			&session.UserAgent,
			&session.IPAddress,
			&session.ExpiresAt,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}
