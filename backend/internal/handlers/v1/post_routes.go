package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// Handle POST {ROOT}/topics/{topic_id}/posts
func handleCreatePost(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// Handle GET {ROOT}/topics/{topic_id}/posts
func handleListPosts(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// Handle GET {ROOT}/posts/{post_id}
func handleGetPostInfo(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// Handle PATCH {ROOT}/posts/{post_id}
func handleUpdatePost(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// Handle DELETE {ROOT}/posts/{post_id}
func handleDeletePost(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
