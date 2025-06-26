// Package config provides configuration management for the expense tracker bot.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	TelegramToken string
	BotID         string
	LogLevel      string
	IsDevMode     bool

	// Database Configuration
	DatabaseURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Parse database connection pool settings
	dbMaxOpenConns := 25 // default
	if val := os.Getenv("DB_MAX_OPEN_CONNS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			dbMaxOpenConns = parsed
		}
	}

	dbMaxIdleConns := 5 // default
	if val := os.Getenv("DB_MAX_IDLE_CONNS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			dbMaxIdleConns = parsed
		}
	}

	dbConnMaxLifetime := 5 * time.Minute // default
	if val := os.Getenv("DB_CONN_MAX_LIFETIME"); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			dbConnMaxLifetime = parsed
		}
	}

	cnfg := &Config{
		TelegramToken:     os.Getenv("TELEGRAM_TOKEN"),
		BotID:             os.Getenv("BOT_ID"),
		LogLevel:          os.Getenv("LOG_LEVEL"),
		IsDevMode:         os.Getenv("IS_DEV_MODE") == "true",
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		DBMaxOpenConns:    dbMaxOpenConns,
		DBMaxIdleConns:    dbMaxIdleConns,
		DBConnMaxLifetime: dbConnMaxLifetime,
	}

	if err := cnfg.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return cnfg, nil
}

// IsValid checks if the configuration is valid
func (cfg *Config) IsValid() error {
	// Validate required configuration
	if cfg.TelegramToken == "" {
		return errors.New("TELEGRAM_TOKEN is required")
	}
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	return nil
}
