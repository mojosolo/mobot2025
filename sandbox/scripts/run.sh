#!/bin/bash

# Simple run script for mobot2025 demo
echo "🚀 Running mobot2025 AEP Parser Demo..."
echo ""

go run demo/main.go

echo ""
echo "💡 To parse a different file, run:"
echo "   cd demo && go run main.go"
echo "   Then modify the testFile variable in main.go"