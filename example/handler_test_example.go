package main

import (
	"fmt"
	"testing"

	"github.com/jrschumacher/adder"
	"github.com/jrschumacher/adder/cmd/adder/generated"
)

// Example of how to test generated handlers using the testing utilities

// TestGenerateHandler demonstrates testing a generated handler
func TestGenerateHandler(t *testing.T) {
	// Create testing utilities
	tu := adder.NewTestingUtils()

	// Create handler
	handler := NewGenerateHandler()

	// Create mock command
	cmd := tu.NewMockCobraCommand("generate")

	t.Run("successful generation", func(t *testing.T) {
		// Create test request using builder pattern  
		req := &generated.GenerateRequest{
			Flags: generated.GenerateRequestFlags{
				Input:   "testdata/commands",
				Output:  "testdata/output", 
				Package: "testpkg",
				Suffix:  "_test.go",
			},
		}

		// This would fail in a real test since directories don't exist,
		// but shows the pattern
		err := handler.HandleGenerate(cmd, req)
		
		// In a real test, you might mock the file system or use temp directories
		_ = err // Ignore for example purposes
	})

	t.Run("validation mode", func(t *testing.T) {
		req := &generated.GenerateRequest{
			Flags: generated.GenerateRequestFlags{
				Input:    "testdata/commands",
				Output:   "testdata/output",
				Package:  "testpkg", 
				Suffix:   "_test.go",
				Validate: true,
			},
		}

		err := handler.HandleGenerate(cmd, req)
		_ = err // Would handle appropriately in real test
	})
}

// TestGenerateRequestBuilder demonstrates using the request builder
func TestGenerateRequestBuilder(t *testing.T) {
	tu := adder.NewTestingUtils()

	// Use the generic builder for complex setup
	builder := adder.NewRequestBuilder().
		WithFlag("input", "docs/commands").
		WithFlag("output", "generated").
		WithFlag("validate", true).
		WithFlag("force", false)

	flags := builder.BuildFlags()

	// Verify flags were set correctly
	if flags["input"] != "docs/commands" {
		t.Errorf("Expected input=docs/commands, got %v", flags["input"])
	}

	if flags["validate"] != true {
		t.Errorf("Expected validate=true, got %v", flags["validate"])
	}

	// Example of asserting no error
	err := (error)(nil)
	tu.AssertNoError(t, err)

	// Example of asserting specific error
	err = fmt.Errorf("file not found: missing.md")
	tu.AssertError(t, err, "file not found")
}

// GenerateRequestBuilder provides type-safe building for GenerateRequest
type GenerateRequestBuilder struct {
	input    string
	output   string
	pkg      string
	suffix   string
	validate bool
	force    bool
}

// NewGenerateRequestBuilder creates a new builder with sensible defaults
func NewGenerateRequestBuilder() *GenerateRequestBuilder {
	return &GenerateRequestBuilder{
		input:  "docs/commands",
		output: "generated", 
		pkg:    "generated",
		suffix: "_generated.go",
	}
}

// WithInput sets the input directory
func (b *GenerateRequestBuilder) WithInput(input string) *GenerateRequestBuilder {
	b.input = input
	return b
}

// WithOutput sets the output directory  
func (b *GenerateRequestBuilder) WithOutput(output string) *GenerateRequestBuilder {
	b.output = output
	return b
}

// WithPackage sets the package name
func (b *GenerateRequestBuilder) WithPackage(pkg string) *GenerateRequestBuilder {
	b.pkg = pkg
	return b
}

// WithSuffix sets the file suffix
func (b *GenerateRequestBuilder) WithSuffix(suffix string) *GenerateRequestBuilder {
	b.suffix = suffix
	return b
}

// WithValidate sets validation mode
func (b *GenerateRequestBuilder) WithValidate(validate bool) *GenerateRequestBuilder {
	b.validate = validate
	return b
}

// WithForce sets force mode
func (b *GenerateRequestBuilder) WithForce(force bool) *GenerateRequestBuilder {
	b.force = force
	return b
}

// Build creates the GenerateRequest
func (b *GenerateRequestBuilder) Build() *generated.GenerateRequest {
	return &generated.GenerateRequest{
		Flags: generated.GenerateRequestFlags{
			Input:    b.input,
			Output:   b.output,
			Package:  b.pkg,
			Suffix:   b.suffix,
			Validate: b.validate,
			Force:    b.force,
		},
	}
}

// Example usage:
// req := NewGenerateRequestBuilder().
//     WithInput("test-docs").
//     WithValidate(true).
//     Build()