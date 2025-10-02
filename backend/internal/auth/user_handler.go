package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// UserHandlers provides HTTP handlers for user management
type UserHandlers struct {
	authService *AuthService
}

// NewUserHandlers creates a new user handlers
func NewUserHandlers(authService *AuthService) *UserHandlers {
	return &UserHandlers{
		authService: authService,
	}
}

// GetUsers handles getting users with pagination and filtering
func (h *UserHandlers) GetUsers(c *gin.Context) {
	// Get current user for audit trail
	_, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse query parameters
	realmName := c.Query("realm_name")
	status := c.Query("status")
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

	// Build request
	req := models.UserRegistrationRequest{
		RealmName: realmName,
		Status:    models.UserStatus(status),
		Page:      page,
		PageSize:  pageSize,
	}

	// Get users
	response, err := h.authService.GetPendingRegistrations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateUser handles creating a new user
func (h *UserHandlers) CreateUser(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	user, err := h.authService.RegisterUser(req, currentUserID)
	if err != nil {
		switch err {
		case ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists in this realm"})
		case ErrInvalidRealm:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid realm"})
		case ErrPasswordTooShort:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password too short"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUser handles getting a user by ID
func (h *UserHandlers) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := h.authService.userService.GetUserByID(userID)
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser handles updating a user
func (h *UserHandlers) UpdateUser(c *gin.Context) {
	// Get current user for audit trail
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get existing user
	user, err := h.authService.userService.GetUserByID(userID)
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user", "details": err.Error()})
		}
		return
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	// Handle password change if requested
	if req.NewPassword != "" {
		// If changing password, current password must be provided (unless admin is changing it)
		if req.CurrentPassword == "" {
			// Check if current user is admin or the user themselves
			if currentUserID != userID {
				// Check if current user has admin role
				roles, hasRoles := GetCurrentRoles(c)
				isAdmin := false
				if hasRoles {
					for _, roleName := range roles {
						if roleName == "admin" || roleName == "super_admin" {
							isAdmin = true
							break
						}
					}
				}

				if !isAdmin {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is required to change password"})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is required to change password"})
				return
			}
		} else {
			// Verify current password
			if !h.authService.VerifyPassword(req.CurrentPassword, user.HashedPassword) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
				return
			}
		}

		// Check new password strength
		if err := h.authService.CheckPasswordStrength(req.NewPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password does not meet requirements", "details": err.Error()})
			return
		}

		// Hash new password
		hashedPassword, err := h.authService.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password", "details": err.Error()})
			return
		}

		user.HashedPassword = hashedPassword
	}

	user.UpdatedBy = currentUserID
	user.UpdatedAt = time.Now()

	// Save updated user
	if err := h.authService.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	// Update user roles if provided (including empty array to remove all roles)
	if req.RoleIDs != nil {
		if err := h.authService.userService.UpdateUserRoles(userID, req.RoleIDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user roles", "details": err.Error()})
			return
		}
	}

	// Get updated user with roles for response
	updatedUser, err := h.authService.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}

// DeleteUser handles deleting a user
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Check if user exists
	_, err := h.authService.userService.GetUserByID(userID)
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user", "details": err.Error()})
		}
		return
	}

	// Delete user
	if err := h.authService.userService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetPendingRegistrations handles getting pending user registrations
func (h *UserHandlers) GetPendingRegistrations(c *gin.Context) {
	var req models.UserRegistrationRequest

	// Parse query parameters
	if realmName := c.Query("realm_name"); realmName != "" {
		req.RealmName = realmName
	}
	if status := c.Query("status"); status != "" {
		req.Status = models.UserStatus(status)
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			req.Page = p
		}
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			req.PageSize = ps
		}
	}

	// Default to pending status if no status specified
	if req.Status == "" {
		req.Status = models.UserStatusPending
	}

	response, err := h.authService.GetPendingRegistrations(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get registrations", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ApproveRegistration handles approving or denying user registration
func (h *UserHandlers) ApproveRegistration(c *gin.Context) {
	var req models.ApproveRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get current user from context (who is approving)
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	response, err := h.authService.ApproveRegistration(req, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration approval", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRegistrationStats handles getting registration statistics
func (h *UserHandlers) GetRegistrationStats(c *gin.Context) {
	realmName := c.Query("realm_name")

	stats, err := h.authService.GetUserRegistrationStats(realmName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get registration stats", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
