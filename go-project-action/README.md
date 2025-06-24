# Go Project CI/CD Action

A reusable GitHub Action that consolidates common Go project workflows including linting, testing, coverage reporting, and automated releases.

## Features

- **Go Setup**: Automatic Go version detection from `go.mod`
- **Caching**: Built-in caching for Go modules and build artifacts
- **Linting**: golangci-lint with customizable linters
- **Testing**: Comprehensive testing with race detection and coverage
- **Coverage**: Codecov integration and artifact uploads
- **Releases**: Automated semantic versioning with release-please
- **Distribution**: GoReleaser integration for binary releases

## Usage

### Basic Setup (Lint & Test)

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: your-username/go-project-action@v1
        with:
          lint-enabled: 'true'
          test-enabled: 'true'
          test-coverage: 'true'
```

### Full CI/CD with Releases

```yaml
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: write
  pull-requests: write

jobs:
  ci-cd:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: your-username/go-project-action@v1
        with:
          # Linting
          lint-enabled: 'true'
          lint-timeout: '10m'
          lint-disable: 'forbidigo,gosec,mnd,nilnil'
          lint-enable: 'errcheck,govet,staticcheck,unused'
          
          # Testing
          test-enabled: 'true'
          test-race: 'true'
          test-coverage: 'true'
          
          # Coverage
          codecov-enabled: 'true'
          codecov-token: ${{ secrets.CODECOV_TOKEN }}
          
          # Release (only on main branch pushes)
          release-enabled: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}
          release-type: 'go'
          
          # GoReleaser
          goreleaser-enabled: ${{ github.event_name == 'push' && github.ref == 'refs/heads/main' }}
          goreleaser-version: '~> v2'
```

### With Custom Test Commands

```yaml
- uses: your-username/go-project-action@v1
  with:
    test-enabled: 'true'
    custom-test-commands: |
      go test -v ./integration/...
      go test -v -run TestGoldenFiles ./...
      go build -o myapp ./cmd/myapp
      ./myapp --version
```

### Separate Workflows

You can also split concerns into separate workflows:

#### lint.yml
```yaml
name: Lint

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: your-username/go-project-action@v1
        with:
          lint-enabled: 'true'
          test-enabled: 'false'
          release-enabled: 'false'
```

#### test.yml
```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: your-username/go-project-action@v1
        with:
          lint-enabled: 'false'
          test-enabled: 'true'
          test-coverage: 'true'
          codecov-enabled: 'true'
          codecov-token: ${{ secrets.CODECOV_TOKEN }}
          release-enabled: 'false'
```

#### release.yml
```yaml
name: Release

on:
  push:
    branches: [main]

permissions:
  contents: write
  pull-requests: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: your-username/go-project-action@v1
        with:
          lint-enabled: 'false'
          test-enabled: 'false'
          release-enabled: 'true'
          goreleaser-enabled: 'true'
```

## Inputs

| Input | Description | Default |
|-------|-------------|---------|
| `go-version-file` | Path to go.mod file | `go.mod` |
| `lint-enabled` | Enable golangci-lint | `true` |
| `lint-timeout` | Timeout for linting | `10m` |
| `lint-disable` | Linters to disable | `forbidigo,gosec,mnd,nilnil` |
| `lint-enable` | Linters to enable | `errcheck,govet,staticcheck,unused` |
| `test-enabled` | Enable testing | `true` |
| `test-race` | Enable race detection | `true` |
| `test-coverage` | Enable coverage | `true` |
| `codecov-enabled` | Upload to Codecov | `false` |
| `codecov-token` | Codecov token | - |
| `release-enabled` | Enable release-please | `false` |
| `release-type` | Release type | `go` |
| `goreleaser-enabled` | Enable GoReleaser | `false` |
| `goreleaser-version` | GoReleaser version | `~> v2` |
| `custom-test-commands` | Additional test commands | - |
| `cache-enabled` | Enable caching | `true` |

## Outputs

| Output | Description |
|--------|-------------|
| `release-created` | Whether a release was created |
| `release-tag` | The release tag if created |
| `coverage-file` | Path to coverage file |

## Publishing Your Action

1. Create a new repository for your action
2. Copy the `action.yml` file to the root
3. Add this README
4. Create a release with a semantic version tag (e.g., `v1.0.0`)
5. Also create/update major version tags (e.g., `v1`) pointing to the latest release

## Examples

See the [examples](examples/) directory for complete workflow configurations.

## License

MIT