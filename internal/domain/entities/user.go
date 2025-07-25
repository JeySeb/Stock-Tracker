package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserTier string
type SubscriptionStatus string

const (
	// User Tiers
	TIER_GUEST   UserTier = "guest"   // Non-registered users
	TIER_BASIC   UserTier = "basic"   // Registered users
	TIER_PREMIUM UserTier = "premium" // Premium subscribers

	// Subscription Status
	SUB_STATUS_ACTIVE    SubscriptionStatus = "active"
	SUB_STATUS_CANCELLED SubscriptionStatus = "cancelled"
	SUB_STATUS_EXPIRED   SubscriptionStatus = "expired"
	SUB_STATUS_PENDING   SubscriptionStatus = "pending"

	// API Rate Limits per hour
	guestRateLimit   = 100
	basicRateLimit   = 500
	premiumRateLimit = 2000
	defaultRateLimit = 50

	// Password requirements
	minPasswordLength = 8
	maxPasswordLength = 128

	// Password cost for bcrypt
	bcryptCost = bcrypt.DefaultCost

	// Account lockout
	maxLoginAttempts = 5
)

// User represents a user entity in the system
type User struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Email      string     `json:"email" db:"email" validate:"required,email"`
	Password   string     `json:"-" db:"password_hash" validate:"required,min=8"`
	FirstName  string     `json:"first_name" db:"first_name" validate:"required,min=1,max=100"`
	LastName   string     `json:"last_name" db:"last_name" validate:"required,min=1,max=100"`
	Tier       UserTier   `json:"tier" db:"tier"`
	IsVerified bool       `json:"is_verified" db:"is_verified"`
	LastLogin  *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user instance with basic tier access
func NewUser(email, password, firstName, lastName string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:         uuid.New(),
		Email:      email,
		Password:   string(hashedPassword),
		FirstName:  firstName,
		LastName:   lastName,
		Tier:       TIER_BASIC,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// ValidatePassword checks if the provided password matches the stored hash
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// SetLastLogin updates the user's last login timestamp
func (u *User) SetLastLogin(lastLogin time.Time) {
	u.LastLogin = &lastLogin
	u.SetUpdatedAt(lastLogin)
}

// SetUpdatedAt updates the last modified timestamp
func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.UpdatedAt = updatedAt
}

// CanAccessExternalAPIs checks if the user has permission to access external APIs
func (u *User) CanAccessExternalAPIs() bool {
	return u.IsVerified && (u.Tier == TIER_BASIC || u.Tier == TIER_PREMIUM)
}

// CanAccessAIFeatures checks if the user has permission to access AI features
func (u *User) CanAccessAIFeatures() bool {
	return u.IsVerified && u.Tier == TIER_PREMIUM
}

// GetAPIRateLimit returns the hourly API rate limit based on user tier
func (u *User) GetAPIRateLimit() int {
	if !u.IsVerified {
		return defaultRateLimit
	}

	switch u.Tier {
	case TIER_GUEST:
		return guestRateLimit
	case TIER_BASIC:
		return basicRateLimit
	case TIER_PREMIUM:
		return premiumRateLimit
	default:
		return defaultRateLimit
	}
}

// ValidatePasswordStrength validates password complexity
func ValidatePasswordStrength(password string) error {
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	}
	if len(password) > maxPasswordLength {
		return fmt.Errorf("password must be no more than %d characters long", maxPasswordLength)
	}

	// Check for at least one uppercase, one lowercase, one digit, and one special character
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case char >= 33 && char <= 126:
			if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
				hasSpecial = true
			}
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// IsAccountLocked checks if the account should be locked due to failed attempts
func (u *User) IsAccountLocked() bool {
	// TODO: Implement account lockout logic with failed login attempts tracking
	// This would require additional fields or a separate entity for tracking login attempts
	return false
}

// SanitizeForJSON returns a user struct safe for JSON serialization (removes sensitive data)
func (u *User) SanitizeForJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":          u.ID,
		"email":       u.Email,
		"first_name":  u.FirstName,
		"last_name":   u.LastName,
		"tier":        u.Tier,
		"is_verified": u.IsVerified,
		"last_login":  u.LastLogin,
		"created_at":  u.CreatedAt,
		"updated_at":  u.UpdatedAt,
	}
}
