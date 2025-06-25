# Use Git for Backup and Rollback

## Status

Accepted

## Context

Users requested backup and rollback functionality for generated files. Proposed approaches included:
- Creating `.bak` files before generation
- Timestamped backup directories
- Custom backup/restore commands

## Decision

We will not implement built-in backup/rollback functionality. Users should use Git or their preferred version control system for:
- Tracking changes to generated files
- Rolling back to previous versions
- Reviewing changes before committing

## Consequences

### Positive
- Leverages existing version control workflows
- No duplicate functionality with Git
- Encourages best practices (version control for generated files)
- Reduces complexity and potential for backup-related bugs
- Users maintain full control over their backup strategy

### Negative
- Users must remember to commit before regenerating
- No safety net for users not using version control

### Recommended Workflow

```bash
# Before regenerating
git add -A && git commit -m "Before regenerating commands"

# Generate
adder generate --input docs/commands --output generated

# Review changes
git diff

# Rollback if needed
git checkout -- generated/

# Or commit if satisfied
git add -A && git commit -m "Regenerate commands"
```