package daily

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// DailyService contains business logic for daily checklist items
type DailyService struct {
	repo *DailyRepository
}

// NewDailyService creates a new daily service
func NewDailyService(db *gorm.DB) *DailyService {
	return &DailyService{
		repo: NewDailyRepository(db),
	}
}

// CreateDailyItemRequest represents the request to create a daily checklist item
type CreateDailyItemRequest struct {
	Title         string     `json:"title" binding:"required" validate:"required,min=1,max=200"`
	Description   string     `json:"description" validate:"max=1000"`
	Priority      string     `json:"priority" validate:"oneof=A B+ B C D"`
	EstimatedTime int        `json:"estimated_time" validate:"min=0,max=1440"` // max 24 hours in minutes
	Deadline      *time.Time `json:"deadline"`
	Context       string     `json:"context" validate:"max=100"`
	Notes         string     `json:"notes" validate:"max=500"`
	InboxItemID   *string    `json:"inbox_item_id"`
	Date          time.Time  `json:"date" binding:"required"`
}

// UpdateDailyItemRequest represents the request to update a daily checklist item
type UpdateDailyItemRequest struct {
	Title         string     `json:"title" validate:"min=1,max=200"`
	Description   string     `json:"description" validate:"max=1000"`
	Priority      string     `json:"priority" validate:"oneof=A B+ B C D"`
	EstimatedTime int        `json:"estimated_time" validate:"min=0,max=1440"`
	Deadline      *time.Time `json:"deadline"`
	Context       string     `json:"context" validate:"max=100"`
	Status        string     `json:"status" validate:"oneof=pending in_progress completed cancelled"`
	ActualTime    int        `json:"actual_time" validate:"min=0,max=1440"`
	Notes         string     `json:"notes" validate:"max=500"`
}

// DailyItemResponse represents the response for daily checklist item operations
type DailyItemResponse struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Priority        string     `json:"priority"`
	EstimatedTime   int        `json:"estimated_time"`
	Deadline        *time.Time `json:"deadline"`
	Context         string     `json:"context"`
	Status          string     `json:"status"`
	CompletionTime  *time.Time `json:"completion_time"`
	ActualTime      int        `json:"actual_time"`
	Notes           string     `json:"notes"`
	InboxItemID     *string    `json:"inbox_item_id"`
	Date            time.Time  `json:"date"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// DailyListResponse represents the response for listing daily checklist items
type DailyListResponse struct {
	Items []DailyItemResponse `json:"items"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

// DailyStatsResponse represents daily statistics
type DailyStatsResponse struct {
	Date                time.Time `json:"date"`
	TotalItems          int64     `json:"total_items"`
	CompletedItems      int64     `json:"completed_items"`
	PendingItems        int64     `json:"pending_items"`
	InProgressItems     int64     `json:"in_progress_items"`
	CancelledItems      int64     `json:"cancelled_items"`
	CompletionRate      float64   `json:"completion_rate"`
	TotalEstimatedTime  int64     `json:"total_estimated_time"`
	TotalActualTime     int64     `json:"total_actual_time"`
	PriorityBreakdown   map[string]int64 `json:"priority_breakdown"`
}

// CreateFromInput creates a new daily checklist item from input data
func (s *DailyService) CreateFromInput(req *CreateDailyItemRequest, realmID, createdBy string) (*DailyItemResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	if req.Priority == "" {
		req.Priority = "B"
	}
	if req.EstimatedTime == 0 {
		req.EstimatedTime = 30 // default 30 minutes
	}

	// Create daily checklist item model
	item := &models.DailyChecklistItem{
		ID:            uuid.New().String(),
		RealmID:       realmID,
		Title:         req.Title,
		Description:   req.Description,
		Priority:      req.Priority,
		EstimatedTime: req.EstimatedTime,
		Deadline:      req.Deadline,
		Context:       req.Context,
		Status:        "pending",
		ActualTime:    0,
		Notes:         req.Notes,
		InboxItemID:   req.InboxItemID,
		Date:          req.Date,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now(),
		UpdatedBy:     createdBy,
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create daily checklist item: %w", err)
	}

	return s.toResponse(item), nil
}

