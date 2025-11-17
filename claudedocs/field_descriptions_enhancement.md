# Field Descriptions Enhancement

## Overview

Enhanced the `fields` parameter descriptions across all tools to include comprehensive, type-based field documentation. This improvement makes the field filtering feature more discoverable and easier to use for LLMs.

## Changes Made

### Before
```go
"fields": map[string]interface{}{
    "type":        "string",
    "description": "Comma-separated list of fields... Example: \"id,name,severity.name\"",
}
```

### After
```go
"fields": map[string]interface{}{
    "type": "string",
    "description": `Comma-separated list of fields to include in response...

Available top-level fields: id, reference, name, summary, permalink, incident_status, severity, incident_type, mode, visibility, created_at, updated_at, slack_team_id, slack_channel_id, slack_channel_name, incident_role_assignments, custom_field_entries, has_debrief

Nested fields (dot notation):
- incident_status: id, name, description, category, rank, created_at, updated_at
- severity: id, name, description, rank, created_at, updated_at
- incident_type: id, name, description, is_default, private_incidents_only, create_in_triage, created_at, updated_at

Examples: "id,name,severity.name" or "id,reference,incident_status.category,severity.name"
Omit to return all fields.`,
}
```

## Benefits

### 1. **Self-Documenting API**
- LLMs can see exactly what fields are available without trial and error
- No need to consult external documentation
- Field lists are authoritative and type-based

### 2. **Improved Usability**
- Clear organization of top-level vs nested fields
- Grouped nested fields by parent object
- Multiple practical examples for different use cases

### 3. **Reduced Errors**
- Less guessing about field names
- Clear indication of nested field syntax
- Prevents typos and invalid field references

### 4. **Better LLM Experience**
- Tool descriptions provide complete context in a single view
- Nested field relationships are immediately clear
- Examples demonstrate real-world usage patterns

## Updated Tools

### Incident Tools

**list_incidents** - Enhanced with:
- 18 top-level incident fields
- 3 nested object breakdowns (incident_status, severity, incident_type)
- Multiple usage examples

**get_incident** - Same field documentation as list_incidents for consistency

### Alert Tools

**list_alerts** - Enhanced with:
- 9 top-level alert fields
- 2 nested object breakdowns (incident, merged_into_alert)
- Alert-specific usage examples

**get_alert** - Same field documentation as list_alerts for consistency

## Field Documentation Structure

Each enhanced description follows this pattern:

1. **Purpose Statement**: What the parameter does
2. **Available Top-Level Fields**: Exhaustive list from type definition
3. **Nested Fields**: Organized by parent object with dot notation syntax
4. **Examples**: Multiple real-world usage patterns
5. **Default Behavior**: What happens when omitted

## Implementation Notes

### Type-Based Generation

Field lists were extracted from Go struct JSON tags:

```go
// From types.go
type Incident struct {
    ID                      string             `json:"id"`
    Reference               string             `json:"reference"`
    Name                    string             `json:"name"`
    // ... more fields
}

// Becomes documentation
Available top-level fields: id, reference, name, ...
```

### Nested Field Organization

Nested objects are documented with clear parent-child relationships:

```
Nested fields (dot notation):
- incident_status: id, name, description, category, rank, created_at, updated_at
```

This makes it obvious that `incident_status.category` is valid but `incident_status.invalid_field` is not.

## Future Enhancements

### Potential Improvements

1. **Automatic Generation**: Script to extract field lists from Go types automatically
2. **Field Type Hints**: Include data types (string, int, bool, timestamp)
3. **Required vs Optional**: Mark which fields are always present vs optional
4. **Field Descriptions**: Add brief descriptions of what each field contains
5. **Common Patterns**: Pre-defined field sets like "minimal", "standard", "full"

### Example Auto-Generated Description

```go
// Hypothetical future enhancement
Available fields:
- id (string, required): Unique incident identifier
- name (string, required): Human-readable incident name
- summary (string, optional): Detailed incident description
- created_at (timestamp, required): Incident creation time
- severity (object): Nested severity information
  └─ severity.name (string): Severity level name
  └─ severity.rank (int): Severity ranking (1=highest)
```

## Migration Path

No breaking changes - this is purely an enhancement to documentation:

✅ Existing code continues to work unchanged
✅ Parameter behavior is identical
✅ Only descriptions were enhanced
✅ Backward compatible with all existing usage

## Testing

Verified with:
```bash
go test ./internal/tools/...     # All tests pass
go build ./cmd/mcp-server        # Build successful
```

No functional changes, only improved documentation in schema descriptions.

## Summary

This enhancement transforms the field filtering feature from a "nice to have" into a truly discoverable, self-documenting API feature. LLMs can now confidently select appropriate fields without guesswork, leading to:

- **Better context management**: More precise field selection
- **Fewer errors**: Clear field availability documentation
- **Improved UX**: Self-contained tool descriptions
- **Professional polish**: Production-quality documentation

The improvement required zero changes to the filtering logic itself - purely enhanced metadata that makes the existing feature significantly more usable.
