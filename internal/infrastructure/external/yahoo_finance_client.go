package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/pkg/logger"
)

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

// NewYahooFinanceClient creates a new Yahoo Finance client instance
func NewYahooFinanceClient(logger logger.Logger) YahooFinanceClient {
	return &yahooFinanceClient{
		baseURL: "https://query1.finance.yahoo.com/v8/finance/chart/",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GetQuote retrieves real-time quote data for a given symbol
func (c *yahooFinanceClient) GetQuote(ctx context.Context, symbol string) (*model.ExternalStockData, error) {
	// Build URL
	requestURL := fmt.Sprintf("%s%s", c.baseURL, url.QueryEscape(symbol))

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for symbol %s", resp.StatusCode, symbol)
	}

	// Parse response
	var yahooResp YahooQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API errors
	if yahooResp.Chart.Error != nil {
		return nil, fmt.Errorf("Yahoo Finance API error: %v", yahooResp.Chart.Error)
	}

	if len(yahooResp.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found for symbol %s", symbol)
	}

	// Extract data
	result := yahooResp.Chart.Result[0].Meta

	// Calculate day change
	dayChange := result.RegularMarketPrice - result.PreviousClose
	dayChangePercent := (dayChange / result.PreviousClose) * 100

	return &model.ExternalStockData{
		CurrentPrice:     result.RegularMarketPrice,
		DayChange:        dayChange,
		DayChangePercent: dayChangePercent,
		Volume:           result.RegularMarketVolume,
		MarketCap:        result.MarketCap,
		PERatio:          &result.TrailingPE,
		DividendYield:    &result.DividendYield,
		Week52High:       &result.FiftyTwoWeekHigh,
		Week52Low:        &result.FiftyTwoWeekLow,
		AvgVolume:        &result.AverageDailyVolume3Month,
		LastUpdated:      time.Now(),
	}, nil
}

// GetHistoricalData retrieves historical price data for trend analysis
func (c *yahooFinanceClient) GetHistoricalData(ctx context.Context, symbol string, period string) ([]HistoricalDataPoint, error) {
	// For now, return empty slice - this can be implemented later for advanced features
	c.logger.Info("Historical data requested", "symbol", symbol, "period", period)
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
