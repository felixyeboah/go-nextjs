# Backend Service

A modern backend service built with Go, featuring:

- Authentication using PASETO tokens
- Email verification with Resend
- Rate limiting and caching with Redis
- Database management with Turso
- OAuth integration (Google & GitHub)
- Swagger documentation
- Comprehensive test coverage
- Security scanning with gosec and nancy

## Prerequisites

- Go 1.21 or later
- Redis
- Turso CLI
- Make (optional, for using Makefile commands)

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fullstack.git
   cd fullstack/backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Generate PASETO keys:
   ```bash
   make generate-keys
   ```

4. Copy environment variables:
   ```bash
   cp .env.example .env
   ```
   Then update the values in `.env` with your configuration.

5. Initialize the database:
   ```bash
   make migrate-up
   ```

## Development

Start the development server:
```bash
make run
```

The API will be available at `http://localhost:8080`.
Swagger documentation can be accessed at `http://localhost:8080/swagger/index.html`.

## Available Commands

- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage report
- `make clean` - Clean build artifacts
- `make migrate` - Create a new migration
- `make migrate-up` - Apply migrations
- `make migrate-down` - Rollback migrations
- `make swagger` - Generate Swagger documentation locally
- `make docker-swagger` - Generate Swagger documentation in Docker container
- `make lint` - Run linter
- `make security-check` - Run security checks on the codebase
- `make sqlc` - Generate SQLc code locally
- `make docker-sqlc` - Generate SQLc code in Docker container
- `make generate-keys` - Generate PASETO keys

## Docker Setup

The backend can be run using Docker:

```bash
# Build the Docker image
docker build -t fullstack-backend .

# Run the container
docker run -p 8080:8080 fullstack-backend
```

For development with hot reloading, use Docker Compose from the root directory:

```bash
docker-compose up -d backend
```

This will start the backend service along with PostgreSQL and Redis, which are required for full functionality.

## Generating Swagger Documentation

After adding new API endpoints or modifying existing ones, you need to update the Swagger documentation:

```bash
# If running locally
make swagger

# If using Docker
make docker-swagger
```

The Swagger documentation is generated using annotations in your handler code. For example:

```go
// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
    // Handler implementation
}
```

## Generating SQLc Code

After modifying SQL queries or database schema, you need to regenerate the SQLc code:

```bash
# If running locally
make sqlc

# If using Docker
make docker-sqlc
```

SQLc generates type-safe Go code from SQL queries. The queries are defined in the `internal/repository/turso/queries` directory.

## Project Structure

```
backend/
├── cmd/
│   └── api/            # Application entrypoint
├── internal/
│   ├── config/         # Configuration
│   ├── errors/         # Custom error types
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Domain models
│   ├── repository/     # Data access layer
│   │   └── turso/     # Turso implementation
│   └── service/        # Business logic
│       ├── auth/       # Authentication
│       ├── cache/      # Caching
│       └── email/      # Email service
├── migrations/         # Database migrations
├── scripts/           # Utility scripts
└── docs/             # Documentation
```

## API Documentation

The API is documented using Swagger. After starting the server, visit:
`http://localhost:8080/swagger/index.html`

## Testing

Run the test suite:
```bash
make test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Security Features

The application includes several security features to protect user accounts and prevent abuse:

### Account Security

- **Rate Limiting**: Prevents brute force attacks and API abuse
- **Account Locking**: Automatically locks accounts after too many failed login attempts
- **Login Notifications**: Alerts users about logins from new devices or locations
- **Suspicious Activity Detection**: Identifies potentially suspicious account activity
- **Password Change Notifications**: Alerts users when their password is changed

### Email Templates

The application includes responsive, accessible email templates for various authentication flows:

- **Verification Emails**: Sent to users to verify their email address during registration
- **Password Reset Emails**: Sent when a user requests a password reset
- **Welcome Emails**: Sent to users after successful registration and verification
- **Login Notification Emails**: Alerts users about new logins to their account
- **Password Changed Emails**: Notifies users when their password has been changed
- **Account Locked Emails**: Informs users when their account has been temporarily locked
- **Suspicious Activity Emails**: Alerts users about potentially suspicious activity

### Configuration

Security features can be configured in the application's configuration:

```
# Security Settings
SECURITY_MAX_LOGIN_ATTEMPTS=5
SECURITY_ACCOUNT_LOCK_DURATION=30m
SECURITY_ENABLE_LOGIN_NOTIFICATIONS=true
SECURITY_ENABLE_SUSPICIOUS_ACTIVITY_DETECTION=true
SECURITY_ENABLE_RATE_LIMITING=true
SECURITY_GLOBAL_RATE_LIMIT=100
SECURITY_AUTH_RATE_LIMIT=5 
```

## Security Scanning

The application includes security scanning tools to identify vulnerabilities:

### Running Security Checks

```bash
make security-check
```

This command runs:

1. **gosec**: Static analysis tool for Go code that finds security issues
2. **nancy**: Checks dependencies for known vulnerabilities

The security check script (`scripts/security_check.sh`) automatically installs the required tools if they're not already available. The script performs the following actions:

- Checks if gosec is installed, and installs it if necessary
- Checks if nancy is installed, and installs it if necessary
- Runs gosec to scan the codebase for security issues
- Runs nancy to check dependencies for known vulnerabilities

### Security Best Practices

The codebase follows these security best practices:

- Input validation for all user-provided data
- Proper error handling to prevent information leakage
- Secure password storage with Argon2id
- Rate limiting to prevent brute force attacks
- Secure token management with PASETO
- Middleware for security headers
- Timeout middleware to prevent long-running requests

## Docker Compose

The application can be run using Docker Compose for local development:

```bash
docker-compose up -d
```

This will start:
- Backend service
- Frontend service
- PostgreSQL database
- Redis cache
- PgAdmin for database management

Access the services:
- Backend API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/
- PgAdmin: http://localhost:5050 (admin@example.com / admin) 