package tools

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentsTool lists incidents from incident.io
type ListIncidentsTool struct {
	client *incidentio.Client
}

func NewListIncidentsTool(client *incidentio.Client) *ListIncidentsTool {
	return &ListIncidentsTool{client: client}
}

func (t *ListIncidentsTool) Name() string {
	return "list_incidents"
}

func (t *ListIncidentsTool) Description() string {
	return `List incidents from incident.io with optional filtering by status category and severity.

IMPORTANT: Use this tool to get BROAD OVERVIEW information about many incidents (returns default essential fields only).
Once you identify specific incidents of interest, use get_incident tool with the incident_id to retrieve COMPLETE DETAILS about individual incidents.

RECOMMENDED WORKFLOW:
1. Use list_incidents to discover and filter incidents (returns: id, reference, name, created_at, updated_at, slack_channel_id)
2. Identify incidents of interest from the list
3. Use get_incident with specific incident_id values to get full details (all fields) for those incidents
4. This two-step approach minimizes context usage while providing comprehensive information when needed

USAGE:
1. Filter incidents using status categories (validated against your org's available categories)
2. Filter by severity using names like "Critical", "sev_1", or full IDs (automatically mapped)
3. Multiple values can be provided to match any of them (OR logic)
4. Default fields provide essential overview; use 'fields' parameter only if you need different fields for the list
5. For manual pagination, use 'after' parameter with the value from pagination_meta.after in previous response

PARAMETERS:
- page_size: Number of results (default 25, max 250). Set to 0 or omit for auto-pagination.
- after: The incident ID to start pagination after. Use the exact value from pagination_meta.after in previous response.
- status: Status values in array OR comma-separated string format. Accepts friendly aliases OR direct API categories:
  * Format: Array ["active", "triage"] OR comma-separated string "active,triage,learning"
  * Aliases: "active" → "live", "open" → "live", "resolved" → "closed", "completed" → "closed"
  * API categories: live, triage, learning, closed, merged, declined, canceled, paused (varies by org)
  * Case-insensitive matching for both aliases and categories
  * Tool validates against your org's exact incident.io configuration
  * Invalid values return helpful error with all available options and aliases
  * Examples: ["active"], ["live"], ["triage", "active"], "active,triage,learning"
- severity: Severity values in array OR comma-separated string format. Tool automatically maps names to IDs:
  * Format: Array ["Critical", "High"] OR comma-separated string "Critical,High,Medium"
  * By name: "Critical", "High", "Medium", "Low", "sev_1", "sev_2", etc.
  * By ID: "01K56QEGAD95K9K5ZQ9CCPF6EF" (full UUID format)
  * Invalid severities will return helpful error with all available options
  * Examples: ["Critical"], ["sev_1", "sev_2"], "Critical,High"
- fields: Comma-separated list of fields to include in response (reduces context usage)
  * Top-level: "id,name,summary,reference"
  * Nested: "severity.name,incident_status.category,incident_type.name"
  * Default: "id,reference,name,created_at,updated_at,slack_channel_id"
  * Omit or leave empty to use default fields

VALIDATION:
- Status categories are validated against your org's incident.io configuration
- Severity names are validated and automatically mapped to IDs
- Both validations fetch live data from the API to ensure accuracy
- Invalid values return helpful errors listing all available options

PAGINATION:
- Auto-pagination: Omit page_size or set to 0 to fetch all results automatically
  * Returns total_record_count = number of incidents fetched
- Manual pagination:
  1. First request: {"page_size": 25}
  2. Response includes pagination_meta.total_record_count (total matching incidents)
  3. Extract pagination_meta.after from response (ID for next page)
  4. Next request: {"page_size": 25, "after": "<value from pagination_meta.after>"}
  5. Repeat until pagination_meta.after is empty (no more pages)
- NOTE: total_record_count shows the total number of incidents matching your filters.

EXAMPLES:
- List all active incidents (uses default fields): {"status": ["active"]} or {"status": "active"}
- List critical incidents: {"severity": ["Critical"]} or {"severity": "Critical"}
- List active high-severity incidents: {"status": ["active"], "severity": ["Critical", "High"]}
- List triaging and active (array): {"status": ["triage", "active"]}
- List triaging and active (string): {"status": "triage,active,learning"}
- List closed incidents: {"status": ["closed"]} or {"status": "closed"}
- Comma-separated severities: {"severity": "Critical,High,Medium"}
- List with custom fields: {"status": "active", "fields": "id,name,severity.name,incident_status.category"}
- Manual pagination: {"page_size": 10, "after": "01K7RPHSXGPM1V07NPW8V6J6RZ"}

NOTE: Both status and severity are validated against live API data. If you receive an error about invalid values, the error message will list all available options for your organization.`
}

