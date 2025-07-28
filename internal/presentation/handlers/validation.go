package handlers

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"stock-tracker/internal/domain/shared/enums"
)

// ValidationConfig holds configuration for validation rules
type ValidationConfig struct {
	// Ticker validation
	TickerMinLength int
	TickerMaxLength int
	TickerPattern   *regexp.Regexp

	// Parameter limits
	MaxLimit      int
	DefaultLimit  int
	MaxOffset     int
	MaxFloatValue float64
	MinFloatValue float64

	// String limits
	MaxStringLength int
	MaxSliceLength  int
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		TickerMinLength: 1,
		TickerMaxLength: 10,
		TickerPattern:   regexp.MustCompile(`^[A-Z0-9.-]+$`),
		MaxLimit:        1000,
		DefaultLimit:    10,
		MaxOffset:       10000,
		MaxFloatValue:   100.0,
		MinFloatValue:   -100.0,
		MaxStringLength: 100,
		MaxSliceLength:  50,
	}
}

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e ValidationErrors) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("validation failed with %d errors", len(e.Errors))
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// AddError adds a validation error
func (e *ValidationErrors) AddError(field, message string, value ...string) {
	err := ValidationError{
		Field:   field,
		Message: message,
	}
	if len(value) > 0 {
		err.Value = value[0]
	}
	e.Errors = append(e.Errors, err)
}

// RequestValidator provides validation for HTTP request parameters
type RequestValidator struct {
	config *ValidationConfig
}

// NewRequestValidator creates a new request validator with optional configuration
func NewRequestValidator(config ...*ValidationConfig) *RequestValidator {
	validationConfig := DefaultValidationConfig()
	if len(config) > 0 && config[0] != nil {
		validationConfig = config[0]
	}

	return &RequestValidator{
		config: validationConfig,
	}
}

// ValidateTicker validates a stock ticker symbol
func (v *RequestValidator) ValidateTicker(ticker string) error {
	if ticker == "" {
		return ValidationError{
			Field:   "ticker",
			Message: "ticker is required",
		}
	}

	// Normalize ticker
	normalizedTicker := strings.ToUpper(strings.TrimSpace(ticker))

	// Check length
	if len(normalizedTicker) < v.config.TickerMinLength {
		return ValidationError{
			Field:   "ticker",
			Message: fmt.Sprintf("ticker must be at least %d characters", v.config.TickerMinLength),
			Value:   ticker,
		}
	}

	if len(normalizedTicker) > v.config.TickerMaxLength {
		return ValidationError{
			Field:   "ticker",
			Message: fmt.Sprintf("ticker must be at most %d characters", v.config.TickerMaxLength),
			Value:   ticker,
		}
	}

	// Check pattern
	if !v.config.TickerPattern.MatchString(normalizedTicker) {
		return ValidationError{
			Field:   "ticker",
			Message: "ticker contains invalid characters (only letters, numbers, dots, and hyphens allowed)",
			Value:   ticker,
		}
	}

	// Additional validation: check for common patterns
	if strings.HasPrefix(normalizedTicker, ".") || strings.HasSuffix(normalizedTicker, ".") {
		return ValidationError{
			Field:   "ticker",
			Message: "ticker cannot start or end with a dot",
			Value:   ticker,
		}
	}

	return nil
}

// ValidateQueryParams validates common query parameters
func (v *RequestValidator) ValidateQueryParams(params url.Values) (*ValidationErrors, *ValidatedParams) {
	errors := &ValidationErrors{}
	validated := &ValidatedParams{}

	// Validate limit
	if limitStr := params.Get("limit"); limitStr != "" {
		if limit, err := v.validateIntParam("limit", limitStr, 1, v.config.MaxLimit); err != nil {
			errors.AddError("limit", err.Error(), limitStr)
		} else {
			validated.Limit = limit
		}
	} else {
		validated.Limit = v.config.DefaultLimit
	}

	// Validate offset
	if offsetStr := params.Get("offset"); offsetStr != "" {
		if offset, err := v.validateIntParam("offset", offsetStr, 0, v.config.MaxOffset); err != nil {
			errors.AddError("offset", err.Error(), offsetStr)
		} else {
			validated.Offset = offset
		}
	}

	// Validate min_score
	if scoreStr := params.Get("min_score"); scoreStr != "" {
		if score, err := v.validateFloatParam("min_score", scoreStr, v.config.MinFloatValue, v.config.MaxFloatValue); err != nil {
			errors.AddError("min_score", err.Error(), scoreStr)
		} else {
			validated.MinScore = &score
		}
	}

	// Validate recommendation type
	if typeStr := params.Get("type"); typeStr != "" {
		if recType, err := v.validateRecommendationType("type", typeStr); err != nil {
			errors.AddError("type", err.Error(), typeStr)
		} else {
			validated.RecommendationType = recType
		}
	}

	// Validate exclude tickers
	if excludeStr := params.Get("exclude"); excludeStr != "" {
		if tickers, err := v.validateTickerList("exclude", excludeStr); err != nil {
			errors.AddError("exclude", err.Error(), excludeStr)
		} else {
			validated.ExcludeTickers = tickers
		}
	}

	// Validate sort_by and sort_order
	if sortBy := params.Get("sort_by"); sortBy != "" {
		if err := v.validateSortBy("sort_by", sortBy); err != nil {
			errors.AddError("sort_by", err.Error(), sortBy)
		} else {
			validated.SortBy = sortBy
		}
	}

	if sortOrder := params.Get("sort_order"); sortOrder != "" {
		if err := v.validateSortOrder("sort_order", sortOrder); err != nil {
			errors.AddError("sort_order", err.Error(), sortOrder)
		} else {
			validated.SortOrder = sortOrder
		}
	}

	return errors, validated
}

