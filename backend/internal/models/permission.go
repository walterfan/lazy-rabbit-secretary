package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// PermissionAction represents the type of action that can be performed
type PermissionAction string

const (
	// CRUD Actions
	ActionCreate PermissionAction = "create"
	ActionRead   PermissionAction = "read"
	ActionUpdate PermissionAction = "update"
	ActionDelete PermissionAction = "delete"

	// Administrative Actions
	ActionManage   PermissionAction = "manage" // Full CRUD + admin operations
	ActionApprove  PermissionAction = "approve"
	ActionReject   PermissionAction = "reject"
	ActionSuspend  PermissionAction = "suspend"
	ActionActivate PermissionAction = "activate"

	// Special Actions
	ActionExecute PermissionAction = "execute" // For commands, scripts, etc.
	ActionExport  PermissionAction = "export"
	ActionImport  PermissionAction = "import"
	ActionBackup  PermissionAction = "backup"
	ActionRestore PermissionAction = "restore"

	// Wildcard
	ActionAll PermissionAction = "*"
)

// PermissionResource represents the resource being accessed
type PermissionResource string

const (
	// Core Resources
	ResourceUsers    PermissionResource = "users"
	ResourceRoles    PermissionResource = "roles"
	ResourcePolicies PermissionResource = "policies"
	ResourceRealms   PermissionResource = "realms"

	// Content Resources
	ResourcePosts    PermissionResource = "posts"
	ResourcePages    PermissionResource = "pages"
	ResourceComments PermissionResource = "comments"
	ResourceMedia    PermissionResource = "media"

	// Knowledge Resources
	ResourceWiki      PermissionResource = "wiki"
	ResourceBooks     PermissionResource = "books"
	ResourceBookmarks PermissionResource = "bookmarks"
	ResourceNews      PermissionResource = "news"

	// Tools Resources
	ResourceTasks     PermissionResource = "tasks"
	ResourceReminders PermissionResource = "reminders"
	ResourceDiagrams  PermissionResource = "diagrams"
	ResourceImages    PermissionResource = "images"
	ResourceCommands  PermissionResource = "commands"

	// System Resources
	ResourceSettings PermissionResource = "settings"
	ResourceLogs     PermissionResource = "logs"
	ResourceBackups  PermissionResource = "backups"
	ResourceReports  PermissionResource = "reports"

	// Wildcard
	ResourceAll PermissionResource = "*"
)

// PermissionLevel represents the granularity of permission
type PermissionLevel string

const (
	LevelReadOnly   PermissionLevel = "readonly"  // Can only read/view
	LevelReadWrite  PermissionLevel = "readwrite" // Can read and modify
	LevelFullAccess PermissionLevel = "full"      // Full access including admin operations
	LevelCustom     PermissionLevel = "custom"    // Custom permissions defined in statements
)

