#!/bin/bash

# Enhanced UX Report Generator Script
# This script generates comprehensive HTML reports with maximum platform utilization

echo "🚀 Enhanced UX Report Generator for mobot2025"
echo "============================================"

# Check if AEP file is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <aep-file>"
    echo "Example: $0 \"../sample-aep/Ai Text Intro.aep\""
    exit 1
fi

AEP_FILE="$1"

# Check if file exists
if [ ! -f "$AEP_FILE" ]; then
    echo "❌ Error: File not found: $AEP_FILE"
    exit 1
fi

echo "📄 Processing: $AEP_FILE"

# Generate enhanced UX report
echo "🎨 Generating enhanced UX report..."
go run generate_enhanced_ux_report.go "$AEP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Enhanced report generated successfully!"
    
    # Find the latest report file
    LATEST_REPORT=$(ls -t *-enhanced-ux-report-*.html 2>/dev/null | head -n 1)
    
    if [ -n "$LATEST_REPORT" ]; then
        echo "📊 Report details:"
        echo "   - File: $LATEST_REPORT"
        echo "   - Size: $(ls -lh "$LATEST_REPORT" | awk '{print $5}')"
        echo ""
        echo "🌐 Opening in browser..."
        
        # Open in default browser
        if command -v open &> /dev/null; then
            open "$LATEST_REPORT"
        elif command -v xdg-open &> /dev/null; then
            xdg-open "$LATEST_REPORT"
        else
            echo "📂 Please open manually: file://$(pwd)/$LATEST_REPORT"
        fi
    fi
else
    echo "❌ Error generating report"
    exit 1
fi