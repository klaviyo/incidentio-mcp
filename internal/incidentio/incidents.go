package incidentio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListIncidentsOptions represents options for listing incidents
type ListIncidentsOptions struct {
	PageSize        int
	After           string
	Status          []string
	Severity        []string
	CreatedAtGTE    string // Greater than or equal to date filter (ISO 8601 format)
	CreatedAtLTE    string // Less than or equal to date filter (ISO 8601 format)
	CreatedAtRange  string // Date range filter (format: "2024-12-02~2024-12-08")
	UpdatedAtGTE    string // Greater than or equal to date filter (ISO 8601 format)
	UpdatedAtLTE    string // Less than or equal to date filter (ISO 8601 format)
	UpdatedAtRange  string // Date range filter (format: "2024-12-02~2024-12-08")
}

// ListIncidentsResponse represents the response from listing incidents
type ListIncidentsResponse struct {
	Incidents []Incident `json:"incidents"`
	ListResponse
}

// ListIncidents retrieves a list of incidents with automatic pagination
func (c *Client) ListIncidents(opts *ListIncidentsOptions) (*ListIncidentsResponse, error) {
	allIncidents := []Incident{}
	pageSize := 250 // Default max page size
	after := ""

	// If a specific page size is requested, respect it and don't paginate
	if opts != nil && opts.PageSize > 0 {
		params := url.Values{}
		params.Set("page_size", strconv.Itoa(opts.PageSize))

		if opts.After != "" {
			params.Set("after", opts.After)
		}

		for _, status := range opts.Status {
			params.Add("status_category[one_of]", status)
		}
		for _, severity := range opts.Severity {
			params.Add("severity[one_of]", severity)
		}

		// Add date filters for created_at
		if opts.CreatedAtGTE != "" {
			params.Set("created_at[gte]", opts.CreatedAtGTE)
		}
		if opts.CreatedAtLTE != "" {
			params.Set("created_at[lte]", opts.CreatedAtLTE)
		}
		if opts.CreatedAtRange != "" {
			params.Set("created_at[date_range]", opts.CreatedAtRange)
		}

		// Add date filters for updated_at
		if opts.UpdatedAtGTE != "" {
			params.Set("updated_at[gte]", opts.UpdatedAtGTE)
		}
		if opts.UpdatedAtLTE != "" {
			params.Set("updated_at[lte]", opts.UpdatedAtLTE)
		}
		if opts.UpdatedAtRange != "" {
			params.Set("updated_at[date_range]", opts.UpdatedAtRange)
		}

		respBody, err := c.doRequest("GET", "/incidents", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListIncidentsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		// API returns total_record_count for single page requests
		return &response, nil
	}

	// Set up base parameters for auto-pagination
	baseParams := url.Values{}
	if opts != nil {
		for _, status := range opts.Status {
			baseParams.Add("status_category[one_of]", status)
		}
		for _, severity := range opts.Severity {
			baseParams.Add("severity[one_of]", severity)
		}

		// Add date filters for created_at
		if opts.CreatedAtGTE != "" {
			baseParams.Set("created_at[gte]", opts.CreatedAtGTE)
		}
		if opts.CreatedAtLTE != "" {
			baseParams.Set("created_at[lte]", opts.CreatedAtLTE)
		}
		if opts.CreatedAtRange != "" {
			baseParams.Set("created_at[date_range]", opts.CreatedAtRange)
		}

		// Add date filters for updated_at
		if opts.UpdatedAtGTE != "" {
			baseParams.Set("updated_at[gte]", opts.UpdatedAtGTE)
		}
		if opts.UpdatedAtLTE != "" {
			baseParams.Set("updated_at[lte]", opts.UpdatedAtLTE)
		}
		if opts.UpdatedAtRange != "" {
			baseParams.Set("updated_at[date_range]", opts.UpdatedAtRange)
		}
	}

	// Paginate through all results
	maxPages := 10 // Safety limit
	for page := 0; page < maxPages; page++ {
		params := url.Values{}
		// Copy base parameters
		for k, v := range baseParams {
			params[k] = v
		}

		params.Set("page_size", strconv.Itoa(pageSize))
		if after != "" {
			params.Set("after", after)
		}

		respBody, err := c.doRequest("GET", "/incidents", params, nil)
		if err != nil {
			return nil, err
		}

		var response ListIncidentsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		allIncidents = append(allIncidents, response.Incidents...)

		// Check if there are more pages
		if response.PaginationMeta.After == "" || len(response.Incidents) == 0 {
			break
		}
		after = response.PaginationMeta.After
	}

	// Return combined results
	// Note: Auto-pagination returns all results, so total_record_count equals the number fetched
	return &ListIncidentsResponse{
		Incidents: allIncidents,
		ListResponse: ListResponse{
			PaginationMeta: struct {
				After            string `json:"after,omitempty"`
				PageSize         int    `json:"page_size"`
				TotalRecordCount int    `json:"total_record_count,omitempty"`
			}{
				PageSize:         pageSize,
				TotalRecordCount: len(allIncidents), // Total count is number of incidents fetched
			},
		},
	}, nil
}

// GetIncident retrieves a specific incident by ID
func (c *Client) GetIncident(id string) (*Incident, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/incidents/%s", id), nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}

// CreateIncident creates a new incident
func (c *Client) CreateIncident(req *CreateIncidentRequest) (*Incident, error) {
	respBody, err := c.doRequest("POST", "/incidents", nil, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}

// UpdateIncident updates an existing incident using V2 actions/edit API
func (c *Client) UpdateIncident(id string, req *UpdateIncidentRequest) (*Incident, error) {
	// Use the correct V2 actions/edit endpoint
	editRequest := map[string]interface{}{
		"notify_incident_channel": true,
	}

	// Build the incident object with only the fields that are being updated
	incident := make(map[string]interface{})

	if req.Name != "" {
		incident["name"] = req.Name
	}
	if req.Summary != "" {
		incident["summary"] = req.Summary
	}
	if req.IncidentStatusID != "" {
		incident["incident_status_id"] = req.IncidentStatusID
	}
	if req.SeverityID != "" {
		incident["severity_id"] = req.SeverityID
	}
	if req.CallURL != "" {
		incident["call_url"] = req.CallURL
	}
	if req.SlackChannelNameOverride != "" {
		incident["slack_channel_name_override"] = req.SlackChannelNameOverride
	}
	if len(req.IncidentRoleAssignments) > 0 {
		incident["incident_role_assignments"] = req.IncidentRoleAssignments
	}
	if len(req.CustomFieldEntries) > 0 {
		incident["custom_field_entries"] = req.CustomFieldEntries
	}
	if len(req.IncidentTimestampValues) > 0 {
		incident["incident_timestamp_values"] = req.IncidentTimestampValues
	}

	// Only include incident object if there are fields to update
	if len(incident) > 0 {
		editRequest["incident"] = incident
	} else {
		return nil, fmt.Errorf("no fields to update")
	}

	respBody, err := c.doRequest("POST", fmt.Sprintf("/incidents/%s/actions/edit", id), nil, editRequest)
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}

// AssignIncidentRoleRequest represents a request to assign a role to a user
type AssignIncidentRoleRequest struct {
	IncidentRoleID string `json:"incident_role_id"`
	UserID         string `json:"user_id"`
}

// AssignIncidentRole assigns a specific role to a user for an incident
func (c *Client) AssignIncidentRole(incidentID string, req *AssignIncidentRoleRequest) (*Incident, error) {
	respBody, err := c.doRequest("PATCH", fmt.Sprintf("/incidents/%s", incidentID), nil, map[string]interface{}{
		"incident_role_assignments": []map[string]interface{}{
			{
				"incident_role_id": req.IncidentRoleID,
				"user_id":          req.UserID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Incident Incident `json:"incident"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response.Incident, nil
}

// GetIncidentDebrief retrieves the debrief/post-mortem document for an incident
// Returns the incident details with has_debrief status and postmortem_document_url if available
func (c *Client) GetIncidentDebrief(id string) (*Incident, error) {
	incident, err := c.GetIncident(id)
	if err != nil {
		return nil, err
	}

	if !incident.HasDebrief {
		return nil, fmt.Errorf("incident %s does not have a debrief document yet", id)
	}

	if incident.PostmortemDocumentURL == "" {
		return nil, fmt.Errorf("incident %s has a debrief but no postmortem_document_url is available", id)
	}

	return incident, nil
}
