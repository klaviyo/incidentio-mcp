# MCP Tools Reference

Complete reference for all MCP tools provided by the incident.io MCP server.

## Table of Contents

- [Tool Interface](#tool-interface)
- [Incident Management Tools](#incident-management-tools)
- [Incident Update Tools](#incident-update-tools)
- [Alert Management Tools](#alert-management-tools)
- [Alert Routing Tools](#alert-routing-tools)
- [Workflow Tools](#workflow-tools)
- [Action Tools](#action-tools)
- [Role & User Tools](#role--user-tools)
- [Configuration Tools](#configuration-tools)
- [Catalog Tools](#catalog-tools)
- [Tool Usage Patterns](#tool-usage-patterns)

---

## Tool Interface

Location: `internal/tools/tool.go:3-8`

All tools implement the `Tool` interface:

```go
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    Execute(args map[string]interface{}) (string, error)
}
```

**MCP Protocol:**
- Tools are registered at server startup
- Tool list available via `tools/list` method
- Tool execution via `tools/call` method with JSON-RPC 2.0

---

## Incident Management Tools

### list_incidents

Location: `internal/tools/incidents.go:13-86`

List incidents with optional filtering.

**Parameters:**
```json
{
  "page_size": 25,
  "status": ["active", "triage"],
  "severity": ["01HXYZ..."]
}
```

**Example:**
```
List all active incidents with high severity
```

**Response:** JSON array of incident objects

---

### get_incident

Location: `internal/tools/incidents.go:88-140`

Get detailed information about a specific incident.

**Parameters:**
```json
{
  "incident_id": "INC-123"
}
```

**Required:** `incident_id`

**Error Handling:**
- Validates `incident_id` is non-empty string
- Returns detailed error with received parameters if validation fails
- HTTP 404 if incident not found

**Example:**
```
Get details for incident INC-123
```

---

### create_incident

Location: `internal/tools/incidents.go:142-284`

Create a new incident with intelligent defaults and suggestions.

**Parameters:**
```json
{
  "name": "Database performance degradation",
  "summary": "Users experiencing slow queries",
  "severity_id": "01HXYZ...",
  "incident_type_id": "01HXYZ...",
  "incident_status_id": "01HXYZ...",
  "mode": "standard",
  "visibility": "public",
  "slack_channel_name_override": "incident-db-perf"
}
```

**Required:** `name`

**Defaults:**
- `mode`: "standard"
- `visibility`: "public"

**Smart Features:**
- Automatic idempotency key generation: `mcp-{timestamp}-{name}`
- Suggests using `list_severities` if `severity_id` missing
- Suggests using `list_incident_types` if `incident_type_id` missing
- Suggests using `list_incident_statuses` if `incident_status_id` missing
- Includes suggestions in error messages when API validation fails

**Example:**
```
Create a new incident called "API Gateway 500 errors" with high severity
```

---

### create_incident_smart

Location: `internal/tools/create_incident_enhanced.go`

Enhanced incident creation with automatic ID resolution.

**Parameters:**
```json
{
  "name": "Service degradation",
  "summary": "High error rates detected",
  "severity": "high",
  "incident_type": "service_issue",
  "incident_status": "triage"
}
```

**Features:**
- Accepts human-readable names instead of IDs
- Automatically looks up severity, type, and status IDs
- Falls back to defaults if lookups fail
- More user-friendly than `create_incident`

---

### update_incident

Location: `internal/tools/incidents.go:286-379`

Update an existing incident's properties.

**Parameters:**
```json
{
  "incident_id": "INC-123",
  "name": "Updated incident name",
  "summary": "Updated summary",
  "incident_status_id": "01HXYZ...",
  "severity_id": "01HXYZ..."
}
```

**Required:** `incident_id`
**Requires:** At least one field to update

**Validation:**
- Ensures `incident_id` is non-empty
- Validates at least one update field is provided

**Example:**
```
Update incident INC-123 summary to "Issue resolved"
```

---

### close_incident

Location: `internal/tools/close_incident.go`

Close an incident with proper workflow handling.

**Parameters:**
```json
{
  "incident_id": "INC-123"
}
```

**Required:** `incident_id`

**Behavior:**
- Fetches incident statuses
- Identifies "closed" status
- Updates incident to closed state
- Returns updated incident details

**Example:**
```
Close incident INC-456
```

---

### list_incident_statuses

Location: `internal/tools/incident_statuses.go`

List all available incident statuses.

**Parameters:** None

**Response:**
```json
{
  "incident_statuses": [
    {
      "id": "01HXYZ...",
      "name": "Triage",
      "category": "triage"
    }
  ]
}
```

**Usage:** Call before creating/updating incidents to get valid status IDs

---

### list_incident_types

Location: `internal/tools/incident_types.go`

List all configured incident types.

**Parameters:** None

**Response:**
```json
{
  "incident_types": [
    {
      "id": "01HXYZ...",
      "name": "Service Outage",
      "description": "Complete service unavailability"
    }
  ]
}
```

**Usage:** Call before creating incidents to get valid type IDs

---

## Incident Update Tools

### list_incident_updates

Location: `internal/tools/incident_updates.go`

List all status updates for an incident.

**Parameters:**
```json
{
  "incident_id": "INC-123",
  "page_size": 25
}
```

**Required:** `incident_id`

---

### get_incident_update

Get a specific incident update.

**Parameters:**
```json
{
  "incident_id": "INC-123",
  "update_id": "01HXYZ..."
}
```

**Required:** `incident_id`, `update_id`

---

### create_incident_update

Post a status update to an incident.

**Parameters:**
```json
{
  "incident_id": "INC-123",
  "message": "Database queries have returned to normal latency"
}
```

**Required:** `incident_id`, `message`

**Example:**
```
Post update to INC-123: "Root cause identified, deploying fix"
```

---

### delete_incident_update

Remove an incident update.

**Parameters:**
```json
{
  "incident_id": "INC-123",
  "update_id": "01HXYZ..."
}
```

**Required:** `incident_id`, `update_id`

---

## Alert Management Tools

### list_alerts

Location: `internal/tools/alerts.go`

List all alerts with optional filtering.

**Parameters:**
```json
{
  "page_size": 50,
  "status": ["firing"]
}
```

---

### get_alert

Get details of a specific alert.

**Parameters:**
```json
{
  "alert_id": "01HXYZ..."
}
```

**Required:** `alert_id`

---

### list_alerts_for_incident

List all alerts associated with an incident.

**Parameters:**
```json
{
  "incident_id": "INC-123"
}
```

**Required:** `incident_id`

**Example:**
```
Show me all alerts for incident INC-123
```

---

### create_alert_event

Location: `internal/tools/alert_events.go`

Create a new alert event from an alert source.

**Parameters:**
```json
{
  "alert_source_id": "01HXYZ...",
  "payload": {
    "title": "High CPU usage",
    "severity": "warning",
    "details": "CPU usage above 90%"
  }
}
```

**Required:** `alert_source_id`, `payload`

---

### list_alert_sources

Location: `internal/tools/alert_sources.go`

List all configured alert sources.

**Parameters:** None

**Usage:** Get alert source IDs for creating alert events

---

## Alert Routing Tools

### list_alert_routes

Location: `internal/tools/alert_routes.go`

List all alert routing rules.

**Parameters:**
```json
{
  "page_size": 25
}
```

---

### get_alert_route

Get details of a specific alert route.

**Parameters:**
```json
{
  "alert_route_id": "01HXYZ..."
}
```

**Required:** `alert_route_id`

---

### create_alert_route

Create a new alert routing rule.

**Parameters:**
```json
{
  "name": "Database Alerts",
  "enabled": true,
  "condition_groups": []
}
```

**Required:** `name`

---

### update_alert_route

Update an existing alert route.

**Parameters:**
```json
{
  "alert_route_id": "01HXYZ...",
  "name": "Updated Route Name",
  "enabled": false
}
```

**Required:** `alert_route_id`

---

## Workflow Tools

### list_workflows

Location: `internal/tools/workflows.go`

List all configured workflows.

**Parameters:**
```json
{
  "page_size": 25
}
```

---

### get_workflow

Get details of a specific workflow.

**Parameters:**
```json
{
  "workflow_id": "01HXYZ..."
}
```

**Required:** `workflow_id`

---

### update_workflow

Update workflow configuration.

**Parameters:**
```json
{
  "workflow_id": "01HXYZ...",
  "name": "Updated Workflow",
  "enabled": true
}
```

**Required:** `workflow_id`

---

## Action Tools

### list_actions

Location: `internal/tools/actions.go`

List all workflow actions.

**Parameters:**
```json
{
  "page_size": 25,
  "incident_id": "INC-123"
}
```

---

### get_action

Get details of a specific action.

**Parameters:**
```json
{
  "action_id": "01HXYZ..."
}
```

**Required:** `action_id`

---

## Role & User Tools

### list_available_incident_roles

Location: `internal/tools/roles.go`

List all configured incident roles.

**Parameters:** None

**Response:**
```json
{
  "incident_roles": [
    {
      "id": "01HXYZ...",
      "name": "Incident Lead",
      "description": "Leads incident response"
    }
  ]
}
```

---

### list_users

List organization users.

**Parameters:**
```json
{
  "page_size": 50,
  "email": "user@example.com"
}
```

**Features:**
- Filter by email for user lookup
- Pagination support

**Example:**
```
Find user with email john.doe@example.com
```

---

### assign_incident_role

Assign a role to a user for an incident.

**Parameters:**
```json
{
  "incident_id": "INC-123",
  "role_id": "01HXYZ...",
  "user_id": "01HXYZ..."
}
```

**Required:** `incident_id`, `role_id`, `user_id`

**Example:**
```
Assign Jane Smith as incident lead for INC-123
```

---

## Configuration Tools

### list_severities

Location: `internal/tools/severities.go`

List all severity levels.

**Parameters:** None

**Response:**
```json
{
  "severities": [
    {
      "id": "01HXYZ...",
      "name": "Critical",
      "description": "Service completely unavailable",
      "rank": 1
    }
  ]
}
```

**Note:** Lower rank = higher severity

---

### get_severity

Get details of a specific severity.

**Parameters:**
```json
{
  "severity_id": "01HXYZ..."
}
```

**Required:** `severity_id`

---

## Catalog Tools

### list_catalog_types

Location: `internal/tools/catalog.go`

List all catalog types (services, teams, etc.).

**Parameters:** None

**Response:**
```json
{
  "catalog_types": [
    {
      "id": "01HXYZ...",
      "name": "Services",
      "type_name": "Service"
    }
  ]
}
```

---

### list_catalog_entries

List entries for a specific catalog type.

**Parameters:**
```json
{
  "catalog_type_id": "01HXYZ...",
  "page_size": 50
}
```

**Required:** `catalog_type_id`

**Example:**
```
List all services in the catalog
```

---

### update_catalog_entry

Update a catalog entry's attributes.

**Parameters:**
```json
{
  "catalog_type_id": "01HXYZ...",
  "entry_id": "01HXYZ...",
  "attribute_values": {}
}
```

**Required:** `catalog_type_id`, `entry_id`

---

## Tool Usage Patterns

### Discovery Pattern

Before creating incidents, discover available options:

```
1. list_severities → Get severity IDs
2. list_incident_types → Get type IDs
3. list_incident_statuses → Get status IDs
4. create_incident → Create with discovered IDs
```

### Role Assignment Pattern

```
1. list_users with email filter → Find user
2. list_available_incident_roles → Get role options
3. assign_incident_role → Assign user to incident
```

### Alert Routing Pattern

```
1. list_alert_sources → Get source IDs
2. list_alert_routes → Check existing routes
3. create_alert_route → Create new routing rule
```

### Incident Lifecycle Pattern

```
1. create_incident → Create new incident
2. create_incident_update → Post status updates
3. assign_incident_role → Assign team members
4. close_incident → Close when resolved
```

### Error Recovery Pattern

When tool execution fails:

```
1. Check error message for missing IDs
2. Call suggested list_* tools
3. Retry with correct IDs
```

---

## Validation Layer

Location: `internal/tools/validation.go`

Common validation utilities used across tools:

- **String validation**: Non-empty string checks
- **ID validation**: Format and existence validation
- **Parameter validation**: Required field enforcement
- **Error formatting**: Consistent error messages with context

**Pattern:**
```go
if id == "" {
    return "", fmt.Errorf("incident_id parameter is required and must be a non-empty string. Received parameters: %+v", args)
}
```

This provides clear feedback to users about what went wrong and what was received.
