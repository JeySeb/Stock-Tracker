package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
	"stock-tracker/internal/domain/marketdata/repositories"
	"stock-tracker/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// marketDataRepositoryImpl implements the MarketDataRepository interface
type marketDataRepositoryImpl struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

// NewMarketDataRepository creates a new market data repository instance
func NewMarketDataRepository(pool *pgxpool.Pool, logger logger.Logger) repositories.MarketDataRepository {
	return &marketDataRepositoryImpl{
		pool:   pool,
		logger: logger,
	}
}

// SaveMarketData saves market data to the database
func (r *marketDataRepositoryImpl) SaveMarketData(ctx context.Context, marketData *model.MarketData) error {
	query := `
		INSERT INTO market_data (
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, market_cap, pe_ratio, dividend_yield, week_52_high, week_52_low, avg_volume,
			collected_at, data_timestamp, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`

	_, err := r.pool.Exec(ctx, query,
		marketData.ID,
		marketData.Ticker,
		marketData.DataSource,
		marketData.DataQuality,
		marketData.CurrentPrice,
		marketData.DayChange,
		marketData.DayChangePercent,
		marketData.Volume,
		marketData.MarketCap,
		marketData.PERatio,
		marketData.DividendYield,
		marketData.Week52High,
		marketData.Week52Low,
		marketData.AvgVolume,
		marketData.CollectedAt,
		marketData.DataTimestamp,
		marketData.CreatedAt,
		marketData.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to save market data",
			"ticker", marketData.Ticker,
			"error", err)
		return fmt.Errorf("failed to save market data: %w", err)
	}

	r.logger.Debug("Market data saved successfully",
		"ticker", marketData.Ticker,
		"data_source", marketData.DataSource)

	return nil
}

// SaveIngestionLog saves an ingestion log entry
func (r *marketDataRepositoryImpl) SaveIngestionLog(ctx context.Context, log *model.MarketDataIngestionLog) error {
	query := `
		INSERT INTO market_data_ingestion_logs (
			id, batch_id, data_source, total_tickers, successful_tickers, failed_tickers,
			skipped_tickers, status, error_details, started_at, completed_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.pool.Exec(ctx, query,
		log.ID,
		log.BatchID,
		log.DataSource,
		log.TotalTickers,
		log.SuccessfulTickers,
		log.FailedTickers,
		log.SkippedTickers,
		log.Status,
		log.ErrorDetails,
		log.StartedAt,
		log.CompletedAt,
		log.CreatedAt,
		log.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to save ingestion log",
			"batch_id", log.BatchID,
			"error", err)
		return fmt.Errorf("failed to save ingestion log: %w", err)
	}

	return nil
}

// UpdateIngestionLog updates an existing ingestion log
func (r *marketDataRepositoryImpl) UpdateIngestionLog(ctx context.Context, log *model.MarketDataIngestionLog) error {
	query := `
		UPDATE market_data_ingestion_logs SET
			successful_tickers = $1,
			failed_tickers = $2,
			skipped_tickers = $3,
			status = $4,
			error_details = $5,
			completed_at = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := r.pool.Exec(ctx, query,
		log.SuccessfulTickers,
		log.FailedTickers,
		log.SkippedTickers,
		log.Status,
		log.ErrorDetails,
		log.CompletedAt,
		log.UpdatedAt,
		log.ID,
	)

	if err != nil {
		r.logger.Error("Failed to update ingestion log",
			"batch_id", log.BatchID,
			"error", err)
		return fmt.Errorf("failed to update ingestion log: %w", err)
	}

	return nil
}

// GetUniqueTickers retrieves unique tickers from the stocks table
func (r *marketDataRepositoryImpl) GetUniqueTickers(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT ticker FROM stocks ORDER BY ticker`

	rows, err := r.pool.Query(ctx, query)
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

// ExistsMarketData checks if market data exists for the given parameters
func (r *marketDataRepositoryImpl) ExistsMarketData(ctx context.Context, ticker string, dataSource model.DataSource, dataTimestamp time.Time) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM market_data 
			WHERE ticker = $1 AND data_source = $2 AND data_timestamp = $3
		)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, ticker, dataSource, dataTimestamp).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check market data existence",
			"ticker", ticker,
			"data_source", dataSource,
			"error", err)
		return false, fmt.Errorf("failed to check market data existence: %w", err)
	}

	return exists, nil
}

// GetLatestMarketData retrieves the latest market data for a ticker
func (r *marketDataRepositoryImpl) GetLatestMarketData(ctx context.Context, ticker string) (*model.MarketData, error) {
	query := `
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, market_cap, pe_ratio, dividend_yield, week_52_high, week_52_low, avg_volume,
			collected_at, data_timestamp, created_at, updated_at
		FROM market_data 
		WHERE ticker = $1 
		ORDER BY data_timestamp DESC 
		LIMIT 1
	`

	var marketData model.MarketData
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&marketData.ID,
		&marketData.Ticker,
		&marketData.DataSource,
		&marketData.DataQuality,
		&marketData.CurrentPrice,
		&marketData.DayChange,
		&marketData.DayChangePercent,
		&marketData.Volume,
		&marketData.MarketCap,
		&marketData.PERatio,
		&marketData.DividendYield,
		&marketData.Week52High,
		&marketData.Week52Low,
		&marketData.AvgVolume,
		&marketData.CollectedAt,
		&marketData.DataTimestamp,
		&marketData.CreatedAt,
		&marketData.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get latest market data",
			"ticker", ticker,
			"error", err)
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	return &marketData, nil
}

// GetMarketDataStats retrieves market data statistics
func (r *marketDataRepositoryImpl) GetMarketDataStats(ctx context.Context, days int) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_records,
			COUNT(DISTINCT ticker) as unique_tickers,
			COUNT(*) FILTER (WHERE data_source = 'yahoo_finance') as yahoo_records,
			COUNT(*) FILTER (WHERE data_source = 'alpha_vantage') as alpha_records,
			AVG(current_price) as avg_price,
			AVG(ABS(day_change_percent)) as avg_volatility
		FROM market_data 
		WHERE data_timestamp > now() - INTERVAL '1 day' * $1
	`

	var stats struct {
		TotalRecords  int64   `db:"total_records"`
		UniqueTickers int64   `db:"unique_tickers"`
		YahooRecords  int64   `db:"yahoo_records"`
		AlphaRecords  int64   `db:"alpha_records"`
		AvgPrice      float64 `db:"avg_price"`
		AvgVolatility float64 `db:"avg_volatility"`
	}

	err := r.pool.QueryRow(ctx, query, days).Scan(
		&stats.TotalRecords,
		&stats.UniqueTickers,
		&stats.YahooRecords,
		&stats.AlphaRecords,
		&stats.AvgPrice,
		&stats.AvgVolatility,
	)

	if err != nil {
		r.logger.Error("Failed to get market data stats", "error", err)
		return nil, fmt.Errorf("failed to get market data stats: %w", err)
	}

	return map[string]interface{}{
		"total_records":  stats.TotalRecords,
		"unique_tickers": stats.UniqueTickers,
		"yahoo_records":  stats.YahooRecords,
		"alpha_records":  stats.AlphaRecords,
		"avg_price":      stats.AvgPrice,
		"avg_volatility": stats.AvgVolatility,
	}, nil
}