// UserPermission represents direct user permissions
type UserPermission struct {
	ID         string          `json:"id" gorm:"primaryKey;type:text"`
	UserID     string          `json:"user_id" gorm:"not null;type:text;index"`
	RealmID    string          `json:"realm_id" gorm:"not null;type:text;index"`
	Resource   string          `json:"resource" gorm:"not null;type:text;index"`
	Actions    string          `json:"actions" gorm:"not null;type:text"` // JSON array of actions
	Level      PermissionLevel `json:"level" gorm:"not null;type:text;default:'readonly'"`
	Conditions string          `json:"conditions" gorm:"type:text"` // JSON string for conditions
	IsActive   bool            `json:"is_active" gorm:"default:true"`
	ExpiresAt  *time.Time      `json:"expires_at" gorm:"index"`
	CreatedBy  string          `json:"created_by" gorm:"type:text"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy  string          `json:"updated_by" gorm:"type:text"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

// RolePermission represents role-based permissions
type RolePermission struct {
	ID         string          `json:"id" gorm:"primaryKey;type:text"`
	RoleID     string          `json:"role_id" gorm:"not null;type:text;index"`
	RealmID    string          `json:"realm_id" gorm:"not null;type:text;index"`
	Resource   string          `json:"resource" gorm:"not null;type:text;index"`
	Actions    string          `json:"actions" gorm:"not null;type:text"` // JSON array of actions
	Level      PermissionLevel `json:"level" gorm:"not null;type:text;default:'readonly'"`
	Conditions string          `json:"conditions" gorm:"type:text"` // JSON string for conditions
	IsActive   bool            `json:"is_active" gorm:"default:true"`
	ExpiresAt  *time.Time      `json:"expires_at" gorm:"index"`
	CreatedBy  string          `json:"created_by" gorm:"type:text"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy  string          `json:"updated_by" gorm:"type:text"`
	UpdatedAt  time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

// PermissionCheckRequest represents a permission check request
type PermissionCheckRequest struct {
	UserID   string                 `json:"user_id"`
	RealmID  string                 `json:"realm_id"`
	Action   string                 `json:"action"`
	Resource string                 `json:"resource"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// PermissionResult represents the result of a permission check
type PermissionResult struct {
	Allowed   bool                   `json:"allowed"`
	Level     PermissionLevel        `json:"level"`
	Actions   []PermissionAction     `json:"actions"`
	Source    string                 `json:"source"` // "user", "role", "policy"
	Reason    string                 `json:"reason,omitempty"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// UserPermissionRequest represents a request to create/update user permissions
type UserPermissionRequest struct {
	UserID     string                 `json:"user_id" binding:"required"`
	RealmID    string                 `json:"realm_id" binding:"required"`
	Resource   string                 `json:"resource" binding:"required"`
	Actions    []PermissionAction     `json:"actions" binding:"required"`
	Level      PermissionLevel        `json:"level" binding:"required"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
}

// RolePermissionRequest represents a request to create/update role permissions
type RolePermissionRequest struct {
	RoleID     string                 `json:"role_id" binding:"required"`
	RealmID    string                 `json:"realm_id" binding:"required"`
	Resource   string                 `json:"resource" binding:"required"`
	Actions    []PermissionAction     `json:"actions" binding:"required"`
	Level      PermissionLevel        `json:"level" binding:"required"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
}

// PermissionSummary represents a summary of user permissions
type PermissionSummary struct {
	UserID    string                    `json:"user_id"`
	RealmID   string                    `json:"realm_id"`
	Resources map[string]ResourceAccess `json:"resources"`
	Roles     []string                  `json:"roles"`
	ExpiresAt *time.Time                `json:"expires_at,omitempty"`
}

// ResourceAccess represents access level for a specific resource
type ResourceAccess struct {
	Resource  string             `json:"resource"`
	Level     PermissionLevel    `json:"level"`
	Actions   []PermissionAction `json:"actions"`
	Source    string             `json:"source"` // "user", "role", "policy"
	ExpiresAt *time.Time         `json:"expires_at,omitempty"`
}

// GetActions returns the actions as a slice
func (up *UserPermission) GetActions() ([]PermissionAction, error) {
	var actions []PermissionAction
	if err := json.Unmarshal([]byte(up.Actions), &actions); err != nil {
		return nil, err
	}
	return actions, nil
}

// SetActions sets the actions from a slice
func (up *UserPermission) SetActions(actions []PermissionAction) error {
	data, err := json.Marshal(actions)
	if err != nil {
		return err
	}
	up.Actions = string(data)
	return nil
}

// GetConditions returns the conditions as a map
func (up *UserPermission) GetConditions() (map[string]interface{}, error) {
	if up.Conditions == "" {
		return make(map[string]interface{}), nil
	}
	var conditions map[string]interface{}
	if err := json.Unmarshal([]byte(up.Conditions), &conditions); err != nil {
		return nil, err
	}
	return conditions, nil
}

// SetConditions sets the conditions from a map
func (up *UserPermission) SetConditions(conditions map[string]interface{}) error {
	if len(conditions) == 0 {
		up.Conditions = ""
		return nil
	}
	data, err := json.Marshal(conditions)
	if err != nil {
		return err
	}
	up.Conditions = string(data)
	return nil
}

// GetActions returns the actions as a slice
func (rp *RolePermission) GetActions() ([]PermissionAction, error) {
	var actions []PermissionAction
	if err := json.Unmarshal([]byte(rp.Actions), &actions); err != nil {
		return nil, err
	}
	return actions, nil
}

// SetActions sets the actions from a slice
func (rp *RolePermission) SetActions(actions []PermissionAction) error {
	data, err := json.Marshal(actions)
	if err != nil {
		return err
	}
	rp.Actions = string(data)
	return nil
}

// GetConditions returns the conditions as a map
func (rp *RolePermission) GetConditions() (map[string]interface{}, error) {
	if rp.Conditions == "" {
		return make(map[string]interface{}), nil
	}
	var conditions map[string]interface{}
	if err := json.Unmarshal([]byte(rp.Conditions), &conditions); err != nil {
		return nil, err
	}
	return conditions, nil
}

// SetConditions sets the conditions from a map
func (rp *RolePermission) SetConditions(conditions map[string]interface{}) error {
	if len(conditions) == 0 {
		rp.Conditions = ""
		return nil
	}
	data, err := json.Marshal(conditions)
	if err != nil {
		return err
	}
	rp.Conditions = string(data)
	return nil
}

// IsExpired checks if the permission has expired
func (up *UserPermission) IsExpired() bool {
	if up.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*up.ExpiresAt)
}

// IsExpired checks if the permission has expired
func (rp *RolePermission) IsExpired() bool {
	if rp.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*rp.ExpiresAt)
}

// HasAction checks if the permission includes a specific action
func (up *UserPermission) HasAction(action PermissionAction) bool {
	actions, err := up.GetActions()
	if err != nil {
		return false
	}

	for _, a := range actions {
		if a == action || a == ActionAll {
			return true
		}
	}
	return false
}

// HasAction checks if the permission includes a specific action
func (rp *RolePermission) HasAction(action PermissionAction) bool {
	actions, err := rp.GetActions()
	if err != nil {
		return false
	}

	for _, a := range actions {
		if a == action || a == ActionAll {
			return true
		}
	}
	return false
}

// CanRead checks if the permission allows read access
func (up *UserPermission) CanRead() bool {
	return up.HasAction(ActionRead) || up.HasAction(ActionAll) || up.Level == LevelReadOnly || up.Level == LevelReadWrite || up.Level == LevelFullAccess
}

// CanWrite checks if the permission allows write access
func (up *UserPermission) CanWrite() bool {
	return up.HasAction(ActionCreate) || up.HasAction(ActionUpdate) || up.HasAction(ActionDelete) || up.HasAction(ActionAll) || up.Level == LevelReadWrite || up.Level == LevelFullAccess
}

// CanManage checks if the permission allows management access
func (up *UserPermission) CanManage() bool {
	return up.HasAction(ActionManage) || up.HasAction(ActionAll) || up.Level == LevelFullAccess
}

// CanRead checks if the permission allows read access
func (rp *RolePermission) CanRead() bool {
	return rp.HasAction(ActionRead) || rp.HasAction(ActionAll) || rp.Level == LevelReadOnly || rp.Level == LevelReadWrite || rp.Level == LevelFullAccess
}

// CanWrite checks if the permission allows write access
func (rp *RolePermission) CanWrite() bool {
	return rp.HasAction(ActionCreate) || rp.HasAction(ActionUpdate) || rp.HasAction(ActionDelete) || rp.HasAction(ActionAll) || rp.Level == LevelReadWrite || rp.Level == LevelFullAccess
}

// CanManage checks if the permission allows management access
func (rp *RolePermission) CanManage() bool {
	return rp.HasAction(ActionManage) || rp.HasAction(ActionAll) || rp.Level == LevelFullAccess
}
