package adder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutputDirectoryConfigBug(t *testing.T) {
	// Reproduce the bug: config file output directory is ignored
	
	// Test the MergeWithFlags logic directly
	configFromFile := &Config{
		BinaryName:          "myapp",
		InputDir:            "docs/commands",
		OutputDir:           "cmd/generated", // Config file specifies this
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		IndexFormat:         "directory",
		PackageStrategy:     "directory",
	}

	// Simulate CLI flags with default values (no --output flag provided)
	binaryName := ""     // No flag provided
	input := "docs/commands"  // Default value
	output := "generated"     // Default value (no --output flag)
	pkg := "generated"        // Default value
	suffix := "_generated.go" // Default value

	// This simulates what happens when user doesn't provide --output flag
	merged := MergeWithFlags(configFromFile, binaryName, input, output, pkg, suffix)

	t.Logf("Config file output: %s", configFromFile.OutputDir)
	t.Logf("CLI flag output: %s", output)
	t.Logf("Merged output: %s", merged.OutputDir)

	expectedOutput := "cmd/generated" // Should use config file value
	if merged.OutputDir != expectedOutput {
		t.Errorf("❌ BUG CONFIRMED: Output directory config ignored!")
		t.Errorf("Expected: %s", expectedOutput)
		t.Errorf("Got: %s", merged.OutputDir)
		t.Errorf("The config file value 'cmd/generated' was overridden by default 'generated'")
	} else {
		t.Logf("✅ Output directory config respected")
	}
}

func TestConfigMergingBehavior(t *testing.T) {
	// Test the merging behavior for all fields
	configFromFile := &Config{
		BinaryName:          "myapp",
		InputDir:            "custom/input",
		OutputDir:           "custom/output",
		Package:             "custompackage",
		GeneratedFileSuffix: "_gen.go",
		IndexFormat:         "directory",
		PackageStrategy:     "path",
	}

	tests := []struct {
		name           string
		binaryName     string
		input          string
		output         string
		pkg            string
		suffix         string
		expectedOutput string
		expectedPkg    string
		description    string
	}{
		{
			name:           "no_flags_provided_should_use_config",
			binaryName:     "",
			input:          "docs/commands", // default
			output:         "generated",     // default
			pkg:            "generated",     // default
			suffix:         "_generated.go", // default
			expectedOutput: "custom/output", // Should use config
			expectedPkg:    "custompackage", // Should use config
			description:    "When no flags provided, should use config file values",
		},
		{
			name:           "explicit_flags_should_override_config",
			binaryName:     "",
			input:          "docs/commands",
			output:         "flag/output",   // explicit flag
			pkg:            "flagpackage",   // explicit flag
			suffix:         "_generated.go",
			expectedOutput: "flag/output",   // Should use flag
			expectedPkg:    "flagpackage",   // Should use flag
			description:    "When flags explicitly provided, should override config",
		},
		{
			name:           "mixed_flags_and_config",
			binaryName:     "",
			input:          "docs/commands",
			output:         "generated",       // default (should use config)
			pkg:            "explicitpackage", // explicit flag
			suffix:         "_generated.go",
			expectedOutput: "custom/output",   // Should use config
			expectedPkg:    "explicitpackage", // Should use flag
			description:    "Should mix config and flags appropriately",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := MergeWithFlags(configFromFile, tt.binaryName, tt.input, tt.output, tt.pkg, tt.suffix)

			t.Logf("Test: %s", tt.description)
			t.Logf("Config output: %s, flag output: %s -> merged: %s", 
				configFromFile.OutputDir, tt.output, merged.OutputDir)
			t.Logf("Config package: %s, flag package: %s -> merged: %s", 
				configFromFile.Package, tt.pkg, merged.Package)

			if merged.OutputDir != tt.expectedOutput {
				t.Errorf("❌ Output directory wrong. Expected: %s, Got: %s", 
					tt.expectedOutput, merged.OutputDir)
			}

			if merged.Package != tt.expectedPkg {
				t.Errorf("❌ Package wrong. Expected: %s, Got: %s", 
					tt.expectedPkg, merged.Package)
			}
		})
	}
}

func TestConfigFileLoadingIntegration(t *testing.T) {
	// Test the full integration: config file loading + flag merging + generation
	tempDir, err := os.MkdirTemp("", "adder-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create a config file with custom output directory
	configContent := `binary_name: testapp
input: docs/commands
output: cmd/generated
package: generated
generated_file_suffix: _generated.go
index_format: directory
package_strategy: directory`

	configPath := filepath.Join(tempDir, ".adder.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config from file
	loadedConfig, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Loaded config output directory: %s", loadedConfig.OutputDir)

	// Simulate command line with default flags (no --output provided)
	mergedConfig := MergeWithFlags(loadedConfig, "", "docs/commands", "generated", "generated", "_generated.go")

	t.Logf("After merging with default flags: %s", mergedConfig.OutputDir)

	expectedOutput := "cmd/generated"
	if mergedConfig.OutputDir != expectedOutput {
		t.Errorf("❌ Integration bug: Config file output directory ignored!")
		t.Errorf("Config file had: %s", loadedConfig.OutputDir)
		t.Errorf("After merge got: %s", mergedConfig.OutputDir)
		t.Errorf("Expected: %s", expectedOutput)
	} else {
		t.Logf("✅ Integration working: Config file output directory respected")
	}
}