# Unified Story Viewer - Complete Guide

## Overview

The **Unified Story Viewer** combines both Easy Mode and Simple Mode into a single application, letting users choose their preferred experience based on their needs and technical expertise.

## Quick Start

```bash
./start-unified-viewer.sh
```

Then open: http://localhost:8080

## Mode Selection

When you first open the viewer, you'll see:

```
┌─────────────────────────────────┐
│     Choose Your Experience      │
├─────────────────────────────────┤
│  [Easy Mode]    [Simple Mode]   │
│                                 │
│  Timeline View   Scene Cards    │
│  Full Details    Just Text      │
│  Chat Interface  Quick Actions  │
│  For Power Users For Everyone   │
└─────────────────────────────────┘
```

## Easy Mode 🔬

**Best for**: Developers, AE engineers, technical users

### Features:
- **Timeline View**: See all text elements in chronological order
- **Full Metadata**: Layer names, composition names, timing info
- **Chat Interface**: Natural language modifications
- **JSON Export**: Complete data structure for developers
- **Technical Details**: Everything you need for debugging

### How to Use:
1. Click "Launch Easy Mode" from the home page
2. Upload your AEP file
3. Browse the timeline of text elements
4. Use the chat to modify text
5. Export JSON for your rendering pipeline

### Sample Chat Commands:
- "Make the story more exciting"
- "Change all placeholder text to be about cats"
- "Make it shorter"
- "Use professional business language"

## Simple Mode ✨

**Best for**: Clients, marketers, content creators, everyone

### Features:
- **Scene Cards**: Maximum 10 clean cards
- **Just Text**: No technical jargon
- **Quick Actions**: One-click modifications
- **Clean Export**: Copy text instantly
- **Mobile Friendly**: Works on any device

### How to Use:
1. Click "Launch Simple Mode" from the home page
2. Click "Open Your Project" or drag & drop
3. Click scenes to expand/collapse
4. Use quick action buttons
5. Export clean text

### Quick Actions:
- **Make Professional**: Business language
- **Make Casual**: Friendly tone
- **Make Shorter**: Brevity mode
- **Fill Placeholders**: Auto-fill empty spots
- **Export Text**: Copy to clipboard

## Mode Switching

You can switch between modes anytime:
- In Easy Mode: Click "Switch to Simple Mode ✨" button (top right)
- In Simple Mode: Click "Switch to Easy Mode 🔬" button (top right)

**Note**: You'll need to re-upload your file when switching modes.

## Technical Details

### Architecture
```
unified_story_viewer.go
├── Landing Page (/)
├── Easy Mode (/easy)
│   ├── Upload endpoint (/upload/easy)
│   └── NLP endpoint (/nlp)
└── Simple Mode (/simple)
    ├── Upload endpoint (/upload/simple)
    └── Action endpoint (/action)
```

### Data Structures

**Easy Mode** uses full `StoryData`:
```json
{
  "projectName": "My Project",
  "totalScenes": 5,
  "totalDuration": 30.0,
  "elements": [...complete element data...]
}
```

**Simple Mode** uses simplified `SimpleStory`:
```json
{
  "projectName": "My Project",
  "scenes": [...max 10 scene cards...],
  "totalText": 27
}
```

## Comparison Table

| Feature | Easy Mode | Simple Mode |
|---------|-----------|-------------|
| View Type | Timeline | Scene Cards |
| Max Elements | All | 10 scenes |
| Technical Data | Yes | No |
| Modification | Chat/NLP | Buttons |
| Export | JSON | Text only |
| Best For | Developers | Everyone |

## Troubleshooting

### Port Already in Use
The launcher automatically tries port 8081 if 8080 is busy.

### Can't Upload File
- Ensure it's a valid .aep file
- File size limit: 10MB
- Check browser console for errors

### Switching Modes Lost My Work
Mode switching requires re-upload. Save your work before switching.

### Text Not Showing
Some AEP files may not have text layers. The viewer will show an empty state.

## Advanced Usage

### Direct Mode URLs
- Easy Mode: http://localhost:8080/easy
- Simple Mode: http://localhost:8080/simple
- Skip the selection page if you know what you want

### Keyboard Shortcuts
- None currently implemented (future enhancement)

### API Endpoints
- `POST /upload/easy` - Upload for Easy Mode
- `POST /upload/simple` - Upload for Simple Mode
- `POST /nlp` - Natural language processing
- `POST /action` - Quick actions

## Design Philosophy

The Unified Viewer follows these principles:

1. **Choice**: Let users pick their complexity level
2. **Clarity**: Clear differences between modes
3. **Consistency**: Shared visual language
4. **Conversion**: Easy switching between modes
5. **Completeness**: Both modes fully functional

## Future Enhancements

1. **Persistent Sessions**: Keep uploads when switching
2. **Unified Data Model**: One parse, two views
3. **Keyboard Shortcuts**: Power user features
4. **Export Options**: PDF, DOCX, etc.
5. **Collaboration**: Share links with mode preference

## Summary

The Unified Story Viewer gives you the best of both worlds:
- **Power** when you need it (Easy Mode)
- **Simplicity** when you don't (Simple Mode)

Choose your experience and start viewing stories!