# Ultimate Story Viewer - Complete Guide

## Overview

The **Ultimate Story Viewer** combines three different viewing modes into a single application, providing the perfect experience for every user type - from clients to developers to technical analysts.

## Quick Start

```bash
./start-ultimate-viewer.sh
```

Then open: http://localhost:8080

## Three Viewing Modes

### 🌟 Mode Selection Screen

When you first open the viewer, you'll see three options:

```
┌───────────────────────────────────────────┐
│         Choose Your Experience            │
├───────────────────────────────────────────┤
│                                           │
│  [Simple Mode ✨]  [Easy Mode 🔬]        │
│     Scene Cards     Timeline View         │
│                                           │
│          [Advanced Mode 🎯]               │
│        Technical Report                   │
│                                           │
└───────────────────────────────────────────┘
```

## Simple Mode ✨

**Best for**: Clients, marketers, content creators, non-technical users

### Features:
- **Scene Cards**: Clean, visual cards (max 10 scenes)
- **Just the Essentials**: Text content without technical details
- **Quick Actions**: One-click text modifications
- **Mobile Friendly**: Works perfectly on any device
- **Clean Export**: Copy text with one click

### How to Use:
1. Click "Simple Mode ✨" from the home page
2. Upload your AEP file (drag & drop supported)
3. View your scenes as clean cards
4. Click cards to expand/collapse
5. Use quick action buttons to modify text
6. Export clean text when done

### Quick Actions Available:
- **Make Professional**: Business-appropriate language
- **Make Casual**: Friendly, approachable tone
- **Make Shorter**: Concise version
- **Fill Placeholders**: Auto-fill empty text
- **Export Text**: Copy all text

## Easy Mode 🔬

**Best for**: Developers, motion designers, technical users

### Features:
- **Timeline View**: All text elements in chronological order
- **Full Metadata**: Layer names, comp names, timing info
- **Chat Interface**: Natural language text modifications
- **JSON Export**: Complete data for rendering pipelines
- **Technical Details**: Everything needed for debugging

### How to Use:
1. Click "Easy Mode 🔬" from the home page
2. Upload your AEP file
3. Browse the complete timeline
4. Use chat to modify text naturally
5. Export JSON for your pipeline

### Sample Chat Commands:
- "Make all text more exciting"
- "Replace placeholders with product names"
- "Shorten everything to 5 words"
- "Use formal business language"
- "Make it sound more friendly"

## Advanced Mode 🎯

**Best for**: Technical directors, project analysts, QA teams

### Features:
- **Complete Analysis**: Full project breakdown
- **Technical Stats**: Compositions, layers, effects, assets
- **Visual Dashboard**: Dark theme with data visualization
- **Detailed Tables**: Sortable composition data
- **Effect Analysis**: See all effects and usage counts
- **Warning System**: Potential issues highlighted
- **Full Export**: Download complete technical report

### How to Use:
1. Click "Advanced Mode 🎯" from the home page
2. Upload your AEP file
3. View comprehensive statistics
4. Browse detailed composition data
5. Review text layers with metadata
6. Check effect usage and warnings
7. Download full JSON report

### Report Sections:
- **Project Statistics**: 6 key metrics at a glance
- **Compositions Table**: Resolution, framerate, duration, layers
- **Text Layers**: All text with layer/comp context
- **Effects Usage**: Count of each effect type
- **Warnings**: Potential issues or optimizations

## Mode Switching

You can switch between modes at any time using the buttons in the top-right corner:

- From Simple Mode → Easy Mode or Advanced Mode
- From Easy Mode → Simple Mode or Advanced Mode  
- From Advanced Mode → Simple Mode or Easy Mode

**Note**: Switching modes requires re-uploading your file.

## Comparison Table

| Feature | Simple Mode | Easy Mode | Advanced Mode |
|---------|-------------|-----------|---------------|
| **View Type** | Scene Cards | Timeline | Dashboard |
| **Max Items** | 10 scenes | All elements | All data |
| **Technical Info** | None | Moderate | Complete |
| **Modification** | Buttons | Chat/NLP | View only |
| **Export Format** | Plain text | JSON | Full JSON |
| **Best For** | Everyone | Developers | Analysts |
| **Theme** | Light | Light | Dark |
| **Learning Curve** | None | Minimal | Moderate |

## Technical Architecture

```
ultimate_story_viewer.go
├── Landing Page (/)
├── Simple Mode (/simple)
│   ├── Upload endpoint (/upload/simple)
│   └── Action endpoint (/action)
├── Easy Mode (/easy)
│   ├── Upload endpoint (/upload/easy)
│   └── NLP endpoint (/nlp)
└── Advanced Mode (/advanced)
    └── Upload endpoint (/upload/advanced)
```

## Use Case Examples

### Marketing Manager Maria
- Uses **Simple Mode** to quickly review video text
- Applies "Make Professional" to ensure brand voice
- Exports clean text for approval

### Developer David  
- Uses **Easy Mode** to see text in timeline order
- Uses chat to batch modify placeholder text
- Exports JSON for rendering pipeline

### Technical Director Tom
- Uses **Advanced Mode** to analyze project complexity
- Reviews effect usage for performance optimization
- Downloads full report for documentation

## Troubleshooting

### Port Already in Use
The launcher automatically tries port 8081 if 8080 is busy.

### File Upload Issues
- Ensure it's a valid .aep file
- Maximum file size: 10MB
- Check browser console for errors

### Mode Switching
Remember that switching modes requires re-uploading your file.

### No Text Showing
Some AEP files may not contain text layers. Check Advanced Mode for confirmation.

## Advanced Tips

### Direct Mode URLs
Skip the selection screen:
- Simple Mode: http://localhost:8080/simple
- Easy Mode: http://localhost:8080/easy
- Advanced Mode: http://localhost:8080/advanced

### Workflow Recommendations
1. Start with **Advanced Mode** to understand the project
2. Use **Easy Mode** for text modifications
3. Switch to **Simple Mode** for client presentation

### Performance Notes
- Simple Mode: Fastest, limits to 10 scenes
- Easy Mode: Fast, shows all elements
- Advanced Mode: May take longer for complex projects

## Summary

The Ultimate Story Viewer provides three distinct experiences:

1. **Simple Mode ✨** - For everyone who just needs the text
2. **Easy Mode 🔬** - For developers who need control
3. **Advanced Mode 🎯** - For analysts who need everything

Choose the mode that fits your needs and enjoy a tailored experience!

## Version History

- **v3.0** - Ultimate Viewer with three modes
- **v2.0** - Unified Viewer with two modes  
- **v1.0** - Separate Easy and Simple viewers

---

*Built with Go and the mobot2025 AEP parser*