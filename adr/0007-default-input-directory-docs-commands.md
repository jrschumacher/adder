# Default Input Directory: docs/commands

## Status

Accepted

## Context

The original default input directory was `docs/man`, which was inherited from another project and doesn't make sense as a general default. Users need a more intuitive default that:
- Clearly indicates the purpose of the directory
- Follows common conventions
- Allows for other documentation in `/docs`

## Decision

Change the default input directory from `docs/man` to `docs/commands` because:
- **Self-documenting**: Immediately clear what the directory contains
- **Descriptive**: More explicit than abbreviated alternatives like `/cmd`
- **Scalable**: Leaves room for other docs like `/docs/api`, `/docs/guides`
- **Follows patterns**: Similar to how cobra uses `/cmd` for Go code

## Consequences

### Positive
- More intuitive for new users
- Clear separation from other documentation types
- Follows established project conventions (adder itself uses docs/commands)
- Self-explanatory directory structure

### Negative
- Breaking change for existing users with `docs/man` directories
- Need to update all documentation and examples

### Migration

Existing users can:
1. Move their files: `mv docs/man docs/commands`
2. Update their `.adder.yaml` config file
3. Continue using `--input docs/man` flag if they prefer

### Alternatives Considered

- `docs/cmd` - Too abbreviated, could be confused with Go code directory
- `docs/cli` - Less specific than "commands"
- `docs/man` - Original choice, but not intuitive for most users