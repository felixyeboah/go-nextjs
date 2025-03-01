#!/bin/bash

# Script to generate Swagger documentation

# Change to the backend directory
cd "$(dirname "$0")/.." || exit

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "swag is not installed. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
    
    # Add GOPATH/bin to PATH temporarily if swag is not found
    export PATH="$PATH:$(go env GOPATH)/bin"
    
    # Check again if swag is installed
    if ! command -v swag &> /dev/null; then
        echo "Failed to install swag. Please install it manually."
        exit 1
    fi
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

echo "Swagger documentation generated successfully!" 