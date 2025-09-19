package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// PostStatus represents the status of a post
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPending   PostStatus = "pending" // awaiting review
	PostStatusPublished PostStatus = "published"
	PostStatusPrivate   PostStatus = "private"
	PostStatusTrash     PostStatus = "trash"
	PostStatusScheduled PostStatus = "scheduled" // scheduled for future publication
)

// PostType represents the type of post content
type PostType string

const (
	PostTypePost       PostType = "post"       // regular blog post
	PostTypePage       PostType = "page"       // static page
	PostTypeAttachment PostType = "attachment" // media file
	PostTypeRevision   PostType = "revision"   // post revision
	PostTypeCustom     PostType = "custom"     // custom post type
)

// PostFormat represents the format of the post content
type PostFormat string

const (
	PostFormatStandard PostFormat = "standard" // default
	PostFormatAside    PostFormat = "aside"    // short note
	PostFormatGallery  PostFormat = "gallery"  // image gallery
	PostFormatLink     PostFormat = "link"     // external link
	PostFormatImage    PostFormat = "image"    // single image
	PostFormatQuote    PostFormat = "quote"    // quotation
	PostFormatStatus   PostFormat = "status"   // short status update
	PostFormatVideo    PostFormat = "video"    // video content
	PostFormatAudio    PostFormat = "audio"    // audio content
	PostFormatChat     PostFormat = "chat"     // chat transcript
)

