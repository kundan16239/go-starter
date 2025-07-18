package logger

import (
	"fmt"
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

// NewLogger creates a new SimpleLogger instance
func NewLogger() Logger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// format logs messages with level prefix
func (l *SimpleLogger) logWithLevel(level string, msg string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", level)
	if len(args) > 0 {
		l.logger.Printf(prefix+msg, args...)
	} else {
		l.logger.Print(prefix + msg)
	}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.logWithLevel("INFO", msg, args...)
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.logWithLevel("ERROR", msg, args...)
}

// Fatal logs a fatal message and exits
func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	l.logWithLevel("FATAL", msg, args...)
	os.Exit(1)
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	l.logWithLevel("DEBUG", msg, args...)
}

// package logger

// import (
// 	"go.uber.org/zap"
// )

// type Logger interface {
// 	Info(msg string, args ...interface{})
// 	Error(msg string, args ...interface{})
// 	Fatal(msg string, args ...interface{})
// 	Debug(msg string, args ...interface{})
// }

// type ZapLogger struct {
// 	sugar *zap.SugaredLogger
// }

// func NewLogger() Logger {
// 	// Use zap.NewProduction() for production environment
// 	logger, err := zap.NewDevelopment()
// 	if err != nil {
// 		panic("failed to create zap logger: " + err.Error())
// 	}
// 	sugar := logger.Sugar()
// 	return &ZapLogger{sugar: sugar}
// }

// func (l *ZapLogger) Info(msg string, args ...interface{}) {
// 	l.sugar.Infof(msg, args...)
// }

// func (l *ZapLogger) Error(msg string, args ...interface{}) {
// 	l.sugar.Errorf(msg, args...)
// }

// func (l *ZapLogger) Fatal(msg string, args ...interface{}) {
// 	l.sugar.Fatalf(msg, args...)
// }

// func (l *ZapLogger) Debug(msg string, args ...interface{}) {
// 	l.sugar.Debugf(msg, args...)
// }
