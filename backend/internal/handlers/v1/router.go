package v1

import (
	gin_jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	service service.Service,
	authMiddleware *gin_jwt.GinJWTMiddleware,
) {
	// Add authentication routes
	r.POST("/auth/login", authMiddleware.LoginHandler)
	r.POST("/auth/logout", authMiddleware.LogoutHandler)
	r.POST("/auth/refresh", authMiddleware.RefreshHandler)
	r.POST("/auth/register", handleUserRegistration(service))

	// Add public routes
	r.GET("/users/:id", handlerPublicUserProfile(service))

	// Add protected routes
	protected := r.Group("", authMiddleware.MiddlewareFunc())
	protected.GET("/users/me", handleCurrentUserProfile(service))
	protected.PUT("/users/me", handleUpdateUserProfile(service))
}
