package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

// Handle GET {ROOT}/users/:user_id
func handlerPublicUserProfile(svc service.Service) gin.HandlerFunc {
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

// Handle GET {ROOT}/users/me
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

// Handle PATCH {ROOT}/users/me
func handleUpdateUserProfile(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve authenticated user from context
		user := mustRetrieveUser(c)
		if user == nil {
			return
		}

		// Parse request body
		var updateReq struct {
			AvatarUrl *string `form:"avatar_url" json:"avatar_url"`
			Bio       *string `form:"bio" json:"bio"`
		}
		if !mustBindReqBody(c, &updateReq) {
			return
		}

		// At least one field must be provided
		if updateReq.AvatarUrl == nil && updateReq.Bio == nil {
			handleError(c, http.StatusUnprocessableEntity, "At least one field (avatar_url or bio) must be provided")
			return
		}
		// Validate avatar URL if provided
		if updateReq.AvatarUrl != nil && !utility.ValidateAvatarUrl(*updateReq.AvatarUrl) {
			handleError(c, http.StatusUnprocessableEntity, "Invalid avatar URL format")
			return
		}
		// Validate bio if provided
		if updateReq.Bio != nil && !utility.ValidateBio(*updateReq.Bio) {
			handleError(c, http.StatusUnprocessableEntity, "Invalid bio format")
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
