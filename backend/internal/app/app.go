package app

import (
	"log/slog"

	gin_jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/config"
	v1 "github.com/nutabi/cvwo-assignment/backend/internal/handlers/v1"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"gorm.io/driver/sqlite"
)

type App struct {
	serverAddress  string
	repo           repository.Repository
	authMiddleware *gin_jwt.GinJWTMiddleware
}

func Initialise(cfg config.Config) App {
	// Initialise repository
	// Only supporting SQLite for now
	repo, err := repository.ConnectSQL(sqlite.Open(cfg.GetDSN()))
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	// Initialise middleware
	auth := middleware.NewAuthConfig(
		repo,
		cfg.IsDebug(),
		cfg.GetServerHostname(),
		cfg.GetJWTSecretKey(),
	)
	return App{repo: repo, authMiddleware: auth}
}

func (a *App) Start() {
	// Initialise router
	r := gin.New()

	// Add routes
	v1.RegisterRoutes(r.Group("/v1"), a.repo, a.authMiddleware)

	// Start listening
	err := r.Run(a.serverAddress)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
