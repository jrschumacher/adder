# Testing Guide

This document describes the comprehensive testing strategy for the adder package.

## Test Categories

### 1. Unit Tests

Unit tests validate individual components in isolation:

- **Parser Tests** (`parser_test.go`): Test YAML frontmatter parsing, command validation, and helper functions
- **Generator Tests** (`generator_test.go`): Test code generation, file output, and validation logic
- **Type Tests**: Test helper methods on types like `GetGoType()` and `GetValidationTag()`

Run unit tests:
```bash
go test -v ./...
```

### 2. Integration Tests

Integration tests validate the CLI commands end-to-end:

- **CLI Command Tests** (`cmd/integration_test.go`): Test the generate and version commands with real file I/O
- **Handler Integration**: Test generated commands with their handlers

Run integration tests:
```bash
go test -v ./cmd/
```

### 3. Golden File Tests

Golden file tests ensure generated code remains consistent:

- **Golden Files** (`testdata/golden/*.golden`): Reference output for code generation
- **Golden Test** (`golden_test.go`): Compares generated output against golden files

Run golden file tests:
```bash
go test -v -run TestGenerator_GoldenFiles
```

Update golden files when generator changes:
```bash
go test -run TestGenerator_UpdateGoldenFiles -update-golden
```

### 4. Example Handler Tests

Example tests demonstrate how to test handlers using generated interfaces:

- **Handler Tests** (`example/handler_test.go`): Shows testing patterns for business logic
- **Mock Integration**: Examples of dependency injection and mocking

Run example tests:
```bash
go test -v ./example/
```

### 5. End-to-End Tests

Complete workflow tests that validate the entire process:

- Generate commands from documentation
- Compile generated code
- Execute commands and verify behavior

## Testing Patterns

### Testing Handlers

Using the adder testing utilities:

```go
func TestMyHandler_HandleCommand(t *testing.T) {
    tu := adder.NewTestingUtils()
    handler := NewMyHandler()
    cmd := tu.NewMockCobraCommand("my-command")
    
    req := &generated.MyRequest{
        Arguments: generated.MyRequestArguments{
            Name: "test",
        },
        Flags: generated.MyRequestFlags{
            Flag: "value",
        },
    }
    
    err := handler.HandleCommand(cmd, req)
    tu.AssertNoError(t, err)
}
```

Using request builders for complex requests:

```go
func TestMyHandler_WithBuilder(t *testing.T) {
    tu := adder.NewTestingUtils()
    
    // Use the generic builder
    builder := adder.NewRequestBuilder().
        WithFlag("input", "testdata/commands").
        WithFlag("output", "testdata/output").
        WithFlag("validate", true).
        WithArg("name", "test-arg")
    
    flags := builder.BuildFlags()
    args := builder.BuildArgs()
    
    req := &generated.MyRequest{
        Flags: generated.MyRequestFlags{
            Input: flags["input"].(string),
            Output: flags["output"].(string),
            Validate: flags["validate"].(bool),
        },
        Arguments: generated.MyRequestArguments{
            Name: args["name"].(string),
        },
    }
    
    handler := NewMyHandler()
    cmd := tu.NewMockCobraCommand("my-command")
    err := handler.HandleMyCommand(cmd, req)
    tu.AssertNoError(t, err)
}
```

Creating type-safe builders for specific commands:

```go
// Create a builder for your specific request type
type MyRequestBuilder struct {
    input    string
    output   string
    validate bool
    name     string
}

func NewMyRequestBuilder() *MyRequestBuilder {
    return &MyRequestBuilder{
        input:  "docs/commands",
        output: "generated",
    }
}

func (b *MyRequestBuilder) WithInput(input string) *MyRequestBuilder {
    b.input = input
    return b
}

func (b *MyRequestBuilder) WithValidate(validate bool) *MyRequestBuilder {
    b.validate = validate
    return b
}

func (b *MyRequestBuilder) WithName(name string) *MyRequestBuilder {
    b.name = name
    return b
}

func (b *MyRequestBuilder) Build() *generated.MyRequest {
    return &generated.MyRequest{
        Flags: generated.MyRequestFlags{
            Input:    b.input,
            Output:   b.output,
            Validate: b.validate,
        },
        Arguments: generated.MyRequestArguments{
            Name: b.name,
        },
    }
}

// Usage in tests:
func TestWithTypeSeafeBuilder(t *testing.T) {
    tu := adder.NewTestingUtils()
    
    req := NewMyRequestBuilder().
        WithInput("testdata/commands").
        WithValidate(true).
        WithName("test-name").
        Build()
    
    handler := NewMyHandler()
    cmd := tu.NewMockCobraCommand("my-command")
    err := handler.HandleMyCommand(cmd, req)
    tu.AssertNoError(t, err)
}
```
```

### Testing CLI Commands

```go
func TestGenerateCommand(t *testing.T) {
    handler := NewGenerateHandler()
    cmd := generated.NewGenerateCommand(handler)
    
    cmd.SetArgs([]string{
        "--input", "testdata",
        "--output", "output",
        "--package", "test",
    })
    
    err := cmd.Execute()
    if err != nil {
        t.Fatalf("Command failed: %v", err)
    }
}
```

### Dependency Injection

For handlers with external dependencies, use dependency injection:

```go
type MyHandler struct {
    service ExternalService
}

func NewMyHandler(service ExternalService) *MyHandler {
    return &MyHandler{service: service}
}

// In tests, inject a mock service
func TestMyHandler_WithMock(t *testing.T) {
    mockService := &MockService{}
    handler := NewMyHandler(mockService)
    // ... test with mock
}
```

## Test Organization

```
/
├── parser_test.go           # Unit tests for parser
├── generator_test.go        # Unit tests for generator  
├── golden_test.go          # Golden file validation tests
├── cmd/
│   └── integration_test.go # CLI integration tests
├── example/
│   └── handler_test.go     # Example handler tests
└── testdata/
    └── golden/             # Golden file test data
        ├── simple.md
        └── simple_generated.go.golden
```

## Running All Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific test categories
go test -v -run TestParser ./...        # Only parser tests
go test -v -run TestGenerator ./...     # Only generator tests
go test -v -run TestCLI ./cmd/          # Only CLI tests
```

## Test Data

Test data is organized in the `testdata/` directory:

- `testdata/golden/`: Golden files for output validation
- Test markdown files with various command configurations
- Expected generated Go code output

## Continuous Integration

Tests should be run in CI with:

```yaml
- name: Run tests
  run: |
    go test -v -race -cover ./...
    
- name: Validate golden files
  run: |
    go test -v -run TestGenerator_GoldenFiles
```

## Writing New Tests

When adding new features:

1. **Add unit tests** for new functions and methods
2. **Update integration tests** if CLI behavior changes  
3. **Add golden files** for new template outputs
4. **Create example tests** for new handler patterns
5. **Update this documentation** with new testing patterns