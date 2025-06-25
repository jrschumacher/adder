// Example demonstrating the new function-based handler approach
package example

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Example using the new function-based approach (current)
func ExampleFunctionBasedHandlers() {
	// Before: Multiple steps required
	// handler := NewProfileCreateHandler()
	// createCmd := generated.NewCreateCommand(handler)

	// Now: Simple inline function
	createCmd := NewCreateCommand(func(cmd *cobra.Command, req *CreateRequest) error {
		fmt.Printf("Creating profile %s at %s\n", req.Arguments.Profile, req.Arguments.Endpoint)
		
		// Business logic directly here
		if req.Flags.SetDefault {
			fmt.Println("Setting as default profile")
		}
		
		return nil
	})

	listCmd := NewListCommand(func(cmd *cobra.Command, req *ListRequest) error {
		fmt.Println("Listing all profiles:")
		// List logic here
		return nil
	})

	// Clean, concise, and maintains full type safety
	fmt.Printf("Commands created: %s, %s\n", createCmd.Name(), listCmd.Name())
}

// Mock types for the example (would be generated)
type CreateRequest struct {
	Arguments struct {
		Profile  string
		Endpoint string
	}
	Flags struct {
		SetDefault bool
	}
}

type ListRequest struct{}

type CreateHandlerFunc func(cmd *cobra.Command, req *CreateRequest) error
type ListHandlerFunc func(cmd *cobra.Command, req *ListRequest) error

func NewCreateCommand(handler CreateHandlerFunc) *cobra.Command {
	return &cobra.Command{Use: "create"}
}

func NewListCommand(handler ListHandlerFunc) *cobra.Command {
	return &cobra.Command{Use: "list"}
}