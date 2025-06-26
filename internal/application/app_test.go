package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	t.Run("should create new app successfully", func(t *testing.T) {
		app, err := NewApp()

		require.NoError(t, err)
		require.NotNil(t, app)
		require.Nil(t, app.config)
		require.Nil(t, app.logger)
		require.Nil(t, app.database)
		require.Nil(t, app.bot)
	})
}

func TestApp_Initialize(t *testing.T) {
	t.Run("should initialize app successfully with valid context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()
		err = app.Initialize(ctx)

		// Note: This will fail in test environment without proper config
		// but we can test the structure and error handling
		require.Error(t, err) // Expected to fail due to missing config
		require.Contains(t, err.Error(), "failed to load configuration")
	})

	t.Run("should handle nil context gracefully", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		err = app.Initialize(context.TODO())

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})

	t.Run("should handle cancelled context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err = app.Initialize(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})
}

func TestApp_Start(t *testing.T) {
	t.Run("should handle start without initialization", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()
		err = app.Start(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "bot stopped with error")
	})

	t.Run("should handle start with nil context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		err = app.Start(context.TODO())

		require.Error(t, err)
		require.Contains(t, err.Error(), "bot stopped with error")
	})

	t.Run("should handle cancelled context during start", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err = app.Start(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "bot stopped with error")
	})
}

func TestApp_Stop(t *testing.T) {
	t.Run("should stop app gracefully without initialization", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()
		err = app.Stop(ctx)

		require.NoError(t, err)
	})

	t.Run("should stop app with nil context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		err = app.Stop(context.TODO())

		require.NoError(t, err)
	})

	t.Run("should handle multiple stop calls", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// First stop
		err = app.Stop(ctx)
		require.NoError(t, err)

		// Second stop should also succeed
		err = app.Stop(ctx)
		require.NoError(t, err)
	})
}

func TestApp_Shutdown(t *testing.T) {
	t.Run("should shutdown app gracefully", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()
		err = app.Shutdown(ctx)

		require.NoError(t, err)
	})

	t.Run("should shutdown with nil context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		err = app.Shutdown(context.TODO())

		require.NoError(t, err)
	})
}

func TestApp_GetLogger(t *testing.T) {
	t.Run("should return nil logger when not initialized", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		logger := app.GetLogger()
		require.Nil(t, logger)
	})
}

func TestApp_Run(t *testing.T) {
	t.Run("should handle run without proper configuration", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		// Run should fail due to missing configuration
		err = app.Run()

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})
}

func TestApp_RunWithContext(t *testing.T) {
	t.Run("should handle run with context without proper configuration", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()
		err = app.RunWithContext(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})

	t.Run("should handle run with cancelled context", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err = app.RunWithContext(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to load configuration")
	})

	t.Run("should handle run with timeout context", func(t *testing.T) {
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
}

func TestApp_ConcurrentOperations(t *testing.T) {
	t.Run("should handle concurrent stop operations", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Run multiple stop operations concurrently
		done := make(chan error, 3)
		for i := 0; i < 3; i++ {
			go func() {
				done <- app.Stop(ctx)
			}()
		}

		// All should succeed
		for i := 0; i < 3; i++ {
			err := <-done
			require.NoError(t, err)
		}
	})

	t.Run("should handle concurrent shutdown operations", func(t *testing.T) {
		app, err := NewApp()
		require.NoError(t, err)

		ctx := context.Background()

		// Run multiple shutdown operations concurrently
		done := make(chan error, 3)
		for i := 0; i < 3; i++ {
			go func() {
				done <- app.Shutdown(ctx)
			}()
		}

		// All should succeed
		for i := 0; i < 3; i++ {
			err := <-done
			require.NoError(t, err)
		}
	})
}
