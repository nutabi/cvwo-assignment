package app

import (
	"log/slog"

	gin_jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/config"
	v1 "github.com/nutabi/cvwo-assignment/backend/internal/handlers/v1"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/service/primary"
	"gorm.io/driver/sqlite"
)

type App struct {
	serverAddress  string
	service        service.Service
	authMiddleware *gin_jwt.GinJWTMiddleware
}

func Initialise(cfg config.Config) App {
	// Set Gin mode
	if cfg.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialise repository
	// Only supporting SQLite for now
	repo, err := repository.ConnectSQL(sqlite.Open(cfg.GetDSN()))
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	// Run migrations
	err = repo.Migrate()
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
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

	// Initialize the JWT middleware
	if err := auth.MiddlewareInit(); err != nil {
		slog.Error("Failed to initialize JWT middleware", "error", err)
		panic(err)
	}

	return App{
		serverAddress:  cfg.GetServerAddress(),
		service:        svc,
		authMiddleware: auth,
	}
}

func (a *App) Start() {
	// Initialise router
	r := gin.New()

	// Add logger and recovery middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add routes
	v1.RegisterRoutes(r.Group("/v1"), a.service, a.authMiddleware)

	// Start listening
	err := r.Run(a.serverAddress)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
