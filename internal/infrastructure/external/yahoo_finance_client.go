package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/pkg/logger"
)

// Custom error types for better error handling
var (
	ErrInvalidSymbol    = errors.New("invalid or empty symbol")
	ErrAPIQuotaExceeded = errors.New("API quota exceeded")
	ErrSymbolNotFound   = errors.New("symbol not found")
	ErrAPIUnavailable   = errors.New("API service unavailable")
	ErrInvalidResponse  = errors.New("invalid API response")
)

// ClientConfig holds configuration for external API clients
type ClientConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	UserAgent      string
	RateLimitDelay time.Duration
}

// DefaultClientConfig returns default configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:        15 * time.Second,
		MaxRetries:     3,
		RetryDelay:     1 * time.Second,
		UserAgent:      "Mozilla/5.0 (compatible; StockTracker/1.0)",
		RateLimitDelay: 100 * time.Millisecond,
	}
}

// YahooFinanceClient defines the interface for Yahoo Finance API integration
type YahooFinanceClient interface {
	GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error)
	GetHistoricalData(ctx context.Context, symbol string, period string) ([]HistoricalDataPoint, error)
}

// yahooFinanceClient implements the Yahoo Finance API client
type yahooFinanceClient struct {
	baseURL    string
	httpClient *http.Client
	logger     logger.Logger
	config     *ClientConfig
}

