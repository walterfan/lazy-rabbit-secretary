package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// PermissionMiddleware provides middleware for permission checking
type PermissionMiddleware struct {
	permissionService *PermissionService
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(permissionService *PermissionService) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionService: permissionService,
	}
}

// RequirePermission checks if the user has permission for a specific action and resource
func (m *PermissionMiddleware) RequirePermission(action, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		realmID, exists := c.Get("realm_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Build context from request
		context := m.buildContext(c)

		// Check permission
		result, err := m.permissionService.CheckPermission(
			userID.(string),
			realmID.(string),
			action,
			resource,
			context,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !result.Allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Insufficient permissions",
				"reason": result.Reason,
			})
			c.Abort()
			return
		}

		// Add permission result to context for use in handlers
		c.Set("permission_result", result)
		c.Next()
	}
}

// RequireReadPermission checks if the user has read permission for a resource
func (m *PermissionMiddleware) RequireReadPermission(resource string) gin.HandlerFunc {
	return m.RequirePermission(string(models.ActionRead), resource)
}

// RequireWritePermission checks if the user has write permission for a resource
func (m *PermissionMiddleware) RequireWritePermission(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		realmID, exists := c.Get("realm_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Build context from request
		context := m.buildContext(c)

		// Check for any write action (create, update, delete)
		writeActions := []string{
			string(models.ActionCreate),
			string(models.ActionUpdate),
			string(models.ActionDelete),
		}

		var hasWritePermission bool
		var lastResult *models.PermissionResult

		for _, action := range writeActions {
			result, err := m.permissionService.CheckPermission(
				userID.(string),
				realmID.(string),
				action,
				resource,
				context,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
				c.Abort()
				return
			}

			lastResult = result
			if result.Allowed {
				hasWritePermission = true
				break
			}
		}

		if !hasWritePermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Insufficient write permissions",
				"reason": lastResult.Reason,
			})
			c.Abort()
			return
		}

		// Add permission result to context
		c.Set("permission_result", lastResult)
		c.Next()
	}
}

// RequireManagePermission checks if the user has management permission for a resource
func (m *PermissionMiddleware) RequireManagePermission(resource string) gin.HandlerFunc {
	return m.RequirePermission(string(models.ActionManage), resource)
}

// RequireAnyPermission checks if the user has any of the specified permissions
func (m *PermissionMiddleware) RequireAnyPermission(permissions []PermissionRequirement) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		realmID, exists := c.Get("realm_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Build context from request
		context := m.buildContext(c)

		// Check each permission requirement
		for _, perm := range permissions {
			result, err := m.permissionService.CheckPermission(
				userID.(string),
				realmID.(string),
				perm.Action,
				perm.Resource,
				context,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
				c.Abort()
				return
			}

			if result.Allowed {
				// User has at least one required permission
				c.Set("permission_result", result)
				c.Next()
				return
			}
		}

		// User doesn't have any of the required permissions
		c.JSON(http.StatusForbidden, gin.H{
			"error":                "Insufficient permissions",
			"required_permissions": permissions,
		})
		c.Abort()
	}
}

// RequireAllPermissions checks if the user has all of the specified permissions
func (m *PermissionMiddleware) RequireAllPermissions(permissions []PermissionRequirement) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		realmID, exists := c.Get("realm_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Build context from request
		context := m.buildContext(c)

		// Check each permission requirement
		for _, perm := range permissions {
			result, err := m.permissionService.CheckPermission(
				userID.(string),
				realmID.(string),
				perm.Action,
				perm.Resource,
				context,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
				c.Abort()
				return
			}

			if !result.Allowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error":              "Insufficient permissions",
					"missing_permission": perm,
					"reason":             result.Reason,
				})
				c.Abort()
				return
			}
		}

		// User has all required permissions
		c.Next()
	}
}

