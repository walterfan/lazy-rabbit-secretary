package image

import (
	"fmt"
	img "image"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/util"
)

// UploadImageRequest represents the request data for image upload
type UploadImageRequest struct {
	Type        string `form:"type" validate:"required"`
	Category    string `form:"category"`
	Description string `form:"description"`
	Tags        string `form:"tags"`
	IsPublic    bool   `form:"is_public"`
	IsShared    bool   `form:"is_shared"`
}

// RegisterRoutes registers image routes with the router
func RegisterRoutes(router *gin.Engine, service *ImageService, middleware *auth.AuthMiddleware) {
	// Public routes (optional authentication - check auth header if present)
	publicGroup := router.Group("/api/v1/images")
	publicGroup.Use(middleware.OptionalAuth())
	{
		// Image access
		publicGroup.GET("/:id", getImageByID(service))
		publicGroup.GET("/", listImages(service))
		publicGroup.GET("/search", searchImages(service))
		publicGroup.GET("/public", listPublicImages(service))
		publicGroup.GET("/shared", listSharedImages(service))

		// Image download
		publicGroup.GET("/:id/download", downloadImage(service))
		publicGroup.GET("/:id/thumbnail", getImageThumbnail(service))

		// Statistics and metadata
		publicGroup.GET("/stats", getImageStats(service))
		publicGroup.GET("/categories", getCategories(service))
		publicGroup.GET("/extensions", getExtensions(service))
		publicGroup.GET("/mime-types", getMimeTypes(service))
	}

	// Admin routes (require authentication)
	adminGroup := router.Group("/api/v1/admin/images")
	adminGroup.Use(middleware.Authenticate())
	{
		// Image management
		adminGroup.POST("/upload", uploadImage(service))
		adminGroup.PUT("/:id", updateImage(service))
		adminGroup.DELETE("/:id", deleteImage(service))

		// Image sharing
		adminGroup.POST("/:id/public", makeImagePublic(service))
		adminGroup.POST("/:id/private", makeImagePrivate(service))
		adminGroup.POST("/:id/share", shareImage(service))
		adminGroup.POST("/:id/unshare", unshareImage(service))
	}
}

// HTTP Handlers

// uploadImage handles image upload
func uploadImage(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get file from form
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "No file uploaded",
				"details": err.Error(),
			})
			return
		}

		// Parse request data
		var req UploadImageRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request data",
				"details": err.Error(),
			})
			return
		}

		// Validate request
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}

		// Get user info from context
		userID, exists := auth.GetCurrentUser(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Upload image
		response, err := service.UploadImage(file, &req, realmID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to upload image",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Image uploaded successfully",
			"data":    response,
		})
	}
}

// getImageByID retrieves an image by ID
func getImageByID(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get image
		response, err := service.GetImage(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get image",
				"details": err.Error(),
			})
			return
		}

		// Increment view count
		if err := service.IncrementViewCount(id); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: failed to increment view count: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// updateImage updates an existing image
func updateImage(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Parse request data
		var req UpdateImageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request data",
				"details": err.Error(),
			})
			return
		}

		// Get user info from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Update image
		response, err := service.UpdateImage(id, &req, userID.(string))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update image",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image updated successfully",
			"data":    response,
		})
	}
}

// deleteImage deletes an image
func deleteImage(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Delete image
		err := service.DeleteImage(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to delete image",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image deleted successfully",
		})
	}
}

