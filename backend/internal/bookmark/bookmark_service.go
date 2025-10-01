package bookmark

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// BookmarkService handles business logic for bookmarks
type BookmarkService struct {
	repo *BookmarkRepository
}

// NewBookmarkService creates a new bookmark service
func NewBookmarkService(db *gorm.DB) *BookmarkService {
	return &BookmarkService{
		repo: NewBookmarkRepository(db),
	}
}

// CreateBookmarkRequest represents the request to create a bookmark
type CreateBookmarkRequest struct {
	URL         string   `json:"url" binding:"required" validate:"required,url"`
	Title       string   `json:"title" binding:"required" validate:"required,min=1,max=200"`
	Description string   `json:"description" validate:"max=1000"`
	Tags        []string `json:"tags" validate:"dive,min=1,max=50"`
	CategoryID  string   `json:"category_id"`
}

// UpdateBookmarkRequest represents the request to update a bookmark
type UpdateBookmarkRequest struct {
	URL         string   `json:"url" validate:"url"`
	Title       string   `json:"title" validate:"min=1,max=200"`
	Description string   `json:"description" validate:"max=1000"`
	Tags        []string `json:"tags" validate:"dive,min=1,max=50"`
	CategoryID  string   `json:"category_id"`
}

// BookmarkListRequest represents the request to list bookmarks
type BookmarkListRequest struct {
	Page       int      `json:"page" form:"page"`
	PageSize   int      `json:"page_size" form:"page_size"`
	Search     string   `json:"search" form:"search"`
	CategoryID string   `json:"category_id" form:"category_id"`
	Tags       []string `json:"tags" form:"tags"`
	SortBy     string   `json:"sort_by" form:"sort_by"`
	SortOrder  string   `json:"sort_order" form:"sort_order"`
}

// BookmarkResponse represents a bookmark response
type BookmarkResponse struct {
	ID          string                `json:"id"`
	RealmID     string                `json:"realm_id"`
	URL         string                `json:"url"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Tags        []BookmarkTagResponse `json:"tags"`
	CategoryID  string                `json:"category_id"`
	CreatedBy   string                `json:"created_by"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedBy   string                `json:"updated_by"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// BookmarkTagResponse represents a bookmark tag response
type BookmarkTagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BookmarkListResponse represents the response for listing bookmarks
type BookmarkListResponse struct {
	Bookmarks []BookmarkResponse `json:"bookmarks"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	Pages     int                `json:"pages"`
}

// CreateBookmarkCategoryRequest represents the request to create a category
type CreateBookmarkCategoryRequest struct {
	Name     string `json:"name" binding:"required" validate:"required,min=1,max=100"`
	ParentID *int64 `json:"parent_id"`
}

// UpdateBookmarkCategoryRequest represents the request to update a category
type UpdateBookmarkCategoryRequest struct {
	Name     string `json:"name" validate:"min=1,max=100"`
	ParentID *int64 `json:"parent_id"`
}

// BookmarkCategoryResponse represents a bookmark category response
type BookmarkCategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateBookmark creates a new bookmark
func (s *BookmarkService) CreateBookmark(req *CreateBookmarkRequest, realmID, userID string) (*BookmarkResponse, error) {
	// Check if bookmark with same URL already exists
	existing, err := s.repo.GetByURL(req.URL, realmID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing bookmark: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("bookmark with URL %s already exists", req.URL)
	}

	// Create bookmark
	bookmark := &models.Bookmark{
		ID:          uuid.New().String(),
		RealmID:     realmID,
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create bookmark tags
	for _, tagName := range req.Tags {
		tag := models.BookmarkTag{
			ID:         uuid.New().String(),
			BookmarkID: bookmark.ID,
			Name:       strings.TrimSpace(strings.ToLower(tagName)),
		}
		bookmark.Tags = append(bookmark.Tags, tag)
	}

	err = s.repo.Create(bookmark)
	if err != nil {
		return nil, fmt.Errorf("failed to create bookmark: %w", err)
	}

	return s.toBookmarkResponse(bookmark), nil
}

// GetBookmark retrieves a bookmark by ID
func (s *BookmarkService) GetBookmark(id string) (*BookmarkResponse, error) {
	bookmark, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bookmark not found")
		}
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}

	return s.toBookmarkResponse(bookmark), nil
}

