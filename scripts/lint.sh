#!/bin/bash

# lint.sh - Comprehensive code quality checks

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 MoBot 2025 Code Quality Check${NC}"
echo "================================="
echo ""

FAILED=false

# 1. Go fmt
echo -e "${YELLOW}Checking code formatting...${NC}"
UNFMT_FILES=$(gofmt -l .)
if [ -n "$UNFMT_FILES" ]; then
    echo -e "${RED}❌ The following files need formatting:${NC}"
    echo "$UNFMT_FILES"
    echo ""
    echo "Run 'make fmt' to fix formatting"
    FAILED=true
else
    echo -e "${GREEN}✓ All files properly formatted${NC}"
fi
echo ""

# 2. Go vet
echo -e "${YELLOW}Running go vet...${NC}"
if go vet ./...; then
    echo -e "${GREEN}✓ No vet issues found${NC}"
else
    echo -e "${RED}❌ go vet found issues${NC}"
    FAILED=true
fi
echo ""

# 3. golangci-lint
echo -e "${YELLOW}Running golangci-lint...${NC}"
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run --timeout=5m; then
        echo -e "${GREEN}✓ No lint issues found${NC}"
    else
        echo -e "${RED}❌ golangci-lint found issues${NC}"
        FAILED=true
    fi
else
    echo -e "${YELLOW}⚠ golangci-lint not installed${NC}"
    echo "Install with: make install-tools"
fi
echo ""

# 4. goimports
echo -e "${YELLOW}Checking imports...${NC}"
if command -v goimports &> /dev/null; then
    IMPORT_ISSUES=$(goimports -l .)
    if [ -n "$IMPORT_ISSUES" ]; then
        echo -e "${RED}❌ The following files have import issues:${NC}"
        echo "$IMPORT_ISSUES"
        echo ""
        echo "Run 'goimports -w .' to fix imports"
        FAILED=true
    else
        echo -e "${GREEN}✓ All imports properly organized${NC}"
    fi
else
    echo -e "${YELLOW}⚠ goimports not installed${NC}"
fi
echo ""

# 5. Check for TODO/FIXME comments
echo -e "${YELLOW}Checking for TODO/FIXME comments...${NC}"
TODO_COUNT=$(grep -r "TODO\|FIXME" --include="*.go" . 2>/dev/null | grep -v vendor | wc -l)
if [ "$TODO_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}⚠ Found $TODO_COUNT TODO/FIXME comments:${NC}"
    grep -r "TODO\|FIXME" --include="*.go" . | grep -v vendor | head -10
    if [ "$TODO_COUNT" -gt 10 ]; then
        echo "... and $((TODO_COUNT - 10)) more"
    fi
else
    echo -e "${GREEN}✓ No TODO/FIXME comments found${NC}"
fi
echo ""

# 6. Check for debugging prints
echo -e "${YELLOW}Checking for debug prints...${NC}"
DEBUG_COUNT=$(grep -r "fmt\.Print\|log\.Print\|println" --include="*.go" . 2>/dev/null | grep -v vendor | grep -v "_test.go" | wc -l)
if [ "$DEBUG_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}⚠ Found $DEBUG_COUNT potential debug prints:${NC}"
    grep -r "fmt\.Print\|log\.Print\|println" --include="*.go" . | grep -v vendor | grep -v "_test.go" | head -5
    echo ""
    echo "Consider using proper logging instead"
fi
echo ""

# 7. Check for large files
echo -e "${YELLOW}Checking for large files...${NC}"
LARGE_FILES=$(find . -name "*.go" -type f -size +500k 2>/dev/null | grep -v vendor)
if [ -n "$LARGE_FILES" ]; then
    echo -e "${YELLOW}⚠ Found large Go files (>500KB):${NC}"
    echo "$LARGE_FILES"
    echo "Consider breaking these into smaller files"
else
    echo -e "${GREEN}✓ No excessively large files${NC}"
fi
echo ""

# 8. Check for security issues
echo -e "${YELLOW}Checking for common security issues...${NC}"
SECURITY_ISSUES=0

# Check for hardcoded credentials
if grep -r "password\s*=\s*\"" --include="*.go" . 2>/dev/null | grep -v vendor | grep -v "_test.go" | grep -v "example"; then
    echo -e "${RED}⚠ Potential hardcoded passwords found${NC}"
    SECURITY_ISSUES=$((SECURITY_ISSUES + 1))
fi

# Check for SQL injection vulnerabilities
if grep -r "fmt\.Sprintf.*SELECT\|fmt\.Sprintf.*INSERT\|fmt\.Sprintf.*UPDATE" --include="*.go" . 2>/dev/null | grep -v vendor; then
    echo -e "${RED}⚠ Potential SQL injection vulnerability${NC}"
    SECURITY_ISSUES=$((SECURITY_ISSUES + 1))
fi

if [ "$SECURITY_ISSUES" -eq 0 ]; then
    echo -e "${GREEN}✓ No obvious security issues found${NC}"
fi
echo ""

# Summary
echo "================================="
if [ "$FAILED" = true ]; then
    echo -e "${RED}❌ Code quality checks failed${NC}"
    echo ""
    echo "Please fix the issues above before committing."
    exit 1
else
    echo -e "${GREEN}✅ All code quality checks passed!${NC}"
fi