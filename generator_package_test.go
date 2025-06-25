package adder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratorPackageConsistency(t *testing.T) {
	// Test actual file generation to ensure package names are consistent
	tempDir, err := os.MkdirTemp("", "adder-generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	// Create the exact structure from the bug report
	files := map[string]string{
		"dev/selectors/_index.md": `---
title: Dev Selectors
command:
  name: selectors
---
# Dev Selectors Root Command`,
		"dev/selectors/generate.md": `---
title: Generate Selector
command:
  name: generate
  flags:
    - name: verbose
      type: bool
      description: Enable verbose output
---
# Generate Command`,
		"dev/selectors/test.md": `---
title: Test Selector
command:
  name: test
  flags:
    - name: dry-run
      type: bool
      description: Perform a dry run
---
# Test Command`,
	}

	// Create input files
	for relPath, content := range files {
		fullPath := filepath.Join(inputDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create dir for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", relPath, err)
		}
	}

	// Create config with directory strategy (default)
	config := &Config{
		BinaryName:          "testapp",
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "directory",
		IndexFormat:         "_index",
	}

	// Generate the files
	adder := New(config)
	if err := adder.Generate(); err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Read the generated files and check package declarations
	expectedOutputDir := filepath.Join(outputDir, "dev", "selectors")
	dirEntries, err := os.ReadDir(expectedOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory %s: %v", expectedOutputDir, err)
	}

	packageNames := make(map[string]string)
	for _, entry := range dirEntries {
		if !strings.HasSuffix(entry.Name(), "_generated.go") {
			continue
		}

		filePath := filepath.Join(expectedOutputDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated file %s: %v", filePath, err)
		}

		// Extract package name from the file
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "package ") {
				packageName := strings.TrimSpace(strings.TrimPrefix(line, "package"))
				packageNames[entry.Name()] = packageName
				break
			}
		}
	}

	// Check that all files have the same package name
	expectedPackage := "dev_selectors"
	t.Logf("Generated files and their package names:")
	for fileName, pkgName := range packageNames {
		t.Logf("  %s: package %s", fileName, pkgName)
		if pkgName != expectedPackage {
			t.Errorf("File %s has package %s, expected %s", fileName, pkgName, expectedPackage)
		}
	}

	if len(packageNames) == 0 {
		t.Error("No generated files found")
	}
}

func TestGeneratorSingleStrategyConsistency(t *testing.T) {
	// Test with single strategy to see if this causes the reported issue
	tempDir, err := os.MkdirTemp("", "adder-single-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	// Create test files
	files := map[string]string{
		"dev/selectors/_index.md": `---
title: Dev Selectors
command:
  name: selectors
---
# Dev Selectors Root Command`,
		"dev/selectors/generate.md": `---
title: Generate Selector
command:
  name: generate
---
# Generate Command`,
	}

	// Create input files
	for relPath, content := range files {
		fullPath := filepath.Join(inputDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create dir for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", relPath, err)
		}
	}

	// Use single strategy - this might reproduce the issue
	config := &Config{
		BinaryName:          "testapp",
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "single", // This might cause inconsistency if incorrectly implemented
		IndexFormat:         "_index",
	}

	// Generate the files
	adder := New(config)
	if err := adder.Generate(); err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Check the generated package names
	expectedOutputDir := filepath.Join(outputDir, "dev", "selectors")
	dirEntries, err := os.ReadDir(expectedOutputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory %s: %v", expectedOutputDir, err)
	}

	packageNames := make(map[string]string)
	for _, entry := range dirEntries {
		if !strings.HasSuffix(entry.Name(), "_generated.go") {
			continue
		}

		filePath := filepath.Join(expectedOutputDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated file %s: %v", filePath, err)
		}

		// Extract package name from the file
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "package ") {
				packageName := strings.TrimSpace(strings.TrimPrefix(line, "package"))
				packageNames[entry.Name()] = packageName
				break
			}
		}
	}

	// With single strategy, all should be "generated"
	expectedPackage := "generated"
	t.Logf("Generated files with single strategy:")
	for fileName, pkgName := range packageNames {
		t.Logf("  %s: package %s", fileName, pkgName)
		if pkgName != expectedPackage {
			t.Errorf("File %s has package %s, expected %s", fileName, pkgName, expectedPackage)
		}
	}
}