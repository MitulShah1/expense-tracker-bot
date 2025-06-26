// Package application provides the main application struct and lifecycle management for the expense tracker bot.
package application

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MitulShah1/expense-tracker-bot/internal/bot"
	"github.com/MitulShah1/expense-tracker-bot/internal/config"
	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
)

// App represents the main application
// It holds configuration, logger, database, and bot dependencies.
type App struct {
	config   *config.Config
	logger   logger.Logger
	database database.Storage
	bot      *bot.Bot
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	app := &App{}
	return app, nil
}

// Initialize sets up all application dependencies
func (a *App) Initialize(ctx context.Context) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	a.config = cfg

	// Initialize logger
	loggerLog, err := logger.New(logger.Config{
		BotID:     cfg.BotID,
		LogLevel:  cfg.LogLevel,
		IsDevMode: cfg.IsDevMode,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	a.logger = loggerLog

	// Initialize database storage
	dbStorage, err := database.NewClient(ctx, cfg.DatabaseURL, loggerLog)
	if err != nil {
		return fmt.Errorf("failed to initialize database storage: %w", err)
	}
	a.database = dbStorage

	// Initialize bot
	botInstance, err := bot.NewBot(ctx, cfg.TelegramToken, dbStorage, loggerLog)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}
	a.bot = botInstance

	return nil
}

// Start begins the application lifecycle
func (a *App) Start(ctx context.Context) error {
	if a.logger != nil {
		a.logger.Info(ctx, "Starting expense tracker bot application")
	}

	// Check if bot is initialized
	if a.bot == nil {
		return errors.New("bot stopped with error: bot not initialized")
	}

	// Start bot
	if err := a.bot.Start(ctx); err != nil {
		return fmt.Errorf("bot stopped with error: %w", err)
	}
	return nil
}

// Stop gracefully shuts down the application
func (a *App) Stop(ctx context.Context) error {
	// Log shutdown if logger is available
	if a.logger != nil {
		a.logger.Info(ctx, "Stopping expense tracker bot application")
	}
	// Close database connection
	if a.database != nil {
		if err := a.database.Close(); err != nil {
			if a.logger != nil {
				a.logger.Error(ctx, "Failed to close database connection", logger.ErrorField(err))
			}
		}
	}
	// Sync logger
	if a.logger != nil {
		if err := a.logger.Sync(); err != nil {
			// Ignore sync errors on stdout/stderr
			if !strings.Contains(err.Error(), "sync /dev/stdout") && !strings.Contains(err.Error(), "sync /dev/stderr") {
				a.logger.Error(ctx, "Failed to sync logger", logger.ErrorField(err))
			}
		}
		a.logger.Info(ctx, "Application stopped successfully")
	}
	return nil
}

// GetLogger returns the application logger
func (a *App) GetLogger() logger.Logger {
	return a.logger
}
