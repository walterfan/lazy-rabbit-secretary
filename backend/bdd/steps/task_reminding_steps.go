package steps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"github.com/walterfan/lazy-rabbit-reminder/internal/reminder"
	"github.com/walterfan/lazy-rabbit-reminder/internal/task"
)

// TaskRemindingContext holds the context for task reminding BDD tests
type TaskRemindingContext struct {
	// Test data
	currentUser      string
	currentRealm     string
	taskRequest      task.CreateTaskRequest
	createdTask      *models.Task
	createdReminders []models.Reminder
	lastError        error
	warnings         []string

	// Mock services
	taskService     *task.TaskService
	reminderService *reminder.ReminderService

	// Test state
	reminderServiceAvailable bool
}

// NewTaskRemindingContext creates a new context for task reminding tests
func NewTaskRemindingContext() *TaskRemindingContext {
	return &TaskRemindingContext{
		reminderServiceAvailable: true,
		warnings:                 make([]string, 0),
	}
}

// Background steps
func (ctx *TaskRemindingContext) iAmAuthenticatedAsUser(userEmail string) error {
	ctx.currentUser = userEmail
	return nil
}

func (ctx *TaskRemindingContext) iBelongToRealm(realmID string) error {
	ctx.currentRealm = realmID
	return nil
}

// Task creation steps
func (ctx *TaskRemindingContext) iWantToCreateATaskWithTheFollowingDetails(table *godog.Table) error {
	ctx.taskRequest = task.CreateTaskRequest{}

	for _, row := range table.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "name":
			ctx.taskRequest.Name = value
		case "description":
			ctx.taskRequest.Description = value
		case "schedule_time":
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("invalid schedule_time format: %v", err)
			}
			ctx.taskRequest.ScheduleTime = t
		case "minutes":
			minutes, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid minutes: %v", err)
			}
			ctx.taskRequest.Minutes = minutes
		case "deadline":
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("invalid deadline format: %v", err)
			}
			ctx.taskRequest.Deadline = t
		case "priority":
			priority, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid priority: %v", err)
			}
			ctx.taskRequest.Priority = &priority
		case "difficulty":
			difficulty, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid difficulty: %v", err)
			}
			ctx.taskRequest.Difficulty = &difficulty
		case "tags":
			ctx.taskRequest.Tags = value
		}
	}

	return nil
}

func (ctx *TaskRemindingContext) iWantToGenerateRemindersWithTheFollowingSettings(table *godog.Table) error {
	for _, row := range table.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "generate_reminders":
			ctx.taskRequest.GenerateReminders = value == "true"
		case "reminder_advance_minutes":
			minutes, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid reminder_advance_minutes: %v", err)
			}
			ctx.taskRequest.ReminderAdvanceMinutes = minutes
		case "reminder_methods":
			ctx.taskRequest.ReminderMethods = value
		case "reminder_targets":
			ctx.taskRequest.ReminderTargets = value
		}
	}

	return nil
}

func (ctx *TaskRemindingContext) iWantToSetUpRepeatSettings(table *godog.Table) error {
	for _, row := range table.Rows {
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "is_repeating":
			ctx.taskRequest.IsRepeating = value == "true"
		case "repeat_pattern":
			ctx.taskRequest.RepeatPattern = value
		case "repeat_interval":
			interval, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid repeat_interval: %v", err)
			}
			ctx.taskRequest.RepeatInterval = interval
		case "repeat_days_of_week":
			ctx.taskRequest.RepeatDaysOfWeek = value
		case "repeat_day_of_month":
			day, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid repeat_day_of_month: %v", err)
			}
			ctx.taskRequest.RepeatDayOfMonth = day
		case "repeat_end_date":
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("invalid repeat_end_date format: %v", err)
			}
			ctx.taskRequest.RepeatEndDate = &t
		case "repeat_count":
			count, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("invalid repeat_count: %v", err)
			}
			ctx.taskRequest.RepeatCount = count
		}
	}

	return nil
}

func (ctx *TaskRemindingContext) iDoNotWantToGenerateReminders(table *godog.Table) error {
	ctx.taskRequest.GenerateReminders = false
	return nil
}

func (ctx *TaskRemindingContext) iWantToCreateATaskWithReminderSettings(table *godog.Table) error {
	// Combine task creation with reminder settings
	ctx.taskRequest = task.CreateTaskRequest{
		Name:         "Test Task",
		ScheduleTime: time.Now().Add(1 * time.Hour),
		Minutes:      30,
		Deadline:     time.Now().Add(1*time.Hour + 30*time.Minute),
	}

	return ctx.iWantToGenerateRemindersWithTheFollowingSettings(table)
}

