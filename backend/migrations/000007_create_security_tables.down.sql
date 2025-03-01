-- Drop indexes
DROP INDEX IF EXISTS idx_login_attempts_user_id;
DROP INDEX IF EXISTS idx_login_attempts_attempted_at;
DROP INDEX IF EXISTS idx_security_events_user_id;
DROP INDEX IF EXISTS idx_security_events_created_at;
DROP INDEX IF EXISTS idx_security_events_event_type;
DROP INDEX IF EXISTS idx_account_locks_user_id;
DROP INDEX IF EXISTS idx_account_locks_unlock_at;

-- Drop tables
DROP TABLE IF EXISTS account_locks;
DROP TABLE IF EXISTS security_events;
DROP TABLE IF EXISTS login_attempts; 