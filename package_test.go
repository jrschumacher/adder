package adder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageNameConsistency(t *testing.T) {
	// Test case to reproduce the package name inconsistency reported
	tests := []struct {
		name                string
		files               map[string]string
		packageStrategy     string
		indexFormat         string
		expectedPackageName string
		expectConsistency   bool
	}{
		{
			name: "directory_strategy_with_mixed_files",
			files: map[string]string{
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
			},
			packageStrategy:     "directory",
			indexFormat:         "_index",
			expectedPackageName: "dev_selectors",
			expectConsistency:   true,
		},
		{
			name: "reported_issue_mixed_packages",
			files: map[string]string{
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
			},
			packageStrategy:     "directory", // This was the reported case
			indexFormat:         "_index", 
			expectedPackageName: "dev_selectors", // All should be consistent
			expectConsistency:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "adder-package-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer func() { _ = os.RemoveAll(tempDir) }()

			inputDir := filepath.Join(tempDir, "input")
			outputDir := filepath.Join(tempDir, "output")

			// Create input files
			for relPath, content := range tt.files {
				fullPath := filepath.Join(inputDir, relPath)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("Failed to create dir for %s: %v", relPath, err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write %s: %v", relPath, err)
				}
			}

			// Create config
			config := &Config{
				BinaryName:          "testapp",
				InputDir:            inputDir,
				OutputDir:           outputDir,
				Package:             "generated",
				GeneratedFileSuffix: "_generated.go",
				PackageStrategy:     tt.packageStrategy,
				IndexFormat:         tt.indexFormat,
			}

			// Generate commands
			adder := New(config)
			commands, err := adder.ParseCommands()
			if err != nil {
				t.Fatalf("Failed to parse commands: %v", err)
			}

			// Check package names for commands
			packageNames := make(map[string]string)
			commandsByDir := make(map[string][]*Command)

			for _, cmd := range commands {
				dir := filepath.Dir(cmd.FilePath)
				commandsByDir[dir] = append(commandsByDir[dir], cmd)
				packageName := config.GetPackageName(cmd.FilePath)
				packageNames[cmd.FilePath] = packageName
			}

			// Check that all commands in the same directory get the same package name
			for dir, dirCommands := range commandsByDir {
				if len(dirCommands) > 1 {
					expectedPackage := config.GetPackageName(dirCommands[0].FilePath)
					for i, cmd := range dirCommands {
						cmdPackageName := config.GetPackageName(cmd.FilePath)
						if cmdPackageName != expectedPackage {
							t.Errorf("Package name inconsistency in directory %s: command %d (%s) has package %s, but first command has %s",
								dir, i, cmd.FilePath, cmdPackageName, expectedPackage)
						}
					}
				}
			}

			// For debugging: print all package names
			t.Logf("Package names generated:")
			for file, pkgName := range packageNames {
				t.Logf("  %s -> %s", file, pkgName)
			}

			// Verify consistency - all files in the same directory should have the same package name
			if tt.expectConsistency {
				expectedPkg := tt.expectedPackageName
				for _, pkgName := range packageNames {
					if pkgName != expectedPkg {
						t.Errorf("Expected all package names to be %s, but got %s", expectedPkg, pkgName)
					}
				}
			}
		})
	}
}

func TestGetPackageNameForIndexFiles(t *testing.T) {
	tests := []struct {
		name            string
		filePath        string
		packageStrategy string
		indexFormat     string
		expected        string
	}{
		{
			name:            "directory_strategy_index_file",
			filePath:        "dev/selectors/_index.md",
			packageStrategy: "directory",
			indexFormat:     "_index",
			expected:        "dev_selectors",
		},
		{
			name:            "path_strategy_index_file_should_use_directory",
			filePath:        "dev/selectors/_index.md", 
			packageStrategy: "path",
			indexFormat:     "_index",
			expected:        "dev_selectors", // Should be directory name, not dev_selectors_index
		},
		{
			name:            "path_strategy_regular_file",
			filePath:        "dev/selectors/generate.md",
			packageStrategy: "path",
			indexFormat:     "_index",
			expected:        "dev_selectors_generate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Package:         "generated",
				PackageStrategy: tt.packageStrategy,
				IndexFormat:     tt.indexFormat,
			}

			result := config.GetPackageName(tt.filePath)
			if result != tt.expected {
				t.Errorf("GetPackageName(%s) = %s, expected %s", tt.filePath, result, tt.expected)
			}
		})
	}
}