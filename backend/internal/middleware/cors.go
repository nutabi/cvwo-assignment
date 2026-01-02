package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCORSConfig sets up CORS middleware based on environment
func NewCORSConfig(isDebug bool, allowedOrigins []string) gin.HandlerFunc {
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if isDebug {
		// Development: Allow all origins for easier testing
		config.AllowAllOrigins = true
		slog.Info("CORS configured for development mode (allowing all origins)")
	} else {
		// Production: Use configured origins or fallback to restrictive default
		if len(allowedOrigins) > 0 {
			config.AllowOrigins = allowedOrigins
			slog.Info("CORS configured for production mode", "allowed_origins", config.AllowOrigins)
		} else {
			// Fallback: very restrictive if no origins configured
			config.AllowOrigins = []string{"https://localhost"}
			slog.Warn("CORS configured with default restrictive origins. Set CORS_ALLOWED_ORIGINS environment variable.")
		}
	}

	return cors.New(config)
}
