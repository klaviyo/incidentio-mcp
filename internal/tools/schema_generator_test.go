package tools

import (
	"reflect"
	"strings"
	"testing"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

func TestGenerateFieldsDescription_Incident(t *testing.T) {
	desc := GetIncidentFieldsDescription()

	// Verify it contains key sections
	if !strings.Contains(desc, "Available top-level fields:") {
		t.Error("Description should contain 'Available top-level fields:'")
	}

	if !strings.Contains(desc, "Nested fields (dot notation):") {
		t.Error("Description should contain 'Nested fields (dot notation):'")
	}

	// Verify key incident fields are present
	requiredFields := []string{"id", "name", "reference", "severity", "incident_status"}
	for _, field := range requiredFields {
		if !strings.Contains(desc, field) {
			t.Errorf("Description should contain field '%s'", field)
		}
	}

	// Verify nested fields are documented
	if !strings.Contains(desc, "severity:") {
		t.Error("Description should document severity nested fields")
	}

	if !strings.Contains(desc, "incident_status:") {
		t.Error("Description should document incident_status nested fields")
	}

	// Verify examples are included
	if !strings.Contains(desc, "Examples:") {
		t.Error("Description should contain examples")
	}
}

func TestGenerateFieldsDescription_Alert(t *testing.T) {
	desc := GetAlertFieldsDescription()

	// Verify key alert fields are present
	requiredFields := []string{"id", "title", "status", "source"}
	for _, field := range requiredFields {
		if !strings.Contains(desc, field) {
			t.Errorf("Description should contain field '%s'", field)
		}
	}

	// Verify nested incident field is documented (alerts can have nested incidents)
	if !strings.Contains(desc, "incident:") {
		t.Error("Description should document incident nested fields")
	}
}

func TestExtractFieldsFromType(t *testing.T) {
	// Test with Incident type
	topLevel, nested := extractFieldsFromType(reflect.TypeOf(incidentio.Incident{}))

	// Verify some top-level fields
	expectedTopLevel := []string{"id", "name", "reference"}
	for _, field := range expectedTopLevel {
		found := false
		for _, f := range topLevel {
			if f == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find top-level field '%s'", field)
		}
	}

	// Verify nested fields exist for severity and incident_status
	if _, ok := nested["severity"]; !ok {
		t.Error("Expected nested fields for 'severity'")
	}

	if _, ok := nested["incident_status"]; !ok {
		t.Error("Expected nested fields for 'incident_status'")
	}

	// Verify severity has expected nested fields
	if severityFields, ok := nested["severity"]; ok {
		expectedSeverityFields := []string{"id", "name", "description", "rank"}
		for _, expected := range expectedSeverityFields {
			found := false
			for _, actual := range severityFields {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected severity to have nested field '%s'", expected)
			}
		}
	}
}

func TestGenerateFieldsDescription_CustomStruct(t *testing.T) {
	// Test with a simple custom struct
	type TestStruct struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Meta struct {
			Version string `json:"version"`
			Author  string `json:"author"`
		} `json:"meta"`
	}

	desc := GenerateFieldsDescription(TestStruct{})

	// Should have top-level fields
	if !strings.Contains(desc, "id") {
		t.Error("Should contain 'id' field")
	}

	if !strings.Contains(desc, "name") {
		t.Error("Should contain 'name' field")
	}

	if !strings.Contains(desc, "meta") {
		t.Error("Should contain 'meta' field")
	}

	// Should have nested fields for meta
	if !strings.Contains(desc, "meta:") {
		t.Error("Should document nested meta fields")
	}

	if !strings.Contains(desc, "version") {
		t.Error("Should contain nested 'version' field")
	}
}

func TestGenerateFieldsDescription_NoNestedFields(t *testing.T) {
	// Test with a flat struct (no nested objects)
	type FlatStruct struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	desc := GenerateFieldsDescription(FlatStruct{})

	// Should have top-level fields
	if !strings.Contains(desc, "id, name, value") {
		t.Error("Should contain all flat fields")
	}

	// Should not have nested fields section when there are none
	// The nested section is only added when there are nested fields
	if strings.Contains(desc, "Nested fields (dot notation):") {
		// Check if it's actually empty (just the header)
		lines := strings.Split(desc, "\n")
		for i, line := range lines {
			if strings.Contains(line, "Nested fields (dot notation):") {
				// Next line should be empty or the Examples section
				if i+1 < len(lines) && !strings.HasPrefix(lines[i+1], "\nExamples:") && strings.TrimSpace(lines[i+1]) != "" {
					// There are actual nested fields listed, which shouldn't happen for flat struct
					t.Error("Should not have nested fields for flat struct")
				}
			}
		}
	}
}
