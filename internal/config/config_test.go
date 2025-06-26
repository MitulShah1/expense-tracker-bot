package config

import (
	"testing"
	"time"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Set up environment variables
	t.Setenv("TELEGRAM_TOKEN", "test_token_123")
	t.Setenv("BOT_ID", "test_bot_id")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("IS_DEV_MODE", "true")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME", "10m")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want no error", err)
	}

	if config == nil {
		t.Fatal("Load() returned nil config")
	}

	// Check all fields
	if config.TelegramToken != "test_token_123" {
		t.Errorf("TelegramToken = %v, want %v", config.TelegramToken, "test_token_123")
	}
	if config.BotID != "test_bot_id" {
		t.Errorf("BotID = %v, want %v", config.BotID, "test_bot_id")
	}
	if config.LogLevel != "debug" {
		t.Errorf("LogLevel = %v, want %v", config.LogLevel, "debug")
	}
	if !config.IsDevMode {
		t.Error("IsDevMode = false, want true")
	}
	if config.DatabaseURL != "postgres://localhost:5432/testdb" {
		t.Errorf("DatabaseURL = %v, want %v", config.DatabaseURL, "postgres://localhost:5432/testdb")
	}
	if config.DBMaxOpenConns != 50 {
		t.Errorf("DBMaxOpenConns = %v, want %v", config.DBMaxOpenConns, 50)
	}
	if config.DBMaxIdleConns != 10 {
		t.Errorf("DBMaxIdleConns = %v, want %v", config.DBMaxIdleConns, 10)
	}
	if config.DBConnMaxLifetime != 10*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want %v", config.DBConnMaxLifetime, 10*time.Minute)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	// Set only required environment variables
	t.Setenv("TELEGRAM_TOKEN", "test_token_123")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want no error", err)
	}

	if config == nil {
		t.Fatal("Load() returned nil config")
	}

	// Check default values
	if config.BotID != "" {
		t.Errorf("BotID = %v, want empty string", config.BotID)
	}
	if config.LogLevel != "" {
		t.Errorf("LogLevel = %v, want empty string", config.LogLevel)
	}
	if config.IsDevMode {
		t.Error("IsDevMode = true, want false")
	}
	if config.DBMaxOpenConns != 25 {
		t.Errorf("DBMaxOpenConns = %v, want %v", config.DBMaxOpenConns, 25)
	}
	if config.DBMaxIdleConns != 5 {
		t.Errorf("DBMaxIdleConns = %v, want %v", config.DBMaxIdleConns, 5)
	}
	if config.DBConnMaxLifetime != 5*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want %v", config.DBConnMaxLifetime, 5*time.Minute)
	}
}

func TestLoad_InvalidDBMaxOpenConns(t *testing.T) {
	// Set required environment variables with invalid DB_MAX_OPEN_CONNS
	t.Setenv("TELEGRAM_TOKEN", "test_token_123")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")
	t.Setenv("DB_MAX_OPEN_CONNS", "invalid")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want no error", err)
	}

	// Should use default value when parsing fails
	if config.DBMaxOpenConns != 25 {
		t.Errorf("DBMaxOpenConns = %v, want %v", config.DBMaxOpenConns, 25)
	}
}

func TestLoad_InvalidDBMaxIdleConns(t *testing.T) {
	// Set required environment variables with invalid DB_MAX_IDLE_CONNS
	t.Setenv("TELEGRAM_TOKEN", "test_token_123")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")
	t.Setenv("DB_MAX_IDLE_CONNS", "invalid")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want no error", err)
	}

	// Should use default value when parsing fails
	if config.DBMaxIdleConns != 5 {
		t.Errorf("DBMaxIdleConns = %v, want %v", config.DBMaxIdleConns, 5)
	}
}

func TestLoad_InvalidDBConnMaxLifetime(t *testing.T) {
	// Set required environment variables with invalid DB_CONN_MAX_LIFETIME
	t.Setenv("TELEGRAM_TOKEN", "test_token_123")
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")
	t.Setenv("DB_CONN_MAX_LIFETIME", "invalid")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want no error", err)
	}

	// Should use default value when parsing fails
	if config.DBConnMaxLifetime != 5*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want %v", config.DBConnMaxLifetime, 5*time.Minute)
	}
}

