# Validate Flag Instead of Dry Run

## Status

Accepted

## Context

Users requested ability to validate their markdown documentation without generating files. The initial suggestion was `--dry-run` flag, following common CLI conventions.

## Decision

Use `--validate` flag instead of `--dry-run` because:
- More descriptive of actual functionality (validation)
- Clearer intent for users
- Aligns with the primary purpose: checking documentation correctness

## Consequences

### Positive
- Clearer communication of feature purpose
- More intuitive for users focused on validation
- Differentiates from typical dry-run behavior (showing what would happen)

### Negative
- Deviates from common `--dry-run` convention
- Might require explanation for users expecting dry-run

### Usage

```bash
# Validate documentation without generating files
adder generate --validate --input docs/commands

# Output shows:
# - Validation errors with file paths and line context
# - Statistics about commands found
# - Clear indication that no files were generated
```