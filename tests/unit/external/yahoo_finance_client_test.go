package external

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"stock-tracker/internal/infrastructure/external"
)

func TestYahooFinanceClient_Constructor(t *testing.T) {
	// Create client
	mockLogger := &MockLogger{}
	client := external.NewYahooFinanceClient(mockLogger)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestYahooFinanceClient_Constructor_NilLogger(t *testing.T) {
	// Create client with nil logger
	client := external.NewYahooFinanceClient(nil)

	// Assert client is created
	assert.NotNil(t, client)
}

// MockLogger implementation
type MockLogger struct{}

func (l *MockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *MockLogger) Error(msg string, keysAndValues ...interface{}) {}
