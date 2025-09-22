package wiki

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// WikiRoutes handles HTTP routes for wiki pages
type WikiRoutes struct {
	service *WikiService
}

// NewWikiRoutes creates new wiki routes
func NewWikiRoutes(db *gorm.DB) *WikiRoutes {
	return &WikiRoutes{
		service: NewWikiService(db),
	}
}

// RegisterRoutes registers wiki routes with the router
func RegisterRoutes(router *gin.Engine, service *WikiService, middleware *auth.AuthMiddleware) {
	// Public routes (optional authentication - check auth header if present)
	publicGroup := router.Group("/api/v1/wiki")
	publicGroup.Use(middleware.OptionalAuth())
	{
		// Page access
		publicGroup.GET("/page/:slug", getWikiPageBySlug(service))
		publicGroup.GET("/pages", listWikiPages(service))
		publicGroup.GET("/search", searchWikiPages(service))
		publicGroup.GET("/category/:category", getPagesByCategory(service))
		publicGroup.GET("/tag/:tag", getPagesByTag(service))

		// Special pages
		publicGroup.GET("/random", getRandomPage(service))
		publicGroup.GET("/recent-changes", getRecentChanges(service))
		publicGroup.GET("/orphaned", getOrphanedPages(service))
		publicGroup.GET("/wanted", getWantedPages(service))
		publicGroup.GET("/dead-end", getDeadEndPages(service))

		// Page history (read-only)
		publicGroup.GET("/page/:slug/history", getPageHistory(service))
		publicGroup.GET("/page/:slug/compare/:from/:to", compareRevisions(service))
	}

	// Admin routes (authentication required)
	adminGroup := router.Group("/api/v1/admin/wiki")
	adminGroup.Use(middleware.Authenticate())
	{
		// Page management
		adminGroup.POST("/pages", createWikiPage(service))
		adminGroup.PUT("/pages/:id", updateWikiPage(service))
		adminGroup.DELETE("/pages/:id", deleteWikiPage(service))
		adminGroup.POST("/pages/:id/lock", lockWikiPage(service))
		adminGroup.POST("/pages/:id/unlock", unlockWikiPage(service))
		adminGroup.POST("/pages/:id/archive", archiveWikiPage(service))
		adminGroup.POST("/pages/:id/publish", publishWikiPage(service))
		adminGroup.POST("/pages/:id/unpublish", unpublishWikiPage(service))

		// Revision management
		adminGroup.POST("/pages/:id/revisions", createRevision(service))
		adminGroup.POST("/pages/:id/revert/:revisionId", revertToRevision(service))

		// Bulk operations
		adminGroup.POST("/pages/bulk-update", bulkUpdatePages(service))
		adminGroup.POST("/pages/bulk-delete", bulkDeletePages(service))
	}
}

// Public handlers

func getWikiPageBySlug(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page slug is required",
			})
			return
		}

		realmID := getRealmID(c)

		page, err := service.GetPageBySlug(slug, realmID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get page",
			})
			return
		}

		// Check if user is authenticated
		userID, isAuthenticated := auth.GetCurrentUser(c)

		// Add authentication status to response
		response := gin.H{
			"page":          page,
			"authenticated": isAuthenticated,
		}

		// Add user info if authenticated
		if isAuthenticated {
			response["user_id"] = userID
			// Add permission flags for authenticated users
			response["can_edit"] = page.CanEdit
			response["can_delete"] = page.CanDelete
		}

		c.JSON(http.StatusOK, response)
	}
}

func listWikiPages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		status := models.WikiPageStatus(c.Query("status"))
		pageType := models.WikiPageType(c.Query("type"))

		// Get realm ID
		realmID := getRealmID(c)

		// Check if user is authenticated
		userID, isAuthenticated := auth.GetCurrentUser(c)

		result, err := service.ListPages(realmID, status, pageType, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to list pages",
			})
			return
		}

		// Add authentication status to response
		response := gin.H{
			"pages":         result.Pages,
			"total":         result.Total,
			"page":          result.Page,
			"limit":         result.Limit,
			"authenticated": isAuthenticated,
		}

		// Add user info if authenticated
		if isAuthenticated {
			response["user_id"] = userID
		}

		c.JSON(http.StatusOK, response)
	}
}

func searchWikiPages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Search query is required",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		realmID := getRealmID(c)

		result, err := service.SearchPages(realmID, query, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to search pages",
			})
			return
		}

		// Check if user is authenticated
		userID, isAuthenticated := auth.GetCurrentUser(c)

		// Add authentication status to response
		response := gin.H{
			"pages":         result.Pages,
			"total":         result.Total,
			"page":          result.Page,
			"limit":         result.Limit,
			"query":         query,
			"authenticated": isAuthenticated,
		}

		// Add user info if authenticated
		if isAuthenticated {
			response["user_id"] = userID
		}

		c.JSON(http.StatusOK, response)
	}
}

func getPagesByCategory(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Param("category")
		if category == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Category is required",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		realmID := getRealmID(c)

		result, err := service.GetPagesByCategory(realmID, category, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get pages by category",
			})
			return
		}

		// Check if user is authenticated
		userID, isAuthenticated := auth.GetCurrentUser(c)

		// Add authentication status to response
		response := gin.H{
			"pages":         result.Pages,
			"total":         result.Total,
			"page":          result.Page,
			"limit":         result.Limit,
			"category":      category,
			"authenticated": isAuthenticated,
		}

		// Add user info if authenticated
		if isAuthenticated {
			response["user_id"] = userID
		}

		c.JSON(http.StatusOK, response)
	}
}

