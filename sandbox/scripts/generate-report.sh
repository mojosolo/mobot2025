#!/bin/bash

# HTML Report Generator for AEP files
echo "📊 AEP HTML Report Generator"
echo "============================"
echo ""

# Check if file argument provided
if [ $# -eq 0 ]; then
    echo "Usage: ./generate-report.sh <aep-file>"
    echo ""
    echo "Examples:"
    echo "  ./generate-report.sh sample-aep/Ai\\ Text\\ Intro.aep"
    echo "  ./generate-report.sh data/Layer-01.aep"
    echo ""
    exit 1
fi

# Generate the report
echo "🔄 Generating HTML report for: $1"
go run demo/generate_html_report.go "$1"

echo ""
echo "💡 Tips:"
echo "  - Report includes interactive stats and charts"
echo "  - Open the HTML file in any web browser"
echo "  - Share the HTML file - it's completely standalone!"