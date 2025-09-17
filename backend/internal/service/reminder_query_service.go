package service

import (
	"time"

	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"gorm.io/gorm"
)

// ReminderQueryService handles database queries for reminders in job processing
// This breaks the import cycle by keeping database operations separate from business logic
type ReminderQueryService struct {
	db *gorm.DB
}

// NewReminderQueryService creates a new reminder query service
func NewReminderQueryService(db *gorm.DB) *ReminderQueryService {
	return &ReminderQueryService{
		db: db,
	}
}

// FindDueReminders retrieves reminders that are due for processing
func (rqs *ReminderQueryService) FindDueReminders(beforeTime time.Time) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	err := rqs.db.Where("status = ? AND remind_time <= ?", "pending", beforeTime).Find(&reminders).Error
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

// UpdateReminderStatus updates the status of a reminder
func (rqs *ReminderQueryService) UpdateReminderStatus(reminder *models.Reminder, status string) error {
	reminder.Status = status
	reminder.UpdatedAt = time.Now()
	return rqs.db.Save(reminder).Error
}

// GetUserByID retrieves a user by ID
func (rqs *ReminderQueryService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	err := rqs.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateReminder creates a new reminder
func (rqs *ReminderQueryService) CreateReminder(reminder *models.Reminder) error {
	return rqs.db.Create(reminder).Error
}
