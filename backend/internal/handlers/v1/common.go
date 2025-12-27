package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func handleServiceError(c *gin.Context, err error) {
	// TODO: Implement error handling logic
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			respondWithError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrUsernameTaken) {
			respondWithError(c, http.StatusConflict, err.Error())
		} else if errors.Is(err, service.ErrEmailInUse) {
			respondWithError(c, http.StatusConflict, err.Error())
		} else {
			respondWithError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
}
