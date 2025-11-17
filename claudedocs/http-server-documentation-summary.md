# HTTP Server Documentation Summary

## Documentation Created

Generated comprehensive documentation for the incident.io MCP HTTP/SSE server implementation.

### Files Created

1. **`docs/HTTP_SERVER.md`** - Comprehensive HTTP server guide (~400 lines)
   - Architecture and design rationale
   - Multi-architecture build instructions with `docker buildx`
   - Runtime configuration and environment variables
   - HTTP/SSE endpoint documentation
   - Protocol flow and state machine
   - Session management and lifecycle
   - Security considerations
   - Debugging and troubleshooting
   - Deployment examples (Docker Compose, Kubernetes, ECS)
   - Performance considerations

2. **`docs/QUICK_REFERENCE.md`** - Quick reference guide (~200 lines)
   - Common build commands
   - Run command examples
   - Environment variables table
   - HTTP endpoint reference
   - Testing procedures
   - Troubleshooting checklist
   - Common deployment patterns
   - Security checklist

3. **`README.md`** - Updated with documentation links
   - Added "Getting Started" section with HTTP Server Guide
   - Added Quick Reference link
   - Reorganized documentation into logical sections

## Key Documentation Topics

### Architecture
- **Proxy Design Pattern**: HTTP proxy translates between HTTP/SSE and stdio MCP server
- **Process Isolation**: Separate `mcp-server` process per session
- **Reference Counting**: Tracks active client connections for cleanup
- **State Machine**: Enforces MCP protocol initialization sequence

### Multi-Architecture Builds

**Single Command Approach:**
```bash
docker buildx build --target http-server \
  --platform linux/amd64,linux/arm64 \
  --tag registry/incidentio-mcp:v0.0.1 --push .
```

**Separate Builds + Manifest:**
```bash
# Build each architecture separately
docker buildx build --platform linux/amd64 --tag registry/image:amd64 --push .
docker buildx build --platform linux/arm64 --tag registry/image:arm64 --push .

# Combine with imagetools
docker buildx imagetools create --tag registry/image:v0.0.1 \
  registry/image:amd64 registry/image:arm64
```

### Runtime Configuration

**Production Example:**
```bash
docker run --rm -p 8080:8080 -u 1001:1001 \
  -e MCP_TRANSPORT_MODE=http \
  -e INCIDENT_IO_API_KEY="$(cat /path/to/secret.txt)" \
  -e MCP_DEBUG=1 \
  -e INCIDENT_IO_DEBUG=1 \
  registry/incidentio-mcp:v0.0.1
```

**Environment Variables:**
- `MCP_TRANSPORT_MODE`: stdio, http, or sse
- `MCP_HTTP_PORT`: HTTP server port (default: 8080)
- `INCIDENT_IO_API_KEY`: API authentication (required)
- `MCP_DEBUG`: Enable MCP protocol debug logs
- `INCIDENT_IO_DEBUG`: Enable API debug logs

### HTTP/SSE Protocol

**Endpoints:**
1. `GET /sse?session=<id>` - Establish SSE stream
2. `POST /message?session=<id>` - Send JSON-RPC messages
3. `GET /health` - Health check

**Protocol Flow:**
1. Client connects to `/sse`, proxy spawns `mcp-server` process
2. Client sends `initialize` request via `/message`
3. Proxy forwards to `mcp-server`, returns response via SSE
4. Client sends `notifications/initialized` notification
5. Protocol state → Ready, all tools available
6. Client sends tool requests, receives responses via SSE

### Security Best Practices

1. **Non-root user**: Run as `1001:1001`
2. **Secret management**: Load API keys from files or secret stores
3. **Network isolation**: Bind to localhost only in development
4. **Resource limits**: Set CPU/memory constraints
5. **TLS/HTTPS**: Use reverse proxy for production

### Deployment Examples

**Docker Compose:**
- Production-ready configuration with secrets
- Health checks and restart policies

**Kubernetes:**
- Security context with non-root user
- Secret references for API keys
- Liveness and readiness probes
- Service exposure

**AWS ECS:**
- Task definition with Secrets Manager integration
- Health checks
- Network configuration

## Usage Patterns

### Development
```bash
# Quick start with debug logging
docker run --rm -p 8080:8080 \
  -e INCIDENT_IO_API_KEY="test-key" \
  -e MCP_DEBUG=1 \
  incidentio-mcp:http-server
```

### Production
```bash
# Full security configuration
docker run --rm -p 8080:8080 -u 1001:1001 \
  --cpus="0.5" --memory="256m" \
  -e INCIDENT_IO_API_KEY="$(cat /secrets/api-key.txt)" \
  registry/incidentio-mcp:v0.0.1
```

### Testing
```bash
# Terminal 1: SSE stream
curl -N http://localhost:8080/sse?session=test

# Terminal 2: Send initialize
curl -X POST http://localhost:8080/message?session=test \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize",...}'
```

## Benefits of Documentation

1. **Comprehensive Coverage**: Architecture through deployment
2. **Multiple Formats**: Detailed guide + quick reference
3. **Practical Examples**: Real-world commands and configurations
4. **Security Focus**: Best practices for production deployments
5. **Troubleshooting**: Common issues and solutions
6. **Multi-Platform**: Covers AMD64 and ARM64 builds

## Next Steps

Users can now:
1. Understand the HTTP server architecture and design
2. Build multi-architecture Docker images
3. Deploy to various platforms (Docker, K8s, ECS)
4. Configure security and resource limits
5. Debug and troubleshoot issues
6. Follow best practices for production

## Documentation Structure

```
docs/
├── HTTP_SERVER.md        # Comprehensive guide (~400 lines)
│   ├── Architecture
│   ├── Multi-arch builds
│   ├── Runtime config
│   ├── Protocol flow
│   ├── Security
│   ├── Debugging
│   └── Deployment examples
│
├── QUICK_REFERENCE.md    # Quick reference (~200 lines)
│   ├── Build commands
│   ├── Run commands
│   ├── Environment variables
│   ├── Testing procedures
│   └── Troubleshooting
│
└── [existing docs]       # Development, testing, etc.

README.md                 # Updated with new doc links
```
