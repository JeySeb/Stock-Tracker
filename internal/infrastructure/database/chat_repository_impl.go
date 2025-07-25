package database

import (
	"context"
	"stock-tracker/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewChatRepository creates a new ChatRepositoryImpl instance
func NewChatRepository(db *pgxpool.Pool) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{db: db}
}

// CreateSession stores a new chat session in the database
func (r *ChatRepositoryImpl) CreateSession(ctx context.Context, session *entities.ChatSession) error {
	query := `
		INSERT INTO chat_sessions (id, user_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.Title,
		session.CreatedAt,
		session.UpdatedAt,
	)
	return err
}

// GetSessionsByUserID retrieves all chat sessions for a given user ID
func (r *ChatRepositoryImpl) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ChatSession, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM chat_sessions 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*entities.ChatSession
	for rows.Next() {
		session := &entities.ChatSession{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Title,
			&session.CreatedAt,
			&session.UpdatedAt,
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

// GetSessionByID retrieves a chat session by its ID
func (r *ChatRepositoryImpl) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*entities.ChatSession, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at
		FROM chat_sessions 
		WHERE id = $1
	`
	session := &entities.ChatSession{}
	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.Title,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// UpdateSession updates an existing chat session
func (r *ChatRepositoryImpl) UpdateSession(ctx context.Context, session *entities.ChatSession) error {
	query := `
		UPDATE chat_sessions 
		SET title = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query,
		session.Title,
		session.UpdatedAt,
		session.ID,
	)
	return err
}

// DeleteSession removes a chat session and its messages
func (r *ChatRepositoryImpl) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM chat_sessions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, sessionID)
	return err
}

// CreateMessage stores a new chat message in the database
func (r *ChatRepositoryImpl) CreateMessage(ctx context.Context, message *entities.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (id, session_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		message.ID,
		message.SessionID,
		message.Role,
		message.Content,
		message.CreatedAt,
	)
	return err
}

// GetMessagesBySessionID retrieves all messages for a given session ID
func (r *ChatRepositoryImpl) GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*entities.ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content, created_at
		FROM chat_messages 
		WHERE session_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entities.ChatMessage
	for rows.Next() {
		message := &entities.ChatMessage{}
		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

// DeleteMessagesBySessionID removes all messages for a given session ID
func (r *ChatRepositoryImpl) DeleteMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM chat_messages WHERE session_id = $1`
	_, err := r.db.Exec(ctx, query, sessionID)
	return err
}
