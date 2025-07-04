package adder

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

// CommandSchema defines the complete schema for command documentation
// This struct can be used to generate JSON Schema for YAML validation
type CommandSchema struct {
	// Metadata about the command documentation
	Title       string `json:"title" jsonschema:"title=Command Title,description=Short title for the command,required"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=Detailed description of the command"`
	
	// The main command definition
	Command CommandDefinition `json:"command" jsonschema:"title=Command Definition,description=The command structure and behavior,required"`
}

// CommandDefinition defines a single command with all its properties
type CommandDefinition struct {
	// Core command properties
	Name    string   `json:"name" jsonschema:"title=Command Name,description=The command name and usage pattern (e.g. 'hello [name]'),required"`
	Aliases []string `json:"aliases,omitempty" jsonschema:"title=Command Aliases,description=Alternative names for this command"`
	
	// Display properties
	Short   string `json:"short,omitempty" jsonschema:"title=Short Description,description=Brief description shown in help output"`
	Long    string `json:"long,omitempty" jsonschema:"title=Long Description,description=Detailed description shown in command-specific help"`
	Example string `json:"example,omitempty" jsonschema:"title=Usage Examples,description=Examples of how to use the command"`
	
	// Behavior properties
	Hidden     bool `json:"hidden,omitempty" jsonschema:"title=Hidden Command,description=Hide this command from help output"`
	Deprecated string `json:"deprecated,omitempty" jsonschema:"title=Deprecated Warning,description=Mark command as deprecated with custom message"`
	
	// Argument validation
	Arguments []ArgumentDefinition `json:"arguments,omitempty" jsonschema:"title=Command Arguments,description=Positional arguments for the command"`
	
	// Flag definitions
	Flags           []FlagDefinition `json:"flags,omitempty" jsonschema:"title=Command Flags,description=Command-specific flags"`
	PersistentFlags []FlagDefinition `json:"persistent_flags,omitempty" jsonschema:"title=Persistent Flags,description=Flags inherited by subcommands"`
	
	// Advanced Cobra features
	GroupID               string            `json:"group_id,omitempty" jsonschema:"title=Command Group,description=Group ID for organizing subcommands"`
	SuggestFor            []string          `json:"suggest_for,omitempty" jsonschema:"title=Suggest For,description=Commands this should be suggested for"`
	ValidArgs             []string          `json:"valid_args,omitempty" jsonschema:"title=Valid Arguments,description=Valid non-flag arguments for shell completion"`
	ArgAliases            []string          `json:"arg_aliases,omitempty" jsonschema:"title=Argument Aliases,description=Aliases for valid arguments"`
	Annotations           map[string]string `json:"annotations,omitempty" jsonschema:"title=Annotations,description=Key-value pairs for application-specific metadata"`
	Version               string            `json:"version,omitempty" jsonschema:"title=Version,description=Version string (adds --version flag if set)"`
	
	// Completion and behavior flags
	TraverseChildren            bool `json:"traverse_children,omitempty" jsonschema:"title=Traverse Children,description=Parse flags on all parents before executing child"`
	DisableFlagParsing          bool `json:"disable_flag_parsing,omitempty" jsonschema:"title=Disable Flag Parsing,description=Pass all flags as arguments"`
	DisableAutoGenTag           bool `json:"disable_auto_gen_tag,omitempty" jsonschema:"title=Disable Auto Gen Tag,description=Disable auto-generated documentation tags"`
	DisableFlagsInUseLine       bool `json:"disable_flags_in_use_line,omitempty" jsonschema:"title=Disable Flags in Use Line,description=Don't show [flags] in usage line"`
	DisableSuggestions          bool `json:"disable_suggestions,omitempty" jsonschema:"title=Disable Suggestions,description=Disable command suggestions for typos"`
	SuggestionsMinimumDistance  int  `json:"suggestions_minimum_distance,omitempty" jsonschema:"title=Suggestions Minimum Distance,description=Minimum edit distance for suggestions,minimum=1"`
	
	// Silence options
	SilenceErrors bool `json:"silence_errors,omitempty" jsonschema:"title=Silence Errors,description=Quiet errors downstream"`
	SilenceUsage  bool `json:"silence_usage,omitempty" jsonschema:"title=Silence Usage,description=Don't show usage when an error occurs"`
}

// ArgumentDefinition defines a positional argument
type ArgumentDefinition struct {
	Name        string `json:"name" jsonschema:"title=Argument Name,description=Name of the argument,required"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=Description of the argument"`
	Required    bool   `json:"required,omitempty" jsonschema:"title=Required,description=Whether this argument is required"`
	Type        string `json:"type,omitempty" jsonschema:"title=Argument Type,description=Type of the argument,enum=string;int;bool,default=string"`
}

