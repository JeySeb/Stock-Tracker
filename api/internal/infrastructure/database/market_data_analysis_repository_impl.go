package database

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
	"stock-tracker/internal/domain/marketdata/repositories"
	"stock-tracker/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// marketDataAnalysisRepositoryImpl implements the MarketDataAnalysisRepository interface
type marketDataAnalysisRepositoryImpl struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

// NewMarketDataAnalysisRepository creates a new market data analysis repository instance
func NewMarketDataAnalysisRepository(pool *pgxpool.Pool, logger logger.Logger) repositories.MarketDataAnalysisRepository {
	return &marketDataAnalysisRepositoryImpl{
		pool:   pool,
		logger: logger,
	}
}

// GetMarketDataByTicker retrieves market data for a specific ticker
func (r *marketDataAnalysisRepositoryImpl) GetMarketDataByTicker(ctx context.Context, ticker string, filters *model.MarketDataFilters) ([]*model.MarketData, error) {
	query := `
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, COALESCE(week_52_high, 0), COALESCE(week_52_low, 0), collected_at, data_timestamp, created_at, updated_at
		FROM market_data 
		WHERE ticker = $1
		ORDER BY data_timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, ticker, filters.Limit, filters.Offset)
	if err != nil {
		r.logger.Error("Failed to get market data by ticker", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to get market data by ticker: %w", err)
	}
	defer rows.Close()

	var marketDataList []*model.MarketData
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan market data", "error", err)
			return nil, fmt.Errorf("failed to scan market data: %w", err)
		}
		marketDataList = append(marketDataList, &data)
	}

	return marketDataList, nil
}

// GetLatestMarketDataByTicker retrieves the latest market data for a ticker
func (r *marketDataAnalysisRepositoryImpl) GetLatestMarketDataByTicker(ctx context.Context, ticker string) (*model.MarketData, error) {
	query := `
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, COALESCE(week_52_high, 0), COALESCE(week_52_low, 0), collected_at, data_timestamp, created_at, updated_at
		FROM market_data 
		WHERE ticker = $1 
		ORDER BY data_timestamp DESC 
		LIMIT 1
	`

	var data model.MarketData
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
		&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
		&data.Volume, &data.Week52High, &data.Week52Low,
		&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get latest market data", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	return &data, nil
}

// GetMarketDataAnalysis creates an analysis object from the latest market data
func (r *marketDataAnalysisRepositoryImpl) GetMarketDataAnalysis(ctx context.Context, ticker string) (*model.MarketDataAnalysis, error) {
	marketData, err := r.GetLatestMarketDataByTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}

	if marketData == nil {
		return nil, fmt.Errorf("no market data found for ticker %s", ticker)
	}

	analysis := model.NewMarketDataAnalysis(
		marketData.Ticker,
		marketData.CurrentPrice,
		marketData.DayChange,
		marketData.DayChangePercent,
		marketData.Volume,
		marketData.Week52High,
		marketData.Week52Low,
		marketData.DataTimestamp,
		marketData.CollectedAt,
		marketData.DataQuality,
		marketData.DataSource,
	)

	return analysis, nil
}

// GetMarketDataTrend calculates trend analysis for a ticker over a period
func (r *marketDataAnalysisRepositoryImpl) GetMarketDataTrend(ctx context.Context, ticker string, period string) (*model.MarketDataTrend, error) {
	var interval string
	var effectivePeriod string
	switch period {
	case "1d":
		interval = "1 day"
		effectivePeriod = "1d"
	case "1w":
		interval = "7 days"
		effectivePeriod = "1w"
	case "1m":
		interval = "30 days"
		effectivePeriod = "1m"
	case "3m":
		interval = "90 days"
		effectivePeriod = "3m"
	default:
		interval = "7 days"
		effectivePeriod = "1w"
	}

	// First, check if we have any data for this ticker in the requested period
	checkQuery := `
		SELECT COUNT(*) 
		FROM market_data 
		WHERE ticker = $1 
		AND data_timestamp >= now() - ($2::INTERVAL)
	`
	var count int
	err := r.pool.QueryRow(ctx, checkQuery, ticker, interval).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to check data availability", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to check data availability: %w", err)
	}

	// If no data in requested period, try with a longer period
	if count == 0 {
		// Try with 7 days if original period was shorter
		if interval == "1 day" {
			interval = "7 days"
			effectivePeriod = "1w"
			err = r.pool.QueryRow(ctx, checkQuery, ticker, interval).Scan(&count)
			if err != nil {
				r.logger.Error("Failed to check data availability with extended period", "ticker", ticker, "error", err)
				return nil, fmt.Errorf("failed to check data availability: %w", err)
			}
		}

		// If still no data, try with 30 days
		if count == 0 {
			interval = "30 days"
			effectivePeriod = "1m"
			err = r.pool.QueryRow(ctx, checkQuery, ticker, interval).Scan(&count)
			if err != nil {
				r.logger.Error("Failed to check data availability with 30-day period", "ticker", ticker, "error", err)
				return nil, fmt.Errorf("failed to check data availability: %w", err)
			}
		}
	}

	if count == 0 {
		return nil, fmt.Errorf("no trend data found for ticker %s in any available period", ticker)
	}

	// Get the first and last prices using a simpler approach
	query := `
		SELECT 
			$1 as ticker,
			MIN(current_price) as min_price,
			MAX(current_price) as max_price,
			ROUND(AVG(volume)::NUMERIC(20,0)) as avg_volume,
			COUNT(*) as data_points,
			(SELECT current_price FROM market_data WHERE ticker = $1 AND data_timestamp >= now() - ($2::INTERVAL) ORDER BY data_timestamp ASC LIMIT 1) as start_price,
			(SELECT current_price FROM market_data WHERE ticker = $1 AND data_timestamp >= now() - ($2::INTERVAL) ORDER BY data_timestamp DESC LIMIT 1) as end_price
		FROM market_data 
		WHERE ticker = $1 
		AND data_timestamp >= now() - ($2::INTERVAL)
	`

	var trend model.MarketDataTrend
	err = r.pool.QueryRow(ctx, query, ticker, interval).Scan(
		&trend.Ticker, &trend.MinPrice, &trend.MaxPrice,
		&trend.AvgVolume, &trend.DataPoints, &trend.StartPrice, &trend.EndPrice,
	)

	if err != nil {
		r.logger.Error("Failed to get market data trend", "ticker", ticker, "error", err)
		return nil, fmt.Errorf("failed to get market data trend: %w", err)
	}

	// Calculate derived fields
	trend.Period = effectivePeriod
	trend.TotalChange = trend.EndPrice - trend.StartPrice
	if trend.StartPrice > 0 {
		trend.TotalChangePercent = (trend.TotalChange / trend.StartPrice) * 100
	}

	// Calculate volatility (standard deviation of price changes)
	volatilityQuery := `
		SELECT COALESCE(STDDEV_SAMP(day_change_percent), 0) as volatility
		FROM market_data 
		WHERE ticker = $1 
		AND data_timestamp >= now() - ($2::INTERVAL)
	`
	err = r.pool.QueryRow(ctx, volatilityQuery, ticker, interval).Scan(&trend.Volatility)
	if err != nil {
		trend.Volatility = 0
	}

	// Determine trend strength and direction
	trend.TrendStrength = r.determineTrendStrength(trend.TotalChangePercent, trend.Volatility)
	trend.Direction = r.determineTrendDirection(trend.TotalChangePercent)

	return &trend, nil
}

// GetMarketDataSummary provides summary statistics for all market data
func (r *marketDataAnalysisRepositoryImpl) GetMarketDataSummary(ctx context.Context, period string) (*model.MarketDataSummary, error) {
	// First, check if we have any data in the last day
	checkQuery := `
		SELECT COUNT(*) 
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
	`
	var totalCount int
	err := r.pool.QueryRow(ctx, checkQuery, "1 day").Scan(&totalCount)
	if err != nil {
		r.logger.Error("Failed to check data availability for summary", "error", err)
		return nil, fmt.Errorf("failed to check data availability: %w", err)
	}

	// If no data in last day, try with longer periods
	var timeInterval string
	if totalCount == 0 {
		// Try with 7 days
		checkQuery7d := `
			SELECT COUNT(*) 
			FROM market_data 
			WHERE data_timestamp >= now() - ($1::INTERVAL)
		`
		err = r.pool.QueryRow(ctx, checkQuery7d, "7 days").Scan(&totalCount)
		if err != nil {
			r.logger.Error("Failed to check data availability for summary with 7 days", "error", err)
			return nil, fmt.Errorf("failed to check data availability: %w", err)
		}
		timeInterval = "7 days"
	} else {
		timeInterval = "1 day"
	}

	// If still no data, try with 30 days
	if totalCount == 0 {
		checkQuery30d := `
			SELECT COUNT(*) 
			FROM market_data 
			WHERE data_timestamp >= now() - ($1::INTERVAL)
		`
		err = r.pool.QueryRow(ctx, checkQuery30d, "30 days").Scan(&totalCount)
		if err != nil {
			r.logger.Error("Failed to check data availability for summary with 30 days", "error", err)
			return nil, fmt.Errorf("failed to check data availability: %w", err)
		}
		timeInterval = "30 days"
	}

	if totalCount == 0 {
		// Return empty summary if no data
		summary := &model.MarketDataSummary{
			TotalRecords:        0,
			UniqueTickers:       0,
			AvgPrice:            0,
			AvgDayChange:        0,
			AvgDayChangePercent: 0,
			TotalVolume:         0,
			AvgVolume:           0,
			MostActiveTicker:    "N/A",
			BestPerformer:       "N/A",
			WorstPerformer:      "N/A",
			BullishCount:        0,
			BearishCount:        0,
			NeutralCount:        0,
			Period:              period,
			LastUpdated:         time.Now(),
		}
		return summary, nil
	}

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_records,
			COUNT(DISTINCT ticker) as unique_tickers,
			COALESCE(AVG(current_price), 0) as avg_price,
			COALESCE(AVG(day_change), 0) as avg_day_change,
			COALESCE(AVG(day_change_percent), 0) as avg_day_change_percent,
			COALESCE(SUM(volume), 0) as total_volume,
			COALESCE(ROUND(AVG(volume)::NUMERIC(20,0)), 0) as avg_volume,
			COUNT(*) FILTER (WHERE day_change_percent > 2.0) as bullish_count,
			COUNT(*) FILTER (WHERE day_change_percent < -2.0) as bearish_count,
			COUNT(*) FILTER (WHERE day_change_percent BETWEEN -2.0 AND 2.0) as neutral_count
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
	`)

	var summary model.MarketDataSummary
	err = r.pool.QueryRow(ctx, query, timeInterval).Scan(
		&summary.TotalRecords, &summary.UniqueTickers, &summary.AvgPrice,
		&summary.AvgDayChange, &summary.AvgDayChangePercent, &summary.TotalVolume,
		&summary.AvgVolume, &summary.BullishCount, &summary.BearishCount, &summary.NeutralCount,
	)

	if err != nil {
		r.logger.Error("Failed to get market data summary", "error", err)
		return nil, fmt.Errorf("failed to get market data summary: %w", err)
	}

	// Get best performer
	bestPerformerQuery := fmt.Sprintf(`
		SELECT ticker 
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
		ORDER BY day_change_percent DESC 
		LIMIT 1
	`)
	err = r.pool.QueryRow(ctx, bestPerformerQuery, timeInterval).Scan(&summary.BestPerformer)
	if err != nil {
		summary.BestPerformer = "N/A"
	}

	// Get worst performer
	worstPerformerQuery := fmt.Sprintf(`
		SELECT ticker 
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
		ORDER BY day_change_percent ASC 
		LIMIT 1
	`)
	err = r.pool.QueryRow(ctx, worstPerformerQuery, timeInterval).Scan(&summary.WorstPerformer)
	if err != nil {
		summary.WorstPerformer = "N/A"
	}

	// Get most active ticker
	mostActiveQuery := fmt.Sprintf(`
		SELECT ticker 
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
		GROUP BY ticker 
		ORDER BY AVG(volume) DESC 
		LIMIT 1
	`)
	err = r.pool.QueryRow(ctx, mostActiveQuery, timeInterval).Scan(&summary.MostActiveTicker)
	if err != nil {
		summary.MostActiveTicker = "N/A"
	}

	summary.Period = period
	summary.LastUpdated = time.Now()

	return &summary, nil
}

