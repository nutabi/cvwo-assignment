package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// Handle GET /topics
func handleListTopics(svc service.Service) func(c *gin.Context) {
	return nil
}

// Handle POST /topics
func handleCreateTopic(svc service.Service) func(c *gin.Context) {
	return nil
}

// Handle GET /topics/:id
func handleGetTopicDetails(svc service.Service) func(c *gin.Context) {
	return nil
}

// Handle PATCH /topics/:id
func handleUpdateTopic(svc service.Service) func(c *gin.Context) {
	return nil
}

// Handle DELETE /topics/:id
func handleDeleteTopic(svc service.Service) func(c *gin.Context) {
	return nil
}
