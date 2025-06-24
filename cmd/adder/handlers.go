// Package main provides the command handlers for the adder CLI application.
package main

import (
	"context"
	"fmt"

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
	fmt.Printf("🐍 Generating CLI commands from %s to %s...\n", req.Flags.Input, req.Flags.Output)

	// Configure the generator
	config := &adder.Config{
		InputDir:   req.Flags.Input,
		OutputDir:  req.Flags.Output,
		Package:    req.Flags.Package,
		FileSuffix: req.Flags.Suffix,
	}

	// Create generator
	generator := adder.New(config)

	// Validate documentation first
	fmt.Println("🔍 Validating documentation...")
	if err := generator.Validate(); err != nil {
		fmt.Printf("⚠️  Validation warnings: %v\n", err)
		fmt.Println("Continuing with generation...")
	}

	// Generate command stubs
	ctx := context.Background()
	if err := generator.GenerateWithContext(ctx); err != nil {
		return fmt.Errorf("❌ Generation failed: %w", err)
	}

	// Show statistics
	stats := generator.GetStats()
	commands := generator.GetCommands()

	fmt.Println("✅ Code generation completed!")
	fmt.Printf("📊 Generated %d commands with %d flags and %d arguments\n",
		stats["total_commands"], stats["total_flags"], stats["total_arguments"])

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
