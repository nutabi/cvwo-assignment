package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/model"
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
		} else if errors.Is(err, service.ErrTopicNotFound) {
			respondWithError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrUnauthorized) {
			respondWithError(c, http.StatusUnauthorized, err.Error())
		} else if errors.Is(err, service.ErrNoUpdateFields) {
			respondWithError(c, http.StatusUnprocessableEntity, err.Error())
		} else {
			// Generic internal server error for unhandled cases
			// as well as the following errors:
			// - ErrCryptoError
			// - ErrDatabaseError
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

func handleUnauthorised(c *gin.Context) {
	respondWithError(c, http.StatusUnauthorized, "unauthorised")
}

func handleInternalError(c *gin.Context) {
	respondWithError(c, http.StatusInternalServerError, "internal server error")
}

func retrieveUser(c *gin.Context) *model.User {
	user, exists := c.Get(middleware.UserIdentityKey)
	if !exists {
		handleUnauthorised(c)
		return nil
	}
	userObj, ok := user.(*model.User)
	if !ok {
		handleInternalError(c)
		return nil
	}
	return userObj
}
