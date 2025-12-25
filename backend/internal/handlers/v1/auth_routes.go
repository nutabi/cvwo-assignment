package v1

import "github.com/gin-gonic/gin"

func (h *handlers) registerAuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", h.authMiddleware.LoginHandler)
	r.POST("/logout", h.authMiddleware.LogoutHandler)
	r.POST("/refresh", h.authMiddleware.RefreshHandler)
	r.POST("/register", h.handleUserRegistration)
}

func (h *handlers) handleUserRegistration(c *gin.Context) {
}
