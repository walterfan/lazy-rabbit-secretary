package post

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// PostService handles business logic for posts
type PostService struct {
	repo *PostRepository
}

// NewPostService creates a new post service
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		repo: NewPostRepository(db),
	}
}

// CreatePostRequest represents the request to create a post
type CreatePostRequest struct {
	Title           string                 `json:"title" binding:"required" validate:"required,min=1,max=200"`
	Slug            string                 `json:"slug" validate:"max=200"`
	Content         string                 `json:"content" validate:"required,min=1"`
	Excerpt         string                 `json:"excerpt" validate:"max=500"`
	Status          models.PostStatus      `json:"status" validate:"oneof=draft pending published private scheduled"`
	Type            models.PostType        `json:"type" validate:"oneof=post page attachment revision custom"`
	Format          models.PostFormat      `json:"format" validate:"oneof=standard aside gallery link image quote status video audio chat"`
	Password        string                 `json:"password,omitempty" validate:"max=100"`
	MetaTitle       string                 `json:"meta_title" validate:"max=200"`
	MetaDescription string                 `json:"meta_description" validate:"max=500"`
	MetaKeywords    string                 `json:"meta_keywords" validate:"max=200"`
	FeaturedImage   string                 `json:"featured_image" validate:"url"`
	Categories      []string               `json:"categories"`
	Tags            []string               `json:"tags"`
	ParentID        string                 `json:"parent_id,omitempty"`
	MenuOrder       int                    `json:"menu_order"`
	IsSticky        bool                   `json:"is_sticky"`
	AllowPings      bool                   `json:"allow_pings"`
	CommentStatus   string                 `json:"comment_status" validate:"oneof=open closed registration_required"`
	ScheduledFor    *time.Time             `json:"scheduled_for,omitempty"`
	CustomFields    map[string]interface{} `json:"custom_fields,omitempty"`
}

// UpdatePostRequest represents the request to update a post
type UpdatePostRequest struct {
	Title           string                 `json:"title" validate:"min=1,max=200"`
	Slug            string                 `json:"slug" validate:"max=200"`
	Content         string                 `json:"content" validate:"min=1"`
	Excerpt         string                 `json:"excerpt" validate:"max=500"`
	Status          models.PostStatus      `json:"status" validate:"oneof=draft pending published private scheduled"`
	Type            models.PostType        `json:"type" validate:"oneof=post page attachment revision custom"`
	Format          models.PostFormat      `json:"format" validate:"oneof=standard aside gallery link image quote status video audio chat"`
	Password        string                 `json:"password,omitempty" validate:"max=100"`
	MetaTitle       string                 `json:"meta_title" validate:"max=200"`
	MetaDescription string                 `json:"meta_description" validate:"max=500"`
	MetaKeywords    string                 `json:"meta_keywords" validate:"max=200"`
	FeaturedImage   string                 `json:"featured_image" validate:"url"`
	Categories      []string               `json:"categories"`
	Tags            []string               `json:"tags"`
	ParentID        string                 `json:"parent_id,omitempty"`
	MenuOrder       int                    `json:"menu_order"`
	IsSticky        bool                   `json:"is_sticky"`
	AllowPings      bool                   `json:"allow_pings"`
	CommentStatus   string                 `json:"comment_status" validate:"oneof=open closed registration_required"`
	ScheduledFor    *time.Time             `json:"scheduled_for,omitempty"`
	CustomFields    map[string]interface{} `json:"custom_fields,omitempty"`
}

