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
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/infrastructure/external"
)

// MockHTTPClient mock implementation
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestYahooFinanceClient_GetQuote_Success(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Contains(t, r.URL.String(), "AAPL")
		assert.Contains(t, r.Header.Get("User-Agent"), "Mozilla")

		// Return successful response
		response := map[string]interface{}{
			"chart": map[string]interface{}{
				"result": []map[string]interface{}{
					{
						"meta": map[string]interface{}{
							"currency":                 "USD",
							"symbol":                   "AAPL",
							"regularMarketPrice":       170.0,
							"previousClose":            168.0,
							"regularMarketVolume":      1000000,
							"marketCap":                2800000000000,
							"trailingPE":               28.5,
							"dividendYield":            0.5,
							"fiftyTwoWeekHigh":         190.0,
							"fiftyTwoWeekLow":          140.0,
							"averageDailyVolume3Month": 800000,
						},
					},
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server URL
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

	// Execute
	result, err := client.GetQuote(context.Background(), "AAPL")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 170.0, result.CurrentPrice)
	assert.Equal(t, int64(1000000), result.Volume)
	assert.Equal(t, int64(2800000000000), result.MarketCap)
	assert.NotNil(t, result.PERatio)
	assert.Equal(t, 28.5, *result.PERatio)
	assert.NotNil(t, result.Week52High)
	assert.Equal(t, 190.0, *result.Week52High)
	assert.NotNil(t, result.Week52Low)
	assert.Equal(t, 140.0, *result.Week52Low)
	assert.NotNil(t, result.AvgVolume)
	assert.Equal(t, int64(800000), *result.AvgVolume)
}

func TestYahooFinanceClient_GetQuote_ErrorHandling(t *testing.T) {
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
			name: "API error response",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"chart": map[string]interface{}{
						"error": map[string]interface{}{
							"code":        "NotFound",
							"description": "No data found",
						},
					},
				})
			},
			expectError:   true,
			errorContains: "Yahoo Finance API error",
		},
		{
			name: "Empty result array",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"chart": map[string]interface{}{
						"result": []interface{}{},
					},
				})
			},
			expectError:   true,
			errorContains: "no data found",
		},
		{
			name: "Missing required fields",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"chart": map[string]interface{}{
						"result": []map[string]interface{}{
							{
								"meta": map[string]interface{}{
									// Missing regularMarketPrice
									"symbol": "AAPL",
								},
							},
						},
					},
				})
			},
			expectError:   true,
			errorContains: "missing required fields",
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
			client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

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

func TestYahooFinanceClient_GetQuote_HTTPErrors(t *testing.T) {
	tests := []struct {
		name          string
		serverHandler func(w http.ResponseWriter, r *http.Request)
		expectError   bool
	}{
		{
			name: "HTTP 429 Too Many Requests",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTooManyRequests)
			},
			expectError: true,
		},
		{
			name: "HTTP 500 Internal Server Error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectError: true,
		},
		{
			name: "HTTP 404 Not Found",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(tt.serverHandler))
			defer server.Close()

			// Create client
			mockLogger := &MockLogger{}
			client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

			// Execute
			result, err := client.GetQuote(context.Background(), "AAPL")

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestYahooFinanceClient_GetHistoricalData_Success(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Contains(t, r.URL.String(), "AAPL")
		assert.Contains(t, r.URL.String(), "1mo") // period parameter

		// Return successful response
		response := map[string]interface{}{
			"chart": map[string]interface{}{
				"result": []map[string]interface{}{
					{
						"timestamp": []int64{1625097600, 1625184000},
						"indicators": map[string]interface{}{
							"quote": []map[string]interface{}{
								{
									"close":  []float64{145.5, 146.8},
									"volume": []int64{80000000, 85000000},
									"high":   []float64{146.0, 147.5},
									"low":    []float64{144.5, 146.0},
									"open":   []float64{144.8, 146.2},
								},
							},
						},
					},
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

	// Execute
	result, err := client.GetHistoricalData(context.Background(), "AAPL", "1mo")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// Verify first data point
	assert.Equal(t, time.Unix(1625097600, 0), result[0].Date)
	assert.Equal(t, 145.5, result[0].Close)
	assert.Equal(t, int64(80000000), result[0].Volume)
	assert.Equal(t, 146.0, result[0].High)
	assert.Equal(t, 144.5, result[0].Low)
	assert.Equal(t, 144.8, result[0].Open)

	// Verify second data point
	assert.Equal(t, time.Unix(1625184000, 0), result[1].Date)
	assert.Equal(t, 146.8, result[1].Close)
}

func TestYahooFinanceClient_GetHistoricalData_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		period         string
		serverResponse func(w http.ResponseWriter)
		expectError    bool
		errorContains  string
	}{
		{
			name:   "Invalid period",
			period: "invalid",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
			},
			expectError:   true,
			errorContains: "invalid period",
		},
		{
			name:   "Missing data",
			period: "1mo",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"chart": map[string]interface{}{
						"result": []map[string]interface{}{},
					},
				})
			},
			expectError:   true,
			errorContains: "no historical data",
		},
		{
			name:   "Invalid data format",
			period: "1mo",
			serverResponse: func(w http.ResponseWriter) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"chart": map[string]interface{}{
						"result": []map[string]interface{}{
							{
								"timestamp": []int64{1625097600},
								"indicators": map[string]interface{}{
									"quote": []map[string]interface{}{
										{
											"close": []string{"invalid"}, // Wrong type
										},
									},
								},
							},
						},
					},
				})
			},
			expectError:   true,
			errorContains: "invalid data format",
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
			client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

			// Execute
			result, err := client.GetHistoricalData(context.Background(), "AAPL", tt.period)

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

func TestYahooFinanceClient_RateLimiting(t *testing.T) {
	// Setup mock server with rate limiting
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount > 5 { // Rate limit after 5 requests
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		// Return successful response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"chart": map[string]interface{}{
				"result": []map[string]interface{}{
					{
						"meta": map[string]interface{}{
							"regularMarketPrice": 170.0,
							"symbol":             "AAPL",
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

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
		} else if strings.Contains(err.Error(), "429") {
			rateLimitCount++
		}
	}

	// Assert
	assert.Equal(t, 5, successCount, "Expected 5 successful requests")
	assert.Equal(t, 5, rateLimitCount, "Expected 5 rate limited requests")
}

func TestYahooFinanceClient_ContextCancellation(t *testing.T) {
	// Setup mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Delay response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"chart": map[string]interface{}{
				"result": []map[string]interface{}{
					{
						"meta": map[string]interface{}{
							"regularMarketPrice": 170.0,
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	// Create client
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(server.URL, server.Client(), mockLogger)

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

// MockLogger implementation
type MockLogger struct{}

func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {}
