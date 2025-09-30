package logging

import (
	"go.uber.org/zap"
)

// Logger is a simple logging interface
type Log struct{}

func NewLogger() *Log {
	return &Log{}
}

// NewLogger creates a new Zap logger.
func ZapLogger() (*zap.Logger, error) {
	// For production, use zap.NewProduction()
	return zap.NewDevelopment()
}
