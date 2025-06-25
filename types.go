package adder

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// Type constants for common data types
const (
	TypeString      = "string"
	TypeBool        = "bool"
	TypeInt         = "int"
	TypeStringArray = "stringArray"
	NilValue        = "nil"
)

// Config represents the generator configuration
type Config struct {
	BinaryName          string            `yaml:"binary_name"`
	InputDir            string            `yaml:"input"`
	OutputDir           string            `yaml:"output"`
	Package             string            `yaml:"package"`
	GeneratedFileSuffix string            `yaml:"generated_file_suffix"`
	IndexFormat         string            `yaml:"index_format,omitempty"`
	PackageStrategy     string            `yaml:"package_strategy,omitempty"` // "single", "directory", "path"
	Validation          ValidationConfig  `yaml:"validation,omitempty"`
}

// ValidationConfig represents validation-specific settings
type ValidationConfig struct {
	Strict bool `yaml:"strict,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		BinaryName:          "", // Required field, must be set in config
		InputDir:            "docs/commands",
		OutputDir:           "generated",
		Package:             "generated",
		GeneratedFileSuffix: "_generated.go",
		IndexFormat:         "directory", // Use directory name (e.g., example/example.md)
		PackageStrategy:     "directory", // Use directory-based package names to avoid conflicts
	}
}

// Command represents a command definition from markdown
type Command struct {
	Title         string     `yaml:"title"`
	Name          string     `yaml:"name"`
	Aliases       []string   `yaml:"aliases"`
	Hidden        bool       `yaml:"hidden"`
	Arguments     []Argument `yaml:"arguments"`
	Flags         []Flag     `yaml:"flags"`
	Description   string     // Markdown content
	FilePath      string     // Source file path
	IsRootCommand bool       // True if this is a root command for subcommands
	CommandPath   string     // The command path (e.g., "example" for "example" root command)
}

// Argument represents a command argument
type Argument struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"`
}

// Flag represents a command flag
type Flag struct {
	Name        string      `yaml:"name"`
	Shorthand   string      `yaml:"shorthand"`
	Description string      `yaml:"description"`
	Type        string      `yaml:"type"`
	Default     interface{} `yaml:"default"`
	Required    bool        `yaml:"required"`
	Enum        []string    `yaml:"enum"`
}

// GetGoType returns the Go type for the flag
func (f *Flag) GetGoType() string {
	switch f.Type {
	case TypeBool:
		return TypeBool
	case TypeInt:
		return TypeInt
	case TypeString:
		return TypeString
	case TypeStringArray:
		return "[]string"
	default:
		return TypeString
	}
}

// GetCobraFlagMethod returns the cobra flag method name
func (f *Flag) GetCobraFlagMethod() string {
	switch f.Type {
	case TypeBool:
		return "Bool"
	case TypeInt:
		return "Int"
	case TypeString:
		return "String"
	case TypeStringArray:
		return "StringArray"
	default:
		return "String"
	}
}

// GetCobraFlagMethodP returns the cobra flag method name with shorthand
func (f *Flag) GetCobraFlagMethodP() string {
	switch f.Type {
	case TypeBool:
		return "BoolP"
	case TypeInt:
		return "IntP"
	case TypeString:
		return "StringP"
	case TypeStringArray:
		return "StringArrayP"
	default:
		return "StringP"
	}
}

// GetDefaultValue returns the default value as a Go literal
func (f *Flag) GetDefaultValue() string {
	if f.Default == nil {
		switch f.Type {
		case TypeBool:
			return "false"
		case TypeInt:
			return "0"
		case TypeString:
			return `""`
		case TypeStringArray:
			return NilValue
		default:
			return `""`
		}
	}

	switch f.Type {
	case TypeBool:
		return fmt.Sprintf("%v", f.Default)
	case TypeInt:
		return fmt.Sprintf("%v", f.Default)
	case TypeString:
		return fmt.Sprintf(`"%s"`, f.Default)
	case TypeStringArray:
		if arr, ok := f.Default.([]interface{}); ok {
			if len(arr) == 0 {
				return NilValue
			}
			result := "[]string{"
			for i, v := range arr {
				if i > 0 {
					result += ", "
				}
				result += fmt.Sprintf(`"%s"`, v)
			}
			result += "}"
			return result
		}
		return "nil"
	default:
		return fmt.Sprintf(`"%s"`, f.Default)
	}
}

