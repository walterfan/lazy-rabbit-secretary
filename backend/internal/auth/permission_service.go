package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// PermissionService handles permission management and checking
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService creates a new permission service
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{
		db: db,
	}
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (s *PermissionService) CheckPermission(userID, realmID, action, resource string, context map[string]interface{}) (*models.PermissionResult, error) {
	// First check direct user permissions
	userPermission, err := s.getUserPermission(userID, realmID, resource)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get user permission: %w", err)
	}

	if userPermission != nil && userPermission.IsActive && !userPermission.IsExpired() {
		if userPermission.HasAction(models.PermissionAction(action)) {
			actions, _ := userPermission.GetActions()
			return &models.PermissionResult{
				Allowed:   true,
				Level:     userPermission.Level,
				Actions:   actions,
				Source:    "user",
				ExpiresAt: userPermission.ExpiresAt,
				Context:   context,
			}, nil
		}
	}

	// Check role-based permissions
	rolePermissions, err := s.getUserRolePermissions(userID, realmID, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	for _, rolePermission := range rolePermissions {
		if rolePermission.IsActive && !rolePermission.IsExpired() {
			if rolePermission.HasAction(models.PermissionAction(action)) {
				actions, _ := rolePermission.GetActions()
				return &models.PermissionResult{
					Allowed:   true,
					Level:     rolePermission.Level,
					Actions:   actions,
					Source:    "role",
					ExpiresAt: rolePermission.ExpiresAt,
					Context:   context,
				}, nil
			}
		}
	}

	// Check if user has any permission for the resource (for read-only access)
	if userPermission != nil && userPermission.IsActive && !userPermission.IsExpired() {
		if userPermission.CanRead() && action == string(models.ActionRead) {
			actions, _ := userPermission.GetActions()
			return &models.PermissionResult{
				Allowed:   true,
				Level:     userPermission.Level,
				Actions:   actions,
				Source:    "user",
				ExpiresAt: userPermission.ExpiresAt,
				Context:   context,
			}, nil
		}
	}

	for _, rolePermission := range rolePermissions {
		if rolePermission.IsActive && !rolePermission.IsExpired() {
			if rolePermission.CanRead() && action == string(models.ActionRead) {
				actions, _ := rolePermission.GetActions()
				return &models.PermissionResult{
					Allowed:   true,
					Level:     rolePermission.Level,
					Actions:   actions,
					Source:    "role",
					ExpiresAt: rolePermission.ExpiresAt,
					Context:   context,
				}, nil
			}
		}
	}

	// No permission found
	return &models.PermissionResult{
		Allowed: false,
		Reason:  "No permission found for this action and resource",
		Context: context,
	}, nil
}

// GetUserPermissions returns all permissions for a user
func (s *PermissionService) GetUserPermissions(userID, realmID string) ([]*models.UserPermission, error) {
	var permissions []*models.UserPermission
	err := s.db.Where("user_id = ? AND realm_id = ? AND is_active = ?", userID, realmID, true).
		Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	return permissions, nil
}

// GetRolePermissions returns all permissions for a role
func (s *PermissionService) GetRolePermissions(roleID, realmID string) ([]*models.RolePermission, error) {
	var permissions []*models.RolePermission
	err := s.db.Where("role_id = ? AND realm_id = ? AND is_active = ?", roleID, realmID, true).
		Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	return permissions, nil
}

// CreateUserPermission creates a new user permission
func (s *PermissionService) CreateUserPermission(req *models.UserPermissionRequest, createdBy string) (*models.UserPermission, error) {
	permission := &models.UserPermission{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		RealmID:   req.RealmID,
		Resource:  req.Resource,
		Level:     req.Level,
		IsActive:  true,
		ExpiresAt: req.ExpiresAt,
		CreatedBy: createdBy,
	}

	if err := permission.SetActions(req.Actions); err != nil {
		return nil, fmt.Errorf("failed to set actions: %w", err)
	}

	if err := permission.SetConditions(req.Conditions); err != nil {
		return nil, fmt.Errorf("failed to set conditions: %w", err)
	}

	if err := s.db.Create(permission).Error; err != nil {
		return nil, fmt.Errorf("failed to create user permission: %w", err)
	}

	return permission, nil
}

// CreateRolePermission creates a new role permission
func (s *PermissionService) CreateRolePermission(req *models.RolePermissionRequest, createdBy string) (*models.RolePermission, error) {
	permission := &models.RolePermission{
		ID:        uuid.New().String(),
		RoleID:    req.RoleID,
		RealmID:   req.RealmID,
		Resource:  req.Resource,
		Level:     req.Level,
		IsActive:  true,
		ExpiresAt: req.ExpiresAt,
		CreatedBy: createdBy,
	}

	if err := permission.SetActions(req.Actions); err != nil {
		return nil, fmt.Errorf("failed to set actions: %w", err)
	}

	if err := permission.SetConditions(req.Conditions); err != nil {
		return nil, fmt.Errorf("failed to set conditions: %w", err)
	}

	if err := s.db.Create(permission).Error; err != nil {
		return nil, fmt.Errorf("failed to create role permission: %w", err)
	}

	return permission, nil
}