// ValidatedParams holds validated and parsed parameters
type ValidatedParams struct {
	Limit              int                       `json:"limit"`
	Offset             int                       `json:"offset"`
	MinScore           *float64                  `json:"min_score,omitempty"`
	RecommendationType *enums.RecommendationType `json:"recommendation_type,omitempty"`
	ExcludeTickers     []string                  `json:"exclude_tickers,omitempty"`
	SortBy             string                    `json:"sort_by,omitempty"`
	SortOrder          string                    `json:"sort_order,omitempty"`
}

// validateIntParam validates an integer parameter within bounds
func (v *RequestValidator) validateIntParam(field, value string, min, max int) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("must be a valid integer")
	}

	if parsed < min {
		return 0, fmt.Errorf("must be at least %d", min)
	}

	if parsed > max {
		return 0, fmt.Errorf("must be at most %d", max)
	}

	return parsed, nil
}

// validateFloatParam validates a float parameter within bounds
func (v *RequestValidator) validateFloatParam(field, value string, min, max float64) (float64, error) {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}

	if parsed < min {
		return 0, fmt.Errorf("must be at least %.2f", min)
	}

	if parsed > max {
		return 0, fmt.Errorf("must be at most %.2f", max)
	}

	return parsed, nil
}

// validateRecommendationType validates recommendation type parameter
func (v *RequestValidator) validateRecommendationType(field, value string) (*enums.RecommendationType, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))

	switch normalized {
	case "strong_buy", "strongbuy", "strong-buy":
		recType := enums.RECOMMENDATION_TYPE_STRONG_BUY
		return &recType, nil
	case "buy":
		recType := enums.RECOMMENDATION_TYPE_BUY
		return &recType, nil
	case "hold":
		recType := enums.RECOMMENDATION_TYPE_HOLD
		return &recType, nil
	case "sell":
		recType := enums.RECOMMENDATION_TYPE_SELL
		return &recType, nil
	case "strong_sell", "strongsell", "strong-sell":
		recType := enums.RECOMMENDATION_TYPE_STRONG_SELL
		return &recType, nil
	default:
		return nil, fmt.Errorf("must be one of: strong_buy, buy, hold, sell, strong_sell")
	}
}

// validateTickerList validates a comma-separated list of tickers
func (v *RequestValidator) validateTickerList(field, value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	tickers := strings.Split(value, ",")
	if len(tickers) > v.config.MaxSliceLength {
		return nil, fmt.Errorf("too many tickers (maximum %d allowed)", v.config.MaxSliceLength)
	}

	var validated []string
	for i, ticker := range tickers {
		ticker = strings.TrimSpace(ticker)
		if ticker == "" {
			continue
		}

		if err := v.ValidateTicker(ticker); err != nil {
			return nil, fmt.Errorf("ticker at position %d: %s", i+1, err.Error())
		}

		validated = append(validated, strings.ToUpper(ticker))
	}

	return validated, nil
}

// validateSortBy validates sort_by parameter
func (v *RequestValidator) validateSortBy(field, value string) error {
	normalized := strings.ToLower(strings.TrimSpace(value))
	validSortFields := []string{
		"ticker", "company", "score", "confidence", "event_time",
		"created_at", "target_price", "recommendation_type",
	}

	for _, validField := range validSortFields {
		if normalized == validField {
			return nil
		}
	}

	return fmt.Errorf("must be one of: %s", strings.Join(validSortFields, ", "))
}

// validateSortOrder validates sort_order parameter
func (v *RequestValidator) validateSortOrder(field, value string) error {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized != "asc" && normalized != "desc" {
		return fmt.Errorf("must be 'asc' or 'desc'")
	}
	return nil
}

// ValidateStringLength validates string length
func (v *RequestValidator) ValidateStringLength(field, value string, maxLength int) error {
	if len(value) > maxLength {
		return ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at most %d characters", maxLength),
			Value:   value,
		}
	}
	return nil
}

// ValidateNoSpecialChars validates that a string contains no special characters
func (v *RequestValidator) ValidateNoSpecialChars(field, value string) error {
	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != ' ' && char != '-' && char != '_' {
			return ValidationError{
				Field:   field,
				Message: "contains invalid characters",
				Value:   value,
			}
		}
	}
	return nil
}

// SanitizeString sanitizes a string by removing potentially harmful characters
func (v *RequestValidator) SanitizeString(input string) string {
	// Remove any null bytes and control characters
	sanitized := strings.ReplaceAll(input, "\x00", "")

	// Remove excessive whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	sanitized = spaceRegex.ReplaceAllString(sanitized, " ")

	return sanitized
}
