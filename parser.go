package adder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// fieldSource tracks where a field name originated from for better error messages
type fieldSource struct {
	fieldType    string // "argument" or "flag"
	originalName string // original name before conversion to PascalCase
}

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

		// Check if this is a root command file
		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		
		// Check if this is the binary's root command file (binary_name.md in root directory)
		if dir == "." && p.config.BinaryName != "" && filename == p.config.BinaryName+".md" {
			cmd, err := p.ParseFile(fsys, path)
			if err != nil {
				return fmt.Errorf("parsing binary root command %s: %w", path, err)
			}
			if cmd != nil {
				cmd.IsRootCommand = true
				cmd.CommandPath = "" // Root command has no path prefix
				commands = append(commands, cmd)
			}
			return nil
		}
		
		// For files in subdirectories, check if it's an index file
		if dir != "." {
			dirName := filepath.Base(dir)
			if p.config.IsIndexFile(filename, dirName) {
				// This is an index file - process it but mark it as such
				cmd, err := p.ParseFile(fsys, path)
				if err != nil {
					return fmt.Errorf("parsing root command %s: %w", path, err)
				}
				if cmd != nil {
					cmd.IsRootCommand = true
					cmd.CommandPath = dirName // Set the command path for root commands
					commands = append(commands, cmd)
				}
				return nil
			}
		}
		
		// Skip other pattern files that are not root commands
		if filename == "_index.md" || filename == "index.md" {
			// Only skip if not in root directory and not a configured root command pattern
			if dir != "." {
				return nil
			}
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
	// Use yaml.v2 directly to parse with more control
	var rawData map[string]interface{}
	
	// Extract frontmatter manually
	if !strings.HasPrefix(content, "---") {
		return nil, nil // No frontmatter, skip this file
	}
	
	parts := strings.SplitN(content[3:], "---", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}
	
	frontmatterContent := strings.TrimSpace(parts[0])
	bodyContent := strings.TrimSpace(parts[1])
	
	// Parse YAML
	err := yaml.Unmarshal([]byte(frontmatterContent), &rawData)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}
	
	// Extract title
	title := ""
	if t, exists := rawData["title"]; exists {
		if titleStr, ok := t.(string); ok {
			title = titleStr
		}
	}
	
	// Extract command section
	commandData, exists := rawData["command"]
	if !exists {
		return nil, nil // No command section
	}
	
	commandMap, ok := commandData.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("command section must be an object")
	}
	
	// Extract command name
	var name string
	if n, exists := commandMap["name"]; exists {
		if nameStr, ok := n.(string); ok {
			name = nameStr
		}
	}
	
	if name == "" {
		return nil, nil // No command name
	}
	
	// Extract aliases
	var aliases []string
	if a, exists := commandMap["aliases"]; exists {
		if aliasSlice, ok := a.([]interface{}); ok {
			for _, alias := range aliasSlice {
				if aliasStr, ok := alias.(string); ok {
					aliases = append(aliases, aliasStr)
				}
			}
		}
	}
	
	// Extract hidden
	var hidden bool
	if h, exists := commandMap["hidden"]; exists {
		if hiddenBool, ok := h.(bool); ok {
			hidden = hiddenBool
		}
	}
	
	// Parse arguments using our custom logic
	var arguments []Argument
	if rawArgs, exists := commandMap["arguments"]; exists {
		args, err := p.parseArguments(rawArgs, filePath)
		if err != nil {
			return nil, fmt.Errorf("parsing arguments: %w", err)
		}
		arguments = args
	}
	
	// Parse flags normally since they don't have the mixed format issue
	var flags []Flag
	if f, exists := commandMap["flags"]; exists {
		if flagsData, ok := f.([]interface{}); ok {
			for i, flagData := range flagsData {
				if flagMap, ok := flagData.(map[interface{}]interface{}); ok {
					flag := Flag{Type: TypeString} // Default type
					
					if name, exists := flagMap["name"]; exists {
						if nameStr, ok := name.(string); ok {
							flag.Name = nameStr
						} else {
							return nil, fmt.Errorf("file %s: flag %d: name must be a string", filePath, i)
						}
					} else {
						return nil, fmt.Errorf("file %s: flag %d: name is required", filePath, i)
					}
					
					if shorthand, exists := flagMap["shorthand"]; exists {
						if shorthandStr, ok := shorthand.(string); ok {
							flag.Shorthand = shorthandStr
						}
					}
					
					if desc, exists := flagMap["description"]; exists {
						if descStr, ok := desc.(string); ok {
							flag.Description = descStr
						}
					}
					
					if typ, exists := flagMap["type"]; exists {
						if typStr, ok := typ.(string); ok {
							flag.Type = typStr
						}
					}
					
					if def, exists := flagMap["default"]; exists {
						flag.Default = def
					}
					
					if req, exists := flagMap["required"]; exists {
						if reqBool, ok := req.(bool); ok {
							flag.Required = reqBool
						}
					}
					
					if enum, exists := flagMap["enum"]; exists {
						if enumSlice, ok := enum.([]interface{}); ok {
							for _, e := range enumSlice {
								if eStr, ok := e.(string); ok {
									flag.Enum = append(flag.Enum, eStr)
								} else {
									// Non-string enum values should cause an error later in validation
									// For now, convert to string to preserve the original value for error reporting
									flag.Enum = append(flag.Enum, fmt.Sprintf("%v", e))
								}
							}
						}
					}
					
					flags = append(flags, flag)
				}
			}
		}
	}

	// Parse persistent flags
	var persistentFlags []Flag
	if f, exists := commandMap["persistent_flags"]; exists {
		if flagsData, ok := f.([]interface{}); ok {
			for i, flagData := range flagsData {
				if flagMap, ok := flagData.(map[interface{}]interface{}); ok {
					flag := Flag{Type: TypeString} // Default type
					
					if name, exists := flagMap["name"]; exists {
						if nameStr, ok := name.(string); ok {
							flag.Name = nameStr
						} else {
							return nil, fmt.Errorf("file %s: persistent_flag %d: name must be a string", filePath, i)
						}
					} else {
						return nil, fmt.Errorf("file %s: persistent_flag %d: name is required", filePath, i)
					}
					
					if shorthand, exists := flagMap["shorthand"]; exists {
						if shorthandStr, ok := shorthand.(string); ok {
							flag.Shorthand = shorthandStr
						}
					}
					
					if desc, exists := flagMap["description"]; exists {
						if descStr, ok := desc.(string); ok {
							flag.Description = descStr
						}
					}
					
					if typ, exists := flagMap["type"]; exists {
						if typStr, ok := typ.(string); ok {
							flag.Type = typStr
						}
					}
					
					if def, exists := flagMap["default"]; exists {
						flag.Default = def
					}
					
					if req, exists := flagMap["required"]; exists {
						if reqBool, ok := req.(bool); ok {
							flag.Required = reqBool
						}
					}
					
					if enum, exists := flagMap["enum"]; exists {
						if enumSlice, ok := enum.([]interface{}); ok {
							for _, e := range enumSlice {
								if eStr, ok := e.(string); ok {
									flag.Enum = append(flag.Enum, eStr)
								} else {
									// Non-string enum values should cause an error later in validation
									// For now, convert to string to preserve the original value for error reporting
									flag.Enum = append(flag.Enum, fmt.Sprintf("%v", e))
								}
							}
						}
					}
					
					persistentFlags = append(persistentFlags, flag)
				}
			}
		}
	}

	cmd := &Command{
		Title:           title,
		Name:            name,
		Aliases:         aliases,
		Hidden:          hidden,
		Arguments:       arguments,
		Flags:           flags,
		PersistentFlags: persistentFlags,
		Description:     bodyContent,
		FilePath:        filePath,
	}

	// Validate command
	if err := p.validateCommand(cmd); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return cmd, nil
}

