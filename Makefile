# MoBot 2025 Makefile
# Comprehensive build and development automation

# Variables
BINARY_NAME := mobot
MAIN_PACKAGE := ./cmd/mobot2025/main.go
GO := go
GOFLAGS := -ldflags="-s -w"
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags with version info
BUILD_FLAGS := -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: all build test clean help

# Default target
all: clean test build

## help: Show this help message
help:
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## build: Build the main binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)✓ Build complete: ./$(BINARY_NAME)$(NC)"

## build-all: Build all command binaries
build-all: build
	@echo "$(GREEN)Building all commands...$(NC)"
	@for cmd in cmd/*/; do \
		if [ -f "$$cmd/main.go" ]; then \
			name=$$(basename $$cmd); \
			echo "  Building $$name..."; \
			$(GO) build $(BUILD_FLAGS) -o bin/$$name $$cmd/main.go; \
		fi \
	done
	@echo "$(GREEN)✓ All builds complete$(NC)"

## test: Run all tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GO) test -v -race ./...

## test-short: Run short tests only
test-short:
	@echo "$(YELLOW)Running short tests...$(NC)"
	$(GO) test -v -short ./...

## test-integration: Run integration tests
test-integration:
	@echo "$(YELLOW)Running integration tests...$(NC)"
	$(GO) test -v -tags=integration ./tests/integration/...

## coverage: Run tests with coverage
coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@./scripts/coverage.sh

## coverage-html: Generate HTML coverage report
coverage-html: coverage
	@echo "$(GREEN)Generating HTML coverage report...$(NC)"
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)✓ Coverage report: $(COVERAGE_HTML)$(NC)"

## coverage-report: Generate detailed coverage analysis
coverage-report:
	@./scripts/coverage-report.sh

## lint: Run linters
lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "$(RED)golangci-lint not installed. Run: make install-tools$(NC)"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(YELLOW)Running go vet...$(NC)"
	$(GO) vet ./...

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -rf reports/
	@find . -name "*.test" -delete
	@echo "$(GREEN)✓ Clean complete$(NC)"

## install: Install the binary
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME)...$(NC)"
	@cp $(BINARY_NAME) $(GOPATH)/bin/
	@echo "$(GREEN)✓ Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

## deps: Download dependencies
deps:
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## mod-update: Update all dependencies
mod-update:
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## install-tools: Install development tools
install-tools:
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@echo "  Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
	@echo "  Installing goimports..."
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "  Installing godoc..."
	@$(GO) install golang.org/x/tools/cmd/godoc@latest
	@echo "  Installing dlv (debugger)..."
	@$(GO) install github.com/go-delve/delve/cmd/dlv@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

## serve: Start the API server
serve: build
	@echo "$(GREEN)Starting API server...$(NC)"
	./$(BINARY_NAME) serve

## run: Run the application
run:
	$(GO) run $(MAIN_PACKAGE)

## docker-build: Build Docker image
docker-build:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t mobot2025:$(VERSION) -t mobot2025:latest .

## release: Create a new release
release:
	@echo "$(YELLOW)Creating release $(VERSION)...$(NC)"
	@echo "  Building binaries..."
	@mkdir -p dist
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" -a "$$arch" = "arm64" ]; then continue; fi; \
			echo "    $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch $(GO) build $(BUILD_FLAGS) \
				-o dist/$(BINARY_NAME)-$(VERSION)-$$os-$$arch \
				$(MAIN_PACKAGE); \
		done \
	done
	@echo "$(GREEN)✓ Release binaries created in dist/$(NC)"

## bench: Run benchmarks
bench:
	@echo "$(YELLOW)Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem ./...

## profile-cpu: Run CPU profiling
profile-cpu:
	@echo "$(YELLOW)Running CPU profiling...$(NC)"
	$(GO) test -cpuprofile=cpu.prof -bench=. ./...
	$(GO) tool pprof cpu.prof

## profile-mem: Run memory profiling
profile-mem:
	@echo "$(YELLOW)Running memory profiling...$(NC)"
	$(GO) test -memprofile=mem.prof -bench=. ./...
	$(GO) tool pprof mem.prof

## check: Run all checks (test, lint, vet)
check: test lint vet
	@echo "$(GREEN)✓ All checks passed$(NC)"

## ci: Run CI pipeline locally
ci: clean deps check build
	@echo "$(GREEN)✓ CI pipeline complete$(NC)"

## dev: Start development mode with file watching
dev:
	@echo "$(GREEN)Starting development mode...$(NC)"
	@if command -v watchexec >/dev/null 2>&1; then \
		watchexec -r -e go -- make run; \
	else \
		echo "$(YELLOW)watchexec not installed. Install with: brew install watchexec$(NC)"; \
		make run; \
	fi

## setup: Initial project setup
setup: deps install-tools
	@echo "$(GREEN)Setting up git hooks...$(NC)"
	@cp scripts/pre-commit .git/hooks/pre-commit 2>/dev/null || true
	@chmod +x .git/hooks/pre-commit 2>/dev/null || true
	@echo "$(GREEN)✓ Project setup complete$(NC)"

## info: Show project information
info:
	@echo "$(GREEN)MoBot 2025 Project Information$(NC)"
	@echo "  Version:    $(VERSION)"
	@echo "  Commit:     $(GIT_COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(shell $(GO) version | cut -d' ' -f3)"
	@echo "  Platform:   $(shell $(GO) env GOOS)/$(shell $(GO) env GOARCH)"

# Custom targets for specific workflows

## parse: Parse an AEP file (requires FILE=path/to/file.aep)
parse: build
	@if [ -z "$(FILE)" ]; then \
		echo "$(RED)Error: FILE not specified. Usage: make parse FILE=path/to/file.aep$(NC)"; \
		exit 1; \
	fi
	./$(BINARY_NAME) parse -file $(FILE)

## scan: Scan directory for AEP files (requires DIR=path/to/dir)
scan: build
	@if [ -z "$(DIR)" ]; then \
		echo "$(RED)Error: DIR not specified. Usage: make scan DIR=path/to/directory$(NC)"; \
		exit 1; \
	fi
	./$(BINARY_NAME) scan -dir $(DIR)