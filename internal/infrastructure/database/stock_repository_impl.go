package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/repositories"
	"stock-tracker/internal/domain/valueObjects"
	"stock-tracker/pkg/logger"
)

type stockRepository struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

// NewStockRepository creates a new instance of stockRepository implementing repositories.StockRepository.
func NewStockRepository(db *pgxpool.Pool, logger logger.Logger) repositories.StockRepository {
	return &stockRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new stock record into the database.
func (r *stockRepository) Create(ctx context.Context, stock *entities.Stock) error {
	query := `
        INSERT INTO stocks (id, ticker, company, broker_id, action, rating_from, rating_to, 
                           target_from, target_to, event_time, price_close, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	_, err := r.db.Exec(ctx, query,
		stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
		stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
		stock.EventTime, stock.PriceClose, stock.CreatedAt, stock.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create stock", "error", err, "ticker", stock.Ticker)
		return fmt.Errorf("failed to create stock: %w", err)
	}

	return nil
}

// BulkCreate inserts multiple stock records in a single transaction.
// If any insert fails, the transaction is rolled back.
func (r *stockRepository) BulkCreate(ctx context.Context, stocks []*entities.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
        INSERT INTO stocks (id, ticker, company, broker_id, action, rating_from, rating_to, 
                           target_from, target_to, event_time, price_close, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        ON CONFLICT (ticker, event_time) DO NOTHING
    `

	for _, stock := range stocks {
		_, err := tx.Exec(ctx, query,
			stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
			stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
			stock.EventTime, stock.PriceClose, stock.CreatedAt, stock.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to insert stock in batch", "error", err, "ticker", stock.Ticker)
			return fmt.Errorf("failed to insert stock %s: %w", stock.Ticker, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Successfully inserted stocks batch", "count", len(stocks))
	return nil
}

// GetAll retrieves stocks from the database based on the provided filters and returns paginated results.
func (r *stockRepository) GetAll(ctx context.Context, filters valueObjects.StockFilters) ([]*entities.Stock, *valueObjects.Pagination, error) {
	filters.SetDefaults()

	whereClause, args := r.buildWhereClause(filters)
	countQuery := "SELECT COUNT(*) FROM stocks s LEFT JOIN brokers b ON s.broker_id = b.id" + whereClause

	r.logger.Info("Counting stocks", "query=%s", countQuery)
	var totalItems int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count stocks: %w", err)
	}

	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
    ` + whereClause + fmt.Sprintf(" ORDER BY s.%s %s LIMIT $%d OFFSET $%d",
		filters.SortBy, strings.ToUpper(filters.SortOrder), len(args)+1, len(args)+2)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	var stocks []*entities.Stock
	for rows.Next() {
		stock := &entities.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	pagination := &valueObjects.Pagination{
		Page:       (filters.Offset / filters.Limit) + 1,
		Limit:      filters.Limit,
		TotalItems: totalItems,
		TotalPages: (totalItems + filters.Limit - 1) / filters.Limit,
	}
	pagination.HasNext = pagination.Page < pagination.TotalPages
	pagination.HasPrev = pagination.Page > 1

	return stocks, pagination, nil
}

// buildWhereClause constructs the SQL WHERE clause and its arguments based on the provided filters.
func (r *stockRepository) buildWhereClause(filters valueObjects.StockFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Ticker != "" {
		conditions = append(conditions, fmt.Sprintf("s.ticker ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Ticker+"%")
		argIndex++
	}

	if filters.Company != "" {
		conditions = append(conditions, fmt.Sprintf("s.company ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Company+"%")
		argIndex++
	}

	if filters.Brokerage != "" {
		conditions = append(conditions, fmt.Sprintf("b.name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Brokerage+"%")
		argIndex++
	}

	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("s.event_time >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("s.event_time <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

// GetRecentByTickers retrieves recent stock records for all tickers since the given time.
func (r *stockRepository) GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*entities.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.event_time >= $1
        ORDER BY s.ticker, s.event_time DESC
    `

	rows, err := r.db.Query(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent stocks: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]*entities.Stock)

	for rows.Next() {
		stock := &entities.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}

		result[stock.Ticker] = append(result[stock.Ticker], stock)
	}

	return result, nil
}

// GetByID retrieves a stock by its ID.
func (r *stockRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.id = $1
    `

	stock := &entities.Stock{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
		&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
		&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
		&stock.BrokerID, &stock.Brokerage,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stock by ID: %w", err)
	}

	return stock, nil
}

// Update updates an existing stock record.
func (r *stockRepository) Update(ctx context.Context, stock *entities.Stock) error {
	query := `
        UPDATE stocks 
        SET ticker = $2, company = $3, broker_id = $4, action = $5, 
            rating_from = $6, rating_to = $7, target_from = $8, target_to = $9,
            event_time = $10, price_close = $11, updated_at = $12
        WHERE id = $1
    `

	_, err := r.db.Exec(ctx, query,
		stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
		stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
		stock.EventTime, stock.PriceClose, stock.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to update stock", "error", err, "ticker", stock.Ticker)
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

// Delete removes a stock record by ID.
func (r *stockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM stocks WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete stock", "error", err, "id", id)
		return fmt.Errorf("failed to delete stock: %w", err)
	}

	return nil
}

// GetByTicker retrieves all stocks for a specific ticker.
func (r *stockRepository) GetByTicker(ctx context.Context, ticker string) ([]*entities.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.ticker ILIKE $1
        ORDER BY s.event_time DESC
    `

	rows, err := r.db.Query(ctx, query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks by ticker: %w", err)
	}
	defer rows.Close()

	var stocks []*entities.Stock
	for rows.Next() {
		stock := &entities.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetLatestByTicker retrieves the most recent stock record for a specific ticker.
func (r *stockRepository) GetLatestByTicker(ctx context.Context, ticker string) (*entities.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.ticker ILIKE $1
        ORDER BY s.event_time DESC
        LIMIT 1
    `

	stock := &entities.Stock{}
	err := r.db.QueryRow(ctx, query, ticker).Scan(
		&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
		&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
		&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
		&stock.BrokerID, &stock.Brokerage,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get latest stock by ticker: %w", err)
	}

	return stock, nil
}

// BulkUpdate updates multiple stock records in a single transaction.
func (r *stockRepository) BulkUpdate(ctx context.Context, stocks []*entities.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
        UPDATE stocks 
        SET ticker = $2, company = $3, broker_id = $4, action = $5, 
            rating_from = $6, rating_to = $7, target_from = $8, target_to = $9,
            event_time = $10, price_close = $11, updated_at = $12
        WHERE id = $1
    `

	for _, stock := range stocks {
		_, err := tx.Exec(ctx, query,
			stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
			stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
			stock.EventTime, stock.PriceClose, stock.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to update stock in batch", "error", err, "ticker", stock.Ticker)
			return fmt.Errorf("failed to update stock %s: %w", stock.Ticker, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Successfully updated stocks batch", "count", len(stocks))
	return nil
}

// GetTopMoversByTarget retrieves stocks with the highest target price changes.
func (r *stockRepository) GetTopMoversByTarget(ctx context.Context, limit int) ([]*entities.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.target_from > 0 AND s.target_to > 0
        ORDER BY ((s.target_to - s.target_from) / s.target_from) DESC
        LIMIT $1
    `

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top movers: %w", err)
	}
	defer rows.Close()

	var stocks []*entities.Stock
	for rows.Next() {
		stock := &entities.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetUniqueTickersCount returns the count of unique tickers in the database.
func (r *stockRepository) GetUniqueTickersCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(DISTINCT ticker) FROM stocks`

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unique tickers: %w", err)
	}

	return count, nil
}

// GetBrokerageStats returns statistics for each brokerage.
func (r *stockRepository) GetBrokerageStats(ctx context.Context) ([]repositories.BrokerageStats, error) {
	query := `
		SELECT b.name as brokerage, COUNT(s.id) as count, AVG(b.credibility_score) as avg_score
		FROM brokers b
		LEFT JOIN stocks s ON b.id = s.broker_id
		GROUP BY b.id, b.name
		ORDER BY count DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query brokerage stats: %w", err)
	}
	defer rows.Close()

	var stats []repositories.BrokerageStats
	for rows.Next() {
		var stat repositories.BrokerageStats
		err := rows.Scan(&stat.Brokerage, &stat.Count, &stat.AvgScore)
		if err != nil {
			r.logger.Error("Failed to scan brokerage stats row", "error", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// GetRecentRecommendations gets recent positive recommendations for the recommendation engine
func (r *stockRepository) GetRecentRecommendations(ctx context.Context, since time.Time, limit int) ([]*entities.Stock, error) {
	query := `
		SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
		       s.target_from, s.target_to, s.event_time, s.price_close, s.created_at, s.updated_at,
		       b.id as broker_id, b.name as brokerage
		FROM stocks s
		LEFT JOIN brokers b ON s.broker_id = b.id
		WHERE s.event_time >= $1 
		  AND (s.action ILIKE '%upgraded%' OR s.action ILIKE '%initiated%' OR s.action ILIKE '%reiterated%')
		ORDER BY s.event_time DESC, b.credibility_score DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent recommendations: %w", err)
	}
	defer rows.Close()

	var stocks []*entities.Stock
	for rows.Next() {
		stock := &entities.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.PriceClose, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}
