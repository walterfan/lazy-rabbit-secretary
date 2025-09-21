package post

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// PostRoutes handles HTTP routes for posts
type PostRoutes struct {
	service *PostService
}

// NewPostRoutes creates new post routes
func NewPostRoutes(db *gorm.DB) *PostRoutes {
	return &PostRoutes{
		service: NewPostService(db),
	}
}

// RegisterRoutes registers post routes with the router
func RegisterRoutes(router *gin.Engine, service *PostService, middleware *auth.AuthMiddleware) {
	// Public routes (no authentication required)
	publicGroup := router.Group("/api/v1/posts")
	{
		publicGroup.GET("/published", listPublishedPosts(service))
		publicGroup.GET("/published/:slug", getPublishedPostBySlug(service))
		publicGroup.GET("/category/:category", getPostsByCategory(service))
		publicGroup.GET("/tag/:tag", getPostsByTag(service))
		publicGroup.GET("/archive/:year", getPostsArchive(service))
		publicGroup.GET("/archive/:year/:month", getPostsArchive(service))
		publicGroup.GET("/popular", getPopularPosts(service))
		publicGroup.GET("/recent", getRecentPosts(service))
		publicGroup.GET("/sticky", getStickyPosts(service))
		publicGroup.GET("/search", searchPublishedPosts(service))
	}

	// Admin routes (authentication required)
	adminGroup := router.Group("/api/v1/admin/posts")
	adminGroup.Use(middleware.Authenticate())
	{
		adminGroup.POST("", createPost(service))
		adminGroup.GET("", listPosts(service))
		adminGroup.GET("/search", searchPosts(service))
		adminGroup.GET("/:id", getPost(service))
		adminGroup.PUT("/:id", updatePost(service))
		adminGroup.DELETE("/:id", deletePost(service))
		adminGroup.POST("/:id/publish", publishPost(service))
		adminGroup.POST("/:id/schedule", schedulePost(service))
		adminGroup.POST("/:id/refine", refinePost(service))
		adminGroup.GET("/author/:authorId", getPostsByAuthor(service))
	}
}

// Public handlers

func listPublishedPosts(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		postType := models.PostType(c.DefaultQuery("type", "post"))

		// Get realm ID (could be from subdomain, header, or default)
		realmID := getRealmID(c)

		result, err := service.ListPublished(realmID, postType, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to list published posts",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getPublishedPostBySlug(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post slug is required",
			})
			return
		}

		realmID := getRealmID(c)

		post, err := service.GetBySlug(slug, realmID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get post",
			})
			return
		}

		// Only return published posts for public access
		if post.Status != models.PostStatusPublished {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func getPostsByCategory(service *PostService) gin.HandlerFunc {
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

		result, err := service.GetByCategory(realmID, category, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get posts by category",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getPostsByTag(service *PostService) gin.HandlerFunc {
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

		result, err := service.GetByTag(realmID, tag, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get posts by tag",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func searchPublishedPosts(service *PostService) gin.HandlerFunc {
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
		postType := models.PostType(c.DefaultQuery("type", "post"))
		realmID := getRealmID(c)

		result, err := service.Search(realmID, query, models.PostStatusPublished, postType, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to search posts",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getPostsArchive(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		yearStr := c.Param("year")
		monthStr := c.Param("month")

		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 2000 || year > 3000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid year",
			})
			return
		}

		month := 0
		if monthStr != "" {
			month, err = strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid month",
				})
				return
			}
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		realmID := getRealmID(c)

		// Use repository directly for archive functionality
		posts, total, err := service.repo.GetArchive(realmID, year, month, (page-1)*limit, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get archive posts",
			})
			return
		}

		result := &PostListResponse{
			Posts: service.toResponseList(posts),
			Total: total,
			Page:  page,
			Limit: limit,
		}

		c.JSON(http.StatusOK, result)
	}
}

