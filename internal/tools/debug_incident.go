package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// DebugIncidentTool returns the raw incident data for debugging
type DebugIncidentTool struct {
	client *incidentio.Client
}

func NewDebugIncidentTool(client *incidentio.Client) *DebugIncidentTool {
	return &DebugIncidentTool{client: client}
}

func (t *DebugIncidentTool) Name() string {
	return "debug_incident"
}

func (t *DebugIncidentTool) Description() string {
	return `DEBUG TOOL: Get complete raw incident data to inspect all fields returned by the API.

This is a diagnostic tool to help understand what fields the incident.io API is actually returning.
Use this when troubleshooting issues with missing fields or unexpected data structures.

IMPORTANT - INTERNAL DEBRIEF LIMITATION:
If has_debrief=true but no postmortem_document_url fields are present, the debrief is stored
internally in incident.io and has NOT been exported. Internal debriefs are NOT accessible via API.

To access internal debriefs, you must export them via the UI first:
1. Visit the incident page in incident.io
2. Go to "Post-incident" or "Debrief" tab
3. Click "Export" to Confluence, Notion, or Google Docs
4. WARNING: Export will MOVE the debrief and make it no longer editable in incident.io

USAGE:
- incident_id: The incident ID, reference (INC-123 or 123), Slack channel ID, or channel name

RETURNS:
- Complete JSON dump of all incident fields
- Diagnostic summary highlighting debrief-related fields
- Field presence/absence analysis
- Actionable guidance if fields are missing

EXAMPLES:
- {"incident_id": "1565"}
- {"incident_id": "INC-1565"}

This tool is useful for:
- Debugging missing postmortem_document_url issues
- Understanding API response structure
- Verifying which optional fields are populated
- Determining if a debrief is internal (needs export) or external (has URL)
- Comparing expected vs actual field names`
}

func (t *DebugIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Incident identifier (ID, reference, Slack channel ID, or channel name)",
			},
		},
		"required":             []interface{}{"incident_id"},
		"additionalProperties": false,
	}
}

func (t *DebugIncidentTool) Execute(args map[string]interface{}) (string, error) {
	identifier, ok := args["incident_id"].(string)
	if !ok || identifier == "" {
		return "", fmt.Errorf("incident_id parameter is required and must be a non-empty string")
	}

	// Resolve identifier to incident ID
	getIncidentTool := NewGetIncidentTool(t.client)
	incidentID, err := getIncidentTool.ResolveIncidentIdentifier(identifier)
	if err != nil {
		return "", err
	}

	// Get the incident
	incident, err := t.client.GetIncident(incidentID)
	if err != nil {
		return "", err
	}

	// Marshal to JSON to see the complete structure
	fullJSON, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal incident: %w", err)
	}

	// Create a diagnostic summary
	diagnostic := map[string]interface{}{
		"incident_id":       incident.ID,
		"incident_reference": incident.Reference,
		"incident_name":     incident.Name,
		"mode":              incident.Mode,
		"has_debrief":       incident.HasDebrief,
		"debrief_fields": map[string]interface{}{
			"postmortem_document_url_present":        incident.PostmortemDocumentURL != "",
			"postmortem_document_url_value":          incident.PostmortemDocumentURL,
			"retrospective_options_present":          incident.RetrospectiveIncidentOptions != nil,
			"debrief_export_id_present":              incident.DebriefExportID != "",
			"debrief_export_id_value":                incident.DebriefExportID,
		},
		"permalink": incident.Permalink,
	}

	// If retrospective options exist, include its details
	if incident.RetrospectiveIncidentOptions != nil {
		diagnostic["retrospective_incident_options"] = map[string]interface{}{
			"external_id":                      incident.RetrospectiveIncidentOptions.ExternalID,
			"postmortem_document_url":          incident.RetrospectiveIncidentOptions.PostmortemDocumentURL,
			"postmortem_document_url_present": incident.RetrospectiveIncidentOptions.PostmortemDocumentURL != "",
			"slack_channel_id":                 incident.RetrospectiveIncidentOptions.SlackChannelID,
		}
	}

	diagnosticJSON, _ := json.MarshalIndent(diagnostic, "", "  ")

	result := fmt.Sprintf(`=== DIAGNOSTIC SUMMARY ===
%s

=== FULL INCIDENT DATA ===
%s

=== ANALYSIS ===
Has Debrief: %v
Postmortem URL at top level: %v
Retrospective Options object: %v
Debrief Export ID: %v

If has_debrief is true but no postmortem_document_url is found:
1. The debrief may only be accessible via the incident.io UI
2. The debrief may require export before a URL is available
3. The field name might be different than expected
4. The debrief might be in a format not exposed via this API endpoint

Consider checking:
- Incident UI: %s
- incident.io API documentation for debrief/postmortem export endpoints
- Whether the debrief was created as an external document (Google Docs, Notion, etc.)`,
		string(diagnosticJSON),
		string(fullJSON),
		incident.HasDebrief,
		incident.PostmortemDocumentURL != "",
		incident.RetrospectiveIncidentOptions != nil,
		incident.DebriefExportID != "",
		incident.Permalink,
	)

	return result, nil
}
