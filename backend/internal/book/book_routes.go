package book

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
)

// RegisterRoutes registers HTTP endpoints for managing books
func RegisterRoutes(router *gin.Engine, service *BookService, middleware *auth.AuthMiddleware) {
	// Create a specific group for books with authentication requirement
	group := router.Group("/api/v1/books")
	group.Use(middleware.Authenticate())

	// GET /api/v1/books - Search/list books
	group.GET("", func(c *gin.Context) {
		query := c.Query("q")
		author := c.Query("author")
		tags := c.Query("tags")
		page := parseIntDefault(c.Query("page"), 1)
		pageSize := parseIntDefault(c.Query("page_size"), 20)
		realmID, _ := auth.GetCurrentRealm(c)

		result, err := service.SearchBooks(realmID, query, author, tags, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"books": result.Books,
			"total": result.Total,
			"page":  result.Page,
			"limit": result.Limit,
		})
	})

	// GET /api/v1/books/borrowed - Get borrowed books
	group.GET("/borrowed", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		books, err := service.GetBorrowedBooks(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"books": books,
		})
	})

	// GET /api/v1/books/overdue - Get overdue books
	group.GET("/overdue", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		books, err := service.GetOverdueBooks(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"books": books,
		})
	})

	// GET /api/v1/books/available - Get available books
	group.GET("/available", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		books, err := service.GetAvailableBooks(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"books": books,
		})
	})

	// GET /api/v1/books/tags - Get all tags
	group.GET("/tags", func(c *gin.Context) {
		realmID, _ := auth.GetCurrentRealm(c)
		tags, err := service.GetAllTags(realmID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tags": tags,
		})
	})

	// GET /api/v1/books/tag/:tagName - Get books by tag
	group.GET("/tag/:tagName", func(c *gin.Context) {
		tagName := c.Param("tagName")
		realmID, _ := auth.GetCurrentRealm(c)
		books, err := service.GetBooksByTag(realmID, tagName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"books": books,
		})
	})

	// POST /api/v1/books - Create a new book
	group.POST("", func(c *gin.Context) {
		var req CreateBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		realmID, _ := auth.GetCurrentRealm(c)
		userID, _ := auth.GetCurrentUser(c)

		book, err := service.CreateFromInput(&req, realmID, userID)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, book)
	})

	// GET /api/v1/books/:id - Get a specific book
	group.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Book ID is required",
			})
			return
		}

		book, err := service.GetBookByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Book not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, book)
	})

	// PUT /api/v1/books/:id - Update a book
	group.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Book ID is required",
			})
			return
		}

		var req UpdateBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		userID, _ := auth.GetCurrentUser(c)

		book, err := service.UpdateFromInput(id, &req, userID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Book not found",
				})
				return
			}
			if strings.Contains(err.Error(), "validation failed") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, book)
	})

	// DELETE /api/v1/books/:id - Delete a book
	group.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Book ID is required",
			})
			return
		}

		err := service.DeleteBook(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Book not found",
				})
				return
			}
			if strings.Contains(err.Error(), "borrowed") {
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

		c.JSON(http.StatusNoContent, nil)
	})

	// POST /api/v1/books/:id/borrow - Borrow a book
	group.POST("/:id/borrow", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Book ID is required",
			})
			return
		}

		var req struct {
			Deadline string `json:"deadline" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			return
		}

		// Parse deadline
		deadline, err := time.Parse("2006-01-02T15:04:05Z07:00", req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid deadline format. Use RFC3339 format (e.g., 2024-12-31T23:59:59Z)",
			})
			return
		}

		book, err := service.BorrowBook(id, deadline)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Book not found",
				})
				return
			}
			if strings.Contains(err.Error(), "already borrowed") {
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

		c.JSON(http.StatusOK, book)
	})

	// POST /api/v1/books/:id/return - Return a book
	group.POST("/:id/return", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Book ID is required",
			})
			return
		}

		book, err := service.ReturnBook(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Book not found",
				})
				return
			}
			if strings.Contains(err.Error(), "not currently borrowed") {
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

		c.JSON(http.StatusOK, book)
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
