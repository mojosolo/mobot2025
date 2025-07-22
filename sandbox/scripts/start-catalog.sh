#!/bin/bash

# MoBot 2025 Catalog System Startup Script

echo "🚀 MoBot 2025 Video Production Pipeline"
echo "======================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}❌ Python 3 is not installed. Please install Python 3 first.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Prerequisites checked${NC}"

# Build the Go parser if needed
if [ ! -f "catalog/bin/aep_parser" ]; then
    echo -e "${YELLOW}Building Go AEP parser...${NC}"
    mkdir -p catalog/bin
    go build -o catalog/bin/aep_parser catalog/cmd/parser/main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Parser built successfully${NC}"
    else
        echo -e "${RED}❌ Failed to build parser${NC}"
        exit 1
    fi
fi

# Build the main command
if [ ! -f "bin/mobot2025" ]; then
    echo -e "${YELLOW}Building mobot2025 command...${NC}"
    mkdir -p bin
    go build -o bin/mobot2025 cmd/mobot2025/main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Command built successfully${NC}"
    else
        echo -e "${RED}❌ Failed to build command${NC}"
        exit 1
    fi
fi

# Create necessary directories
mkdir -p reports renders catalog/reports

echo ""
echo "Available Commands:"
echo "==================="
echo ""
echo "1. Parse a single AEP file:"
echo "   ./bin/mobot2025 parse -file <your-file.aep>"
echo ""
echo "2. Deep analysis with dangerous mode:"
echo "   ./bin/mobot2025 analyze -file <your-file.aep> -deep -report"
echo ""
echo "3. Catalog a directory:"
echo "   ./bin/mobot2025 catalog -dir <templates-folder> -report"
echo ""
echo "4. Import MoBot templates:"
echo "   ./bin/mobot2025 import -dir ../mobot -output import_report.md"
echo ""
echo "5. Start API server:"
echo "   ./bin/mobot2025 serve -port 8080"
echo ""
echo "6. Submit render job:"
echo "   ./bin/mobot2025 render -file <your-file.aep> -config render.json"
echo ""
echo "🔍 API Endpoints (when server is running):"
echo "   Search: curl 'http://localhost:8080/api/v1/search?q=logo&limit=10'"
echo "   Filter: curl -X POST http://localhost:8080/api/v1/filter \\"
echo "           -H 'Content-Type: application/json' \\"
echo "           -d '{\"categories\":[\"HD\",\"Text Animation\"],\"limit\":20}'"
echo ""

# Quick demo option
if [ "$1" == "demo" ]; then
    echo -e "${YELLOW}Running demo analysis...${NC}"
    
    # Find a sample AEP file
    SAMPLE_AEP=$(find . -name "*.aep" -type f | head -1)
    
    if [ -n "$SAMPLE_AEP" ]; then
        echo -e "Found sample: $SAMPLE_AEP"
        ./bin/mobot2025 analyze -file "$SAMPLE_AEP" -deep -report
    else
        echo -e "${RED}No AEP files found for demo${NC}"
    fi
fi

# Start API server option
if [ "$1" == "serve" ]; then
    echo -e "${GREEN}Starting API server on port 8080...${NC}"
    ./bin/mobot2025 serve -port 8080
fi

# Catalog option
if [ "$1" == "catalog" ] && [ -n "$2" ]; then
    echo -e "${GREEN}Cataloging directory: $2${NC}"
    ./bin/mobot2025 catalog -dir "$2" -report
fi

# Import option
if [ "$1" == "import" ] && [ -n "$2" ]; then
    echo -e "${GREEN}Importing from MoBot directory: $2${NC}"
    ./bin/mobot2025 import -dir "$2" -output import_report.md -export
fi