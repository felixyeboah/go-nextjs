#!/bin/bash

# Script to generate SQLc code

# Change to the backend directory
cd "$(dirname "$0")/.." || exit

# Check if sqlc is installed
if ! command -v sqlc &> /dev/null; then
    echo "sqlc is not installed. Installing..."
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    
    # Add GOPATH/bin to PATH temporarily if sqlc is not found
    export PATH="$PATH:$(go env GOPATH)/bin"
    
    # Check again if sqlc is installed
    if ! command -v sqlc &> /dev/null; then
        echo "Failed to install sqlc. Please install it manually."
        exit 1
    fi
fi

# Generate SQLc code
echo "Generating SQLc code..."
sqlc generate

echo "SQLc code generated successfully!" 