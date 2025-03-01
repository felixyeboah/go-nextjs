#!/bin/sh

# Create data directory if it doesn't exist
mkdir -p /app/data

# Generate PASETO keys if they don't exist
echo "Generating PASETO keys..."
if [ -f /app/scripts/generate_keys.go ] && command -v go >/dev/null 2>&1; then
    cd /app && go run scripts/generate_keys.go
else
    echo "Warning: Cannot run go command to generate keys. Using default keys."
    # Create default keys if they don't exist
    mkdir -p /app/keys
    if [ ! -f /app/keys/public_key.pem ]; then
        echo "Creating default public key..."
        echo "-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAuTrsCa+y/l54FQaAAiMPRmCXqKBT7PvSgZSCxyYg0ic=
-----END PUBLIC KEY-----" > /app/keys/public_key.pem
    fi
    if [ ! -f /app/keys/private_key.pem ]; then
        echo "Creating default private key..."
        echo "-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIAVVxlx9EjeE9Tx1sRdwK4by/31yXzdcW/0s1i6/XYZt
-----END PRIVATE KEY-----" > /app/keys/private_key.pem
    fi
fi

# Debug environment variables
echo "DEBUG: Environment variables:"
echo "DATABASE_URL: $DATABASE_URL"
echo "DATABASE_AUTH_TOKEN: $DATABASE_AUTH_TOKEN"
env | grep DATABASE

# Skip migrations for Turso database
if echo "$DATABASE_URL" | grep -q "turso\|libsql"; then
    echo "Skipping migrations for Turso database..."
else
    echo "Running migrations..."
    migrate -path /app/migrations -database "$DATABASE_URL" up
fi

# Start the application
echo "Starting the application..."
exec /app/api 