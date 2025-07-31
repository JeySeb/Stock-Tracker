package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockLogger implements logger.Logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	arguments := []interface{}{msg}
	arguments = append(arguments, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	arguments := []interface{}{msg}
	arguments = append(arguments, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	arguments := []interface{}{msg}
	arguments = append(arguments, args...)
	m.Called(arguments...)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	arguments := []interface{}{msg}
	arguments = append(arguments, args...)
	m.Called(arguments...)
}
