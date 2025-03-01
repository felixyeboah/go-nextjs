# Docker Compose Setup

This document provides detailed information about the Docker Compose configuration for the fullstack application.

## Overview

The Docker Compose setup includes the following services:

- **Backend**: Go API server with hot-reloading for development
- **Frontend**: Next.js application with hot-reloading
- **PostgreSQL**: Database server
- **Redis**: In-memory cache and message broker
- **PgAdmin**: Web-based PostgreSQL administration tool

## Prerequisites

- Docker and Docker Compose installed on your system
- Git to clone the repository

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fullstack.git
   cd fullstack
   ```

2. Create necessary environment files:
   ```bash
   # Backend environment file
   cp backend/.env.example backend/.env.docker
   
   # Frontend environment file
   cp frontend/.env.example frontend/.env
   ```

3. Start the services:
   ```bash
   docker-compose up -d
   ```

4. Access the services:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/
   - PgAdmin: http://localhost:5050 (admin@example.com / admin)

5. Stop the services:
   ```bash
   docker-compose down
   ```

## Service Details

### Backend Service

```yaml
backend:
  container_name: fullstack-backend
  build:
    context: ./backend
    dockerfile: Dockerfile
  ports:
    - "8080:8080"
  volumes:
    - ./backend:/app
  env_file:
    - ./backend/.env.docker
  depends_on:
    - postgres
    - redis
  networks:
    - fullstack-network
```

- **Build Context**: The `./backend` directory
- **Ports**: Maps port 8080 on the host to port 8080 in the container
- **Volumes**: Mounts the backend directory for hot-reloading
- **Environment**: Uses variables from `.env.docker`
- **Dependencies**: Requires PostgreSQL and Redis services

### Frontend Service

```yaml
frontend:
  container_name: fullstack-frontend
  build:
    context: ./frontend
    dockerfile: Dockerfile
  ports:
    - "3000:3000"
  volumes:
    - ./frontend:/app
    - /app/node_modules
  env_file:
    - ./frontend/.env
  depends_on:
    - backend
  networks:
    - fullstack-network
```

- **Build Context**: The `./frontend` directory
- **Ports**: Maps port 3000 on the host to port 3000 in the container
- **Volumes**: Mounts the frontend directory for hot-reloading and preserves node_modules
- **Environment**: Uses variables from `.env`
- **Dependencies**: Requires the backend service

### PostgreSQL Service

```yaml
postgres:
  container_name: fullstack-postgres
  image: postgres:15-alpine
  ports:
    - "5432:5432"
  volumes:
    - postgres_data:/var/lib/postgresql/data
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: fullstack
  networks:
    - fullstack-network
```

- **Image**: Uses the official PostgreSQL 15 Alpine image
- **Ports**: Maps port 5432 on the host to port 5432 in the container
- **Volumes**: Uses a named volume for data persistence
- **Environment**: Sets up a default database with credentials

### Redis Service

```yaml
redis:
  container_name: fullstack-redis
  image: redis:7-alpine
  ports:
    - "6379:6379"
  volumes:
    - redis_data:/data
  networks:
    - fullstack-network
```

- **Image**: Uses the official Redis 7 Alpine image
- **Ports**: Maps port 6379 on the host to port 6379 in the container
- **Volumes**: Uses a named volume for data persistence

### PgAdmin Service

```yaml
pgadmin:
  container_name: fullstack-pgadmin
  image: dpage/pgadmin4
  ports:
    - "5050:80"
  environment:
    PGADMIN_DEFAULT_EMAIL: admin@example.com
    PGADMIN_DEFAULT_PASSWORD: admin
  volumes:
    - pgadmin_data:/var/lib/pgadmin
  depends_on:
    - postgres
  networks:
    - fullstack-network
```

- **Image**: Uses the official pgAdmin4 image
- **Ports**: Maps port 5050 on the host to port 80 in the container
- **Environment**: Sets up default admin credentials
- **Volumes**: Uses a named volume for data persistence
- **Dependencies**: Requires the PostgreSQL service

## Networks and Volumes

```yaml
networks:
  fullstack-network:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
  pgadmin_data:
```

- **Networks**: Creates a bridge network for service communication
- **Volumes**: Creates named volumes for data persistence

## Development Workflow

1. **Making changes to the backend**:
   - Edit files in the `./backend` directory
   - Changes will be automatically detected and the server will reload

2. **Making changes to the frontend**:
   - Edit files in the `./frontend` directory
   - Changes will be automatically detected and the application will reload

3. **Accessing the database**:
   - Use pgAdmin at http://localhost:5050
   - Login with admin@example.com / admin
   - Add a new server with:
     - Name: fullstack
     - Host: postgres
     - Port: 5432
     - Username: postgres
     - Password: postgres

4. **Viewing logs**:
   ```bash
   # View logs for all services
   docker-compose logs -f
   
   # View logs for a specific service
   docker-compose logs -f backend
   ```

5. **Rebuilding services**:
   ```bash
   # Rebuild a specific service
   docker-compose build backend
   
   # Rebuild and restart a service
   docker-compose up -d --build backend
   ```

## Troubleshooting

- **Database connection issues**: Ensure the PostgreSQL service is running and the connection details in the backend environment file are correct.
- **Redis connection issues**: Ensure the Redis service is running and the connection details in the backend environment file are correct.
- **Frontend not connecting to backend**: Check that the `NEXT_PUBLIC_API_URL` in the frontend environment file is set correctly.
- **Container not starting**: Check the logs with `docker-compose logs -f [service_name]` to identify the issue. 