func (t *ListIncidentsTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250, default 25). Set to 0 or omit for automatic pagination through all results.",
				"default":     25,
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to start pagination after. IMPORTANT: Use the EXACT value from pagination_meta.after field in the previous response (e.g., \"01K7RPHSXGPM1V07NPW8V6J6RZ\"). This tells the API to return incidents after this ID. Only used with manual pagination when page_size > 0.",
			},
			"status": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by incident status. Accepts BOTH array format [\"active\", \"triage\"] AND comma-separated string \"active,triage,learning\". Accepts aliases (\"active\" → \"live\", \"resolved\" → \"closed\") OR direct categories (live, triage, learning, closed, merged, declined, canceled, paused). Case-insensitive. Validated against your org's configuration. Invalid values return helpful errors with available options and aliases. Multiple values match any of them (OR logic). Examples: [\"active\"], [\"live\"], [\"triage\", \"active\"], \"active,triage,learning\"",
			},
			"severity": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Filter by severity. Accepts BOTH array format [\"Critical\", \"High\"] AND comma-separated string \"Critical,High,Medium\". Accepts severity names (\"Critical\", \"High\", \"sev_1\", etc.) AND full IDs. Tool automatically maps names to IDs. Multiple values will match any of them (OR logic). Examples: [\"Critical\"], [\"sev_1\", \"sev_2\"], [\"Critical\", \"High\"], \"Critical,High\"",
			},
			"fields": map[string]interface{}{
				"type":        "string",
				"description": GetIncidentFieldsDescription(),
				"default":     "id,reference,name,created_at,updated_at,slack_channel_id",
			},
		},
	}
}

func (t *ListIncidentsTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListIncidentsOptions{}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	if after, ok := args["after"].(string); ok {
		opts.After = after
	}

	// Handle status parameter - supports both array and comma-separated string
	var statusInputs []string
	if statuses, ok := args["status"].([]interface{}); ok {
		// Array format: ["active", "triage", "learning"]
		for _, s := range statuses {
			if str, ok := s.(string); ok {
				statusInputs = append(statusInputs, str)
			}
		}
	} else if statusStr, ok := args["status"].(string); ok {
		// Comma-separated string format: "active,triage,learning"
		for _, s := range strings.Split(statusStr, ",") {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				statusInputs = append(statusInputs, trimmed)
			}
		}
	}

	// Validate status categories against API
	if len(statusInputs) > 0 {
		validatedStatuses, err := t.validateStatusCategories(statusInputs)
		if err != nil {
			return "", fmt.Errorf("failed to validate status categories: %w", err)
		}
		opts.Status = validatedStatuses
	}

	// Handle severity parameter - supports both array and comma-separated string
	var severityInputs []string
	if severities, ok := args["severity"].([]interface{}); ok {
		// Array format: ["Critical", "High"]
		for _, s := range severities {
			if str, ok := s.(string); ok {
				severityInputs = append(severityInputs, str)
			}
		}
	} else if severityStr, ok := args["severity"].(string); ok {
		// Comma-separated string format: "Critical,High"
		for _, s := range strings.Split(severityStr, ",") {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				severityInputs = append(severityInputs, trimmed)
			}
		}
	}

	// Map severity names to IDs
	if len(severityInputs) > 0 {
		mappedSeverities, err := t.mapSeveritiesToIDs(severityInputs)
		if err != nil {
			return "", fmt.Errorf("failed to map severities: %w", err)
		}
		opts.Severity = mappedSeverities
	}

	resp, err := t.client.ListIncidents(opts)
	if err != nil {
		return "", err
	}

	// Apply field filtering with default fields if not specified
	fieldsStr, ok := args["fields"].(string)
	if !ok || fieldsStr == "" {
		fieldsStr = "id,reference,name,created_at,updated_at,slack_channel_id"
	}
	return FilterFields(resp, fieldsStr)
}

