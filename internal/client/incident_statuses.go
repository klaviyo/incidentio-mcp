package client

import (
	"encoding/json"
	"fmt"
)

// Using IncidentStatus type from types.go

// ListIncidentStatusesResponse represents the response from listing incident statuses
type ListIncidentStatusesResponse struct {
	IncidentStatuses []IncidentStatus `json:"incident_statuses"`
}

// ListIncidentStatuses returns all incident statuses
func (c *Client) ListIncidentStatuses() (*ListIncidentStatusesResponse, error) {
	// Note: Incident statuses are under V1 API, not V2
	// We need to temporarily change the base URL for this request
	originalBaseURL := c.BaseURL()
	c.SetBaseURL("https://api.incident.io/v1")
	defer func() { c.SetBaseURL(originalBaseURL) }()

	respBody, err := c.doRequest("GET", "/incident_statuses", nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListIncidentStatusesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
