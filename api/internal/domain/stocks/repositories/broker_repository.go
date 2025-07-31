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
}
