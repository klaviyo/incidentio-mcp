package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// CloseIncidentTool closes an incident by setting it to "Closed" status
type CloseIncidentTool struct {
	client *incidentio.Client
}

func NewCloseIncidentTool(client *incidentio.Client) *CloseIncidentTool {
	return &CloseIncidentTool{client: client}
}

func (t *CloseIncidentTool) Name() string {
	return "close_incident"
}

func (t *CloseIncidentTool) Description() string {
	return `Close an incident by transitioning it to 'Closed' status with automatic workflow handling.

USAGE WORKFLOW:
1. Get incident ID from list_incidents or get_incident
2. Call this tool with the incident ID
3. Tool checks if already closed to avoid errors
4. Attempts direct closure or provides guidance if workflow restrictions apply

PARAMETERS:
- id: Required. The incident ID to close

EXAMPLES:
- Close incident: {"id": "01HXYZ..."}

IMPORTANT: incident.io may require incidents to go through specific status transitions before closing (e.g., Triage → Active → Monitoring → Closed). This tool attempts direct closure and provides helpful guidance if workflow restrictions prevent it.`
}

func (t *CloseIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to close",
			},
		},
		"required": []string{"id"},
	}
}

func (t *CloseIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok {
		return "", fmt.Errorf("id parameter is required")
	}

	// Get the current incident first
	incident, err := t.client.GetIncident(id)
	if err != nil {
		return "", fmt.Errorf("failed to get incident: %w", err)
	}

	// Check if it's already closed
	if incident.IncidentStatus.Category == "closed" {
		return fmt.Sprintf("Incident %s (%s) is already closed with status: %s",
			incident.ID, incident.Name, incident.IncidentStatus.Name), nil
	}

	// Try to close the incident using the update API
	// incident.io has workflow restrictions, so we might need to go through intermediate steps
	closedStatusID := "01JAR1BCBHSK633DVJSFC16RPY"

	req := &incidentio.UpdateIncidentRequest{
		IncidentStatusID: closedStatusID,
	}

	updatedIncident, err := t.client.UpdateIncident(id, req)
	if err != nil {
		// If direct closure fails, provide helpful guidance
		return fmt.Sprintf(`Failed to close incident directly: %v

This might be due to workflow restrictions. incident.io often requires incidents to go through specific states before closing.

Current status: %s (%s)
Suggested workflow:
1. First move to "Monitoring" status if fixing is complete
2. Then move to "Closed" status

You can also close manually:
- Incident page: %s
- Slack channel: %s

Use the update_incident tool with incident_status_id: %s`,
			err,
			incident.IncidentStatus.Name,
			incident.IncidentStatus.Category,
			incident.Permalink,
			incident.SlackChannelName,
			closedStatusID), nil
	}

	// Success! Return the updated incident
	result, err := json.MarshalIndent(map[string]interface{}{
		"message": fmt.Sprintf("Successfully updated incident %s to status: %s",
			updatedIncident.Name, updatedIncident.IncidentStatus.Name),
		"incident": updatedIncident,
	}, "", "  ")

	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
