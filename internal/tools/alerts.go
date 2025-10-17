package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListAlertsTool lists alerts from incident.io
type ListAlertsTool struct {
	client *incidentio.Client
}

func NewListAlertsTool(client *incidentio.Client) *ListAlertsTool {
	return &ListAlertsTool{client: client}
}

func (t *ListAlertsTool) Name() string {
	return "list_alerts"
}

func (t *ListAlertsTool) Description() string {
	return `List alerts from incident.io with optional status filtering.

USAGE WORKFLOW:
1. Call without filters to see all alerts
2. Filter by status to see alerts in specific states
3. Use alert IDs with get_alert for detailed information

PARAMETERS:
- page_size: Number of results (default 25, max 250). Set to 0 or omit for auto-pagination.
- status: Array of status values to filter by - Multiple values match any (OR logic)

EXAMPLES:
- List all alerts: {}
- List firing alerts: {"status": ["firing"]}
- List resolved alerts: {"status": ["resolved"]}`
}

func (t *ListAlertsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
			"status": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by alert status",
			},
		},
	}
}

func (t *ListAlertsTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListAlertsOptions{}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	if statuses, ok := args["status"].([]interface{}); ok {
		for _, s := range statuses {
			if str, ok := s.(string); ok {
				opts.Status = append(opts.Status, str)
			}
		}
	}

	resp, err := t.client.ListAlerts(opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// GetAlertTool retrieves a specific alert
type GetAlertTool struct {
	client *incidentio.Client
}

func NewGetAlertTool(client *incidentio.Client) *GetAlertTool {
	return &GetAlertTool{client: client}
}

func (t *GetAlertTool) Name() string {
	return "get_alert"
}

func (t *GetAlertTool) Description() string {
	return `Get detailed information about a specific alert.

USAGE WORKFLOW:
1. Get alert ID from list_alerts or list_alerts_for_incident
2. Call this tool for complete alert details
3. Review status, metadata, and routing information

PARAMETERS:
- id: Required. The alert ID to retrieve

EXAMPLES:
- Get alert: {"id": "alert_123"}`
}

func (t *GetAlertTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The alert ID",
			},
		},
		"required": []string{"id"},
	}
}

func (t *GetAlertTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok {
		return "", fmt.Errorf("id parameter is required")
	}

	alert, err := t.client.GetAlert(id)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(alert, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// ListAlertsForIncidentTool lists alerts for a specific incident
type ListAlertsForIncidentTool struct {
	client *incidentio.Client
}

func NewListAlertsForIncidentTool(client *incidentio.Client) *ListAlertsForIncidentTool {
	return &ListAlertsForIncidentTool{client: client}
}

func (t *ListAlertsForIncidentTool) Name() string {
	return "list_alerts_for_incident"
}

func (t *ListAlertsForIncidentTool) Description() string {
	return `List all alerts that have been associated with a specific incident.

USAGE WORKFLOW:
1. Get incident ID from list_incidents or get_incident
2. Call this tool to see all alerts linked to that incident
3. Review alert details to understand what triggered the incident

PARAMETERS:
- incident_id: Required. The incident ID to list alerts for
- page_size: Number of results (default 25, max 250)

EXAMPLES:
- List alerts for incident: {"incident_id": "01HXYZ..."}
- List with pagination: {"incident_id": "01HXYZ...", "page_size": 50}`
}

func (t *ListAlertsForIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
		},
		"required": []interface{}{"incident_id"},
	}
}

func (t *ListAlertsForIncidentTool) Execute(args map[string]interface{}) (string, error) {
	incidentID, ok := args["incident_id"].(string)
	if !ok || incidentID == "" {
		return "", fmt.Errorf("incident_id parameter is required")
	}

	opts := &incidentio.ListAlertsOptions{}
	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	resp, err := t.client.ListAlertsForIncident(incidentID, opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
