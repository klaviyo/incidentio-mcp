# Dynamic Schema Generation

## Overview

Implemented runtime reflection-based schema generation that automatically extracts and documents available fields from Go struct types. This eliminates manual maintenance of field lists and ensures documentation always matches the actual type definitions.

## Problem Solved

**Before**: Field descriptions were manually written and could become outdated when types changed:
```go
"description": "Available fields: id, name, summary..." // Could drift from actual type
```

**After**: Field descriptions are generated dynamically from struct definitions:
```go
"description": GetIncidentFieldsDescription() // Always accurate, never drifts
```

## Implementation

### Core Components

#### 1. Schema Generator (`internal/tools/schema_generator.go`)

**`GenerateFieldsDescription(exampleType interface{}) string`**
- Takes any Go struct type
- Uses reflection to extract JSON field tags
- Identifies nested structures automatically
- Returns formatted description string

**`GetIncidentFieldsDescription() string`**
- Convenience function for Incident type
- Returns: Complete field documentation

**`GetAlertFieldsDescription() string`**
- Convenience function for Alert type
- Returns: Complete field documentation

#### 2. Field Extraction Logic

**`extractFieldsFromType(t reflect.Type) ([]string, map[string][]string)`**
- Walks through struct fields using reflection
- Reads `json:` struct tags to get field names
- Handles `omitempty` and other tag modifiers
- Identifies nested struct fields (excluding time.Time)
- Returns top-level fields and nested field mappings

### How It Works

```go
// 1. Define a type with JSON tags
type Incident struct {
    ID       string         `json:"id"`
    Name     string         `json:"name"`
    Severity Severity       `json:"severity"`
    // ...
}

type Severity struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Rank int    `json:"rank"`
}

// 2. Generate description at runtime
desc := GenerateFieldsDescription(Incident{})

// 3. Result includes all fields automatically
// "Available top-level fields: id, name, severity, ...
//  Nested fields (dot notation):
//  - severity: id, name, rank"
```

### Reflection-Based Extraction

The system uses Go's `reflect` package to inspect types:

1. **Get Type**: `reflect.TypeOf(Incident{})`
2. **Iterate Fields**: `t.NumField()` and `t.Field(i)`
3. **Extract JSON Tags**: `field.Tag.Get("json")`
4. **Parse Tag**: Split by comma, extract field name
5. **Check Nested**: If field type is struct, recurse
6. **Build Documentation**: Format as human-readable description

## Usage

### In Tool Definitions

```go
func (t *ListIncidentsTool) InputSchema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "fields": map[string]interface{}{
                "type":        "string",
                "description": GetIncidentFieldsDescription(), // ← Dynamic!
            },
        },
    }
}
```

### Generated Output Example

For `Incident` type, generates:

```
Comma-separated list of fields to include in response to reduce context usage. Supports nested fields with dot notation.

Available top-level fields: id, reference, name, summary, permalink, incident_status, severity, incident_type, mode, visibility, created_at, updated_at, slack_team_id, slack_channel_id, slack_channel_name, incident_role_assignments, custom_field_entries, has_debrief

Nested fields (dot notation):
- incident_status: id, name, description, category, rank, created_at, updated_at
- severity: id, name, description, rank, created_at, updated_at
- incident_type: id, name, description, is_default, private_incidents_only, create_in_triage, created_at, updated_at

Examples: "id,name" or with nested fields: "id,name,severity.name,incident_status.category"
Omit to return all fields.
```

## Benefits

### 1. **Zero Maintenance**
- Field lists update automatically when types change
- No risk of documentation drift
- Single source of truth (the struct definition)

### 2. **Type Safety**
- Compile-time verification that types exist
- Can't document non-existent fields
- Catches field renames/removals automatically

### 3. **Consistency**
- All tools use identical generation logic
- Uniform documentation format across tools
- Predictable structure for LLMs

