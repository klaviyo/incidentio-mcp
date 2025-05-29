#!/bin/bash

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

echo "Testing incident.io API connection..."

# Test API connection with curl
curl -H "Authorization: Bearer $INCIDENT_IO_API_KEY" \
     -H "Content-Type: application/json" \
     https://api.incident.io/v2/incidents?page_size=1

echo -e "\n\nBuilding MCP server..."

# Build the server
go mod download
go mod tidy
go build -o bin/mcp-server cmd/mcp-server/main.go

echo "Build complete! Server binary created at bin/mcp-server"
echo "To run: ./start-with-env.sh"