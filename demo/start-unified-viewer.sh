#!/bin/bash

# Unified Story Viewer Launcher
echo "🎯 Starting Unified Story Viewer..."
echo "===================================="

# Kill any existing story viewers
echo "🔧 Cleaning up existing processes..."
pkill -f "story_viewer" 2>/dev/null
pkill -f "easy_mode" 2>/dev/null
pkill -f "simple_story" 2>/dev/null

# Check if port 8080 is available
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Port 8080 is already in use. Trying port 8081..."
    PORT=8081
else
    PORT=8080
fi

# Start the server
echo ""
echo "🚀 Server starting on http://localhost:$PORT"
echo ""
echo "📖 Features:"
echo "   • Easy Mode - Full timeline view with all details"
echo "   • Simple Mode - Clean scene cards for everyone"
echo "   • Switch between modes anytime!"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Set port in environment for the Go app to read
export PORT=$PORT

# Run the unified server
go run unified_story_viewer.go