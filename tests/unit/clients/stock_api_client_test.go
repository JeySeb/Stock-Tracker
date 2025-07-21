package clients_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"stock-tracker/internal/infrastructure/clients"
	"stock-tracker/tests/mocks"
)

func TestStockAPIClient_FetchPage_Success(t *testing.T) {
	// Arrange
	mockResponse := clients.StockAPIResponse{
		Items: []clients.StockAPIItem{
			{
				Ticker:     "AAPL",
				TargetFrom: "$150.00",
				TargetTo:   "$180.00",
				Company:    "Apple Inc.",
				Action:     "upgraded by",
				Brokerage:  "Goldman Sachs",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				Time:       "2024-01-15T10:30:00Z",
			},
			{
				Ticker:     "GOOGL",
				TargetFrom: "$2800.00",
				TargetTo:   "$3000.00",
				Company:    "Alphabet Inc.",
				Action:     "initiated by",
				Brokerage:  "Morgan Stanley",
				RatingFrom: "Neutral",
				RatingTo:   "Outperform",
				Time:       "2024-01-15T11:00:00Z",
			},
		},
		NextPage: "page2",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Stock-Tracker/1.0", r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	logger.On("Warn", mock.AnythingOfType("string"), mock.Anything).Maybe()

	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	stocks, nextPage, err := client.FetchPage(context.Background(), "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "page2", nextPage)
	assert.Len(t, stocks, 2)

	// Verify first stock
	stock1 := stocks[0]
	assert.Equal(t, "AAPL", stock1.Ticker)
	assert.Equal(t, "Apple Inc.", stock1.Company)
	assert.Equal(t, "Goldman Sachs", stock1.Brokerage)
	assert.Equal(t, "upgraded by", stock1.Action)
	assert.Equal(t, "Hold", stock1.RatingFrom)
	assert.Equal(t, "Buy", stock1.RatingTo)
	assert.Equal(t, 150.0, stock1.TargetFrom)
	assert.Equal(t, 180.0, stock1.TargetTo)

	// Verify timestamp parsing
	expectedTime, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
	assert.Equal(t, expectedTime, stock1.EventTime)
}

func TestStockAPIClient_FetchPage_WithNextPage(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify next_page parameter
		nextPage := r.URL.Query().Get("next_page")
		assert.Equal(t, "test-page", nextPage)

		response := clients.StockAPIResponse{
			Items:    []clients.StockAPIItem{},
			NextPage: "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	_, nextPage, err := client.FetchPage(context.Background(), "test-page")

	// Assert
	require.NoError(t, err)
	assert.Empty(t, nextPage)
}

func TestStockAPIClient_FetchPage_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	stocks, nextPage, err := client.FetchPage(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "failed to make request")
	assert.Empty(t, nextPage)
	assert.Nil(t, stocks)
}

func TestStockAPIClient_FetchPage_InvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	stocks, nextPage, err := client.FetchPage(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response")
	assert.Empty(t, nextPage)
	assert.Nil(t, stocks)
}

func TestStockAPIClient_FetchAllStocks_Success(t *testing.T) {
	// Arrange
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		var response clients.StockAPIResponse

		switch callCount {
		case 1:
			// First page
			response = clients.StockAPIResponse{
				Items: []clients.StockAPIItem{
					{
						Ticker:    "AAPL",
						Company:   "Apple Inc.",
						Brokerage: "Goldman Sachs",
						Action:    "upgraded by",
						Time:      "2024-01-15T10:30:00Z",
					},
				},
				NextPage: "page2",
			}
		case 2:
			// Second page
			response = clients.StockAPIResponse{
				Items: []clients.StockAPIItem{
					{
						Ticker:    "GOOGL",
						Company:   "Alphabet Inc.",
						Brokerage: "Morgan Stanley",
						Action:    "initiated by",
						Time:      "2024-01-15T11:00:00Z",
					},
				},
				NextPage: "", // Last page
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	logger.On("Info", mock.AnythingOfType("string"), mock.Anything).Maybe()

	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	stocks, err := client.FetchAllStocks(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, stocks, 2)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Equal(t, "GOOGL", stocks[1].Ticker)
	assert.Equal(t, 2, callCount) // Verify pagination worked
}

func TestStockAPIClient_FetchAllStocks_ContextCancellation(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	logger.On("Info", mock.AnythingOfType("string"), mock.Anything).Maybe()
	logger.On("Error", mock.AnythingOfType("string"), mock.Anything).Maybe()

	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Act
	stocks, err := client.FetchAllStocks(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "context deadline exceeded")
	assert.Nil(t, stocks)
}

func TestStockAPIClient_ConvertAPIItemToStock_InvalidTime(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := clients.StockAPIResponse{
			Items: []clients.StockAPIItem{
				{
					Ticker:    "AAPL",
					Company:   "Apple Inc.",
					Brokerage: "Goldman Sachs",
					Action:    "upgraded by",
					Time:      "invalid-time-format",
				},
			},
			NextPage: "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	logger.On("Warn", mock.AnythingOfType("string"), mock.Anything).Return()

	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	_, nextPage, err := client.FetchPage(context.Background(), "")

	// Assert
	require.NoError(t, err) // FetchPage should not fail
	assert.Empty(t, nextPage)

	// Verify warning was logged for invalid time format
	logger.AssertCalled(t, "Warn", mock.AnythingOfType("string"), mock.Anything)
}

func TestStockAPIClient_ParsePrice(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "Simple price",
			input:    "150.50",
			expected: 150.50,
		},
		{
			name:     "Price with dollar sign",
			input:    "$150.50",
			expected: 150.50,
		},
		{
			name:     "Price with commas",
			input:    "$1,250.75",
			expected: 1250.75,
		},
		{
			name:     "Price with spaces",
			input:    " $150.50 ",
			expected: 150.50,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "Invalid format",
			input:    "not-a-number",
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We need to test this indirectly through the API response
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := clients.StockAPIResponse{
					Items: []clients.StockAPIItem{
						{
							Ticker:     "TEST",
							TargetFrom: tc.input,
							TargetTo:   tc.input,
							Company:    "Test Company",
							Action:     "test action",
							Brokerage:  "Test Brokerage",
							Time:       "2024-01-15T10:30:00Z",
						},
					},
					NextPage: "",
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			logger := &mocks.MockLogger{}
			logger.On("Warn", mock.AnythingOfType("string"), mock.Anything).Maybe()

			client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

			// Act
			stocks, _, err := client.FetchPage(context.Background(), "")

			// Assert
			require.NoError(t, err)
			if tc.expected == 0.0 && tc.input != "" && tc.input != "0" {
				// For invalid formats, we might get empty stocks or zero values
				if len(stocks) > 0 {
					assert.Equal(t, tc.expected, stocks[0].TargetFrom)
					assert.Equal(t, tc.expected, stocks[0].TargetTo)
				}
			} else if len(stocks) > 0 {
				assert.Equal(t, tc.expected, stocks[0].TargetFrom)
				assert.Equal(t, tc.expected, stocks[0].TargetTo)
			}
		})
	}
}

func TestStockAPIClient_RateLimit(t *testing.T) {
	// Arrange
	requestTimes := []time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTimes = append(requestTimes, time.Now())

		// Return different responses to trigger pagination
		callCount := len(requestTimes)
		var response clients.StockAPIResponse

		if callCount == 1 {
			response = clients.StockAPIResponse{
				Items: []clients.StockAPIItem{
					{
						Ticker:    "AAPL",
						Company:   "Apple Inc.",
						Brokerage: "Goldman Sachs",
						Action:    "upgraded by",
						Time:      "2024-01-15T10:30:00Z",
					},
				},
				NextPage: "page2",
			}
		} else {
			response = clients.StockAPIResponse{
				Items:    []clients.StockAPIItem{},
				NextPage: "",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	logger.On("Info", mock.AnythingOfType("string"), mock.Anything).Maybe()

	client := clients.NewStockAPIClient(server.URL, "test-api-key", logger)

	// Act
	start := time.Now()
	_, err := client.FetchAllStocks(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, requestTimes, 2) // Should have made 2 requests

	// Verify rate limiting - there should be a delay between requests
	if len(requestTimes) >= 2 {
		timeDiff := requestTimes[1].Sub(requestTimes[0])
		assert.True(t, timeDiff >= 100*time.Millisecond, "Expected rate limiting delay, got %v", timeDiff)
	}

	// Total time should be reasonable (not too long due to excessive delays)
	totalTime := time.Since(start)
	assert.True(t, totalTime < 5*time.Second, "Total time too long: %v", totalTime)
}

// Benchmark tests
func BenchmarkStockAPIClient_ConvertAPIItem(b *testing.B) {
	// Create a test server for benchmarking
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := clients.StockAPIResponse{
			Items: []clients.StockAPIItem{
				{
					Ticker:     "AAPL",
					TargetFrom: "$150.50",
					TargetTo:   "$180.75",
					Company:    "Apple Inc.",
					Action:     "upgraded by",
					Brokerage:  "Goldman Sachs",
					RatingFrom: "Hold",
					RatingTo:   "Buy",
					Time:       "2024-01-15T10:30:00Z",
				},
			},
			NextPage: "",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	logger := &mocks.MockLogger{}
	logger.On("Warn", mock.AnythingOfType("string"), mock.Anything).Maybe()

	client := clients.NewStockAPIClient(server.URL, "test-key", logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.FetchPage(context.Background(), "")
	}
}
