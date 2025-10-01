package secret

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
)

// RegisterRoutes registers HTTP endpoints for managing secrets
func RegisterRoutes(router *gin.Engine, service *SecretService, middleware *auth.AuthMiddleware) {
	// Create a specific group for secrets with admin/super_admin restriction
	group := router.Group("/api/v1/secrets")
	group.Use(middleware.Authenticate())
	group.Use(middleware.RequireRole("admin", "super_admin"))

	group.GET("", func(c *gin.Context) {

		query := c.Query("q")
		grp := c.Query("group")
		path := c.Query("path")
		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)
		realmID, realmExists := auth.GetCurrentRealm(c)
		if !realmExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		items, total, err := service.SearchSecrets(realmID, query, grp, path, page, pageSize)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
	})

	group.POST("", func(c *gin.Context) {
		var req CreateSecretRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		realmID, realmExists := auth.GetCurrentRealm(c)
		if !realmExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		creator, _ := auth.GetCurrentUsername(c)
		created, err := service.CreateFromInput(req, realmID, creator)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, created)
	})

	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		item, err := service.GetSecret(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	group.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req UpdateSecretRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updater, _ := auth.GetCurrentUsername(c)
		updated, err := service.UpdateFromInput(id, req, updater)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	group.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeleteSecret(id); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	group.POST("/:id/decrypt", func(c *gin.Context) {
		id := c.Param("id")
		value, err := service.DecryptSecret(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	group.POST("/:id/decrypt-with-kek", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			KEK string `json:"kek" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		value, err := service.DecryptSecretWithKEK(id, req.KEK)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	// Version management endpoints
	group.GET("/:id/versions", func(c *gin.Context) {
		id := c.Param("id")
		versions, err := service.GetSecretVersions(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"versions": versions})
	})

	group.POST("/:id/versions/:version/decrypt", func(c *gin.Context) {
		id := c.Param("id")
		version := parseIntDefault(c.Param("version"), 0)
		if version <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
			return
		}
		value, err := service.DecryptSecretVersion(id, version)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	group.POST("/:id/versions/:version/decrypt-with-kek", func(c *gin.Context) {
		id := c.Param("id")
		version := parseIntDefault(c.Param("version"), 0)
		if version <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
			return
		}
		var req struct {
			KEK string `json:"kek" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		value, err := service.DecryptSecretVersionWithKEK(id, version, req.KEK)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	group.POST("/:id/versions/:version/activate", func(c *gin.Context) {
		id := c.Param("id")
		version := parseIntDefault(c.Param("version"), 0)
		if version <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
			return
		}
		updater, _ := auth.GetCurrentUsername(c)
		if err := service.ActivateSecretVersion(id, version, updater); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	group.DELETE("/:id/versions/:version", func(c *gin.Context) {
		id := c.Param("id")
		version := parseIntDefault(c.Param("version"), 0)
		if version <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
			return
		}
		if err := service.DeleteSecretVersion(id, version); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	group.POST("/:id/versions/pending", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Value string `json:"value" binding:"required"`
			KEK   string `json:"kek"` // Optional custom KEK
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		creator, _ := auth.GetCurrentUsername(c)
		version, err := service.CreatePendingVersion(id, req.Value, req.KEK, creator)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, version)
	})
}

func parseIntDefault(value string, defaultVal int) int {
	if value == "" {
		return defaultVal
	}
	var out int
	_, err := fmt.Sscanf(value, "%d", &out)
	if err != nil || out <= 0 {
		return defaultVal
	}
	return out
}
