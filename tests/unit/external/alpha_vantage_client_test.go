package external

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"stock-tracker/internal/domain/recommendation/model"
	"stock-tracker/internal/infrastructure/external"
)

func TestAlphaVantageClient_GetQuote_Success(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Contains(t, r.URL.String(), "AAPL")
		assert.Contains(t, r.URL.String(), "apikey=test_key")

		// Return successful response
		response := map[string]interface{}{
			"Global Quote": map[string]interface{}{
				"01. symbol":             "AAPL",
				"02. open":               "170.50",
				"03. high":               "172.00",
				"04. low":                "169.50",
				"05. price":              "171.00",
				"06. volume":             "1000000",
				"07. latest trading day": "2023-12-01",
				"08. previous close":     "169.00",
				"09. change":             "2.00",
				"10. change percent":     "1.1834%",
			},
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

	// Execute
	result, err := client.GetQuote(context.Background(), "AAPL")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 171.0, result.CurrentPrice)
	assert.Equal(t, int64(1000000), result.Volume)
	assert.Equal(t, 2.0, result.DayChange)
	assert.InDelta(t, 1.1834, result.DayChangePercent, 0.0001)
}

func TestAlphaVantageClient_GetQuote_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		serverResponse func(w http.ResponseWriter)
		expectError    bool
		errorContains  string
	}{
		{
			name:   "Missing API key",
			apiKey: "",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusUnauthorized)
			},
			expectError:   true,
			errorContains: "API key is required",
		},
		{
			name:   "Invalid API key",
			apiKey: "invalid_key",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"Note": "Invalid API key",
				})
			},
			expectError:   true,
			errorContains: "invalid API key",
		},
		{
			name:   "Invalid JSON response",
			apiKey: "test_key",
			serverResponse: func(w http.ResponseWriter) {
				w.Write([]byte("invalid json"))
			},
			expectError:   true,
			errorContains: "failed to decode response",
		},
		{
			name:   "Empty response",
			apiKey: "test_key",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{})
			},
			expectError:   true,
			errorContains: "empty response",
		},
		{
			name:   "Invalid price format",
			apiKey: "test_key",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"Global Quote": map[string]interface{}{
						"05. price": "invalid",
					},
				})
			},
			expectError:   true,
			errorContains: "invalid price format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.serverResponse(w)
			}))
			defer server.Close()

			// Create client
			mockLogger := &MockLogger{}
			client := external.NewAlphaVantageClient(server.URL, tt.apiKey, server.Client(), mockLogger)

			// Execute
			result, err := client.GetQuote(context.Background(), "AAPL")

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestAlphaVantageClient_GetCompanyOverview_Success(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Contains(t, r.URL.String(), "AAPL")
		assert.Contains(t, r.URL.String(), "apikey=test_key")
		assert.Contains(t, r.URL.String(), "function=OVERVIEW")

		// Return successful response
		response := map[string]interface{}{
			"Symbol":               "AAPL",
			"Name":                 "Apple Inc",
			"MarketCapitalization": "2800000000000",
			"PERatio":              "28.5",
			"DividendYield":        "0.5",
			"EPS":                  "6.05",
			"RevenuePerShareTTM":   "24.30",
			"ProfitMargin":         "0.25",
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

	// Execute
	result, err := client.GetCompanyOverview(context.Background(), "AAPL")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "AAPL", result.Symbol)
	assert.Equal(t, "Apple Inc", result.Name)
	assert.Equal(t, int64(2800000000000), result.MarketCap)
	assert.Equal(t, 28.5, result.PERatio)
	assert.Equal(t, 0.5, result.DividendYield)
	assert.Equal(t, 6.05, result.EPS)
	assert.Equal(t, 24.30, result.RevenuePerShare)
	assert.Equal(t, 0.25, result.ProfitMargin)
}