// OptionalPermission checks permission but allows access even without it (useful for read-only features)
func (m *PermissionMiddleware) OptionalPermission(action, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// No user, continue without permission check
			c.Next()
			return
		}

		realmID, exists := c.Get("realm_id")
		if !exists {
			// No realm, continue without permission check
			c.Next()
			return
		}

		// Build context from request
		context := m.buildContext(c)

		// Check permission
		result, err := m.permissionService.CheckPermission(
			userID.(string),
			realmID.(string),
			action,
			resource,
			context,
		)

		if err != nil {
			// Log error but continue
			c.Next()
			return
		}

		// Add permission result to context (even if not allowed)
		c.Set("permission_result", result)
		c.Next()
	}
}

// CheckReadOnlyAccess checks if the user only has read-only access to a resource
func (m *PermissionMiddleware) CheckReadOnlyAccess(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// No user, assume read-only
			c.Set("read_only", true)
			c.Next()
			return
		}

		realmID, exists := c.Get("realm_id")
		if !exists {
			// No realm, assume read-only
			c.Set("read_only", true)
			c.Next()
			return
		}

		// Build context from request
		context := m.buildContext(c)

		// Check for write permissions
		writeActions := []string{
			string(models.ActionCreate),
			string(models.ActionUpdate),
			string(models.ActionDelete),
		}

		hasWritePermission := false
		for _, action := range writeActions {
			result, err := m.permissionService.CheckPermission(
				userID.(string),
				realmID.(string),
				action,
				resource,
				context,
			)

			if err == nil && result.Allowed {
				hasWritePermission = true
				break
			}
		}

		// Set read-only flag in context
		c.Set("read_only", !hasWritePermission)
		c.Next()
	}
}

// buildContext builds a context map from the gin context
func (m *PermissionMiddleware) buildContext(c *gin.Context) map[string]interface{} {
	context := make(map[string]interface{})

	// Add path parameters
	for _, param := range c.Params {
		context[param.Key] = param.Value
	}

	// Add query parameters
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			context[key] = values[0]
		}
	}

	// Add user context
	if userID, exists := c.Get("user_id"); exists {
		context["user:id"] = userID
	}
	if realmID, exists := c.Get("realm_id"); exists {
		context["user:realm_id"] = realmID
	}
	if username, exists := c.Get("username"); exists {
		context["user:username"] = username
	}
	if email, exists := c.Get("email"); exists {
		context["user:email"] = email
	}
	if roles, exists := c.Get("roles"); exists {
		context["user:roles"] = roles
	}

	// Add request method and path
	context["request:method"] = c.Request.Method
	context["request:path"] = c.Request.URL.Path

	return context
}

// PermissionRequirement represents a permission requirement
type PermissionRequirement struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

// GetPermissionResult extracts the permission result from context
func GetPermissionResult(c *gin.Context) (*models.PermissionResult, bool) {
	result, exists := c.Get("permission_result")
	if !exists {
		return nil, false
	}
	permResult, ok := result.(*models.PermissionResult)
	return permResult, ok
}

// IsReadOnly checks if the current request is read-only
func IsReadOnly(c *gin.Context) bool {
	readOnly, exists := c.Get("read_only")
	if !exists {
		return false
	}
	isReadOnly, ok := readOnly.(bool)
	return ok && isReadOnly
}

// GetResourceFromPath extracts the resource name from the request path
func GetResourceFromPath(path string) string {
	// Remove leading slash and split by '/'
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	// Look for API version pattern (e.g., /api/v1/users -> users)
	for i, part := range parts {
		if part == "api" && i+2 < len(parts) {
			return parts[i+2] // Return the resource after /api/v1/
		}
	}

	// If no API pattern found, return the first part
	if len(parts) > 0 {
		return parts[0]
	}

	return "unknown"
}

// GetActionFromMethod maps HTTP methods to permission actions
func GetActionFromMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return string(models.ActionRead)
	case "POST":
		return string(models.ActionCreate)
	case "PUT", "PATCH":
		return string(models.ActionUpdate)
	case "DELETE":
		return string(models.ActionDelete)
	default:
		return string(models.ActionRead)
	}
}
