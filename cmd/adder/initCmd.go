// Package main provides the command handlers for the adder CLI application.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jrschumacher/adder"
	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/spf13/cobra"
)

// initCmd processes the init command to create a configuration file
func initCmd(_ *cobra.Command, req *generated.InitRequest) error {
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
