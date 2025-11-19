package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

func TestDebugIncidentTool_Name(t *testing.T) {
	tool := NewDebugIncidentTool(nil)
	if tool.Name() != "debug_incident" {
		t.Errorf("Expected name 'debug_incident', got '%s'", tool.Name())
	}
}

func TestDebugIncidentTool_Description(t *testing.T) {
	tool := NewDebugIncidentTool(nil)
	desc := tool.Description()

	// Check for key phrases in the description
	expectedPhrases := []string{
		"DEBUG TOOL",
		"raw incident data",
		"inspect all fields",
		"debrief-related fields",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(desc, phrase) {
			t.Errorf("Description missing expected phrase: %s", phrase)
		}
	}
}

func TestDebugIncidentTool_InputSchema(t *testing.T) {
	tool := NewDebugIncidentTool(nil)
	schema := tool.InputSchema()

	// Verify required fields
	if schema["type"] != "object" {
		t.Errorf("Expected type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check incident_id property
	incidentIDProp, ok := properties["incident_id"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected incident_id property")
	}

	if incidentIDProp["type"] != "string" {
		t.Errorf("Expected incident_id type 'string', got '%v'", incidentIDProp["type"])
	}

	// Check required fields
	required, ok := schema["required"].([]interface{})
	if !ok || len(required) != 1 || required[0] != "incident_id" {
		t.Error("Expected incident_id to be required")
	}
}

func TestDebugIncidentTool_Execute(t *testing.T) {
	tests := []struct {
		name             string
		incidentID       string
		mockResponse     string
		mockStatusCode   int
		wantError        bool
		expectedInResult []string
	}{
		{
			name:           "successful debug with all debrief fields",
			incidentID:     "01HXYZ1234567890ABCDEFGH",
			mockStatusCode: http.StatusOK,
			mockResponse: `{
				"incident": {
					"id": "01HXYZ1234567890ABCDEFGH",
					"reference": "INC-123",
					"name": "Test Incident",
					"permalink": "https://app.incident.io/incidents/inc_123",
					"mode": "standard",
					"has_debrief": true,
					"postmortem_document_url": "https://docs.example.com/postmortem/inc-123",
					"debrief_export_id": "export_123",
					"retrospective_incident_options": {
						"external_id": 12345,
						"postmortem_document_url": "https://docs.example.com/postmortem/inc-123",
						"slack_channel_id": "C123456789"
					},
					"incident_status": {
						"id": "status_123",
						"name": "Resolved",
						"description": "Incident resolved",
						"category": "closed",
						"rank": 1,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"severity": {
						"id": "sev_123",
						"name": "SEV1",
						"description": "Critical",
						"rank": 1,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"incident_type": {
						"id": "type_123",
						"name": "Production",
						"description": "Production incident",
						"is_default": true,
						"private_incidents_only": false,
						"create_in_triage": "optional",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z"
				}
			}`,
			wantError: false,
			expectedInResult: []string{
				"DIAGNOSTIC SUMMARY",
				"FULL INCIDENT DATA",
				"ANALYSIS",
				"01HXYZ1234567890ABCDEFGH",
				"INC-123",
				"Test Incident",
				"has_debrief",
				"postmortem_document_url",
				"retrospective_incident_options",
				"debrief_export_id",
			},
		},
		{
			name:           "debug incident without debrief",
			incidentID:     "01HABC9876543210ZYXWVUTS",
			mockStatusCode: http.StatusOK,
			mockResponse: `{
				"incident": {
					"id": "01HABC9876543210ZYXWVUTS",
					"reference": "INC-456",
					"name": "Test Incident No Debrief",
					"permalink": "https://app.incident.io/incidents/inc_456",
					"mode": "standard",
					"has_debrief": false,
					"incident_status": {
						"id": "status_123",
						"name": "Active",
						"description": "Active incident",
						"category": "active",
						"rank": 1,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"severity": {
						"id": "sev_123",
						"name": "SEV2",
						"description": "Major",
						"rank": 2,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"incident_type": {
						"id": "type_123",
						"name": "Production",
						"description": "Production incident",
						"is_default": true,
						"private_incidents_only": false,
						"create_in_triage": "optional",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z"
				}
			}`,
			wantError: false,
			expectedInResult: []string{
				"DIAGNOSTIC SUMMARY",
				"FULL INCIDENT DATA",
				"01HABC9876543210ZYXWVUTS",
				"INC-456",
				"has_debrief",
			},
		},
		{
			name:           "debug with retrospective incident",
			incidentID:     "01HPQR1122334455MNOPQRST",
			mockStatusCode: http.StatusOK,
			mockResponse: `{
				"incident": {
					"id": "01HPQR1122334455MNOPQRST",
					"reference": "INC-789",
					"name": "Retrospective Incident",
					"permalink": "https://app.incident.io/incidents/inc_789",
					"mode": "retrospective",
					"has_debrief": true,
					"retrospective_incident_options": {
						"external_id": 98765,
						"postmortem_document_url": "https://docs.example.com/retro/inc-789",
						"slack_channel_id": "C987654321"
					},
					"incident_status": {
						"id": "status_123",
						"name": "Closed",
						"description": "Closed incident",
						"category": "closed",
						"rank": 1,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"severity": {
						"id": "sev_123",
						"name": "SEV3",
						"description": "Minor",
						"rank": 3,
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"incident_type": {
						"id": "type_123",
						"name": "Retrospective",
						"description": "Retrospective incident",
						"is_default": false,
						"private_incidents_only": false,
						"create_in_triage": "optional",
						"created_at": "2024-01-01T00:00:00Z",
						"updated_at": "2024-01-01T00:00:00Z"
					},
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z"
				}
			}`,
			wantError: false,
			expectedInResult: []string{
				"DIAGNOSTIC SUMMARY",
				"retrospective_incident_options",
				"mode",
				"retrospective",
			},
		},
		{
			name:           "API error - 404 not found",
			incidentID:     "invalid_id",
			mockStatusCode: http.StatusNotFound,
			mockResponse: `{
				"type": "not_found",
				"status": 404,
				"message": "The requested incident could not be found"
			}`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				fmt.Fprint(w, tt.mockResponse)
			}))
			defer server.Close()

			// Create client with mock server
			// Set environment variable for API key
			t.Setenv("INCIDENT_IO_API_KEY", "test-key")
			t.Setenv("INCIDENT_IO_BASE_URL", server.URL)
			client, err := incidentio.NewClient()
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Create tool
			tool := NewDebugIncidentTool(client)

			// Execute tool
			result, err := tool.Execute(map[string]interface{}{
				"incident_id": tt.incidentID,
			})

			// Check error expectation
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify result contains expected strings
			for _, expected := range tt.expectedInResult {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain '%s', but it didn't. Result: %s", expected, result)
				}
			}
		})
	}
}

