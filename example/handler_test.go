package main

import (
	"fmt"
	"testing"

	"github.com/opentdf/adder/example/generated"
	"github.com/spf13/cobra"
)

func TestHelloHandler_HandleHello(t *testing.T) {
	tests := []struct {
		name     string
		req      *generated.HelloRequest
		wantErr  bool
		validate func(t *testing.T)
	}{
		{
			name: "simple greeting",
			req: &generated.HelloRequest{
				Arguments: generated.HelloRequestArguments{
					Name: "World",
				},
				Flags: generated.HelloRequestFlags{
					Capitalize: false,
					Repeat:     1,
				},
			},
			wantErr: false,
		},
		{
			name: "capitalized greeting",
			req: &generated.HelloRequest{
				Arguments: generated.HelloRequestArguments{
					Name: "Alice",
				},
				Flags: generated.HelloRequestFlags{
					Capitalize: true,
					Repeat:     1,
				},
			},
			wantErr: false,
		},
		{
			name: "repeated greeting",
			req: &generated.HelloRequest{
				Arguments: generated.HelloRequestArguments{
					Name: "Bob",
				},
				Flags: generated.HelloRequestFlags{
					Capitalize: false,
					Repeat:     3,
				},
			},
			wantErr: false,
		},
		{
			name: "empty name should use default",
			req: &generated.HelloRequest{
				Arguments: generated.HelloRequestArguments{
					Name: "",
				},
				Flags: generated.HelloRequestFlags{
					Capitalize: false,
					Repeat:     1,
				},
			},
			wantErr: false,
		},
	}

	handler := NewHelloHandler()
	cmd := &cobra.Command{} // Mock command for testing

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.HandleHello(cmd, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("HelloHandler.HandleHello() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Additional validation if provided
			if tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

func TestHelloHandler_Integration(t *testing.T) {
	// This test demonstrates full integration testing
	// Create handler
	handler := NewHelloHandler()

	// Create command using generated interface
	cmd := generated.NewHelloCommand(handler)

	// Test various argument combinations
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "basic hello",
			args: []string{"Alice"},
		},
		{
			name: "hello with capitalize flag",
			args: []string{"Bob", "--capitalize"},
		},
		{
			name: "hello with repeat flag",
			args: []string{"Charlie", "--repeat", "2"},
		},
		{
			name: "hello with both flags",
			args: []string{"Diana", "--capitalize", "--repeat", "3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set command arguments
			cmd.SetArgs(tc.args)

			// Execute command
			err := cmd.Execute()
			if err != nil {
				t.Errorf("Command execution failed: %v", err)
			}

			// Reset command for next test
			cmd.SetArgs(nil)
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

func TestHelloHandler_WithMockDependency(t *testing.T) {
	// This example shows how you might test with injected dependencies
	// if your handler had external dependencies

	_ = &MockGreeter{} // Example mock for documentation

	// In a real scenario, you might inject the mock into your handler
	// handler := NewHelloHandlerWithGreeter(mock)

	handler := NewHelloHandler()
	cmd := &cobra.Command{}

	req := &generated.HelloRequest{
		Arguments: generated.HelloRequestArguments{
			Name: "TestUser",
		},
		Flags: generated.HelloRequestFlags{
			Capitalize: true,
			Repeat:     2,
		},
	}

	err := handler.HandleHello(cmd, req)
	if err != nil {
		t.Fatalf("HandleHello failed: %v", err)
	}

	// In a real test with dependency injection, you would verify
	// that the mock was called with expected parameters
	// expectedCall := "Greet(TestUser, true, 2)"
	// if len(mock.calls) != 1 || mock.calls[0] != expectedCall {
	//     t.Errorf("Expected call %s, got %v", expectedCall, mock.calls)
	// }
}
