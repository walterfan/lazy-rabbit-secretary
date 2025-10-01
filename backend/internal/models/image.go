package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ImageFormat represents the format of an image
type ImageFormat string

const (
	ImageFormatSVG  ImageFormat = "svg"
	ImageFormatPNG  ImageFormat = "png"
	ImageFormatJPG  ImageFormat = "jpg"
	ImageFormatJPEG ImageFormat = "jpeg"
	ImageFormatGIF  ImageFormat = "gif"
	ImageFormatWEBP ImageFormat = "webp"
)

// ImageSource represents the source of an image
type ImageSource string

const (
	ImageSourceUploaded  ImageSource = "uploaded"  // User uploaded image
	ImageSourceGenerated ImageSource = "generated" // Generated image (e.g., from diagram)
)

// ImageType represents the type of image
type ImageType string

const (
	ImageTypeAvatar     ImageType = "avatar"     // User avatar
	ImageTypeProfile    ImageType = "profile"    // Profile image
	ImageTypeCover      ImageType = "cover"      // Cover image
	ImageTypeThumbnail  ImageType = "thumbnail"  // Thumbnail image
	ImageTypeGallery    ImageType = "gallery"    // Gallery image
	ImageTypeAttachment ImageType = "attachment" // File attachment
	ImageTypeCustom     ImageType = "custom"     // Custom image type
	ImageTypeDiagram    ImageType = "diagram"    // Diagram generated image
)

// ImageStatus represents the status of an image
type ImageStatus string

const (
	ImageStatusUploading  ImageStatus = "uploading"
	ImageStatusUploaded   ImageStatus = "uploaded"
	ImageStatusProcessing ImageStatus = "processing"
	ImageStatusProcessed  ImageStatus = "processed"
	ImageStatusFailed     ImageStatus = "failed"
	ImageStatusDeleted    ImageStatus = "deleted"
)

// Image represents an image (uploaded or generated)
type Image struct {
	ID      string      `json:"id" gorm:"primaryKey;type:text"`
	RealmID string      `json:"realm_id" gorm:"not null;type:text;index"`
	UserID  string      `json:"user_id" gorm:"type:text;index"`
	Source  ImageSource `json:"source" gorm:"not null;type:text;index"` // uploaded or generated
	Type    ImageType   `json:"type" gorm:"not null;type:text;index"`
	Status  ImageStatus `json:"status" gorm:"default:'uploaded';type:text;index"`

	// File information
	OriginalName string `json:"original_name" gorm:"type:text"`
	FileName     string `json:"file_name" gorm:"not null;type:text;uniqueIndex"`
	FilePath     string `json:"file_path" gorm:"not null;type:text"`
	FileSize     int64  `json:"file_size" gorm:"default:0"`
	MimeType     string `json:"mime_type" gorm:"type:text"`
	Extension    string `json:"extension" gorm:"type:text"`

	// Image metadata
	Width      int         `json:"width" gorm:"default:0"`
	Height     int         `json:"height" gorm:"default:0"`
	Format     ImageFormat `json:"format" gorm:"type:text"`
	Quality    int         `json:"quality" gorm:"default:90"`
	ColorSpace string      `json:"color_space" gorm:"type:text"`
	HasAlpha   bool        `json:"has_alpha" gorm:"default:false"`

	// Processing information
	ProcessedAt *time.Time `json:"processed_at"`
	ErrorMsg    string     `json:"error_msg" gorm:"type:text"`

	// Access control
	Public      bool   `json:"public" gorm:"default:false"`
	Shared      bool   `json:"shared" gorm:"default:false"`
	ShareToken  string `json:"share_token,omitempty" gorm:"type:text"`
	Permissions string `json:"permissions,omitempty" gorm:"type:text"` // JSON permissions

	// Organization
	Tags        string `json:"tags" gorm:"type:text"` // comma-separated tags
	Category    string `json:"category" gorm:"type:text"`
	Description string `json:"description" gorm:"type:text"`

	// Usage tracking
	ViewCount     int64 `json:"view_count" gorm:"default:0"`
	DownloadCount int64 `json:"download_count" gorm:"default:0"`

	// Diagram-specific fields (for generated images)
	DiagramID string `json:"diagram_id" gorm:"type:text;index"`
	IsPrimary bool   `json:"is_primary" gorm:"default:false"` // Primary image for the diagram

	// Audit fields
	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationship (optional, only for diagram images)
	Diagram *Diagram `json:"diagram,omitempty" gorm:"foreignKey:DiagramID;references:ID"`
}

// TableName specifies the table name for GORM
func (Image) TableName() string {
	return "images"
}

// IsSVG returns true if the image is in SVG format
func (i *Image) IsSVG() bool {
	return i.Format == ImageFormatSVG
}

// IsPNG returns true if the image is in PNG format
func (i *Image) IsPNG() bool {
	return i.Format == ImageFormatPNG
}

// IsJPG returns true if the image is in JPG format
func (i *Image) IsJPG() bool {
	return i.Format == ImageFormatJPG || i.Format == ImageFormatJPEG
}

// GetURL returns the URL for the image
func (i *Image) GetURL() string {
	if i.IsGeneratedSource() && i.DiagramID != "" {
		return "/api/v1/diagrams/" + i.DiagramID + "/images/" + i.ID
	}
	return "/api/v1/images/" + i.ID
}

