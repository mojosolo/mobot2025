#!/bin/bash

# coverage.sh - Run tests with coverage and generate reports

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_THRESHOLD=80
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"

echo "🧪 Running tests with coverage..."

# Run tests with coverage
go test -v -race -coverprofile=$COVERAGE_FILE -covermode=atomic ./... 2>&1 | tee test.log

# Check if tests passed
if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All tests passed${NC}"

# Generate coverage report
echo ""
echo "📊 Coverage Report:"
echo "=================="
go tool cover -func=$COVERAGE_FILE | tee coverage-summary.txt

# Extract total coverage percentage
COVERAGE=$(go tool cover -func=$COVERAGE_FILE | grep total | awk '{print $3}' | sed 's/%//')

# Convert to integer for comparison
COVERAGE_INT=$(echo $COVERAGE | awk '{print int($1)}')

echo ""
echo "Total Coverage: ${COVERAGE}%"
echo "Threshold: ${COVERAGE_THRESHOLD}%"

# Check if coverage meets threshold
if [ $COVERAGE_INT -lt $COVERAGE_THRESHOLD ]; then
    echo -e "${RED}❌ Coverage ${COVERAGE}% is below threshold ${COVERAGE_THRESHOLD}%${NC}"
    echo ""
    echo "Areas with low coverage:"
    go tool cover -func=$COVERAGE_FILE | grep -E "^[^[:space:]].*[0-9]+\.[0-9]%$" | awk '$3 < 50.0' | sort -k3 -n
    exit 1
else
    echo -e "${GREEN}✅ Coverage ${COVERAGE}% meets threshold ${COVERAGE_THRESHOLD}%${NC}"
fi

# Generate HTML coverage report
echo ""
echo "📄 Generating HTML coverage report..."
go tool cover -html=$COVERAGE_FILE -o=$COVERAGE_HTML
echo "HTML report generated: $COVERAGE_HTML"

# Generate coverage badge (if running in CI)
if [ -n "$CI" ]; then
    echo ""
    echo "🏷️  Generating coverage badge..."
    
    # Determine badge color based on coverage
    if [ $COVERAGE_INT -ge 90 ]; then
        COLOR="brightgreen"
    elif [ $COVERAGE_INT -ge 80 ]; then
        COLOR="green"
    elif [ $COVERAGE_INT -ge 70 ]; then
        COLOR="yellow"
    elif [ $COVERAGE_INT -ge 60 ]; then
        COLOR="orange"
    else
        COLOR="red"
    fi
    
    # Create badge URL
    BADGE_URL="https://img.shields.io/badge/coverage-${COVERAGE}%25-${COLOR}"
    echo "Badge URL: $BADGE_URL"
    
    # Save badge URL for README update
    echo $BADGE_URL > coverage-badge-url.txt
fi

# Package-specific coverage report
echo ""
echo "📦 Package Coverage:"
echo "==================="
go tool cover -func=$COVERAGE_FILE | grep -E "^github.com/mojosolo/mobot2025/[^/]+[[:space:]]" | sort -k3 -nr

# Find uncovered lines
echo ""
echo "🔍 Top 10 files with lowest coverage:"
echo "===================================="
go tool cover -func=$COVERAGE_FILE | grep -E "\.go[[:space:]]" | sort -k3 -n | head -10

# Summary
echo ""
echo "📈 Coverage Summary:"
echo "==================="
echo "- Total Coverage: ${COVERAGE}%"
echo "- Threshold: ${COVERAGE_THRESHOLD}%"
echo "- Status: $([ $COVERAGE_INT -ge $COVERAGE_THRESHOLD ] && echo -e "${GREEN}PASSED${NC}" || echo -e "${RED}FAILED${NC}")"
echo "- HTML Report: $COVERAGE_HTML"

# Clean up
rm -f test.log

exit 0