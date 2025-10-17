package tools

import (
	"strings"
	"testing"
)

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func TestCreateIncidentTool_Execute(t *testing.T) {
	tool := &CreateIncidentTool{}

	// Test missing required name parameter
	t.Run("missing required name", func(t *testing.T) {
		args := map[string]interface{}{
			"summary": "Test Summary",
		}

		_, err := tool.Execute(args)
		if err == nil {
			t.Error("Expected error for missing name parameter")
		}
		if err.Error() != "name parameter is required" {
			t.Errorf("Expected 'name parameter is required' error, got: %v", err)
		}
	})

	// Test name parameter with wrong type
	t.Run("name parameter wrong type", func(t *testing.T) {
		args := map[string]interface{}{
			"name": 123, // Not a string
		}

		_, err := tool.Execute(args)
		if err == nil {
			t.Error("Expected error for wrong type name parameter")
		}
		if err.Error() != "name parameter is required" {
			t.Errorf("Expected 'name parameter is required' error, got: %v", err)
		}
	})

	// Note: We can't test the full execution without a real client,
	// but we can test the parameter validation and schema
}

func TestCreateIncidentTool_Schema(t *testing.T) {
	tool := &CreateIncidentTool{}

	// Test Name
	if tool.Name() != "create_incident" {
		t.Errorf("Expected name 'create_incident', got %s", tool.Name())
	}

	// Test Description - verify it's not empty and contains key information
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
	// Check for key elements in the description
	if !contains(desc, "Create") && !contains(desc, "incident") {
		t.Errorf("Description should mention creating incidents: %s", desc)
	}

	// Test InputSchema
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}

	properties := schema["properties"].(map[string]interface{})
	if _, ok := properties["name"]; !ok {
		t.Error("Schema should have 'name' property")
	}

	required := schema["required"].([]interface{})
	if len(required) != 1 || required[0] != "name" {
		t.Error("Schema should require only 'name'")
	}
}
