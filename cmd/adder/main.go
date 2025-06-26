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

	rootCmd *cobra.Command
)

func main() {
	// Use the generated root command with persistent flags
	rootCmd = generated.NewAdderCommand(adderCmd)
	rootCmd.Long = `Adder generates type-safe CLI commands from markdown documentation.

It processes markdown files with YAML frontmatter to create:
- Type-safe command interfaces
- Request/response structures  
- Handler interfaces
- Argument and flag validation`

	// Add generated commands with handler functions
	rootCmd.AddCommand(generated.NewGenerateCommand(generateCmd))
	rootCmd.AddCommand(generated.NewVersionCommand(versionCmd))
	rootCmd.AddCommand(generated.NewInitCommand(initCmd))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// adderCmd processes the root adder command request
func adderCmd(cmd *cobra.Command, req *generated.AdderRequest) error {
	// For now, just show help when no subcommand is provided
	return cmd.Help()
}
