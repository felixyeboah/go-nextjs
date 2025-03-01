.PHONY: build run stop clean prune help

# Docker compose commands
DC = docker-compose

# Default target
all: help

# Build all containers
build:
	@echo "Building Docker containers..."
	$(DC) build

# Run the application
run:
	@echo "Starting the application..."
	$(DC) up -d

# Run the application with logs
run-logs:
	@echo "Starting the application with logs..."
	$(DC) up

# Stop the application
stop:
	@echo "Stopping the application..."
	$(DC) down

# Clean up containers
clean:
	@echo "Cleaning up containers..."
	$(DC) down --rmi local

# Prune Docker system
prune:
	@echo "Pruning Docker system..."
	docker system prune -af --volumes

# Backend specific commands
backend-build:
	@echo "Building backend container..."
	$(DC) build backend

backend-logs:
	@echo "Showing backend logs..."
	$(DC) logs -f backend

# Frontend specific commands
frontend-build:
	@echo "Building frontend container..."
	$(DC) build frontend

frontend-logs:
	@echo "Showing frontend logs..."
	$(DC) logs -f frontend

# Help command
help:
	@echo "Available commands:"
	@echo "  make build         - Build all Docker containers"
	@echo "  make run           - Start the application in detached mode"
	@echo "  make run-logs      - Start the application with logs"
	@echo "  make stop          - Stop the application"
	@echo "  make clean         - Clean up containers and images"
	@echo "  make prune         - Prune Docker system (remove all unused containers, networks, images)"
	@echo "  make backend-build - Build only the backend container"
	@echo "  make backend-logs  - Show backend logs"
	@echo "  make frontend-build - Build only the frontend container"
	@echo "  make frontend-logs  - Show frontend logs" 