// listImages lists images with pagination
func listImages(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		imageType := c.Query("type")
		status := c.Query("status")
		category := c.Query("category")
		userID := c.Query("user_id")

		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// List images
		response, err := service.ListImages(realmID, page, limit, imageType, status, category, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list images",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// searchImages searches images
func searchImages(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		query := c.Query("q")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// Search images
		response, err := service.SearchImages(query, realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to search images",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// listPublicImages lists public images
func listPublicImages(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// List public images
		response, err := service.GetPublicImages(realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list public images",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// listSharedImages lists shared images
func listSharedImages(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// List shared images
		response, err := service.GetSharedImages(realmID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to list shared images",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// downloadImage downloads an image file
func downloadImage(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get image
		image, err := service.DownloadImage(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to download image",
				"details": err.Error(),
			})
			return
		}

		// Check if file exists
		if _, err := os.Stat(image.FilePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Image file not found",
			})
			return
		}

		// Set headers for file download
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+image.OriginalName)
		c.Header("Content-Type", "application/octet-stream")

		// Serve the file
		c.File(image.FilePath)
	}
}

// getImageThumbnail retrieves an image thumbnail
func getImageThumbnail(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.GetLogger()
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get image
		imageResp, err := service.GetImage(id)
		if err != nil {
			logger.Errorf("Failed to get image thumbnail id=%s, error=%v", id, err)
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get image",
				"details": err.Error(),
			})
			return
		}

		// Load the original image from file
		imgFile, err := os.Open(imageResp.FilePath)
		if err != nil {
			logger.Errorf("Failed to open image file id=%s, error=%v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to open image file",
				"details": err.Error(),
			})
			return

		}
		defer imgFile.Close()

		// Decode the image
		srcImg, _, err := img.Decode(imgFile)
		if err != nil {
			logger.Errorf("Failed to decode image id=%s, error=%v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to decode image",
				"details": err.Error(),
			})
			return
		}

		// Create thumbnail (max 200x200 pixels)
		thumbnail := util.CreateThumbnail(srcImg, 200, 200)

		// Create thumbnail filename
		ext := filepath.Ext(imageResp.FilePath)
		thumbnailPath := strings.TrimSuffix(imageResp.FilePath, ext) + "_thumb.png"

		// Save thumbnail to file
		thumbFile, err := os.Create(thumbnailPath)
		if err != nil {
			logger.Errorf("Failed to create thumbnail file id=%s, error=%v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create thumbnail file",
				"details": err.Error(),
			})
			return
		}
		defer thumbFile.Close()

		// Encode thumbnail as PNG (thumbnails are always PNG)
		err = png.Encode(thumbFile, thumbnail)
		if err != nil {
			logger.Errorf("Failed to encode thumbnail id=%s, error=%v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to encode thumbnail",
				"details": err.Error(),
			})
			return
		}

		logger.Infof("Created thumbnail for image id=%s at path=%s", id, thumbnailPath)

		// Return the thumbnail URL for client access
		c.JSON(http.StatusOK, gin.H{
			"thumbnail_url":  imageResp.ThumbnailURL,
			"original_url":   imageResp.URL,
			"thumbnail_path": thumbnailPath,
		})
	}
}

// getImageStats returns storage statistics
func getImageStats(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// Get statistics
		response, err := service.GetImageStats(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get image statistics",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
	}
}

// getCategories returns all categories
func getCategories(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// Get categories
		categories, err := service.GetCategories(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get categories",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": categories,
		})
	}
}

// getExtensions returns all file extensions
func getExtensions(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// Get extensions
		extensions, err := service.GetExtensions(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get extensions",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": extensions,
		})
	}
}

// getMimeTypes returns all MIME types
func getMimeTypes(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get realm ID from context
		realmID, exists := auth.GetCurrentRealm(c)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Realm not specified",
			})
			return
		}

		// Get MIME types
		mimeTypes, err := service.GetMimeTypes(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get MIME types",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": mimeTypes,
		})
	}
}

// makeImagePublic makes an image public
func makeImagePublic(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get user info from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Update image to make it public
		req := &UpdateImageRequest{
			IsPublic: true,
		}

		response, err := service.UpdateImage(id, req, userID.(string))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to make image public",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image made public successfully",
			"data":    response,
		})
	}
}

// makeImagePrivate makes an image private
func makeImagePrivate(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get user info from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Update image to make it private
		req := &UpdateImageRequest{
			IsPublic: false,
			IsShared: false,
		}

		response, err := service.UpdateImage(id, req, userID.(string))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to make image private",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image made private successfully",
			"data":    response,
		})
	}
}

// shareImage enables sharing for an image
func shareImage(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get user info from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Update image to enable sharing
		req := &UpdateImageRequest{
			IsShared: true,
		}

		response, err := service.UpdateImage(id, req, userID.(string))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to share image",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image shared successfully",
			"data":    response,
		})
	}
}

// unshareImage disables sharing for an image
func unshareImage(service *ImageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Image ID is required",
			})
			return
		}

		// Get user info from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// Update image to disable sharing
		req := &UpdateImageRequest{
			IsShared: false,
		}

		response, err := service.UpdateImage(id, req, userID.(string))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Image not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to unshare image",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Image unshared successfully",
			"data":    response,
		})
	}
}
