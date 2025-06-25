package adder

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestingUtils provides utilities for testing generated command handlers
type TestingUtils struct{}

// NewTestingUtils creates a new instance of testing utilities
func NewTestingUtils() *TestingUtils {
	return &TestingUtils{}
}

// AssertNoError fails the test if err is not nil
func (tu *TestingUtils) AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError fails the test if err is nil or doesn't contain expectedMsg
func (tu *TestingUtils) AssertError(t *testing.T, err error, expectedMsg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected error containing %q, got nil", expectedMsg)
	}
	if expectedMsg != "" && !stringContains(err.Error(), expectedMsg) {
		t.Fatalf("Expected error containing %q, got: %v", expectedMsg, err)
	}
}

// NewMockCobraCommand creates a mock cobra command for testing
func (tu *TestingUtils) NewMockCobraCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Mock command for testing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}

// Helper function to check if string contains substring (simple implementation)
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) <= len(s) && s[:len(substr)] == substr) ||
		stringContains(s[1:], substr))
}

// RequestBuilder provides a fluent interface for building test requests
type RequestBuilder struct {
	flags map[string]interface{}
	args  map[string]interface{}
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		flags: make(map[string]interface{}),
		args:  make(map[string]interface{}),
	}
}

// WithFlag adds a flag value to the request
func (rb *RequestBuilder) WithFlag(name string, value interface{}) *RequestBuilder {
	rb.flags[name] = value
	return rb
}

// WithArg adds an argument value to the request
func (rb *RequestBuilder) WithArg(name string, value interface{}) *RequestBuilder {
	rb.args[name] = value
	return rb
}

// BuildFlags returns the flags map
func (rb *RequestBuilder) BuildFlags() map[string]interface{} {
	return rb.flags
}

// BuildArgs returns the arguments map
func (rb *RequestBuilder) BuildArgs() map[string]interface{} {
	return rb.args
}