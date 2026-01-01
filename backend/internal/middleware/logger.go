package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// NewRequestLogger returns a Gin middleware that logs endpoint usage at debug level.
// It logs the HTTP method, path, status code, latency, and client IP.
func NewRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Build log attributes
		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
		}

		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		// Log at debug level for endpoint usage
		slog.Debug("request processed", attrs...)
	}
}
