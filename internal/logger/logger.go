// Package logger provides logging functionality for the application.
package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Debugf(ctx context.Context, format string, args ...any)
	Infof(ctx context.Context, format string, args ...any)
	Warnf(ctx context.Context, format string, args ...any)
	Errorf(ctx context.Context, format string, args ...any)
	Fatalf(ctx context.Context, format string, args ...any)
	Sync() error
}

// ContextKey is a type for context keys
type ContextKey string

const (
	// BotIDKey is the context key for bot ID
	BotIDKey ContextKey = "bot_id"
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
)

// Config holds the logger configuration
type Config struct {
	BotID     string
	LogLevel  string
	IsDevMode bool
}

// logger implements the Logger interface
type logger struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
	botID string
}

// New creates a new logger instance
func New(cfg Config) (Logger, error) {
	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Set log level
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level = zap.InfoLevel
	}

	// Create core with stdout output
	var core zapcore.Core
	if cfg.IsDevMode {
		// Development mode - console output with colors
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	} else {
		// Production mode - JSON output
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	}

	// Create logger
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := log.Sugar()

	return &logger{
		log:   log,
		sugar: sugar,
		botID: cfg.BotID,
	}, nil
}

// Sync flushes any buffered log entries
func (l *logger) Sync() error {
	return l.log.Sync()
}

// Fatal logs a fatal message and exits
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Fatal(msg, fields...)
}

// Fatalf logs a fatal message with formatting and exits
func (l *logger) Fatalf(ctx context.Context, format string, args ...any) {
	l.withContextSugar(ctx).Fatalf(format, args...)
}

// Error logs an error message
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Error(msg, fields...)
}

// Errorf logs an error message with formatting
func (l *logger) Errorf(ctx context.Context, format string, args ...any) {
	l.withContextSugar(ctx).Errorf(format, args...)
}

// Warn logs a warning message
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Warn(msg, fields...)
}

// Warnf logs a warning message with formatting
func (l *logger) Warnf(ctx context.Context, format string, args ...any) {
	l.withContextSugar(ctx).Warnf(format, args...)
}

// Debug logs a debug message with context
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Debug(msg, fields...)
}

// Info logs an info message with context
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.withContext(ctx).Info(msg, fields...)
}

// Debugf logs a debug message with formatting and context
func (l *logger) Debugf(ctx context.Context, format string, args ...any) {
	l.withContextSugar(ctx).Debugf(format, args...)
}

// Infof logs an info message with formatting and context
func (l *logger) Infof(ctx context.Context, format string, args ...any) {
	l.withContextSugar(ctx).Infof(format, args...)
}

// String creates a zap.Field for a string value
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates a zap.Field for an int value
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Float64 creates a zap.Field for a float64 value
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool creates a zap.Field for a bool value
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// ErrorField creates a zap.Field for an error value
func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

// Any creates a zap.Field for any value
func Any(key string, val any) zap.Field {
	return zap.Any(key, val)
}

// MockLogger implements the Logger interface for testing
type MockLogger struct{}

// NewMockLogger creates a new mock logger instance
func NewMockLogger() Logger {
	return &MockLogger{}
}

// Debug implements Logger.Debug (no-op for testing)
func (m *MockLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {}

// Info implements Logger.Info (no-op for testing)
func (m *MockLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {}

// Warn implements Logger.Warn (no-op for testing)
func (m *MockLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {}

// Error implements Logger.Error (no-op for testing)
func (m *MockLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}

// Fatal implements Logger.Fatal (no-op for testing)
func (m *MockLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {}

// Debugf implements Logger.Debugf (no-op for testing)
func (m *MockLogger) Debugf(ctx context.Context, format string, args ...any) {}

// Infof implements Logger.Infof (no-op for testing)
func (m *MockLogger) Infof(ctx context.Context, format string, args ...any) {}

// Warnf implements Logger.Warnf (no-op for testing)
func (m *MockLogger) Warnf(ctx context.Context, format string, args ...any) {}

// Errorf implements Logger.Errorf (no-op for testing)
func (m *MockLogger) Errorf(ctx context.Context, format string, args ...any) {}

// Fatalf implements Logger.Fatalf (no-op for testing)
func (m *MockLogger) Fatalf(ctx context.Context, format string, args ...any) {}

// Sync implements Logger.Sync (no-op for testing)
func (m *MockLogger) Sync() error {
	return nil
}

// Private methods
func (l *logger) withContext(ctx context.Context) *zap.Logger {
	fields := []zap.Field{
		zap.String("bot_id", l.botID),
	}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields = append(fields, zap.String("request_id", reqID))
	}

	return l.log.With(fields...)
}

func (l *logger) withContextSugar(ctx context.Context) *zap.SugaredLogger {
	fields := []any{
		"bot_id", l.botID,
	}

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields = append(fields, "request_id", reqID)
	}

	return l.sugar.With(fields...)
}
