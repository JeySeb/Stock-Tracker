package database

import (
	"context"
	"fmt"
	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrokerRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewBrokerRepository(db *pgxpool.Pool) repositories.BrokerRepository {
	return &BrokerRepositoryImpl{db: db}
}

func (r *BrokerRepositoryImpl) Create(ctx context.Context, broker *entities.Broker) error {
	query := `
		INSERT INTO brokers (id, name, credibility_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		broker.ID, broker.Name, broker.CredibilityScore,
		broker.CreatedAt, broker.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create broker: %w", err)
	}

	return nil
}

func (r *BrokerRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.Broker, error) {
	query := `
		SELECT id, name, credibility_score, created_at, updated_at
		FROM brokers WHERE id = $1
	`

	broker := &entities.Broker{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&broker.ID, &broker.Name, &broker.CredibilityScore,
		&broker.CreatedAt, &broker.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get broker by ID: %w", err)
	}

	return broker, nil
}

func (r *BrokerRepositoryImpl) GetByName(ctx context.Context, name string) (*entities.Broker, error) {
	query := `
		SELECT id, name, credibility_score, created_at, updated_at
		FROM brokers WHERE name = $1
	`

	broker := &entities.Broker{}
	err := r.db.QueryRow(ctx, query, name).Scan(
		&broker.ID, &broker.Name, &broker.CredibilityScore,
		&broker.CreatedAt, &broker.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get broker by name: %w", err)
	}

	return broker, nil
}

func (r *BrokerRepositoryImpl) GetAll(ctx context.Context) ([]*entities.Broker, error) {
	query := `
		SELECT id, name, credibility_score, created_at, updated_at
		FROM brokers ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all brokers: %w", err)
	}
	defer rows.Close()

	var brokers []*entities.Broker
	for rows.Next() {
		broker := &entities.Broker{}
		err := rows.Scan(
			&broker.ID, &broker.Name, &broker.CredibilityScore,
			&broker.CreatedAt, &broker.UpdatedAt,
		)
		if err != nil {
			continue
		}
		brokers = append(brokers, broker)
	}

	return brokers, nil
}

func (r *BrokerRepositoryImpl) Update(ctx context.Context, broker *entities.Broker) error {
	query := `
		UPDATE brokers 
		SET name = $2, credibility_score = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		broker.ID, broker.Name, broker.CredibilityScore, broker.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update broker: %w", err)
	}

	return nil
}

func (r *BrokerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM brokers WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete broker: %w", err)
	}

	return nil
}

func (r *BrokerRepositoryImpl) UpsertByName(ctx context.Context, broker *entities.Broker) error {
	query := `
		INSERT INTO brokers (id, name, credibility_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name) DO UPDATE SET
			credibility_score = EXCLUDED.credibility_score,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(ctx, query,
		broker.ID, broker.Name, broker.CredibilityScore,
		broker.CreatedAt, broker.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert broker: %w", err)
	}

	return nil
}
