package adder

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// LoadConfig loads configuration from file and merges with defaults
func LoadConfig(dir string) (*Config, error) {
	config := DefaultConfig()

	// Look for config file in the specified directory
	configPath, err := findConfigFile(dir)
	if err != nil {
		return nil, fmt.Errorf("error finding config file: %w", err)
	}

	// If no config file found, return defaults
	if configPath == "" {
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", configPath, err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", configPath, err)
	}

	// Validate required fields
	if config.BinaryName == "" {
		return nil, fmt.Errorf("binary_name is required in config file %s", configPath)
	}

	return config, nil
}

// findConfigFile looks for .adder.yaml or .adder.yml in the given directory
func findConfigFile(dir string) (string, error) {
	// Check .adder.yaml first
	yamlPath := filepath.Join(dir, ".adder.yaml")
	if _, err := os.Stat(yamlPath); err == nil {
		return yamlPath, nil
	}

	// Check .adder.yml second
	ymlPath := filepath.Join(dir, ".adder.yml")
	if _, err := os.Stat(ymlPath); err == nil {
		return ymlPath, nil
	}

	// No config file found
	return "", nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	// Add helpful comments
	configWithComments := `# Adder configuration file
# See https://github.com/jrschumacher/adder for documentation

` + string(data)

	if err := os.WriteFile(path, []byte(configWithComments), 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// MergeWithFlags merges command line flags into config, with flags taking precedence
func MergeWithFlags(config *Config, binaryName, input, output, pkg, suffix string) *Config {
	merged := &Config{
		BinaryName:          config.BinaryName,
		InputDir:            config.InputDir,
		OutputDir:           config.OutputDir,
		Package:             config.Package,
		GeneratedFileSuffix: config.GeneratedFileSuffix,
		IndexFormat:         config.IndexFormat,
		PackageStrategy:     config.PackageStrategy,
		Validation:          config.Validation,
	}

	// Override with flags if provided
	if binaryName != "" {
		merged.BinaryName = binaryName
	}
	if input != "" && input != "docs/commands" { // Check against default
		merged.InputDir = input
	}
	if output != "" && output != "generated" { // Check against default
		merged.OutputDir = output
	}
	if pkg != "" && pkg != "generated" { // Check against default
		merged.Package = pkg
	}
	if suffix != "" && suffix != "_generated.go" { // Check against default
		merged.GeneratedFileSuffix = suffix
	}

	return merged
}