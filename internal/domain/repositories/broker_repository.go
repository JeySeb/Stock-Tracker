package repositories

import (
	"context"
	"stock-tracker/internal/domain/entities"

	"github.com/google/uuid"
)

type BrokerRepository interface {
	Create(ctx context.Context, broker *entities.Broker) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Broker, error)
	GetByName(ctx context.Context, name string) (*entities.Broker, error)
	GetAll(ctx context.Context) ([]*entities.Broker, error)
	Update(ctx context.Context, broker *entities.Broker) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpsertByName(ctx context.Context, broker *entities.Broker) error
}
