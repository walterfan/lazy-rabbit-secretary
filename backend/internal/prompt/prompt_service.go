package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// PromptService handles business logic for prompts
type PromptService struct {
	repo *PromptRepository
}

// NewPromptService creates a new prompt service
func NewPromptService(db *gorm.DB) *PromptService {
	return &PromptService{
		repo: NewPromptRepository(db),
	}
}

// CreatePromptRequest represents the request to create a prompt
type CreatePromptRequest struct {
	Name         string `json:"name" binding:"required" validate:"required,min=1,max=100"`
	Description  string `json:"description" validate:"max=500"`
	SystemPrompt string `json:"system_prompt" binding:"required" validate:"required,min=1"`
	UserPrompt   string `json:"user_prompt" binding:"required" validate:"required,min=1"`
	Tags         string `json:"tags" validate:"max=200"`
}

// UpdatePromptRequest represents the request to update a prompt
type UpdatePromptRequest struct {
	Name         string `json:"name" validate:"min=1,max=100"`
	Description  string `json:"description" validate:"max=500"`
	SystemPrompt string `json:"system_prompt" validate:"min=1"`
	UserPrompt   string `json:"user_prompt" validate:"min=1"`
	Tags         string `json:"tags" validate:"max=200"`
}

// PromptResponse represents the response for prompt operations
type PromptResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	SystemPrompt string    `json:"system_prompt"`
	UserPrompt   string    `json:"user_prompt"`
	Tags         string    `json:"tags"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedBy    string    `json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PromptListResponse represents the response for listing prompts
type PromptListResponse struct {
	Prompts []PromptResponse `json:"prompts"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
}

// CreateFromInput creates a new prompt from input data
func (s *PromptService) CreateFromInput(req *CreatePromptRequest, realmID, createdBy string) (*PromptResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if prompt name already exists
	exists, err := s.repo.Exists(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check prompt existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("prompt with name '%s' already exists", req.Name)
	}

	// Create prompt model
	prompt := &models.Prompt{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Description:  req.Description,
		SystemPrompt: req.SystemPrompt,
		UserPrompt:   req.UserPrompt,
		Tags:         req.Tags,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
		UpdatedBy:    createdBy,
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := s.repo.Create(prompt); err != nil {
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}

	return s.toResponse(prompt), nil
}

// UpdateFromInput updates a prompt from input data
func (s *PromptService) UpdateFromInput(id string, req *UpdatePromptRequest, updatedBy string) (*PromptResponse, error) {
	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing prompt
	prompt, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("prompt not found")
		}
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	// Check if name is being changed and already exists
	if req.Name != "" && req.Name != prompt.Name {
		exists, err := s.repo.Exists(req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check prompt existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("prompt with name '%s' already exists", req.Name)
		}
		prompt.Name = req.Name
	}

	// Update fields if provided
	if req.Description != "" {
		prompt.Description = req.Description
	}
	if req.SystemPrompt != "" {
		prompt.SystemPrompt = req.SystemPrompt
	}
	if req.UserPrompt != "" {
		prompt.UserPrompt = req.UserPrompt
	}
	if req.Tags != "" {
		prompt.Tags = req.Tags
	}

	prompt.UpdatedBy = updatedBy
	prompt.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(prompt); err != nil {
		return nil, fmt.Errorf("failed to update prompt: %w", err)
	}

	return s.toResponse(prompt), nil
}

// GetByID retrieves a prompt by ID
func (s *PromptService) GetByID(id string) (*PromptResponse, error) {
	prompt, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("prompt not found")
		}
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	return s.toResponse(prompt), nil
}

// GetByName retrieves a prompt by name
func (s *PromptService) GetByName(name string) (*PromptResponse, error) {
	prompt, err := s.repo.GetByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("prompt not found")
		}
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	return s.toResponse(prompt), nil
}

// List retrieves prompts with pagination
func (s *PromptService) List(page, limit int) (*PromptListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	prompts, total, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	return &PromptListResponse{
		Prompts: s.toResponseList(prompts),
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// Search searches prompts by query
func (s *PromptService) Search(query string, page, limit int) (*PromptListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	prompts, total, err := s.repo.Search(query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search prompts: %w", err)
	}

	return &PromptListResponse{
		Prompts: s.toResponseList(prompts),
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// GetByTags retrieves prompts by tags
func (s *PromptService) GetByTags(tags []string, page, limit int) (*PromptListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	prompts, total, err := s.repo.GetByTags(tags, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompts by tags: %w", err)
	}

	return &PromptListResponse{
		Prompts: s.toResponseList(prompts),
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// Delete deletes a prompt
func (s *PromptService) Delete(id string) error {
	// Check if prompt exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("prompt not found")
		}
		return fmt.Errorf("failed to get prompt: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete prompt: %w", err)
	}

	return nil
}

// validateCreateRequest validates create prompt request
func (s *PromptService) validateCreateRequest(req *CreatePromptRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.SystemPrompt) == "" {
		return fmt.Errorf("system prompt is required")
	}
	if strings.TrimSpace(req.UserPrompt) == "" {
		return fmt.Errorf("user prompt is required")
	}
	if len(req.Name) > 100 {
		return fmt.Errorf("name must be less than 100 characters")
	}
	if len(req.Description) > 500 {
		return fmt.Errorf("description must be less than 500 characters")
	}
	if len(req.Tags) > 200 {
		return fmt.Errorf("tags must be less than 200 characters")
	}
	return nil
}

// validateUpdateRequest validates update prompt request
func (s *PromptService) validateUpdateRequest(req *UpdatePromptRequest) error {
	if req.Name != "" && len(req.Name) > 100 {
		return fmt.Errorf("name must be less than 100 characters")
	}
	if req.Description != "" && len(req.Description) > 500 {
		return fmt.Errorf("description must be less than 500 characters")
	}
	if req.Tags != "" && len(req.Tags) > 200 {
		return fmt.Errorf("tags must be less than 200 characters")
	}
	return nil
}

// toResponse converts a prompt model to response
func (s *PromptService) toResponse(prompt *models.Prompt) *PromptResponse {
	return &PromptResponse{
		ID:           prompt.ID,
		Name:         prompt.Name,
		Description:  prompt.Description,
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   prompt.UserPrompt,
		Tags:         prompt.Tags,
		CreatedBy:    prompt.CreatedBy,
		CreatedAt:    prompt.CreatedAt,
		UpdatedBy:    prompt.UpdatedBy,
		UpdatedAt:    prompt.UpdatedAt,
	}
}

// toResponseList converts a list of prompt models to responses
func (s *PromptService) toResponseList(prompts []*models.Prompt) []PromptResponse {
	responses := make([]PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = *s.toResponse(prompt)
	}
	return responses
}
