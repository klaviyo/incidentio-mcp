# Deployment Guide

## Prerequisites

1. Go 1.21 or higher installed
2. Valid incident.io API key from https://app.incident.io/settings/api-keys
3. Claude Desktop installed

## Quick Start

### 1. Clone and Build

```bash
git clone https://github.com/incident-io/incidentio-mcp-golang.git
cd incidentio-mcp-golang
go build -o bin/mcp-server ./cmd/mcp-server
```

### 2. Configure Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "incidentio-golang": {
      "command": "/path/to/incidentio-mcp-golang/start-mcp-server.sh",
      "env": {
        "INCIDENT_IO_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

### 3. Restart Claude Desktop

The MCP server will now be available in Claude with access to incident.io tools.

## Environment Variables

- `INCIDENT_IO_API_KEY` (required) - Your incident.io API key
- `INCIDENT_IO_BASE_URL` (optional) - Custom API endpoint (defaults to https://api.incident.io/v2)

## Available Tools

Once deployed, Claude will have access to:

- `list_incidents` - List and filter incidents
- `get_incident` - Get incident details
- `create_incident` - Create new incidents
- `update_incident` - Update incident properties
- `close_incident` - Close an incident
- `list_alerts` - List alerts
- `list_severities` - List severity levels
- `list_incident_statuses` - List incident statuses
- `assign_incident_role` - Assign roles to users
- And more...

## Troubleshooting

1. **Server disconnects immediately**
   - Check that your API key is valid
   - Ensure the binary has execute permissions: `chmod +x bin/mcp-server`

2. **No tools available**
   - Verify the INCIDENT_IO_API_KEY environment variable is set
   - Check Claude Desktop logs for errors

3. **Authentication errors**
   - Generate a new API key from incident.io settings
   - Update the key in your Claude Desktop configuration

## Security Notes

- Never commit your API key to version control
- The `.env` file is gitignored for local development
- API keys should be stored securely in production environments