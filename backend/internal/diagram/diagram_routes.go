package diagram

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
)

// RegisterRoutes registers diagram routes with the router
func RegisterRoutes(router *gin.Engine, service *DiagramService, middleware *auth.AuthMiddleware) {
	// Public routes (optional authentication - check auth header if present)
	publicGroup := router.Group("/api/v1/diagrams")
	publicGroup.Use(middleware.OptionalAuth())
	{
		// Diagram access
		publicGroup.GET("/:id", getDiagram(service))
		publicGroup.GET("/", listDiagrams(service))
		publicGroup.GET("/search", searchDiagrams(service))
		publicGroup.GET("/public", getPublicDiagrams(service))
		publicGroup.GET("/shared", getSharedDiagrams(service))

		// Diagram generation (draw)
		publicGroup.POST("/draw", drawDiagram(service))

		// Image access
		publicGroup.GET("/:id/images", getDiagramImages(service))
		publicGroup.GET("/:id/images/:imageId", getImage(service))
		publicGroup.GET("/:id/images/:imageId/thumbnail", getImageThumbnail(service))

		// Tag access
		publicGroup.GET("/tags", getAllTags(service))
		publicGroup.GET("/tags/:id", getTag(service))
	}

	// Admin routes (require authentication)
	adminGroup := router.Group("/api/v1/admin/diagrams")
	adminGroup.Use(middleware.Authenticate())
	{
		// Diagram management
		adminGroup.POST("/", createDiagram(service))
		adminGroup.PUT("/:id", updateDiagram(service))
		adminGroup.DELETE("/:id", deleteDiagram(service))
		adminGroup.POST("/:id/publish", publishDiagram(service))
		adminGroup.POST("/:id/archive", archiveDiagram(service))
		adminGroup.POST("/:id/public", makeDiagramPublic(service))
		adminGroup.POST("/:id/private", makeDiagramPrivate(service))
		adminGroup.POST("/:id/share", shareDiagram(service))
		adminGroup.POST("/:id/unshare", unshareDiagram(service))

		// Image management
		adminGroup.POST("/:id/images", createImage(service))
		adminGroup.PUT("/:id/images/:imageId", updateImage(service))
		adminGroup.DELETE("/:id/images/:imageId", deleteImage(service))
		adminGroup.POST("/:id/images/:imageId/primary", setPrimaryImage(service))

		// Tag management
		adminGroup.POST("/tags", createTag(service))
		adminGroup.PUT("/tags/:id", updateTag(service))
		adminGroup.DELETE("/tags/:id", deleteTag(service))
	}
}

// Helper function to get realm ID from context
func getRealmID(c *gin.Context) string {
	realmID := c.GetHeader("X-Realm-ID")
	if realmID == "" {
		realmID = "default" // Default realm
	}
	return realmID
}

// Helper function to get current user ID
func getCurrentUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}

// Public handlers

func getDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		// Check if user is authenticated
		userID, isAuthenticated := auth.GetCurrentUser(c)

		// Get diagram with images
		diagram, err := service.GetDiagram(id, true)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get diagram",
			})
			return
		}

		// Check permissions
		if !diagram.Public && !isAuthenticated {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}

		// Increment view count if authenticated
		if isAuthenticated {
			service.IncrementViewCount(id)
		}

		// Add authentication status to response
		response := gin.H{
			"diagram":       diagram,
			"authenticated": isAuthenticated,
		}

		// Add user info if authenticated
		if isAuthenticated {
			response["user_id"] = userID
			// Add permission flags for authenticated users
			response["can_edit"] = diagram.CreatedBy == userID
			response["can_delete"] = diagram.CreatedBy == userID
		}

		c.JSON(http.StatusOK, response)
	}
}

func listDiagrams(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		diagramType := c.Query("type")
		scriptType := c.Query("script_type")
		status := c.Query("status")
		userID := c.Query("user_id")

		// Check if user is authenticated
		_, isAuthenticated := auth.GetCurrentUser(c)

		// If not authenticated, only show public diagrams
		if !isAuthenticated {
			diagrams, err := service.GetPublicDiagrams(realmID, page, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to get diagrams",
				})
				return
			}
			c.JSON(http.StatusOK, diagrams)
			return
		}

		// Get diagrams based on filters
		diagrams, err := service.ListDiagrams(realmID, page, limit, diagramType, scriptType, status, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get diagrams",
			})
			return
		}

		c.JSON(http.StatusOK, diagrams)
	}
}

func searchDiagrams(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)
		query := c.Query("q")

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Check if user is authenticated
		_, isAuthenticated := auth.GetCurrentUser(c)

		// If not authenticated, only search public diagrams
		if !isAuthenticated {
			// TODO: Implement public diagram search
			c.JSON(http.StatusOK, gin.H{
				"diagrams": []DiagramResponse{},
				"total":    0,
				"page":     page,
				"limit":    limit,
				"has_more": false,
			})
			return
		}

		// Search diagrams
		diagrams, err := service.SearchDiagrams(query, realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to search diagrams",
			})
			return
		}

		c.JSON(http.StatusOK, diagrams)
	}
}

