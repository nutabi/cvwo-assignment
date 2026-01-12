// @title           CVWO Forum API
// @version         1.0
// @description     API for CVWO forum assignment backend
// @description
// @description     ## Error Responses
// @description     All error responses follow a consistent structure with two fields:
// @description     - `error_code`: Machine-readable error code (e.g., INVALID_INPUT, USER_NOT_FOUND)
// @description     - `error_message`: Human-readable error message
// @description
// @description     Example: `{"error_code": "INVALID_INPUT", "error_message": "invalid input"}`
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/app"
	"github.com/nutabi/cvwo-assignment/backend/internal/config"
)

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
