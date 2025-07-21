package entities

import (
	"time"
	"github.com/google/uuid"
)

type Broker struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name" validate:"required,min=1,max=255"`
	CredibilityScore float64   `json:"credibility_score" db:"credibility_score" validate:"min=0,max=1"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}


func NewBroker(name string, credibilityScore float64) *Broker {
	return &Broker{
		ID:               uuid.New(),
		Name:             name,
		CredibilityScore: credibilityScore,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}