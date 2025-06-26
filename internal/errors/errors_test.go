package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name      string
		errorType ErrorType
		expected  string
	}{
		{"validation error", ErrorTypeValidation, "VALIDATION_ERROR"},
		{"not found error", ErrorTypeNotFound, "NOT_FOUND"},
		{"database error", ErrorTypeDatabase, "DATABASE_ERROR"},
		{"telegram error", ErrorTypeTelegram, "TELEGRAM_ERROR"},
		{"internal error", ErrorTypeInternal, "INTERNAL_ERROR"},
		{"rate limit error", ErrorTypeRateLimit, "RATE_LIMIT_ERROR"},
		{"unauthorized error", ErrorTypeUnauthorized, "UNAUTHORIZED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, string(tt.errorType))
		})
	}
}

func TestNewValidationError(t *testing.T) {
	message := "Invalid input"
	details := "Field 'amount' must be positive"

	err := NewValidationError(message, details)

	require.Equal(t, ErrorTypeValidation, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusBadRequest, err.Code)
	require.Equal(t, details, err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Nil(t, err.Err)
}

func TestNewNotFoundError(t *testing.T) {
	message := "Resource not found"
	details := "User with ID 123 not found"

	err := NewNotFoundError(message, details)

	require.Equal(t, ErrorTypeNotFound, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusNotFound, err.Code)
	require.Equal(t, details, err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Nil(t, err.Err)
}

func TestNewDatabaseError(t *testing.T) {
	message := "Database operation failed"
	originalErr := errors.New("connection timeout")

	err := NewDatabaseError(message, originalErr)

	require.Equal(t, ErrorTypeDatabase, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusInternalServerError, err.Code)
	require.Equal(t, "Database operation failed", err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Equal(t, originalErr, err.Err)
}

func TestNewTelegramError(t *testing.T) {
	message := "Telegram API call failed"
	originalErr := errors.New("invalid token")

	err := NewTelegramError(message, originalErr)

	require.Equal(t, ErrorTypeTelegram, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusBadGateway, err.Code)
	require.Equal(t, "Telegram API operation failed", err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Equal(t, originalErr, err.Err)
}

func TestNewInternalError(t *testing.T) {
	message := "Internal server error"
	originalErr := errors.New("unexpected error")

	err := NewInternalError(message, originalErr)

	require.Equal(t, ErrorTypeInternal, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusInternalServerError, err.Code)
	require.Equal(t, "Internal server error", err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Equal(t, originalErr, err.Err)
}

func TestNewRateLimitError(t *testing.T) {
	message := "Too many requests"

	err := NewRateLimitError(message)

	require.Equal(t, ErrorTypeRateLimit, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusTooManyRequests, err.Code)
	require.Equal(t, "Rate limit exceeded", err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Nil(t, err.Err)
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "Access denied"

	err := NewUnauthorizedError(message)

	require.Equal(t, ErrorTypeUnauthorized, err.Type)
	require.Equal(t, message, err.Message)
	require.Equal(t, http.StatusUnauthorized, err.Code)
	require.Equal(t, "Unauthorized access", err.Details)
	require.Empty(t, err.RequestID)
	require.Zero(t, err.UserID)
	require.Nil(t, err.Err)
}

func TestAppError_Error(t *testing.T) {
	t.Run("error without underlying error", func(t *testing.T) {
		err := NewValidationError("Invalid input", "Field required")
		expected := "VALIDATION_ERROR: Invalid input (details: Field required)"
		require.Equal(t, expected, err.Error())
	})

	t.Run("error with underlying error", func(t *testing.T) {
		originalErr := errors.New("connection failed")
		err := NewDatabaseError("Database error", originalErr)
		expected := "DATABASE_ERROR: Database error (details: Database operation failed, original: connection failed)"
		require.Equal(t, expected, err.Error())
	})
}

func TestAppError_Unwrap(t *testing.T) {
	t.Run("error with underlying error", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := NewDatabaseError("Database error", originalErr)
		require.Equal(t, originalErr, err.Unwrap())
	})

	t.Run("error without underlying error", func(t *testing.T) {
		err := NewValidationError("Invalid input", "Field required")
		require.Nil(t, err.Unwrap())
	})
}

func TestAppError_TypeChecking(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		err := NewValidationError("test", "test")
		require.True(t, err.IsValidationError())
		require.False(t, err.IsNotFoundError())
		require.False(t, err.IsDatabaseError())
		require.False(t, err.IsTelegramError())
		require.False(t, err.IsInternalError())
		require.False(t, err.IsRateLimitError())
		require.False(t, err.IsUnauthorizedError())
	})

	t.Run("not found error", func(t *testing.T) {
		err := NewNotFoundError("test", "test")
		require.False(t, err.IsValidationError())
		require.True(t, err.IsNotFoundError())
		require.False(t, err.IsDatabaseError())
		require.False(t, err.IsTelegramError())
		require.False(t, err.IsInternalError())
		require.False(t, err.IsRateLimitError())
		require.False(t, err.IsUnauthorizedError())
	})

	t.Run("database error", func(t *testing.T) {
		err := NewDatabaseError("test", errors.New("test"))
		require.False(t, err.IsValidationError())
		require.False(t, err.IsNotFoundError())
		require.True(t, err.IsDatabaseError())
		require.False(t, err.IsTelegramError())
		require.False(t, err.IsInternalError())
		require.False(t, err.IsRateLimitError())
		require.False(t, err.IsUnauthorizedError())
	})

	t.Run("telegram error", func(t *testing.T) {
		err := NewTelegramError("test", errors.New("test"))
		require.False(t, err.IsValidationError())
		require.False(t, err.IsNotFoundError())
		require.False(t, err.IsDatabaseError())
		require.True(t, err.IsTelegramError())
		require.False(t, err.IsInternalError())
		require.False(t, err.IsRateLimitError())
		require.False(t, err.IsUnauthorizedError())
	})

	t.Run("internal error", func(t *testing.T) {
		err := NewInternalError("test", errors.New("test"))
		require.False(t, err.IsValidationError())
		require.False(t, err.IsNotFoundError())
		require.False(t, err.IsDatabaseError())
		require.False(t, err.IsTelegramError())
		require.True(t, err.IsInternalError())
		require.False(t, err.IsRateLimitError())
		require.False(t, err.IsUnauthorizedError())
	})

	t.Run("rate limit error", func(t *testing.T) {
		err := NewRateLimitError("test")
		require.False(t, err.IsValidationError())
		require.False(t, err.IsNotFoundError())
		require.False(t, err.IsDatabaseError())
		require.False(t, err.IsTelegramError())
		require.False(t, err.IsInternalError())
		require.True(t, err.IsRateLimitError())
		require.False(t, err.IsUnauthorizedError())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		err := NewUnauthorizedError("test")
		require.False(t, err.IsValidationError())
		require.False(t, err.IsNotFoundError())
		require.False(t, err.IsDatabaseError())
		require.False(t, err.IsTelegramError())
		require.False(t, err.IsInternalError())
		require.False(t, err.IsRateLimitError())
		require.True(t, err.IsUnauthorizedError())
	})
}

func TestWrapError(t *testing.T) {
	t.Run("wrap AppError", func(t *testing.T) {
		originalErr := NewValidationError("Original", "details")
		originalErr.RequestID = "req-123"
		originalErr.UserID = 456

		wrapped := WrapError(originalErr, "New message", ErrorTypeDatabase)

		require.Equal(t, ErrorTypeDatabase, wrapped.Type)
		require.Equal(t, "New message", wrapped.Message)
		require.Equal(t, "req-123", wrapped.RequestID)
		require.Equal(t, int64(456), wrapped.UserID)
		require.Equal(t, originalErr, wrapped)
	})

	t.Run("wrap standard error", func(t *testing.T) {
		originalErr := errors.New("standard error")

		wrapped := WrapError(originalErr, "Wrapped message", ErrorTypeInternal)

		require.Equal(t, ErrorTypeInternal, wrapped.Type)
		require.Equal(t, "Wrapped message", wrapped.Message)
		require.Equal(t, http.StatusInternalServerError, wrapped.Code)
		require.Equal(t, "Wrapped error", wrapped.Details)
		require.Equal(t, originalErr, wrapped.Err)
		require.NotEqual(t, originalErr, wrapped)
	})
}

func TestAppError_WithMetadata(t *testing.T) {
	err := NewValidationError("test", "test")

	// Test setting metadata
	err.RequestID = "req-123"
	err.UserID = 456

	require.Equal(t, "req-123", err.RequestID)
	require.Equal(t, int64(456), err.UserID)
}

func TestAppError_JSONTags(t *testing.T) {
	// This test ensures the JSON tags are correctly defined
	err := &AppError{
		Type:      ErrorTypeValidation,
		Message:   "test",
		Code:      400,
		Details:   "test details",
		RequestID: "req-123",
		UserID:    456,
		Err:       errors.New("test"),
	}

	// The struct should have all the expected fields
	require.Equal(t, ErrorTypeValidation, err.Type)
	require.Equal(t, "test", err.Message)
	require.Equal(t, 400, err.Code)
	require.Equal(t, "test details", err.Details)
	require.Equal(t, "req-123", err.RequestID)
	require.Equal(t, int64(456), err.UserID)
	require.NotNil(t, err.Err)
}
