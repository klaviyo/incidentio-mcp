#!/bin/bash
set -e
cd /Users/tomwentworth/incidentio-mcp-golang
export INCIDENT_IO_API_KEY=***REMOVED***

# Ensure the binary exists
if [ ! -f "./bin/mcp-server-clean" ]; then
    echo "Building MCP server..." >&2
    go build -o bin/mcp-server-clean cmd/mcp-server-clean/main.go
fi

# Run the server
exec ./bin/mcp-server-clean