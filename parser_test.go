package adder

import (
	"testing"
)

func TestParser_ParseContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filePath string
		want     *Command
		wantErr  bool
	}{
		{
			name: "valid simple command",
			content: `---
title: Test Command
command:
  name: test
---

# Test Command

This is a test command.`,
			filePath: "test.md",
			want: &Command{
				Title:       "Test Command",
				Name:        "test",
				Arguments:   []Argument{},
				Flags:       []Flag{},
				Description: "# Test Command\n\nThis is a test command.",
				FilePath:    "test.md",
			},
			wantErr: false,
		},
		{
			name: "command with arguments and flags",
			content: `---
title: Complex Command
command:
  name: complex [arg]
  arguments:
    - name: arg
      description: An argument
      required: true
      type: string
  flags:
    - name: flag
      shorthand: f
      description: A flag
      default: false
      type: bool
---

# Complex Command

This command has arguments and flags.`,
			filePath: "complex.md",
			want: &Command{
				Title: "Complex Command",
				Name:  "complex [arg]",
				Arguments: []Argument{
					{
						Name:        "arg",
						Description: "An argument",
						Required:    true,
						Type:        "string",
					},
				},
				Flags: []Flag{
					{
						Name:        "flag",
						Shorthand:   "f",
						Description: "A flag",
						Default:     "false",
						Type:        "bool",
					},
				},
				Description: "# Complex Command\n\nThis command has arguments and flags.",
				FilePath:    "complex.md",
			},
			wantErr: false,
		},
		{
			name: "invalid yaml frontmatter",
			content: `---
title: Invalid Command
command:
  name: invalid
  flags:
    - name: flag
      invalid_field: true
---

# Invalid Command`,
			filePath: "invalid.md",
			want: &Command{
				Title:     "Invalid Command",
				Name:      "invalid",
				Arguments: []Argument{},
				Flags: []Flag{
					{
						Name: "flag",
						Type: "string", // Default type is set
					},
				},
				Description: "# Invalid Command",
				FilePath:    "invalid.md",
			},
			wantErr: false, // Parser doesn't validate unknown fields, just ignores them
		},
		{
			name: "no command section",
			content: `---
title: No Command
---

# No Command

This has no command section.`,
			filePath: "nocommand.md",
			want:     nil,
			wantErr:  false, // Returns nil for files without commands
		},
	}

	parser := NewParser(&Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseContent(tt.content, tt.filePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.ParseContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil && got != nil {
				t.Errorf("Parser.ParseContent() = %v, want nil", got)
				return
			}

			if tt.want != nil && got == nil {
				t.Errorf("Parser.ParseContent() = nil, want %v", tt.want)
				return
			}

			if tt.want != nil && got != nil {
				if got.Title != tt.want.Title {
					t.Errorf("Parser.ParseContent() Title = %v, want %v", got.Title, tt.want.Title)
				}
				if got.Name != tt.want.Name {
					t.Errorf("Parser.ParseContent() Name = %v, want %v", got.Name, tt.want.Name)
				}
				if len(got.Arguments) != len(tt.want.Arguments) {
					t.Errorf("Parser.ParseContent() Arguments count = %v, want %v", len(got.Arguments), len(tt.want.Arguments))
				}
				if len(got.Flags) != len(tt.want.Flags) {
					t.Errorf("Parser.ParseContent() Flags count = %v, want %v", len(got.Flags), len(tt.want.Flags))
				}
			}
		})
	}
}

func TestParser_cleanCommandName(t *testing.T) {
	parser := NewParser(&Config{})

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple command", "hello", "hello"},
		{"command with argument", "hello [name]", "hello"},
		{"command with multiple args", "deploy [app] [env]", "deploy"},
		{"command with spaces", "hello world [name]", "hello"},
		{"command with brackets in name", "test[bracket]", "testbracket"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.cleanCommandName(tt.input)
			if got != tt.want {
				t.Errorf("cleanCommandName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFlag_GetGoType(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want string
	}{
		{"string flag", Flag{Type: "string"}, "string"},
		{"bool flag", Flag{Type: "bool"}, "bool"},
		{"int flag", Flag{Type: "int"}, "int"},
		{"default type", Flag{Type: ""}, "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.flag.GetGoType()
			if got != tt.want {
				t.Errorf("Flag.GetGoType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFlag_GetValidationTag(t *testing.T) {
	tests := []struct {
		name string
		flag Flag
		want string
	}{
		{
			"required flag",
			Flag{Required: true},
			`validate:"required"`,
		},
		{
			"enum flag",
			Flag{Enum: []string{"small", "big", "huge"}},
			`validate:"oneof=small big huge"`,
		},
		{
			"required enum flag",
			Flag{Required: true, Enum: []string{"a", "b"}},
			`validate:"required,oneof=a b"`,
		},
		{
			"no validation",
			Flag{},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.flag.GetValidationTag()
			if got != tt.want {
				t.Errorf("Flag.GetValidationTag() = %q, want %q", got, tt.want)
			}
		})
	}
}
