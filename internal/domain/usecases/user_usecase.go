package usecases

import (
	"context"
	"fmt"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/repositories"
	"stock-tracker/internal/infrastructure/auth"
	"stock-tracker/pkg/logger"
)

type UserUseCase struct {
	userRepo         repositories.UserRepository
	subscriptionRepo repositories.SubscriptionRepository
	sessionRepo      repositories.SessionRepository
	jwtService       auth.JWTService
	logger           logger.Logger
}

func NewUserUseCase(
	userRepo repositories.UserRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	sessionRepo repositories.SessionRepository,
	jwtService auth.JWTService,
	logger logger.Logger,
) *UserUseCase {
	return &UserUseCase{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		sessionRepo:      sessionRepo,
		jwtService:       jwtService,
		logger:           logger,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (uc *UserUseCase) Register(ctx context.Context, req RegisterRequest) (*entities.User, *auth.TokenPair, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Create new user
	user, err := entities.NewUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		uc.logger.Error("Failed to create user", "error", err)
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to save user", "error", err)
		return nil, nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Generate tokens
	tokens, err := uc.jwtService.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session, err := entities.NewSession(user.ID, tokens.RefreshToken, "", "")
	if err != nil {
		uc.logger.Error("Failed to create session", "error", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Warn("Failed to save session", "user_id", user.ID, "error", err)
	}

	uc.logger.Info("User registered successfully", "user_id", user.ID, "email", user.Email)
	return user, tokens, nil
}

func (uc *UserUseCase) Login(ctx context.Context, req LoginRequest) (*entities.User, *auth.TokenPair, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		uc.logger.Info("Login attempt failed - user not found", "email", req.Email)
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Validate password
	if !user.ValidatePassword(req.Password) {
		uc.logger.Info("Login attempt failed - invalid password", "user_id", user.ID)
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		uc.logger.Warn("Failed to update last login", "user_id", user.ID, "error", err)
	}

	// Generate tokens
	tokens, err := uc.jwtService.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session, err := entities.NewSession(user.ID, tokens.RefreshToken, "", "")
	if err != nil {
		uc.logger.Error("Failed to create session", "error", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Warn("Failed to save session", "user_id", user.ID, "error", err)
	}

	uc.logger.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)
	return user, tokens, nil
}

func (uc *UserUseCase) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	// Get session by refresh token
	session, err := uc.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		uc.logger.Info("Invalid refresh token attempt", "error", err)
		return nil, fmt.Errorf("invalid refresh token")
	}

	if session.IsExpired() {
		uc.logger.Info("Expired refresh token attempt", "session_id", session.ID)
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		uc.logger.Error("User not found for valid session", "user_id", session.UserID)
		return nil, fmt.Errorf("user not found")
	}

	// Generate new tokens
	tokens, err := uc.jwtService.GenerateTokenPair(user)
	if err != nil {
		uc.logger.Error("Failed to generate tokens", "error", err)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Delete old session and create new one
	if err := uc.sessionRepo.DeleteByRefreshToken(ctx, refreshToken); err != nil {
		uc.logger.Warn("Failed to delete old session", "error", err)
	}

	newSession, err := entities.NewSession(user.ID, tokens.RefreshToken, session.UserAgent, session.IPAddress)
	if err != nil {
		uc.logger.Error("Failed to create new session", "error", err)
		return nil, fmt.Errorf("failed to create new session: %w", err)
	}

	if err := uc.sessionRepo.Create(ctx, newSession); err != nil {
		uc.logger.Warn("Failed to save new session", "error", err)
	}

	return tokens, nil
}
