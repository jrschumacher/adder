package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jrschumacher/adder/example/generated"
	"github.com/jrschumacher/adder/example/generated/hello"
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

	// Create hello group command
	helloCmd := generated.NewHelloCommand(func(cmd *cobra.Command, req *generated.HelloRequest) error {
		// Show help when no subcommand is provided
		return cmd.Help()
	})
	
	// Add greet subcommand
	greetCmd := hello.NewGreetCommand(func(cmd *cobra.Command, req *hello.GreetRequest) error {
		return handleGreet(cmd, req)
	})
	helloCmd.AddCommand(greetCmd)
	
	// Add debug subcommand
	debugCmd := hello.NewDebugCommand(func(cmd *cobra.Command, req *hello.DebugRequest) error {
		return handleDebug(cmd, req)
	})
	helloCmd.AddCommand(debugCmd)
	
	rootCmd.AddCommand(helloCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// handleGreet implements the business logic for the greet command
func handleGreet(cmd *cobra.Command, req *hello.GreetRequest) error {
	greeting := fmt.Sprintf("%s, %s!", req.Flags.Prefix, req.Arguments.Name)

	if req.Flags.Capitalize {
		greeting = strings.ToUpper(greeting)
	}

	// Apply ASCII art styling
	styledGreeting := applyAsciiArt(greeting, req.Flags.AsciiArt)

	// Print the greeting the specified number of times
	for i := 0; i < req.Flags.Repeat; i++ {
		if i > 0 {
			fmt.Print("\n") // Add spacing between repeats
		}
		
		// Output in requested format
		switch req.Flags.Format {
		case "json":
			fmt.Printf(`{"greeting": %q, "iteration": %d}`, styledGreeting, i+1)
		case "yaml":
			fmt.Printf("greeting: %q\niteration: %d", styledGreeting, i+1)
		default:
			if !req.Flags.Quiet {
				fmt.Printf("Greeting %d: ", i+1)
			}
			fmt.Print(styledGreeting)
		}
		
		if i < req.Flags.Repeat-1 {
			fmt.Print("\n")
		}
	}

	// Additional languages
	if len(req.Flags.Languages) > 0 {
		fmt.Print("\n\nAdditional greetings:\n")
		for _, lang := range req.Flags.Languages {
			greeting := getGreetingInLanguage(req.Arguments.Name, lang)
			fmt.Printf("- %s: %s\n", lang, greeting)
		}
	}

	return nil
}

// handleDebug implements the debug command
func handleDebug(cmd *cobra.Command, req *hello.DebugRequest) error {
	if req.Flags.Trace {
		fmt.Println("🔍 Debug trace enabled")
	}
	
	if req.Flags.DumpConfig {
		fmt.Println("📋 Configuration dump:")
		fmt.Printf("  Trace: %t\n", req.Flags.Trace)
		fmt.Printf("  Test Enum: %s\n", req.Flags.TestEnum)
	}
	
	fmt.Printf("Debug level set to: %s\n", req.Flags.TestEnum)
	return nil
}

// getGreetingInLanguage returns a greeting in the specified language
func getGreetingInLanguage(name, language string) string {
	greetings := map[string]string{
		"spanish": "Hola, " + name + "!",
		"french":  "Bonjour, " + name + "!",
		"german":  "Hallo, " + name + "!",
		"italian": "Ciao, " + name + "!",
		"portuguese": "Olá, " + name + "!",
	}
	
	if greeting, exists := greetings[strings.ToLower(language)]; exists {
		return greeting
	}
	return "Hello, " + name + "! (language '" + language + "' not supported)"
}

// applyAsciiArt applies the specified ASCII art style to the text
func applyAsciiArt(text, style string) string {
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
