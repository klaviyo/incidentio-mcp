package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

// ListAlertRoutesTool lists alert routes from incident.io
type ListAlertRoutesTool struct {
	apiClient *client.Client
}

func NewListAlertRoutesTool(c *client.Client) *ListAlertRoutesTool {
	return &ListAlertRoutesTool{apiClient: c}
}

func (t *ListAlertRoutesTool) Name() string {
	return "list_alert_routes"
}

func (t *ListAlertRoutesTool) Description() string {
	return `List alert routes that define how alerts are routed and escalated.

USAGE WORKFLOW:
1. Call to see all configured alert routes
2. Review conditions, escalations, and grouping keys for each route
3. Use route IDs with get_alert_route for detailed configuration

PARAMETERS:
- page_size: Number of results per page (1-250)
- after: Pagination cursor for next page

EXAMPLES:
- List all routes: {}
- List with custom page size: {"page_size": 50}`
}

func (t *ListAlertRoutesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page",
				"minimum":     1,
				"maximum":     250,
			},
			"after": map[string]interface{}{
				"type":        "string",
				"description": "Pagination cursor for next page",
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListAlertRoutesTool) Execute(args map[string]interface{}) (string, error) {
	params := &client.ListAlertRoutesParams{}

	if pageSize, ok := args["page_size"].(float64); ok {
		params.PageSize = int(pageSize)
	}
	if after, ok := args["after"].(string); ok {
		params.After = after
	}

	result, err := t.apiClient.ListAlertRoutes(params)
	if err != nil {
		return "", fmt.Errorf("failed to list alert routes: %w", err)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}

// GetAlertRouteTool gets details of a specific alert route
type GetAlertRouteTool struct {
	apiClient *client.Client
}

func NewGetAlertRouteTool(c *client.Client) *GetAlertRouteTool {
	return &GetAlertRouteTool{apiClient: c}
}

func (t *GetAlertRouteTool) Name() string {
	return "get_alert_route"
}

func (t *GetAlertRouteTool) Description() string {
	return `Get detailed configuration of a specific alert route.

USAGE WORKFLOW:
1. Get route ID from list_alert_routes
2. Call this tool for complete route configuration
3. Review conditions, escalations, and grouping settings

PARAMETERS:
- id: Required. The alert route ID to retrieve

EXAMPLES:
- Get route: {"id": "route_123"}`
}

func (t *GetAlertRouteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The alert route ID",
				"minLength":   1,
			},
		},
		"required":             []string{"id"},
		"additionalProperties": false,
	}
}

func (t *GetAlertRouteTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("alert route ID is required")
	}

	alertRoute, err := t.apiClient.GetAlertRoute(id)
	if err != nil {
		return "", fmt.Errorf("failed to get alert route: %w", err)
	}

	output, err := json.MarshalIndent(alertRoute, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}

// CreateAlertRouteTool creates a new alert route
type CreateAlertRouteTool struct {
	apiClient *client.Client
}

func NewCreateAlertRouteTool(c *client.Client) *CreateAlertRouteTool {
	return &CreateAlertRouteTool{apiClient: c}
}

func (t *CreateAlertRouteTool) Name() string {
	return "create_alert_route"
}

func (t *CreateAlertRouteTool) Description() string {
	return `Create a new alert route to define how alerts are routed and escalated based on conditions.

USAGE WORKFLOW:
1. Define routing conditions (field, operation, value)
2. Configure escalation bindings (which escalation paths to use)
3. Optional: Set grouping keys to group similar alerts
4. Optional: Configure incident template for auto-creation

PARAMETERS:
- name: Required. Name for the alert route
- enabled: Optional. Whether route is active (default: true)
- conditions: Required. Array of condition objects with field, operation, value
- escalations: Required. Array of escalation bindings with id and level
- grouping_keys: Optional. Array of field names to group alerts by
- template: Optional. Incident template for auto-creating incidents

EXAMPLES:
- Basic route: {"name": "Production alerts", "conditions": [{"field": "severity", "operation": "equals", "value": "critical"}], "escalations": [{"id": "esc_123", "level": 1}]}
- With grouping: {"name": "Service alerts", "conditions": [...], "escalations": [...], "grouping_keys": ["service_name", "environment"]}`
}

