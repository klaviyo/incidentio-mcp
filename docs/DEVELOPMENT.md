# Development Guide

This document covers development setup, testing, and contribution guidelines for the incident.io MCP server.

## Prerequisites

- Go 1.21 or higher
- Access to incident.io API (for testing)

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your API key:
   ```bash
   cp .env.example .env
   # Edit .env and add your INCIDENT_IO_API_KEY
   ```

## Building

Build the MCP server:
```bash
go build -o bin/mcp-server ./cmd/mcp-server
```

## Running

### Quick Start
```bash
./start-mcp-server.sh
```

### With Environment Validation
```bash
./start-with-env.sh
```

This script will:
- Load environment variables from `.env`
- Validate that required API keys are set
- Build the server if needed
- Start the MCP server

### Testing API Connection
```bash
./tests/test_api.sh
```

## Testing

### Unit Tests
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/client/...

# Run tests with verbose output
go test -v ./...
```

### Integration Tests
```bash
# Test API endpoints
./tests/test_endpoints.sh
```

## Code Quality

### Formatting
```bash
go fmt ./...
```

### Linting
```bash
# If golangci-lint is available
golangci-lint run
```

## Project Structure

- `/cmd/` - Entry points for different server configurations
- `/internal/client/` - incident.io API client implementation
- `/internal/handlers/` - MCP tool implementations
- `/internal/server/` - MCP server logic
- `/tests/` - Integration test scripts

## Development Standards

- Use Go 1.21+ features and idioms
- Follow standard Go formatting (gofmt)
- Use meaningful variable and function names
- Handle errors explicitly, don't ignore them
- Write tests for new functionality
- Update documentation for API changes

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the development standards
4. Add tests for new functionality
5. Run the full test suite
6. Submit a pull request

## Debugging

For debugging issues:
1. Check that your API key is correctly set
2. Verify API connectivity with `./tests/test_api.sh`
3. Check server logs for detailed error messages
4. Use Go's built-in debugging tools for development