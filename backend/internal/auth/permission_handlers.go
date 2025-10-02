package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// PermissionHandlers provides HTTP handlers for permission management
type PermissionHandlers struct {
	permissionService *PermissionService
}

// NewPermissionHandlers creates a new permission handlers
func NewPermissionHandlers(permissionService *PermissionService) *PermissionHandlers {
	return &PermissionHandlers{
		permissionService: permissionService,
	}
}

// CheckPermission checks if the current user has permission for an action
func (h *PermissionHandlers) CheckPermission(c *gin.Context) {
	userID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	realmID, exists := GetCurrentRealm(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	action := c.Query("action")
	if action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action parameter required"})
		return
	}

	resource := c.Query("resource")
	if resource == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource parameter required"})
		return
	}

	// Build context from query parameters
	context := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if key != "action" && key != "resource" && len(values) > 0 {
			context[key] = values[0]
		}
	}

	// Add user context
	context["user:id"] = userID
	context["user:realm_id"] = realmID
	context["user:username"] = c.GetString("username")
	context["user:email"] = c.GetString("email")
	context["user:roles"] = c.GetStringSlice("roles")

	result, err := h.permissionService.CheckPermission(userID, realmID, action, resource, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserPermissions returns all permissions for a user
func (h *PermissionHandlers) GetUserPermissions(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	realmID, exists := GetCurrentRealm(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	permissions, err := h.permissionService.GetUserPermissions(userID, realmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// GetRolePermissions returns all permissions for a role
func (h *PermissionHandlers) GetRolePermissions(c *gin.Context) {
	roleID := c.Param("role_id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	realmID, exists := GetCurrentRealm(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	permissions, err := h.permissionService.GetRolePermissions(roleID, realmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// CreateUserPermission creates a new user permission
func (h *PermissionHandlers) CreateUserPermission(c *gin.Context) {
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	permission, err := h.permissionService.CreateUserPermission(&req, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user permission", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User permission created successfully",
		"permission": permission,
	})
}

// CreateRolePermission creates a new role permission
func (h *PermissionHandlers) CreateRolePermission(c *gin.Context) {
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.RolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	permission, err := h.permissionService.CreateRolePermission(&req, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role permission", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Role permission created successfully",
		"permission": permission,
	})
}

// UpdateUserPermission updates an existing user permission
func (h *PermissionHandlers) UpdateUserPermission(c *gin.Context) {
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	permissionID := c.Param("id")
	if permissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission ID is required"})
		return
	}

	var req models.UserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	permission, err := h.permissionService.UpdateUserPermission(permissionID, &req, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user permission", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User permission updated successfully",
		"permission": permission,
	})
}

// UpdateRolePermission updates an existing role permission
func (h *PermissionHandlers) UpdateRolePermission(c *gin.Context) {
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	permissionID := c.Param("id")
	if permissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission ID is required"})
		return
	}

	var req models.RolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	permission, err := h.permissionService.UpdateRolePermission(permissionID, &req, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role permission", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Role permission updated successfully",
		"permission": permission,
	})
}

// DeleteUserPermission deletes a user permission
func (h *PermissionHandlers) DeleteUserPermission(c *gin.Context) {
	permissionID := c.Param("id")
	if permissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission ID is required"})
		return
	}

	if err := h.permissionService.DeleteUserPermission(permissionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user permission", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User permission deleted successfully"})
}

// DeleteRolePermission deletes a role permission
func (h *PermissionHandlers) DeleteRolePermission(c *gin.Context) {
	permissionID := c.Param("id")
	if permissionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission ID is required"})
		return
	}

	if err := h.permissionService.DeleteRolePermission(permissionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role permission", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role permission deleted successfully"})
}

// GetPermissionSummary returns a summary of all permissions for a user
func (h *PermissionHandlers) GetPermissionSummary(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	realmID, exists := GetCurrentRealm(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summary, err := h.permissionService.GetPermissionSummary(userID, realmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permission summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetCurrentUserPermissions returns permissions for the current user
func (h *PermissionHandlers) GetCurrentUserPermissions(c *gin.Context) {
	userID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	realmID, exists := GetCurrentRealm(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summary, err := h.permissionService.GetPermissionSummary(userID, realmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permission summary", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetAvailableActions returns all available permission actions
func (h *PermissionHandlers) GetAvailableActions(c *gin.Context) {
	actions := []models.PermissionAction{
		models.ActionCreate,
		models.ActionRead,
		models.ActionUpdate,
		models.ActionDelete,
		models.ActionManage,
		models.ActionApprove,
		models.ActionReject,
		models.ActionSuspend,
		models.ActionActivate,
		models.ActionExecute,
		models.ActionExport,
		models.ActionImport,
		models.ActionBackup,
		models.ActionRestore,
		models.ActionAll,
	}

	c.JSON(http.StatusOK, gin.H{"actions": actions})
}

// GetAvailableResources returns all available permission resources
func (h *PermissionHandlers) GetAvailableResources(c *gin.Context) {
	resources := []models.PermissionResource{
		models.ResourceUsers,
		models.ResourceRoles,
		models.ResourcePolicies,
		models.ResourceRealms,
		models.ResourcePosts,
		models.ResourcePages,
		models.ResourceComments,
		models.ResourceMedia,
		models.ResourceWiki,
		models.ResourceBooks,
		models.ResourceBookmarks,
		models.ResourceNews,
		models.ResourceTasks,
		models.ResourceReminders,
		models.ResourceDiagrams,
		models.ResourceImages,
		models.ResourceCommands,
		models.ResourceSettings,
		models.ResourceLogs,
		models.ResourceBackups,
		models.ResourceReports,
		models.ResourceAll,
	}

	c.JSON(http.StatusOK, gin.H{"resources": resources})
}

// GetAvailableLevels returns all available permission levels
func (h *PermissionHandlers) GetAvailableLevels(c *gin.Context) {
	levels := []models.PermissionLevel{
		models.LevelReadOnly,
		models.LevelReadWrite,
		models.LevelFullAccess,
		models.LevelCustom,
	}

	c.JSON(http.StatusOK, gin.H{"levels": levels})
}

// CleanupExpiredPermissions removes expired permissions
func (h *PermissionHandlers) CleanupExpiredPermissions(c *gin.Context) {
	if err := h.permissionService.CleanupExpiredPermissions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup expired permissions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expired permissions cleaned up successfully"})
}

// InitializeDefaultPermissions creates default permissions for admin and user roles
func (h *PermissionHandlers) InitializeDefaultPermissions(c *gin.Context) {
	if err := h.permissionService.InitializeDefaultPermissions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize default permissions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default permissions initialized successfully"})
}
