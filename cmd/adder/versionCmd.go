package main

import (
	"fmt"

	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/spf13/cobra"
)

// versionCmd processes the version command request to display version information.
func versionCmd(_ *cobra.Command, _ *generated.VersionRequest) error {
	fmt.Printf("adder version %s\n", version)
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("built at: %s\n", date)
	return nil
}
