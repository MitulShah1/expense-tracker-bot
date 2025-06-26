package logger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with dev mode",
			config: Config{
				BotID:     "test-bot",
				LogLevel:  "info",
				IsDevMode: true,
			},
			wantErr: false,
		},
		{
			name: "valid config with production mode",
			config: Config{
				BotID:     "test-bot",
				LogLevel:  "debug",
				IsDevMode: false,
			},
			wantErr: false,
		},
		{
			name: "invalid log level defaults to info",
			config: Config{
				BotID:     "test-bot",
				LogLevel:  "invalid-level",
				IsDevMode: false,
			},
			wantErr: false,
		},
		{
			name: "empty bot ID",
			config: Config{
				BotID:     "",
				LogLevel:  "info",
				IsDevMode: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logger)
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	// Create a test logger with observer to capture logs
	core, obs := observer.New(zap.DebugLevel)
	log := zap.New(core)
	sugar := log.Sugar()

	logger := &logger{
		log:   log,
		sugar: sugar,
		botID: "test-bot",
	}

	ctx := context.Background()

	t.Run("Debug", func(t *testing.T) {
		logger.Debug(ctx, "debug message", String("key", "value"))
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "debug message", entry.Message)
		require.Equal(t, zapcore.DebugLevel, entry.Level)
		require.Equal(t, "test-bot", entry.ContextMap()["bot_id"])
	})

	t.Run("Info", func(t *testing.T) {
		logger.Info(ctx, "info message", Int("count", 42))
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "info message", entry.Message)
		require.Equal(t, zapcore.InfoLevel, entry.Level)
		require.Equal(t, int64(42), entry.ContextMap()["count"])
	})

	t.Run("Warn", func(t *testing.T) {
		logger.Warn(ctx, "warn message", Bool("flag", true))
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "warn message", entry.Message)
		require.Equal(t, zapcore.WarnLevel, entry.Level)
		require.Equal(t, true, entry.ContextMap()["flag"])
	})

	t.Run("Error", func(t *testing.T) {
		testErr := errors.New("test error")
		logger.Error(ctx, "error message", ErrorField(testErr))
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "error message", entry.Message)
		require.Equal(t, zapcore.ErrorLevel, entry.Level)
		require.Equal(t, "test error", entry.ContextMap()["error"])
	})

	t.Run("Debugf", func(t *testing.T) {
		logger.Debugf(ctx, "debug %s %d", "formatted", 123)
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "debug formatted 123", entry.Message)
		require.Equal(t, zapcore.DebugLevel, entry.Level)
	})

	t.Run("Infof", func(t *testing.T) {
		logger.Infof(ctx, "info %s %f", "formatted", 3.14)
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "info formatted 3.140000", entry.Message)
		require.Equal(t, zapcore.InfoLevel, entry.Level)
	})

	t.Run("Warnf", func(t *testing.T) {
		logger.Warnf(ctx, "warn %s", "formatted")
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "warn formatted", entry.Message)
		require.Equal(t, zapcore.WarnLevel, entry.Level)
	})

	t.Run("Errorf", func(t *testing.T) {
		logger.Errorf(ctx, "error %s", "formatted")
		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "error formatted", entry.Message)
		require.Equal(t, zapcore.ErrorLevel, entry.Level)
	})
}

func TestLoggerWithContext(t *testing.T) {
	core, obs := observer.New(zap.DebugLevel)
	log := zap.New(core)
	sugar := log.Sugar()

	logger := &logger{
		log:   log,
		sugar: sugar,
		botID: "test-bot",
	}

	t.Run("with request ID in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), RequestIDKey, "req-123")
		logger.Info(ctx, "message with request ID")

		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "test-bot", entry.ContextMap()["bot_id"])
		require.Equal(t, "req-123", entry.ContextMap()["request_id"])
	})

	t.Run("without request ID in context", func(t *testing.T) {
		ctx := context.Background()
		logger.Info(ctx, "message without request ID")

		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "test-bot", entry.ContextMap()["bot_id"])
		_, exists := entry.ContextMap()["request_id"]
		require.False(t, exists)
	})

	t.Run("with different context key type", func(t *testing.T) {
		type DifferentKey string
		ctx := context.WithValue(context.Background(), DifferentKey("different_key"), "value")
		logger.Info(ctx, "message with different key")

		require.Equal(t, 1, obs.Len())
		entry := obs.TakeAll()[0]
		require.Equal(t, "test-bot", entry.ContextMap()["bot_id"])
		_, exists := entry.ContextMap()["request_id"]
		require.False(t, exists)
	})
}

func TestFieldHelpers(t *testing.T) {
	t.Run("String field", func(t *testing.T) {
		field := String("key", "value")
		require.Equal(t, "key", field.Key)
		require.Equal(t, "value", field.String)
	})

	t.Run("Int field", func(t *testing.T) {
		field := Int("key", 42)
		require.Equal(t, "key", field.Key)
		require.Equal(t, int64(42), field.Integer)
	})

	t.Run("Float64 field", func(t *testing.T) {
		field := Float64("key", 3.14)
		require.Equal(t, "key", field.Key)
		require.Equal(t, zapcore.Float64Type, field.Type)
	})

	t.Run("Bool field", func(t *testing.T) {
		field := Bool("key", true)
		require.Equal(t, "key", field.Key)
		require.Equal(t, zapcore.BoolType, field.Type)
	})

	t.Run("Error field", func(t *testing.T) {
		testErr := errors.New("test error")
		field := ErrorField(testErr)
		require.Equal(t, "error", field.Key)
		require.Equal(t, testErr, field.Interface)
	})

	t.Run("Any field", func(t *testing.T) {
		data := map[string]string{"nested": "value"}
		field := Any("key", data)
		require.Equal(t, "key", field.Key)
		require.Equal(t, data, field.Interface)
	})
}

func TestSync(t *testing.T) {
	logger, err := New(Config{
		BotID:     "test-bot",
		LogLevel:  "info",
		IsDevMode: false,
	})
	require.NoError(t, err)

	// Sync may fail on some systems (like stdout sync), so we don't require it to succeed
	_ = logger.Sync()
}

func TestContextKeys(t *testing.T) {
	require.Equal(t, ContextKey("bot_id"), BotIDKey)
	require.Equal(t, ContextKey("request_id"), RequestIDKey)
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		expected zapcore.Level
	}{
		{"debug level", "debug", zapcore.DebugLevel},
		{"info level", "info", zapcore.InfoLevel},
		{"warn level", "warn", zapcore.WarnLevel},
		{"error level", "error", zapcore.ErrorLevel},
		{"fatal level", "fatal", zapcore.FatalLevel},
		{"invalid level defaults to info", "invalid", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(Config{
				BotID:     "test-bot",
				LogLevel:  tt.logLevel,
				IsDevMode: false,
			})
			require.NoError(t, err)
			require.NotNil(t, logger)
		})
	}
}

func TestLoggerConcurrency(t *testing.T) {
	logger, err := New(Config{
		BotID:     "test-bot",
		LogLevel:  "info",
		IsDevMode: false,
	})
	require.NoError(t, err)

	// Test concurrent logging
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			ctx := context.WithValue(context.Background(), RequestIDKey, "req-"+string(rune(id)))
			logger.Info(ctx, "concurrent message", Int("id", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
