package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func handleServiceError(c *gin.Context, err error) {
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

func retrievePaginationParams(c *gin.Context) (limit, offset int) {
	var err error

	// Parse limit
	limitStr := c.Query("limit")
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Fallback to default on error
	}

	// Parse offset
	offsetStr := c.Query("offset")
	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Fallback to default on error
	}

	return limit, offset
}

func getBoolParam(c *gin.Context, paramName string, defaultValue bool) (value bool) {
	switch c.Query(paramName) {
	case "1", "true":
		value = true
	case "0", "false":
		value = false
	default:
		value = defaultValue
	}
	return value
}
