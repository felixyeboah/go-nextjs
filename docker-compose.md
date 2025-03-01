# Docker Compose Setup for Fullstack Application

This document provides instructions for setting up and running the fullstack application using Docker Compose.

## Overview

The Docker Compose configuration includes the following services:

- **Backend**: Go API server
- **Frontend**: Next.js web application
- **Redis**: In-memory data store for caching and session management

**Note on Database**: This application uses Turso (a remote SQLite database service) rather than a containerized database. The Turso connection details must be configured in the backend's `.env` file.

## Prerequisites

- Docker and Docker Compose installed
- Git (to clone the repository)
- A Turso database account and credentials

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/felixyeboah/go-nextjs.git
   cd go-nextjs
   ```

2. Create environment files:
   - Create `.env` files for both backend and frontend based on their respective `.env.example` files
   - Ensure your backend `.env` file includes the Turso database URL and authentication token:
   ```
   DATABASE_URL=libsql://your-database-name.turso.io
   DATABASE_AUTH_TOKEN=your-turso-auth-token
   ```

3. Start the services:
   ```bash
   docker-compose up -d
   ```

## Service Configuration

### Backend

- **Image**: Built from `./backend/Dockerfile`
- **Port**: 8080
- **Environment**: Variables loaded from `./backend/.env`
- **Database**: Uses Turso (remote SQLite database)
- **Dependencies**: Redis

### Frontend

- **Image**: Built from `./frontend/Dockerfile`
- **Port**: 3000
- **Environment**: Variables loaded from `./frontend/.env`
- **Dependencies**: Backend

### Redis

- **Image**: redis:alpine
- **Port**: 6379
- **Volumes**: Persistent data stored in `redis_data` volume
- **Usage**: Used for caching, session management, and rate limiting

## Accessing Services

- Backend API: http://localhost:8080
- Frontend: http://localhost:3000

## Stopping Services

To stop all services:
```bash
docker-compose down
```

To stop and remove volumes (will delete Redis data):
```bash
docker-compose down -v
```

## Troubleshooting

### Database Connection Issues

If the backend cannot connect to the Turso database:

1. Verify your Turso credentials in the backend's `.env` file:
   ```
   DATABASE_URL=libsql://your-database-name.turso.io
   DATABASE_AUTH_TOKEN=your-turso-auth-token
   ```

2. Check backend logs for connection errors:
   ```bash
   docker-compose logs backend
   ```

3. Ensure your IP address is allowed in Turso's connection settings.

### Frontend Cannot Connect to Backend

1. Ensure the backend service is running:
   ```bash
   docker-compose ps backend
   ```

2. Check frontend logs for connection errors:
   ```bash
   docker-compose logs frontend
   ```

3. Verify that the API URL in the frontend's `.env` file is correctly pointing to the backend service.

## Security Notes

- Never commit `.env` files to version control
- Always use strong, unique passwords and API tokens
- Regularly rotate your Turso authentication tokens
- Consider using Docker secrets for production deployments
- Regularly update Docker images to get security patches

# Docker Compose Setup Issues and Fixes

## Current Issues

1. **Backend Container**:
   - The database driver for libsql (Turso) is not found: `error: failed to open database: database driver: unknown driver libsql (forgotten import?)`
   - The main executable is not found: `/app/scripts/docker_init.sh: exec: line 79: /app/main: not found`

2. **Frontend Container**:
   - The Next.js application hasn't been built yet: `[Error: Could not find a production build in the '.next' directory. Try building your app with 'next build' before starting the production server.]`

## Required Fixes

### Backend Container

1. **Fix the libsql driver import**:
   - Ensure the libsql driver is properly imported in the main.go file or wherever the database connection is established.
   - Add `_ "github.com/tursodatabase/libsql-client-go/libsql"` to the imports.

2. **Fix the missing executable**:
   - Update the Dockerfile to properly build the Go application.
   - Ensure the build output is placed in the correct location (/app/main).

### Frontend Container

1. **Build the Next.js application**:
   - Update the Dockerfile to build the Next.js application before starting it.
   - Change the command from `next start` to `next build && next start`.

## Docker Compose Configuration

The docker-compose.yml file has been updated to:
- Remove PostgreSQL and PgAdmin services as the application now uses Turso database.
- Keep only the necessary services: backend, frontend, and Redis.

## Next Steps

1. Fix the backend Dockerfile to properly build the Go application and import the libsql driver.
2. Fix the frontend Dockerfile to build the Next.js application before starting it.
3. Test the Docker Compose setup again after making these changes.

## Alternative Approach

Until the Docker setup is fixed, users can run the application locally using the Development Setup instructions in the README. 