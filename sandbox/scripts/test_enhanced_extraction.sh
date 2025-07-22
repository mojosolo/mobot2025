#!/bin/bash

# Build and test enhanced text extraction

echo "🔨 Building enhanced text extraction test..."

# Navigate to demo directory
cd demo

# Build the test program
go build -o test_text_extraction test_text_extraction_final.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"
echo ""

# Test with sample files if they exist
if [ -f "../samples/Ai Text Intro.aep" ]; then
    echo "🧪 Testing with Ai Text Intro.aep..."
    ./test_text_extraction "../samples/Ai Text Intro.aep"
elif [ -f "Ai Text Intro.aep" ]; then
    echo "🧪 Testing with Ai Text Intro.aep..."
    ./test_text_extraction "Ai Text Intro.aep"
else
    echo "⚠️  No sample AEP files found. Usage:"
    echo "   ./test_text_extraction <path-to-aep-file>"
fi

# Clean up
rm -f test_text_extraction