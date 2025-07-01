#!/bin/bash
set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load environment variables from .env file if it exists
if [ -f ".env" ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Check if API key is set
if [ -z "$INCIDENT_IO_API_KEY" ]; then
    echo "Error: INCIDENT_IO_API_KEY environment variable is not set" >&2
    echo "Please create a .env file from .env.example and add your API key" >&2
    exit 1
fi

# Ensure the binary exists
if [ ! -f "./bin/mcp-server" ]; then
    echo "Building MCP server..." >&2
    go build -o bin/mcp-server cmd/mcp-server/main.go
fi

# Run the server
exec "$SCRIPT_DIR/bin/mcp-server"