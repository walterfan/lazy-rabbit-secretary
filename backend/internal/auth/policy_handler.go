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

// PolicyHandlers provides HTTP handlers for policy management
type PolicyHandlers struct {
	authService *AuthService
}

// NewPolicyHandlers creates a new policy handlers
func NewPolicyHandlers(authService *AuthService) *PolicyHandlers {
	return &PolicyHandlers{
		authService: authService,
	}
}

// GetPolicies handles getting policies with pagination and filtering
func (h *PolicyHandlers) GetPolicies(c *gin.Context) {
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
	var policies []models.Policy
	var total int64

	// Build query
	query := db.Model(&models.Policy{})

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count policies", "details": err.Error()})
		return
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&policies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get policies", "details": err.Error()})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"policies":    policies,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// CreatePolicy handles creating a new policy
func (h *PolicyHandlers) CreatePolicy(c *gin.Context) {
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
		Version     string `json:"version"`
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

	// Check if policy name already exists in realm
	db := database.GetDB()
	var existingPolicy models.Policy
	if err := db.Where("name = ? AND realm_id = ?", req.Name, realm.ID).First(&existingPolicy).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Policy name already exists in this realm"})
		return
	}

	// Set default version if not provided
	version := req.Version
	if version == "" {
		version = "2012-10-17" // AWS policy version format
	}

	// Create policy
	policy := &models.Policy{
		ID:          uuid.New().String(),
		RealmID:     realm.ID,
		Name:        req.Name,
		Description: req.Description,
		Version:     version,
		CreatedBy:   currentUserID,
		CreatedAt:   time.Now(),
		UpdatedBy:   currentUserID,
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(policy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create policy", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Policy created successfully",
		"policy":  policy,
	})
}

// GetPolicy handles getting a policy by ID
func (h *PolicyHandlers) GetPolicy(c *gin.Context) {
	policyID := c.Param("id")
	if policyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Policy ID is required"})
		return
	}

	db := database.GetDB()
	var policy models.Policy

	if err := db.Where("id = ?", policyID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get policy", "details": err.Error()})
		}
		return
	}

	// Get policy statements
	var statements []models.Statement
	if err := db.Where("policy_id = ?", policyID).Find(&statements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get policy statements", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"policy":     policy,
		"statements": statements,
	})
}

// UpdatePolicy handles updating a policy
func (h *PolicyHandlers) UpdatePolicy(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	policyID := c.Param("id")
	if policyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Policy ID is required"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	db := database.GetDB()
	var policy models.Policy

	// Get existing policy
	if err := db.Where("id = ?", policyID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get policy", "details": err.Error()})
		}
		return
	}

	// Update fields
	if req.Name != "" {
		// Check if new name already exists in realm
		var existingPolicy models.Policy
		if err := db.Where("name = ? AND realm_id = ? AND id != ?", req.Name, policy.RealmID, policyID).First(&existingPolicy).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Policy name already exists in this realm"})
			return
		}
		policy.Name = req.Name
	}
	if req.Description != "" {
		policy.Description = req.Description
	}
	if req.Version != "" {
		policy.Version = req.Version
	}

	policy.UpdatedBy = currentUserID
	policy.UpdatedAt = time.Now()

	// Save updated policy
	if err := db.Save(&policy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update policy", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Policy updated successfully",
		"policy":  policy,
	})
}

// DeletePolicy handles deleting a policy
func (h *PolicyHandlers) DeletePolicy(c *gin.Context) {
	policyID := c.Param("id")
	if policyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Policy ID is required"})
		return
	}

	db := database.GetDB()

	// Check if policy exists
	var policy models.Policy
	if err := db.Where("id = ?", policyID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get policy", "details": err.Error()})
		}
		return
	}

	// Check if policy is assigned to any roles or users
	var rolePolicyCount int64
	var userPolicyCount int64

	if err := db.Model(&models.RolePolicy{}).Where("policy_id = ?", policyID).Count(&rolePolicyCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role policy assignments", "details": err.Error()})
		return
	}

	if err := db.Model(&models.UserPolicy{}).Where("policy_id = ?", policyID).Count(&userPolicyCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user policy assignments", "details": err.Error()})
		return
	}

	if rolePolicyCount > 0 || userPolicyCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete policy that is assigned to roles or users"})
		return
	}

	// Delete policy statements first
	if err := db.Where("policy_id = ?", policyID).Delete(&models.Statement{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete policy statements", "details": err.Error()})
		return
	}

	// Delete policy
	if err := db.Delete(&policy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete policy", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Policy deleted successfully"})
}
