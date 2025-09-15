package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/database"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/email"
	"gorm.io/gorm"
)

// AuthService handles authentication and user management
type AuthService struct {
	userService      UserService
	passwordManager  *PasswordManager
	jwtManager       *JWTManager
	permissionEngine *PermissionEngine
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userService UserService,
	passwordManager *PasswordManager,
	jwtManager *JWTManager,
	permissionEngine *PermissionEngine,
) *AuthService {
	return &AuthService{
		userService:      userService,
		passwordManager:  passwordManager,
		jwtManager:       jwtManager,
		permissionEngine: permissionEngine,
	}
}

// Login authenticates a user and returns JWT tokens
func (a *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// Get realm by name
	realm, err := a.userService.GetRealmByName(req.RealmName)
	if err != nil {
		return nil, ErrInvalidRealm
	}
	realmID := realm.ID

	// Get user by username and realm
	user, err := a.userService.GetUserByUsername(req.Username, realmID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if !a.passwordManager.VerifyPassword(req.Password, user.HashedPassword) {
		return nil, ErrInvalidCredentials
	}

	// Get user roles
	roles, err := a.userService.GetUserRoles(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate tokens
	accessToken, err := a.jwtManager.GenerateToken(
		user.ID,
		user.RealmID,
		user.Username,
		user.Email,
		roleNames,
		15*time.Minute, // Access token expires in 15 minutes
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := a.jwtManager.GenerateToken(
		user.ID,
		user.RealmID,
		user.Username,
		user.Email,
		roleNames,
		7*24*time.Hour, // Refresh token expires in 7 days
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    60 * 60, // 60 minutes in seconds
		User:         *user,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (a *AuthService) RefreshToken(req models.RefreshRequest) (*models.LoginResponse, error) {
	// Validate refresh token
	claims, err := a.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := a.userService.GetUserByID(claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Get user roles
	roles, err := a.userService.GetUserRoles(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// Generate new access token
	accessToken, err := a.jwtManager.GenerateToken(
		user.ID,
		user.RealmID,
		user.Username,
		user.Email,
		roleNames,
		15*time.Minute,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // Keep the same refresh token
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60,
		User:         *user,
	}, nil
}

// RegisterUser creates a new user account
func (a *AuthService) RegisterUser(req models.CreateUserRequest, createdBy string) (*models.User, error) {
	// Parse realm ID
	realm, err := a.userService.GetRealmByName(req.RealmName)
	if err != nil {
		return nil, ErrInvalidRealm
	}

	// Check password strength
	if err := a.passwordManager.CheckPasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if username already exists in the realm
	existingUser, _ := a.userService.GetUserByUsername(req.Username, realm.ID)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := a.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:             uuid.New().String(),
		RealmID:        realm.ID,
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		IsActive:       false,
		Status:         models.UserStatusPending,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
		UpdatedBy:      createdBy,
		UpdatedAt:      time.Now(),
	}

	// Generate email confirmation token
	if err := user.GenerateEmailConfirmationToken(); err != nil {
		return nil, fmt.Errorf("failed to generate confirmation token: %w", err)
	}

	if err := a.userService.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send email confirmation to user
	go func() {
		if err := a.sendEmailConfirmation(user); err != nil {
			log.Printf("Failed to send email confirmation: %v", err)
		}
	}()

	// Send notification email to admin about new registration
	go func() {
		if err := a.sendNewRegistrationNotification(user); err != nil {
			log.Printf("Failed to send new registration notification: %v", err)
		}
	}()

	return user, nil
}

// CheckPermission checks if a user has permission to perform an action
func (a *AuthService) CheckPermission(userID, realmID, action, resource string, context map[string]interface{}) (bool, error) {
	check := models.PermissionCheck{
		Action:   action,
		Resource: resource,
		UserID:   userID,
		RealmID:  realmID,
		Context:  context,
	}

	return a.permissionEngine.CheckPermission(check)
}

// ValidateToken validates a JWT token and returns the claims
func (a *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	return a.jwtManager.ValidateToken(tokenString)
}

// GetPendingRegistrations retrieves users with pending registration status
func (a *AuthService) GetPendingRegistrations(req models.UserRegistrationRequest) (*models.UserRegistrationResponse, error) {
	return a.userService.GetUserRegistrations(req)
}

// ApproveRegistration approves or denies a user registration
func (a *AuthService) ApproveRegistration(req models.ApproveRegistrationRequest, approvedBy string) (*models.ApproveRegistrationResponse, error) {
	// Get the user
	user, err := a.userService.GetUserByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user can be approved (must be confirmed first)
	if !user.CanBeApproved() {
		if user.Status == models.UserStatusPending {
			return nil, fmt.Errorf("user email must be confirmed before approval")
		}
		return nil, fmt.Errorf("user cannot be approved, current status: %s", user.Status)
	}

	// Update user status
	if req.Approved {
		user.Status = models.UserStatusApproved
		user.IsActive = true
	} else {
		user.Status = models.UserStatusDenied
		user.IsActive = false
	}

	user.UpdatedBy = approvedBy
	user.UpdatedAt = time.Now()

	// Save the updated user
	if err := a.userService.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}

	// Send email notification
	go func() {
		if err := a.sendRegistrationStatusEmail(user, req.Approved, req.Reason); err != nil {
			log.Printf("Failed to send registration status email to %s: %v", user.Email, err)
		}
	}()

	var message string
	if req.Approved {
		message = "User registration approved successfully"
	} else {
		message = "User registration denied"
		if req.Reason != "" {
			message += ": " + req.Reason
		}
	}

	return &models.ApproveRegistrationResponse{
		Message: message,
		User:    *user,
	}, nil
}

// GetUserRegistrationStats returns statistics about user registrations
func (a *AuthService) GetUserRegistrationStats(realmName string) (map[string]int64, error) {
	return a.userService.GetUserRegistrationStats(realmName)
}

// sendRegistrationStatusEmail sends an email notification when registration status changes
func (a *AuthService) sendRegistrationStatusEmail(user *models.User, approved bool, reason string) error {
	sender, err := email.NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	// Get template manager
	templateManager := email.GetGlobalTemplateManager()
	if templateManager == nil {
		return fmt.Errorf("email template manager not initialized")
	}

	// Prepare template data
	templateData := map[string]interface{}{
		"Username": user.Username,
		"Email":    user.Email,
		"Approved": approved,
		"Status":   map[bool]string{true: "Approved", false: "Denied"}[approved],
		"Reason":   reason,
		"ToAddr":   []string{user.Email},
	}

	// Render email template
	message, err := templateManager.RenderTemplate("registration_approval", templateData)
	if err != nil {
		return fmt.Errorf("failed to render registration approval template: %w", err)
	}

	return sender.SendEmail(message)
}

// sendNewRegistrationNotification sends an email to admin when a new user registers
func (a *AuthService) sendNewRegistrationNotification(user *models.User) error {
	sender, err := email.NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	// Get template manager
	templateManager := email.GetGlobalTemplateManager()
	if templateManager == nil {
		return fmt.Errorf("email template manager not initialized")
	}

	// Get app config for admin panel URL
	appConfig := templateManager.GetAppConfig()

	// Prepare template data
	templateData := map[string]interface{}{
		"Username":         user.Username,
		"Email":            user.Email,
		"RealmID":          user.RealmID,
		"RegistrationTime": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"Status":           string(user.Status),
		"AdminPanelURL":    appConfig.AdminPanelURL,
	}

	// Render email template
	message, err := templateManager.RenderTemplate("new_registration_notification", templateData)
	if err != nil {
		return fmt.Errorf("failed to render new registration notification template: %w", err)
	}

	// Send to default admin email (configured in MAIL_RECEIVER)
	return sender.SendEmail(message)
}

// sendEmailConfirmation sends a confirmation email to the user
func (a *AuthService) sendEmailConfirmation(user *models.User) error {
	sender, err := email.NewEmailSender()
	if err != nil {
		return fmt.Errorf("failed to create email sender: %w", err)
	}

	// Get template manager
	templateManager := email.GetGlobalTemplateManager()
	if templateManager == nil {
		return fmt.Errorf("email template manager not initialized")
	}

	// Get app config for base URL
	appConfig := templateManager.GetAppConfig()
	confirmURL := fmt.Sprintf("%s/api/v1/auth/confirm?token=%s", appConfig.BaseURL, user.EmailConfirmationToken)

	// Prepare template data
	templateData := map[string]interface{}{
		"Username":   user.Username,
		"Email":      user.Email,
		"ConfirmURL": confirmURL,
		"ExpiryTime": user.ConfirmationExpiresAt.Format("2006-01-02 15:04:05 MST"),
		"ToAddr":     []string{user.Email},
	}

	// Render email template
	message, err := templateManager.RenderTemplate("email_confirmation", templateData)
	if err != nil {
		return fmt.Errorf("failed to render email confirmation template: %w", err)
	}

	return sender.SendEmail(message)
}

// ConfirmEmail confirms a user's email address using the provided token
func (a *AuthService) ConfirmEmail(token string) error {
	// Get user by confirmation token
	user, err := a.userService.GetUserByConfirmationToken(token)
	if err != nil {
		return fmt.Errorf("invalid confirmation token: %w", err)
	}

	// Validate the token
	if !user.IsEmailConfirmationValid(token) {
		return fmt.Errorf("confirmation token is invalid or expired")
	}

	// Confirm the user's email
	if err := a.userService.ConfirmUserEmail(user.ID); err != nil {
		return fmt.Errorf("failed to confirm email: %w", err)
	}

	return nil
}

// UserService interface for dependency injection
type UserService interface {
	GetUserByID(userID string) (*models.User, error)
	GetUserByUsername(username string, realmID string) (*models.User, error)
	GetUserRoles(userID string) ([]*models.Role, error)
	GetRealmByName(realmName string) (*models.Realm, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
	GetUserRegistrations(req models.UserRegistrationRequest) (*models.UserRegistrationResponse, error)
	GetUserRegistrationStats(realmName string) (map[string]int64, error)
	GetUserByConfirmationToken(token string) (*models.User, error)
	ConfirmUserEmail(userID string) error
}

// SimpleUserService implements the UserService interface with database operations
type SimpleUserService struct{}

func (s *SimpleUserService) GetUserByID(userID string) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	result := db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}

	// Return the user directly since models.User already uses string IDs
	return &user, nil
}

func (s *SimpleUserService) GetUserByUsername(username string, realmID string) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	result := db.Where("username = ? AND realm_id = ?", username, realmID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", result.Error)
	}

	// Return the user directly since models.User already uses string IDs
	return &user, nil
}

func (s *SimpleUserService) GetUserRoles(userID string) ([]*models.Role, error) {
	db := database.GetDB()
	var userRoles []models.UserRole
	var roles []*models.Role

	// Get user roles through join table
	result := db.Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", result.Error)
	}

	// Get role details for each user role
	for _, userRole := range userRoles {
		var role models.Role
		if err := db.Where("id = ?", userRole.RoleID).First(&role).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to get role details: %w", err)
			}
			continue // Skip if role not found
		}

		roles = append(roles, &models.Role{
			ID:          role.ID,
			RealmID:     role.RealmID,
			Name:        role.Name,
			Description: role.Description,
			CreatedBy:   role.CreatedBy,
			CreatedAt:   role.CreatedAt,
			UpdatedBy:   role.UpdatedBy,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	return roles, nil
}

func (s *SimpleUserService) CreateUser(user *models.User) error {
	db := database.GetDB()

	// Convert auth.User to models.User
	modelUser := models.User{
		ID:             user.ID,
		RealmID:        user.RealmID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		IsActive:       user.IsActive,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedBy:      user.UpdatedBy,
		UpdatedAt:      user.UpdatedAt,
	}

	result := db.Create(&modelUser)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

func (s *SimpleUserService) UpdateUser(user *models.User) error {
	db := database.GetDB()

	// Convert auth.User to models.User
	modelUser := models.User{
		ID:             user.ID,
		RealmID:        user.RealmID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		IsActive:       user.IsActive,
		CreatedBy:      user.CreatedBy,
		CreatedAt:      user.CreatedAt,
		UpdatedBy:      user.UpdatedBy,
		UpdatedAt:      user.UpdatedAt,
	}

	result := db.Save(&modelUser)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	return nil
}

func (s *SimpleUserService) DeleteUser(userID string) error {
	db := database.GetDB()

	result := db.Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	return nil
}

func (s *SimpleUserService) GetRealmByName(realmName string) (*models.Realm, error) {
	db := database.GetDB()
	var realm models.Realm

	result := db.Where("name = ?", realmName).First(&realm)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrInvalidRealm
		}
		return nil, fmt.Errorf("failed to get realm by name: %w", result.Error)
	}

	// Return the realm directly since models.Realm already uses string IDs
	return &realm, nil
}

// SimplePolicyService implements the PolicyService interface with database operations
type SimplePolicyService struct{}

func (s *SimplePolicyService) GetUserPolicies(userID, realmID string) ([]*models.Policy, error) {
	db := database.GetDB()
	var policies []*models.Policy

	// Get direct user policies
	var userPolicies []models.UserPolicy
	result := db.Where("user_id = ?", userID).Find(&userPolicies)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user policies: %w", result.Error)
	}

	// Get policy details for direct user policies
	for _, userPolicy := range userPolicies {
		var policy models.Policy
		if err := db.Where("id = ? AND realm_id = ?", userPolicy.PolicyID, realmID).First(&policy).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("failed to get policy details: %w", err)
			}
			continue // Skip if policy not found
		}

		policies = append(policies, &models.Policy{
			ID:          policy.ID,
			RealmID:     policy.RealmID,
			Name:        policy.Name,
			Description: policy.Description,
			Version:     policy.Version,
			CreatedBy:   policy.CreatedBy,
			CreatedAt:   policy.CreatedAt,
			UpdatedBy:   policy.UpdatedBy,
			UpdatedAt:   policy.UpdatedAt,
		})
	}

	// Get user roles
	var userRoles []models.UserRole
	result = db.Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", result.Error)
	}

	// Get role policies for each user role
	for _, userRole := range userRoles {
		var rolePolicies []models.RolePolicy
		if err := db.Where("role_id = ?", userRole.RoleID).Find(&rolePolicies).Error; err != nil {
			return nil, fmt.Errorf("failed to get role policies: %w", err)
		}

		// Get policy details for each role policy
		for _, rolePolicy := range rolePolicies {
			var policy models.Policy
			if err := db.Where("id = ? AND realm_id = ?", rolePolicy.PolicyID, realmID).First(&policy).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return nil, fmt.Errorf("failed to get role policy details: %w", err)
				}
				continue // Skip if policy not found
			}

			// Check if we already have this policy (avoid duplicates)
			policyExists := false
			for _, existingPolicy := range policies {
				if existingPolicy.ID == policy.ID {
					policyExists = true
					break
				}
			}
			if policyExists {
				continue
			}

			policies = append(policies, &models.Policy{
				ID:          policy.ID,
				RealmID:     policy.RealmID,
				Name:        policy.Name,
				Description: policy.Description,
				Version:     policy.Version,
				CreatedBy:   policy.CreatedBy,
				CreatedAt:   policy.CreatedAt,
				UpdatedBy:   policy.UpdatedBy,
				UpdatedAt:   policy.UpdatedAt,
			})
		}
	}

	return policies, nil
}

