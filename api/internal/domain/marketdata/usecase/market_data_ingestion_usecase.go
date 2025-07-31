package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"stock-tracker/internal/domain/marketdata/model"
	"stock-tracker/internal/domain/marketdata/repositories"
	"stock-tracker/internal/infrastructure/external"
	"stock-tracker/pkg/logger"

	"github.com/google/uuid"
)

// MarketDataIngestionUseCase handles the market data ingestion process
type MarketDataIngestionUseCase struct {
	marketDataRepo repositories.MarketDataRepository
	yahooClient    external.YahooFinanceClient
	logger         logger.Logger
}

// NewMarketDataIngestionUseCase creates a new market data ingestion use case
func NewMarketDataIngestionUseCase(
	marketDataRepo repositories.MarketDataRepository,
	yahooClient external.YahooFinanceClient,
	logger logger.Logger,
) *MarketDataIngestionUseCase {
	return &MarketDataIngestionUseCase{
		marketDataRepo: marketDataRepo,
		yahooClient:    yahooClient,
		logger:         logger,
	}
}

// IngestMarketData performs the market data ingestion process
func (uc *MarketDataIngestionUseCase) IngestMarketData(ctx context.Context) error {
	// Generate batch ID for this ingestion run
	batchID := uuid.New().String()

	uc.logger.Info("Starting market data ingestion",
		"batch_id", batchID,
		"data_source", model.DataSourceYahooFinance)

	// Get unique tickers from stocks table
	tickers, err := uc.marketDataRepo.GetUniqueTickers(ctx)
	if err != nil {
		uc.logger.Error("Failed to get unique tickers", "error", err)
		return fmt.Errorf("failed to get unique tickers: %w", err)
	}

	if len(tickers) == 0 {
		uc.logger.Warn("No tickers found for market data ingestion")
		return nil
	}

	// Create ingestion log
	ingestionLog := model.NewMarketDataIngestionLog(batchID, model.DataSourceYahooFinance, len(tickers))

	// Save initial log entry
	if err := uc.marketDataRepo.SaveIngestionLog(ctx, ingestionLog); err != nil {
		uc.logger.Error("Failed to save initial ingestion log", "error", err)
		return fmt.Errorf("failed to save initial ingestion log: %w", err)
	}

	// Process tickers with concurrency control
	uc.logger.Info("Processing tickers for market data ingestion",
		"batch_id", batchID,
		"total_tickers", len(tickers))

	// Use worker pool for concurrent processing
	results := uc.processTickersConcurrently(ctx, tickers, batchID)

	// Update ingestion log with results
	ingestionLog.SuccessfulTickers = results.successful
	ingestionLog.FailedTickers = results.failed
	ingestionLog.SkippedTickers = results.skipped
	ingestionLog.Status = "completed"
	now := time.Now()
	ingestionLog.CompletedAt = &now
	ingestionLog.UpdatedAt = now

	if results.failed > 0 {
		ingestionLog.Status = "failed"
		ingestionLog.ErrorDetails = map[string]interface{}{
			"failed_tickers": results.failedTickers,
			"errors":         results.errors,
		}
	}

	// Update final log entry
	if err := uc.marketDataRepo.UpdateIngestionLog(ctx, ingestionLog); err != nil {
		uc.logger.Error("Failed to update final ingestion log", "error", err)
		return fmt.Errorf("failed to update final ingestion log: %w", err)
	}

	uc.logger.Info("Market data ingestion completed",
		"batch_id", batchID,
		"successful", results.successful,
		"failed", results.failed,
		"skipped", results.skipped)

	return nil
}

// processResult holds the results of processing tickers
type processResult struct {
	successful    int
	failed        int
	skipped       int
	failedTickers []string
	errors        []string
}

// processTickersConcurrently processes tickers with concurrency control
func (uc *MarketDataIngestionUseCase) processTickersConcurrently(ctx context.Context, tickers []string, batchID string) processResult {
	const maxWorkers = 10 // Limit concurrent API calls

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxWorkers)

	// Channels for collecting results
	successCh := make(chan string, len(tickers))
	failCh := make(chan struct {
		ticker string
		error  string
	}, len(tickers))
	skipCh := make(chan string, len(tickers))

	// Start workers
	for _, ticker := range tickers {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			if err := uc.processTicker(ctx, t, batchID); err != nil {
				failCh <- struct {
					ticker string
					error  string
				}{t, err.Error()}
			} else {
				successCh <- t
			}
		}(ticker)
	}

	// Wait for all workers to complete
	wg.Wait()
	close(successCh)
	close(failCh)
	close(skipCh)

	// Collect results
	var successful, failed, skipped int
	var failedTickers []string
	var errors []string

	// Count successful
	for range successCh {
		successful++
	}

	// Count failed
	for result := range failCh {
		failed++
		failedTickers = append(failedTickers, result.ticker)
		errors = append(errors, result.error)
	}

	// Count skipped
	for range skipCh {
		skipped++
	}

	return processResult{
		successful:    successful,
		failed:        failed,
		skipped:       skipped,
		failedTickers: failedTickers,
		errors:        errors,
	}
}

// processTicker processes a single ticker
func (uc *MarketDataIngestionUseCase) processTicker(ctx context.Context, ticker, batchID string) error {
	uc.logger.Debug("Processing ticker for market data",
		"batch_id", batchID,
		"ticker", ticker)

	// Get market data from Yahoo Finance (returns recommendation model)
	externalData, err := uc.yahooClient.GetQuote(ctx, ticker)
	if err != nil {
		uc.logger.Error("Failed to get quote from Yahoo Finance",
			"batch_id", batchID,
			"ticker", ticker,
			"error", err)
		return fmt.Errorf("failed to get quote for %s: %w", ticker, err)
	}

	// Check if data already exists for this timestamp
	exists, err := uc.marketDataRepo.ExistsMarketData(ctx, ticker, model.DataSourceYahooFinance, externalData.LastUpdated)
	if err != nil {
		uc.logger.Error("Failed to check existing market data",
			"batch_id", batchID,
			"ticker", ticker,
			"error", err)
		return fmt.Errorf("failed to check existing market data for %s: %w", ticker, err)
	}

	if exists {
		uc.logger.Debug("Market data already exists, skipping",
			"batch_id", batchID,
			"ticker", ticker,
			"timestamp", externalData.LastUpdated)
		return nil
	}

	// Convert recommendation model to market data model
	marketData := model.ConvertFromRecommendationModel(ticker, model.DataSourceYahooFinance, externalData)

	// Save market data
	if err := uc.marketDataRepo.SaveMarketData(ctx, marketData); err != nil {
		uc.logger.Error("Failed to save market data",
			"batch_id", batchID,
			"ticker", ticker,
			"error", err)
		return fmt.Errorf("failed to save market data for %s: %w", ticker, err)
	}

	uc.logger.Debug("Successfully processed ticker",
		"batch_id", batchID,
		"ticker", ticker,
		"price", externalData.CurrentPrice)

	return nil
}

// GetMarketDataStats retrieves market data statistics
func (uc *MarketDataIngestionUseCase) GetMarketDataStats(ctx context.Context, days int) (map[string]interface{}, error) {
	return uc.marketDataRepo.GetMarketDataStats(ctx, days)
}
