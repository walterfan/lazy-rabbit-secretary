package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// WikiPageStatus represents the status of a wiki page
type WikiPageStatus string

const (
	WikiPageStatusDraft     WikiPageStatus = "draft"
	WikiPageStatusPublished WikiPageStatus = "published"
	WikiPageStatusArchived  WikiPageStatus = "archived"
	WikiPageStatusProtected WikiPageStatus = "protected" // requires login
)

// WikiPageType represents the type of wiki content
type WikiPageType string

const (
	WikiPageTypeArticle  WikiPageType = "article"  // regular wiki article
	WikiPageTypeTemplate WikiPageType = "template" // reusable template
	WikiPageTypeCategory WikiPageType = "category" // category page
	WikiPageTypeRedirect WikiPageType = "redirect" // redirect page
	WikiPageTypeStub     WikiPageType = "stub"     // stub article
)

// WikiPage represents a wiki page/article
type WikiPage struct {
	ID      string `json:"id" gorm:"primaryKey;type:text"`
	RealmID string `json:"realm_id" gorm:"not null;type:text;index"`
	Title   string `json:"title" gorm:"not null;type:text"`
	Slug    string `json:"slug" gorm:"not null;type:text;uniqueIndex:idx_wiki_realm_slug"`
	Content string `json:"content" gorm:"type:text"`
	Summary string `json:"summary" gorm:"type:text"` // short summary

	// Wiki-specific fields
	Status    WikiPageStatus `json:"status" gorm:"default:'draft';type:text;index"`
	Type      WikiPageType   `json:"type" gorm:"default:'article';type:text;index"`
	Protected bool           `json:"is_protected" gorm:"default:false"` // requires login
	Locked    bool           `json:"is_locked" gorm:"default:false"`    // prevents editing

	// Organization
	Categories string `json:"categories" gorm:"type:text"` // comma-separated
	Tags       string `json:"tags" gorm:"type:text"`       // comma-separated

	// Hierarchy
	ParentID  *string `json:"parent_id" gorm:"type:text;index"` // for sub-pages
	MenuOrder int     `json:"menu_order" gorm:"default:0"`

	// Redirect support
	RedirectTo *string `json:"redirect_to" gorm:"type:text"` // for redirect pages

	// Statistics
	ViewCount int64 `json:"view_count" gorm:"default:0"`
	EditCount int64 `json:"edit_count" gorm:"default:0"`

	// Version control
	CurrentVersion int `json:"current_version" gorm:"default:1"`

	// Multi-language support
	Language      string  `json:"language" gorm:"default:'en';type:text"`
	TranslationOf *string `json:"translation_of,omitempty" gorm:"type:text;index"`

	// Audit fields
	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (*WikiPage) TableName() string {
	return "wiki_pages"
}

// WikiRevision represents a version of a wiki page
type WikiRevision struct {
	ID         string `json:"id" gorm:"primaryKey;type:text"`
	PageID     string `json:"page_id" gorm:"not null;type:text;index"`
	Version    int    `json:"version" gorm:"not null;index"`
	Title      string `json:"title" gorm:"type:text"`
	Content    string `json:"content" gorm:"type:text"`
	Summary    string `json:"summary" gorm:"type:text"`
	ChangeNote string `json:"change_note" gorm:"type:text"` // edit summary

	// Change tracking
	ContentSize int `json:"content_size" gorm:"default:0"` // bytes
	LineCount   int `json:"line_count" gorm:"default:0"`

	// Audit fields
	CreatedBy string    `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Foreign key relationship
	Page WikiPage `json:"page" gorm:"foreignKey:PageID;references:ID"`
}

// TableName specifies the table name for GORM
func (*WikiRevision) TableName() string {
	return "wiki_revisions"
}

// IsPublished returns true if the wiki page is published
func (w *WikiPage) IsPublished() bool {
	return w.Status == WikiPageStatusPublished
}

// IsProtected returns true if the wiki page requires authentication
func (w *WikiPage) IsProtected() bool {
	return w.Protected
}

// IsLocked returns true if the wiki page is locked from editing
func (w *WikiPage) IsLocked() bool {
	return w.Locked
}

// IsRedirect returns true if this is a redirect page
func (w *WikiPage) IsRedirect() bool {
	return w.Type == WikiPageTypeRedirect && w.RedirectTo != nil
}

// IsStub returns true if this is a stub article
func (w *WikiPage) IsStub() bool {
	return w.Type == WikiPageTypeStub
}

// HasParent returns true if this page has a parent
func (w *WikiPage) HasParent() bool {
	return w.ParentID != nil
}

// GetCategories returns the categories as a slice
func (w *WikiPage) GetCategories() []string {
	if w.Categories == "" {
		return []string{}
	}
	categories := strings.Split(w.Categories, ",")
	result := make([]string, 0, len(categories))
	for _, category := range categories {
		if trimmed := strings.TrimSpace(category); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetTags returns the tags as a slice
func (w *WikiPage) GetTags() []string {
	if w.Tags == "" {
		return []string{}
	}
	tags := strings.Split(w.Tags, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// SetCategories sets categories from a slice
func (w *WikiPage) SetCategories(categories []string) {
	if len(categories) == 0 {
		w.Categories = ""
		return
	}
	w.Categories = strings.Join(categories, ",")
}

// SetTags sets tags from a slice
func (w *WikiPage) SetTags(tags []string) {
	if len(tags) == 0 {
		w.Tags = ""
		return
	}
	w.Tags = strings.Join(tags, ",")
}

// IncrementViewCount increments the view count
func (w *WikiPage) IncrementViewCount() {
	w.ViewCount++
}

// IncrementEditCount increments the edit count
func (w *WikiPage) IncrementEditCount() {
	w.EditCount++
	w.CurrentVersion++
}

// CanBeEdited returns true if the page can be edited
func (w *WikiPage) CanBeEdited() bool {
	return !w.Locked && w.Status != WikiPageStatusArchived
}

// CanBeDeleted returns true if the page can be deleted
func (w *WikiPage) CanBeDeleted() bool {
	return w.Status != WikiPageStatusArchived
}

// Lock locks the page from editing
func (w *WikiPage) Lock() {
	w.Locked = true
}

// Unlock unlocks the page for editing
func (w *WikiPage) Unlock() {
	w.Locked = false
}

// Archive archives the page
func (w *WikiPage) Archive() {
	w.Status = WikiPageStatusArchived
}

// Publish publishes the page
func (w *WikiPage) Publish() {
	w.Status = WikiPageStatusPublished
}

// Unpublish sets the page back to draft status
func (w *WikiPage) Unpublish() {
	w.Status = WikiPageStatusDraft
}

// Protect makes the page require authentication
func (w *WikiPage) Protect() {
	w.Protected = true
}

// Unprotect removes authentication requirement
func (w *WikiPage) Unprotect() {
	w.Protected = false
}

// GetContentSize returns the size of the content in bytes
func (r *WikiRevision) GetContentSize() int {
	return len(r.Content)
}

// GetLineCount returns the number of lines in the content
func (r *WikiRevision) GetLineCount() int {
	if r.Content == "" {
		return 0
	}
	return strings.Count(r.Content, "\n") + 1
}
