package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// AuthHandlers provides HTTP handlers for authentication
type AuthHandlers struct {
	authService *AuthService
}

// NewAuthHandlers creates new authentication handlers
func NewAuthHandlers(authService *AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

// Login handles user authentication
func (h *AuthHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		case ErrUserInactive:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is inactive"})
		case ErrInvalidRealm:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid realm"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		case ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		case ErrUserInactive:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is inactive"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token refresh failed"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register handles user registration
func (h *AuthHandlers) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get current user from context (for created_by field)
	currentUserID, exists := GetCurrentUser(c)
	if !exists {
		// For public registration, use a system user ID or empty string
		currentUserID = ""
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User registration failed"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email to confirm your account.",
		"user":    user,
	})
}

// ConfirmEmail confirms a user's email address
func (h *AuthHandlers) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Confirmation token is required"})
		return
	}

	if err := h.authService.ConfirmEmail(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email confirmed successfully. Your registration is now pending admin approval.",
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandlers) GetProfile(c *gin.Context) {
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

	username, exists := GetCurrentUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roles, exists := GetCurrentRoles(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile := gin.H{
		"user_id":  userID,
		"realm_id": realmID,
		"username": username,
		"email":    email,
		"roles":    roles,
	}

	c.JSON(http.StatusOK, profile)
}

// CheckPermission checks if the current user has permission for an action
func (h *AuthHandlers) CheckPermission(c *gin.Context) {
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

	allowed, err := h.authService.CheckPermission(userID, realmID, action, resource, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"action":   action,
		"resource": resource,
		"allowed":  allowed,
	})
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandlers) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// The server can't invalidate tokens, but we can return success
	// In production, you might want to implement a token blacklist
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// HealthCheck provides a health check endpoint
func (h *AuthHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth",
	})
}
