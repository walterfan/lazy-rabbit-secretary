package image

import (
	"database/sql"
	"strings"

	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// ImageRepository handles database operations for uploaded images
type ImageRepository struct {
	db *gorm.DB
}

// NewImageRepository creates a new image repository
func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// Create creates a new uploaded image
func (r *ImageRepository) Create(image *models.Image) error {
	return r.db.Create(image).Error
}

// GetByID retrieves an uploaded image by ID
func (r *ImageRepository) GetByID(id string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByFileName retrieves an uploaded image by file name
func (r *ImageRepository) GetByFileName(fileName string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("file_name = ?", fileName).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByRealmID retrieves images by realm ID
func (r *ImageRepository) GetByRealmID(realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("realm_id = ?", realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByUserID retrieves images by user ID
func (r *ImageRepository) GetByUserID(userID string, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("user_id = ? AND realm_id = ?", userID, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByType retrieves images by type
func (r *ImageRepository) GetByType(imageType models.ImageType, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("type = ? AND realm_id = ?", imageType, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByStatus retrieves images by status
func (r *ImageRepository) GetByStatus(status models.ImageStatus, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("status = ? AND realm_id = ?", status, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetPublic retrieves public images
func (r *ImageRepository) GetPublic(realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("is_public = ? AND realm_id = ?", true, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetShared retrieves shared images
func (r *ImageRepository) GetShared(realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("is_shared = ? AND realm_id = ?", true, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// Search searches images by name, description, or tags
func (r *ImageRepository) Search(query string, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	searchQuery := r.db.Where("realm_id = ?", realmID)

	if query != "" {
		searchPattern := "%" + strings.ToLower(query) + "%"
		searchQuery = searchQuery.Where(
			"LOWER(original_name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if limit > 0 {
		searchQuery = searchQuery.Limit(limit)
	}
	if offset > 0 {
		searchQuery = searchQuery.Offset(offset)
	}

	err := searchQuery.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByCategory retrieves images by category
func (r *ImageRepository) GetByCategory(category string, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("category = ? AND realm_id = ?", category, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByMimeType retrieves images by MIME type
func (r *ImageRepository) GetByMimeType(mimeType string, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("mime_type = ? AND realm_id = ?", mimeType, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByExtension retrieves images by file extension
func (r *ImageRepository) GetByExtension(extension string, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("extension = ? AND realm_id = ?", extension, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// Update updates an existing uploaded image
func (r *ImageRepository) Update(image *models.Image) error {
	return r.db.Save(image).Error
}

// Delete soft deletes an uploaded image
func (r *ImageRepository) Delete(id string) error {
	return r.db.Delete(&models.Image{}, "id = ?", id).Error
}

// Count returns the total count of images
func (r *ImageRepository) Count(realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Image{}).Where("realm_id = ?", realmID).Count(&count).Error
	return count, err
}

// CountByType returns the count of images by type
func (r *ImageRepository) CountByType(imageType models.ImageType, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Image{}).Where("type = ? AND realm_id = ?", imageType, realmID).Count(&count).Error
	return count, err
}

// CountByStatus returns the count of images by status
func (r *ImageRepository) CountByStatus(status models.ImageStatus, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Image{}).Where("status = ? AND realm_id = ?", status, realmID).Count(&count).Error
	return count, err
}

// CountByUser returns the count of images by user
func (r *ImageRepository) CountByUser(userID string, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Image{}).Where("user_id = ? AND realm_id = ?", userID, realmID).Count(&count).Error
	return count, err
}

// CountByCategory returns the count of images by category
func (r *ImageRepository) CountByCategory(category string, realmID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Image{}).Where("category = ? AND realm_id = ?", category, realmID).Count(&count).Error
	return count, err
}

// GetRecent retrieves recently uploaded images
func (r *ImageRepository) GetRecent(realmID string, limit int) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("realm_id = ?", realmID).
		Order("created_at DESC").
		Limit(limit).
		Find(&images).Error
	return images, err
}

// GetMostViewed retrieves most viewed images
func (r *ImageRepository) GetMostViewed(realmID string, limit int) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("realm_id = ?", realmID).
		Order("view_count DESC").
		Limit(limit).
		Find(&images).Error
	return images, err
}

// GetMostDownloaded retrieves most downloaded images
func (r *ImageRepository) GetMostDownloaded(realmID string, limit int) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("realm_id = ?", realmID).
		Order("download_count DESC").
		Limit(limit).
		Find(&images).Error
	return images, err
}

// GetLargest retrieves largest images by file size
func (r *ImageRepository) GetLargest(realmID string, limit int) ([]*models.Image, error) {
	var images []*models.Image
	err := r.db.Where("realm_id = ?", realmID).
		Order("file_size DESC").
		Limit(limit).
		Find(&images).Error
	return images, err
}

// GetBySizeRange retrieves images within a size range
func (r *ImageRepository) GetBySizeRange(minSize, maxSize int64, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("file_size >= ? AND file_size <= ? AND realm_id = ?", minSize, maxSize, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// GetByDateRange retrieves images within a date range
func (r *ImageRepository) GetByDateRange(startDate, endDate string, realmID string, limit, offset int) ([]*models.Image, error) {
	var images []*models.Image
	query := r.db.Where("created_at >= ? AND created_at <= ? AND realm_id = ?", startDate, endDate, realmID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&images).Error
	return images, err
}

// IncrementViewCount increments the view count for an image
func (r *ImageRepository) IncrementViewCount(id string) error {
	return r.db.Model(&models.Image{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementDownloadCount increments the download count for an image
func (r *ImageRepository) IncrementDownloadCount(id string) error {
	return r.db.Model(&models.Image{}).Where("id = ?", id).Update("download_count", gorm.Expr("download_count + 1")).Error
}

// GetCategories retrieves all unique categories
func (r *ImageRepository) GetCategories(realmID string) ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Image{}).
		Where("realm_id = ? AND category != ''", realmID).
		Distinct("category").
		Pluck("category", &categories).Error
	return categories, err
}

// GetExtensions retrieves all unique file extensions
func (r *ImageRepository) GetExtensions(realmID string) ([]string, error) {
	var extensions []string
	err := r.db.Model(&models.Image{}).
		Where("realm_id = ? AND extension != ''", realmID).
		Distinct("extension").
		Pluck("extension", &extensions).Error
	return extensions, err
}

// GetMimeTypes retrieves all unique MIME types
func (r *ImageRepository) GetMimeTypes(realmID string) ([]string, error) {
	var mimeTypes []string
	err := r.db.Model(&models.Image{}).
		Where("realm_id = ? AND mime_type != ''", realmID).
		Distinct("mime_type").
		Pluck("mime_type", &mimeTypes).Error
	return mimeTypes, err
}

// GetStorageStats returns storage statistics
func (r *ImageRepository) GetStorageStats(realmID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total count
	var totalCount int64
	err := r.db.Model(&models.Image{}).Where("realm_id = ?", realmID).Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	stats["total_count"] = totalCount

	// Total size
	var totalSize sql.NullInt64
	err = r.db.Model(&models.Image{}).Where("realm_id = ?", realmID).Select("SUM(file_size)").Scan(&totalSize).Error
	if err != nil {
		return nil, err
	}
	if totalSize.Valid {
		stats["total_size"] = totalSize.Int64
	} else {
		stats["total_size"] = int64(0)
	}

	// Average size
	if totalCount > 0 && totalSize.Valid {
		stats["average_size"] = totalSize.Int64 / totalCount
	} else {
		stats["average_size"] = int64(0)
	}

	// Count by type
	typeCounts := make(map[string]int64)
	var results []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	err = r.db.Model(&models.Image{}).
		Where("realm_id = ?", realmID).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		typeCounts[result.Type] = result.Count
	}
	stats["type_counts"] = typeCounts

	// Count by status
	statusCounts := make(map[string]int64)
	var statusResults []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	err = r.db.Model(&models.Image{}).
		Where("realm_id = ?", realmID).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusResults).Error
	if err != nil {
		return nil, err
	}

	for _, result := range statusResults {
		statusCounts[result.Status] = result.Count
	}
	stats["status_counts"] = statusCounts

	return stats, nil
}
