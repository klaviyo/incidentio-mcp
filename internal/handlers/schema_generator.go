package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// GenerateFieldsDescription generates a comprehensive fields description from a Go type
func GenerateFieldsDescription(exampleType interface{}) string {
	topLevel, nested := extractFieldsFromType(reflect.TypeOf(exampleType))

	var desc strings.Builder
	desc.WriteString("Comma-separated list of fields to include in response to reduce context usage. Supports nested fields with dot notation.\n\n")

	// Top-level fields
	desc.WriteString("Available top-level fields: ")
	desc.WriteString(strings.Join(topLevel, ", "))
	desc.WriteString("\n")

	// Nested fields
	if len(nested) > 0 {
		desc.WriteString("\nNested fields (dot notation):\n")
		for parentField, childFields := range nested {
			desc.WriteString(fmt.Sprintf("- %s: %s\n", parentField, strings.Join(childFields, ", ")))
		}
	}

	desc.WriteString("\nExamples: \"id,name\" or with nested fields: \"id,name,severity.name,incident_status.category\"\n")
	desc.WriteString("Omit to return all fields.")

	return desc.String()
}

// extractFieldsFromType extracts field names from a struct type using reflection
func extractFieldsFromType(t reflect.Type) (topLevel []string, nested map[string][]string) {
	nested = make(map[string][]string)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Only process struct types
	if t.Kind() != reflect.Struct {
		return topLevel, nested
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON tag (remove omitempty, etc.)
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			continue
		}

		topLevel = append(topLevel, jsonName)

		// Check if field is a struct (nested object)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Process nested struct fields
		if fieldType.Kind() == reflect.Struct && !isTimeType(fieldType) {
			nestedFields := extractNestedFields(fieldType)
			if len(nestedFields) > 0 {
				nested[jsonName] = nestedFields
			}
		}
	}

	return topLevel, nested
}

// extractNestedFields extracts field names from a nested struct
func extractNestedFields(t reflect.Type) []string {
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON tag
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName != "" {
			fields = append(fields, jsonName)
		}
	}

	return fields
}

// isTimeType checks if a type is time.Time
func isTimeType(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

// GetIncidentFieldsDescription returns the fields description for Incident type
func GetIncidentFieldsDescription() string {
	return GenerateFieldsDescription(client.Incident{})
}

// GetAlertFieldsDescription returns the fields description for Alert type
func GetAlertFieldsDescription() string {
	return GenerateFieldsDescription(client.Alert{})
}
