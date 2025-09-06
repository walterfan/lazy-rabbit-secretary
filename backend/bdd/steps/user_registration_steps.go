package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/walterfan/lazy-rabbit-reminder/bdd/support"
	"github.com/walterfan/lazy-rabbit-reminder/internal/auth"
	"github.com/walterfan/lazy-rabbit-reminder/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRegistrationContext holds the test context for user registration scenarios
type UserRegistrationContext struct {
	db           *gorm.DB
	authService  *auth.AuthService
	authHandlers *auth.AuthHandlers
	router       *gin.Engine
	response     *httptest.ResponseRecorder
	lastEmail    string
	lastPassword string
}

// NewUserRegistrationContext creates a new test context
func NewUserRegistrationContext() *UserRegistrationContext {
	return &UserRegistrationContext{}
}

// InitializeScenario sets up the test environment for each scenario
func (ctx *UserRegistrationContext) InitializeScenario(sc *godog.ScenarioContext) {
	// Given steps
	sc.Step(`^a clean user repository$`, ctx.aCleanUserRepository)
	sc.Step(`^an email "([^"]*)" doesn't exist$`, ctx.anEmailDoesntExist)
	sc.Step(`^an email "([^"]*)" already exists$`, ctx.anEmailAlreadyExists)

	// When steps
	sc.Step(`^I register with email "([^"]*)" and password "([^"]*)"$`, ctx.iRegisterWithEmailAndPassword)

	// Then steps
	sc.Step(`^the response status should be (\d+)$`, ctx.theResponseStatusShouldBe)
	sc.Step(`^the user "([^"]*)" should exist with status "([^"]*)"$`, ctx.theUserShouldExistWithStatus)

	// Before each scenario
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return ctx, nil
	})

	// After each scenario
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})
}

// aCleanUserRepository sets up a clean test database
func (ctx *UserRegistrationContext) aCleanUserRepository() error {
	// Setup test database using helper
	config := support.DefaultTestConfig()
	db, err := support.SetupTestDatabase(config)
	if err != nil {
		return fmt.Errorf("failed to setup test database: %w", err)
	}
	ctx.db = db

	// Clean up existing test data
	if err := support.CleanupTestDatabase(ctx.db); err != nil {
		return fmt.Errorf("failed to cleanup test database: %w", err)
	}

	// Initialize auth service and handlers
	userService := support.NewTestUserService(ctx.db)
	passwordManager := auth.NewPasswordManager(bcrypt.DefaultCost)

	// Create a test JWT manager using helper
	jwtManager, err := support.CreateTestJWTManager()
	if err != nil {
		return fmt.Errorf("failed to create test JWT manager: %w", err)
	}

	// Create a simple policy service for testing
	policyService := &auth.SimplePolicyService{}
	permissionEngine := auth.NewPermissionEngine(policyService)
	ctx.authService = auth.NewAuthService(userService, passwordManager, jwtManager, permissionEngine)
	ctx.authHandlers = auth.NewAuthHandlers(ctx.authService)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	ctx.router = gin.New()

	// Set up routes
	public := ctx.router.Group("/api/v1/public")
	public.POST("/register", ctx.authHandlers.Register)

	return nil
}

// anEmailDoesntExist verifies that an email doesn't exist in the database
func (ctx *UserRegistrationContext) anEmailDoesntExist(email string) error {
	return support.AssertUserNotExists(ctx.db, email)
}

// anEmailAlreadyExists creates a user with the given email
func (ctx *UserRegistrationContext) anEmailAlreadyExists(email string) error {
	_, err := support.CreateTestUser(ctx.db, email, email, models.UserStatusApproved)
	return err
}

// iRegisterWithEmailAndPassword performs a registration request
func (ctx *UserRegistrationContext) iRegisterWithEmailAndPassword(email, password string) error {
	ctx.lastEmail = email
	ctx.lastPassword = password

	// Create registration request
	registerReq := models.CreateUserRequest{
		Username:  email,
		Email:     email,
		Password:  password,
		RealmName: "default", // Default realm
	}

	// Convert to JSON
	jsonData, err := json.Marshal(registerReq)
	if err != nil {
		return fmt.Errorf("failed to marshal registration request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "/api/v1/public/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record response
	ctx.response = httptest.NewRecorder()
	ctx.router.ServeHTTP(ctx.response, req)

	return nil
}

// theResponseStatusShouldBe verifies the HTTP response status
func (ctx *UserRegistrationContext) theResponseStatusShouldBe(expectedStatus int) error {
	if ctx.response.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d. Response body: %s",
			expectedStatus, ctx.response.Code, ctx.response.Body.String())
	}
	return nil
}

// theUserShouldExistWithStatus verifies that a user exists with the expected status
func (ctx *UserRegistrationContext) theUserShouldExistWithStatus(email, expectedStatus string) error {
	return support.AssertUserExists(ctx.db, email, models.UserStatus(expectedStatus))
}
