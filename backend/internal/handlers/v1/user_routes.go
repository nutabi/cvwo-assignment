package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

// Handle GET /{ROOT}/users/:id
func handlerPublicUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse user ID from path
		userId, ok := utility.TryParseID(c.Param("id"))
		if !ok {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid ID format")
			return
		}

		// Delegate to service layer to get user profile
		userProfile, err := svc.FetchUserByID(c.Request.Context(), userId)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with user profile
		c.JSON(http.StatusOK, userProfile)
	}
}

// Handle GET /{ROOT}/users/me
func handleCurrentUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Since this is an authenticated route, retrieve user from context
		user, exists := c.Get(middleware.UserIdentityKey)
		if !exists {
			// This should not happen as the middleware ensures authentication.
			// But handle it gracefully just in case.
			respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Gracefully type assert the user
		userObj, ok := user.(*model.User)
		if !ok {
			respondWithError(c, http.StatusInternalServerError, "Failed to retrieve user from context")
			return
		}

		// Delegate to service layer to get user profile
		userProfile, err := svc.FetchCurrentUser(c.Request.Context(), userObj)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with user profile
		c.JSON(http.StatusOK, userProfile)
	}
}

// Handle PATCH /{ROOT}/users/me
func handleUpdateUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Since this is an authenticated route, retrieve user from context
		user, exists := c.Get(middleware.UserIdentityKey)
		if !exists {
			// This should not happen as the middleware ensures authentication.
			// But handle it gracefully just in case.
			respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Gracefully type assert the user
		userObj, ok := user.(*model.User)
		if !ok {
			respondWithError(c, http.StatusInternalServerError, "Failed to retrieve user from context")
			return
		}

		// Parse request body
		var updateReq struct {
			AvatarUrl *string `form:"avatar_url" json:"avatar_url"`
			Bio       *string `form:"bio" json:"bio"`
		}
		if err := c.ShouldBind(&updateReq); err != nil {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid request body")
			return
		}

		// At least one field must be provided
		if updateReq.AvatarUrl == nil && updateReq.Bio == nil {
			respondWithError(c, http.StatusUnprocessableEntity, "At least one field (avatar_url or bio) must be provided")
			return
		}
		// Validate avatar URL if provided
		if updateReq.AvatarUrl != nil && !utility.ValidateAvatarUrl(*updateReq.AvatarUrl) {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid avatar URL format")
			return
		}
		// Validate bio if provided
		if updateReq.Bio != nil && !utility.ValidateBio(*updateReq.Bio) {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid bio format")
			return
		}

		// Delegate to service layer to update user profile
		err := svc.UpdateCurrentUser(c.Request.Context(), userObj, updateReq.AvatarUrl, updateReq.Bio)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond
		c.Status(http.StatusNoContent)
	}
}
