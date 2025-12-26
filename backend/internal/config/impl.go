package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	serverHostname string
	serverPort     int
	databaseUrl    string
	jwtSecret      string
	debug          bool
}

// Return the application configuration.
//
// This function loads environment variables from an .env file if present and constructs a Config
// object with the appropriate settings. It can be used in both development and production environments,
// since .env files do not override existing environment variables.
func LoadConfig() (Config, error) {
	// Load .env file if present
	// Log any error but continue since .env is optional
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found or error loading .env file", "error", err)
	}

	// Helper function to get required environment variable
	var isBad bool
	getRequired := func(key string) string {
		value := os.Getenv(key)
		if value == "" {
			slog.Error("Missing required environment variable", "key", key)
			isBad = true
		}
		return value
	}
	getRequiredInt := func(key string) int {
		value := os.Getenv(key)
		if result, err := strconv.Atoi(os.Getenv(key)); err != nil {
			slog.Error("Invalid integer for required environment variable", "key", key, "value", value)
			isBad = true
			return 0
		} else {
			return result
		}
	}

	cfg := config{
		serverHostname: getRequired("SERVER_HOSTNAME"),
		serverPort:     getRequiredInt("SERVER_PORT"),
		databaseUrl:    getRequired("DATABASE_URL"),
		jwtSecret:      getRequired("JWT_SECRET"),
		debug:          getRequired("DEBUG") == "true",
	}

	if isBad {
		return nil, fmt.Errorf("missing or invalid environment variables")
	}

	return &cfg, nil
}

func (c *config) GetServerHostname() string {
	return c.serverHostname
}

func (c *config) GetServerPort() int {
	return c.serverPort
}

func (c *config) GetServerAddress() string {
	return c.serverHostname + ":" + fmt.Sprintf("%d", c.serverPort)
}

func (c *config) GetDSN() string {
	return c.databaseUrl
}

func (c *config) IsDebug() bool {
	return c.debug
}

func (c *config) GetJWTSecretKey() string {
	return c.jwtSecret
}