// GetTopPerformers retrieves the top performing tickers
func (r *marketDataAnalysisRepositoryImpl) GetTopPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	var interval string
	//var effectivePeriod string
	switch period {
	case "1d":
		interval = "1 day"
		//effectivePeriod = "1d"
	case "1w":
		interval = "7 days"
		//effectivePeriod = "1w"
	case "1m":
		interval = "30 days"
		//effectivePeriod = "1m"
	default:
		interval = "1 day"
		//effectivePeriod = "1d"
	}

	// First, check if we have any data in the requested period
	checkQuery := `
		SELECT COUNT(*) 
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
	`
	var count int
	err := r.pool.QueryRow(ctx, checkQuery, interval).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to check data availability for top performers", "error", err)
		return nil, fmt.Errorf("failed to check data availability: %w", err)
	}

	// If no data in requested period, try with longer periods
	if count == 0 {
		// Try with 7 days if original period was shorter
		if interval == "1 day" {
			interval = "7 days"
			//effectivePeriod = "1w"
			err = r.pool.QueryRow(ctx, checkQuery, interval).Scan(&count)
			if err != nil {
				r.logger.Error("Failed to check data availability for top performers with extended period", "error", err)
				return nil, fmt.Errorf("failed to check data availability: %w", err)
			}
		}

		// If still no data, try with 30 days
		if count == 0 {
			interval = "30 days"
			//effectivePeriod = "1m"
			err = r.pool.QueryRow(ctx, checkQuery, interval).Scan(&count)
			if err != nil {
				r.logger.Error("Failed to check data availability for top performers with 30-day period", "error", err)
				return nil, fmt.Errorf("failed to check data availability: %w", err)
			}
		}
	}

	if count == 0 {
		return []*model.MarketDataAnalysis{}, nil
	}

	// Use a simpler query that's compatible with CockroachDB
	query := `
		WITH latest_data AS (
			SELECT ticker, MAX(data_timestamp) as max_ts
			FROM market_data 
			WHERE data_timestamp >= now() - ($1::INTERVAL)
			GROUP BY ticker
		)
		SELECT 
			md.id, md.ticker, md.data_source, md.data_quality, md.current_price, md.day_change, md.day_change_percent,
			md.volume, COALESCE(md.week_52_high, 0), COALESCE(md.week_52_low, 0), md.collected_at, md.data_timestamp, md.created_at, md.updated_at
		FROM market_data md
		JOIN latest_data ld ON md.ticker = ld.ticker AND md.data_timestamp = ld.max_ts
		ORDER BY md.day_change_percent DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, interval, limit)
	if err != nil {
		r.logger.Error("Failed to get top performers", "error", err)
		return nil, fmt.Errorf("failed to get top performers: %w", err)
	}
	defer rows.Close()

	var analyses []*model.MarketDataAnalysis
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan top performer data", "error", err)
			continue
		}

		analysis := model.NewMarketDataAnalysis(
			data.Ticker, data.CurrentPrice, data.DayChange, data.DayChangePercent,
			data.Volume, data.Week52High, data.Week52Low, data.DataTimestamp,
			data.CollectedAt, data.DataQuality, data.DataSource,
		)
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetWorstPerformers retrieves the worst performing tickers
func (r *marketDataAnalysisRepositoryImpl) GetWorstPerformers(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	var interval string
	//var effectivePeriod string
	switch period {
	case "1d":
		interval = "1 day"
		//effectivePeriod = "1d"
	case "1w":
		interval = "7 days"
		//effectivePeriod = "1w"
	case "1m":
		interval = "30 days"
		//effectivePeriod = "1m"
	default:
		interval = "1 day"
		//effectivePeriod = "1d"
	}

	// First, check if we have any data in the requested period
	checkQuery := `
		SELECT COUNT(*) 
		FROM market_data 
		WHERE data_timestamp >= now() - ($1::INTERVAL)
	`
	var count int
	err := r.pool.QueryRow(ctx, checkQuery, interval).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to check data availability for worst performers", "error", err)
		return nil, fmt.Errorf("failed to check data availability: %w", err)
	}

	// If no data in requested period, try with longer periods
	if count == 0 {
		// Try with 7 days if original period was shorter
		if interval == "1 day" {
			interval = "7 days"
			//effectivePeriod = "1w"
			err = r.pool.QueryRow(ctx, checkQuery, interval).Scan(&count)
			if err != nil {
				r.logger.Error("Failed to check data availability for worst performers with extended period", "error", err)
				return nil, fmt.Errorf("failed to check data availability: %w", err)
			}
		}

		// If still no data, try with 30 days
		if count == 0 {
			interval = "30 days"
			//effectivePeriod = "1m"
			err = r.pool.QueryRow(ctx, checkQuery, interval).Scan(&count)
			if err != nil {
				r.logger.Error("Failed to check data availability for worst performers with 30-day period", "error", err)
				return nil, fmt.Errorf("failed to check data availability: %w", err)
			}
		}
	}

	if count == 0 {
		return []*model.MarketDataAnalysis{}, nil
	}

	// Use a simpler query that's compatible with CockroachDB
	query := `
		WITH latest_data AS (
			SELECT ticker, MAX(data_timestamp) as max_ts
			FROM market_data 
			WHERE data_timestamp >= now() - ($1::INTERVAL)
			GROUP BY ticker
		)
		SELECT 
			md.id, md.ticker, md.data_source, md.data_quality, md.current_price, md.day_change, md.day_change_percent,
			md.volume, COALESCE(md.week_52_high, 0), COALESCE(md.week_52_low, 0), md.collected_at, md.data_timestamp, md.created_at, md.updated_at
		FROM market_data md
		JOIN latest_data ld ON md.ticker = ld.ticker AND md.data_timestamp = ld.max_ts
		ORDER BY md.day_change_percent ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, interval, limit)
	if err != nil {
		r.logger.Error("Failed to get worst performers", "error", err)
		return nil, fmt.Errorf("failed to get worst performers: %w", err)
	}
	defer rows.Close()

	var analyses []*model.MarketDataAnalysis
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan worst performer data", "error", err)
			continue
		}

		analysis := model.NewMarketDataAnalysis(
			data.Ticker, data.CurrentPrice, data.DayChange, data.DayChangePercent,
			data.Volume, data.Week52High, data.Week52Low, data.DataTimestamp,
			data.CollectedAt, data.DataQuality, data.DataSource,
		)
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// Helper methods for trend analysis
func (r *marketDataAnalysisRepositoryImpl) determineTrendStrength(changePercent, volatility float64) string {
	// Use z-score based thresholds
	zScore := math.Abs(changePercent / volatility)
	if zScore > 2.0 {
		return "strong"
	} else if zScore > 1.0 {
		return "moderate"
	}
	return "weak"
}

