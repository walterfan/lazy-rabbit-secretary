
# permission system with read-only support

- AWS IAM-style policies with Allow/Deny
- Realm-based multi-tenancy
- Role-based access control (RBAC)
- JWT authentication
- Missing: granular read/write permissions and module-level access control


## Overview

Redesigned the user/role/policy/statement management module into a permission system with read-only support.

### **Backend Changes**

#### 1. **Enhanced Permission Models** (`backend/internal/models/permission.go`)
- **Permission Actions**: CRUD, administrative, and special actions
- **Permission Resources**: Core, content, knowledge, tools, and system resources
- **Permission Levels**: `readonly`, `readwrite`, `full`, `custom`
- **UserPermission & RolePermission**: Direct and role-based permissions with expiration
- **PermissionResult**: Detailed results with source and context

#### 2. **Permission Service** (`backend/internal/auth/permission_service.go`)
- **Permission Checking**: Evaluates user and role permissions
- **Read-Only Support**: Enforces read-only access
- **Permission Management**: CRUD for user and role permissions
- **Permission Summary**: Overview of user access
- **Default Permissions**: Initializes admin and user permissions

#### 3. **Permission Middleware** (`backend/internal/auth/permission_middleware.go`)
- **RequirePermission**: Enforces specific action/resource permissions
- **RequireReadPermission**: Read-only access
- **RequireWritePermission**: Write access (create, update, delete)
- **RequireManagePermission**: Administrative access
- **OptionalPermission**: Optional checks for read-only features
- **CheckReadOnlyAccess**: Detects read-only mode

#### 4. **Permission Handlers** (`backend/internal/auth/permission_handlers.go`)
- **Permission Management**: CRUD for user and role permissions
- **Permission Checking**: Real-time checks
- **Permission Summary**: User access overview
- **Available Options**: Lists actions, resources, and levels

#### 5. **Permission Routes** (`backend/internal/auth/permission_routes.go`)
- **Admin Routes**: Permission management (requires manage permissions)
- **User Routes**: Current user permission access
- **Protected Endpoints**: Permission-based access control

### **Frontend Changes**

#### 1. **Permission Management UI** (`frontend/src/views/PermissionManagementView.vue`)
- **Permission Summary**: Current user permissions
- **User Permissions**: Manage direct user permissions
- **Role Permissions**: Manage role-based permissions
- **Permission Checker**: Test permissions
- **Create/Edit/Delete**: Full CRUD
- **Real-time Validation**: Immediate feedback

#### 2. **Navigation Integration**
- **Admin Menu**: Added to admin dropdown
- **Route Protection**: Requires authentication
- **Internationalization**: Added translation keys

### **Features**

#### **Read-Only Support**
- **Read-Only Level**: Users can view but not modify
- **Automatic Detection**: Middleware detects read-only access
- **UI Adaptation**: Frontend can adapt to read-only mode
- **Granular Control**: Per-resource read-only permissions

#### **Permission Levels**
- **Read-Only**: View only
- **Read-Write**: View and modify
- **Full Access**: All operations including admin
- **Custom**: Specific action combinations

#### **Permission Sources**
- **User Permissions**: Direct assignments
- **Role Permissions**: Inherited from roles
- **Policy Permissions**: AWS-style policies (existing)
- **Priority System**: User > Role > Policy

#### **Expiration Support**
- **Time-Limited Permissions**: Expiration dates
- **Automatic Cleanup**: Removes expired permissions
- **Flexible Management**: Temporary access grants

### **Usage Examples**

#### **Backend Middleware Usage**
```go
// Require read permission
router.GET("/posts", middleware.RequireReadPermission("posts"), handler)

// Require write permission  
router.POST("/posts", middleware.RequireWritePermission("posts"), handler)

// Check read-only access
router.GET("/posts", middleware.CheckReadOnlyAccess("posts"), handler)
```

#### **Frontend Permission Check**
```typescript
// Check if user has permission
const result = await checkPermission("read", "posts");
if (result.allowed) {
  // Show content
} else {
  // Show read-only or hide content
}
```

### **Security Benefits**
1. **Principle of Least Privilege**: Users get only required access
2. **Read-Only by Default**: Safe default for sensitive data
3. **Granular Control**: Per-resource permissions
4. **Audit Trail**: Tracks permission changes
5. **Expiration Support**: Time-limited access
6. **Multi-Source Permissions**: Flexible assignment methods

The system supports read-only access and granular permissions across the application.


## example: How to Check Permissions for `/api/v1/admin/users`

The system uses role-based access control (RBAC) and permission-based access control (PBAC). Here's how to check permissions for the admin users endpoint:

### **Current Protection Method**

The `/api/v1/admin/users` endpoint is currently protected using **role-based access control**:

```go
// From backend/internal/auth/routes.go lines 29-31
admin := router.Group("/api/v1/admin")
admin.Use(middleware.Authenticate())
admin.Use(middleware.RequireRole("admin", "super_admin"))
{
    admin.GET("/users", userHandlers.GetUsers)  // This endpoint
    // ... other admin endpoints
}
```

### **Permission Checking Methods**

#### **1. Role-Based Check (Current Implementation)**

```go
// Check if user has admin or super_admin role
admin.Use(middleware.RequireRole("admin", "super_admin"))
```

**How it works:**
- User must be authenticated (JWT token)
- User must have either "admin" or "super_admin" role
- Roles are checked from JWT claims

#### **2. Permission-Based Check (Enhanced System)**

To use the new permission system, you would change the route to:

