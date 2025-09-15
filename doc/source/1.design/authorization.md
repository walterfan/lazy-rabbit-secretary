# Authentication & Authorization Module

A comprehensive, production-ready authentication and authorization system built with Go, Gin, and JWT, featuring multi-tenant support and AWS-style policy-based access control.

## 🏗️ Architecture Overview

This system implements a modern, secure authentication architecture with the following components:

- **JWT-based Authentication**: Secure token-based authentication with RSA key pairs
- **Multi-tenant Support**: Realm-based isolation for different organizations
- **Policy-based Authorization**: AWS-style IAM policies with fine-grained permissions
- **OAuth 2.0 Support**: Standard OAuth 2.0 flows for third-party integrations
- **Role-based Access Control**: Hierarchical role management within realms

## 🔐 Core Components

### 1. JWT Manager (`jwt.go`)
Handles JWT token generation, validation, and management using RSA key pairs.

**Features:**
- RSA-256 signing for enhanced security
- Configurable token expiration
- Token refresh capabilities
- Custom claims with user context

**Usage:**
```go
jwtManager, err := auth.NewJWTManager(
    "private.pem", 
    "public.pem", 
    "your-app", 
    "your-audience"
)

token, err := jwtManager.GenerateToken(
    userID, realmID, username, email, roles, 
    15*time.Minute
)
```

### 2. Password Manager (`password.go`)
Secure password handling using bcrypt with configurable cost.

**Features:**
- Bcrypt hashing with configurable cost
- Password strength validation
- Secure password verification

**Usage:**
```go
passwordManager := auth.NewPasswordManager(12) // bcrypt cost 12
hash, err := passwordManager.HashPassword("secure-password")
isValid := passwordManager.VerifyPassword("secure-password", hash)
```

### 3. Permission Engine (`permission.go`)
AWS-style policy evaluation engine for fine-grained access control.

**Features:**
- Policy statement evaluation (Allow/Deny)
- Wildcard pattern matching
- Conditional access control
- Context-aware permission checking

**Policy Example:**
```json
{
  "Effect": "Allow",
  "Action": ["read:project", "write:project"],
  "Resource": ["project:${user:project_id}"],
  "Condition": {
    "StringEquals": {
      "project:owner": "${user:id}"
    }
  }
}
```

### 4. Authentication Service (`service.go`)
Core service layer handling authentication logic and user management.

**Features:**
- User login/logout
- Token refresh
- User registration
- Permission checking

**Usage:**
```go
authService := auth.NewAuthService(
    userService, 
    passwordManager, 
    jwtManager, 
    permissionEngine
)

response, err := authService.Login(loginRequest)
```

### 5. Middleware (`middleware.go`)
Gin middleware for protecting routes and enforcing permissions.

**Available Middleware:**
- `Authenticate()`: JWT token validation
- `RequireRole(roles...)`: Role-based access control
- `RequirePermission(action, resource)`: Policy-based access control
- `RequireRealm(realmID)`: Realm isolation enforcement

**Usage:**
```go
// Protected route with authentication
router.GET("/protected", middleware.Authenticate(), handler)

// Route requiring specific role
router.GET("/admin", middleware.RequireRole("admin"), handler)

// Route requiring specific permission
router.GET("/projects/:id", 
    middleware.RequirePermission("read:project", "project:*"), 
    handler
)
```

### 6. OAuth 2.0 Support (`oauth.go`)
Full OAuth 2.0 implementation for third-party integrations.

**Supported Flows:**
- Authorization Code Grant
- Client Credentials
- User Info endpoint

**Endpoints:**
- `GET /oauth2/authorize` - Authorization endpoint
- `POST /oauth2/token` - Token endpoint
- `GET /oauth2/userinfo` - User info endpoint

## 🚀 Quick Start

### 1. Generate RSA Key Pair
```bash
# Generate private key
openssl genrsa -out private.pem 2048

# Generate public key
openssl rsa -in private.pem -pubout -out public.pem
```

