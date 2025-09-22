package wiki

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// WikiService handles business logic for wiki pages
type WikiService struct {
	repo *WikiRepository
}

// NewWikiService creates a new wiki service
func NewWikiService(db *gorm.DB) *WikiService {
	return &WikiService{
		repo: NewWikiRepository(db),
	}
}

// Request/Response structures

// CreateWikiPageRequest represents the request to create a wiki page
type CreateWikiPageRequest struct {
	Title       string                `json:"title" binding:"required" validate:"required,min=1,max=200"`
	Slug        string                `json:"slug" validate:"max=200"`
	Content     string                `json:"content" validate:"required,min=1"`
	Summary     string                `json:"summary" validate:"max=500"`
	Status      models.WikiPageStatus `json:"status" validate:"oneof=draft published archived protected"`
	Type        models.WikiPageType   `json:"type" validate:"oneof=article template category redirect stub"`
	IsProtected bool                  `json:"is_protected"`
	Categories  []string              `json:"categories"`
	Tags        []string              `json:"tags"`
	ParentID    string                `json:"parent_id,omitempty"`
	RedirectTo  string                `json:"redirect_to,omitempty"` // for redirect pages
	Language    string                `json:"language" validate:"max=10"`
	ChangeNote  string                `json:"change_note" validate:"max=200"` // for initial revision
}

// UpdateWikiPageRequest represents the request to update a wiki page
type UpdateWikiPageRequest struct {
	Title       string                `json:"title" validate:"min=1,max=200"`
	Slug        string                `json:"slug" validate:"max=200"`
	Content     string                `json:"content" validate:"min=1"`
	Summary     string                `json:"summary" validate:"max=500"`
	Status      models.WikiPageStatus `json:"status" validate:"oneof=draft published archived protected"`
	Type        models.WikiPageType   `json:"type" validate:"oneof=article template category redirect stub"`
	IsProtected bool                  `json:"is_protected"`
	Categories  []string              `json:"categories"`
	Tags        []string              `json:"tags"`
	ParentID    string                `json:"parent_id,omitempty"`
	RedirectTo  string                `json:"redirect_to,omitempty"`
	Language    string                `json:"language" validate:"max=10"`
	ChangeNote  string                `json:"change_note" validate:"max=200"`
}

// CreateRevisionRequest represents the request to create a revision
type CreateRevisionRequest struct {
	Title      string `json:"title" validate:"required,min=1,max=200"`
	Content    string `json:"content" validate:"required,min=1"`
	Summary    string `json:"summary" validate:"max=500"`
	ChangeNote string `json:"change_note" validate:"max=200"`
}

