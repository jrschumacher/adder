---
title: Say hello to someone

command:
  name: hello [name]
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

## Examples

```bash
# Simple greeting
otdfctl hello Alice

# Capitalized greeting with big ASCII art
otdfctl hello Alice --capitalize --ascii-art=big

# Repeat the greeting multiple times with banner style
otdfctl hello Bob --ascii-art=banner --repeat=2
```