// UpdateBookmark updates an existing bookmark
func (s *BookmarkService) UpdateBookmark(id string, req *UpdateBookmarkRequest, userID string) (*BookmarkResponse, error) {
	bookmark, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bookmark not found")
		}
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}

	// Update fields if provided
	if req.URL != "" {
		// Check if URL conflicts with existing bookmark (excluding current one)
		existing, err := s.repo.GetByURL(req.URL, bookmark.RealmID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check existing bookmark: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("bookmark with URL %s already exists", req.URL)
		}
		bookmark.URL = req.URL
	}
	if req.Title != "" {
		bookmark.Title = req.Title
	}
	if req.Description != "" {
		bookmark.Description = req.Description
	}
	if req.CategoryID != "" {
		bookmark.CategoryID = req.CategoryID
	}

	bookmark.UpdatedBy = userID
	bookmark.UpdatedAt = time.Now()

	// Update tags if provided
	if req.Tags != nil {
		// Delete existing tags
		err = s.repo.DeleteTagsByBookmarkID(bookmark.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing tags: %w", err)
		}

		// Create new tags
		bookmark.Tags = nil
		for _, tagName := range req.Tags {
			tag := models.BookmarkTag{
				ID:         uuid.New().String(),
				BookmarkID: bookmark.ID,
				Name:       strings.TrimSpace(strings.ToLower(tagName)),
			}
			bookmark.Tags = append(bookmark.Tags, tag)
			err = s.repo.CreateTag(&tag)
			if err != nil {
				return nil, fmt.Errorf("failed to create tag: %w", err)
			}
		}
	}

	err = s.repo.Update(bookmark)
	if err != nil {
		return nil, fmt.Errorf("failed to update bookmark: %w", err)
	}

	// Reload bookmark with tags
	updatedBookmark, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload bookmark: %w", err)
	}

	return s.toBookmarkResponse(updatedBookmark), nil
}

// DeleteBookmark deletes a bookmark
func (s *BookmarkService) DeleteBookmark(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("bookmark not found")
		}
		return fmt.Errorf("failed to get bookmark: %w", err)
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	return nil
}

