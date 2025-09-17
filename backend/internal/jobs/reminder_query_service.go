package jobs

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
	// Convert beforeTime to UTC for consistent database queries
	//beforeTimeUTC := beforeTime.UTC()
	err := rqs.db.Where("status != ? AND remind_time <= ?", "completed", beforeTime).Find(&reminders).Error
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

// UpdateReminderStatus updates the status of a reminder
func (rqs *ReminderQueryService) UpdateReminderStatus(reminder *models.Reminder, status string, updatedBy string) error {
	reminder.Status = status
	reminder.UpdatedBy = updatedBy
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

// CreateTaskReminder creates a new task-reminder mapping
func (rqs *ReminderQueryService) CreateTaskReminder(taskReminder *models.TaskReminder) error {
	return rqs.db.Create(taskReminder).Error
}

// FindRemindersByTaskID retrieves all reminders associated with a specific task
func (rqs *ReminderQueryService) FindRemindersByTaskID(taskID string) ([]*models.Reminder, error) {
	var reminders []*models.Reminder
	err := rqs.db.Table("reminders").
		Joins("JOIN task_reminders ON reminders.id = task_reminders.reminder_id").
		Where("task_reminders.task_id = ?", taskID).
		Find(&reminders).Error
	if err != nil {
		return nil, err
	}
	return reminders, nil
}