// FlagDefinition defines a command flag with comprehensive Cobra support
type FlagDefinition struct {
	// Core flag properties
	Name        string      `json:"name" jsonschema:"title=Flag Name,description=Name of the flag (without dashes),required"`
	Shorthand   string      `json:"shorthand,omitempty" jsonschema:"title=Shorthand,description=Single character shorthand"`
	Description string      `json:"description,omitempty" jsonschema:"title=Description,description=Description of the flag"`
	Type        string      `json:"type,omitempty" jsonschema:"title=Flag Type,description=Type of the flag,enum=string;bool;int;stringArray,default=string"`
	Default     interface{} `json:"default,omitempty" jsonschema:"title=Default Value,description=Default value for the flag"`
	Required    bool        `json:"required,omitempty" jsonschema:"title=Required,description=Whether this flag is required"`
	
	// Validation
	Enum []string `json:"enum,omitempty" jsonschema:"title=Enum Values,description=Valid values for string flags"`
	
	// Advanced flag features
	Hidden     bool `json:"hidden,omitempty" jsonschema:"title=Hidden Flag,description=Hide this flag from help output"`
	Deprecated string `json:"deprecated,omitempty" jsonschema:"title=Deprecated Warning,description=Mark flag as deprecated with custom message"`
	
	// Shell completion
	CompletionFunc string   `json:"completion_func,omitempty" jsonschema:"title=Completion Function,description=Custom completion function name"`
	ValidValues    []string `json:"valid_values,omitempty" jsonschema:"title=Valid Values,description=Valid values for shell completion"`
	
	// File/directory completion
	IsFilename   bool     `json:"is_filename,omitempty" jsonschema:"title=Is Filename,description=Flag accepts filenames"`
	IsDirname    bool     `json:"is_dirname,omitempty" jsonschema:"title=Is Directory,description=Flag accepts directory names"`
	FileExts     []string `json:"file_extensions,omitempty" jsonschema:"title=File Extensions,description=Valid file extensions for filename completion"`
	
	// Flag relationships (mutual exclusion, etc.)
	MutuallyExclusive []string `json:"mutually_exclusive,omitempty" jsonschema:"title=Mutually Exclusive,description=Flags that cannot be used together"`
	RequiredTogether  []string `json:"required_together,omitempty" jsonschema:"title=Required Together,description=Flags that must be used together"`
	OneRequired       []string `json:"one_required,omitempty" jsonschema:"title=One Required,description=At least one of these flags must be provided"`
}

// SupportedTypes defines the types supported by adder
var SupportedTypes = struct {
	String      string
	Bool        string
	Int         string
	StringArray string
}{
	String:      "string",
	Bool:        "bool", 
	Int:         "int",
	StringArray: "stringArray",
}

// ValidationRules defines validation rules that can be checked
type ValidationRules struct {
	// Type validation rules
	TypeDefaultConsistency bool `json:"type_default_consistency" jsonschema:"title=Type-Default Consistency,description=Validate that default values match their declared types"`
	EnumStringOnly         bool `json:"enum_string_only" jsonschema:"title=Enum String Only,description=Only allow enum on string type flags"`
	EnumDefaultInList      bool `json:"enum_default_in_list" jsonschema:"title=Enum Default In List,description=Validate that default values are in enum list"`
	
	// Required field validation
	RequiredFields []string `json:"required_fields" jsonschema:"title=Required Fields,description=Fields that must be present"`
	
	// Naming validation
	DuplicateFieldNames bool `json:"duplicate_field_names" jsonschema:"title=No Duplicate Field Names,description=Prevent duplicate argument/flag names"`
	
	// Command validation
	CommandTitleRequired bool `json:"command_title_required" jsonschema:"title=Command Title Required,description=Require command title"`
	CommandNameRequired  bool `json:"command_name_required" jsonschema:"title=Command Name Required,description=Require command name"`
}

// GetDefaultValidationRules returns the default validation rules used by adder
func GetDefaultValidationRules() ValidationRules {
	return ValidationRules{
		TypeDefaultConsistency: true,
		EnumStringOnly:         true,
		EnumDefaultInList:      true,
		RequiredFields:         []string{"name", "title"},
		DuplicateFieldNames:    true,
		CommandTitleRequired:   true,
		CommandNameRequired:    true,
	}
}

