package logger

import (
	"log"
	"os"
)

// Logger defines the interface for logging operations
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// SimpleLogger implements the Logger interface with basic logging
type SimpleLogger struct {
	logger *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() Logger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Printf("[INFO] "+msg, args...)
	} else {
		l.logger.Printf("[INFO] " + msg)
	}
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Printf("[ERROR] "+msg, args...)
	} else {
		l.logger.Printf("[ERROR] " + msg)
	}
}

// Fatal logs a fatal message and exits
func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Printf("[FATAL] "+msg, args...)
	} else {
		l.logger.Printf("[FATAL] " + msg)
	}
	os.Exit(1)
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Printf("[DEBUG] "+msg, args...)
	} else {
		l.logger.Printf("[DEBUG] " + msg)
	}
}