func TestAlphaVantageClient_GetCompanyOverview_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter)
		expectError    bool
		errorContains  string
	}{
		{
			name: "Invalid JSON response",
			serverResponse: func(w http.ResponseWriter) {
				w.Write([]byte("invalid json"))
			},
			expectError:   true,
			errorContains: "failed to decode response",
		},
		{
			name: "Empty response",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{})
			},
			expectError:   true,
			errorContains: "empty response",
		},
		{
			name: "Missing required fields",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"Symbol": "AAPL",
					// Missing other required fields
				})
			},
			expectError:   true,
			errorContains: "missing required fields",
		},
		{
			name: "Invalid numeric format",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"Symbol":               "AAPL",
					"Name":                 "Apple Inc",
					"MarketCapitalization": "invalid",
					"PERatio":              "28.5",
				})
			},
			expectError:   true,
			errorContains: "invalid numeric format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.serverResponse(w)
			}))
			defer server.Close()

			// Create client
			mockLogger := &MockLogger{}
			client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

			// Execute
			result, err := client.GetCompanyOverview(context.Background(), "AAPL")

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestAlphaVantageClient_RateLimiting(t *testing.T) {
	// Setup mock server with rate limiting
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount > 5 { // Rate limit after 5 requests
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Note": "Thank you for using Alpha Vantage! Our standard API rate limit is 5 requests per minute.",
			})
			return
		}

		// Return successful response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Global Quote": map[string]interface{}{
				"05. price": "171.00",
			},
		})
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

	// Make concurrent requests
	var wg sync.WaitGroup
	results := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.GetQuote(context.Background(), "AAPL")
			results <- err
		}()
	}

	wg.Wait()
	close(results)

	// Count successes and failures
	successCount := 0
	rateLimitCount := 0

	for err := range results {
		if err == nil {
			successCount++
		} else if strings.Contains(err.Error(), "rate limit") {
			rateLimitCount++
		}
	}

	// Assert
	assert.Equal(t, 5, successCount, "Expected 5 successful requests")
	assert.Equal(t, 5, rateLimitCount, "Expected 5 rate limited requests")
}

func TestAlphaVantageClient_ContextCancellation(t *testing.T) {
	// Setup mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Delay response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Global Quote": map[string]interface{}{
				"05. price": "171.00",
			},
		})
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Execute
	result, err := client.GetQuote(ctx, "AAPL")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
	assert.Nil(t, result)
}

func TestAlphaVantageClient_InvalidSymbol(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return empty quote for invalid symbol
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Global Quote": map[string]interface{}{},
		})
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

	// Execute with invalid symbol
	result, err := client.GetQuote(context.Background(), "INVALID")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no data found for symbol")
	assert.Nil(t, result)
}

func TestAlphaVantageClient_RetryBehavior(t *testing.T) {
	// Setup mock server that fails first request
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount == 1 {
			// First attempt fails
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Subsequent attempts succeed
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Global Quote": map[string]interface{}{
				"05. price": "171.00",
			},
		})
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

	// Execute
	result, err := client.GetQuote(context.Background(), "AAPL")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, attemptCount, "Expected one retry after initial failure")
}

func TestAlphaVantageClient_DataValidation(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse map[string]interface{}
		expectError    bool
		validate       func(t *testing.T, result *model.ExternalStockData)
	}{
		{
			name: "Valid positive price change",
			serverResponse: map[string]interface{}{
				"Global Quote": map[string]interface{}{
					"05. price":          "171.00",
					"08. previous close": "169.00",
					"09. change":         "2.00",
					"10. change percent": "1.1834%",
				},
			},
			expectError: false,
			validate: func(t *testing.T, result *model.ExternalStockData) {
				assert.Equal(t, 171.0, result.CurrentPrice)
				assert.Equal(t, 2.0, result.DayChange)
				assert.InDelta(t, 1.1834, result.DayChangePercent, 0.0001)
			},
		},
		{
			name: "Valid negative price change",
			serverResponse: map[string]interface{}{
				"Global Quote": map[string]interface{}{
					"05. price":          "169.00",
					"08. previous close": "171.00",
					"09. change":         "-2.00",
					"10. change percent": "-1.1696%",
				},
			},
			expectError: false,
			validate: func(t *testing.T, result *model.ExternalStockData) {
				assert.Equal(t, 169.0, result.CurrentPrice)
				assert.Equal(t, -2.0, result.DayChange)
				assert.InDelta(t, -1.1696, result.DayChangePercent, 0.0001)
			},
		},
		{
			name: "Zero price change",
			serverResponse: map[string]interface{}{
				"Global Quote": map[string]interface{}{
					"05. price":          "171.00",
					"08. previous close": "171.00",
					"09. change":         "0.00",
					"10. change percent": "0.0000%",
				},
			},
			expectError: false,
			validate: func(t *testing.T, result *model.ExternalStockData) {
				assert.Equal(t, 171.0, result.CurrentPrice)
				assert.Equal(t, 0.0, result.DayChange)
				assert.Equal(t, 0.0, result.DayChangePercent)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			// Create client
			mockLogger := &MockLogger{}
			client := external.NewAlphaVantageClient(server.URL, "test_key", server.Client(), mockLogger)

			// Execute
			result, err := client.GetQuote(context.Background(), "AAPL")

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				tt.validate(t, result)
			}
		})
	}
}
