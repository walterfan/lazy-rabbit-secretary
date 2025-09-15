package book

import (
	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-reminder/internal/auth"
)

// RegisterRoutes registers authentication routes with the Gin router
func RegisterRoutes(router *gin.Engine, handlers *auth.AuthHandlers, middleware *auth.AuthMiddleware) {


	// Protected routes (authentication required)
	protected := router.Group("/api/v1/books")
	protected.Use(middleware.Authenticate())
	{
		//protected.GET("/", handlers.GetProfile)
		//protected.GET("/permissions/check", handlers.CheckPermission)
		//protected.POST("/logout", handlers.Logout)
	}
}