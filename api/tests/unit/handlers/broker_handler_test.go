package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"stock-tracker/internal/domain/shared/enums"
	stockRepos "stock-tracker/internal/domain/stocks/repositories"
	"stock-tracker/internal/infrastructure/middleware"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/tests/mocks"

	"github.com/google/uuid"
)

type mockBrokerUseCase struct {
	mock.Mock
}

func (m *mockBrokerUseCase) GetBrokersWithScores(ctx context.Context, limit *int, orderBy string) ([]*stockRepos.BrokerWithScore, error) {
	args := m.Called(ctx, limit, orderBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*stockRepos.BrokerWithScore), args.Error(1)
}

func TestBrokerHandler_GetBrokersWithScores_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	testBrokers := []*stockRepos.BrokerWithScore{
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
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "desc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_GUEST,
		"count", 2,
		"limit", (*int)(nil),
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["meta"])

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	meta := response["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["count"])
	assert.Equal(t, "guest", meta["user_tier"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_WithLimit(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	limit := 5
	testBrokers := []*stockRepos.BrokerWithScore{
		{
			ID:               uuid.New(),
			Name:             "Goldman Sachs",
			CredibilityScore: 0.85,
			ReportCount:      150,
			CalculatedScore:  0.92,
			CreatedAt:        "2024-01-01T00:00:00Z",
			UpdatedAt:        "2024-01-15T00:00:00Z",
		},
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), &limit, "desc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_GUEST,
		"count", 1,
		"limit", &limit,
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores?limit=5", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])
	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_WithOrderAsc(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	testBrokers := []*stockRepos.BrokerWithScore{
		{
			ID:               uuid.New(),
			Name:             "Small Broker",
			CredibilityScore: 0.45,
			ReportCount:      10,
			CalculatedScore:  0.52,
			CreatedAt:        "2024-01-01T00:00:00Z",
			UpdatedAt:        "2024-01-15T00:00:00Z",
		},
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "asc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_GUEST,
		"count", 1,
		"limit", (*int)(nil),
		"order", "asc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores?order=asc", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_WithOrderDesc(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	testBrokers := []*stockRepos.BrokerWithScore{
		{
			ID:               uuid.New(),
			Name:             "Goldman Sachs",
			CredibilityScore: 0.95,
			ReportCount:      200,
			CalculatedScore:  0.98,
			CreatedAt:        "2024-01-01T00:00:00Z",
			UpdatedAt:        "2024-01-15T00:00:00Z",
		},
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "desc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_GUEST,
		"count", 1,
		"limit", (*int)(nil),
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores?order=desc", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_WithInvalidLimit(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	testBrokers := []*stockRepos.BrokerWithScore{
		{
			ID:               uuid.New(),
			Name:             "Goldman Sachs",
			CredibilityScore: 0.85,
			ReportCount:      150,
			CalculatedScore:  0.92,
			CreatedAt:        "2024-01-01T00:00:00Z",
			UpdatedAt:        "2024-01-15T00:00:00Z",
		},
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "desc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_GUEST,
		"count", 1,
		"limit", (*int)(nil),
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores?limit=invalid", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_WithInvalidOrder(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	testBrokers := []*stockRepos.BrokerWithScore{
		{
			ID:               uuid.New(),
			Name:             "Goldman Sachs",
			CredibilityScore: 0.85,
			ReportCount:      150,
			CalculatedScore:  0.92,
			CreatedAt:        "2024-01-01T00:00:00Z",
			UpdatedAt:        "2024-01-15T00:00:00Z",
		},
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "desc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_GUEST,
		"count", 1,
		"limit", (*int)(nil),
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores?order=invalid", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])

	mockUseCase.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_UseCaseError(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "desc").
		Return(nil, errors.New("database error"))

	// Set up mock logger expectations for error case
	mockLogger.On("Error", "Failed to get brokers with scores",
		"error", errors.New("database error"),
		"user_tier", enums.TIER_GUEST,
		"limit", (*int)(nil),
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Failed to retrieve brokers with scores", response["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestBrokerHandler_GetBrokersWithScores_WithUserTier(t *testing.T) {
	// Arrange
	mockUseCase := &mockBrokerUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewBrokerHandler(mockUseCase, mockLogger)

	testBrokers := []*stockRepos.BrokerWithScore{
		{
			ID:               uuid.New(),
			Name:             "Goldman Sachs",
			CredibilityScore: 0.85,
			ReportCount:      150,
			CalculatedScore:  0.92,
			CreatedAt:        "2024-01-01T00:00:00Z",
			UpdatedAt:        "2024-01-15T00:00:00Z",
		},
	}

	mockUseCase.On("GetBrokersWithScores", mock.MatchedBy(func(ctx context.Context) bool { return true }), (*int)(nil), "desc").
		Return(testBrokers, nil)

	// Set up mock logger expectations
	mockLogger.On("Info", "Brokers with scores served",
		"user_tier", enums.TIER_BASIC,
		"count", 1,
		"limit", (*int)(nil),
		"order", "desc").Return()

	req := httptest.NewRequest("GET", "/brokers/scores", nil)
	// Add user tier to context
	ctx := context.WithValue(req.Context(), middleware.UserTierContextKey, enums.TIER_BASIC)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Act
	handler.GetBrokersWithScores(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["meta"])

	meta := response["meta"].(map[string]interface{})
	assert.Equal(t, "basic", meta["user_tier"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
