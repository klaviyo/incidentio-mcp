package tools

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListCatalogTypesTool lists available catalog types
type ListCatalogTypesTool struct {
	client *incidentio.Client
}

func NewListCatalogTypesTool(client *incidentio.Client) *ListCatalogTypesTool {
	return &ListCatalogTypesTool{client: client}
}

func (t *ListCatalogTypesTool) Name() string {
	return "list_catalog_types"
}

func (t *ListCatalogTypesTool) Description() string {
	return `List available catalog types in incident.io (automatically filtered to Custom* types only).

USAGE WORKFLOW:
1. Call to see all custom catalog types configured in your organization
2. Review type IDs, names, and attributes for each catalog type
3. Use catalog type IDs with list_catalog_entries to see entries

PARAMETERS:
- None required

EXAMPLES:
- List all custom catalog types: {}

IMPORTANT: This tool automatically filters to show only catalog types with TypeName starting with 'Custom' (case-insensitive). This filtering helps focus on user-defined catalogs rather than system catalogs.`
}

func (t *ListCatalogTypesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":                 "object",
		"properties":           map[string]interface{}{},
		"additionalProperties": false,
	}
}

func (t *ListCatalogTypesTool) Execute(args map[string]interface{}) (string, error) {
	result, err := t.client.ListCatalogTypes()
	if err != nil {
		return "", fmt.Errorf("failed to list catalog types: %w", err)
	}

	// Filter catalog types to only include those with TypeName starting with "Custom" (case-insensitive)
	var filteredTypes []incidentio.CatalogType
	for _, catalogType := range result.CatalogTypes {
		if strings.HasPrefix(strings.ToLower(catalogType.TypeName), "custom") {
			filteredTypes = append(filteredTypes, catalogType)
		}
	}

	output := fmt.Sprintf("Found %d catalog types (filtered for Custom* names):\n\n", len(filteredTypes))

	for _, catalogType := range filteredTypes {
		output += fmt.Sprintf("ID: %s\n", catalogType.ID)
		output += fmt.Sprintf("Name: %s\n", catalogType.Name)
		output += fmt.Sprintf("Type Name: %s\n", catalogType.TypeName)
		if catalogType.Description != "" {
			output += fmt.Sprintf("Description: %s\n", catalogType.Description)
		}
		if catalogType.Color != "" {
			output += fmt.Sprintf("Color: %s\n", catalogType.Color)
		}
		if catalogType.Icon != "" {
			output += fmt.Sprintf("Icon: %s\n", catalogType.Icon)
		}
		if len(catalogType.Attributes) > 0 {
			output += fmt.Sprintf("Attributes (%d):\n", len(catalogType.Attributes))
			for _, attr := range catalogType.Attributes {
				output += fmt.Sprintf("  - %s (%s): %s\n", attr.Name, attr.Type, attr.ID)
			}
		}
		output += fmt.Sprintf("Created: %s\n", catalogType.CreatedAt.Format("2006-01-02 15:04:05"))
		output += fmt.Sprintf("Updated: %s\n", catalogType.UpdatedAt.Format("2006-01-02 15:04:05"))
		output += "\n"
	}

	// Also return the raw JSON (only filtered types)
	filteredResult := &incidentio.ListCatalogTypesResponse{
		CatalogTypes: filteredTypes,
		ListResponse: result.ListResponse,
	}
	jsonOutput, err := json.MarshalIndent(filteredResult, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}

// ListCatalogEntriesTool lists catalog entries for a given type
type ListCatalogEntriesTool struct {
	client *incidentio.Client
}

func NewListCatalogEntriesTool(client *incidentio.Client) *ListCatalogEntriesTool {
	return &ListCatalogEntriesTool{client: client}
}

func (t *ListCatalogEntriesTool) Name() string {
	return "list_catalog_entries"
}

func (t *ListCatalogEntriesTool) Description() string {
	return `List catalog entries for a specific catalog type.

USAGE WORKFLOW:
1. First call 'list_catalog_types' to get available catalog type IDs
2. Call this tool with a catalog_type_id to see all entries
3. Optional: Filter by identifier to find specific entries
4. Review entry IDs, names, and attribute values

PARAMETERS:
- catalog_type_id: Required. The catalog type ID (from list_catalog_types)
- page_size: Optional. Number of entries per page (default: 25)
- after: Optional. Pagination cursor for next page
- identifier: Optional. Filter by identifier string

EXAMPLES:
- List all entries: {"catalog_type_id": "01HXYZ..."}
- Filter by identifier: {"catalog_type_id": "01HXYZ...", "identifier": "prod-api"}
- Paginated list: {"catalog_type_id": "01HXYZ...", "page_size": 50}`
}

func (t *ListCatalogEntriesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"catalog_type_id": map[string]interface{}{
				"type":        "string",
				"description": "The catalog type ID to list entries for",
			},
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of entries to return per page (default: 25)",
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor for next page",
			},
			"identifier": map[string]interface{}{
				"type":        "string",
				"description": "Filter by identifier",
			},
		},
		"required":             []interface{}{"catalog_type_id"},
		"additionalProperties": false,
	}
}

