package models

import (
	"time"

	"gorm.io/gorm"
)

// InboxItem represents a quick note/idea captured in the inbox
type InboxItem struct {
	ID          string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID     string         `json:"realm_id" gorm:"not null;type:text;index"`
	Title       string         `json:"title" gorm:"not null;type:text"`
	Description string         `json:"description" gorm:"type:text"`
	Priority    string         `json:"priority" gorm:"type:text;default:'normal'"` // low, normal, high, urgent
	Status      string         `json:"status" gorm:"type:text;default:'pending'"`  // pending, processing, completed, archived
	Tags        string         `json:"tags" gorm:"type:text"`                      // comma-separated tags
	Context     string         `json:"context" gorm:"type:text"`                   // @home, @office, @phone, etc.
	CreatedBy   string         `json:"created_by" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy   string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// DailyChecklistItem represents a planned task for the day
type DailyChecklistItem struct {
	ID             string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID        string         `json:"realm_id" gorm:"not null;type:text;index"`
	Title          string         `json:"title" gorm:"not null;type:text"`
	Description    string         `json:"description" gorm:"type:text"`
	Priority       string         `json:"priority" gorm:"type:text;default:'B'"`    // A, B+, B, C, D
	EstimatedTime  int            `json:"estimated_time" gorm:"type:int;default:0"` // in minutes
	Deadline       *time.Time     `json:"deadline" gorm:"type:text"`
	Context        string         `json:"context" gorm:"type:text"`                  // @home, @office, @phone, etc.
	Status         string         `json:"status" gorm:"type:text;default:'pending'"` // pending, in_progress, completed, cancelled
	CompletionTime *time.Time     `json:"completion_time" gorm:"type:text"`
	ActualTime     int            `json:"actual_time" gorm:"type:int;default:0"` // actual time spent in minutes
	Notes          string         `json:"notes" gorm:"type:text"`                // notes about execution
	InboxItemID    *string        `json:"inbox_item_id" gorm:"type:text;index"`  // reference to source inbox item
	Date           time.Time      `json:"date" gorm:"type:text;index"`           // the date this task is planned for
	CreatedBy      string         `json:"created_by" gorm:"type:text"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy      string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// InboxItemResponse represents the response format for inbox items
type InboxItemResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	Tags        string    `json:"tags"`
	Context     string    `json:"context"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DailyChecklistItemResponse represents the response format for daily checklist items
type DailyChecklistItemResponse struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`
	EstimatedTime  int        `json:"estimated_time"`
	Deadline       *time.Time `json:"deadline"`
	Context        string     `json:"context"`
	Status         string     `json:"status"`
	CompletionTime *time.Time `json:"completion_time"`
	ActualTime     int        `json:"actual_time"`
	Notes          string     `json:"notes"`
	InboxItemID    *string    `json:"inbox_item_id"`
	Date           time.Time  `json:"date"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
