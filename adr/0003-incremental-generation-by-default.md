# Incremental Generation by Default

## Status

Accepted

## Context

The code generator was regenerating all files on every run, regardless of whether source files had changed. This was:
- Inefficient for large projects
- Triggering unnecessary rebuilds in watch mode
- Making it harder to track actual changes in version control

## Decision

Implement incremental generation by default that:
- Compares source file modification time with generated file modification time
- Only regenerates when source is newer than output
- Provides `--force` flag to override and regenerate all files
- Reports number of skipped files in output

## Consequences

### Positive
- Faster generation for large projects
- Better integration with build tools and watch mode
- Clearer git diffs (only actual changes)
- Reduced CPU usage and disk I/O

### Negative
- Potential edge cases where timestamp comparison isn't sufficient
- Users might be confused why some files aren't regenerating
- Need to use `--force` when changing templates or generator logic

### Implementation Details

The generator checks modification times by:
1. For each output file, finding all source markdown files that contribute to it
2. Comparing newest source modification time with output modification time
3. Skipping generation if output is newer than all sources
4. Counting and reporting skipped files

```bash
# Normal run - only regenerates changed files
adder generate --input docs --output generated

# Force regeneration of all files
adder generate --input docs --output generated --force
```