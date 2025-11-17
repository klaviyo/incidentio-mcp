# Quick Start Guide

Get up and running with the incident.io MCP server quickly.

## Prerequisites

- Go 1.21 or later
- incident.io account with API access
- API key from incident.io dashboard
- (Optional) Claude Desktop or another MCP client

---

## Installation

### Option 1: Direct Binary

```bash
# Clone the repository
git clone https://github.com/incident-io/incidentio-mcp-golang.git
cd incidentio-mcp-golang

# Set up environment
cp .env.example .env
# Edit .env and add your API key: INCIDENT_IO_API_KEY=your-key-here

# Build
go build -o bin/mcp-server ./cmd/mcp-server

# Run
./start-mcp-server.sh
```

### Option 2: Docker

```bash
# Clone the repository
git clone https://github.com/incident-io/incidentio-mcp-golang.git
cd incidentio-mcp-golang

# Build Docker image
docker-compose build

# Run
docker-compose up
```

---

## Configuration

### Claude Desktop Setup

Add to your Claude Desktop configuration:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "incidentio": {
      "command": "/absolute/path/to/incidentio-mcp-golang/start-mcp-server.sh",
      "env": {
        "INCIDENT_IO_API_KEY": "your-api-key-here"
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
      "args": ["-f", "/absolute/path/to/docker-compose.yml", "run", "--rm", "-T", "mcp-server"],
      "env": {
        "INCIDENT_IO_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `INCIDENT_IO_API_KEY` | Yes | Your incident.io API key |
| `INCIDENT_IO_BASE_URL` | No | Custom API endpoint (defaults to production) |
| `MCP_DEBUG` | No | Enable debug logging (set to 1) |

---

## Your First Tool Call

Once configured with Claude Desktop, you can use natural language:

```
"List all active incidents"
```

Claude will use the `list_incidents` tool with appropriate parameters.

### Behind the Scenes

What happens:
1. Claude interprets your request
2. Selects `list_incidents` tool
3. Calls MCP server with parameters: `{"status": ["active"]}`
4. Server calls incident.io API
5. Results returned to Claude
6. Claude presents information naturally

---

## Common Workflows

### 1. Creating Your First Incident

**Natural Language:**
```
"What severity levels are available?"
"What incident types exist?"
"Create a new incident called 'API Gateway Timeout' with high severity"
```

**What Happens:**
1. `list_severities` - Discovers available severity options
2. `list_incident_types` - Discovers incident type options
3. `create_incident` - Creates incident with discovered IDs

**Direct Tool Usage (if building your own MCP client):**

```json
// Step 1: List severities
{
  "tool": "list_severities"
}

// Step 2: Create incident
{
  "tool": "create_incident",
  "arguments": {
    "name": "API Gateway Timeout",
    "summary": "Users experiencing 504 gateway timeout errors",
    "severity_id": "01HXYZ...",
    "incident_type_id": "01HXYZ...",
    "incident_status_id": "01HXYZ..."
  }
}
```

### 2. Managing Existing Incidents

**Natural Language:**
```
"Show me details for incident INC-123"
"Update INC-123 summary to 'Issue resolved after cache flush'"
"Post an update to INC-123 saying 'Investigating root cause'"
```

**Tools Used:**
- `get_incident` - View incident details
- `update_incident` - Modify incident properties
- `create_incident_update` - Post status updates

### 3. Team Collaboration

**Natural Language:**
```
"Find user with email john.doe@company.com"
"What incident roles are available?"
"Assign John Doe as incident lead for INC-123"
```

**Tools Used:**
- `list_users` with email filter
- `list_available_incident_roles`
- `assign_incident_role`

### 4. Alert Management

**Natural Language:**
```
"Show me all alerts for incident INC-123"
"List all alert sources"
"What alert routes are configured?"
```

**Tools Used:**
- `list_alerts_for_incident`
- `list_alert_sources`
- `list_alert_routes`

---

## Discovering Available Options

Before creating incidents, it's helpful to discover what options are available.

### Severities

```
"What severities are available?"
```

**Tool**: `list_severities`

**Response**: List of severities with IDs, names, descriptions, and ranks (lower rank = higher severity)

### Incident Types

```
"What incident types exist?"
```

**Tool**: `list_incident_types`

**Response**: List of incident types with IDs, names, and descriptions

### Incident Statuses

```
"What incident statuses are there?"
```

**Tool**: `list_incident_statuses`

**Response**: List of statuses like triage, active, resolved, closed

### Users

```
"Find user with email jane@company.com"
```

**Tool**: `list_users` with email filter

**Response**: User details including ID for role assignment

### Roles

```
"What incident roles can be assigned?"
```

**Tool**: `list_available_incident_roles`

**Response**: Available roles like Incident Lead, Communications Lead, etc.

---

## Common Patterns

### Pattern 1: Full Incident Lifecycle

```
1. "Create incident: Database performance degradation, high severity"
2. "Assign Alice as incident lead for INC-123"
3. "Post update to INC-123: Investigating slow queries"
4. "Update INC-123 severity to medium"
5. "Post update to INC-123: Issue resolved, monitoring"
6. "Close incident INC-123"
```

### Pattern 2: Alert Investigation

```
1. "List all active alerts"
2. "Show details for alert ALT-456"
3. "List alerts for incident INC-123"
4. "Create incident from these alerts"
```

### Pattern 3: Workflow Management

```
1. "List all workflows"
2. "Show details for workflow WF-789"
3. "Update workflow WF-789 to enable it"
```

### Pattern 4: Catalog Updates

```
1. "List catalog types"
2. "List all services in the catalog"
3. "Update service 'payments-api' with new team information"
```

---

## Tips and Best Practices

### 1. Always Use Correct Parameter Names

❌ **Wrong**: `"id": "INC-123"`
✅ **Right**: `"incident_id": "INC-123"`

The parameter is always `incident_id`, never just `id`.

### 2. Discover Before Creating

Before creating incidents, discover available options:
- List severities first
- List incident types first
- List incident statuses first

This ensures you use valid IDs.

### 3. Use Smart Creation

For easier incident creation, use `create_incident_smart` which accepts human-readable names:

```json
{
  "tool": "create_incident_smart",
  "arguments": {
    "name": "Service Down",
    "severity": "high",
    "incident_type": "outage"
  }
}
```

The tool will automatically look up the correct IDs.

### 4. Check Error Messages

Error messages include helpful suggestions:

```
"Error: severity_id is required.
Suggestion: Use list_severities to see available options."
```

Follow these suggestions to resolve issues quickly.

### 5. Use Filtering

Most list tools support filtering:

```json
{
  "tool": "list_incidents",
  "arguments": {
    "status": ["active", "triage"],
    "page_size": 50
  }
}
```

### 6. Pagination

For large result sets, use `page_size` parameter:

```json
{
  "tool": "list_incidents",
  "arguments": {
    "page_size": 100
  }
}
```

Maximum: 250 results per page.

---

## Troubleshooting

### "Tool not found" Error

**Cause**: Server failed to initialize due to missing API key

**Solution**: Ensure `INCIDENT_IO_API_KEY` is set in environment

```bash
# Check if set
echo $INCIDENT_IO_API_KEY

# Set if missing
export INCIDENT_IO_API_KEY="your-api-key"
```

### "incident_id parameter is required" Error

**Cause**: Missing or incorrectly named parameter

**Solution**: Use `incident_id`, not `id`:

❌ `{"id": "INC-123"}`
✅ `{"incident_id": "INC-123"}`

### "404 Not Found" Error

**Cause**: Incident ID doesn't exist

**Solution**:
1. Verify incident ID is correct
2. Check if incident exists: `list_incidents`
3. Ensure you have access permissions

### "401 Unauthorized" Error

**Cause**: Invalid or missing API key

**Solution**:
1. Verify API key is correct
2. Check API key has proper permissions
3. Regenerate API key if needed

### "422 Validation Error"

**Cause**: Missing required fields or invalid values

**Solution**:
1. Check error message for specific field
2. Use list tools to discover valid values
3. Follow suggestions in error message

### Server Not Starting

**Symptoms**: Claude shows "Server failed to start"

**Debug Steps**:

```bash
# Enable debug mode
export MCP_DEBUG=1
export INCIDENT_IO_DEBUG=1

# Run directly to see errors
./start-mcp-server.sh
```

Check stderr output for specific error messages.

---

## Next Steps

### Learn More

- **[Tools Reference](TOOLS_REFERENCE.md)** - Complete tool documentation
- **[API Reference](API_REFERENCE.md)** - API client library details
- **[Architecture](ARCHITECTURE.md)** - System design and internals

### Extend the Server

- **Add New Tools**: See [Architecture - Extension Points](ARCHITECTURE.md#extension-points)
- **Contribute**: See [CONTRIBUTING.md](../docs/CONTRIBUTING.md)
- **Test**: See [TESTING.md](../docs/TESTING.md)

### Advanced Usage

- **Custom API Endpoints**: Set `INCIDENT_IO_BASE_URL` for testing
- **Debug Logging**: Enable `MCP_DEBUG=1` for troubleshooting
- **Docker Deployment**: See [DEPLOYMENT.md](../docs/DEPLOYMENT.md)

---

## Example: Complete Incident Workflow

Here's a complete example workflow from creation to closure:

### Step 1: Discover Options

```
User: "What severity levels and incident types are available?"
```

**Tools Called**:
- `list_severities`
- `list_incident_types`

### Step 2: Create Incident

```
User: "Create a new incident: 'Payment API returning 500 errors' with high severity"
```

**Tool Called**: `create_incident`
**Result**: Incident INC-123 created

### Step 3: Assign Team

```
User: "Find user alice@company.com and assign her as incident lead for INC-123"
```

**Tools Called**:
- `list_users` (with email filter)
- `assign_incident_role`

### Step 4: Monitor and Update

```
User: "Post update to INC-123: 'Identified database connection pool exhaustion as root cause'"
```

**Tool Called**: `create_incident_update`

### Step 5: Resolve

```
User: "Post update to INC-123: 'Connection pool size increased, errors stopped. Monitoring for 30 minutes.'"
```

**Tool Called**: `create_incident_update`

### Step 6: Close

```
User: "Close incident INC-123"
```

**Tool Called**: `close_incident`

---

## Getting Help

- **Documentation**: Start with [INDEX.md](INDEX.md) for navigation
- **Issues**: Report bugs at [GitHub Issues](https://github.com/incident-io/incidentio-mcp-golang/issues)
- **Community**: incident.io community forums
- **API Docs**: [incident.io API Documentation](https://api-docs.incident.io/)

---

**Ready to start?** Try creating your first incident with Claude or dive deeper into the [Tools Reference](TOOLS_REFERENCE.md)!
