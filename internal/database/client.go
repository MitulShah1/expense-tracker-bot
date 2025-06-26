// Package database provides PostgreSQL database operations for the expense tracker bot.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// contextKey is a custom type for context keys in this package
// to avoid collisions with other context keys.
type contextKey string

// Client represents a PostgreSQL database client implementing the Storage interface
type Client struct {
	db     *sqlx.DB
	logger logger.Logger
}

// Ensure Client implements Storage interface
var _ Storage = (*Client)(nil)

// Storage defines the main storage interface that combines all sub-interfaces
type Storage interface {
	UserStorage
	CategoryStorage
	ExpenseStorage

	// Connection management
	Close() error
	GetDB() *sqlx.DB
}

// NewClient creates a new database client
func NewClient(ctx context.Context, databaseURL string, logger logger.Logger) (Storage, error) {
	// Create context with request ID
	reqCtx := context.WithValue(ctx, contextKey("request_id"), "db_init")

	// Connect to database
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info(reqCtx, "Database connection established")

	client := &Client{
		db:     db,
		logger: logger,
	}

	return client, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}

// GetDB returns the underlying database connection
func (c *Client) GetDB() *sqlx.DB {
	return c.db
}
