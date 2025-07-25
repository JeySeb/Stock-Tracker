package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/internal/presentation/handlers"
	"stock-tracker/tests/mocks"
)

type mockUserUseCase struct {
	mock.Mock
}

func (m *mockUserUseCase) Register(ctx context.Context, req usecases.RegisterRequest) (*entities.User, *auth.TokenPair, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entities.User), args.Get(1).(*auth.TokenPair), args.Error(2)
}

func (m *mockUserUseCase) Login(ctx context.Context, req usecases.LoginRequest) (*entities.User, *auth.TokenPair, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entities.User), args.Get(1).(*auth.TokenPair), args.Error(2)
}

func (m *mockUserUseCase) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenPair), args.Error(1)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	userID := uuid.New()
	testUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Tier:      entities.TIER_BASIC,
	}

	testTokens := &auth.TokenPair{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    900,
	}

	registerReq := usecases.RegisterRequest{
		Email:     "test@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	mockUseCase.On("Register", mock.Anything, registerReq).Return(testUser, testTokens, nil)
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	requestBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Register(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "tokens")

	mockUseCase.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	mockLogger.On("Error", "Failed to decode register request", "error", mock.Anything)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Register(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])

	mockLogger.AssertExpectations(t)
}

func TestAuthHandler_Register_ValidationFailure(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	mockLogger.On("Info", "Invalid register request", "error", mock.Anything)

	registerReq := usecases.RegisterRequest{
		Email:     "invalid-email",
		Password:  "123", // too short
		FirstName: "",    // empty
		LastName:  "User",
	}

	requestBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Register(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Validation failed")

	mockLogger.AssertExpectations(t)
}

func TestAuthHandler_Register_UseCaseError(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	registerReq := usecases.RegisterRequest{
		Email:     "test@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	mockUseCase.On("Register", mock.Anything, registerReq).Return(nil, nil, errors.New("user already exists"))
	mockLogger.On("Error", "Registration failed", "error", mock.Anything, "email", registerReq.Email)

	requestBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Register(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user already exists", response["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	userID := uuid.New()
	testUser := &entities.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Tier:      entities.TIER_BASIC,
	}

	testTokens := &auth.TokenPair{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    900,
	}

	loginReq := usecases.LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	mockUseCase.On("Login", mock.Anything, loginReq).Return(testUser, testTokens, nil)

	requestBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "tokens")

	mockUseCase.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	loginReq := usecases.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	mockUseCase.On("Login", mock.Anything, loginReq).Return(nil, nil, errors.New("invalid credentials"))
	mockLogger.On("Info", "Login failed", "error", mock.Anything, "email", loginReq.Email)

	requestBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.Login(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid credentials", response["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	testTokens := &auth.TokenPair{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresIn:    900,
	}

	refreshReq := map[string]string{
		"refresh_token": "valid-refresh-token",
	}

	mockUseCase.On("RefreshToken", mock.Anything, "valid-refresh-token").Return(testTokens, nil)

	requestBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.RefreshToken(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "tokens")

	mockUseCase.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_ExpiredToken(t *testing.T) {
	// Arrange
	mockUseCase := &mockUserUseCase{}
	mockLogger := &mocks.MockLogger{}
	handler := handlers.NewAuthHandler(mockUseCase, mockLogger)

	refreshReq := map[string]string{
		"refresh_token": "expired-refresh-token",
	}

	mockUseCase.On("RefreshToken", mock.Anything, "expired-refresh-token").Return(nil, errors.New("refresh token expired"))
	mockLogger.On("Info", "Token refresh failed", "error", mock.Anything)

	requestBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.RefreshToken(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid or expired refresh token", response["error"])

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
