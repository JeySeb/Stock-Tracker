package repositories

import (
	"context"
	chatModel "stock-tracker/internal/domain/chat/model"

	"github.com/google/uuid"
)

type ChatRepository interface {
	// Chat Sessions
	CreateSession(ctx context.Context, session *chatModel.ChatSession) error
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*chatModel.ChatSession, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*chatModel.ChatSession, error)
	UpdateSession(ctx context.Context, session *chatModel.ChatSession) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error

	// Chat Messages
	CreateMessage(ctx context.Context, message *chatModel.ChatMessage) error
	GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*chatModel.ChatMessage, error)
	DeleteMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) error
}
