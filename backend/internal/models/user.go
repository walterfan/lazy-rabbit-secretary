package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// UserStatus represents the status of a user registration
type UserStatus string

const (
	UserStatusPending   UserStatus = "pending"   // Registration pending email confirmation
	UserStatusConfirmed UserStatus = "confirmed" // Email confirmed, pending approval
	UserStatusApproved  UserStatus = "approved"  // Registration approved, user active
	UserStatusDenied    UserStatus = "denied"    // Registration denied
	UserStatusSuspended UserStatus = "suspended" // User suspended
)

// User represents the app_user table in the database
type User struct {
	ID                     string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID                string         `json:"realm_id" gorm:"not null;type:text;index"`
	Username               string         `json:"username" gorm:"not null;type:text;index"`
	Email                  string         `json:"email" gorm:"unique;type:text;index"`
	HashedPassword         string         `json:"-" gorm:"not null;column:hashed_password;type:text"` // Don't expose password in JSON
	IsActive               bool           `json:"is_active" gorm:"default:false"`
	Status                 UserStatus     `json:"status" gorm:"default:'pending';type:text;index"`
	EmailConfirmationToken string         `json:"-" gorm:"type:text;index"` // Don't expose token in JSON
	EmailConfirmedAt       *time.Time     `json:"email_confirmed_at"`
	ConfirmationExpiresAt  *time.Time     `json:"-" gorm:"index"` // Don't expose expiry in JSON
	CreatedBy              string         `json:"created_by" gorm:"type:text"`
	CreatedAt              time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy              string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt              time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
	Roles                  []Role         `json:"roles,omitempty" gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID"` // User roles
}

func (User) TableName() string {
	return "app_user"
}

// GenerateEmailConfirmationToken generates a secure random token for email confirmation
func (u *User) GenerateEmailConfirmationToken() error {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}

	// Convert to hex string
	u.EmailConfirmationToken = hex.EncodeToString(bytes)

	// Set expiration to 24 hours from now
	expiresAt := time.Now().Add(24 * time.Hour)
	u.ConfirmationExpiresAt = &expiresAt

	return nil
}

// IsEmailConfirmationValid checks if the email confirmation token is valid and not expired
func (u *User) IsEmailConfirmationValid(token string) bool {
	// Check if token matches
	if u.EmailConfirmationToken != token {
		return false
	}

	// Check if token has expired
	if u.ConfirmationExpiresAt == nil || time.Now().After(*u.ConfirmationExpiresAt) {
		return false
	}

	return true
}

// ConfirmEmail marks the user's email as confirmed
func (u *User) ConfirmEmail() {
	now := time.Now()
	u.EmailConfirmedAt = &now
	u.Status = UserStatusConfirmed
	u.EmailConfirmationToken = "" // Clear the token
	u.ConfirmationExpiresAt = nil // Clear the expiration
}

// IsEmailConfirmed checks if the user's email has been confirmed
func (u *User) IsEmailConfirmed() bool {
	return u.EmailConfirmedAt != nil && u.Status == UserStatusConfirmed
}

// CanBeApproved checks if the user can be approved (must be confirmed first)
func (u *User) CanBeApproved() bool {
	return u.Status == UserStatusConfirmed && u.IsEmailConfirmed()
}
