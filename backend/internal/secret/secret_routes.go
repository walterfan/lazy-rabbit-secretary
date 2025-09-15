package secret

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-reminder/internal/auth"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
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
		var s models.Secret
		if err := c.ShouldBindJSON(&s); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s.ID = id
		if err := service.UpdateSecret(&s); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
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
