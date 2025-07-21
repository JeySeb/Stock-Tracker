package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

type SimpleLogger struct {
	logger *log.Logger
}

func NewSimpleLogger() Logger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.logger.Printf("[INFO] "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.logger.Printf("[WARN] "+msg, args...)
}

// New creates a logger with the specified level (for compatibility)
func New(level string) Logger {
	return NewSimpleLogger()
}
