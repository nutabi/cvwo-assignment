package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// @Summary      List topics
// @Description  Get a paginated list of topics, optionally filtered by user ID
// @Tags         topics
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of items per page" default(20)
// @Param        offset query int false "Number of items to skip" default(0)
// @Param        user_id query int false "Filter by user ID"
// @Success      200 {array} service.TopicInfo
// @Failure      500 {object} ErrorResponse
// @Router       /topics [get]
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

// @Summary      Create a new topic
// @Description  Create a new discussion topic
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTopicRequest true "Topic creation details"
// @Success      201 {object} service.TopicInfo
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      409 {object} ErrorResponse "Topic title already exists"
// @Failure      500 {object} ErrorResponse
// @Router       /topics [post]
func handleCreateTopic(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse request body
		var req CreateTopicRequest
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

// @Summary      Get topic by ID
// @Description  Retrieve detailed information about a specific topic
// @Tags         topics
// @Accept       json
// @Produce      json
// @Param        topic_id path int true "Topic ID"
// @Success      200 {object} service.TopicInfo
// @Failure      400 {object} ErrorResponse "Invalid topic ID"
// @Failure      404 {object} ErrorResponse "Topic not found"
// @Failure      500 {object} ErrorResponse
// @Router       /topics/{topic_id} [get]
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

// @Summary      Update a topic
// @Description  Update a topic's title or description (author only)
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic_id path int true "Topic ID"
// @Param        request body UpdateTopicRequest true "Topic update details"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Forbidden - not the author"
// @Failure      404 {object} ErrorResponse "Topic not found"
// @Failure      422 {object} ErrorResponse "Invalid input"
// @Failure      500 {object} ErrorResponse
// @Router       /topics/{topic_id} [patch]
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
		var req UpdateTopicRequest
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

// @Summary      Delete a topic
// @Description  Delete a topic (author only)
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic_id path int true "Topic ID"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Forbidden - not the author"
// @Failure      404 {object} ErrorResponse "Topic not found"
// @Failure      500 {object} ErrorResponse
// @Router       /topics/{topic_id} [delete]
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
