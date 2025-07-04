---
title: Generate JSON Schema for command validation

command:
  name: schema
  short: Generate JSON Schema for YAML validation
  long: Generate a comprehensive JSON Schema that can be used to validate adder command documentation in YAML frontmatter
  flags:
    - name: output
      shorthand: o
      description: Output file path for the schema
      type: string
    - name: format
      shorthand: f
      description: Output format
      type: string
      default: json
      enum:
        - json
        - yaml
---

# Generate JSON Schema

Generate a comprehensive JSON Schema for validating adder command documentation.

This command produces a JSON Schema that defines the complete structure and validation rules
for YAML frontmatter in adder markdown files. The schema includes:

## Supported Features

- **Command Properties**: All major Cobra command properties including hidden, deprecated, grouping
- **Flag Types**: string, bool, int, stringArray with full validation
- **Argument Types**: Positional arguments with type and requirement validation  
- **Advanced Features**: Persistent flags, enum validation, mutual exclusion, completion
- **Validation Rules**: Type consistency, default value validation, required fields

## Usage Examples

```bash
# Output schema to stdout
adder schema

# Save schema to file
adder schema --output command-schema.json

# Generate YAML format schema
adder schema --format yaml --output command-schema.yaml
```

## Integration Examples

### IDE Integration (VS Code)
```json
{
  "yaml.schemas": {
    "./command-schema.json": "docs/commands/*.md"
  }
}
```

### CI/CD Validation
```bash
# Generate schema in CI
adder schema --output /tmp/schema.json

# Validate all command files
find docs/commands -name "*.md" -exec yajsv -s /tmp/schema.json {} \;
```