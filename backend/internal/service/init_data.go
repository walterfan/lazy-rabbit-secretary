package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DEFAULT_REALM_ID = "65e52554-40f7-4afa-b41f-e897b94d31b2"
	ADMIN_USER_ID    = "7ec73634-2569-48fe-a2e5-d124b7d922ae"
	ADMIN_ROLE_ID    = "6c12df67-5e77-44be-b6ca-3ccc1cc61c00"
	USER_ROLE_ID     = "8cb242bb-0f71-443d-ad78-b51dbdafc4ab"
	ADMIN_POLICY_ID  = "1012b695-9209-42ae-bd45-8f9dad559466"
	USER_POLICY_ID   = "5a63416b-a889-4906-a5c8-c02855c14d14"
)

// InitCompleteData initializes database with comprehensive default data
func InitCompleteData(db *gorm.DB) error {
	if err := initRealms(db); err != nil {
		return fmt.Errorf("failed to initialize realms: %w", err)
	}

	if err := initRoles(db); err != nil {
		return fmt.Errorf("failed to initialize roles: %w", err)
	}

	if err := initPolicies(db); err != nil {
		return fmt.Errorf("failed to initialize policies: %w", err)
	}

	if err := initStatements(db); err != nil {
		return fmt.Errorf("failed to initialize statements: %w", err)
	}

	if err := initUsers(db); err != nil {
		return fmt.Errorf("failed to initialize users: %w", err)
	}

	if err := initUserRoles(db); err != nil {
		return fmt.Errorf("failed to initialize user roles: %w", err)
	}

	if err := initRolePolicies(db); err != nil {
		return fmt.Errorf("failed to initialize role policies: %w", err)
	}

	log.Println("Database initialization completed successfully")
	return nil
}

// initRealms creates default realm
func initRealms(db *gorm.DB) error {
	var count int64
	db.Model(&models.Realm{}).Count(&count)
	if count > 0 {
		log.Println("Realms already exist, skipping realm initialization")
		return nil
	}

	realms := []models.Realm{
		{
			ID:          DEFAULT_REALM_ID,
			Name:        "default",
			Description: "Default organizational realm for the application",
			CreatedBy:   "system",
		},
	}

	result := db.Create(&realms)
	if result.Error != nil {
		return fmt.Errorf("failed to create realms: %w", result.Error)
	}

	log.Printf("Created %d default realms", len(realms))
	return nil
}

// initRoles creates default roles
func initRoles(db *gorm.DB) error {
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		log.Println("Roles already exist, skipping role initialization")
		return nil
	}

	roles := []models.Role{
		{
			ID:          ADMIN_ROLE_ID,
			RealmID:     DEFAULT_REALM_ID,
			Name:        "admin",
			Description: "Administrator role with full system access",
			CreatedBy:   "system",
		},
		{
			ID:          USER_ROLE_ID,
			RealmID:     DEFAULT_REALM_ID,
			Name:        "user",
			Description: "Standard user role with basic access",
			CreatedBy:   "system",
		},
	}

	result := db.Create(&roles)
	if result.Error != nil {
		return fmt.Errorf("failed to create roles: %w", result.Error)
	}

	log.Printf("Created %d default roles", len(roles))
	return nil
}

// initPolicies creates default policies
func initPolicies(db *gorm.DB) error {
	var count int64
	db.Model(&models.Policy{}).Count(&count)
	if count > 0 {
		log.Println("Policies already exist, skipping policy initialization")
		return nil
	}

	policies := []models.Policy{
		{
			ID:          ADMIN_POLICY_ID,
			RealmID:     DEFAULT_REALM_ID,
			Name:        "AdminFullAccess",
			Description: "Full administrative access to all resources",
			Version:     "v1.0",
			CreatedBy:   "system",
		},
		{
			ID:          USER_POLICY_ID,
			RealmID:     DEFAULT_REALM_ID,
			Name:        "UserBasicAccess",
			Description: "Basic user access to personal resources",
			Version:     "v1.0",
			CreatedBy:   "system",
		},
	}

	result := db.Create(&policies)
	if result.Error != nil {
		return fmt.Errorf("failed to create policies: %w", result.Error)
	}

	log.Printf("Created %d default policies", len(policies))
	return nil
}

