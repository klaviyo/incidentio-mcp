package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentStatusesTool lists available incident statuses
type ListIncidentStatusesTool struct {
	client *incidentio.Client
}

func NewListIncidentStatusesTool(client *incidentio.Client) *ListIncidentStatusesTool {
	return &ListIncidentStatusesTool{client: client}
}

func (t *ListIncidentStatusesTool) Name() string {
	return "list_incident_statuses"
}

func (t *ListIncidentStatusesTool) Description() string {
	return `List all available incident statuses configured in your organization.

USAGE WORKFLOW:
1. Call to see all status options (triage, active, monitoring, closed, etc.)
2. Use status IDs when creating or updating incidents
3. Review categories and names to understand workflow stages

PARAMETERS:
- None required

EXAMPLES:
- List all statuses: {}

IMPORTANT: Status IDs from this tool are required for create_incident and update_incident tools. Each status has a category (triage, active, closed, etc.) that determines workflow stage.`
}

func (t *ListIncidentStatusesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *ListIncidentStatusesTool) Execute(args map[string]interface{}) (string, error) {
	// Use V1 API to get incident statuses
	originalBaseURL := t.client.BaseURL()
	t.client.SetBaseURL("https://api.incident.io/v1")
	defer t.client.SetBaseURL(originalBaseURL)

	respBody, err := t.client.DoRequest("GET", "/incident_statuses", nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch incident statuses: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	result, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
