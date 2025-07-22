#!/bin/bash

# Test script for mobot2025
echo "🧪 Running mobot2025 tests..."
echo ""

# Run all tests with coverage
go test -v -cover ./...

echo ""
echo "To run specific tests:"
echo "   go test -v -run TestName"
echo ""
echo "To see detailed coverage:"
echo "   go test -coverprofile=coverage.out ./..."
echo "   go tool cover -html=coverage.out"