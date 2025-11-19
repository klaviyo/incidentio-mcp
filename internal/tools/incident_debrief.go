package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

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
	return `Get the debrief/post-mortem document information and content for a specific incident.

DOCUMENT CONTENT RETRIEVAL:
When a postmortem_document_url is available (exported debriefs), this tool will automatically
attempt to fetch and include the document content in the response. Supported platforms:
- Google Docs (fetches as plain text via export URL)
- Other URLs (attempts to fetch HTML content)

CRITICAL LIMITATION - INTERNAL DEBRIEFS:
incident.io supports TWO types of post-mortems/debriefs:
1. INTERNAL: Written directly in the incident.io UI (not accessible via API)
2. EXTERNAL: Written in or exported to Confluence, Notion, Google Docs, etc. (accessible via API)

This tool can ONLY retrieve debriefs that have been exported to external platforms.
If has_debrief=true but no postmortem_document_url is returned, the debrief is internal-only.

TO ACCESS INTERNAL DEBRIEFS:
You must export them via the incident.io UI first:
1. Navigate to the incident page in incident.io
2. Go to the "Post-incident" or "Debrief" tab
3. Click "Export" button
4. Choose destination (Confluence, Notion, Google Docs)
5. IMPORTANT: Once exported, the post-mortem will be moved to the destination and
   no longer be editable within incident.io

After export, the postmortem_document_url field will be populated in the API response.

POSTMORTEM URL LOCATIONS:
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
5. If available and exported, returns the incident details including postmortem_document_url
6. If has_debrief=true but no URL, you need to export it via the UI first

PARAMETERS:
- incident_id: Required. Can be any of these formats:
  * Full incident ID: "01FDAG4SAP5TYPT98WGR2N7"
  * Incident reference: "INC-123" or "123"
  * Slack channel ID: "C123456789"
  * Slack channel name: "20251020-aws-outage-ci-impaired"
- include_content: Optional boolean (default: true). Set to false to skip document content fetching

EXAMPLES:
- Get by full ID: {"incident_id": "01HXYZ..."}
- Get by reference: {"incident_id": "INC-123"} or {"incident_id": "123"}
- Get by Slack channel ID: {"incident_id": "C123456789"}
- Get by Slack channel name: {"incident_id": "20251020-aws-outage-ci-impaired"}

ERROR HANDLING:
- Returns an error if the incident doesn't have a debrief document yet
- Returns an error if the incident has a debrief but no postmortem_document_url is available
  (means it's an internal debrief that needs to be exported first)
- Returns helpful error messages with actionable next steps

NOTES:
- The postmortem_document_url can be used to download the debrief document
- Check the has_debrief field to verify if a debrief is available
- Use debug_incident tool to see raw API response and diagnose field availability
- Use get_incident tool if you need all incident details, not just debrief information`
}

func (t *GetIncidentDebriefTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "Incident identifier in any of these formats: full ID (01FDAG4SAP5TYPT98WGR2N7), reference (INC-123 or 123), Slack channel ID (C123456789), or Slack channel name (20251107-machine-zitadel-high-cpu). Tool automatically resolves to incident ID.",
			},
			"include_content": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether to fetch and include the document content (default: true). Set to false to only get metadata.",
				"default":     true,
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

	// Check if we should fetch document content
	includeContent := true
	if includeContentArg, ok := args["include_content"].(bool); ok {
		includeContent = includeContentArg
	}

	// Fetch document content if URL is available and content is requested
	if includeContent && postmortemURL != "" {
		content, fetchErr := fetchDocumentContent(postmortemURL)
		if fetchErr != nil {
			// Don't fail the entire request if content fetch fails
			response["content_fetch_error"] = fetchErr.Error()
			response["content_note"] = "Document URL is available but content could not be fetched automatically"
		} else {
			response["document_content"] = content
			response["content_length"] = len(content)
		}
	}

	result, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// fetchDocumentContent attempts to fetch the content from a document URL
func fetchDocumentContent(docURL string) (string, error) {
	// Try to convert Google Docs URLs to export format
	exportURL := convertToExportURL(docURL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", exportURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent to avoid potential blocking
	req.Header.Set("User-Agent", "incidentio-mcp-server/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("document fetch returned status %d", resp.StatusCode)
	}

	// Read the content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read document content: %w", err)
	}

	return string(content), nil
}

// convertToExportURL converts various document URLs to their export/plain text formats
func convertToExportURL(docURL string) string {
	// Google Docs: convert to plain text export
	// Format: https://docs.google.com/document/d/{id}/edit -> https://docs.google.com/document/d/{id}/export?format=txt
	// Also handles format without /edit suffix
	googleDocsRegex := regexp.MustCompile(`https://docs\.google\.com/document/d/([a-zA-Z0-9_-]+)`)
	if matches := googleDocsRegex.FindStringSubmatch(docURL); len(matches) > 1 {
		docID := matches[1]
		return fmt.Sprintf("https://docs.google.com/document/d/%s/export?format=txt", docID)
	}

	// For other URLs, return as-is and let the HTTP client handle it
	return docURL
}