// WikiPageResponse represents the response for a wiki page
type WikiPageResponse struct {
	ID             string                `json:"id"`
	Title          string                `json:"title"`
	Slug           string                `json:"slug"`
	Content        string                `json:"content"`
	Summary        string                `json:"summary"`
	Status         models.WikiPageStatus `json:"status"`
	Type           models.WikiPageType   `json:"type"`
	IsProtected    bool                  `json:"is_protected"`
	IsLocked       bool                  `json:"is_locked"`
	Categories     []string              `json:"categories"`
	Tags           []string              `json:"tags"`
	ParentID       *string               `json:"parent_id"`
	RedirectTo     *string               `json:"redirect_to"`
	ViewCount      int64                 `json:"view_count"`
	EditCount      int64                 `json:"edit_count"`
	CurrentVersion int                   `json:"current_version"`
	Language       string                `json:"language"`
	TranslationOf  *string               `json:"translation_of"`
	CreatedBy      string                `json:"created_by"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedBy      string                `json:"updated_by"`
	UpdatedAt      time.Time             `json:"updated_at"`

	// Additional metadata for responses
	Authenticated bool   `json:"authenticated"`
	UserID        string `json:"user_id,omitempty"`
	CanEdit       bool   `json:"can_edit,omitempty"`
	CanDelete     bool   `json:"can_delete,omitempty"`
}

// WikiRevisionResponse represents the response for a wiki revision
type WikiRevisionResponse struct {
	ID          string    `json:"id"`
	PageID      string    `json:"page_id"`
	Version     int       `json:"version"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	ChangeNote  string    `json:"change_note"`
	ContentSize int       `json:"content_size"`
	LineCount   int       `json:"line_count"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// WikiPageListResponse represents the response for listing wiki pages
type WikiPageListResponse struct {
	Pages []WikiPageResponse `json:"pages"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

// WikiRevisionListResponse represents the response for listing revisions
type WikiRevisionListResponse struct {
	Revisions []WikiRevisionResponse `json:"revisions"`
	Total     int64                  `json:"total"`
	Page      int                    `json:"page"`
	Limit     int                    `json:"limit"`
}

// WikiDiffResponse represents the response for comparing revisions
type WikiDiffResponse struct {
	PageID      string `json:"page_id"`
	FromVersion int    `json:"from_version"`
	ToVersion   int    `json:"to_version"`
	Diff        string `json:"diff"`    // simplified diff representation
	Changes     int    `json:"changes"` // number of changes
}

// Core CRUD operations

// CreatePage creates a new wiki page
func (s *WikiService) CreatePage(req *CreateWikiPageRequest, realmID, createdBy string) (*WikiPageResponse, error) {
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

	// Create wiki page model
	page := &models.WikiPage{
		ID:             uuid.New().String(),
		RealmID:        realmID,
		Title:          req.Title,
		Slug:           uniqueSlug,
		Content:        req.Content,
		Summary:        req.Summary,
		Status:         req.Status,
		Type:           req.Type,
		Protected:      req.IsProtected,
		CurrentVersion: 1,
		Language:       req.Language,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
		UpdatedBy:      createdBy,
		UpdatedAt:      time.Now(),
	}

	// Set categories and tags
	page.SetCategories(req.Categories)
	page.SetTags(req.Tags)

	// Set parent ID if provided
	if req.ParentID != "" {
		page.ParentID = &req.ParentID
	}

	// Set redirect if provided
	if req.RedirectTo != "" {
		page.RedirectTo = &req.RedirectTo
	}

	// Create the page
	if err := s.repo.Create(page); err != nil {
		return nil, fmt.Errorf("failed to create wiki page: %w", err)
	}

	// Create initial revision
	revision := &models.WikiRevision{
		ID:         uuid.New().String(),
		PageID:     page.ID,
		Version:    1,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		ChangeNote: req.ChangeNote,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	// Calculate content metrics
	revision.ContentSize = len(req.Content)
	revision.LineCount = strings.Count(req.Content, "\n") + 1

	if err := s.repo.CreateRevision(revision); err != nil {
		// If revision creation fails, we should still return the page
		// but log the error
		fmt.Printf("Warning: failed to create initial revision for page %s: %v\n", page.ID, err)
	}

	return s.toPageResponse(page), nil
}

// GetPageBySlug retrieves a wiki page by slug
func (s *WikiService) GetPageBySlug(slug, realmID string) (*WikiPageResponse, error) {
	page, err := s.repo.GetBySlug(slug, realmID)
	if err != nil {
		return nil, err
	}

	// Increment view count
	s.repo.IncrementViewCount(page.ID)

	return s.toPageResponse(page), nil
}

// GetPageByID retrieves a wiki page by ID
func (s *WikiService) GetPageByID(id string) (*WikiPageResponse, error) {
	page, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toPageResponse(page), nil
}

// UpdatePage updates an existing wiki page
func (s *WikiService) UpdatePage(id string, req *UpdateWikiPageRequest, updatedBy string) (*WikiPageResponse, error) {
	// Get existing page
	page, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Check if page can be edited
	if !page.CanBeEdited() {
		return nil, fmt.Errorf("page cannot be edited (locked or archived)")
	}

	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if slug needs to be updated and is unique
	if req.Slug != "" && req.Slug != page.Slug {
		uniqueSlug, err := s.ensureUniqueSlug(req.Slug, page.RealmID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique slug: %w", err)
		}
		page.Slug = uniqueSlug
	}

	// Update fields
	if req.Title != "" {
		page.Title = req.Title
	}
	if req.Content != "" {
		page.Content = req.Content
	}
	if req.Summary != "" {
		page.Summary = req.Summary
	}
	if req.Status != "" {
		page.Status = req.Status
	}
	if req.Type != "" {
		page.Type = req.Type
	}
	page.Protected = req.IsProtected

	// Update categories and tags
	if req.Categories != nil {
		page.SetCategories(req.Categories)
	}
	if req.Tags != nil {
		page.SetTags(req.Tags)
	}

	// Update parent ID
	if req.ParentID != "" {
		page.ParentID = &req.ParentID
	} else {
		page.ParentID = nil
	}

	// Update redirect
	if req.RedirectTo != "" {
		page.RedirectTo = &req.RedirectTo
	} else {
		page.RedirectTo = nil
	}

	// Update language
	if req.Language != "" {
		page.Language = req.Language
	}

	// Increment version and edit count
	page.IncrementEditCount()
	page.UpdatedBy = updatedBy
	page.UpdatedAt = time.Now()

	// Save the page
	if err := s.repo.Update(page); err != nil {
		return nil, fmt.Errorf("failed to update wiki page: %w", err)
	}

	// Create new revision if content changed
	if req.Content != "" && req.Content != page.Content {
		revision := &models.WikiRevision{
			ID:         uuid.New().String(),
			PageID:     page.ID,
			Version:    page.CurrentVersion,
			Title:      page.Title,
			Content:    page.Content,
			Summary:    page.Summary,
			ChangeNote: req.ChangeNote,
			CreatedBy:  updatedBy,
			CreatedAt:  time.Now(),
		}

		revision.ContentSize = len(page.Content)
		revision.LineCount = strings.Count(page.Content, "\n") + 1

		if err := s.repo.CreateRevision(revision); err != nil {
			fmt.Printf("Warning: failed to create revision for page %s: %v\n", page.ID, err)
		}
	}

	return s.toPageResponse(page), nil
}

// DeletePage deletes a wiki page
func (s *WikiService) DeletePage(id string) error {
	page, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("page not found: %w", err)
	}

	if !page.CanBeDeleted() {
		return fmt.Errorf("page cannot be deleted (archived)")
	}

	return s.repo.Delete(id)
}

// Wiki-specific operations

// CreateRevision creates a new revision for a page
func (s *WikiService) CreateRevision(pageID string, req *CreateRevisionRequest, createdBy string) (*WikiRevisionResponse, error) {
	// Get the page
	page, err := s.repo.GetByID(pageID)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Check if page can be edited
	if !page.CanBeEdited() {
		return nil, fmt.Errorf("page cannot be edited (locked or archived)")
	}

	// Validate request
	if err := s.validateRevisionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create revision
	revision := &models.WikiRevision{
		ID:         uuid.New().String(),
		PageID:     pageID,
		Version:    page.CurrentVersion + 1,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		ChangeNote: req.ChangeNote,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	revision.ContentSize = len(req.Content)
	revision.LineCount = strings.Count(req.Content, "\n") + 1

	// Create the revision
	if err := s.repo.CreateRevision(revision); err != nil {
		return nil, fmt.Errorf("failed to create revision: %w", err)
	}

	// Update the page
	page.Title = req.Title
	page.Content = req.Content
	page.Summary = req.Summary
	page.IncrementEditCount()
	page.UpdatedBy = createdBy
	page.UpdatedAt = time.Now()

	if err := s.repo.Update(page); err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	return s.toRevisionResponse(revision), nil
}

// GetPageHistory retrieves the revision history for a page
func (s *WikiService) GetPageHistory(pageID string, page, limit int) (*WikiRevisionListResponse, error) {
	revisions, total, err := s.repo.GetRevisions(pageID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get page history: %w", err)
	}

	responseRevisions := make([]WikiRevisionResponse, len(revisions))
	for i, revision := range revisions {
		responseRevisions[i] = *s.toRevisionResponse(&revision)
	}

	return &WikiRevisionListResponse{
		Revisions: responseRevisions,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// RevertToRevision reverts a page to a specific revision
func (s *WikiService) RevertToRevision(pageID string, revisionID string, revertedBy string) (*WikiPageResponse, error) {
	// Get the page
	page, err := s.repo.GetByID(pageID)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	// Check if page can be edited
	if !page.CanBeEdited() {
		return nil, fmt.Errorf("page cannot be edited (locked or archived)")
	}

	// Get the revision
	revision, err := s.repo.GetRevisionByVersion(pageID, 0) // We'll need to implement this method
	if err != nil {
		return nil, fmt.Errorf("revision not found: %w", err)
	}

	// Create a new revision with the reverted content
	revertRevision := &models.WikiRevision{
		ID:         uuid.New().String(),
		PageID:     pageID,
		Version:    page.CurrentVersion + 1,
		Title:      revision.Title,
		Content:    revision.Content,
		Summary:    revision.Summary,
		ChangeNote: fmt.Sprintf("Reverted to revision %d", revision.Version),
		CreatedBy:  revertedBy,
		CreatedAt:  time.Now(),
	}

	revertRevision.ContentSize = len(revision.Content)
	revertRevision.LineCount = strings.Count(revision.Content, "\n") + 1

	// Create the revert revision
	if err := s.repo.CreateRevision(revertRevision); err != nil {
		return nil, fmt.Errorf("failed to create revert revision: %w", err)
	}

	// Update the page
	page.Title = revision.Title
	page.Content = revision.Content
	page.Summary = revision.Summary
	page.IncrementEditCount()
	page.UpdatedBy = revertedBy
	page.UpdatedAt = time.Now()

	if err := s.repo.Update(page); err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	return s.toPageResponse(page), nil
}

// CompareRevisions compares two revisions of a page
func (s *WikiService) CompareRevisions(pageID string, fromVersion, toVersion int) (*WikiDiffResponse, error) {
	// Get both revisions
	fromRevision, err := s.repo.GetRevisionByVersion(pageID, fromVersion)
	if err != nil {
		return nil, fmt.Errorf("from revision not found: %w", err)
	}

	toRevision, err := s.repo.GetRevisionByVersion(pageID, toVersion)
	if err != nil {
		return nil, fmt.Errorf("to revision not found: %w", err)
	}

	// Simple diff implementation (in a real system, you'd use a proper diff library)
	diff := s.generateSimpleDiff(fromRevision.Content, toRevision.Content)
	changes := s.countChanges(fromRevision.Content, toRevision.Content)

	return &WikiDiffResponse{
		PageID:      pageID,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Diff:        diff,
		Changes:     changes,
	}, nil
}

// Search and discovery

// SearchPages searches for wiki pages
func (s *WikiService) SearchPages(realmID, query string, page, limit int) (*WikiPageListResponse, error) {
	pages, total, err := s.repo.Search(realmID, query, models.WikiPageStatusPublished, "", page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search pages: %w", err)
	}

	responsePages := make([]WikiPageResponse, len(pages))
	for i, page := range pages {
		responsePages[i] = *s.toPageResponse(&page)
	}

	return &WikiPageListResponse{
		Pages: responsePages,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetPagesByCategory retrieves pages by category
func (s *WikiService) GetPagesByCategory(realmID, category string, page, limit int) (*WikiPageListResponse, error) {
	pages, total, err := s.repo.GetByCategory(realmID, category, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pages by category: %w", err)
	}

	responsePages := make([]WikiPageResponse, len(pages))
	for i, page := range pages {
		responsePages[i] = *s.toPageResponse(&page)
	}

	return &WikiPageListResponse{
		Pages: responsePages,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// GetRecentChanges retrieves recent changes across all pages
func (s *WikiService) GetRecentChanges(realmID string, page, limit int) (*WikiRevisionListResponse, error) {
	revisions, err := s.repo.GetRecentRevisions(realmID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent changes: %w", err)
	}

	responseRevisions := make([]WikiRevisionResponse, len(revisions))
	for i, revision := range revisions {
		responseRevisions[i] = *s.toRevisionResponse(&revision)
	}

	return &WikiRevisionListResponse{
		Revisions: responseRevisions,
		Total:     int64(len(responseRevisions)),
		Page:      page,
		Limit:     limit,
	}, nil
}

// GetRandomPage retrieves a random published page
func (s *WikiService) GetRandomPage(realmID string) (*WikiPageResponse, error) {
	page, err := s.repo.GetRandomPage(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get random page: %w", err)
	}

	return s.toPageResponse(page), nil
}

// GetOrphanedPages retrieves pages with no incoming links
func (s *WikiService) GetOrphanedPages(realmID string) (*WikiPageListResponse, error) {
	pages, err := s.repo.GetOrphanedPages(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orphaned pages: %w", err)
	}

	responsePages := make([]WikiPageResponse, len(pages))
	for i, page := range pages {
		responsePages[i] = *s.toPageResponse(&page)
	}

	return &WikiPageListResponse{
		Pages: responsePages,
		Total: int64(len(responsePages)),
		Page:  1,
		Limit: len(responsePages),
	}, nil
}

// GetWantedPages retrieves pages that are linked but don't exist
func (s *WikiService) GetWantedPages(realmID string) ([]string, error) {
	return s.repo.GetWantedPages(realmID)
}

// GetDeadEndPages retrieves pages with no outgoing links
func (s *WikiService) GetDeadEndPages(realmID string) (*WikiPageListResponse, error) {
	pages, err := s.repo.GetDeadEndPages(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dead end pages: %w", err)
	}

	responsePages := make([]WikiPageResponse, len(pages))
	for i, page := range pages {
		responsePages[i] = *s.toPageResponse(&page)
	}

	return &WikiPageListResponse{
		Pages: responsePages,
		Total: int64(len(responsePages)),
		Page:  1,
		Limit: len(responsePages),
	}, nil
}

// ListPages lists wiki pages with filtering
func (s *WikiService) ListPages(realmID string, status models.WikiPageStatus, pageType models.WikiPageType, page, limit int) (*WikiPageListResponse, error) {
	pages, total, err := s.repo.List(realmID, status, pageType, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}

	responsePages := make([]WikiPageResponse, len(pages))
	for i, page := range pages {
		responsePages[i] = *s.toPageResponse(&page)
	}

	return &WikiPageListResponse{
		Pages: responsePages,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// Helper methods

// validateCreateRequest validates a create request
func (s *WikiService) validateCreateRequest(req *CreateWikiPageRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	if req.Type == models.WikiPageTypeRedirect && req.RedirectTo == "" {
		return fmt.Errorf("redirect_to is required for redirect pages")
	}
	return nil
}

// validateUpdateRequest validates an update request
func (s *WikiService) validateUpdateRequest(req *UpdateWikiPageRequest) error {
	if req.Type == models.WikiPageTypeRedirect && req.RedirectTo == "" {
		return fmt.Errorf("redirect_to is required for redirect pages")
	}
	return nil
}

// validateRevisionRequest validates a revision request
func (s *WikiService) validateRevisionRequest(req *CreateRevisionRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}

// generateSlugFromTitle generates a URL-friendly slug from a title
func (s *WikiService) generateSlugFromTitle(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 200 {
		slug = slug[:200]
		slug = strings.Trim(slug, "-")
	}

	return slug
}

// ensureUniqueSlug ensures a slug is unique within a realm
func (s *WikiService) ensureUniqueSlug(slug, realmID, excludeID string) (string, error) {
	baseSlug := slug
	counter := 1

	for {
		exists, err := s.repo.CheckSlugExists(slug, realmID, excludeID)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		// Prevent infinite loop
		if counter > 1000 {
			return "", fmt.Errorf("unable to generate unique slug")
		}
	}
}

// toPageResponse converts a WikiPage model to a WikiPageResponse
func (s *WikiService) toPageResponse(page *models.WikiPage) *WikiPageResponse {
	return &WikiPageResponse{
		ID:             page.ID,
		Title:          page.Title,
		Slug:           page.Slug,
		Content:        page.Content,
		Summary:        page.Summary,
		Status:         page.Status,
		Type:           page.Type,
		IsProtected:    page.IsProtected(),
		IsLocked:       page.IsLocked(),
		Categories:     page.GetCategories(),
		Tags:           page.GetTags(),
		ParentID:       page.ParentID,
		RedirectTo:     page.RedirectTo,
		ViewCount:      page.ViewCount,
		EditCount:      page.EditCount,
		CurrentVersion: page.CurrentVersion,
		Language:       page.Language,
		TranslationOf:  page.TranslationOf,
		CreatedBy:      page.CreatedBy,
		CreatedAt:      page.CreatedAt,
		UpdatedBy:      page.UpdatedBy,
		UpdatedAt:      page.UpdatedAt,
	}
}

// toRevisionResponse converts a WikiRevision model to a WikiRevisionResponse
func (s *WikiService) toRevisionResponse(revision *models.WikiRevision) *WikiRevisionResponse {
	return &WikiRevisionResponse{
		ID:          revision.ID,
		PageID:      revision.PageID,
		Version:     revision.Version,
		Title:       revision.Title,
		Content:     revision.Content,
		Summary:     revision.Summary,
		ChangeNote:  revision.ChangeNote,
		ContentSize: revision.ContentSize,
		LineCount:   revision.LineCount,
		CreatedBy:   revision.CreatedBy,
		CreatedAt:   revision.CreatedAt,
	}
}

// generateSimpleDiff generates a simple diff between two content strings
func (s *WikiService) generateSimpleDiff(from, to string) string {
	// This is a very basic diff implementation
	// In a real system, you'd use a proper diff library like go-diff
	if from == to {
		return "No changes"
	}

	fromLines := strings.Split(from, "\n")
	toLines := strings.Split(to, "\n")

	var diff strings.Builder
	diff.WriteString("Content changed:\n")

	// Simple line-by-line comparison
	maxLines := len(fromLines)
	if len(toLines) > maxLines {
		maxLines = len(toLines)
	}

	for i := 0; i < maxLines; i++ {
		var fromLine, toLine string
		if i < len(fromLines) {
			fromLine = fromLines[i]
		}
		if i < len(toLines) {
			toLine = toLines[i]
		}

		if fromLine != toLine {
			diff.WriteString(fmt.Sprintf("Line %d: %s -> %s\n", i+1, fromLine, toLine))
		}
	}

	return diff.String()
}

// countChanges counts the number of changes between two content strings
func (s *WikiService) countChanges(from, to string) int {
	if from == to {
		return 0
	}

	fromLines := strings.Split(from, "\n")
	toLines := strings.Split(to, "\n")

	changes := 0
	maxLines := len(fromLines)
	if len(toLines) > maxLines {
		maxLines = len(toLines)
	}

	for i := 0; i < maxLines; i++ {
		var fromLine, toLine string
		if i < len(fromLines) {
			fromLine = fromLines[i]
		}
		if i < len(toLines) {
			toLine = toLines[i]
		}

		if fromLine != toLine {
			changes++
		}
	}

	return changes
}
