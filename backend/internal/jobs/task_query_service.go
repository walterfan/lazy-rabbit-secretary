package jobs

import (
	"time"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// TaskQueryService handles database queries for tasks in job processing
// This breaks the import cycle by keeping database operations separate from business logic
type TaskQueryService struct {
	db *gorm.DB
}

// NewTaskQueryService creates a new task query service
func NewTaskQueryService(db *gorm.DB) *TaskQueryService {
	return &TaskQueryService{
		db: db,
	}
}

// FindParentRepeatingTasks retrieves parent repeating tasks that may need new instances
func (tqs *TaskQueryService) FindParentRepeatingTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	err := tqs.db.Where("is_repeating = ? AND parent_task_id IS NULL", true).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// CountExistingInstances counts how many instances exist for a parent task
func (tqs *TaskQueryService) CountExistingInstances(parentTaskID string) (int64, error) {
	var count int64
	err := tqs.db.Model(&models.Task{}).Where("parent_task_id = ?", parentTaskID).Count(&count).Error
	return count, err
}

// CreateTaskInstance creates a new task instance
func (tqs *TaskQueryService) CreateTaskInstance(task *models.Task) error {
	return tqs.db.Create(task).Error
}

// FindTaskInstancesAfterTime finds task instances scheduled after a specific time
func (tqs *TaskQueryService) FindTaskInstancesAfterTime(parentTaskID string, afterTime time.Time) ([]*models.Task, error) {
	var tasks []*models.Task
	err := tqs.db.Where("parent_task_id = ? AND schedule_time > ?", parentTaskID, afterTime).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindTasksDueForReminders finds tasks that are scheduled soon and may need reminders
func (tqs *TaskQueryService) FindTasksDueForReminders(beforeTime time.Time) ([]*models.Task, error) {
	var tasks []*models.Task
	// Convert times to UTC for consistent database queries
	nowUTC := time.Now().UTC()
	beforeTimeUTC := beforeTime.UTC()

	// Find tasks that:
	// 1. Are scheduled between now and the specified time
	// 2. Have generate_reminders = true
	// 3. Don't already have reminders created (check via task_reminders mapping table)
	err := tqs.db.Where(`
		schedule_time > ? AND schedule_time <= ?
		AND generate_reminders = true
		AND reminder_methods != ''
		AND NOT EXISTS (
			SELECT 1 FROM task_reminders tr
			JOIN reminders r ON tr.reminder_id = r.id
			WHERE tr.task_id = tasks.id
			AND r.status IN ('pending', 'sent', 'active')
		)`, nowUTC, beforeTimeUTC).Find(&tasks).Error

	if err != nil {
		return nil, err
	}
	return tasks, nil
}