// PostResponse represents the response for post operations
type PostResponse struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Slug            string                 `json:"slug"`
	Content         string                 `json:"content"`
	Excerpt         string                 `json:"excerpt"`
	Status          models.PostStatus      `json:"status"`
	Type            models.PostType        `json:"type"`
	Format          models.PostFormat      `json:"format"`
	Password        string                 `json:"password,omitempty"`
	MetaTitle       string                 `json:"meta_title"`
	MetaDescription string                 `json:"meta_description"`
	MetaKeywords    string                 `json:"meta_keywords"`
	FeaturedImage   string                 `json:"featured_image"`
	Categories      []string               `json:"categories"`
	Tags            []string               `json:"tags"`
	PublishedAt     *time.Time             `json:"published_at"`
	ScheduledFor    *time.Time             `json:"scheduled_for"`
	ViewCount       int64                  `json:"view_count"`
	CommentCount    int64                  `json:"comment_count"`
	ParentID        *string                `json:"parent_id"`
	MenuOrder       int                    `json:"menu_order"`
	IsSticky        bool                   `json:"is_sticky"`
	AllowPings      bool                   `json:"allow_pings"`
	CommentStatus   string                 `json:"comment_status"`
	Language        string                 `json:"language"`
	CustomFields    map[string]interface{} `json:"custom_fields,omitempty"`
	CreatedBy       string                 `json:"created_by"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedBy       string                 `json:"updated_by"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// PostListResponse represents the response for listing posts
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// CreateFromInput creates a new post from input data
func (s *PostService) CreateFromInput(req *CreatePostRequest, realmID, createdBy string) (*PostResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = s.generateSlugFromTitle(req.Title)
	}

	// Ensure slug is unique
	uniqueSlug, err := s.ensureUniqueSlug(slug, realmID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique slug: %w", err)
	}

	// Create post model
	post := &models.Post{
		ID:              uuid.New().String(),
		RealmID:         realmID,
		Title:           req.Title,
		Slug:            uniqueSlug,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		Status:          req.Status,
		Type:            req.Type,
		Format:          req.Format,
		Password:        req.Password,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		MetaKeywords:    req.MetaKeywords,
		FeaturedImage:   req.FeaturedImage,
		MenuOrder:       req.MenuOrder,
		IsSticky:        req.IsSticky,
		AllowPings:      req.AllowPings,
		CommentStatus:   req.CommentStatus,
		ScheduledFor:    req.ScheduledFor,
		Language:        "en", // Default language
		CreatedBy:       createdBy,
		CreatedAt:       time.Now(),
		UpdatedBy:       createdBy,
		UpdatedAt:       time.Now(),
	}

	// Set categories and tags
	post.SetCategories(req.Categories)
	post.SetTags(req.Tags)

	// Set parent ID if provided
	if req.ParentID != "" {
		post.ParentID = &req.ParentID
	}

	// Handle custom fields
	if len(req.CustomFields) > 0 {
		customFieldsJSON, err := s.serializeCustomFields(req.CustomFields)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize custom fields: %w", err)
		}
		post.CustomFields = customFieldsJSON
	}

	// Handle publishing logic
	if req.Status == models.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	} else if req.Status == models.PostStatusScheduled && req.ScheduledFor != nil {
		post.ScheduledFor = req.ScheduledFor
	}

	// Save to database
	if err := s.repo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return s.toResponse(post), nil
}