func (r *marketDataAnalysisRepositoryImpl) determineTrendDirection(changePercent float64) string {
	if changePercent > 2.0 {
		return "up"
	} else if changePercent < -2.0 {
		return "down"
	}
	return "sideways"
}

// GetMarketDataByDateRange retrieves market data within a date range
func (r *marketDataAnalysisRepositoryImpl) GetMarketDataByDateRange(ctx context.Context, startDate, endDate time.Time, filters *model.MarketDataFilters) ([]*model.MarketData, error) {
	query := `
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, COALESCE(week_52_high, 0), COALESCE(week_52_low, 0), collected_at, data_timestamp, created_at, updated_at
		FROM market_data 
		WHERE data_timestamp BETWEEN $1 AND $2
		ORDER BY data_timestamp DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, startDate, endDate, filters.Limit, filters.Offset)
	if err != nil {
		r.logger.Error("Failed to get market data by date range", "error", err)
		return nil, fmt.Errorf("failed to get market data by date range: %w", err)
	}
	defer rows.Close()

	var marketDataList []*model.MarketData
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan market data", "error", err)
			return nil, fmt.Errorf("failed to scan market data: %w", err)
		}
		marketDataList = append(marketDataList, &data)
	}

	return marketDataList, nil
}

// GetMarketDataComparison compares market data between two tickers
func (r *marketDataAnalysisRepositoryImpl) GetMarketDataComparison(ctx context.Context, ticker1, ticker2 string, date time.Time) (*model.MarketDataComparison, error) {
	query := `
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, COALESCE(week_52_high, 0), COALESCE(week_52_low, 0), collected_at, data_timestamp, created_at, updated_at
		FROM market_data 
		WHERE ticker IN ($1, $2)
		AND DATE(data_timestamp) = DATE($3)
		ORDER BY data_timestamp DESC
	`

	rows, err := r.pool.Query(ctx, query, ticker1, ticker2, date)
	if err != nil {
		r.logger.Error("Failed to get market data comparison", "error", err)
		return nil, fmt.Errorf("failed to get market data comparison: %w", err)
	}
	defer rows.Close()

	var data1, data2 *model.MarketData
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan market data", "error", err)
			return nil, fmt.Errorf("failed to scan market data: %w", err)
		}

		if data.Ticker == ticker1 {
			data1 = &data
		} else {
			data2 = &data
		}
	}

	if data1 == nil || data2 == nil {
		return nil, fmt.Errorf("market data not found for one or both tickers")
	}

	comparison := &model.MarketDataComparison{
		Ticker1:          ticker1,
		Ticker2:          ticker2,
		ComparisonDate:   date,
		Price1:           data1.CurrentPrice,
		Price2:           data2.CurrentPrice,
		Change1:          data1.DayChange,
		Change2:          data2.DayChange,
		ChangePercent1:   data1.DayChangePercent,
		ChangePercent2:   data2.DayChangePercent,
		Volume1:          data1.Volume,
		Volume2:          data2.Volume,
		RelativeStrength: data1.DayChangePercent / data2.DayChangePercent,
		Correlation:      0, // Would need historical data to calculate correlation
	}

	return comparison, nil
}

// GetMostVolatile retrieves the most volatile tickers
func (r *marketDataAnalysisRepositoryImpl) GetMostVolatile(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	var interval string
	switch period {
	case "1d":
		interval = "1 day"
	case "1w":
		interval = "1 week"
	case "1m":
		interval = "1 month"
	default:
		interval = "1 day"
	}

	query := `
		WITH volatility_calc AS (
			SELECT 
				ticker,
				STDDEV(day_change_percent) as volatility,
				MAX(data_timestamp) as latest_timestamp
			FROM market_data
			WHERE data_timestamp >= now() - INTERVAL $1
			GROUP BY ticker
			ORDER BY volatility DESC
			LIMIT $2
		)
		SELECT 
			md.id, md.ticker, md.data_source, md.data_quality, md.current_price, 
			md.day_change, md.day_change_percent, md.volume, md.week_52_high, md.week_52_low,
			md.collected_at, md.data_timestamp, md.created_at, md.updated_at
		FROM market_data md
		INNER JOIN volatility_calc vc ON md.ticker = vc.ticker AND md.data_timestamp = vc.latest_timestamp
	`

	rows, err := r.pool.Query(ctx, query, interval, limit)
	if err != nil {
		r.logger.Error("Failed to get most volatile tickers", "error", err)
		return nil, fmt.Errorf("failed to get most volatile tickers: %w", err)
	}
	defer rows.Close()

	var analyses []*model.MarketDataAnalysis
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan volatile ticker data", "error", err)
			continue
		}

		analysis := model.NewMarketDataAnalysis(
			data.Ticker, data.CurrentPrice, data.DayChange, data.DayChangePercent,
			data.Volume, data.Week52High, data.Week52Low, data.DataTimestamp,
			data.CollectedAt, data.DataQuality, data.DataSource,
		)
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetMostActive retrieves the most active tickers by volume
func (r *marketDataAnalysisRepositoryImpl) GetMostActive(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	var interval string
	switch period {
	case "1d":
		interval = "1 day"
	case "1w":
		interval = "1 week"
	case "1m":
		interval = "1 month"
	default:
		interval = "1 day"
	}

	query := `
		WITH volume_calc AS (
			SELECT 
				ticker,
				AVG(volume) as avg_volume,
				MAX(data_timestamp) as latest_timestamp
			FROM market_data
			WHERE data_timestamp >= now() - INTERVAL $1
			GROUP BY ticker
			ORDER BY avg_volume DESC
			LIMIT $2
		)
		SELECT 
			md.id, md.ticker, md.data_source, md.data_quality, md.current_price, 
			md.day_change, md.day_change_percent, md.volume, md.week_52_high, md.week_52_low,
			md.collected_at, md.data_timestamp, md.created_at, md.updated_at
		FROM market_data md
		INNER JOIN volume_calc vc ON md.ticker = vc.ticker AND md.data_timestamp = vc.latest_timestamp
	`

	rows, err := r.pool.Query(ctx, query, interval, limit)
	if err != nil {
		r.logger.Error("Failed to get most active tickers", "error", err)
		return nil, fmt.Errorf("failed to get most active tickers: %w", err)
	}
	defer rows.Close()

	var analyses []*model.MarketDataAnalysis
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan active ticker data", "error", err)
			continue
		}

		analysis := model.NewMarketDataAnalysis(
			data.Ticker, data.CurrentPrice, data.DayChange, data.DayChangePercent,
			data.Volume, data.Week52High, data.Week52Low, data.DataTimestamp,
			data.CollectedAt, data.DataQuality, data.DataSource,
		)
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GetDataQualityStats retrieves statistics about data quality
func (r *marketDataAnalysisRepositoryImpl) GetDataQualityStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			data_quality,
			COUNT(*) as count,
			COUNT(DISTINCT ticker) as unique_tickers,
			AVG(EXTRACT(EPOCH FROM (now() - data_timestamp))) as avg_age_seconds
		FROM market_data
		GROUP BY data_quality
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.logger.Error("Failed to get data quality stats", "error", err)
		return nil, fmt.Errorf("failed to get data quality stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]interface{})
	for rows.Next() {
		var quality string
		var count, uniqueTickers int
		var avgAge float64

		err := rows.Scan(&quality, &count, &uniqueTickers, &avgAge)
		if err != nil {
			r.logger.Error("Failed to scan data quality stats", "error", err)
			continue
		}

		stats[quality] = map[string]interface{}{
			"count":          count,
			"unique_tickers": uniqueTickers,
			"avg_age":        avgAge,
		}
	}
	return stats, nil
}

func (r *marketDataAnalysisRepositoryImpl) GetCorrelatedTickers(ctx context.Context, ticker string, threshold float64, period string) ([]string, error) {
	// Implementation for correlated tickers
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetPriceBreakouts(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	// Implementation for price breakouts
	return nil, fmt.Errorf("not implemented")
}

// GetDataSourceStats retrieves statistics about data sources
func (r *marketDataAnalysisRepositoryImpl) GetDataSourceStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			data_source,
			COUNT(*) as count,
			COUNT(DISTINCT ticker) as unique_tickers,
			AVG(EXTRACT(EPOCH FROM (now() - data_timestamp))) as avg_age_seconds,
			COUNT(*) FILTER (WHERE data_quality = 'high') as high_quality_count
		FROM market_data
		GROUP BY data_source
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.logger.Error("Failed to get data source stats", "error", err)
		return nil, fmt.Errorf("failed to get data source stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]interface{})
	for rows.Next() {
		var source string
		var count, uniqueTickers, highQualityCount int
		var avgAge float64

		err := rows.Scan(&source, &count, &uniqueTickers, &avgAge, &highQualityCount)
		if err != nil {
			r.logger.Error("Failed to scan data source stats", "error", err)
			continue
		}

		stats[source] = map[string]interface{}{
			"count":              count,
			"unique_tickers":     uniqueTickers,
			"avg_age":            avgAge,
			"high_quality_count": highQualityCount,
		}
	}

	return stats, nil
}

func (r *marketDataAnalysisRepositoryImpl) GetVolumeSurges(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	// Implementation for volume surges
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetHighRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error) {
	query := `
		WITH latest_data AS (
			SELECT DISTINCT ON (ticker)
				id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
				volume, week_52_high, week_52_low, collected_at, data_timestamp, created_at, updated_at
			FROM market_data 
			WHERE data_timestamp >= now() - INTERVAL '1 day'
			ORDER BY ticker, data_timestamp DESC
		),
		risk_calc AS (
			SELECT *,
				ABS(day_change_percent) + 
				CASE 
					WHEN current_price > 0 THEN (ABS(current_price - ((week_52_high + week_52_low) / 2)) / current_price) * 100
					ELSE 0
				END as risk_score
			FROM latest_data
		)
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, week_52_high, week_52_low, collected_at, data_timestamp, created_at, updated_at
		FROM risk_calc
		ORDER BY risk_score DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		r.logger.Error("Failed to get high risk tickers", "error", err)
		return nil, fmt.Errorf("failed to get high risk tickers: %w", err)
	}
	defer rows.Close()

	var analyses []*model.MarketDataAnalysis
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan high risk ticker data", "error", err)
			continue
		}

		analysis := model.NewMarketDataAnalysis(
			data.Ticker, data.CurrentPrice, data.DayChange, data.DayChangePercent,
			data.Volume, data.Week52High, data.Week52Low, data.DataTimestamp,
			data.CollectedAt, data.DataQuality, data.DataSource,
		)
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

