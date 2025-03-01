-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash,
    full_name,
    avatar_url
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE(?, email),
    full_name = COALESCE(?, full_name),
    avatar_url = COALESCE(?, avatar_url),
    email_verified = COALESCE(?, email_verified),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    refresh_token,
    user_agent,
    client_ip,
    expires_at
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = ? AND is_blocked = FALSE
LIMIT 1;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE refresh_token = ? AND is_blocked = FALSE
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;

-- name: BlockSession :exec
UPDATE sessions
SET is_blocked = TRUE
WHERE id = ?;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = ?;

-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (
    id,
    user_id,
    provider,
    provider_user_id,
    access_token,
    refresh_token,
    expires_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetOAuthAccount :one
SELECT * FROM oauth_accounts
WHERE provider = ? AND provider_user_id = ?
LIMIT 1;

-- name: UpdateOAuthAccount :one
UPDATE oauth_accounts
SET 
    access_token = ?,
    refresh_token = COALESCE(?, refresh_token),
    expires_at = COALESCE(?, expires_at),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteOAuthAccount :exec
DELETE FROM oauth_accounts
WHERE id = ?;

-- name: CreateVerificationToken :one
INSERT INTO verification_tokens (
    id,
    user_id,
    token,
    type,
    expires_at
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetVerificationToken :one
SELECT * FROM verification_tokens
WHERE token = ? AND type = ? AND expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: DeleteVerificationToken :exec
DELETE FROM verification_tokens
WHERE id = ?;

-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    id,
    user_id,
    action,
    entity_type,
    entity_id,
    metadata
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetAuditLogs :many
SELECT * FROM audit_logs
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?; 