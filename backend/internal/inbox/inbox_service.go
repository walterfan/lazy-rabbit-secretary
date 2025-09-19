package inbox

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// InboxService contains business logic for inbox items
type InboxService struct {
	repo *InboxRepository
}

// NewInboxService creates a new inbox service
func NewInboxService(db *gorm.DB) *InboxService {
	return &InboxService{
		repo: NewInboxRepository(db),
	}
}

// CreateInboxItemRequest represents the request to create an inbox item
type CreateInboxItemRequest struct {
	Title       string `json:"title" binding:"required" validate:"required,min=1,max=200"`
	Description string `json:"description" validate:"max=1000"`
	Priority    string `json:"priority" validate:"oneof=low normal high urgent"`
	Tags        string `json:"tags" validate:"max=200"`
	Context     string `json:"context" validate:"max=100"`
}

// UpdateInboxItemRequest represents the request to update an inbox item
type UpdateInboxItemRequest struct {
	Title       string `json:"title" validate:"min=1,max=200"`
	Description string `json:"description" validate:"max=1000"`
	Priority    string `json:"priority" validate:"oneof=low normal high urgent"`
	Status      string `json:"status" validate:"oneof=pending processing completed archived"`
	Tags        string `json:"tags" validate:"max=200"`
	Context     string `json:"context" validate:"max=100"`
}

// InboxItemResponse represents the response for inbox item operations
type InboxItemResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	Tags        string    `json:"tags"`
	Context     string    `json:"context"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InboxListResponse represents the response for listing inbox items
type InboxListResponse struct {
	Items []InboxItemResponse `json:"items"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

// CreateFromInput creates a new inbox item from input data
func (s *InboxService) CreateFromInput(req *CreateInboxItemRequest, realmID, createdBy string) (*InboxItemResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	if req.Priority == "" {
		req.Priority = "normal"
	}

	// Create inbox item model
	item := &models.InboxItem{
		ID:          uuid.New().String(),
		RealmID:     realmID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      "pending",
		Tags:        req.Tags,
		Context:     req.Context,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedBy:   createdBy,
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create inbox item: %w", err)
	}

	return s.toResponse(item), nil
}

// GetByID retrieves an inbox item by ID
func (s *InboxService) GetByID(id string) (*InboxItemResponse, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("inbox item not found")
		}
		return nil, fmt.Errorf("failed to get inbox item: %w", err)
	}

	return s.toResponse(item), nil
}

// UpdateFromInput updates an inbox item from input data
func (s *InboxService) UpdateFromInput(id string, req *UpdateInboxItemRequest, updatedBy string) (*InboxItemResponse, error) {
	// Get existing item
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("inbox item not found")
		}
		return nil, fmt.Errorf("failed to get inbox item: %w", err)
	}

	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update fields
	if req.Title != "" {
		item.Title = req.Title
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Priority != "" {
		item.Priority = req.Priority
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.Tags != "" {
		item.Tags = req.Tags
	}
	if req.Context != "" {
		item.Context = req.Context
	}

	item.UpdatedBy = updatedBy
	item.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update inbox item: %w", err)
	}

	return s.toResponse(item), nil
}

// Delete deletes an inbox item
func (s *InboxService) Delete(id string) error {
	// Check if item exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("inbox item not found")
		}
		return fmt.Errorf("failed to get inbox item: %w", err)
	}

	// Delete item
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete inbox item: %w", err)
	}

	return nil
}

// List retrieves inbox items with pagination and filtering
func (s *InboxService) List(realmID string, params ListParams) (*InboxListResponse, error) {
	items, total, err := s.repo.List(realmID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list inbox items: %w", err)
	}

	// Convert to response format
	responseItems := make([]InboxItemResponse, len(items))
	for i, item := range items {
		responseItems[i] = *s.toResponse(&item)
	}

	return &InboxListResponse{
		Items: responseItems,
		Total: total,
		Page:  params.Page,
		Limit: params.PageSize,
	}, nil
}

// GetPendingItems retrieves all pending inbox items
func (s *InboxService) GetPendingItems(realmID string) ([]InboxItemResponse, error) {
	items, err := s.repo.GetPendingItems(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending items: %w", err)
	}

	responseItems := make([]InboxItemResponse, len(items))
	for i, item := range items {
		responseItems[i] = *s.toResponse(&item)
	}

	return responseItems, nil
}

// GetUrgentItems retrieves all urgent inbox items
func (s *InboxService) GetUrgentItems(realmID string) ([]InboxItemResponse, error) {
	items, err := s.repo.GetUrgentItems(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get urgent items: %w", err)
	}

	responseItems := make([]InboxItemResponse, len(items))
	for i, item := range items {
		responseItems[i] = *s.toResponse(&item)
	}

	return responseItems, nil
}

// UpdateStatus updates the status of an inbox item
func (s *InboxService) UpdateStatus(id, status string) (*InboxItemResponse, error) {
	// Validate status
	validStatuses := []string{"pending", "processing", "completed", "archived"}
	if !contains(validStatuses, status) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	// Update status
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Get updated item
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated item: %w", err)
	}

	return s.toResponse(item), nil
}

// BulkUpdateStatus updates status for multiple items
func (s *InboxService) BulkUpdateStatus(ids []string, status string) error {
	// Validate status
	validStatuses := []string{"pending", "processing", "completed", "archived"}
	if !contains(validStatuses, status) {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Validate IDs
	if len(ids) == 0 {
		return fmt.Errorf("no items to update")
	}

	// Update status
	if err := s.repo.BulkUpdateStatus(ids, status); err != nil {
		return fmt.Errorf("failed to bulk update status: %w", err)
	}

	return nil
}

// GetStats retrieves inbox statistics
func (s *InboxService) GetStats(realmID string) (map[string]int64, error) {
	stats, err := s.repo.GetStats(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}

// Helper methods

func (s *InboxService) validateCreateRequest(req *CreateInboxItemRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) > 200 {
		return fmt.Errorf("title too long (max 200 characters)")
	}
	if len(req.Description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}
	if req.Priority != "" && !contains([]string{"low", "normal", "high", "urgent"}, req.Priority) {
		return fmt.Errorf("invalid priority")
	}
	if len(req.Tags) > 200 {
		return fmt.Errorf("tags too long (max 200 characters)")
	}
	if len(req.Context) > 100 {
		return fmt.Errorf("context too long (max 100 characters)")
	}
	return nil
}

func (s *InboxService) validateUpdateRequest(req *UpdateInboxItemRequest) error {
	if req.Title != "" && len(req.Title) > 200 {
		return fmt.Errorf("title too long (max 200 characters)")
	}
	if req.Description != "" && len(req.Description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}
	if req.Priority != "" && !contains([]string{"low", "normal", "high", "urgent"}, req.Priority) {
		return fmt.Errorf("invalid priority")
	}
	if req.Status != "" && !contains([]string{"pending", "processing", "completed", "archived"}, req.Status) {
		return fmt.Errorf("invalid status")
	}
	if req.Tags != "" && len(req.Tags) > 200 {
		return fmt.Errorf("tags too long (max 200 characters)")
	}
	if req.Context != "" && len(req.Context) > 100 {
		return fmt.Errorf("context too long (max 100 characters)")
	}
	return nil
}

func (s *InboxService) toResponse(item *models.InboxItem) *InboxItemResponse {
	return &InboxItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		Priority:    item.Priority,
		Status:      item.Status,
		Tags:        item.Tags,
		Context:     item.Context,
		CreatedBy:   item.CreatedBy,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
