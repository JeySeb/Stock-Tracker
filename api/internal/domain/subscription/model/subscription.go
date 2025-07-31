package model

import (
	"time"

	"github.com/google/uuid"

	"stock-tracker/internal/domain/shared/enums"
)

// Subscription represents a user's subscription details
type Subscription struct {
	ID               uuid.UUID                `json:"id" db:"id"`
	UserID           uuid.UUID                `json:"user_id" db:"user_id" validate:"required"`
	Plan             enums.SubscriptionPlan   `json:"plan" db:"plan" validate:"required"`
	Status           enums.SubscriptionStatus `json:"status" db:"status"`
	Price            float64                  `json:"price" db:"price"`
	Currency         string                   `json:"currency" db:"currency"`
	StartDate        time.Time                `json:"start_date" db:"start_date"`
	EndDate          time.Time                `json:"end_date" db:"end_date"`
	PaymentReference string                   `json:"payment_reference" db:"payment_reference"`
	CreatedAt        time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at" db:"updated_at"`
}

// NewSubscription creates a new subscription instance
func NewSubscription(userID uuid.UUID, plan enums.SubscriptionPlan) *Subscription {
	now := time.Now()
	var endDate time.Time
	var price float64

	switch plan {
	case enums.PLAN_MONTHLY:
		endDate = now.AddDate(0, 1, 0) // 1 month
		price = float64(enums.PRICE_MONTHLY)
	case enums.PLAN_YEARLY:
		endDate = now.AddDate(1, 0, 0) // 1 year
		price = float64(enums.PRICE_YEARLY)
	default:
		// Default to monthly if invalid plan provided
		endDate = now.AddDate(0, 1, 0)
		price = float64(enums.PRICE_MONTHLY)
		plan = enums.PLAN_MONTHLY
	}

	return &Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		Plan:      plan,
		Status:    enums.STATUS_PENDING,
		Price:     price,
		Currency:  "USD",
		StartDate: now,
		EndDate:   endDate,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsActive checks if the subscription is currently active
func (s *Subscription) IsActive() bool {
	now := time.Now()
	return s.Status == enums.STATUS_ACTIVE && s.EndDate.After(now)
}

// Activate sets the subscription status to active with payment reference
func (s *Subscription) Activate(paymentReference string) {
	s.Status = enums.STATUS_ACTIVE
	s.PaymentReference = paymentReference
	s.UpdatedAt = time.Now()
}

// Cancel sets the subscription status to cancelled
func (s *Subscription) Cancel() {
	s.Status = enums.STATUS_CANCELLED
	s.UpdatedAt = time.Now()
}

// Expire sets the subscription status to expired
func (s *Subscription) Expire() {
	s.Status = enums.STATUS_EXPIRED
	s.UpdatedAt = time.Now()
}

// RenewSubscription extends the subscription period based on the current plan
func (s *Subscription) RenewSubscription() {
	now := time.Now()
	switch s.Plan {
	case enums.PLAN_MONTHLY:
		s.EndDate = now.AddDate(0, 1, 0)
	case enums.PLAN_YEARLY:
		s.EndDate = now.AddDate(1, 0, 0)
	}
	s.Status = enums.STATUS_ACTIVE
	s.UpdatedAt = now
}

// GetRemainingDays returns the number of days remaining in the subscription
func (s *Subscription) GetRemainingDays() int {
	now := time.Now()
	if now.After(s.EndDate) {
		return 0
	}
	return int(s.EndDate.Sub(now).Hours() / 24)
}
