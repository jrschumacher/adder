package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jrschumacher/adder/cmd/generated"
	"github.com/spf13/cobra"
)

func TestGenerateHandler_HandleGenerate(t *testing.T) {
	// Create temporary directories
	tempDir, err := os.MkdirTemp("", "adder-cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	// Create input directory
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input dir: %v", err)
	}

	// Create test markdown file
	testMarkdown := `---
title: CLI Test Command
command:
  name: clitest
  flags:
    - name: verbose
      type: bool
      description: Enable verbose output
---

# CLI Test Command

This is a test command for CLI integration testing.`

	if err := os.WriteFile(filepath.Join(inputDir, "clitest.md"), []byte(testMarkdown), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create handler and request
	handler := NewGenerateHandler()
	req := &generated.GenerateRequest{
		Flags: generated.GenerateRequestFlags{
			Input:   inputDir,
			Output:  outputDir,
			Package: "testcli",
			Suffix:  "_generated.go",
		},
	}

	// Create mock command for context
	cmd := &cobra.Command{}

	// Test the handler
	err = handler.HandleGenerate(cmd, req)
	if err != nil {
		t.Fatalf("HandleGenerate failed: %v", err)
	}

	// Verify output file was created
	expectedFile := filepath.Join(outputDir, "clitest_generated.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s was not created", expectedFile)
	}

	// Verify content
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package testcli") {
		t.Errorf("Generated file missing expected package declaration")
	}
	if !strings.Contains(contentStr, "ClitestHandler") {
		t.Errorf("Generated file missing expected handler interface")
	}
}

func TestVersionHandler_HandleVersion(t *testing.T) {
	// Create handler
	handler := NewVersionHandler()
	req := &generated.VersionRequest{}

	// Create mock command
	cmd := &cobra.Command{}

	// Test the handler
	err := handler.HandleVersion(cmd, req)
	if err != nil {
		t.Fatalf("HandleVersion failed: %v", err)
	}

	// Note: This test would be improved by injecting output writers
	// For now, we just verify it doesn't error
}

func TestCLI_EndToEnd(t *testing.T) {
	// Create temporary directories
	tempDir, err := os.MkdirTemp("", "adder-e2e-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	inputDir := filepath.Join(tempDir, "docs")
	outputDir := filepath.Join(tempDir, "generated")

	// Create input directory
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input dir: %v", err)
	}

	// Create test command documentation
	testDocs := map[string]string{
		"simple.md": `---
title: Simple Command
command:
  name: simple
---

# Simple Command

A simple test command.`,
		"complex.md": `---
title: Complex Command
command:
  name: complex [arg]
  arguments:
    - name: arg
      description: Test argument
      required: true
      type: string
  flags:
    - name: flag
      shorthand: f
      description: Test flag
      type: bool
---

# Complex Command

A complex test command with arguments and flags.`,
	}

	for filename, content := range testDocs {
		if err := os.WriteFile(filepath.Join(inputDir, filename), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	// Create and execute generate command
	generateHandler := NewGenerateHandler()
	generateCmd := generated.NewGenerateCommand(generateHandler)

	// Set up command arguments
	generateCmd.SetArgs([]string{
		"--input", inputDir,
		"--output", outputDir,
		"--package", "e2etest",
	})

	// Execute command
	if err := generateCmd.Execute(); err != nil {
		t.Fatalf("Generate command failed: %v", err)
	}

	// Verify outputs
	expectedFiles := []string{
		filepath.Join(outputDir, "simple_generated.go"),
		filepath.Join(outputDir, "complex_generated.go"),
	}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Errorf("Expected output file %s was not created", expectedFile)
		}
	}

	// Verify content of complex command (has more structures to check)
	complexContent, err := os.ReadFile(filepath.Join(outputDir, "complex_generated.go"))
	if err != nil {
		t.Fatalf("Failed to read complex output file: %v", err)
	}

	complexStr := string(complexContent)
	expectedInComplex := []string{
		"package e2etest",
		"ComplexRequestArguments",
		"ComplexRequestFlags",
		"ComplexRequest",
		"ComplexHandler",
		"HandleComplex",
	}

	for _, expected := range expectedInComplex {
		if !strings.Contains(complexStr, expected) {
			t.Errorf("Complex generated file missing expected string: %q", expected)
		}
	}
}

func TestCLI_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		req         *generated.GenerateRequest
		expectError bool
	}{
		{
			name: "nonexistent input directory",
			req: &generated.GenerateRequest{
				Flags: generated.GenerateRequestFlags{
					Input:   "/nonexistent/directory",
					Output:  "/tmp/test-output",
					Package: "test",
					Suffix:  "_generated.go",
				},
			},
			expectError: true,
		},
		{
			name: "empty package name",
			req: &generated.GenerateRequest{
				Flags: generated.GenerateRequestFlags{
					Input:   "testdata",
					Output:  "/tmp/test-output",
					Package: "",
					Suffix:  "_generated.go",
				},
			},
			expectError: false, // Should use default
		},
	}

	handler := NewGenerateHandler()
	cmd := &cobra.Command{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.HandleGenerate(cmd, tt.req)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
