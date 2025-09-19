package daily

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
)

// RegisterDailyRoutes registers HTTP endpoints for managing daily checklist items
func RegisterDailyRoutes(router *gin.Engine, service *DailyService, middleware *auth.AuthMiddleware) {
	// Create a specific group for daily checklist with authentication requirement
	group := router.Group("/api/v1/daily")
	group.Use(middleware.Authenticate())

	// GET /api/v1/daily - List daily checklist items
	group.GET("", func(c *gin.Context) {
		params := ListParams{
			Page:        parseIntDefault(c.Query("page"), 1),
			PageSize:    parseIntDefault(c.Query("page_size"), 20),
			Status:      c.Query("status"),
			Priority:    c.Query("priority"),
			Context:     c.Query("context"),
			SearchQuery: c.Query("q"),
		}

		// Parse date if provided
		if dateStr := c.Query("date"); dateStr != "" {
			if date, err := time.Parse("2006-01-02", dateStr); err == nil {
				params.Date = date
			}
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

	// GET /api/v1/daily/date/:date - Get items for specific date
	group.GET("/date/:date", func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
			return
		}

		realmID, _ := auth.GetCurrentRealm(c)
		items, err := service.GetByDate(realmID, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": items,
			"date":  date.Format("2006-01-02"),
		})
	})

	// GET /api/v1/daily/stats/:date - Get daily statistics
	group.GET("/stats/:date", func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
			return
		}

		realmID, _ := auth.GetCurrentRealm(c)
		stats, err := service.GetStats(realmID, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, stats)
	})

	// POST /api/v1/daily - Create a new daily checklist item
	group.POST("", func(c *gin.Context) {
		var req CreateDailyItemRequest
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

	// GET /api/v1/daily/:id - Get a specific daily checklist item
	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Daily checklist item ID is required",
			})
			return
		}

		item, err := service.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Daily checklist item not found",
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

	// PUT /api/v1/daily/:id - Update a daily checklist item
	group.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Daily checklist item ID is required",
			})
			return
		}

		var req UpdateDailyItemRequest
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
					"error": "Daily checklist item not found",
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

	// DELETE /api/v1/daily/:id - Delete a daily checklist item
	group.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Daily checklist item ID is required",
			})
			return
		}

		err := service.Delete(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Daily checklist item not found",
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

	// PUT /api/v1/daily/:id/status - Update daily checklist item status
	group.PUT("/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Daily checklist item ID is required",
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
					"error": "Daily checklist item not found",
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

	// PUT /api/v1/daily/:id/time - Update actual time spent
	group.PUT("/:id/time", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Daily checklist item ID is required",
			})
			return
		}

		var req struct {
			ActualTime int `json:"actual_time" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		item, err := service.UpdateActualTime(id, req.ActualTime)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Daily checklist item not found",
				})
				return
			}
			if strings.Contains(err.Error(), "invalid actual time") {
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

	// PUT /api/v1/daily/bulk/status - Bulk update status
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