// validateStatusCategories validates status categories against API and uses exact API values
func (t *ListIncidentsTool) validateStatusCategories(inputs []string) ([]string, error) {
	// Fetch all incident statuses to get valid categories
	statuses, err := t.client.ListIncidentStatuses()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch incident statuses for validation: %w", err)
	}

	// Build map of lowercase category to actual category value from API
	// This preserves the exact case/format the API uses
	categoryMap := make(map[string]string)
	for _, status := range statuses.IncidentStatuses {
		categoryLower := strings.ToLower(status.Category)
		// Store the actual category value from API response
		categoryMap[categoryLower] = status.Category
	}

	// Define common aliases that map to actual API categories
	aliasMap := map[string]string{
		"active":      "live",
		"open":        "live",
		"ongoing":     "live",
		"in_progress": "live",
		"resolved":    "closed",
		"completed":   "closed",
		"done":        "closed",
	}

	// Validate each input and normalize to API format
	var result []string
	for _, input := range inputs {
		inputLower := strings.ToLower(input)

		// First, check if it's an alias
		if aliasTarget, isAlias := aliasMap[inputLower]; isAlias {
			// Verify the alias target exists in this org's categories
			if apiCategory, ok := categoryMap[aliasTarget]; ok {
				result = append(result, apiCategory)
				continue
			}
			// Alias target not available in this org, fall through to error
		}

		// Check if it matches a valid category directly (case-insensitive lookup)
		if apiCategory, ok := categoryMap[inputLower]; ok {
			result = append(result, apiCategory)
			continue
		}

		// If not found, return error with helpful message including aliases
		return nil, fmt.Errorf("status category '%s' not found. Available categories: %s. You can also use aliases: 'active' → 'live', 'resolved' → 'closed'. Call list_incident_statuses to see all status options", input, t.formatAvailableCategories(statuses.IncidentStatuses))
	}

	return result, nil
}

// formatAvailableCategories formats category list for error messages
func (t *ListIncidentsTool) formatAvailableCategories(statuses []incidentio.IncidentStatus) string {
	// Build unique set of categories
	categorySet := make(map[string]bool)
	for _, status := range statuses {
		categorySet[status.Category] = true
	}

	// Convert to sorted list
	var categories []string
	for category := range categorySet {
		categories = append(categories, category)
	}

	return strings.Join(categories, ", ")
}

// mapSeveritiesToIDs fetches the severity list and maps names to IDs
func (t *ListIncidentsTool) mapSeveritiesToIDs(inputs []string) ([]string, error) {
	// Fetch all severities
	severities, err := t.client.ListSeverities()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch severities for mapping: %w", err)
	}

	// Build name-to-ID and ID-to-ID maps
	nameToID := make(map[string]string)
	idToID := make(map[string]string)
	for _, sev := range severities.Severities {
		// Map by name (case-insensitive)
		nameToID[strings.ToLower(sev.Name)] = sev.ID
		// Map by ID (for passthrough)
		idToID[sev.ID] = sev.ID
	}

	// Map each input
	var result []string
	for _, input := range inputs {
		inputLower := strings.ToLower(input)

		// Try as ID first (direct match)
		if id, ok := idToID[input]; ok {
			result = append(result, id)
			continue
		}

		// Try as name (case-insensitive)
		if id, ok := nameToID[inputLower]; ok {
			result = append(result, id)
			continue
		}

		// If not found, return error with helpful message
		return nil, fmt.Errorf("severity '%s' not found. Available severities: %s. Call list_severities to see all options", input, t.formatAvailableSeverities(severities.Severities))
	}

	return result, nil
}

