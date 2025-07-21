package entities

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Stock struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Ticker     string    `json:"ticker" db:"ticker" validate:"required,min=1,max=10"`
	Company    string    `json:"company" db:"company" validate:"required,min=1,max=255"`
	BrokerID   uuid.UUID `json:"broker_id" db:"broker_id"`
	Brokerage  string    `json:"brokerage" db:"brokerage"`
	Action     string    `json:"action" db:"action" validate:"required"`
	RatingFrom string    `json:"rating_from" db:"rating_from"`
	RatingTo   string    `json:"rating_to" db:"rating_to"`
	TargetFrom float64   `json:"target_from" db:"target_from"`
	TargetTo   float64   `json:"target_to" db:"target_to"`
	EventTime  time.Time `json:"event_time" db:"event_time"`
	PriceClose *float64  `json:"price_close,omitempty" db:"price_close"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewStock(ticker, company, brokerage, action string, eventTime time.Time) *Stock {
	return &Stock{
		ID:        uuid.New(),
		Ticker:    ticker,
		Company:   company,
		Brokerage: brokerage,
		Action:    action,
		EventTime: eventTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Stock) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *Stock) IsUpgrade() bool {
	upgradeActions := []string{"upgraded by", "raised to", "initiated by"}

	for _, action := range upgradeActions {
		if strings.Contains(strings.ToLower(s.Action), action) {
			return true
		}
	}

	return false

}

func (s *Stock) GetPriceTargetChange() float64 {
	if s.TargetFrom <= 0 || s.TargetTo <= 0 {
		return 0
	}
	return (s.TargetTo - s.TargetFrom) / s.TargetFrom
}

func (s *Stock) GetRatingScore() (fromScore, toScore float64) {
	ratingScores := map[string]float64{
		"strong buy":   1.0,
		"buy":          0.8,
		"outperform":   0.75,
		"hold":         0.5,
		"neutral":      0.4,
		"underperform": 0.25,
		"sell":         0.2,
		"strong sell":  0.0,
	}

	fromScore = ratingScores[strings.ToLower(s.RatingFrom)]
	toScore = ratingScores[strings.ToLower(s.RatingTo)]

	return
}

func (s *Stock) GetPriceChange() float64 {
	if s.PriceClose == nil {
		return 0
	}
	return *s.PriceClose
}

// GetRatingChangeScore calculates the improvement/degradation of rating
func (s *Stock) GetRatingChangeScore() float64 {
	fromScore, toScore := s.GetRatingScore()
	return toScore - fromScore
}

// IsRecommendation determines if this is a positive recommendation
func (s *Stock) IsRecommendation() bool {
	// Check for positive actions
	positiveActions := []string{"upgraded", "initiated", "reiterated"}
	action := strings.ToLower(s.Action)

	for _, positive := range positiveActions {
		if strings.Contains(action, positive) {
			return true
		}
	}

	// Check for rating improvement
	return s.GetRatingChangeScore() > 0
}