// UpdateFromInput updates a post from input data
func (s *PostService) UpdateFromInput(id string, req *UpdatePostRequest, updatedBy string) (*PostResponse, error) {
	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing post
	post, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Update fields if provided
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		// Ensure slug is unique
		uniqueSlug, err := s.ensureUniqueSlug(req.Slug, post.RealmID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique slug: %w", err)
		}
		post.Slug = uniqueSlug
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.Status != "" {
		oldStatus := post.Status
		post.Status = req.Status

		// Update realm ID based on new status and type
		if req.Status == "published" && post.Type == "post" {
			// Public posts use PUBLIC_REALM_ID
			post.RealmID = models.PUBLIC_REALM_ID
		} else if oldStatus == "published" && req.Status != "published" {
			// If changing from published to non-published, use DEFAULT_REALM_ID
			post.RealmID = models.DEFAULT_REALM_ID
		}

		// Handle status transitions
		if err := s.handleStatusTransition(post, oldStatus, req.Status, req.ScheduledFor); err != nil {
			return nil, fmt.Errorf("failed to handle status transition: %w", err)
		}
	}
	if req.Type != "" {
		post.Type = req.Type

		// Update realm ID based on type and status
		if post.Status == "published" && req.Type == "post" {
			// Public posts use PUBLIC_REALM_ID
			post.RealmID = models.PUBLIC_REALM_ID
		} else if post.Status == "published" && req.Type != "post" {
			// Published pages/other types use DEFAULT_REALM_ID
			post.RealmID = models.DEFAULT_REALM_ID
		}
	}
	if req.Format != "" {
		post.Format = req.Format
	}
	if req.Password != "" {
		post.Password = req.Password
	}
	if req.MetaTitle != "" {
		post.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		post.MetaDescription = req.MetaDescription
	}
	if req.MetaKeywords != "" {
		post.MetaKeywords = req.MetaKeywords
	}
	if req.FeaturedImage != "" {
		post.FeaturedImage = req.FeaturedImage
	}
	if len(req.Categories) > 0 {
		post.SetCategories(req.Categories)
	}
	if len(req.Tags) > 0 {
		post.SetTags(req.Tags)
	}
	if req.ParentID != "" {
		post.ParentID = &req.ParentID
	}
	if req.MenuOrder != 0 {
		post.MenuOrder = req.MenuOrder
	}
	post.IsSticky = req.IsSticky
	post.AllowPings = req.AllowPings
	if req.CommentStatus != "" {
		post.CommentStatus = req.CommentStatus
	}

	// Handle custom fields
	if len(req.CustomFields) > 0 {
		customFieldsJSON, err := s.serializeCustomFields(req.CustomFields)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize custom fields: %w", err)
		}
		post.CustomFields = customFieldsJSON
	}

	post.UpdatedBy = updatedBy
	post.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return s.toResponse(post), nil
}

// GetByID retrieves a post by ID
func (s *PostService) GetByID(id string) (*PostResponse, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return s.toResponse(post), nil
}

