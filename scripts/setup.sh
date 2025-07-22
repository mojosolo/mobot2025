#!/bin/bash

# setup.sh - Initial development environment setup

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "🚀 MoBot 2025 Development Setup"
echo "==============================="

# Check Go version
echo ""
echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go 1.19 or higher from https://go.dev"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓ Go ${GO_VERSION} installed${NC}"

# Check minimum Go version
MIN_VERSION="1.19"
if [ "$(printf '%s\n' "$MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]; then
    echo -e "${RED}❌ Go version ${GO_VERSION} is too old. Need ${MIN_VERSION} or higher${NC}"
    exit 1
fi

# Install dependencies
echo ""
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Install development tools
echo ""
echo "Installing development tools..."

# golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "  Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
else
    echo "  golangci-lint already installed"
fi

# goimports
echo "  Installing goimports..."
go install golang.org/x/tools/cmd/goimports@latest

# dlv debugger
echo "  Installing delve debugger..."
go install github.com/go-delve/delve/cmd/dlv@latest

# godoc
echo "  Installing godoc..."
go install golang.org/x/tools/cmd/godoc@latest

echo -e "${GREEN}✓ Development tools installed${NC}"

# Set up git hooks
echo ""
echo "Setting up git hooks..."
if [ -f scripts/pre-commit ]; then
    cp scripts/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    echo -e "${GREEN}✓ Git hooks installed${NC}"
else
    echo -e "${YELLOW}⚠ No pre-commit hook found${NC}"
fi

# Create directories
echo ""
echo "Creating project directories..."
mkdir -p bin
mkdir -p reports/coverage
mkdir -p test/fixtures/{aep,compositions,layers,assets}
echo -e "${GREEN}✓ Directories created${NC}"

# Build the project
echo ""
echo "Building the project..."
if make build; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Run tests
echo ""
echo "Running tests..."
if go test -short ./...; then
    echo -e "${GREEN}✓ Tests passed${NC}"
else
    echo -e "${YELLOW}⚠ Some tests failed${NC}"
fi

# Final instructions
echo ""
echo "================================"
echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Run 'make help' to see available commands"
echo "  2. Run 'make test' to run all tests"
echo "  3. Run 'make serve' to start the API server"
echo "  4. Read docs/DEVELOPER_GUIDE.md for more info"
echo ""
echo "Happy coding! 🎉"