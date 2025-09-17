package service

import (
	"time"

	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
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
