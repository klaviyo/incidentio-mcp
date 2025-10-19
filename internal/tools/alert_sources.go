package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListAlertSourcesTool lists alert sources from incident.io
type ListAlertSourcesTool struct {
	client *incidentio.Client
}

func NewListAlertSourcesTool(client *incidentio.Client) *ListAlertSourcesTool {
	return &ListAlertSourcesTool{client: client}
}

func (t *ListAlertSourcesTool) Name() string {
	return "list_alert_sources"
}

func (t *ListAlertSourcesTool) Description() string {
	return `List available alert sources that can receive and process alert events.

USAGE WORKFLOW:
1. Call to see all configured alert sources
2. Use alert source IDs when creating alert events with create_alert_event

PARAMETERS:
- page_size: Number of results per page (1-250)
- after: Pagination cursor for next page

EXAMPLES:
- List all sources: {}
- List with pagination: {"page_size": 50, "after": "cursor_abc"}

IMPORTANT: Alert source IDs from this tool are required for the create_alert_event tool.`
}

func (t *ListAlertSourcesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page",
				"minimum":     1,
				"maximum":     250,
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor for next page",
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListAlertSourcesTool) Execute(args map[string]interface{}) (string, error) {
	params := &incidentio.ListAlertSourcesParams{}

	if pageSize, ok := args["page_size"].(float64); ok {
		params.PageSize = int(pageSize)
	}
	if after, ok := args["after"].(string); ok {
		params.After = after
	}

	result, err := t.client.ListAlertSources(params)
	if err != nil {
		return "", fmt.Errorf("failed to list alert sources: %w", err)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}
