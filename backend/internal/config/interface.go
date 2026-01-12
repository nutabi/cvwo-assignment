package config

import "log/slog"

// Config interface defines methods to access configuration settings
type Config interface {
	// Get server hostname for cookies/domain, like localhost or example.com.
	GetServerHostname() string

	// Get server listening port, like 8080.
	GetServerPort() int

	// Get full server listen address (always 0.0.0.0:port for Docker).
	GetServerAddress() string

	// Get database connection string.
	GetDSN() string

	// Check if the application is running in debug mode.
	IsDebug() bool

	// Get JWT secret key for authentication.
	GetJWTSecretKey() string

	// Get allowed CORS origins for production.
	GetCORSOrigins() []string

	// Get log level for structured logging.
	GetLogLevel() slog.Level
}
