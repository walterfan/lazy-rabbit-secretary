package models

import (
	"time"

	"gorm.io/gorm"
)

// Prompt represents the prompt table in the database
type Prompt struct {
	ID           string         `json:"id" gorm:"primaryKey;type:text"`
	Name         string         `json:"name" gorm:"index;not null;type:text"`
	Description  string         `json:"description" gorm:"type:text"`
	SystemPrompt string         `json:"system_prompt" gorm:"column:system_prompt;type:text"`
	UserPrompt   string         `json:"user_prompt" gorm:"column:user_prompt;type:text"`
	Tags         string         `json:"tags" gorm:"type:text"`
	CreatedBy    string         `json:"created_by" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy    string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
