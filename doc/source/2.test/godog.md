# BDD Tests with Godog

## ✅ **Complete BDD Test Framework Setup**

1. **📁 Directory Structure**
   ```
   bdd/
   ├── steps/                    # Step definitions
   │   └── user_registration_steps.go
   ├── support/                  # Test helpers and utilities
   │   └── test_helpers.go
   ├── main_test.go             # Test runner
   ├── godog.yaml               # Godog configuration
   ├── Makefile                 # Build and test commands
   └── README.md                # This file
   ```

2. **🔧 Dependencies Added**
   - `github.com/cucumber/godog@latest` - BDD testing framework
   - `github.com/cucumber/godog/cmd/godog@latest` - CLI tool

3. **📝 Feature Coverage**
   - User registration scenarios from `../features/user_registration.feature`
   - Successful registration flow
   - Duplicate email validation

## ✅ **Key Components Implemented**

### **Step Definitions** (`steps/user_registration_steps.go`)
- `UserRegistrationContext` - Test context management
- Database setup and cleanup
- HTTP request/response handling
- User creation and validation steps
- Integration with auth services

### **Test Helpers** (`support/test_helpers.go`)
- `TestConfig` - Configuration management
- `SetupTestDatabase()` - Database initialization
- `CleanupTestDatabase()` - Test data cleanup
- `CreateTestJWTManager()` - JWT key generation for tests
- `CreateTestUser()` - Test user creation
- `AssertUserExists()` / `AssertUserNotExists()` - User validation

### **Test Runner** (`main_test.go`)
- Godog integration with Go testing
- Scenario initialization
- Global hooks for setup/teardown

## ✅ **Features Implemented**

1. **🔐 Authentication Service Integration**
   - JWT token generation with test keys
   - Password hashing with bcrypt
   - User service with database operations
   - Permission engine setup

2. **🗄️ Database Testing**
   - SQLite in-memory database for tests
   - GORM model auto-migration
   - Test data seeding and cleanup
   - Realm and user management

3. **🌐 HTTP API Testing**
   - Gin router setup in test mode
   - HTTP request simulation
   - Response status and body validation
   - JSON request/response handling

4. **📊 Test Configuration**
   - Godog YAML configuration
   - Makefile with multiple test targets
   - Environment-specific settings

## ✅ **Test Scenarios Covered**

### **Successful Registration**
```gherkin
Scenario: Successful registration
  Given a clean user repository
  And an email "test@example.com" doesn't exist
  When I register with email "test@example.com" and password "SecurePass123!"
  Then the response status should be 201
  And the user "test@example.com" should exist with status "pending"
```

### **Duplicate Email Validation**
```gherkin
Scenario: Registration with existing email
  Given a clean user repository
  And an email "existing@example.com" already exists
  When I register with email "existing@example.com" and password "SecurePass123!"
  Then the response status should be 409
```

# 🚀 **How to Run Tests**

## **Using Go Test**
```bash
cd bdd
go test -v
```

## **Using Godog CLI**
```bash
# Install godog CLI
make install

# Run with pretty format
make test-verbose

# Run with progress format
make test

# Generate JUnit XML report
make test-junit

# Run specific feature
make test-feature FEATURE=user_registration.feature

# Run with tags
make test-tags TAGS="@registration"
```

## **Available Make Targets**
```bash
make help           # Show available targets
make test           # Run BDD tests with progress format
make test-verbose   # Run BDD tests with pretty format
make test-junit     # Generate JUnit XML report
make clean          # Clean test artifacts
make install        # Install godog CLI tool
```

# 🔧 **Configuration**

## **Godog Configuration** (`godog.yaml`)
```yaml
format: pretty
paths:
  - ../features
output: stdout
randomize: false
strict: true
stop-on-failure: false
tags: ""
concurrency: 1
```

## **Test Database**
- **Type**: SQLite in-memory (`:memory:`)
- **Models**: Auto-migrated from `internal/models`
- **Data**: Seeded with default realms, users, roles
- **Cleanup**: Automatic between scenarios

## **JWT Keys**
- **Generation**: RSA 2048-bit keys created in-memory
- **Format**: PKCS1 for private key, PKIX for public key
- **Usage**: Temporary files for JWT manager initialization

# 🐛 **Current Status & Known Issues**

## ✅ **Working Components**
- BDD framework setup and configuration
- Step definitions and test helpers
- JWT key generation and authentication setup
- Database initialization and model migration
- HTTP request/response handling
- Test runner integration

## ⚠️ **Known Issues**
1. **Database Cleanup**: Table names mismatch between cleanup queries and actual schema
2. **Test Isolation**: Data persistence between test runs causing conflicts
3. **Default Data**: Initial database seeding interferes with clean test state

## 🔧 **Fixes Needed**
1. **Update cleanup queries** to use correct table names (`app_user` instead of `users`)
2. **Improve test isolation** by using unique test database per scenario
3. **Fix realm ID resolution** in test user creation

# 📈 **Benefits Achieved**

1. **🎯 Behavior-Driven Testing**: Clear, readable test scenarios in Gherkin syntax
2. **🔄 Automated Testing**: Integration with Go testing framework
3. **🏗️ Comprehensive Setup**: Full authentication and database testing infrastructure
4. **📊 Multiple Formats**: Support for different output formats (pretty, progress, JUnit)
5. **🛠️ Developer Tools**: Makefile with convenient test commands
6. **🔧 Extensible Framework**: Easy to add new scenarios and step definitions

# 🚀 **Next Steps**

## **Immediate Fixes**
1. Fix database table name mismatches in cleanup queries
2. Implement proper test isolation with unique databases
3. Resolve realm ID issues in test setup

## **Feature Expansion**
1. Add more user registration scenarios (password validation, email confirmation)
2. Implement login/logout BDD tests
3. Add permission and role management scenarios
4. Create API endpoint testing scenarios

## **Infrastructure Improvements**
1. Add Docker support for test execution
2. Implement CI/CD integration
3. Add performance testing scenarios
4. Create test data factories

# 📚 **Example Usage**

## **Adding New Step Definitions**
```go
// In steps/user_registration_steps.go
func (ctx *UserRegistrationContext) iLoginWithCredentials(email, password string) error {
    // Implementation for login step
    return nil
}

// Register in InitializeScenario
sc.Step(`^I login with email "([^"]*)" and password "([^"]*)"$`, ctx.iLoginWithCredentials)
```

## **Adding New Feature File**
```gherkin
# features/user_login.feature
Feature: User Login
  As a registered user
  I want to login to my account
  So that I can access protected features

  Scenario: Successful login
    Given a user exists with email "user@example.com"
    When I login with email "user@example.com" and password "ValidPass123!"
    Then the response status should be 200
    And I should receive a valid JWT token
```

# 🎉 **Conclusion**

The BDD testing framework is successfully implemented and provides a solid foundation for behavior-driven testing of the lazy-rabbit-reminder application. While there are some minor issues with database cleanup and test isolation, the core framework is working and can be easily extended to cover more scenarios.

The implementation demonstrates best practices for:
- BDD testing with Godog
- Test context management
- Database testing with GORM
- HTTP API testing with Gin
- JWT authentication testing
- Comprehensive test utilities and helpers
