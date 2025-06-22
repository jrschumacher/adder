---
title: Generate CLI commands from markdown documentation
command:
  name: generate
  flags:
    - name: input
      shorthand: i
      description: Input directory containing markdown files
      default: docs/man
      type: string
    - name: output
      shorthand: o
      description: Output directory for generated files
      default: generated
      type: string
    - name: package
      shorthand: p
      description: Go package name for generated files
      default: generated
      type: string
    - name: suffix
      description: File suffix for generated files
      default: _generated.go
      type: string
---

# Generate CLI Commands

Generate type-safe CLI commands from markdown files with YAML frontmatter.

The generator reads markdown files from the input directory and creates
Go code with command definitions, request structures, and handler interfaces.

## Usage

```bash
adder generate [flags]
```

## Examples

```bash
# Generate from docs/man to generated/ package
adder generate

# Custom input and output directories
adder generate -i documentation -o src/cli -p commands
```

## Output

The generator preserves directory structure from input to output and creates:

- Type-safe request structures
- Handler interfaces
- Command constructors
- Automatic validation