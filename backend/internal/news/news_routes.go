package news

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// NewsHandler handles HTTP requests for news
type NewsHandler struct {
	service *NewsService
}

// NewNewsHandler creates a new news handler
func NewNewsHandler(service *NewsService) *NewsHandler {
	return &NewsHandler{
		service: service,
	}
}

// RegisterRoutes registers all news routes
func (h *NewsHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *auth.AuthMiddleware) {
	// Public routes (no authentication required)
	publicGroup := router.Group("/news")
	publicGroup.Use(authMiddleware.OptionalAuth())
	{
		publicGroup.GET("/", h.getNews)
		publicGroup.GET("/:id", h.getNewsByID)
	}

	// Admin routes (authentication required)
	adminGroup := router.Group("/admin/news")
	adminGroup.Use(authMiddleware.Authenticate())
	{
		adminGroup.GET("/", h.getAllNews)
		adminGroup.POST("/", h.createNews)
		adminGroup.PUT("/:id", h.updateNews)
		adminGroup.DELETE("/:id", h.deleteNews)
	}
}

// Public API handlers

// getNews retrieves news items
func (h *NewsHandler) getNews(c *gin.Context) {
	// Get realm ID from context (set by OptionalAuth middleware)
	realmID, exists := auth.GetCurrentRealm(c)
	if !exists || realmID == "" {
		realmID = models.PUBLIC_REALM_ID
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get user info for personalized response
	userID, isAuthenticated := auth.GetCurrentUser(c)

	// Fetch news
	response, err := h.service.GetByRealm(c.Request.Context(), realmID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add authentication info to response
	result := gin.H{
		"news":          response.News,
		"total":         response.Total,
		"page":          response.Page,
		"limit":         response.Limit,
		"has_more":      response.HasMore,
		"authenticated": isAuthenticated,
	}

	if isAuthenticated {
		result["user_id"] = userID
	}

	c.JSON(http.StatusOK, result)
}

// getNewsByID retrieves a news item by ID (public)
func (h *NewsHandler) getNewsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "News ID is required"})
		return
	}

	news, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"news": news})
}

// Admin API handlers

// getAllNews retrieves all news items (admin)
func (h *NewsHandler) getAllNews(c *gin.Context) {
	realmID, exists := auth.GetCurrentRealm(c)
	if !exists || realmID == "" {
		realmID = models.PUBLIC_REALM_ID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	response, err := h.service.GetByRealm(c.Request.Context(), realmID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// createNews creates a new news item
func (h *NewsHandler) createNews(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set realm ID from context
	req.RealmID = c.GetString("realm_id")
	if req.RealmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Realm ID is required"})
		return
	}

	// Set author from user context
	userID, _ := auth.GetCurrentUser(c)
	req.Author = userID

	news, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"news": news})
}

// updateNews updates a news item
func (h *NewsHandler) updateNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "News ID is required"})
		return
	}

	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	news, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"news": news})
}

// deleteNews deletes a news item
func (h *NewsHandler) deleteNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "News ID is required"})
		return
	}

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News item deleted successfully"})
}
