package app

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gin_jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/config"
	v1 "github.com/nutabi/cvwo-assignment/backend/internal/handlers/v1"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository/sql"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/service/primary"
	"gorm.io/driver/sqlite"
)

type App struct {
	serverAddress  string
	service        service.Service
	authMiddleware *gin_jwt.GinJWTMiddleware
	isDebug        bool
	corsOrigins    []string
	logFile        *os.File // Track log file for cleanup
}

func Initialise(cfg config.Config) App {
	// Set up structured logging
	var logWriter io.Writer = os.Stdout
	var logFile *os.File

	// In debug mode, always log to console; otherwise use configured destination
	if !cfg.IsDebug() && cfg.GetLogRoot() != "" {
		file, err := os.OpenFile(cfg.GetLogRoot(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			slog.Warn("Failed to open log file, falling back to stdout", "path", cfg.GetLogRoot(), "error", err)
		} else {
			logWriter = file
			logFile = file
		}
	}

	// Create handler with configured log level
	handlerOpts := &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	}

	var handler slog.Handler
	if cfg.IsDebug() {
		// Use text handler for readable console output in debug mode
		handler = slog.NewTextHandler(logWriter, handlerOpts)
	} else {
		// Use JSON handler for structured logging in production
		handler = slog.NewJSONHandler(logWriter, handlerOpts)
	}

	// Set as default logger
	slog.SetDefault(slog.New(handler))

	// Set Gin to release mode to suppress route printing
	// Application logging is controlled separately via slog
	gin.SetMode(gin.ReleaseMode)

	// Initialise repository
	// Only supporting SQLite for now
	repo, err := sql.Connect(sqlite.Open(cfg.GetDSN()))
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	// Run migrations
	err = repo.Migrate()
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		panic(err)
	}

	// Initialise service layer
	svc := primary.NewService(repo)

	// Initialise middleware
	auth := middleware.NewAuthConfig(
		repo,
		cfg.IsDebug(),
		cfg.GetServerHostname(),
		cfg.GetJWTSecretKey(),
	)

	// Initialise the JWT middleware
	if err := auth.MiddlewareInit(); err != nil {
		slog.Error("Failed to initialise JWT middleware", "error", err)
		panic(err)
	}

	return App{
		serverAddress:  cfg.GetServerAddress(),
		service:        svc,
		authMiddleware: auth,
		isDebug:        cfg.IsDebug(),
		corsOrigins:    cfg.GetCORSOrigins(),
		logFile:        logFile,
	}
}

func (a *App) Start() {
	// Initialise router
	r := gin.New()

	// Add structured request logger middleware (logs at debug level)
	r.Use(middleware.NewRequestLogger())

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add CORS middleware
	r.Use(middleware.NewCORSConfig(a.isDebug, a.corsOrigins))

	// Add general rate limiting middleware
	r.Use(middleware.NewGeneralRateLimiter())

	// Add routes with authentication rate limiter
	authRateLimiter := middleware.NewAuthRateLimiter()
	v1.RegisterRoutes(r.Group("/api/v1"), a.service, a.authMiddleware, authRateLimiter)

	// Create HTTP server for graceful shutdown
	srv := &http.Server{
		Addr:    a.serverAddress,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("server starting", "address", a.serverAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("server shutting down")

	// Give outstanding requests a deadline to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	// Close log file if opened
	if a.logFile != nil {
		if err := a.logFile.Close(); err != nil {
			slog.Error("failed to close log file", "error", err)
		}
	}

	slog.Info("server stopped")
}
