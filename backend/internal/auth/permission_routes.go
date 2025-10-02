package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterPermissionRoutes registers permission management routes
func RegisterPermissionRoutes(router *gin.Engine, handlers *PermissionHandlers, middleware *AuthMiddleware, permissionMiddleware *PermissionMiddleware) {
	// Permission management routes (require admin permissions)
	permissionGroup := router.Group("/api/v1/permissions")
	permissionGroup.Use(middleware.Authenticate())
	permissionGroup.Use(permissionMiddleware.RequirePermission("manage", "permissions"))
	{
		// Permission checking
		permissionGroup.GET("/check", handlers.CheckPermission)

		// User permissions
		permissionGroup.GET("/users/:user_id", handlers.GetUserPermissions)
		permissionGroup.POST("/users", handlers.CreateUserPermission)
		permissionGroup.PUT("/users/:id", handlers.UpdateUserPermission)
		permissionGroup.DELETE("/users/:id", handlers.DeleteUserPermission)

		// Role permissions
		permissionGroup.GET("/roles/:role_id", handlers.GetRolePermissions)
		permissionGroup.POST("/roles", handlers.CreateRolePermission)
		permissionGroup.PUT("/roles/:id", handlers.UpdateRolePermission)
		permissionGroup.DELETE("/roles/:id", handlers.DeleteRolePermission)

		// Permission summary
		permissionGroup.GET("/summary/:user_id", handlers.GetPermissionSummary)

		// Available options
		permissionGroup.GET("/actions", handlers.GetAvailableActions)
		permissionGroup.GET("/resources", handlers.GetAvailableResources)
		permissionGroup.GET("/levels", handlers.GetAvailableLevels)

		// Maintenance
		permissionGroup.POST("/cleanup", handlers.CleanupExpiredPermissions)
		permissionGroup.POST("/initialize", handlers.InitializeDefaultPermissions)
	}

	// Current user permission routes (accessible to authenticated users)
	userPermissionGroup := router.Group("/api/v1/user-permissions")
	userPermissionGroup.Use(middleware.Authenticate())
	{
		// Get current user's permissions
		userPermissionGroup.GET("/", handlers.GetCurrentUserPermissions)

		// Check specific permission
		userPermissionGroup.GET("/check", handlers.CheckPermission)
	}
}