func getPagesByTag(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := c.Param("tag")
		if tag == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag is required",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		realmID := getRealmID(c)

		result, err := service.GetPagesByCategory(realmID, tag, page, limit) // Using category method for now
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get pages by tag",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pages": result.Pages,
			"total": result.Total,
			"page":  result.Page,
			"limit": result.Limit,
			"tag":   tag,
		})
	}
}

func getRandomPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		page, err := service.GetRandomPage(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get random page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page": page,
		})
	}
}

func getRecentChanges(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		realmID := getRealmID(c)

		result, err := service.GetRecentChanges(realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get recent changes",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getOrphanedPages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		result, err := service.GetOrphanedPages(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get orphaned pages",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getWantedPages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		pages, err := service.GetWantedPages(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get wanted pages",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pages": pages,
		})
	}
}

func getDeadEndPages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		result, err := service.GetDeadEndPages(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get dead end pages",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getPageHistory(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page slug is required",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		realmID := getRealmID(c)

		// First get the page to get its ID
		pageResponse, err := service.GetPageBySlug(slug, realmID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get page",
			})
			return
		}

		result, err := service.GetPageHistory(pageResponse.ID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get page history",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func compareRevisions(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		fromStr := c.Param("from")
		toStr := c.Param("to")

		if slug == "" || fromStr == "" || toStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page slug, from version, and to version are required",
			})
			return
		}

		fromVersion, err := strconv.Atoi(fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid from version",
			})
			return
		}

		toVersion, err := strconv.Atoi(toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid to version",
			})
			return
		}

		realmID := getRealmID(c)

		// First get the page to get its ID
		pageResponse, err := service.GetPageBySlug(slug, realmID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get page",
			})
			return
		}

		result, err := service.CompareRevisions(pageResponse.ID, fromVersion, toVersion)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to compare revisions",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// Admin handlers

func createWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateWikiPageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)
		realmID, _ := auth.GetCurrentRealm(c)

		page, err := service.CreatePage(&req, realmID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create page",
			})
			return
		}

		c.JSON(http.StatusCreated, page)
	}
}

func updateWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		var req UpdateWikiPageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		page, err := service.UpdatePage(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page not found",
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
				"error": "Failed to update page",
			})
			return
		}

		c.JSON(http.StatusOK, page)
	}
}

func deleteWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		err := service.DeletePage(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Page deleted successfully",
		})
	}
}

func lockWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		// Get the page
		_, err := service.GetPageByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Page not found",
			})
			return
		}

		// Update the page to lock it
		updateReq := &UpdateWikiPageRequest{
			ChangeNote: "Page locked",
		}

		_, err = service.UpdatePage(id, updateReq, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to lock page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Page locked successfully",
		})
	}
}

func unlockWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		// Update the page to unlock it
		updateReq := &UpdateWikiPageRequest{
			ChangeNote: "Page unlocked",
		}

		_, err := service.UpdatePage(id, updateReq, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unlock page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Page unlocked successfully",
		})
	}
}

func archiveWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		// Update the page to archive it
		updateReq := &UpdateWikiPageRequest{
			Status:     models.WikiPageStatusArchived,
			ChangeNote: "Page archived",
		}

		_, err := service.UpdatePage(id, updateReq, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to archive page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Page archived successfully",
		})
	}
}

func publishWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		// Update the page to publish it
		updateReq := &UpdateWikiPageRequest{
			Status:     models.WikiPageStatusPublished,
			ChangeNote: "Page published",
		}

		_, err := service.UpdatePage(id, updateReq, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to publish page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Page published successfully",
		})
	}
}

func unpublishWikiPage(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		// Update the page to unpublish it
		updateReq := &UpdateWikiPageRequest{
			Status:     models.WikiPageStatusDraft,
			ChangeNote: "Page unpublished",
		}

		_, err := service.UpdatePage(id, updateReq, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unpublish page",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Page unpublished successfully",
		})
	}
}

func createRevision(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID is required",
			})
			return
		}

		var req CreateRevisionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		revision, err := service.CreateRevision(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page not found",
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
				"error": "Failed to create revision",
			})
			return
		}

		c.JSON(http.StatusCreated, revision)
	}
}

func revertToRevision(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		revisionID := c.Param("revisionId")

		if id == "" || revisionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Page ID and revision ID are required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		page, err := service.RevertToRevision(id, revisionID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Page or revision not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to revert to revision",
			})
			return
		}

		c.JSON(http.StatusOK, page)
	}
}

func bulkUpdatePages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would implement bulk update functionality
		// For now, return not implemented
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "Bulk update not implemented yet",
		})
	}
}

func bulkDeletePages(service *WikiService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would implement bulk delete functionality
		// For now, return not implemented
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "Bulk delete not implemented yet",
		})
	}
}

// Helper functions

func getRealmID(c *gin.Context) string {
	// For authenticated routes, get realm ID from auth context
	if realmID, exists := auth.GetCurrentRealm(c); exists {
		return realmID
	}

	// For public wiki routes, return empty string to indicate "any realm"
	// This will be handled in the repository layer
	return models.PUBLIC_REALM_ID
}
