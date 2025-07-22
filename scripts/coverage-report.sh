#!/bin/bash

# coverage-report.sh - Generate detailed coverage reports with visualization

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Create reports directory
mkdir -p reports/coverage

echo "🔬 Generating detailed coverage analysis..."

# Run tests with coverage if coverage.out doesn't exist
if [ ! -f "coverage.out" ]; then
    echo "Running tests to generate coverage data..."
    go test -coverprofile=coverage.out -covermode=atomic ./...
fi

# Generate different report formats
echo ""
echo "📊 Generating coverage reports..."

# 1. Text report
echo "  • Text report..."
go tool cover -func=coverage.out > reports/coverage/coverage.txt

# 2. HTML report
echo "  • HTML report..."
go tool cover -html=coverage.out -o reports/coverage/coverage.html

# 3. Function-level report
echo "  • Function-level analysis..."
go tool cover -func=coverage.out | grep -v "100.0%" | sort -k3 -n > reports/coverage/uncovered-functions.txt

# 4. Package summary
echo "  • Package summary..."
cat > reports/coverage/package-summary.txt << 'EOF'
Package Coverage Summary
========================

EOF

# Extract package-level coverage
go list ./... | grep -v vendor | while read pkg; do
    coverage=$(go test -cover $pkg 2>&1 | grep -o '[0-9]*\.[0-9]*%' | head -1)
    if [ -n "$coverage" ]; then
        printf "%-60s %s\n" "$pkg" "$coverage" >> reports/coverage/package-summary.txt
    fi
done

# 5. Generate coverage heatmap (ASCII)
echo "  • Coverage heatmap..."
cat > reports/coverage/coverage-heatmap.txt << 'EOF'
Coverage Heatmap
================

Legend: ■ = >90% | ▣ = 70-90% | ▤ = 50-70% | ▥ = 30-50% | □ = <30%

EOF

# Process each package for heatmap
go list ./... | grep -v vendor | while read pkg; do
    coverage=$(go test -cover $pkg 2>&1 | grep -o '[0-9]*\.[0-9]*' | head -1)
    if [ -n "$coverage" ]; then
        cov_int=$(echo $coverage | awk '{print int($1)}')
        
        # Determine symbol based on coverage
        if [ $cov_int -ge 90 ]; then
            symbol="■"
            color=$GREEN
        elif [ $cov_int -ge 70 ]; then
            symbol="▣"
            color=$GREEN
        elif [ $cov_int -ge 50 ]; then
            symbol="▤"
            color=$YELLOW
        elif [ $cov_int -ge 30 ]; then
            symbol="▥"
            color=$YELLOW
        else
            symbol="□"
            color=$RED
        fi
        
        # Short package name
        short_pkg=$(echo $pkg | sed 's|github.com/mojosolo/mobot2025/||')
        printf "${color}%s${NC} %-40s %5.1f%%\n" "$symbol" "$short_pkg" "$coverage" >> reports/coverage/coverage-heatmap.txt
    fi
done

# 6. Critical paths analysis
echo "  • Critical paths analysis..."
cat > reports/coverage/critical-paths.txt << 'EOF'
Critical Code Paths - Coverage Analysis
=======================================

Core Parser:
EOF
go tool cover -func=coverage.out | grep -E "(parser|Parser)" | grep -v test >> reports/coverage/critical-paths.txt

echo -e "\nCatalog System:" >> reports/coverage/critical-paths.txt
go tool cover -func=coverage.out | grep -E "(catalog|Catalog)" | grep -v test >> reports/coverage/critical-paths.txt

echo -e "\nAgents:" >> reports/coverage/critical-paths.txt
go tool cover -func=coverage.out | grep -E "(agent|Agent)" | grep -v test >> reports/coverage/critical-paths.txt

# 7. Generate recommendations
echo "  • Generating recommendations..."
cat > reports/coverage/recommendations.txt << 'EOF'
Coverage Improvement Recommendations
====================================

Based on the coverage analysis, here are the recommended areas for improvement:

1. Low Coverage Files (Priority: High)
--------------------------------------
EOF

go tool cover -func=coverage.out | grep -E "\.go[[:space:]]" | awk '$3 < 50.0' | sort -k3 -n | head -10 | while read line; do
    echo "  - $line" >> reports/coverage/recommendations.txt
done

cat >> reports/coverage/recommendations.txt << 'EOF'

2. Uncovered Critical Functions
--------------------------------
EOF

go tool cover -func=coverage.out | grep -E "(Parse|Process|Analyze|Execute)" | grep "0.0%" | head -10 | while read line; do
    echo "  - $line" >> reports/coverage/recommendations.txt
done

cat >> reports/coverage/recommendations.txt << 'EOF'

3. Suggested Test Cases
-----------------------
- Error handling paths
- Edge cases for parsers
- Concurrent operations
- Database transaction failures
- Network timeout scenarios

4. Quick Wins
-------------
- Add tests for utility functions
- Cover error return paths
- Test configuration loading
- Add validation tests
EOF

# Summary
echo ""
echo "✅ Coverage reports generated in reports/coverage/"
echo ""
echo "📁 Available reports:"
echo "  • coverage.txt          - Function coverage listing"
echo "  • coverage.html         - Interactive HTML report"
echo "  • package-summary.txt   - Package-level summary"
echo "  • coverage-heatmap.txt  - Visual coverage heatmap"
echo "  • critical-paths.txt    - Coverage of critical code"
echo "  • uncovered-functions.txt - Functions needing tests"
echo "  • recommendations.txt   - Improvement recommendations"

# Display quick summary
echo ""
echo "📈 Quick Summary:"
total_coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "  Total Coverage: $total_coverage"

# Count of packages below threshold
low_coverage_count=$(go tool cover -func=coverage.out | grep -E "\.go[[:space:]]" | awk '$3 < 80.0' | wc -l)
echo "  Files below 80%: $low_coverage_count"

# Open HTML report if requested
if [ "$1" == "--open" ]; then
    echo ""
    echo "📖 Opening HTML report..."
    open reports/coverage/coverage.html 2>/dev/null || xdg-open reports/coverage/coverage.html 2>/dev/null || echo "Please open reports/coverage/coverage.html manually"
fi