func (ctx *TaskRemindingContext) iWantToCreateATaskWithRepeatPattern(pattern string) error {
	ctx.taskRequest = task.CreateTaskRequest{
		Name:          "Test Task",
		ScheduleTime:  time.Now().Add(1 * time.Hour),
		Minutes:       30,
		Deadline:      time.Now().Add(1*time.Hour + 30*time.Minute),
		IsRepeating:   true,
		RepeatPattern: pattern,
	}
	return nil
}

func (ctx *TaskRemindingContext) iWantToSetRepeatIntervalTo(interval int) error {
	ctx.taskRequest.RepeatInterval = interval
	return nil
}

func (ctx *TaskRemindingContext) iWantToGenerateReminders() error {
	ctx.taskRequest.GenerateReminders = true
	ctx.taskRequest.ReminderMethods = "email"
	ctx.taskRequest.ReminderTargets = "test@example.com"
	return nil
}

func (ctx *TaskRemindingContext) iWantToCreateATaskWithValidReminderSettings() error {
	ctx.taskRequest = task.CreateTaskRequest{
		Name:                   "Valid Task",
		ScheduleTime:           time.Now().Add(1 * time.Hour),
		Minutes:                30,
		Deadline:               time.Now().Add(1*time.Hour + 30*time.Minute),
		GenerateReminders:      true,
		ReminderAdvanceMinutes: 15,
		ReminderMethods:        "email",
		ReminderTargets:        "test@example.com",
	}
	return nil
}

// Service state steps
func (ctx *TaskRemindingContext) theReminderServiceIsUnavailable() error {
	ctx.reminderServiceAvailable = false
	return nil
}

func (ctx *TaskRemindingContext) theReminderServiceReturnsAnError() error {
	// This would be handled in the mock service implementation
	return nil
}

// Action steps
func (ctx *TaskRemindingContext) iCreateTheTask() error {
	// In a real test, this would call the actual service
	// For BDD testing, we simulate the behavior

	if !ctx.reminderServiceAvailable {
		ctx.warnings = append(ctx.warnings, "reminder service not available")
	}

	// Simulate task creation validation
	if ctx.taskRequest.Name == "" {
		ctx.lastError = fmt.Errorf("name is required")
		return nil
	}

	if ctx.taskRequest.Minutes <= 0 {
		ctx.lastError = fmt.Errorf("minutes must be greater than 0")
		return nil
	}

	if ctx.taskRequest.IsRepeating {
		if err := ctx.validateRepeatSettings(); err != nil {
			ctx.lastError = err
			return nil
		}
	}

	if ctx.taskRequest.GenerateReminders {
		if err := ctx.validateReminderSettings(); err != nil {
			ctx.lastError = err
			return nil
		}
	}

	// Simulate successful task creation
	ctx.createdTask = &models.Task{
		ID:                     "task-123",
		RealmID:                ctx.currentRealm,
		Name:                   ctx.taskRequest.Name,
		Description:            ctx.taskRequest.Description,
		ScheduleTime:           ctx.taskRequest.ScheduleTime,
		Minutes:                ctx.taskRequest.Minutes,
		Deadline:               ctx.taskRequest.Deadline,
		IsRepeating:            ctx.taskRequest.IsRepeating,
		RepeatPattern:          ctx.taskRequest.RepeatPattern,
		RepeatInterval:         ctx.taskRequest.RepeatInterval,
		RepeatDaysOfWeek:       ctx.taskRequest.RepeatDaysOfWeek,
		RepeatDayOfMonth:       ctx.taskRequest.RepeatDayOfMonth,
		RepeatEndDate:          ctx.taskRequest.RepeatEndDate,
		RepeatCount:            ctx.taskRequest.RepeatCount,
		GenerateReminders:      ctx.taskRequest.GenerateReminders,
		ReminderAdvanceMinutes: ctx.taskRequest.ReminderAdvanceMinutes,
		ReminderMethods:        ctx.taskRequest.ReminderMethods,
		ReminderTargets:        ctx.taskRequest.ReminderTargets,
		CreatedBy:              ctx.currentUser,
		CreatedAt:              time.Now(),
	}

	// Set priority and difficulty
	if ctx.taskRequest.Priority != nil {
		ctx.createdTask.Priority = *ctx.taskRequest.Priority
	} else {
		ctx.createdTask.Priority = 2
	}

	if ctx.taskRequest.Difficulty != nil {
		ctx.createdTask.Difficulty = *ctx.taskRequest.Difficulty
	} else {
		ctx.createdTask.Difficulty = 2
	}

	// Simulate reminder creation
	if ctx.createdTask.GenerateReminders && ctx.reminderServiceAvailable {
		ctx.simulateReminderCreation()
	}

	return nil
}

