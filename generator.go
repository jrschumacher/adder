package adder

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

// Generator handles code generation from parsed commands
type Generator struct {
	config   *Config
	parser   *Parser
	commands []*Command
}

// NewGenerator creates a new generator instance
func NewGenerator(config *Config) *Generator {
	return &Generator{
		config: config,
		parser: NewParser(config),
	}
}

// Generate processes the input directory and generates code
func (g *Generator) Generate(_ context.Context, inputFS fs.FS) error {
	// Parse all commands from input directory
	commands, err := g.parser.ParseDirectory(inputFS)
	if err != nil {
		return fmt.Errorf("parsing directory: %w", err)
	}

	g.commands = commands

	// Group commands by output file
	fileGroups := g.groupCommandsByFile()

	// Generate code for each file
	for filename, cmds := range fileGroups {
		if err := g.generateFile(filename, cmds); err != nil {
			return fmt.Errorf("generating %s: %w", filename, err)
		}
	}

	return nil
}

// groupCommandsByFile groups commands by their output file
func (g *Generator) groupCommandsByFile() map[string][]*Command {
	groups := make(map[string][]*Command)

	for _, cmd := range g.commands {
		filename := g.getOutputFilename(cmd)
		groups[filename] = append(groups[filename], cmd)
	}

	return groups
}

// getOutputFilename returns the output filename for a command
func (g *Generator) getOutputFilename(cmd *Command) string {
	// FilePath is already relative to InputDir from filesystem walk
	// e.g., "auth/login.md" -> "adder/auth/login_generated.go"
	
	// Replace .md extension with configured suffix
	dir := filepath.Dir(cmd.FilePath)
	base := filepath.Base(cmd.FilePath)
	nameWithoutExt := base[:len(base)-len(filepath.Ext(base))]
	filename := nameWithoutExt + g.config.FileSuffix
	
	// Handle root directory case
	if dir == "." {
		return filepath.Join(g.config.OutputDir, filename)
	}
	
	return filepath.Join(g.config.OutputDir, dir, filename)
}

// generateFile generates code for a single output file
func (g *Generator) generateFile(filename string, commands []*Command) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Generate file content
	content, err := g.generateFileContent(commands)
	if err != nil {
		return fmt.Errorf("generating content: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// generateFileContent generates the content for a file
func (g *Generator) generateFileContent(commands []*Command) (string, error) {
	var buf bytes.Buffer

	// Check if any command needs fmt (for enum validation)
	needsFmt := false
	for _, cmd := range commands {
		for _, flag := range cmd.Flags {
			if len(flag.Enum) > 0 {
				needsFmt = true
				break
			}
		}
		if needsFmt {
			break
		}
	}

	// Generate package header
	packageData := struct {
		Package  string
		NeedsFmt bool
	}{
		Package:  g.config.Package,
		NeedsFmt: needsFmt,
	}

	tmpl := template.Must(template.New("package").Parse(Templates.Package))
	if err := tmpl.Execute(&buf, packageData); err != nil {
		return "", fmt.Errorf("executing package template: %w", err)
	}

	// Generate each command
	for _, cmd := range commands {
		cmdContent, err := g.generateCommand(cmd)
		if err != nil {
			return "", fmt.Errorf("generating command %s: %w", cmd.Name, err)
		}
		buf.WriteString(cmdContent)
	}

	return buf.String(), nil
}

// generateCommand generates code for a single command
func (g *Generator) generateCommand(cmd *Command) (string, error) {
	// Prepare template data
	data := struct {
		Command      *Command
		StructName   string
		HandlerName  string
		MethodName   string
		FunctionName string
	}{
		Command:      cmd,
		StructName:   g.parser.GetStructName(cmd),
		HandlerName:  g.parser.GetHandlerName(cmd),
		MethodName:   g.parser.GetMethodName(cmd),
		FunctionName: g.parser.GetFunctionName(cmd),
	}

	// Execute template
	tmpl := template.Must(template.New("command").
		Funcs(TemplateFunctions()).
		Parse(Templates.Command))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}


// GenerateFromDirectory is a convenience function to generate from a directory path
func GenerateFromDirectory(_, inputDir string) error {
	// Load config (for now use default)
	config := DefaultConfig()
	if inputDir != "" {
		config.InputDir = inputDir
	}

	// Create filesystem
	inputFS := os.DirFS(config.InputDir)

	// Generate
	generator := NewGenerator(config)
	return generator.Generate(context.Background(), inputFS)
}

// ListCommands returns a list of all parsed commands
func (g *Generator) ListCommands() []*Command {
	return g.commands
}

// GetCommand returns a command by name
func (g *Generator) GetCommand(name string) *Command {
	for _, cmd := range g.commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

// ValidateCommands validates all parsed commands
func (g *Generator) ValidateCommands() error {
	for _, cmd := range g.commands {
		if err := g.parser.validateCommand(cmd); err != nil {
			return fmt.Errorf("command %s: %w", cmd.Name, err)
		}
	}
	return nil
}

// GetStats returns generation statistics
func (g *Generator) GetStats() map[string]int {
	stats := make(map[string]int)
	stats["total_commands"] = len(g.commands)
	
	for _, cmd := range g.commands {
		stats["total_arguments"] += len(cmd.Arguments)
		stats["total_flags"] += len(cmd.Flags)
	}
	
	return stats
}