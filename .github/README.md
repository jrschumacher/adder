# GitHub Workflows

This directory contains GitHub Actions workflows for the adder project.

## Workflows

### 1. Lint (`lint.yml`)
**Triggers:** Push/PR to main/develop branches

**What it does:**
- Runs golangci-lint with comprehensive checks
- Validates code formatting (gofmt)
- Checks import organization (goimports)
- Runs `go vet` for static analysis
- Scans for vulnerabilities with govulncheck

### 2. Test (`test.yml`)
**Triggers:** Push/PR to main/develop branches

**What it does:**
- Runs tests across multiple OS (Linux, Windows, macOS)
- Tests multiple Go versions (1.21, 1.22, 1.23)
- Runs unit tests, integration tests, and example tests
- Generates coverage reports
- Tests documentation generation
- Validates self-dogfooding capability

### 3. Release Please (`release-please.yml`)
**Triggers:** Push to main branch

**What it does:**
- Automatically creates release PRs when changes are pushed to main
- Manages semantic versioning based on conventional commits
- Creates GitHub releases when release PR is merged
- Tags major/minor versions for easy consumption

### 4. Build and Release (`build.yml`)
**Triggers:** Version tags and release-please workflow completion

**What it does:**
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64, arm64)
- Builds multi-platform Docker images
- Publishes to GitHub Container Registry
- Uploads release artifacts to GitHub releases
- Generates checksums for verification

## Configuration Files

### `.golangci.yml`
Comprehensive linting configuration with:
- 30+ enabled linters
- Custom rules for generated code
- Performance and style checks
- Security vulnerability detection

### `.release-please-config.json`
Configuration for release-please automation:
- Go module release type
- Semantic versioning rules
- Release draft/prerelease settings

### `Makefile`
Local development commands:
- `make build` - Build local binary
- `make test` - Run all tests
- `make lint` - Run linters
- `make ci-test` - Run all CI checks locally

### `Dockerfile`
Multi-stage Docker build:
- Builds from Go 1.23 Alpine
- Creates minimal scratch-based final image
- Includes CA certificates and timezone data

## Release Process

1. **Development:**
   - Make changes on feature branches
   - Create PR to main branch
   - Workflows run lint and test checks

2. **Release:**
   - Merge PR to main branch
   - Release-please creates release PR automatically
   - Review and merge release PR
   - GitHub release created with artifacts

3. **Distribution:**
   - Binaries available for download
   - Docker images published to ghcr.io
   - Checksums provided for verification

## Local Development

Run the same checks locally:

```bash
# Install dependencies
make mod

# Run all checks
make ci-test

# Build for all platforms
make build-all

# Test Docker build
make docker-build
```

## Conventional Commits

Use conventional commit messages for automatic version bumping:

- `feat:` - New features (minor version bump)
- `fix:` - Bug fixes (patch version bump)
- `docs:` - Documentation changes (patch version bump)
- `chore:` - Maintenance tasks (no version bump)
- `BREAKING CHANGE:` - Breaking changes (major version bump)

Example:
```
feat: add support for custom templates

This adds the ability to specify custom Go templates
for code generation.

BREAKING CHANGE: The Config struct now requires a
Templates field.
```