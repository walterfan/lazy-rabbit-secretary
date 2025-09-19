package inbox

import (
	"strings"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// InboxRepository provides data access for InboxItem entities
type InboxRepository struct {
	db *gorm.DB
}

// NewInboxRepository creates a new inbox repository
func NewInboxRepository(db *gorm.DB) *InboxRepository {
	return &InboxRepository{db: db}
}

// Create creates a new inbox item
func (r *InboxRepository) Create(item *models.InboxItem) error {
	return r.db.Create(item).Error
}

// GetByID retrieves an inbox item by ID
func (r *InboxRepository) GetByID(id string) (*models.InboxItem, error) {
	var item models.InboxItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update updates an inbox item
func (r *InboxRepository) Update(item *models.InboxItem) error {
	return r.db.Save(item).Error
}

// Delete soft deletes an inbox item
func (r *InboxRepository) Delete(id string) error {
	return r.db.Delete(&models.InboxItem{}, "id = ?", id).Error
}

// List retrieves inbox items with pagination and filtering
func (r *InboxRepository) List(realmID string, params ListParams) ([]models.InboxItem, int64, error) {
	var items []models.InboxItem
	var total int64

	query := r.db.Model(&models.InboxItem{}).Where("realm_id = ?", realmID)

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Priority != "" {
		query = query.Where("priority = ?", params.Priority)
	}
	if params.Context != "" {
		query = query.Where("context = ?", params.Context)
	}
	if params.SearchQuery != "" {
		searchTerm := "%" + strings.ToLower(params.SearchQuery) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (params.Page - 1) * params.PageSize
	err := query.Order("created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

// GetByStatus retrieves inbox items by status
func (r *InboxRepository) GetByStatus(realmID, status string) ([]models.InboxItem, error) {
	var items []models.InboxItem
	err := r.db.Where("realm_id = ? AND status = ?", realmID, status).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// GetByPriority retrieves inbox items by priority
func (r *InboxRepository) GetByPriority(realmID, priority string) ([]models.InboxItem, error) {
	var items []models.InboxItem
	err := r.db.Where("realm_id = ? AND priority = ?", realmID, priority).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// GetByContext retrieves inbox items by context
func (r *InboxRepository) GetByContext(realmID, context string) ([]models.InboxItem, error) {
	var items []models.InboxItem
	err := r.db.Where("realm_id = ? AND context = ?", realmID, context).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// GetPendingItems retrieves all pending inbox items
func (r *InboxRepository) GetPendingItems(realmID string) ([]models.InboxItem, error) {
	return r.GetByStatus(realmID, "pending")
}

// GetUrgentItems retrieves all urgent inbox items
func (r *InboxRepository) GetUrgentItems(realmID string) ([]models.InboxItem, error) {
	return r.GetByPriority(realmID, "urgent")
}

// UpdateStatus updates the status of an inbox item
func (r *InboxRepository) UpdateStatus(id, status string) error {
	return r.db.Model(&models.InboxItem{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// BulkUpdateStatus updates status for multiple items
func (r *InboxRepository) BulkUpdateStatus(ids []string, status string) error {
	return r.db.Model(&models.InboxItem{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// GetStats retrieves inbox statistics
func (r *InboxRepository) GetStats(realmID string) (map[string]int64, error) {
	stats := make(map[string]int64)

	var count int64

	// Count by status
	statuses := []string{"pending", "processing", "completed", "archived"}
	for _, status := range statuses {
		err := r.db.Model(&models.InboxItem{}).
			Where("realm_id = ? AND status = ?", realmID, status).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		stats["status_"+status] = count
	}

	// Count by priority
	priorities := []string{"urgent", "high", "normal", "low"}
	for _, priority := range priorities {
		err := r.db.Model(&models.InboxItem{}).
			Where("realm_id = ? AND priority = ?", realmID, priority).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		stats["priority_"+priority] = count
	}

	// Total count
	err := r.db.Model(&models.InboxItem{}).
		Where("realm_id = ?", realmID).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	stats["total"] = count

	return stats, nil
}

// ListParams represents parameters for listing inbox items
type ListParams struct {
	Page        int
	PageSize    int
	Status      string
	Priority    string
	Context     string
	SearchQuery string
}
