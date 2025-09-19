package support

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestConfig holds configuration for BDD tests
type TestConfig struct {
	DatabaseURL    string
	JWTPrivateKey  string
	JWTPublicKey   string
	TestDataPath   string
	CleanupEnabled bool
}

// DefaultTestConfig returns a default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL:    ":memory:", // SQLite in-memory database for tests
		TestDataPath:   "./testdata",
		CleanupEnabled: true,
	}
}

// SetupTestDatabase initializes a test database with required tables
func SetupTestDatabase(config *TestConfig) (*gorm.DB, error) {
	// Create a direct SQLite in-memory database connection for testing
	// This bypasses the database.InitDB() which calls InitData() and tries to create default users

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(models.GetAllModels()...); err != nil {
		return nil, fmt.Errorf("failed to migrate test database: %w", err)
	}

	// Create default realm for testing
	defaultRealm := &models.Realm{
		ID:          "default-realm-id",
		Name:        "default",
		Description: "Default test realm",
		CreatedBy:   "test",
		UpdatedBy:   "test",
	}

	if err := db.Create(defaultRealm).Error; err != nil {
		return nil, fmt.Errorf("failed to create default realm: %w", err)
	}

	return db, nil
}

// CleanupTestDatabase removes test data from the database
func CleanupTestDatabase(db *gorm.DB) error {
	// Since we're using fresh in-memory databases for each scenario,
	// minimal cleanup is needed. Just remove any test users created during scenarios.

	// Clean test scenario users (keep the default realm)
	if err := db.Exec("DELETE FROM app_user WHERE email LIKE '%@example.com'").Error; err != nil {
		fmt.Printf("Warning: failed to clean test users: %v\n", err)
	}

	return nil
}

// CreateTestJWTKeys generates RSA key pair for testing
func CreateTestJWTKeys() (privateKeyPEM, publicKeyPEM []byte, err error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Encode private key to PEM (PKCS1 format for compatibility)
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}

// CreateTestJWTManager creates a JWT manager with test keys
func CreateTestJWTManager() (*auth.JWTManager, error) {
	privateKeyPEM, publicKeyPEM, err := CreateTestJWTKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to create test keys: %w", err)
	}

	// Create temporary key files
	tempDir := os.TempDir()
	privateKeyPath := filepath.Join(tempDir, "test_jwt_private.pem")
	publicKeyPath := filepath.Join(tempDir, "test_jwt_public.pem")

	// Write keys to temporary files
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return nil, fmt.Errorf("failed to write private key: %w", err)
	}

	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return nil, fmt.Errorf("failed to write public key: %w", err)
	}

	// Create JWT manager
	jwtManager, err := auth.NewJWTManager(privateKeyPath, publicKeyPath, "test-issuer", "test-audience")
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT manager: %w", err)
	}

	// Clean up temporary files
	defer func() {
		os.Remove(privateKeyPath)
		os.Remove(publicKeyPath)
	}()

	return jwtManager, nil
}

// CreateTestUser creates a test user in the database
func CreateTestUser(db *gorm.DB, email, username string, status models.UserStatus) (*models.User, error) {
	user := &models.User{
		ID:             fmt.Sprintf("test-user-%s", username),
		RealmID:        "default-realm-id",
		Username:       username,
		Email:          email,
		HashedPassword: "$2a$10$test.hash.for.test.user.password",
		IsActive:       status == models.UserStatusApproved,
		Status:         status,
		CreatedBy:      "test",
		UpdatedBy:      "test",
	}

	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create test user: %w", err)
	}

	return user, nil
}

// AssertUserExists checks if a user exists with the given email and status
func AssertUserExists(db *gorm.DB, email string, expectedStatus models.UserStatus) error {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with email %s not found", email)
		}
		return fmt.Errorf("failed to query user: %w", result.Error)
	}

	if user.Status != expectedStatus {
		return fmt.Errorf("expected user status %s, got %s", expectedStatus, user.Status)
	}

	return nil
}

// AssertUserNotExists checks if a user does not exist with the given email
func AssertUserNotExists(db *gorm.DB, email string) error {
	var count int64
	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return fmt.Errorf("user with email %s should not exist, but found %d records", email, count)
	}
	return nil
}

// TestUserService implements UserService interface for testing with a specific database connection
type TestUserService struct {
	db *gorm.DB
}

// NewTestUserService creates a new test user service with the given database connection
func NewTestUserService(db *gorm.DB) *TestUserService {
	return &TestUserService{db: db}
}

func (s *TestUserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}
	return &user, nil
}

func (s *TestUserService) GetUserByUsername(username string, realmID string) (*models.User, error) {
	var user models.User
	result := s.db.Where("username = ? AND realm_id = ?", username, realmID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", result.Error)
	}
	return &user, nil
}

func (s *TestUserService) GetUserRoles(userID string) ([]*models.Role, error) {
	var userRoles []models.UserRole
	var roles []*models.Role

	result := s.db.Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", result.Error)
	}

	for _, userRole := range userRoles {
		var role models.Role
		if err := s.db.Where("id = ?", userRole.RoleID).First(&role).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to get role details: %w", err)
			}
			continue
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (s *TestUserService) GetRealmByName(realmName string) (*models.Realm, error) {
	var realm models.Realm
	result := s.db.Where("name = ?", realmName).First(&realm)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, auth.ErrInvalidRealm
		}
		return nil, fmt.Errorf("failed to get realm by name: %w", result.Error)
	}
	return &realm, nil
}

func (s *TestUserService) CreateUser(user *models.User) error {
	result := s.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (s *TestUserService) UpdateUser(user *models.User) error {
	result := s.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

func (s *TestUserService) DeleteUser(userID string) error {
	result := s.db.Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	return nil
}

func (s *TestUserService) GetUserRegistrations(req models.UserRegistrationRequest) (*models.UserRegistrationResponse, error) {
	// Simplified implementation for testing
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &models.UserRegistrationResponse{
		Users: users,
		Total: total,
	}, nil
}

func (s *TestUserService) GetUserRegistrationStats(realmName string) (map[string]int64, error) {
	stats := make(map[string]int64)
	var total int64
	s.db.Model(&models.User{}).Count(&total)
	stats["total"] = total
	return stats, nil
}

func (s *TestUserService) GetUserByConfirmationToken(token string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email_confirmation_token = ?", token).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid confirmation token")
		}
		return nil, fmt.Errorf("failed to get user by confirmation token: %w", err)
	}
	return &user, nil
}

func (s *TestUserService) ConfirmUserEmail(userID string) error {
	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email_confirmed_at": time.Now(),
		"status":             models.UserStatusConfirmed,
	})
	return result.Error
}
