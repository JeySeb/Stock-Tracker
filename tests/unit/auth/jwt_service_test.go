package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/infrastructure/auth"
)

func TestJWTService_GenerateTokenPair_Success(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Tier:      entities.TIER_BASIC,
	}

	// Act
	tokens, err := jwtService.GenerateTokenPair(user)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Greater(t, tokens.ExpiresIn, int64(0))
}

func TestJWTService_GenerateTokenPair_NilUser(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	// Act
	tokens, err := jwtService.GenerateTokenPair(nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "user is nil")
}

func TestJWTService_ValidateAccessToken_Success(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Tier:      entities.TIER_PREMIUM,
	}

	tokens, err := jwtService.GenerateTokenPair(user)
	require.NoError(t, err)

	// Act
	claims, err := jwtService.ValidateAccessToken(tokens.AccessToken)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Tier, claims.Tier)
	assert.Equal(t, "stock-tracker", claims.Issuer)
}

func TestJWTService_ValidateAccessToken_EmptyToken(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	// Act
	claims, err := jwtService.ValidateAccessToken("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is empty")
}

func TestJWTService_ValidateAccessToken_InvalidToken(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	// Act
	claims, err := jwtService.ValidateAccessToken("invalid.jwt.token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_ValidateAccessToken_WrongSigningKey(t *testing.T) {
	// Arrange
	jwtService1 := auth.NewJWTService("first-secret-key-minimum-32-chars")
	jwtService2 := auth.NewJWTService("second-secret-key-minimum-32-chars")

	user := &entities.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Tier:  entities.TIER_BASIC,
	}

	tokens, err := jwtService1.GenerateTokenPair(user)
	require.NoError(t, err)

	// Act - try to validate with different secret
	claims, err := jwtService2.ValidateAccessToken(tokens.AccessToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTService_GenerateRefreshToken_Success(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	// Act
	refreshToken1, err1 := jwtService.GenerateRefreshToken()
	refreshToken2, err2 := jwtService.GenerateRefreshToken()

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, refreshToken1)
	assert.NotEmpty(t, refreshToken2)
	assert.NotEqual(t, refreshToken1, refreshToken2) // Should be unique
}

func TestJWTService_TokenPair_Integration(t *testing.T) {
	// Arrange
	jwtService := auth.NewJWTService("test-secret-key-minimum-32-characters")

	user := &entities.User{
		ID:        uuid.New(),
		Email:     "integration@example.com",
		FirstName: "Integration",
		LastName:  "Test",
		Tier:      entities.TIER_PREMIUM,
	}

	// Act - Generate tokens
	tokens, err := jwtService.GenerateTokenPair(user)
	require.NoError(t, err)

	// Act - Validate access token
	claims, err := jwtService.ValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)

	// Assert - Complete flow
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Tier, claims.Tier)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.True(t, time.Until(claims.ExpiresAt.Time) > 0, "Token should not be expired")
}
