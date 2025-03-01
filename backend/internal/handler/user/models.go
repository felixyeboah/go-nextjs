package user

// UserProfileResponse represents a user profile
type UserProfileResponse struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	CreatedAt string `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-02T12:00:00Z"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"currentsecurepassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"newsecurepassword123"`
}

// ChangePasswordResponse represents a password change response
type ChangePasswordResponse struct {
	Message string `json:"message" example:"Password changed successfully"`
}

// DeleteAccountResponse represents an account deletion response
type DeleteAccountResponse struct {
	Message string `json:"message" example:"Account deleted successfully"`
}

// UserActivityResponse represents a user activity response
type UserActivityResponse struct {
	Activities []ActivityItem `json:"activities"`
	Pagination PaginationInfo `json:"pagination"`
}

// ActivityItem represents a single activity item
type ActivityItem struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Action    string `json:"action" example:"login"`
	Timestamp string `json:"timestamp" example:"2023-01-01T12:00:00Z"`
	Details   string `json:"details" example:"Login from Chrome on Windows"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	CurrentPage int `json:"current_page" example:"1"`
	PageSize    int `json:"page_size" example:"10"`
	TotalItems  int `json:"total_items" example:"42"`
	TotalPages  int `json:"total_pages" example:"5"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Message string `json:"message" example:"Current password is incorrect"`
}
