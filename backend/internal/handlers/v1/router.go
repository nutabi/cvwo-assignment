package v1

import (
	gin_jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
)

type handlers struct {
	repo           repository.Repository
	authMiddleware *gin_jwt.GinJWTMiddleware
}

func RegisterRoutes(
	r *gin.RouterGroup,
	repo repository.Repository,
	authMiddleware *gin_jwt.GinJWTMiddleware,
) {
	h := handlers{
		repo:           repo,
		authMiddleware: authMiddleware,
	}

	h.registerAuthRoutes(r.Group("/auth"))
	h.registerUserRoutes(r.Group("/users"))
}
