# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Go variables
GO_VERSION = 1.23
BINARY_NAME = adder
DIST_DIR = dist

# LDFLAGS for version information
LDFLAGS = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the binary
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd

.PHONY: build-all
build-all: clean ## Build binaries for all platforms
	@echo "Building $(BINARY_NAME) $(VERSION) for all platforms..."
	@mkdir -p $(DIST_DIR)
	
	# Linux AMD64
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd
	
	# Linux ARM64
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd
	
	# macOS AMD64
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd
	
	# macOS ARM64
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd
	
	# Windows AMD64
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd
	
	# Windows ARM64
	@GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe ./cmd
	
	@echo "Built binaries:"
	@ls -la $(DIST_DIR)/

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -cover ./...
	@cd example && go test -v -race .

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run --timeout=10m

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

.PHONY: mod
mod: ## Update and tidy go modules
	@echo "Updating go modules..."
	@go mod download
	@go mod tidy

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html

.PHONY: install
install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_NAME) $(shell go env GOPATH)/bin/


.PHONY: generate-example
generate-example: build ## Generate example commands
	@echo "Generating example commands..."
	@./$(BINARY_NAME) generate --input example/docs/man --output example/generated --package generated

.PHONY: generate-self
generate-self: build ## Generate self commands (dogfooding)
	@echo "Generating self commands..."
	@./$(BINARY_NAME) generate --input docs/commands --output cmd/generated --package generated

.PHONY: ci-test
ci-test: mod lint vet test ## Run all CI checks
	@echo "All CI checks passed!"

.PHONY: release-dry-run
release-dry-run: clean build-all test-coverage ## Prepare for release (dry run)
	@echo "Release dry run completed for version $(VERSION)"
	@echo "Binaries:"
	@ls -la $(DIST_DIR)/

.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"