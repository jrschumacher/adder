# Adder

🐍 **Documentation-driven CLI generator for Cobra based Go applications**

Adder generates type-safe CLI commands from markdown documentation, providing a clean separation between command structure and business logic.

[![Test](https://github.com/jrschumacher/adder/actions/workflows/test.yml/badge.svg)](https://github.com/jrschumacher/adder/actions/workflows/test.yml)
[![Lint](https://github.com/jrschumacher/adder/actions/workflows/lint.yml/badge.svg)](https://github.com/jrschumacher/adder/actions/workflows/lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jrschumacher/adder.svg)](https://pkg.go.dev/github.com/jrschumacher/adder)
[![Go Report Card](https://goreportcard.com/badge/github.com/jrschumacher/adder)](https://goreportcard.com/report/github.com/jrschumacher/adder)
[![Release](https://img.shields.io/github/v/release/jrschumacher/adder)](https://github.com/jrschumacher/adder/releases)

## ✨ Features

- **📝 Documentation-First** - Commands defined in readable markdown
- **🔒 Type-Safe** - Compile-time validation and structured requests  
- **🏗️ Clean Architecture** - Separation of CLI structure and business logic
- **⚡ Performance** - No runtime parsing overhead
- **🎯 Handler Interfaces** - Easy testing and dependency injection
- **📁 Organized Output** - Preserves directory structure to avoid naming conflicts
- **🧪 Comprehensive Testing** - Unit, integration, golden file, and example tests
- **🚀 Production Ready** - Full CI/CD pipeline with automated releases
- **🔧 Self-Dogfooding** - Adder generates its own CLI commands

## 🚀 Quick Start

### 1. Install

**Via Go Install:**
```bash
go install github.com/jrschumacher/adder/cmd@latest
```

**Via GitHub Releases:**
```bash
# Download binary for your platform from:
# https://github.com/jrschumacher/adder/releases
```


### 2. Define Command

Create `docs/man/hello.md`:

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

### 3. Generate Code

```bash
adder generate -i docs/man -o generated -p generated
```

### 4. Implement Handler

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

### 5. Wire It Up

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
}

type HelloRequest struct {
    Arguments HelloRequestArguments `json:"arguments"`
    Flags     HelloRequestFlags     `json:"flags"`
}

// Handler receives full command access
type HelloHandler interface {
    HandleHello(cmd *cobra.Command, req *HelloRequest) error
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

# Build for all platforms  
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

Adder uses automated releases with [release-please](https://github.com/googleapis/release-please):

1. **Merge to main** triggers release PR creation
2. **Merge release PR** creates GitHub release
3. **Automated builds** create multi-platform binaries

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