package external

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"stock-tracker/internal/infrastructure/external"
)

func TestAlphaVantageClient_Constructor(t *testing.T) {
	// Create client
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient("test_key", mockLogger)

	// Assert client is created
	assert.NotNil(t, client)
}

func TestAlphaVantageClient_Constructor_EmptyAPIKey(t *testing.T) {
	// Create client with empty API key
	mockLogger := &MockLogger{}
	client := external.NewAlphaVantageClient("", mockLogger)

	// Assert client is still created (validation might be done at request time)
	assert.NotNil(t, client)
}

func TestAlphaVantageClient_Constructor_NilLogger(t *testing.T) {
	// Create client with nil logger
	client := external.NewAlphaVantageClient("test_key", nil)

	// Assert client is created
	assert.NotNil(t, client)
}
