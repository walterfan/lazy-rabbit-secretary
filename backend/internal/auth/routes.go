package auth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers authentication routes with the Gin router
func RegisterRoutes(router *gin.Engine, authHandlers *AuthHandlers, userHandlers *UserHandlers, roleHandlers *RoleHandlers, policyHandlers *PolicyHandlers, realmHandlers *RealmHandlers, middleware *AuthMiddleware) {
	// Public routes (no authentication required)
	public := router.Group("/api/v1/auth")
	{
		public.POST("/login", authHandlers.Login)
		public.POST("/register", authHandlers.Register)
		public.GET("/confirm", authHandlers.ConfirmEmail)
		public.POST("/refresh", authHandlers.RefreshToken)
		public.GET("/health", authHandlers.HealthCheck)
	}

	// Protected routes (authentication required)
	protected := router.Group("/api/v1/auth")
	protected.Use(middleware.Authenticate())
	{
		protected.GET("/profile", authHandlers.GetProfile)
		protected.GET("/permissions/check", authHandlers.CheckPermission)
		protected.POST("/logout", authHandlers.Logout)
	}

	// Admin routes (require admin role)
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.Authenticate())
	admin.Use(middleware.RequireRole("admin", "super_admin"))
	{
		// User management
		admin.GET("/users", userHandlers.GetUsers)          // List users
		admin.POST("/users", userHandlers.CreateUser)       // Create user
		admin.GET("/users/:id", userHandlers.GetUser)       // Get user
		admin.PUT("/users/:id", userHandlers.UpdateUser)    // Update user
		admin.DELETE("/users/:id", userHandlers.DeleteUser) // Delete user

		// Registration management
		admin.GET("/registrations", userHandlers.GetPendingRegistrations)      // List pending registrations
		admin.POST("/registrations/approve", userHandlers.ApproveRegistration) // Approve/deny registration
		admin.GET("/registrations/stats", userHandlers.GetRegistrationStats)   // Get registration statistics

		// Role management
		admin.GET("/roles", roleHandlers.GetRoles)          // List roles
		admin.POST("/roles", roleHandlers.CreateRole)       // Create role
		admin.GET("/roles/:id", roleHandlers.GetRole)       // Get role
		admin.PUT("/roles/:id", roleHandlers.UpdateRole)    // Update role
		admin.DELETE("/roles/:id", roleHandlers.DeleteRole) // Delete role

		// Policy management
		admin.GET("/policies", policyHandlers.GetPolicies)         // List policies
		admin.POST("/policies", policyHandlers.CreatePolicy)       // Create policy
		admin.GET("/policies/:id", policyHandlers.GetPolicy)       // Get policy
		admin.PUT("/policies/:id", policyHandlers.UpdatePolicy)    // Update policy
		admin.DELETE("/policies/:id", policyHandlers.DeletePolicy) // Delete policy

		// Realm management
		admin.GET("/realms", realmHandlers.GetRealms)          // List realms
		admin.POST("/realms", realmHandlers.CreateRealm)       // Create realm
		admin.GET("/realms/:id", realmHandlers.GetRealm)       // Get realm
		admin.PUT("/realms/:id", realmHandlers.UpdateRealm)    // Update realm
		admin.DELETE("/realms/:id", realmHandlers.DeleteRealm) // Delete realm
	}
}

// RegisterProtectedRoutes registers routes that require authentication
// This can be used by other parts of the application to add protected routes
func RegisterProtectedRoutes(router *gin.RouterGroup, middleware *AuthMiddleware) {
	// Example of how to add protected routes with specific permissions
	router.Use(middleware.Authenticate())

	// Routes that require specific permissions
	projects := router.Group("/projects")
	{
		projects.GET("", middleware.RequirePermission("read:projects", "project:*"))
		projects.POST("", middleware.RequirePermission("create:project", "project:*"))
		projects.GET("/:id", middleware.RequirePermission("read:project", "project:*"))
		projects.PUT("/:id", middleware.RequirePermission("update:project", "project:*"))
		projects.DELETE("/:id", middleware.RequirePermission("delete:project", "project:*"))
	}

	prompts := router.Group("/prompts")
	{
		prompts.GET("", middleware.RequirePermission("read:prompts", "prompt:*"))
		prompts.POST("", middleware.RequirePermission("create:prompt", "prompt:*"))
		prompts.GET("/:id", middleware.RequirePermission("read:prompt", "prompt:*"))
		prompts.PUT("/:id", middleware.RequirePermission("update:prompt", "prompt:*"))
		prompts.DELETE("/:id", middleware.RequirePermission("delete:prompt", "prompt:*"))
	}

	// Routes that require specific roles
	adminOnly := router.Group("/admin")
	adminOnly.Use(middleware.RequireRole("admin", "super_admin"))
	{
		adminOnly.GET("/stats", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin stats"})
		})
	}
}

// Example of how to use the middleware in your existing routes
func ExampleUsage(router *gin.Engine, middleware *AuthMiddleware) {
	// Public routes
	router.GET("/public", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Public endpoint"})
	})

	// Protected routes with authentication
	protected := router.Group("/protected")
	protected.Use(middleware.Authenticate())
	{
		protected.GET("/user-info", func(c *gin.Context) {
			userID, _ := GetCurrentUser(c)
			username, _ := GetCurrentUsername(c)
			c.JSON(200, gin.H{
				"user_id":  userID,
				"username": username,
			})
		})
	}

	// Routes with role-based access control
	admin := router.Group("/admin")
	admin.Use(middleware.Authenticate())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin dashboard"})
		})
	}

	// Routes with permission-based access control
	projects := router.Group("/projects")
	projects.Use(middleware.Authenticate())
	{
		projects.GET("", middleware.RequirePermission("read:projects", "project:*"))
		projects.POST("", middleware.RequirePermission("create:project", "project:*"))
		projects.GET("/:id", middleware.RequirePermission("read:project", "project:*"))
		projects.PUT("/:id", middleware.RequirePermission("update:project", "project:*"))
		projects.DELETE("/:id", middleware.RequirePermission("delete:project", "project:*"))
	}
}
