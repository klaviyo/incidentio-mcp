package tools

import (
	"encoding/json"
	"fmt"

	"github.com/incident-io/incidentio-mcp-golang/internal/incidentio"
)

// ListIncidentRolesTool lists available incident roles
type ListIncidentRolesTool struct {
	client *incidentio.Client
}

func NewListIncidentRolesTool(client *incidentio.Client) *ListIncidentRolesTool {
	return &ListIncidentRolesTool{client: client}
}

func (t *ListIncidentRolesTool) Name() string {
	return "list_available_incident_roles"
}

func (t *ListIncidentRolesTool) Description() string {
	return `List all available incident roles that can be assigned to users during incidents.

USAGE WORKFLOW:
1. Call to see all role types (incident lead, communications lead, etc.)
2. Use role IDs when assigning roles with assign_incident_role
3. Review role names, descriptions, and types

PARAMETERS:
- page_size: Number of results (default 25, max 250)

EXAMPLES:
- List all roles: {}
- List with pagination: {"page_size": 50}

IMPORTANT: Role IDs from this tool are required for the assign_incident_role tool.`
}

func (t *ListIncidentRolesTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     25,
			},
		},
	}
}

func (t *ListIncidentRolesTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListIncidentRolesOptions{}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	resp, err := t.client.ListIncidentRoles(opts)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}

// ListUsersTool lists available users for role assignment
type ListUsersTool struct {
	client *incidentio.Client
}

func NewListUsersTool(client *incidentio.Client) *ListUsersTool {
	return &ListUsersTool{client: client}
}

func (t *ListUsersTool) Name() string {
	return "list_users"
}

func (t *ListUsersTool) Description() string {
	return `List all users available for incident role assignment (automatically paginated).

USAGE WORKFLOW:
1. Call to see all users in your organization
2. Optional: Filter by email to find specific user
3. Use user IDs when assigning roles with assign_incident_role

PARAMETERS:
- page_size: Number of results (default 250, max 250)
- email: Optional. Filter users by email address

EXAMPLES:
- List all users: {}
- Find by email: {"email": "user@example.com"}

IMPORTANT: User IDs from this tool are required for the assign_incident_role tool.`
}

func (t *ListUsersTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"page_size": map[string]interface{}{
				"type":        "integer",
				"description": "Number of results per page (max 250)",
				"default":     250,
			},
			"email": map[string]interface{}{
				"type":        "string",
				"description": "Filter users by email address",
			},
		},
		"additionalProperties": false,
	}
}

func (t *ListUsersTool) Execute(args map[string]interface{}) (string, error) {
	opts := &incidentio.ListUsersOptions{}

	if pageSize, ok := args["page_size"].(float64); ok {
		opts.PageSize = int(pageSize)
	}

	if email, ok := args["email"].(string); ok && email != "" {
		opts.Email = email
	}

	resp, err := t.client.ListUsers(opts)
	if err != nil {
		return "", err
	}

	// Add a helpful message about the results
	var output string
	if opts.Email != "" {
		output = fmt.Sprintf("Users matching email '%s':\n", opts.Email)
	} else {
		output = fmt.Sprintf("Found %d users:\n", len(resp.Users))
	}

	// Format users in a more readable way
	for _, user := range resp.Users {
		output += fmt.Sprintf("\n- Name: %s\n  Email: %s\n  ID: %s\n  Role: %s\n",
			user.Name, user.Email, user.ID, user.Role)
	}

	// Also include the raw JSON
	jsonResult, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return output, nil // Return readable output even if JSON fails
	}

	output += "\n\nRaw JSON response:\n" + string(jsonResult)

	return output, nil
}

// AssignIncidentRoleTool assigns a role to a user for an incident
type AssignIncidentRoleTool struct {
	client *incidentio.Client
}

func NewAssignIncidentRoleTool(client *incidentio.Client) *AssignIncidentRoleTool {
	return &AssignIncidentRoleTool{client: client}
}

func (t *AssignIncidentRoleTool) Name() string {
	return "assign_incident_role"
}

func (t *AssignIncidentRoleTool) Description() string {
	return `Assign a specific incident role to a user for an incident.

USAGE WORKFLOW:
1. First call 'list_available_incident_roles' to get role IDs
2. Call 'list_users' to get user IDs (or filter by email)
3. Call this tool with incident ID, role ID, and user ID

PARAMETERS:
- id: Required. The incident ID to assign role for
- incident_role_id: Required. The role ID (from list_available_incident_roles)
- user_id: Required. The user ID (from list_users)

EXAMPLES:
- Assign lead: {"id": "01HXYZ...", "incident_role_id": "role_123", "user_id": "user_456"}

IMPORTANT: Use list_available_incident_roles and list_users to discover valid IDs before calling this tool.`
}

func (t *AssignIncidentRoleTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The incident ID",
			},
			"incident_role_id": map[string]interface{}{
				"type":        "string",
				"description": "The incident role ID to assign",
			},
			"user_id": map[string]interface{}{
				"type":        "string",
				"description": "The user ID to assign the role to",
			},
		},
		"required":             []interface{}{"id", "incident_role_id", "user_id"},
		"additionalProperties": false,
	}
}

func (t *AssignIncidentRoleTool) Execute(args map[string]interface{}) (string, error) {
	argDetails := make(map[string]interface{})
	for key, value := range args {
		argDetails[key] = value
	}

	if len(args) == 0 {
		return "", fmt.Errorf("no parameters provided")
	}

	incidentID, ok := args["id"].(string)
	if !ok || incidentID == "" {
		return "", fmt.Errorf("id parameter is required and must be a non-empty string. Received parameters: %+v", argDetails)
	}

	roleID, ok := args["incident_role_id"].(string)
	if !ok {
		return "", fmt.Errorf("incident_role_id parameter is required")
	}

	userID, ok := args["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id parameter is required")
	}

	// Create role assignment request using UpdateIncident
	req := &incidentio.UpdateIncidentRequest{
		IncidentRoleAssignments: []incidentio.CreateRoleAssignmentRequest{
			{
				IncidentRoleID: roleID,
				UserID:         userID,
			},
		},
	}

	incident, err := t.client.UpdateIncident(incidentID, req)
	if err != nil {
		return "", err
	}

	// Return just the role assignments part for clarity
	roleAssignments := make([]map[string]interface{}, 0)
	for _, assignment := range incident.IncidentRoleAssignments {
		roleData := map[string]interface{}{
			"role": map[string]interface{}{
				"id":          assignment.Role.ID,
				"name":        assignment.Role.Name,
				"description": assignment.Role.Description,
				"role_type":   assignment.Role.RoleType,
			},
		}

		if assignment.Assignee != nil {
			roleData["assignee"] = map[string]interface{}{
				"id":    assignment.Assignee.ID,
				"name":  assignment.Assignee.Name,
				"email": assignment.Assignee.Email,
			}
		}

		roleAssignments = append(roleAssignments, roleData)
	}

	response := map[string]interface{}{
		"message":          fmt.Sprintf("Successfully assigned role to user for incident %s", incident.Name),
		"incident_id":      incident.ID,
		"incident_name":    incident.Name,
		"role_assignments": roleAssignments,
	}

	result, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format response: %w", err)
	}

	return string(result), nil
}
