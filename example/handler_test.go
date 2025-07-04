package main

import (
	"fmt"
	"testing"

	"github.com/jrschumacher/adder/example/generated"
	"github.com/jrschumacher/adder/example/generated/hello"
	"github.com/spf13/cobra"
)

func TestGreetHandler_HandleGreet(t *testing.T) {
	tests := []struct {
		name     string
		req      *hello.GreetRequest
		wantErr  bool
		validate func(t *testing.T)
	}{
		{
			name: "simple greeting",
			req: &hello.GreetRequest{
				Arguments: hello.GreetRequestArguments{
					Name: "World",
				},
				Flags: hello.GreetRequestFlags{
					Capitalize: false,
					Repeat:     1,
					Prefix:     "Hello",
					Format:     "text",
				},
			},
			wantErr: false,
		},
		{
			name: "capitalized greeting",
			req: &hello.GreetRequest{
				Arguments: hello.GreetRequestArguments{
					Name: "Alice",
				},
				Flags: hello.GreetRequestFlags{
					Capitalize: true,
					Repeat:     1,
					Prefix:     "Hello",
					Format:     "text",
				},
			},
			wantErr: false,
		},
		{
			name: "repeated greeting",
			req: &hello.GreetRequest{
				Arguments: hello.GreetRequestArguments{
					Name: "Bob",
				},
				Flags: hello.GreetRequestFlags{
					Capitalize: false,
					Repeat:     3,
					Prefix:     "Hello",
					Format:     "text",
				},
			},
			wantErr: false,
		},
		{
			name: "JSON format greeting",
			req: &hello.GreetRequest{
				Arguments: hello.GreetRequestArguments{
					Name: "Charlie",
				},
				Flags: hello.GreetRequestFlags{
					Capitalize: false,
					Repeat:     1,
					Prefix:     "Hello",
					Format:     "json",
				},
			},
			wantErr: false,
		},
	}

	cmd := &cobra.Command{} // Mock command for testing

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleGreet(cmd, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GreetHandler.HandleGreet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Additional validation if provided
			if tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

func TestGreetHandler_Integration(t *testing.T) {
	// This test demonstrates full integration testing
	// Create the parent hello command
	helloCmd := generated.NewHelloCommand(func(cmd *cobra.Command, req *generated.HelloRequest) error {
		return cmd.Help()
	})

	// Add the greet subcommand
	greetCmd := hello.NewGreetCommand(func(cmd *cobra.Command, req *hello.GreetRequest) error {
		return handleGreet(cmd, req)
	})
	helloCmd.AddCommand(greetCmd)

	// Test various argument combinations
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "basic greet",
			args: []string{"greet", "Alice"},
		},
		{
			name: "greet with capitalize flag",
			args: []string{"greet", "Bob", "--capitalize"},
		},
		{
			name: "greet with repeat flag",
			args: []string{"greet", "Charlie", "--repeat", "2"},
		},
		{
			name: "greet with multiple flags",
			args: []string{"greet", "Diana", "--capitalize", "--repeat", "3", "--format", "json"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set command arguments
			helloCmd.SetArgs(tc.args)

			// Execute command
			err := helloCmd.Execute()
			if err != nil {
				t.Errorf("Command execution failed: %v", err)
			}

			// Reset command for next test
			helloCmd.SetArgs(nil)
		})
	}
}

// Example of testing with dependency injection
type MockGreeter struct {
	calls []string
}

func (m *MockGreeter) Greet(name string, capitalize bool, repeat int) {
	call := fmt.Sprintf("Greet(%s, %t, %d)", name, capitalize, repeat)
	m.calls = append(m.calls, call)
}

func TestGreetHandler_WithMockDependency(t *testing.T) {
	// This example shows how you might test with injected dependencies
	// if your handler had external dependencies

	_ = &MockGreeter{} // Example mock for documentation

	// In a real scenario, you might inject the mock into your handler
	// handler := NewGreetHandlerWithGreeter(mock)

	cmd := &cobra.Command{}

	req := &hello.GreetRequest{
		Arguments: hello.GreetRequestArguments{
			Name: "TestUser",
		},
		Flags: hello.GreetRequestFlags{
			Capitalize: true,
			Repeat:     2,
			Prefix:     "Hello",
			Format:     "text",
		},
	}

	err := handleGreet(cmd, req)
	if err != nil {
		t.Fatalf("HandleGreet failed: %v", err)
	}

	// In a real test with dependency injection, you would verify
	// that the mock was called with expected parameters
	// expectedCall := "Greet(TestUser, true, 2)"
	// if len(mock.calls) != 1 || mock.calls[0] != expectedCall {
	//     t.Errorf("Expected call %s, got %v", expectedCall, mock.calls)
	// }
}