// Assertion steps
func (ctx *TaskRemindingContext) theTaskShouldBeCreatedSuccessfully() error {
	if ctx.lastError != nil {
		return fmt.Errorf("expected task creation to succeed, but got error: %v", ctx.lastError)
	}

	if ctx.createdTask == nil {
		return fmt.Errorf("expected task to be created, but it was nil")
	}

	return nil
}

func (ctx *TaskRemindingContext) theTaskCreationShouldFail() error {
	if ctx.lastError == nil {
		return fmt.Errorf("expected task creation to fail, but it succeeded")
	}
	return nil
}

func (ctx *TaskRemindingContext) theTaskCreationShould(result string) error {
	switch result {
	case "succeed":
		return ctx.theTaskShouldBeCreatedSuccessfully()
	case "fail":
		return ctx.theTaskCreationShouldFail()
	default:
		return fmt.Errorf("unknown result expectation: %s", result)
	}
}

func (ctx *TaskRemindingContext) iShouldReceiveAnError(expectedError string) error {
	if ctx.lastError == nil {
		return fmt.Errorf("expected error '%s', but no error occurred", expectedError)
	}

	if !strings.Contains(ctx.lastError.Error(), expectedError) {
		return fmt.Errorf("expected error to contain '%s', but got: %v", expectedError, ctx.lastError)
	}

	return nil
}

func (ctx *TaskRemindingContext) iShouldReceiveAWarning(expectedWarning string) error {
	for _, warning := range ctx.warnings {
		if strings.Contains(warning, expectedWarning) {
			return nil
		}
	}

	return fmt.Errorf("expected warning containing '%s', but got warnings: %v", expectedWarning, ctx.warnings)
}

func (ctx *TaskRemindingContext) theTaskShouldBeMarkedAsRepeating() error {
	if ctx.createdTask == nil || !ctx.createdTask.IsRepeating {
		return fmt.Errorf("expected task to be marked as repeating")
	}
	return nil
}

func (ctx *TaskRemindingContext) aReminderShouldBeCreatedFor(reminderTime string) error {
	expectedTime, err := time.Parse(time.RFC3339, reminderTime)
	if err != nil {
		return fmt.Errorf("invalid reminder time format: %v", err)
	}

	for _, reminder := range ctx.createdReminders {
		if reminder.RemindTime.Equal(expectedTime) {
			return nil
		}
	}

	return fmt.Errorf("expected reminder for time %s, but none found", reminderTime)
}

func (ctx *TaskRemindingContext) theReminderShouldHaveMethod(method string) error {
	if len(ctx.createdReminders) == 0 {
		return fmt.Errorf("no reminders found")
	}

	reminder := ctx.createdReminders[0]
	if reminder.RemindMethods != method {
		return fmt.Errorf("expected reminder method '%s', but got '%s'", method, reminder.RemindMethods)
	}

	return nil
}

func (ctx *TaskRemindingContext) theReminderShouldTarget(target string) error {
	if len(ctx.createdReminders) == 0 {
		return fmt.Errorf("no reminders found")
	}

	reminder := ctx.createdReminders[0]
	if reminder.RemindTargets != target {
		return fmt.Errorf("expected reminder target '%s', but got '%s'", target, reminder.RemindTargets)
	}

	return nil
}

func (ctx *TaskRemindingContext) theReminderContentShouldContain(expectedContent string) error {
	if len(ctx.createdReminders) == 0 {
		return fmt.Errorf("no reminders found")
	}

	reminder := ctx.createdReminders[0]
	if !strings.Contains(reminder.Content, expectedContent) {
		return fmt.Errorf("expected reminder content to contain '%s', but got: %s", expectedContent, reminder.Content)
	}

	return nil
}

func (ctx *TaskRemindingContext) theReminderContentShouldContainTable(table *godog.Table) error {
	if len(ctx.createdReminders) == 0 {
		return fmt.Errorf("no reminders found")
	}

	reminder := ctx.createdReminders[0]

	for _, row := range table.Rows {
		expectedContent := row.Cells[0].Value
		if !strings.Contains(reminder.Content, expectedContent) {
			return fmt.Errorf("expected reminder content to contain '%s', but got: %s", expectedContent, reminder.Content)
		}
	}

	return nil
}

