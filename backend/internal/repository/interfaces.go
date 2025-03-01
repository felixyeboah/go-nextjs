package repository

import (
	"context"
	"time"

	"github.com/nanayaw/fullstack/internal/models"
)

type UserRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error)

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
	DeleteSession(ctx context.Context, id string) error
	BlockSession(ctx context.Context, id string) error
	DeleteUserSessions(ctx context.Context, userID string) error

	// OAuth operations
	CreateOAuthAccount(ctx context.Context, account *models.OAuthAccount) error
	GetOAuthAccount(ctx context.Context, provider, providerUserID string) (*models.OAuthAccount, error)
	UpdateOAuthAccount(ctx context.Context, account *models.OAuthAccount) error
	DeleteOAuthAccount(ctx context.Context, id string) error

	// Verification token operations
	CreateVerificationToken(ctx context.Context, token *models.VerificationToken) error
	GetVerificationToken(ctx context.Context, token, tokenType string) (*models.VerificationToken, error)
	DeleteVerificationToken(ctx context.Context, id string) error

	// Audit log operations
	CreateAuditLog(ctx context.Context, log *models.AuditLog) error
	GetAuditLogs(ctx context.Context, userID string, offset, limit int) ([]*models.AuditLog, error)
}

type CacheRepository interface {
	// Key-value operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error

	// Rate limiting
	IncrementCounter(ctx context.Context, key string, expiration time.Duration) (int64, error)
	ResetCounter(ctx context.Context, key string) error

	// Session management
	StoreSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error
	GetSession(ctx context.Context, sessionID string) (string, error)
	InvalidateSession(ctx context.Context, sessionID string) error

	// Cache operations
	FlushAll(ctx context.Context) error
	Ping(ctx context.Context) error
}
