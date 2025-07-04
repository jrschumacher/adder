// Code generated by adder. DO NOT EDIT.

package generated

import (
	"github.com/jrschumacher/adder"
	"github.com/spf13/cobra"
)

// GenerateRequestFlags represents the flags for the generate command
type GenerateRequestFlags struct {
	BinaryName      string `json:"binaryName"`      // Name of the binary/CLI (required unless set in config)
	Input           string `json:"input"`           // Input directory containing markdown files
	Output          string `json:"output"`          // Output directory for generated files
	Package         string `json:"pkg"`             // Go package name for generated files
	Suffix          string `json:"suffix"`          // File suffix for generated files
	Validate        bool   `json:"validate"`        // Validate documentation without generating files
	Force           bool   `json:"force"`           // Force regeneration of all files regardless of modification time
	PackageStrategy string `json:"packageStrategy"` // Package naming strategy (single, directory, path)
}

// GenerateRequest represents the parameters for the generate command
type GenerateRequest struct {
	Flags        GenerateRequestFlags `json:"flags"`
	RawArguments []string             `json:"raw_arguments"` // Raw command line arguments passed to the command
}

// GetRawArguments implements the adder.Request interface
func (r *GenerateRequest) GetRawArguments() []string {
	return r.RawArguments
}

// Ensure GenerateRequest implements adder.Request interface at compile time
var _ adder.Request = (*GenerateRequest)(nil)

// GenerateHandler defines the function type for handling generate commands
type GenerateHandler func(cmd *cobra.Command, req *GenerateRequest) error

// NewGenerateCommand creates a new generate command with the provided handler function
func NewGenerateCommand(handler GenerateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate CLI commands from markdown documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd, args, handler)
		},
	}

	// Register persistent flags

	// Register flags
	cmd.Flags().StringP("binary-name", "b", "", "Name of the binary/CLI (required unless set in config)")
	cmd.Flags().StringP("input", "i", "docs/commands", "Input directory containing markdown files")
	cmd.Flags().StringP("output", "o", "generated", "Output directory for generated files")
	cmd.Flags().StringP("package", "p", "generated", "Go package name for generated files")
	cmd.Flags().String("suffix", "_generated.go", "File suffix for generated files")
	cmd.Flags().Bool("validate", false, "Validate documentation without generating files")
	cmd.Flags().BoolP("force", "f", false, "Force regeneration of all files regardless of modification time")
	cmd.Flags().String("package-strategy", "directory", "Package naming strategy (single, directory, path)")

	return cmd
}

// runGenerate handles argument and flag extraction
func runGenerate(cmd *cobra.Command, args []string, handler GenerateHandler) error {
	binaryName, _ := cmd.Flags().GetString("binary-name")
	input, _ := cmd.Flags().GetString("input")
	output, _ := cmd.Flags().GetString("output")
	pkg, _ := cmd.Flags().GetString("package")
	suffix, _ := cmd.Flags().GetString("suffix")
	validate, _ := cmd.Flags().GetBool("validate")
	force, _ := cmd.Flags().GetBool("force")
	packageStrategy, _ := cmd.Flags().GetString("package-strategy")

	// Create request
	req := &GenerateRequest{
		Flags: GenerateRequestFlags{
			BinaryName:      binaryName,
			Input:           input,
			Output:          output,
			Package:         pkg,
			Suffix:          suffix,
			Validate:        validate,
			Force:           force,
			PackageStrategy: packageStrategy,
		},
		RawArguments: args,
	}

	// Call handler
	return handler(cmd, req)
}
