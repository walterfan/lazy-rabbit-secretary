package auth

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	authService *AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// isAuthDisabled checks if authentication is disabled via DISABLE_AUTH environment variable
func isAuthDisabled() bool {
	disableAuth := os.Getenv("DISABLE_AUTH")
	if disableAuth == "" {
		return false // Default to enabled
	}

	disabled, err := strconv.ParseBool(disableAuth)
	if err != nil {
		return false // Default to enabled if parsing fails
	}

	return disabled
}

// Authenticate validates JWT tokens and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if authentication is disabled
		if isAuthDisabled() {
			// Set default context values when auth is disabled
			c.Set("user_id", "disabled-auth-user")
			c.Set("realm_id", "default-realm")
			c.Set("username", "disabled-auth")
			c.Set("email", "disabled@auth.local")
			c.Set("roles", []string{"admin"})
			c.Next()
			return
		}

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("realm_id", claims.RealmID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// OptionalAuth attempts to authenticate but allows access even without authentication
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if authentication is disabled
		if isAuthDisabled() {
			// Set default context values when auth is disabled
			c.Set("user_id", "disabled-auth-user")
			c.Set("realm_id", "default-realm")
			c.Set("username", "disabled-auth")
			c.Set("email", "disabled@auth.local")
			c.Set("roles", []string{"admin"})
			c.Next()
			return
		}

		// Try to get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without authentication
			c.Next()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// Try to validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Token is valid, set user context
		c.Set("user_id", claims.UserID)
		c.Set("realm_id", claims.RealmID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// RequireRealm ensures the user belongs to the specified realm
func (m *AuthMiddleware) RequireRealm(realmID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRealmID, exists := c.Get("realm_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if userRealmID.(uuid.UUID) != realmID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this realm"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole ensures the user has at least one of the specified roles
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userRolesList := userRoles.([]string)
		hasRole := false

		for _, requiredRole := range roles {
			for _, userRole := range userRolesList {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission checks if the user has permission to perform an action on a resource
func (m *AuthMiddleware) RequirePermission(action, resource string) gin.HandlerFunc {
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

		// Build context from request parameters
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
		context["user:id"] = userID.(string)
		context["user:realm_id"] = realmID.(string)
		context["user:username"] = c.GetString("username")
		context["user:email"] = c.GetString("email")
		context["user:roles"] = c.GetStringSlice("roles")

		// Check permission
		allowed, err := m.authService.CheckPermission(
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

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser extracts the current user ID from context
func GetCurrentUser(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetCurrentRealm extracts the current realm ID from context
func GetCurrentRealm(c *gin.Context) (string, bool) {
	realmID, exists := c.Get("realm_id")
	if !exists {
		return "", false
	}
	return realmID.(string), true
}

// GetCurrentUsername extracts the current username from context
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// GetCurrentRoles extracts the current user roles from context
func GetCurrentRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get("roles")
	if !exists {
		return nil, false
	}
	return roles.([]string), true
}
