# Docker Compose Setup for Fullstack Application

This document provides instructions for setting up and running the fullstack application using Docker Compose.

## Overview

The Docker Compose configuration includes the following services:

- **Backend**: Go API server
- **Frontend**: Next.js web application
- **PostgreSQL**: Database server
- **Redis**: In-memory data store
- **PgAdmin**: PostgreSQL administration tool

## Prerequisites

- Docker and Docker Compose installed
- Git (to clone the repository)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/felixyeboah/go-nextjs.git
   cd go-nextjs
   ```

2. Create environment files:
   - Create `.env` files for both backend and frontend based on their respective `.env.example` files
   - Create a `.env` file in the root directory for Docker Compose environment variables:
   ```
   # PostgreSQL
   POSTGRES_USER=your_postgres_user
   POSTGRES_PASSWORD=your_secure_password
   POSTGRES_DB=your_database_name
   
   # PgAdmin
   PGADMIN_DEFAULT_EMAIL=your_email@example.com
   PGADMIN_DEFAULT_PASSWORD=your_secure_pgadmin_password
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
- **Dependencies**: PostgreSQL, Redis

### Frontend

- **Image**: Built from `./frontend/Dockerfile`
- **Port**: 3000
- **Environment**: Variables loaded from `./frontend/.env`
- **Dependencies**: Backend

### PostgreSQL

- **Image**: postgres:16-alpine
- **Port**: 5432
- **Environment Variables**:
  - `POSTGRES_USER`: Database username (from environment or default "postgres")
  - `POSTGRES_PASSWORD`: Database password (from environment or default "postgres")
  - `POSTGRES_DB`: Database name (from environment or default "fullstack")
- **Volumes**: Persistent data stored in `postgres_data` volume

### Redis

- **Image**: redis:alpine
- **Port**: 6379
- **Volumes**: Persistent data stored in `redis_data` volume

### PgAdmin

- **Image**: dpage/pgadmin4
- **Port**: 5050
- **Environment Variables**:
  - `PGADMIN_DEFAULT_EMAIL`: Admin email (from environment or default "admin@admin.com")
  - `PGADMIN_DEFAULT_PASSWORD`: Admin password (from environment or default "admin")
- **Dependencies**: PostgreSQL

## Accessing Services

- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- PgAdmin: http://localhost:5050

## Stopping Services

To stop all services:
```bash
docker-compose down
```

To stop and remove volumes (will delete all data):
```bash
docker-compose down -v
```

## Troubleshooting

### Database Connection Issues

If the backend cannot connect to the database:

1. Ensure PostgreSQL service is running:
   ```bash
   docker-compose ps postgres
   ```

2. Check backend logs for connection errors:
   ```bash
   docker-compose logs backend
   ```

3. Verify that the database connection string in the backend's `.env` file matches the PostgreSQL service configuration.

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
- Always use strong, unique passwords for database and admin accounts
- Consider using Docker secrets for production deployments
- Regularly update Docker images to get security patches 