func (t *CreateAlertRouteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the alert route",
				"minLength":   1,
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the alert route should be enabled",
				"default":     true,
			},
			"conditions": map[string]interface{}{
				"type":        "array",
				"description": "Conditions for routing alerts",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"field": map[string]interface{}{
							"type":        "string",
							"description": "Field to match on",
						},
						"operation": map[string]interface{}{
							"type":        "string",
							"description": "Operation to perform (equals, contains, etc)",
						},
						"value": map[string]interface{}{
							"type":        "string",
							"description": "Value to match against",
						},
					},
					"required": []string{"field", "operation", "value"},
				},
			},
			"escalations": map[string]interface{}{
				"type":        "array",
				"description": "Escalation bindings",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Escalation ID",
						},
						"level": map[string]interface{}{
							"type":        "integer",
							"description": "Escalation level",
						},
					},
					"required": []string{"id", "level"},
				},
			},
			"grouping_keys": map[string]interface{}{
				"type":        "array",
				"description": "Keys to group alerts by",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"template": map[string]interface{}{
				"type":        "object",
				"description": "Template for creating incidents from alerts",
			},
		},
		"required":             []string{"name", "conditions", "escalations"},
		"additionalProperties": false,
	}
}

func (t *CreateAlertRouteTool) Execute(args map[string]interface{}) (string, error) {
	req := &client.CreateAlertRouteRequest{}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("name is required")
	}
	req.Name = name

	if enabled, ok := args["enabled"].(bool); ok {
		req.Enabled = enabled
	} else {
		req.Enabled = true // default to enabled
	}

	// Parse conditions
	if conditions, ok := args["conditions"].([]interface{}); ok {
		for _, c := range conditions {
			if cond, ok := c.(map[string]interface{}); ok {
				condition := client.AlertCondition{
					Field:     cond["field"].(string),
					Operation: cond["operation"].(string),
					Value:     cond["value"].(string),
				}
				req.Conditions = append(req.Conditions, condition)
			}
		}
	}

	// Parse escalations
	if escalations, ok := args["escalations"].([]interface{}); ok {
		for _, e := range escalations {
			if esc, ok := e.(map[string]interface{}); ok {
				escalation := client.EscalationBinding{
					ID:    esc["id"].(string),
					Level: int(esc["level"].(float64)),
				}
				req.Escalations = append(req.Escalations, escalation)
			}
		}
	}

	// Parse grouping keys
	if groupingKeys, ok := args["grouping_keys"].([]interface{}); ok {
		for _, k := range groupingKeys {
			if key, ok := k.(string); ok {
				req.GroupingKeys = append(req.GroupingKeys, key)
			}
		}
	}

	// Parse template
	if template, ok := args["template"].(map[string]interface{}); ok {
		req.Template = template
	}

	alertRoute, err := t.apiClient.CreateAlertRoute(req)
	if err != nil {
		return "", fmt.Errorf("failed to create alert route: %w", err)
	}

	output, err := json.MarshalIndent(alertRoute, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}

// UpdateAlertRouteTool updates an alert route
type UpdateAlertRouteTool struct {
	apiClient *client.Client
}

func NewUpdateAlertRouteTool(c *client.Client) *UpdateAlertRouteTool {
	return &UpdateAlertRouteTool{apiClient: c}
}

func (t *UpdateAlertRouteTool) Name() string {
	return "update_alert_route"
}

