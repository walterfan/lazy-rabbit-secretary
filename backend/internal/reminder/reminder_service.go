package reminder

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// ReminderService contains business logic for reminders
type ReminderService struct {
	repo *ReminderRepository
}

func NewReminderService(repo *ReminderRepository) *ReminderService {
	return &ReminderService{repo: repo}
}

// CreateReminderRequest defines the allowed input for creating a reminder
type CreateReminderRequest struct {
	Name          string    `json:"name" binding:"required"`
	Content       string    `json:"content" binding:"required"`
	RemindTime    time.Time `json:"remind_time" binding:"required"`
	Tags          string    `json:"tags"`
	RemindMethods string    `json:"remind_methods"` // Comma-separated: email,im,webhook
	RemindTargets string    `json:"remind_targets"` // JSON or comma-separated targets
}

// UpdateReminderRequest defines the allowed input for updating a reminder
type UpdateReminderRequest struct {
	Name          string     `json:"name"`
	Content       string     `json:"content"`
	Status        string     `json:"status"`
	RemindTime    *time.Time `json:"remind_time"`
	Tags          string     `json:"tags"`
	RemindMethods string     `json:"remind_methods"` // Comma-separated: email,im,webhook
	RemindTargets string     `json:"remind_targets"` // JSON or comma-separated targets
}

func (s *ReminderService) CreateFromInput(req CreateReminderRequest, realmID, createdBy string) (*models.Reminder, error) {
	if strings.TrimSpace(realmID) == "" || strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("realm_id and name are required")
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.New("content is required")
	}

	// Validate that remind_time is in the future
	if req.RemindTime.Before(time.Now()) {
		return nil, errors.New("remind_time must be in the future")
	}

	reminder := &models.Reminder{
		ID:            uuid.NewString(),
		RealmID:       realmID,
		Name:          req.Name,
		Content:       req.Content,
		RemindTime:    req.RemindTime,
		Status:        "pending", // default status
		Tags:          req.Tags,
		RemindMethods: req.RemindMethods,
		RemindTargets: req.RemindTargets,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		UpdatedBy:     createdBy,
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(reminder); err != nil {
		return nil, err
	}
	return reminder, nil
}

func (s *ReminderService) UpdateFromInput(id string, req UpdateReminderRequest, updatedBy string) (*models.Reminder, error) {
	// Get existing reminder
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.Status != "" {
		// Validate status values
		if req.Status != "pending" && req.Status != "active" && req.Status != "completed" && req.Status != "cancelled" {
			return nil, errors.New("status must be one of: pending, active, completed, cancelled")
		}
		existing.Status = req.Status
	}
	if req.RemindTime != nil {
		// Only allow future times for pending reminders
		if existing.Status == "pending" && req.RemindTime.Before(time.Now()) {
			return nil, errors.New("remind_time must be in the future for pending reminders")
		}
		existing.RemindTime = *req.RemindTime
	}
	if req.Tags != "" {
		existing.Tags = req.Tags
	}
	if req.RemindMethods != "" {
		existing.RemindMethods = req.RemindMethods
	}
	if req.RemindTargets != "" {
		existing.RemindTargets = req.RemindTargets
	}

	existing.UpdatedBy = updatedBy
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *ReminderService) GetReminder(id string) (*models.Reminder, error) {
	return s.repo.GetByID(id)
}

func (s *ReminderService) DeleteReminder(id string) error {
	// Check if reminder exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// SearchReminders searches reminders with filters and pagination
func (s *ReminderService) SearchReminders(realmID, query, status, tags string, page, pageSize int) ([]models.Reminder, int64, error) {
	return s.repo.Search(realmID, query, status, tags, page, pageSize)
}

// GetRemindersByStatus returns reminders filtered by status
func (s *ReminderService) GetRemindersByStatus(realmID string, status string, page, pageSize int) ([]models.Reminder, int64, error) {
	return s.repo.GetByStatus(realmID, status, page, pageSize)
}

// GetUpcomingReminders returns upcoming reminders
func (s *ReminderService) GetUpcomingReminders(realmID string, limit int) ([]models.Reminder, error) {
	return s.repo.GetUpcoming(realmID, limit)
}

// GetOverdueReminders returns overdue reminders
func (s *ReminderService) GetOverdueReminders(realmID string, limit int) ([]models.Reminder, error) {
	return s.repo.GetOverdue(realmID, limit)
}

// GetRemindersByTimeRange returns reminders within a time range
func (s *ReminderService) GetRemindersByTimeRange(realmID string, startTime, endTime time.Time, page, pageSize int) ([]models.Reminder, int64, error) {
	if startTime.After(endTime) {
		return nil, 0, errors.New("start_time cannot be after end_time")
	}
	return s.repo.GetByTimeRange(realmID, startTime, endTime, page, pageSize)
}

// MarkAsCompleted marks a reminder as completed
func (s *ReminderService) MarkAsCompleted(id, updatedBy string) (*models.Reminder, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existing.Status = "completed"
	existing.UpdatedBy = updatedBy
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// SnoozeReminder postpones a reminder to a new time
func (s *ReminderService) SnoozeReminder(id string, newRemindTime time.Time, updatedBy string) (*models.Reminder, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing.Status == "completed" || existing.Status == "cancelled" {
		return nil, errors.New("cannot snooze completed or cancelled reminders")
	}

	if newRemindTime.Before(time.Now()) {
		return nil, errors.New("new remind_time must be in the future")
	}

	existing.RemindTime = newRemindTime
	existing.UpdatedBy = updatedBy
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}
