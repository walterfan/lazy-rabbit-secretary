package diagram

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DiagramService handles business logic for diagrams
type DiagramService struct {
	repo   *DiagramRepository
	logger *zap.SugaredLogger
}

// NewDiagramService creates a new diagram service
func NewDiagramService(db *gorm.DB) *DiagramService {
	return &DiagramService{
		repo:   NewDiagramRepository(db),
		logger: log.GetLogger(),
	}
}

// Request/Response structures

// CreateDiagramRequest represents the request to create a diagram
type CreateDiagramRequest struct {
	Name        string                   `json:"name" binding:"required" validate:"required,min=1,max=256"`
	Type        models.DiagramType       `json:"type" validate:"required,oneof=flowchart sequence class mindmap architecture custom"`
	ScriptType  models.DiagramScriptType `json:"script_type" validate:"required,oneof=mermaid plantuml graphviz"`
	Status      models.DiagramStatus     `json:"status" validate:"oneof=draft published private archived"`
	Description string                   `json:"description" validate:"max=1000"`
	Script      string                   `json:"script" validate:"required,min=1"`
	Theme       string                   `json:"theme" validate:"max=50"`
	Tags        []string                 `json:"tags"`
	Public      bool                     `json:"public"`
	Shared      bool                     `json:"shared"`
	Language    string                   `json:"language" validate:"max=10"`
}

// UpdateDiagramRequest represents the request to update a diagram
type UpdateDiagramRequest struct {
	Name        string                   `json:"name" validate:"min=1,max=256"`
	Type        models.DiagramType       `json:"type" validate:"oneof=flowchart sequence class mindmap architecture custom"`
	ScriptType  models.DiagramScriptType `json:"script_type" validate:"oneof=mermaid plantuml graphviz"`
	Status      models.DiagramStatus     `json:"status" validate:"oneof=draft published private archived"`
	Description string                   `json:"description" validate:"max=1000"`
	Script      string                   `json:"script" validate:"min=1"`
	Theme       string                   `json:"theme" validate:"max=50"`
	Tags        []string                 `json:"tags"`
	Public      bool                     `json:"public"`
	Shared      bool                     `json:"shared"`
	Language    string                   `json:"language" validate:"max=10"`
}

