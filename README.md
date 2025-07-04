# Adder

🐍 **Documentation-driven CLI generator for Cobra based Go applications**

Adder generates type-safe CLI commands from markdown documentation, providing a clean separation between command structure and business logic.

[![Test](https://github.com/jrschumacher/adder/actions/workflows/test.yml/badge.svg)](https://github.com/jrschumacher/adder/actions/workflows/test.yml)
[![Lint](https://github.com/jrschumacher/adder/actions/workflows/lint.yml/badge.svg)](https://github.com/jrschumacher/adder/actions/workflows/lint.yml)
[![Release](https://github.com/jrschumacher/adder/actions/workflows/release.yml/badge.svg)](https://github.com/jrschumacher/adder/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jrschumacher/adder.svg)](https://pkg.go.dev/github.com/jrschumacher/adder)
[![Go Report Card](https://goreportcard.com/badge/github.com/jrschumacher/adder)](https://github.com/jrschumacher/adder/report/github.com/jrschumacher/adder)
[![Release](https://img.shields.io/github/v/release/jrschumacher/adder)](https://github.com/jrschumacher/adder/releases)

## ✨ Features

- **📝 Documentation-First** - Commands defined in readable markdown
- **🔒 Type-Safe** - Compile-time validation and structured requests  
- **🏗️ Clean Architecture** - Separation of CLI structure and business logic
- **⚡ Performance** - No runtime parsing overhead
- **🎯 Handler Interfaces** - Easy testing and dependency injection
- **📁 Organized Output** - Preserves directory structure to avoid naming conflicts
- **✅ Enhanced Validation** - Markdown linter with strict type checking and enum validation
- **🔌 Request Interface** - All generated requests implement `adder.Request` for consistency
- **🛡️ Centralized Enum Validation** - Clean `adder.ValidateEnum()` function for runtime validation
- **🧪 Comprehensive Testing** - Unit, integration, golden file, and example tests
- **🚀 Production Ready** - Full CI/CD pipeline with automated releases
- **🔧 Self-Dogfooding** - Adder generates its own CLI commands

## 🚀 Quick Start

### 1. Install

**Via Go Install:**
```bash
go install github.com/jrschumacher/adder/cmd/adder@latest
```

**Via GitHub Releases:**
```bash
# Download binary for your platform from:
# https://github.com/jrschumacher/adder/releases
```


### 2. Configure Project

Create `.adder.yaml`:

```yaml
binary_name: myapp              # Required: name of your CLI binary
input: docs/commands            # Default: docs/commands  
output: generated               # Default: generated
package: generated              # Default: generated
suffix: _generated.go           # Default: _generated.go
package_strategy: directory     # Default: directory
index_format: directory         # Default: directory
```

### 3. Define Command

Create `docs/commands/hello.md`:

```yaml
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
      type: bool
---

# Say hello to someone

Greet someone with a friendly hello message.
```

### 4. Generate Code

```bash
adder generate
# Or with flags: adder generate -i docs/commands -o generated -p generated
```

### 5. Implement Handler

```go
func (h *HelloHandler) HandleHello(cmd *cobra.Command, req *generated.HelloRequest) error {
    greeting := fmt.Sprintf("Hello, %s!", req.Arguments.Name)
    if req.Flags.Capitalize {
        greeting = strings.ToUpper(greeting)
    }
    fmt.Println(greeting)
    return nil
}
```

### 6. Wire It Up

```go
handler := &HelloHandler{}
helloCmd := generated.NewHelloCommand(handler)
rootCmd.AddCommand(helloCmd)
```

## 🏗️ Generated Structure

Adder creates clean, type-safe structures:

```go
// Separate arguments from flags for clarity
type HelloRequestArguments struct {
    Name string `json:"name" validate:"required"`
}

type HelloRequestFlags struct {
    Capitalize bool `json:"capitalize"`
    Style      string `json:"style" validate:"oneof=normal bold italic"`
}

// All requests implement adder.Request interface
type HelloRequest struct {
    Arguments    HelloRequestArguments `json:"arguments"`
    Flags        HelloRequestFlags     `json:"flags"`
    RawArguments []string              `json:"raw_arguments"`
}

// GetRawArguments implements adder.Request interface
func (r *HelloRequest) GetRawArguments() []string {
    return r.RawArguments
}

// Handler receives full command access
type HelloHandler func(cmd *cobra.Command, req *HelloRequest) error

// Clean enum validation in generated code
if err := adder.ValidateEnum("style", style, []string{"normal", "bold", "italic"}); err != nil {
    return err
}
```

## 📁 Directory Organization

Adder preserves your documentation structure:

```
docs/man/              generated/
├── auth/              ├── auth/
│   └── login.md  →    │   └── login_generated.go
└── policy/            └── policy/
    └── create.md  →       └── create_generated.go
```

This prevents naming conflicts between commands like `auth create` and `policy create`.

## ⚙️ Configuration

Create `.adder.yaml` in your project root:

```yaml
# Required: Name of your CLI binary (used to detect root command)
binary_name: myapp

# Optional: Input/output directories (can be overridden with flags)
input: docs/commands
output: generated
package: generated
suffix: _generated.go

# Optional: Package naming strategy
package_strategy: directory  # single, directory, path

# Optional: Index file format for subcommands
index_format: directory      # directory, index, _index, hugo
```

**Root Command Detection:**
- The parser looks for `{binary_name}.md` in the input directory
- This file becomes your CLI's root command
- Example: `binary_name: myapp` → looks for `myapp.md`

## ✅ Enhanced Validation

Adder acts as a comprehensive markdown linter, catching configuration errors early:

### **Type Consistency Validation**
```bash
# ❌ This will fail validation
flags:
  - name: count
    type: int
    default: "not-a-number"  # Error: must be integer for type 'int'

# ✅ This is correct
flags:
  - name: count
    type: int
    default: 42
```

### **Enum Validation**
```bash
# ❌ This will fail validation
flags:
  - name: level
    type: int                    # Error: enum only supported on string type
    enum: ["debug", "info"]

# ❌ This will also fail
flags:
  - name: level
    type: string
    enum: ["debug", "info"]
    default: "invalid"           # Error: default must be in enum values

# ✅ This is correct
flags:
  - name: level
    type: string
    enum: ["debug", "info", "warn"]
    default: "info"
```

### **Validation Commands**
```bash
# Validate without generating
adder generate --validate

# Example validation errors
❌ Validation failed: flag count: default value must be an integer for type 'int', got 'not-a-number'
❌ Validation failed: flag level: enum validation only supported on string type, got 'int'
❌ Validation failed: flag level: default value 'invalid' must be one of: [debug info warn]
```

## 🎯 Key Benefits

| Feature                 | Benefit                                         |
|-------------------------|-------------------------------------------------|
| **Type Safety**         | Compile-time validation prevents runtime errors |
| **Documentation-First** | Single source of truth in readable format       |
| **Performance**         | No runtime parsing of embedded docs             |
| **Clean Architecture**  | Handler interfaces promote testability          |
| **Organized Output**    | Directory structure prevents naming conflicts   |
| **Command Access**      | Full `*cobra.Command` capabilities available    |

## 🏆 Production Ready

Adder is built with production use in mind:

### ✅ **Quality Assurance**
- **19 Active Linters** - golangci-lint with comprehensive checks
- **4 Test Categories** - Unit, integration, golden file, and example tests  
- **90%+ Test Coverage** - Comprehensive test suite
- **Automated Quality Gates** - CI/CD pipeline enforces standards

### ✅ **Reliability**
- **Multi-Platform Support** - Linux, macOS, Windows (AMD64 + ARM64)
- **Semantic Versioning** - Automated releases with conventional commits
- **Backward Compatibility** - Careful API evolution

### ✅ **Developer Experience**
- **Self-Dogfooding** - Tool generates its own CLI commands
- **Comprehensive Documentation** - Examples, guides, and API reference
- **Local Development Tools** - Makefile with all common tasks
- **IDE Integration** - Works with any Go-compatible editor

## 🧪 Testing

Adder includes comprehensive testing:

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run linting
make lint

# Run all CI checks locally
make ci-test
```

**Test Categories:**
- **Unit Tests** - Core parser and generator logic
- **Integration Tests** - CLI command testing  
- **Golden File Tests** - Generated code validation
- **Example Tests** - Handler testing patterns

## 📚 Documentation

- **[Complete Example](example/)** - Working hello world demo
- **[Testing Guide](docs/testing.md)** - Comprehensive testing patterns
- **[GitHub Workflows](.github/README.md)** - CI/CD pipeline documentation
- **[API Reference](https://pkg.go.dev/github.com/jrschumacher/adder)** - Full documentation

## 🏗️ Development

**Local Development:**
```bash
# Build the CLI
make build

# Build for all platforms (uses GitHub Actions)
make build-all

# Generate example commands
make generate-example

# Self-generate CLI commands (dogfooding)
make generate-self
```


## 🎬 Example Output

```bash
$ hello-example hello "Adder" --capitalize
HELLO, ADDER!

$ adder generate --input docs/man --output generated --package generated
🐍 Generating CLI commands from docs/man to generated...
🔍 Validating documentation...
✅ Code generation completed!
📊 Generated 3 commands with 5 flags and 2 arguments
```

## 🚀 Release Process

Adder uses automated releases with [release-please](https://github.com/googleapis/release-please) and [GoReleaser](https://goreleaser.com/):

1. **Merge to main** triggers release PR creation
2. **Merge release PR** creates GitHub release  
3. **GoReleaser** automatically builds and publishes multi-platform binaries

Use [conventional commits](https://www.conventionalcommits.org/) for automatic versioning:
- `feat:` → minor version bump
- `fix:` → patch version bump  
- `BREAKING CHANGE:` → major version bump

## 🤝 Contributing

We welcome contributions! 

**Getting Started:**
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make ci-test` to validate
5. Submit a pull request

**Code Quality:**
- All tests must pass
- golangci-lint must pass
- Include tests for new features
- Update documentation as needed

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.