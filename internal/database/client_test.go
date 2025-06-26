package database

import (
	"context"
	"os"
	"testing"

	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestNewClient_InvalidURL(t *testing.T) {
	ctx := context.Background()
	log := logger.NewMockLogger()
	_, err := NewClient(ctx, "invalid-url", log)
	require.Error(t, err)
}

func TestNewClient_ValidURL(t *testing.T) {
	// This test is skipped by default because it requires a real PostgreSQL instance.
	// To enable, set TEST_DATABASE_URL env var to a valid PostgreSQL URL.
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping integration test")
	}
	ctx := context.Background()
	log := logger.NewMockLogger()
	client, err := NewClient(ctx, dbURL, log)
	require.NoError(t, err)
	require.NotNil(t, client)

	db := client.GetDB()
	require.NotNil(t, db)

	require.NoError(t, client.Close())
}

func TestClient_InterfaceCompliance(t *testing.T) {
	var _ Storage = &Client{}
}
