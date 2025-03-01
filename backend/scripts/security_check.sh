#!/bin/bash

# Exit on error
set -e

echo "Running security checks on the codebase..."

# Get GOPATH and set default if not set
GOPATH=${GOPATH:-$HOME/go}
GOBIN=${GOBIN:-$GOPATH/bin}

# Check if gosec is installed
if ! command -v gosec &> /dev/null; then
    echo "gosec is not installed. Installing..."
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

# Check if nancy is installed
if ! command -v nancy &> /dev/null; then
    echo "nancy is not installed. Installing..."
    go install github.com/sonatype-nexus-community/nancy@latest
fi

echo "Running gosec to check for security issues in the code..."
# Use full path to gosec
$GOBIN/gosec -quiet -exclude-dir=vendor ./...

echo "Running nancy to check for vulnerabilities in dependencies..."
# Use full path to nancy
go list -json -deps ./... | $GOBIN/nancy sleuth

echo "Security checks completed." 