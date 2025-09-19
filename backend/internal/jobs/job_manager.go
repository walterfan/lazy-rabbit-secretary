package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/email"
	"gorm.io/gorm"
)

// =============================================================================
// TYPES AND STRUCTURES
// =============================================================================

// CronJob represents a scheduled task configuration from YAML
type CronJob struct {
	Name       string `yaml:"name"`
	Schedule   string `yaml:"schedule"`
	Function   string `yaml:"function"`
	Parameters string `yaml:"parameters"`
	Deadline   string `yaml:"deadline"`
}

// Config holds the task configuration
type Config struct {
	Jobs []CronJob `yaml:"jobs"`
}

// No duplicate structs needed - we use models.Reminder, models.User, models.Task directly
// with dedicated query services to avoid import cycles

// JobManager manages scheduled tasks and reminders
type JobManager struct {
	// Core dependencies
	config      *Config
	logger      *zap.SugaredLogger
	ctx         context.Context
	rdb         *redis.Client
	db          *gorm.DB
	emailSender *email.EmailSender

	// Query services to avoid import cycles
	reminderQueryService *ReminderQueryService
	taskQueryService     *TaskQueryService

	// Runtime state
	cronScheduler *cron.Cron
}

// =============================================================================
// CONSTRUCTOR AND INITIALIZATION
// =============================================================================

// NewJobManager creates a new JobManager instance
func NewJobManager(logger *zap.Logger, redisClient *redis.Client, db *gorm.DB) *JobManager {
	// Initialize email sender with graceful degradation
	emailSender, err := email.NewEmailSender()
	if err != nil {
		logger.Sugar().Warnf("Failed to initialize email sender: %v. Email notifications will be disabled.", err)
	}

	jm := &JobManager{
		config:               nil, // Will be loaded later
		logger:               logger.Sugar(),
		ctx:                  context.Background(),
		rdb:                  redisClient,
		db:                   db,
		emailSender:          emailSender,
		reminderQueryService: NewReminderQueryService(db),
		taskQueryService:     NewTaskQueryService(db),
	}

	err = jm.loadConfig()
	if err != nil {
		logger.Sugar().Fatalf("Failed to load configuration: %v", err)
		return nil
	}

	return jm
}

// =============================================================================
// CONFIGURATION MANAGEMENT
// =============================================================================

// loadConfig loads task configuration from YAML file
func (jm *JobManager) loadConfig() error {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	jm.config = &config
	jm.logger.Infof("Loaded %d tasks from configuration", len(config.Jobs))

	// register job handlers: checkTask, writeBlog, generateCalendar
	RegisterJobHandler("checkTask", &TaskCheckHandler{jobManager: jm})
	RegisterJobHandler("remindTask", &TaskRemindHandler{jobManager: jm})
	RegisterJobHandler("writeBlog", &BlogWriteHandler{jobManager: jm})
	RegisterJobHandler("generateCalendar", &CalendarGenerateHandler{jobManager: jm})

	return nil
}

// =============================================================================
// TASK EXECUTION FUNCTIONS
// =============================================================================

// ExecuteFunction is a public method to execute job handlers by name
func (jm *JobManager) ExecuteFunction(functionName, param string) error {
	// Clean up function name
	functionName = strings.TrimSuffix(functionName, "()")

	jm.logger.Infof("Executing function '%s' with parameter '%s'", functionName, param)

	// Handle plugin functions
	handler, exists := JobHandlers[functionName]
	if !exists {
		jm.logger.Warnf("No handler found for function: %s", functionName)
		return fmt.Errorf("no handler found for function: %s", functionName)
	}

	if err := handler.Execute(param); err != nil {
		jm.logger.Errorf("Error executing plugin %s: %v", functionName, err)
		return fmt.Errorf("error executing plugin %s: %w", functionName, err)
	}

	return nil
}

// executeFunction maps function names to actual Go functions and executes them (private method)
func (jm *JobManager) executeFunction(functionName, param string) {
	// Use the public method but ignore the error (for backward compatibility)
	jm.ExecuteFunction(functionName, param)
}

// =============================================================================
// REDIS TASK MANAGEMENT
// =============================================================================

// readUnfinishedTasks retrieves all unfinished task keys from Redis
func (jm *JobManager) readUnfinishedTasks() ([]string, error) {
	keys, err := jm.rdb.Keys(jm.ctx, "task:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read task keys from Redis: %w", err)
	}
	return keys, nil
}

