# incident.io MCP Server - Documentation Index

Comprehensive documentation index for the incident.io MCP server project.

## Quick Navigation

| Documentation | Purpose | Audience |
|---------------|---------|----------|
| [API Reference](API_REFERENCE.md) | incident.io API client library documentation | Developers integrating the API client |
| [Tools Reference](TOOLS_REFERENCE.md) | MCP tools documentation and usage patterns | MCP client users and developers |
| [Architecture](ARCHITECTURE.md) | System design and technical architecture | Architects and senior developers |
| [Quick Start](QUICK_START.md) | Get started quickly with common tasks | New users |

---

## Documentation by Role

### For MCP Client Users

**Goal**: Use the MCP server effectively with Claude or other MCP clients

1. **Start Here**: [Quick Start Guide](QUICK_START.md)
2. **Tool Reference**: [Tools Reference](TOOLS_REFERENCE.md)
   - Learn available tools
   - Understand parameters and responses
   - See usage examples
3. **Common Workflows**: [Tools Reference - Usage Patterns](TOOLS_REFERENCE.md#tool-usage-patterns)

**Key Concepts:**
- Tools are accessed through natural language via MCP clients
- Tool names and parameters are consistent across all tools
- Always use `incident_id` parameter (not `id` or `incidentId`)

---

### For Go Developers

**Goal**: Extend the API client or understand implementation

1. **Architecture Overview**: [Architecture Documentation](ARCHITECTURE.md)
   - System layers and components
   - Design patterns used
   - Extension points
2. **API Client**: [API Reference](API_REFERENCE.md)
   - All API methods and types
   - Error handling patterns
   - HTTP client configuration
3. **Code Structure**:
   ```
   internal/incidentio/     → API client implementation
   internal/tools/          → MCP tool implementations
   internal/server/         → MCP protocol handler
   pkg/mcp/                 → MCP type definitions
   ```

**Development Workflow:**
```bash
# Run tests
go test ./...

# Build
go build -o bin/mcp-server ./cmd/mcp-server

# Format
go fmt ./...

# Lint (if available)
golangci-lint run
```

---

### For MCP Tool Developers

**Goal**: Create new MCP tools or modify existing ones

1. **Tool Interface**: [Architecture - Tools Layer](ARCHITECTURE.md#layer-3-tools-layer)
2. **Tool Examples**: `internal/tools/incidents.go`, `internal/tools/alerts.go`
3. **Adding New Tools**: [Architecture - Extension Points](ARCHITECTURE.md#adding-new-tools)

**Tool Development Pattern:**
```go
// 1. Define tool struct
type MyTool struct {
    client *incidentio.Client
}

// 2. Implement Tool interface
func (t *MyTool) Name() string { return "my_tool" }
func (t *MyTool) Description() string { return "..." }
func (t *MyTool) InputSchema() map[string]interface{} { ... }
func (t *MyTool) Execute(args map[string]interface{}) (string, error) { ... }

// 3. Register in server
s.tools["my_tool"] = NewMyTool(client)
```

---

### For System Architects

**Goal**: Understand system design and make architectural decisions

1. **System Overview**: [Architecture - System Overview](ARCHITECTURE.md#system-overview)
2. **Layer Architecture**: [Architecture - Architecture Layers](ARCHITECTURE.md#architecture-layers)
3. **Data Flow**: [Architecture - Data Flow](ARCHITECTURE.md#data-flow)
4. **Design Patterns**: [Architecture - Design Patterns](ARCHITECTURE.md#design-patterns)
5. **Deployment**: [Architecture - Deployment Architecture](ARCHITECTURE.md#deployment-architecture)

**Key Design Decisions:**
- **Layered Architecture**: Clear separation between protocol, tools, and API client
- **Interface-Based Design**: Tool interface enables extensibility and testing
- **Stateless Server**: Each request is independent, no session management
- **stdio Communication**: MCP protocol over stdin/stdout for process isolation

---

## Documentation by Task

### Getting Started

1. [Installation and Setup](../README.md#-quick-start)
2. [Configuration](../docs/CONFIGURATION.md)
3. [First Tool Call](QUICK_START.md#your-first-tool-call)

### Creating Incidents

1. [Discover Available Options](QUICK_START.md#discovering-available-options)
   - List severities
   - List incident types
   - List incident statuses
2. [Create Basic Incident](TOOLS_REFERENCE.md#create_incident)
3. [Create Smart Incident](TOOLS_REFERENCE.md#create_incident_smart) (with auto-ID resolution)
4. [Post Status Updates](TOOLS_REFERENCE.md#create_incident_update)

### Managing Alerts

1. [List Alerts](TOOLS_REFERENCE.md#list_alerts)
2. [View Alert Details](TOOLS_REFERENCE.md#get_alert)
3. [Create Alert Events](TOOLS_REFERENCE.md#create_alert_event)
4. [Configure Alert Routes](TOOLS_REFERENCE.md#create_alert_route)

### Team Management

1. [Find Users](TOOLS_REFERENCE.md#list_users)
2. [View Available Roles](TOOLS_REFERENCE.md#list_available_incident_roles)
3. [Assign Roles](TOOLS_REFERENCE.md#assign_incident_role)

### Workflow Automation

1. [List Workflows](TOOLS_REFERENCE.md#list_workflows)
2. [View Workflow Details](TOOLS_REFERENCE.md#get_workflow)
3. [Update Workflow](TOOLS_REFERENCE.md#update_workflow)

### Catalog Management

1. [View Catalog Types](TOOLS_REFERENCE.md#list_catalog_types)
2. [List Catalog Entries](TOOLS_REFERENCE.md#list_catalog_entries)
3. [Update Catalog Entries](TOOLS_REFERENCE.md#update_catalog_entry)

---

## API Coverage Matrix

| API Category | Supported Operations | Documentation |
|--------------|---------------------|---------------|
| **Incidents** | List, Get, Create, Update | [API Reference](API_REFERENCE.md#incidents-api) |
| **Incident Updates** | List, Get, Create, Delete | [API Reference](API_REFERENCE.md#incident-updates-api) |
| **Incident Types** | List | [API Reference](API_REFERENCE.md#incident-types-api) |
| **Alerts** | List, Get, List for Incident | [API Reference](API_REFERENCE.md#alerts-api) |
| **Alert Routes** | List, Get, Create, Update | [API Reference](API_REFERENCE.md#alert-routes-api) |
| **Alert Sources** | List | [API Reference](API_REFERENCE.md#alert-sources--events-api) |
| **Alert Events** | Create | [API Reference](API_REFERENCE.md#alert-sources--events-api) |
| **Workflows** | List, Get, Update | [API Reference](API_REFERENCE.md#workflows-api) |
| **Actions** | List, Get | [API Reference](API_REFERENCE.md#actions-api) |
| **Roles** | List, Assign | [API Reference](API_REFERENCE.md#roles--users-api) |
| **Users** | List (with email filter) | [API Reference](API_REFERENCE.md#roles--users-api) |
| **Severities** | List, Get | [API Reference](API_REFERENCE.md#severities-api) |
| **Catalog** | List Types, List Entries, Update | [API Reference](API_REFERENCE.md#catalog-api) |

---

## Code Reference Index

### Core Interfaces

| Interface | Location | Purpose |
|-----------|----------|---------|
| `Tool` | `internal/tools/tool.go:3-8` | Base interface for all MCP tools |
| `Client` | `internal/incidentio/client.go:20-24` | HTTP client for incident.io API |
| `Message` | `pkg/mcp/types.go` | MCP protocol message structure |

### Key Implementations

| Component | Location | Description |
|-----------|----------|-------------|
| MCP Server | `cmd/mcp-server/main.go` | Main entry point and server lifecycle |
| Server Handler | `internal/server/server.go` | MCP protocol message handling |
| Tool Registry | `internal/server/server.go:59-118` | Tool registration system |
| HTTP Client | `internal/incidentio/client.go:66-115` | HTTP request/response handling |

### Tool Implementations

| Tool Category | Location | Tools |
|---------------|----------|-------|
| Incident Tools | `internal/tools/incidents.go` | list_incidents, get_incident, create_incident, update_incident |
| Enhanced Tools | `internal/tools/create_incident_enhanced.go` | create_incident_smart |
| Alert Tools | `internal/tools/alerts.go` | list_alerts, get_alert, list_alerts_for_incident |
| Alert Route Tools | `internal/tools/alert_routes.go` | list_alert_routes, get_alert_route, create_alert_route |
| Workflow Tools | `internal/tools/workflows.go` | list_workflows, get_workflow, update_workflow |
| Role Tools | `internal/tools/roles.go` | list_available_incident_roles, list_users, assign_incident_role |
| Config Tools | `internal/tools/severities.go` | list_severities, get_severity |
| Catalog Tools | `internal/tools/catalog.go` | list_catalog_types, list_catalog_entries, update_catalog_entry |

---

## Testing Index

### Unit Tests

| Test File | Coverage |
|-----------|----------|
| `internal/incidentio/client_test.go` | HTTP client functionality |
| `internal/incidentio/incidents_test.go` | Incident API methods |
| `internal/incidentio/workflows_test.go` | Workflow API methods |
| `internal/incidentio/alert_routes_test.go` | Alert route API methods |
| `internal/incidentio/alert_sources_test.go` | Alert source API methods |
| `internal/incidentio/alert_events_test.go` | Alert event API methods |
| `internal/tools/incidents_test.go` | Incident tools |
| `internal/tools/create_incident_enhanced_test.go` | Enhanced incident creation |

### Integration Tests

| Test Script | Purpose |
|-------------|---------|
| `tests/test_api.sh` | Basic API connectivity |
| `tests/test_endpoints.sh` | Endpoint availability |
| `tests/test_mcp.py` | MCP protocol compliance |
| `tests/test_create_incident.py` | Incident creation workflow |
| `tests/test_severities.py` | Severity listing and lookup |

---

## Configuration Reference

### Environment Variables

| Variable | Required | Default | Purpose |
|----------|----------|---------|---------|
| `INCIDENT_IO_API_KEY` | Yes | - | Authentication token |
| `INCIDENT_IO_BASE_URL` | No | `https://api.incident.io/v2` | API endpoint |
| `MCP_DEBUG` | No | - | Enable debug logging |
| `INCIDENT_IO_DEBUG` | No | - | Enable API client debug logging |

### Claude Desktop Configuration

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "incidentio": {
      "command": "/path/to/start-mcp-server.sh",
      "env": {
        "INCIDENT_IO_API_KEY": "your-api-key"
      }
    }
  }
}
```

**Docker Alternative**:

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

---

## Troubleshooting Guide

### Common Issues

| Issue | Solution | Reference |
|-------|----------|-----------|
| 404 errors | Verify incident ID exists | [README - Troubleshooting](../README.md#common-issues) |
| Authentication errors | Check API key validity | [Configuration](../docs/CONFIGURATION.md) |
| Parameter errors | Use `incident_id` not `id` | [README - Troubleshooting](../README.md#common-issues) |
| Tool not found | Verify tool registration | [Architecture - Tool Registration](ARCHITECTURE.md#tool-registration-system) |
| API timeout | Check network and API status | [API Reference - HTTP Client](API_REFERENCE.md#http-client-configuration) |

### Debug Mode

```bash
export MCP_DEBUG=1
export INCIDENT_IO_DEBUG=1
./start-mcp-server.sh
```

Logs to stderr for debugging without breaking MCP protocol.

---

## Related Documentation

### Project Documentation

- [README.md](../README.md) - Project overview and quick start
- [CONTRIBUTING.md](../docs/CONTRIBUTING.md) - Contribution guidelines
- [DEVELOPMENT.md](../docs/DEVELOPMENT.md) - Development setup and workflow
- [TESTING.md](../docs/TESTING.md) - Testing guidelines
- [DEPLOYMENT.md](../docs/DEPLOYMENT.md) - Deployment instructions
- [CONFIGURATION.md](../docs/CONFIGURATION.md) - Configuration reference

### External Resources

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [incident.io API V2 Documentation](https://api-docs.incident.io/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)

---

## Version History

### Current Version: 0.1.0

**Features:**
- Complete incident.io V2 API coverage
- 35+ MCP tools
- MCP protocol 2024-11-05 compliance
- Docker support
- Comprehensive test suite

**Known Limitations:**
- No connection pooling
- No built-in rate limiting
- Sequential request processing only
- No metrics/monitoring

**Roadmap:**
- Concurrent request handling
- Connection pooling
- Built-in rate limiting
- Metrics and observability
- Additional workflow automation tools

---

## Documentation Maintenance

### Documentation Structure

```
claudedocs/
├── INDEX.md              # This file - navigation hub
├── QUICK_START.md        # Getting started guide
├── API_REFERENCE.md      # API client documentation
├── TOOLS_REFERENCE.md    # MCP tools documentation
└── ARCHITECTURE.md       # System architecture
```

### Update Guidelines

When modifying the codebase:

1. **New Tools**: Update TOOLS_REFERENCE.md and QUICK_START.md
2. **API Changes**: Update API_REFERENCE.md
3. **Architecture Changes**: Update ARCHITECTURE.md
4. **New Features**: Update README.md and relevant guides
5. **All Changes**: Update this INDEX.md if navigation changes

### Documentation Standards

- Use code references with `location:line_number` format
- Include practical examples for all features
- Maintain cross-references between documents
- Keep table of contents updated
- Use consistent terminology throughout

---

## Quick Reference Commands

### Build and Run

```bash
# Build
go build -o bin/mcp-server ./cmd/mcp-server

# Run directly
./start-mcp-server.sh

# Run with Docker
docker-compose up

# Run tests
go test ./...

# Run specific test
go test ./internal/tools/...
```

### Development

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter
golangci-lint run

# Check test coverage
go test -cover ./...
```

### Integration Testing

```bash
# API tests
./tests/test_api.sh

# Endpoint tests
./tests/test_endpoints.sh

# Python tests
python tests/test_mcp.py
```

---

## Support and Feedback

- **Issues**: [GitHub Issues](https://github.com/incident-io/incidentio-mcp-golang/issues)
- **Contributing**: See [CONTRIBUTING.md](../docs/CONTRIBUTING.md)
- **Community**: Check the incident.io community forums

---

**Last Updated**: 2025-10-13
**Documentation Version**: 1.0
**Project Version**: 0.1.0
