package adder

import (
	"context"
	"fmt"
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
	InputDir   string `yaml:"input_dir"`
	OutputDir  string `yaml:"output_dir"`
	Package    string `yaml:"package"`
	FileSuffix string `yaml:"file_suffix"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		InputDir:   "docs",
		OutputDir:  "cmd",
		Package:    "cmd",
		FileSuffix: "_generated.go",
	}
}

// Command represents a command definition from markdown
type Command struct {
	Title       string     `yaml:"title"`
	Name        string     `yaml:"name"`
	Aliases     []string   `yaml:"aliases"`
	Hidden      bool       `yaml:"hidden"`
	Arguments   []Argument `yaml:"arguments"`
	Flags       []Flag     `yaml:"flags"`
	Description string     // Markdown content
	FilePath    string     // Source file path
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
	Name       string
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