### 2. Initialize the System
```go
package main

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/auth"
)

func main() {
    // Initialize components
    passwordManager := auth.NewPasswordManager(12)
    jwtManager, _ := auth.NewJWTManager(
        "private.pem", "public.pem", 
        "your-app", "your-audience"
    )
    
    // Create services (you'll need to implement these)
    userService := &YourUserService{}
    policyService := &YourPolicyService{}
    
    permissionEngine := auth.NewPermissionEngine(policyService)
    authService := auth.NewAuthService(
        userService, passwordManager, 
        jwtManager, permissionEngine
    )
    
    // Create handlers and middleware
    handlers := auth.NewAuthHandlers(authService)
    middleware := auth.NewAuthMiddleware(authService)
    
    // Setup routes
    router := gin.Default()
    auth.RegisterRoutes(router, handlers, middleware)
    
    router.Run(":8080")
}
```

### 3. Protect Your Routes
```go
// Public routes
router.GET("/public", publicHandler)

// Protected routes
protected := router.Group("/api")
protected.Use(middleware.Authenticate())
{
    protected.GET("/profile", profileHandler)
    protected.GET("/projects", 
        middleware.RequirePermission("read:projects", "project:*"), 
        projectsHandler
    )
}

// Admin routes
admin := router.Group("/admin")
admin.Use(middleware.Authenticate())
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/dashboard", dashboardHandler)
}
```

## 🔑 API Endpoints

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Token refresh
- `GET /api/v1/auth/profile` - Get user profile
- `POST /api/v1/auth/logout` - User logout

### Permission Endpoints
- `GET /api/v1/auth/permissions/check` - Check user permissions

### Admin Endpoints
- `GET /api/v1/admin/users` - List users
- `POST /api/v1/admin/users` - Create user
- `GET /api/v1/admin/roles` - List roles
- `POST /api/v1/admin/policies` - Create policy

### OAuth 2.0 Endpoints
- `GET /oauth2/authorize` - Authorization endpoint
- `POST /oauth2/token` - Token endpoint
- `GET /oauth2/userinfo` - User info endpoint

## 🛡️ Security Features

### JWT Security
- RSA-256 signing for tamper-proof tokens
- Configurable token expiration
- Secure token refresh mechanism
- Custom claims validation

### Password Security
- Bcrypt hashing with configurable cost
- Password strength requirements
- Secure password verification

### Access Control
- Multi-tenant isolation
- Role-based access control
- Policy-based permissions
- Conditional access control
- Default deny principle

### OAuth 2.0 Security
- Authorization code flow
- Secure client registration
- Token expiration and validation
- Scope-based access control

## 🏢 Multi-Tenant Architecture

The system supports multiple organizations (realms) with complete data isolation:

- **Realm Isolation**: Users can only access resources within their realm
- **Cross-Realm Policies**: Support for cross-realm access when needed
- **Realm Management**: Admin tools for realm creation and management

## 📋 Policy Examples

### Basic Read Access
```json
{
  "Effect": "Allow",
  "Action": ["read:project"],
  "Resource": ["project:*"]
}
```

### Owner-Only Access
```json
{
  "Effect": "Allow",
  "Action": ["update:project", "delete:project"],
  "Resource": ["project:${user:project_id}"],
  "Condition": {
    "StringEquals": {
      "project:owner": "${user:id}"
    }
  }
}
```

### Time-Based Access
```json
{
  "Effect": "Allow",
  "Action": ["read:reports"],
  "Resource": ["report:*"],
  "Condition": {
    "DateGreaterThan": {
      "current:time": "2024-01-01T00:00:00Z"
    }
  }
}
```

### IP-Based Access
```json
{
  "Effect": "Allow",
  "Action": ["admin:*"],
  "Resource": ["*"],
  "Condition": {
    "IpAddress": {
      "source:ip": "192.168.1.0/24"
    }
  }
}
```

## 🔧 Configuration

### Environment Variables
```bash
JWT_PRIVATE_KEY_PATH=./private.pem
JWT_PUBLIC_KEY_PATH=./public.pem
JWT_ISSUER=your-app
JWT_AUDIENCE=your-audience
BCRYPT_COST=12
```

