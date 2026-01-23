package app

import (
	"context"
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
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
	"gorm.io/driver/sqlite"
)

type App struct {
	serverAddress  string
	service        service.Service
	authMiddleware *gin_jwt.GinJWTMiddleware
	isDebug        bool
	corsOrigins    []string
}

func Initialise(cfg config.Config) App {
	// Set up structured logging - always to console with human-readable format
	handlerOpts := &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	}
	
	// Use custom pretty handler for readable, colorized console output
	handler := utility.NewPrettyHandler(os.Stdout, handlerOpts)

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

	slog.Info("server stopped")
}
