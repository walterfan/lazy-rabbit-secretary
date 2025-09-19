package post

import (
	"strings"
	"time"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// PostRepository handles database operations for posts
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create creates a new post
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// GetByID retrieves a post by ID
func (r *PostRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetBySlug retrieves a post by slug
func (r *PostRepository) GetBySlug(slug string, realmID string) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("slug = ? AND realm_id = ?", slug, realmID).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update updates an existing post
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete soft deletes a post
func (r *PostRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Post{}).Error
}

// List retrieves posts with pagination and filtering
func (r *PostRepository) List(realmID string, status models.PostStatus, postType models.PostType, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("realm_id = ? or realm_id = ?", realmID, models.PUBLIC_REALM_ID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if postType != "" {
		query = query.Where("type = ?", postType)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// ListPublished retrieves published posts with pagination
func (r *PostRepository) ListPublished(realmID string, postType models.PostType, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// Build base query for published posts
	query := r.db.Model(&models.Post{}).Where("status = ?", models.PostStatusPublished)
	
	// Add realm filter only if realmID is specified
	if realmID != "" {
		query = query.Where("realm_id = ?", realmID)
	}

	if postType != "" {
		query = query.Where("type = ?", postType)
	}

	// Only show posts that are published and not scheduled for future
	query = query.Where("published_at IS NULL OR published_at <= ?", time.Now())

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with sticky posts first
	err := query.Offset(offset).Limit(limit).Order("is_sticky DESC, published_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Search searches posts by title, content, or excerpt
func (r *PostRepository) Search(realmID string, query string, status models.PostStatus, postType models.PostType, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	// Prepare search query
	searchPattern := "%" + strings.ToLower(query) + "%"

	baseQuery := r.db.Model(&models.Post{}).Where("realm_id = ?", realmID).Where(
		"LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ?",
		searchPattern, searchPattern, searchPattern,
	)

	if status != "" {
		baseQuery = baseQuery.Where("status = ?", status)
	}
	if postType != "" {
		baseQuery = baseQuery.Where("type = ?", postType)
	}

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := baseQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetByCategory retrieves posts by category
func (r *PostRepository) GetByCategory(realmID string, category string, status models.PostStatus, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	categoryPattern := "%" + strings.ToLower(category) + "%"
	
	query := r.db.Model(&models.Post{}).Where("LOWER(categories) LIKE ?", categoryPattern)
	
	// Add realm filter only if realmID is specified
	if realmID != "" {
		query = query.Where("realm_id = ?", realmID)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// Default to published posts for public category browsing
		query = query.Where("status = ?", models.PostStatusPublished)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("is_sticky DESC, published_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetByTag retrieves posts by tag
func (r *PostRepository) GetByTag(realmID string, tag string, status models.PostStatus, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	tagPattern := "%" + strings.ToLower(tag) + "%"
	
	query := r.db.Model(&models.Post{}).Where("LOWER(tags) LIKE ?", tagPattern)
	
	// Add realm filter only if realmID is specified
	if realmID != "" {
		query = query.Where("realm_id = ?", realmID)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// Default to published posts for public tag browsing
		query = query.Where("status = ?", models.PostStatusPublished)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("is_sticky DESC, published_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetByAuthor retrieves posts by author
func (r *PostRepository) GetByAuthor(realmID string, authorID string, status models.PostStatus, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("realm_id = ? AND created_by = ?", realmID, authorID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetChildren retrieves child posts (for hierarchical content like pages)
func (r *PostRepository) GetChildren(parentID string, status models.PostStatus) ([]*models.Post, error) {
	var posts []*models.Post

	query := r.db.Where("parent_id = ?", parentID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("menu_order ASC, title ASC").Find(&posts).Error
	return posts, err
}

// GetScheduledPosts retrieves posts scheduled for publication
func (r *PostRepository) GetScheduledPosts(beforeTime time.Time) ([]*models.Post, error) {
	var posts []*models.Post

	err := r.db.Where("status = ? AND scheduled_for IS NOT NULL AND scheduled_for <= ?",
		models.PostStatusScheduled, beforeTime).Find(&posts).Error

	return posts, err
}

// IncrementViewCount increments the view count for a post
func (r *PostRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateCommentCount updates the comment count for a post
func (r *PostRepository) UpdateCommentCount(postID string, count int64) error {
	return r.db.Model(&models.Post{}).Where("id = ?", postID).UpdateColumn("comment_count", count).Error
}

// SlugExists checks if a slug already exists in the realm
func (r *PostRepository) SlugExists(slug string, realmID string, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&models.Post{}).Where("slug = ? AND realm_id = ?", slug, realmID)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetPopularPosts retrieves posts ordered by view count
func (r *PostRepository) GetPopularPosts(realmID string, postType models.PostType, limit int, days int) ([]*models.Post, error) {
	var posts []*models.Post
	
	query := r.db.Model(&models.Post{}).Where("status = ?", models.PostStatusPublished)
	
	// Add realm filter only if realmID is specified
	if realmID != "" {
		query = query.Where("realm_id = ?", realmID)
	}
	
	if postType != "" {
		query = query.Where("type = ?", postType)
	}
	
	if days > 0 {
		since := time.Now().AddDate(0, 0, -days)
		query = query.Where("published_at >= ?", since)
	}

	err := query.Order("view_count DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

// GetRecentPosts retrieves recent posts
func (r *PostRepository) GetRecentPosts(realmID string, postType models.PostType, limit int) ([]*models.Post, error) {
	var posts []*models.Post
	
	query := r.db.Model(&models.Post{}).Where("status = ?", models.PostStatusPublished)
	
	// Add realm filter only if realmID is specified
	if realmID != "" {
		query = query.Where("realm_id = ?", realmID)
	}
	
	if postType != "" {
		query = query.Where("type = ?", postType)
	}

	err := query.Order("published_at DESC").Limit(limit).Find(&posts).Error
	return posts, err
}

// GetStickyPosts retrieves sticky posts
func (r *PostRepository) GetStickyPosts(realmID string, postType models.PostType) ([]*models.Post, error) {
	var posts []*models.Post
	
	query := r.db.Model(&models.Post{}).Where("status = ? AND is_sticky = ?", 
		models.PostStatusPublished, true)
	
	// Add realm filter only if realmID is specified
	if realmID != "" {
		query = query.Where("realm_id = ?", realmID)
	}
	
	if postType != "" {
		query = query.Where("type = ?", postType)
	}

	err := query.Order("published_at DESC").Find(&posts).Error
	return posts, err
}

// GetArchive retrieves posts for archive view (by year/month)
func (r *PostRepository) GetArchive(realmID string, year int, month int, offset, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	query := r.db.Model(&models.Post{}).Where("realm_id = ? AND status = ?", realmID, models.PostStatusPublished)

	if month > 0 {
		// Specific month
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)
		query = query.Where("published_at >= ? AND published_at < ?", startDate, endDate)
	} else {
		// Entire year
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(1, 0, 0)
		query = query.Where("published_at >= ? AND published_at < ?", startDate, endDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Order("published_at DESC").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
