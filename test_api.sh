#!/bin/bash

# Source environment variables
export INCIDENT_IO_API_KEY=***REMOVED***

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
echo "To run: export INCIDENT_IO_API_KEY=$INCIDENT_IO_API_KEY && ./bin/mcp-server"