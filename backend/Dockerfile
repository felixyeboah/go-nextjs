FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git gcc musl-dev

# Install swag for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Install sqlc for SQL code generation
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install migrate for database migrations
RUN go install -tags 'libsql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy source code
COPY . .

# Create docs directory if it doesn't exist
RUN mkdir -p docs

# Generate Swagger docs
RUN swag init -g cmd/api/main.go -o docs

# Create .env.docker file if it doesn't exist
RUN touch .env.docker

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -tags 'libsql' -o api ./cmd/api

FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env.docker ./.env

# Create docs directory and copy swagger files
COPY --from=builder /app/docs ./docs

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./api"] 