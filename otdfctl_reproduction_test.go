package adder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOtdfctlScenarioReproduction(t *testing.T) {
	// Test using the exact same configuration that otdfctl might be using
	tempDir, err := os.MkdirTemp("", "adder-otdfctl-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	inputDir := filepath.Join(tempDir, "docs", "man")
	outputDir := filepath.Join(tempDir, "cmd", "generated")

	// Create the exact structure from otdfctl
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

	// Test multiple different configurations to see which one causes the bug
	testConfigs := []struct {
		name   string
		config *Config
	}{
		{
			name: "otdfctl_likely_config_v1",
			config: &Config{
				BinaryName:          "otdfctl",
				InputDir:            inputDir,
				OutputDir:           outputDir,
				Package:             "generated",
				GeneratedFileSuffix: "_generated.go",
				PackageStrategy:     "directory",
				IndexFormat:         "_index",
			},
		},
		{
			name: "otdfctl_possible_config_v2",
			config: &Config{
				BinaryName:          "otdfctl",
				InputDir:            inputDir,
				OutputDir:           outputDir,
				Package:             "generated",
				GeneratedFileSuffix: "_generated.go",
				PackageStrategy:     "single", // Maybe they're using single?
				IndexFormat:         "_index",
			},
		},
		{
			name: "otdfctl_possible_config_v3",
			config: &Config{
				BinaryName:          "otdfctl",
				InputDir:            inputDir,
				OutputDir:           outputDir,
				Package:             "generated",
				GeneratedFileSuffix: "_generated.go",
				PackageStrategy:     "", // Default behavior
				IndexFormat:         "_index",
			},
		},
	}

	for _, tc := range testConfigs {
		t.Run(tc.name, func(t *testing.T) {
			// Clean output directory
			os.RemoveAll(outputDir)
			
			// Generate with this config
			adder := New(tc.config)
			if err := adder.Generate(); err != nil {
				t.Fatalf("Failed to generate with %s: %v", tc.name, err)
			}

			// Check package names in the generated files
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

			// Log results
			t.Logf("Config %s generated files:", tc.name)
			for fileName, pkgName := range packageNames {
				t.Logf("  %s: package %s", fileName, pkgName)
			}

			// Check for inconsistencies
			if len(packageNames) > 1 {
				firstPackage := ""
				for _, pkgName := range packageNames {
					if firstPackage == "" {
						firstPackage = pkgName
					} else if pkgName != firstPackage {
						t.Errorf("INCONSISTENCY FOUND with config %s!", tc.name)
						t.Errorf("Files in same directory have different package names:")
						for fn, pn := range packageNames {
							t.Errorf("  %s: %s", fn, pn)
						}
						return
					}
				}
				t.Logf("✅ All files have consistent package name: %s", firstPackage)
			}
		})
	}
}

func TestEmptyPackageStrategyBehavior(t *testing.T) {
	// Test what happens when PackageStrategy is empty (default behavior)
	config := &Config{
		Package:         "generated",
		PackageStrategy: "", // Empty - should default to directory
		IndexFormat:     "_index",
	}

	testFiles := []string{
		"dev/selectors/_index.md",
		"dev/selectors/generate.md", 
		"dev/selectors/test.md",
	}

	t.Logf("Testing empty PackageStrategy behavior:")
	for _, filePath := range testFiles {
		packageName := config.GetPackageName(filePath)
		t.Logf("  %s -> package %s", filePath, packageName)
	}

	// Check if they're all the same
	packageNames := make(map[string]bool)
	for _, filePath := range testFiles {
		packageName := config.GetPackageName(filePath)
		packageNames[packageName] = true
	}

	if len(packageNames) > 1 {
		t.Errorf("Empty PackageStrategy causes inconsistent package names!")
		for pkgName := range packageNames {
			t.Errorf("  Found package: %s", pkgName)
		}
	}
}