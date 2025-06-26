// Package application provides the main application struct and lifecycle management for the expense tracker bot.
package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
)

// Run starts the application and handles graceful shutdown
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize application
	if err := a.Initialize(ctx); err != nil {
		return err
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to receive errors from the bot
	errChan := make(chan error, 1)

	// Start the application in a goroutine
	go func() {
		if err := a.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	// Wait for either a signal or an error
	select {
	case sig := <-sigChan:
		if a.logger != nil {
			a.logger.Info(ctx, "Received shutdown signal", logger.String("signal", sig.String()))
		}
	case err := <-errChan:
		if a.logger != nil {
			a.logger.Error(ctx, "Application error", logger.ErrorField(err))
		}
		// Still perform graceful shutdown even on error
	}

	// Perform graceful shutdown
	return a.Stop(ctx)
}

// RunWithContext runs the application with a custom context
func (a *App) RunWithContext(ctx context.Context) error {
	// Initialize application
	if err := a.Initialize(ctx); err != nil {
		return err
	}
	// Start the application
	return a.Start(ctx)
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	return a.Stop(ctx)
}
