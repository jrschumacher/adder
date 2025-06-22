# Adder

🐍 **Documentation-driven CLI generator for Cobra based Go applications**

Adder generates type-safe CLI commands from markdown documentation, providing a clean separation between command structure and business logic.

[![Go Reference](https://pkg.go.dev/badge/github.com/jrschumacher/adder.svg)](https://pkg.go.dev/github.com/jrschumacher/adder)
[![Go Report Card](https://goreportcard.com/badge/github.com/jrschumacher/adder)](https://goreportcard.com/report/github.com/jrschumacher/adder)

## ✨ Features

- **📝 Documentation-First** - Commands defined in readable markdown
- **🔒 Type-Safe** - Compile-time validation and structured requests  
- **🏗️ Clean Architecture** - Separation of CLI structure and business logic
- **⚡ Performance** - No runtime parsing overhead
- **🎯 Handler Interfaces** - Easy testing and dependency injection
- **📁 Organized Output** - Preserves directory structure to avoid naming conflicts

## 🚀 Quick Start

### 1. Install

```bash
go install github.com/jrschumacher/adder/cmd@latest
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

## 📚 Documentation

- **[Quick Start Guide](docs/)** - Get up and running
- **[Complete Example](example/)** - Working hello world demo
- **[API Reference](https://pkg.go.dev/github.com/jrschumacher/adder)** - Full documentation
- **[Command Reference](docs/commands/)** - CLI tool usage

## 🎬 Example Output

```bash
$ hello-example hello "Adder" --capitalize
HELLO, ADDER!
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.