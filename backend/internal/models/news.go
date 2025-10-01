package models

import (
	"time"

	"gorm.io/gorm"
)

// NewsItem represents a simple message board item
type NewsItem struct {
	ID        string         `json:"id" gorm:"type:text;primaryKey"`
	RealmID   string         `json:"realm_id" gorm:"type:text;not null;index:idx_news_realm"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	Author    string         `json:"author" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:datetime;index:idx_news_created"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:datetime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index:idx_news_deleted"`
}

// TableName returns the table name for NewsItem
func (NewsItem) TableName() string {
	return "news_items"
}
