package security

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nanayaw/fullstack/internal/model"
)

// Repository implements the security.Repository interface
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new security repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// RecordLoginAttempt records a login attempt for a user
func (r *Repository) RecordLoginAttempt(ctx context.Context, attempt *model.LoginAttempt) error {
	if attempt.ID == "" {
		attempt.ID = uuid.New().String()
	}

	query := `
		INSERT INTO login_attempts (
			id, user_id, ip_address, user_agent, location, successful, attempted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		attempt.ID,
		attempt.UserID,
		attempt.IPAddress,
		attempt.UserAgent,
		attempt.Location,
		attempt.Successful,
		attempt.AttemptedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to record login attempt: %w", err)
	}

	return nil
}

// GetRecentLoginAttempts gets recent login attempts for a user
func (r *Repository) GetRecentLoginAttempts(ctx context.Context, userID string, limit int) ([]*model.LoginAttempt, error) {
	query := `
		SELECT id, user_id, ip_address, user_agent, location, successful, attempted_at
		FROM login_attempts
		WHERE user_id = $1
		ORDER BY attempted_at DESC
		LIMIT $2
	`

	var attempts []*model.LoginAttempt
	err := r.db.SelectContext(ctx, &attempts, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent login attempts: %w", err)
	}

	return attempts, nil
}

// LockAccount locks a user account
func (r *Repository) LockAccount(ctx context.Context, userID string, until time.Time, reason string) error {
	// First check if the account is already locked
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM account_locks WHERE user_id = $1)", userID)
	if err != nil {
		return fmt.Errorf("failed to check if account is locked: %w", err)
	}

	// If the account is already locked, update the lock
	if exists {
		query := `
			UPDATE account_locks
			SET unlock_at = $1, reason = $2, locked_at = NOW()
			WHERE user_id = $3
		`

		_, err := r.db.ExecContext(ctx, query, until, reason, userID)
		if err != nil {
			return fmt.Errorf("failed to update account lock: %w", err)
		}

		return nil
	}

	// Otherwise, create a new lock
	query := `
		INSERT INTO account_locks (
			id, user_id, locked_at, unlock_at, reason, created_by
		) VALUES (
			$1, $2, NOW(), $3, $4, $5
		)
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		uuid.New().String(),
		userID,
		until,
		reason,
		"system", // Created by system
	)

	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	return nil
}

// UnlockAccount unlocks a user account
func (r *Repository) UnlockAccount(ctx context.Context, userID string) error {
	query := `DELETE FROM account_locks WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	return nil
}

// IsAccountLocked checks if a user account is locked
func (r *Repository) IsAccountLocked(ctx context.Context, userID string) (bool, time.Time, string, error) {
	query := `
		SELECT unlock_at, reason
		FROM account_locks
		WHERE user_id = $1 AND unlock_at > NOW()
	`

	var unlockAt time.Time
	var reason string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&unlockAt, &reason)
	if err != nil {
		if err == sql.ErrNoRows {
			// Account is not locked
			return false, time.Time{}, "", nil
		}
		return false, time.Time{}, "", fmt.Errorf("failed to check if account is locked: %w", err)
	}

	// Account is locked
	return true, unlockAt, reason, nil
}

// RecordSecurityEvent records a security event
func (r *Repository) RecordSecurityEvent(ctx context.Context, event *model.SecurityEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	query := `
		INSERT INTO security_events (
			id, user_id, event_type, ip_address, user_agent, location, description, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.UserID,
		event.EventType,
		event.IPAddress,
		event.UserAgent,
		event.Location,
		event.Description,
		event.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to record security event: %w", err)
	}

	return nil
}

// GetUserSecurityEvents gets security events for a user
func (r *Repository) GetUserSecurityEvents(ctx context.Context, userID string, limit int) ([]*model.SecurityEvent, error) {
	query := `
		SELECT id, user_id, event_type, ip_address, user_agent, location, description, created_at
		FROM security_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	var events []*model.SecurityEvent
	err := r.db.SelectContext(ctx, &events, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user security events: %w", err)
	}

	return events, nil
}
