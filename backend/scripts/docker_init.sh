#!/bin/sh

# Create data directory if it doesn't exist
mkdir -p /app/data

# Generate PASETO keys if they don't exist
echo "Checking for PASETO keys..."
mkdir -p /app/keys

# Try to generate keys using Go script if available
if [ -f /app/scripts/generate_keys.go ] && command -v go >/dev/null 2>&1; then
    echo "Generating PASETO keys using Go script..."
    cd /app && go run scripts/generate_keys.go
    # Copy generated keys to the keys directory if needed
    if [ -f /app/.env.keys ]; then
        # Extract keys from .env.keys and save them to the proper files
        grep "PASETO_PUBLIC_KEY" /app/.env.keys | cut -d '=' -f2- > /app/keys/public_key.pem
        grep "PASETO_PRIVATE_KEY" /app/.env.keys | cut -d '=' -f2- > /app/keys/private_key.pem
        chmod 600 /app/keys/private_key.pem
        chmod 644 /app/keys/public_key.pem
    fi
else
    echo "Warning: Cannot run go command to generate keys."
    echo "Checking for existing keys or environment variables..."
    
    # Check if keys are provided via environment variables
    if [ -n "$PASETO_PUBLIC_KEY" ] && [ -n "$PASETO_PRIVATE_KEY" ]; then
        echo "Using keys from environment variables..."
        echo "$PASETO_PUBLIC_KEY" > /app/keys/public_key.pem
        echo "$PASETO_PRIVATE_KEY" > /app/keys/private_key.pem
        chmod 600 /app/keys/private_key.pem
        chmod 644 /app/keys/public_key.pem
    elif [ ! -f /app/keys/public_key.pem ] || [ ! -f /app/keys/private_key.pem ]; then
        echo "WARNING: No keys found. Generating temporary keys for development only."
        echo "DO NOT USE THESE KEYS IN PRODUCTION!"
        
        # Generate random keys for development only
        openssl genrsa -out /app/keys/private_key.pem 2048 2>/dev/null
        openssl rsa -in /app/keys/private_key.pem -pubout -out /app/keys/public_key.pem 2>/dev/null
        chmod 600 /app/keys/private_key.pem
        chmod 644 /app/keys/public_key.pem
    else
        echo "Using existing keys found in /app/keys directory."
    fi
fi

# Set environment variables for keys if not already set
if [ -z "$PASETO_PUBLIC_KEY_PATH" ]; then
    export PASETO_PUBLIC_KEY_PATH="/app/keys/public_key.pem"
fi

if [ -z "$PASETO_PRIVATE_KEY_PATH" ]; then
    export PASETO_PRIVATE_KEY_PATH="/app/keys/private_key.pem"
fi

# Debug information (don't show sensitive data)
echo "Environment:"
echo "DATABASE_TYPE: $DATABASE_TYPE"
echo "DATABASE_URL: [REDACTED]"
echo "DATABASE_AUTH_TOKEN: [REDACTED]"
echo "PASETO_PUBLIC_KEY_PATH: $PASETO_PUBLIC_KEY_PATH"
echo "PASETO_PRIVATE_KEY_PATH: $PASETO_PRIVATE_KEY_PATH"

# Skip migrations for certain database types
if [ "$DATABASE_TYPE" = "sqlite" ] || [ "$DATABASE_TYPE" = "sqlite3" ] || [ "$DATABASE_TYPE" = "memory" ]; then
    echo "Skipping migrations for database type: $DATABASE_TYPE"
else
    # Run migrations
    echo "Running database migrations..."
    if command -v migrate >/dev/null 2>&1; then
        migrate -path /app/migrations -database "$DATABASE_URL" up
    else
        echo "Warning: migrate command not found, skipping migrations"
    fi
fi

# Start the application
echo "Starting application..."
exec /app/main 