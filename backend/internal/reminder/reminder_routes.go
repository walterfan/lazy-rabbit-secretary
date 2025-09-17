package reminder

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-reminder/internal/auth"
)

// RegisterRoutes registers HTTP endpoints for managing reminders
func RegisterRoutes(router *gin.Engine, service *ReminderService, middleware *auth.AuthMiddleware) {
	// Create a specific group for reminders with authentication requirement
	group := router.Group("/api/v1/reminders")
	group.Use(middleware.Authenticate())

	// GET /api/v1/reminders - Search/list reminders
	group.GET("", func(c *gin.Context) {
		query := c.Query("q")
		status := c.Query("status")
		tags := c.Query("tags")
		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)

		realmID, _ := auth.GetCurrentRealm(c)
		items, total, err := service.SearchReminders(realmID, query, status, tags, page, pageSize)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":       items,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		})
	})

	// POST /api/v1/reminders - Create a new reminder
	group.POST("", func(c *gin.Context) {
		var req CreateReminderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		realmID, _ := auth.GetCurrentRealm(c)
		creator, _ := auth.GetCurrentUser(c)
		created, err := service.CreateFromInput(req, realmID, creator)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, created)
	})

	// GET /api/v1/reminders/:id - Get a specific reminder
	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		item, err := service.GetReminder(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// PUT /api/v1/reminders/:id - Update a reminder
	group.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req UpdateReminderRequest
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

	// DELETE /api/v1/reminders/:id - Delete a reminder
	group.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeleteReminder(id); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted successfully"})
	})

	// GET /api/v1/reminders/status/:status - Get reminders by status
	group.GET("/status/:status", func(c *gin.Context) {
		status := c.Param("status")
		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)

		realmID, _ := auth.GetCurrentRealm(c)
		items, total, err := service.GetRemindersByStatus(realmID, status, page, pageSize)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":       items,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		})
	})

	// GET /api/v1/reminders/upcoming - Get upcoming reminders
	group.GET("/upcoming", func(c *gin.Context) {
		limit := parseIntDefault(c.Query("limit"), 10)

		realmID, _ := auth.GetCurrentRealm(c)
		items, err := service.GetUpcomingReminders(realmID, limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	// GET /api/v1/reminders/overdue - Get overdue reminders
	group.GET("/overdue", func(c *gin.Context) {
		limit := parseIntDefault(c.Query("limit"), 10)

		realmID, _ := auth.GetCurrentRealm(c)
		items, err := service.GetOverdueReminders(realmID, limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	// GET /api/v1/reminders/range - Get reminders by time range
	group.GET("/range", func(c *gin.Context) {
		startTimeStr := c.Query("start_time")
		endTimeStr := c.Query("end_time")
		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)

		if startTimeStr == "" || endTimeStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "start_time and end_time are required"})
			return
		}

		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format, use RFC3339"})
			return
		}

		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format, use RFC3339"})
			return
		}

		realmID, _ := auth.GetCurrentRealm(c)
		items, total, err := service.GetRemindersByTimeRange(realmID, startTime, endTime, page, pageSize)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items":       items,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			"start_time":  startTime,
			"end_time":    endTime,
		})
	})

	// POST /api/v1/reminders/:id/complete - Mark reminder as completed
	group.POST("/:id/complete", func(c *gin.Context) {
		id := c.Param("id")
		updater, _ := auth.GetCurrentUsername(c)

		updated, err := service.MarkAsCompleted(id, updater)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	// POST /api/v1/reminders/:id/snooze - Snooze a reminder
	group.POST("/:id/snooze", func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			RemindTime time.Time `json:"remind_time" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updater, _ := auth.GetCurrentUsername(c)
		updated, err := service.SnoozeReminder(id, req.RemindTime, updater)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})
}

// parseIntDefault parses string to int with default value
func parseIntDefault(str string, defaultValue int) int {
	if str == "" {
		return defaultValue
	}
	if val, err := strconv.Atoi(str); err == nil {
		return val
	}
	return defaultValue
}
