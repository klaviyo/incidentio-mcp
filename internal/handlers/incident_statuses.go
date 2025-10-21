package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListIncidentStatusesTool lists available incident statuses
type ListIncidentStatusesTool struct {
	apiClient *client.Client
}

func NewListIncidentStatusesTool(c *client.Client) *ListIncidentStatusesTool {
	return &ListIncidentStatusesTool{apiClient: c}
}

func (t *ListIncidentStatusesTool) Name() string {
	return "list_incident_statuses"
}

func (t *ListIncidentStatusesTool) Description() string {
	return "List all available incident statuses (useful for updating incident status)"
}

func (t *ListIncidentStatusesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func (t *ListIncidentStatusesTool) Execute(args map[string]interface{}) (string, error) {
	// Use V1 API to get incident statuses
	originalBaseURL := t.apiClient.BaseURL()
	t.apiClient.SetBaseURL("https://api.incident.io/v1")
	defer t.apiClient.SetBaseURL(originalBaseURL)

	respBody, err := t.apiClient.DoRequest("GET", "/incident_statuses", nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch incident statuses: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return FormatJSONResponse(response)
}
