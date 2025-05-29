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

exec /Users/tomwentworth/incidentio-mcp-golang/bin/mcp-server