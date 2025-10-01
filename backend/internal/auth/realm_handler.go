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

// RealmHandlers provides HTTP handlers for realm management
type RealmHandlers struct {
	authService *AuthService
}

// NewRealmHandlers creates a new realm handlers
func NewRealmHandlers(authService *AuthService) *RealmHandlers {
	return &RealmHandlers{
		authService: authService,
	}
}

// GetRealms handles getting realms with pagination
func (h *RealmHandlers) GetRealms(c *gin.Context) {
	// Get current user for audit trail
	_, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse query parameters
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
	var realms []models.Realm
	var total int64

	// Count total records
	if err := db.Model(&models.Realm{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count realms", "details": err.Error()})
		return
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&realms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get realms", "details": err.Error()})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, gin.H{
		"realms":      realms,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// CreateRealm handles creating a new realm
func (h *RealmHandlers) CreateRealm(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Check if realm name already exists
	db := database.GetDB()
	var existingRealm models.Realm
	if err := db.Where("name = ?", req.Name).First(&existingRealm).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Realm name already exists"})
		return
	}

	// Create realm
	realm := &models.Realm{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   currentUserID,
		CreatedAt:   time.Now(),
		UpdatedBy:   currentUserID,
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(realm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create realm", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Realm created successfully",
		"realm":   realm,
	})
}

// GetRealm handles getting a realm by ID
func (h *RealmHandlers) GetRealm(c *gin.Context) {
	realmID := c.Param("id")
	if realmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Realm ID is required"})
		return
	}

	db := database.GetDB()
	var realm models.Realm

	if err := db.Where("id = ?", realmID).First(&realm).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Realm not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get realm", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"realm": realm})
}

// UpdateRealm handles updating a realm
func (h *RealmHandlers) UpdateRealm(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	realmID := c.Param("id")
	if realmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Realm ID is required"})
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
	var realm models.Realm

	// Get existing realm
	if err := db.Where("id = ?", realmID).First(&realm).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Realm not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get realm", "details": err.Error()})
		}
		return
	}

	// Update fields
	if req.Name != "" {
		// Check if new name already exists
		var existingRealm models.Realm
		if err := db.Where("name = ? AND id != ?", req.Name, realmID).First(&existingRealm).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Realm name already exists"})
			return
		}
		realm.Name = req.Name
	}
	if req.Description != "" {
		realm.Description = req.Description
	}

	realm.UpdatedBy = currentUserID
	realm.UpdatedAt = time.Now()

	// Save updated realm
	if err := db.Save(&realm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update realm", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Realm updated successfully",
		"realm":   realm,
	})
}

// DeleteRealm handles deleting a realm
func (h *RealmHandlers) DeleteRealm(c *gin.Context) {
	realmID := c.Param("id")
	if realmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Realm ID is required"})
		return
	}

	db := database.GetDB()

	// Check if realm exists
	var realm models.Realm
	if err := db.Where("id = ?", realmID).First(&realm).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Realm not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get realm", "details": err.Error()})
		}
		return
	}

	// Check if realm has any users
	var userCount int64
	if err := db.Model(&models.User{}).Where("realm_id = ?", realmID).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check realm users", "details": err.Error()})
		return
	}

	if userCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete realm that contains users"})
		return
	}

	// Check if realm has any roles
	var roleCount int64
	if err := db.Model(&models.Role{}).Where("realm_id = ?", realmID).Count(&roleCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check realm roles", "details": err.Error()})
		return
	}

	if roleCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete realm that contains roles"})
		return
	}

	// Check if realm has any policies
	var policyCount int64
	if err := db.Model(&models.Policy{}).Where("realm_id = ?", realmID).Count(&policyCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check realm policies", "details": err.Error()})
		return
	}

	if policyCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete realm that contains policies"})
		return
	}

	// Delete realm
	if err := db.Delete(&realm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete realm", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Realm deleted successfully"})
}
