package task

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/internal/reminder"
)

// TaskService contains business logic for tasks
type TaskService struct {
	repo            *TaskRepository
	reminderService *reminder.ReminderService
}

func NewTaskService(repo *TaskRepository, reminderService *reminder.ReminderService) *TaskService {
	return &TaskService{
		repo:            repo,
		reminderService: reminderService,
	}
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

	// Repeat task fields
	IsRepeating      bool       `json:"is_repeating"`
	RepeatPattern    string     `json:"repeat_pattern"`      // daily, weekly, monthly, yearly
	RepeatInterval   int        `json:"repeat_interval"`     // every N days/weeks/months
	RepeatDaysOfWeek string     `json:"repeat_days_of_week"` // comma-separated: mon,tue,wed
	RepeatDayOfMonth int        `json:"repeat_day_of_month"` // for monthly: 1-31, 0=same day
	RepeatEndDate    *time.Time `json:"repeat_end_date"`     // when to stop repeating
	RepeatCount      int        `json:"repeat_count"`        // max occurrences, 0=infinite

	// Reminder generation settings
	GenerateReminders      bool   `json:"generate_reminders"`
	ReminderAdvanceMinutes int    `json:"reminder_advance_minutes"` // remind N minutes before task
	ReminderMethods        string `json:"reminder_methods"`         // email,webhook
	ReminderTargets        string `json:"reminder_targets"`         // notification targets
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

	// Validate repeat settings
	if req.IsRepeating {
		if err := s.validateRepeatSettings(req); err != nil {
			return nil, err
		}
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

		// Repeat fields
		IsRepeating:      req.IsRepeating,
		RepeatPattern:    req.RepeatPattern,
		RepeatInterval:   req.RepeatInterval,
		RepeatDaysOfWeek: req.RepeatDaysOfWeek,
		RepeatDayOfMonth: req.RepeatDayOfMonth,
		RepeatEndDate:    req.RepeatEndDate,
		RepeatCount:      req.RepeatCount,

		// Reminder fields
		GenerateReminders:      req.GenerateReminders,
		ReminderAdvanceMinutes: req.ReminderAdvanceMinutes,
		ReminderMethods:        req.ReminderMethods,
		ReminderTargets:        req.ReminderTargets,

		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedBy: createdBy,
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	// If this is a repeating task, generate initial instances
	if task.IsRepeating {
		instances, err := s.GenerateTaskInstances(task, 5) // Generate first 5 instances
		if err != nil {
			// Log error but don't fail the creation of the parent task
			// In a real system, you might want to use a logger here
			fmt.Printf("Warning: Failed to generate task instances: %v\n", err)
		} else {
			fmt.Printf("Generated %d task instances for repeating task: %s\n", len(instances), task.Name)

			// Generate reminders for each instance if requested
			if task.GenerateReminders {
				for _, instance := range instances {
					if err := s.generateReminderForTask(instance); err != nil {
						fmt.Printf("Warning: Failed to generate reminder for task instance %s: %v\n", instance.ID, err)
					}
				}
			}
		}
	} else if task.GenerateReminders {
		// For non-repeating tasks, generate a single reminder
		if err := s.generateReminderForTask(task); err != nil {
			fmt.Printf("Warning: Failed to generate reminder for task %s: %v\n", task.ID, err)
		}
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

// validateRepeatSettings validates repeat task configuration
func (s *TaskService) validateRepeatSettings(req CreateTaskRequest) error {
	validPatterns := map[string]bool{
		"daily": true, "weekly": true, "monthly": true, "yearly": true,
	}

	if !validPatterns[req.RepeatPattern] {
		return errors.New("repeat_pattern must be one of: daily, weekly, monthly, yearly")
	}

	if req.RepeatInterval < 1 {
		return errors.New("repeat_interval must be at least 1")
	}

	if req.RepeatCount < 0 {
		return errors.New("repeat_count cannot be negative")
	}

	if req.RepeatEndDate != nil && req.RepeatEndDate.Before(req.ScheduleTime) {
		return errors.New("repeat_end_date cannot be before schedule_time")
	}

	// Validate weekly settings
	if req.RepeatPattern == "weekly" && req.RepeatDaysOfWeek != "" {
		validDays := map[string]bool{
			"sun": true, "mon": true, "tue": true, "wed": true,
			"thu": true, "fri": true, "sat": true,
		}

		days := strings.Split(req.RepeatDaysOfWeek, ",")
		for _, day := range days {
			day = strings.TrimSpace(strings.ToLower(day))
			if !validDays[day] {
				return fmt.Errorf("invalid day of week: %s. Use: sun,mon,tue,wed,thu,fri,sat", day)
			}
		}
	}

	// Validate monthly settings
	if req.RepeatPattern == "monthly" && req.RepeatDayOfMonth != 0 {
		if req.RepeatDayOfMonth < 1 || req.RepeatDayOfMonth > 31 {
			return errors.New("repeat_day_of_month must be between 1 and 31")
		}
	}

	// Validate reminder settings
	if req.GenerateReminders {
		if req.ReminderMethods == "" {
			return errors.New("reminder_methods is required when generate_reminders is true")
		}

		validMethods := map[string]bool{"email": true, "webhook": true}
		methods := strings.Split(req.ReminderMethods, ",")
		for _, method := range methods {
			method = strings.TrimSpace(strings.ToLower(method))
			if !validMethods[method] {
				return fmt.Errorf("invalid reminder method: %s. Use: email, webhook", method)
			}
		}

		if req.ReminderAdvanceMinutes < 0 {
			return errors.New("reminder_advance_minutes cannot be negative")
		}
	}

	return nil
}

// GenerateTaskInstances creates task instances for a repeating task
func (s *TaskService) GenerateTaskInstances(parentTask *models.Task, maxInstances int) ([]*models.Task, error) {
	if !parentTask.IsParentTask() {
		return nil, errors.New("task is not a repeating parent task")
	}

	var instances []*models.Task
	currentDate := parentTask.ScheduleTime

	for i := 0; i < maxInstances; i++ {
		nextDate := parentTask.GetNextOccurrence(currentDate)
		if nextDate == nil {
			break // No more occurrences
		}

		// Calculate deadline based on the original task's duration
		duration := parentTask.Deadline.Sub(parentTask.ScheduleTime)
		newDeadline := nextDate.Add(duration)

		instance := &models.Task{
			ID:           uuid.NewString(),
			RealmID:      parentTask.RealmID,
			Name:         parentTask.Name,
			Description:  parentTask.Description,
			Priority:     parentTask.Priority,
			Difficulty:   parentTask.Difficulty,
			Status:       models.TaskStatusPending,
			ScheduleTime: *nextDate,
			Minutes:      parentTask.Minutes,
			Deadline:     newDeadline,
			Tags:         parentTask.Tags,

			// Mark as instance
			IsRepeating:  false, // Instances are not repeating
			ParentTaskID: &parentTask.ID,

			CreatedBy: parentTask.CreatedBy,
			CreatedAt: time.Now(),
			UpdatedBy: parentTask.CreatedBy,
			UpdatedAt: time.Now(),
		}

		if err := s.repo.Create(instance); err != nil {
			return instances, fmt.Errorf("failed to create task instance: %w", err)
		}

		instances = append(instances, instance)
		currentDate = *nextDate
	}

	// Update parent task instance count
	parentTask.InstanceCount += len(instances)
	if err := s.repo.Update(parentTask); err != nil {
		return instances, fmt.Errorf("failed to update parent task instance count: %w", err)
	}

	return instances, nil
}

// GenerateRemindersForTask creates reminders for a task (public method)
func (s *TaskService) GenerateRemindersForTask(task *models.Task) error {
	return s.generateReminderForTask(task)
}

// generateReminderForTask creates a reminder for a task (internal method)
func (s *TaskService) generateReminderForTask(task *models.Task) error {
	if !task.ShouldGenerateReminders() {
		return nil
	}

	if s.reminderService == nil {
		return errors.New("reminder service not available")
	}

	// Calculate reminder time (task schedule time minus advance minutes)
	reminderTime := task.ScheduleTime.Add(-time.Duration(task.ReminderAdvanceMinutes) * time.Minute)

	// Don't create reminders for past times
	if reminderTime.Before(time.Now()) {
		return fmt.Errorf("skipping reminder creation for task %s - reminder time is in the past", task.ID)
	}

	// Create reminder request
	reminderReq := reminder.CreateReminderRequest{
		Name:          fmt.Sprintf("Task Reminder: %s", task.Name),
		Content:       s.formatReminderContent(task),
		RemindTime:    reminderTime,
		Tags:          s.formatReminderTags(task),
		RemindMethods: task.ReminderMethods,
		RemindTargets: task.ReminderTargets,
	}

	// Create the reminder
	createdReminder, err := s.reminderService.CreateFromInput(reminderReq, task.RealmID, task.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to create reminder for task %s: %w", task.ID, err)
	}

	fmt.Printf("Successfully created reminder %s for task %s\n", createdReminder.ID, task.ID)
	return nil
}

// formatReminderContent creates a formatted reminder message for a task
func (s *TaskService) formatReminderContent(task *models.Task) string {
	content := fmt.Sprintf(`You have an upcoming task scheduled:

📋 Task: %s
📝 Description: %s
⏰ Scheduled: %s
⏱️  Duration: %d minutes
📅 Deadline: %s
🎯 Priority: %d/5
🔧 Difficulty: %d/5`,
		task.Name,
		task.Description,
		task.ScheduleTime.Format("2006-01-02 15:04:05"),
		task.Minutes,
		task.Deadline.Format("2006-01-02 15:04:05"),
		task.Priority,
		task.Difficulty)

	if task.Tags != "" {
		content += fmt.Sprintf("\n🏷️  Tags: %s", task.Tags)
	}

	if task.IsTaskInstance() {
		content += "\n\n🔄 This is a recurring task instance."
	}

	content += "\n\n---\nGenerated automatically by Lazy Rabbit Secretary System"

	return content
}

// formatReminderTags creates tags for the reminder based on the task
func (s *TaskService) formatReminderTags(task *models.Task) string {
	tags := []string{"task", "auto-generated"}

	if task.IsTaskInstance() {
		tags = append(tags, "recurring")
	}

	if task.Tags != "" {
		// Add original task tags
		taskTags := strings.Split(task.Tags, ",")
		for _, tag := range taskTags {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	return strings.Join(tags, ",")
}
