package service

import (
	"context"
	"time"

	"github.com/nanayaw/fullstack/internal/models"
)

type AuthService interface {
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

type UserService interface {
	// User management
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*models.User, int, error)

	// Profile management
	UpdateProfile(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error)
	UploadAvatar(ctx context.Context, userID string, fileData []byte, fileType string) (string, error)
	RemoveAvatar(ctx context.Context, userID string) error

	// Activity tracking
	GetUserActivity(ctx context.Context, userID string, page, pageSize int) ([]*models.AuditLog, int, error)
}

type EmailService interface {
	// Email sending
	SendVerificationEmail(ctx context.Context, to string, token string) error
	SendPasswordResetEmail(ctx context.Context, to string, token string) error
	SendWelcomeEmail(ctx context.Context, to string, userName string) error
	SendLoginNotificationEmail(ctx context.Context, to string, deviceInfo string, location string) error
	SendPasswordChangedEmail(ctx context.Context, to string) error

	// Template management
	ParseTemplate(templateName string, data interface{}) (string, error)
	ValidateEmailAddress(email string) bool
}

type CacheService interface {
	// Rate limiting
	CheckRateLimit(ctx context.Context, key string, limit int, duration int) (bool, error)
	ResetRateLimit(ctx context.Context, key string) error

	// Caching
	CacheData(ctx context.Context, key string, data interface{}, ttl int) error
	GetCachedData(ctx context.Context, key string, dest interface{}) error
	InvalidateCache(ctx context.Context, key string) error
	InvalidateCachePattern(ctx context.Context, pattern string) error

	// Session management
	StoreSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error
	GetSession(ctx context.Context, sessionID string) (string, error)
	InvalidateSession(ctx context.Context, sessionID string) error
}
