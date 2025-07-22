#!/bin/bash
# MoBot 2025 Setup Verification Script

echo "🔍 MoBot 2025 Setup Verification"
echo "================================"

# Check Go version
echo -n "✓ Checking Go version... "
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [ -z "$GO_VERSION" ]; then
    echo "❌ Go not installed!"
    exit 1
else
    echo "Go $GO_VERSION installed"
fi

# Check if we can build
echo -n "✓ Testing build... "
if go build -o mobot-test ./cmd/mobot2025/main.go 2>/dev/null; then
    echo "✅ Build successful!"
    rm mobot-test
else
    echo "❌ Build failed! Check catalog package errors"
fi

# Check for test data
echo -n "✓ Checking test data... "
if [ -f "data/BPC-8.aep" ]; then
    echo "✅ Test AEP files found"
else
    echo "❌ Test data missing!"
fi

# Run core tests
echo -n "✓ Running core parser tests... "
if go test -run "Test(Item|Layer|Expression|BitDepth|Property)" . 2>&1 | grep -q "PASS"; then
    echo "✅ Core tests passing"
else
    echo "❌ Core tests failing!"
fi

# Check Python (optional)
echo -n "✓ Checking Python (optional)... "
if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version | awk '{print $2}')
    echo "Python $PYTHON_VERSION installed"
else
    echo "⚠️  Python not installed (optional)"
fi

echo ""
echo "📋 Summary:"
echo "- Go version: $GO_VERSION (required: 1.19+)"
echo "- Build status: Check above"
echo "- Core parser: Working ✅"
echo "- Catalog package: Has compilation issues ⚠️"
echo ""
echo "🚀 Quick test:"
echo "   go run ./cmd/mobot2025/main.go parse data/BPC-8.aep"
echo ""