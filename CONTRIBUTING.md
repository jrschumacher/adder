# Contributing to Adder

Thank you for your interest in contributing to Adder! This guide will help you get started.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Docker (optional, for container testing)
- golangci-lint (for linting)

### Setting Up Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/adder.git
   cd adder
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Verify Setup**
   ```bash
   make ci-test
   ```

## 🏗️ Development Workflow

### 1. Create a Branch

Use descriptive branch names:
```bash
git checkout -b feat/add-validation-support
git checkout -b fix/parser-error-handling
git checkout -b docs/improve-readme
```

### 2. Make Changes

**Code Style:**
- Follow Go conventions and idioms
- Add comments for exported functions and types
- Keep functions focused and testable
- Use meaningful variable names

**Documentation:**
- Update README.md if adding features
- Add godoc comments for public APIs
- Update example code if interfaces change

### 3. Add Tests

**Required for all changes:**
- Unit tests for new functions
- Integration tests for CLI changes
- Golden file tests for template updates
- Example tests for new patterns

**Run tests locally:**
```bash
make test           # All tests
make test-coverage  # With coverage report
```

### 4. Validate Code Quality

**Run all quality checks:**
```bash
make ci-test        # Runs all CI checks locally
```

**Individual checks:**
```bash
make lint           # golangci-lint
make fmt            # Format code
make vet            # Go vet
```

## 📝 Commit Guidelines

Use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

### Commit Types

- `feat:` - New features (minor version bump)
- `fix:` - Bug fixes (patch version bump)
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring without feature changes
- `chore:` - Maintenance tasks
- `BREAKING CHANGE:` - Breaking changes (major version bump)

### Examples

```bash
feat: add support for enum validation in flags

fix: handle empty markdown files gracefully

docs: add examples for complex command structures

test: add golden file tests for template generation

BREAKING CHANGE: change handler signature to include context

The handler interface now receives context.Context as the first
parameter to support cancellation and timeouts.
```

## 🧪 Testing Guidelines

### Test Categories

1. **Unit Tests** (`*_test.go`)
   - Test individual functions and methods
   - Use table-driven tests for multiple scenarios
   - Mock external dependencies

2. **Integration Tests** (`cmd/integration_test.go`)
   - Test CLI command execution
   - Verify file generation and output
   - Test error handling and edge cases

3. **Golden File Tests** (`golden_test.go`)
   - Ensure generated code consistency
   - Compare against reference outputs
   - Update with `go test -run TestGenerator_UpdateGoldenFiles -update-golden`

4. **Example Tests** (`example/handler_test.go`)
   - Demonstrate testing patterns
   - Test generated interfaces
   - Show dependency injection examples

### Writing Good Tests

```go
func TestParser_ParseContent(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        want     *Command
        wantErr  bool
    }{
        {
            name: "valid simple command",
            content: `---
title: Test Command
command:
  name: test
---
# Test Command`,
            want: &Command{
                Title: "Test Command",
                Name:  "test",
                // ...
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewParser(&Config{})
            got, err := parser.ParseContent(tt.content, "test.md")
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseContent() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            // Validate results...
        })
    }
}
```

## 🔄 Pull Request Process

### 1. Pre-submission Checklist

- [ ] All tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] Branch is up to date with main

### 2. Pull Request Template

```markdown
## Description
Brief description of the changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change (fix or feature causing existing functionality to change)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

### 3. Review Process

1. **Automated Checks** - CI must pass
2. **Code Review** - At least one maintainer approval
3. **Testing** - Verify functionality works as expected
4. **Documentation** - Ensure docs are clear and complete

## 🐛 Bug Reports

### Before Reporting

1. Check existing [issues](https://github.com/jrschumacher/adder/issues)
2. Try with the latest version
3. Provide minimal reproduction case

### Bug Report Template

```markdown
**Describe the Bug**
Clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. Create file with content '...'
2. Run command '...'
3. See error

**Expected Behavior**
What you expected to happen.

**Environment**
- OS: [e.g., macOS, Linux, Windows]
- Go version: [e.g., 1.21.0]
- Adder version: [e.g., v1.0.0]

**Additional Context**
Any other context about the problem.
```

## 💡 Feature Requests

### Before Requesting

1. Check if it aligns with project goals
2. Consider if it can be implemented as a plugin
3. Search existing issues and discussions

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
Clear description of the problem.

**Describe the solution you'd like**
Clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or features considered.

**Additional context**
Any other context, mockups, or examples.
```

## 📞 Getting Help

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - Questions and general discussion
- **Code Review** - Ask for feedback on your contributions

## 🏆 Recognition

Contributors are recognized in:
- Release notes for significant contributions
- README contributors section
- Git commit history

Thank you for contributing to Adder! 🎉