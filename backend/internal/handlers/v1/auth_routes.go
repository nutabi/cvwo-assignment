package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "User login details"
// @Success      200 {object} LoginResponse
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      500 {object} ErrorResponse
// @Router       /auth/login [post]
func handleLogin(h gin.HandlerFunc) gin.HandlerFunc {
	// Wrap login handler to add documentation
	return h
}

// @Summary      User logout
// @Description  Invalidate user JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]string "Logout successful message"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      500 {object} ErrorResponse
// @Router       /auth/logout [post]
func handleLogout(h gin.HandlerFunc) gin.HandlerFunc {
	// Wrap logout handler to add documentation
	return h
}

// @Summary      Refresh JWT token
// @Description  Refresh user JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} LoginResponse
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      500 {object} ErrorResponse
// @Router       /auth/refresh [post]
func handleTokenRefresh(h gin.HandlerFunc) gin.HandlerFunc {
	// Wrap token refresh handler to add documentation
	return h
}

// @Summary      Register a new user
// @Description  Create a new user account with username, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "User registration details"
// @Success      201 {object} service.UserProfile
// @Failure      422 {object} ErrorResponse "Invalid input"
// @Failure      409 {object} ErrorResponse "Username or email already exists"
// @Failure      500 {object} ErrorResponse
// @Router       /auth/register [post]
func handleUserRegistration(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var reqBody RegisterRequest
		if !mustBindReqBody(c, &reqBody) {
			return
		}

		// Validate input
		if !utility.ValidateUsername(reqBody.Username) {
			handleError(c, http.StatusUnprocessableEntity, "invalid or missing username")
			return
		} else if !utility.ValidateEmail(reqBody.Email) {
			handleError(c, http.StatusUnprocessableEntity, "invalid or missing email")
			return
		} else if !utility.ValidatePassword(reqBody.Password) {
			handleError(c, http.StatusUnprocessableEntity, "invalid or missing password")
			return
		}

		// Delegate to service layer for registration
		userProfile, err := svc.RegisterUser(c.Request.Context(), reqBody.Username, reqBody.Email, reqBody.Password)
		if err != nil {
			handleServiceError(c, err)
			return
		}

		// Respond with the created user profile
		c.JSON(http.StatusCreated, userProfile)
	}
}
