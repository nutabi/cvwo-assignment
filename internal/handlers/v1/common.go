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
	ErrorCode    string `json:"error_code" example:"INVALID_INPUT"`
	ErrorMessage string `json:"error_message" example:"The provided input is invalid"`
}

// Error codes
const (
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeUnprocessableEntity = "UNPROCESSABLE_ENTITY"
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeMissingLoginValues  = "MISSING_LOGIN_VALUES"
	ErrCodeTokenExpired        = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid        = "TOKEN_INVALID"
	ErrCodeCommentNotFound     = "COMMENT_NOT_FOUND"
	ErrCodeUserNotFound        = "USER_NOT_FOUND"
	ErrCodeUsernameTaken       = "USERNAME_TAKEN"
	ErrCodeEmailInUse          = "EMAIL_IN_USE"
	ErrCodeTopicNotFound       = "TOPIC_NOT_FOUND"
	ErrCodeTopicTitleTaken     = "TOPIC_TITLE_TAKEN"
	ErrCodePostNotFound        = "POST_NOT_FOUND"
	ErrCodeNoUpdateFields      = "NO_UPDATE_FIELDS"
)

// errorInfo holds HTTP status code and message for an error code
type errorInfo struct {
	httpStatus int
	message    string
}

// errorCodeMap maps error codes to their HTTP status and messages
var errorCodeMap = map[string]errorInfo{
	ErrCodeUnauthorized:        {http.StatusUnauthorized, "unauthorised"},
	ErrCodeForbidden:           {http.StatusForbidden, "forbidden"},
	ErrCodeInternalServer:      {http.StatusInternalServerError, "internal server error"},
	ErrCodeBadRequest:          {http.StatusBadRequest, "bad request"},
	ErrCodeNotFound:            {http.StatusNotFound, "not found"},
	ErrCodeConflict:            {http.StatusConflict, "conflict"},
	ErrCodeUnprocessableEntity: {http.StatusUnprocessableEntity, "unprocessable entity"},
	ErrCodeInvalidInput:        {http.StatusBadRequest, "invalid input"},
	ErrCodeInvalidCredentials:  {http.StatusUnauthorized, "invalid username or password"},
	ErrCodeMissingLoginValues:  {http.StatusBadRequest, "missing username or password"},
	ErrCodeTokenExpired:        {http.StatusUnauthorized, "token has expired"},
	ErrCodeTokenInvalid:        {http.StatusUnauthorized, "invalid token"},
	ErrCodeCommentNotFound:     {http.StatusNotFound, "comment not found"},
	ErrCodeUserNotFound:        {http.StatusNotFound, "user not found"},
	ErrCodeUsernameTaken:       {http.StatusConflict, "username taken"},
	ErrCodeEmailInUse:          {http.StatusConflict, "email in use"},
	ErrCodeTopicNotFound:       {http.StatusNotFound, "topic not found"},
	ErrCodeTopicTitleTaken:     {http.StatusConflict, "topic title already exists"},
	ErrCodePostNotFound:        {http.StatusNotFound, "post not found"},
	ErrCodeNoUpdateFields:      {http.StatusUnprocessableEntity, "no fields to update"},
}

// LoginResponse represents user login response with JWT token
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in" example:"300"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
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

// Helper function to respond with an error based on error code.
// The HTTP status code and message are determined automatically from the error code.
func handleError(c *gin.Context, errorCode string) {
	info, ok := errorCodeMap[errorCode]
	if !ok {
		// Fallback to internal server error if error code is unknown
		info = errorCodeMap[ErrCodeInternalServer]
		errorCode = ErrCodeInternalServer
	}

	c.AbortWithStatusJSON(info.httpStatus, ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: info.message,
	})
}

// Responds with a generic 401 Unauthorized error.
func handleUnauthorised(c *gin.Context) {
	handleError(c, ErrCodeUnauthorized)
}

// Responds with a generic 500 Internal Server Error.
func handleInternalError(c *gin.Context) {
	handleError(c, ErrCodeInternalServer)
}

// Responds with appropriate HTTP status codes based on the provided service error.
func handleServiceError(c *gin.Context, err error) {
	if err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			handleError(c, ErrCodeCommentNotFound)
		} else if errors.Is(err, service.ErrUserNotFound) {
			handleError(c, ErrCodeUserNotFound)
		} else if errors.Is(err, service.ErrUsernameTaken) {
			handleError(c, ErrCodeUsernameTaken)
		} else if errors.Is(err, service.ErrEmailInUse) {
			handleError(c, ErrCodeEmailInUse)
		} else if errors.Is(err, service.ErrTopicNotFound) {
			handleError(c, ErrCodeTopicNotFound)
		} else if errors.Is(err, service.ErrTopicTitleTaken) {
			handleError(c, ErrCodeTopicTitleTaken)
		} else if errors.Is(err, service.ErrPostNotFound) {
			handleError(c, ErrCodePostNotFound)
		} else if errors.Is(err, service.ErrForbidden) {
			handleError(c, ErrCodeForbidden)
		} else if errors.Is(err, service.ErrNoUpdateFields) {
			handleError(c, ErrCodeNoUpdateFields)
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
			handleError(c, ErrCodeInternalServer)
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
		handleError(c, ErrCodeBadRequest)
	}
	return paramValue
}

// Tries to retrieve a required integer parameter from the URL path.
// If the parameter is missing or invalid, responds with a 400 Bad Request error.
func mustGetIntParam(c *gin.Context, paramName string) (int, bool) {
	paramValueStr := mustGetStrParam(c, paramName)
	paramValue, err := strconv.Atoi(paramValueStr)
	if err != nil {
		handleError(c, ErrCodeBadRequest)
		return 0, false
	}
	return paramValue, true
}

func mustGetIDParam(c *gin.Context, paramName string) (uint, bool) {
	id, ok := mustGetIntParam(c, paramName)
	if ok && id <= 0 {
		handleError(c, ErrCodeBadRequest)
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
		handleError(c, ErrCodeInvalidInput)
		return false
	}
	return true
}
