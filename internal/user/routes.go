package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user routes with the given Gin router
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Create user route group
	users := router.Group("/users")
	{
		// Public routes (no authentication required)
		users.POST("", h.CreateUser)
		users.POST("/login", h.LoginUser)

		// Protected routes (authentication required)
		protected := users.Group("")
		// TODO: Add auth middleware here
		// protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("", h.ListUsers)
			protected.GET("/:id", h.GetUser)
			protected.PUT("/:id", h.UpdateUser)
			protected.DELETE("/:id", h.DeleteUser)
		}
	}
}
