package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// @Summary      Create a new post
// @Description  Create a new post in a specific topic
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic_id path int true "Topic ID"
// @Param        request body CreatePostRequest true "Post creation details"
// @Success      201 {object} service.PostInfo
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      404 {object} ErrorResponse "Topic not found"
// @Failure      500 {object} ErrorResponse
// @Router       /topics/{topic_id}/posts [post]
func handleCreatePost(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse topic ID from URL parameters
		topicID, ok := mustGetIDParam(c, "topic_id")
		if !ok {
			return
		}

		// Parse request body
		var req CreatePostRequest
		if !mustBindReqBody(c, &req) {
			return
		}

		// Delegate to service layer to create post
		post, err := svc.CreatePost(c.Request.Context(), user.ID, topicID, req.Title, req.Content)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusCreated, post)
	}
}

// @Summary      List posts
// @Description  Get a paginated list of posts, optionally filtered by topic ID or user ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of items per page" default(20)
// @Param        offset query int false "Number of items to skip" default(0)
// @Param        topic_id query int false "Filter by topic ID"
// @Param        user_id query int false "Filter by user ID"
// @Success      200 {array} service.PostInfo
// @Failure      500 {object} ErrorResponse
// @Router       /posts [get]
func handleListPosts(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse pagination parameters
		limit, offset := tryGetPagingParams(c)

		// Parse filter parameters
		topicID := tryGetIDQuery(c, "topic_id")
		userID := tryGetIDQuery(c, "user_id")

		// Delegate to service layer to get list of posts
		posts, err := svc.FetchPosts(
			c.Request.Context(),
			limit,
			offset,
			topicID,
			userID,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with list of posts
		c.JSON(http.StatusOK, posts)
	}
}

// @Summary      Get post by ID
// @Description  Retrieve detailed information about a specific post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post_id path int true "Post ID"
// @Success      200 {object} service.PostInfo
// @Failure      400 {object} ErrorResponse "Invalid post ID"
// @Failure      404 {object} ErrorResponse "Post not found"
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{post_id} [get]
func handleGetPostInfo(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse post ID from path
		postID, ok := mustGetIDParam(c, "post_id")
		if !ok {
			return
		}

		// Delegate to service layer to get post info
		post, err := svc.FetchPostByID(
			c.Request.Context(),
			postID,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with post info
		c.JSON(http.StatusOK, post)
	}
}

// @Summary      Update a post
// @Description  Update a post's title or content (author only)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Param        request body UpdatePostRequest true "Post update details"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Forbidden - not the author"
// @Failure      404 {object} ErrorResponse "Post not found"
// @Failure      422 {object} ErrorResponse "Invalid input"
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{post_id} [patch]
func handleUpdatePost(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse post ID from path
		postID, ok := mustGetIDParam(c, "post_id")
		if !ok {
			return
		}

		// Parse request body
		var req UpdatePostRequest
		if !mustBindReqBody(c, &req) {
			return
		}

		// Make sure at least one field is being updated
		if req.Title == nil && req.Content == nil {
			handleError(c, ErrCodeInvalidInput)
			return
		}

		// Delegate to service layer to update post
		if err := svc.UpdatePost(
			c.Request.Context(),
			postID,
			user.ID,
			req.Title,
			req.Content,
		); err != nil {
			handleServiceError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary      Delete a post
// @Description  Delete a post (author only)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id path int true "Post ID"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Forbidden - not the author"
// @Failure      404 {object} ErrorResponse "Post not found"
// @Failure      500 {object} ErrorResponse
// @Router       /posts/{post_id} [delete]
func handleDeletePost(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse post ID from path
		postID, ok := mustGetIDParam(c, "post_id")
		if !ok {
			return
		}

		// Delegate to service layer to delete post
		if err := svc.DeletePost(
			c.Request.Context(),
			postID,
			user.ID,
		); err != nil {
			handleServiceError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
