package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// Handle GET {ROOT}/topics
func handleListTopics(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters for pagination
		limit, offset := tryGetPagingParams(c)

		// Parse filter parameters
		userID := tryGetIDQuery(c, "user_id")

		// Delegate to service layer to get list of topics
		topics, err := svc.FetchTopics(
			c.Request.Context(),
			limit,
			offset,
			userID,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with list of topics
		c.JSON(http.StatusOK, topics)
	}
}

// Handle POST {ROOT}/topics
func handleCreateTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse request body
		var req struct {
			Title       string  `json:"title" form:"title" binding:"required"`
			Description *string `json:"description" form:"description"`
		}
		if !mustBindReqBody(c, &req) {
			return
		}

		// Delegate to service layer to create topic
		topic, err := svc.CreateTopic(
			c.Request.Context(),
			user.ID,
			req.Title,
			req.Description,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with created topic
		c.JSON(http.StatusCreated, topic)
	}
}

// Handle GET {ROOT}/topics/:topic_id
func handleGetTopicInfo(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse topic ID from path
		topicID, ok := mustGetIDParam(c, "topic_id")
		if !ok {
			return
		}

		// Delegate to service layer to get topic details
		topic, err := svc.FetchTopicByID(
			c.Request.Context(),
			topicID,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with topic details
		c.JSON(http.StatusOK, topic)
	}
}

// Handle PATCH {ROOT}/topics/:topic_id
func handleUpdateTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse topic ID from path
		topicID, ok := mustGetIDParam(c, "topic_id")
		if !ok {
			return
		}

		// Parse request body
		var req struct {
			Title       *string `json:"title" form:"title"`
			Description *string `json:"description" form:"description"`
		}
		if !mustBindReqBody(c, &req) {
			return
		}

		// Make sure at least one field is being updated
		if req.Title == nil && req.Description == nil {
			handleError(c, http.StatusUnprocessableEntity, "at least one field must be provided")
			return
		}

		// Delegate to service layer to update topic
		if err := svc.UpdateTopic(
			c.Request.Context(),
			topicID,
			user.ID,
			req.Title,
			req.Description,
		); err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with no content
		c.Status(http.StatusNoContent)
	}
}

// Handle DELETE {ROOT}/topics/:id
func handleDeleteTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse topic ID from path
		topicID, ok := mustGetIDParam(c, "topic_id")
		if !ok {
			return
		}

		// Delegate to service layer to delete topic
		if err := svc.DeleteTopic(
			c.Request.Context(),
			topicID,
			user.ID,
		); err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with no content
		c.Status(http.StatusNoContent)
	}
}
