package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func handleServiceError(c *gin.Context, err error) {
	// TODO: Implement error handling logic
	respondWithError(c, http.StatusInternalServerError, err.Error())
}
