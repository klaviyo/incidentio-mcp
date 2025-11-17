package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListIncidentTypesTool lists available incident types
type ListIncidentTypesTool struct {
	apiClient *client.Client
}

func NewListIncidentTypesTool(c *client.Client) *ListIncidentTypesTool {
	return &ListIncidentTypesTool{apiClient: c}
}

func (t *ListIncidentTypesTool) Name() string {
	return "list_incident_types"
}

func (t *ListIncidentTypesTool) Description() string {
	return `List all available incident types configured in your organization.

USAGE WORKFLOW:
1. Call to see all incident type options
2. Use incident type IDs when creating incidents with create_incident
3. Review names and descriptions to choose appropriate type

PARAMETERS:
- None required

EXAMPLES:
- List all types: {}

IMPORTANT: Incident type IDs from this tool are required for the create_incident tool. One type is typically marked as default.`
}

func (t *ListIncidentTypesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":                 "object",
		"properties":           map[string]interface{}{},
		"additionalProperties": false,
	}
}

func (t *ListIncidentTypesTool) Execute(args map[string]interface{}) (string, error) {
	result, err := t.apiClient.ListIncidentTypes()
	if err != nil {
		return "", fmt.Errorf("failed to list incident types: %w", err)
	}

	// Format the output to be more readable
	output := fmt.Sprintf("Found %d incident types:\n\n", len(result.IncidentTypes))

	for _, incidentType := range result.IncidentTypes {
		output += fmt.Sprintf("ID: %s\n", incidentType.ID)
		output += fmt.Sprintf("Name: %s\n", incidentType.Name)
		if incidentType.Description != "" {
			output += fmt.Sprintf("Description: %s\n", incidentType.Description)
		}
		if incidentType.IsDefault {
			output += "Default: Yes\n"
		}
		output += "\n"
	}

	// Also return the raw JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}
