package inbox

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
)

// RegisterInboxRoutes registers HTTP endpoints for managing inbox items
func RegisterInboxRoutes(router *gin.Engine, service *InboxService, middleware *auth.AuthMiddleware) {
	// Create a specific group for inbox with authentication requirement
	group := router.Group("/api/v1/inbox")
	group.Use(middleware.Authenticate())

	// GET /api/v1/inbox - List inbox items
	group.GET("", func(c *gin.Context) {
		params := ListParams{
			Page:        parseIntDefault(c.Query("page"), 1),
			PageSize:    parseIntDefault(c.Query("page_size"), 20),
			Status:      c.Query("status"),
			Priority:    c.Query("priority"),
			Context:     c.Query("context"),
			SearchQuery: c.Query("q"),
		}

		realmID, _ := auth.GetCurrentRealm(c)

		result, err := service.List(realmID, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	// GET /api/v1/inbox/pending - Get pending items
	group.GET("/pending", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		items, err := service.GetPendingItems(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": items,
		})
	})

	// GET /api/v1/inbox/urgent - Get urgent items
	group.GET("/urgent", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		items, err := service.GetUrgentItems(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": items,
		})
	})

	// GET /api/v1/inbox/stats - Get inbox statistics
	group.GET("/stats", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		stats, err := service.GetStats(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, stats)
	})

	// POST /api/v1/inbox - Create a new inbox item
	group.POST("", func(c *gin.Context) {
		var req CreateInboxItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		realmID, _ := auth.GetCurrentRealm(c)
		userID, _ := auth.GetCurrentUser(c)

		item, err := service.CreateFromInput(&req, realmID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, item)
	})

	// GET /api/v1/inbox/:id - Get a specific inbox item
	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Inbox item ID is required",
			})
			return
		}

		item, err := service.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Inbox item not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, item)
	})

	// PUT /api/v1/inbox/:id - Update an inbox item
	group.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Inbox item ID is required",
			})
			return
		}

		var req UpdateInboxItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		userID, _ := auth.GetCurrentUser(c)

		item, err := service.UpdateFromInput(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Inbox item not found",
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
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, item)
	})

	// DELETE /api/v1/inbox/:id - Delete an inbox item
	group.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Inbox item ID is required",
			})
			return
		}

		err := service.Delete(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Inbox item not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	// PUT /api/v1/inbox/:id/status - Update inbox item status
	group.PUT("/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Inbox item ID is required",
			})
			return
		}

		var req struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		item, err := service.UpdateStatus(id, req.Status)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Inbox item not found",
				})
				return
			}
			if strings.Contains(err.Error(), "invalid status") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, item)
	})

	// PUT /api/v1/inbox/bulk/status - Bulk update status
	group.PUT("/bulk/status", func(c *gin.Context) {
		var req struct {
			IDs    []string `json:"ids" binding:"required"`
			Status string   `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		err := service.BulkUpdateStatus(req.IDs, req.Status)
		if err != nil {
			if strings.Contains(err.Error(), "invalid status") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Status updated successfully",
		})
	})
}

// Helper function to parse integer with default value
func parseIntDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}
