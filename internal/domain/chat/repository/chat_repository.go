package repositories

import (
	"context"
	"stock-tracker/internal/domain/model"

	"github.com/google/uuid"
)

type ChatRepository interface {
	// Chat Sessions
	CreateSession(ctx context.Context, session *model.ChatSession) error
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*model.ChatSession, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*model.ChatSession, error)
	UpdateSession(ctx context.Context, session *model.ChatSession) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error

	// Chat Messages
	CreateMessage(ctx context.Context, message *model.ChatMessage) error
	GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.ChatMessage, error)
	DeleteMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) error
}