// GetBySlug retrieves a post by slug
func (s *PostService) GetBySlug(slug string, realmID string) (*PostResponse, error) {
	post, err := s.repo.GetBySlug(slug, realmID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Increment view count for published posts
	if post.Status == models.PostStatusPublished {
		_ = s.repo.IncrementViewCount(post.ID)
		post.ViewCount++ // Update local copy
	}

	return s.toResponse(post), nil
}

// List retrieves posts with pagination and filtering
func (s *PostService) List(realmID string, status models.PostStatus, postType models.PostType, page, limit int) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.List(realmID, status, postType, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return &PostListResponse{
		Posts: s.toResponseList(posts),
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// ListPublished retrieves published posts
func (s *PostService) ListPublished(realmID string, postType models.PostType, page, limit int) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.ListPublished(realmID, postType, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list published posts: %w", err)
	}

	return &PostListResponse{
		Posts: s.toResponseList(posts),
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// Search searches posts
func (s *PostService) Search(realmID string, query string, status models.PostStatus, postType models.PostType, page, limit int) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.Search(realmID, query, status, postType, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	return &PostListResponse{
		Posts: s.toResponseList(posts),
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetByCategory retrieves posts by category
func (s *PostService) GetByCategory(realmID string, category string, page, limit int) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.GetByCategory(realmID, category, models.PostStatusPublished, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by category: %w", err)
	}

	return &PostListResponse{
		Posts: s.toResponseList(posts),
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetByTag retrieves posts by tag
func (s *PostService) GetByTag(realmID string, tag string, page, limit int) (*PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	posts, total, err := s.repo.GetByTag(realmID, tag, models.PostStatusPublished, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by tag: %w", err)
	}

	return &PostListResponse{
		Posts: s.toResponseList(posts),
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// Delete deletes a post
func (s *PostService) Delete(id string) error {
	// Check if post exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("post not found")
		}
		return fmt.Errorf("failed to get post: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// Publish publishes a post
func (s *PostService) Publish(id string, updatedBy string) (*PostResponse, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if !post.CanBePublished() {
		return nil, fmt.Errorf("post cannot be published in its current status: %s", post.Status)
	}

	post.Publish()
	post.UpdatedBy = updatedBy
	post.UpdatedAt = time.Now()

	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}

	return s.toResponse(post), nil
}

// Schedule schedules a post for publication
func (s *PostService) Schedule(id string, scheduledTime time.Time, updatedBy string) (*PostResponse, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if !post.CanBePublished() {
		return nil, fmt.Errorf("post cannot be scheduled in its current status: %s", post.Status)
	}

	post.Schedule(scheduledTime)
	post.UpdatedBy = updatedBy
	post.UpdatedAt = time.Now()

	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to schedule post: %w", err)
	}

	return s.toResponse(post), nil
}

// Helper methods

func (s *PostService) validateCreateRequest(req *CreatePostRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if req.Status == "" {
		req.Status = models.PostStatusDraft
	}
	if req.Type == "" {
		req.Type = models.PostTypePost
	}
	if req.Format == "" {
		req.Format = models.PostFormatStandard
	}
	if req.CommentStatus == "" {
		req.CommentStatus = "open"
	}
	return nil
}

func (s *PostService) validateUpdateRequest(req *UpdatePostRequest) error {
	if req.Title != "" && strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if req.Content != "" && strings.TrimSpace(req.Content) == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

func (s *PostService) generateSlugFromTitle(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Basic character filtering
	var cleanSlug strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			cleanSlug.WriteRune(char)
		}
	}

	return cleanSlug.String()
}

func (s *PostService) ensureUniqueSlug(slug string, realmID string, excludeID string) (string, error) {
	originalSlug := slug
	counter := 1

	for {
		exists, err := s.repo.SlugExists(slug, realmID, excludeID)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}

		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++

		// Prevent infinite loop
		if counter > 1000 {
			return "", fmt.Errorf("unable to generate unique slug")
		}
	}
}

func (s *PostService) handleStatusTransition(post *models.Post, oldStatus, newStatus models.PostStatus, scheduledFor *time.Time) error {
	switch newStatus {
	case models.PostStatusPublished:
		if oldStatus != models.PostStatusPublished {
			now := time.Now()
			post.PublishedAt = &now
			post.ScheduledFor = nil
		}
	case models.PostStatusScheduled:
		if scheduledFor != nil {
			post.ScheduledFor = scheduledFor
			post.PublishedAt = nil
		} else {
			return fmt.Errorf("scheduled_for time is required when setting status to scheduled")
		}
	case models.PostStatusDraft, models.PostStatusPending:
		// Keep existing dates for drafts and pending
	}

	return nil
}

func (s *PostService) serializeCustomFields(fields map[string]interface{}) (string, error) {
	// In a real implementation, you'd use JSON marshaling
	// For simplicity, returning empty string here
	return "", nil
}

func (s *PostService) toResponse(post *models.Post) *PostResponse {
	response := &PostResponse{
		ID:              post.ID,
		Title:           post.Title,
		Slug:            post.Slug,
		Content:         post.Content,
		Excerpt:         post.Excerpt,
		Status:          post.Status,
		Type:            post.Type,
		Format:          post.Format,
		Password:        post.Password,
		MetaTitle:       post.MetaTitle,
		MetaDescription: post.MetaDescription,
		MetaKeywords:    post.MetaKeywords,
		FeaturedImage:   post.FeaturedImage,
		Categories:      post.GetCategories(),
		Tags:            post.GetTags(),
		PublishedAt:     post.PublishedAt,
		ScheduledFor:    post.ScheduledFor,
		ViewCount:       post.ViewCount,
		CommentCount:    post.CommentCount,
		ParentID:        post.ParentID,
		MenuOrder:       post.MenuOrder,
		IsSticky:        post.IsSticky,
		AllowPings:      post.AllowPings,
		CommentStatus:   post.CommentStatus,
		Language:        post.Language,
		CreatedBy:       post.CreatedBy,
		CreatedAt:       post.CreatedAt,
		UpdatedBy:       post.UpdatedBy,
		UpdatedAt:       post.UpdatedAt,
	}

	// Don't expose password in response for security
	if response.Password != "" {
		response.Password = "[PROTECTED]"
	}

	return response
}

func (s *PostService) toResponseList(posts []*models.Post) []PostResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = *s.toResponse(post)
	}
	return responses
}
