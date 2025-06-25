# ADR-0010: Function-Based Handlers

## Status

Accepted

## Context

The current handler generation approach requires developers to create structs that implement interfaces for each command handler. This results in significant boilerplate:

1. **Interface definition** (generated)
2. **Handler struct** with constructor
3. **Method implementation** on the struct
4. **Manual wiring** in init functions

For a typical command, this means ~15-20 lines of repetitive structural code before any business logic. In real usage (like otdfctl), developers create 6-7 handlers with identical patterns, resulting in 100+ lines of boilerplate.

## Decision

Replace interface-based handlers with function-based handlers, following Go's `http.HandlerFunc` pattern.

### Before (Interface-based)
```go
// Generated interface
type CreateHandler interface {
    HandleCreate(cmd *cobra.Command, req *CreateRequest) error
}

// User implements
type ProfileCreateHandler struct{}
func NewProfileCreateHandler() *ProfileCreateHandler { return &ProfileCreateHandler{} }
func (h *ProfileCreateHandler) HandleCreate(cmd *cobra.Command, req *CreateRequest) error {
    // business logic
}

// User wires up
createHandler := NewProfileCreateHandler()
createCmd := NewCreateCommand(createHandler)
```

### After (Function-based)
```go
// Generated function type
type CreateHandlerFunc func(cmd *cobra.Command, req *CreateRequest) error

// User implements inline
createCmd := NewCreateCommand(func(cmd *cobra.Command, req *CreateRequest) error {
    // business logic directly here
    return nil
})
```

## Consequences

### Positive
- **90% reduction in boilerplate** - eliminates structs, constructors, and interface implementations
- **Maintains full type safety** - still strongly typed with IDE autocompletion
- **Follows Go idioms** - similar to `http.HandlerFunc`, `sort.Interface` patterns
- **Inline business logic** - handler logic colocated with command registration
- **Easier testing** - function types are simpler to mock/test than interfaces

### Negative
- **Breaking change** - existing handler implementations need migration
- **Less structured** - some developers prefer explicit struct organization
- **Longer init functions** - business logic moves into init() instead of separate methods

### Neutral
- Generated code complexity remains similar
- Type safety and validation behavior unchanged
- Template modification required but straightforward

## Implementation Plan

1. Update generator templates to emit `HandlerFunc` types instead of interfaces
2. Modify command constructors to accept function parameters
3. Update documentation and examples
4. Provide migration guide for existing codebases

## Migration Path

Existing interface-based handlers can be easily converted:

```go
// Old way
func (h *ProfileCreateHandler) HandleCreate(cmd *cobra.Command, req *CreateRequest) error {
    // business logic
}

// New way  
createCmd := NewCreateCommand(func(cmd *cobra.Command, req *CreateRequest) error {
    // same business logic
    return nil
})
```

No changes to command definitions, request types, or validation - only the handler implementation pattern changes.