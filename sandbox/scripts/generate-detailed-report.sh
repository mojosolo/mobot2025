#!/bin/bash

# Detailed HTML Report Generator for AEP files
echo "📊 AEP Detailed Report Generator"
echo "================================"
echo ""

# Check if file argument provided
if [ $# -eq 0 ]; then
    echo "Usage: ./generate-detailed-report.sh <aep-file>"
    echo ""
    echo "Examples:"
    echo "  ./generate-detailed-report.sh sample-aep/Ai\\ Text\\ Intro.aep"
    echo "  ./generate-detailed-report.sh data/Layer-01.aep"
    echo ""
    echo "Features:"
    echo "  📊 Overview - Project statistics and main compositions"
    echo "  🎬 Compositions - All comps with layer details"
    echo "  📑 Layers - Complete layer listing with properties"
    echo "  📹 Media - All assets with usage tracking"
    echo "  📝 Text - Text layer detection"
    echo "  ⚙️ Attributes - Layer properties and effects"
    echo "  📁 Hierarchy - Full project structure"
    echo ""
    exit 1
fi

# Generate the detailed report
echo "🔄 Generating detailed report for: $1"
echo "   This may take a moment for large projects..."
echo ""

go run demo/generate_detailed_report.go "$1"

echo ""
echo "💡 Features:"
echo "  - 7 different tabs for comprehensive analysis"
echo "  - Search functionality on each tab"
echo "  - Usage tracking shows which comps use each asset"
echo "  - Fully responsive design"
echo "  - Standalone HTML - share with anyone!"