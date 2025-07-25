package middleware

import (
	"net/http"
	"sync"
	"time"

	"stock-tracker/internal/domain/entities"
	"stock-tracker/pkg/logger"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
	logger   logger.Logger
}

func NewRateLimiter(logger logger.Logger) *RateLimiter {
	if logger == nil {
		panic("logger cannot be nil")
	}

	rl := &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		logger:   logger,
	}

	// Clean up old visitors every hour
	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user tier from context (set by auth middleware)
		tier, ok := r.Context().Value(UserTierContextKey).(entities.UserTier)
		if !ok {
			tier = entities.TIER_GUEST
		}

		// Create identifier (IP + tier for logged users, just IP for guests)
		identifier := r.RemoteAddr
		if tier != entities.TIER_GUEST {
			if userID, ok := r.Context().Value(UserIDContextKey).(uuid.UUID); ok {
				identifier = userID.String()
			}
		}

		limiter := rl.getLimiter(identifier, tier)

		if !limiter.Allow() {
			rl.logger.Warn("Rate limit exceeded",
				"identifier", identifier,
				"tier", tier,
				"path", r.URL.Path)
			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, map[string]string{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getLimiter(identifier string, tier entities.UserTier) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[identifier]
	if !exists {
		// Set rate limits based on user tier
		var limit rate.Limit
		switch tier {
		case entities.TIER_GUEST:
			limit = rate.Every(36 * time.Second) // ~100 requests per hour
		case entities.TIER_BASIC:
			limit = rate.Every(7200 * time.Millisecond) // ~500 requests per hour
		case entities.TIER_PREMIUM:
			limit = rate.Every(1800 * time.Millisecond) // ~2000 requests per hour
		default:
			limit = rate.Every(72 * time.Second) // ~50 requests per hour
		}

		limiter = rate.NewLimiter(limit, 10) // Burst of 10
		rl.visitors[identifier] = limiter
	}

	return limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Hour)
		rl.mu.Lock()
		for id, limiter := range rl.visitors {
			if limiter.Tokens() == float64(limiter.Burst()) {
				delete(rl.visitors, id)
			}
		}
		rl.mu.Unlock()
	}
}
