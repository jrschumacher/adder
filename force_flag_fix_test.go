package adder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGeneratedFilesAutoRegenerate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "adder-force-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create input directory with a test command
	inputDir := filepath.Join(tempDir, "docs", "commands")
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input dir: %v", err)
	}

	// Also create output directory
	outputDir := filepath.Join(tempDir, "generated")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Create a simple test command
	testCommand := `---
title: Test Command
command:
  name: test
  arguments:
    - name: arg1
      description: Test argument
      required: true
      type: string
---

# Test Command

A test command for validating force flag behavior.
`

	testFile := filepath.Join(inputDir, "test.md")
	if err := os.WriteFile(testFile, []byte(testCommand), 0644); err != nil {
		t.Fatalf("Failed to write test command: %v", err)
	}

	// Create config
	config := &Config{
		BinaryName:          "testapp",
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "directory",
	}

	// Create generator (without force flag)
	generator := NewGenerator(config)
	generator.SetForceRegeneration(false) // Explicitly disable force

	// First generation - should create the file
	inputFS := os.DirFS(inputDir)
	if err := generator.Generate(context.Background(), inputFS); err != nil {
		t.Fatalf("First generation failed: %v", err)
	}

	// Debug: List what was actually generated
	if entries, err := os.ReadDir(outputDir); err == nil {
		t.Logf("Generated files in %s:", outputDir)
		for _, entry := range entries {
			t.Logf("  - %s", entry.Name())
		}
	}

	// Check that the file was created
	expectedOutput := filepath.Join(config.OutputDir, "test_generated.go")
	if _, err := os.Stat(expectedOutput); err != nil {
		// List all files in the output directory for debugging
		if entries, err2 := os.ReadDir(outputDir); err2 == nil {
			t.Logf("Files in output directory %s:", outputDir)
			for _, entry := range entries {
				t.Logf("  %s", entry.Name())
			}
		}
		t.Fatalf("Generated file doesn't exist: %v", err)
	}

	t.Logf("✅ First generation successful: %s", expectedOutput)

	// Get the modification time of the generated file
	firstGenInfo, err := os.Stat(expectedOutput)
	if err != nil {
		t.Fatalf("Failed to stat generated file: %v", err)
	}

	// Wait a bit to ensure different modification times
	time.Sleep(10 * time.Millisecond)

	// Touch the source file to make it newer
	now := time.Now()
	if err := os.Chtimes(testFile, now, now); err != nil {
		t.Fatalf("Failed to touch source file: %v", err)
	}

	// Second generation - should regenerate without --force flag
	if err := generator.Generate(context.Background(), inputFS); err != nil {
		t.Fatalf("Second generation failed: %v", err)
	}

	// Check that the file was regenerated (modification time changed)
	secondGenInfo, err := os.Stat(expectedOutput)
	if err != nil {
		t.Fatalf("Failed to stat regenerated file: %v", err)
	}

	if !secondGenInfo.ModTime().After(firstGenInfo.ModTime()) {
		t.Errorf("❌ Generated file was NOT regenerated automatically!")
		t.Errorf("First gen time: %v", firstGenInfo.ModTime())
		t.Errorf("Second gen time: %v", secondGenInfo.ModTime())
		t.Errorf("Expected second generation to update the file without --force flag")
	} else {
		t.Logf("✅ Generated file was automatically regenerated without --force flag")
		t.Logf("First gen time: %v", firstGenInfo.ModTime())
		t.Logf("Second gen time: %v", secondGenInfo.ModTime())
	}
}

func TestNonGeneratedFilesStillRequireForceOrSourceChange(t *testing.T) {
	// Test that our change doesn't break the existing behavior for non-generated files
	
	tempDir, err := os.MkdirTemp("", "adder-force-test-non-gen-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	config := &Config{
		GeneratedFileSuffix: "_generated.go", // Our test file won't have this suffix
	}

	generator := NewGenerator(config)
	generator.SetForceRegeneration(false)

	// Create a fake source and output file
	sourceFile := filepath.Join(tempDir, "source.md")
	outputFile := filepath.Join(tempDir, "output.go") // Note: doesn't have _generated.go suffix
	
	// Create source file
	if err := os.WriteFile(sourceFile, []byte("source content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create output file that's newer than source
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(outputFile, []byte("output content"), 0644); err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}

	// Check that non-generated file is not regenerated
	shouldRegen, err := generator.shouldRegenerateFile(sourceFile, outputFile)
	if err != nil {
		t.Fatalf("shouldRegenerateFile failed: %v", err)
	}

	if shouldRegen {
		t.Errorf("❌ Non-generated file should not be regenerated when output is newer")
	} else {
		t.Logf("✅ Non-generated files still respect timestamp logic")
	}

	// But generated files should always be regenerated
	generatedFile := filepath.Join(tempDir, "output_generated.go")
	if err := os.WriteFile(generatedFile, []byte("generated content"), 0644); err != nil {
		t.Fatalf("Failed to create generated file: %v", err)
	}

	shouldRegenGenerated, err := generator.shouldRegenerateFile(sourceFile, generatedFile)
	if err != nil {
		t.Fatalf("shouldRegenerateFile failed for generated file: %v", err)
	}

	if !shouldRegenGenerated {
		t.Errorf("❌ Generated file should always be regenerated")
	} else {
		t.Logf("✅ Generated files are always regenerated")
	}
}