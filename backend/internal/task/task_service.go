package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
)

// TaskService contains business logic for tasks
type TaskService struct {
	repo *TaskRepository
}

func NewTaskService(repo *TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTaskRequest defines the allowed input for creating a task
type CreateTaskRequest struct {
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description"`
	Priority     *int      `json:"priority"`   // Optional, defaults to 2 if not provided
	Difficulty   *int      `json:"difficulty"` // Optional, defaults to 2 if not provided
	ScheduleTime time.Time `json:"schedule_time" binding:"required"`
	Minutes      int       `json:"minutes" binding:"required,min=1"`
	Deadline     time.Time `json:"deadline" binding:"required"`
	Tags         string    `json:"tags"`
}

// UpdateTaskRequest defines the allowed input for updating a task
type UpdateTaskRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Priority     *int              `json:"priority"`
	Difficulty   *int              `json:"difficulty"`
	Status       models.TaskStatus `json:"status"`
	ScheduleTime *time.Time        `json:"schedule_time"`
	Minutes      *int              `json:"minutes"`
	Deadline     *time.Time        `json:"deadline"`
	Tags         string            `json:"tags"`
}

func (s *TaskService) CreateFromInput(req CreateTaskRequest, realmID, createdBy string) (*models.Task, error) {
	if strings.TrimSpace(realmID) == "" || strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("realm_id and name are required")
	}

	if req.Minutes <= 0 {
		return nil, errors.New("minutes must be greater than 0")
	}

	if req.ScheduleTime.After(req.Deadline) {
		return nil, errors.New("schedule_time cannot be after deadline")
	}

	// Validate priority and difficulty ranges (1-5)
	if req.Priority != nil && (*req.Priority < 1 || *req.Priority > 5) {
		return nil, errors.New("priority must be between 1 and 5")
	}
	if req.Difficulty != nil && (*req.Difficulty < 1 || *req.Difficulty > 5) {
		return nil, errors.New("difficulty must be between 1 and 5")
	}

	// Set default values for priority and difficulty if not provided
	priority := 2 // default (matches GORM default)
	if req.Priority != nil {
		priority = *req.Priority
	}

	difficulty := 2 // default (matches GORM default)
	if req.Difficulty != nil {
		difficulty = *req.Difficulty
	}

	task := &models.Task{
		ID:           uuid.NewString(),
		RealmID:      realmID,
		Name:         req.Name,
		Description:  req.Description,
		Priority:     priority,
		Difficulty:   difficulty,
		Status:       models.TaskStatusPending,
		ScheduleTime: req.ScheduleTime,
		Minutes:      req.Minutes,
		Deadline:     req.Deadline,
		Tags:         req.Tags,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedBy:    createdBy,
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

// Legacy simple create retained for compatibility
func (s *TaskService) CreateTask(task *models.Task) error {
	if task.ID == "" {
		task.ID = uuid.NewString()
	}
	if strings.TrimSpace(task.RealmID) == "" || strings.TrimSpace(task.Name) == "" {
		return errors.New("realm_id and name are required")
	}
	if task.Status == "" {
		task.Status = models.TaskStatusPending
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	return s.repo.Create(task)
}

func (s *TaskService) GetTask(id string) (*models.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) UpdateTask(id string, req UpdateTaskRequest, updatedBy string) (*models.Task, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get existing task
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != nil {
		if *req.Priority < 1 || *req.Priority > 5 {
			return nil, errors.New("priority must be between 1 and 5")
		}
		task.Priority = *req.Priority
	}
	if req.Difficulty != nil {
		if *req.Difficulty < 1 || *req.Difficulty > 5 {
			return nil, errors.New("difficulty must be between 1 and 5")
		}
		task.Difficulty = *req.Difficulty
	}
	if req.Status != "" {
		// Validate status transition
		if err := s.validateStatusTransition(task.Status, req.Status); err != nil {
			return nil, err
		}
		task.Status = req.Status

		// Set timestamps based on status
		now := time.Now()
		switch req.Status {
		case models.TaskStatusRunning:
			if task.StartTime == nil {
				task.StartTime = &now
			}
		case models.TaskStatusCompleted, models.TaskStatusFailed:
			if task.StartTime == nil {
				task.StartTime = &now
			}
			task.EndTime = &now
		}
	}
	if req.ScheduleTime != nil {
		task.ScheduleTime = *req.ScheduleTime
	}
	if req.Minutes != nil && *req.Minutes > 0 {
		task.Minutes = *req.Minutes
	}
	if req.Deadline != nil {
		task.Deadline = *req.Deadline
	}
	if req.Tags != "" {
		task.Tags = req.Tags
	}

	task.UpdatedBy = updatedBy
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

func (s *TaskService) SearchTasks(realmID, query, status, tags string, priority, difficulty, page, pageSize int) ([]models.Task, int64, error) {
	return s.repo.Search(realmID, query, status, tags, priority, difficulty, page, pageSize)
}

func (s *TaskService) GetTasksByStatus(realmID string, status models.TaskStatus, page, pageSize int) ([]models.Task, int64, error) {
	return s.repo.GetByStatus(realmID, status, page, pageSize)
}

func (s *TaskService) GetUpcomingTasks(realmID string, limit int) ([]models.Task, error) {
	return s.repo.GetUpcoming(realmID, limit)
}

func (s *TaskService) GetOverdueTasks(realmID string, limit int) ([]models.Task, error) {
	return s.repo.GetOverdue(realmID, limit)
}

// StartTask marks a task as running
func (s *TaskService) StartTask(id string, startedBy string) (*models.Task, error) {
	return s.UpdateTask(id, UpdateTaskRequest{
		Status: models.TaskStatusRunning,
	}, startedBy)
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(id string, completedBy string) (*models.Task, error) {
	return s.UpdateTask(id, UpdateTaskRequest{
		Status: models.TaskStatusCompleted,
	}, completedBy)
}

// FailTask marks a task as failed
func (s *TaskService) FailTask(id string, failedBy string) (*models.Task, error) {
	return s.UpdateTask(id, UpdateTaskRequest{
		Status: models.TaskStatusFailed,
	}, failedBy)
}

// --- helpers ---

func (s *TaskService) validateStatusTransition(currentStatus, newStatus models.TaskStatus) error {
	// Define valid status transitions
	validTransitions := map[models.TaskStatus][]models.TaskStatus{
		models.TaskStatusPending: {
			models.TaskStatusRunning,
			models.TaskStatusFailed,
		},
		models.TaskStatusRunning: {
			models.TaskStatusCompleted,
			models.TaskStatusFailed,
		},
		models.TaskStatusCompleted: {
			// No transitions from completed
		},
		models.TaskStatusFailed: {
			models.TaskStatusPending, // Allow retry
		},
	}

	allowed, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	for _, allowedStatus := range allowed {
		if newStatus == allowedStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}
