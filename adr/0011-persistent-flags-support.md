# ADR-0011: Persistent Flags Support

## Status

Accepted

## Context

CLI applications often need global flags that are available to all subcommands. Common examples include `--verbose`, `--quiet`, `--config`, `--debug`, etc. These flags should be accessible from any subcommand without requiring repetitive definitions.

In Cobra, this is achieved through `cmd.PersistentFlags()` which registers flags that are inherited by all subcommands. However, adder only supported local flags via `cmd.Flags()`, limiting users to command-specific flags only.

### Current Limitations

```yaml
# Before: Only local flags supported
command:
  name: myapp
  flags:
    - name: verbose  # Only available to this specific command
      type: bool
```

### User Needs

- Global `--verbose`/`--quiet` flags for CI/CD pipelines
- Shared `--config` flag across all subcommands  
- `--debug` flag available everywhere for troubleshooting
- Type-safe access to persistent flags in handlers

## Decision

Add explicit `persistent_flags:` section to YAML frontmatter, separate from regular `flags:`.

### Design Choice: Separate Section vs Attribute

**Option A: Separate `persistent_flags:` section** ✅ **CHOSEN**
```yaml
command:
  name: myapp
  persistent_flags:
    - name: verbose
      type: bool
  flags:
    - name: input
      type: string
```

**Option B: `persistent: true` attribute** ❌ **REJECTED**
```yaml
command:
  name: myapp
  flags:
    - name: verbose
      type: bool
      persistent: true
    - name: input
      type: string
```

**Rationale for Choice A:**
- **Explicit intent** - clear separation between local and persistent flags
- **Easier templating** - separate iteration loops in templates
- **Better organization** - mirrors Cobra's conceptual distinction
- **Consistent patterns** - follows existing `arguments:` and `flags:` structure

## Implementation

### YAML Structure
```yaml
---
title: Root CLI command
command:
  name: mycli
  persistent_flags:
    - name: verbose
      shorthand: v
      description: Enable verbose output for debugging and CI
      type: bool
      default: false
    - name: config
      shorthand: c
      description: Path to configuration file
      type: string
  flags:
    - name: input
      description: Local flag only for this command
      type: string
---
```

### Generated Code Structure

**Request Types:**
```go
// Separate persistent flags struct
type MycliPersistentFlags struct {
    Verbose bool   `json:"verbose"`
    Config  string `json:"config"`
}

// Main request includes both local and persistent
type MycliRequest struct {
    PersistentFlags MycliPersistentFlags `json:"persistent_flags"`
    Flags          MycliFlags           `json:"flags"`
}
```

**Flag Registration:**
```go
// Register persistent flags first
cmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output for debugging and CI")
cmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file")

// Then register local flags
cmd.Flags().String("input", "", "Local flag only for this command")
```

**Handler Access:**
```go
func myHandler(cmd *cobra.Command, req *MycliRequest) error {
    if req.PersistentFlags.Verbose {
        // Handle verbose mode
    }
    
    config := req.PersistentFlags.Config
    input := req.Flags.Input
    
    return nil
}
```

## Consequences

### Positive

- **Global flag inheritance** - persistent flags available to all subcommands
- **Type safety maintained** - strongly typed access to persistent flags
- **Clear separation** - explicit distinction between local and persistent flags
- **CI/CD friendly** - common `--verbose`/`--quiet` patterns supported
- **Backward compatible** - existing `flags:` sections unchanged
- **Template consistency** - follows existing patterns for code generation

### Negative

- **Template complexity** - additional template logic for persistent flags
- **Request struct growth** - larger request types with separate flag sections
- **Learning curve** - users need to understand persistent vs local flag concepts
- **Documentation burden** - need to explain when to use each type

### Neutral

- **Breaking change avoided** - additive feature, existing code unaffected
- **Generated code size** - slightly larger but still readable
- **Validation complexity** - persistent flags follow same validation patterns

## Usage Patterns

### Root Command Global Flags
```yaml
# Root command defines global flags
command:
  name: myapp
  persistent_flags:
    - name: verbose
      type: bool
    - name: config
      type: string
```

### Subcommand-Specific Persistent Flags
```yaml
# Subcommand can also define persistent flags for its own children
command:
  name: database
  persistent_flags:
    - name: connection-string
      type: string
  flags:
    - name: timeout
      type: int
```

### Usage: `myapp --verbose --config=dev.yaml database --connection-string=postgres://... migrate --timeout=30`

## Testing Strategy

- **Golden file tests** - ensure correct code generation
- **Integration tests** - verify flag inheritance behavior
- **Type safety tests** - confirm request struct generation
- **CLI tests** - validate actual Cobra flag behavior
- **Self-dogfooding** - adder uses its own persistent flags (`--verbose`, `--quiet`)

## Migration Path

**Existing users:** No migration required - additive feature
**New users:** Can immediately use `persistent_flags:` for global flags
**Best practices:** Use persistent flags for cross-cutting concerns (logging, config, debug)

## Examples

**Before:**
```bash
# No way to have global verbose flag
myapp subcommand --input=file.txt  # No global flags available
```

**After:**
```bash
# Global flags available everywhere
myapp --verbose subcommand --input=file.txt
myapp --quiet --config=prod.yaml subcommand --input=file.txt
```

This implementation provides the foundation for building more sophisticated CLI applications with proper global flag support while maintaining adder's type safety and code generation benefits.