package config

// Config interface defines methods to access configuration settings
type Config interface {
	// Get server hostname, like localhost or example.com.
	GetServerHostname() string

	// Get server listening port, like 8080.
	GetServerPort() int

	// Get full server address in the format hostname:port.
	GetServerAddress() string
}