// formatAvailableSeverities formats severity list for error messages
func (t *ListIncidentsTool) formatAvailableSeverities(severities []incidentio.Severity) string {
	var names []string
	for _, sev := range severities {
		names = append(names, fmt.Sprintf("%s (ID: %s)", sev.Name, sev.ID))
	}
	return strings.Join(names, ", ")
}

// GetIncidentTool retrieves a specific incident
type GetIncidentTool struct {
	client *incidentio.Client
}

func NewGetIncidentTool(client *incidentio.Client) *GetIncidentTool {
	return &GetIncidentTool{client: client}
}

func (t *GetIncidentTool) Name() string {
	return "get_incident"
}

func (t *GetIncidentTool) Description() string {
	return `Get COMPLETE, DETAILED information about a specific incident (returns all fields by default).

IMPORTANT: This tool returns ALL incident data and should be used AFTER list_incidents to get full details about specific incidents.

RECOMMENDED WORKFLOW:
1. First use list_incidents to discover incidents (returns only essential fields: id, reference, name, timestamps, slack_channel_id)
2. Identify specific incident(s) of interest from the list
3. Use THIS TOOL (get_incident) with the incident_id to retrieve COMPLETE information including:
   - Full incident details (status, severity, timeline, assignments, custom fields)
   - Related entities (incident type, status details, severity details)
   - All timestamps and metadata
   - Complete incident history and context
4. Optionally use 'fields' parameter to limit response if you only need specific fields

USAGE:
1. Get incident_id from list_incidents results
2. Call this tool with the incident_id for complete details
3. Review comprehensive information including status, severity, timeline, assignments, and custom fields
4. Use 'fields' parameter only if you need to reduce context by selecting specific fields (otherwise returns everything)

PARAMETERS:
- incident_id: Required. The incident ID to retrieve
- fields: Comma-separated list of fields to include in response (reduces context usage)
  * Top-level: "id,name,summary,reference"
  * Nested: "severity.name,incident_status.category,incident_type.name"
  * Omit to return all fields

EXAMPLES:
- Get incident: {"incident_id": "01HXYZ..."}
- Get with selected fields: {"incident_id": "01HXYZ...", "fields": "id,name,severity.name,incident_status.category"}`
}

func (t *GetIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID",
			},
			"fields": map[string]interface{}{
				"type":        "string",
				"description": GetIncidentFieldsDescription(),
			},
		},
		"required":             []interface{}{"incident_id"},
		"additionalProperties": false,
	}
}

func (t *GetIncidentTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["incident_id"].(string)
	if !ok || id == "" {
		argDetails := make(map[string]interface{})
		for key, value := range args {
			argDetails[key] = value
		}
		return "", fmt.Errorf("incident_id parameter is required and must be a non-empty string. Received parameters: %+v", argDetails)
	}

	incident, err := t.client.GetIncident(id)
	if err != nil {
		return "", err
	}

	// Apply field filtering if requested
	fieldsStr, _ := args["fields"].(string)
	return FilterFields(incident, fieldsStr)
}

// CreateIncidentTool creates a new incident
type CreateIncidentTool struct {
	client *incidentio.Client
}

func NewCreateIncidentTool(client *incidentio.Client) *CreateIncidentTool {
	return &CreateIncidentTool{client: client}
}

func (t *CreateIncidentTool) Name() string {
	return "create_incident"
}

