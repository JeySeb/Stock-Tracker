package repositories

import (
	"context"
	"stock-tracker/internal/domain/entities"

	"github.com/google/uuid"
)

type ChatRepository interface {
	// Chat Sessions
	CreateSession(ctx context.Context, session *entities.ChatSession) error
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.ChatSession, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*entities.ChatSession, error)
	UpdateSession(ctx context.Context, session *entities.ChatSession) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error

	// Chat Messages
	CreateMessage(ctx context.Context, message *entities.ChatMessage) error
	GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*entities.ChatMessage, error)
	DeleteMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) error
}
