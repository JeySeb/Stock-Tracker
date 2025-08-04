package repositories

import (
	"context"
	stockModel "stock-tracker/internal/domain/stocks/model"

	"github.com/google/uuid"
)

type BrokerRepository interface {
	Create(ctx context.Context, broker *stockModel.Broker) error
	GetByID(ctx context.Context, id uuid.UUID) (*stockModel.Broker, error)
	GetByName(ctx context.Context, name string) (*stockModel.Broker, error)
	GetAll(ctx context.Context) ([]*stockModel.Broker, error)
	Update(ctx context.Context, broker *stockModel.Broker) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpsertByName(ctx context.Context, broker *stockModel.Broker) error

	// New methods for broker scoring endpoints
	GetBrokersWithScores(ctx context.Context, limit *int, orderBy string) ([]*BrokerWithScore, error)
}

// BrokerWithScore represents a broker with its calculated score
type BrokerWithScore struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	CredibilityScore float64   `json:"credibility_score" db:"credibility_score"`
	ReportCount      int       `json:"report_count" db:"report_count"`
	CalculatedScore  float64   `json:"calculated_score" db:"calculated_score"`
	CreatedAt        string    `json:"created_at" db:"created_at"`
	UpdatedAt        string    `json:"updated_at" db:"updated_at"`
}
