package adder

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCobraFlagDefaultBehavior(t *testing.T) {
	// This test demonstrates the core issue: Cobra always returns default values
	// even when flags aren't provided, making it impossible to distinguish
	// between "user provided default value" vs "user didn't provide flag at all"

	// Create a command similar to our generate command
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// Register flag with same default as our generate command
	cmd.Flags().StringP("output", "o", "generated", "Output directory")

	tests := []struct {
		name        string
		args        []string
		description string
	}{
		{
			name:        "no_flag_provided",
			args:        []string{}, // No --output flag
			description: "User doesn't provide --output flag at all",
		},
		{
			name:        "explicit_default_value",
			args:        []string{"--output", "generated"}, // Explicit --output generated
			description: "User explicitly provides --output generated",
		},
		{
			name:        "custom_value",
			args:        []string{"--output", "custom/path"}, // Custom value
			description: "User provides custom --output value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			cmd.ResetFlags()
			cmd.Flags().StringP("output", "o", "generated", "Output directory")

			// Set args and parse
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Command execution failed: %v", err)
			}

			// Get flag value
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				t.Fatalf("Failed to get flag: %v", err)
			}

			t.Logf("%s: args=%v -> output=%q", tt.description, tt.args, output)

			// The problem: we can't distinguish case 1 from case 2!
			if tt.name == "no_flag_provided" || tt.name == "explicit_default_value" {
				if output != "generated" {
					t.Errorf("Expected 'generated', got %q", output)
				}
				if tt.name == "no_flag_provided" {
					t.Logf("❌ PROBLEM: Can't tell this case apart from explicit default!")
				}
			}
		})
	}
}

func TestRealWorldConfigScenario(t *testing.T) {
	// Simulate the real issue: config file has custom output, but merge logic ignores it
	
	// Simulate config file content
	configFromFile := &Config{
		BinaryName:          "myapp",
		InputDir:            "docs/commands", 
		OutputDir:           "cmd/generated", // Custom output in config
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		PackageStrategy:     "directory",
	}

	// Simulate what happens when user runs: adder generate
	// (no --output flag provided, so Cobra returns default)
	
	// Simulate what the generate command would receive
	type mockFlags struct {
		BinaryName      string
		Input           string
		Output          string
		Package         string
		Suffix          string
		PackageStrategy string
	}
	
	req := mockFlags{
		BinaryName:      "", // Not provided
		Input:           "docs/commands", // Default 
		Output:          "generated",     // ❌ PROBLEM: Cobra default, but we can't tell it wasn't provided
		Package:         "generated",     // Default
		Suffix:          "_generated.go", // Default
		PackageStrategy: "directory",     // Default
	}

	// Merge with current logic
	merged := MergeWithFlags(configFromFile, 
		req.BinaryName,
		req.Input,
		req.Output,    // This is "generated" even though user didn't provide --output
		req.Package,
		req.Suffix)

	t.Logf("Config file output: %s", configFromFile.OutputDir)
	t.Logf("Flag value (default): %s", req.Output)
	t.Logf("Merged result: %s", merged.OutputDir)

	// The bug: if config has non-default value, it gets ignored
	if configFromFile.OutputDir != "generated" && merged.OutputDir == "generated" {
		t.Errorf("❌ BUG CONFIRMED: Config file output directory ignored!")
		t.Errorf("Config file specified: %s", configFromFile.OutputDir)
		t.Errorf("But merge result is: %s", merged.OutputDir)
		t.Errorf("This happens because we can't distinguish 'no flag' from 'explicit default'")
	}
}

func TestFixedMergeLogic(t *testing.T) {
	// Test a potential fix: change the merge logic
	
	// Instead of checking against default values, we could:
	// 1. Use flag.Changed() to detect if flag was actually set
	// 2. Or change the merge logic to be more permissive
	
	configFromFile := &Config{
		OutputDir: "cmd/generated", // Non-default config value
	}

	// Current (broken) logic
	output := "generated" // Default from Cobra
	var mergedCurrent string
	if output != "" && output != "generated" {
		mergedCurrent = output // Flag value
	} else {
		mergedCurrent = configFromFile.OutputDir // Config value
	}

	// Fixed logic: always prefer config unless flag is non-default
	// OR better: always prefer config unless we know flag was explicitly set
	var mergedFixed string
	if output != "generated" { // Only use flag if it's non-default
		mergedFixed = output
	} else {
		mergedFixed = configFromFile.OutputDir // Use config value
	}

	t.Logf("Config: %s, Flag: %s", configFromFile.OutputDir, output)
	t.Logf("Current logic result: %s", mergedCurrent)
	t.Logf("Fixed logic result: %s", mergedFixed)

	if mergedFixed != configFromFile.OutputDir {
		t.Errorf("Even the fixed logic doesn't work perfectly")
	} else {
		t.Logf("✅ Fixed logic respects config file")
	}
}