func (s *SimplePolicyService) GetPolicyStatements(policyID string) ([]*models.Statement, error) {
	db := database.GetDB()
	var statements []*models.Statement

	// Get statements for the policy
	var modelStatements []models.Statement
	result := db.Where("policy_id = ?", policyID).Find(&modelStatements)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get policy statements: %w", result.Error)
	}

	// Convert model statements to auth statements
	for _, stmt := range modelStatements {
		// Parse JSON fields
		var actions []string
		var resources []string
		var conditions map[string]interface{}

		// Parse actions (JSON array as string)
		if err := json.Unmarshal([]byte(stmt.Actions), &actions); err != nil {
			return nil, fmt.Errorf("failed to parse statement actions: %w", err)
		}

		// Parse resources (JSON array as string)
		if err := json.Unmarshal([]byte(stmt.Resources), &resources); err != nil {
			return nil, fmt.Errorf("failed to parse statement resources: %w", err)
		}

		// Parse conditions (JSON string, optional)
		if stmt.Conditions != "" {
			if err := json.Unmarshal([]byte(stmt.Conditions), &conditions); err != nil {
				return nil, fmt.Errorf("failed to parse statement conditions: %w", err)
			}
		}

		statements = append(statements, &models.Statement{
			ID:         stmt.ID,
			PolicyID:   policyID,
			SID:        stmt.SID,
			Effect:     stmt.Effect,
			Actions:    marshalJSON(actions),
			Resources:  marshalJSON(resources),
			Conditions: marshalJSON(conditions),
			CreatedBy:  stmt.CreatedBy,
			CreatedAt:  stmt.CreatedAt,
			UpdatedBy:  stmt.UpdatedBy,
			UpdatedAt:  stmt.UpdatedAt,
		})
	}

	return statements, nil
}

