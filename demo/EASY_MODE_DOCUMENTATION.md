# Easy Mode Story Viewer - Documentation

## Overview

The Easy Mode Story Viewer is a user-friendly web interface designed for non-developers to view and interact with After Effects projects as stories. It extracts text elements from AEP files, presents them in narrative order, and allows users to chat with AI to modify the content.

## Features

### 1. **Simple Upload Interface**
- Drag-and-drop or click-to-upload AEP files
- No technical knowledge required
- Clean, modern UI with gradient design

### 2. **Story Extraction**
- Automatically extracts all text layers from AEP projects
- Orders content by scene/composition sequence
- Shows text in chronological narrative order
- Identifies placeholder text vs. actual content

### 3. **Visual Timeline**
- Beautiful timeline view with numbered elements
- Click any element to select it
- Shows composition name, layer name, and timing
- Visual badges for placeholder text

### 4. **AI Chat Interface**
- Natural language processing for story modifications
- Supports commands like:
  - "Make the story more exciting"
  - "Change all placeholder text to be about cats"
  - "Make it shorter"
  - "Use professional business language"
- Generates JSON output for backend rendering

### 5. **Export Options**
- Copy JSON to clipboard
- Download JSON file
- Ready for integration with rendering systems

## Usage

### Starting the Server

```bash
# Option 1: Use the launcher script
./start-easy-mode.sh

# Option 2: Run directly
go run easy_mode_story_viewer.go

# Option 3: Test mode with a file
go run easy_mode_story_viewer.go test "path/to/file.aep"
```

### Web Interface

1. **Access**: Open http://localhost:8080 in your browser
2. **Upload**: Drag an AEP file onto the upload area
3. **View**: See your story in chronological order
4. **Chat**: Enter natural language requests to modify the story
5. **Export**: Copy or download the JSON for rendering

## Technical Details

### Story Extraction Algorithm

1. **Composition Ordering**: 
   - Sorts by scene numbers (Scene 1, S01, etc.)
   - Falls back to alphabetical order

2. **Text Extraction**:
   - Uses layer names as primary text source
   - Falls back to TextDocument data if available
   - Identifies placeholders by brackets [] or "placeholder" keyword

3. **Timeline Calculation**:
   - Calculates cumulative timeline positions
   - Shows start/end times for each text element

### NLP Processing

The system includes a simple NLP processor that handles:
- **Tone changes**: exciting, dramatic, professional
- **Content replacement**: themed replacements (cats, etc.)
- **Length modifications**: shortening content
- **Style adjustments**: business language, formality

### API Endpoints

- `GET /` - Main web interface
- `POST /upload` - File upload endpoint
- `POST /nlp` - Natural language processing

### Data Structure

```json
{
  "projectName": "My Project",
  "totalScenes": 5,
  "totalDuration": 30.0,
  "elements": [
    {
      "order": 1,
      "text": "Welcome to our presentation",
      "layerName": "Title Text",
      "compName": "Scene 1 - Intro",
      "timeStart": 0,
      "timeEnd": 5,
      "duration": 5,
      "isPlaceholder": false
    }
  ],
  "extractedAt": "2025-07-20 15:30:00"
}
```

## Design Philosophy

### For Non-Developers
- **No technical jargon**: Everything is explained in simple terms
- **Visual feedback**: See the story as it would appear
- **Natural language**: Chat with AI using plain English
- **One-click actions**: Simple buttons for all functions

### For Integration
- **Clean JSON output**: Ready for any rendering system
- **RESTful API**: Easy to integrate with other tools
- **Self-contained**: No external dependencies
- **Extensible**: Easy to add new NLP features

## Example Use Cases

1. **Marketing Teams**: Review and modify promotional video text
2. **Content Creators**: Quickly see all text in a project
3. **Translators**: Extract text for translation workflows
4. **Project Managers**: Review narrative flow without After Effects
5. **Clients**: Provide feedback on placeholder content

## Future Enhancements

1. **Database Integration**: Load projects from database
2. **Advanced AI**: GPT-4 integration for smarter modifications
3. **Multi-language**: Support for non-English content
4. **Collaboration**: Real-time multi-user editing
5. **Version Control**: Track changes over time
6. **Templates**: Save and apply story templates

## Troubleshooting

### Port Already in Use
The launcher script automatically tries port 8081 if 8080 is busy.

### Upload Fails
- Ensure file is a valid .aep file
- Check file size (limit: 10MB)
- Verify server is running

### No Text Found
- Some AEP files may not contain text layers
- Text might be in nested compositions
- Check if text is stored as expressions

## Summary

Easy Mode Story Viewer transforms complex After Effects projects into simple, readable stories that anyone can understand and modify. It bridges the gap between technical motion graphics files and everyday content creation needs.