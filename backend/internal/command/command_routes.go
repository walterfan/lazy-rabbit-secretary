package command

import (
	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
)

// RegisterRoutes registers command routes with the router
func RegisterRoutes(router *gin.Engine, service *CommandService, middleware *auth.AuthMiddleware) {
	group := router.Group("/api/v1/commands")
	group.Use(middleware.Authenticate())
	{
		group.GET("", service.GetCommands)
		group.POST("", service.ExecuteCommand)
		group.GET("/ws", service.HandleWebSocket)
	}
}
