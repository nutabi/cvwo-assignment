package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

// Handle GET /topics
func handleListTopics(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters for pagination
		limit, offset := retrievePaginationParams(c)

		// Parse filter parameters
		userID, ok := utility.TryParseID(c.DefaultQuery("user_id", "0"))
		if !ok {
			respondWithError(c, http.StatusUnprocessableEntity, "invalid user_id format")
			return
		}

		// Parse preload posts
		withPosts := getBoolParam(c, "with_posts", false)

		// Delegate to service layer to get list of topics
		topics, err := svc.FetchTopics(
			c.Request.Context(),
			limit,
			offset,
			userID,
			withPosts,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with list of topics
		c.JSON(http.StatusOK, topics)
	}
}

// Handle POST /topics
func handleCreateTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user, exists := c.Get(middleware.UserIdentityKey)
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Gracefully type assert the user
		userObj, ok := user.(*model.User)
		if !ok {
			respondWithError(c, http.StatusInternalServerError, "failed to retrieve user from context")
			return
		}

		// Parse request body
		var req struct {
			Title       string  `json:"title" form:"title" binding:"required"`
			Description *string `json:"description" form:"description"`
		}
		if err := c.ShouldBind(&req); err != nil {
			respondWithError(c, http.StatusUnprocessableEntity, "invalid request body")
			return
		}

		// Delegate to service layer to create topic
		topic, err := svc.CreateTopic(
			c.Request.Context(),
			userObj,
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

// Handle GET /topics/:id
func handleGetTopicInfo(svc service.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse topic ID from path
		topicId, ok := utility.TryParseID(ctx.Param("id"))
		if !ok {
			respondWithError(ctx, http.StatusUnprocessableEntity, "invalid topic ID")
			return
		}

		// Delegate to service layer to get topic details
		topic, err := svc.FetchTopicByID(
			ctx.Request.Context(),
			topicId,
		)
		if err != nil {
			handleServiceError(ctx, err)
			return
		}

		// Respond with topic details
		ctx.JSON(http.StatusOK, topic)
	}
}

// Handle PATCH /topics/:id
func handleUpdateTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user, exists := c.Get(middleware.UserIdentityKey)
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Gracefully type assert the user
		userObj, ok := user.(*model.User)
		if !ok {
			respondWithError(c, http.StatusInternalServerError, "failed to retrieve user from context")
			return
		}

		// Parse topic ID from path
		topicId, ok := utility.TryParseID(c.Param("id"))
		if !ok {
			respondWithError(c, http.StatusUnprocessableEntity, "invalid topic ID")
			return
		}

		// Parse request body
		var req struct {
			Title       *string `json:"title" form:"title"`
			Description *string `json:"description" form:"description"`
		}
		if err := c.ShouldBind(&req); err != nil {
			respondWithError(c, http.StatusUnprocessableEntity, "invalid request body")
			return
		}

		// Delegate to service layer to update topic
		if err := svc.UpdateTopic(
			c.Request.Context(),
			topicId,
			userObj.ID,
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

// Handle DELETE /topics/:id
func handleDeleteTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user, exists := c.Get(middleware.UserIdentityKey)
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Gracefully type assert the user
		userObj, ok := user.(*model.User)
		if !ok {
			respondWithError(c, http.StatusInternalServerError, "failed to retrieve user from context")
			return
		}

		// Parse topic ID from path
		topicId, ok := utility.TryParseID(c.Param("id"))
		if !ok {
			respondWithError(c, http.StatusUnprocessableEntity, "invalid topic ID")
			return
		}

		// Delegate to service layer to delete topic
		if err := svc.DeleteTopic(
			c.Request.Context(),
			topicId,
			userObj.ID,
		); err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with no content
		c.Status(http.StatusNoContent)
	}
}
