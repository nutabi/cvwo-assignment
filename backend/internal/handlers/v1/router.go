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

	// Add user public routes
	r.GET("/users/:id", handlerPublicUserProfile(service))

	// Add topic public routes
	r.GET("/topics", handleListTopics(service))
	r.GET("/topics/:id", handleGetTopicDetails(service))

	protected := r.Group("", authMiddleware.MiddlewareFunc())

	// Add user protected routes
	protected.GET("/users/me", handleCurrentUserProfile(service))
	protected.PUT("/users/me", handleUpdateUserProfile(service))

	// Add topic protected routes
	protected.POST("/topics", handleCreateTopic(service))
	protected.PATCH("/topics/:id", handleUpdateTopic(service))
	protected.DELETE("/topics/:id", handleDeleteTopic(service))
}