func getPublicDiagrams(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Get public diagrams
		diagrams, err := service.GetPublicDiagrams(realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get public diagrams",
			})
			return
		}

		c.JSON(http.StatusOK, diagrams)
	}
}

func getSharedDiagrams(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)

		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Get shared diagrams
		diagrams, err := service.GetSharedDiagrams(realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get shared diagrams",
			})
			return
		}

		c.JSON(http.StatusOK, diagrams)
	}
}

func drawDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DrawDiagramRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Generate diagram
		result, err := service.DrawDiagram(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate diagram",
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func getDiagramImages(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		diagramID := c.Param("id")
		if diagramID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		// Get images for diagram
		images, err := service.GetImagesByDiagramID(diagramID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get images",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"images": images,
		})
	}
}

func getImage(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("imageId")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get image
		image, err := service.GetImage(imageID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get image",
			})
			return
		}

		c.JSON(http.StatusOK, image)
	}
}

func getImageThumbnail(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("imageId")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get image
		image, err := service.GetImage(imageID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get image",
			})
			return
		}

		// TODO: Implement thumbnail generation
		// For now, return the same image
		c.JSON(http.StatusOK, gin.H{
			"thumbnail_url": image.ThumbnailURL,
		})
	}
}

func getAllTags(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

		// Get all tags
		tags, err := service.GetAllTags(page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get tags",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tags": tags,
		})
	}
}

func getTag(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID := c.Param("id")
		if tagID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag ID is required",
			})
			return
		}

		// Get tag
		tag, err := service.GetTag(tagID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Tag not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get tag",
			})
			return
		}

		c.JSON(http.StatusOK, tag)
	}
}

// Admin handlers

func createDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		realmID := getRealmID(c)
		userID := getCurrentUserID(c)

		var req CreateDiagramRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Create diagram
		diagram, err := service.CreateDiagram(&req, realmID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create diagram",
			})
			return
		}

		c.JSON(http.StatusCreated, diagram)
	}
}

func updateDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		var req UpdateDiagramRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Update diagram
		diagram, err := service.UpdateDiagram(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update diagram",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

func deleteDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		// Delete diagram
		err := service.DeleteDiagram(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete diagram",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Diagram deleted successfully",
		})
	}
}

func publishDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		// Publish diagram
		diagram, err := service.PublishDiagram(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to publish diagram",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

func archiveDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		// Archive diagram
		diagram, err := service.ArchiveDiagram(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to archive diagram",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

func makeDiagramPublic(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		// Make diagram public
		diagram, err := service.MakePublic(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to make diagram public",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

func makeDiagramPrivate(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		// Make diagram private
		diagram, err := service.MakePrivate(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to make diagram private",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

func shareDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		// Share diagram
		diagram, err := service.ShareDiagram(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to share diagram",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

func unshareDiagram(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		// Unshare diagram
		diagram, err := service.UnshareDiagram(id, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Diagram not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unshare diagram",
			})
			return
		}

		c.JSON(http.StatusOK, diagram)
	}
}

// Image handlers

func createImage(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		diagramID := c.Param("id")
		if diagramID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		var req CreateImageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Set diagram ID from URL
		req.DiagramID = diagramID

		// Create image
		image, err := service.CreateImage(&req, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create image",
			})
			return
		}

		c.JSON(http.StatusCreated, image)
	}
}

func updateImage(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("imageId")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		userID := getCurrentUserID(c)

		var req UpdateImageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Update image
		image, err := service.UpdateImage(imageID, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update image",
			})
			return
		}

		c.JSON(http.StatusOK, image)
	}
}

func deleteImage(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		imageID := c.Param("imageId")
		if imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Delete image
		err := service.DeleteImage(imageID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete image",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image deleted successfully",
		})
	}
}

func setPrimaryImage(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		diagramID := c.Param("id")
		imageID := c.Param("imageId")
		if diagramID == "" || imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Diagram ID and Image ID are required",
			})
			return
		}

		// Set primary image
		err := service.SetPrimaryImage(diagramID, imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to set primary image",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Primary image set successfully",
		})
	}
}

// Tag handlers

func createTag(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Create tag
		tag, err := service.CreateTag(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create tag",
			})
			return
		}

		c.JSON(http.StatusCreated, tag)
	}
}

func updateTag(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID := c.Param("id")
		if tagID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag ID is required",
			})
			return
		}

		var req UpdateTagRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Update tag
		tag, err := service.UpdateTag(tagID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Tag not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update tag",
			})
			return
		}

		c.JSON(http.StatusOK, tag)
	}
}

func deleteTag(service *DiagramService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID := c.Param("id")
		if tagID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag ID is required",
			})
			return
		}

		// Delete tag
		err := service.DeleteTag(tagID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Tag not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete tag",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Tag deleted successfully",
		})
	}
}