// publishTaskExpiryEvent publishes a task expiry event to Redis
func (jm *JobManager) publishTaskExpiryEvent(taskID string) error {
	err := jm.rdb.Publish(jm.ctx, "task_expiry_channel", taskID).Err()
	if err != nil {
		return fmt.Errorf("failed to publish task expiry event: %w", err)
	}
	return nil
}

// checkTaskExpiry checks for expired tasks and publishes events
func (jm *JobManager) checkTaskExpiry() {
	jm.logger.Debug("Checking for expired tasks...")

	tasks, err := jm.readUnfinishedTasks()
	if err != nil {
		jm.logger.Errorf("Failed to read unfinished tasks: %v", err)
		return
	}

	currentTime := time.Now().Unix()
	expiredCount := 0

	for _, taskKey := range tasks {
		expiryTime, err := jm.rdb.Get(jm.ctx, taskKey).Int64()
		if err != nil {
			jm.logger.Warnf("Failed to get expiry time for task %s: %v", taskKey, err)
			continue
		}

		jm.logger.Debugf("Task %s: expiry=%d, current=%d", taskKey, expiryTime, currentTime)

		if currentTime >= expiryTime {
			taskID := strings.TrimPrefix(taskKey, "task:")
			if err := jm.publishTaskExpiryEvent(taskID); err != nil {
				jm.logger.Errorf("Failed to publish expiry event for task %s: %v", taskID, err)
			} else {
				expiredCount++
				jm.logger.Infof("Published expiry event for task: %s", taskID)
			}
		}
	}

	if expiredCount > 0 {
		jm.logger.Infof("Processed %d expired tasks", expiredCount)
	}
}

// =============================================================================
// REMINDER MANAGEMENT
// =============================================================================

// checkReminders finds and processes due reminders
func (jm *JobManager) checkReminders() error {
	if jm.reminderQueryService == nil {
		jm.logger.Warn("Reminder query service not initialized, skipping reminder check")
		return fmt.Errorf("reminder query service not initialized")
	}

	jm.logger.Debug("Checking for due reminders...")

	// Query for due reminders using the query service
	now := time.Now().UTC()
	dueReminders, err := jm.reminderQueryService.FindDueReminders(now)
	if err != nil {
		jm.logger.Errorf("Failed to fetch due reminders: %v", err)
		return fmt.Errorf("failed to fetch due reminders: %w", err)
	}

	if len(dueReminders) == 0 {
		jm.logger.Debug("No due reminders found")
		return nil
	}

	jm.logger.Infof("Found %d due reminders to process", len(dueReminders))

	// Process each reminder
	successCount := 0
	for _, reminder := range dueReminders {
		if err := jm.processReminder(reminder); err != nil {
			jm.logger.Errorf("Failed to process reminder %s (%s): %v", reminder.ID, reminder.Name, err)
		} else {
			successCount++
		}
	}

	jm.logger.Infof("Successfully processed %d/%d reminders", successCount, len(dueReminders))
	return nil
}

// processReminder handles a single reminder: sends email and updates status
func (jm *JobManager) processReminder(reminder *models.Reminder) error {
	jm.logger.Infof("Processing reminder: %s - %s", reminder.ID, reminder.Name)

	// Get user information using query service
	user, err := jm.reminderQueryService.GetUserByID(reminder.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to get user info for %s reminder %s: %w", reminder.CreatedBy, reminder.ID, err)
	}

	// Handle case where user has no email
	if user.Email == "" {
		jm.logger.Warnf("User %s has no email address, marking reminder %s as active without notification", user.ID, reminder.ID)
		return fmt.Errorf("user %s has no email address, marking reminder %s as active without notification", user.ID, reminder.ID)
	}

	// Send email notification
	if jm.emailSender != nil {
		if err := jm.sendReminderEmail(reminder, user); err != nil {
			jm.logger.Errorf("Failed to send reminder email for %s: %v", reminder.ID, err)
			// Continue to mark as active even if email fails
			return fmt.Errorf("failed to send reminder email for %s: %w", reminder.ID, err)
		} else {
			jm.logger.Infof("Successfully sent reminder email for %s to %s", reminder.ID, user.Email)
			// Update reminder status
			return jm.markReminderAsActive(reminder)
		}
	} else {
		jm.logger.Warn("Email sender not available, skipping email notification")
		return fmt.Errorf("email sender not available, skipping email notification")
	}
}

