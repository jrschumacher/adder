---
title: Hello command group

command:
  name: hello
  short: Greeting commands and utilities
  long: |
    The hello command group provides various greeting functionality
    including personalized messages, formatting options, and debugging tools.
    
    This serves as the parent command for hello-related subcommands and
    demonstrates command grouping in adder.
  flags:
    - name: version
      description: Show version information
      default: false
      type: bool
  persistent_flags:
    - name: verbose
      shorthand: v
      description: Enable verbose output for all hello commands
      default: false
      type: bool
    - name: config
      shorthand: c
      description: Configuration file path
      type: string
      default: ~/.hello.yaml
---

# Hello Commands

The hello command group provides greeting functionality with various customization options.

## Available Commands

- `hello greet` - Say hello to someone (main greeting command)
- `hello debug` - Debug greeting functionality (hidden)

## Global Flags

All hello commands support these persistent flags:

- `--verbose, -v` - Enable verbose output
- `--config, -c` - Specify configuration file path

## Examples

```bash
# Show help for hello commands
hello --help

# Use verbose mode for any hello command
hello greet Alice --verbose

# Use custom config file
hello greet Bob --config /path/to/config.yaml
```