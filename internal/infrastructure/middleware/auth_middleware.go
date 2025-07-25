package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/render"
)

type contextKey string

const (
	UserContextKey     contextKey = "user"
	UserIDContextKey   contextKey = "user_id"
	UserTierContextKey contextKey = "user_tier"
)

type AuthMiddleware struct {
	jwtService auth.JWTService
	logger     logger.Logger
}

func NewAuthMiddleware(jwtService auth.JWTService, logger logger.Logger) *AuthMiddleware {
	if jwtService == nil {
		panic("jwtService cannot be nil")
	}
	if logger == nil {
		panic("logger cannot be nil")
	}
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// RequireAuth middleware - requires valid JWT token
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.extractAndValidateToken(r)
		if err != nil {
			m.logger.Warn("Authentication failed", "error", err, "path", r.URL.Path)
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Authentication required"})
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UserTierContextKey, claims.Tier)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePremium middleware - requires premium subscription
func (m *AuthMiddleware) RequirePremium(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.extractAndValidateToken(r)
		if err != nil {
			m.logger.Warn("Authentication failed", "error", err, "path", r.URL.Path)
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Authentication required"})
			return
		}

		if claims.Tier != entities.UserTier("premium") {
			m.logger.Info("Non-premium user attempted to access premium feature",
				"user_id", claims.UserID,
				"tier", claims.Tier,
				"path", r.URL.Path)
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "Premium subscription required"})
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, UserTierContextKey, claims.Tier)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware - adds user info if token is present
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.extractAndValidateToken(r)
		if err == nil {
			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
			ctx = context.WithValue(ctx, UserTierContextKey, claims.Tier)
			r = r.WithContext(ctx)
		} else {
			// Set guest tier for non-authenticated users
			ctx := context.WithValue(r.Context(), UserTierContextKey, entities.UserTier("guest"))
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// extractAndValidateToken validates the Authorization header and returns the token claims
func (m *AuthMiddleware) extractAndValidateToken(r *http.Request) (*auth.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("missing authorization header")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], "bearer") {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	claims, err := m.jwtService.ValidateAccessToken(tokenParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
