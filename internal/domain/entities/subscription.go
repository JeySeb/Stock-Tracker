package entities

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionPlan string

const (
	// Subscription Plans
	PlanMonthly SubscriptionPlan = "monthly"
	PlanYearly  SubscriptionPlan = "yearly"

	// Plan Prices in USD
	monthlyPrice = 29.99
	yearlyPrice  = 299.99
)

// Subscription represents a user's subscription details
type Subscription struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	UserID           uuid.UUID          `json:"user_id" db:"user_id" validate:"required"`
	Plan             SubscriptionPlan   `json:"plan" db:"plan" validate:"required"`
	Status           SubscriptionStatus `json:"status" db:"status"`
	Price            float64            `json:"price" db:"price"`
	Currency         string             `json:"currency" db:"currency"`
	StartDate        time.Time          `json:"start_date" db:"start_date"`
	EndDate          time.Time          `json:"end_date" db:"end_date"`
	PaymentReference string             `json:"payment_reference" db:"payment_reference"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" db:"updated_at"`
}

// NewSubscription creates a new subscription instance
func NewSubscription(userID uuid.UUID, plan SubscriptionPlan) *Subscription {
	now := time.Now()
	var endDate time.Time
	var price float64

	switch plan {
	case PlanMonthly:
		endDate = now.AddDate(0, 1, 0) // 1 month
		price = monthlyPrice
	case PlanYearly:
		endDate = now.AddDate(1, 0, 0) // 1 year
		price = yearlyPrice
	default:
		// Default to monthly if invalid plan provided
		endDate = now.AddDate(0, 1, 0)
		price = monthlyPrice
		plan = PlanMonthly
	}

	return &Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		Plan:      plan,
		Status:    SUB_STATUS_PENDING,
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
	return s.Status == SUB_STATUS_ACTIVE && s.EndDate.After(now)
}

// Activate sets the subscription status to active with payment reference
func (s *Subscription) Activate(paymentReference string) {
	s.Status = SUB_STATUS_ACTIVE
	s.PaymentReference = paymentReference
	s.UpdatedAt = time.Now()
}

// Cancel sets the subscription status to cancelled
func (s *Subscription) Cancel() {
	s.Status = SUB_STATUS_CANCELLED
	s.UpdatedAt = time.Now()
}

// Expire sets the subscription status to expired
func (s *Subscription) Expire() {
	s.Status = SUB_STATUS_EXPIRED
	s.UpdatedAt = time.Now()
}

// RenewSubscription extends the subscription period based on the current plan
func (s *Subscription) RenewSubscription() {
	now := time.Now()
	switch s.Plan {
	case PlanMonthly:
		s.EndDate = now.AddDate(0, 1, 0)
	case PlanYearly:
		s.EndDate = now.AddDate(1, 0, 0)
	}
	s.Status = SUB_STATUS_ACTIVE
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