func (t *ListCatalogEntriesTool) Execute(args map[string]interface{}) (string, error) {
	catalogTypeID, ok := args["catalog_type_id"].(string)
	if !ok || catalogTypeID == "" {
		return "", fmt.Errorf("catalog_type_id parameter is required")
	}

	opts := incidentio.ListCatalogEntriesOptions{
		CatalogTypeID: catalogTypeID,
	}

	if pageSize, ok := args["page_size"]; ok {
		if ps, ok := pageSize.(float64); ok {
			opts.PageSize = int(ps)
		} else if ps, ok := pageSize.(string); ok {
			if parsed, err := strconv.Atoi(ps); err == nil {
				opts.PageSize = parsed
			}
		}
	}

	if after, ok := args["after"].(string); ok {
		opts.After = after
	}

	if identifier, ok := args["identifier"].(string); ok {
		opts.Identifier = identifier
	}

	result, err := t.client.ListCatalogEntries(opts)
	if err != nil {
		return "", fmt.Errorf("failed to list catalog entries: %w", err)
	}

	output := fmt.Sprintf("Found %d catalog entries for type %s:\n\n", len(result.CatalogEntries), catalogTypeID)

	for _, entry := range result.CatalogEntries {
		output += fmt.Sprintf("ID: %s\n", entry.ID)
		output += fmt.Sprintf("Name: %s\n", entry.Name)
		if len(entry.Aliases) > 0 {
			output += fmt.Sprintf("Aliases: %v\n", entry.Aliases)
		}
		if entry.ExternalID != "" {
			output += fmt.Sprintf("External ID: %s\n", entry.ExternalID)
		}
		output += fmt.Sprintf("Rank: %d\n", entry.Rank)
		if len(entry.AttributeValues) > 0 {
			output += "Attributes:\n"
			for key, value := range entry.AttributeValues {
				if value.Value != nil {
					if value.Value.Literal != "" {
						output += fmt.Sprintf("  %s: %s\n", key, value.Value.Literal)
					} else if value.Value.ID != "" {
						output += fmt.Sprintf("  %s: %s (ID)\n", key, value.Value.ID)
					}
				}
				if len(value.ArrayValue) > 0 {
					output += fmt.Sprintf("  %s: [", key)
					for i, v := range value.ArrayValue {
						if i > 0 {
							output += ", "
						}
						if v.Literal != "" {
							output += v.Literal
						} else if v.ID != "" {
							output += v.ID + " (ID)"
						}
					}
					output += "]\n"
				}
			}
		}
		output += fmt.Sprintf("Created: %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))
		output += fmt.Sprintf("Updated: %s\n", entry.UpdatedAt.Format("2006-01-02 15:04:05"))
		output += "\n"
	}

	// Add pagination info
	if result.PaginationMeta.After != "" {
		output += fmt.Sprintf("Pagination: Next page available (after: %s)\n", result.PaginationMeta.After)
	}
	output += fmt.Sprintf("Total entries: %d\n", result.PaginationMeta.TotalCount)

	// Also return the raw JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}

// UpdateCatalogEntryTool updates a catalog entry
type UpdateCatalogEntryTool struct {
	client *incidentio.Client
}

func NewUpdateCatalogEntryTool(client *incidentio.Client) *UpdateCatalogEntryTool {
	return &UpdateCatalogEntryTool{client: client}
}

func (t *UpdateCatalogEntryTool) Name() string {
	return "update_catalog_entry"
}

func (t *UpdateCatalogEntryTool) Description() string {
	return `Update an existing catalog entry's properties and attribute values.

USAGE WORKFLOW:
1. First call 'list_catalog_entries' to find the entry you want to update
2. Prepare updated values (name, aliases, rank, attribute_values)
3. Call this tool with the entry ID and new values
4. For attributes, specify which attributes to update in update_attributes array

PARAMETERS:
- id: Required. The catalog entry ID to update
- name: Optional. New name for the entry
- aliases: Optional. Array of alias strings
- external_id: Optional. External system ID
- rank: Optional. Sort order/rank integer
- attribute_values: Optional. Object mapping attribute IDs to values
- update_attributes: Optional. Array of attribute IDs to update

EXAMPLES:
- Update name: {"id": "entry_123", "name": "New Name"}
- Update attributes: {"id": "entry_123", "attribute_values": {"attr_abc": {"value": {"literal": "new value"}}}, "update_attributes": ["attr_abc"]}
- Update rank: {"id": "entry_123", "rank": 10}`
}

func (t *UpdateCatalogEntryTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The catalog entry ID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "New name for the catalog entry",
			},
			"aliases": map[string]interface{}{
				"type":        "array",
				"description": "List of aliases for the catalog entry",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"external_id": map[string]interface{}{
				"type":        "string",
				"description": "External ID for the catalog entry",
			},
			"rank": map[string]interface{}{
				"type":        "integer",
				"description": "Rank/order of the catalog entry",
			},
			"attribute_values": map[string]interface{}{
				"type":        "object",
				"description": "Attribute values as a JSON object",
			},
			"update_attributes": map[string]interface{}{
				"type":        "array",
				"description": "List of attribute IDs to update",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required":             []interface{}{"id"},
		"additionalProperties": false,
	}
}

func (t *UpdateCatalogEntryTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("id parameter is required")
	}

	req := incidentio.UpdateCatalogEntryRequest{}

	if name, ok := args["name"].(string); ok {
		req.Name = name
	}

	if aliases, ok := args["aliases"].([]interface{}); ok {
		req.Aliases = make([]string, len(aliases))
		for i, alias := range aliases {
			if s, ok := alias.(string); ok {
				req.Aliases[i] = s
			}
		}
	}

	if externalID, ok := args["external_id"].(string); ok {
		req.ExternalID = externalID
	}

	if rank, ok := args["rank"]; ok {
		if r, ok := rank.(float64); ok {
			req.Rank = int(r)
		} else if r, ok := rank.(string); ok {
			if parsed, err := strconv.Atoi(r); err == nil {
				req.Rank = parsed
			}
		}
	}

	if attrValues, ok := args["attribute_values"].(map[string]interface{}); ok {
		req.AttributeValues = make(map[string]incidentio.CatalogEntryAttributeValue)
		for key, value := range attrValues {
			if valueMap, ok := value.(map[string]interface{}); ok {
				attrValue := incidentio.CatalogEntryAttributeValue{}

				// Handle single value
				if v, ok := valueMap["value"].(map[string]interface{}); ok {
					attrValue.Value = &incidentio.CatalogEntryAttributeValueItem{}
					if literal, ok := v["literal"].(string); ok {
						attrValue.Value.Literal = literal
					}
					if id, ok := v["id"].(string); ok {
						attrValue.Value.ID = id
					}
				}

				// Handle array value
				if arrayValue, ok := valueMap["array_value"].([]interface{}); ok {
					attrValue.ArrayValue = make([]incidentio.CatalogEntryAttributeValueItem, len(arrayValue))
					for i, item := range arrayValue {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if literal, ok := itemMap["literal"].(string); ok {
								attrValue.ArrayValue[i].Literal = literal
							}
							if id, ok := itemMap["id"].(string); ok {
								attrValue.ArrayValue[i].ID = id
							}
						}
					}
				}

				req.AttributeValues[key] = attrValue
			}
		}
	}

	if updateAttrs, ok := args["update_attributes"].([]interface{}); ok {
		req.UpdateAttributes = make([]string, len(updateAttrs))
		for i, attr := range updateAttrs {
			if s, ok := attr.(string); ok {
				req.UpdateAttributes[i] = s
			}
		}
	}

	result, err := t.client.UpdateCatalogEntry(id, req)
	if err != nil {
		return "", fmt.Errorf("failed to update catalog entry: %w", err)
	}

	output := "Updated catalog entry:\n\n"
	output += fmt.Sprintf("ID: %s\n", result.ID)
	output += fmt.Sprintf("Name: %s\n", result.Name)
	if len(result.Aliases) > 0 {
		output += fmt.Sprintf("Aliases: %v\n", result.Aliases)
	}
	if result.ExternalID != "" {
		output += fmt.Sprintf("External ID: %s\n", result.ExternalID)
	}
	output += fmt.Sprintf("Rank: %d\n", result.Rank)
	output += fmt.Sprintf("Created: %s\n", result.CreatedAt.Format("2006-01-02 15:04:05"))
	output += fmt.Sprintf("Updated: %s\n", result.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Also return the raw JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return output, nil
	}

	return output + "\nRaw JSON:\n" + string(jsonOutput), nil
}
