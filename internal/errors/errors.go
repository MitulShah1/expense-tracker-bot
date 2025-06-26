// Package errors provides custom error types for the expense tracker bot.
package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeDatabase represents database errors
	ErrorTypeDatabase ErrorType = "DATABASE_ERROR"
	// ErrorTypeTelegram represents Telegram API errors
	ErrorTypeTelegram ErrorType = "TELEGRAM_ERROR"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
	// ErrorTypeRateLimit represents rate limiting errors
	ErrorTypeRateLimit ErrorType = "RATE_LIMIT_ERROR"
	// ErrorTypeUnauthorized represents unauthorized access errors
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
)

// AppError represents an application error
type AppError struct {
	Type      ErrorType `json:"type"`
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Details   string    `json:"details,omitempty"`
	RequestID string    `json:"requestId,omitempty"`
	UserID    int64     `json:"userId,omitempty"`
	Err       error     `json:"-"`
}

// NewValidationError creates a new validation error
func NewValidationError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    http.StatusBadRequest,
		Details: details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    http.StatusNotFound,
		Details: details,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Code:    http.StatusInternalServerError,
		Details: "Database operation failed",
		Err:     err,
	}
}

// NewTelegramError creates a new Telegram API error
func NewTelegramError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeTelegram,
		Message: message,
		Code:    http.StatusBadGateway,
		Details: "Telegram API operation failed",
		Err:     err,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    http.StatusInternalServerError,
		Details: "Internal server error",
		Err:     err,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeRateLimit,
		Message: message,
		Code:    http.StatusTooManyRequests,
		Details: "Rate limit exceeded",
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Code:    http.StatusUnauthorized,
		Details: "Unauthorized access",
	}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (details: %s, original: %v)", e.Type, e.Message, e.Details, e.Err)
	}
	return fmt.Sprintf("%s: %s (details: %s)", e.Type, e.Message, e.Details)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// IsValidationError checks if the error is a validation error
func (e *AppError) IsValidationError() bool {
	return e.Type == ErrorTypeValidation
}

// IsNotFoundError checks if the error is a not found error
func (e *AppError) IsNotFoundError() bool {
	return e.Type == ErrorTypeNotFound
}

// IsDatabaseError checks if the error is a database error
func (e *AppError) IsDatabaseError() bool {
	return e.Type == ErrorTypeDatabase
}

// IsTelegramError checks if the error is a Telegram API error
func (e *AppError) IsTelegramError() bool {
	return e.Type == ErrorTypeTelegram
}

// IsInternalError checks if the error is an internal error
func (e *AppError) IsInternalError() bool {
	return e.Type == ErrorTypeInternal
}

// IsRateLimitError checks if the error is a rate limit error
func (e *AppError) IsRateLimitError() bool {
	return e.Type == ErrorTypeRateLimit
}

// IsUnauthorizedError checks if the error is an unauthorized error
func (e *AppError) IsUnauthorizedError() bool {
	return e.Type == ErrorTypeUnauthorized
}

// WrapError wraps an existing error with additional context
func WrapError(err error, message string, errorType ErrorType) *AppError {
	appErr, ok := err.(*AppError)
	if ok {
		appErr.Message = message
		appErr.Type = errorType
		return appErr
	}

	return &AppError{
		Type:    errorType,
		Message: message,
		Code:    http.StatusInternalServerError,
		Details: "Wrapped error",
		Err:     err,
	}
}
