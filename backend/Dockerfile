FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Install tools
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install -tags 'libsql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Build the application with Turso driver
RUN CGO_ENABLED=1 GOOS=linux go build -tags 'libsql' -o api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates and SQLite for HTTPS and database
RUN apk --no-cache add ca-certificates sqlite

# Install Go for development tools
RUN apk add --no-cache go git

# Install development tools
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy the binary and other necessary files from the builder stage
COPY --from=builder /app/api /app/api
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env.docker ./.env
COPY --from=builder /app/sqlc.yaml .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy SQL queries for sqlc
COPY --from=builder /app/internal/repository/turso/queries ./internal/repository/turso/queries

# Create directories
RUN mkdir -p /app/data /app/keys /app/internal/repository/turso/db

# Make the initialization script executable
RUN chmod +x /app/scripts/docker_init.sh

# Expose the application port
EXPOSE 8080

# Use the initialization script as entrypoint
ENTRYPOINT ["/app/scripts/docker_init.sh"]

# Command to run the application
CMD ["./api"] 