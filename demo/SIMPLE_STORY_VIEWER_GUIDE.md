# Simple Story Viewer - Quick Guide

## What's New?

This is the **simplified version** of the Easy Mode Story Viewer, designed with these personas in mind:
- 📊 **Marketing Managers** - Quick copy review
- 🎨 **Content Creators** - Easy text variations  
- 👔 **Clients** - Simple approval interface
- 🌍 **Translators** - Clean text export
- 📋 **Project Managers** - Clear overview

## Key Improvements

### Before vs After

| Old Version | New Simple Version |
|------------|-------------------|
| 27 text elements in timeline | 5-10 scene cards |
| Technical metadata visible | Only text content |
| Complex chat interface | One-click actions |
| JSON output visible | Hidden by default |
| Gradients and animations | Clean, simple cards |

## How to Use

### 1. Start the Server
```bash
./start-simple-viewer.sh
```

### 2. Open Your Browser
Navigate to: http://localhost:8080

### 3. Upload Your Project
- Click "Open Your Project" button
- Or drag & drop your .aep file

### 4. View Your Story
- See up to 10 scenes maximum
- Click any scene to expand/collapse
- Only the main text is shown by default

### 5. Quick Actions
Instead of typing, just click:
- **Make Professional** - Business language
- **Make Casual** - Friendly tone
- **Make Shorter** - Brevity mode
- **Fill Placeholders** - Auto-fill empty spots
- **Export Text** - Copy clean text

## What's Hidden?

To keep it simple, we've hidden:
- ❌ Timeline view
- ❌ Layer names
- ❌ Composition names  
- ❌ Timing information
- ❌ JSON output
- ❌ Technical metadata
- ❌ Complex animations

## Mobile Friendly

✅ Works great on phones
✅ Touch-friendly buttons
✅ Responsive design
✅ Fast loading

## Export Options

Click "Export Text" to get:
- Clean text only (no JSON)
- Ready for email
- Perfect for translation
- No technical jargon

## Success Metrics

- 📊 View story in **<3 seconds**
- ✏️ Modify text in **<10 seconds**
- 📤 Export in **<5 seconds**
- 📱 Works on any device

## Troubleshooting

**Port already in use?**
The script automatically tries port 8081 if 8080 is busy.

**Can't see all text?**
Click on a scene card to expand and see more.

**Need the technical version?**
Use `./start-easy-mode.sh` for the full-featured version.

## Design Philosophy

> "Make it so simple that a client can use it without instructions."

We removed 80% of features to focus on what matters:
1. See the text
2. Change the text
3. Export the text

That's it. No complexity. Just results.