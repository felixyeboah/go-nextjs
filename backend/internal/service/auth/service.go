package auth

import (
	"context"

	"github.com/nanayaw/fullstack/internal/models"
)

type Service interface {
	// Registration and login
	Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error)
	Logout(ctx context.Context, sessionID string) error

	// Email verification
	SendVerificationEmail(ctx context.Context, userID string) error
	VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest) error

	// Password management
	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error

	// OAuth
	HandleOAuthLogin(ctx context.Context, req *models.OAuthLoginRequest) (*models.LoginResponse, error)
	LinkOAuthAccount(ctx context.Context, userID string, req *models.OAuthLoginRequest) error
	UnlinkOAuthAccount(ctx context.Context, userID string, provider string) error

	// Session management
	ValidateSession(ctx context.Context, sessionID string) (*models.Session, error)
	InvalidateAllSessions(ctx context.Context, userID string) error
}
