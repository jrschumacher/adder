---
title: Say hello to someone

command:
  name: greet [name]
  short: Greet someone with a personalized message
  long: |
    A comprehensive greeting command that demonstrates various adder features
    including arguments, flags, enums, and validation.
  example: |
    # Simple greeting
    hello greet Alice
    
    # Fancy greeting with options
    hello greet Bob --capitalize --ascii-art=banner --repeat=3
    
    # Quiet greeting
    hello greet Charlie --quiet --format=json
  arguments:
    - name: name
      description: Name of the person to greet
      required: true
      type: string
  flags:
    - name: capitalize
      description: Capitalize the greeting
      default: false
      type: bool
    - name: ascii-art
      shorthand: a
      description: ASCII art style for the greeting
      default: small
      type: string
      enum:
        - small
        - big
        - banner
    - name: repeat
      shorthand: r
      description: Number of times to repeat the greeting
      default: 1
      type: int
    - name: format
      shorthand: f
      description: Output format for the greeting
      default: text
      type: string
      enum:
        - text
        - json
        - yaml
    - name: quiet
      shorthand: q
      description: Suppress extra output
      default: false
      type: bool
    - name: prefix
      description: Prefix to add before the greeting
      type: string
      default: "Hello"
    - name: languages
      description: Additional languages to greet in
      type: stringArray
---

# Say hello to someone

Greet someone with a friendly hello message.

This command demonstrates the adder package's ability to generate type-safe CLI commands from markdown documentation.

## Arguments

- `name` - Name of the person to greet (required)

## Flags

- `--capitalize` - Capitalize the greeting
- `--ascii-art, -a` - ASCII art style for the greeting (small|big|banner)
- `--repeat, -r` - Number of times to repeat the greeting
- `--format, -f` - Output format (text|json|yaml)
- `--quiet, -q` - Suppress extra output
- `--prefix` - Prefix to add before the greeting
- `--languages` - Additional languages to greet in (string array)

## Examples

```bash
# Simple greeting
hello greet Alice

# Capitalized greeting with big ASCII art
hello greet Alice --capitalize --ascii-art=big

# Repeat the greeting multiple times with banner style
hello greet Bob --ascii-art=banner --repeat=2

# JSON output format
hello greet Charlie --format=json

# Multiple languages
hello greet Diana --languages=spanish,french,german
```

## Features Demonstrated

- **Required Arguments**: The `name` argument is required
- **Enum Validation**: `ascii-art` and `format` flags have predefined valid values
- **Type Variety**: Demonstrates string, bool, int, and stringArray types
- **Default Values**: Most flags have sensible defaults
- **Shorthand Flags**: Several flags have single-character shortcuts