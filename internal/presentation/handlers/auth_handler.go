package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/usecases"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// UserUseCaseInterface defines the contract for user use cases
type UserUseCaseInterface interface {
	Register(ctx context.Context, req usecases.RegisterRequest) (*entities.User, *auth.TokenPair, error)
	Login(ctx context.Context, req usecases.LoginRequest) (*entities.User, *auth.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error)
}

type AuthHandler struct {
	userUC    UserUseCaseInterface
	validator *validator.Validate
	logger    logger.Logger
}

func NewAuthHandler(userUC UserUseCaseInterface, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		userUC:    userUC,
		validator: validator.New(),
		logger:    logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req usecases.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode register request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		h.logger.Info("Invalid register request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Validation failed: " + err.Error()})
		return
	}

	user, tokens, err := h.userUC.Register(r.Context(), req)
	if err != nil {
		h.logger.Error("Registration failed", "error", err, "email", req.Email)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req usecases.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode login request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		h.logger.Info("Invalid login request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Validation failed: " + err.Error()})
		return
	}

	user, tokens, err := h.userUC.Login(r.Context(), req)
	if err != nil {
		h.logger.Info("Login failed", "error", err, "email", req.Email)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid credentials"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode refresh token request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request body"})
		return
	}
	defer r.Body.Close()

	if err := h.validator.Struct(req); err != nil {
		h.logger.Info("Invalid refresh token request", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Validation failed: " + err.Error()})
		return
	}

	tokens, err := h.userUC.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Info("Token refresh failed", "error", err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid or expired refresh token"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"tokens": tokens,
	})
}
