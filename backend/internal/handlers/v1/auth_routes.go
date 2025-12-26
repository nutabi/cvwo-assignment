package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

func handleUserRegistration(svc service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Parse request body
		var reqBody struct {
			Username string `form:"username" json:"username" binding:"required"`
			Email    string `form:"email" json:"email" binding:"required"`
			Password string `form:"password" json:"password" binding:"required"`
		}
		if err := c.ShouldBind(&reqBody); err != nil {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid request body")
			return
		}

		// Validate input
		if !utility.ValidateUsername(reqBody.Username) {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid username")
		} else if !utility.ValidateEmail(reqBody.Email) {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid email")
		} else if !utility.ValidatePassword(reqBody.Password) {
			respondWithError(c, http.StatusUnprocessableEntity, "Invalid password")
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
