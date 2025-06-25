package adder

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "update golden files")

func TestGenerator_GoldenFiles(t *testing.T) {
	// Create temporary output directory
	tempDir, err := os.MkdirTemp("", "adder-golden-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	outputDir := filepath.Join(tempDir, "output")

	// Configure generator with golden test data
	config := &Config{
		InputDir:   "testdata/golden",
		OutputDir:  outputDir,
		Package:    "testpkg",
		GeneratedFileSuffix: "_generated.go",
	}

	generator := New(config)

	// Generate code
	ctx := context.Background()
	if err := generator.GenerateWithContext(ctx); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Compare generated output with golden files
	goldenDir := "testdata/golden"

	err = filepath.Walk(goldenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-golden files
		if info.IsDir() || filepath.Ext(path) != ".golden" {
			return nil
		}

		// Read golden file
		goldenContent, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read golden file %s: %v", path, err)
			return nil
		}

		// Determine corresponding generated file
		relPath, err := filepath.Rel(goldenDir, path)
		if err != nil {
			t.Errorf("Failed to get relative path for %s: %v", path, err)
			return nil
		}

		// Remove .golden extension to get generated filename
		generatedFilename := relPath[:len(relPath)-len(".golden")]
		generatedPath := filepath.Join(outputDir, generatedFilename)

		// Read generated file
		generatedContent, err := os.ReadFile(generatedPath)
		if err != nil {
			t.Errorf("Failed to read generated file %s: %v", generatedPath, err)
			return nil
		}

		// Compare content
		if string(goldenContent) != string(generatedContent) {
			t.Errorf("Generated file %s does not match golden file %s", generatedPath, path)
			t.Logf("Golden content:\n%s", string(goldenContent))
			t.Logf("Generated content:\n%s", string(generatedContent))
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk golden directory: %v", err)
	}
}

func TestGenerator_UpdateGoldenFiles(t *testing.T) {
	// This test can be used to update golden files when the generator changes
	// Run with: go test -run TestGenerator_UpdateGoldenFiles -update-golden

	if !*updateGolden {
		t.Skip("Skipping golden file update test (use -update-golden to run)")
	}

	// Create temporary output directory
	tempDir, err := os.MkdirTemp("", "adder-update-golden-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	outputDir := filepath.Join(tempDir, "output")

	// Configure generator
	config := &Config{
		InputDir:   "testdata/golden",
		OutputDir:  outputDir,
		Package:    "testpkg",
		GeneratedFileSuffix: "_generated.go",
	}

	generator := New(config)

	// Generate code
	ctx := context.Background()
	if err := generator.GenerateWithContext(ctx); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Copy generated files to golden directory
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Read generated file
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Determine golden file path
		relPath, err := filepath.Rel(outputDir, path)
		if err != nil {
			return err
		}

		goldenPath := filepath.Join("testdata/golden", relPath+".golden")

		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
			return err
		}

		// Write golden file
		return os.WriteFile(goldenPath, content, 0644)
	})

	if err != nil {
		t.Fatalf("Failed to update golden files: %v", err)
	}

	t.Log("Golden files updated successfully")
}
