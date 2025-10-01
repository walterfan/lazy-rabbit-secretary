package image

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/gorm"
)

// ImageService handles business logic for uploaded images
type ImageService struct {
	repo      *ImageRepository
	uploadDir string
}

// NewImageService creates a new image service
func NewImageService(db *gorm.DB, uploadDir string) *ImageService {
	return &ImageService{
		repo:      NewImageRepository(db),
		uploadDir: uploadDir,
	}
}

// UpdateImageRequest represents the request to update an image
type UpdateImageRequest struct {
	Type        models.ImageType `json:"type" validate:"oneof=avatar profile cover thumbnail gallery attachment custom"`
	Category    string           `json:"category" validate:"max=100"`
	Description string           `json:"description" validate:"max=500"`
	Tags        []string         `json:"tags"`
	IsPublic    bool             `json:"is_public"`
	IsShared    bool             `json:"is_shared"`
}

// ImageResponse represents the response for an uploaded image
type ImageResponse struct {
	ID            string             `json:"id"`
	RealmID       string             `json:"realm_id"`
	UserID        string             `json:"user_id"`
	Source        string             `json:"source"`
	Type          models.ImageType   `json:"type"`
	Status        models.ImageStatus `json:"status"`
	OriginalName  string             `json:"original_name"`
	FileName      string             `json:"file_name"`
	FilePath      string             `json:"file_path"`
	FileSize      int64              `json:"file_size"`
	MimeType      string             `json:"mime_type"`
	Extension     string             `json:"extension"`
	Width         int                `json:"width"`
	Height        int                `json:"height"`
	Format        string             `json:"format"`
	Quality       int                `json:"quality"`
	ColorSpace    string             `json:"color_space"`
	HasAlpha      bool               `json:"has_alpha"`
	ProcessedAt   *time.Time         `json:"processed_at"`
	ErrorMsg      string             `json:"error_msg"`
	IsPublic      bool               `json:"is_public"`
	IsShared      bool               `json:"is_shared"`
	ShareToken    string             `json:"share_token,omitempty"`
	Tags          []string           `json:"tags"`
	Category      string             `json:"category"`
	Description   string             `json:"description"`
	ViewCount     int64              `json:"view_count"`
	DownloadCount int64              `json:"download_count"`
	CreatedBy     string             `json:"created_by"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedBy     string             `json:"updated_by"`
	UpdatedAt     time.Time          `json:"updated_at"`
	DiagramID     string             `json:"diagram_id,omitempty"`

	// Computed fields
	URL           string  `json:"url"`
	ThumbnailURL  string  `json:"thumbnail_url"`
	DownloadURL   string  `json:"download_url"`
	AspectRatio   float64 `json:"aspect_ratio"`
	FormattedSize string  `json:"formatted_size"`
}

// ImageListResponse represents a paginated list of images
type ImageListResponse struct {
	Images  []ImageResponse `json:"images"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
	HasMore bool            `json:"has_more"`
}

// ImageStatsResponse represents storage statistics
type ImageStatsResponse struct {
	TotalCount   int64            `json:"total_count"`
	TotalSize    int64            `json:"total_size"`
	AverageSize  int64            `json:"average_size"`
	TypeCounts   map[string]int64 `json:"type_counts"`
	StatusCounts map[string]int64 `json:"status_counts"`
}

// Service Methods