// GetThumbnailURL returns the URL for the image thumbnail
func (i *Image) GetThumbnailURL() string {
	if i.IsGeneratedSource() && i.DiagramID != "" {
		return "/api/v1/diagrams/" + i.DiagramID + "/images/" + i.ID + "/thumbnail"
	}
	return "/api/v1/images/" + i.ID + "/thumbnail"
}

// GetAspectRatio returns the aspect ratio of the image
func (i *Image) GetAspectRatio() float64 {
	if i.Height == 0 {
		return 0
	}
	return float64(i.Width) / float64(i.Height)
}

// IsLandscape returns true if the image is landscape orientation
func (i *Image) IsLandscape() bool {
	return i.Width > i.Height
}

// IsPortrait returns true if the image is portrait orientation
func (i *Image) IsPortrait() bool {
	return i.Height > i.Width
}

// IsSquare returns true if the image is square
func (i *Image) IsSquare() bool {
	return i.Width == i.Height
}

// GetFormattedSize returns the file size in human-readable format
func (i *Image) GetFormattedSize() string {
	const unit = 1024
	if i.FileSize < unit {
		return fmt.Sprintf("%d B", i.FileSize)
	}
	div, exp := int64(unit), 0
	for n := i.FileSize / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(i.FileSize)/float64(div), "KMGTPE"[exp])
}

// Source and Status Methods

// IsUploaded returns true if the image is successfully uploaded
func (i *Image) IsUploaded() bool {
	return i.Status == ImageStatusUploaded || i.Status == ImageStatusProcessed
}

// IsProcessing returns true if the image is being processed
func (i *Image) IsProcessing() bool {
	return i.Status == ImageStatusProcessing
}

// IsFailed returns true if the image upload/processing failed
func (i *Image) IsFailed() bool {
	return i.Status == ImageStatusFailed
}

// IsUploadedSource returns true if the image is from user upload
func (i *Image) IsUploadedSource() bool {
	return i.Source == ImageSourceUploaded
}

// IsGeneratedSource returns true if the image is generated (e.g., from diagram)
func (i *Image) IsGeneratedSource() bool {
	return i.Source == ImageSourceGenerated
}

// IsPublic returns true if the image is publicly accessible
func (i *Image) IsPublic() bool {
	return i.Public
}

// IsShared returns true if the image is shared
func (i *Image) IsShared() bool {
	return i.Shared
}

// IsImage returns true if the file is an image
func (i *Image) IsImage() bool {
	return i.MimeType != "" && (i.MimeType[:5] == "image" || i.Extension != "")
}

// Type checking methods
func (i *Image) IsAvatar() bool {
	return i.Type == ImageTypeAvatar
}

func (i *Image) IsProfile() bool {
	return i.Type == ImageTypeProfile
}

func (i *Image) IsCover() bool {
	return i.Type == ImageTypeCover
}

func (i *Image) IsThumbnail() bool {
	return i.Type == ImageTypeThumbnail
}

func (i *Image) IsGallery() bool {
	return i.Type == ImageTypeGallery
}

func (i *Image) IsAttachment() bool {
	return i.Type == ImageTypeAttachment
}

func (i *Image) IsDiagram() bool {
	return i.Type == ImageTypeDiagram
}

// GetTags returns the tags as a slice
func (i *Image) GetTags() []string {
	if i.Tags == "" {
		return []string{}
	}
	tags := strings.Split(i.Tags, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// SetTags sets tags from a slice
func (i *Image) SetTags(tags []string) {
	if len(tags) == 0 {
		i.Tags = ""
		return
	}
	// Clean and join tags
	cleaned := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	i.Tags = strings.Join(cleaned, ",")
}

// IncrementViewCount increments the view count
func (i *Image) IncrementViewCount() {
	i.ViewCount++
}

// IncrementDownloadCount increments the download count
func (i *Image) IncrementDownloadCount() {
	i.DownloadCount++
}

// Status management methods
func (i *Image) MarkAsUploaded() {
	i.Status = ImageStatusUploaded
	now := time.Now()
	i.ProcessedAt = &now
}

func (i *Image) MarkAsProcessing() {
	i.Status = ImageStatusProcessing
}

func (i *Image) MarkAsProcessed() {
	i.Status = ImageStatusProcessed
	now := time.Now()
	i.ProcessedAt = &now
}

func (i *Image) MarkAsFailed(errorMsg string) {
	i.Status = ImageStatusFailed
	i.ErrorMsg = errorMsg
}

// Access control methods
func (i *Image) MakePublic() {
	i.Public = true
}

func (i *Image) MakePrivate() {
	i.Public = false
	i.Shared = false
	i.ShareToken = ""
}

func (i *Image) Share() {
	i.Shared = true
}

func (i *Image) Unshare() {
	i.Shared = false
	i.ShareToken = ""
}

// Permission methods
func (i *Image) CanBeEdited() bool {
	return i.Status != ImageStatusDeleted
}

func (i *Image) CanBeDeleted() bool {
	return i.Status != ImageStatusDeleted
}

// Utility methods
func (i *Image) HasDimensions() bool {
	return i.Width > 0 && i.Height > 0
}

// GetDownloadURL returns the URL for downloading the image
func (i *Image) GetDownloadURL() string {
	if i.IsGeneratedSource() && i.DiagramID != "" {
		return "/api/v1/diagrams/" + i.DiagramID + "/images/" + i.ID + "/download"
	}
	return "/api/v1/images/" + i.ID + "/download"
}
