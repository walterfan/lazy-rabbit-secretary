package task

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// RegisterRoutes registers HTTP endpoints for managing tasks
func RegisterRoutes(router *gin.Engine, service *TaskService, middleware *auth.AuthMiddleware) {
	// Create a specific group for tasks with authentication requirement
	group := router.Group("/api/v1/tasks")
	group.Use(middleware.Authenticate())

	// GET /api/v1/tasks - Search/list tasks
	group.GET("", func(c *gin.Context) {

		query := c.Query("q")
		status := c.Query("status")
		tags := c.Query("tags")
		priority := parseIntDefault(c.Query("priority"), 0)
		difficulty := parseIntDefault(c.Query("difficulty"), 0)
		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)
		realmID, _ := auth.GetCurrentRealm(c)
		items, total, err := service.SearchTasks(realmID, query, status, tags, priority, difficulty, page, pageSize)
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

	// POST /api/v1/tasks - Create a new task
	group.POST("", func(c *gin.Context) {
		var req CreateTaskRequest
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

	// GET /api/v1/tasks/:id - Get a specific task
	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		item, err := service.GetTask(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// PUT /api/v1/tasks/:id - Update a task
	group.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req UpdateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updater, _ := auth.GetCurrentUsername(c)
		updated, err := service.UpdateTask(id, req, updater)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	// DELETE /api/v1/tasks/:id - Delete a task
	group.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeleteTask(id); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// GET /api/v1/tasks/status/:status - Get tasks by status
	group.GET("/status/:status", func(c *gin.Context) {
		status := c.Param("status")

		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)

		realmID, realmExists := auth.GetCurrentRealm(c)
		if !realmExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Validate status
		var taskStatus models.TaskStatus
		switch status {
		case "pending":
			taskStatus = models.TaskStatusPending
		case "running":
			taskStatus = models.TaskStatusRunning
		case "completed":
			taskStatus = models.TaskStatusCompleted
		case "failed":
			taskStatus = models.TaskStatusFailed
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}

		items, total, err := service.GetTasksByStatus(realmID, taskStatus, page, pageSize)
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

	// GET /api/v1/tasks/upcoming - Get upcoming tasks
	group.GET("/upcoming", func(c *gin.Context) {

		limit := parseIntDefault(c.Query("limit"), 10)

		realmID, realmExists := auth.GetCurrentRealm(c)
		if !realmExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		items, err := service.GetUpcomingTasks(realmID, limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
	})

	// GET /api/v1/tasks/overdue - Get overdue tasks
	group.GET("/overdue", func(c *gin.Context) {

		limit := parseIntDefault(c.Query("limit"), 10)
		realmID, realmExists := auth.GetCurrentRealm(c)
		if !realmExists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		items, err := service.GetOverdueTasks(realmID, limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
	})

	// POST /api/v1/tasks/:id/start - Start a task
	group.POST("/:id/start", func(c *gin.Context) {
		id := c.Param("id")
		starter, _ := auth.GetCurrentUsername(c)

		updated, err := service.StartTask(id, starter)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	// POST /api/v1/tasks/:id/complete - Complete a task
	group.POST("/:id/complete", func(c *gin.Context) {
		id := c.Param("id")
		completer, _ := auth.GetCurrentUsername(c)

		updated, err := service.CompleteTask(id, completer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	// POST /api/v1/tasks/:id/fail - Mark a task as failed
	group.POST("/:id/fail", func(c *gin.Context) {
		id := c.Param("id")
		failer, _ := auth.GetCurrentUsername(c)

		updated, err := service.FailTask(id, failer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
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
