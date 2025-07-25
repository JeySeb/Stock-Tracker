package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"stock-tracker/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenGeneration   = errors.New("failed to generate token")
	ErrUnexpectedSigning = errors.New("unexpected signing method")
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Claims struct {
	UserID uuid.UUID         `json:"user_id"`
	Email  string            `json:"email"`
	Tier   entities.UserTier `json:"tier"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateTokenPair(user *entities.User) (*TokenPair, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	GenerateRefreshToken() (string, error)
}

type jwtService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  15 * time.Minute,   // 15 minutes
		refreshTokenTTL: 7 * 24 * time.Hour, // 7 days
		issuer:          "stock-tracker",
	}
}

func (s *jwtService) GenerateTokenPair(user *entities.User) (*TokenPair, error) {
	if user == nil {
		return nil, fmt.Errorf("%w: user is nil", ErrTokenGeneration)
	}

	// Generate access token
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Tier:   user.Tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	// Generate refresh token
	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}, nil
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("%w: token is empty", ErrInvalidToken)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigning, token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Additional validation
	if time.Until(claims.ExpiresAt.Time) < 0 {
		return nil, fmt.Errorf("%w: token has expired", ErrInvalidToken)
	}

	if claims.Issuer != s.issuer {
		return nil, fmt.Errorf("%w: invalid issuer", ErrInvalidToken)
	}

	return claims, nil
}

func (s *jwtService) GenerateRefreshToken() (string, error) {
	const tokenLength = 32
	bytes := make([]byte, tokenLength)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
