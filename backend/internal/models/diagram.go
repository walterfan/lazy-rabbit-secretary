package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// DiagramScriptType represents the type of diagram script
type DiagramScriptType string

const (
	DiagramScriptTypeMermaid  DiagramScriptType = "mermaid"  // Mermaid diagrams
	DiagramScriptTypePlantUML DiagramScriptType = "plantuml" // PlantUML diagrams
	DiagramScriptTypeGraphviz DiagramScriptType = "graphviz" // Graphviz diagrams
)

// DiagramType represents the type of diagram
type DiagramType string

const (
	DiagramTypeFlowchart    DiagramType = "flowchart"    // Flowchart diagrams
	DiagramTypeSequence     DiagramType = "sequence"     // Sequence diagrams
	DiagramTypeClass        DiagramType = "class"        // Class diagrams
	DiagramTypeER           DiagramType = "mindmap"      // Mindmap diagrams
	DiagramTypeArchitecture DiagramType = "architecture" // Architecture diagrams
	DiagramTypeCustom       DiagramType = "custom"       // Custom diagram types
)

// DiagramStatus represents the status of a diagram
type DiagramStatus string

const (
	DiagramStatusDraft     DiagramStatus = "draft"
	DiagramStatusPublished DiagramStatus = "published"
	DiagramStatusPrivate   DiagramStatus = "private"
	DiagramStatusArchived  DiagramStatus = "archived"
)

