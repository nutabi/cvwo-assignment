package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	v1 "github.com/nutabi/cvwo-assignment/backend/internal/handlers/v1"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository/sql"
	"github.com/nutabi/cvwo-assignment/backend/internal/service/primary"
	"gorm.io/driver/sqlite"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	// Initialise repository, service, middleware, etc. as needed for tests
	repo, err := sql.Connect(sqlite.Open("file::memory:"))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations to create tables
	err = repo.Migrate()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	svc := primary.NewService(repo)
	auth := middleware.NewAuthConfig(repo, true, "localhost", "testsecret")
	err = auth.MiddlewareInit()
	if err != nil {
		t.Fatalf("Failed to initialise auth middleware: %v", err)
	}

	// Initialise Gin router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Register routes
	v1.RegisterRoutes(r.Group("/v1"), svc, auth)
	return r
}
