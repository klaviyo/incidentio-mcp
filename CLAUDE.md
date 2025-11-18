# Project Overview
- This is an MCP (Model Context Protocol) server for incident.io API integration
- Written in Go, following standard Go conventions
- Provides tools for incident management, alerts, and workflow automation

# Development Standards
- Use Go 1.21+ features and idioms
- Follow standard Go formatting (gofmt)
- Use meaningful variable and function names
- Handle errors explicitly, don't ignore them

# Testing Commands
- Run tests: `go test ./...`
- Run specific package tests: `go test ./internal/incidentio/...`
- Lint code: `golangci-lint run` (if available)
- Format code: `go fmt ./...`

# Build and Run
- Build server: `go build -o bin/mcp-server ./cmd/mcp-server`
- Run server: `./start-mcp-server.sh`
- Debug mode: `./debug-mcp.sh`

# Project Structure
- `/cmd/` - Entry points for different server configurations
- `/internal/incidentio/` - incident.io API client implementation
- `/internal/tools/` - MCP tool implementations
- `/internal/server/` - MCP server logic

# Environment Variables
- INCIDENT_IO_API_KEY - Required for API authentication
- INCIDENT_IO_BASE_URL - Optional, defaults to https://api.incident.io/v2