// YahooQuoteResponse represents the response structure from Yahoo Finance API
type YahooQuoteResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency                 string  `json:"currency"`
				Symbol                   string  `json:"symbol"`
				RegularMarketPrice       float64 `json:"regularMarketPrice"`
				PreviousClose            float64 `json:"previousClose"`
				RegularMarketVolume      int64   `json:"regularMarketVolume"`
				MarketCap                int64   `json:"marketCap"`
				TrailingPE               float64 `json:"trailingPE"`
				DividendYield            float64 `json:"dividendYield"`
				FiftyTwoWeekHigh         float64 `json:"fiftyTwoWeekHigh"`
				FiftyTwoWeekLow          float64 `json:"fiftyTwoWeekLow"`
				AverageDailyVolume3Month int64   `json:"averageDailyVolume3Month"`
			} `json:"meta"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

// HistoricalDataPoint represents a single historical data point
type HistoricalDataPoint struct {
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int64     `json:"volume"`
}

// NewYahooFinanceClient creates a new Yahoo Finance client instance with optional configuration
func NewYahooFinanceClient(logger logger.Logger, config ...*ClientConfig) YahooFinanceClient {
	clientConfig := DefaultClientConfig()
	if len(config) > 0 && config[0] != nil {
		clientConfig = config[0]
	}

	return &yahooFinanceClient{
		baseURL: "https://query1.finance.yahoo.com/v8/finance/chart/",
		httpClient: &http.Client{
			Timeout: clientConfig.Timeout,
		},
		logger: logger,
		config: clientConfig,
	}
}

// validateSymbol validates the input symbol
func (c *yahooFinanceClient) validateSymbol(symbol string) error {
	if symbol == "" {
		return ErrInvalidSymbol
	}

	symbol = strings.TrimSpace(symbol)
	if len(symbol) == 0 || len(symbol) > 10 {
		return ErrInvalidSymbol
	}

	// Basic validation for special characters
	for _, char := range symbol {
		if (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '.' && char != '-' {
			return ErrInvalidSymbol
		}
	}

	return nil
}

// executeWithRetry executes an HTTP request with retry logic
func (c *yahooFinanceClient) executeWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Rate limiting delay between retries
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.config.RetryDelay * time.Duration(attempt)):
			}

			c.logger.Debug("Retrying request",
				"attempt", attempt,
				"url", req.URL.String())
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request attempt %d failed: %w", attempt+1, err)
			continue
		}

		// Success cases
		if resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		// Handle specific HTTP status codes
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			_ = resp.Body.Close()
			lastErr = ErrAPIQuotaExceeded
			continue
		case http.StatusNotFound:
			_ = resp.Body.Close()
			return nil, ErrSymbolNotFound
		case http.StatusServiceUnavailable, http.StatusBadGateway:
			_ = resp.Body.Close()
			lastErr = ErrAPIUnavailable
			continue
		default:
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("API returned status %d", resp.StatusCode)
			continue
		}
	}

	return nil, fmt.Errorf("all retry attempts failed: %w", lastErr)
}

// GetQuote retrieves real-time quote data for a given symbol
func (c *yahooFinanceClient) GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error) {
	// Validate input
	if err := c.validateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("symbol validation failed: %w", err)
	}

	// Build URL
	requestURL := fmt.Sprintf("%s%s", c.baseURL, url.QueryEscape(strings.ToUpper(strings.TrimSpace(symbol))))

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "application/json")

	c.logger.Debug("Making Yahoo Finance API request",
		"symbol", symbol,
		"url", requestURL)

	// Execute request with retry logic
	resp, err := c.executeWithRetry(ctx, req)
	if err != nil {
		c.logger.Error("Failed to execute Yahoo Finance request",
			"symbol", symbol,
			"error", err)
		return nil, fmt.Errorf("failed to execute request for symbol %s: %w", symbol, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()

	// Parse response
	var yahooResp YahooQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		return nil, fmt.Errorf("failed to decode response for symbol %s: %w", symbol, err)
	}

	// Validate response structure
	if err := c.validateYahooResponse(&yahooResp, symbol); err != nil {
		return nil, err
	}

	// Extract and convert data
	result := yahooResp.Chart.Result[0].Meta
	return c.convertYahooDataToModel(&result, symbol)
}

// validateYahooResponse validates the Yahoo Finance API response
func (c *yahooFinanceClient) validateYahooResponse(resp *YahooQuoteResponse, symbol string) error {
	// Check for API errors
	if resp.Chart.Error != nil {
		c.logger.Error("Yahoo Finance API returned error",
			"symbol", symbol,
			"error", resp.Chart.Error)
		return fmt.Errorf("yahoo Finance API error for symbol %s: %v", symbol, resp.Chart.Error)
	}

	if len(resp.Chart.Result) == 0 {
		c.logger.Warn("No data found for symbol", "symbol", symbol)
		return fmt.Errorf("%w: %s", ErrSymbolNotFound, symbol)
	}

	return nil
}

// convertYahooDataToModel converts Yahoo Finance response to internal model
func (c *yahooFinanceClient) convertYahooDataToModel(result *struct {
	Currency                 string  `json:"currency"`
	Symbol                   string  `json:"symbol"`
	RegularMarketPrice       float64 `json:"regularMarketPrice"`
	PreviousClose            float64 `json:"previousClose"`
	RegularMarketVolume      int64   `json:"regularMarketVolume"`
	MarketCap                int64   `json:"marketCap"`
	TrailingPE               float64 `json:"trailingPE"`
	DividendYield            float64 `json:"dividendYield"`
	FiftyTwoWeekHigh         float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow          float64 `json:"fiftyTwoWeekLow"`
	AverageDailyVolume3Month int64   `json:"averageDailyVolume3Month"`
}, symbol string) (*model.ExternalStockData, error) {

	// Validate essential data
	if result.RegularMarketPrice <= 0 {
		return nil, fmt.Errorf("invalid price data for symbol %s", symbol)
	}

	if result.PreviousClose <= 0 {
		c.logger.Warn("Invalid previous close price",
			"symbol", symbol,
			"previous_close", result.PreviousClose)
		// Don't fail completely, but set day change to 0
		result.PreviousClose = result.RegularMarketPrice
	}

	// Calculate day change safely
	dayChange := result.RegularMarketPrice - result.PreviousClose
	dayChangePercent := 0.0
	if result.PreviousClose > 0 {
		dayChangePercent = (dayChange / result.PreviousClose) * 100
	}

	// Convert optional fields safely
	var pERatio, dividendYield, week52High, week52Low *float64
	var avgVolume *int64

	if result.TrailingPE > 0 {
		pERatio = &result.TrailingPE
	}
	if result.DividendYield > 0 {
		dividendYield = &result.DividendYield
	}
	if result.FiftyTwoWeekHigh > 0 {
		week52High = &result.FiftyTwoWeekHigh
	}
	if result.FiftyTwoWeekLow > 0 {
		week52Low = &result.FiftyTwoWeekLow
	}
	if result.AverageDailyVolume3Month > 0 {
		avgVolume = &result.AverageDailyVolume3Month
	}

	return &model.ExternalStockData{
		CurrentPrice:     result.RegularMarketPrice,
		DayChange:        dayChange,
		DayChangePercent: dayChangePercent,
		Volume:           result.RegularMarketVolume,
		MarketCap:        result.MarketCap,
		PERatio:          pERatio,
		DividendYield:    dividendYield,
		Week52High:       week52High,
		Week52Low:        week52Low,
		AvgVolume:        avgVolume,
		LastUpdated:      time.Now(),
	}, nil
}

// GetHistoricalData retrieves historical price data for trend analysis
func (c *yahooFinanceClient) GetHistoricalData(ctx context.Context, symbol string, period string) ([]HistoricalDataPoint, error) {
	// Validate input
	if err := c.validateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("symbol validation failed: %w", err)
	}

	if period == "" {
		return nil, errors.New("period cannot be empty")
	}

	// For now, return empty slice - this can be implemented later for advanced features
	c.logger.Info("Historical data requested",
		"symbol", symbol,
		"period", period)

	return []HistoricalDataPoint{}, nil
}

// AlphaVantageClient defines the interface for Alpha Vantage API integration
type AlphaVantageClient interface {
	GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error)
	GetCompanyOverview(ctx context.Context, symbol string) (*CompanyOverview, error)
}

// alphaVantageClient implements the Alpha Vantage API client
type alphaVantageClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     logger.Logger
}

// CompanyOverview represents fundamental data from Alpha Vantage
type CompanyOverview struct {
	Symbol          string  `json:"Symbol"`
	Name            string  `json:"Name"`
	MarketCap       int64   `json:"MarketCapitalization,string"`
	PERatio         float64 `json:"PERatio,string"`
	DividendYield   float64 `json:"DividendYield,string"`
	EPS             float64 `json:"EPS,string"`
	RevenuePerShare float64 `json:"RevenuePerShareTTM,string"`
	ProfitMargin    float64 `json:"ProfitMargin,string"`
}

// NewAlphaVantageClient creates a new Alpha Vantage client instance
func NewAlphaVantageClient(apiKey string, logger logger.Logger) AlphaVantageClient {
	return &alphaVantageClient{
		baseURL: "https://www.alphavantage.co/query",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: logger,
	}
}

// GetQuote retrieves quote data from Alpha Vantage
func (c *alphaVantageClient) GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error) {
	// Build URL
	params := url.Values{}
	params.Set("function", "GLOBAL_QUOTE")
	params.Set("symbol", symbol)
	params.Set("apikey", c.apiKey)

	requestURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for symbol %s", resp.StatusCode, symbol)
	}

	// For now, return a basic implementation
	// Full implementation would parse Alpha Vantage response format
	c.logger.Info("Alpha Vantage quote requested", "symbol", symbol)

	return &model.ExternalStockData{
		LastUpdated: time.Now(),
	}, nil
}

// GetCompanyOverview retrieves fundamental company data
func (c *alphaVantageClient) GetCompanyOverview(ctx context.Context, symbol string) (*CompanyOverview, error) {
	// Build URL
	params := url.Values{}
	params.Set("function", "OVERVIEW")
	params.Set("symbol", symbol)
	params.Set("apikey", c.apiKey)

	requestURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for symbol %s", resp.StatusCode, symbol)
	}

	// Parse response
	var overview CompanyOverview
	if err := json.NewDecoder(resp.Body).Decode(&overview); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &overview, nil
}
