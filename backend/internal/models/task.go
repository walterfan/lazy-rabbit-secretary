package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID           string     `json:"id" gorm:"primaryKey;type:text"`
	RealmID      string     `json:"realm_id" gorm:"not null;type:text;index"`
	Name         string     `json:"name" gorm:"not null;type:text"`
	Description  string     `json:"description" gorm:"type:text"`
	Priority     int        `json:"priority" gorm:"not null;default:2"`
	Difficulty   int        `json:"difficulty" gorm:"not null;default:2"`
	Status       TaskStatus `json:"status" gorm:"default:'pending';type:text;index"`
	ScheduleTime time.Time  `json:"schedule_time" gorm:"not null"`
	Minutes      int        `json:"minutes" gorm:"not null"`
	Deadline     time.Time  `json:"deadline" gorm:"not null"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Tags         string     `json:"tags" gorm:"type:text"`

	// Repeat task fields
	IsRepeating      bool       `json:"is_repeating" gorm:"default:false"`
	RepeatPattern    string     `json:"repeat_pattern" gorm:"type:text"`      // daily, weekly, monthly, yearly, custom
	RepeatInterval   int        `json:"repeat_interval" gorm:"default:1"`     // every N days/weeks/months
	RepeatDaysOfWeek string     `json:"repeat_days_of_week" gorm:"type:text"` // comma-separated: mon,tue,wed
	RepeatDayOfMonth int        `json:"repeat_day_of_month" gorm:"default:0"` // for monthly: 1-31, 0=same day
	RepeatEndDate    *time.Time `json:"repeat_end_date"`                      // when to stop repeating
	RepeatCount      int        `json:"repeat_count" gorm:"default:0"`        // max occurrences, 0=infinite

	// Reminder generation settings
	GenerateReminders      bool   `json:"generate_reminders" gorm:"default:false"`
	ReminderAdvanceMinutes int    `json:"reminder_advance_minutes" gorm:"default:60"` // remind N minutes before task
	ReminderMethods        string `json:"reminder_methods" gorm:"type:text"`          // email,webhook,message
	ReminderTargets        string `json:"reminder_targets" gorm:"type:text"`          // notification targets

	// Tracking
	ParentTaskID  *string `json:"parent_task_id" gorm:"type:text;index"` // for generated task instances
	InstanceCount int     `json:"instance_count" gorm:"default:0"`       // how many instances generated

	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (Task) TableName() string {
	return "tasks"
}

// IsParentTask returns true if this is a repeating task template
func (t *Task) IsParentTask() bool {
	return t.IsRepeating && t.ParentTaskID == nil
}

// IsTaskInstance returns true if this is an instance of a repeating task
func (t *Task) IsTaskInstance() bool {
	return t.ParentTaskID != nil
}

// GetRepeatDaysOfWeek returns the days of week as a slice
func (t *Task) GetRepeatDaysOfWeek() []string {
	if t.RepeatDaysOfWeek == "" {
		return []string{}
	}
	days := strings.Split(t.RepeatDaysOfWeek, ",")
	result := make([]string, 0, len(days))
	for _, day := range days {
		if trimmed := strings.TrimSpace(day); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ShouldGenerateReminders returns true if this task should generate reminders
func (t *Task) ShouldGenerateReminders() bool {
	return t.GenerateReminders && t.ReminderMethods != ""
}

// GetNextOccurrence calculates the next occurrence date based on repeat pattern
func (t *Task) GetNextOccurrence(fromDate time.Time) *time.Time {
	if !t.IsRepeating {
		return nil
	}

	var nextDate time.Time

	switch t.RepeatPattern {
	case "daily":
		nextDate = fromDate.AddDate(0, 0, t.RepeatInterval)
	case "weekly":
		if len(t.GetRepeatDaysOfWeek()) > 0 {
			nextDate = t.getNextWeeklyOccurrence(fromDate)
		} else {
			nextDate = fromDate.AddDate(0, 0, 7*t.RepeatInterval)
		}
	case "monthly":
		nextDate = t.getNextMonthlyOccurrence(fromDate)
	case "yearly":
		nextDate = fromDate.AddDate(t.RepeatInterval, 0, 0)
	default:
		return nil
	}

	// Check if we've exceeded the end date or count
	if t.RepeatEndDate != nil && nextDate.After(*t.RepeatEndDate) {
		return nil
	}
	if t.RepeatCount > 0 && t.InstanceCount >= t.RepeatCount {
		return nil
	}

	return &nextDate
}

// getNextWeeklyOccurrence handles weekly recurrence with specific days
func (t *Task) getNextWeeklyOccurrence(fromDate time.Time) time.Time {
	daysOfWeek := t.GetRepeatDaysOfWeek()
	if len(daysOfWeek) == 0 {
		return fromDate.AddDate(0, 0, 7*t.RepeatInterval)
	}

	// Map day names to weekday numbers
	dayMap := map[string]time.Weekday{
		"sun": time.Sunday, "mon": time.Monday, "tue": time.Tuesday,
		"wed": time.Wednesday, "thu": time.Thursday, "fri": time.Friday, "sat": time.Saturday,
	}

	currentWeekday := fromDate.Weekday()
	var nextDate time.Time

	// Find the next occurrence
	for _, dayName := range daysOfWeek {
		if targetDay, exists := dayMap[strings.ToLower(dayName)]; exists {
			daysUntil := int(targetDay - currentWeekday)
			if daysUntil <= 0 {
				daysUntil += 7 // Move to next week
			}
			candidateDate := fromDate.AddDate(0, 0, daysUntil)
			if nextDate.IsZero() || candidateDate.Before(nextDate) {
				nextDate = candidateDate
			}
		}
	}

	return nextDate
}

// getNextMonthlyOccurrence handles monthly recurrence
func (t *Task) getNextMonthlyOccurrence(fromDate time.Time) time.Time {
	if t.RepeatDayOfMonth > 0 {
		// Use specific day of month
		year, month, _ := fromDate.Date()
		nextMonth := month + time.Month(t.RepeatInterval)
		nextYear := year
		if nextMonth > 12 {
			nextYear += int(nextMonth-1) / 12
			nextMonth = ((nextMonth - 1) % 12) + 1
		}
		return time.Date(nextYear, nextMonth, t.RepeatDayOfMonth, fromDate.Hour(), fromDate.Minute(), fromDate.Second(), fromDate.Nanosecond(), fromDate.Location())
	} else {
		// Use same day of month as original
		return fromDate.AddDate(0, t.RepeatInterval, 0)
	}
}

type Reminder struct {
	ID            string         `db:"id" json:"id" gorm:"primaryKey;type:text"`
	RealmID       string         `json:"realm_id" gorm:"not null;type:text;index"`
	Name          string         `json:"name" gorm:"not null;type:text"`
	Content       string         `db:"content" json:"content" gorm:"not null;type:text"`
	RemindTime    time.Time      `db:"remind_time" json:"remind_time" gorm:"not null"`
	Status        string         `db:"status" json:"status" gorm:"not null;type:text"`
	Tags          string         `json:"tags" gorm:"type:text"`
	RemindMethods string         `json:"remind_methods" gorm:"type:text"` // Comma-separated: email,im,webhook
	RemindTargets string         `json:"remind_targets" gorm:"type:text"` // JSON or comma-separated targets
	CreatedBy     string         `json:"created_by" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy     string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TaskReminder represents the mapping between tasks and reminders
type TaskReminder struct {
	ID         string         `json:"id" gorm:"primaryKey;type:text"`
	TaskID     string         `json:"task_id" gorm:"not null;type:text;index"`
	ReminderID string         `json:"reminder_id" gorm:"not null;type:text;index"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationships
	Task     Task     `json:"task" gorm:"foreignKey:TaskID;references:ID"`
	Reminder Reminder `json:"reminder" gorm:"foreignKey:ReminderID;references:ID"`
}

// TableName returns the table name for TaskReminder
func (TaskReminder) TableName() string {
	return "task_reminders"
}