// GetByID retrieves a daily checklist item by ID
func (s *DailyService) GetByID(id string) (*DailyItemResponse, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("daily checklist item not found")
		}
		return nil, fmt.Errorf("failed to get daily checklist item: %w", err)
	}

	return s.toResponse(item), nil
}

// UpdateFromInput updates a daily checklist item from input data
func (s *DailyService) UpdateFromInput(id string, req *UpdateDailyItemRequest, updatedBy string) (*DailyItemResponse, error) {
	// Get existing item
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("daily checklist item not found")
		}
		return nil, fmt.Errorf("failed to get daily checklist item: %w", err)
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
	if req.EstimatedTime > 0 {
		item.EstimatedTime = req.EstimatedTime
	}
	if req.Deadline != nil {
		item.Deadline = req.Deadline
	}
	if req.Context != "" {
		item.Context = req.Context
	}
	if req.Status != "" {
		item.Status = req.Status
		if req.Status == "completed" && item.CompletionTime == nil {
			now := time.Now()
			item.CompletionTime = &now
		}
	}
	if req.ActualTime > 0 {
		item.ActualTime = req.ActualTime
	}
	if req.Notes != "" {
		item.Notes = req.Notes
	}

	item.UpdatedBy = updatedBy
	item.UpdatedAt = time.Now()

	// Save to database
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update daily checklist item: %w", err)
	}

	return s.toResponse(item), nil
}

// Delete deletes a daily checklist item
func (s *DailyService) Delete(id string) error {
	// Check if item exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("daily checklist item not found")
		}
		return fmt.Errorf("failed to get daily checklist item: %w", err)
	}

	// Delete item
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete daily checklist item: %w", err)
	}

	return nil
}

// List retrieves daily checklist items with pagination and filtering
func (s *DailyService) List(realmID string, params ListParams) (*DailyListResponse, error) {
	items, total, err := s.repo.List(realmID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list daily checklist items: %w", err)
	}

	// Convert to response format
	responseItems := make([]DailyItemResponse, len(items))
	for i, item := range items {
		responseItems[i] = *s.toResponse(&item)
	}

	return &DailyListResponse{
		Items: responseItems,
		Total: total,
		Page:  params.Page,
		Limit: params.PageSize,
	}, nil
}

// GetByDate retrieves daily checklist items for a specific date
func (s *DailyService) GetByDate(realmID string, date time.Time) ([]DailyItemResponse, error) {
	items, err := s.repo.GetByDate(realmID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get items by date: %w", err)
	}

	responseItems := make([]DailyItemResponse, len(items))
	for i, item := range items {
		responseItems[i] = *s.toResponse(&item)
	}

	return responseItems, nil
}

// UpdateStatus updates the status of a daily checklist item
func (s *DailyService) UpdateStatus(id, status string) (*DailyItemResponse, error) {
	// Validate status
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
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

// UpdateActualTime updates the actual time spent on a task
func (s *DailyService) UpdateActualTime(id string, actualTime int) (*DailyItemResponse, error) {
	// Validate actual time
	if actualTime < 0 || actualTime > 1440 {
		return nil, fmt.Errorf("invalid actual time: must be between 0 and 1440 minutes")
	}

	// Update actual time
	if err := s.repo.UpdateActualTime(id, actualTime); err != nil {
		return nil, fmt.Errorf("failed to update actual time: %w", err)
	}

	// Get updated item
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated item: %w", err)
	}

	return s.toResponse(item), nil
}