// UpdateUserPermission updates an existing user permission
func (s *PermissionService) UpdateUserPermission(permissionID string, req *models.UserPermissionRequest, updatedBy string) (*models.UserPermission, error) {
	var permission models.UserPermission
	if err := s.db.Where("id = ?", permissionID).First(&permission).Error; err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	permission.Resource = req.Resource
	permission.Level = req.Level
	permission.ExpiresAt = req.ExpiresAt
	permission.UpdatedBy = updatedBy

	if err := permission.SetActions(req.Actions); err != nil {
		return nil, fmt.Errorf("failed to set actions: %w", err)
	}

	if err := permission.SetConditions(req.Conditions); err != nil {
		return nil, fmt.Errorf("failed to set conditions: %w", err)
	}

	if err := s.db.Save(&permission).Error; err != nil {
		return nil, fmt.Errorf("failed to update user permission: %w", err)
	}

	return &permission, nil
}

// UpdateRolePermission updates an existing role permission
func (s *PermissionService) UpdateRolePermission(permissionID string, req *models.RolePermissionRequest, updatedBy string) (*models.RolePermission, error) {
	var permission models.RolePermission
	if err := s.db.Where("id = ?", permissionID).First(&permission).Error; err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	permission.Resource = req.Resource
	permission.Level = req.Level
	permission.ExpiresAt = req.ExpiresAt
	permission.UpdatedBy = updatedBy

	if err := permission.SetActions(req.Actions); err != nil {
		return nil, fmt.Errorf("failed to set actions: %w", err)
	}

	if err := permission.SetConditions(req.Conditions); err != nil {
		return nil, fmt.Errorf("failed to set conditions: %w", err)
	}

	if err := s.db.Save(&permission).Error; err != nil {
		return nil, fmt.Errorf("failed to update role permission: %w", err)
	}

	return &permission, nil
}

// DeleteUserPermission deletes a user permission
func (s *PermissionService) DeleteUserPermission(permissionID string) error {
	if err := s.db.Where("id = ?", permissionID).Delete(&models.UserPermission{}).Error; err != nil {
		return fmt.Errorf("failed to delete user permission: %w", err)
	}
	return nil
}

// DeleteRolePermission deletes a role permission
func (s *PermissionService) DeleteRolePermission(permissionID string) error {
	if err := s.db.Where("id = ?", permissionID).Delete(&models.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to delete role permission: %w", err)
	}
	return nil
}

