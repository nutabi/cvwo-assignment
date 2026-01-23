package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// @Summary      Create a new comment
// @Description  Create a new comment on a specific post
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Param        request body CreateCommentRequest true "Comment creation details"
// @Success      201 {object} service.CommentInfo
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      404 {object} ErrorResponse "Post not found"
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{post_id}/comments [post]
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
		var req CreateCommentRequest
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

// @Summary      List comments
// @Description  Get a paginated list of comments, optionally filtered by post ID or user ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of items per page" default(20)
// @Param        offset query int false "Number of items to skip" default(0)
// @Param        post_id query int false "Filter by post ID"
// @Param        user_id query int false "Filter by user ID"
// @Success      200 {array} service.CommentInfo
// @Failure      500 {object} ErrorResponse
// @Router       /comments [get]
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

// @Summary      Get comment by ID
// @Description  Retrieve detailed information about a specific comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        comment_id path int true "Comment ID"
// @Success      200 {object} service.CommentInfo
// @Failure      400 {object} ErrorResponse "Invalid comment ID"
// @Failure      404 {object} ErrorResponse "Comment not found"
// @Failure      500 {object} ErrorResponse
// @Router       /comments/{comment_id} [get]
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

// @Summary      Update a comment
// @Description  Update a comment's content (author only)
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id path int true "Comment ID"
// @Param        request body UpdateCommentRequest true "Comment update details"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Forbidden - not the author"
// @Failure      404 {object} ErrorResponse "Comment not found"
// @Failure      422 {object} ErrorResponse "Invalid input"
// @Failure      500 {object} ErrorResponse
// @Router       /comments/{comment_id} [patch]
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
		var req UpdateCommentRequest
		if !mustBindReqBody(c, &req) {
			return
		}

		// Make sure at least one field is being updated
		if req.Content == nil {
			handleError(c, ErrCodeInvalidInput)
			return
		}

		// Delegate to service layer to update comment
		err := svc.UpdateComment(
			c.Request.Context(),
			commentID,
			user.ID,
			*req.Content,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary      Delete a comment
// @Description  Delete a comment (author only)
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id path int true "Comment ID"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Forbidden - not the author"
// @Failure      404 {object} ErrorResponse "Comment not found"
// @Failure      500 {object} ErrorResponse
// @Router       /comments/{comment_id} [delete]
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
