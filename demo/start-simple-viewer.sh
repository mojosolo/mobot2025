#!/bin/bash

# Simple Story Viewer Launcher
echo "🎯 Starting Simple Story Viewer..."
echo "================================"

# Check if port 8080 is available
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Port 8080 is already in use. Trying port 8081..."
    PORT=8081
else
    PORT=8080
fi

# Kill any existing story viewers
pkill -f "story_viewer" 2>/dev/null

# Start the server
echo "📖 Starting server on http://localhost:$PORT"
echo "🚀 Simple, clean interface for everyone!"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Set port in environment for the Go app to read
export PORT=$PORT

# Run the simplified server
go run simple_story_viewer.go