func getPopularPosts(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
		postType := models.PostType(c.DefaultQuery("type", "post"))
		realmID := getRealmID(c)

		posts, err := service.repo.GetPopularPosts(realmID, postType, limit, days)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get popular posts",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"posts": service.toResponseList(posts),
		})
	}
}

func getRecentPosts(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		postType := models.PostType(c.DefaultQuery("type", "post"))
		realmID := getRealmID(c)

		posts, err := service.repo.GetRecentPosts(realmID, postType, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get recent posts",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"posts": service.toResponseList(posts),
		})
	}
}

func getStickyPosts(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postType := models.PostType(c.DefaultQuery("type", "post"))
		realmID := getRealmID(c)

		posts, err := service.repo.GetStickyPosts(realmID, postType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get sticky posts",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"posts": service.toResponseList(posts),
		})
	}
}

// Admin handlers

func createPost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		// Determine realm ID based on post status and type
		var realmID string
		if req.Status == "published" && req.Type == "post" {
			// Public posts use PUBLIC_REALM_ID
			realmID = models.PUBLIC_REALM_ID
		} else {
			// Private posts (drafts, pages, etc.) use realm from auth context
			realmID, _ = auth.GetCurrentRealm(c)
		}

		post, err := service.CreateFromInput(&req, realmID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create post",
			})
			return
		}

		c.JSON(http.StatusCreated, post)
	}
}

func getPost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
			return
		}

		post, err := service.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get post",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func updatePost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
			return
		}

		var req UpdatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		post, err := service.UpdateFromInput(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
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
				"error": "Failed to update post",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func refinePost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
			return
		}

		var req RefineRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		post, err := service.RefinePost(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
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
				"error": "Failed to refine post",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func deletePost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
			return
		}

		err := service.Delete(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete post",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post deleted successfully",
		})
	}
}

func listPosts(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		status := models.PostStatus(c.Query("status"))
		postType := models.PostType(c.DefaultQuery("type", ""))

		// Get realm ID from auth context
		realmID, _ := auth.GetCurrentRealm(c)

		result, err := service.List(realmID, status, postType, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to list posts",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func searchPosts(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Search query is required",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		status := models.PostStatus(c.Query("status"))
		postType := models.PostType(c.DefaultQuery("type", ""))

		// Get realm ID from auth context
		realmID, _ := auth.GetCurrentRealm(c)

		result, err := service.Search(realmID, query, status, postType, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to search posts",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func publishPost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		post, err := service.Publish(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
				})
				return
			}
			if strings.Contains(err.Error(), "cannot be published") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to publish post",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func schedulePost(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Post ID is required",
			})
			return
		}

		var req struct {
			ScheduledFor time.Time `json:"scheduled_for" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Get user ID from auth context
		userID, _ := auth.GetCurrentUser(c)

		post, err := service.Schedule(id, req.ScheduledFor, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Post not found",
				})
				return
			}
			if strings.Contains(err.Error(), "cannot be scheduled") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to schedule post",
			})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func getPostsByAuthor(service *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID := c.Param("authorId")
		if authorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Author ID is required",
			})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		status := models.PostStatus(c.Query("status"))

		// Get realm ID from auth context
		realmID, _ := auth.GetCurrentRealm(c)

		offset := (page - 1) * limit
		posts, total, err := service.repo.GetByAuthor(realmID, authorID, status, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get posts by author",
			})
			return
		}

		result := &PostListResponse{
			Posts: service.toResponseList(posts),
			Total: total,
			Page:  page,
			Limit: limit,
		}

		c.JSON(http.StatusOK, result)
	}
}

// Helper functions

func getRealmID(c *gin.Context) string {
	// For authenticated routes, get realm ID from auth context
	if realmID, exists := auth.GetCurrentRealm(c); exists {
		return realmID
	}

	// For public blog routes, return empty string to indicate "any realm"
	// This will be handled in the repository layer
	return models.PUBLIC_REALM_ID
}
