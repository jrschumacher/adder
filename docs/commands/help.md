# adder help

Get help information for the adder CLI tool and its commands.

## Usage

```bash
# Show general help
adder help

# Show help for specific command  
adder help generate
adder help version

# Alternative syntax
adder generate --help
adder version --help
```

## Available Commands

- `generate` - Generate CLI commands from markdown documentation
- `version` - Print version information
- `help` - Show help information

## Getting Started

1. **Create command documentation** in markdown files with YAML frontmatter
2. **Generate code** using `adder generate`
3. **Implement handlers** for your business logic
4. **Build and run** your CLI application

For a complete example, see the [Hello World Example](../../example/).