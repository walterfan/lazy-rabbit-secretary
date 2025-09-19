package models

// models.go - Overview of all models in the system
//
// This package implements an AWS-style IAM (Identity and Access Management) system
// with support for realms, users, roles, policies, and resource management.
//
// =====================
// 🔐 Authentication & Authorization Models
// =====================
//
// Realm - Top-level organizational unit (tenant)
// User - Application users (stored as app_user in DB to avoid conflicts)
// Role - Named collections of permissions
// Policy - Documents that define permissions using AWS IAM format
// Statement - Individual permission statements within policies
//
// Join Tables:
// UserRole - Many-to-many relationship between users and roles
// UserPolicy - Direct user-to-policy assignments
// RolePolicy - Role-to-policy assignments
// ResourcePolicy - Resource-based policy attachments
//
// =====================
// 📁 Project & Content Models
// =====================
//
// Project - Code projects with Git integration
// Code - Source code files with vector embeddings for search
// Document - Documentation files with vector embeddings
// Prompt - AI prompt templates (existing functionality)
//
// =====================
// 📝 Blog & CMS Models (WordPress-style)
// =====================
//
// Post - Blog posts and pages with WordPress-like features
// PostMeta - Additional metadata for posts (custom fields)
// PostCategory - Hierarchical categories for content organization
// PostTag - Flat tags for content labeling
// PostCategoryRelation - Many-to-many posts ↔ categories
// PostTagRelation - Many-to-many posts ↔ tags
// Comment - Threaded comments on posts
//
// =====================
// 🔗 Key Relationships
// =====================
//
// Realm (1) → (many) Users, Roles, Policies, Projects, Posts
// User (many) ← → (many) Roles (via UserRole)
// User (many) ← → (many) Policies (via UserPolicy)
// Role (many) ← → (many) Policies (via RolePolicy)
// Policy (1) → (many) Statements
// Project (1) → (many) Code, Documents
// Post (1) → (many) PostMeta, Comments
// Post (many) ← → (many) PostCategory (via PostCategoryRelation)
// Post (many) ← → (many) PostTag (via PostTagRelation)
// PostCategory (1) → (many) PostCategory (hierarchical)
// Comment (1) → (many) Comment (threaded replies)
//
// =====================
// 🎯 Permission Evaluation (AWS-style)
// =====================
//
// 1. Explicit DENY always wins
// 2. Explicit ALLOW is required (default is implicit deny)
// 3. Evaluation order: Resource-based policies → User policies → Role policies
// 4. Conditions are evaluated using variable substitution
//
// =====================
// 📝 Database Compatibility
// =====================
//
// All models use string IDs (TEXT in SQLite) for maximum compatibility
// JSON fields are stored as strings and parsed at application level
// Timestamps use Go's time.Time with proper GORM auto-management
// Soft deletes supported via gorm.DeletedAt

// GetAllModels returns a slice of all model types for GORM AutoMigrate
func GetAllModels() []interface{} {
	return []interface{}{
		// Authentication & Authorization
		&Realm{},
		&User{},
		&Role{},
		&Policy{},
		&Statement{},
		&UserRole{},
		&RolePolicy{},
		&UserPolicy{},
		&ResourcePolicy{},

		// Project & Content Management
		&Project{},
		&Book{},
		&BookTag{},
		&Code{},
		&Document{},
		&Prompt{},
		&Secret{},

		// Task & Reminder System
		&Task{},
		&Reminder{},
		&TaskReminder{},

		// GTD System
		&InboxItem{},
		&DailyChecklistItem{},

		// Blog & CMS (WordPress-style)
		&Post{},
		&PostMeta{},
		&PostCategory{},
		&PostTag{},
		&PostCategoryRelation{},
		&PostTagRelation{},
		&Comment{},
	}
}

// DatabaseModels represents the complete model structure for reference
type DatabaseModels struct {
	// Authentication & Authorization
	Realms           []Realm
	Users            []User
	Roles            []Role
	Policies         []Policy
	Statements       []Statement
	UserRoles        []UserRole
	RolePolicies     []RolePolicy
	UserPolicies     []UserPolicy
	ResourcePolicies []ResourcePolicy

	// Project & Content Management
	Projects  []Project
	Books     []Book
	BookTags  []BookTag
	Codes     []Code
	Documents []Document
	Prompts   []Prompt
	Secrets   []Secret

	// Task & Reminder System
	Tasks         []Task
	Reminders     []Reminder
	TaskReminders []TaskReminder

	// Blog & CMS (WordPress-style)
	Posts                 []Post
	PostMetas             []PostMeta
	PostCategories        []PostCategory
	PostTags              []PostTag
	PostCategoryRelations []PostCategoryRelation
	PostTagRelations      []PostTagRelation
	Comments              []Comment
}