### JWT Configuration
```go
type JWTConfig struct {
    PrivateKeyPath string
    PublicKeyPath  string
    Issuer         string
    Audience       string
    AccessTokenTTL time.Duration
    RefreshTokenTTL time.Duration
}
```

## Refresh Token Mechanism


### **Modules:**

1. **HTTP Interceptor** (`/frontend/src/utils/httpInterceptor.ts`):
   - `HttpInterceptor` class that wraps all API calls
   - Automatic 401 detection and token refresh
   - Request queuing during refresh to prevent race conditions
   - Automatic retry of failed requests after successful refresh
   - Fallback to login page if refresh fails

2. **Enhanced Auth Store** (`authStore.ts`):
   - Improved `refreshAuth()` method with better error handling
   - Support for rotating refresh tokens
   - Better logging for debugging

### **How It Works:**

1. **API Call Made**: Any store method calls `makeAuthenticatedRequest()`
2. **Request Sent**: Interceptor adds Authorization header and sends request
3. **401 Detected**: If response is 401, interceptor triggers refresh flow
4. **Token Refresh**: Calls `/api/v1/auth/refresh` with refresh token
5. **Queue Management**: Other requests wait in queue during refresh
6. **Retry Request**: Original request is retried with new access token
7. **Success**: User continues working seamlessly
8. **Failure**: User is redirected to login page

### **Updated Stores:**

**taskStore.ts** - All methods now use automatic refresh:
- `fetchTasks()`, `getTask()`, `addTask()`, `updateTask()`, `deleteTask()`
- `startTask()`, `completeTask()`, `failTask()`
- `getUpcomingTasks()`, `getOverdueTasks()`

**secretStore.ts** - All methods now use automatic refresh:
- `fetchSecrets()`, `createSecret()`, `updateSecret()`, `deleteSecret()`
- `copySecretValue()`

**authStore.ts** - Admin methods now use automatic refresh:
- `getRegistrations()`, `getRegistrationStats()`, `approveRegistration()`

### **Key Features:**

✅ **Seamless Experience**: Users never see "Please log in again" errors
✅ **Race Condition Prevention**: Multiple requests don't trigger multiple refreshes
✅ **Automatic Retry**: Failed requests are automatically retried after refresh
✅ **Fallback Handling**: If refresh fails, user is redirected to login
✅ **No Code Changes Required**: Existing components work without modification
✅ **Debugging Support**: Console logging for troubleshooting

### **Benefits:**

- **Better UX**: No interruption when tokens expire
- **Reduced Support**: Fewer "session expired" complaints
- **Automatic Recovery**: System handles token lifecycle automatically
- **Security**: Tokens are refreshed before they fully expire
- **Performance**: Failed requests are retried instead of requiring user action

### **Usage:**

The system is now fully automatic! All existing API calls will automatically:
1. Refresh tokens when they expire
2. Retry failed requests
3. Redirect to login only when refresh fails

No changes needed in components - the automatic refresh works transparently!

## 🧪 Testing

### Unit Tests
```bash
go test ./internal/auth/...
```

### Integration Tests
```bash
go test -tags=integration ./internal/auth/...
```

### Test Coverage
```bash
go test -cover ./internal/auth/...
```

## 🚀 Production Deployment

### Security Checklist
- [ ] Use strong RSA keys (2048+ bits)
- [ ] Store private keys securely
- [ ] Enable HTTPS/TLS
- [ ] Set appropriate token expiration times
- [ ] Implement rate limiting
- [ ] Enable audit logging
- [ ] Regular security updates

### Performance Considerations
- [ ] Use connection pooling for database
- [ ] Implement caching for policies
- [ ] Consider Redis for session storage
- [ ] Monitor token validation performance

### Monitoring
- [ ] Token generation/validation metrics
- [ ] Permission check performance
- [ ] Failed authentication attempts
- [ ] Policy evaluation metrics

## 📚 Additional Resources

- [JWT.io](https://jwt.io/) - JWT debugging and validation
- [OAuth 2.0 RFC](https://tools.ietf.org/html/rfc6749) - OAuth 2.0 specification
- [AWS IAM Policies](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html) - Policy examples
- [Gin Framework](https://gin-gonic.com/) - Web framework documentation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
