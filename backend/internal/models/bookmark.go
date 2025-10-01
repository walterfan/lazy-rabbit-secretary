package models

import (
	"time"

	"gorm.io/gorm"
)

// Bookmark represents a bookmark entity with tags and borrowing information
type Bookmark struct {
	ID          string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID     string         `json:"realm_id" gorm:"not null;type:text;index"`
	URL         string         `json:"url" gorm:"not null;type:text"`
	Title       string         `json:"title" gorm:"not null;type:text"`
	Description string         `json:"description" gorm:"type:text"`
	Tags        []BookmarkTag  `json:"tags" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CategoryID  string         `json:"category_id" gorm:"type:text;index"`
	CreatedBy   string         `json:"created_by" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy   string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BookmarkTag represents a tag for a bookmark
type BookmarkTag struct {
	ID         string `json:"id" gorm:"primaryKey;type:text"`
	BookmarkID string `json:"bookmark_id" gorm:"not null;type:text;index"`
	Name       string `json:"name" gorm:"not null;type:text;index"`
}

// BookmarkCategory represents a category for a bookmark
type BookmarkCategory struct {
	ID       int64  `json:"id" gorm:"primaryKey;type:text"`
	Name     string `json:"name" gorm:"not null;type:text"`
	ParentID *int64 `json:"parent_id" gorm:"type:text;index"`

	CreatedBy string `json:"created_by" gorm:"type:text"`
	UpdatedBy string `json:"updated_by" gorm:"type:text"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
