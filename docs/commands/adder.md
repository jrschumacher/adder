---
title: A documentation-driven CLI generator
command:
  name: adder
  persistent_flags:
    - name: verbose
      shorthand: v
      description: Enable verbose output for debugging and CI
      type: bool
      default: false
    - name: quiet
      shorthand: q
      description: Suppress all output except errors
      type: bool
      default: false
---

# Adder

A documentation-driven CLI generator that generates type-safe CLI commands from markdown documentation.

Adder processes markdown files with YAML frontmatter to create:
- Type-safe command interfaces
- Request/response structures  
- Handler interfaces
- Argument and flag validation

## Usage

```bash
adder [command] [flags]
```

## Available Commands

- `generate` - Generate CLI commands from markdown documentation
- `init` - Initialize a new adder project
- `version` - Show version information
- `help` - Help about any command

## Examples

```bash
# Generate commands from documentation
adder generate

# Initialize a new project
adder init

# Show version information
adder version
```

## Getting Started

1. Create markdown files with YAML frontmatter defining your commands
2. Run `adder generate` to create Go code
3. Implement the generated handler interfaces
4. Build and run your CLI application