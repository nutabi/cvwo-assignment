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
	r.GET("/users/:user_id", handlePublicUserProfile(service))

	// Add topic public routes
	r.GET("/topics", handleListTopics(service))
	r.GET("/topics/:topic_id", handleGetTopicInfo(service))

	// Add post public routes
	r.GET("/posts", handleListPosts(service))
	r.GET("/posts/:post_id", handleGetPostInfo(service))

	// Add comment public routes
	r.GET("/posts/:post_id/comments", handleListComments(service))
	r.GET("/comments/:comment_id", handleGetCommentInfo(service))

	protected := r.Group("", authMiddleware.MiddlewareFunc())

	// Add user protected routes
	protected.GET("/users/me", handleCurrentUserProfile(service))
	protected.PATCH("/users/me", handleUpdateUserProfile(service))

	// Add topic protected routes
	protected.POST("/topics", handleCreateTopic(service))
	protected.PATCH("/topics/:topic_id", handleUpdateTopic(service))
	protected.DELETE("/topics/:topic_id", handleDeleteTopic(service))

	// Add post protected routes
	protected.POST("/topics/:topic_id/posts", handleCreatePost(service))
	protected.PATCH("/posts/:post_id", handleUpdatePost(service))
	protected.DELETE("/posts/:post_id", handleDeletePost(service))

	// Add comment protected routes
	protected.POST("/posts/:post_id/comments", handleCreateComment(service))
	protected.PATCH("/comments/:comment_id", handleUpdateComment(service))
	protected.DELETE("/comments/:comment_id", handleDeleteComment(service))
}
