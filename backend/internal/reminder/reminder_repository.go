package reminder

import (
	"errors"
	"time"

	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/database"
	"gorm.io/gorm"
)

// ReminderRepository provides data access for Reminder entities
type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository() *ReminderRepository {
	return &ReminderRepository{db: database.GetDB()}
}

func (r *ReminderRepository) Create(reminder *models.Reminder) error {
	return r.db.Create(reminder).Error
}

func (r *ReminderRepository) GetByID(id string) (*models.Reminder, error) {
	var reminder models.Reminder
	if err := r.db.First(&reminder, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &reminder, nil
}

func (r *ReminderRepository) Update(reminder *models.Reminder) error {
	if reminder.ID == "" {
		return errors.New("missing id for update")
	}
	return r.db.Save(reminder).Error
}

func (r *ReminderRepository) Delete(id string) error {
	return r.db.Delete(&models.Reminder{}, "id = ?", id).Error
}

// Search returns reminders under a realm with optional filters and pagination
func (r *ReminderRepository) Search(realmID, query, status, tags string, page, pageSize int) ([]models.Reminder, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var reminders []models.Reminder
	var total int64

	// Enable debug mode for this query to see SQL statements
	q := r.db.Debug().Model(&models.Reminder{}).Where("realm_id = ?", realmID)

	if status != "" {
		q = q.Where("status = ?", status)
	}
	if tags != "" {
		q = q.Where("tags LIKE ?", "%"+tags+"%")
	}
	if query != "" {
		like := "%" + query + "%"
		q = q.Where("name LIKE ? OR content LIKE ? OR tags LIKE ?", like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("remind_time ASC").Find(&reminders).Error; err != nil {
		return nil, 0, err
	}

	return reminders, total, nil
}

// GetByStatus returns reminders filtered by status
func (r *ReminderRepository) GetByStatus(realmID string, status string, page, pageSize int) ([]models.Reminder, int64, error) {
	return r.Search(realmID, "", status, "", page, pageSize)
}

// GetUpcoming returns reminders scheduled within a time range
func (r *ReminderRepository) GetUpcoming(realmID string, limit int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("realm_id = ? AND status IN (?, ?)", realmID, "pending", "active").
		Order("remind_time ASC").
		Limit(limit).
		Find(&reminders).Error
	return reminders, err
}

// GetOverdue returns reminders past their remind time
func (r *ReminderRepository) GetOverdue(realmID string, limit int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("realm_id = ? AND status = ? AND remind_time < NOW()", realmID, "pending").
		Order("remind_time ASC").
		Limit(limit).
		Find(&reminders).Error
	return reminders, err
}

// GetByTimeRange returns reminders within a specific time range
func (r *ReminderRepository) GetByTimeRange(realmID string, startTime, endTime time.Time, page, pageSize int) ([]models.Reminder, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var reminders []models.Reminder
	var total int64

	q := r.db.Model(&models.Reminder{}).
		Where("realm_id = ? AND remind_time BETWEEN ? AND ?", realmID, startTime, endTime)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("remind_time ASC").Find(&reminders).Error; err != nil {
		return nil, 0, err
	}

	return reminders, total, nil
}