// parseArguments handles both string array and object array formats for arguments
func (p *Parser) parseArguments(rawArgs interface{}, filePath string) ([]Argument, error) {
	switch args := rawArgs.(type) {
	case []interface{}:
		// Check if first element is a string or object to determine format
		if len(args) == 0 {
			return []Argument{}, nil
		}

		// Check the type of the first element
		switch args[0].(type) {
		case string:
			// Array of strings format: ["file", "input"]
			var arguments []Argument
			for i, arg := range args {
				argStr, ok := arg.(string)
				if !ok {
					return nil, fmt.Errorf("file %s: argument %d: mixed array formats not supported (expected all strings)", filePath, i)
				}
				arguments = append(arguments, Argument{
					Name:        argStr,
					Description: "",
					Required:    true, // Default to required for string-only format
					Type:        TypeString, // Default type
				})
			}
			return arguments, nil

		case map[interface{}]interface{}:
			// Array of objects format: [{name: "file", type: "string", ...}]
			var arguments []Argument
			for i, arg := range args {
				argMap, ok := arg.(map[interface{}]interface{})
				if !ok {
					return nil, fmt.Errorf("file %s: argument %d: mixed array formats not supported (expected all objects)", filePath, i)
				}

				// Convert map to Argument struct
				argument := Argument{
					Type: TypeString, // Default type
				}

				if name, exists := argMap["name"]; exists {
					if nameStr, ok := name.(string); ok {
						argument.Name = nameStr
					} else {
						return nil, fmt.Errorf("file %s: argument %d: name must be a string", filePath, i)
					}
				} else {
					return nil, fmt.Errorf("file %s: argument %d: name is required", filePath, i)
				}

				if desc, exists := argMap["description"]; exists {
					if descStr, ok := desc.(string); ok {
						argument.Description = descStr
					}
				}

				if req, exists := argMap["required"]; exists {
					if reqBool, ok := req.(bool); ok {
						argument.Required = reqBool
					}
				}

				if typ, exists := argMap["type"]; exists {
					if typStr, ok := typ.(string); ok {
						argument.Type = typStr
					}
				}

				arguments = append(arguments, argument)
			}
			return arguments, nil

		default:
			return nil, fmt.Errorf("file %s: unsupported argument format - expected string or object", filePath)
		}

	default:
		return nil, fmt.Errorf("file %s: arguments must be an array", filePath)
	}
}

