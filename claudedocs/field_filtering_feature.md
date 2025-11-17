# Field Filtering Feature

## Overview

The field filtering feature allows MCP tool users to selectively retrieve only the fields they need from API responses, significantly reducing context window usage when working with LLMs.

## Problem Statement

Many incident.io API responses contain extensive nested data structures with 25+ fields per object. When listing multiple incidents or retrieving detailed information, this can quickly consume large portions of an LLM's context window with unnecessary data.

## Solution

A reusable field filtering system that:
- Accepts a comma-separated string of field names
- Supports nested field selection using dot notation
- Filters JSON responses before returning them to the LLM
- Works consistently across all tools with large outputs

## Implementation

### Core Components

#### 1. Enhanced Field Parameter Descriptions

Each tool's `fields` parameter now includes comprehensive documentation:
- **Available top-level fields**: Complete list of all root-level fields from the response type
- **Nested field paths**: Organized breakdown of nested object fields with dot notation
- **Practical examples**: Real-world usage patterns for common scenarios
- **Type-based documentation**: Field lists are derived directly from Go struct JSON tags

This self-documenting approach ensures LLMs always have accurate, up-to-date information about available fields.

#### 2. Field Filter Utility (`internal/tools/fieldfilter.go`)

**`FilterFields(data interface{}, fieldsStr string) (string, error)`**
- Main entry point for field filtering
- Returns formatted JSON string containing only requested fields
- If `fieldsStr` is empty, returns all fields (backward compatible)

**Key Features:**
- **Top-level field filtering**: `"id,name,summary"`
- **Nested field filtering**: `"severity.name,incident_status.category"`
- **Array handling**: Filters applied recursively to array elements
- **Whitespace tolerance**: Handles spaces in field lists
- **JSON formatting**: Returns indented JSON for readability

### 2. Tools with Field Filtering Support

The following tools have been enhanced with field filtering:

#### Incident Tools
- **`list_incidents`**: Filter incident list responses
  - Example: `{"status": ["active"], "fields": "id,name,severity.name"}`

- **`get_incident`**: Filter single incident details
  - Example: `{"incident_id": "01HXYZ...", "fields": "id,name,summary,severity.name"}`

#### Alert Tools
- **`list_alerts`**: Filter alert list responses
  - Example: `{"fields": "id,title,status,incident.id"}`

- **`get_alert`**: Filter single alert details
  - Example: `{"id": "alert_123", "fields": "id,title,status"}`

## Usage Examples

### Basic Field Selection

```json
{
  "incident_id": "01HXYZ123",
  "fields": "id,name,reference"
}
```

**Before (full response ~500 tokens):**
```json
{
  "id": "01HXYZ123",
  "reference": "INC-123",
  "name": "Production Outage",
  "summary": "Database connection pool exhausted...",
  "permalink": "https://...",
  "incident_status": {
    "id": "status_1",
    "name": "Active",
    "description": "Incident is active",
    "category": "triage",
    "rank": 1,
    "created_at": "...",
    "updated_at": "..."
  },
  "severity": { /* ... */ },
  "incident_role_assignments": [ /* ... */ ],
  "custom_field_entries": [ /* ... */ ],
  // ... 20+ more fields
}
```

**After (filtered response ~50 tokens):**
```json
{
  "id": "01HXYZ123",
  "name": "Production Outage",
  "reference": "INC-123"
}
```

### Nested Field Selection

```json
{
  "incident_id": "01HXYZ123",
  "fields": "id,name,severity.name,incident_status.category"
}
```

**Result:**
```json
{
  "id": "01HXYZ123",
  "name": "Production Outage",
  "severity": {
    "name": "Critical"
  },
  "incident_status": {
    "category": "triage"
  }
}
```

### List Operations with Filtering

```json
{
  "status": ["active"],
  "fields": "incidents.id,incidents.name,incidents.severity.name"
}
```

Filters are applied to each incident in the list, reducing context usage proportionally.

## Performance Impact

### Context Window Savings

