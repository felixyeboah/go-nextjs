# Fullstack Go + Next.js Application

A modern, production-ready fullstack application built with Go (backend) and Next.js (frontend).

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Development Setup](#development-setup)
  - [Docker Setup](#docker-setup)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Authentication](#authentication)
- [Email Service](#email-service)
- [Database](#database)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)
- [Security](#security)
- [Testing](#testing)
- [Additional Documentation](#additional-documentation)
- [Recent Changes](#recent-changes)

## Overview

This project is a fullstack application template that combines a Go backend with a Next.js frontend. It includes user authentication, OAuth integration, email services, and a modern UI built with Tailwind CSS.

## Features

- **Backend (Go)**
  - RESTful API with Echo framework
  - Authentication with PASETO tokens
  - OAuth integration (Google, GitHub)
  - Email service integration (Resend, Upstash Workflow)
  - Redis caching
  - Turso database (distributed SQLite)
  - Swagger API documentation
  - Structured logging
  - Graceful shutdown
  - Comprehensive test coverage
  - Security scanning with gosec and nancy

- **Frontend (Next.js)**
  - Modern React with Next.js
  - Type-safe with TypeScript
  - Tailwind CSS for styling
  - Form validation with Zod
  - React Query for data fetching
  - Authentication with JWT
  - Responsive design

## Project Structure

### Backend Structure

```
backend/
├── cmd/                  # Application entry points
│   └── api/              # Main API server
├── docs/                 # Swagger documentation
├── internal/             # Private application code
│   ├── config/           # Configuration management
│   ├── handler/          # HTTP handlers
│   │   ├── auth/         # Authentication handlers
│   │   ├── middleware/   # HTTP middleware
│   │   └── user/         # User handlers
│   ├── repository/       # Data access layer
│   │   └── turso/        # Turso DB implementation
│   └── service/          # Business logic
│       ├── auth/         # Authentication service
│       ├── cache/        # Caching service
│       ├── email/        # Email service
│       ├── oauth/        # OAuth service
│       └── user/         # User service
├── migrations/           # Database migrations
├── pkg/                  # Public libraries
│   ├── auth/             # Authentication utilities
│   ├── email/            # Email utilities
│   └── validator/        # Validation utilities
├── scripts/              # Utility scripts
├── .env                  # Environment variables (development)
├── .env.example          # Example environment variables
├── Dockerfile            # Docker configuration
├── Makefile              # Build and development commands
├── go.mod                # Go dependencies
└── sqlc.yaml             # SQLC configuration
```

### Frontend Structure

```
frontend/
├── app/                  # Next.js app directory
│   ├── api/              # API routes
│   ├── auth/             # Authentication pages
│   └── users/            # User pages
├── public/               # Static assets
├── src/                  # Source code
│   ├── components/       # React components
│   │   ├── auth/         # Authentication components
│   │   ├── layout/       # Layout components
│   │   └── ui/           # UI components
│   ├── lib/              # Utility functions
│   │   ├── api/          # API client
│   │   └── utils/        # Helper functions
│   └── types/            # TypeScript type definitions
├── .env.local            # Environment variables
├── Dockerfile            # Docker configuration
├── next.config.ts        # Next.js configuration
├── package.json          # NPM dependencies
└── tsconfig.json         # TypeScript configuration
```

## Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 20+ (for local development)
- Make (optional, for using Makefile commands)

## Getting Started

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fullstack.git
   cd fullstack
   ```

2. Set up the backend:
   ```bash
   cd backend
   cp .env.example .env
   go mod download
   make run
   ```

3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### Docker Setup

The easiest way to run the entire application is using Docker Compose:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fullstack.git
   cd fullstack
   ```

2. Configure environment variables:
   ```bash
   # For backend
   cp backend/.env.docker.example backend/.env.docker
   # Edit backend/.env.docker with your Turso database credentials and other settings
   
   # For frontend
   cp frontend/.env.example frontend/.env
   # Edit frontend/.env with your API URL and other settings
   ```

3. Start the containers:
   ```bash
   docker-compose up -d
   ```

4. Access the application:
   - Frontend: http://localhost:3001
   - Backend API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/

5. Stop the application:
   ```bash
   docker-compose down
   ```

The Docker Compose setup includes the following services:
- **Backend**: Go API server with Turso database integration
- **Frontend**: Next.js application
- **Redis**: In-memory cache and message broker

Note: This application uses Turso database (distributed SQLite) instead of a local database container.

Each service is configured with appropriate volumes for data persistence and connected through a dedicated network.

## Configuration

### Backend Configuration

The backend is configured using environment variables. See `.env.example` for all available options.

Key configuration sections:
- Server settings
- Database connection
- Redis connection
- Authentication settings
- Email service
- OAuth providers

### Frontend Configuration

The frontend is configured using environment variables in `.env.local`.

Key configuration:
- API URL
- Authentication settings
- Feature flags

## API Documentation

The API is documented using Swagger. When the backend is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/
```

## Authentication

The application uses PASETO tokens for authentication. The authentication flow includes:

1. User registration
2. Email verification
3. Login (with rate limiting)
4. OAuth login (Google, GitHub)
5. Password reset
6. Session management

## Email Service

The application includes a flexible email service that supports multiple providers:

### Resend

[Resend](https://resend.com/) is the default email provider. It offers a modern API for sending transactional emails.

### Upstash Workflow

[Upstash Workflow](https://upstash.com/docs/workflow/overall/getstarted) is supported as an alternative email provider. It's particularly useful for scheduling emails and handling complex email workflows.

To configure the email service, update the following environment variables in `.env`:

```env
# Email (Resend)
RESEND_API_KEY=your-resend-api-key
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
EMAIL_FROM_NAME=Your App Name
EMAIL_VERIFICATION_URL=http://localhost:3000/verify-email
EMAIL_PASSWORD_RESET_URL=http://localhost:3000/reset-password
EMAIL_LOGIN_NOTIFICATION=true

# Upstash Workflow (for email workflows)
UPSTASH_WORKFLOW_URL=your-upstash-workflow-url
UPSTASH_WORKFLOW_TOKEN=your-upstash-workflow-token
```

The application will automatically choose the appropriate email provider based on the configuration.

## Database

The application uses Turso, a distributed SQLite database, for data storage. Turso provides:

- Global distribution with low latency
- SQLite compatibility
- Serverless operation
- Built-in replication and high availability

The backend connects to Turso using the libsql-client-go driver. Database migrations are managed using the `migrate` tool.

For local development, you'll need to:
1. Install the Turso CLI: https://docs.turso.tech/reference/turso-cli
2. Create a database and get your connection URL and auth token
3. Configure these in your `.env` file

## Deployment

### Production Deployment

For production deployment, consider:

1. Using a managed database service
2. Setting up proper TLS certificates
3. Configuring environment variables for production
4. Setting up monitoring and logging
5. Using a container orchestration system like Kubernetes

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Security

The application includes several security features to protect user data and prevent abuse:

### Security Features

- **Authentication**: PASETO tokens for secure authentication
- **Rate Limiting**: Prevents brute force attacks and API abuse
- **Input Validation**: All user input is validated using strict schemas
- **Password Security**: Passwords are hashed using Argon2id
- **HTTPS**: All communication is encrypted using TLS
- **Security Headers**: HTTP security headers are set by default
- **CSRF Protection**: Cross-Site Request Forgery protection
- **XSS Protection**: Cross-Site Scripting protection
- **SQL Injection Protection**: Parameterized queries prevent SQL injection

### Security Scanning

The application includes automated security scanning tools:

```bash
# Run security checks on the backend
cd backend
make security-check
```

This command runs:
1. **gosec**: Static analysis tool for Go code that finds security issues
2. **nancy**: Checks dependencies for known vulnerabilities

The security check script automatically installs the required tools if they're not already available.

## Testing

The application includes comprehensive tests for both backend and frontend.

### Backend Tests

Run backend tests with:

```bash
cd backend
make test
```

Generate test coverage report with:

```bash
cd backend
make test-coverage
```

### Frontend Tests

Run frontend tests with:

```bash
cd frontend
npm test
```

## Additional Documentation

For more detailed information, please refer to the following documentation:

- [Backend Documentation](backend/README.md) - Detailed information about the backend service
- [Frontend Documentation](frontend/README.md) - Detailed information about the frontend application
- [Docker Compose Setup](docker-compose.md) - Detailed information about the Docker Compose configuration
- [Security Documentation](SECURITY.md) - Detailed information about security features and best practices

## Recent Changes

### March 1, 2024

1. **Database Migration**:
   - Migrated from PostgreSQL to Turso database (distributed SQLite)
   - Updated database connection code to use libsql-client-go
   - Removed PostgreSQL and PgAdmin services from Docker Compose

2. **Docker Setup**:
   - Updated Docker Compose configuration to reflect the use of Turso database
   - Created documentation for Docker setup issues and fixes (see docker-compose.md)
   - Added note to Docker Setup section about ongoing updates

3. **Documentation**:
   - Updated README to reflect the use of Turso database
   - Added instructions for setting up Turso for local development
   - Updated Features section to mention Turso instead of SQLite

For more details on the Docker setup issues and fixes, see the [docker-compose.md](docker-compose.md) file. 