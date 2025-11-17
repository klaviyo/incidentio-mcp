# Dockerfile Consolidation

## Summary
Consolidated `http.Dockerfile` functionality into the main `Dockerfile` with a unified image that supports both stdio and HTTP modes.

## Changes Made

### Dockerfile Structure
The `Dockerfile` now creates a single unified image:

- **Unified image** - Supports both stdio and HTTP modes via entrypoint script
  - Includes both `mcp-server` and `mcp-http-proxy` binaries
  - Runtime dependencies: ca-certificates, bash
  - Non-root user setup (incidentio:1001)
  - Entrypoint script for flexible mode switching
  - Exposed port 8080 (only used in HTTP mode)
  - Environment variables: `MCP_TRANSPORT_MODE=http` (default), `MCP_HTTP_PORT=8080`

### Build Commands

**Single build for both modes:**
```bash
docker build -t incidentio-mcp:latest .
```

**Running in different modes:**
```bash
# HTTP mode (default)
docker run --rm -p 8080:8080 -e INCIDENT_IO_API_KEY="key" incidentio-mcp:latest

# stdio mode
docker run --rm -e MCP_TRANSPORT_MODE=stdio -e INCIDENT_IO_API_KEY="key" incidentio-mcp:latest
```

### Benefits

1. **Single Source of Truth**: All Docker builds defined in one file
2. **Better Caching**: Shared builder stage improves build times
3. **Maintainability**: Changes to build process only need to happen in one place
4. **Flexibility**: Users can choose the appropriate target for their needs

### Migration Notes

- The original `http.Dockerfile` can now be deprecated
- The consolidated approach eliminates the need for git cloning during image build
- Build context is now the repository root, allowing direct access to all source files
