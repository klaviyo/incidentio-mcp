package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FilterFields filters a JSON object to only include the specified fields.
// Fields can be specified as:
// - Top-level fields: "id", "name", "summary"
// - Nested fields with dot notation: "severity.name", "incident_status.category"
// - Array elements are filtered recursively
//
// Example:
//   fields := "id,name,severity.name,incident_status.category"
//   filtered, err := FilterFields(data, fields)
func FilterFields(data interface{}, fieldsStr string) (string, error) {
	if fieldsStr == "" {
		// No filtering requested, return original data
		result, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal data: %w", err)
		}
		return string(result), nil
	}

	// Parse field list
	fields := parseFieldList(fieldsStr)

	// Marshal to JSON first to get map representation
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	var rawData interface{}
	if err := json.Unmarshal(jsonBytes, &rawData); err != nil {
		return "", fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Filter the data
	filtered := filterObject(rawData, fields)

	// Marshal the filtered result
	result, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal filtered data: %w", err)
	}

	return string(result), nil
}

// parseFieldList parses a comma-separated field list into a hierarchical structure
func parseFieldList(fieldsStr string) map[string]interface{} {
	fields := make(map[string]interface{})

	for _, field := range strings.Split(fieldsStr, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		// Split by dot for nested fields
		parts := strings.Split(field, ".")
		current := fields

		for i, part := range parts {
			if i == len(parts)-1 {
				// Leaf node - mark as included
				current[part] = true
			} else {
				// Intermediate node - create nested map if needed
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				// Move to nested map
				if nested, ok := current[part].(map[string]interface{}); ok {
					current = nested
				}
			}
		}
	}

	return fields
}

// filterObject recursively filters an object based on the field specification
func filterObject(data interface{}, fields map[string]interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		return filterMap(v, fields)
	case []interface{}:
		return filterArray(v, fields)
	default:
		return v
	}
}

// filterMap filters a map object
func filterMap(data map[string]interface{}, fields map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		// No fields specified, include everything
		return data
	}

	result := make(map[string]interface{})

	for key, value := range data {
		if fieldSpec, exists := fields[key]; exists {
			switch spec := fieldSpec.(type) {
			case bool:
				// Simple field - include as-is
				if spec {
					result[key] = value
				}
			case map[string]interface{}:
				// Nested field - recursively filter
				filtered := filterObject(value, spec)
				result[key] = filtered
			}
		}
	}

	return result
}

// filterArray filters an array by applying the same filter to each element
func filterArray(data []interface{}, fields map[string]interface{}) []interface{} {
	result := make([]interface{}, len(data))

	for i, item := range data {
		result[i] = filterObject(item, fields)
	}

	return result
}
