package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// @Summary      Get user profile by ID
// @Description  Retrieve public profile information for a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id path int true "User ID"
// @Success      200 {object} service.UserProfile
// @Failure      400 {object} ErrorResponse "Invalid user ID"
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse
// @Router       /users/{user_id} [get]
func handlePublicUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse user ID from path
		userID, ok := mustGetIDParam(c, "user_id")
		if !ok {
			return
		}

		// Delegate to service layer to get user profile
		userProfile, err := svc.FetchUserByID(c.Request.Context(), userID)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with user profile
		c.JSON(http.StatusOK, userProfile)
	}
}

// @Summary      Get current user profile
// @Description  Retrieve the authenticated user's profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} service.UserProfile
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      500 {object} ErrorResponse
// @Router       /users/me [get]
func handleCurrentUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Delegate to service layer to get user profile
		userProfile, err := svc.FetchCurrentUser(c.Request.Context(), user)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with user profile
		c.JSON(http.StatusOK, userProfile)
	}
}

// @Summary      Update current user profile
// @Description  Update the authenticated user's avatar URL or bio
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateUserRequest true "Profile update details"
// @Success      204 "No Content"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      422 {object} ErrorResponse "Invalid input"
// @Failure      500 {object} ErrorResponse
// @Router       /users/me [patch]
func handleUpdateUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse request body
		var updateReq UpdateUserRequest
		if !mustBindReqBody(c, &updateReq) {
			return
		}

		// At least one field must be provided
		if updateReq.AvatarUrl == nil && updateReq.Bio == nil {
			handleError(c, http.StatusUnprocessableEntity, "at least one field must be provided")
			return
		}

		// Delegate to service layer to update user profile
		err := svc.UpdateCurrentUser(c.Request.Context(), user, updateReq.AvatarUrl, updateReq.Bio)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond
		c.Status(http.StatusNoContent)
	}
}
