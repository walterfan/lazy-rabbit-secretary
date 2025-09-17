package jobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
)

// TaskRemindHandler implements JobHandler for processing reminders
type TaskRemindHandler struct {
	jobManager *JobManager
}

// Execute processes due reminders and sends notifications
func (h *TaskRemindHandler) Execute(params string) error {
	if h.jobManager.reminderQueryService == nil {
		h.jobManager.logger.Warn("Reminder query service not initialized, skipping reminder check")
		return fmt.Errorf("reminder query service not initialized")
	}

	h.jobManager.logger.Info("Processing due reminders...")

	// This will delegate to the existing checkReminders functionality in JobManager
	// which is already implemented and working
	return h.remindTask()
}

// remindTask processes due reminders with proper implementation
func (h *TaskRemindHandler) remindTask() error {
	h.jobManager.logger.Info("Checking and processing due reminders...")

	if h.jobManager.reminderQueryService == nil {
		h.jobManager.logger.Warn("Reminder query service not initialized, skipping reminder processing")
		return fmt.Errorf("reminder query service not initialized")
	}

	// Find due reminders
	now := time.Now()
	dueReminders, err := h.jobManager.reminderQueryService.FindDueReminders(now)
	if err != nil {
		h.jobManager.logger.Errorf("Failed to fetch due reminders: %v", err)
		return fmt.Errorf("failed to fetch due reminders: %w", err)
	}

	if len(dueReminders) == 0 {
		h.jobManager.logger.Info("No due reminders found")
		return nil
	}

	h.jobManager.logger.Infof("Found %d due reminders to process", len(dueReminders))

	// Process each reminder
	successCount := 0
	for _, reminder := range dueReminders {
		if err := h.processReminder(reminder); err != nil {
			h.jobManager.logger.Errorf("Failed to process reminder %s (%s): %v", reminder.ID, reminder.Name, err)
		} else {
			successCount++
		}
	}

	h.jobManager.logger.Infof("Successfully processed %d/%d reminders", successCount, len(dueReminders))
	return nil
}

// processReminder handles a single reminder: sends notifications and updates status
func (h *TaskRemindHandler) processReminder(reminder *models.Reminder) error {
	h.jobManager.logger.Infof("Processing reminder: %s - %s", reminder.ID, reminder.Name)

	// Get user information using query service
	user, err := h.jobManager.reminderQueryService.GetUserByID(reminder.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to get user info for reminder %s: %w", reminder.ID, err)
	}

	// Process different reminder methods
	if err := h.sendNotifications(reminder, user); err != nil {
		h.jobManager.logger.Errorf("Failed to send notifications for reminder %s: %v", reminder.ID, err)
		// Continue to mark as processed even if notifications fail
		return fmt.Errorf("failed to send notifications for reminder %s: %w", reminder.ID, err)
	} else {
		h.jobManager.logger.Infof("Successfully sent notifications for reminder %s", reminder.ID)
		// Update reminder status to 'active' (updated by the original creator)
		return h.jobManager.reminderQueryService.UpdateReminderStatus(reminder, "completed", reminder.CreatedBy)
	}
}

// sendNotifications sends notifications based on reminder methods
func (h *TaskRemindHandler) sendNotifications(reminder *models.Reminder, user *models.User) error {
	methods := strings.Split(reminder.RemindMethods, ",")
	var errors []string

	for _, method := range methods {
		method = strings.TrimSpace(method)
		switch method {
		case "email":
			if err := h.sendEmailNotification(reminder, user); err != nil {
				errors = append(errors, fmt.Sprintf("email: %v", err))
			}
		case "webhook":
			if err := h.sendWebhookNotification(reminder, user); err != nil {
				errors = append(errors, fmt.Sprintf("webhook: %v", err))
			}
		case "im", "message":
			if err := h.sendIMNotification(reminder, user); err != nil {
				errors = append(errors, fmt.Sprintf("im: %v", err))
			}
		default:
			h.jobManager.logger.Warnf("Unknown reminder method: %s", method)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errors, ", "))
	}
	return nil
}

// sendEmailNotification sends email notification
func (h *TaskRemindHandler) sendEmailNotification(reminder *models.Reminder, user *models.User) error {
	// Handle case where user has no email
	if user.Email == "" {
		h.jobManager.logger.Warnf("User %s has no email address, skipping email notification", user.ID)
		return fmt.Errorf("user %s has no email address, skipping email notification", user.ID)
	}

	if h.jobManager.emailSender == nil {
		h.jobManager.logger.Warn("Email sender not available, skipping email notification")
		return fmt.Errorf("email sender not available, skipping email notification")
	}

	subject := fmt.Sprintf("🔔 Reminder: %s", reminder.Name)
	body := h.formatReminderEmailBody(reminder, user)
	h.jobManager.logger.Infof("Sending email notification to %s: %s", user.Email, subject)
	return h.jobManager.emailSender.SendEmailToRecipients(subject, body, []string{user.Email}, nil)
}

// sendWebhookNotification sends webhook notification
func (h *TaskRemindHandler) sendWebhookNotification(reminder *models.Reminder, user *models.User) error {
	// TODO: Implement webhook notification
	h.jobManager.logger.Infof("Webhook notification for reminder %s (not implemented yet)", reminder.ID)
	return nil
}

// sendIMNotification sends instant message notification
func (h *TaskRemindHandler) sendIMNotification(reminder *models.Reminder, user *models.User) error {
	// TODO: Implement IM notification
	h.jobManager.logger.Infof("IM notification for reminder %s (not implemented yet)", reminder.ID)
	return nil
}

// formatReminderEmailBody creates the email body content
func (h *TaskRemindHandler) formatReminderEmailBody(reminder *models.Reminder, user *models.User) string {
	return fmt.Sprintf(`Hello %s,

This is a reminder for you:

**%s**

%s

Reminder Time: %s
Methods: %s
Targets: %s

Best regards,
Lazy Rabbit Reminder System`,
		user.Username,
		reminder.Name,
		reminder.Content,
		reminder.RemindTime.Format("2006-01-02 15:04:05"),
		reminder.RemindMethods,
		reminder.RemindTargets,
	)
}
