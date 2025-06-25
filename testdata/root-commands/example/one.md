---
title: Example One
command:
  name: one
  arguments:
    - name: target
      description: Target to process
      required: true
      type: string
  flags:
    - name: force
      shorthand: f
      description: Force the operation
      type: bool
---

# Example One

First example subcommand with arguments and flags.