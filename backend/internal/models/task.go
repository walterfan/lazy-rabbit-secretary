package models

import (
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
	ID           string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID      string         `json:"realm_id" gorm:"not null;type:text;index"`
	Name         string         `json:"name" gorm:"not null;type:text"`
	Description  string         `json:"description" gorm:"type:text"`
	Priority     int            `json:"priority" gorm:"not null;default:2"`
	Difficulty   int            `json:"difficulty" gorm:"not null;default:2"`
	Status       TaskStatus     `json:"status" gorm:"default:'pending';type:text;index"`
	ScheduleTime time.Time      `json:"schedule_time" gorm:"not null"`
	Minutes      int            `json:"minutes" gorm:"not null"`
	Deadline     time.Time      `json:"deadline" gorm:"not null"`
	StartTime    *time.Time     `json:"start_time"`
	EndTime      *time.Time     `json:"end_time"`
	Tags         string         `json:"tags" gorm:"type:text"`
	CreatedBy    string         `json:"created_by" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy    string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Reminder struct {
	ID         string         `db:"id" json:"id"`
	RealmID    string         `json:"realm_id" gorm:"not null;type:text;index"`
	Name       string         `json:"name" gorm:"not null;type:text"`
	Content    string         `db:"content" json:"content"`
	RemindTime time.Time      `db:"remind_time" json:"remind_time"`
	Status     string         `db:"status" json:"status"`
	Tags       string         `json:"tags" gorm:"type:text"`
	CreatedBy  string         `json:"created_by" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy  string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
