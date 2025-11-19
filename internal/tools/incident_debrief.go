package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// GetIncidentDebriefTool retrieves the incident debrief/post-mortem document
type GetIncidentDebriefTool struct {
	client *incidentio.Client
}

func NewGetIncidentDebriefTool(client *incidentio.Client) *GetIncidentDebriefTool {
	return &GetIncidentDebriefTool{client: client}
}

func (t *GetIncidentDebriefTool) Name() string {
	return "get_incident_debrief"
}

func (t *GetIncidentDebriefTool) Description() string {
	return `Get the debrief/post-mortem document information for a specific incident.

IMPORTANT: This tool checks if an incident has a debrief document and returns the postmortem_document_url if available.
The tool checks multiple locations for the postmortem URL:
1. Top-level postmortem_document_url field
2. retrospective_incident_options.postmortem_document_url (nested object for retrospective incidents)
3. debrief_export_id field (indicates debrief exists but may need export)

IDENTIFIER FORMATS SUPPORTED:
This tool accepts multiple identifier formats for flexible incident lookup:
1. Full incident ID: "01FDAG4SAP5TYPT98WGR2N7" (direct API call)
2. Incident reference: "INC-123" or just "123" (direct API call - most efficient)
3. Slack channel ID: "C123456789" (looks up via list_incidents)
4. Slack channel name: "20251020-aws-outage-ci-impaired" (looks up via list_incidents, case-insensitive)

USAGE WORKFLOW:
1. Get incident identifier from list_incidents results or from Slack/reference
2. Call this tool with the identifier (supports ID, reference, Slack channel ID/name)
3. Tool automatically resolves the identifier to the incident ID
4. Tool checks if the incident has a debrief document
5. If available, returns the incident details including postmortem_document_url

PARAMETERS:
- incident_id: Required. Can be any of these formats:
  * Full incident ID: "01FDAG4SAP5TYPT98WGR2N7"
  * Incident reference: "INC-123" or "123"
  * Slack channel ID: "C123456789"
  * Slack channel name: "20251020-aws-outage-ci-impaired"

EXAMPLES:
- Get by full ID: {"incident_id": "01HXYZ..."}
- Get by reference: {"incident_id": "INC-123"} or {"incident_id": "123"}
- Get by Slack channel ID: {"incident_id": "C123456789"}
- Get by Slack channel name: {"incident_id": "20251020-aws-outage-ci-impaired"}

ERROR HANDLING:
- Returns an error if the incident doesn't have a debrief document yet
- Returns an error if the incident has a debrief but no postmortem_document_url is available
- Returns helpful error messages to guide users

NOTES:
- The postmortem_document_url can be used to download the debrief document
- Check the has_debrief field to verify if a debrief is available
- Use get_incident tool if you need all incident details, not just debrief information`
}

func (t *GetIncidentDebriefTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Incident identifier in any of these formats: full ID (01FDAG4SAP5TYPT98WGR2N7), reference (INC-123 or 123), Slack channel ID (C123456789), or Slack channel name (20251020-aws-outage-ci-impaired). Tool automatically resolves to incident ID.",
			},
		},
		"required":             []interface{}{"incident_id"},
		"additionalProperties": false,
	}
}

func (t *GetIncidentDebriefTool) Execute(args map[string]interface{}) (string, error) {
	identifier, ok := args["incident_id"].(string)
	if !ok || identifier == "" {
		argDetails := make(map[string]interface{})
		for key, value := range args {
			argDetails[key] = value
		}
		return "", fmt.Errorf("incident_id parameter is required and must be a non-empty string. Received parameters: %+v", argDetails)
	}

	// Create a temporary GetIncidentTool to reuse the identifier resolution logic
	getIncidentTool := NewGetIncidentTool(t.client)
	incidentID, err := getIncidentTool.ResolveIncidentIdentifier(identifier)
	if err != nil {
		return "", err
	}

	// Get the incident debrief using the client method
	incident, err := t.client.GetIncidentDebrief(incidentID)
	if err != nil {
		return "", err
	}

	// Determine where the postmortem URL is located
	postmortemURL := incident.PostmortemDocumentURL
	urlLocation := "top_level"

	// Check if URL is in nested retrospective options
	if postmortemURL == "" && incident.RetrospectiveIncidentOptions != nil {
		postmortemURL = incident.RetrospectiveIncidentOptions.PostmortemDocumentURL
		if postmortemURL != "" {
			urlLocation = "retrospective_incident_options"
		}
	}

	// Format the response with relevant debrief information
	response := map[string]interface{}{
		"incident_id":             incident.ID,
		"incident_reference":      incident.Reference,
		"incident_name":           incident.Name,
		"has_debrief":             incident.HasDebrief,
		"postmortem_document_url": postmortemURL,
		"url_location":            urlLocation,
		"permalink":               incident.Permalink,
	}

	// Include debrief_export_id if available
	if incident.DebriefExportID != "" {
		response["debrief_export_id"] = incident.DebriefExportID
	}

	// Include mode if set
	if incident.Mode != "" {
		response["mode"] = incident.Mode
	}

	result, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
