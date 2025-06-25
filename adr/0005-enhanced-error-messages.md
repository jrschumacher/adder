# Enhanced Error Messages with File Context

## Status

Accepted

## Context

Original error messages lacked context about where errors occurred:
- No file paths in error messages
- No indication of which field caused conflicts
- Difficult to debug in multi-file projects

## Decision

Enhance all validation error messages to include:
- Full file path where error occurred
- Specific field names and types in conflicts
- Original field names before transformation
- Clear indication of conflict source (argument vs flag)

## Consequences

### Positive
- Faster debugging for users
- Clear actionable error messages
- Better experience for large projects
- Easier to fix validation issues

### Negative
- Slightly more complex error handling code
- Longer error messages

### Examples

Before:
```
duplicate field name after conversion: Output (from flag output)
```

After:
```
file docs/commands/example.md: duplicate field name 'Output' - flag 'output' conflicts with argument 'output'
```

The enhanced format provides:
- Exact file location
- Transformed field name that conflicts
- Both original names
- Clear indication of conflict type