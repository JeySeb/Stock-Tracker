package database

import (
	"context"
	"fmt"
	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/repositories"
	"stock-tracker/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewUserRepository(db *pgxpool.Pool, logger logger.Logger) repositories.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, first_name, last_name, tier, is_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.Tier, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create user", "error", err, "email", user.Email)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, tier, is_verified, last_login, created_at, updated_at
        FROM users WHERE id = $1
    `

	user := &entities.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Tier, &user.IsVerified, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, tier, is_verified, last_login, created_at, updated_at
        FROM users WHERE email = $1
    `

	user := &entities.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Tier, &user.IsVerified, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
        UPDATE users 
        SET email = $2, password_hash = $3, first_name = $4, last_name = $5,
            tier = $6, is_verified = $7, updated_at = $8
        WHERE id = $1
    `

	result, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.Tier, user.IsVerified, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) VerifyUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_verified = true, updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login = NOW(), updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) GetUserCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}

	return count, nil
}

func (r *userRepository) GetUsersByTier(ctx context.Context, tier entities.UserTier) ([]*entities.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, tier, is_verified, last_login, created_at, updated_at
        FROM users WHERE tier = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(ctx, query, tier)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by tier: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
			&user.Tier, &user.IsVerified, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan user row", "error", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