func (t *CreateIncidentTool) Description() string {
	return `Create a new incident in incident.io with automatic Slack channel creation.

USAGE WORKFLOW:
1. Prepare incident details (name is required, other fields optional)
2. Optional but recommended: Call list_severities, list_incident_types, and list_incident_statuses to get valid IDs
3. Create incident with desired configuration
4. Tool provides helpful suggestions if required IDs are missing

PARAMETERS:
- name: Required. The incident title/name
- summary: Optional. Detailed incident description
- severity_id: Optional. Severity ID (from list_severities)
- incident_type_id: Optional. Type ID (from list_incident_types)
- incident_status_id: Optional. Status ID (from list_incident_statuses)
- mode: Optional. Incident mode (standard, retrospective, tutorial), default: standard
- visibility: Optional. Visibility (public, private), default: public
- slack_channel_name_override: Optional. Custom Slack channel name

EXAMPLES:
- Minimal incident: {"name": "API outage in production"}
- Full configuration: {"name": "Database unavailable", "severity_id": "01HXYZ...", "incident_type_id": "01HABC...", "incident_status_id": "01HDEF...", "summary": "Primary database not responding"}

IMPORTANT: Tool automatically generates idempotency keys. If severity, type, or status IDs are not provided, helpful error messages suggest using list_severities, list_incident_types, and list_incident_statuses.`
}

func (t *CreateIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "The incident name/title",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "A summary of the incident",
			},
			"status": map[string]interface{}{
				"type":        "string",
				"description": "Initial status (triage, active, resolved, closed)",
				"default":     "triage",
			},
			"severity_id": map[string]interface{}{
				"type":        "string",
				"description": "The severity ID",
			},
			"incident_type_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident type ID",
			},
			"incident_status_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident status ID",
			},
			"mode": map[string]interface{}{
				"type":        "string",
				"description": "The incident mode (standard, retrospective, tutorial)",
				"enum":        []string{"standard", "retrospective", "tutorial"},
				"default":     "standard",
			},
			"visibility": map[string]interface{}{
				"type":        "string",
				"description": "The incident visibility (public, private)",
				"enum":        []string{"public", "private"},
				"default":     "public",
			},
			"slack_channel_name_override": map[string]interface{}{
				"type":        "string",
				"description": "Override the auto-generated Slack channel name",
			},
		},
		"required":             []interface{}{"name"},
		"additionalProperties": false,
	}
}

func (t *CreateIncidentTool) Execute(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("name parameter is required")
	}

	// Generate idempotency key using timestamp and name
	idempotencyKey := fmt.Sprintf("mcp-%d-%s", time.Now().UnixNano(), name)

	req := &incidentio.CreateIncidentRequest{
		IdempotencyKey: idempotencyKey,
		Name:           name,
		Mode:           "standard", // Default to standard mode
		Visibility:     "public",   // Default to public visibility
	}

	if summary, ok := args["summary"].(string); ok {
		req.Summary = summary
	}
	if statusID, ok := args["incident_status_id"].(string); ok {
		req.IncidentStatusID = statusID
	}
	if severityID, ok := args["severity_id"].(string); ok {
		req.SeverityID = severityID
	}
	if typeID, ok := args["incident_type_id"].(string); ok {
		req.IncidentTypeID = typeID
	}
	if mode, ok := args["mode"].(string); ok {
		req.Mode = mode
	}
	if visibility, ok := args["visibility"].(string); ok {
		req.Visibility = visibility
	}
	if slackOverride, ok := args["slack_channel_name_override"].(string); ok {
		req.SlackChannelNameOverride = slackOverride
	}

	// Check if critical fields are missing and provide helpful suggestions
	var suggestions []string

	if req.SeverityID == "" {
		suggestions = append(suggestions, "severity_id is not set. Use list_severities to see available options.")
	}

	if req.IncidentTypeID == "" {
		suggestions = append(suggestions, "incident_type_id is not set. Use list_incident_types to see available options.")
	}

	if req.IncidentStatusID == "" {
		suggestions = append(suggestions, "incident_status_id is not set. Use list_incident_statuses to see available options.")
	}

	incident, err := t.client.CreateIncident(req)
	if err != nil {
		// If the error is related to missing required fields, provide more helpful error message
		errMsg := err.Error()
		if len(suggestions) > 0 && (strings.Contains(errMsg, "severity") || strings.Contains(errMsg, "incident_type") || strings.Contains(errMsg, "incident_status")) {
			return "", fmt.Errorf("%s\n\nSuggestions:\n%s", errMsg, strings.Join(suggestions, "\n"))
		}
		return "", err
	}

	// Include suggestions in successful response if fields were missing
	result, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	if len(suggestions) > 0 {
		return fmt.Sprintf("%s\n\nNote: Incident created with defaults. %s", result, strings.Join(suggestions, " ")), nil
	}

	return string(result), nil
}

