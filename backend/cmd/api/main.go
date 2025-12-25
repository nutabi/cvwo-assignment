package api

import (
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/app"
	"github.com/nutabi/cvwo-assignment/backend/internal/config"
)

// Entry point for the API server
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		panic(err)
	}

	// Initialise application with configuration
	app := app.Initialise(cfg)

	// Start the API server
	app.Start()
}
