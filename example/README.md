# Adder Example - Hello World

This directory contains a complete example of using the adder package to generate a CLI application.

## Structure

```
example/
├── main.go                  # CLI application entry point
├── go.mod                   # Go module with adder dependency
├── docs/man/hello.md        # Command documentation (source)
├── generated/               # Generated code (output)
└── README.md                # This file
```

## How It Works

1. **Documentation First**: The `hello` command is defined in `docs/man/hello.md` with YAML frontmatter
2. **Code Generation**: Run the adder generator to create type-safe command interfaces
3. **Implementation**: Implement the handler interface in `main.go`
4. **Integration**: Wire the generated command to your CLI

## Usage

### Generate Code

From the example directory:

```bash
# Generate the hello command
go run ../cmd/main.go generate -i docs/man -o generated -p generated

# Or using the adder CLI (once built)
adder generate -i docs/man -o generated -p generated
```

### Build and Run

```bash
# Build the example
go build -o hello-example

# Run it
./hello-example hello Alice --capitalize --ascii-art=big
./hello-example hello Bob --ascii-art=banner --repeat=2
```

## Key Benefits Demonstrated

- **Type Safety**: Arguments and flags are validated at compile time
- **Separation of Concerns**: CLI structure vs business logic
- **Documentation-Driven**: Commands defined in readable markdown
- **Code Generation**: No runtime parsing overhead
- **Clean Architecture**: Handler interfaces promote testability

## Example Output

```
$ ./hello-example hello "Adder" --capitalize --ascii-art=big
╔═══════════════╗
║  HELLO, ADDER!  ║
╚═══════════════╝
```