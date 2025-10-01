package news

import (
	"context"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// NewsRepository handles database operations for news items
type NewsRepository struct {
	db *gorm.DB
}

// NewNewsRepository creates a new news repository
func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

// Create creates a new news item
func (r *NewsRepository) Create(ctx context.Context, newsItem *models.NewsItem) error {
	return r.db.WithContext(ctx).Create(newsItem).Error
}

// GetByID retrieves a news item by ID
func (r *NewsRepository) GetByID(ctx context.Context, id string) (*models.NewsItem, error) {
	var newsItem models.NewsItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&newsItem).Error
	if err != nil {
		return nil, err
	}
	return &newsItem, nil
}

// GetByRealm retrieves news items by realm ID
func (r *NewsRepository) GetByRealm(ctx context.Context, realmID string, limit int, offset int) ([]*models.NewsItem, error) {
	var newsItems []*models.NewsItem
	query := r.db.WithContext(ctx).Where("realm_id = ?", realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&newsItems).Error
	return newsItems, err
}

// Update updates a news item
func (r *NewsRepository) Update(ctx context.Context, newsItem *models.NewsItem) error {
	return r.db.WithContext(ctx).Save(newsItem).Error
}

// Delete soft deletes a news item
func (r *NewsRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.NewsItem{}, "id = ?", id).Error
}

// Count returns the total count of news items for a realm
func (r *NewsRepository) Count(ctx context.Context, realmID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.NewsItem{}).Where("realm_id = ?", realmID).Count(&count).Error
	return count, err
}
