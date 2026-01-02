package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nutabi/cvwo-assignment/backend/internal/middleware"
	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

// Swagger model definitions

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// UpdateUserRequest represents user profile update request
type UpdateUserRequest struct {
	AvatarUrl *string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Bio       *string `json:"bio" example:"Software developer"`
}

// CreateTopicRequest represents topic creation request
type CreateTopicRequest struct {
	Title       string  `json:"title" binding:"required" example:"General Discussion"`
	Description *string `json:"description" example:"A place for general conversations"`
}

// UpdateTopicRequest represents topic update request
type UpdateTopicRequest struct {
	Title       *string `json:"title" example:"Updated Topic Title"`
	Description *string `json:"description" example:"Updated description"`
}

// CreatePostRequest represents post creation request
type CreatePostRequest struct {
	Title   string  `json:"title" binding:"required" example:"My First Post"`
	Content *string `json:"content" example:"This is the content of my post"`
}

// UpdatePostRequest represents post update request
type UpdatePostRequest struct {
	Title   *string `json:"title" example:"Updated Post Title"`
	Content *string `json:"content" example:"Updated post content"`
}

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required" example:"This is a comment"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Content *string `json:"content" example:"Updated comment content"`
}

// Helper function to respond with an error message and HTTP status code.
func handleError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// Responds with a generic 401 Unauthorized error.
func handleUnauthorised(c *gin.Context) {
	handleError(c, http.StatusUnauthorized, "unauthorised")
}

// Responds with a generic 500 Internal Server Error.
func handleInternalError(c *gin.Context) {
	handleError(c, http.StatusInternalServerError, "internal server error")
}

// Responds with appropriate HTTP status codes based on the provided service error.
func handleServiceError(c *gin.Context, err error) {
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			handleError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrUserNotFound) {
			handleError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrUsernameTaken) {
			handleError(c, http.StatusConflict, err.Error())
		} else if errors.Is(err, service.ErrEmailInUse) {
			handleError(c, http.StatusConflict, err.Error())
		} else if errors.Is(err, service.ErrTopicNotFound) {
			handleError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrTopicTitleTaken) {
			handleError(c, http.StatusConflict, err.Error())
		} else if errors.Is(err, service.ErrPostNotFound) {
			handleError(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, service.ErrForbidden) {
			handleError(c, http.StatusForbidden, err.Error())
		} else if errors.Is(err, service.ErrNoUpdateFields) {
			handleError(c, http.StatusUnprocessableEntity, err.Error())
		} else {
			// Generic internal server error for unhandled cases
			// as well as the following errors:
			// - ErrCryptoError
			// - ErrDatabaseError
			slog.Error("internal server error in handler",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"error", err,
			)
			handleError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
}

// Tries to retrieve an unsigned integer ID query parameter from the request context.
// If the parameter is not present or invalid, returns 0.
func tryGetIDQuery(c *gin.Context, paramName string) uint {
	paramStr, exists := c.GetQuery(paramName)
	if !exists {
		return 0
	}
	parsedValue, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		return 0
	}
	return uint(parsedValue)
}

// Tries to retrieve pagination parameters (limit and offset) from the request context.
// If parameters are not present or invalid, returns default values (limit=20, offset=0).
func tryGetPagingParams(c *gin.Context) (limit, offset int) {
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

// Tries to retrieve a required string parameter from the URL path.
// If the parameter is missing, responds with a 400 Bad Request error.
func mustGetStrParam(c *gin.Context, paramName string) string {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		handleError(c, http.StatusBadRequest, "missing required parameter: "+paramName)
	}
	return paramValue
}

// Tries to retrieve a required integer parameter from the URL path.
// If the parameter is missing or invalid, responds with a 400 Bad Request error.
func mustGetIntParam(c *gin.Context, paramName string) (int, bool) {
	paramValueStr := mustGetStrParam(c, paramName)
	paramValue, err := strconv.Atoi(paramValueStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "invalid integer parameter: "+paramName)
		return 0, false
	}
	return paramValue, true
}

func mustGetIDParam(c *gin.Context, paramName string) (uint, bool) {
	id, ok := mustGetIntParam(c, paramName)
	if ok && id <= 0 {
		handleError(c, http.StatusBadRequest, "invalid ID parameter: "+paramName)
		return 0, false
	}
	return uint(id), true
}

// Tries to retrieve the authenticated user from the context.
// If the user is not found or of incorrect type, responds with an error
// and returns nil.
func mustRetrieveUser(c *gin.Context) *model.User {
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

// Tries to bind the request body to the given object.
// If binding fails, responds with a 400 Bad Request error.
// Returns true if binding was successful, false otherwise.
func mustBindReqBody(c *gin.Context, obj any) bool {
	if err := c.ShouldBind(obj); err != nil {
		handleError(c, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}
