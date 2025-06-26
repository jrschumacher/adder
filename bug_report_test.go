package adder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageNameInconsistencyBugReport(t *testing.T) {
	// This test reproduces the exact bug reported from otdfctl
	// Issue: Files in the same directory get different package names
	// Expected: All files in same directory should have same package name
	
	tempDir, err := os.MkdirTemp("", "adder-bug-report-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	inputDir := filepath.Join(tempDir, "docs", "man")
	outputDir := filepath.Join(tempDir, "cmd", "generated")

	// Create the exact file structure that causes the issue
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
    - name: output
      type: string
      description: Output format
---
# Generate Command`,
		"dev/selectors/test.md": `---
title: Test Selector
command:
  name: test
  flags:
    - name: verbose
      type: bool
      description: Verbose output
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

	// Use configuration that should produce consistent package names
	config := &Config{
		BinaryName:          "testapp",
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "directory", // Should make all files use directory-based names
		IndexFormat:         "_index",
	}

	// Generate files (test both with and without --force behavior)
	adder := New(config)
	
	// First generation
	if err := adder.Generate(); err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Check what was generated
	expectedOutputDir := filepath.Join(outputDir, "dev", "selectors")
	packageNames := getPackageNamesFromDir(t, expectedOutputDir)

	t.Logf("Generated files and their package names:")
	for fileName, pkgName := range packageNames {
		t.Logf("  %s: package %s", fileName, pkgName)
	}

	// Verify consistency
	expectedPackage := "dev_selectors" // With directory strategy, should be this
	inconsistentFiles := []string{}
	
	for fileName, pkgName := range packageNames {
		if pkgName != expectedPackage {
			inconsistentFiles = append(inconsistentFiles, fileName)
		}
	}

	if len(inconsistentFiles) > 0 {
		t.Errorf("🐛 BUG CONFIRMED: Package name inconsistency in same directory!")
		t.Errorf("Directory: dev/selectors")
		t.Errorf("Expected all files to have package '%s'", expectedPackage)
		t.Errorf("Files with incorrect package names:")
		for _, fileName := range inconsistentFiles {
			t.Errorf("  ❌ %s has package '%s'", fileName, packageNames[fileName])
		}
		
		t.Errorf("\nBug Details:")
		t.Errorf("- All files are in the same directory: dev/selectors")
		t.Errorf("- Using package_strategy: 'directory'")
		t.Errorf("- Go requires all files in same directory to have same package name")
		t.Errorf("- This causes compilation errors in the generated code")
	} else {
		t.Logf("✅ All files have consistent package name: %s", expectedPackage)
	}

	// Test regeneration without --force (should work automatically)
	t.Logf("\n🔄 Testing regeneration (should work without --force)")
	if err := adder.Generate(); err != nil {
		t.Errorf("❌ Regeneration failed - this is a UX issue!")
		t.Errorf("Generated files should be automatically overwritten")
		t.Errorf("Users shouldn't need --force for generated code")
	} else {
		t.Logf("✅ Regeneration succeeded without --force")
	}
}

func TestForceNotNeededForGeneratedFiles(t *testing.T) {
	// Test that --force is not needed for generated files (UX improvement)
	tempDir, err := os.MkdirTemp("", "adder-force-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	// Create a simple test file
	testFile := filepath.Join(inputDir, "hello.md")
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input dir: %v", err)
	}
	
	content := `---
title: Hello Command
command:
  name: hello
  arguments:
    - name: name
      type: string
      required: true
---
# Hello Command`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	config := &Config{
		BinaryName:          "testapp",
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "directory",
		IndexFormat:         "directory",
	}

	adder := New(config)

	// First generation
	if err := adder.Generate(); err != nil {
		t.Fatalf("First generation failed: %v", err)
	}

	// Verify file was created
	expectedFile := filepath.Join(outputDir, "hello_generated.go")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Generated file was not created: %s", expectedFile)
	}

	// Second generation (should work without --force)
	if err := adder.Generate(); err != nil {
		t.Errorf("❌ UX ISSUE: Second generation failed without --force")
		t.Errorf("Users should not need --force to regenerate generated code")
		t.Errorf("Error: %v", err)
	} else {
		t.Logf("✅ Generated files can be regenerated without --force")
	}
}

// Helper function to extract package names from generated files
func getPackageNamesFromDir(t *testing.T, dir string) map[string]string {
	packageNames := make(map[string]string)
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), "_generated.go") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", filePath, err)
			continue
		}

		// Extract package name
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

	return packageNames
}