// UploadImage uploads a new image file
func (s *ImageService) UploadImage(file *multipart.FileHeader, req *UploadImageRequest, realmID, userID string) (*ImageResponse, error) {
	// Validate request
	if err := s.validateUploadRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Validate file
	if err := s.validateFile(file); err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	// Generate unique file name
	fileName := s.generateFileName(file.Filename)
	filePath := filepath.Join(s.uploadDir, fileName)

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save file
	if err := s.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create image record
	image := &models.Image{
		ID:           uuid.New().String(),
		RealmID:      realmID,
		UserID:       userID,
		Source:       models.ImageSourceUploaded,
		Type:         models.ImageType(req.Type),
		Status:       models.ImageStatusUploading,
		OriginalName: file.Filename,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     fileInfo.Size(),
		MimeType:     file.Header.Get("Content-Type"),
		Extension:    strings.ToLower(filepath.Ext(file.Filename)),
		Format:       models.ImageFormat(strings.ToLower(filepath.Ext(file.Filename)[1:])), // Remove the dot
		Category:     req.Category,
		Description:  req.Description,
		Public:       req.IsPublic,
		Shared:       req.IsShared,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	// Set tags
	image.SetTags(strings.Split(req.Tags, ","))

	// Create in database
	if err := s.repo.Create(image); err != nil {
		// Clean up file if database creation fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create image record: %w", err)
	}

	// Mark as uploaded
	image.MarkAsUploaded()
	if err := s.repo.Update(image); err != nil {
		return nil, fmt.Errorf("failed to update image status: %w", err)
	}

	return s.toImageResponse(image), nil
}

// GetImage retrieves an image by ID
func (s *ImageService) GetImage(id string) (*ImageResponse, error) {
	image, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return s.toImageResponse(image), nil
}

// UpdateImage updates an existing image
func (s *ImageService) UpdateImage(id string, req *UpdateImageRequest, userID string) (*ImageResponse, error) {
	// Get existing image
	image, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	// Check if image can be edited
	if !image.CanBeEdited() {
		return nil, fmt.Errorf("image cannot be edited")
	}

	// Update fields
	if req.Type != "" {
		image.Type = req.Type
	}
	if req.Category != "" {
		image.Category = req.Category
	}
	if req.Description != "" {
		image.Description = req.Description
	}

	// Update tags if provided
	if req.Tags != nil {
		image.SetTags(req.Tags)
	}

	// Update sharing settings
	image.Public = req.IsPublic
	image.Shared = req.IsShared

	image.UpdatedBy = userID
	image.UpdatedAt = time.Now()

	// Save the image
	if err := s.repo.Update(image); err != nil {
		return nil, fmt.Errorf("failed to update image: %w", err)
	}

	return s.toImageResponse(image), nil
}

// DeleteImage deletes an image
func (s *ImageService) DeleteImage(id string) error {
	// Get existing image
	image, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("image not found")
		}
		return fmt.Errorf("failed to get image: %w", err)
	}

	// Check if image can be deleted
	if !image.CanBeDeleted() {
		return fmt.Errorf("image cannot be deleted")
	}

	// Delete file
	if err := os.Remove(image.FilePath); err != nil {
		// Log error but don't fail the deletion
		fmt.Printf("Warning: failed to delete file %s: %v\n", image.FilePath, err)
	}

	// Delete the image record
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

