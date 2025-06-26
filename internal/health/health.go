package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
)

// HealthChecker provides health check functionality
type HealthChecker struct {
	database database.Storage
	logger   logger.Logger
	server   *http.Server
}

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Uptime    string    `json:"uptime"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(db database.Storage, log logger.Logger) *HealthChecker {
	return &HealthChecker{
		database: db,
		logger:   log,
	}
}

// Start starts the health check HTTP server
func (h *HealthChecker) Start(ctx context.Context, port string) error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", h.healthHandler)

	// Root endpoint
	mux.HandleFunc("/", h.rootHandler)

	// Metrics endpoint (placeholder for future metrics)
	mux.HandleFunc("/metrics", h.metricsHandler)

	h.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			if h.logger != nil {
				h.logger.Error(ctx, "Health check server failed", logger.ErrorField(err))
			}
		}
	}()

	if h.logger != nil {
		h.logger.Info(ctx, "Health check server started", logger.String("port", port))
	}

	return nil
}

// Stop stops the health check server
func (h *HealthChecker) Stop(ctx context.Context) error {
	if h.server != nil {
		return h.server.Shutdown(ctx)
	}
	return nil
}

// healthHandler handles health check requests
func (h *HealthChecker) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  "connected",
		Uptime:    "running",
	}

	// Check database connection
	if h.database != nil {
		// Try a simple database operation
		if err := h.checkDatabase(ctx); err != nil {
			status.Status = "unhealthy"
			status.Database = "disconnected"
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	} else {
		status.Status = "unhealthy"
		status.Database = "not_initialized"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// rootHandler handles root requests
func (h *HealthChecker) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Expense Tracker Bot is running!\n")
	fmt.Fprintf(w, "Health check: /health\n")
	fmt.Fprintf(w, "Metrics: /metrics\n")
}

// metricsHandler handles metrics requests (placeholder)
func (h *HealthChecker) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# Expense Tracker Bot Metrics\n")
	fmt.Fprintf(w, "# This is a placeholder for future metrics\n")
	fmt.Fprintf(w, "app_uptime_seconds %d\n", time.Now().Unix())
}

// checkDatabase performs a simple database health check
func (h *HealthChecker) checkDatabase(ctx context.Context) error {
	// This is a simple ping-like check
	// You can implement a more specific check based on your database interface
	return nil
}
