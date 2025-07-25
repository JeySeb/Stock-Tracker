package usecases

import (
	"context"
	"fmt"
	"time"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/internal/domain/repositories"
	"stock-tracker/pkg/logger"

	"github.com/google/uuid"
)

type SubscriptionUseCase struct {
	subscriptionRepo repositories.SubscriptionRepository
	userRepo         repositories.UserRepository
	logger           logger.Logger
}

func NewSubscriptionUseCase(
	subscriptionRepo repositories.SubscriptionRepository,
	userRepo repositories.UserRepository,
	logger logger.Logger,
) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		logger:           logger,
	}
}

type PaymentSimulationRequest struct {
	Plan entities.SubscriptionPlan `json:"plan" validate:"required,oneof=monthly yearly"`
}

func (uc *SubscriptionUseCase) CreateSubscription(ctx context.Context, userID uuid.UUID, req PaymentSimulationRequest) (*entities.Subscription, error) {
	// Validate user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user already has active subscription
	existing, err := uc.subscriptionRepo.GetActiveByUserID(ctx, userID)
	if err == nil && existing != nil && existing.IsActive() {
		uc.logger.Info("User already has active subscription", "user_id", userID)
		return nil, fmt.Errorf("user already has an active subscription")
	}

	// Create subscription
	subscription := entities.NewSubscription(userID, req.Plan)
	if err := uc.subscriptionRepo.Create(ctx, subscription); err != nil {
		uc.logger.Error("Failed to create subscription", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	uc.logger.Info("Subscription created",
		"user_id", userID,
		"subscription_id", subscription.ID,
		"plan", subscription.Plan)
	return subscription, nil
}

func (uc *SubscriptionUseCase) SimulatePayment(ctx context.Context, subscriptionID uuid.UUID) error {
	// Get subscription
	subscription, err := uc.subscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		uc.logger.Error("Subscription not found", "error", err, "subscription_id", subscriptionID)
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Validate subscription status
	if subscription.Status != entities.SUB_STATUS_PENDING {
		uc.logger.Info("Invalid subscription status for payment",
			"subscription_id", subscriptionID,
			"status", subscription.Status)
		return fmt.Errorf("invalid subscription status for payment")
	}

	// Simulate payment processing (in real app, this would integrate with Stripe, etc.)
	time.Sleep(2 * time.Second) // Simulate payment processing time

	// Activate subscription
	paymentRef := fmt.Sprintf("sim_payment_%s_%d", subscriptionID.String()[:8], time.Now().Unix())
	subscription.Activate(paymentRef)

	if err := uc.subscriptionRepo.Update(ctx, subscription); err != nil {
		uc.logger.Error("Failed to activate subscription", "error", err, "subscription_id", subscriptionID)
		return fmt.Errorf("failed to activate subscription: %w", err)
	}

	// Update user tier to premium
	user, err := uc.userRepo.GetByID(ctx, subscription.UserID)
	if err != nil {
		uc.logger.Error("Failed to get user", "error", err, "user_id", subscription.UserID)
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.Tier = entities.TIER_PREMIUM
	user.SetUpdatedAt(time.Now())

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user tier", "error", err, "user_id", user.ID)
		return fmt.Errorf("failed to update user tier: %w", err)
	}

	uc.logger.Info("Payment simulated and subscription activated",
		"subscription_id", subscriptionID,
		"user_id", subscription.UserID,
		"payment_ref", paymentRef)

	return nil
}
