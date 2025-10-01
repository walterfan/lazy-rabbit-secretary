package post

import (
	"fmt"

	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/internal/prompt"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/llm"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PostService handles business logic for posts
type PostService struct {
	repo          *PostRepository
	promptService *prompt.PromptService
	logger        *zap.SugaredLogger
}

// NewPostService creates a new post service
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		repo:          NewPostRepository(db),
		promptService: prompt.NewPromptService(db),
		logger:        log.GetLogger(),
	}
}

// BasePostRequest contains common fields for all post operations
type BasePostRequest struct {
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

// CreatePostRequest represents the request to create a post
type CreatePostRequest struct {
	BasePostRequest
	// Override validation for required fields in create
	Title   string `json:"title" binding:"required" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// UpdatePostRequest represents the request to update a post
type UpdatePostRequest struct {
	BasePostRequest
}

// RefineRequest represents the request to refine a post
type RefineRequest struct {
	BasePostRequest
	Action      string `json:"action" binding:"required" validate:"required,oneof=improve_writing make_shorter make_longer change_tone translate"`
	Requirement string `json:"requirement,omitempty" validate:"max=1000"`
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
	slug := req.BasePostRequest.Slug
	if slug == "" {
		slug = s.generateSlugFromTitle(req.BasePostRequest.Title)
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
		Title:           req.BasePostRequest.Title,
		Slug:            uniqueSlug,
		Content:         req.BasePostRequest.Content,
		Excerpt:         req.BasePostRequest.Excerpt,
		Status:          req.BasePostRequest.Status,
		Type:            req.BasePostRequest.Type,
		Format:          req.BasePostRequest.Format,
		Password:        req.BasePostRequest.Password,
		MetaTitle:       req.BasePostRequest.MetaTitle,
		MetaDescription: req.BasePostRequest.MetaDescription,
		MetaKeywords:    req.BasePostRequest.MetaKeywords,
		FeaturedImage:   req.BasePostRequest.FeaturedImage,
		MenuOrder:       req.BasePostRequest.MenuOrder,
		IsSticky:        req.BasePostRequest.IsSticky,
		AllowPings:      req.BasePostRequest.AllowPings,
		CommentStatus:   req.BasePostRequest.CommentStatus,
		ScheduledFor:    req.BasePostRequest.ScheduledFor,
		Language:        "en", // Default language
		CreatedBy:       createdBy,
		CreatedAt:       time.Now(),
		UpdatedBy:       createdBy,
		UpdatedAt:       time.Now(),
	}

	// Set categories and tags
	post.SetCategories(req.BasePostRequest.Categories)
	post.SetTags(req.BasePostRequest.Tags)

	// Set parent ID if provided
	if req.BasePostRequest.ParentID != "" {
		post.ParentID = &req.BasePostRequest.ParentID
	}

	// Handle custom fields
	if len(req.BasePostRequest.CustomFields) > 0 {
		customFieldsJSON, err := s.serializeCustomFields(req.BasePostRequest.CustomFields)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize custom fields: %w", err)
		}
		post.CustomFields = customFieldsJSON
	}

	// Handle publishing logic
	if req.BasePostRequest.Status == models.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	} else if req.BasePostRequest.Status == models.PostStatusScheduled && req.BasePostRequest.ScheduledFor != nil {
		post.ScheduledFor = req.BasePostRequest.ScheduledFor
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
	if req.BasePostRequest.Title != "" {
		post.Title = req.BasePostRequest.Title
	}
	if req.BasePostRequest.Slug != "" {
		// Ensure slug is unique
		uniqueSlug, err := s.ensureUniqueSlug(req.BasePostRequest.Slug, post.RealmID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique slug: %w", err)
		}
		post.Slug = uniqueSlug
	}
	if req.BasePostRequest.Content != "" {
		post.Content = req.BasePostRequest.Content
	}
	if req.BasePostRequest.Excerpt != "" {
		post.Excerpt = req.BasePostRequest.Excerpt
	}
	if req.BasePostRequest.Status != "" {
		oldStatus := post.Status
		post.Status = req.BasePostRequest.Status

		// Handle status transitions
		if err := s.handleStatusTransition(post, oldStatus, req.BasePostRequest.Status, req.BasePostRequest.ScheduledFor); err != nil {
			return nil, fmt.Errorf("failed to handle status transition: %w", err)
		}
	}
	if req.BasePostRequest.Type != "" {
		post.Type = req.BasePostRequest.Type
	}
	if req.BasePostRequest.Format != "" {
		post.Format = req.BasePostRequest.Format
	}
	if req.BasePostRequest.Password != "" {
		post.Password = req.BasePostRequest.Password
	}
	if req.BasePostRequest.MetaTitle != "" {
		post.MetaTitle = req.BasePostRequest.MetaTitle
	}
	if req.BasePostRequest.MetaDescription != "" {
		post.MetaDescription = req.BasePostRequest.MetaDescription
	}
	if req.BasePostRequest.MetaKeywords != "" {
		post.MetaKeywords = req.BasePostRequest.MetaKeywords
	}
	if req.BasePostRequest.FeaturedImage != "" {
		post.FeaturedImage = req.BasePostRequest.FeaturedImage
	}
	if len(req.BasePostRequest.Categories) > 0 {
		post.SetCategories(req.BasePostRequest.Categories)
	}
	if len(req.BasePostRequest.Tags) > 0 {
		post.SetTags(req.BasePostRequest.Tags)
	}
	if req.BasePostRequest.ParentID != "" {
		post.ParentID = &req.BasePostRequest.ParentID
	}
	if req.BasePostRequest.MenuOrder != 0 {
		post.MenuOrder = req.BasePostRequest.MenuOrder
	}
	post.IsSticky = req.BasePostRequest.IsSticky
	post.AllowPings = req.BasePostRequest.AllowPings
	if req.BasePostRequest.CommentStatus != "" {
		post.CommentStatus = req.BasePostRequest.CommentStatus
	}

	// Handle custom fields
	if len(req.BasePostRequest.CustomFields) > 0 {
		customFieldsJSON, err := s.serializeCustomFields(req.BasePostRequest.CustomFields)
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

// RefinePost refines a post based on the specified action
func (s *PostService) RefinePost(id string, req *RefineRequest, updatedBy string) (*PostResponse, error) {
	// Validate request
	if err := s.validateRefineRequest(req); err != nil {
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

	// Apply refinement based on action using LLM
	refinedContent, err := s.refineContentWithLLM(post.Content, req.Action, req.Requirement)
	if err != nil {
		return nil, fmt.Errorf("failed to refine content: %w", err)
	}
	post.Content = refinedContent

	// Update other fields if provided
	if req.BasePostRequest.Title != "" {
		post.Title = req.BasePostRequest.Title
	}
	if req.BasePostRequest.Excerpt != "" {
		post.Excerpt = req.BasePostRequest.Excerpt
	}
	if req.BasePostRequest.MetaTitle != "" {
		post.MetaTitle = req.BasePostRequest.MetaTitle
	}
	if req.BasePostRequest.MetaDescription != "" {
		post.MetaDescription = req.BasePostRequest.MetaDescription
	}
	if req.BasePostRequest.MetaKeywords != "" {
		post.MetaKeywords = req.BasePostRequest.MetaKeywords
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
	if strings.TrimSpace(req.BasePostRequest.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(req.BasePostRequest.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if req.BasePostRequest.Status == "" {
		req.BasePostRequest.Status = models.PostStatusDraft
	}
	if req.BasePostRequest.Type == "" {
		req.BasePostRequest.Type = models.PostTypePost
	}
	if req.BasePostRequest.Format == "" {
		req.BasePostRequest.Format = models.PostFormatStandard
	}
	if req.BasePostRequest.CommentStatus == "" {
		req.BasePostRequest.CommentStatus = "open"
	}
	return nil
}

func (s *PostService) validateUpdateRequest(req *UpdatePostRequest) error {
	if req.BasePostRequest.Title != "" && strings.TrimSpace(req.BasePostRequest.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if req.BasePostRequest.Content != "" && strings.TrimSpace(req.BasePostRequest.Content) == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

func (s *PostService) validateRefineRequest(req *RefineRequest) error {
	if req.Action == "" {
		return fmt.Errorf("action is required")
	}

	// Validate action values
	validActions := []string{"improve_writing", "make_shorter", "make_longer", "change_tone", "translate"}
	validAction := false
	for _, action := range validActions {
		if req.Action == action {
			validAction = true
			break
		}
	}
	if !validAction {
		return fmt.Errorf("invalid action: %s", req.Action)
	}

	// Validate requirement for specific actions
	if req.Action == "change_tone" && req.Requirement != "" {
		validTones := []string{"formal", "informal", "friendly", "persuasive", "serious"}
		validTone := false
		for _, tone := range validTones {
			if req.Requirement == tone {
				validTone = true
				break
			}
		}
		if !validTone {
			return fmt.Errorf("invalid tone: %s. Valid tones are: formal, informal, friendly, persuasive, serious", req.Requirement)
		}
	}

	if req.Action == "translate" && req.Requirement != "" {
		validLanguages := []string{"chinese", "english", "japanese", "spanish", "french", "german", "italian", "portuguese", "russian", "korean"}
		validLanguage := false
		for _, lang := range validLanguages {
			if req.Requirement == lang {
				validLanguage = true
				break
			}
		}
		if !validLanguage {
			return fmt.Errorf("invalid language: %s. Valid languages are: chinese, english, japanese, spanish, french, german, italian, portuguese, russian, korean", req.Requirement)
		}
	}

	return nil
}

// Refinement helper methods
func (s *PostService) improveWriting(content, requirement string) string {
	// This is a placeholder implementation
	// In a real application, you would integrate with an AI service like OpenAI
	return fmt.Sprintf("[IMPROVED] %s", content)
}

func (s *PostService) makeShorter(content, requirement string) string {
	// This is a placeholder implementation
	// In a real application, you would integrate with an AI service
	return fmt.Sprintf("[SHORTENED] %s", content)
}

func (s *PostService) makeLonger(content, requirement string) string {
	// This is a placeholder implementation
	// In a real application, you would integrate with an AI service
	return fmt.Sprintf("[LENGTHENED] %s", content)
}

func (s *PostService) changeTone(content, requirement string) string {
	// This is a placeholder implementation
	// In a real application, you would integrate with an AI service
	return fmt.Sprintf("[%s TONE] %s", strings.ToUpper(requirement), content)
}

func (s *PostService) translateContent(content, requirement string) string {
	// This is a placeholder implementation
	// In a real application, you would integrate with a translation service
	return fmt.Sprintf("[TRANSLATED TO %s] %s", strings.ToUpper(requirement), content)
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

// refineContentWithLLM uses LLM to refine content based on action and requirement
func (s *PostService) refineContentWithLLM(content, action, requirement string) (string, error) {
	// Try to get prompt from repository first
	systemPrompt, userPrompt, err := s.getRefinePrompts(action, requirement, content)
	if err != nil {
		return "", fmt.Errorf("failed to get refine prompts: %w", err)
	}
	s.logger.Infof("Refining content with LLM - Action: %s, systemPrompt: %s, userPrompt: %s", action, systemPrompt, userPrompt)
	// Call LLM
	refinedContent, err := llm.AskLLM(systemPrompt, userPrompt)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return refinedContent, nil
}

// getRefinePrompts retrieves or creates prompts for content refinement
func (s *PostService) getRefinePrompts(promptName, requirement, content string) (systemPrompt, userPrompt string, err error) {
	// Map frontend actions to backend prompt names and extract requirements

	var extractedRequirement string

	// Use extracted requirement if available, otherwise use provided requirement
	if extractedRequirement != "" {
		requirement = extractedRequirement
	}

	// Try to get the prompt from repository
	promptResponse, err := s.promptService.GetByName(promptName)
	if err == nil {
		// Found prompt in repository, use it
		userPrompt = promptResponse.UserPrompt
		if requirement != "" {
			// Replace placeholders in user prompt
			userPrompt = strings.ReplaceAll(userPrompt, "{{requirement}}", requirement)
			userPrompt = strings.ReplaceAll(userPrompt, "{{tone}}", requirement)
			userPrompt = strings.ReplaceAll(userPrompt, "{{language}}", requirement)
		}
		userPrompt = strings.ReplaceAll(userPrompt, "{{content}}", content)
		return promptResponse.SystemPrompt, userPrompt, nil
	}

	// Prompt not found in repository, return error
	return "", "", fmt.Errorf("prompt '%s' not found in repository", promptName)
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
