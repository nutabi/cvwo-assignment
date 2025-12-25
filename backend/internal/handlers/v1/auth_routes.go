package v1

import "github.com/gin-gonic/gin"

func (h *handlers) registerAuthRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")

	authRoutes.POST("/login", h.authMiddleware.LoginHandler)
	authRoutes.POST("/logout", h.authMiddleware.LogoutHandler)
	authRoutes.POST("/refresh", h.authMiddleware.RefreshHandler)
}
