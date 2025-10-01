package bookmark

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"gorm.io/gorm"
)

// BookmarkRoutes handles HTTP routes for bookmarks
type BookmarkRoutes struct {
	service *BookmarkService
}

// NewBookmarkRoutes creates new bookmark routes
func NewBookmarkRoutes(db *gorm.DB) *BookmarkRoutes {
	return &BookmarkRoutes{
		service: NewBookmarkService(db),
	}
}

// RegisterRoutes registers bookmark routes with the router
func RegisterRoutes(router *gin.Engine, service *BookmarkService, middleware *auth.AuthMiddleware) {
	// Public routes (optional authentication - check auth header if present)
	publicGroup := router.Group("/api/v1/bookmarks")
	publicGroup.Use(middleware.OptionalAuth())
	{
		publicGroup.GET("/search", searchBookmarks(service))
		publicGroup.GET("/recent", getRecentBookmarks(service))
		publicGroup.GET("/tags", getAllTags(service))
		publicGroup.GET("/tags/popular", getPopularTags(service))
		publicGroup.GET("/category/:categoryId", getBookmarksByCategory(service))
		publicGroup.GET("/tag/:tagName", getBookmarksByTag(service))
		publicGroup.GET("/stats", getBookmarkStats(service))
	}

	// Admin routes (authentication required)
	adminGroup := router.Group("/api/v1/admin/bookmarks")
	adminGroup.Use(middleware.Authenticate())
	{
		adminGroup.POST("", createBookmark(service))
		adminGroup.GET("", listBookmarks(service))
		adminGroup.GET("/:id", getBookmark(service))
		adminGroup.PUT("/:id", updateBookmark(service))
		adminGroup.DELETE("/:id", deleteBookmark(service))

		// Bulk operations
		adminGroup.POST("/bulk/import", bulkImportBookmarks(service))
		adminGroup.POST("/bulk/export", bulkExportBookmarks(service))
		adminGroup.DELETE("/bulk", bulkDeleteBookmarks(service))

		// Category management
		adminGroup.POST("/categories", createCategory(service))
		adminGroup.GET("/categories", listCategories(service))
		adminGroup.GET("/categories/:id", getCategory(service))
		adminGroup.PUT("/categories/:id", updateCategory(service))
		adminGroup.DELETE("/categories/:id", deleteCategory(service))
	}
}

// Public handlers

func searchBookmarks(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return
		}

		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		result, err := service.SearchBookmarks(query, realmID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getRecentBookmarks(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		bookmarks, err := service.GetRecentBookmarks(realmID, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})
	}
}

func getAllTags(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		tags, err := service.GetAllTags(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}

func getPopularTags(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		tags, err := service.GetPopularTags(realmID, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}

func getBookmarksByCategory(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("categoryId")
		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		result, err := service.GetBookmarksByCategory(categoryID, realmID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getBookmarksByTag(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagName := c.Param("tagName")
		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

		result, err := service.GetBookmarksByTag(tagName, realmID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getBookmarkStats(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		stats, err := service.GetBookmarkStats(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

// Admin handlers

func createBookmark(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateBookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		bookmark, err := service.CreateBookmark(&req, realmID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, bookmark)
	}
}

func listBookmarks(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BookmarkListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse tags from query parameter
		if tagsParam := c.Query("tags"); tagsParam != "" {
			req.Tags = strings.Split(tagsParam, ",")
			for i, tag := range req.Tags {
				req.Tags[i] = strings.TrimSpace(tag)
			}
		}

		realmID := c.GetString("realm_id")
		if realmID == "" {
			realmID = "default"
		}

		result, err := service.ListBookmarks(&req, realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getBookmark(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		bookmark, err := service.GetBookmark(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bookmark)
	}
}

func updateBookmark(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req UpdateBookmarkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		bookmark, err := service.UpdateBookmark(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bookmark)
	}
}

func deleteBookmark(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := service.DeleteBookmark(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bookmark deleted successfully"})
	}
}

// Category handlers

func createCategory(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateBookmarkCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		category, err := service.CreateCategory(&req, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, category)
	}
}

func listCategories(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := service.ListCategories()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}

func getCategory(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		category, err := service.GetCategory(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func updateCategory(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		var req UpdateBookmarkCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		category, err := service.UpdateCategory(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func deleteCategory(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}

		err = service.DeleteCategory(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
	}
}

// Bulk operation handlers (placeholder implementations)

func bulkImportBookmarks(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement bulk import functionality
		// This could accept various formats like CSV, JSON, or browser bookmark exports
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Bulk import not yet implemented"})
	}
}

func bulkExportBookmarks(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement bulk export functionality
		// This could export to various formats like CSV, JSON, or HTML
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Bulk export not yet implemented"})
	}
}

func bulkDeleteBookmarks(service *BookmarkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement bulk delete functionality
		// This could accept an array of bookmark IDs to delete
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Bulk delete not yet implemented"})
	}
}
