---
title: Initialize adder configuration
command:
  name: init
  flags:
    - name: binary-name
      shorthand: b
      description: Name of the binary/CLI (required)
      type: string
    - name: force
      shorthand: f
      description: Overwrite existing configuration file
      default: false
      type: bool
---

# Initialize Adder Configuration

Create a `.adder.yaml` configuration file through an interactive setup process.

## Usage

```bash
adder init [flags]
```

## Description

The init command guides you through creating an adder configuration file for your project. It will:

1. Require binary name (CLI name) - either via flag or prompt
2. Ask for input directory (where your markdown files are)
3. Ask for output directory (where generated code should go)
4. Ask for Go package name
5. Ask for file suffix for generated files
6. Set index format to 'directory' by default
7. Create `.adder.yaml` with your choices

## Examples

```bash
# Interactive setup with binary name flag
adder init --binary-name myapp

# Interactive setup (will prompt for binary name)
adder init

# Overwrite existing config
adder init --binary-name myapp --force
```

## Configuration File

The generated `.adder.yaml` file contains:

```yaml
binary_name: myapp
input: docs/commands
output: generated
package: generated
suffix: _generated.go
package_strategy: directory
index_format: directory
```

Once created, you can run `adder generate` without any flags and it will use the configuration file.