// UpdateIncidentTool updates an existing incident
type UpdateIncidentTool struct {
	client *incidentio.Client
}

func NewUpdateIncidentTool(client *incidentio.Client) *UpdateIncidentTool {
	return &UpdateIncidentTool{client: client}
}

func (t *UpdateIncidentTool) Name() string {
	return "update_incident"
}

func (t *UpdateIncidentTool) Description() string {
	return `Update an existing incident's properties (name, summary, status, severity).

USAGE WORKFLOW:
1. Get incident ID from list_incidents or get_incident
2. Prepare updated values for fields you want to change
3. Call this tool with incident ID and new values
4. At least one field must be updated

PARAMETERS:
- incident_id: Required. The incident ID to update
- name: Optional. New incident name
- summary: Optional. New incident summary
- incident_status_id: Optional. New status ID (from list_incident_statuses)
- severity_id: Optional. New severity ID (from list_severities)

EXAMPLES:
- Update status: {"incident_id": "01HXYZ...", "incident_status_id": "status_456"}
- Update severity: {"incident_id": "01HXYZ...", "severity_id": "sev_789"}
- Update multiple fields: {"incident_id": "01HXYZ...", "name": "Updated name", "summary": "Updated summary"}

IMPORTANT: At least one field to update must be provided.`
}

func (t *UpdateIncidentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"incident_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Update the incident name",
			},
			"summary": map[string]interface{}{
				"type":        "string",
				"description": "Update the incident summary",
			},
			"incident_status_id": map[string]interface{}{
				"type":        "string",
				"description": "Update the incident status ID",
			},
			"severity_id": map[string]interface{}{
				"type":        "string",
				"description": "Update the severity ID",
			},
		},
		"required":             []interface{}{"incident_id"},
		"additionalProperties": false,
	}
}

func (t *UpdateIncidentTool) Execute(args map[string]interface{}) (string, error) {

	id, ok := args["incident_id"].(string)
	if !ok || id == "" {
		argDetails := make(map[string]interface{})
		for key, value := range args {
			argDetails[key] = value
		}
		return "", fmt.Errorf("incident_id parameter is required and must be a non-empty string. Received parameters: %+v", argDetails)
	}

	req := &incidentio.UpdateIncidentRequest{}
	hasUpdate := false

	if name, ok := args["name"].(string); ok {
		req.Name = name
		hasUpdate = true
	}
	if summary, ok := args["summary"].(string); ok {
		req.Summary = summary
		hasUpdate = true
	}
	if statusID, ok := args["incident_status_id"].(string); ok {
		req.IncidentStatusID = statusID
		hasUpdate = true
	}
	if severityID, ok := args["severity_id"].(string); ok {
		req.SeverityID = severityID
		hasUpdate = true
	}

	if !hasUpdate {
		return "", fmt.Errorf("at least one field to update must be provided")
	}

	incident, err := t.client.UpdateIncident(id, req)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
