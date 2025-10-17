package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFilterFields_NoFieldsSpecified(t *testing.T) {
	data := map[string]interface{}{
		"id":   "123",
		"name": "Test",
		"value": 42,
	}

	result, err := FilterFields(data, "")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Should return all fields
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	if len(parsed) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(parsed))
	}
}

func TestFilterFields_TopLevelFields(t *testing.T) {
	data := map[string]interface{}{
		"id":      "123",
		"name":    "Test",
		"summary": "A test summary",
		"value":   42,
	}

	result, err := FilterFields(data, "id,name")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	if len(parsed) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(parsed))
	}

	if parsed["id"] != "123" {
		t.Errorf("Expected id='123', got %v", parsed["id"])
	}

	if parsed["name"] != "Test" {
		t.Errorf("Expected name='Test', got %v", parsed["name"])
	}

	if _, exists := parsed["summary"]; exists {
		t.Error("Expected summary to be filtered out")
	}

	if _, exists := parsed["value"]; exists {
		t.Error("Expected value to be filtered out")
	}
}

func TestFilterFields_NestedFields(t *testing.T) {
	data := map[string]interface{}{
		"id":   "123",
		"name": "Test",
		"severity": map[string]interface{}{
			"id":          "sev_1",
			"name":        "Critical",
			"description": "Very important",
			"rank":        1,
		},
		"status": map[string]interface{}{
			"id":   "status_1",
			"name": "Active",
		},
	}

	result, err := FilterFields(data, "id,name,severity.name,status.id")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	// Check top-level fields
	if parsed["id"] != "123" {
		t.Errorf("Expected id='123', got %v", parsed["id"])
	}

	// Check nested severity field
	severity, ok := parsed["severity"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected severity to be a map")
	}

	if severity["name"] != "Critical" {
		t.Errorf("Expected severity.name='Critical', got %v", severity["name"])
	}

	if _, exists := severity["description"]; exists {
		t.Error("Expected severity.description to be filtered out")
	}

	// Check nested status field
	status, ok := parsed["status"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected status to be a map")
	}

	if status["id"] != "status_1" {
		t.Errorf("Expected status.id='status_1', got %v", status["id"])
	}

	if _, exists := status["name"]; exists {
		t.Error("Expected status.name to be filtered out")
	}
}

func TestFilterFields_ArrayOfObjects(t *testing.T) {
	data := map[string]interface{}{
		"incidents": []interface{}{
			map[string]interface{}{
				"id":      "inc_1",
				"name":    "Incident 1",
				"summary": "Summary 1",
			},
			map[string]interface{}{
				"id":      "inc_2",
				"name":    "Incident 2",
				"summary": "Summary 2",
			},
		},
	}

	// Note: The current implementation filters arrays, but field filtering
	// for array elements is tricky. For now, we wrap it with incidents.
	result, err := FilterFields(data, "incidents")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	incidents, ok := parsed["incidents"].([]interface{})
	if !ok {
		t.Fatal("Expected incidents to be an array")
	}

	if len(incidents) != 2 {
		t.Errorf("Expected 2 incidents, got %d", len(incidents))
	}
}

func TestFilterFields_WithSpaces(t *testing.T) {
	data := map[string]interface{}{
		"id":   "123",
		"name": "Test",
		"value": 42,
	}

	result, err := FilterFields(data, " id , name ")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	if len(parsed) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(parsed))
	}
}

func TestFilterFields_EmptyFields(t *testing.T) {
	data := map[string]interface{}{
		"id":   "123",
		"name": "Test",
	}

	result, err := FilterFields(data, ",,  ,")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Empty or whitespace-only fields should be ignored
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	// Should return all fields since no valid fields were specified
	if len(parsed) != 2 {
		t.Errorf("Expected 2 fields (all), got %d", len(parsed))
	}
}

func TestFilterFields_ComplexNesting(t *testing.T) {
	data := map[string]interface{}{
		"id":   "inc_123",
		"name": "Production Outage",
		"incident_status": map[string]interface{}{
			"id":          "status_1",
			"name":        "Active",
			"category":    "triage",
			"description": "In triage",
		},
		"severity": map[string]interface{}{
			"id":   "sev_1",
			"name": "Critical",
			"rank": 1,
		},
		"incident_role_assignments": []interface{}{
			map[string]interface{}{
				"role": map[string]interface{}{
					"id":   "role_1",
					"name": "Incident Lead",
				},
				"assignee": map[string]interface{}{
					"id":   "user_1",
					"name": "John Doe",
				},
			},
		},
	}

	result, err := FilterFields(data, "id,name,severity.name,incident_status.category")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	// Verify correct fields are present
	if parsed["id"] != "inc_123" {
		t.Errorf("Expected id='inc_123', got %v", parsed["id"])
	}

	severity, ok := parsed["severity"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected severity to be a map")
	}

	if severity["name"] != "Critical" {
		t.Errorf("Expected severity.name='Critical', got %v", severity["name"])
	}

	if _, exists := severity["rank"]; exists {
		t.Error("Expected severity.rank to be filtered out")
	}

	status, ok := parsed["incident_status"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected incident_status to be a map")
	}

	if status["category"] != "triage" {
		t.Errorf("Expected incident_status.category='triage', got %v", status["category"])
	}

	if _, exists := status["description"]; exists {
		t.Error("Expected incident_status.description to be filtered out")
	}

	// incident_role_assignments should not be present since we didn't request it
	if _, exists := parsed["incident_role_assignments"]; exists {
		t.Error("Expected incident_role_assignments to be filtered out")
	}
}

func TestFilterFields_JSONFormatting(t *testing.T) {
	data := map[string]interface{}{
		"id":   "123",
		"name": "Test",
	}

	result, err := FilterFields(data, "id,name")
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Verify result is valid, formatted JSON
	if !strings.Contains(result, "\n") {
		t.Error("Expected formatted JSON with newlines")
	}

	if !strings.Contains(result, "  ") {
		t.Error("Expected indented JSON")
	}
}
