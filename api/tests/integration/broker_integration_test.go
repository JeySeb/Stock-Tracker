package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stock-tracker/internal/domain/shared/enums"
	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/tests/mocks"

	"github.com/google/uuid"
)

func TestBrokerIntegration_GetBrokersWithScores(t *testing.T) {
	// This is a basic integration test that verifies the broker endpoint works
	// In a real scenario, you would set up a test database with sample data

	// Arrange
	mockLogger := &mocks.MockLogger{}

	// Create a mock broker usecase that returns sample data
	mockBrokerUC := &mockBrokerUseCase{
		brokers: []*stockRepos.BrokerWithScore{
			{
				ID:               uuid.New(),
				Name:             "Goldman Sachs",
				CredibilityScore: 0.85,
				ReportCount:      150,
				CalculatedScore:  0.92,
				CreatedAt:        "2024-01-01T00:00:00Z",
				UpdatedAt:        "2024-01-15T00:00:00Z",
			},
			{
				ID:               uuid.New(),
				Name:             "Morgan Stanley",
				CredibilityScore: 0.78,
				ReportCount:      120,
				CalculatedScore:  0.88,
				CreatedAt:        "2024-01-01T00:00:00Z",
				UpdatedAt:        "2024-01-15T00:00:00Z",
			},
		},
	}

	handler := handlers.NewBrokerHandler(mockBrokerUC, mockLogger)

	// Test cases
	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Get all brokers",
			url:            "/brokers/scores",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get brokers with limit",
			url:            "/brokers/scores?limit=1",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Get brokers with ascending order",
			url:            "/brokers/scores?order=asc",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get brokers with descending order",
			url:            "/brokers/scores?order=desc",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get brokers with invalid limit",
			url:            "/brokers/scores?limit=invalid",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Get brokers with invalid order",
			url:            "/brokers/scores?order=invalid",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", tc.url, nil)

			// Add user tier to context
			ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_GUEST)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			// Act
			handler.GetBrokersWithScores(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)

				assert.NotNil(t, response["data"])
				assert.NotNil(t, response["meta"])

				data := response["data"].([]interface{})
				assert.Len(t, data, tc.expectedCount)

				meta := response["meta"].(map[string]interface{})
				assert.Equal(t, "TIER_GUEST", meta["user_tier"])
				assert.Equal(t, float64(tc.expectedCount), meta["count"])
			}
		})
	}
}

func TestBrokerIntegration_WithUserTiers(t *testing.T) {
	// Test broker endpoint with different user tiers
	mockLogger := &mocks.MockLogger{}
	mockBrokerUC := &mockBrokerUseCase{
		brokers: []*stockRepos.BrokerWithScore{
			{
				ID:               uuid.New(),
				Name:             "Goldman Sachs",
				CredibilityScore: 0.85,
				ReportCount:      150,
				CalculatedScore:  0.92,
				CreatedAt:        "2024-01-01T00:00:00Z",
				UpdatedAt:        "2024-01-15T00:00:00Z",
			},
		},
	}

	handler := handlers.NewBrokerHandler(mockBrokerUC, mockLogger)

	userTiers := []enums.UserTier{
		enums.TIER_GUEST,
		enums.TIER_BASIC,
		enums.TIER_PREMIUM,
	}

	for _, tier := range userTiers {
		t.Run("UserTier_"+string(tier), func(t *testing.T) {
			req := httptest.NewRequest("GET", "/brokers/scores", nil)
			ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, tier)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.GetBrokersWithScores(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			meta := response["meta"].(map[string]interface{})
			assert.Equal(t, string(tier), meta["user_tier"])
		})
	}
}

// Mock implementations for testing

type mockBrokerUseCase struct {
	brokers []*stockRepos.BrokerWithScore
}

func (m *mockBrokerUseCase) GetBrokersWithScores(ctx context.Context, limit *int, orderBy string) ([]*stockRepos.BrokerWithScore, error) {
	// Simple mock implementation
	if limit != nil && *limit < len(m.brokers) {
		return m.brokers[:*limit], nil
	}
	return m.brokers, nil
}
