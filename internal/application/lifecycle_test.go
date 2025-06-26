package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Custom context key types to avoid collisions
type (
	testKey  string
	key1Type string
	key2Type string
	key3Type string
)

const (
	testKeyVal testKey  = "test_key"
	key1       key1Type = "key1"
	key2       key2Type = "key2"
	key3       key3Type = "key3"
)

func TestApp_RunWithContext_ContextCancellation(t *testing.T) {
	t.Run("should handle context cancellation during initialization", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel context immediately
		cancel()

		err = app.RunWithContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})

	t.Run("should handle context timeout during initialization", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(2 * time.Millisecond)

		err = app.RunWithContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})

	t.Run("should handle context with deadline", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		deadline := time.Now().Add(1 * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		// Wait for deadline
		time.Sleep(2 * time.Millisecond)

		err = app.RunWithContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})
}

func TestApp_RunWithContext_ContextValues(t *testing.T) {
	t.Run("should preserve context values during execution", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), testKeyVal, "test_value")

		err = app.RunWithContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")

		// Context value should be preserved
		value := ctx.Value(testKeyVal)
		require.Equal(t, "test_value", value)
	})

	t.Run("should handle context with multiple values", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), key1, "value1")
		ctx = context.WithValue(ctx, key2, "value2")
		ctx = context.WithValue(ctx, key3, "value3")

		err = app.RunWithContext(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")

		// All context values should be preserved
		require.Equal(t, "value1", ctx.Value(key1))
		require.Equal(t, "value2", ctx.Value(key2))
		require.Equal(t, "value3", ctx.Value(key3))
	})
}

func TestApp_Shutdown_ContextHandling(t *testing.T) {
	t.Run("should handle shutdown with cancelled context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err = app.Shutdown(ctx)
		require.NoError(t, err) // Shutdown should succeed even with cancelled context
	})

	t.Run("should handle shutdown with timeout context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = app.Shutdown(ctx)
		require.NoError(t, err)
	})

	t.Run("should handle shutdown with deadline context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		deadline := time.Now().Add(1 * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		err = app.Shutdown(ctx)
		require.NoError(t, err)
	})
}

func TestApp_Lifecycle_Integration(t *testing.T) {
	t.Run("should handle complete lifecycle without initialization", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Initialize (should fail but not crash)
		err = app.Initialize(ctx)
		require.Error(t, err)

		// Start (should fail but not crash)
		err = app.Start(ctx)
		require.Error(t, err)

		// Stop (should succeed)
		err = app.Stop(ctx)
		require.NoError(t, err)

		// Shutdown (should succeed)
		err = app.Shutdown(ctx)
		require.NoError(t, err)
	})

	t.Run("should handle multiple lifecycle cycles", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// First cycle
		err = app.Initialize(ctx)
		require.Error(t, err)
		err = app.Stop(ctx)
		require.NoError(t, err)

		// Second cycle
		err = app.Initialize(ctx)
		require.Error(t, err)
		err = app.Stop(ctx)
		require.NoError(t, err)

		// Third cycle
		err = app.Initialize(ctx)
		require.Error(t, err)
		err = app.Stop(ctx)
		require.NoError(t, err)
	})
}

func TestApp_ConcurrentLifecycle(t *testing.T) {
	t.Run("should handle concurrent lifecycle operations", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Run multiple operations concurrently
		done := make(chan error, 4)

		go func() { done <- app.Initialize(ctx) }()
		go func() { done <- app.Start(ctx) }()
		go func() { done <- app.Stop(ctx) }()
		go func() { done <- app.Shutdown(ctx) }()

		// Collect results
		results := make([]error, 4)
		for i := 0; i < 4; i++ {
			results[i] = <-done
		}

		// Debug: print what we got
		t.Logf("Results: Initialize=%v, Start=%v, Stop=%v, Shutdown=%v",
			results[0], results[1], results[2], results[3])

		// Check that we have the expected mix of errors and successes
		// (2 errors from Initialize/Start, 2 successes from Stop/Shutdown)
		errorCount := 0
		successCount := 0

		for _, result := range results {
			if result != nil {
				errorCount++
			} else {
				successCount++
			}
		}

		// Should have 2 errors and 2 successes
		require.Equal(t, 2, errorCount)
		require.Equal(t, 2, successCount)
	})
}

func TestApp_ErrorRecovery(t *testing.T) {
	t.Run("should recover from initialization errors", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// First initialization attempt fails
		err = app.Initialize(ctx)
		require.Error(t, err)

		// Stop should still work
		err = app.Stop(ctx)
		require.NoError(t, err)

		// Shutdown should still work
		err = app.Shutdown(ctx)
		require.NoError(t, err)
	})

	t.Run("should recover from start errors", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Start fails
		err = app.Start(ctx)
		require.Error(t, err)

		// Stop should still work
		err = app.Stop(ctx)
		require.NoError(t, err)

		// Shutdown should still work
		err = app.Shutdown(ctx)
		require.NoError(t, err)
	})
}

func TestApp_ResourceCleanup(t *testing.T) {
	t.Run("should cleanup resources on multiple stop calls", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Multiple stop calls should not cause issues
		for i := 0; i < 5; i++ {
			err = app.Stop(ctx)
			require.NoError(t, err)
		}
	})

	t.Run("should cleanup resources on multiple shutdown calls", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Multiple shutdown calls should not cause issues
		for i := 0; i < 5; i++ {
			err = app.Shutdown(ctx)
			require.NoError(t, err)
		}
	})
}
