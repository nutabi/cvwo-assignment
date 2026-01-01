package config

// Config interface defines methods to access configuration settings
type Config interface {
	// Get server hostname, like localhost or example.com.
	GetServerHostname() string

	// Get server listening port, like 8080.
	GetServerPort() int

	// Get full server address in the format hostname:port.
	GetServerAddress() string

	// Get database connection string.
	GetDSN() string

	// Check if the application is running in debug mode.
	IsDebug() bool

	// Get JWT secret key for authentication.
	GetJWTSecretKey() string

	// Get allowed CORS origins for production.
	GetCORSOrigins() []string
}
