package task

import (
	"errors"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/database"
	"gorm.io/gorm"
)

// TaskRepository provides data access for Task entities
type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{db: database.GetDB()}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetByID(id string) (*models.Task, error) {
	var task models.Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	if task.ID == "" {
		return errors.New("missing id for update")
	}
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id string) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}

// Search returns tasks under a realm with optional filters and pagination
func (r *TaskRepository) Search(realmID, query, status, tags string, priority, difficulty, page, pageSize int) ([]models.Task, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var tasks []models.Task
	var total int64

	// Enable debug mode for this query to see SQL statements
	q := r.db.Debug().Model(&models.Task{}).Where("realm_id = ?", realmID)

	if status != "" {
		q = q.Where("status = ?", status)
	}
	if tags != "" {
		q = q.Where("tags LIKE ?", "%"+tags+"%")
	}
	if priority > 0 {
		q = q.Where("priority = ?", priority)
	}
	if difficulty > 0 {
		q = q.Where("difficulty = ?", difficulty)
	}
	if query != "" {
		like := "%" + query + "%"
		q = q.Where("name LIKE ? OR description LIKE ? OR tags LIKE ?", like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("schedule_time ASC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetByStatus returns tasks filtered by status
func (r *TaskRepository) GetByStatus(realmID string, status models.TaskStatus, page, pageSize int) ([]models.Task, int64, error) {
	return r.Search(realmID, "", string(status), "", 0, 0, page, pageSize)
}

// GetUpcoming returns tasks scheduled within a time range
func (r *TaskRepository) GetUpcoming(realmID string, limit int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("realm_id = ? AND status IN (?, ?)", realmID, models.TaskStatusPending, models.TaskStatusRunning).
		Order("schedule_time ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// GetOverdue returns tasks past their deadline
func (r *TaskRepository) GetOverdue(realmID string, limit int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("realm_id = ? AND status = ? AND deadline < NOW()", realmID, models.TaskStatusPending).
		Order("deadline ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}
