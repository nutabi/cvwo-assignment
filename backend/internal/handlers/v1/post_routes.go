package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// Handle POST {ROOT}/topics/:topic_id/posts
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
		var req struct {
			Title   string  `json:"title" form:"title" binding:"required"`
			Content *string `json:"content" form:"content"`
		}
		if !mustBindReqBody(c, &req) {
			return
		}

		// Delegate to service layer to create post
		post, err := svc.CreatePost(c.Request.Context(), user.ID, topicID, req.Title, *req.Content)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		c.JSON(http.StatusCreated, post)
	}
}

// Handle GET {ROOT}/posts
func handleListPosts(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse pagination parameters
		limit, offset := tryGetPagingParams(c)

		// Parse filter parameters
		postID := tryGetIDQuery(c, "post_id")
		userID := tryGetIDQuery(c, "user_id")

		// Parse preload comments
		withComments := tryGetBoolQuery(c, "with_comments", false)

		// Delegate to service layer to get list of posts
		posts, err := svc.FetchPosts(
			c.Request.Context(),
			limit,
			offset,
			postID,
			userID,
			withComments,
		)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with list of posts
		c.JSON(http.StatusOK, posts)
	}
}

// Handle GET {ROOT}/posts/:post_id
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

// Handle PATCH {ROOT}/posts/:post_id
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
		var req struct {
			Title   *string `json:"title" form:"title"`
			Content *string `json:"content" form:"content"`
		}
		if !mustBindReqBody(c, &req) {
			return
		}

		// Make sure at least one field is being updated
		if req.Title == nil && req.Content == nil {
			handleError(c, http.StatusBadRequest, "no fields to update")
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

// Handle DELETE {ROOT}/posts/:post_id
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
