package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/database"
	"gorm.io/gorm"
)

// RoleHandlers provides HTTP handlers for role management
type RoleHandlers struct {
	authService *AuthService
}

// NewRoleHandlers creates a new role handlers
func NewRoleHandlers(authService *AuthService) *RoleHandlers {
	return &RoleHandlers{
		authService: authService,
	}
}

// GetRoles handles getting roles with pagination and filtering
func (h *RoleHandlers) GetRoles(c *gin.Context) {
	// Get current user for audit trail
	_, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse query parameters
	realmName := c.Query("realm_name")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	// Set defaults
	page := 1
	pageSize := 10

	// Parse pagination
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	db := database.GetDB()
	var roles []models.Role
	var total int64

	// Build query
	query := db.Model(&models.Role{})

	// Filter by realm if specified
	if realmName != "" {
		realm, err := h.authService.userService.GetRealmByName(realmName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid realm"})
			return
		}
		query = query.Where("realm_id = ?", realm.ID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count roles", "details": err.Error()})
		return
	}

	// Apply pagination and preload policies
	offset := (page - 1) * pageSize
	if err := query.Preload("Policies").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get roles", "details": err.Error()})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"roles":       roles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// CreateRole handles creating a new role
func (h *RoleHandlers) CreateRole(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req struct {
		RealmName   string `json:"realm_name" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get realm
	realm, err := h.authService.userService.GetRealmByName(req.RealmName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid realm"})
		return
	}

	// Check if role name already exists in realm
	db := database.GetDB()
	var existingRole models.Role
	if err := db.Where("name = ? AND realm_id = ?", req.Name, realm.ID).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists in this realm"})
		return
	}

	// Create role
	role := &models.Role{
		ID:          uuid.New().String(),
		RealmID:     realm.ID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   currentUserID,
		CreatedAt:   time.Now(),
		UpdatedBy:   currentUserID,
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"role":    role,
	})
}

// GetRole handles getting a role by ID
func (h *RoleHandlers) GetRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	db := database.GetDB()
	var role models.Role

	if err := db.Where("id = ?", roleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

// UpdateRole handles updating a role
func (h *RoleHandlers) UpdateRole(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	db := database.GetDB()
	var role models.Role

	// Get existing role
	if err := db.Where("id = ?", roleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role", "details": err.Error()})
		}
		return
	}

	// Update fields
	if req.Name != "" {
		// Check if new name already exists in realm
		var existingRole models.Role
		if err := db.Where("name = ? AND realm_id = ? AND id != ?", req.Name, role.RealmID, roleID).First(&existingRole).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Role name already exists in this realm"})
			return
		}
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	role.UpdatedBy = currentUserID
	role.UpdatedAt = time.Now()

	// Save updated role
	if err := db.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"role":    role,
	})
}

// DeleteRole handles deleting a role
func (h *RoleHandlers) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role ID is required"})
		return
	}

	db := database.GetDB()

	// Check if role exists
	var role models.Role
	if err := db.Where("id = ?", roleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role", "details": err.Error()})
		}
		return
	}

	// Check if role is assigned to any users
	var userRoleCount int64
	if err := db.Model(&models.UserRole{}).Where("role_id = ?", roleID).Count(&userRoleCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role assignments", "details": err.Error()})
		return
	}

	if userRoleCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete role that is assigned to users"})
		return
	}

	// Delete role
	if err := db.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
