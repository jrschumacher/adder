package adder

import (
	"testing"
)

func TestConfig_GetIndexPatterns(t *testing.T) {
	tests := []struct {
		name             string
		format           string
		dirName          string
		expectedPatterns []string
	}{
		{
			name:             "directory format",
			format:           "directory",
			dirName:          "example",
			expectedPatterns: []string{"example.md", "index.md"},
		},
		{
			name:             "index format",
			format:           "index",
			dirName:          "example",
			expectedPatterns: []string{"index.md"},
		},
		{
			name:             "_index format",
			format:           "_index",
			dirName:          "example",
			expectedPatterns: []string{"_index.md"},
		},
		{
			name:             "hugo alias",
			format:           "hugo",
			dirName:          "example",
			expectedPatterns: []string{"_index.md"},
		},
		{
			name:             "unknown format defaults to directory",
			format:           "unknown",
			dirName:          "example",
			expectedPatterns: []string{"example.md", "index.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{IndexFormat: tt.format}
			patterns := config.GetIndexPatterns(tt.dirName)

			if len(patterns) != len(tt.expectedPatterns) {
				t.Errorf("Expected %d patterns, got %d", len(tt.expectedPatterns), len(patterns))
				return
			}

			for i, expected := range tt.expectedPatterns {
				if patterns[i] != expected {
					t.Errorf("Pattern %d: expected %q, got %q", i, expected, patterns[i])
				}
			}
		})
	}
}

func TestConfig_GetPackageName(t *testing.T) {
	tests := []struct {
		name            string
		strategy        string
		basePackage     string
		filePath        string
		expectedPackage string
	}{
		{
			name:            "single strategy - root file",
			strategy:        "single",
			basePackage:     "generated",
			filePath:        "login.md",
			expectedPackage: "generated",
		},
		{
			name:            "single strategy - nested file",
			strategy:        "single",
			basePackage:     "generated",
			filePath:        "auth/admin/user.md",
			expectedPackage: "generated",
		},
		{
			name:            "directory strategy - root file",
			strategy:        "directory",
			basePackage:     "generated",
			filePath:        "login.md",
			expectedPackage: "generated",
		},
		{
			name:            "directory strategy - single level",
			strategy:        "directory",
			basePackage:     "generated",
			filePath:        "auth/login.md",
			expectedPackage: "auth",
		},
		{
			name:            "directory strategy - multiple levels",
			strategy:        "directory",
			basePackage:     "generated",
			filePath:        "dev/selectors/generate.md",
			expectedPackage: "dev_selectors",
		},
		{
			name:            "directory strategy - with dashes",
			strategy:        "directory",
			basePackage:     "generated",
			filePath:        "auth-service/admin-users/create.md",
			expectedPackage: "auth_service_admin_users",
		},
		{
			name:            "path strategy - root file",
			strategy:        "path",
			basePackage:     "generated",
			filePath:        "login.md",
			expectedPackage: "login",
		},
		{
			name:            "path strategy - nested file",
			strategy:        "path",
			basePackage:     "generated",
			filePath:        "auth/login.md",
			expectedPackage: "auth_login",
		},
		{
			name:            "path strategy - deep nesting",
			strategy:        "path",
			basePackage:     "generated",
			filePath:        "dev/selectors/generate.md",
			expectedPackage: "dev_selectors_generate",
		},
		{
			name:            "default strategy fallback",
			strategy:        "unknown",
			basePackage:     "generated",
			filePath:        "auth/login.md",
			expectedPackage: "auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Package:         tt.basePackage,
				PackageStrategy: tt.strategy,
			}

			result := config.GetPackageName(tt.filePath)
			if result != tt.expectedPackage {
				t.Errorf("GetPackageName() = %q, want %q", result, tt.expectedPackage)
			}
		})
	}
}

func TestConfig_IsIndexFile(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		filename string
		dirName  string
		expected bool
	}{
		{
			name:     "directory format - matching file",
			format:   "directory",
			filename: "example.md",
			dirName:  "example",
			expected: true,
		},
		{
			name:     "directory format - index fallback",
			format:   "directory",
			filename: "index.md",
			dirName:  "example",
			expected: true,
		},
		{
			name:     "directory format - non-matching",
			format:   "directory",
			filename: "other.md",
			dirName:  "example",
			expected: false,
		},
		{
			name:     "index format - matching",
			format:   "index",
			filename: "index.md",
			dirName:  "example",
			expected: true,
		},
		{
			name:     "index format - non-matching",
			format:   "index",
			filename: "example.md",
			dirName:  "example",
			expected: false,
		},
		{
			name:     "_index format - matching",
			format:   "_index",
			filename: "_index.md",
			dirName:  "example",
			expected: true,
		},
		{
			name:     "_index format - non-matching",
			format:   "_index",
			filename: "index.md",
			dirName:  "example",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{IndexFormat: tt.format}
			result := config.IsIndexFile(tt.filename, tt.dirName)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}