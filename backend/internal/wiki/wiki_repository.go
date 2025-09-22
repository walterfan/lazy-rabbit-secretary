package wiki

import (
	"strings"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// WikiRepository handles database operations for wiki pages
type WikiRepository struct {
	db *gorm.DB
}

// NewWikiRepository creates a new wiki repository
func NewWikiRepository(db *gorm.DB) *WikiRepository {
	return &WikiRepository{db: db}
}

// Create creates a new wiki page
func (r *WikiRepository) Create(page *models.WikiPage) error {
	return r.db.Create(page).Error
}

// GetByID retrieves a wiki page by ID
func (r *WikiRepository) GetByID(id string) (*models.WikiPage, error) {
	var page models.WikiPage
	err := r.db.Where("id = ?", id).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetBySlug retrieves a wiki page by slug
func (r *WikiRepository) GetBySlug(slug string, realmID string) (*models.WikiPage, error) {
	var page models.WikiPage
	err := r.db.Where("slug = ? AND realm_id = ?", slug, realmID).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Update updates an existing wiki page
func (r *WikiRepository) Update(page *models.WikiPage) error {
	return r.db.Save(page).Error
}

// Delete soft deletes a wiki page
func (r *WikiRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.WikiPage{}).Error
}

// List retrieves wiki pages with pagination and filtering
func (r *WikiRepository) List(realmID string, status models.WikiPageStatus, pageType models.WikiPageType, page, limit int) ([]models.WikiPage, int64, error) {
	var pages []models.WikiPage
	var total int64

	query := r.db.Model(&models.WikiPage{}).Where("realm_id = ?", realmID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if pageType != "" {
		query = query.Where("type = ?", pageType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&pages).Error
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// Search searches wiki pages by title and content
func (r *WikiRepository) Search(realmID, query string, status models.WikiPageStatus, pageType models.WikiPageType, page, limit int) ([]models.WikiPage, int64, error) {
	var pages []models.WikiPage
	var total int64

	dbQuery := r.db.Model(&models.WikiPage{}).Where("realm_id = ?", realmID)

	// Add search conditions
	searchTerm := "%" + strings.ToLower(query) + "%"
	dbQuery = dbQuery.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(summary) LIKE ?", searchTerm, searchTerm, searchTerm)

	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}

	if pageType != "" {
		dbQuery = dbQuery.Where("type = ?", pageType)
	}

	// Count total
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := dbQuery.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&pages).Error
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// GetByCategory retrieves wiki pages by category
func (r *WikiRepository) GetByCategory(realmID, category string, page, limit int) ([]models.WikiPage, int64, error) {
	var pages []models.WikiPage
	var total int64

	query := r.db.Model(&models.WikiPage{}).Where("realm_id = ? AND categories LIKE ?", realmID, "%"+category+"%")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&pages).Error
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// GetByTag retrieves wiki pages by tag
func (r *WikiRepository) GetByTag(realmID, tag string, page, limit int) ([]models.WikiPage, int64, error) {
	var pages []models.WikiPage
	var total int64

	query := r.db.Model(&models.WikiPage{}).Where("realm_id = ? AND tags LIKE ?", realmID, "%"+tag+"%")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&pages).Error
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// GetRecentPages retrieves recently updated pages
func (r *WikiRepository) GetRecentPages(realmID string, limit int) ([]models.WikiPage, error) {
	var pages []models.WikiPage
	err := r.db.Where("realm_id = ? AND status = ?", realmID, models.WikiPageStatusPublished).
		Order("updated_at DESC").Limit(limit).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// GetRandomPage retrieves a random published page
func (r *WikiRepository) GetRandomPage(realmID string) (*models.WikiPage, error) {
	var page models.WikiPage
	err := r.db.Where("realm_id = ? AND status = ?", realmID, models.WikiPageStatusPublished).
		Order("RANDOM()").First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// GetOrphanedPages retrieves pages with no incoming links
func (r *WikiRepository) GetOrphanedPages(realmID string) ([]models.WikiPage, error) {
	var pages []models.WikiPage

	// This is a simplified implementation - in a real wiki, you'd need to track links
	// For now, we'll consider pages with no categories or tags as potentially orphaned
	err := r.db.Where("realm_id = ? AND status = ? AND (categories = '' OR categories IS NULL) AND (tags = '' OR tags IS NULL)",
		realmID, models.WikiPageStatusPublished).Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// GetWantedPages retrieves pages that are linked but don't exist
func (r *WikiRepository) GetWantedPages(realmID string) ([]string, error) {
	// This would require analyzing content for links and checking if target pages exist
	// For now, return empty slice - this would need more complex implementation
	return []string{}, nil
}

// GetDeadEndPages retrieves pages with no outgoing links
func (r *WikiRepository) GetDeadEndPages(realmID string) ([]models.WikiPage, error) {
	// This would require analyzing content for links
	// For now, return empty slice - this would need more complex implementation
	return []models.WikiPage{}, nil
}

// IncrementViewCount increments the view count for a page
func (r *WikiRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.WikiPage{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// CheckSlugExists checks if a slug already exists for a realm
func (r *WikiRepository) CheckSlugExists(slug, realmID, excludeID string) (bool, error) {
	var count int64
	query := r.db.Model(&models.WikiPage{}).Where("slug = ? AND realm_id = ?", slug, realmID)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPageHierarchy retrieves pages in hierarchical order
func (r *WikiRepository) GetPageHierarchy(realmID string) ([]models.WikiPage, error) {
	var pages []models.WikiPage
	err := r.db.Where("realm_id = ? AND status = ?", realmID, models.WikiPageStatusPublished).
		Order("parent_id ASC, menu_order ASC, title ASC").Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// GetPagesByParent retrieves child pages of a parent
func (r *WikiRepository) GetPagesByParent(realmID, parentID string) ([]models.WikiPage, error) {
	var pages []models.WikiPage
	err := r.db.Where("realm_id = ? AND parent_id = ? AND status = ?", realmID, parentID, models.WikiPageStatusPublished).
		Order("menu_order ASC, title ASC").Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// Revision operations

// CreateRevision creates a new revision
func (r *WikiRepository) CreateRevision(revision *models.WikiRevision) error {
	return r.db.Create(revision).Error
}

// GetRevisions retrieves revisions for a page
func (r *WikiRepository) GetRevisions(pageID string, page, limit int) ([]models.WikiRevision, int64, error) {
	var revisions []models.WikiRevision
	var total int64

	query := r.db.Model(&models.WikiRevision{}).Where("page_id = ?", pageID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("version DESC").Offset(offset).Limit(limit).Find(&revisions).Error
	if err != nil {
		return nil, 0, err
	}

	return revisions, total, nil
}

// GetRevisionByVersion retrieves a specific revision by version number
func (r *WikiRepository) GetRevisionByVersion(pageID string, version int) (*models.WikiRevision, error) {
	var revision models.WikiRevision
	err := r.db.Where("page_id = ? AND version = ?", pageID, version).First(&revision).Error
	if err != nil {
		return nil, err
	}
	return &revision, nil
}

// GetLatestRevision retrieves the latest revision for a page
func (r *WikiRepository) GetLatestRevision(pageID string) (*models.WikiRevision, error) {
	var revision models.WikiRevision
	err := r.db.Where("page_id = ?", pageID).Order("version DESC").First(&revision).Error
	if err != nil {
		return nil, err
	}
	return &revision, nil
}

// GetRecentRevisions retrieves recent revisions across all pages
func (r *WikiRepository) GetRecentRevisions(realmID string, limit int) ([]models.WikiRevision, error) {
	var revisions []models.WikiRevision

	// Join with wiki_pages to filter by realm
	err := r.db.Table("wiki_revisions").
		Select("wiki_revisions.*").
		Joins("JOIN wiki_pages ON wiki_revisions.page_id = wiki_pages.id").
		Where("wiki_pages.realm_id = ?", realmID).
		Order("wiki_revisions.created_at DESC").
		Limit(limit).
		Find(&revisions).Error
	if err != nil {
		return nil, err
	}

	return revisions, nil
}

// DeleteRevision deletes a revision
func (r *WikiRepository) DeleteRevision(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.WikiRevision{}).Error
}

// GetRevisionCount returns the number of revisions for a page
func (r *WikiRepository) GetRevisionCount(pageID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.WikiRevision{}).Where("page_id = ?", pageID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