// sendReminderEmail sends a formatted email notification for a reminder
func (jm *JobManager) sendReminderEmail(reminder *models.Reminder, user *models.User) error {
	subject := fmt.Sprintf("🔔 Reminder: %s", reminder.Name)

	body := jm.formatReminderEmailBody(reminder, user)

	return jm.emailSender.SendEmailToRecipients(subject, body, []string{user.Email}, nil)
}

// formatReminderEmailBody creates the email body content
func (jm *JobManager) formatReminderEmailBody(reminder *models.Reminder, user *models.User) string {
	return fmt.Sprintf(`Hello %s,

This is your scheduled reminder:

📋 Title: %s
📝 Content: %s
⏰ Scheduled Time: %s
🏷️ Tags: %s
📬 Notification Methods: %s
🎯 Targets: %s

---
Sent by Lazy Rabbit Secretary System
Time: %s`,
		user.Username,
		reminder.Name,
		reminder.Content,
		reminder.RemindTime.Format("2006-01-02 15:04:05"),
		reminder.Tags,
		reminder.RemindMethods,
		reminder.RemindTargets,
		time.Now().Format("2006-01-02 15:04:05"))
}

// markReminderAsActive updates a reminder's status to 'active'
func (jm *JobManager) markReminderAsActive(reminder *models.Reminder) error {
	err := jm.reminderQueryService.UpdateReminderStatus(reminder, "active", reminder.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to update reminder status: %w", err)
	}

	jm.logger.Infof("Marked reminder %s as active", reminder.ID)
	return nil
}

// generateRepeatTaskInstances creates new instances for repeating tasks
func (jm *JobManager) generateRepeatTaskInstances() {
	if jm.taskQueryService == nil {
		jm.logger.Warn("Task query service not initialized, skipping repeat task generation")
		return
	}

	jm.logger.Debug("Checking for repeat tasks that need new instances...")

	// Find all parent repeating tasks using query service
	parentTasks, err := jm.taskQueryService.FindParentRepeatingTasks()
	if err != nil {
		jm.logger.Errorf("Failed to fetch repeat tasks: %v", err)
		return
	}

	if len(parentTasks) == 0 {
		jm.logger.Debug("No repeat tasks found")
		return
	}

	jm.logger.Infof("Found %d repeat tasks to process", len(parentTasks))

	successCount := 0
	for _, task := range parentTasks {
		if err := jm.processRepeatTask(task); err != nil {
			jm.logger.Errorf("Failed to process repeat task %s: %v", task.ID, err)
		} else {
			successCount++
		}
	}

	jm.logger.Infof("Successfully processed %d/%d repeat tasks", successCount, len(parentTasks))
}

