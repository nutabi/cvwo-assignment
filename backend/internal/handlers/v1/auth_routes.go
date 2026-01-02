package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

func handleUserRegistration(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var reqBody struct {
			Username string `form:"username" json:"username"`
			Email    string `form:"email" json:"email"`
			Password string `form:"password" json:"password"`
		}
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
