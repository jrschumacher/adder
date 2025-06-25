package adder

import (
	"errors"
	"testing"
)

func TestTestingUtils_AssertNoError(t *testing.T) {
	tu := NewTestingUtils()

	// Should not fail
	tu.AssertNoError(t, nil)

	// Test that it would fail with an error (can't test this directly without sub-test)
	// This is just to show the pattern
}

func TestTestingUtils_AssertError(t *testing.T) {
	tu := NewTestingUtils()

	// Test error with expected message
	err := errors.New("file not found: missing.md")
	tu.AssertError(t, err, "file not found")

	// Test error without checking message
	tu.AssertError(t, err, "")
}

func TestTestingUtils_NewMockCobraCommand(t *testing.T) {
	tu := NewTestingUtils()

	cmd := tu.NewMockCobraCommand("test")
	if cmd.Use != "test" {
		t.Errorf("Expected command Use=test, got %s", cmd.Use)
	}

	if cmd.Short != "Mock command for testing" {
		t.Errorf("Expected mock short description, got %s", cmd.Short)
	}

	// Test that RunE doesn't return error
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Expected mock command to return no error, got %v", err)
	}
}

func TestRequestBuilder(t *testing.T) {
	builder := NewRequestBuilder().
		WithFlag("input", "testdata/commands").
		WithFlag("output", "testdata/output").
		WithFlag("validate", true).
		WithArg("name", "test")

	flags := builder.BuildFlags()
	args := builder.BuildArgs()

	// Test flags
	if flags["input"] != "testdata/commands" {
		t.Errorf("Expected input=testdata/commands, got %v", flags["input"])
	}

	if flags["validate"] != true {
		t.Errorf("Expected validate=true, got %v", flags["validate"])
	}

	// Test args
	if args["name"] != "test" {
		t.Errorf("Expected name=test, got %v", args["name"])
	}
}

func TestStringContains(t *testing.T) {
	tests := []struct {
		s       string
		substr  string
		want    bool
	}{
		{"hello world", "world", true},
		{"hello world", "foo", false},
		{"", "", true},
		{"test", "", true},
		{"", "test", false},
	}

	for _, tt := range tests {
		got := stringContains(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("stringContains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}