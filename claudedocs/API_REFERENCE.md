# incident.io API Client Reference

Comprehensive documentation for the incident.io V2 API Go client library.

## Table of Contents

- [Client Initialization](#client-initialization)
- [Incidents API](#incidents-api)
- [Alerts API](#alerts-api)
- [Alert Routes API](#alert-routes-api)
- [Alert Sources & Events API](#alert-sources--events-api)
- [Workflows API](#workflows-api)
- [Actions API](#actions-api)
- [Roles & Users API](#roles--users-api)
- [Severities API](#severities-api)
- [Incident Types API](#incident-types-api)
- [Incident Updates API](#incident-updates-api)
- [Catalog API](#catalog-api)
- [Error Handling](#error-handling)

---

## Client Initialization

### Creating a Client

Location: `internal/incidentio/client.go:26-49`

```go
client, err := incidentio.NewClient()
```

**Environment Variables:**
- `INCIDENT_IO_API_KEY` (required) - Your incident.io API key
- `INCIDENT_IO_BASE_URL` (optional) - Custom API base URL (defaults to `https://api.incident.io/v2`)

**Returns:**
- `*Client` - Configured HTTP client with authentication
- `error` - Error if API key is missing

**Configuration:**
- Timeout: 30 seconds
- TLS: Minimum TLS 1.2
- User-Agent: `incidentio-mcp-server/0.1.0`

---

## Incidents API

Location: `internal/incidentio/incidents.go`

### List Incidents

```go
func (c *Client) ListIncidents(opts *ListIncidentsOptions) (*ListIncidentsResponse, error)
```

**Options:**
```go
type ListIncidentsOptions struct {
    PageSize int      // Max 250, default 25
    Status   []string // Filter by status: triage, active, resolved, closed
    Severity []string // Filter by severity IDs
}
```

**Response:**
```go
type ListIncidentsResponse struct {
    Incidents []Incident `json:"incidents"`
}
```

**Example:**
```go
resp, err := client.ListIncidents(&incidentio.ListIncidentsOptions{
    PageSize: 50,
    Status:   []string{"active", "triage"},
})
```

### Get Incident

```go
func (c *Client) GetIncident(id string) (*Incident, error)
```

**Parameters:**
- `id` - Incident ID (e.g., "INC-123")

**Returns:**
- `*Incident` - Full incident details
- `error` - HTTP 404 if incident not found

### Create Incident

```go
func (c *Client) CreateIncident(req *CreateIncidentRequest) (*Incident, error)
```

**Request:**
```go
type CreateIncidentRequest struct {
    IdempotencyKey           string `json:"idempotency_key"`
    Name                     string `json:"name"`
    Summary                  string `json:"summary,omitempty"`
    IncidentStatusID         string `json:"incident_status_id,omitempty"`
    SeverityID               string `json:"severity_id,omitempty"`
    IncidentTypeID           string `json:"incident_type_id,omitempty"`
    Mode                     string `json:"mode"` // standard, retrospective, tutorial
    Visibility               string `json:"visibility"` // public, private
    SlackChannelNameOverride string `json:"slack_channel_name_override,omitempty"`
}
```

**Required Fields:**
- `IdempotencyKey` - Unique key to prevent duplicate creation
- `Name` - Incident title

**Example:**
```go
incident, err := client.CreateIncident(&incidentio.CreateIncidentRequest{
    IdempotencyKey: "mcp-1234567890-incident-name",
    Name:           "Database performance degradation",
    Summary:        "Users experiencing slow query responses",
    SeverityID:     "severity-01HXYZ...",
    IncidentTypeID: "type-01HXYZ...",
    Mode:           "standard",
    Visibility:     "public",
})
```

### Update Incident

```go
func (c *Client) UpdateIncident(id string, req *UpdateIncidentRequest) (*Incident, error)
```

**Request:**
```go
type UpdateIncidentRequest struct {
    Name             string `json:"name,omitempty"`
    Summary          string `json:"summary,omitempty"`
    IncidentStatusID string `json:"incident_status_id,omitempty"`
    SeverityID       string `json:"severity_id,omitempty"`
}
```

**Example:**
```go
updated, err := client.UpdateIncident("INC-123", &incidentio.UpdateIncidentRequest{
    Summary:    "Issue resolved after restarting database",
    SeverityID: "severity-low",
})
```

### Incident Type

```go
type Incident struct {
    ID             string          `json:"id"`
    Name           string          `json:"name"`
    Summary        string          `json:"summary"`
    Status         IncidentStatus  `json:"incident_status"`
    Severity       Severity        `json:"severity"`
    IncidentType   IncidentType    `json:"incident_type"`
    Mode           string          `json:"mode"`
    Visibility     string          `json:"visibility"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
    SlackChannelID string          `json:"slack_channel_id"`
    Permalink      string          `json:"permalink"`
}
```

---

## Alerts API

Location: `internal/incidentio/alerts.go`

### List Alerts

```go
func (c *Client) ListAlerts(opts *ListAlertsOptions) (*ListAlertsResponse, error)
```

**Options:**
```go
type ListAlertsOptions struct {
    PageSize int      // Max 250
    Status   []string // Filter by status
}
```

### Get Alert

```go
func (c *Client) GetAlert(id string) (*Alert, error)
```

### List Alerts for Incident

```go
func (c *Client) ListAlertsForIncident(incidentID string) (*ListAlertsResponse, error)
```

**Example:**
```go
alerts, err := client.ListAlertsForIncident("INC-123")
```

---

## Alert Routes API

Location: `internal/incidentio/alert_routes.go`

### List Alert Routes

```go
func (c *Client) ListAlertRoutes(opts *ListAlertRoutesOptions) (*ListAlertRoutesResponse, error)
```

### Get Alert Route

```go
func (c *Client) GetAlertRoute(id string) (*AlertRoute, error)
```

### Create Alert Route

```go
func (c *Client) CreateAlertRoute(req *CreateAlertRouteRequest) (*AlertRoute, error)
```

### Update Alert Route

```go
func (c *Client) UpdateAlertRoute(id string, req *UpdateAlertRouteRequest) (*AlertRoute, error)
```

---

## Alert Sources & Events API

Location: `internal/incidentio/alert_sources.go`, `internal/incidentio/alert_events.go`

### List Alert Sources

```go
func (c *Client) ListAlertSources(opts *ListAlertSourcesOptions) (*ListAlertSourcesResponse, error)
```

### Create Alert Event

```go
func (c *Client) CreateAlertEvent(req *CreateAlertEventRequest) (*CreateAlertEventResponse, error)
```

**Request:**
```go
type CreateAlertEventRequest struct {
    AlertSourceID string                 `json:"alert_source_id"`
    Payload       map[string]interface{} `json:"payload"`
}
```

---

## Workflows API

Location: `internal/incidentio/workflows.go`

### List Workflows

```go
func (c *Client) ListWorkflows(opts *ListWorkflowsOptions) (*ListWorkflowsResponse, error)
```

### Get Workflow

```go
func (c *Client) GetWorkflow(id string) (*Workflow, error)
```

### Update Workflow

```go
func (c *Client) UpdateWorkflow(id string, req *UpdateWorkflowRequest) (*Workflow, error)
```

---

## Actions API

Location: `internal/incidentio/actions.go`

### List Actions

```go
func (c *Client) ListActions(opts *ListActionsOptions) (*ListActionsResponse, error)
```

### Get Action

```go
func (c *Client) GetAction(id string) (*Action, error)
```

---

## Roles & Users API

Location: `internal/incidentio/roles.go`

### List Incident Roles

```go
func (c *Client) ListIncidentRoles(opts *ListIncidentRolesOptions) (*ListIncidentRolesResponse, error)
```

### List Users

```go
func (c *Client) ListUsers(opts *ListUsersOptions) (*ListUsersResponse, error)
```

**Options:**
```go
type ListUsersOptions struct {
    PageSize int
    Email    string // Filter by email
}
```

### Assign Incident Role

```go
func (c *Client) AssignIncidentRole(incidentID string, req *AssignIncidentRoleRequest) (*AssignIncidentRoleResponse, error)
```

**Request:**
```go
type AssignIncidentRoleRequest struct {
    RoleID   string `json:"role_id"`
    Assignee struct {
        ID string `json:"id"`
    } `json:"assignee"`
}
```

---

## Severities API

Location: `internal/incidentio/severities.go`

### List Severities

```go
func (c *Client) ListSeverities(opts *ListSeveritiesOptions) (*ListSeveritiesResponse, error)
```

### Get Severity

```go
func (c *Client) GetSeverity(id string) (*Severity, error)
```

**Severity Type:**
```go
type Severity struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Rank        int    `json:"rank"` // Lower rank = higher severity
}
```

---

## Incident Types API

Location: `internal/incidentio/incident_types.go`

### List Incident Types

```go
func (c *Client) ListIncidentTypes(opts *ListIncidentTypesOptions) (*ListIncidentTypesResponse, error)
```

**Response:**
```go
type IncidentType struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

---

## Incident Updates API

Location: `internal/incidentio/incident_updates.go`

### List Incident Updates

```go
func (c *Client) ListIncidentUpdates(incidentID string, opts *ListIncidentUpdatesOptions) (*ListIncidentUpdatesResponse, error)
```

### Get Incident Update

```go
func (c *Client) GetIncidentUpdate(incidentID, updateID string) (*IncidentUpdate, error)
```

### Create Incident Update

```go
func (c *Client) CreateIncidentUpdate(incidentID string, req *CreateIncidentUpdateRequest) (*IncidentUpdate, error)
```

**Request:**
```go
type CreateIncidentUpdateRequest struct {
    Message string `json:"message"`
}
```

### Delete Incident Update

```go
func (c *Client) DeleteIncidentUpdate(incidentID, updateID string) error
```

---

## Catalog API

Location: `internal/incidentio/catalog.go`

### List Catalog Types

```go
func (c *Client) ListCatalogTypes(opts *ListCatalogTypesOptions) (*ListCatalogTypesResponse, error)
```

### List Catalog Entries

```go
func (c *Client) ListCatalogEntries(catalogTypeID string, opts *ListCatalogEntriesOptions) (*ListCatalogEntriesResponse, error)
```

### Update Catalog Entry

```go
func (c *Client) UpdateCatalogEntry(catalogTypeID, entryID string, req *UpdateCatalogEntryRequest) (*CatalogEntry, error)
```

---

## Error Handling

Location: `internal/incidentio/client.go:117-122`

All API methods return errors with the following format:

```go
type ErrorResponse struct {
    Error struct {
        Message string `json:"message"`
        Code    string `json:"code"`
    } `json:"error"`
}
```

**Common Error Patterns:**

```go
// API errors (HTTP 4xx/5xx)
if err != nil {
    // Format: "API error: <message> (HTTP <status>)"
    return err
}

// 404 Not Found
incident, err := client.GetIncident("invalid-id")
// Returns: "API error: Incident not found (HTTP 404)"

// Authentication errors
// Returns: "API error: Unauthorized (HTTP 401)"

// Validation errors
// Returns: "API error: severity_id is required (HTTP 422)"
```

**Best Practices:**

1. Always check for nil errors after API calls
2. Use specific error messages to guide users
3. HTTP 404 errors indicate resource not found
4. HTTP 422 errors indicate validation failures
5. HTTP 401/403 errors indicate authentication/authorization issues

**Example Error Handling:**

```go
incident, err := client.GetIncident(id)
if err != nil {
    if strings.Contains(err.Error(), "404") {
        return fmt.Errorf("incident %s not found", id)
    }
    return fmt.Errorf("failed to get incident: %w", err)
}
```

---

## HTTP Client Configuration

The client uses a custom HTTP transport with:

- **Timeout**: 30 seconds per request
- **TLS**: Minimum version TLS 1.2
- **Headers**:
  - `Authorization: Bearer <api_key>`
  - `Content-Type: application/json`
  - `User-Agent: incidentio-mcp-server/0.1.0`

**Extending the Client:**

```go
// Access base URL
baseURL := client.BaseURL()

// Change base URL (for testing)
client.SetBaseURL("https://test.incident.io/v2")

// Make custom requests
resp, err := client.DoRequest("GET", "/custom/endpoint", nil, nil)
```