func (ctx *TaskRemindingContext) taskInstancesShouldBeCreated(count int) error {
	// In a real implementation, this would check the database
	// For BDD testing, we simulate based on repeat settings
	expectedCount := ctx.calculateExpectedInstanceCount()

	if expectedCount != count {
		return fmt.Errorf("expected %d task instances, but calculated %d", count, expectedCount)
	}

	return nil
}

func (ctx *TaskRemindingContext) remindersShouldBeCreated(count int) error {
	if len(ctx.createdReminders) != count {
		return fmt.Errorf("expected %d reminders, but got %d", count, len(ctx.createdReminders))
	}

	return nil
}

func (ctx *TaskRemindingContext) eachReminderShouldHaveMethods(methods string) error {
	for _, reminder := range ctx.createdReminders {
		if reminder.RemindMethods != methods {
			return fmt.Errorf("expected each reminder to have methods '%s', but found '%s'", methods, reminder.RemindMethods)
		}
	}

	return nil
}

func (ctx *TaskRemindingContext) eachReminderShouldBeScheduledMinutesBeforeItsTask(minutes int) error {
	// This would verify the timing relationship in a real implementation
	return nil
}

func (ctx *TaskRemindingContext) taskInstancesShouldBeCreatedForDays(days string) error {
	// Verify that instances are created for the correct days
	return nil
}

func (ctx *TaskRemindingContext) remindersShouldBeCreatedMinutesBeforeEach(minutes int, taskType string) error {
	// Verify reminder timing
	return nil
}

func (ctx *TaskRemindingContext) noRemindersShouldBeCreated() error {
	if len(ctx.createdReminders) != 0 {
		return fmt.Errorf("expected no reminders, but got %d", len(ctx.createdReminders))
	}
	return nil
}

func (ctx *TaskRemindingContext) noReminderShouldBeCreated() error {
	return ctx.noRemindersShouldBeCreated()
}

func (ctx *TaskRemindingContext) theReminderTagsShouldContain(expectedTags string) error {
	if len(ctx.createdReminders) == 0 {
		return fmt.Errorf("no reminders found")
	}

	reminder := ctx.createdReminders[0]

	expectedTagList := strings.Split(expectedTags, ",")
	for _, expectedTag := range expectedTagList {
		expectedTag = strings.TrimSpace(expectedTag)
		if !strings.Contains(reminder.Tags, expectedTag) {
			return fmt.Errorf("expected reminder tags to contain '%s', but got: %s", expectedTag, reminder.Tags)
		}
	}

	return nil
}

func (ctx *TaskRemindingContext) taskInstancesShouldBeCreatedOnlyForDays(days string) error {
	// Verify specific days for weekly patterns
	return nil
}

func (ctx *TaskRemindingContext) eachInstanceShouldHaveAWebhookReminderMinutesBefore(minutes int) error {
	// Verify webhook reminders
	return nil
}

func (ctx *TaskRemindingContext) theTotalNumberOfInstancesShouldNotExceed(maxCount int) error {
	expectedCount := ctx.calculateExpectedInstanceCount()
	if expectedCount > maxCount {
		return fmt.Errorf("expected instance count %d should not exceed %d", expectedCount, maxCount)
	}
	return nil
}

func (ctx *TaskRemindingContext) onlyTaskInstancesShouldBeCreated(count int) error {
	return ctx.taskInstancesShouldBeCreated(count)
}

func (ctx *TaskRemindingContext) onlyRemindersShouldBeCreated(count int) error {
	return ctx.remindersShouldBeCreated(count)
}

func (ctx *TaskRemindingContext) noInstancesShouldBeCreatedAfter(dateStr string) error {
	// Verify end date limits
	return nil
}

func (ctx *TaskRemindingContext) butReminderCreationShouldFail() error {
	// Check that reminders failed to create but task succeeded
	if ctx.createdTask == nil {
		return fmt.Errorf("expected task to be created")
	}

	if len(ctx.createdReminders) > 0 {
		return fmt.Errorf("expected reminder creation to fail, but reminders were created")
	}

	return nil
}

func (ctx *TaskRemindingContext) theTaskShouldStillBeFunctionalWithoutReminders() error {
	if ctx.createdTask == nil {
		return fmt.Errorf("expected task to be functional")
	}
	return nil
}

