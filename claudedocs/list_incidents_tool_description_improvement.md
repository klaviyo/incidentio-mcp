# List Incidents Tool Description Improvement

## Date: 2025-10-13

## Overview
Enhanced the `list_incidents` tool description and parameter documentation to provide clear, comprehensive guidance on proper tool usage, particularly for severity filtering.

## Problem Statement
The original tool description was too brief and didn't provide adequate guidance on:
- The proper workflow for filtering by severity
- The requirement to use severity IDs rather than names
- The need to call `list_severities` first
- How filters combine (OR logic)
- Pagination behavior

This led to potential confusion when using the tool, especially for AI agents that need clear instructions.

## Changes Made

### 1. Enhanced Tool Description

**File:** `/internal/tools/incidents.go` - `ListIncidentsTool.Description()`

**Before:**
```go
return "List incidents from incident.io with optional filters"
```

**After:**
```go
return `List incidents from incident.io with optional filtering by status and severity.

USAGE WORKFLOW:
1. To filter by severity, first call 'list_severities' to get available severity IDs
2. Use the severity ID (not the name) in the 'severity' parameter
3. Multiple severity IDs can be provided to match any of them (OR logic)
4. Status filters can be combined with severity filters

PARAMETERS:
- page_size: Number of results (default 25, max 250). Set to 0 or omit for auto-pagination.
- status: Array of status values (triage, active, resolved, closed)
- severity: Array of severity IDs (e.g., ["sev_1", "01HXYZ..."]) - Use list_severities first to get valid IDs

EXAMPLES:
- List all active incidents: {"status": ["active"]}
- List critical incidents: First call list_severities, then use severity ID like {"severity": ["sev_1"]}
- List active high-severity incidents: {"status": ["active"], "severity": ["sev_1", "sev_2"]}

IMPORTANT: Severity parameter requires severity IDs, not severity names. Always call list_severities first to discover available severity IDs.`
```

### 2. Enhanced Parameter Descriptions

**File:** `/internal/tools/incidents.go` - `ListIncidentsTool.InputSchema()`

#### page_size Parameter
**Before:**
```go
"description": "Number of results per page (max 250)"
```

**After:**
```go
"description": "Number of results per page (max 250, default 25). Set to 0 or omit for automatic pagination through all results."
```

#### status Parameter
**Before:**
```go
"description": "Filter by incident status (e.g., triage, active, resolved, closed)"
```

**After:**
```go
"description": "Filter by incident status values. Valid values: triage, active, resolved, closed. Multiple values will match any of them (OR logic). Example: [\"active\", \"triage\"]"
```

#### severity Parameter
**Before:**
```go
"description": "Filter by severity"
```

**After:**
```go
"description": "Filter by severity IDs (NOT severity names). IMPORTANT: Call 'list_severities' tool first to discover available severity IDs. Example: [\"sev_1\", \"sev_2\"] or [\"01HXYZ...\"]. Multiple IDs will match any of them (OR logic)."
```

## Benefits

### For AI Agents
1. **Clear Workflow**: Step-by-step instructions on how to use severity filtering
2. **Prevents Errors**: Explicitly states to use IDs, not names
3. **Discovery Pattern**: Directs to call `list_severities` first
4. **Examples**: Concrete examples for common use cases

### For Developers
1. **Self-Documenting**: Tool description serves as inline documentation
2. **Less Confusion**: Clear expectations about parameter values
3. **Better Testing**: Examples provide test case ideas
4. **Easier Debugging**: Clear descriptions help identify issues

### For Users
1. **Better Error Messages**: LLMs can provide clearer guidance based on detailed descriptions
2. **Fewer Failed Calls**: Proper guidance reduces trial-and-error
3. **Consistent Experience**: All tool descriptions follow same detailed pattern

## Technical Details

### Description Format
- Multi-line string using backticks for readability
- Structured sections: USAGE WORKFLOW, PARAMETERS, EXAMPLES, IMPORTANT
- Uppercase section headers for easy scanning
- Numbered steps for workflow clarity

### Parameter Descriptions
- Explicit value types and formats
- Range and default information
- Practical examples with actual syntax
- Warning about common mistakes (IDs vs names)
- OR logic explicitly stated

### Backward Compatibility
✅ **Fully backward compatible**
- No changes to function signatures
- No changes to parameter types
- No changes to execution logic
- Only description strings enhanced
- Existing code continues to work unchanged

## Verification

### Compilation
```bash
go build ./internal/tools/...
# ✓ Code compiles successfully
```

### Tests
```bash
go test ./internal/tools/...
# ✓ All tests pass
```

### String Length
The enhanced descriptions are significantly longer but remain reasonable:
- Main description: ~800 characters
- Parameter descriptions: 100-200 characters each
- Well within MCP protocol limits

## Usage Examples

### Example 1: Filter by Severity (Correct Workflow)
```
Step 1: AI agent reads tool description
Step 2: Sees "Call list_severities first"
Step 3: Calls list_severities tool
Step 4: Gets severity IDs (e.g., sev_1, sev_2, sev_3)
Step 5: Calls list_incidents with {"severity": ["sev_1"]}
Result: ✓ Only critical incidents returned
```

### Example 2: Combined Filters
```
Tool call: list_incidents
Parameters: {
  "status": ["active", "triage"],
  "severity": ["sev_1", "sev_2"]
}
Result: ✓ Returns incidents that are (active OR triage) AND (sev_1 OR sev_2)
```

### Example 3: Auto-Pagination
```
Tool call: list_incidents
Parameters: {
  "status": ["active"]
  // page_size omitted
}
Result: ✓ All active incidents returned via automatic pagination
```

## Best Practices Applied

### 1. User-Centric Documentation
- Written from the perspective of someone using the tool
- Anticipates common questions and mistakes
- Provides practical, actionable guidance

### 2. Progressive Disclosure
- Summary at top (one-line overview)
- Detailed workflow for those who need it
- Examples for quick reference

### 3. Defensive Design
- Explicitly warns about common mistakes
- Points to related tools when needed
- Clarifies OR vs AND logic

### 4. Concrete Examples
- Real JSON syntax
- Actual parameter values
- Multiple use case scenarios

## Impact on Related Tools

### Tools Following Similar Pattern
This improvement can serve as a template for other tools:

1. **create_incident** - Already has good suggestions, could be enhanced further
2. **update_incident** - Could benefit from similar detailed descriptions
3. **list_alerts** - Similar filtering logic, should follow same pattern
4. **list_workflows** - Could use enhanced parameter descriptions

### Recommendation
Apply similar detailed description pattern to all filtering-related tools for consistency.

## Maintenance

### When to Update
- When adding new filter parameters
- When filter logic changes (e.g., AND vs OR)
- When related tools are renamed
- When common usage patterns change

### Related Files
- `/internal/tools/incidents.go` - Tool implementation
- `/claudedocs/TOOLS_REFERENCE.md` - External documentation
- `/internal/incidentio/incidents.go` - API client (uses same parameters)

## Summary

The tool description has been significantly enhanced to provide:
- ✅ Clear step-by-step usage workflow
- ✅ Explicit parameter requirements (IDs vs names)
- ✅ Guidance to use `list_severities` first
- ✅ Concrete examples for common scenarios
- ✅ Clear documentation of OR logic
- ✅ Pagination behavior explanation

This improvement ensures that both AI agents and human developers have clear, comprehensive guidance on proper tool usage, particularly for the severity filtering feature.

---

**Changes:**
- 1 file modified: `/internal/tools/incidents.go`
- 0 breaking changes
- Backward compatible
- Code compiles and tests pass
