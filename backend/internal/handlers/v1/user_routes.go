package v1

import "github.com/gin-gonic/gin"

func (h *handlers) registerUserRoutes(r *gin.RouterGroup) {
	r.GET("/me", h.handleCurrentUserProfile)
	r.PATCH("/me", h.handleUpdateUserProfile)
	r.GET("/:user_id", h.handlerPublicUserProfile)
}

func (h *handlers) handleCurrentUserProfile(c *gin.Context) {
}

func (h *handlers) handlerPublicUserProfile(c *gin.Context) {
}

func (h *handlers) handleUpdateUserProfile(c *gin.Context) {
}
