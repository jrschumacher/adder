// Code generated by adder. DO NOT EDIT.

package hello

import (
	"github.com/jrschumacher/adder"
	"github.com/spf13/cobra"
)

// DebugRequestFlags represents the flags for the debug command
type DebugRequestFlags struct {
	Trace bool `json:"trace"` // Enable detailed tracing
	DumpConfig bool `json:"dumpConfig"` // Dump current configuration
	TestEnum string `json:"testEnum" validate:"oneof=debug info warn error"` // Test enum validation
}

// DebugRequest represents the parameters for the debug command
type DebugRequest struct {
	Flags DebugRequestFlags `json:"flags"`
	RawArguments []string `json:"raw_arguments"` // Raw command line arguments passed to the command
}

// GetRawArguments implements the adder.Request interface
func (r *DebugRequest) GetRawArguments() []string {
	return r.RawArguments
}

// Ensure DebugRequest implements adder.Request interface at compile time
var _ adder.Request = (*DebugRequest)(nil)

// DebugHandler defines the function type for handling debug commands
type DebugHandler func(cmd *cobra.Command, req *DebugRequest) error

// NewDebugCommand creates a new debug command with the provided handler function
func NewDebugCommand(handler DebugHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debug",
		Short:   "Debug greeting functionality (hidden)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebug(cmd, args, handler)
		},
	}

	// Register persistent flags

	// Register flags
	cmd.Flags().Bool("trace", false, "Enable detailed tracing")
	cmd.Flags().Bool("dump-config", false, "Dump current configuration")
	cmd.Flags().String("test-enum", "info", "Test enum validation")

	return cmd
}

// runDebug handles argument and flag extraction
func runDebug(cmd *cobra.Command, args []string, handler DebugHandler) error {
	trace, _ := cmd.Flags().GetBool("trace")
	dumpConfig, _ := cmd.Flags().GetBool("dump-config")
	testEnum, _ := cmd.Flags().GetString("test-enum")
	// Validate enum for test-enum
	if err := adder.ValidateEnum("test-enum", testEnum, []string{"debug", "info", "warn", "error"}); err != nil {
		return err
	}

	// Create request
	req := &DebugRequest{
		Flags: DebugRequestFlags{
			Trace: trace,
			DumpConfig: dumpConfig,
			TestEnum: testEnum,
		},
		RawArguments: args,
	}

	// Call handler
	return handler(cmd, req)
}
