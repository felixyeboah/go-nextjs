# Security Service

This package provides security-related functionality for the application, including:

- Login attempt tracking and account locking
- Suspicious activity detection
- Security event logging
- Login notifications for new devices/locations
- Password change notifications

## Features

### Account Security

- **Login Attempt Tracking**: Records all login attempts (successful and failed)
- **Account Locking**: Automatically locks accounts after too many failed login attempts
- **Device Fingerprinting**: Detects logins from new devices and locations
- **Suspicious Activity Detection**: Identifies potentially suspicious account activity

### Notifications

- **Login Notifications**: Sends email notifications for logins from new devices or locations
- **Password Change Notifications**: Alerts users when their password is changed
- **Account Lock Notifications**: Informs users when their account is locked
- **Suspicious Activity Alerts**: Warns users about potentially suspicious activity

### Geolocation

- **IP Geolocation**: Identifies the location of login attempts and other activities
- **Private IP Detection**: Properly handles private and local IP addresses

## Configuration

The security service is configured through the application's main configuration:

```go
type SecurityConfig struct {
    MaxLoginAttempts                int           // Maximum failed login attempts before locking
    AccountLockDuration             time.Duration // How long accounts remain locked
    EnableLoginNotifications        bool          // Whether to send login notifications
    EnableSuspiciousActivityDetection bool        // Whether to detect suspicious activity
    EnableRateLimiting              bool          // Whether to enable rate limiting
    GlobalRateLimit                 int           // Global rate limit (requests per minute)
    AuthRateLimit                   int           // Auth rate limit (login attempts per 15 minutes)
}
```

## Usage

### Basic Usage

```go
// Create a new security service
securityService := security.NewService(
    securityRepo,    // Repository implementation
    emailService,    // Email service for notifications
    appConfig,       // Application configuration
    logger,          // Logger
)

// Record a login attempt
err := securityService.RecordLoginAttempt(
    ctx,
    userID,
    userEmail,
    ipAddress,
    userAgent,
    wasSuccessful,
)

// Check if an account is locked
isLocked, unlockTime, reason, err := securityService.IsAccountLocked(ctx, userID)

// Detect suspicious activity
err := securityService.DetectSuspiciousActivity(
    ctx,
    userID,
    userEmail,
    ipAddress,
    userAgent,
    "password_reset_requested",
)

// Notify about password change
err := securityService.NotifyPasswordChanged(
    ctx,
    userID,
    userEmail,
    ipAddress,
    userAgent,
)
```

### Custom GeoIP Lookup

You can provide a custom implementation of the GeoIP lookup interface:

```go
// Implement the GeoIPLookup interface
type MyGeoIPLookup struct {
    // Your fields here
}

func (g *MyGeoIPLookup) GetLocation(ip string) (*security.Location, error) {
    // Your implementation here
}

// Set the custom implementation
securityService.SetGeoIPLookup(&MyGeoIPLookup{})
```

## Database Schema

The security service requires the following database tables:

### Login Attempts

```sql
CREATE TABLE login_attempts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    location VARCHAR(255),
    successful BOOLEAN NOT NULL,
    attempted_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_login_attempts_user_id ON login_attempts(user_id);
CREATE INDEX idx_login_attempts_attempted_at ON login_attempts(attempted_at);
```

### Security Events

```sql
CREATE TABLE security_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    location VARCHAR(255),
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_security_events_user_id ON security_events(user_id);
CREATE INDEX idx_security_events_created_at ON security_events(created_at);
CREATE INDEX idx_security_events_event_type ON security_events(event_type);
```

### Account Locks

```sql
CREATE TABLE account_locks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    locked_at TIMESTAMP NOT NULL,
    unlock_at TIMESTAMP NOT NULL,
    reason TEXT NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_account_locks_user_id ON account_locks(user_id);
CREATE INDEX idx_account_locks_unlock_at ON account_locks(unlock_at);
```

## Best Practices

1. **Always check if an account is locked** before validating credentials
2. **Record all security-related events** for audit purposes
3. **Notify users about security events** that affect their account
4. **Use rate limiting** to prevent brute force attacks
5. **Implement proper IP geolocation** for better security insights 