func TestDebugIncidentTool_Execute_ParameterValidation(t *testing.T) {
	tool := NewDebugIncidentTool(nil)

	tests := []struct {
		name      string
		args      map[string]interface{}
		wantError bool
	}{
		{
			name:      "missing incident_id parameter",
			args:      map[string]interface{}{},
			wantError: true,
		},
		{
			name: "empty incident_id parameter",
			args: map[string]interface{}{
				"incident_id": "",
			},
			wantError: true,
		},
		{
			name: "invalid incident_id type",
			args: map[string]interface{}{
				"incident_id": 123,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tool.Execute(tt.args)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestDebugIncidentTool_Execute_DiagnosticFields(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"incident": map[string]interface{}{
				"id":                      "01HSTU7788990011FGHIJKLM",
				"reference":               "INC-999",
				"name":                    "Diagnostic Test",
				"permalink":               "https://app.incident.io/incidents/inc_999",
				"mode":                    "standard",
				"has_debrief":             true,
				"postmortem_document_url": "https://docs.example.com/postmortem/inc-999",
				"debrief_export_id":       "export_999",
				"retrospective_incident_options": map[string]interface{}{
					"external_id":            99999,
					"postmortem_document_url": "https://docs.example.com/nested/inc-999",
					"slack_channel_id":        "C999999999",
				},
				"incident_status": map[string]interface{}{
					"id":          "status_123",
					"name":        "Resolved",
					"description": "Incident resolved",
					"category":    "closed",
					"rank":        1,
					"created_at":  "2024-01-01T00:00:00Z",
					"updated_at":  "2024-01-01T00:00:00Z",
				},
				"severity": map[string]interface{}{
					"id":          "sev_123",
					"name":        "SEV1",
					"description": "Critical",
					"rank":        1,
					"created_at":  "2024-01-01T00:00:00Z",
					"updated_at":  "2024-01-01T00:00:00Z",
				},
				"incident_type": map[string]interface{}{
					"id":                      "type_123",
					"name":                    "Production",
					"description":             "Production incident",
					"is_default":              true,
					"private_incidents_only":  false,
					"create_in_triage":        "optional",
					"created_at":              "2024-01-01T00:00:00Z",
					"updated_at":              "2024-01-01T00:00:00Z",
				},
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	t.Setenv("INCIDENT_IO_API_KEY", "test-key")
	t.Setenv("INCIDENT_IO_BASE_URL", server.URL)
	client, err := incidentio.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tool := NewDebugIncidentTool(client)
	result, err := tool.Execute(map[string]interface{}{
		"incident_id": "01HSTU7788990011FGHIJKLM",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify diagnostic summary sections are present
	requiredSections := []string{
		"=== DIAGNOSTIC SUMMARY ===",
		"=== FULL INCIDENT DATA ===",
		"=== ANALYSIS ===",
	}

	for _, section := range requiredSections {
		if !strings.Contains(result, section) {
			t.Errorf("Expected result to contain section '%s'", section)
		}
	}

	// Verify all diagnostic fields are present
	diagnosticFields := []string{
		"incident_id",
		"incident_reference",
		"incident_name",
		"mode",
		"has_debrief",
		"debrief_fields",
		"postmortem_document_url_present",
		"postmortem_document_url_value",
		"retrospective_options_present",
		"debrief_export_id_present",
		"debrief_export_id_value",
		"permalink",
		"retrospective_incident_options",
	}

	for _, field := range diagnosticFields {
		if !strings.Contains(result, field) {
			t.Errorf("Expected result to contain diagnostic field '%s'", field)
		}
	}

	// Verify analysis output
	analysisChecks := []string{
		"Has Debrief: true",
		"Postmortem URL at top level: true",
		"Retrospective Options object: true",
		"Debrief Export ID: true",
	}

	for _, check := range analysisChecks {
		if !strings.Contains(result, check) {
			t.Errorf("Expected result to contain analysis check '%s'", check)
		}
	}
}
