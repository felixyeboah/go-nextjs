package auth

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email" example:"user@example.com"`
	Password  string `json:"password" validate:"required,min=8" example:"securepassword123"`
	FirstName string `json:"first_name" validate:"required" example:"John"`
	LastName  string `json:"last_name" validate:"required" example:"Doe"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"securepassword123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
}

// LogoutRequest represents the logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// LogoutResponse represents the logout response
type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
}

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// VerifyEmailResponse represents the email verification response
type VerifyEmailResponse struct {
	Message string `json:"message" example:"Email verified successfully"`
}

// ResetPasswordRequest represents the password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	NewPassword string `json:"new_password" validate:"required,min=8" example:"newsecurepassword123"`
}

// ResetPasswordResponse represents the password reset response
type ResetPasswordResponse struct {
	Message string `json:"message" example:"Password reset successful"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// ForgotPasswordResponse represents the forgot password response
type ForgotPasswordResponse struct {
	Message string `json:"message" example:"Password reset email sent"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Message string `json:"message" example:"Email is required"`
}
