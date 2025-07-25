package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// MessageRole defines the type of chat message role
type MessageRole string

const (
	// RoleUser represents a message from the user
	RoleUser MessageRole = "user"
	// RoleAssistant represents a message from the AI assistant
	RoleAssistant MessageRole = "assistant"
)

// ChatMessage represents a single message in a chat session
type ChatMessage struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	SessionID uuid.UUID   `json:"session_id" db:"session_id" validate:"required"`
	Role      MessageRole `json:"role" db:"role" validate:"required,oneof=user assistant"`
	Content   string      `json:"content" db:"content" validate:"required,min=1"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}

// NewChatMessage creates a new chat message instance with validation
func NewChatMessage(sessionID uuid.UUID, role string, content string) (*ChatMessage, error) {
	if sessionID == uuid.Nil {
		return nil, errors.New("session ID is required")
	}
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	messageRole := MessageRole(role)
	if messageRole != RoleUser && messageRole != RoleAssistant {
		return nil, errors.New("invalid message role")
	}

	return &ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		Role:      messageRole,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

// GetRole returns the message role as a string
func (m *ChatMessage) GetRole() string {
	return string(m.Role)
}

// IsFromUser checks if the message is from a user
func (m *ChatMessage) IsFromUser() bool {
	return m.Role == RoleUser
}

// IsFromAssistant checks if the message is from the assistant
func (m *ChatMessage) IsFromAssistant() bool {
	return m.Role == RoleAssistant
}

// GetContent returns the message content
func (m *ChatMessage) GetContent() string {
	return m.Content
}
