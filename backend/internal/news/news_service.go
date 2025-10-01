package news

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
)

// NewsService handles business logic for news items
type NewsService struct {
	repo *NewsRepository
}

// NewNewsService creates a new news service
func NewNewsService(repo *NewsRepository) *NewsService {
	return &NewsService{
		repo: repo,
	}
}

// Request/Response models
type CreateNewsRequest struct {
	RealmID string `json:"realm_id"`
	Message string `json:"message" binding:"required"`
	Author  string `json:"author"`
}

type UpdateNewsRequest struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

type NewsResponse struct {
	ID        string `json:"id"`
	RealmID   string `json:"realm_id"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NewsListResponse struct {
	News    []NewsResponse `json:"news"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
	HasMore bool           `json:"has_more"`
}

// Create creates a new news item
func (s *NewsService) Create(ctx context.Context, req *CreateNewsRequest) (*NewsResponse, error) {
	// Generate ID
	id := uuid.New().String()
	
	// Use default realm if not provided
	realmID := req.RealmID
	if realmID == "" {
		realmID = "default-realm"
	}
	
	// Create news item
	newsItem := &models.NewsItem{
		ID:        id,
		RealmID:   realmID,
		Message:   req.Message,
		Author:    req.Author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := s.repo.Create(ctx, newsItem)
	if err != nil {
		return nil, err
	}
	
	return s.toResponse(newsItem), nil
}

// GetByID retrieves a news item by ID
func (s *NewsService) GetByID(ctx context.Context, id string) (*NewsResponse, error) {
	newsItem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(newsItem), nil
}

// GetByRealm retrieves news items by realm ID
func (s *NewsService) GetByRealm(ctx context.Context, realmID string, page, limit int) (*NewsListResponse, error) {
	offset := (page - 1) * limit
	newsItems, err := s.repo.GetByRealm(ctx, realmID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx, realmID)
	if err != nil {
		return nil, err
	}

	responses := make([]NewsResponse, len(newsItems))
	for i, item := range newsItems {
		responses[i] = *s.toResponse(item)
	}

	return &NewsListResponse{
		News:    responses,
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(page*limit) < total,
	}, nil
}

// Update updates a news item
func (s *NewsService) Update(ctx context.Context, id string, req *UpdateNewsRequest) (*NewsResponse, error) {
	newsItem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Message != "" {
		newsItem.Message = req.Message
	}
	if req.Author != "" {
		newsItem.Author = req.Author
	}
	newsItem.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, newsItem)
	if err != nil {
		return nil, err
	}

	return s.toResponse(newsItem), nil
}

// Delete deletes a news item
func (s *NewsService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// Helper method to convert model to response
func (s *NewsService) toResponse(newsItem *models.NewsItem) *NewsResponse {
	return &NewsResponse{
		ID:        newsItem.ID,
		RealmID:   newsItem.RealmID,
		Message:   newsItem.Message,
		Author:    newsItem.Author,
		CreatedAt: newsItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt: newsItem.UpdatedAt.Format(time.RFC3339),
	}
}
