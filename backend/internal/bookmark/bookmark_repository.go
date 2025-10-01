package bookmark

import (
	"strings"
	"time"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// BookmarkRepository handles database operations for bookmarks
type BookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository creates a new bookmark repository
func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

// Create creates a new bookmark
func (r *BookmarkRepository) Create(bookmark *models.Bookmark) error {
	return r.db.Create(bookmark).Error
}

// GetByID retrieves a bookmark by ID
func (r *BookmarkRepository) GetByID(id string) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	err := r.db.Preload("Tags").Where("id = ?", id).First(&bookmark).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

// GetByURL retrieves a bookmark by URL and realm
func (r *BookmarkRepository) GetByURL(url string, realmID string) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	err := r.db.Preload("Tags").Where("url = ? AND realm_id = ?", url, realmID).First(&bookmark).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

// Update updates an existing bookmark
func (r *BookmarkRepository) Update(bookmark *models.Bookmark) error {
	return r.db.Save(bookmark).Error
}

// Delete deletes a bookmark by ID
func (r *BookmarkRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Bookmark{}).Error
}

// List retrieves bookmarks with pagination and filtering
func (r *BookmarkRepository) List(realmID string, offset, limit int, search, categoryID string, tags []string) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	query := r.db.Model(&models.Bookmark{}).Where("realm_id = ?", realmID)

	// Apply search filter
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(url) LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	// Apply category filter
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Apply tag filter
	if len(tags) > 0 {
		subQuery := r.db.Table("bookmark_tags").
			Select("bookmark_id").
			Where("name IN ?", tags).
			Group("bookmark_id").
			Having("COUNT(DISTINCT name) = ?", len(tags))

		query = query.Where("id IN (?)", subQuery)
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with tags preloaded
	err = query.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).Error

	if err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

// GetByCategory retrieves bookmarks by category
func (r *BookmarkRepository) GetByCategory(categoryID string, realmID string, offset, limit int) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	query := r.db.Model(&models.Bookmark{}).
		Where("category_id = ? AND realm_id = ?", categoryID, realmID)

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get results
	err = query.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

// GetByTag retrieves bookmarks by tag
func (r *BookmarkRepository) GetByTag(tagName string, realmID string, offset, limit int) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	subQuery := r.db.Table("bookmark_tags").
		Select("bookmark_id").
		Where("name = ?", tagName)

	query := r.db.Model(&models.Bookmark{}).
		Where("id IN (?) AND realm_id = ?", subQuery, realmID)

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get results
	err = query.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

// GetRecentBookmarks retrieves recently created bookmarks
func (r *BookmarkRepository) GetRecentBookmarks(realmID string, limit int) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	err := r.db.Preload("Tags").
		Where("realm_id = ?", realmID).
		Order("created_at DESC").
		Limit(limit).
		Find(&bookmarks).Error
	return bookmarks, err
}

// SearchBookmarks performs full-text search on bookmarks
func (r *BookmarkRepository) SearchBookmarks(realmID string, query string, offset, limit int) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	searchTerm := "%" + strings.ToLower(query) + "%"
	dbQuery := r.db.Model(&models.Bookmark{}).
		Where("realm_id = ? AND (LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(url) LIKE ?)",
			realmID, searchTerm, searchTerm, searchTerm)

	// Count total
	err := dbQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get results
	err = dbQuery.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

// CreateTag creates a new bookmark tag
func (r *BookmarkRepository) CreateTag(tag *models.BookmarkTag) error {
	return r.db.Create(tag).Error
}

// DeleteTagsByBookmarkID deletes all tags for a bookmark
func (r *BookmarkRepository) DeleteTagsByBookmarkID(bookmarkID string) error {
	return r.db.Where("bookmark_id = ?", bookmarkID).Delete(&models.BookmarkTag{}).Error
}

// GetAllTags retrieves all unique tags for a realm
func (r *BookmarkRepository) GetAllTags(realmID string) ([]string, error) {
	var tags []string
	err := r.db.Table("bookmark_tags").
		Joins("JOIN bookmarks ON bookmark_tags.bookmark_id = bookmarks.id").
		Where("bookmarks.realm_id = ?", realmID).
		Distinct("bookmark_tags.name").
		Pluck("bookmark_tags.name", &tags).Error
	return tags, err
}

// GetPopularTags retrieves most used tags for a realm
func (r *BookmarkRepository) GetPopularTags(realmID string, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("bookmark_tags").
		Select("bookmark_tags.name, COUNT(*) as count").
		Joins("JOIN bookmarks ON bookmark_tags.bookmark_id = bookmarks.id").
		Where("bookmarks.realm_id = ?", realmID).
		Group("bookmark_tags.name").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

// CreateCategory creates a new bookmark category
func (r *BookmarkRepository) CreateCategory(category *models.BookmarkCategory) error {
	return r.db.Create(category).Error
}

// GetCategoryByID retrieves a category by ID
func (r *BookmarkRepository) GetCategoryByID(id int64) (*models.BookmarkCategory, error) {
	var category models.BookmarkCategory
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category
func (r *BookmarkRepository) UpdateCategory(category *models.BookmarkCategory) error {
	return r.db.Save(category).Error
}

// DeleteCategory deletes a category by ID
func (r *BookmarkRepository) DeleteCategory(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.BookmarkCategory{}).Error
}

// ListCategories retrieves all categories
func (r *BookmarkRepository) ListCategories() ([]models.BookmarkCategory, error) {
	var categories []models.BookmarkCategory
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

// GetBookmarkCount returns the total count of bookmarks for a realm
func (r *BookmarkRepository) GetBookmarkCount(realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Bookmark{}).Where("realm_id = ?", realmID).Count(&count).Error
	return count, err
}

// GetBookmarksByDateRange retrieves bookmarks created within a date range
func (r *BookmarkRepository) GetBookmarksByDateRange(realmID string, startDate, endDate time.Time, offset, limit int) ([]models.Bookmark, int64, error) {
	var bookmarks []models.Bookmark
	var total int64

	query := r.db.Model(&models.Bookmark{}).
		Where("realm_id = ? AND created_at >= ? AND created_at <= ?", realmID, startDate, endDate)

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get results
	err = query.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).Error

	return bookmarks, total, err
}
