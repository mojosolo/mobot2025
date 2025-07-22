#!/bin/bash

# Parse your sample AEP file
echo "🎬 Parsing your Ai Text Intro project..."
echo ""

# Make sure we're in the right directory
cd "$(dirname "$0")"

# Run the parser
go run demo/parse_sample.go

echo ""
echo "📝 To parse a different AEP file:"
echo "   1. Edit demo/parse_sample.go"
echo "   2. Change the sampleFile variable to your AEP path"
echo "   3. Run this script again"