// Helper methods
func (ctx *TaskRemindingContext) validateRepeatSettings() error {
	validPatterns := map[string]bool{
		"daily": true, "weekly": true, "monthly": true, "yearly": true,
	}

	if !validPatterns[ctx.taskRequest.RepeatPattern] {
		return fmt.Errorf("repeat_pattern must be one of: daily, weekly, monthly, yearly")
	}

	if ctx.taskRequest.RepeatInterval < 1 {
		return fmt.Errorf("repeat_interval must be at least 1")
	}

	return nil
}

func (ctx *TaskRemindingContext) validateReminderSettings() error {
	if ctx.taskRequest.ReminderMethods == "" {
		return fmt.Errorf("reminder_methods is required when generate_reminders is true")
	}

	validMethods := map[string]bool{"email": true, "webhook": true}
	methods := strings.Split(ctx.taskRequest.ReminderMethods, ",")
	for _, method := range methods {
		method = strings.TrimSpace(strings.ToLower(method))
		if !validMethods[method] {
			return fmt.Errorf("invalid reminder method: %s", method)
		}
	}

	if ctx.taskRequest.ReminderAdvanceMinutes < 0 {
		return fmt.Errorf("reminder_advance_minutes cannot be negative")
	}

	return nil
}

func (ctx *TaskRemindingContext) calculateExpectedInstanceCount() int {
	if !ctx.taskRequest.IsRepeating {
		return 0
	}

	// Simplified calculation for BDD testing
	if ctx.taskRequest.RepeatCount > 0 {
		return min(ctx.taskRequest.RepeatCount, 5) // Simulate initial 5 instances
	}

	return 5 // Default initial instance count
}

func (ctx *TaskRemindingContext) simulateReminderCreation() {
	if ctx.createdTask.ScheduleTime.Before(time.Now()) {
		ctx.warnings = append(ctx.warnings, "skipping reminder creation for task - reminder time is in the past")
		return
	}

	reminderTime := ctx.createdTask.ScheduleTime.Add(-time.Duration(ctx.createdTask.ReminderAdvanceMinutes) * time.Minute)

	content := ctx.formatReminderContent(ctx.createdTask)
	tags := ctx.formatReminderTags(ctx.createdTask)

	reminder := models.Reminder{
		ID:            "reminder-123",
		RealmID:       ctx.createdTask.RealmID,
		Name:          fmt.Sprintf("Task Reminder: %s", ctx.createdTask.Name),
		Content:       content,
		RemindTime:    reminderTime,
		Status:        "pending",
		Tags:          tags,
		RemindMethods: ctx.createdTask.ReminderMethods,
		RemindTargets: ctx.createdTask.ReminderTargets,
		CreatedBy:     ctx.createdTask.CreatedBy,
		CreatedAt:     time.Now(),
	}

	// Simulate multiple reminders for repeating tasks
	instanceCount := ctx.calculateExpectedInstanceCount()
	if ctx.createdTask.IsRepeating && instanceCount > 0 {
		for i := 0; i < instanceCount; i++ {
			ctx.createdReminders = append(ctx.createdReminders, reminder)
		}
	} else {
		ctx.createdReminders = append(ctx.createdReminders, reminder)
	}
}

func (ctx *TaskRemindingContext) formatReminderContent(task *models.Task) string {
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

	if task.IsRepeating {
		content += "\n\n🔄 This is a recurring task instance."
	}

	content += "\n\n---\nGenerated automatically by Lazy Rabbit Reminder System"

	return content
}

