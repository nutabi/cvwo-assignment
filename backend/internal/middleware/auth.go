package middleware

import (
	"log/slog"
	"net/http"
	"time"

	gin_jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"github.com/nutabi/cvwo-assignment/backend/internal/utility"
)

const (
	JwtRealm            = "cvwo-assignment"
	JwtSigningAlgorithm = "HS512"
	IdentityKey         = "sub"

	JwtIssuer   = "cvwo-assignment-backend"
	JWTAudience = "cvwo-assignment-frontend"

	CookieName = "jwt"
)

func NewAuthConfig(
	repo repository.Repository,
	isDebug bool,
	serverDomain,
	jwtSecret string) *gin_jwt.GinJWTMiddleware {
	return &gin_jwt.GinJWTMiddleware{
		Realm:            JwtRealm,
		SigningAlgorithm: JwtSigningAlgorithm,
		Key:              []byte(jwtSecret),
		Timeout:          5 * time.Minute,    // Token valid for 5 minutes
		MaxRefresh:       7 * 24 * time.Hour, // Refreshable for 7 days
		Authenticator: func(c *gin.Context) (any, error) {
			// Parse login credentials from request
			var loginVals struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBind(&loginVals); err != nil {
				slog.Debug("Failed to bind login values", "error", err)
				return nil, gin_jwt.ErrMissingLoginValues
			}

			// Fetch user from repository
			user, err := repo.GetUserByUsername(c.Request.Context(), loginVals.Username)

			// Validate user credentials
			if err == nil && utility.VerifyPassword(loginVals.Password, user.PHC) {
				return user, nil
			}
			return nil, gin_jwt.ErrFailedAuthentication
		},
		PayloadFunc: func(data any) jwt.MapClaims {
			if user, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					"sub": user.ID,
					"iss": JwtIssuer,
					"aud": JWTAudience,
					"nbf": time.Now().Unix(),
					"iat": time.Now().Unix(),
					"jti": uuid.NewString(),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityKey:    IdentityKey,
		TokenLookup:    "header:Authorization,cookie:jwt",
		SendCookie:     true,
		SecureCookie:   !isDebug,
		CookieHTTPOnly: true,
		CookieDomain:   serverDomain,
		CookieName:     CookieName,
		CookieSameSite: http.SameSiteDefaultMode,
	}
}
