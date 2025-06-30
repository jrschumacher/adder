---
title: Debug greeting functionality (hidden)

command:
  name: debug
  short: Debug greeting functionality
  long: |
    Internal debugging command for the hello functionality.
    This command is hidden from help output but can be used
    for troubleshooting and development.
  hidden: true
  flags:
    - name: trace
      description: Enable detailed tracing
      default: false
      type: bool
    - name: dump-config
      description: Dump current configuration
      default: false
      type: bool
    - name: test-enum
      description: Test enum validation
      type: string
      enum:
        - debug
        - info
        - warn
        - error
      default: info
---

# Debug Command (Hidden)

This is a hidden debugging command that demonstrates:

- Hidden commands (not shown in help)
- Advanced flag configurations
- Enum validation with defaults
- Development/debugging patterns

## Usage

```bash
# This won't show in help output
hello help

# But can be called directly
hello debug --trace --dump-config

# Test enum validation
hello debug --test-enum=debug
```

## Features Demonstrated

- **Hidden Command**: Uses `hidden: true` to hide from help
- **Development Tools**: Tracing and config dumping
- **Enum Validation**: Demonstrates enum with default values