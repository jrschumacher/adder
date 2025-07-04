package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jrschumacher/adder/cmd/adder/generated"
	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// schemaCmd handles the schema command
func schemaCmd(cmd *cobra.Command, req *generated.SchemaRequest) error {
	// Generate the JSON schema
	var schemaData []byte
	var err error
	
	switch req.Flags.Format {
	case "json":
		schemaData, err = generateJSONSchema()
		if err != nil {
			return fmt.Errorf("failed to generate JSON schema: %w", err)
		}
	case "yaml":
		// First generate JSON schema, then convert to YAML
		jsonSchema, err := generateJSONSchema()
		if err != nil {
			return fmt.Errorf("failed to generate JSON schema: %w", err)
		}
		
		// Parse JSON and convert to YAML
		var schemaObj interface{}
		if err := json.Unmarshal(jsonSchema, &schemaObj); err != nil {
			return fmt.Errorf("failed to parse JSON schema for YAML conversion: %w", err)
		}
		
		schemaData, err = yaml.Marshal(schemaObj)
		if err != nil {
			return fmt.Errorf("failed to convert schema to YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", req.Flags.Format)
	}
	
	// Output to file or stdout
	if req.Flags.Output != "" {
		err := os.WriteFile(req.Flags.Output, schemaData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write schema to file %s: %w", req.Flags.Output, err)
		}
		fmt.Printf("✅ Schema written to %s\n", req.Flags.Output)
	} else {
		fmt.Print(string(schemaData))
	}
	
	return nil
}

// CommandSchema defines the complete schema for command documentation
type CommandSchema struct {
	Title       string `json:"title" jsonschema:"title=Command Title,description=Short title for the command,required"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=Detailed description of the command"`
	Command     CommandDefinition `json:"command" jsonschema:"title=Command Definition,description=The command structure and behavior,required"`
}

// CommandDefinition defines a single command with all its properties
type CommandDefinition struct {
	Name    string   `json:"name" jsonschema:"title=Command Name,description=The command name and usage pattern,required"`
	Aliases []string `json:"aliases,omitempty" jsonschema:"title=Command Aliases,description=Alternative names for this command"`
	Short   string   `json:"short,omitempty" jsonschema:"title=Short Description,description=Brief description shown in help output"`
	Long    string   `json:"long,omitempty" jsonschema:"title=Long Description,description=Detailed description shown in command-specific help"`
	Example string   `json:"example,omitempty" jsonschema:"title=Usage Examples,description=Examples of how to use the command"`
	Hidden     bool `json:"hidden,omitempty" jsonschema:"title=Hidden Command,description=Hide this command from help output"`
	Deprecated string `json:"deprecated,omitempty" jsonschema:"title=Deprecated Warning,description=Mark command as deprecated with custom message"`
	Arguments []ArgumentDefinition `json:"arguments,omitempty" jsonschema:"title=Command Arguments,description=Positional arguments for the command"`
	Flags           []FlagDefinition `json:"flags,omitempty" jsonschema:"title=Command Flags,description=Command-specific flags"`
	PersistentFlags []FlagDefinition `json:"persistent_flags,omitempty" jsonschema:"title=Persistent Flags,description=Flags inherited by subcommands"`
}

// ArgumentDefinition defines a positional argument
type ArgumentDefinition struct {
	Name        string `json:"name" jsonschema:"title=Argument Name,description=Name of the argument,required"`
	Description string `json:"description,omitempty" jsonschema:"title=Description,description=Description of the argument"`
	Required    bool   `json:"required,omitempty" jsonschema:"title=Required,description=Whether this argument is required"`
	Type        string `json:"type,omitempty" jsonschema:"title=Argument Type,description=Type of the argument,enum=string;int;bool,default=string"`
}

// FlagDefinition defines a command flag
type FlagDefinition struct {
	Name        string      `json:"name" jsonschema:"title=Flag Name,description=Name of the flag,required"`
	Shorthand   string      `json:"shorthand,omitempty" jsonschema:"title=Shorthand,description=Single character shorthand"`
	Description string      `json:"description,omitempty" jsonschema:"title=Description,description=Description of the flag"`
	Type        string      `json:"type,omitempty" jsonschema:"title=Flag Type,description=Type of the flag,enum=string;bool;int;stringArray,default=string"`
	Default     interface{} `json:"default,omitempty" jsonschema:"title=Default Value,description=Default value for the flag"`
	Required    bool        `json:"required,omitempty" jsonschema:"title=Required,description=Whether this flag is required"`
	Enum        []string    `json:"enum,omitempty" jsonschema:"title=Enum Values,description=Valid values for string flags"`
	Hidden      bool        `json:"hidden,omitempty" jsonschema:"title=Hidden Flag,description=Hide this flag from help output"`
	Deprecated  string      `json:"deprecated,omitempty" jsonschema:"title=Deprecated Warning,description=Mark flag as deprecated with custom message"`
}

// generateJSONSchema generates a JSON Schema from the CommandSchema struct
func generateJSONSchema() ([]byte, error) {
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