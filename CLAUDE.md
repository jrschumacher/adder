# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Core Development Commands

Since the Makefile referenced in README.md is not present, use these Go commands directly:

```bash
# Testing
go test -v ./...                              # Run all tests
go test -v -race -cover ./...                 # Tests with race detection and coverage
go test -v ./cmd/                             # Integration tests only
go test -v -run TestGenerator_GoldenFiles     # Golden file tests
go test -v ./example/                         # Example tests

# Building
go build -o adder ./cmd/adder                 # Build the CLI tool
go build -o hello-example ./example           # Build the example

# Code Generation (after building)
./adder generate                              # Uses .adder.yaml config
./adder generate --binary-name myapp          # Override binary name
./adder generate --input docs/man --output generated --package generated

# Linting
golangci-lint run --timeout=10m --disable=forbidigo,gosec,mnd,nilnil --enable=errcheck,govet,staticcheck,unused

# Run a single test
go test -v -run TestName ./path/to/package
```

## Configuration

The project uses `.adder.yaml` for configuration:

```yaml
binary_name: adder          # Required: CLI binary name (detects root command)
input: docs/commands        # Input directory for markdown files
output: cmd/adder/generated # Output directory for generated Go files
package: generated          # Go package name for generated files
suffix: _generated.go       # File suffix for generated files
```

**Root Command Detection:**
- `binary_name` is required and determines root command file
- Parser looks for `{binary_name}.md` in input directory
- Example: `binary_name: adder` → looks for `adder.md`

### Self-Dogfooding Commands

Adder uses itself to generate its own CLI commands:

```bash
# Build adder first
go build -o adder ./cmd/adder

# Generate adder's own commands from docs/commands
./adder generate --input docs/commands --output cmd/adder/generated --package generated
```

## Architecture

### Project Overview

Adder is a documentation-driven CLI generator for Cobra-based Go applications that parses markdown files with YAML frontmatter to generate type-safe command structures.

### Key Components

1. **Parser** (`pkg/parser/`)
   - Parses markdown files with YAML frontmatter
   - Validates command definitions
   - Builds internal command representation

2. **Generator** (`pkg/generator/`)
   - Generates Go code from parsed commands
   - Creates handler interfaces and request types
   - Uses Go templates for code generation

3. **CLI** (`cmd/adder/`)
   - Self-generated using adder (dogfooding)
   - Commands defined in `docs/commands/`
   - Generated code in `cmd/adder/generated/`

4. **Templates** (`pkg/generator/templates/`)
   - Go templates for code generation
   - Separate templates for commands, types, and helpers

### Generated Code Structure

For each markdown command file, adder generates:

1. **Request Types**:
   - Separate `Arguments` and `Flags` structs
   - Combined `Request` struct with validation tags

2. **Handler Interface**:
   - Method name: `Handle<CommandName>`
   - Receives `*cobra.Command` and typed request
   - Returns error

3. **Command Constructor**:
   - `New<CommandName>Command(handler)` function
   - Wires up Cobra command with validation

### Directory Preservation

Input directory structure is preserved in output to avoid naming conflicts:
- `docs/man/auth/login.md` → `generated/auth/login_generated.go`
- `docs/man/policy/create.md` → `generated/policy/create_generated.go`

### Testing Strategy

1. **Unit Tests**: Test parser and generator logic in isolation
2. **Integration Tests**: Test CLI commands end-to-end
3. **Golden File Tests**: Ensure consistent code generation
4. **Example Tests**: Demonstrate handler testing patterns

When updating templates or generator logic, run golden file tests with `-update-golden` flag to update reference files.

### CI/CD Pipeline

- **Test Workflow**: Runs all tests, builds binaries, validates examples
- **Lint Workflow**: Enforces code quality with golangci-lint
- **Release Workflow**: Automated releases with release-please and GoReleaser

Uses conventional commits for automatic versioning:
- `feat:` → minor version
- `fix:` → patch version
- `BREAKING CHANGE:` → major version