// GetPermissionSummary returns a summary of all permissions for a user
func (s *PermissionService) GetPermissionSummary(userID, realmID string) (*models.PermissionSummary, error) {
	// Get user permissions
	userPermissions, err := s.GetUserPermissions(userID, realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Get role permissions
	rolePermissions, err := s.getUserRolePermissions(userID, realmID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	// Get user roles
	roles, err := s.getUserRoles(userID, realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Build resource access map
	resourceAccess := make(map[string]models.ResourceAccess)

	// Process user permissions
	for _, perm := range userPermissions {
		if perm.IsActive && !perm.IsExpired() {
			actions, _ := perm.GetActions()
			resourceAccess[perm.Resource] = models.ResourceAccess{
				Resource:  perm.Resource,
				Level:     perm.Level,
				Actions:   actions,
				Source:    "user",
				ExpiresAt: perm.ExpiresAt,
			}
		}
	}

	// Process role permissions (role permissions can override user permissions)
	for _, perm := range rolePermissions {
		if perm.IsActive && !perm.IsExpired() {
			actions, _ := perm.GetActions()
			resourceAccess[perm.Resource] = models.ResourceAccess{
				Resource:  perm.Resource,
				Level:     perm.Level,
				Actions:   actions,
				Source:    "role",
				ExpiresAt: perm.ExpiresAt,
			}
		}
	}

	return &models.PermissionSummary{
		UserID:    userID,
		RealmID:   realmID,
		Resources: resourceAccess,
		Roles:     roles,
	}, nil
}

// getUserPermission gets a specific user permission for a resource
func (s *PermissionService) getUserPermission(userID, realmID, resource string) (*models.UserPermission, error) {
	var permission models.UserPermission
	err := s.db.Where("user_id = ? AND realm_id = ? AND resource = ?", userID, realmID, resource).
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// getUserRolePermissions gets all role permissions for a user
func (s *PermissionService) getUserRolePermissions(userID, realmID, resource string) ([]*models.RolePermission, error) {
	query := `
		SELECT rp.* FROM role_permissions rp
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ? AND rp.realm_id = ? AND rp.is_active = ?
	`
	args := []interface{}{userID, realmID, true}

	if resource != "" {
		query += " AND rp.resource = ?"
		args = append(args, resource)
	}

	var permissions []*models.RolePermission
	err := s.db.Raw(query, args...).Scan(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// getUserRoles gets all roles for a user
func (s *PermissionService) getUserRoles(userID, realmID string) ([]string, error) {
	var roles []string
	err := s.db.Table("user_roles").
		Select("role_id").
		Where("user_id = ?", userID).
		Pluck("role_id", &roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// CleanupExpiredPermissions removes expired permissions
func (s *PermissionService) CleanupExpiredPermissions() error {
	now := time.Now()

	// Cleanup expired user permissions
	if err := s.db.Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Delete(&models.UserPermission{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired user permissions: %w", err)
	}

	// Cleanup expired role permissions
	if err := s.db.Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Delete(&models.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired role permissions: %w", err)
	}

	return nil
}

// InitializeDefaultPermissions creates default permissions for admin and user roles
func (s *PermissionService) InitializeDefaultPermissions() error {
	// Get default realm and roles
	var realm models.Realm
	if err := s.db.Where("name = ?", "default").First(&realm).Error; err != nil {
		return fmt.Errorf("default realm not found: %w", err)
	}

	var adminRole models.Role
	if err := s.db.Where("name = ? AND realm_id = ?", "admin", realm.ID).First(&adminRole).Error; err != nil {
		return fmt.Errorf("admin role not found: %w", err)
	}

	var userRole models.Role
	if err := s.db.Where("name = ? AND realm_id = ?", "user", realm.ID).First(&userRole).Error; err != nil {
		return fmt.Errorf("user role not found: %w", err)
	}

	// Create admin permissions (full access to all resources)
	adminResources := []string{
		string(models.ResourceUsers),
		string(models.ResourceRoles),
		string(models.ResourcePolicies),
		string(models.ResourceRealms),
		string(models.ResourcePosts),
		string(models.ResourcePages),
		string(models.ResourceComments),
		string(models.ResourceMedia),
		string(models.ResourceWiki),
		string(models.ResourceBooks),
		string(models.ResourceBookmarks),
		string(models.ResourceNews),
		string(models.ResourceTasks),
		string(models.ResourceReminders),
		string(models.ResourceDiagrams),
		string(models.ResourceImages),
		string(models.ResourceCommands),
		string(models.ResourceSettings),
		string(models.ResourceLogs),
		string(models.ResourceBackups),
		string(models.ResourceReports),
	}

	for _, resource := range adminResources {
		req := &models.RolePermissionRequest{
			RoleID:   adminRole.ID,
			RealmID:  realm.ID,
			Resource: resource,
			Actions:  []models.PermissionAction{models.ActionAll},
			Level:    models.LevelFullAccess,
		}

		if _, err := s.CreateRolePermission(req, "system"); err != nil {
			return fmt.Errorf("failed to create admin permission for %s: %w", resource, err)
		}
	}

	// Create user permissions (read-only access to most resources, read-write for personal resources)
	userReadOnlyResources := []string{
		string(models.ResourcePosts),
		string(models.ResourcePages),
		string(models.ResourceComments),
		string(models.ResourceMedia),
		string(models.ResourceWiki),
		string(models.ResourceBooks),
		string(models.ResourceBookmarks),
		string(models.ResourceNews),
		string(models.ResourceDiagrams),
		string(models.ResourceImages),
	}

	userReadWriteResources := []string{
		string(models.ResourceTasks),
		string(models.ResourceReminders),
	}

	for _, resource := range userReadOnlyResources {
		req := &models.RolePermissionRequest{
			RoleID:   userRole.ID,
			RealmID:  realm.ID,
			Resource: resource,
			Actions:  []models.PermissionAction{models.ActionRead},
			Level:    models.LevelReadOnly,
		}

		if _, err := s.CreateRolePermission(req, "system"); err != nil {
			return fmt.Errorf("failed to create user read-only permission for %s: %w", resource, err)
		}
	}

	for _, resource := range userReadWriteResources {
		req := &models.RolePermissionRequest{
			RoleID:   userRole.ID,
			RealmID:  realm.ID,
			Resource: resource,
			Actions:  []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
			Level:    models.LevelReadWrite,
		}

		if _, err := s.CreateRolePermission(req, "system"); err != nil {
			return fmt.Errorf("failed to create user read-write permission for %s: %w", resource, err)
		}
	}

	return nil
}
