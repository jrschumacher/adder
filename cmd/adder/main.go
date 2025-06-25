package main

import (
	"fmt"
	"os"

	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/spf13/cobra"
)

var (
	// Version information (set by build)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "adder",
		Short: "A documentation-driven CLI generator",
		Long: `Adder generates type-safe CLI commands from markdown documentation.

It processes markdown files with YAML frontmatter to create:
- Type-safe command interfaces
- Request/response structures  
- Handler interfaces
- Argument and flag validation`,
	}

	// Create handlers
	generateHandler := NewGenerateHandler()
	versionHandler := NewVersionHandler()
	initHandler := NewInitHandler()

	// Add generated commands
	rootCmd.AddCommand(generated.NewGenerateCommand(generateHandler))
	rootCmd.AddCommand(generated.NewVersionCommand(versionHandler))
	rootCmd.AddCommand(generated.NewInitCommand(initHandler))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
