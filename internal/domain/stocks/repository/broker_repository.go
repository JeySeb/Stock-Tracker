package repositories

import (
	"context"
	"stock-tracker/internal/domain/model"

	"github.com/google/uuid"
)

type BrokerRepository interface {
	Create(ctx context.Context, broker *model.Broker) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Broker, error)
	GetByName(ctx context.Context, name string) (*model.Broker, error)
	GetAll(ctx context.Context) ([]*model.Broker, error)
	Update(ctx context.Context, broker *model.Broker) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpsertByName(ctx context.Context, broker *model.Broker) error
}
