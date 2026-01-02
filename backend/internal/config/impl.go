package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	serverHostname string
	serverPort     int
	databaseUrl    string
	jwtSecret      string
	debug          bool
	corsOrigins    []string
	logLevel       slog.Level
	logRoot        string
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

	// Parse CORS origins (comma-separated, optional)
	corsOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	var corsOrigins []string
	if corsOriginsStr != "" {
		for origin := range strings.SplitSeq(corsOriginsStr, ",") {
			if trimmed := strings.TrimSpace(origin); trimmed != "" {
				corsOrigins = append(corsOrigins, trimmed)
			}
		}
	}

	// Parse log level (optional, defaults to WARN)
	// In debug mode, always use DEBUG level
	logLevel := slog.LevelWarn
	isDebug := getRequired("DEBUG") == "true"
	if isDebug {
		logLevel = slog.LevelDebug
	} else {
		logLevelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))
		switch logLevelStr {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN", "WARNING":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		}
	}

	// Parse log root (optional, defaults to stdout)
	// Ignored in debug mode
	logRoot := os.Getenv("LOG_ROOT")

	cfg := config{
		serverHostname: getRequired("SERVER_HOSTNAME"),
		serverPort:     getRequiredInt("SERVER_PORT"),
		databaseUrl:    getRequired("DATABASE_URL"),
		jwtSecret:      getRequired("JWT_SECRET"),
		debug:          isDebug,
		corsOrigins:    corsOrigins,
		logLevel:       logLevel,
		logRoot:        logRoot,
	}

	if isBad {
		return nil, fmt.Errorf("missing or invalid environment variables")
	}

	// Log all config values at debug level
	slog.Debug("configuration loaded",
		"server_hostname", cfg.serverHostname,
		"server_port", cfg.serverPort,
		"database_url", cfg.databaseUrl,
		"jwt_secret", "***REDACTED***",
		"debug", cfg.debug,
		"cors_origins", cfg.corsOrigins,
		"log_level", cfg.logLevel,
		"log_root", cfg.logRoot,
	)

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

func (c *config) GetCORSOrigins() []string {
	return c.corsOrigins
}

func (c *config) GetLogLevel() slog.Level {
	return c.logLevel
}

func (c *config) GetLogRoot() string {
	return c.logRoot
}
