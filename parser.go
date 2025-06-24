package adder

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/adrg/frontmatter"
)

// Parser handles parsing markdown files with YAML frontmatter
type Parser struct {
	config *Config
}

// NewParser creates a new parser instance
func NewParser(config *Config) *Parser {
	return &Parser{
		config: config,
	}
}

// ParseDirectory parses all markdown files in the input directory
func (p *Parser) ParseDirectory(fsys fs.FS) ([]*Command, error) {
	var commands []*Command

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Skip index files for now
		if strings.HasSuffix(path, "_index.md") {
			return nil
		}

		cmd, err := p.ParseFile(fsys, path)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		if cmd != nil {
			commands = append(commands, cmd)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return commands, nil
}

// ParseFile parses a single markdown file
func (p *Parser) ParseFile(fsys fs.FS, path string) (*Command, error) {
	content, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return p.ParseContent(string(content), path)
}

// ParseContent parses markdown content with YAML frontmatter
func (p *Parser) ParseContent(content, filePath string) (*Command, error) {
	var matter struct {
		Title   string `yaml:"title"`
		Command struct {
			Name      string     `yaml:"name"`
			Aliases   []string   `yaml:"aliases"`
			Hidden    bool       `yaml:"hidden"`
			Arguments []Argument `yaml:"arguments"`
			Flags     []Flag     `yaml:"flags"`
		} `yaml:"command"`
	}

	body, err := frontmatter.Parse(strings.NewReader(content), &matter)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	// Skip files without command definitions
	if matter.Command.Name == "" {
		return nil, nil
	}

	cmd := &Command{
		Title:       matter.Title,
		Name:        matter.Command.Name,
		Aliases:     matter.Command.Aliases,
		Hidden:      matter.Command.Hidden,
		Arguments:   matter.Command.Arguments,
		Flags:       matter.Command.Flags,
		Description: string(body),
		FilePath:    filePath,
	}

	// Validate command
	if err := p.validateCommand(cmd); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return cmd, nil
}

// validateCommand validates a parsed command
func (p *Parser) validateCommand(cmd *Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("command name is required")
	}

	if cmd.Title == "" {
		return fmt.Errorf("command title is required")
	}

	// Track field names to prevent duplicates
	fieldNames := make(map[string]bool)

	// Validate arguments
	for i, arg := range cmd.Arguments {
		if arg.Name == "" {
			return fmt.Errorf("argument %d: name is required", i)
		}
		if arg.Type == "" {
			cmd.Arguments[i].Type = TypeString // Default type
		}

		// Check for duplicate field names (after conversion to PascalCase)
		fieldName := pascalCase(arg.Name)
		if fieldNames[fieldName] {
			return fmt.Errorf("duplicate field name after conversion: %s (from argument %s)", fieldName, arg.Name)
		}
		fieldNames[fieldName] = true
	}

	// Validate flags and set defaults
	for i, flag := range cmd.Flags {
		if flag.Name == "" {
			return fmt.Errorf("flag %d: name is required", i)
		}
		if flag.Type == "" {
			cmd.Flags[i].Type = TypeString // Default type
			flag.Type = TypeString         // Update local copy for validation
		}

		// Check for duplicate field names (after conversion to PascalCase)
		fieldName := pascalCase(flag.Name)
		if fieldNames[fieldName] {
			return fmt.Errorf("duplicate field name after conversion: %s (from flag %s)", fieldName, flag.Name)
		}
		fieldNames[fieldName] = true

		// Validate enum values
		if len(flag.Enum) > 0 && flag.Type != TypeString {
			return fmt.Errorf("flag %s: enum is only supported for string flags", flag.Name)
		}
	}

	return nil
}

// GetCommandPath returns the command path from the file path
func (p *Parser) GetCommandPath(filePath string) string {
	// Remove file extension
	path := strings.TrimSuffix(filePath, ".md")

	// Convert to command path (replace / with spaces for now)
	return strings.ReplaceAll(path, "/", " ")
}

// GetPackageName returns the Go package name for the command
func (p *Parser) GetPackageName(_ *Command) string {
	return p.config.Package
}

// GetStructName returns the request struct name for the command
func (p *Parser) GetStructName(cmd *Command) string {
	return pascalCase(p.cleanCommandName(cmd.Name)) + "Request"
}

// GetHandlerName returns the handler interface name for the command
func (p *Parser) GetHandlerName(cmd *Command) string {
	return pascalCase(p.cleanCommandName(cmd.Name)) + "Handler"
}

// GetMethodName returns the handler method name for the command
func (p *Parser) GetMethodName(cmd *Command) string {
	return "Handle" + pascalCase(p.cleanCommandName(cmd.Name))
}

// GetFunctionName returns the command constructor function name
func (p *Parser) GetFunctionName(cmd *Command) string {
	return "New" + pascalCase(p.cleanCommandName(cmd.Name)) + "Command"
}

// pascalCase converts a string to PascalCase
func pascalCase(s string) string {
	// Split by common separators
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})

	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}

	return result
}

// cleanCommandName removes arguments from command names (e.g., "hello [name]" -> "hello")
func (p *Parser) cleanCommandName(name string) string {
	// Remove everything from the first space or bracket
	parts := strings.Fields(name)
	if len(parts) > 0 {
		// Also remove any remaining brackets from the first part
		cleaned := strings.TrimSpace(parts[0])
		cleaned = strings.ReplaceAll(cleaned, "[", "")
		cleaned = strings.ReplaceAll(cleaned, "]", "")
		return cleaned
	}
	return name
}

// camelCase converts a string to camelCase
func camelCase(s string) string {
	pascal := pascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}
