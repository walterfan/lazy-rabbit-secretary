package daily

import (
	"strings"
	"time"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// DailyRepository provides data access for DailyChecklistItem entities
type DailyRepository struct {
	db *gorm.DB
}

// NewDailyRepository creates a new daily repository
func NewDailyRepository(db *gorm.DB) *DailyRepository {
	return &DailyRepository{db: db}
}

// Create creates a new daily checklist item
func (r *DailyRepository) Create(item *models.DailyChecklistItem) error {
	return r.db.Create(item).Error
}

// GetByID retrieves a daily checklist item by ID
func (r *DailyRepository) GetByID(id string) (*models.DailyChecklistItem, error) {
	var item models.DailyChecklistItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update updates a daily checklist item
func (r *DailyRepository) Update(item *models.DailyChecklistItem) error {
	return r.db.Save(item).Error
}

// Delete soft deletes a daily checklist item
func (r *DailyRepository) Delete(id string) error {
	return r.db.Delete(&models.DailyChecklistItem{}, "id = ?", id).Error
}

// List retrieves daily checklist items with pagination and filtering
func (r *DailyRepository) List(realmID string, params ListParams) ([]models.DailyChecklistItem, int64, error) {
	var items []models.DailyChecklistItem
	var total int64

	query := r.db.Model(&models.DailyChecklistItem{}).Where("realm_id = ?", realmID)

	// Apply filters
	if !params.Date.IsZero() {
		query = query.Where("date = ?", params.Date.Format("2006-01-02"))
	}
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
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (params.Page - 1) * params.PageSize
	err := query.Order("priority ASC, created_at ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&items).Error

	return items, total, err
}

// GetByDate retrieves daily checklist items for a specific date
func (r *DailyRepository) GetByDate(realmID string, date time.Time) ([]models.DailyChecklistItem, error) {
	var items []models.DailyChecklistItem
	err := r.db.Where("realm_id = ? AND date = ?", realmID, date.Format("2006-01-02")).
		Order("priority ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

// GetByStatus retrieves daily checklist items by status
func (r *DailyRepository) GetByStatus(realmID, status string) ([]models.DailyChecklistItem, error) {
	var items []models.DailyChecklistItem
	err := r.db.Where("realm_id = ? AND status = ?", realmID, status).
		Order("priority ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

// GetByPriority retrieves daily checklist items by priority
func (r *DailyRepository) GetByPriority(realmID, priority string) ([]models.DailyChecklistItem, error) {
	var items []models.DailyChecklistItem
	err := r.db.Where("realm_id = ? AND priority = ?", realmID, priority).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

// GetByContext retrieves daily checklist items by context
func (r *DailyRepository) GetByContext(realmID, context string) ([]models.DailyChecklistItem, error) {
	var items []models.DailyChecklistItem
	err := r.db.Where("realm_id = ? AND context = ?", realmID, context).
		Order("priority ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

// GetPendingItems retrieves all pending daily checklist items
func (r *DailyRepository) GetPendingItems(realmID string) ([]models.DailyChecklistItem, error) {
	return r.GetByStatus(realmID, "pending")
}

// GetInProgressItems retrieves all in-progress daily checklist items
func (r *DailyRepository) GetInProgressItems(realmID string) ([]models.DailyChecklistItem, error) {
	return r.GetByStatus(realmID, "in_progress")
}

// GetCompletedItems retrieves all completed daily checklist items
func (r *DailyRepository) GetCompletedItems(realmID string) ([]models.DailyChecklistItem, error) {
	return r.GetByStatus(realmID, "completed")
}

// UpdateStatus updates the status of a daily checklist item
func (r *DailyRepository) UpdateStatus(id, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// If completing, set completion time
	if status == "completed" {
		now := time.Now()
		updates["completion_time"] = &now
	}

	return r.db.Model(&models.DailyChecklistItem{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateActualTime updates the actual time spent on a task
func (r *DailyRepository) UpdateActualTime(id string, actualTime int) error {
	return r.db.Model(&models.DailyChecklistItem{}).
		Where("id = ?", id).
		Update("actual_time", actualTime).Error
}

// BulkUpdateStatus updates status for multiple items
func (r *DailyRepository) BulkUpdateStatus(ids []string, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// If completing, set completion time
	if status == "completed" {
		now := time.Now()
		updates["completion_time"] = &now
	}

	return r.db.Model(&models.DailyChecklistItem{}).
		Where("id IN ?", ids).
		Updates(updates).Error
}

// GetStats retrieves daily checklist statistics
func (r *DailyRepository) GetStats(realmID string, date time.Time) (map[string]int64, error) {
	stats := make(map[string]int64)

	var count int64

	// Count by status for the specific date
	statuses := []string{"pending", "in_progress", "completed", "cancelled"}
	for _, status := range statuses {
		err := r.db.Model(&models.DailyChecklistItem{}).
			Where("realm_id = ? AND date = ? AND status = ?", realmID, date.Format("2006-01-02"), status).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		stats["status_"+status] = count
	}

	// Count by priority for the specific date
	priorities := []string{"A", "B+", "B", "C", "D"}
	for _, priority := range priorities {
		err := r.db.Model(&models.DailyChecklistItem{}).
			Where("realm_id = ? AND date = ? AND priority = ?", realmID, date.Format("2006-01-02"), priority).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		stats["priority_"+priority] = count
	}

	// Total count for the date
	err := r.db.Model(&models.DailyChecklistItem{}).
		Where("realm_id = ? AND date = ?", realmID, date.Format("2006-01-02")).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	stats["total"] = count

	// Calculate total estimated time
	var totalEstimated int64
	err = r.db.Model(&models.DailyChecklistItem{}).
		Where("realm_id = ? AND date = ?", realmID, date.Format("2006-01-02")).
		Select("COALESCE(SUM(estimated_time), 0)").
		Scan(&totalEstimated).Error
	if err != nil {
		return nil, err
	}
	stats["total_estimated_time"] = totalEstimated

	// Calculate total actual time
	var totalActual int64
	err = r.db.Model(&models.DailyChecklistItem{}).
		Where("realm_id = ? AND date = ?", realmID, date.Format("2006-01-02")).
		Select("COALESCE(SUM(actual_time), 0)").
		Scan(&totalActual).Error
	if err != nil {
		return nil, err
	}
	stats["total_actual_time"] = totalActual

	return stats, nil
}

// GetCompletionRate calculates completion rate for a date
func (r *DailyRepository) GetCompletionRate(realmID string, date time.Time) (float64, error) {
	var total, completed int64

	// Count total items
	err := r.db.Model(&models.DailyChecklistItem{}).
		Where("realm_id = ? AND date = ?", realmID, date.Format("2006-01-02")).
		Count(&total).Error
	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}

	// Count completed items
	err = r.db.Model(&models.DailyChecklistItem{}).
		Where("realm_id = ? AND date = ? AND status = ?", realmID, date.Format("2006-01-02"), "completed").
		Count(&completed).Error
	if err != nil {
		return 0, err
	}

	return float64(completed) / float64(total) * 100, nil
}

// ListParams represents parameters for listing daily checklist items
type ListParams struct {
	Page        int
	PageSize    int
	Date        time.Time
	Status      string
	Priority    string
	Context     string
	SearchQuery string
}
