package entities

import (
	"time"

	"github.com/google/uuid"
)

// ChatSession represents a chat conversation between a user and the AI assistant
type ChatSession struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	UserID    uuid.UUID     `json:"user_id" db:"user_id" validate:"required"`
	Title     string        `json:"title" db:"title" validate:"required,min=1,max=200"`
	Status    string        `json:"status" db:"status"`
	Messages  []ChatMessage `json:"messages,omitempty" db:"-"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

// NewChatSession creates a new chat session instance
func NewChatSession(userID uuid.UUID, title string) *ChatSession {
	now := time.Now()
	return &ChatSession{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     title,
		Status:    "active",
		Messages:  make([]ChatMessage, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddMessage adds a new message to the chat session
func (cs *ChatSession) AddMessage(role, content string) (*ChatMessage, error) {
	msg, err := NewChatMessage(cs.ID, role, content)
	if err != nil {
		return nil, err
	}
	cs.Messages = append(cs.Messages, *msg)
	cs.UpdatedAt = time.Now()
	return msg, nil
}

// GetMessages returns all messages in the chat session
func (cs *ChatSession) GetMessages() []ChatMessage {
	return cs.Messages
}

// Close marks the chat session as closed
func (cs *ChatSession) Close() {
	cs.Status = "closed"
	cs.UpdatedAt = time.Now()
}

// Reopen marks the chat session as active
func (cs *ChatSession) Reopen() {
	cs.Status = "active"
	cs.UpdatedAt = time.Now()
}

// UpdateTitle updates the chat session title
func (cs *ChatSession) UpdateTitle(newTitle string) {
	cs.Title = newTitle
	cs.UpdatedAt = time.Now()
}