func TestLoad_MissingTelegramToken(t *testing.T) {
	// Set only DATABASE_URL, missing TELEGRAM_TOKEN
	t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")

	config, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if config != nil {
		t.Fatal("Load() returned config, want nil")
	}

	// Check error message
	expectedError := "invalid configuration: TELEGRAM_TOKEN is required"
	if err.Error() != expectedError {
		t.Errorf("Load() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	// Set only TELEGRAM_TOKEN, missing DATABASE_URL
	t.Setenv("TELEGRAM_TOKEN", "test_token_123")

	config, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if config != nil {
		t.Fatal("Load() returned config, want nil")
	}

	// Check error message
	expectedError := "invalid configuration: DATABASE_URL is required"
	if err.Error() != expectedError {
		t.Errorf("Load() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestLoad_MissingBothRequired(t *testing.T) {
	// Don't set any environment variables
	config, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}

	if config != nil {
		t.Fatal("Load() returned config, want nil")
	}

	// Check error message (should mention TELEGRAM_TOKEN first)
	expectedError := "invalid configuration: TELEGRAM_TOKEN is required"
	if err.Error() != expectedError {
		t.Errorf("Load() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestConfig_IsValid_ValidConfig(t *testing.T) {
	config := &Config{
		TelegramToken: "test_token_123",
		DatabaseURL:   "postgres://localhost:5432/testdb",
	}

	err := config.IsValid()
	if err != nil {
		t.Errorf("IsValid() error = %v, want no error", err)
	}
}

func TestConfig_IsValid_MissingTelegramToken(t *testing.T) {
	config := &Config{
		DatabaseURL: "postgres://localhost:5432/testdb",
	}

	err := config.IsValid()
	if err == nil {
		t.Fatal("IsValid() error = nil, want error")
	}

	expectedError := "TELEGRAM_TOKEN is required"
	if err.Error() != expectedError {
		t.Errorf("IsValid() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestConfig_IsValid_MissingDatabaseURL(t *testing.T) {
	config := &Config{
		TelegramToken: "test_token_123",
	}

	err := config.IsValid()
	if err == nil {
		t.Fatal("IsValid() error = nil, want error")
	}

	expectedError := "DATABASE_URL is required"
	if err.Error() != expectedError {
		t.Errorf("IsValid() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestConfig_IsValid_MissingBoth(t *testing.T) {
	config := &Config{}

	err := config.IsValid()
	if err == nil {
		t.Fatal("IsValid() error = nil, want error")
	}

	// Should mention TELEGRAM_TOKEN first
	expectedError := "TELEGRAM_TOKEN is required"
	if err.Error() != expectedError {
		t.Errorf("IsValid() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestConfig_IsValid_EmptyStrings(t *testing.T) {
	config := &Config{
		TelegramToken: "",
		DatabaseURL:   "",
	}

	err := config.IsValid()
	if err == nil {
		t.Fatal("IsValid() error = nil, want error")
	}

	// Should mention TELEGRAM_TOKEN first
	expectedError := "TELEGRAM_TOKEN is required"
	if err.Error() != expectedError {
		t.Errorf("IsValid() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestLoad_IsDevModeVariations(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     bool
	}{
		{"true value", "true", true},
		{"false value", "false", false},
		{"empty string", "", false},
		{"random string", "random", false},
		{"TRUE uppercase", "TRUE", false},
		{"True mixed case", "True", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set required environment variables
			t.Setenv("TELEGRAM_TOKEN", "test_token_123")
			t.Setenv("DATABASE_URL", "postgres://localhost:5432/testdb")
			t.Setenv("IS_DEV_MODE", tt.envValue)

			config, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v, want no error", err)
			}

			if config.IsDevMode != tt.want {
				t.Errorf("IsDevMode = %v, want %v", config.IsDevMode, tt.want)
			}
		})
	}
}
