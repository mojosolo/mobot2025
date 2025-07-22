#!/bin/bash

# Comprehensive test runner for mobot2025
# This script runs all tests in a methodical, human-friendly way

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Print header
echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}  mobot2025 Comprehensive Test Suite  ${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Function to run a test suite
run_test_suite() {
    local suite_name=$1
    local test_command=$2
    
    echo -e "${YELLOW}Running $suite_name...${NC}"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ $suite_name passed${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ $suite_name failed${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""
}

# Function to check prerequisites
check_prerequisites() {
    echo -e "${BLUE}Checking prerequisites...${NC}"
    
    # Check Go installation
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed${NC}"
        exit 1
    fi
    
    # Check Python installation
    if ! command -v python3 &> /dev/null; then
        echo -e "${RED}Error: Python 3 is not installed${NC}"
        exit 1
    fi
    
    # Check test data exists
    if [ ! -d "data" ]; then
        echo -e "${RED}Error: Test data directory not found${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ All prerequisites met${NC}"
    echo ""
}

# Function to setup test environment
setup_test_env() {
    echo -e "${BLUE}Setting up test environment...${NC}"
    
    # Create test directories if needed
    mkdir -p tests/fixtures
    mkdir -p tests/output
    mkdir -p coverage
    
    # Download Go dependencies
    echo "Installing Go dependencies..."
    go mod download
    
    # Install Python test dependencies
    if [ -f "requirements-test.txt" ]; then
        echo "Installing Python test dependencies..."
        pip3 install -r requirements-test.txt --quiet
    fi
    
    echo -e "${GREEN}✓ Test environment ready${NC}"
    echo ""
}

# Main test execution
main() {
    # Start timer
    start_time=$(date +%s)
    
    # Check prerequisites
    check_prerequisites
    
    # Setup environment
    setup_test_env
    
    # 1. Run unit tests
    echo -e "${BLUE}=== Unit Tests ===${NC}"
    echo ""
    
    run_test_suite "Core Parser Tests" "go test -v ./... -tags=unit"
    run_test_suite "Agent Tests" "go test -v ./catalog -run TestAgent"
    run_test_suite "Database Tests" "go test -v ./catalog -run TestDatabase"
    
    # 2. Run integration tests
    echo -e "${BLUE}=== Integration Tests ===${NC}"
    echo ""
    
    if [ -d "tests/integration" ]; then
        run_test_suite "Integration Tests" "go test -v ./tests/integration/..."
    else
        echo -e "${YELLOW}⚠ Integration tests directory not found, skipping${NC}"
        SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
    fi
    
    # 3. Run demo viewer tests
    echo -e "${BLUE}=== Demo Viewer Tests ===${NC}"
    echo ""
    
    if [ -d "tests/demo" ]; then
        run_test_suite "Demo Viewer Tests" "go test -v ./tests/demo/..."
    else
        echo -e "${YELLOW}⚠ Demo tests directory not found, skipping${NC}"
        SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
    fi
    
    # 4. Run Python bridge tests
    echo -e "${BLUE}=== Python Bridge Tests ===${NC}"
    echo ""
    
    if [ -f "tests/python/test_bridge.py" ]; then
        run_test_suite "Python Bridge Tests" "python3 -m pytest tests/python/test_bridge.py -v"
    else
        run_test_suite "Python Bridge Tests" "python3 tests/python/test_bridge.py"
    fi
    
    # 5. Run performance tests (optional)
    if [ "$1" == "--with-performance" ]; then
        echo -e "${BLUE}=== Performance Tests ===${NC}"
        echo ""
        
        run_test_suite "Benchmarks" "go test -bench=. -benchmem ./..."
    fi
    
    # 6. Generate coverage report
    echo -e "${BLUE}=== Coverage Report ===${NC}"
    echo ""
    
    echo "Generating coverage report..."
    go test -coverprofile=coverage/coverage.out ./... > /dev/null 2>&1
    go tool cover -func=coverage/coverage.out | tail -1
    
    # Generate HTML coverage report
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    echo -e "${GREEN}✓ Coverage report saved to coverage/coverage.html${NC}"
    echo ""
    
    # 7. Run linting (optional)
    if [ "$1" == "--with-lint" ] || [ "$2" == "--with-lint" ]; then
        echo -e "${BLUE}=== Linting ===${NC}"
        echo ""
        
        if command -v golangci-lint &> /dev/null; then
            run_test_suite "Go Linting" "golangci-lint run ./..."
        else
            echo -e "${YELLOW}⚠ golangci-lint not installed, skipping${NC}"
            SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
        fi
    fi
    
    # Calculate elapsed time
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))
    
    # Print summary
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}           Test Summary               ${NC}"
    echo -e "${BLUE}======================================${NC}"
    echo ""
    echo -e "Total Tests:    $TOTAL_TESTS"
    echo -e "${GREEN}Passed Tests:   $PASSED_TESTS${NC}"
    echo -e "${RED}Failed Tests:   $FAILED_TESTS${NC}"
    echo -e "${YELLOW}Skipped Tests:  $SKIPPED_TESTS${NC}"
    echo -e "Time Elapsed:   ${elapsed}s"
    echo ""
    
    # Exit with appropriate code
    if [ $FAILED_TESTS -gt 0 ]; then
        echo -e "${RED}❌ Test suite failed${NC}"
        exit 1
    else
        echo -e "${GREEN}✅ All tests passed!${NC}"
        exit 0
    fi
}

# Show usage
if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --with-performance    Run performance benchmarks"
    echo "  --with-lint          Run linting checks"
    echo "  --help, -h           Show this help message"
    echo ""
    exit 0
fi

# Run main function
main "$@"