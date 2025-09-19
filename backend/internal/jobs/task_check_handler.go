package jobs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
)

// TaskCheckHandler implements JobHandler for checking tasks
type TaskCheckHandler struct {
	jobManager *JobManager
}

// Execute checks tasks for reminder generation
func (h *TaskCheckHandler) Execute(params string) error {
	if h.jobManager.taskQueryService == nil {
		h.jobManager.logger.Warn("Task query service not initialized, skipping task check")
		return fmt.Errorf("task query service not initialized")
	}

	h.jobManager.logger.Info("Checking tasks for reminder generation...")

	// Find tasks that are due within the next hour and need reminders
	oneHourLater := time.Now().UTC().Add(1 * time.Hour)

	// Find tasks that are scheduled soon and should generate reminders
	tasks, err := h.jobManager.taskQueryService.FindTasksDueForReminders(oneHourLater)
	if err != nil {
		h.jobManager.logger.Errorf("Failed to find tasks due for reminders: %v", err)
		return fmt.Errorf("failed to find tasks due for reminders: %w", err)
	}

	h.jobManager.logger.Infof("Found %d tasks due for reminders", len(tasks))

	for _, task := range tasks {
		if err := h.generateReminderForTask(task); err != nil {
			h.jobManager.logger.Errorf("Failed to generate reminder for task %s: %v", task.ID, err)
			continue
		}
	}

	h.jobManager.logger.Info("Task check completed successfully")
	return nil
}

// generateReminderForTask generates a reminder for a specific task
func (h *TaskCheckHandler) generateReminderForTask(task *models.Task) error {
	logger := log.GetLogger()

	// Check if the task should generate reminders
	if !task.GenerateReminders {
		return nil
	}

	// Calculate reminder time (task deadline minus advance minutes)
	reminderTime := task.Deadline.Add(-time.Duration(task.ReminderAdvanceMinutes) * time.Minute)

	// Create reminder
	reminder := &models.Reminder{
		ID:            uuid.New().String(),
		Name:          h.formatReminderTitle(task),
		Content:       h.formatReminderContent(task),
		RemindTime:    reminderTime,
		Status:        "pending",
		RealmID:       task.RealmID,
		RemindMethods: task.ReminderMethods,
		RemindTargets: task.ReminderTargets,
		CreatedBy:     task.CreatedBy, // Use the task's CreatedBy (should be user ID)
	}

	// Save reminder using the reminder query service
	if err := h.jobManager.reminderQueryService.CreateReminder(reminder); err != nil {
		logger.Errorf("Failed to create reminder for task %s: %v", task.ID, err)
		return fmt.Errorf("failed to create reminder: %w", err)
	}

	// Create task-reminder mapping
	taskReminder := &models.TaskReminder{
		ID:         uuid.New().String(),
		TaskID:     task.ID,
		ReminderID: reminder.ID,
	}

	// Save the mapping (we'll need to add this method to the query service)
	if err := h.jobManager.reminderQueryService.CreateTaskReminder(taskReminder); err != nil {
		logger.Errorf("Failed to create task-reminder mapping for task %s: %v", task.ID, err)
		return fmt.Errorf("failed to create task-reminder mapping: %w", err)
	}

	logger.Infof("Generated reminder for task: %s (due: %s)", task.Name, task.Deadline.Format("2006-01-02 15:04"))
	return nil
}

// formatReminderTitle formats the reminder title
func (h *TaskCheckHandler) formatReminderTitle(task *models.Task) string {
	return fmt.Sprintf("Reminder: %s", task.Name)
}

// formatReminderContent formats the reminder content with task details
func (h *TaskCheckHandler) formatReminderContent(task *models.Task) string {
	content := fmt.Sprintf("Task: %s\n", task.Name)
	if task.Description != "" {
		content += fmt.Sprintf("Description: %s\n", task.Description)
	}
	content += fmt.Sprintf("Deadline: %s\n", task.Deadline.Format("2006-01-02 15:04"))
	content += fmt.Sprintf("Priority: %d\n", task.Priority)

	if task.Tags != "" {
		content += fmt.Sprintf("Tags: %s\n", h.formatReminderTags(task.Tags))
	}

	return content
}

// formatReminderTags formats task tags for reminder content
func (h *TaskCheckHandler) formatReminderTags(tags string) string {
	// Simple tag formatting - can be enhanced later
	return tags
}
