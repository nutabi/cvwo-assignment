package app

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/config"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"gorm.io/driver/sqlite"
)

type App struct {
	repo           repository.Repository
	authMiddleware gin.HandlerFunc
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
	return App{repo: repo, authMiddleware: auth.MiddlewareFunc()}
}

func (a *App) Start() {
	// TODO: Initilise router

	// TODO: Add routes

	// TODO: Start listening
}