// GetSupportedCobraFeatures returns a list of Cobra features supported by this schema
func GetSupportedCobraFeatures() []string {
	return []string{
		// Core command features
		"Use", "Aliases", "Short", "Long", "Example",
		
		// Display and behavior
		"Hidden", "Deprecated", "GroupID",
		
		// Advanced features  
		"SuggestFor", "ValidArgs", "ArgAliases", "Annotations", "Version",
		
		// Behavior flags
		"TraverseChildren", "DisableFlagParsing", "DisableAutoGenTag",
		"DisableFlagsInUseLine", "DisableSuggestions", "SuggestionsMinimumDistance",
		
		// Silence options
		"SilenceErrors", "SilenceUsage",
		
		// Flag features
		"Required flags", "Hidden flags", "Deprecated flags", "Flag shortcuts",
		"Mutual exclusion", "Required together", "One required",
		"File/directory completion", "Custom completion functions",
		
		// Validation
		"Enum validation", "Type validation", "Default value validation",
	}
}

// GetUnsupportedCobraFeatures returns Cobra features not yet implemented
func GetUnsupportedCobraFeatures() []string {
	return []string{
		// Hook functions (can't be defined in YAML)
		"PersistentPreRun", "PreRun", "Run", "PostRun", "PersistentPostRun",
		"PersistentPreRunE", "PreRunE", "RunE", "PostRunE", "PersistentPostRunE",
		
		// Complex completion features
		"ValidArgsFunction", "BashCompletionFunction", "RegisterFlagCompletionFunc",
		
		// Runtime configuration
		"FParseErrWhitelist", "CompletionOptions",
		
		// Dynamic features
		"Context", "SetArgs", "Custom help/usage templates",
	}
}

// GenerateJSONSchema generates a JSON Schema from the CommandSchema struct
func GenerateJSONSchema() ([]byte, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		RequiredFromJSONSchemaTags: true,
	}
	
	schema := reflector.Reflect(&CommandSchema{})
	
	// Add custom schema properties
	schema.Title = "Adder Command Schema"
	schema.Description = "JSON Schema for validating adder command documentation in YAML frontmatter"
	schema.ID = "https://github.com/jrschumacher/adder/schema/command.json"
	
	return json.MarshalIndent(schema, "", "  ")
}

// GenerateJSONSchemaString generates a JSON Schema as a formatted string
func GenerateJSONSchemaString() (string, error) {
	schema, err := GenerateJSONSchema()
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON schema: %w", err)
	}
	return string(schema), nil
}

// ValidateAgainstSchema validates a CommandSchema instance against the generated JSON Schema
func ValidateAgainstSchema(cmd *CommandSchema) error {
	// This is a placeholder for actual validation
	// In practice, you'd use a JSON Schema validator library
	// For now, we'll do basic validation using the struct tags
	
	if cmd.Title == "" {
		return fmt.Errorf("title is required")
	}
	
	if cmd.Command.Name == "" {
		return fmt.Errorf("command.name is required")
	}
	
	// Validate flag types
	for _, flag := range cmd.Command.Flags {
		if flag.Name == "" {
			return fmt.Errorf("flag name is required")
		}
		
		if flag.Type != "" {
			validTypes := []string{"string", "bool", "int", "stringArray"}
			valid := false
			for _, validType := range validTypes {
				if flag.Type == validType {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid flag type '%s' for flag '%s', must be one of: %v", flag.Type, flag.Name, validTypes)
			}
		}
		
		// Validate enum is only used with string types
		if len(flag.Enum) > 0 && flag.Type != "string" {
			return fmt.Errorf("enum validation only supported on string type flags, flag '%s' has type '%s'", flag.Name, flag.Type)
		}
	}
	
	// Validate persistent flags
	for _, flag := range cmd.Command.PersistentFlags {
		if flag.Name == "" {
			return fmt.Errorf("persistent flag name is required")
		}
		
		if flag.Type != "" {
			validTypes := []string{"string", "bool", "int", "stringArray"}
			valid := false
			for _, validType := range validTypes {
				if flag.Type == validType {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid persistent flag type '%s' for flag '%s', must be one of: %v", flag.Type, flag.Name, validTypes)
			}
		}
		
		// Validate enum is only used with string types
		if len(flag.Enum) > 0 && flag.Type != "string" {
			return fmt.Errorf("enum validation only supported on string type persistent flags, flag '%s' has type '%s'", flag.Name, flag.Type)
		}
	}
	
	// Validate arguments
	for _, arg := range cmd.Command.Arguments {
		if arg.Name == "" {
			return fmt.Errorf("argument name is required")
		}
		
		if arg.Type != "" {
			validTypes := []string{"string", "int", "bool"}
			valid := false
			for _, validType := range validTypes {
				if arg.Type == validType {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid argument type '%s' for argument '%s', must be one of: %v", arg.Type, arg.Name, validTypes)
			}
		}
	}
	
	return nil
}