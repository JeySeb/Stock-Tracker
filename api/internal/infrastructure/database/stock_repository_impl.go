package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	stockModel "stock-tracker/internal/domain/stocks/model"
	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	stockValidation "stock-tracker/internal/domain/stocks/validation"
	"stock-tracker/pkg/logger"
)

type stockRepository struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

// NewStockRepository creates a new instance of stockRepository implementing repositories.StockRepository.
func NewStockRepository(db *pgxpool.Pool, logger logger.Logger) stockRepos.StockRepository {
	return &stockRepository{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new stock record into the database.
func (r *stockRepository) Create(ctx context.Context, stock *stockModel.Stock) error {
	query := `
        INSERT INTO stocks (id, ticker, company, broker_id, action, rating_from, rating_to, 
                           target_from, target_to, event_time, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	_, err := r.db.Exec(ctx, query,
		stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
		stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
		stock.EventTime, stock.CreatedAt, stock.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create stock", "error", err, "ticker", stock.Ticker)
		return fmt.Errorf("failed to create stock: %w", err)
	}

	return nil
}

// BulkCreate inserts multiple stock records in a single transaction.
// If any insert fails, the transaction is rolled back.
func (r *stockRepository) BulkCreate(ctx context.Context, stocks []*stockModel.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			// Log rollback error but don't override the main error
			// Note: Rollback errors during cleanup are typically non-critical
			r.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	query := `
        INSERT INTO stocks (id, ticker, company, broker_id, action, rating_from, rating_to, 
                           target_from, target_to, event_time,  created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        ON CONFLICT (ticker, event_time) DO NOTHING
    `

	for _, stock := range stocks {
		_, err := tx.Exec(ctx, query,
			stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
			stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
			stock.EventTime, stock.CreatedAt, stock.UpdatedAt,
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
func (r *stockRepository) GetAll(ctx context.Context, filters stockValidation.StockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error) {
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
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
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

	var stocks []*stockModel.Stock
	for rows.Next() {
		stock := &stockModel.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	pagination := &stockValidation.Pagination{
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
func (r *stockRepository) buildWhereClause(filters stockValidation.StockFilters) (string, []interface{}) {
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
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

// GetRecentByTickers retrieves recent stock records for all tickers since the given time.
func (r *stockRepository) GetRecentByTickers(ctx context.Context, since time.Time) (map[string][]*stockModel.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
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

	result := make(map[string][]*stockModel.Stock)

	for rows.Next() {
		stock := &stockModel.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
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
func (r *stockRepository) GetByID(ctx context.Context, id uuid.UUID) (*stockModel.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.id = $1
    `

	stock := &stockModel.Stock{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
		&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
		&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
		&stock.BrokerID, &stock.Brokerage,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stock by ID: %w", err)
	}

	return stock, nil
}

// Update updates an existing stock record.
func (r *stockRepository) Update(ctx context.Context, stock *stockModel.Stock) error {
	query := `
        UPDATE stocks 
        SET ticker = $2, company = $3, broker_id = $4, action = $5, 
            rating_from = $6, rating_to = $7, target_from = $8, target_to = $9,
            event_time = $10, updated_at = $11
        WHERE id = $1
    `

	_, err := r.db.Exec(ctx, query,
		stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
		stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
		stock.EventTime, stock.UpdatedAt,
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
func (r *stockRepository) GetByTicker(ctx context.Context, ticker string) ([]*stockModel.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
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

	var stocks []*stockModel.Stock
	for rows.Next() {
		stock := &stockModel.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
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
func (r *stockRepository) GetLatestByTicker(ctx context.Context, ticker string) (*stockModel.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
               b.id as broker_id, b.name as brokerage
        FROM stocks s
        LEFT JOIN brokers b ON s.broker_id = b.id
        WHERE s.ticker ILIKE $1
        ORDER BY s.event_time DESC
        LIMIT 1
    `

	stock := &stockModel.Stock{}
	err := r.db.QueryRow(ctx, query, ticker).Scan(
		&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
		&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
		&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
		&stock.BrokerID, &stock.Brokerage,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get latest stock by ticker: %w", err)
	}

	return stock, nil
}

// BulkUpdate updates multiple stock records in a single transaction.
func (r *stockRepository) BulkUpdate(ctx context.Context, stocks []*stockModel.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			// Log rollback error but don't override the main error
			// Note: Rollback errors during cleanup are typically non-critical
			r.logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	query := `
        UPDATE stocks 
        SET ticker = $2, company = $3, broker_id = $4, action = $5, 
            rating_from = $6, rating_to = $7, target_from = $8, target_to = $9,
            event_time = $10, updated_at = $11
        WHERE id = $1
    `

	for _, stock := range stocks {
		_, err := tx.Exec(ctx, query,
			stock.ID, stock.Ticker, stock.Company, stock.BrokerID, stock.Action,
			stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
			stock.EventTime, stock.UpdatedAt,
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
func (r *stockRepository) GetTopMoversByTarget(ctx context.Context, limit int) ([]*stockModel.Stock, error) {
	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
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

	var stocks []*stockModel.Stock
	for rows.Next() {
		stock := &stockModel.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
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

// GetUniqueTickers returns all unique tickers from the database.
func (r *stockRepository) GetUniqueTickers(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT ticker FROM stocks ORDER BY ticker`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("Failed to get unique tickers", "error", err)
		return nil, fmt.Errorf("failed to get unique tickers: %w", err)
	}
	defer rows.Close()

	var tickers []string
	for rows.Next() {
		var ticker string
		if err := rows.Scan(&ticker); err != nil {
			r.logger.Error("Failed to scan ticker", "error", err)
			return nil, fmt.Errorf("failed to scan ticker: %w", err)
		}
		tickers = append(tickers, ticker)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over tickers", "error", err)
		return nil, fmt.Errorf("error iterating over tickers: %w", err)
	}

	r.logger.Debug("Retrieved unique tickers", "count", len(tickers))
	return tickers, nil
}

// GetUniqueCompanies returns all unique companies from the database.
func (r *stockRepository) GetUniqueCompanies(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT company FROM stocks ORDER BY company`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("Failed to get unique companies", "error", err)
		return nil, fmt.Errorf("failed to get unique companies: %w", err)
	}
	defer rows.Close()

	var companies []string
	for rows.Next() {
		var company string
		if err := rows.Scan(&company); err != nil {
			r.logger.Error("Failed to scan company", "error", err)
			return nil, fmt.Errorf("failed to scan company: %w", err)
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over companies", "error", err)
		return nil, fmt.Errorf("error iterating over companies: %w", err)
	}

	r.logger.Debug("Retrieved unique companies", "count", len(companies))
	return companies, nil
}

// GetBrokerageStats returns statistics for each brokerage.
func (r *stockRepository) GetBrokerageStats(ctx context.Context) ([]stockRepos.BrokerageStats, error) {
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

	var stats []stockRepos.BrokerageStats
	for rows.Next() {
		var stat stockRepos.BrokerageStats
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
func (r *stockRepository) GetRecentRecommendations(ctx context.Context, since time.Time, limit int) ([]*stockModel.Stock, error) {
	query := `
		SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
		       s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
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

	var stocks []*stockModel.Stock
	for rows.Next() {
		stock := &stockModel.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
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

// GetAllWithEnhancedFilters retrieves stocks from the database based on the provided enhanced filters and returns paginated results.
func (r *stockRepository) GetAllWithEnhancedFilters(ctx context.Context, filters stockValidation.EnhancedStockFilters) ([]*stockModel.Stock, *stockValidation.Pagination, error) {
	filters.SetDefaults()

	whereClause, args := r.buildEnhancedWhereClause(filters)
	countQuery := "SELECT COUNT(*) FROM stocks s LEFT JOIN brokers b ON s.broker_id = b.id" + whereClause

	r.logger.Info("Counting stocks with enhanced filters", "query=%s", countQuery)
	var totalItems int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count stocks: %w", err)
	}

	query := `
        SELECT s.id, s.ticker, s.company, s.action, s.rating_from, s.rating_to,
               s.target_from, s.target_to, s.event_time,  s.created_at, s.updated_at,
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

	var stocks []*stockModel.Stock
	for rows.Next() {
		stock := &stockModel.Stock{}
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Action,
			&stock.RatingFrom, &stock.RatingTo, &stock.TargetFrom, &stock.TargetTo,
			&stock.EventTime, &stock.CreatedAt, &stock.UpdatedAt,
			&stock.BrokerID, &stock.Brokerage,
		)
		if err != nil {
			r.logger.Error("Failed to scan stock row", "error", err)
			continue
		}
		stocks = append(stocks, stock)
	}

	pagination := &stockValidation.Pagination{
		Page:       (filters.Offset / filters.Limit) + 1,
		Limit:      filters.Limit,
		TotalItems: totalItems,
		TotalPages: (totalItems + filters.Limit - 1) / filters.Limit,
	}
	pagination.HasNext = pagination.Page < pagination.TotalPages
	pagination.HasPrev = pagination.Page > 1

	return stocks, pagination, nil
}

// buildEnhancedWhereClause constructs the SQL WHERE clause and its arguments based on the provided enhanced filters.
func (r *stockRepository) buildEnhancedWhereClause(filters stockValidation.EnhancedStockFilters) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Handle multiple tickers
	if len(filters.Tickers) > 0 {
		placeholders := make([]string, len(filters.Tickers))
		for i := range filters.Tickers {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, filters.Tickers[i])
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("s.ticker IN (%s)", strings.Join(placeholders, ",")))
	}

	// Handle multiple companies
	if len(filters.Companies) > 0 {
		placeholders := make([]string, len(filters.Companies))
		for i := range filters.Companies {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, "%"+filters.Companies[i]+"%")
			argIndex++
		}
		companyConditions := make([]string, len(placeholders))
		for i, placeholder := range placeholders {
			companyConditions[i] = fmt.Sprintf("s.company ILIKE %s", placeholder)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(companyConditions, " OR ")))
	}

	// Handle multiple brokerages
	if len(filters.Brokerages) > 0 {
		placeholders := make([]string, len(filters.Brokerages))
		for i := range filters.Brokerages {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, "%"+filters.Brokerages[i]+"%")
			argIndex++
		}
		brokerageConditions := make([]string, len(placeholders))
		for i, placeholder := range placeholders {
			brokerageConditions[i] = fmt.Sprintf("b.name ILIKE %s", placeholder)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(brokerageConditions, " OR ")))
	}

	// Handle multiple actions
	if len(filters.Actions) > 0 {
		placeholders := make([]string, len(filters.Actions))
		for i := range filters.Actions {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, "%"+filters.Actions[i]+"%")
			argIndex++
		}
		actionConditions := make([]string, len(placeholders))
		for i, placeholder := range placeholders {
			actionConditions[i] = fmt.Sprintf("s.action ILIKE %s", placeholder)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(actionConditions, " OR ")))
	}

	// Handle rating filters
	if filters.RatingFrom != "" {
		conditions = append(conditions, fmt.Sprintf("s.rating_from ILIKE $%d", argIndex))
		args = append(args, "%"+filters.RatingFrom+"%")
		argIndex++
	}

	if filters.RatingTo != "" {
		conditions = append(conditions, fmt.Sprintf("s.rating_to ILIKE $%d", argIndex))
		args = append(args, "%"+filters.RatingTo+"%")
		argIndex++
	}

	// Handle target price filters
	if filters.TargetFrom != nil {
		conditions = append(conditions, fmt.Sprintf("s.target_from >= $%d", argIndex))
		args = append(args, *filters.TargetFrom)
		argIndex++
	}

	if filters.TargetTo != nil {
		conditions = append(conditions, fmt.Sprintf("s.target_to <= $%d", argIndex))
		args = append(args, *filters.TargetTo)
		argIndex++
	}

	// Handle advanced filters
	if filters.HasTargetPrice != nil {
		if *filters.HasTargetPrice {
			conditions = append(conditions, "s.target_from IS NOT NULL AND s.target_to IS NOT NULL")
		} else {
			conditions = append(conditions, "(s.target_from IS NULL OR s.target_to IS NULL)")
		}
	}

	if filters.HasRating != nil {
		if *filters.HasRating {
			conditions = append(conditions, "s.rating_from IS NOT NULL AND s.rating_to IS NOT NULL")
		} else {
			conditions = append(conditions, "(s.rating_from IS NULL OR s.rating_to IS NULL)")
		}
	}

	// Handle target change percentage filters
	if filters.MinTargetChange != nil || filters.MaxTargetChange != nil {
		changeCondition := "CASE WHEN s.target_from > 0 AND s.target_to > 0 THEN ((s.target_to - s.target_from) / s.target_from * 100) ELSE 0 END"

		if filters.MinTargetChange != nil {
			conditions = append(conditions, fmt.Sprintf("%s >= $%d", changeCondition, argIndex))
			args = append(args, *filters.MinTargetChange)
			argIndex++
		}

		if filters.MaxTargetChange != nil {
			conditions = append(conditions, fmt.Sprintf("%s <= $%d", changeCondition, argIndex))
			args = append(args, *filters.MaxTargetChange)
			argIndex++
		}
	}

	// Handle brokerage credibility score filters
	if filters.MinBrokerScore != nil {
		conditions = append(conditions, fmt.Sprintf("b.credibility_score >= $%d", argIndex))
		args = append(args, *filters.MinBrokerScore)
		argIndex++
	}

	if filters.MaxBrokerScore != nil {
		conditions = append(conditions, fmt.Sprintf("b.credibility_score <= $%d", argIndex))
		args = append(args, *filters.MaxBrokerScore)
		argIndex++
	}

	// Handle enhanced date filters
	dateConditions := r.buildEnhancedDateConditions(filters, &argIndex, &args)
	if len(dateConditions) > 0 {
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(dateConditions, " OR ")))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

// buildEnhancedDateConditions builds date-related WHERE conditions for enhanced filters
func (r *stockRepository) buildEnhancedDateConditions(filters stockValidation.EnhancedStockFilters, argIndex *int, args *[]interface{}) []string {
	var conditions []string

	// Handle basic date range
	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("s.event_time >= $%d", *argIndex))
		*args = append(*args, *filters.DateFrom)
		*argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("s.event_time <= $%d", *argIndex))
		*args = append(*args, *filters.DateTo)
		*argIndex++
	}

	// Handle multiple date ranges
	for _, dateRange := range filters.DateRanges {
		rangeCondition := fmt.Sprintf("(s.event_time >= $%d AND s.event_time <= $%d)", *argIndex, *argIndex+1)
		conditions = append(conditions, rangeCondition)
		*args = append(*args, dateRange.From, dateRange.To)
		*argIndex += 2
	}

	// Handle time-based filters
	now := time.Now()
	if filters.LastHours != nil {
		timeFrom := now.Add(-time.Duration(*filters.LastHours) * time.Hour)
		conditions = append(conditions, fmt.Sprintf("s.event_time >= $%d", *argIndex))
		*args = append(*args, timeFrom)
		*argIndex++
	}

	if filters.LastDays != nil {
		timeFrom := now.Add(-time.Duration(*filters.LastDays) * 24 * time.Hour)
		conditions = append(conditions, fmt.Sprintf("s.event_time >= $%d", *argIndex))
		*args = append(*args, timeFrom)
		*argIndex++
	}

	if filters.LastWeeks != nil {
		timeFrom := now.Add(-time.Duration(*filters.LastWeeks) * 7 * 24 * time.Hour)
		conditions = append(conditions, fmt.Sprintf("s.event_time >= $%d", *argIndex))
		*args = append(*args, timeFrom)
		*argIndex++
	}

	if filters.LastMonths != nil {
		timeFrom := now.Add(-time.Duration(*filters.LastMonths) * 30 * 24 * time.Hour)
		conditions = append(conditions, fmt.Sprintf("s.event_time >= $%d", *argIndex))
		*args = append(*args, timeFrom)
		*argIndex++
	}

	return conditions
}
