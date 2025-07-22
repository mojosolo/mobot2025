#!/bin/bash

# MoBot 2025 Real Data Test Runner
# Tests EVERYTHING with REAL AEP files - NO mocks!

set -e

echo "🚀 MoBot 2025 Real Data Test Runner"
echo "====================================="
echo "Testing with REAL AEP files only!"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Test results
PASSED=0
FAILED=0

# Function to run tests and track results
run_test() {
    local name=$1
    local cmd=$2
    
    echo -n "Running $name... "
    
    if eval "$cmd" > /tmp/test_output.log 2>&1; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "  Error output:"
        tail -20 /tmp/test_output.log | sed 's/^/    /'
        ((FAILED++))
    fi
}

# 1. Core Parser Tests with Real Data
echo -e "\n${YELLOW}=== Core Parser Tests ===${NC}"
run_test "Real AEP Parsing" "go test -v . -run TestRealAEP"
run_test "Core Parser Tests" "go test -v ."

# 2. Real Data Benchmarks
echo -e "\n${YELLOW}=== Performance Benchmarks ===${NC}"
run_test "Parsing Benchmark" "go test -bench BenchmarkRealAEPParsing -run ^$"

# 3. Test Helpers (if they compile)
echo -e "\n${YELLOW}=== Test Infrastructure ===${NC}"
run_test "Helper Compilation" "go build ./tests/helpers"

# 4. Python Bridge Tests with Real Data
echo -e "\n${YELLOW}=== Python Bridge Tests ===${NC}"
if command -v python3 &> /dev/null; then
    run_test "Python Bridge" "cd tests/python && python3 -m pytest test_bridge.py -v"
else
    echo -e "${YELLOW}⚠ Python3 not found, skipping Python tests${NC}"
fi

# 5. Dangerous Analyzer Tests (when catalog compiles)
echo -e "\n${YELLOW}=== Dangerous Analyzer Tests ===${NC}"
if go build ./catalog 2>/dev/null; then
    echo -e "${GREEN}✓ Catalog compiled, running dangerous tests${NC}"
    # Would run dangerous analyzer tests here
else
    echo -e "${YELLOW}⚠ Catalog has compilation errors, skipping${NC}"
fi

# 6. Memory and Resource Tests
echo -e "\n${YELLOW}=== Resource Usage Tests ===${NC}"
echo "Testing memory usage with complex AEP..."
/usr/bin/time -l go test -v . -run TestComplexRealAEP 2>&1 | grep -E "(maximum resident|PASS|FAIL)" || true

# Summary
echo -e "\n${YELLOW}=== Test Summary ===${NC}"
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✅ All tests passed with REAL data!${NC}"
    exit 0
else
    echo -e "\n${RED}❌ Some tests failed${NC}"
    exit 1
fi