package adder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageNameBugReproduction(t *testing.T) {
	// This test reproduces the exact bug reported from otdfctl
	tempDir, err := os.MkdirTemp("", "adder-bug-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	// Create the exact structure that causes the bug
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
		"dev/selectors/test.md": `---
title: Test Selector  
command:
  name: test
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

	// Create config with directory strategy (the reported configuration)
	config := &Config{
		BinaryName:          "testapp",
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             "generated", // Base package name
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "directory",
		IndexFormat:         "_index",
	}

	// Generate the files
	adder := New(config)
	if err := adder.Generate(); err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Read the generated files and extract package names
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

	// Log what we found (this will show the bug)
	t.Logf("Generated files and their package names:")
	for fileName, pkgName := range packageNames {
		t.Logf("  %s: package %s", fileName, pkgName)
	}

	// The bug: All files should have the same package name since they're in the same directory
	expectedPackage := "dev_selectors" // directory strategy should use this for ALL files
	inconsistentFiles := []string{}
	
	for fileName, pkgName := range packageNames {
		if pkgName != expectedPackage {
			inconsistentFiles = append(inconsistentFiles, fileName)
		}
	}

	if len(inconsistentFiles) > 0 {
		t.Errorf("BUG REPRODUCED: Files in same directory have inconsistent package names!")
		t.Errorf("Expected all files to have package '%s'", expectedPackage)
		t.Errorf("Files with wrong package names: %v", inconsistentFiles)
		
		// Show the exact issue
		for _, fileName := range inconsistentFiles {
			t.Errorf("  %s has package '%s' but should have '%s'", 
				fileName, packageNames[fileName], expectedPackage)
		}
	}
}

func TestDirectoryStrategyLogic(t *testing.T) {
	// Test the GetPackageName logic directly to understand the issue
	config := &Config{
		Package:         "generated",
		PackageStrategy: "directory",
		IndexFormat:     "_index",
	}

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "index_file_should_use_directory_name",
			filePath: "dev/selectors/_index.md",
			expected: "dev_selectors",
		},
		{
			name:     "regular_file_should_also_use_directory_name",
			filePath: "dev/selectors/generate.md",
			expected: "dev_selectors", // This should be the same as index file!
		},
		{
			name:     "another_regular_file_should_also_use_directory_name",
			filePath: "dev/selectors/test.md",
			expected: "dev_selectors", // This should be the same as index file!
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetPackageName(tt.filePath)
			if result != tt.expected {
				t.Errorf("GetPackageName(%s) = %s, expected %s", tt.filePath, result, tt.expected)
				t.Errorf("This shows the bug: files in same directory get different package names")
			}
		})
	}
}