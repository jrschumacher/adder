package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/adder/example/generated"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "hello-example",
		Short: "A simple hello world CLI built with adder",
		Long: `This is a demonstration of the adder package.

The hello command is generated from markdown documentation
in docs/man/hello.md and demonstrates type-safe CLI generation.`,
	}

	// Create handler and add generated command
	helloHandler := NewHelloHandler()
	helloCmd := generated.NewHelloCommand(helloHandler)
	rootCmd.AddCommand(helloCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// HelloHandler implements the generated HelloHandler interface
type HelloHandler struct{}

// NewHelloHandler creates a new HelloHandler
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// HandleHello implements the business logic for the hello command
func (h *HelloHandler) HandleHello(cmd *cobra.Command, req *generated.HelloRequest) error {
	greeting := fmt.Sprintf("Hello, %s!", req.Arguments.Name)

	if req.Flags.Capitalize {
		greeting = strings.ToUpper(greeting)
	}

	// Apply ASCII art styling
	styledGreeting := h.applyAsciiArt(greeting, req.Flags.AsciiArt)

	// Print the greeting the specified number of times
	for i := 0; i < req.Flags.Repeat; i++ {
		if i > 0 {
			fmt.Print("\n") // Add spacing between repeats
		}
		fmt.Print(styledGreeting)
	}

	return nil
}

// applyAsciiArt applies the specified ASCII art style to the text
func (h *HelloHandler) applyAsciiArt(text, style string) string {
	switch style {
	case "small":
		return text

	case "big":
		border := strings.Repeat("═", len(text)+4)
		return fmt.Sprintf(`
╔%s╗
║  %s  ║
╚%s╝`, border, text, border)

	case "banner":
		border := strings.Repeat("*", len(text)+8)
		padding := strings.Repeat(" ", len(text)+8)
		return fmt.Sprintf(`
%s
*%s*
*   %s   *
*%s*
%s`, border, padding, text, padding, border)

	default:
		return text
	}
}