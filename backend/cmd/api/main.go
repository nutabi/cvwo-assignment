package api

import (
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/config"
)

// Entry point for the API server
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		panic(err)
	}

	// TODO: Initialize application with configuration

	// TODO: Start the API server
}