// Post represents a blog post or page in WordPress style
type Post struct {
	ID       string     `json:"id" gorm:"primaryKey;type:text"`
	RealmID  string     `json:"realm_id" gorm:"not null;type:text;index"`
	Title    string     `json:"title" gorm:"not null;type:text"`
	Slug     string     `json:"slug" gorm:"not null;type:text;uniqueIndex:idx_realm_slug"`
	Content  string     `json:"content" gorm:"type:text"`
	Excerpt  string     `json:"excerpt" gorm:"type:text"` // short summary
	Status   PostStatus `json:"status" gorm:"default:'draft';type:text;index"`
	Type     PostType   `json:"type" gorm:"default:'post';type:text;index"`
	Format   PostFormat `json:"format" gorm:"default:'standard';type:text"`
	Password string     `json:"password,omitempty" gorm:"type:text"` // for password-protected posts

	// SEO and metadata
	MetaTitle       string `json:"meta_title" gorm:"type:text"`       // SEO title
	MetaDescription string `json:"meta_description" gorm:"type:text"` // SEO description
	MetaKeywords    string `json:"meta_keywords" gorm:"type:text"`    // SEO keywords
	FeaturedImage   string `json:"featured_image" gorm:"type:text"`   // featured image URL

	// Organization
	Categories string `json:"categories" gorm:"type:text"` // comma-separated category names
	Tags       string `json:"tags" gorm:"type:text"`       // comma-separated tags

	// Publishing
	PublishedAt  *time.Time `json:"published_at" gorm:"index"`      // when published
	ScheduledFor *time.Time `json:"scheduled_for" gorm:"index"`     // scheduled publication time
	ViewCount    int64      `json:"view_count" gorm:"default:0"`    // page views
	CommentCount int64      `json:"comment_count" gorm:"default:0"` // number of comments

	// Hierarchy (for pages and nested content)
	ParentID   *string `json:"parent_id" gorm:"type:text;index"`  // parent post ID
	MenuOrder  int     `json:"menu_order" gorm:"default:0"`       // display order
	IsSticky   bool    `json:"is_sticky" gorm:"default:false"`    // sticky post
	AllowPings bool    `json:"allow_pings" gorm:"default:true"`   // allow pingbacks/trackbacks
	PingStatus string  `json:"ping_status" gorm:"default:'open'"` // open, closed

	// Comments
	CommentStatus string `json:"comment_status" gorm:"default:'open'"` // open, closed, registration_required

	// Content management
	ContentFiltered string `json:"content_filtered,omitempty" gorm:"type:text"` // filtered content for search
	GUID            string `json:"guid,omitempty" gorm:"type:text"`             // globally unique identifier

	// Revision tracking
	ParentRevisionID *string `json:"parent_revision_id,omitempty" gorm:"type:text;index"` // for revisions
	RevisionCount    int     `json:"revision_count" gorm:"default:0"`                     // number of revisions

	// WordPress-style custom fields support
	CustomFields string `json:"custom_fields,omitempty" gorm:"type:text"` // JSON string of custom fields

	// Multi-language support
	Language      string  `json:"language" gorm:"default:'en';type:text"`          // content language
	TranslationOf *string `json:"translation_of,omitempty" gorm:"type:text;index"` // original post ID for translations

	// Social sharing
	ShareCount int64  `json:"share_count" gorm:"default:0"`           // social shares
	LikeCount  int64  `json:"like_count" gorm:"default:0"`            // likes
	SocialMeta string `json:"social_meta,omitempty" gorm:"type:text"` // JSON for social media metadata

	// Analytics and tracking
	AnalyticsData string `json:"analytics_data,omitempty" gorm:"type:text"` // JSON for analytics data

	// Audit fields (matching task.go pattern)
	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (Post) TableName() string {
	return "posts"
}

// IsPublished returns true if the post is published
func (p *Post) IsPublished() bool {
	return p.Status == PostStatusPublished &&
		(p.PublishedAt == nil || p.PublishedAt.Before(time.Now()))
}

// IsScheduled returns true if the post is scheduled for future publication
func (p *Post) IsScheduled() bool {
	return p.Status == PostStatusScheduled &&
		p.ScheduledFor != nil &&
		p.ScheduledFor.After(time.Now())
}

// IsPasswordProtected returns true if the post requires a password
func (p *Post) IsPasswordProtected() bool {
	return p.Password != ""
}

// IsPage returns true if this is a page rather than a post
func (p *Post) IsPage() bool {
	return p.Type == PostTypePage
}

// IsRevision returns true if this is a revision of another post
func (p *Post) IsRevision() bool {
	return p.Type == PostTypeRevision || p.ParentRevisionID != nil
}

// HasParent returns true if this post has a parent
func (p *Post) HasParent() bool {
	return p.ParentID != nil
}

// GetCategories returns the categories as a slice
func (p *Post) GetCategories() []string {
	if p.Categories == "" {
		return []string{}
	}
	categories := strings.Split(p.Categories, ",")
	result := make([]string, 0, len(categories))
	for _, category := range categories {
		if trimmed := strings.TrimSpace(category); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetTags returns the tags as a slice
func (p *Post) GetTags() []string {
	if p.Tags == "" {
		return []string{}
	}
	tags := strings.Split(p.Tags, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// SetCategories sets categories from a slice
func (p *Post) SetCategories(categories []string) {
	if len(categories) == 0 {
		p.Categories = ""
		return
	}
	// Clean and join categories
	cleaned := make([]string, 0, len(categories))
	for _, category := range categories {
		if trimmed := strings.TrimSpace(category); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	p.Categories = strings.Join(cleaned, ",")
}

// SetTags sets tags from a slice
func (p *Post) SetTags(tags []string) {
	if len(tags) == 0 {
		p.Tags = ""
		return
	}
	// Clean and join tags
	cleaned := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	p.Tags = strings.Join(cleaned, ",")
}

// GenerateSlug creates a URL-friendly slug from the title
func (p *Post) GenerateSlug() {
	if p.Title == "" {
		return
	}

	slug := strings.ToLower(p.Title)
	// Replace spaces and special characters with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Basic character filtering (keep only alphanumeric and hyphens)
	var cleanSlug strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			cleanSlug.WriteRune(char)
		}
	}

	p.Slug = cleanSlug.String()
}

// IncrementViewCount increments the view count
func (p *Post) IncrementViewCount() {
	p.ViewCount++
}

// IncrementCommentCount increments the comment count
func (p *Post) IncrementCommentCount() {
	p.CommentCount++
}

// DecrementCommentCount decrements the comment count
func (p *Post) DecrementCommentCount() {
	if p.CommentCount > 0 {
		p.CommentCount--
	}
}

// CanBePublished returns true if the post can be published
func (p *Post) CanBePublished() bool {
	return p.Status == PostStatusDraft ||
		p.Status == PostStatusPending ||
		p.Status == PostStatusScheduled
}

// Publish sets the post status to published and sets the published date
func (p *Post) Publish() {
	if p.CanBePublished() {
		p.Status = PostStatusPublished
		now := time.Now()
		p.PublishedAt = &now
		p.ScheduledFor = nil // Clear scheduled time
	}
}

// Schedule schedules the post for future publication
func (p *Post) Schedule(scheduledTime time.Time) {
	if p.CanBePublished() {
		p.Status = PostStatusScheduled
		p.ScheduledFor = &scheduledTime
		p.PublishedAt = nil // Clear published time
	}
}

// Unpublish sets the post back to draft status
func (p *Post) Unpublish() {
	if p.Status == PostStatusPublished {
		p.Status = PostStatusDraft
		p.PublishedAt = nil
		p.ScheduledFor = nil
	}
}

// PostMeta represents additional metadata for posts (WordPress-style meta)
type PostMeta struct {
	ID        string `json:"id" gorm:"primaryKey;type:text"`
	PostID    string `json:"post_id" gorm:"not null;type:text;index"`
	MetaKey   string `json:"meta_key" gorm:"not null;type:text;index"`
	MetaValue string `json:"meta_value" gorm:"type:text"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationship
	Post Post `json:"post" gorm:"foreignKey:PostID;references:ID"`
}

// TableName specifies the table name for GORM
func (PostMeta) TableName() string {
	return "post_meta"
}

// PostCategory represents post categories (WordPress-style taxonomy)
type PostCategory struct {
	ID          string  `json:"id" gorm:"primaryKey;type:text"`
	Name        string  `json:"name" gorm:"not null;type:text"`
	Slug        string  `json:"slug" gorm:"not null;type:text;uniqueIndex"`
	Description string  `json:"description" gorm:"type:text"`
	ParentID    *string `json:"parent_id" gorm:"type:text;index"` // for hierarchical categories
	Count       int64   `json:"count" gorm:"default:0"`           // number of posts in this category

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (PostCategory) TableName() string {
	return "post_categories"
}

// PostTag represents post tags (WordPress-style taxonomy)
type PostTag struct {
	ID          string `json:"id" gorm:"primaryKey;type:text"`
	Name        string `json:"name" gorm:"not null;type:text"`
	Slug        string `json:"slug" gorm:"not null;type:text;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	Count       int64  `json:"count" gorm:"default:0"` // number of posts with this tag

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (PostTag) TableName() string {
	return "post_tags"
}

// PostCategoryRelation represents the many-to-many relationship between posts and categories
type PostCategoryRelation struct {
	ID         string `json:"id" gorm:"primaryKey;type:text"`
	PostID     string `json:"post_id" gorm:"not null;type:text;index"`
	CategoryID string `json:"category_id" gorm:"not null;type:text;index"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationships
	Post     Post         `json:"post" gorm:"foreignKey:PostID;references:ID"`
	Category PostCategory `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
}

// TableName specifies the table name for GORM
func (PostCategoryRelation) TableName() string {
	return "post_category_relations"
}

// PostTagRelation represents the many-to-many relationship between posts and tags
type PostTagRelation struct {
	ID     string `json:"id" gorm:"primaryKey;type:text"`
	PostID string `json:"post_id" gorm:"not null;type:text;index"`
	TagID  string `json:"tag_id" gorm:"not null;type:text;index"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationships
	Post Post    `json:"post" gorm:"foreignKey:PostID;references:ID"`
	Tag  PostTag `json:"tag" gorm:"foreignKey:TagID;references:ID"`
}

// TableName specifies the table name for GORM
func (PostTagRelation) TableName() string {
	return "post_tag_relations"
}

// Comment represents a comment on a post (WordPress-style)
type Comment struct {
	ID       string  `json:"id" gorm:"primaryKey;type:text"`
	PostID   string  `json:"post_id" gorm:"not null;type:text;index"`
	ParentID *string `json:"parent_id" gorm:"type:text;index"` // for threaded comments

	Author      string `json:"author" gorm:"not null;type:text"`
	AuthorEmail string `json:"author_email" gorm:"type:text"`
	AuthorURL   string `json:"author_url" gorm:"type:text"`
	AuthorIP    string `json:"author_ip" gorm:"type:text"`

	Content   string `json:"content" gorm:"not null;type:text"`
	Approved  bool   `json:"approved" gorm:"default:false"`
	Type      string `json:"type" gorm:"default:'comment';type:text"` // comment, pingback, trackback
	UserAgent string `json:"user_agent" gorm:"type:text"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationship
	Post Post `json:"post" gorm:"foreignKey:PostID;references:ID"`
}

// TableName specifies the table name for GORM
func (Comment) TableName() string {
	return "comments"
}

// IsApproved returns true if the comment is approved
func (c *Comment) IsApproved() bool {
	return c.Approved
}

// IsReply returns true if this is a reply to another comment
func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}

// Approve approves the comment
func (c *Comment) Approve() {
	c.Approved = true
}

// Unapprove unapproves the comment
func (c *Comment) Unapprove() {
	c.Approved = false
}
