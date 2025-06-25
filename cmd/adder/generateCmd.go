package main

import (
	"context"
	"fmt"

	"github.com/jrschumacher/adder"
	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/spf13/cobra"
)

// generateCmd processes the generate command request to create CLI command stubs.
func generateCmd(_ *cobra.Command, req *generated.GenerateRequest) error {
	// Load config from file if it exists
	fileConfig, err := adder.LoadConfig(".")
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not load config file: %v\n", err)
		fileConfig = adder.DefaultConfig()
	}

	// Merge command line flags with config file (flags take precedence)
	config := adder.MergeWithFlags(fileConfig, req.Flags.BinaryName, req.Flags.Input, req.Flags.Output, req.Flags.Package, req.Flags.Suffix)

	// Override package strategy if provided
	if req.Flags.PackageStrategy != "" && req.Flags.PackageStrategy != "directory" {
		config.PackageStrategy = req.Flags.PackageStrategy
	}

	// Validate that binary_name is set (either from config or flag)
	if config.BinaryName == "" {
		return fmt.Errorf("binary_name is required. Set it in .adder.yaml or use --binary-name flag")
	}

	if req.Flags.Validate {
		fmt.Printf("🔍 Validating documentation in %s...\n", config.InputDir)
	} else {
		fmt.Printf("🐍 Generating CLI commands from %s to %s...\n", config.InputDir, config.OutputDir)
	}

	// Create generator
	generator := adder.New(config)

	// Validate documentation
	fmt.Println("🔍 Validating documentation...")
	if err := generator.Validate(); err != nil {
		if req.Flags.Validate {
			return fmt.Errorf("❌ Validation failed: %w", err)
		}
		fmt.Printf("⚠️  Validation warnings: %v\n", err)
		fmt.Println("Continuing with generation...")
	}

	// If validate-only, stop here after validation
	if req.Flags.Validate {
		// Parse commands to show statistics without generating
		commands, err := generator.ParseCommands()
		if err != nil {
			return fmt.Errorf("❌ Parsing failed: %w", err)
		}

		// Calculate stats manually since we didn't generate
		stats := map[string]int{
			"total_commands":  len(commands),
			"total_arguments": 0,
			"total_flags":     0,
		}
		for _, cmd := range commands {
			stats["total_arguments"] += len(cmd.Arguments)
			stats["total_flags"] += len(cmd.Flags)
		}

		fmt.Println("✅ Validation completed!")
		fmt.Printf("📊 Found %d commands with %d flags and %d arguments\n",
			stats["total_commands"], stats["total_flags"], stats["total_arguments"])

		fmt.Println("\n📋 Commands found:")
		for _, command := range commands {
			if command.Name != "" {
				fmt.Printf("  - %s: %s\n", command.Name, command.Title)
			}
		}

		fmt.Println("\n💡 Run without --validate to generate code")
		return nil
	}

	// Generate command stubs
	ctx := context.Background()
	opts := adder.GenerateOptions{
		Force: req.Flags.Force,
	}
	if err := generator.GenerateWithOptions(ctx, opts); err != nil {
		return fmt.Errorf("❌ Generation failed: %w", err)
	}

	// Show statistics
	stats := generator.GetStats()
	commands := generator.GetCommands()

	fmt.Println("✅ Code generation completed!")
	if stats["skipped_files"] > 0 {
		fmt.Printf("📊 Generated %d commands with %d flags and %d arguments (%d files skipped, up-to-date)\n",
			stats["total_commands"], stats["total_flags"], stats["total_arguments"], stats["skipped_files"])
	} else {
		fmt.Printf("📊 Generated %d commands with %d flags and %d arguments\n",
			stats["total_commands"], stats["total_flags"], stats["total_arguments"])
	}

	fmt.Println("\n📋 Generated commands:")
	for _, command := range commands {
		if command.Name != "" {
			fmt.Printf("  - %s: %s\n", command.Name, command.Title)
		}
	}

	fmt.Println("\n💡 Next steps:")
	fmt.Println("  1. Implement handler interfaces in your handlers package")
	fmt.Println("  2. Wire commands using the generated constructors")
	fmt.Println("  3. Add commands to your root command")

	return nil
}
