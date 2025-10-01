package diagram

import (
	"strings"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// DiagramRepository handles database operations for diagrams
type DiagramRepository struct {
	db *gorm.DB
}

// NewDiagramRepository creates a new diagram repository
func NewDiagramRepository(db *gorm.DB) *DiagramRepository {
	return &DiagramRepository{db: db}
}

// Create creates a new diagram
func (r *DiagramRepository) Create(diagram *models.Diagram) error {
	return r.db.Create(diagram).Error
}

// GetByID retrieves a diagram by ID
func (r *DiagramRepository) GetByID(id string) (*models.Diagram, error) {
	var diagram models.Diagram
	err := r.db.Where("id = ?", id).First(&diagram).Error
	if err != nil {
		return nil, err
	}
	return &diagram, nil
}

// GetByIDWithImages retrieves a diagram by ID with associated images
func (r *DiagramRepository) GetByIDWithImages(id string) (*models.Diagram, error) {
	var diagram models.Diagram
	err := r.db.Preload("Images").Where("id = ?", id).First(&diagram).Error
	if err != nil {
		return nil, err
	}
	return &diagram, nil
}

// GetByRealmID retrieves diagrams by realm ID
func (r *DiagramRepository) GetByRealmID(realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("realm_id = ?", realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// GetByType retrieves diagrams by type
func (r *DiagramRepository) GetByType(diagramType models.DiagramType, realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("type = ? AND realm_id = ?", diagramType, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// GetByScriptType retrieves diagrams by script type
func (r *DiagramRepository) GetByScriptType(scriptType models.DiagramScriptType, realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("script_type = ? AND realm_id = ?", scriptType, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// GetByStatus retrieves diagrams by status
func (r *DiagramRepository) GetByStatus(status models.DiagramStatus, realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("status = ? AND realm_id = ?", status, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// GetByUser retrieves diagrams by user
func (r *DiagramRepository) GetByUser(userID string, realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("created_by = ? AND realm_id = ?", userID, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// Search searches diagrams by name, description, or tags
func (r *DiagramRepository) Search(query string, realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	searchQuery := r.db.Where("realm_id = ?", realmID)

	if query != "" {
		searchPattern := "%" + strings.ToLower(query) + "%"
		searchQuery = searchQuery.Where(
			"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if limit > 0 {
		searchQuery = searchQuery.Limit(limit)
	}
	if offset > 0 {
		searchQuery = searchQuery.Offset(offset)
	}

	err := searchQuery.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// GetPublic retrieves public diagrams
func (r *DiagramRepository) GetPublic(realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("public = ? AND realm_id = ?", true, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// GetShared retrieves shared diagrams
func (r *DiagramRepository) GetShared(realmID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Where("shared = ? AND realm_id = ?", true, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// Update updates an existing diagram
func (r *DiagramRepository) Update(diagram *models.Diagram) error {
	return r.db.Save(diagram).Error
}

// Delete soft deletes a diagram
func (r *DiagramRepository) Delete(id string) error {
	return r.db.Delete(&models.Diagram{}, "id = ?", id).Error
}

// Count returns the total count of diagrams
func (r *DiagramRepository) Count(realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Diagram{}).Where("realm_id = ?", realmID).Count(&count).Error
	return count, err
}

// CountByType returns the count of diagrams by type
func (r *DiagramRepository) CountByType(diagramType models.DiagramType, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Diagram{}).Where("type = ? AND realm_id = ?", diagramType, realmID).Count(&count).Error
	return count, err
}

// CountByScriptType returns the count of diagrams by script type
func (r *DiagramRepository) CountByScriptType(scriptType models.DiagramScriptType, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Diagram{}).Where("script_type = ? AND realm_id = ?", scriptType, realmID).Count(&count).Error
	return count, err
}

// CountByStatus returns the count of diagrams by status
func (r *DiagramRepository) CountByStatus(status models.DiagramStatus, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Diagram{}).Where("status = ? AND realm_id = ?", status, realmID).Count(&count).Error
	return count, err
}

// CountByUser returns the count of diagrams by user
func (r *DiagramRepository) CountByUser(userID string, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Diagram{}).Where("created_by = ? AND realm_id = ?", userID, realmID).Count(&count).Error
	return count, err
}

// GetRecent retrieves recently created diagrams
func (r *DiagramRepository) GetRecent(realmID string, limit int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	err := r.db.Where("realm_id = ?", realmID).
		Order("created_at DESC").
		Limit(limit).
		Find(&diagrams).Error
	return diagrams, err
}

// GetMostViewed retrieves most viewed diagrams
func (r *DiagramRepository) GetMostViewed(realmID string, limit int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	err := r.db.Where("realm_id = ?", realmID).
		Order("view_count DESC").
		Limit(limit).
		Find(&diagrams).Error
	return diagrams, err
}

// GetMostEdited retrieves most edited diagrams
func (r *DiagramRepository) GetMostEdited(realmID string, limit int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	err := r.db.Where("realm_id = ?", realmID).
		Order("edit_count DESC").
		Limit(limit).
		Find(&diagrams).Error
	return diagrams, err
}

// IncrementViewCount increments the view count for a diagram
func (r *DiagramRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.Diagram{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementEditCount increments the edit count for a diagram
func (r *DiagramRepository) IncrementEditCount(id string) error {
	return r.db.Model(&models.Diagram{}).Where("id = ?", id).Update("edit_count", gorm.Expr("edit_count + 1")).Error
}

// Image Repository Methods

// CreateImage creates a new image
func (r *DiagramRepository) CreateImage(image *models.Image) error {
	return r.db.Create(image).Error
}

// GetImageByID retrieves an image by ID
func (r *DiagramRepository) GetImageByID(id string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImagesByDiagramID retrieves all images for a diagram
func (r *DiagramRepository) GetImagesByDiagramID(diagramID string) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("diagram_id = ?", diagramID).Order("created_at ASC").Find(&images).Error
	return images, err
}

// GetPrimaryImage retrieves the primary image for a diagram
func (r *DiagramRepository) GetPrimaryImage(diagramID string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("diagram_id = ? AND is_primary = ?", diagramID, true).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// UpdateImage updates an existing image
func (r *DiagramRepository) UpdateImage(image *models.Image) error {
	return r.db.Save(image).Error
}

// DeleteImage soft deletes an image
func (r *DiagramRepository) DeleteImage(id string) error {
	return r.db.Delete(&models.Image{}, "id = ?", id).Error
}

// SetPrimaryImage sets an image as primary and unsets others
func (r *DiagramRepository) SetPrimaryImage(diagramID, imageID string) error {
	// First, unset all primary images for this diagram
	err := r.db.Model(&models.Image{}).Where("diagram_id = ?", diagramID).Update("is_primary", false).Error
	if err != nil {
		return err
	}

	// Then set the specified image as primary
	return r.db.Model(&models.Image{}).Where("id = ?", imageID).Update("is_primary", true).Error
}

// DiagramTag Repository Methods

// CreateTag creates a new diagram tag
func (r *DiagramRepository) CreateTag(tag *models.DiagramTag) error {
	return r.db.Create(tag).Error
}

// GetTagByID retrieves a tag by ID
func (r *DiagramRepository) GetTagByID(id string) (*models.DiagramTag, error) {
	var tag models.DiagramTag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetTagBySlug retrieves a tag by slug
func (r *DiagramRepository) GetTagBySlug(slug string) (*models.DiagramTag, error) {
	var tag models.DiagramTag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAllTags retrieves all tags
func (r *DiagramRepository) GetAllTags(limit, offset int) ([]*models.DiagramTag, error) {
	var tags []*models.DiagramTag
	query := r.db

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("name ASC").Find(&tags).Error
	return tags, err
}

// UpdateTag updates an existing tag
func (r *DiagramRepository) UpdateTag(tag *models.DiagramTag) error {
	return r.db.Save(tag).Error
}

// DeleteTag soft deletes a tag
func (r *DiagramRepository) DeleteTag(id string) error {
	return r.db.Delete(&models.DiagramTag{}, "id = ?", id).Error
}

// DiagramTagRelation Repository Methods

// CreateTagRelation creates a new diagram-tag relation
func (r *DiagramRepository) CreateTagRelation(relation *models.DiagramTagRelation) error {
	return r.db.Create(relation).Error
}

// GetTagsByDiagramID retrieves all tags for a diagram
func (r *DiagramRepository) GetTagsByDiagramID(diagramID string) ([]*models.DiagramTag, error) {
	var tags []*models.DiagramTag
	err := r.db.Joins("JOIN diagram_tag_relations ON diagram_tags.id = diagram_tag_relations.tag_id").
		Where("diagram_tag_relations.diagram_id = ?", diagramID).
		Find(&tags).Error
	return tags, err
}

// GetDiagramsByTagID retrieves all diagrams for a tag
func (r *DiagramRepository) GetDiagramsByTagID(tagID string, limit, offset int) ([]*models.Diagram, error) {
	var diagrams []*models.Diagram
	query := r.db.Joins("JOIN diagram_tag_relations ON diagrams.id = diagram_tag_relations.diagram_id").
		Where("diagram_tag_relations.tag_id = ?", tagID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&diagrams).Error
	return diagrams, err
}

// DeleteTagRelation deletes a diagram-tag relation
func (r *DiagramRepository) DeleteTagRelation(diagramID, tagID string) error {
	return r.db.Where("diagram_id = ? AND tag_id = ?", diagramID, tagID).Delete(&models.DiagramTagRelation{}).Error
}

// DeleteAllTagRelations deletes all tag relations for a diagram
func (r *DiagramRepository) DeleteAllTagRelations(diagramID string) error {
	return r.db.Where("diagram_id = ?", diagramID).Delete(&models.DiagramTagRelation{}).Error
}
