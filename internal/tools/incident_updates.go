package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentUpdatesTool lists incident updates
type ListIncidentUpdatesTool struct {
	client *incidentio.Client
}

func NewListIncidentUpdatesTool(client *incidentio.Client) *ListIncidentUpdatesTool {
	return &ListIncidentUpdatesTool{client: client}
}

func (t *ListIncidentUpdatesTool) Name() string {
	return "list_incident_updates"
}

func (t *ListIncidentUpdatesTool) Description() string {
	return `List incident updates (status messages and communications posted during incidents).

USAGE WORKFLOW:
1. Call without filter to see all updates across incidents
2. Filter by incident_id to see timeline for specific incident
3. Review updates to understand incident progression

PARAMETERS:
- incident_id: Optional. Filter updates by specific incident ID
- page_size: Number of results (default 25, max 250)

EXAMPLES:
- List all updates: {}
- List for incident: {"incident_id": "01HXYZ..."}
- Paginated list: {"incident_id": "01HXYZ...", "page_size": 50}`
}

func (t *ListIncidentUpdatesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Filter updates by incident ID",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListIncidentUpdatesTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListIncidentUpdatesOptions{}

	if incidentID, ok := args["incident_id"].(string); ok {
		opts.IncidentID = incidentID
	}
	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	resp, err := t.client.ListIncidentUpdates(opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// GetIncidentUpdateTool gets a specific incident update
type GetIncidentUpdateTool struct {
	client *incidentio.Client
}

func NewGetIncidentUpdateTool(client *incidentio.Client) *GetIncidentUpdateTool {
	return &GetIncidentUpdateTool{client: client}
}

func (t *GetIncidentUpdateTool) Name() string {
	return "get_incident_update"
}

func (t *GetIncidentUpdateTool) Description() string {
	return `Get detailed information about a specific incident update.

USAGE WORKFLOW:
1. Get update ID from list_incident_updates
2. Call this tool for complete update details
3. Review message content, author, and timestamp

PARAMETERS:
- id: Required. The incident update ID to retrieve

EXAMPLES:
- Get update: {"id": "update_123"}`
}

func (t *GetIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident update ID",
			},
		},
		"required":             []interface{}{"id"},
		"additionalProperties": false,
	}
}

func (t *GetIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required")
	}

	update, err := t.client.GetIncidentUpdate(id)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// CreateIncidentUpdateTool creates a new incident update
type CreateIncidentUpdateTool struct {
	client *incidentio.Client
}

func NewCreateIncidentUpdateTool(client *incidentio.Client) *CreateIncidentUpdateTool {
	return &CreateIncidentUpdateTool{client: client}
}

func (t *CreateIncidentUpdateTool) Name() string {
	return "create_incident_update"
}

func (t *CreateIncidentUpdateTool) Description() string {
	return `Create a new incident update (status message) to communicate progress during an incident.

USAGE WORKFLOW:
1. Get incident ID from list_incidents or get_incident
2. Compose status message describing current state or actions taken
3. Post update to incident timeline and notifications

PARAMETERS:
- incident_id: Required. The incident ID to post update to
- message: Required. The status message text

EXAMPLES:
- Post update: {"incident_id": "01HXYZ...", "message": "Database failover completed. Services recovering."}
- Brief update: {"incident_id": "01HXYZ...", "message": "Investigating root cause"}`
}

func (t *CreateIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to post the update to",
			},
			"message": map[string]interface{}{
				"type":        "string",
				"description": "The update message to post",
			},
		},
		"required":             []interface{}{"incident_id", "message"},
		"additionalProperties": false,
	}
}

func (t *CreateIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	incidentID, ok := args["incident_id"].(string)
	if !ok || incidentID == "" {
		return "", fmt.Errorf("incident_id parameter is required")
	}

	message, ok := args["message"].(string)
	if !ok || message == "" {
		return "", fmt.Errorf("message parameter is required")
	}

	req := &incidentio.CreateIncidentUpdateRequest{
		IncidentID: incidentID,
		Message:    message,
	}

	update, err := t.client.CreateIncidentUpdate(req)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// DeleteIncidentUpdateTool deletes an incident update
type DeleteIncidentUpdateTool struct {
	client *incidentio.Client
}

func NewDeleteIncidentUpdateTool(client *incidentio.Client) *DeleteIncidentUpdateTool {
	return &DeleteIncidentUpdateTool{client: client}
}

func (t *DeleteIncidentUpdateTool) Name() string {
	return "delete_incident_update"
}

func (t *DeleteIncidentUpdateTool) Description() string {
	return `Delete an incident update (removes it from incident timeline).

USAGE WORKFLOW:
1. Get update ID from list_incident_updates
2. Call this tool to remove the update
3. Update will be deleted from timeline and notifications

PARAMETERS:
- id: Required. The incident update ID to delete

EXAMPLES:
- Delete update: {"id": "update_123"}`
}

func (t *DeleteIncidentUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident update ID to delete",
			},
		},
		"required":             []interface{}{"id"},
		"additionalProperties": false,
	}
}

func (t *DeleteIncidentUpdateTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required")
	}

	if err := t.client.DeleteIncidentUpdate(id); err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully deleted incident update %s", id), nil
}
