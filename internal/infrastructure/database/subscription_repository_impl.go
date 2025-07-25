package database

import (
	"context"
	"fmt"
	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/repositories"
	"stock-tracker/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewSubscriptionRepository(db *pgxpool.Pool, logger logger.Logger) repositories.SubscriptionRepository {
	return &subscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *entities.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, plan, status, price, currency, 
			start_date, end_date, payment_reference, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		subscription.ID, subscription.UserID, subscription.Plan,
		subscription.Status, subscription.Price, subscription.Currency,
		subscription.StartDate, subscription.EndDate, subscription.PaymentReference,
		subscription.CreatedAt, subscription.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create subscription", "error", err, "userID", subscription.UserID)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Subscription, error) {
	query := `
		SELECT id, user_id, plan, status, price, currency,
			start_date, end_date, payment_reference, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`

	subscription := &entities.Subscription{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.ID, &subscription.UserID, &subscription.Plan,
		&subscription.Status, &subscription.Price, &subscription.Currency,
		&subscription.StartDate, &subscription.EndDate, &subscription.PaymentReference,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get subscription by ID: %w", err)
	}

	return subscription, nil
}

func (r *subscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Subscription, error) {
	query := `
		SELECT id, user_id, plan, status, price, currency,
			start_date, end_date, payment_reference, created_at, updated_at
		FROM subscriptions WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by user ID: %w", err)
	}
	defer rows.Close()

	var subscriptions []*entities.Subscription
	for rows.Next() {
		subscription := &entities.Subscription{}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.Plan,
			&subscription.Status, &subscription.Price, &subscription.Currency,
			&subscription.StartDate, &subscription.EndDate, &subscription.PaymentReference,
			&subscription.CreatedAt, &subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription row: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*entities.Subscription, error) {
	query := `
		SELECT id, user_id, plan, status, price, currency,
			start_date, end_date, payment_reference, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1 AND status = 'active' AND end_date > NOW()
		ORDER BY end_date DESC
		LIMIT 1
	`

	subscription := &entities.Subscription{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.Plan,
		&subscription.Status, &subscription.Price, &subscription.Currency,
		&subscription.StartDate, &subscription.EndDate, &subscription.PaymentReference,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get active subscription: %w", err)
	}

	return subscription, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *entities.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET plan = $2, status = $3, price = $4, currency = $5,
			start_date = $6, end_date = $7, payment_reference = $8, updated_at = $9
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		subscription.ID, subscription.Plan, subscription.Status,
		subscription.Price, subscription.Currency,
		subscription.StartDate, subscription.EndDate,
		subscription.PaymentReference, subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *subscriptionRepository) GetExpiring(ctx context.Context, within time.Duration) ([]*entities.Subscription, error) {
	query := `
		SELECT id, user_id, plan, status, price, currency,
			start_date, end_date, payment_reference, created_at, updated_at
		FROM subscriptions 
		WHERE status = 'active' 
		AND end_date BETWEEN NOW() AND NOW() + $1
		ORDER BY end_date ASC
	`

	rows, err := r.db.Query(ctx, query, within)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*entities.Subscription
	for rows.Next() {
		subscription := &entities.Subscription{}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.Plan,
			&subscription.Status, &subscription.Price, &subscription.Currency,
			&subscription.StartDate, &subscription.EndDate, &subscription.PaymentReference,
			&subscription.CreatedAt, &subscription.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription row: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *subscriptionRepository) GetSubscriptionCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM subscriptions`

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get subscription count: %w", err)
	}

	return count, nil
}
