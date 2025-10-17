package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// FilterFields filters a JSON object to only include the specified fields.
// Fields can be specified as:
// - Top-level fields: "id", "name", "summary"
// - Nested fields with dot notation: "severity.name", "incident_status.category"
// - Array elements are filtered recursively
//
// For API responses with collection fields (incidents, alerts), the field filter
// is automatically applied to the items in the collection, not the response wrapper.
//
// Example:
//   fields := "id,name,severity.name,incident_status.category"
//   filtered, err := FilterFields(data, fields)
func FilterFields(data interface{}, fieldsStr string) (string, error) {
	log.Printf("[FilterFields] START - fieldsStr=%q", fieldsStr)

	if fieldsStr == "" {
		log.Printf("[FilterFields] No fields specified, returning all data")
		// No filtering requested, return original data
		result, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal data: %w", err)
		}
		return string(result), nil
	}

	// Parse field list
	fields := parseFieldList(fieldsStr)
	log.Printf("[FilterFields] Parsed fields structure: %+v", fields)

	// Marshal to JSON first to get map representation
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}
	log.Printf("[FilterFields] JSON bytes length: %d", len(jsonBytes))

	var rawData interface{}
	if err := json.Unmarshal(jsonBytes, &rawData); err != nil {
		return "", fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Check if this is a collection response (incidents, alerts, etc.)
	// If so, apply filtering to the collection items, not the response wrapper
	if dataMap, ok := rawData.(map[string]interface{}); ok {
		log.Printf("[FilterFields] Data is a map with keys: %v", getKeys(dataMap))

		if incidents, hasIncidents := dataMap["incidents"]; hasIncidents {
			log.Printf("[FilterFields] Found incidents collection")
			if incArray, ok := incidents.([]interface{}); ok {
				log.Printf("[FilterFields] Incidents array has %d items", len(incArray))
				if len(incArray) > 0 {
					firstInc, _ := json.Marshal(incArray[0])
					log.Printf("[FilterFields] First incident sample: %s", string(firstInc))
				}
			}

			// Filter the incidents array
			filteredIncidents := filterObject(incidents, fields)
			log.Printf("[FilterFields] Filtered incidents: %+v", filteredIncidents)

			// Preserve the response structure with filtered incidents
			filtered := map[string]interface{}{
				"incidents": filteredIncidents,
			}
			// Include pagination_meta if present
			if paginationMeta, hasPagination := dataMap["pagination_meta"]; hasPagination {
				filtered["pagination_meta"] = paginationMeta
			}
			result, err := json.MarshalIndent(filtered, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal filtered data: %w", err)
			}
			log.Printf("[FilterFields] END - returning %d bytes", len(result))
			return string(result), nil
		}

		if alerts, hasAlerts := dataMap["alerts"]; hasAlerts {
			log.Printf("[FilterFields] Found alerts collection")
			// Filter the alerts array
			filteredAlerts := filterObject(alerts, fields)
			// Preserve the response structure with filtered alerts
			filtered := map[string]interface{}{
				"alerts": filteredAlerts,
			}
			// Include pagination_meta if present
			if paginationMeta, hasPagination := dataMap["pagination_meta"]; hasPagination {
				filtered["pagination_meta"] = paginationMeta
			}
			result, err := json.MarshalIndent(filtered, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal filtered data: %w", err)
			}
			log.Printf("[FilterFields] END - returning %d bytes", len(result))
			return string(result), nil
		}

		log.Printf("[FilterFields] No collection found, filtering data map directly")
	}

	// Default behavior: filter the object directly
	filtered := filterObject(rawData, fields)
	log.Printf("[FilterFields] Filtered result: %+v", filtered)

	// Marshal the filtered result
	result, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal filtered data: %w", err)
	}

	log.Printf("[FilterFields] END - returning %d bytes", len(result))
	return string(result), nil
}

// Helper function to get map keys for logging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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
	log.Printf("[filterMap] Filtering map with %d keys, fields spec: %+v", len(data), fields)

	if len(fields) == 0 {
		// No fields specified, include everything
		log.Printf("[filterMap] No fields specified, returning all %d keys", len(data))
		return data
	}

	result := make(map[string]interface{})
	log.Printf("[filterMap] Data keys: %v", getKeys(data))
	log.Printf("[filterMap] Requested fields: %v", getKeys(fields))

	for key, value := range data {
		if fieldSpec, exists := fields[key]; exists {
			log.Printf("[filterMap] Field %q exists in spec: %+v", key, fieldSpec)
			switch spec := fieldSpec.(type) {
			case bool:
				// Simple field - include as-is
				if spec {
					log.Printf("[filterMap] Including simple field %q", key)
					result[key] = value
				}
			case map[string]interface{}:
				// Nested field - recursively filter
				log.Printf("[filterMap] Recursively filtering nested field %q", key)
				filtered := filterObject(value, spec)
				result[key] = filtered
			}
		} else {
			log.Printf("[filterMap] Field %q NOT in spec, skipping", key)
		}
	}

	log.Printf("[filterMap] Result has %d keys: %v", len(result), getKeys(result))
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
