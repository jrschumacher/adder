# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the adder project.

## What are ADRs?

ADRs are short documents that capture important architectural decisions made during the development of a project. They help future developers (including ourselves) understand why certain choices were made.

## ADR Format

We use a simplified version of the [MADR](https://adr.github.io/madr/) template:

- **Title**: Short descriptive title
- **Status**: Accepted/Rejected/Deprecated/Superseded
- **Context**: What prompted this decision?
- **Decision**: What did we decide?
- **Consequences**: What are the positive and negative outcomes?

## Current ADRs

1. [Use External Tools for Watch Mode](0001-use-external-tools-for-watch-mode.md) - Why we don't build watch mode into adder
2. [Use Git for Backup and Rollback](0002-use-git-for-backup-and-rollback.md) - Why we rely on version control instead of custom backup
3. [Incremental Generation by Default](0003-incremental-generation-by-default.md) - How we optimize regeneration performance
4. [Validate Flag Instead of Dry Run](0004-validate-flag-instead-of-dry-run.md) - Why we chose --validate over --dry-run
5. [Enhanced Error Messages](0005-enhanced-error-messages.md) - How we improved validation error reporting
6. [Config File and Init Command](0006-config-file-and-init-command.md) - Supporting .adder.yaml/.adder.yml configuration files
7. [Default Input Directory](0007-default-input-directory-docs-commands.md) - Changing default from docs/man to docs/commands