func marshalJSON(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// GetUserRegistrations retrieves users based on registration criteria with pagination
func (s *SimpleUserService) GetUserRegistrations(req models.UserRegistrationRequest) (*models.UserRegistrationResponse, error) {
	db := database.GetDB()
	var users []models.User
	var total int64

	// Build query
	query := db.Model(&models.User{})

	// Filter by realm if specified
	if req.RealmName != "" {
		realm, err := s.GetRealmByName(req.RealmName)
		if err != nil {
			return nil, fmt.Errorf("failed to get realm: %w", err)
		}
		query = query.Where("realm_id = ?", realm.ID)
	}

	// Filter by status if specified
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Set default pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Limit max page size
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &models.UserRegistrationResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserRegistrationStats returns statistics about user registrations by status
func (s *SimpleUserService) GetUserRegistrationStats(realmName string) (map[string]int64, error) {
	db := database.GetDB()
	stats := make(map[string]int64)

	query := db.Model(&models.User{})

	// Filter by realm if specified
	if realmName != "" {
		realm, err := s.GetRealmByName(realmName)
		if err != nil {
			return nil, fmt.Errorf("failed to get realm: %w", err)
		}
		query = query.Where("realm_id = ?", realm.ID)
	}

	// Get counts for each status
	statuses := []models.UserStatus{
		models.UserStatusPending,
		models.UserStatusConfirmed,
		models.UserStatusApproved,
		models.UserStatusDenied,
		models.UserStatusSuspended,
	}

	for _, status := range statuses {
		var count int64
		if err := query.Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count users with status %s: %w", status, err)
		}
		stats[string(status)] = count
	}

	// Get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}
	stats["total"] = totalCount

	return stats, nil
}

// GetUserByConfirmationToken retrieves a user by their email confirmation token
func (s *SimpleUserService) GetUserByConfirmationToken(token string) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("email_confirmation_token = ? AND confirmation_expires_at > ?", token, time.Now()).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired confirmation token")
		}
		return nil, fmt.Errorf("failed to get user by confirmation token: %w", err)
	}

	return &user, nil
}

// ConfirmUserEmail confirms a user's email address
func (s *SimpleUserService) ConfirmUserEmail(userID string) error {
	db := database.GetDB()

	// Get the user first to call the model method
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Use the model method to confirm email
	user.ConfirmEmail()

	// Update the user in the database
	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