// ListBookmarks lists bookmarks with pagination and filtering
func (s *BookmarkService) ListBookmarks(req *BookmarkListRequest, realmID string) (*BookmarkListResponse, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	bookmarks, total, err := s.repo.List(realmID, offset, req.PageSize, req.Search, req.CategoryID, req.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookmarks: %w", err)
	}

	// Convert to response format
	bookmarkResponses := make([]BookmarkResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		bookmarkResponses[i] = *s.toBookmarkResponse(&bookmark)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &BookmarkListResponse{
		Bookmarks: bookmarkResponses,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Pages:     pages,
	}, nil
}

// GetBookmarksByCategory retrieves bookmarks by category
func (s *BookmarkService) GetBookmarksByCategory(categoryID string, realmID string, page, pageSize int) (*BookmarkListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	bookmarks, total, err := s.repo.GetByCategory(categoryID, realmID, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks by category: %w", err)
	}

	bookmarkResponses := make([]BookmarkResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		bookmarkResponses[i] = *s.toBookmarkResponse(&bookmark)
	}

	pages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &BookmarkListResponse{
		Bookmarks: bookmarkResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Pages:     pages,
	}, nil
}

// GetBookmarksByTag retrieves bookmarks by tag
func (s *BookmarkService) GetBookmarksByTag(tagName string, realmID string, page, pageSize int) (*BookmarkListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	bookmarks, total, err := s.repo.GetByTag(tagName, realmID, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks by tag: %w", err)
	}

	bookmarkResponses := make([]BookmarkResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		bookmarkResponses[i] = *s.toBookmarkResponse(&bookmark)
	}

	pages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &BookmarkListResponse{
		Bookmarks: bookmarkResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Pages:     pages,
	}, nil
}

// SearchBookmarks searches bookmarks
func (s *BookmarkService) SearchBookmarks(query string, realmID string, page, pageSize int) (*BookmarkListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	bookmarks, total, err := s.repo.SearchBookmarks(realmID, query, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search bookmarks: %w", err)
	}

	bookmarkResponses := make([]BookmarkResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		bookmarkResponses[i] = *s.toBookmarkResponse(&bookmark)
	}

	pages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &BookmarkListResponse{
		Bookmarks: bookmarkResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Pages:     pages,
	}, nil
}

// GetRecentBookmarks retrieves recent bookmarks
func (s *BookmarkService) GetRecentBookmarks(realmID string, limit int) ([]BookmarkResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	bookmarks, err := s.repo.GetRecentBookmarks(realmID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent bookmarks: %w", err)
	}

	responses := make([]BookmarkResponse, len(bookmarks))
	for i, bookmark := range bookmarks {
		responses[i] = *s.toBookmarkResponse(&bookmark)
	}

	return responses, nil
}

// GetAllTags retrieves all tags for a realm
func (s *BookmarkService) GetAllTags(realmID string) ([]string, error) {
	return s.repo.GetAllTags(realmID)
}

// GetPopularTags retrieves popular tags for a realm
func (s *BookmarkService) GetPopularTags(realmID string, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return s.repo.GetPopularTags(realmID, limit)
}

// CreateCategory creates a new bookmark category
func (s *BookmarkService) CreateCategory(req *CreateBookmarkCategoryRequest, userID string) (*BookmarkCategoryResponse, error) {
	category := &models.BookmarkCategory{
		Name:      req.Name,
		ParentID:  req.ParentID,
		CreatedBy: userID,
		UpdatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.repo.CreateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return s.toCategoryResponse(category), nil
}

// GetCategory retrieves a category by ID
func (s *BookmarkService) GetCategory(id int64) (*BookmarkCategoryResponse, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return s.toCategoryResponse(category), nil
}

// UpdateCategory updates an existing category
func (s *BookmarkService) UpdateCategory(id int64, req *UpdateBookmarkCategoryRequest, userID string) (*BookmarkCategoryResponse, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}

	category.UpdatedBy = userID
	category.UpdatedAt = time.Now()

	err = s.repo.UpdateCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return s.toCategoryResponse(category), nil
}

// DeleteCategory deletes a category
func (s *BookmarkService) DeleteCategory(id int64) error {
	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("category not found")
		}
		return fmt.Errorf("failed to get category: %w", err)
	}

	err = s.repo.DeleteCategory(id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// ListCategories lists all categories
func (s *BookmarkService) ListCategories() ([]BookmarkCategoryResponse, error) {
	categories, err := s.repo.ListCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	responses := make([]BookmarkCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.toCategoryResponse(&category)
	}

	return responses, nil
}

// GetBookmarkStats retrieves bookmark statistics for a realm
func (s *BookmarkService) GetBookmarkStats(realmID string) (map[string]interface{}, error) {
	total, err := s.repo.GetBookmarkCount(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark count: %w", err)
	}

	tags, err := s.repo.GetAllTags(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	categories, err := s.repo.ListCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return map[string]interface{}{
		"total_bookmarks":  total,
		"total_tags":       len(tags),
		"total_categories": len(categories),
	}, nil
}

// Helper methods

func (s *BookmarkService) toBookmarkResponse(bookmark *models.Bookmark) *BookmarkResponse {
	tags := make([]BookmarkTagResponse, len(bookmark.Tags))
	for i, tag := range bookmark.Tags {
		tags[i] = BookmarkTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return &BookmarkResponse{
		ID:          bookmark.ID,
		RealmID:     bookmark.RealmID,
		URL:         bookmark.URL,
		Title:       bookmark.Title,
		Description: bookmark.Description,
		Tags:        tags,
		CategoryID:  bookmark.CategoryID,
		CreatedBy:   bookmark.CreatedBy,
		CreatedAt:   bookmark.CreatedAt,
		UpdatedBy:   bookmark.UpdatedBy,
		UpdatedAt:   bookmark.UpdatedAt,
	}
}

func (s *BookmarkService) toCategoryResponse(category *models.BookmarkCategory) *BookmarkCategoryResponse {
	return &BookmarkCategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		ParentID:  category.ParentID,
		CreatedBy: category.CreatedBy,
		UpdatedBy: category.UpdatedBy,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