func (t *UpdateAlertRouteTool) Description() string {
	return `Update an existing alert route's configuration (name, enabled status, conditions, escalations).

USAGE WORKFLOW:
1. First call 'get_alert_route' to see current configuration
2. Modify desired fields
3. Call update with route ID and new configuration

PARAMETERS:
- id: Required. The alert route ID to update
- name: Optional. New name for the route
- enabled: Optional. Enable or disable the route
- conditions: Optional. New array of routing conditions
- escalations: Optional. New array of escalation bindings
- grouping_keys: Optional. New array of grouping keys
- template: Optional. New incident template

EXAMPLES:
- Disable route: {"id": "route_123", "enabled": false}
- Update conditions: {"id": "route_123", "conditions": [{"field": "severity", "operation": "equals", "value": "high"}]}`
}

func (t *UpdateAlertRouteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The alert route ID to update",
				"minLength":   1,
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "New name for the alert route",
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the alert route should be enabled",
			},
			"conditions": map[string]interface{}{
				"type":        "array",
				"description": "New conditions for routing alerts",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"field": map[string]interface{}{
							"type":        "string",
							"description": "Field to match on",
						},
						"operation": map[string]interface{}{
							"type":        "string",
							"description": "Operation to perform",
						},
						"value": map[string]interface{}{
							"type":        "string",
							"description": "Value to match against",
						},
					},
					"required": []string{"field", "operation", "value"},
				},
			},
			"escalations": map[string]interface{}{
				"type":        "array",
				"description": "New escalation bindings",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Escalation ID",
						},
						"level": map[string]interface{}{
							"type":        "integer",
							"description": "Escalation level",
						},
					},
					"required": []string{"id", "level"},
				},
			},
			"grouping_keys": map[string]interface{}{
				"type":        "array",
				"description": "Keys to group alerts by",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"template": map[string]interface{}{
				"type":        "object",
				"description": "Template for creating incidents from alerts",
			},
		},
		"required":             []string{"id"},
		"additionalProperties": false,
	}
}

func (t *UpdateAlertRouteTool) Execute(args map[string]interface{}) (string, error) {
	id, ok := args["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("alert route ID is required")
	}

	req := &client.UpdateAlertRouteRequest{}

	if name, ok := args["name"].(string); ok {
		req.Name = name
	}

	if enabled, ok := args["enabled"].(bool); ok {
		req.Enabled = &enabled
	}

	// Parse conditions
	if conditions, ok := args["conditions"].([]interface{}); ok {
		req.Conditions = []client.AlertCondition{}
		for _, c := range conditions {
			if cond, ok := c.(map[string]interface{}); ok {
				condition := client.AlertCondition{
					Field:     cond["field"].(string),
					Operation: cond["operation"].(string),
					Value:     cond["value"].(string),
				}
				req.Conditions = append(req.Conditions, condition)
			}
		}
	}

	// Parse escalations
	if escalations, ok := args["escalations"].([]interface{}); ok {
		req.Escalations = []client.EscalationBinding{}
		for _, e := range escalations {
			if esc, ok := e.(map[string]interface{}); ok {
				escalation := client.EscalationBinding{
					ID:    esc["id"].(string),
					Level: int(esc["level"].(float64)),
				}
				req.Escalations = append(req.Escalations, escalation)
			}
		}
	}

	// Parse grouping keys
	if groupingKeys, ok := args["grouping_keys"].([]interface{}); ok {
		req.GroupingKeys = []string{}
		for _, k := range groupingKeys {
			if key, ok := k.(string); ok {
				req.GroupingKeys = append(req.GroupingKeys, key)
			}
		}
	}

	// Parse template
	if template, ok := args["template"].(map[string]interface{}); ok {
		req.Template = template
	}

	alertRoute, err := t.apiClient.UpdateAlertRoute(id, req)
	if err != nil {
		return "", fmt.Errorf("failed to update alert route: %w", err)
	}

	output, err := json.MarshalIndent(alertRoute, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(output), nil
}
