#!/bin/bash

# Build script for mobot2025
echo "🔨 Building mobot2025 AEP Parser..."

# Create build directory
mkdir -p build

# Build the demo binary
cd demo
go build -o ../build/aep-parser-demo main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "Binary created at: build/aep-parser-demo"
    echo ""
    echo "To run the binary:"
    echo "   ./build/aep-parser-demo"
else
    echo "❌ Build failed!"
    exit 1
fi