# Adder Documentation

Welcome to the Adder documentation! Adder is a documentation-driven CLI generator that creates type-safe command interfaces from markdown files.

## Quick Start

1. **Install Adder**
   ```bash
   go install github.com/opentdf/adder/cmd@latest
   ```

2. **Create Command Documentation**
   ```yaml
   ---
   title: My Command
   command:
     name: mycommand [arg]
     arguments:
       - name: arg
         description: An argument
         required: true
         type: string
     flags:
       - name: flag
         description: A flag
         type: bool
   ---
   
   # My Command
   
   This is my command description.
   ```

3. **Generate Code**
   ```bash
   adder generate -i docs -o generated -p myapp
   ```

4. **Implement Handler**
   ```go
   func (h *Handler) HandleMycommand(cmd *cobra.Command, req *MycommandRequest) error {
       fmt.Printf("Hello %s!\n", req.Arguments.Arg)
       return nil
   }
   ```

## Core Concepts

### Documentation-First Approach

Adder follows a documentation-first philosophy:
- Commands are defined in markdown files with YAML frontmatter
- Documentation serves as the single source of truth
- Generated code matches the documented interface exactly

### Type Safety

All generated code is compile-time type-safe:
- Request structures with proper types
- Automatic validation for enums and required fields
- Clear separation between arguments and flags

### Handler Interfaces

Business logic is separated through clean interfaces:
- Handlers receive `*cobra.Command` for full access
- Request structures provide typed access to parameters
- Easy to test and mock

## Architecture

```
docs/man/           # Command documentation (input)
├── auth/
│   └── login.md
├── policy/
│   └── create.md
└── ...

generated/          # Generated code (output)  
├── auth/
│   └── login_generated.go
├── policy/
│   └── create_generated.go
└── ...

handlers/           # Your implementations
├── auth.go
├── policy.go
└── ...
```

## Command Reference

- [`adder generate`](commands/generate.md) - Generate CLI commands from documentation
- [`adder version`](commands/version.md) - Show version information
- [`adder help`](commands/help.md) - Show help information

## Examples

- [Hello World Example](../example/) - Complete working example
- [Complex CLI](examples/complex.md) - Advanced patterns and techniques
- [Testing](examples/testing.md) - How to test generated commands

## Advanced Topics

- [Custom Templates](advanced/templates.md) - Customizing code generation
- [Validation](advanced/validation.md) - Advanced validation patterns  
- [Error Handling](advanced/errors.md) - Best practices for error handling
- [Integration](advanced/integration.md) - Integrating with existing CLIs