| Operation | Without Filtering | With Filtering | Savings |
|-----------|-------------------|----------------|---------|
| Single incident (full) | ~500 tokens | ~50 tokens | 90% |
| List 10 incidents (full) | ~5,000 tokens | ~500 tokens | 90% |
| List 25 incidents (full) | ~12,500 tokens | ~1,250 tokens | 90% |

### Real-World Example

**Scenario**: List active high-severity incidents, show only ID, name, and severity name

```json
{
  "status": ["active"],
  "severity": ["sev_1", "sev_2"],
  "fields": "incidents.id,incidents.name,incidents.severity.name"
}
```

- **Before**: 25 incidents × ~500 tokens = ~12,500 tokens
- **After**: 25 incidents × ~50 tokens = ~1,250 tokens
- **Saved**: ~11,250 tokens (90% reduction)

## Best Practices

### When to Use Field Filtering

✅ **DO use field filtering when:**
- Working with lists of resources
- Only specific fields are needed for the task
- Context window usage is a concern
- Building dashboards or reports with focused data

❌ **DON'T use field filtering when:**
- You need complete object details
- Debugging or investigating unknown issues
- First-time exploration of available data
- Small single-object queries where overhead isn't worth it

### Recommended Field Combinations

#### Quick Status Check
```
"id,name,incident_status.category,severity.name"
```

#### Summary Dashboard
```
"id,reference,name,permalink,created_at,severity.name,incident_status.name"
```

#### Detailed Investigation
```
"id,name,summary,severity,incident_status,incident_role_assignments"
```

## Technical Details

### Field Path Resolution

The filtering system uses dot notation to navigate nested structures:

1. **Parse**: Split field string by commas: `"id,severity.name"` → `["id", "severity.name"]`
2. **Build Tree**: Create hierarchical structure from dot notation
3. **Filter**: Recursively traverse and filter JSON based on tree structure
4. **Output**: Marshal filtered structure to formatted JSON

### Array Handling

When filtering arrays, the filter is applied to each element:

```json
// Input
{
  "incidents": [
    {"id": "1", "name": "A", "summary": "..."},
    {"id": "2", "name": "B", "summary": "..."}
  ]
}

// Filter: "incidents.id,incidents.name"
// Output
{
  "incidents": [
    {"id": "1", "name": "A"},
    {"id": "2", "name": "B"}
  ]
}
```

### Backward Compatibility

- **Omitting `fields` parameter**: Returns complete response (default behavior)
- **Empty `fields` string**: Returns complete response
- **Invalid field names**: Silently ignored (fields not present won't appear in output)
- **Existing tool parameters**: Unchanged and fully compatible

## Testing

Comprehensive test coverage in `internal/tools/fieldfilter_test.go`:

- ✅ No fields specified (returns all)
- ✅ Top-level field filtering
- ✅ Nested field filtering with dot notation
- ✅ Array filtering
- ✅ Whitespace handling
- ✅ Complex nested structures
- ✅ JSON formatting validation

Run tests:
```bash
go test -v ./internal/tools -run TestFilterFields
```

## Future Enhancements

Potential improvements for future iterations:

1. **Wildcard Support**: `"severity.*"` to include all severity fields
2. **Exclusion Syntax**: `"*,!summary,!custom_field_entries"` to exclude specific fields
3. **Preset Profiles**: Named field sets like `"@minimal"`, `"@dashboard"`, `"@full"`
4. **Performance Metrics**: Track and report filtering time and memory usage
5. **Field Discovery**: Auto-suggest commonly used field combinations

## Migration Guide

No breaking changes were introduced. To adopt field filtering:

1. **Update Tool Calls**: Add optional `fields` parameter to tool invocations
2. **Test Gradually**: Start with non-critical operations
3. **Monitor Context Usage**: Track context window improvements
4. **Adjust Field Sets**: Refine field selections based on actual needs

## Summary

The field filtering feature provides:
- **90% context window reduction** for typical use cases
- **Consistent API** across all supported tools
- **Zero breaking changes** to existing functionality
- **Simple syntax** with powerful nested field support
- **Production-ready** with comprehensive test coverage

This feature empowers LLM applications to work more efficiently with incident.io data while maintaining full backward compatibility.
