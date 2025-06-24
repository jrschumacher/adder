---
title: Quote Test Command
command:
  name: quote-test
  flags:
    - name: example
      description: 'Use format "key=value" for this flag'
      type: string
    - name: pattern
      description: 'Example: pattern="*.go" or pattern="test_*"'
      type: string
---

# Quote Test Command

This command tests proper quote escaping in generated code.