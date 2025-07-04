package adder

import (
	"fmt"
	"strconv"
)

// validateFlagConfiguration validates a flag's configuration for consistency
func validateFlagConfiguration(flag *Flag, filePath string, index int) error {
	// Validate required fields
	if flag.Name == "" {
		return fmt.Errorf("file %s: flag %d: name is required", filePath, index)
	}

	// Validate type if specified
	if flag.Type != "" {
		validTypes := []string{"string", "bool", "int", "stringArray"}
		isValidType := false
		for _, validType := range validTypes {
			if flag.Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return fmt.Errorf("file %s: flag %s: invalid type '%s' (must be one of: string, bool, int, stringArray)", filePath, flag.Name, flag.Type)
		}
	}

	// Validate default value matches type
	if flag.Default != nil {
		if err := validateDefaultValueType(flag.Name, flag.Type, flag.Default, filePath); err != nil {
			return err
		}
	}

	// Validate enum configuration
	if len(flag.Enum) > 0 {
		if flag.Type != "string" && flag.Type != "" { // empty type defaults to string
			return fmt.Errorf("file %s: flag %s: enum validation only supported on string type, got '%s'", filePath, flag.Name, flag.Type)
		}
		
		// Validate enum values are strings
		for i, enumValue := range flag.Enum {
			if enumValue == "" {
				return fmt.Errorf("file %s: flag %s: enum value %d cannot be empty", filePath, flag.Name, i)
			}
		}
		
		// Validate default value is in enum (if both specified)
		if flag.Default != nil {
			defaultStr, ok := flag.Default.(string)
			if !ok {
				return fmt.Errorf("file %s: flag %s: default value must be string when enum is specified", filePath, flag.Name)
			}
			
			isValidDefault := false
			for _, enumValue := range flag.Enum {
				if defaultStr == enumValue {
					isValidDefault = true
					break
				}
			}
			if !isValidDefault {
				return fmt.Errorf("file %s: flag %s: default value '%s' must be one of the enum values: %v", filePath, flag.Name, defaultStr, flag.Enum)
			}
		}
	}

	return nil
}

// validateArgumentConfiguration validates an argument's configuration
func validateArgumentConfiguration(arg *Argument, filePath string, index int) error {
	// Validate required fields
	if arg.Name == "" {
		return fmt.Errorf("file %s: argument %d: name is required", filePath, index)
	}

	// For object-style arguments, validate type is specified or defaults correctly
	if arg.Type == "" {
		arg.Type = "string" // Default type
	}

	// Validate type
	validTypes := []string{"string", "int", "bool"}
	isValidType := false
	for _, validType := range validTypes {
		if arg.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("file %s: argument %s: invalid type '%s' (must be one of: string, int, bool)", filePath, arg.Name, arg.Type)
	}

	return nil
}

// validateDefaultValueType checks if a default value matches its declared type
func validateDefaultValueType(fieldName, fieldType string, defaultValue interface{}, filePath string) error {
	if fieldType == "" {
		fieldType = "string" // Default type
	}

	switch fieldType {
	case "string":
		if _, ok := defaultValue.(string); !ok {
			return fmt.Errorf("file %s: %s: default value must be a string for type 'string', got %T", filePath, fieldName, defaultValue)
		}
	case "bool":
		if _, ok := defaultValue.(bool); !ok {
			return fmt.Errorf("file %s: %s: default value must be a boolean for type 'bool', got %T", filePath, fieldName, defaultValue)
		}
	case "int":
		// Accept both int and float64 (YAML numbers), but validate it's a whole number
		switch v := defaultValue.(type) {
		case int:
			// Valid
		case int64:
			// Valid
		case float64:
			if v != float64(int64(v)) {
				return fmt.Errorf("file %s: %s: default value must be a whole number for type 'int', got %v", filePath, fieldName, v)
			}
		default:
			// Try to parse as string
			if str, ok := defaultValue.(string); ok {
				if _, err := strconv.Atoi(str); err != nil {
					return fmt.Errorf("file %s: %s: default value must be an integer for type 'int', got '%s'", filePath, fieldName, str)
				}
			} else {
				return fmt.Errorf("file %s: %s: default value must be an integer for type 'int', got %T", filePath, fieldName, defaultValue)
			}
		}
	case "stringArray":
		// Must be an array of strings
		if arr, ok := defaultValue.([]interface{}); ok {
			for i, item := range arr {
				if _, ok := item.(string); !ok {
					return fmt.Errorf("file %s: %s: default value array item %d must be a string for type 'stringArray', got %T", filePath, fieldName, i, item)
				}
			}
		} else if _, ok := defaultValue.([]string); ok {
			// Already a string array - valid
		} else {
			return fmt.Errorf("file %s: %s: default value must be an array of strings for type 'stringArray', got %T", filePath, fieldName, defaultValue)
		}
	default:
		return fmt.Errorf("file %s: %s: unsupported type '%s'", filePath, fieldName, fieldType)
	}

	return nil
}

// validateCommandConfiguration validates a command's overall configuration
func validateCommandConfiguration(cmd *Command, filePath string) error {
	// Validate required command fields
	if cmd.Name == "" {
		return fmt.Errorf("file %s: command name is required", filePath)
	}
	
	if cmd.Title == "" {
		return fmt.Errorf("file %s: command title is required", filePath)
	}

	// Validate flags
	for i, flag := range cmd.Flags {
		if err := validateFlagConfiguration(&flag, filePath, i); err != nil {
			return err
		}
	}

	// Validate persistent flags
	for i, flag := range cmd.PersistentFlags {
		if err := validateFlagConfiguration(&flag, filePath, i); err != nil {
			return err
		}
	}

	// Validate arguments
	for i, arg := range cmd.Arguments {
		if err := validateArgumentConfiguration(&arg, filePath, i); err != nil {
			return err
		}
	}

	// Validate no duplicate field names (args + flags)
	fieldNames := make(map[string]bool)
	
	for _, arg := range cmd.Arguments {
		if fieldNames[arg.Name] {
			return fmt.Errorf("file %s: duplicate field name '%s'", filePath, arg.Name)
		}
		fieldNames[arg.Name] = true
	}
	
	for _, flag := range cmd.Flags {
		if fieldNames[flag.Name] {
			return fmt.Errorf("file %s: duplicate field name '%s'", filePath, flag.Name)
		}
		fieldNames[flag.Name] = true
	}
	
	for _, flag := range cmd.PersistentFlags {
		if fieldNames[flag.Name] {
			return fmt.Errorf("file %s: duplicate field name '%s'", filePath, flag.Name)
		}
		fieldNames[flag.Name] = true
	}

	return nil
}