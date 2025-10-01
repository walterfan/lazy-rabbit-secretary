package prompt

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PromptRoutes handles HTTP routes for prompts
type PromptRoutes struct {
	service *PromptService
}

// NewPromptRoutes creates new prompt routes
func NewPromptRoutes(db *gorm.DB) *PromptRoutes {
	return &PromptRoutes{
		service: NewPromptService(db),
	}
}

// RegisterRoutes registers prompt routes with the router
func (r *PromptRoutes) RegisterRoutes(router *gin.Engine, authMiddleware interface{}) {
	// Type assert the middleware
	middleware, ok := authMiddleware.(interface {
		Authenticate() gin.HandlerFunc
	})
	if !ok {
		panic("authMiddleware does not implement required interface")
	}

	prompts := router.Group("/api/v1/prompts")
	prompts.Use(middleware.Authenticate()) // Apply authentication middleware
	{
		prompts.POST("", r.createPrompt)
		prompts.GET("", r.listPrompts)
		prompts.GET("/search", r.searchPrompts)
		prompts.GET("/tags/:tags", r.getPromptsByTags)
		prompts.GET("/:id", r.getPrompt)
		prompts.PUT("/:id", r.updatePrompt)
		prompts.DELETE("/:id", r.deletePrompt)
	}
}

// createPrompt handles POST /prompts
func (r *PromptRoutes) createPrompt(c *gin.Context) {
	var req CreatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get realm ID and user ID from context (set by auth middleware)
	realmID, exists := c.Get("realm_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Realm ID not found in context",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	prompt, err := r.service.CreateFromInput(&req, realmID.(string), userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create prompt",
		})
		return
	}

	c.JSON(http.StatusCreated, prompt)
}

// getPrompt handles GET /prompts/:id
func (r *PromptRoutes) getPrompt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prompt ID is required",
		})
		return
	}

	prompt, err := r.service.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Prompt not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get prompt",
		})
		return
	}

	c.JSON(http.StatusOK, prompt)
}

// updatePrompt handles PUT /prompts/:id
func (r *PromptRoutes) updatePrompt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prompt ID is required",
		})
		return
	}

	var req UpdatePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	prompt, err := r.service.UpdateFromInput(id, &req, userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Prompt not found",
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update prompt",
		})
		return
	}

	c.JSON(http.StatusOK, prompt)
}

// deletePrompt handles DELETE /prompts/:id
func (r *PromptRoutes) deletePrompt(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prompt ID is required",
		})
		return
	}

	err := r.service.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Prompt not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete prompt",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prompt deleted successfully",
	})
}

// listPrompts handles GET /prompts
func (r *PromptRoutes) listPrompts(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := r.service.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list prompts",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// searchPrompts handles GET /prompts/search
func (r *PromptRoutes) searchPrompts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := r.service.Search(query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search prompts",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// getPromptsByTags handles GET /prompts/tags/:tags
func (r *PromptRoutes) getPromptsByTags(c *gin.Context) {
	tagsParam := c.Param("tags")
	if tagsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tags parameter is required",
		})
		return
	}

	// Split tags by comma
	tags := strings.Split(tagsParam, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := r.service.GetByTags(tags, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get prompts by tags",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