// validateCommand validates a parsed command
func (p *Parser) validateCommand(cmd *Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("file %s: command name is required", cmd.FilePath)
	}

	if cmd.Title == "" {
		return fmt.Errorf("file %s: command title is required", cmd.FilePath)
	}

	// Track field names to prevent duplicates with their source information
	fieldNames := make(map[string]fieldSource)

	// Validate arguments
	for i, arg := range cmd.Arguments {
		if arg.Name == "" {
			return fmt.Errorf("file %s: argument %d: name is required", cmd.FilePath, i)
		}
		if arg.Type == "" {
			cmd.Arguments[i].Type = TypeString // Default type
		}

		// Check for duplicate field names (after conversion to PascalCase)
		fieldName := pascalCase(arg.Name)
		if existing, exists := fieldNames[fieldName]; exists {
			return fmt.Errorf("file %s: duplicate field name '%s' - argument '%s' conflicts with %s '%s'", 
				cmd.FilePath, fieldName, arg.Name, existing.fieldType, existing.originalName)
		}
		fieldNames[fieldName] = fieldSource{
			fieldType:    "argument",
			originalName: arg.Name,
		}
	}

	// Validate flags and set defaults
	for i, flag := range cmd.Flags {
		if flag.Name == "" {
			return fmt.Errorf("file %s: flag %d: name is required", cmd.FilePath, i)
		}
		if flag.Type == "" {
			cmd.Flags[i].Type = TypeString // Default type
			flag.Type = TypeString         // Update local copy for validation
		}

		// Check for duplicate field names (after conversion to PascalCase)
		fieldName := pascalCase(flag.Name)
		if existing, exists := fieldNames[fieldName]; exists {
			return fmt.Errorf("file %s: duplicate field name '%s' - flag '%s' conflicts with %s '%s'", 
				cmd.FilePath, fieldName, flag.Name, existing.fieldType, existing.originalName)
		}
		fieldNames[fieldName] = fieldSource{
			fieldType:    "flag",
			originalName: flag.Name,
		}

		// Validate enum values
		if len(flag.Enum) > 0 && flag.Type != TypeString {
			return fmt.Errorf("file %s: flag %s: enum is only supported for string flags", cmd.FilePath, flag.Name)
		}
	}

	// Run enhanced validation (type consistency, defaults, etc.)
	return validateCommandConfiguration(cmd, cmd.FilePath)
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