```go
// Using the new permission middleware
admin.Use(permissionMiddleware.RequirePermission("read", "users"))
// or
admin.Use(permissionMiddleware.RequireReadPermission("users"))
```

**How it works:**
- User must be authenticated
- System checks for specific permission: `read` action on `users` resource
- Supports multiple permission sources: User → Role → Policy
- Supports permission levels: `readonly`, `readwrite`, `full`, `custom`

### **Permission Check Process**

#### **Step 1: Authentication**
```go
// JWT token validation
claims, err := m.authService.ValidateToken(tokenString)
// Sets user context: user_id, realm_id, username, email, roles
```

#### **Step 2: Permission Evaluation**
```go
// Check permission using the permission service
result, err := m.permissionService.CheckPermission(
    userID,     // Current user ID
    realmID,    // User's realm
    "read",     // Action (read, create, update, delete, manage)
    "users",    // Resource (users, roles, policies, etc.)
    context,    // Request context (path params, query params, etc.)
)
```

#### **Step 3: Permission Sources (Priority Order)**
1. **Direct User Permissions** - `UserPermission` table
2. **Role Permissions** - `RolePermission` table (via user roles)
3. **Policy Permissions** - AWS-style policies (existing system)

### **Available Permission Actions**

```go
// From backend/internal/models/permission.go
const (
    ActionCreate   PermissionAction = "create"
    ActionRead     PermissionAction = "read"      // For GET /api/v1/admin/users
    ActionUpdate   PermissionAction = "update"
    ActionDelete   PermissionAction = "delete"
    ActionManage   PermissionAction = "manage"    // Full CRUD + admin operations
    ActionApprove  PermissionAction = "approve"
    ActionReject   PermissionAction = "reject"
    ActionSuspend  PermissionAction = "suspend"
    ActionActivate PermissionAction = "activate"
    ActionAll      PermissionAction = "*"         // Wildcard
)
```

### **Available Permission Resources**

```go
const (
    ResourceUsers     PermissionResource = "users"     // For /api/v1/admin/users
    ResourceRoles     PermissionResource = "roles"
    ResourcePolicies  PermissionResource = "policies"
    ResourceRealms    PermissionResource = "realms"
    // ... other resources
)
```

### **Permission Levels**

```go
const (
    LevelReadOnly PermissionLevel = "readonly"   // Can only read/view
    LevelReadWrite PermissionLevel = "readwrite" // Can read and modify
    LevelFullAccess PermissionLevel = "full"     // Full access including admin operations
    LevelCustom PermissionLevel = "custom"       // Custom permissions defined in statements
)
```

### **How to Check Permissions Programmatically**

#### **Backend API Call**
```bash
# Check if current user can read users
GET /api/v1/user-permissions/check?action=read&resource=users
Authorization: Bearer <jwt_token>

# Response
{
  "allowed": true,
  "level": "readonly",
  "actions": ["read"],
  "source": "role",
  "reason": "",
  "expires_at": null,
  "context": {...}
}
```

#### **Frontend Permission Check**
```typescript
// Check permission in Vue component
const checkUserPermission = async () => {
  try {
    const response = await fetch('/api/v1/user-permissions/check?action=read&resource=users', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    });
    const result = await response.json();
    
    if (result.allowed) {
      // User can access the endpoint
      console.log('Permission level:', result.level);
      console.log('Permission source:', result.source);
    } else {
      // User cannot access
      console.log('Access denied:', result.reason);
    }
  } catch (error) {
    console.error('Permission check failed:', error);
  }
};
```

### **Default Permissions**

The system initializes default permissions:

#### **Admin Role Permissions**
```go
// Admin gets full access to all resources
for _, resource := range adminResources {
    req := &models.RolePermissionRequest{
        RoleID:   adminRole.ID,
        RealmID:  realm.ID,
        Resource: resource,
        Actions:  []models.PermissionAction{models.ActionAll},
        Level:    models.LevelFullAccess,
    }
    // Creates permission
}
```

#### **User Role Permissions**
```go
// Regular users get read-only access to most resources
req := &models.RolePermissionRequest{
    RoleID:   userRole.ID,
    RealmID:  realm.ID,
    Resource: "users",
    Actions:  []models.PermissionAction{models.ActionRead},
    Level:    models.LevelReadOnly,
}
```

### **Migration from Role-Based to Permission-Based**

To migrate the admin users endpoint to use the new permission system:

#### **Step 1: Update Route Definition**
```go
// Change from:
admin.Use(middleware.RequireRole("admin", "super_admin"))

// To:
admin.Use(permissionMiddleware.RequireReadPermission("users"))
```

#### **Step 2: Set Up Permissions**
```go
// Create permission for admin role
permissionService.CreateRolePermission(&models.RolePermissionRequest{
    RoleID:   adminRoleID,
    RealmID:  realmID,
    Resource: "users",
    Actions:  []models.PermissionAction{models.ActionRead, models.ActionCreate, models.ActionUpdate, models.ActionDelete},
    Level:    models.LevelFullAccess,
}, "system")
```

#### **Step 3: Test Permission Check**
```bash
# Test the permission check
curl -H "Authorization: Bearer <token>" \
     "https://localhost:5173/api/v1/user-permissions/check?action=read&resource=users"
```

### **Benefits of Permission-Based System**

1. **Granular Control**: Per-action permissions (read, create, update, delete)
2. **Read-Only Support**: Users can have view-only access
3. **Expiration Support**: Time-limited permissions
4. **Multi-Source**: User, role, and policy permissions
5. **Audit Trail**: Track permission changes
6. **Flexible**: Easy to modify without code changes

The current system uses role-based access, but the enhanced permission system provides more granular control and better security.