# Configuration Guide

This document covers environment variables, server configuration, and deployment options for the incident.io MCP server.

## Environment Variables

### Required

- **`INCIDENT_IO_API_KEY`** - Your incident.io API key for authentication
  - Obtain from your incident.io account settings
  - Required for all API operations

### Optional

- **`INCIDENT_IO_BASE_URL`** - Base URL for incident.io API
  - Default: `https://api.incident.io/v2`
  - Only change if using a different incident.io instance

## Configuration Files

### `.env` File

Create a `.env` file in the project root for local development:

```bash
# Copy from example
cp .env.example .env

# Edit with your values
INCIDENT_IO_API_KEY=your_api_key_here
# INCIDENT_IO_BASE_URL=https://api.incident.io/v2  # Optional
```

## MCP Client Configuration

### Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "incident-io": {
      "command": "/path/to/your/project/start-mcp-server.sh",
      "env": {
        "INCIDENT_IO_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

### Other MCP Clients

For other MCP clients, configure them to:

1. Execute the `start-mcp-server.sh` script
2. Set the `INCIDENT_IO_API_KEY` environment variable
3. Optionally set `INCIDENT_IO_BASE_URL` if needed

## Server Options

### Standard Server

```bash
./start-mcp-server.sh
```

- Minimal startup validation
- Graceful handling of missing API key

### Validated Server

```bash
./start-with-env.sh
```

- Validates API key is set before starting
- Auto-builds server if binary is missing
- Exits with error if API key is missing

## Deployment Considerations

### Docker

If deploying with Docker, ensure environment variables are passed:

```dockerfile
ENV INCIDENT_IO_API_KEY=${INCIDENT_IO_API_KEY}
```

### System Service

When running as a system service, ensure:

1. Environment variables are properly set
2. Working directory is set to the project root
3. Binary has appropriate permissions

### Security Notes

- **Never commit API keys to version control**
- Use environment variables or secure secret management
- Restrict file permissions on configuration files containing secrets
- Consider using dedicated service accounts with minimal required permissions

## API Rate Limits

The incident.io API has rate limits. The server handles these gracefully, but be aware:

- Plan API usage according to your incident.io plan limits
- Consider caching strategies for frequently accessed data
- Monitor API usage through incident.io dashboard

## Troubleshooting Configuration

### Missing API Key

```
Error: INCIDENT_IO_API_KEY environment variable is not set
```

**Solution**: Set the API key in your `.env` file or environment

### Invalid API Key

```
Authentication failed
```

**Solution**: Verify your API key is correct and has necessary permissions

### Connection Issues

```
Failed to connect to incident.io API
```

**Solutions**:

- Check network connectivity
- Verify `INCIDENT_IO_BASE_URL` if using custom endpoint
- Check firewall settings