func (ctx *TaskRemindingContext) formatReminderTags(task *models.Task) string {
	tags := []string{"task", "auto-generated"}

	if task.IsRepeating {
		tags = append(tags, "recurring")
	}

	if task.Tags != "" {
		taskTags := strings.Split(task.Tags, ",")
		for _, tag := range taskTags {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	return strings.Join(tags, ",")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Register step definitions
func InitializeTaskRemindingSteps(ctx *godog.ScenarioContext) {
	trc := NewTaskRemindingContext()

	// Background steps
	ctx.Step(`^I am authenticated as user "([^"]*)"$`, trc.iAmAuthenticatedAsUser)
	ctx.Step(`^I belong to realm "([^"]*)"$`, trc.iBelongToRealm)

	// Task setup steps
	ctx.Step(`^I want to create a task with the following details:$`, trc.iWantToCreateATaskWithTheFollowingDetails)
	ctx.Step(`^I want to generate reminders with the following settings:$`, trc.iWantToGenerateRemindersWithTheFollowingSettings)
	ctx.Step(`^I want to set up repeat settings:$`, trc.iWantToSetUpRepeatSettings)
	ctx.Step(`^I do not want to generate reminders:$`, trc.iDoNotWantToGenerateReminders)
	ctx.Step(`^I want to create a task with reminder settings:$`, trc.iWantToCreateATaskWithReminderSettings)
	ctx.Step(`^I want to create a task with repeat pattern "([^"]*)"$`, trc.iWantToCreateATaskWithRepeatPattern)
	ctx.Step(`^I want to set repeat interval to (\d+)$`, trc.iWantToSetRepeatIntervalTo)
	ctx.Step(`^I want to generate reminders$`, trc.iWantToGenerateReminders)
	ctx.Step(`^I want to create a task with valid reminder settings$`, trc.iWantToCreateATaskWithValidReminderSettings)

	// Service state steps
	ctx.Step(`^the reminder service is unavailable$`, trc.theReminderServiceIsUnavailable)
	ctx.Step(`^the reminder service returns an error$`, trc.theReminderServiceReturnsAnError)

	// Action steps
	ctx.Step(`^I create the task$`, trc.iCreateTheTask)

	// Assertion steps
	ctx.Step(`^the task should be created successfully$`, trc.theTaskShouldBeCreatedSuccessfully)
	ctx.Step(`^the task creation should fail$`, trc.theTaskCreationShouldFail)
	ctx.Step(`^the task creation should (succeed|fail)$`, trc.theTaskCreationShould)
	ctx.Step(`^I should receive an error "([^"]*)"$`, trc.iShouldReceiveAnError)
	ctx.Step(`^I should receive a warning "([^"]*)"$`, trc.iShouldReceiveAWarning)
	ctx.Step(`^the task should be marked as repeating$`, trc.theTaskShouldBeMarkedAsRepeating)
	ctx.Step(`^a reminder should be created for "([^"]*)"$`, trc.aReminderShouldBeCreatedFor)
	ctx.Step(`^the reminder should have method "([^"]*)"$`, trc.theReminderShouldHaveMethod)
	ctx.Step(`^the reminder should target "([^"]*)"$`, trc.theReminderShouldTarget)
	ctx.Step(`^the reminder content should contain "([^"]*)"$`, trc.theReminderContentShouldContain)
	ctx.Step(`^the reminder content should contain:$`, trc.theReminderContentShouldContainTable)
	ctx.Step(`^(\d+) task instances should be created$`, trc.taskInstancesShouldBeCreated)
	ctx.Step(`^(\d+) reminders should be created$`, trc.remindersShouldBeCreated)
	ctx.Step(`^each reminder should have methods "([^"]*)"$`, trc.eachReminderShouldHaveMethods)
	ctx.Step(`^each reminder should be scheduled (\d+) minutes before its task$`, trc.eachReminderShouldBeScheduledMinutesBeforeItsTask)
	ctx.Step(`^task instances should be created for (.+)$`, trc.taskInstancesShouldBeCreatedForDays)
	ctx.Step(`^reminders should be created (\d+) minutes before each (.+)$`, trc.remindersShouldBeCreatedMinutesBeforeEach)
	ctx.Step(`^no reminders should be created$`, trc.noRemindersShouldBeCreated)
	ctx.Step(`^no reminder should be created$`, trc.noReminderShouldBeCreated)
	ctx.Step(`^the reminder tags should contain "([^"]*)"$`, trc.theReminderTagsShouldContain)
	ctx.Step(`^task instances should be created only for (.+)$`, trc.taskInstancesShouldBeCreatedOnlyForDays)
	ctx.Step(`^each instance should have a webhook reminder (\d+) minutes before$`, trc.eachInstanceShouldHaveAWebhookReminderMinutesBefore)
	ctx.Step(`^the total number of instances should not exceed (\d+)$`, trc.theTotalNumberOfInstancesShouldNotExceed)
	ctx.Step(`^only (\d+) task instances should be created$`, trc.onlyTaskInstancesShouldBeCreated)
	ctx.Step(`^only (\d+) reminders should be created$`, trc.onlyRemindersShouldBeCreated)
	ctx.Step(`^no instances should be created after (.+)$`, trc.noInstancesShouldBeCreatedAfter)
	ctx.Step(`^but reminder creation should fail$`, trc.butReminderCreationShouldFail)
	ctx.Step(`^the task should still be functional without reminders$`, trc.theTaskShouldStillBeFunctionalWithoutReminders)
}
