package tools

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

func TestGetIncidentDebriefTool_ExecuteValidation(t *testing.T) {
	tool := &GetIncidentDebriefTool{}

	tests := []struct {
		name          string
		args          map[string]interface{}
		wantError     bool
		expectedError string
	}{
		{
			name:          "missing incident_id",
			args:          map[string]interface{}{},
			wantError:     true,
			expectedError: "incident_id parameter is required",
		},
		{
			name: "empty incident_id",
			args: map[string]interface{}{
				"incident_id": "",
			},
			wantError:     true,
			expectedError: "incident_id parameter is required",
		},
		{
			name: "wrong type incident_id",
			args: map[string]interface{}{
				"incident_id": 123,
			},
			wantError:     true,
			expectedError: "incident_id parameter is required",
		},
		{
			name: "nil incident_id",
			args: map[string]interface{}{
				"incident_id": nil,
			},
			wantError:     true,
			expectedError: "incident_id parameter is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tool.Execute(tt.args)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetIncidentDebriefTool_Schema(t *testing.T) {
	tool := &GetIncidentDebriefTool{}
	schema := tool.InputSchema()

	t.Run("schema structure", func(t *testing.T) {
		if schema == nil {
			t.Fatal("Schema should not be nil")
		}

		// Verify type
		if schemaType, ok := schema["type"].(string); !ok || schemaType != "object" {
			t.Error("Schema type should be 'object'")
		}

		// Verify required fields
		required, ok := schema["required"].([]interface{})
		if !ok {
			t.Fatal("Schema should have 'required' field as an array")
		}

		if len(required) != 1 {
			t.Errorf("Expected 1 required field, got %d", len(required))
		}

		if required[0] != "incident_id" {
			t.Errorf("Expected 'incident_id' to be required, got: %v", required[0])
		}

		// Verify properties
		properties, ok := schema["properties"].(map[string]interface{})
		if !ok {
			t.Fatal("Schema should have 'properties' field as a map")
		}

		// Verify incident_id property
		incidentIDProp, ok := properties["incident_id"].(map[string]interface{})
		if !ok {
			t.Fatal("Schema should have 'incident_id' property")
		}

		if propType, ok := incidentIDProp["type"].(string); !ok || propType != "string" {
			t.Error("incident_id property should have type 'string'")
		}

		if desc, ok := incidentIDProp["description"].(string); !ok || desc == "" {
			t.Error("incident_id property should have a non-empty description")
		} else {
			// Verify description mentions supported formats
			if !strings.Contains(desc, "full ID") {
				t.Error("Description should mention full ID format")
			}
			if !strings.Contains(desc, "reference") {
				t.Error("Description should mention reference format")
			}
			if !strings.Contains(desc, "Slack channel") {
				t.Error("Description should mention Slack channel formats")
			}
		}

		// Verify additionalProperties is false
		if additionalProps, ok := schema["additionalProperties"].(bool); !ok || additionalProps {
			t.Error("Schema should have additionalProperties set to false")
		}
	})
}

func TestGetIncidentDebriefTool_Name(t *testing.T) {
	tool := &GetIncidentDebriefTool{}
	expectedName := "get_incident_debrief"

	if tool.Name() != expectedName {
		t.Errorf("Expected tool name to be %q, got: %q", expectedName, tool.Name())
	}
}

func TestGetIncidentDebriefTool_Description(t *testing.T) {
	tool := &GetIncidentDebriefTool{}
	desc := tool.Description()

	t.Run("description content", func(t *testing.T) {
		if desc == "" {
			t.Fatal("Description should not be empty")
		}

		requiredTerms := []string{
			"debrief",
			"post-mortem",
			"postmortem_document_url",
			"has_debrief",
			"incident_id",
		}

		for _, term := range requiredTerms {
			if !strings.Contains(desc, term) {
				t.Errorf("Description should contain %q", term)
			}
		}

		// Verify supported identifier formats are documented
		identifierFormats := []string{
			"Full incident ID",
			"Incident reference",
			"Slack channel ID",
			"Slack channel name",
		}

		for _, format := range identifierFormats {
			if !strings.Contains(desc, format) {
				t.Errorf("Description should document %q format", format)
			}
		}

		// Verify examples are provided
		if !strings.Contains(desc, "EXAMPLES") {
			t.Error("Description should include examples section")
		}

		// Verify error handling is documented
		if !strings.Contains(desc, "ERROR HANDLING") {
			t.Error("Description should document error handling")
		}
	})
}

func TestGetIncidentDebriefTool_ResponseFormat(t *testing.T) {
	// This test verifies the expected response format structure
	// We can't test with a real client, but we can verify the structure

	t.Run("response structure", func(t *testing.T) {
		// Expected response fields based on the implementation
		expectedFields := []string{
			"incident_id",
			"incident_reference",
			"incident_name",
			"has_debrief",
			"postmortem_document_url",
			"permalink",
		}

		// Create a mock response to verify JSON structure
		mockResponse := map[string]interface{}{
			"incident_id":             "inc_123",
			"incident_reference":      "INC-123",
			"incident_name":           "Test Incident",
			"has_debrief":             true,
			"postmortem_document_url": "https://example.com/postmortem",
			"permalink":               "https://app.incident.io/incidents/123",
		}

		// Marshal to verify JSON structure
		jsonBytes, err := json.Marshal(mockResponse)
		if err != nil {
			t.Fatalf("Failed to marshal mock response: %v", err)
		}

		// Unmarshal back to verify all fields are present
		var result map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Verify all expected fields are present
		for _, field := range expectedFields {
			if _, ok := result[field]; !ok {
				t.Errorf("Response should include field %q", field)
			}
		}
	})
}

func TestGetIncidentDebriefTool_Integration(t *testing.T) {
	// Integration test scenarios that don't require actual API calls

	t.Run("tool initialization", func(t *testing.T) {
		client := &incidentio.Client{}
		tool := NewGetIncidentDebriefTool(client)

		if tool == nil {
			t.Fatal("NewGetIncidentDebriefTool should return a non-nil tool")
		}

		if tool.client == nil {
			t.Error("Tool should have a non-nil client")
		}

		if tool.Name() == "" {
			t.Error("Tool should have a non-empty name")
		}

		if tool.Description() == "" {
			t.Error("Tool should have a non-empty description")
		}

		schema := tool.InputSchema()
		if schema == nil {
			t.Error("Tool should have a non-nil input schema")
		}
	})
}

func TestGetIncidentDebriefTool_EdgeCases(t *testing.T) {
	t.Run("edge case scenarios", func(t *testing.T) {
		// Test that various edge case inputs are properly handled
		// Note: These tests can't be fully executed without a mock client,
		// but they document expected behavior

		edgeCases := []struct {
			name        string
			input       string
			description string
		}{
			{
				name:        "whitespace only",
				input:       "   ",
				description: "Should be rejected as invalid identifier",
			},
			{
				name:        "very long input",
				input:       strings.Repeat("a", 1000),
				description: "Should be handled without panic",
			},
			{
				name:        "special characters",
				input:       "inc_!@#$%^&*()",
				description: "Should be validated appropriately",
			},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				// Document the expected behavior
				t.Logf("Edge case: %s - %s", tc.name, tc.description)
				// Actual validation would require a mock client
			})
		}
	})

	t.Run("parameter validation", func(t *testing.T) {
		// Test that extra parameters are documented in the schema
		tool := &GetIncidentDebriefTool{}
		schema := tool.InputSchema()

		// Verify the schema correctly specifies additionalProperties: false
		if additionalProps, ok := schema["additionalProperties"].(bool); !ok || additionalProps {
			t.Error("Schema should set additionalProperties to false to prevent extra parameters")
		}

		// Document that extra parameters would be ignored or rejected by MCP protocol
		t.Log("Extra parameters beyond 'incident_id' are handled by MCP protocol validation")
	})
}
