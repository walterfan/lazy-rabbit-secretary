package prompt

import (
	"strings"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// PromptRepository handles database operations for prompts
type PromptRepository struct {
	db *gorm.DB
}

// NewPromptRepository creates a new prompt repository
func NewPromptRepository(db *gorm.DB) *PromptRepository {
	return &PromptRepository{db: db}
}

// Create creates a new prompt
func (r *PromptRepository) Create(prompt *models.Prompt) error {
	return r.db.Create(prompt).Error
}

// GetByID retrieves a prompt by ID
func (r *PromptRepository) GetByID(id string) (*models.Prompt, error) {
	var prompt models.Prompt
	err := r.db.Where("id = ?", id).First(&prompt).Error
	if err != nil {
		return nil, err
	}
	return &prompt, nil
}

// GetByName retrieves a prompt by name
func (r *PromptRepository) GetByName(name string) (*models.Prompt, error) {
	var prompt models.Prompt
	err := r.db.Where("name = ?", name).First(&prompt).Error
	if err != nil {
		return nil, err
	}
	return &prompt, nil
}

// Update updates an existing prompt
func (r *PromptRepository) Update(prompt *models.Prompt) error {
	return r.db.Save(prompt).Error
}

// Delete soft deletes a prompt
func (r *PromptRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Prompt{}).Error
}

// List retrieves prompts with pagination
func (r *PromptRepository) List(offset, limit int) ([]*models.Prompt, int64, error) {
	var prompts []*models.Prompt
	var total int64

	// Get total count
	if err := r.db.Model(&models.Prompt{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&prompts).Error
	if err != nil {
		return nil, 0, err
	}

	return prompts, total, nil
}

// Search searches prompts by name, description, or tags
func (r *PromptRepository) Search(query string, offset, limit int) ([]*models.Prompt, int64, error) {
	var prompts []*models.Prompt
	var total int64

	// Prepare search query
	searchPattern := "%" + strings.ToLower(query) + "%"

	baseQuery := r.db.Model(&models.Prompt{}).Where(
		"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?",
		searchPattern, searchPattern, searchPattern,
	)

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := baseQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&prompts).Error
	if err != nil {
		return nil, 0, err
	}

	return prompts, total, nil
}

// GetByTags retrieves prompts by tags
func (r *PromptRepository) GetByTags(tags []string, offset, limit int) ([]*models.Prompt, int64, error) {
	var prompts []*models.Prompt
	var total int64

	if len(tags) == 0 {
		return r.List(offset, limit)
	}

	// Build query for tags
	var conditions []string
	var args []interface{}

	for _, tag := range tags {
		conditions = append(conditions, "LOWER(tags) LIKE ?")
		args = append(args, "%"+strings.ToLower(tag)+"%")
	}

	whereClause := strings.Join(conditions, " OR ")
	baseQuery := r.db.Model(&models.Prompt{}).Where(whereClause, args...)

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := baseQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&prompts).Error
	if err != nil {
		return nil, 0, err
	}

	return prompts, total, nil
}

// GetByCreatedBy retrieves prompts by creator
func (r *PromptRepository) GetByCreatedBy(createdBy string, offset, limit int) ([]*models.Prompt, int64, error) {
	var prompts []*models.Prompt
	var total int64

	baseQuery := r.db.Model(&models.Prompt{}).Where("created_by = ?", createdBy)

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := baseQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&prompts).Error
	if err != nil {
		return nil, 0, err
	}

	return prompts, total, nil
}

// Exists checks if a prompt with the given name exists
func (r *PromptRepository) Exists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Prompt{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
