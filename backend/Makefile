.PHONY: all build run test test-coverage clean migrate swagger lint security-check

# Variables
BINARY_NAME=api
BUILD_DIR=build
MIGRATION_DIR=migrations
SWAGGER_DIR=docs
MODULE_NAME=github.com/nanayaw/fullstack

all: clean build

build:
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

run:
	@echo "Running..."
	go run ./cmd/api

test:
	@echo "Running tests..."
	go test -v -race -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

migrate:
	@echo "Running migrations..."
	@if [ -z "$(m)" ]; then \
		echo "Please specify migration name: make migrate m=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(m)

migrate-up:
	@echo "Applying migrations..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) down

swagger:
	@echo "Generating Swagger documentation..."
	./scripts/generate_swagger.sh

docker-swagger:
	@echo "Generating Swagger documentation in Docker container..."
	docker exec fullstack-backend swag init -g /app/cmd/api/main.go -o /app/docs --parseDependency --parseInternal

lint:
	@echo "Running linter..."
	golangci-lint run

security-check:
	@echo "Running security checks..."
	./scripts/security_check.sh

sqlc:
	@echo "Generating SQLC code..."
	./scripts/generate_sqlc.sh

docker-sqlc:
	@echo "Generating SQLC code in Docker container..."
	docker exec fullstack-backend sqlc generate

generate-keys:
	@echo "Generating PASETO keys..."
	@go run scripts/generate_keys.go

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make migrate      - Create a new migration (use: make migrate m=migration_name)"
	@echo "  make migrate-up   - Apply migrations"
	@echo "  make migrate-down - Rollback migrations"
	@echo "  make swagger      - Generate Swagger documentation locally"
	@echo "  make docker-swagger - Generate Swagger documentation in Docker container"
	@echo "  make lint         - Run linter"
	@echo "  make security-check - Run security checks"
	@echo "  make sqlc        - Generate SQLC code locally"
	@echo "  make docker-sqlc - Generate SQLC code in Docker container"
	@echo "  make generate-keys - Generate PASETO keys" 