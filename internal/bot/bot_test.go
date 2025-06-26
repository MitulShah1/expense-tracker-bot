// Package bot implements the Telegram bot functionality for expense tracking.
package bot

import (
	"context"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/require"
)

func TestBot_GetMetrics(t *testing.T) {
	t.Run("should return current metrics", func(t *testing.T) {
		mockDB := database.NewMockStorage()
		mockLogger := logger.NewMockLogger()

		// Create a minimal bot instance for testing
		bot := &Bot{
			db:              mockDB,
			logger:          mockLogger,
			expenseService:  nil, // Not needed for this test
			categoryService: nil, // Not needed for this test
			userService:     nil, // Not needed for this test
			states:          make(map[int64]*models.UserState),
			stateTimeout:    30 * time.Minute,
			rateLimiter:     nil, // Not needed for this test
		}

		// Initialize metrics
		bot.metrics.lastUpdateTime = time.Now()

		// Increment some metrics
		bot.incrementMetric(&bot.metrics.messageCount)
		bot.incrementMetric(&bot.metrics.commandCount)
		bot.incrementMetric(&bot.metrics.errorCount)

		metrics := bot.GetMetrics()
		require.Equal(t, int64(1), metrics["message_count"])
		require.Equal(t, int64(1), metrics["command_count"])
		require.Equal(t, int64(1), metrics["error_count"])
		require.Equal(t, int64(0), metrics["expense_count"])
		require.Equal(t, int64(0), metrics["active_users"])
		require.NotNil(t, metrics["last_update_time"])
	})

	t.Run("should handle concurrent metric access", func(t *testing.T) {
		mockDB := database.NewMockStorage()
		mockLogger := logger.NewMockLogger()

		// Create a minimal bot instance for testing
		bot := &Bot{
			db:              mockDB,
			logger:          mockLogger,
			expenseService:  nil,
			categoryService: nil,
			userService:     nil,
			states:          make(map[int64]*models.UserState),
			stateTimeout:    30 * time.Minute,
			rateLimiter:     nil,
		}

		// Initialize metrics
		bot.metrics.lastUpdateTime = time.Now()

		// Test concurrent access
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				bot.incrementMetric(&bot.metrics.messageCount)
				bot.GetMetrics()
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		metrics := bot.GetMetrics()
		require.Equal(t, int64(10), metrics["message_count"])
	})
}

func TestBot_StateManagement(t *testing.T) {
	t.Run("should get and set state correctly", func(t *testing.T) {
		mockDB := database.NewMockStorage()
		mockLogger := logger.NewMockLogger()

		// Create a minimal bot instance for testing
		bot := &Bot{
			db:              mockDB,
			logger:          mockLogger,
			expenseService:  nil,
			categoryService: nil,
			userService:     nil,
			states:          make(map[int64]*models.UserState),
			stateTimeout:    30 * time.Minute,
			rateLimiter:     nil,
		}

		userID := int64(123)
		state := models.NewUserState()
		state.Step = models.StepCategory

		// Initially no state
		retrievedState := bot.getState(userID)
		require.Nil(t, retrievedState)

		// Set state
		bot.setState(userID, state)

		// Retrieve state
		retrievedState = bot.getState(userID)
		require.NotNil(t, retrievedState)
		require.Equal(t, models.StepCategory, retrievedState.Step)
	})

	t.Run("should handle concurrent state access", func(t *testing.T) {
		mockDB := database.NewMockStorage()
		mockLogger := logger.NewMockLogger()

		// Create a minimal bot instance for testing
		bot := &Bot{
			db:              mockDB,
			logger:          mockLogger,
			expenseService:  nil,
			categoryService: nil,
			userService:     nil,
			states:          make(map[int64]*models.UserState),
			stateTimeout:    30 * time.Minute,
			rateLimiter:     nil,
		}

		userID := int64(123)
		done := make(chan bool, 10)

		// Test concurrent state access
		for i := 0; i < 10; i++ {
			go func() {
				state := models.NewUserState()
				state.Step = models.StepCategory
				bot.setState(userID, state)
				bot.getState(userID)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify state exists
		retrievedState := bot.getState(userID)
		require.NotNil(t, retrievedState)
	})
}

func TestBot_CleanupRoutine(t *testing.T) {
	t.Run("should cleanup expired states", func(t *testing.T) {
		ctx := context.Background()
		mockDB := database.NewMockStorage()
		mockLogger := logger.NewMockLogger()

		// Create a minimal bot instance for testing
		bot := &Bot{
			db:              mockDB,
			logger:          mockLogger,
			expenseService:  nil,
			categoryService: nil,
			userService:     nil,
			states:          make(map[int64]*models.UserState),
			stateTimeout:    30 * time.Minute,
			rateLimiter:     nil,
		}

		// Set up expired state (more than 30 minutes ago)
		userID := int64(123)
		state := models.NewUserState()
		state.LastActivity = time.Now().Add(-35 * time.Minute) // Expired

		// Set state directly to avoid LastActivity being updated
		bot.stateMutex.Lock()
		bot.states[userID] = state
		bot.stateMutex.Unlock()

		// Verify state exists
		retrievedState := bot.getState(userID)
		require.NotNil(t, retrievedState)

		// Run cleanup
		bot.cleanupExpiredStates(ctx)

		// Verify state was cleaned up
		retrievedState = bot.getState(userID)
		require.Nil(t, retrievedState)
	})

	t.Run("should not cleanup active states", func(t *testing.T) {
		ctx := context.Background()
		mockDB := database.NewMockStorage()
		mockLogger := logger.NewMockLogger()

		// Create a minimal bot instance for testing
		bot := &Bot{
			db:              mockDB,
			logger:          mockLogger,
			expenseService:  nil,
			categoryService: nil,
			userService:     nil,
			states:          make(map[int64]*models.UserState),
			stateTimeout:    30 * time.Minute,
			rateLimiter:     nil,
		}

		// Set up active state (less than 30 minutes ago)
		userID := int64(123)
		state := models.NewUserState()
		state.LastActivity = time.Now().Add(-10 * time.Minute) // Still active

		// Set state directly to avoid LastActivity being updated
		bot.stateMutex.Lock()
		bot.states[userID] = state
		bot.stateMutex.Unlock()

		// Verify state exists
		retrievedState := bot.getState(userID)
		require.NotNil(t, retrievedState)

		// Run cleanup
		bot.cleanupExpiredStates(ctx)

		// Verify state still exists
		retrievedState = bot.getState(userID)
		require.NotNil(t, retrievedState)
	})
}