// initStatements creates default policy statements
func initStatements(db *gorm.DB) error {
	var count int64
	db.Model(&models.Statement{}).Count(&count)
	if count > 0 {
		log.Println("Statements already exist, skipping statement initialization")
		return nil
	}

	// Admin statements - full access
	adminActions := []string{"*"}
	adminResources := []string{"*"}
	adminActionsJSON, _ := json.Marshal(adminActions)
	adminResourcesJSON, _ := json.Marshal(adminResources)

	// User statements - limited access
	userActions := []string{"read", "create", "update"}
	userResources := []string{"tasks", "books", "profile"}
	userActionsJSON, _ := json.Marshal(userActions)
	userResourcesJSON, _ := json.Marshal(userResources)

	statements := []models.Statement{
		{
			ID:         uuid.New().String(),
			PolicyID:   ADMIN_POLICY_ID,
			SID:        "AdminFullAccessStatement",
			Effect:     "Allow",
			Actions:    string(adminActionsJSON),
			Resources:  string(adminResourcesJSON),
			Conditions: "",
			CreatedBy:  "system",
		},
		{
			ID:         uuid.New().String(),
			PolicyID:   USER_POLICY_ID,
			SID:        "UserBasicAccessStatement",
			Effect:     "Allow",
			Actions:    string(userActionsJSON),
			Resources:  string(userResourcesJSON),
			Conditions: "",
			CreatedBy:  "system",
		},
		{
			ID:         uuid.New().String(),
			PolicyID:   USER_POLICY_ID,
			SID:        "UserDenyAdminResources",
			Effect:     "Deny",
			Actions:    `["*"]`,
			Resources:  `["admin", "users", "roles", "policies"]`,
			Conditions: "",
			CreatedBy:  "system",
		},
	}

	result := db.Create(&statements)
	if result.Error != nil {
		return fmt.Errorf("failed to create statements: %w", result.Error)
	}

	log.Printf("Created %d default statements", len(statements))
	return nil
}

// initUsers creates default users
func initUsers(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already exist, skipping user initialization")
		return nil
	}

	// Get admin credentials from environment
	adminUsername := getEnvOrDefault("ADMIN_USERNAME", "admin")
	adminPassword := getEnvOrDefault("ADMIN_PASSWORD", "admin123")
	adminEmail := getEnvOrDefault("ADMIN_EMAIL", "admin@fanyamin.com")

	// Hash passwords
	adminPwdHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	testPwdHash, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash user password: %w", err)
	}

	users := []models.User{
		{
			ID:             ADMIN_USER_ID,
			RealmID:        DEFAULT_REALM_ID,
			Username:       adminUsername,
			Email:          adminEmail,
			HashedPassword: string(adminPwdHash),
			IsActive:       true,
			CreatedBy:      "system",
		},
		{
			ID:             uuid.New().String(),
			RealmID:        DEFAULT_REALM_ID,
			Username:       "testuser",
			Email:          "test@example.com",
			HashedPassword: string(testPwdHash),
			IsActive:       true,
			CreatedBy:      "system",
		},
	}

	result := db.Create(&users)
	if result.Error != nil {
		return fmt.Errorf("failed to create users: %w", result.Error)
	}

	log.Printf("Created %d default users", len(users))
	return nil
}

// initUserRoles assigns roles to users
func initUserRoles(db *gorm.DB) error {
	var count int64
	db.Model(&models.UserRole{}).Count(&count)
	if count > 0 {
		log.Println("User roles already exist, skipping user role initialization")
		return nil
	}

	// Get the test user ID
	var testUser models.User
	result := db.Where("username = ? AND realm_id = ?", "testuser", DEFAULT_REALM_ID).First(&testUser)
	if result.Error != nil {
		return fmt.Errorf("failed to find test user: %w", result.Error)
	}

	userRoles := []models.UserRole{
		{
			UserID: ADMIN_USER_ID,
			RoleID: ADMIN_ROLE_ID,
		},
		{
			UserID: testUser.ID,
			RoleID: USER_ROLE_ID,
		},
	}

	result = db.Create(&userRoles)
	if result.Error != nil {
		return fmt.Errorf("failed to create user roles: %w", result.Error)
	}

	log.Printf("Created %d user role assignments", len(userRoles))
	return nil
}

// initRolePolicies assigns policies to roles
func initRolePolicies(db *gorm.DB) error {
	var count int64
	db.Model(&models.RolePolicy{}).Count(&count)
	if count > 0 {
		log.Println("Role policies already exist, skipping role policy initialization")
		return nil
	}

	rolePolicies := []models.RolePolicy{
		{
			RoleID:   ADMIN_ROLE_ID,
			PolicyID: ADMIN_POLICY_ID,
		},
		{
			RoleID:   USER_ROLE_ID,
			PolicyID: USER_POLICY_ID,
		},
	}

	result := db.Create(&rolePolicies)
	if result.Error != nil {
		return fmt.Errorf("failed to create role policies: %w", result.Error)
	}

	log.Printf("Created %d role policy assignments", len(rolePolicies))
	return nil
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
