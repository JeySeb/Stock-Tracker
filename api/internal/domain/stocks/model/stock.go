package model

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

// GetRatingChangeScore calculates the improvement/degradation of rating
func (s *Stock) GetRatingChangeScore() float64 {
	fromScore, toScore := s.GetRatingScore()
	return toScore - fromScore
}

// IsRecommendation determines if this is a positive recommendation
// A recommendation is positive if the final rating suggests buying (Strong Buy, Buy, Outperform)
func (s *Stock) IsRecommendation() bool {
	_, toScore := s.GetRatingScore()

	// Consider positive if final rating is above neutral (0.5)
	// This means Strong Buy (1.0), Buy (0.8), Outperform (0.75) are positive
	// Hold (0.5), Neutral (0.4), Underperform (0.25), Sell (0.2), Strong Sell (0.0) are negative
	return toScore > 0.5
}

// IsPositiveChange determines if this represents a positive rating change
func (s *Stock) IsPositiveChange() bool {
	return s.GetRatingChangeScore() > 0
}

// IsNegativeChange determines if this represents a negative rating change
func (s *Stock) IsNegativeChange() bool {
	return s.GetRatingChangeScore() < 0
}

// GetRecommendationStrength returns the strength of the recommendation (0-1)
// regardless of direction, for confidence-based calculations
func (s *Stock) GetRecommendationStrength() float64 {
	_, toScore := s.GetRatingScore()

	// Calculate distance from neutral (0.5)
	// Strong opinions (very high or very low scores) have high strength
	neutralPoint := 0.5
	distance := abs(toScore - neutralPoint)

	// Normalize to 0-1 scale (max distance from neutral is 0.5)
	return distance / 0.5
}

// GetDirectionalCertainty returns directional certainty [-1, +1]
// Positive values = confident buy recommendations
// Negative values = confident sell recommendations
// Zero = neutral/uncertain
func (s *Stock) GetDirectionalCertainty() float64 {
	_, toScore := s.GetRatingScore()
	strength := s.GetRecommendationStrength()

	// Apply direction to certainty
	if toScore > 0.5 {
		return strength // Positive certainty for buy recommendations
	} else if toScore < 0.5 {
		return -strength // Negative certainty for sell recommendations
	} else {
		return 0.0 // Neutral
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
