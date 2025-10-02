package models

// Auth-specific models that are not part of the core database schema
// These are used for API requests/responses and business logic

// LoginRequest represents a login request
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	RealmName string `json:"realm_name" binding:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         User   `json:"user"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateUserRequest represents a user creation request
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	RealmName string `json:"realm_name" binding:"required"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Username        string     `json:"username"`
	Email           string     `json:"email" binding:"omitempty,email"`
	Status          UserStatus `json:"status"`
	RoleIDs         []string   `json:"role_ids"`
	CurrentPassword string     `json:"current_password"`
	NewPassword     string     `json:"new_password" binding:"omitempty,min=8"`
}

// PermissionCheck represents a permission check request
type PermissionCheck struct {
	Action   string                 `json:"action" binding:"required"`
	Resource string                 `json:"resource" binding:"required"`
	UserID   string                 `json:"user_id" binding:"required"`
	RealmID  string                 `json:"realm_id" binding:"required"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// UserRegistrationRequest represents a request to get pending registrations
type UserRegistrationRequest struct {
	RealmName string     `json:"realm_name,omitempty"`
	Status    UserStatus `json:"status,omitempty"`
	Page      int        `json:"page,omitempty"`
	PageSize  int        `json:"page_size,omitempty"`
}

// UserRegistrationResponse represents the response for listing registrations
type UserRegistrationResponse struct {
	Users      []User `json:"users"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// ApproveRegistrationRequest represents a request to approve/deny registration
type ApproveRegistrationRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Approved bool   `json:"approved" binding:"required"`
	Reason   string `json:"reason,omitempty"`
}

// ApproveRegistrationResponse represents the response for registration approval
type ApproveRegistrationResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
