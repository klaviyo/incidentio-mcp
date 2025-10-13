# incident.io MCP Server

[![CI](https://github.com/incident-io/incidentio-mcp-golang/actions/workflows/ci.yml/badge.svg)](https://github.com/incident-io/incidentio-mcp-golang/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/incident-io/incidentio-mcp-golang)](https://goreportcard.com/report/github.com/incident-io/incidentio-mcp-golang)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://go.dev/dl/)

A GoLang implementation of an MCP (Model Context Protocol) server for incident.io, providing tools to interact with the incident.io V2 API.

> ⚠️ **Fair warning!** ⚠️  
> This repository is largely vibe-coded and unsupported. Built by our CMO and an enterprising Solutions Engineer with questionable coding practices but undeniable enthusiasm. Use at your own risk! 🚀

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/incident-io/incidentio-mcp-golang.git
cd incidentio-mcp-golang

# Set up environment
cp .env.example .env
# Edit .env and add your incident.io API key

# Build and run
go build -o bin/mcp-server ./cmd/mcp-server
./start-mcp-server.sh
```

## 📋 Features

- ✅ Complete incident.io V2 API coverage
- ✅ Workflow automation and management
- ✅ Alert routing and event handling
- ✅ Comprehensive test suite
- ✅ MCP protocol compliant
- ✅ Clean, modular architecture

## 🤖 Using with Claude

Add to your Claude Desktop configuration:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "incidentio": {
      "command": "/path/to/incidentio-mcp-golang/start-mcp-server.sh",
      "env": {
        "INCIDENT_IO_API_KEY": "your-api-key"
      }
    }
  }
}
```
Or, if you'd prefer to run everything in Docker:

```json
{
    "mcpServers": {
      "incidentio": {
        "command": "docker-compose",
        "args": ["-f", "/path/to/docker-compose.yml", "run", "--rm", "-T", "mcp-server"],
        "env": {
          "INCIDENT_IO_API_KEY": "your-api-key"
        }
      }
    }
}
```


## Available Tools

### Incident Management

- `list_incidents` - List incidents with optional filters
- `get_incident` - Get details of a specific incident
- `create_incident` - Create a new incident
- `update_incident` - Update an existing incident
- `close_incident` - Close an incident with proper workflow
- `create_incident_update` - Post status updates to incidents

### Alert Management

- `list_alerts` - List alerts with optional filters
- `get_alert` - Get details of a specific alert
- `list_alerts_for_incident` - List alerts for an incident
- `create_alert_event` - Create an alert event
- `list_alert_routes` - List and manage alert routes

### Workflow & Automation

- `list_workflows` - List available workflows
- `get_workflow` - Get workflow details
- `update_workflow` - Update workflow configuration

### Team & Roles

- `list_users` - List organization users
- `list_available_incident_roles` - List available incident roles
- `assign_incident_role` - Assign roles to users

### Catalog Management

- `list_catalog_types` - List available catalog types
- `list_catalog_entries` - List catalog entries
- `update_catalog_entry` - Update catalog entries

## 📝 Example Usage

```bash
# Through Claude or another MCP client
"Show me all active incidents"
"Create a new incident called 'Database performance degradation' with severity high"
"List alerts for incident INC-123"
"Assign John Doe as incident lead for INC-456"
"Update the Payments service catalog entry with new team information"
```

## 📚 Documentation

### Getting Started
- **[HTTP Server Guide](docs/HTTP_SERVER.md)** - HTTP/SSE server architecture, multi-arch builds, and deployment
- **[Quick Reference](docs/QUICK_REFERENCE.md)** - Common commands and configuration patterns
- **[Development Guide](docs/DEVELOPMENT.md)** - Setup, testing, and contribution guidelines
- **[Configuration Guide](docs/CONFIGURATION.md)** - Environment variables and deployment options

### Advanced Topics
- **[Testing Guide](docs/TESTING.md)** - Testing documentation and best practices
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Deployment instructions and considerations
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to the project
- **[Code of Conduct](docs/CODE_OF_CONDUCT.md)** - Community guidelines and standards

## 🔧 Troubleshooting

### Common Issues

- **404 errors**: Ensure incident IDs are valid and exist in your instance
- **Authentication errors**: Verify your API key is correct and has proper permissions
- **Parameter errors**: All incident-related tools use `incident_id` as the parameter name

### Debug Mode

Enable debug logging by setting environment variables:

```bash
export MCP_DEBUG=1
export INCIDENT_IO_DEBUG=1
./start-mcp-server.sh
```

## 🤝 Contributing

Contributions are welcome! Please see our [Development Guide](docs/DEVELOPMENT.md) for details on setup, testing, and contribution guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with the [Model Context Protocol](https://modelcontextprotocol.io/) specification
- Powered by [incident.io](https://incident.io/) API
- Created with assistance from Claude
