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

```go
func TestMyHandler_HandleCommand(t *testing.T) {
    handler := NewMyHandler()
    cmd := &cobra.Command{}
    
    req := &generated.MyRequest{
        Arguments: generated.MyRequestArguments{
            Name: "test",
        },
        Flags: generated.MyRequestFlags{
            Flag: "value",
        },
    }
    
    err := handler.HandleCommand(cmd, req)
    if err != nil {
        t.Fatalf("Handler failed: %v", err)
    }
}
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