// Package main provides the command handlers for the adder CLI application.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jrschumacher/adder"
	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/spf13/cobra"
)

// GenerateHandler implements the generated GenerateHandler interface
type GenerateHandler struct{}

// NewGenerateHandler creates a new instance of GenerateHandler.
func NewGenerateHandler() *GenerateHandler {
	return &GenerateHandler{}
}

// HandleGenerate processes the generate command request to create CLI command stubs.
func (h *GenerateHandler) HandleGenerate(_ *cobra.Command, req *generated.GenerateRequest) error {
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

// VersionHandler implements the generated VersionHandler interface
type VersionHandler struct{}

// NewVersionHandler creates a new instance of VersionHandler.
func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

// HandleVersion processes the version command request to display version information.
func (h *VersionHandler) HandleVersion(_ *cobra.Command, _ *generated.VersionRequest) error {
	fmt.Printf("adder version %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("built at: %s\n", date)
	return nil
}

// InitHandler implements the generated InitHandler interface
type InitHandler struct{}

// NewInitHandler creates a new instance of InitHandler
func NewInitHandler() *InitHandler {
	return &InitHandler{}
}

// HandleInit processes the init command to create a configuration file
func (h *InitHandler) HandleInit(_ *cobra.Command, req *generated.InitRequest) error {
	fmt.Println("🚀 Welcome to adder! Let's create a configuration file.")
	fmt.Println()

	// Check if config file already exists
	configPath := ".adder.yaml"
	if _, err := os.Stat(configPath); err == nil && !req.Flags.Force {
		fmt.Printf("❌ Configuration file %s already exists. Use --force to overwrite.\n", configPath)
		return fmt.Errorf("configuration file already exists")
	}

	// Create interactive prompts
	reader := bufio.NewReader(os.Stdin)

	// Binary name (required)
	binaryName := req.Flags.BinaryName
	if binaryName == "" {
		fmt.Print("🔧 Binary name (required, e.g., 'myapp', 'cli-tool'): ")
		input, _ := reader.ReadString('\n')
		binaryName = strings.TrimSpace(input)
		if binaryName == "" {
			return fmt.Errorf("❌ Binary name is required")
		}
	}

	// Input directory
	fmt.Print("📁 Input directory for markdown files (default: docs/commands): ")
	inputDir, _ := reader.ReadString('\n')
	inputDir = strings.TrimSpace(inputDir)
	if inputDir == "" {
		inputDir = "docs/commands"
	}

	// Output directory
	fmt.Print("📦 Output directory for generated files (default: generated): ")
	outputDir, _ := reader.ReadString('\n')
	outputDir = strings.TrimSpace(outputDir)
	if outputDir == "" {
		outputDir = "generated"
	}

	// Package name
	fmt.Print("📝 Go package name (default: generated): ")
	packageName, _ := reader.ReadString('\n')
	packageName = strings.TrimSpace(packageName)
	if packageName == "" {
		packageName = "generated"
	}

	// File suffix
	fmt.Print("🏷️  File suffix for generated files (default: _generated.go): ")
	suffix, _ := reader.ReadString('\n')
	suffix = strings.TrimSpace(suffix)
	if suffix == "" {
		suffix = "_generated.go"
	}

	// Index format
	fmt.Print("📁 Index format - directory (example/example.md), index (example/index.md), or _index (example/_index.md) (default: directory): ")
	indexFormat, _ := reader.ReadString('\n')
	indexFormat = strings.TrimSpace(indexFormat)
	if indexFormat == "" {
		indexFormat = "directory"
	}

	// Create config
	config := &adder.Config{
		BinaryName:          binaryName,
		InputDir:            inputDir,
		OutputDir:           outputDir,
		Package:             packageName,
		GeneratedFileSuffix: suffix,
		IndexFormat:         indexFormat,
		PackageStrategy:     "directory", // Default to directory strategy
	}

	// Save config
	if err := adder.SaveConfig(config, configPath); err != nil {
		return fmt.Errorf("❌ Failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Printf("✅ Configuration saved to %s\n", configPath)
	fmt.Println()
	fmt.Println("📋 Configuration summary:")
	fmt.Printf("  Binary name:   %s\n", binaryName)
	fmt.Printf("  Input:         %s\n", inputDir)
	fmt.Printf("  Output:        %s\n", outputDir)
	fmt.Printf("  Package:       %s\n", packageName)
	fmt.Printf("  Suffix:        %s\n", suffix)
	fmt.Printf("  Index format:  %s\n", indexFormat)
	fmt.Printf("  Package strategy: directory\n")
	fmt.Println()
	fmt.Printf("💡 Next step: Create %s/%s.md for your root command, then run 'adder generate'!\n", inputDir, binaryName)

	return nil
}
