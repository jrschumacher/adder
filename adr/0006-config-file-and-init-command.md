# Config File and Init Command

## Status

Accepted

## Context

Users need to repeatedly specify the same flags (input, output, package) when running adder. This creates:
- Verbose command lines
- Potential for inconsistency across team members
- No project-specific defaults

## Decision

Implement:
1. Config file support for project-specific settings (`.adder.yaml` or `.adder.yml`)
2. `adder init` command for interactive setup and config file creation
3. Config file values override defaults but are overridden by CLI flags

Priority order (highest to lowest):
1. Command line flags
2. Config file (`.adder.yaml` or `.adder.yml`, checked in that order)
3. Built-in defaults

## Consequences

### Positive
- Simpler command usage after initial setup
- Project-wide consistency
- Better onboarding experience with guided init
- Follows common CLI tool patterns (like .eslintrc, .prettierrc)

### Negative
- Another file to maintain in projects
- Potential confusion about which settings are active
- Need to handle config file validation

### Config File Structure

```yaml
# .adder.yaml
input: docs/commands
output: internal/cli/generated
package: commands
suffix: _gen.go

# Optional: validation rules
validation:
  strict: true
```

### Usage

```bash
# Interactive setup
adder init

# Uses config file automatically
adder generate

# Override config file with flags
adder generate --input other/docs
```