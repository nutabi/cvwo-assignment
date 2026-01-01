package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// Handle POST {ROOT}/posts/:post_id/comments
func handleCreateComment(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse post ID from URL parameters
		postID, ok := mustGetIDParam(c, "post_id")
		if !ok {
			return
		}

		// Parse request body
		var req struct {
			Content string `json:"content" form:"content" binding:"required"`
		}
		if !mustBindReqBody(c, &req) {
			return
		}

		// Delegate to service layer to create comment
		comment, err := svc.CreateComment(c.Request.Context(), user.ID, postID, req.Content)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusCreated, comment)
	}
}

// Handle GET {ROOT}/comments
func handleListComments(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse pagination parameters
		limit, offset := tryGetPagingParams(c)

		// Parse filter parameters
		postID := tryGetIDQuery(c, "post_id")
		userID := tryGetIDQuery(c, "user_id")

		// Delegate to service layer to get list of comments
		comments, err := svc.FetchComments(
			c.Request.Context(),
			limit,
			offset,
			postID,
			userID,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusOK, comments)
	}
}

// Handle GET {ROOT}/comments/:comment_id
func handleGetCommentInfo(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse comment ID from path
		commentID, ok := mustGetIDParam(c, "comment_id")
		if !ok {
			return
		}

		// Delegate to service layer to get comment info
		comment, err := svc.FetchCommentByID(c.Request.Context(), commentID)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}

func handleUpdateComment(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse comment ID from path
		commentID, ok := mustGetIDParam(c, "comment_id")
		if !ok {
			return
		}

		// Parse request body
		var req struct {
			Content *string `json:"content" form:"content"`
		}
		if !mustBindReqBody(c, &req) {
			return
		}

		// Make sure at least one field is being updated
		if req.Content == nil {
			handleError(c, http.StatusUnprocessableEntity, "content must be provided")
			return
		}

		// Delegate to service layer to update comment
		err := svc.UpdateComment(
			c.Request.Context(),
			user.ID,
			commentID,
			*req.Content,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func handleDeleteComment(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse comment ID from path
		commentID, ok := mustGetIDParam(c, "comment_id")
		if !ok {
			return
		}

		// Delegate to service layer to delete comment
		err := svc.DeleteComment(
			c.Request.Context(),
			commentID,
			user.ID,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
