// Package adder provides a documentation-driven CLI generator for Go applications.
// It processes markdown files with YAML frontmatter to generate type-safe command interfaces,
// request/response structures, and handler interfaces for Cobra CLI applications.
package adder

import (
	"context"
	"fmt"
	"io/fs"
	"os"
)

// Adder is the main interface for the code generator
type Adder struct {
	config    *Config
	generator *Generator
}

// New creates a new Adder instance with the given configuration
func New(config *Config) *Adder {
	if config == nil {
		config = DefaultConfig()
	}

	return &Adder{
		config:    config,
		generator: NewGenerator(config),
	}
}

// NewWithDefaults creates a new Adder instance with default configuration
func NewWithDefaults() *Adder {
	return New(DefaultConfig())
}

// Generate processes the input directory and generates command stubs
func (a *Adder) Generate() error {
	return a.GenerateWithContext(context.Background())
}

// GenerateWithContext processes the input directory and generates command stubs with context
func (a *Adder) GenerateWithContext(ctx context.Context) error {
	// Check if input directory exists
	if _, err := os.Stat(a.config.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory does not exist: %s", a.config.InputDir)
	}

	// Create filesystem from input directory
	inputFS := os.DirFS(a.config.InputDir)

	// Generate code
	if err := a.generator.Generate(ctx, inputFS); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	return nil
}

// GenerateFromFS generates code from the provided filesystem
func (a *Adder) GenerateFromFS(ctx context.Context, inputFS fs.FS) error {
	return a.generator.Generate(ctx, inputFS)
}

// Validate validates the configuration and input files
func (a *Adder) Validate() error {
	// Validate config
	if err := a.validateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Parse and validate commands without generating
	inputFS := os.DirFS(a.config.InputDir)
	commands, err := a.generator.parser.ParseDirectory(inputFS)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	// Validate each command
	for _, cmd := range commands {
		if err := a.generator.parser.validateCommand(cmd); err != nil {
			return fmt.Errorf("command %s validation failed: %w", cmd.Name, err)
		}
	}

	return nil
}

// validateConfig validates the adder configuration
func (a *Adder) validateConfig() error {
	if a.config.InputDir == "" {
		return fmt.Errorf("input_dir is required")
	}

	if a.config.OutputDir == "" {
		return fmt.Errorf("output_dir is required")
	}

	if a.config.Package == "" {
		return fmt.Errorf("package is required")
	}

	if a.config.FileSuffix == "" {
		a.config.FileSuffix = "_generated.go"
	}

	return nil
}

// GetCommands returns all parsed commands (requires Generate to be called first)
func (a *Adder) GetCommands() []*Command {
	return a.generator.ListCommands()
}

// GetCommand returns a specific command by name
func (a *Adder) GetCommand(name string) *Command {
	return a.generator.GetCommand(name)
}

// GetStats returns generation statistics
func (a *Adder) GetStats() map[string]int {
	return a.generator.GetStats()
}

// SetConfig updates the configuration
func (a *Adder) SetConfig(config *Config) {
	a.config = config
	a.generator = NewGenerator(config)
}

// GetConfig returns the current configuration
func (a *Adder) GetConfig() *Config {
	return a.config
}

// Clean removes all generated files
func (a *Adder) Clean() error {
	// This would remove files matching the pattern in the output directory
	// For now, just return nil as it's not implemented
	return nil
}

// Version returns the version of the adder package
func Version() string {
	return "0.1.0"
}
