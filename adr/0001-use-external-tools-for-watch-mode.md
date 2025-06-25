# Use External Tools for Watch Mode

## Status

Accepted

## Context

Users requested a watch mode feature that would automatically regenerate code when markdown files change. This would require:
- File system monitoring implementation
- Debouncing logic to avoid excessive regeneration
- Cross-platform compatibility concerns
- Additional dependencies (like fsnotify)

## Decision

We will not implement a built-in watch mode. Users should use existing, specialized tools like:
- `air` - Live reload for Go apps
- `watchexec` - Execute commands when files change
- `entr` - Run arbitrary commands when files change
- Shell scripts with `inotify-tools`

## Consequences

### Positive
- Follows Unix philosophy of doing one thing well
- Avoids reinventing existing robust solutions
- Reduces codebase complexity and maintenance burden
- Users can choose their preferred watch tool
- Better integration with existing development workflows

### Negative
- Users need to install/configure external tools
- Slightly less convenient for newcomers

### Example Usage

Users can easily implement watch mode with existing tools:

```bash
# Using air
air --build.cmd "adder generate --input docs/commands --output generated"

# Using watchexec
watchexec -e md -- adder generate --input docs/commands --output generated

# Using entr
find docs/commands -name "*.md" | entr adder generate --input docs/commands --output generated
```