// processRepeatTask generates instances for a single repeat task
func (jm *JobManager) processRepeatTask(parentTask *models.Task) error {

	// Check if we need to generate more instances
	// Look ahead for the next 7 days to ensure we always have upcoming instances

	// Count existing future instances using query service
	existingCount, err := jm.taskQueryService.CountExistingInstances(parentTask.ID)
	if err != nil {
		return fmt.Errorf("failed to count existing instances: %w", err)
	}

	// Generate instances if we don't have enough in the look-ahead window
	maxInstancesToGenerate := 10 // Limit to prevent runaway generation
	if existingCount < 3 {       // Keep at least 3 instances in the pipeline
		instancesToCreate := maxInstancesToGenerate - int(existingCount)

		// Find the last instance date to continue from
		var lastInstanceTime time.Time
		instances, err := jm.taskQueryService.FindTaskInstancesAfterTime(parentTask.ID, time.Time{})

		if err != nil || len(instances) == 0 {
			// No instances yet, start from parent schedule time
			lastInstanceTime = parentTask.ScheduleTime
		} else {
			// Find the latest instance
			latestTime := instances[0].ScheduleTime
			for _, instance := range instances {
				if instance.ScheduleTime.After(latestTime) {
					latestTime = instance.ScheduleTime
				}
			}
			lastInstanceTime = latestTime
		}

		generatedCount := 0
		currentDate := lastInstanceTime

		for i := 0; i < instancesToCreate; i++ {
			nextDate := parentTask.GetNextOccurrence(currentDate)
			if nextDate == nil {
				break // No more occurrences
			}

			// Don't create instances too far in the past
			if nextDate.Before(time.Now().AddDate(0, 0, -1)) {
				currentDate = *nextDate
				continue
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
				Status:       "pending",
				ScheduleTime: *nextDate,
				Minutes:      parentTask.Minutes,
				Deadline:     newDeadline,
				Tags:         parentTask.Tags,

				// Mark as instance
				IsRepeating:  false,
				ParentTaskID: &parentTask.ID,

				CreatedBy: parentTask.CreatedBy,
				CreatedAt: time.Now(),
				UpdatedBy: parentTask.CreatedBy,
				UpdatedAt: time.Now(),
			}

			if err := jm.taskQueryService.CreateTaskInstance(instance); err != nil {
				jm.logger.Errorf("Failed to create task instance: %v", err)
				break
			}

			// Generate reminder if needed
			if parentTask.GenerateReminders {
				if err := jm.generateReminderForTaskInstance(instance); err != nil {
					jm.logger.Errorf("Failed to generate reminder for task instance %s: %v", instance.ID, err)
				}
			}

			generatedCount++
			currentDate = *nextDate
		}

		if generatedCount > 0 {
			// Update parent task instance count
			parentTask.InstanceCount += generatedCount
			if err := jm.db.Save(parentTask).Error; err != nil {
				jm.logger.Errorf("Failed to update parent task instance count: %v", err)
			}

			jm.logger.Infof("Generated %d new instances for repeat task: %s", generatedCount, parentTask.Name)
		}
	}

	return nil
}

// convertToModelTask is no longer needed - we use models.Task directly

// generateReminderForTaskInstance creates a reminder for a task instance
func (jm *JobManager) generateReminderForTaskInstance(taskInstance *models.Task) error {
	if !taskInstance.GenerateReminders || taskInstance.ReminderMethods == "" {
		return nil
	}

	// Calculate reminder time
	reminderTime := taskInstance.ScheduleTime.Add(-time.Duration(taskInstance.ReminderAdvanceMinutes) * time.Minute)

	// Don't create reminders for past times
	if reminderTime.Before(time.Now()) {
		return nil
	}

	// Create reminder
	reminder := &models.Reminder{
		ID:            uuid.NewString(),
		RealmID:       taskInstance.RealmID,
		Name:          fmt.Sprintf("Task Reminder: %s", taskInstance.Name),
		Content:       jm.formatTaskReminderContent(taskInstance),
		RemindTime:    reminderTime,
		Status:        "pending",
		Tags:          fmt.Sprintf("task,auto-generated,%s", taskInstance.Tags),
		RemindMethods: taskInstance.ReminderMethods,
		RemindTargets: taskInstance.ReminderTargets,
		CreatedBy:     taskInstance.CreatedBy,
		CreatedAt:     time.Now(),
		UpdatedBy:     taskInstance.CreatedBy,
		UpdatedAt:     time.Now(),
	}

	if err := jm.reminderQueryService.CreateReminder(reminder); err != nil {
		return fmt.Errorf("failed to create reminder: %w", err)
	}

	jm.logger.Infof("Created reminder %s for task instance %s", reminder.ID, taskInstance.ID)
	return nil
}

// formatTaskReminderContent creates formatted content for task reminders
func (jm *JobManager) formatTaskReminderContent(task *models.Task) string {
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

	if task.ParentTaskID != nil {
		content += "\n\n🔄 This is a recurring task instance."
	}

	content += "\n\n---\nGenerated automatically by Lazy Rabbit Secretary System"

	return content
}

// =============================================================================
// CRON SCHEDULER MANAGEMENT
// =============================================================================

// setupCronScheduler initializes and configures the cron scheduler
func (jm *JobManager) setupCronScheduler() *cron.Cron {
	c := cron.New(cron.WithSeconds()) // Enable seconds precision

	// Add system cron jobs
	jm.addSystemCronJobs(c)

	// Add configured tasks
	jm.addConfiguredTasks(c)

	return c
}

