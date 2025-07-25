package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/tests/mocks"
)

func TestUserUseCase_Register_Success(t *testing.T) {
	// Arrange
	userRepo := &mocks.MockUserRepository{}
	subscriptionRepo := &mocks.MockSubscriptionRepository{}
	sessionRepo := &mocks.MockSessionRepository{}
	jwtService := &mocks.MockJWTService{}
	logger := &mocks.MockLogger{}

	useCase := usecases.NewUserUseCase(userRepo, subscriptionRepo, sessionRepo, jwtService, logger)

	req := usecases.RegisterRequest{
		Email:     "test@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations
	userRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, assert.AnError)
	userRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	tokenPair := &auth.TokenPair{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
	}
	jwtService.On("GenerateTokenPair", mock.Anything).Return(tokenPair, nil)
	sessionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe()

	// Act
	user, tokens, err := useCase.Register(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, tokens)
}
