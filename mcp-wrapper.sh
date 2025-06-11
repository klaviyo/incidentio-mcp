#!/bin/bash

# Load environment variables from .env file if it exists
if [ -f ".env" ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Pass through to the MCP server
# The server will handle missing API key gracefully without exiting
exec "$(dirname "$0")/bin/mcp-server"