// addSystemCronJobs adds built-in system cron jobs
func (jm *JobManager) addSystemCronJobs(c *cron.Cron) {
	// Task expiry check (every minute)
	_, err := c.AddFunc("@every 1m", jm.checkTaskExpiry)
	if err != nil {
		jm.logger.Fatalf("Failed to add task expiry check cron job: %v", err)
	} else {
		jm.logger.Info("Scheduled task expiry check (every minute)")
	}

	// Reminder check (every minute)
	_, err = c.AddFunc("@every 1m", func() {
		if err := jm.checkReminders(); err != nil {
			jm.logger.Errorf("Reminder check failed: %v", err)
		}
	})
	if err != nil {
		jm.logger.Fatalf("Failed to add reminder check cron job: %v", err)
	} else {
		jm.logger.Info("Scheduled reminder check (every minute)")
	}

	// Repeat task instance generation (every hour)
	_, err = c.AddFunc("@every 1h", jm.generateRepeatTaskInstances)
	if err != nil {
		jm.logger.Fatalf("Failed to add repeat task generation cron job: %v", err)
	} else {
		jm.logger.Info("Scheduled repeat task instance generation (every hour)")
	}
}

// addConfiguredTasks adds tasks from configuration to the scheduler
func (jm *JobManager) addConfiguredTasks(c *cron.Cron) {
	if jm.config == nil {
		jm.logger.Warn("No configuration loaded, skipping configured tasks")
		return
	}

	for _, task := range jm.config.Jobs {
		jm.addSingleTask(c, task)
		jm.setTaskDeadline(task)
	}
}

// addSingleTask adds a single task to the cron scheduler
func (jm *JobManager) addSingleTask(c *cron.Cron, theJob CronJob) {
	if theJob.Function == "" || theJob.Schedule == "" {
		jm.logger.Warnf("Task %s has missing function or schedule, skipping", theJob.Name)
		return
	}

	// Parse function and parameters
	functionName, param := jm.parseFunctionCall(theJob.Function)

	// Create the cron job
	_, err := c.AddFunc(theJob.Schedule, func() {
		jm.executeFunction(functionName, param)
	})

	if err != nil {
		jm.logger.Errorf("Failed to add task %s: %v", theJob.Name, err)
	} else {
		jm.logger.Infof("Scheduled task '%s' with schedule '%s'", theJob.Name, theJob.Schedule)
	}
}

// parseFunctionCall parses a function call string into function name and parameter
func (jm *JobManager) parseFunctionCall(functionCall string) (string, string) {
	functionCall = strings.TrimSpace(functionCall)

	// Check if it's a function call with parameters
	if strings.Contains(functionCall, "(") && strings.Contains(functionCall, ")") {
		startIdx := strings.Index(functionCall, "(") + 1
		endIdx := strings.Index(functionCall, ")")
		if startIdx < endIdx {
			param := functionCall[startIdx:endIdx]
			functionName := strings.TrimSpace(functionCall[:startIdx-1])
			return functionName, param
		}
	}

	return functionCall, ""
}

// setTaskDeadline sets a task deadline in Redis if specified
func (jm *JobManager) setTaskDeadline(theJob CronJob) {
	if theJob.Deadline == "" {
		return
	}

	// Parse deadline
	deadlineTime, err := time.Parse(time.RFC3339, theJob.Deadline)
	if err != nil {
		jm.logger.Errorf("Failed to parse deadline for task %s: %v", theJob.Name, err)
		return
	}

	// Store in Redis
	key := fmt.Sprintf("task:%s:expiry", strings.ReplaceAll(theJob.Name, " ", "_"))
	expiryTime := deadlineTime.Unix()

	err = jm.rdb.Set(jm.ctx, key, expiryTime, 0).Err()
	if err != nil {
		jm.logger.Errorf("Failed to set expiry time for task %s: %v", theJob.Name, err)
		return
	}

	jm.logger.Infof("Set expiry time for task %s: %d", theJob.Name, expiryTime)
}

// =============================================================================
// MAIN ENTRY POINT
// =============================================================================

// CheckTasks is the main entry point that starts the job manager
func (jm *JobManager) CheckTasks() {
	jm.logger.Info("Starting Job Manager...")

	// Setup and start cron scheduler
	jm.cronScheduler = jm.setupCronScheduler()
	jm.cronScheduler.Start()

	jm.logger.Info("Job Manager started successfully")
}

// Stop gracefully stops the job manager
func (jm *JobManager) Stop() {
	if jm.cronScheduler != nil {
		jm.cronScheduler.Stop()
		jm.logger.Info("Job Manager stopped")
	}
}
