package main

import (
	"log"

	"github.com/MitulShah1/expense-tracker-bot/internal/application"
)

func main() {
	// Create new application instance
	app, err := application.NewApp()
	if err != nil {
		log.Fatalf("Failed to create application: %v\n", err)
	}

	// Run the application (handles initialization, startup, and graceful shutdown)
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v\n", err)
	}
}