// Diagram represents a diagram in the system
type Diagram struct {
	ID          string        `json:"id" gorm:"primaryKey;type:text"`
	RealmID     string        `json:"realm_id" gorm:"not null;type:text;index"`
	Name        string        `json:"name" gorm:"not null;type:text"`
	Type        DiagramType   `json:"type" gorm:"not null;type:text;index"`
	Status      DiagramStatus `json:"status" gorm:"default:'draft';type:text;index"`
	Description string        `json:"description" gorm:"type:text"`

	// Diagram content
	Script     string            `json:"script" gorm:"type:text"`            // The diagram script/source code
	ScriptType DiagramScriptType `json:"script_type" gorm:"type:text;index"` // The diagram type
	ScriptPath string            `json:"script_path" gorm:"type:text"`       // Path to script file

	// Organization
	Tags string `json:"tags" gorm:"type:text"` // comma-separated tags

	// Metadata
	Theme   string `json:"theme" gorm:"default:'default'"` // Diagram theme
	Version int    `json:"version" gorm:"default:1"`       // Version number

	// Statistics
	ViewCount int64 `json:"view_count" gorm:"default:0"` // Number of views
	EditCount int64 `json:"edit_count" gorm:"default:0"` // Number of edits

	// Sharing and collaboration
	Public      bool   `json:"public" gorm:"default:false"`            // Public visibility
	Shared      bool   `json:"shared" gorm:"default:false"`            // Shared with others
	ShareToken  string `json:"share_token,omitempty" gorm:"type:text"` // Share token for public access
	Permissions string `json:"permissions,omitempty" gorm:"type:text"` // JSON permissions

	// Audit fields
	CreatedBy string         `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (Diagram) TableName() string {
	return "diagrams"
}

// IsPublished returns true if the diagram is published
func (d *Diagram) IsPublished() bool {
	return d.Status == DiagramStatusPublished
}

// IsPublic returns true if the diagram is publicly accessible
func (d *Diagram) IsPublic() bool {
	return d.Public
}

// IsShared returns true if the diagram is shared
func (d *Diagram) IsShared() bool {
	return d.Shared
}

// IsMermaid returns true if this is a Mermaid diagram
func (d *Diagram) IsMermaid() bool {
	return d.ScriptType == DiagramScriptTypeMermaid
}

// IsPlantUML returns true if this is a PlantUML diagram
func (d *Diagram) IsPlantUML() bool {
	return d.ScriptType == DiagramScriptTypePlantUML
}

// IsGraphviz returns true if this is a Graphviz diagram
func (d *Diagram) IsGraphviz() bool {
	return d.ScriptType == DiagramScriptTypeGraphviz
}

// GetTags returns the tags as a slice
func (d *Diagram) GetTags() []string {
	if d.Tags == "" {
		return []string{}
	}
	tags := strings.Split(d.Tags, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// SetTags sets tags from a slice
func (d *Diagram) SetTags(tags []string) {
	if len(tags) == 0 {
		d.Tags = ""
		return
	}
	// Clean and join tags
	cleaned := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	d.Tags = strings.Join(cleaned, ",")
}

// IncrementViewCount increments the view count
func (d *Diagram) IncrementViewCount() {
	d.ViewCount++
}

// IncrementEditCount increments the edit count and version
func (d *Diagram) IncrementEditCount() {
	d.EditCount++
	d.Version++
}

// Publish sets the diagram status to published
func (d *Diagram) Publish() {
	d.Status = DiagramStatusPublished
}

// Unpublish sets the diagram back to draft status
func (d *Diagram) Unpublish() {
	d.Status = DiagramStatusDraft
}

// Archive archives the diagram
func (d *Diagram) Archive() {
	d.Status = DiagramStatusArchived
}

// MakePublic makes the diagram publicly accessible
func (d *Diagram) MakePublic() {
	d.Public = true
}

// MakePrivate makes the diagram private
func (d *Diagram) MakePrivate() {
	d.Public = false
	d.Shared = false
	d.ShareToken = ""
}

// Share enables sharing for the diagram
func (d *Diagram) Share() {
	d.Shared = true
}

// Unshare disables sharing for the diagram
func (d *Diagram) Unshare() {
	d.Shared = false
	d.ShareToken = ""
}

// CanBeEdited returns true if the diagram can be edited
func (d *Diagram) CanBeEdited() bool {
	return d.Status != DiagramStatusArchived
}

// CanBeDeleted returns true if the diagram can be deleted
func (d *Diagram) CanBeDeleted() bool {
	return d.Status != DiagramStatusArchived
}

// HasScript returns true if the diagram has a script
func (d *Diagram) HasScript() bool {
	return d.Script != ""
}

// GetScriptURL returns the URL for the diagram script
func (d *Diagram) GetScriptURL() string {
	if d.ScriptPath == "" {
		return ""
	}
	// This would typically be constructed based on your file serving setup
	return "/api/v1/diagrams/" + d.ID + "/script"
}

// DiagramTag represents tags for diagrams
type DiagramTag struct {
	ID          string `json:"id" gorm:"primaryKey;type:text"`
	Name        string `json:"name" gorm:"not null;type:text"`
	Slug        string `json:"slug" gorm:"not null;type:text;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	Color       string `json:"color" gorm:"type:text"` // Hex color for UI
	Count       int64  `json:"count" gorm:"default:0"` // Number of diagrams with this tag

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (DiagramTag) TableName() string {
	return "diagram_tags"
}

// DiagramTagRelation represents the many-to-many relationship between diagrams and tags
type DiagramTagRelation struct {
	ID        string `json:"id" gorm:"primaryKey;type:text"`
	DiagramID string `json:"diagram_id" gorm:"not null;type:text;index"`
	TagID     string `json:"tag_id" gorm:"not null;type:text;index"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key relationships
	Diagram Diagram    `json:"diagram" gorm:"foreignKey:DiagramID;references:ID"`
	Tag     DiagramTag `json:"tag" gorm:"foreignKey:TagID;references:ID"`
}

// TableName specifies the table name for GORM
func (DiagramTagRelation) TableName() string {
	return "diagram_tag_relations"
}

// DiagramVersion represents version history for diagrams
type DiagramVersion struct {
	ID        string `json:"id" gorm:"primaryKey;type:text"`
	DiagramID string `json:"diagram_id" gorm:"not null;type:text;index"`
	Version   int    `json:"version" gorm:"not null;index"`

	// Version content
	Script     string `json:"script" gorm:"type:text"`
	ScriptPath string `json:"script_path" gorm:"type:text"`

	// Version metadata
	ChangeNote string `json:"change_note" gorm:"type:text"` // Description of changes
	Size       int    `json:"size" gorm:"default:0"`        // Size in bytes

	// Audit fields
	CreatedBy string    `json:"created_by" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Foreign key relationship
	Diagram Diagram `json:"diagram" gorm:"foreignKey:DiagramID;references:ID"`
}

// TableName specifies the table name for GORM
func (DiagramVersion) TableName() string {
	return "diagram_versions"
}
