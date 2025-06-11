#!/bin/bash

# Check if API key is set
if [ -z "$INCIDENT_IO_API_KEY" ]; then
    # Don't output errors to stderr when running under Claude
    # The server will handle missing API key gracefully
    :
fi

exec /Users/tomwentworth/incidentio-mcp-golang/bin/mcp-server