// DiagramResponse represents the response for a diagram
type DiagramResponse struct {
	ID          string                   `json:"id"`
	RealmID     string                   `json:"realm_id"`
	Name        string                   `json:"name"`
	Type        models.DiagramType       `json:"type"`
	ScriptType  models.DiagramScriptType `json:"script_type"`
	Status      models.DiagramStatus     `json:"status"`
	Description string                   `json:"description"`
	Script      string                   `json:"script"`
	Theme       string                   `json:"theme"`
	Tags        []string                 `json:"tags"`
	Version     int                      `json:"version"`
	ViewCount   int64                    `json:"view_count"`
	EditCount   int64                    `json:"edit_count"`
	Public      bool                     `json:"public"`
	Shared      bool                     `json:"shared"`
	ShareToken  string                   `json:"share_token,omitempty"`
	Language    string                   `json:"language"`
	CreatedBy   string                   `json:"created_by"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedBy   string                   `json:"updated_by"`
	UpdatedAt   time.Time                `json:"updated_at"`

	// Associated data
	Images     []ImageResponse      `json:"images,omitempty"`
	TagObjects []DiagramTagResponse `json:"tag_objects,omitempty"`
}

// ImageResponse represents the response for an image
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

// DiagramTagResponse represents the response for a diagram tag
type DiagramTagResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Count       int64     `json:"count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DiagramListResponse represents a paginated list of diagrams
type DiagramListResponse struct {
	Diagrams []DiagramResponse `json:"diagrams"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	HasMore  bool              `json:"has_more"`
}

// CreateImageRequest represents the request to create an image
type CreateImageRequest struct {
	DiagramID string             `json:"diagram_id" binding:"required"`
	Path      string             `json:"path" binding:"required"`
	Format    models.ImageFormat `json:"format" binding:"required"`
	Width     int                `json:"width" binding:"required"`
	Height    int                `json:"height" binding:"required"`
	Size      int64              `json:"size"`
	IsPrimary bool               `json:"is_primary"`
	Quality   int                `json:"quality" validate:"min=1,max=100"`
}

// UpdateImageRequest represents the request to update an image
type UpdateImageRequest struct {
	Path      string             `json:"path"`
	Format    models.ImageFormat `json:"format"`
	Width     int                `json:"width"`
	Height    int                `json:"height"`
	Size      int64              `json:"size"`
	IsPrimary bool               `json:"is_primary"`
	Quality   int                `json:"quality" validate:"min=1,max=100"`
}

// CreateTagRequest represents the request to create a tag
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"max=7"` // Hex color
}

// UpdateTagRequest represents the request to update a tag
type UpdateTagRequest struct {
	Name        string `json:"name" validate:"min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"max=7"` // Hex color
}

// DrawDiagramRequest represents the request to draw/generate a diagram
type DrawDiagramRequest struct {
	Script     string                   `json:"script" binding:"required"`
	ScriptType models.DiagramScriptType `json:"script_type" binding:"required"`
	Format     models.ImageFormat       `json:"format" binding:"required"`
	Width      int                      `json:"width" validate:"min=100,max=4000"`
	Height     int                      `json:"height" validate:"min=100,max=4000"`
	Theme      string                   `json:"theme" validate:"max=50"`
	Quality    int                      `json:"quality" validate:"min=1,max=100"`
}

// DrawDiagramResponse represents the response for drawing a diagram
type DrawDiagramResponse struct {
	ImageData string             `json:"image_data"` // Base64 encoded image
	Format    models.ImageFormat `json:"format"`
	Width     int                `json:"width"`
	Height    int                `json:"height"`
	Size      int64              `json:"size"`
	URL       string             `json:"url,omitempty"` // If saved to file
}

// Service Methods

// CreateDiagram creates a new diagram
func (s *DiagramService) CreateDiagram(req *CreateDiagramRequest, realmID, userID string) (*DiagramResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create diagram
	diagram := &models.Diagram{
		ID:          uuid.New().String(),
		RealmID:     realmID,
		Name:        req.Name,
		Type:        req.Type,
		ScriptType:  req.ScriptType,
		Status:      req.Status,
		Description: req.Description,
		Script:      req.Script,
		Theme:       req.Theme,
		Public:      req.Public,
		Shared:      req.Shared,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	// Set tags
	diagram.SetTags(req.Tags)

	// Set default values
	if diagram.Status == "" {
		diagram.Status = models.DiagramStatusDraft
	}
	if diagram.Theme == "" {
		diagram.Theme = "default"
	}

	// Create in database
	if err := s.repo.Create(diagram); err != nil {
		return nil, fmt.Errorf("failed to create diagram: %w", err)
	}

	// Create tag relations
	if err := s.createTagRelations(diagram.ID, req.Tags); err != nil {
		// Log error but don't fail the creation
		fmt.Printf("Warning: failed to create tag relations: %v\n", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// GetDiagram retrieves a diagram by ID
func (s *DiagramService) GetDiagram(id string, includeImages bool) (*DiagramResponse, error) {
	var diagram *models.Diagram
	var err error

	if includeImages {
		diagram, err = s.repo.GetByIDWithImages(id)
	} else {
		diagram, err = s.repo.GetByID(id)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// UpdateDiagram updates an existing diagram
func (s *DiagramService) UpdateDiagram(id string, req *UpdateDiagramRequest, userID string) (*DiagramResponse, error) {
	// Get existing diagram
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	// Check if diagram can be edited
	if !diagram.CanBeEdited() {
		return nil, fmt.Errorf("diagram cannot be edited (archived)")
	}

	// Validate request
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update fields
	if req.Name != "" {
		diagram.Name = req.Name
	}
	if req.Type != "" {
		diagram.Type = req.Type
	}
	if req.ScriptType != "" {
		diagram.ScriptType = req.ScriptType
	}
	if req.Status != "" {
		diagram.Status = req.Status
	}
	if req.Description != "" {
		diagram.Description = req.Description
	}
	if req.Script != "" {
		diagram.Script = req.Script
	}
	if req.Theme != "" {
		diagram.Theme = req.Theme
	}

	// Update tags if provided
	if req.Tags != nil {
		diagram.SetTags(req.Tags)
		// Update tag relations
		if err := s.updateTagRelations(diagram.ID, req.Tags); err != nil {
			// Log error but don't fail the update
			fmt.Printf("Warning: failed to update tag relations: %v\n", err)
		}
	}

	// Update sharing settings
	diagram.Public = req.Public
	diagram.Shared = req.Shared

	// Increment version and edit count
	diagram.IncrementEditCount()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	// Save the diagram
	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to update diagram: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// DeleteDiagram deletes a diagram
func (s *DiagramService) DeleteDiagram(id string) error {
	// Get existing diagram
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("diagram not found")
		}
		return fmt.Errorf("failed to get diagram: %w", err)
	}

	// Check if diagram can be deleted
	if !diagram.CanBeDeleted() {
		return fmt.Errorf("diagram cannot be deleted (archived)")
	}

	// Delete tag relations
	if err := s.repo.DeleteAllTagRelations(id); err != nil {
		// Log error but don't fail the deletion
		fmt.Printf("Warning: failed to delete tag relations: %v\n", err)
	}

	// Delete the diagram
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete diagram: %w", err)
	}

	return nil
}

// ListDiagrams retrieves a list of diagrams with pagination
func (s *DiagramService) ListDiagrams(realmID string, page, limit int, diagramType, scriptType, status, userID string) (*DiagramListResponse, error) {
	offset := (page - 1) * limit

	var diagrams []*models.Diagram
	var total int64
	var err error

	// Determine query type
	switch {
	case diagramType != "":
		diagrams, err = s.repo.GetByType(models.DiagramType(diagramType), realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByType(models.DiagramType(diagramType), realmID)
		}
	case scriptType != "":
		diagrams, err = s.repo.GetByScriptType(models.DiagramScriptType(scriptType), realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByScriptType(models.DiagramScriptType(scriptType), realmID)
		}
	case status != "":
		diagrams, err = s.repo.GetByStatus(models.DiagramStatus(status), realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByStatus(models.DiagramStatus(status), realmID)
		}
	case userID != "":
		diagrams, err = s.repo.GetByUser(userID, realmID, limit, offset)
		if err == nil {
			total, err = s.repo.CountByUser(userID, realmID)
		}
	default:
		diagrams, err = s.repo.GetByRealmID(realmID, limit, offset)
		if err == nil {
			total, err = s.repo.Count(realmID)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list diagrams: %w", err)
	}

	// Convert to response
	response := &DiagramListResponse{
		Diagrams: make([]DiagramResponse, len(diagrams)),
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasMore:  int64(offset+limit) < total,
	}

	for i, diagram := range diagrams {
		response.Diagrams[i] = *s.toDiagramResponse(diagram)
	}

	return response, nil
}

// SearchDiagrams searches diagrams by query
func (s *DiagramService) SearchDiagrams(query, realmID string, page, limit int) (*DiagramListResponse, error) {
	offset := (page - 1) * limit

	diagrams, err := s.repo.Search(query, realmID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search diagrams: %w", err)
	}

	total, err := s.repo.Count(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to count diagrams: %w", err)
	}

	// Convert to response
	response := &DiagramListResponse{
		Diagrams: make([]DiagramResponse, len(diagrams)),
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasMore:  int64(offset+limit) < total,
	}

	for i, diagram := range diagrams {
		response.Diagrams[i] = *s.toDiagramResponse(diagram)
	}

	return response, nil
}

// GetPublicDiagrams retrieves public diagrams
func (s *DiagramService) GetPublicDiagrams(realmID string, page, limit int) (*DiagramListResponse, error) {
	offset := (page - 1) * limit

	diagrams, err := s.repo.GetPublic(realmID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public diagrams: %w", err)
	}

	total, err := s.repo.Count(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to count diagrams: %w", err)
	}

	// Convert to response
	response := &DiagramListResponse{
		Diagrams: make([]DiagramResponse, len(diagrams)),
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasMore:  int64(offset+limit) < total,
	}

	for i, diagram := range diagrams {
		response.Diagrams[i] = *s.toDiagramResponse(diagram)
	}

	return response, nil
}

// GetSharedDiagrams retrieves shared diagrams
func (s *DiagramService) GetSharedDiagrams(realmID string, page, limit int) (*DiagramListResponse, error) {
	offset := (page - 1) * limit

	diagrams, err := s.repo.GetShared(realmID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared diagrams: %w", err)
	}

	total, err := s.repo.Count(realmID)
	if err != nil {
		return nil, fmt.Errorf("failed to count diagrams: %w", err)
	}

	// Convert to response
	response := &DiagramListResponse{
		Diagrams: make([]DiagramResponse, len(diagrams)),
		Total:    total,
		Page:     page,
		Limit:    limit,
		HasMore:  int64(offset+limit) < total,
	}

	for i, diagram := range diagrams {
		response.Diagrams[i] = *s.toDiagramResponse(diagram)
	}

	return response, nil
}

// IncrementViewCount increments the view count for a diagram
func (s *DiagramService) IncrementViewCount(id string) error {
	return s.repo.IncrementViewCount(id)
}

// PublishDiagram publishes a diagram
func (s *DiagramService) PublishDiagram(id string, userID string) (*DiagramResponse, error) {
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	diagram.Publish()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to publish diagram: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// ArchiveDiagram archives a diagram
func (s *DiagramService) ArchiveDiagram(id string, userID string) (*DiagramResponse, error) {
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	diagram.Archive()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to archive diagram: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// MakePublic makes a diagram public
func (s *DiagramService) MakePublic(id string, userID string) (*DiagramResponse, error) {
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	diagram.MakePublic()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to make diagram public: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// MakePrivate makes a diagram private
func (s *DiagramService) MakePrivate(id string, userID string) (*DiagramResponse, error) {
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	diagram.MakePrivate()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to make diagram private: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// ShareDiagram shares a diagram
func (s *DiagramService) ShareDiagram(id string, userID string) (*DiagramResponse, error) {
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	diagram.Share()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to share diagram: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// UnshareDiagram unshares a diagram
func (s *DiagramService) UnshareDiagram(id string, userID string) (*DiagramResponse, error) {
	diagram, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("diagram not found")
		}
		return nil, fmt.Errorf("failed to get diagram: %w", err)
	}

	diagram.Unshare()
	diagram.UpdatedBy = userID
	diagram.UpdatedAt = time.Now()

	if err := s.repo.Update(diagram); err != nil {
		return nil, fmt.Errorf("failed to unshare diagram: %w", err)
	}

	return s.toDiagramResponse(diagram), nil
}

// DrawDiagram generates an image from a diagram script
func (s *DiagramService) DrawDiagram(req *DrawDiagramRequest) (*DrawDiagramResponse, error) {

	s.logger.Infof("Drawing diagram: %+v", req)
	// Validate request
	if err := s.validateDrawRequest(req); err != nil {
		s.logger.Errorf("Validation failed: %v", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Handle different script types
	switch req.ScriptType {
	case models.DiagramScriptTypePlantUML:
		return s.drawPlantUMLDiagram(req)
	case models.DiagramScriptTypeMermaid:
		return s.drawMermaidDiagram(req)
	case models.DiagramScriptTypeGraphviz:
		return s.drawGraphvizDiagram(req)
	default:
		s.logger.Errorf("Unsupported script type: %s", req.ScriptType)
		return nil, fmt.Errorf("unsupported script type: %s", req.ScriptType)
	}
}

// drawPlantUMLDiagram generates a PlantUML diagram
func (s *DiagramService) drawPlantUMLDiagram(req *DrawDiagramRequest) (*DrawDiagramResponse, error) {

	// Get PlantUML server URL from environment variable
	plantUMLURL := os.Getenv("PLANTUML_URL")
	if plantUMLURL == "" {
		// Default to public PlantUML server if not configured
		plantUMLURL = "http://www.plantuml.com/plantuml"
	}

	// Create PlantUML client
	plantUMLClient := util.NewPlantUMLClient(plantUMLURL)

	// Generate PlantUML URL
	imageURL, err := plantUMLClient.GeneratePngUrl(req.Script)
	if err != nil {
		s.logger.Errorf("Failed to generate PlantUML URL: %v", err)
		return nil, fmt.Errorf("failed to generate PlantUML URL: %w", err)
	}
	s.logger.Infof("PlantUML image URL: %s", imageURL)
	// Download the image
	resp, err := http.Get(imageURL)
	if err != nil {
		s.logger.Errorf("Failed to download PlantUML image: %v", err)
		return nil, fmt.Errorf("failed to download PlantUML image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("PlantUML server returned status: %d", resp.StatusCode)
		return nil, fmt.Errorf("PlantUML server returned status: %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Failed to read image data: %v", err)
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Data)

	// Create response
	response := &DrawDiagramResponse{
		ImageData: dataURL,
		Format:    models.ImageFormatPNG, // PlantUML generates PNG by default
		Width:     req.Width,
		Height:    req.Height,
		Size:      int64(len(imageData)),
		URL:       imageURL,
	}

	return response, nil
}

// drawMermaidDiagram generates a Mermaid diagram (placeholder)
func (s *DiagramService) drawMermaidDiagram(req *DrawDiagramRequest) (*DrawDiagramResponse, error) {
	// TODO: Implement Mermaid diagram generation
	// This would typically use Mermaid CLI or a Mermaid service

	// For now, return a mock response
	response := &DrawDiagramResponse{
		ImageData: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48dGV4dCB4PSI1MCIgeT0iMTAwIiBmb250LWZhbWlseT0iQXJpYWwiIGZvbnQtc2l6ZT0iMTQiPk1lcm1haWQgRGlhZ3JhbTwvdGV4dD48L3N2Zz4=",
		Format:    req.Format,
		Width:     req.Width,
		Height:    req.Height,
		Size:      1024,
	}

	return response, nil
}

// drawGraphvizDiagram generates a Graphviz diagram (placeholder)
func (s *DiagramService) drawGraphvizDiagram(req *DrawDiagramRequest) (*DrawDiagramResponse, error) {
	// TODO: Implement Graphviz diagram generation
	// This would typically use Graphviz CLI or a Graphviz service

	// For now, return a mock response
	response := &DrawDiagramResponse{
		ImageData: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48dGV4dCB4PSI1MCIgeT0iMTAwIiBmb250LWZhbWlseT0iQXJpYWwiIGZvbnQtc2l6ZT0iMTQiPkdyYXBodml6IERpYWdyYW08L3RleHQ+PC9zdmc+",
		Format:    req.Format,
		Width:     req.Width,
		Height:    req.Height,
		Size:      1024,
	}

	return response, nil
}

// Image Service Methods

// CreateImage creates a new image for a diagram
func (s *DiagramService) CreateImage(req *CreateImageRequest, userID string) (*ImageResponse, error) {
	// Validate request
	if err := s.validateCreateImageRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create image
	image := &models.Image{
		ID:        uuid.New().String(),
		RealmID:   "default", // TODO: Get from context
		UserID:    userID,
		Source:    models.ImageSourceGenerated,
		Type:      models.ImageTypeDiagram,
		Status:    models.ImageStatusProcessed,
		FileName:  req.Path, // Use path as filename for now
		FilePath:  req.Path,
		FileSize:  req.Size,
		Width:     req.Width,
		Height:    req.Height,
		Format:    req.Format,
		Quality:   req.Quality,
		DiagramID: req.DiagramID,
		IsPrimary: req.IsPrimary,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	// Set default quality
	if image.Quality == 0 {
		image.Quality = 90
	}

	// If this is set as primary, unset others
	if image.IsPrimary {
		if err := s.repo.SetPrimaryImage(req.DiagramID, ""); err != nil {
			// Log error but don't fail the creation
			fmt.Printf("Warning: failed to unset primary images: %v\n", err)
		}
	}

	// Create in database
	if err := s.repo.CreateImage(image); err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	return s.toImageResponse(image), nil
}

// GetImage retrieves an image by ID
func (s *DiagramService) GetImage(id string) (*ImageResponse, error) {
	image, err := s.repo.GetImageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return s.toImageResponse(image), nil
}

// GetImagesByDiagramID retrieves all images for a diagram
func (s *DiagramService) GetImagesByDiagramID(diagramID string) ([]ImageResponse, error) {
	images, err := s.repo.GetImagesByDiagramID(diagramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}

	response := make([]ImageResponse, len(images))
	for i, image := range images {
		response[i] = *s.toImageResponse(image)
	}

	return response, nil
}

// UpdateImage updates an existing image
func (s *DiagramService) UpdateImage(id string, req *UpdateImageRequest, userID string) (*ImageResponse, error) {
	// Get existing image
	image, err := s.repo.GetImageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	// Update fields
	if req.Path != "" {
		image.FilePath = req.Path
	}
	if req.Format != "" {
		image.Format = req.Format
	}
	if req.Width > 0 {
		image.Width = req.Width
	}
	if req.Height > 0 {
		image.Height = req.Height
	}
	if req.Size > 0 {
		image.FileSize = req.Size
	}
	if req.Quality > 0 {
		image.Quality = req.Quality
	}

	// Handle primary image setting
	if req.IsPrimary && !image.IsPrimary {
		// Unset other primary images
		if err := s.repo.SetPrimaryImage(image.DiagramID, ""); err != nil {
			// Log error but don't fail the update
			fmt.Printf("Warning: failed to unset primary images: %v\n", err)
		}
		image.IsPrimary = true
	} else if !req.IsPrimary && image.IsPrimary {
		image.IsPrimary = false
	}

	image.UpdatedBy = userID
	image.UpdatedAt = time.Now()

	// Save the image
	if err := s.repo.UpdateImage(image); err != nil {
		return nil, fmt.Errorf("failed to update image: %w", err)
	}

	return s.toImageResponse(image), nil
}

// DeleteImage deletes an image
func (s *DiagramService) DeleteImage(id string) error {
	// Get existing image
	image, err := s.repo.GetImageByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("image not found")
		}
		return fmt.Errorf("failed to get image: %w", err)
	}

	// Delete the image
	if err := s.repo.DeleteImage(id); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	// If this was the primary image, set another one as primary
	if image.IsPrimary {
		images, err := s.repo.GetImagesByDiagramID(image.DiagramID)
		if err == nil && len(images) > 0 {
			// Set the first remaining image as primary
			s.repo.SetPrimaryImage(image.DiagramID, images[0].ID)
		}
	}

	return nil
}

// SetPrimaryImage sets an image as primary
func (s *DiagramService) SetPrimaryImage(diagramID, imageID string) error {
	return s.repo.SetPrimaryImage(diagramID, imageID)
}

// Tag Service Methods

// CreateTag creates a new tag
func (s *DiagramService) CreateTag(req *CreateTagRequest) (*DiagramTagResponse, error) {
	// Validate request
	if err := s.validateCreateTagRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create tag
	tag := &models.DiagramTag{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Slug:        s.generateSlug(req.Name),
		Description: req.Description,
		Color:       req.Color,
	}

	// Create in database
	if err := s.repo.CreateTag(tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return s.toTagResponse(tag), nil
}

// GetTag retrieves a tag by ID
func (s *DiagramService) GetTag(id string) (*DiagramTagResponse, error) {
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return s.toTagResponse(tag), nil
}

// GetAllTags retrieves all tags
func (s *DiagramService) GetAllTags(page, limit int) ([]DiagramTagResponse, error) {
	offset := (page - 1) * limit

	tags, err := s.repo.GetAllTags(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	response := make([]DiagramTagResponse, len(tags))
	for i, tag := range tags {
		response[i] = *s.toTagResponse(tag)
	}

	return response, nil
}

// UpdateTag updates an existing tag
func (s *DiagramService) UpdateTag(id string, req *UpdateTagRequest) (*DiagramTagResponse, error) {
	// Get existing tag
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	// Update fields
	if req.Name != "" {
		tag.Name = req.Name
		tag.Slug = s.generateSlug(req.Name)
	}
	if req.Description != "" {
		tag.Description = req.Description
	}
	if req.Color != "" {
		tag.Color = req.Color
	}

	// Save the tag
	if err := s.repo.UpdateTag(tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return s.toTagResponse(tag), nil
}

// DeleteTag deletes a tag
func (s *DiagramService) DeleteTag(id string) error {
	// Check if tag exists
	_, err := s.repo.GetTagByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("tag not found")
		}
		return fmt.Errorf("failed to get tag: %w", err)
	}

	// Delete tag relations first
	if err := s.repo.DeleteAllTagRelations(""); err != nil {
		// Log error but don't fail the deletion
		fmt.Printf("Warning: failed to delete tag relations: %v\n", err)
	}

	// Delete the tag
	if err := s.repo.DeleteTag(id); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// Helper Methods

// toDiagramResponse converts a Diagram model to a DiagramResponse
func (s *DiagramService) toDiagramResponse(diagram *models.Diagram) *DiagramResponse {
	response := &DiagramResponse{
		ID:          diagram.ID,
		RealmID:     diagram.RealmID,
		Name:        diagram.Name,
		Type:        diagram.Type,
		ScriptType:  diagram.ScriptType,
		Status:      diagram.Status,
		Description: diagram.Description,
		Script:      diagram.Script,
		Theme:       diagram.Theme,
		Tags:        diagram.GetTags(),
		Version:     diagram.Version,
		ViewCount:   diagram.ViewCount,
		EditCount:   diagram.EditCount,
		Public:      diagram.Public,
		Shared:      diagram.Shared,
		ShareToken:  diagram.ShareToken,
		CreatedBy:   diagram.CreatedBy,
		CreatedAt:   diagram.CreatedAt,
		UpdatedBy:   diagram.UpdatedBy,
		UpdatedAt:   diagram.UpdatedAt,
	}

	// Images will be loaded separately if needed

	return response
}

// toImageResponse converts an Image model to an ImageResponse
func (s *DiagramService) toImageResponse(image *models.Image) *ImageResponse {
	return &ImageResponse{
		ID:            image.ID,
		RealmID:       image.RealmID,
		UserID:        image.UserID,
		Source:        string(image.Source),
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

// toTagResponse converts a DiagramTag model to a DiagramTagResponse
func (s *DiagramService) toTagResponse(tag *models.DiagramTag) *DiagramTagResponse {
	return &DiagramTagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		Color:       tag.Color,
		Count:       tag.Count,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}
}

// generateSlug generates a URL-friendly slug from a name
func (s *DiagramService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// Validation Methods

// validateCreateRequest validates a create diagram request
func (s *DiagramService) validateCreateRequest(req *CreateDiagramRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Script == "" {
		return fmt.Errorf("script is required")
	}
	if req.ScriptType == "" {
		return fmt.Errorf("script_type is required")
	}
	return nil
}

// validateUpdateRequest validates an update diagram request
func (s *DiagramService) validateUpdateRequest(req *UpdateDiagramRequest) error {
	// At least one field must be provided for update
	if req.Name == "" && req.Script == "" && req.Description == "" && req.Tags == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}
	return nil
}

// validateDrawRequest validates a draw diagram request
func (s *DiagramService) validateDrawRequest(req *DrawDiagramRequest) error {
	if req.Script == "" {
		return fmt.Errorf("script is required")
	}
	if req.ScriptType == "" {
		return fmt.Errorf("script_type is required")
	}
	if req.Format == "" {
		return fmt.Errorf("format is required")
	}
    /*
	if req.Width < 100 || req.Width > 4000 {
		return fmt.Errorf("width must be between 100 and 4000")
	}
	if req.Height < 100 || req.Height > 4000 {
		return fmt.Errorf("height must be between 100 and 4000")
	} */
	return nil
}

// validateCreateImageRequest validates a create image request
func (s *DiagramService) validateCreateImageRequest(req *CreateImageRequest) error {
	if req.DiagramID == "" {
		return fmt.Errorf("diagram_id is required")
	}
	if req.Path == "" {
		return fmt.Errorf("path is required")
	}
	if req.Format == "" {
		return fmt.Errorf("format is required")
	}
	if req.Width <= 0 {
		return fmt.Errorf("width must be greater than 0")
	}
	if req.Height <= 0 {
		return fmt.Errorf("height must be greater than 0")
	}
	return nil
}

// validateCreateTagRequest validates a create tag request
func (s *DiagramService) validateCreateTagRequest(req *CreateTagRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// Tag Relation Helper Methods

// createTagRelations creates tag relations for a diagram
func (s *DiagramService) createTagRelations(diagramID string, tagNames []string) error {
	for _, tagName := range tagNames {
		// Get or create tag
		tag, err := s.repo.GetTagBySlug(s.generateSlug(tagName))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new tag
				tag = &models.DiagramTag{
					ID:   uuid.New().String(),
					Name: tagName,
					Slug: s.generateSlug(tagName),
				}
				if err := s.repo.CreateTag(tag); err != nil {
					return fmt.Errorf("failed to create tag: %w", err)
				}
			} else {
				return fmt.Errorf("failed to get tag: %w", err)
			}
		}

		// Create relation
		relation := &models.DiagramTagRelation{
			ID:        uuid.New().String(),
			DiagramID: diagramID,
			TagID:     tag.ID,
		}
		if err := s.repo.CreateTagRelation(relation); err != nil {
			return fmt.Errorf("failed to create tag relation: %w", err)
		}
	}
	return nil
}

// updateTagRelations updates tag relations for a diagram
func (s *DiagramService) updateTagRelations(diagramID string, tagNames []string) error {
	// Delete existing relations
	if err := s.repo.DeleteAllTagRelations(diagramID); err != nil {
		return fmt.Errorf("failed to delete existing tag relations: %w", err)
	}

	// Create new relations
	return s.createTagRelations(diagramID, tagNames)
}