func (r *marketDataAnalysisRepositoryImpl) GetLowRiskTickers(ctx context.Context, limit int) ([]*model.MarketDataAnalysis, error) {
	query := `
		WITH latest_data AS (
			SELECT DISTINCT ON (ticker)
				id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
				volume, week_52_high, week_52_low, collected_at, data_timestamp, created_at, updated_at
			FROM market_data 
			WHERE data_timestamp >= now() - INTERVAL '1 day'
			ORDER BY ticker, data_timestamp DESC
		),
		risk_calc AS (
			SELECT *,
				ABS(day_change_percent) + 
				CASE 
					WHEN current_price > 0 THEN (ABS(current_price - ((week_52_high + week_52_low) / 2)) / current_price) * 100
					ELSE 0
				END as risk_score
			FROM latest_data
		)
		SELECT 
			id, ticker, data_source, data_quality, current_price, day_change, day_change_percent,
			volume, week_52_high, week_52_low, collected_at, data_timestamp, created_at, updated_at
		FROM risk_calc
		ORDER BY risk_score ASC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		r.logger.Error("Failed to get low risk tickers", "error", err)
		return nil, fmt.Errorf("failed to get low risk tickers: %w", err)
	}
	defer rows.Close()

	var analyses []*model.MarketDataAnalysis
	for rows.Next() {
		var data model.MarketData
		err := rows.Scan(
			&data.ID, &data.Ticker, &data.DataSource, &data.DataQuality,
			&data.CurrentPrice, &data.DayChange, &data.DayChangePercent,
			&data.Volume, &data.Week52High, &data.Week52Low,
			&data.CollectedAt, &data.DataTimestamp, &data.CreatedAt, &data.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan low risk ticker data", "error", err)
			continue
		}

		analysis := model.NewMarketDataAnalysis(
			data.Ticker, data.CurrentPrice, data.DayChange, data.DayChangePercent,
			data.Volume, data.Week52High, data.Week52Low, data.DataTimestamp,
			data.CollectedAt, data.DataQuality, data.DataSource,
		)
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

func (r *marketDataAnalysisRepositoryImpl) GetRiskDistribution(ctx context.Context) (map[string]int, error) {
	// Implementation for risk distribution
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetTrendingTickers(ctx context.Context, direction string, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	// Implementation for trending tickers
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetTrendStrength(ctx context.Context, ticker string, period string) (string, error) {
	// Implementation for trend strength
	return "", fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetTrendReversal(ctx context.Context, limit int, period string) ([]*model.MarketDataAnalysis, error) {
	// Implementation for trend reversals
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetPriceStatistics(ctx context.Context, ticker string, period string) (map[string]float64, error) {
	// Implementation for price statistics
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetVolumeStatistics(ctx context.Context, ticker string, period string) (map[string]int64, error) {
	// Implementation for volume statistics
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetVolatilityStatistics(ctx context.Context, ticker string, period string) (map[string]float64, error) {
	// Implementation for volatility statistics
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetMarketAlerts(ctx context.Context, alertType string, severity string, limit int) ([]*model.MarketDataAlert, error) {
	// Implementation for market alerts
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) CreateMarketAlert(ctx context.Context, alert *model.MarketDataAlert) error {
	// Implementation for creating market alerts
	return fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) UpdateMarketAlert(ctx context.Context, alert *model.MarketDataAlert) error {
	// Implementation for updating market alerts
	return fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) DeleteMarketAlert(ctx context.Context, alertID string) error {
	// Implementation for deleting market alerts
	return fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetDataFreshness(ctx context.Context) (map[string]time.Time, error) {
	// Implementation for data freshness
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetMarketDataWithStockAnalysis(ctx context.Context, ticker string) (map[string]interface{}, error) {
	// Implementation for cross-domain analysis
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetCorrelationWithBrokerActions(ctx context.Context, ticker string, period string) (map[string]interface{}, error) {
	// Implementation for broker action correlation
	return nil, fmt.Errorf("not implemented")
}

func (r *marketDataAnalysisRepositoryImpl) GetMarketDataImpactOnRecommendations(ctx context.Context, ticker string) (map[string]interface{}, error) {
	// Implementation for recommendation impact analysis
	return nil, fmt.Errorf("not implemented")
}