// GetValidationTag returns the validation tag for the field
func (f *Flag) GetValidationTag() string {
	var tags []string

	if f.Required {
		tags = append(tags, "required")
	}

	if len(f.Enum) > 0 {
		enumStr := "oneof="
		for i, val := range f.Enum {
			if i > 0 {
				enumStr += " "
			}
			enumStr += val
		}
		tags = append(tags, enumStr)
	}

	if len(tags) == 0 {
		return ""
	}

	return fmt.Sprintf(`validate:"%s"`, joinStrings(tags, ","))
}

// GetValidationTag returns the validation tag for the argument
func (a *Argument) GetValidationTag() string {
	if a.Required {
		return `validate:"required"`
	}
	return ""
}

// GetGoType returns the Go type for the argument
func (a *Argument) GetGoType() string {
	switch a.Type {
	case TypeBool:
		return TypeBool
	case TypeInt:
		return TypeInt
	case TypeString:
		return TypeString
	case TypeStringArray:
		return "[]string"
	default:
		return TypeString
	}
}

// HandlerInterface represents a handler interface to be generated
type HandlerInterface struct {
	Name    string
	Methods []HandlerMethod
}

// HandlerMethod represents a method in a handler interface
type HandlerMethod struct {
	Name        string
	RequestType string
}

// GeneratorInterface defines the interface for code generators
type GeneratorInterface interface {
	Generate(ctx context.Context) error
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// GetPackageName returns the package name for a command based on the strategy
func (c *Config) GetPackageName(filePath string) string {
	switch c.PackageStrategy {
	case "single":
		// Always use the base package name (old behavior)
		return c.Package
		
	case "directory":
		// Use directory structure for package names
		dir := filepath.Dir(filePath)
		if dir == "." {
			// Root directory - use base package name
			return c.Package
		}
		
		// Convert directory path to valid Go package name
		// e.g., "auth/admin" -> "auth_admin", "dev/selectors" -> "dev_selectors"
		packageName := strings.ReplaceAll(dir, "/", "_")
		packageName = strings.ReplaceAll(packageName, "-", "_")
		
		// Ensure it starts with a letter (Go package name requirements)
		if len(packageName) > 0 && !isLetter(packageName[0]) {
			packageName = "pkg_" + packageName
		}
		
		return packageName
		
	case "path":
		// Use full path including filename for maximum uniqueness
		// e.g., "auth/login.md" -> "auth_login"
		// Special case: Index files should use directory name only
		dir := filepath.Dir(filePath)
		filename := filepath.Base(filePath)
		
		if dir != "." && c.IsIndexFile(filename, filepath.Base(dir)) {
			// This is an index file - use directory name only (like directory strategy)
			packageName := strings.ReplaceAll(dir, "/", "_")
			packageName = strings.ReplaceAll(packageName, "-", "_")
			
			// Ensure it starts with a letter
			if len(packageName) > 0 && !isLetter(packageName[0]) {
				packageName = "pkg_" + packageName
			}
			
			return packageName
		}
		
		// Regular file - use full path including filename
		fullPath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
		packageName := strings.ReplaceAll(fullPath, "/", "_")
		packageName = strings.ReplaceAll(packageName, "-", "_")
		
		// Ensure it starts with a letter
		if len(packageName) > 0 && !isLetter(packageName[0]) {
			packageName = "pkg_" + packageName
		}
		
		return packageName
		
	default:
		// Default to directory strategy
		c.PackageStrategy = "directory"
		return c.GetPackageName(filePath)
	}
}

// isLetter checks if a byte is a letter (for Go package name validation)
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// GetIndexPatterns returns the possible filenames for index files in a directory
func (c *Config) GetIndexPatterns(dirName string) []string {
	switch c.IndexFormat {
	case "index":
		return []string{"index.md"}
	case "_index":
		return []string{"_index.md"}
	case "directory":
		return []string{dirName + ".md", "index.md"} // Try dirName.md first, then index.md
	case "hugo":
		return []string{"_index.md"} // Alias for _index
	default:
		// Default to directory name format
		return []string{dirName + ".md", "index.md"}
	}
}

// IsIndexFile checks if a filename matches the index pattern for a directory
func (c *Config) IsIndexFile(filename, dirName string) bool {
	patterns := c.GetIndexPatterns(dirName)
	for _, pattern := range patterns {
		if filename == pattern {
			return true
		}
	}
	return false
}
