package models

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a book entity with tags and borrowing information
type Book struct {
	ID          string  `json:"id" gorm:"primaryKey;type:text"`
	RealmID     string  `json:"realm_id" gorm:"not null;type:text;index"`
	ISBN        string  `json:"isbn" gorm:"uniqueIndex;not null;type:text"`
	Name        string  `json:"name" gorm:"not null;type:text"`
	Author      string  `json:"author" gorm:"type:text"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price"`

	BorrowTime *time.Time `json:"borrow_time"`
	ReturnTime *time.Time `json:"return_time"`
	Deadline   *time.Time `json:"deadline"`

	// ensure book tags are deleted if the book is deleted
	Tags []BookTag `json:"tags" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BookTag represents a tag for a book
type BookTag struct {
	ID     string `json:"id" gorm:"primaryKey;type:text"`
	BookID string `json:"book_id" gorm:"not null;type:text;index"`
	Name   string `json:"name" gorm:"not null;type:text;index"`

	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