// ListImages retrieves a list of images with pagination
func (s *ImageService) ListImages(realmID string, page, limit int, imageType, status, category, userID string) (*ImageListResponse, error) {
	offset := (page - 1) * limit

	var images []*models.Image
	var total int64
	var err error

	// Determine query type
	switch {
	case imageType != "":
		images, err = s.repo.GetByType(models.ImageType(imageType), realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByType(models.ImageType(imageType), realmID)
		}
	case status != "":
		images, err = s.repo.GetByStatus(models.ImageStatus(status), realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByStatus(models.ImageStatus(status), realmID)
		}
	case category != "":
		images, err = s.repo.GetByCategory(category, realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByCategory(category, realmID)
		}
	case userID != "":
		images, err = s.repo.GetByUserID(userID, realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByUser(userID, realmID)
		}
	default:
		images, err = s.repo.GetByRealmID(realmID, limit, offset)
		if err == nil {
			total, err = s.repo.Count(realmID)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	// Convert to response
	response := &ImageListResponse{
		Images:  make([]ImageResponse, len(images)),
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(offset+limit) < total,
	}

	for i, image := range images {
		response.Images[i] = *s.toImageResponse(image)
	}

	return response, nil
}

// SearchImages searches images by query
func (s *ImageService) SearchImages(query, realmID string, page, limit int) (*ImageListResponse, error) {
	offset := (page - 1) * limit

	images, err := s.repo.Search(query, realmID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search images: %w", err)
	}

	total, err := s.repo.Count(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to count images: %w", err)
	}

	// Convert to response
	response := &ImageListResponse{
		Images:  make([]ImageResponse, len(images)),
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(offset+limit) < total,
	}

	for i, image := range images {
		response.Images[i] = *s.toImageResponse(image)
	}

	return response, nil
}

// GetPublicImages retrieves public images
func (s *ImageService) GetPublicImages(realmID string, page, limit int) (*ImageListResponse, error) {
	offset := (page - 1) * limit

	images, err := s.repo.GetPublic(realmID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public images: %w", err)
	}

	total, err := s.repo.Count(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to count images: %w", err)
	}

	// Convert to response
	response := &ImageListResponse{
		Images:  make([]ImageResponse, len(images)),
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(offset+limit) < total,
	}

	for i, image := range images {
		response.Images[i] = *s.toImageResponse(image)
	}

	return response, nil
}

// GetSharedImages retrieves shared images
func (s *ImageService) GetSharedImages(realmID string, page, limit int) (*ImageListResponse, error) {
	offset := (page - 1) * limit

	images, err := s.repo.GetShared(realmID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared images: %w", err)
	}

	total, err := s.repo.Count(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to count images: %w", err)
	}

	// Convert to response
	response := &ImageListResponse{
		Images:  make([]ImageResponse, len(images)),
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(offset+limit) < total,
	}

	for i, image := range images {
		response.Images[i] = *s.toImageResponse(image)
	}

	return response, nil
}

// DownloadImage downloads an image file
func (s *ImageService) DownloadImage(id string) (*models.Image, error) {
	image, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(image.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found")
	}

	// Increment download count
	if err := s.repo.IncrementDownloadCount(id); err != nil {
		// Log error but don't fail the download
		fmt.Printf("Warning: failed to increment download count: %v\n", err)
	}

	return image, nil
}

// IncrementViewCount increments the view count for an image
func (s *ImageService) IncrementViewCount(id string) error {
	return s.repo.IncrementViewCount(id)
}

// GetImageStats returns storage statistics
func (s *ImageService) GetImageStats(realmID string) (*ImageStatsResponse, error) {
	stats, err := s.repo.GetStorageStats(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage stats: %w", err)
	}

	response := &ImageStatsResponse{
		TotalCount:   stats["total_count"].(int64),
		TotalSize:    stats["total_size"].(int64),
		AverageSize:  stats["average_size"].(int64),
		TypeCounts:   stats["type_counts"].(map[string]int64),
		StatusCounts: stats["status_counts"].(map[string]int64),
	}

	return response, nil
}

// GetCategories retrieves all categories
func (s *ImageService) GetCategories(realmID string) ([]string, error) {
	return s.repo.GetCategories(realmID)
}

// GetExtensions retrieves all file extensions
func (s *ImageService) GetExtensions(realmID string) ([]string, error) {
	return s.repo.GetExtensions(realmID)
}

// GetMimeTypes retrieves all MIME types
func (s *ImageService) GetMimeTypes(realmID string) ([]string, error) {
	return s.repo.GetMimeTypes(realmID)
}

// Helper Methods

// toImageResponse converts an Image model to an ImageResponse
func (s *ImageService) toImageResponse(image *models.Image) *ImageResponse {
	return &ImageResponse{
		ID:            image.ID,
		RealmID:       image.RealmID,
		UserID:        image.UserID,
		Type:          image.Type,
		Status:        image.Status,
		OriginalName:  image.OriginalName,
		FileName:      image.FileName,
		FilePath:      image.FilePath,
		FileSize:      image.FileSize,
		MimeType:      image.MimeType,
		Extension:     image.Extension,
		Width:         image.Width,
		Height:        image.Height,
		Format:        string(image.Format),
		Quality:       image.Quality,
		ColorSpace:    image.ColorSpace,
		HasAlpha:      image.HasAlpha,
		ProcessedAt:   image.ProcessedAt,
		ErrorMsg:      image.ErrorMsg,
		IsPublic:      image.Public,
		IsShared:      image.Shared,
		Source:        string(image.Source),
		ShareToken:    image.ShareToken,
		Tags:          image.GetTags(),
		Category:      image.Category,
		Description:   image.Description,
		ViewCount:     image.ViewCount,
		DownloadCount: image.DownloadCount,
		CreatedBy:     image.CreatedBy,
		CreatedAt:     image.CreatedAt,
		UpdatedBy:     image.UpdatedBy,
		UpdatedAt:     image.UpdatedAt,
		DiagramID:     image.DiagramID,
		URL:           image.GetURL(),
		ThumbnailURL:  image.GetThumbnailURL(),
		DownloadURL:   image.GetDownloadURL(),
		AspectRatio:   image.GetAspectRatio(),
		FormattedSize: image.GetFormattedSize(),
	}
}

// generateFileName generates a unique file name
func (s *ImageService) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	// Generate MD5 hash of original name + timestamp
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d", name, time.Now().UnixNano())))

	return fmt.Sprintf("%x%s", hash, ext)
}

// saveFile saves the uploaded file to disk
func (s *ImageService) saveFile(file *multipart.FileHeader, filePath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// Validation Methods

// validateUploadRequest validates an upload request
func (s *ImageService) validateUploadRequest(req *UploadImageRequest) error {
	if req.Type == "" {
		return fmt.Errorf("type is required")
	}
	return nil
}

// validateFile validates an uploaded file
func (s *ImageService) validateFile(file *multipart.FileHeader) error {
	if file == nil {
		return fmt.Errorf("file is required")
	}

	if file.Size == 0 {
		return fmt.Errorf("file cannot be empty")
	}

	// Check file size (10MB limit)
	if file.Size > 10*1024*1024 {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}

	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	return nil
}
