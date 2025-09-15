package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents the role table in the database
type Role struct {
	ID          string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID     string         `json:"realm_id" gorm:"not null;type:text;index"`
	Name        string         `json:"name" gorm:"not null;type:text"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedBy   string         `json:"created_by" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy   string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