// BulkUpdateStatus updates status for multiple items
func (s *DailyService) BulkUpdateStatus(ids []string, status string) error {
	// Validate status
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
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

// GetStats retrieves daily statistics
func (s *DailyService) GetStats(realmID string, date time.Time) (*DailyStatsResponse, error) {
	stats, err := s.repo.GetStats(realmID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	completionRate, err := s.repo.GetCompletionRate(realmID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get completion rate: %w", err)
	}

	// Build priority breakdown
	priorityBreakdown := make(map[string]int64)
	priorities := []string{"A", "B+", "B", "C", "D"}
	for _, priority := range priorities {
		if count, exists := stats["priority_"+priority]; exists {
			priorityBreakdown[priority] = count
		} else {
			priorityBreakdown[priority] = 0
		}
	}

	return &DailyStatsResponse{
		Date:                date,
		TotalItems:          stats["total"],
		CompletedItems:      stats["status_completed"],
		PendingItems:        stats["status_pending"],
		InProgressItems:     stats["status_in_progress"],
		CancelledItems:      stats["status_cancelled"],
		CompletionRate:      completionRate,
		TotalEstimatedTime:  stats["total_estimated_time"],
		TotalActualTime:     stats["total_actual_time"],
		PriorityBreakdown:   priorityBreakdown,
	}, nil
}

// Helper methods

func (s *DailyService) validateCreateRequest(req *CreateDailyItemRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) > 200 {
		return fmt.Errorf("title too long (max 200 characters)")
	}
	if len(req.Description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}
	if req.Priority != "" && !contains([]string{"A", "B+", "B", "C", "D"}, req.Priority) {
		return fmt.Errorf("invalid priority")
	}
	if req.EstimatedTime < 0 || req.EstimatedTime > 1440 {
		return fmt.Errorf("estimated time must be between 0 and 1440 minutes")
	}
	if len(req.Context) > 100 {
		return fmt.Errorf("context too long (max 100 characters)")
	}
	if len(req.Notes) > 500 {
		return fmt.Errorf("notes too long (max 500 characters)")
	}
	if req.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	return nil
}

func (s *DailyService) validateUpdateRequest(req *UpdateDailyItemRequest) error {
	if req.Title != "" && len(req.Title) > 200 {
		return fmt.Errorf("title too long (max 200 characters)")
	}
	if req.Description != "" && len(req.Description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}
	if req.Priority != "" && !contains([]string{"A", "B+", "B", "C", "D"}, req.Priority) {
		return fmt.Errorf("invalid priority")
	}
	if req.EstimatedTime > 0 && (req.EstimatedTime < 0 || req.EstimatedTime > 1440) {
		return fmt.Errorf("estimated time must be between 0 and 1440 minutes")
	}
	if req.Status != "" && !contains([]string{"pending", "in_progress", "completed", "cancelled"}, req.Status) {
		return fmt.Errorf("invalid status")
	}
	if req.ActualTime > 0 && (req.ActualTime < 0 || req.ActualTime > 1440) {
		return fmt.Errorf("actual time must be between 0 and 1440 minutes")
	}
	if req.Context != "" && len(req.Context) > 100 {
		return fmt.Errorf("context too long (max 100 characters)")
	}
	if req.Notes != "" && len(req.Notes) > 500 {
		return fmt.Errorf("notes too long (max 500 characters)")
	}
	return nil
}

func (s *DailyService) toResponse(item *models.DailyChecklistItem) *DailyItemResponse {
	return &DailyItemResponse{
		ID:              item.ID,
		Title:           item.Title,
		Description:     item.Description,
		Priority:        item.Priority,
		EstimatedTime:   item.EstimatedTime,
		Deadline:        item.Deadline,
		Context:         item.Context,
		Status:          item.Status,
		CompletionTime:  item.CompletionTime,
		ActualTime:      item.ActualTime,
		Notes:           item.Notes,
		InboxItemID:     item.InboxItemID,
		Date:            item.Date,
		CreatedBy:       item.CreatedBy,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
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
