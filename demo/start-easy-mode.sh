#!/bin/bash

# Easy Mode Story Viewer Launcher
echo "🚀 Starting Easy Mode Story Viewer..."
echo "=================================="

# Check if port 8080 is available
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Port 8080 is already in use. Trying port 8081..."
    PORT=8081
else
    PORT=8080
fi

# Start the server
echo "📖 Starting server on http://localhost:$PORT"
echo "📤 Upload your AEP file to see the story!"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Run the server
go run easy_mode_story_viewer.go

# Alternative: Test mode with a specific file
# go run easy_mode_story_viewer.go test "../sample-aep/Ai Text Intro.aep"