### 4. **Extensibility**
- Easy to add new types
- Simple to customize formatting
- Can extend to support more metadata

## Testing

Comprehensive test suite in `schema_generator_test.go`:

- ✅ Incident type field extraction
- ✅ Alert type field extraction
- ✅ Custom struct handling
- ✅ Flat struct (no nested fields) handling
- ✅ Nested field identification
- ✅ JSON tag parsing
- ✅ Field organization

Run tests:
```bash
go test -v ./internal/tools -run TestGenerateFieldsDescription
```

## Technical Details

### Handling Edge Cases

**Time Fields**: `time.Time` is treated as a primitive, not a nested struct
```go
func isTimeType(t reflect.Type) bool {
    return t.PkgPath() == "time" && t.Name() == "Time"
}
```

**Pointer Types**: Automatically unwrapped to access underlying type
```go
if t.Kind() == reflect.Ptr {
    t = t.Elem()
}
```

**Optional Fields**: `omitempty` tag modifier is stripped, field still documented
```go
jsonName := strings.Split(jsonTag, ",")[0]  // Takes name, ignores modifiers
```

**Unexported Fields**: Skipped automatically
```go
if !field.IsExported() {
    continue
}
```

### Performance Considerations

- **Initialization Time**: Generation happens once per tool initialization
- **Runtime Cost**: Zero - descriptions are generated at startup, cached
- **Memory**: Minimal - just string storage for descriptions
- **Reflection Overhead**: One-time cost during tool schema creation

## Future Enhancements

### Potential Improvements

1. **Field Type Annotations**
```go
Available fields:
- id (string): Unique identifier
- rank (int): Priority ranking
- created_at (timestamp): Creation time
```

2. **Required vs Optional Markers**
```go
- id (required)
- summary (optional)
```

3. **Field Descriptions from Comments**
```go
type Incident struct {
    // The unique incident identifier
    ID string `json:"id" desc:"Unique incident identifier"`
}
```

4. **Validation Rules**
```go
- page_size (int, max: 250): Results per page
```

5. **Example Values**
```go
- status (string, example: "active"): Current status
```

### Extensibility Example

```go
type FieldMetadata struct {
    Name        string
    Type        string
    Required    bool
    Description string
    Example     interface{}
}

func GenerateEnhancedDescription(t interface{}) string {
    metadata := extractFieldMetadata(reflect.TypeOf(t))
    return formatEnhancedDescription(metadata)
}
```

## Migration Impact

### Zero Breaking Changes

✅ Existing tools continue to work
✅ Generated descriptions match previous format
✅ All tests pass without modification
✅ Build successful with no errors

### Maintenance Improvement

**Before**:
- ❌ Manual field list updates when types change
- ❌ Risk of documentation drift
- ❌ Inconsistent format across tools

**After**:
- ✅ Automatic updates when types change
- ✅ Guaranteed accuracy
- ✅ Consistent format everywhere

## Example: Adding a New Field

**Manual Approach (Old)**:
1. Add field to `Incident` struct
2. Update `ListIncidentsTool` description
3. Update `GetIncidentTool` description
4. Update documentation
5. Hope you didn't miss any spots

**Dynamic Approach (New)**:
1. Add field to `Incident` struct
2. Done! All descriptions update automatically

```go
// Add new field to type
type Incident struct {
    // ... existing fields ...
    Priority string `json:"priority"`  // ← New field added
}

// All tool descriptions automatically include "priority" now!
// No manual updates needed anywhere
```

## Summary

Dynamic schema generation transforms field documentation from a maintenance burden into a zero-cost benefit:

- **Automatic**: Fields extracted via reflection at runtime
- **Accurate**: Always matches actual type definitions
- **Maintainable**: Single source of truth in struct definitions
- **Extensible**: Easy to add new types and enhance metadata
- **Tested**: Comprehensive test coverage ensures correctness

This is the difference between "documentation that works" and "documentation that can't be wrong."
