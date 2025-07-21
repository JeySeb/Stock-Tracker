package clients

import (
	"context"
	"encoding/json"

	"fmt"
	"net/http"
	"strings"
	"time"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/pkg/logger"

	"github.com/hashicorp/go-retryablehttp"
)

type StockAPIResponse struct {
	Items    []StockAPIItem `json:"items"`
	NextPage string         `json:"next_page"`
}

type StockAPIItem struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

type StockAPIClient interface {
	FetchAllStocks(ctx context.Context) ([]*entities.Stock, error)
	FetchPage(ctx context.Context, nextPage string) ([]*entities.Stock, string, error)
}

type stockAPIClient struct {
	client    *retryablehttp.Client
	baseURL   string
	apiKey    string
	logger    logger.Logger
	rateLimit time.Duration
}

func NewStockAPIClient(baseURL, apiKey string, logger logger.Logger) StockAPIClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second
	retryClient.Logger = nil // Disable logging for retryablehttp, because we have our own logger

	return &stockAPIClient{
		client:    retryClient,
		baseURL:   baseURL,
		apiKey:    apiKey,
		logger:    logger,
		rateLimit: 100 * time.Millisecond,
	}
}

func (c *stockAPIClient) FetchAllStocks(ctx context.Context) ([]*entities.Stock, error) {
	var allStocks []*entities.Stock
	nextPage := ""
	pageCount := 0

	c.logger.Info("Starting to fetch all stocks from the API")

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		pageStocks, next, err := c.FetchPage(ctx, nextPage)
		if err != nil {
			c.logger.Error("Failed to fetch page", "page", pageCount, "nextPage", next, "error", err)
			return nil, fmt.Errorf("failed to fetch page %d: %w", pageCount, err)
		}

		allStocks = append(allStocks, pageStocks...)
		pageCount++

		c.logger.Info(fmt.Sprintf("Successfully fetched page %d with %d stocks. Next page: %s", pageCount, len(pageStocks), next))

		if next == "" {
			break
		}
		nextPage = next

		//Rate limit
		time.Sleep(c.rateLimit)

	}

	c.logger.Info("Completed stock data ingestion", "total_stocks", len(allStocks), "pages", pageCount)
	return allStocks, nil
}

func (c *stockAPIClient) FetchPage(ctx context.Context, nextPage string) ([]*entities.Stock, string, error) {
	url := c.baseURL
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Stock-Tracker/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResponse StockAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	stocks := make([]*entities.Stock, 0, len(apiResponse.Items))
	for _, item := range apiResponse.Items {
		stock, err := c.convertAPIItemToStock(item)
		if err != nil {
			c.logger.Warn("Failed to convert API item to stock", "ticker", item.Ticker, "error", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, apiResponse.NextPage, nil
}

func (c *stockAPIClient) convertAPIItemToStock(item StockAPIItem) (*entities.Stock, error) {

	eventTime, err := time.Parse(time.RFC3339, item.Time)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event time %s: %w", item.Time, err)
	}

	stock := entities.NewStock(item.Ticker, item.Company, item.Brokerage, item.Action, eventTime)
	stock.RatingFrom = item.RatingFrom
	stock.RatingTo = item.RatingTo

	// Parse the target prices
	if targetFrom := c.parsePrice(item.TargetFrom); targetFrom > 0 {
		stock.TargetFrom = targetFrom
	}
	if targetTo := c.parsePrice(item.TargetTo); targetTo > 0 {
		stock.TargetTo = targetTo
	}

	return stock, nil
}

func (c *stockAPIClient) parsePrice(priceStr string) float64 {
	if priceStr == "" {
		return 0
	}

	// Remove currency symbols and spaces
	cleaned := strings.ReplaceAll(priceStr, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	var price float64
	if _, err := fmt.Sscanf(cleaned, "%f", &price); err != nil {
		return 0
	}

	return price
}
