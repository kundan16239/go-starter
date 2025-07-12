package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// New creates a new application error
func New(code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewWithDetails creates a new application error with details
func NewWithDetails(code int, message, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewBadRequest creates a bad request error
func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:       1001,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewUnauthorized creates an unauthorized error
func NewUnauthorized(message string) *AppError {
	return &AppError{
		Code:       1002,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbidden creates a forbidden error
func NewForbidden(message string) *AppError {
	return &AppError{
		Code:       1003,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewNotFound creates a not found error
func NewNotFound(message string) *AppError {
	return &AppError{
		Code:       1004,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewConflict creates a conflict error
func NewConflict(message string) *AppError {
	return &AppError{
		Code:       1005,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       1006,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string) *AppError {
	return &AppError{
		Code:       2001,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewExternalServiceError creates an external service error
func NewExternalServiceError(message string) *AppError {
	return &AppError{
		Code:       3001,
		Message:    message,
		StatusCode: http.StatusServiceUnavailable,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError returns the AppError if the error is an AppError, otherwise nil
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
