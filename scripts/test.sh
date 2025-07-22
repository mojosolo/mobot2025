#!/bin/bash

# test.sh - Enhanced test runner with options

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Default values
VERBOSE=false
COVERAGE=false
SHORT=false
PATTERN=""
PACKAGE="./..."
TIMEOUT="10m"

# Usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Run tests with various options.

OPTIONS:
    -v, --verbose      Show verbose test output
    -c, --coverage     Generate coverage report
    -s, --short        Run only short tests
    -p, --pattern      Run tests matching pattern
    -k, --package      Test specific package (default: ./...)
    -t, --timeout      Set timeout (default: 10m)
    -h, --help         Show this help message

EXAMPLES:
    $0                     # Run all tests
    $0 -v -c              # Verbose with coverage
    $0 -p TestParser      # Run tests matching TestParser
    $0 -k ./catalog       # Test only catalog package
    $0 -s                 # Run short tests only

EOF
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -s|--short)
            SHORT=true
            shift
            ;;
        -p|--pattern)
            PATTERN="$2"
            shift 2
            ;;
        -k|--package)
            PACKAGE="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# Build test command
TEST_CMD="go test"

# Add timeout
TEST_CMD="$TEST_CMD -timeout=$TIMEOUT"

# Add verbose flag
if [ "$VERBOSE" = true ]; then
    TEST_CMD="$TEST_CMD -v"
fi

# Add coverage
if [ "$COVERAGE" = true ]; then
    TEST_CMD="$TEST_CMD -coverprofile=coverage.out -covermode=atomic"
fi

# Add short flag
if [ "$SHORT" = true ]; then
    TEST_CMD="$TEST_CMD -short"
fi

# Add pattern
if [ -n "$PATTERN" ]; then
    TEST_CMD="$TEST_CMD -run $PATTERN"
fi

# Add race detector
TEST_CMD="$TEST_CMD -race"

# Add package
TEST_CMD="$TEST_CMD $PACKAGE"

# Header
echo -e "${BLUE}🧪 MoBot 2025 Test Runner${NC}"
echo "=========================="
echo ""
echo "Configuration:"
echo "  Package:  $PACKAGE"
echo "  Timeout:  $TIMEOUT"
echo "  Verbose:  $VERBOSE"
echo "  Coverage: $COVERAGE"
echo "  Short:    $SHORT"
[ -n "$PATTERN" ] && echo "  Pattern:  $PATTERN"
echo ""

# Run tests
echo -e "${YELLOW}Running tests...${NC}"
echo "Command: $TEST_CMD"
echo ""

START_TIME=$(date +%s)

if eval $TEST_CMD; then
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    
    echo ""
    echo -e "${GREEN}✅ All tests passed!${NC}"
    echo "Duration: ${DURATION}s"
    
    # Show coverage if enabled
    if [ "$COVERAGE" = true ]; then
        echo ""
        echo -e "${YELLOW}Coverage Summary:${NC}"
        go tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $3}'
        
        echo ""
        echo "To view detailed coverage:"
        echo "  HTML report: go tool cover -html=coverage.out"
        echo "  Function list: go tool cover -func=coverage.out"
    fi
else
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    
    echo ""
    echo -e "${RED}❌ Tests failed!${NC}"
    echo "Duration: ${DURATION}s"
    exit 1
fi

# Check for race conditions
if grep -q "WARNING: DATA RACE" test.log 2>/dev/null; then
    echo ""
    echo -e "${RED}⚠️  Race condition detected!${NC}"
    echo "Review test output for details."
fi