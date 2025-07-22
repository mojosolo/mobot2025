# 👤 MoBot 2025 User Guide

## Overview

MoBot 2025 provides powerful tools for working with Adobe After Effects templates, from simple text extraction to complex automation workflows. This guide covers all user-facing features and viewers.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Template Import](#template-import)
3. [Story Viewers](#story-viewers)
4. [Search and Discovery](#search-and-discovery)
5. [Automation Workflows](#automation-workflows)
6. [Export Options](#export-options)
7. [Best Practices](#best-practices)

## Getting Started

### Quick Setup

```bash
# Start the API server
./mobot api

# Access the web interface
open http://localhost:8080
```

### First Steps

1. **Import a Template**: Upload your AEP file
2. **View Content**: Choose a viewer based on your needs
3. **Search Templates**: Find templates by content or properties
4. **Create Workflows**: Automate repetitive tasks

## Template Import

### Web Upload

1. Navigate to `http://localhost:8080`
2. Click "Import Template"
3. Drag and drop your AEP file or click to browse
4. Wait for processing (typically 5-30 seconds)

### Command Line Import

```bash
# Import single file
./mobot import template.aep

# Import directory
./mobot import ./templates/

# Import with metadata
./mobot import template.aep --name "Marketing Template" --tags "social,animated"
```

### API Import

```bash
curl -X POST http://localhost:8080/api/templates/import \
  -F "file=@template.aep" \
  -F "name=My Template"
```

## Story Viewers

MoBot 2025 includes four different viewers, each designed for specific use cases:

### 🎯 Ultimate Story Viewer (Recommended)

**Best for**: Organizations needing a complete solution

```bash
./start-ultimate-viewer.sh
```

The Ultimate Viewer includes three modes:

#### Simple Mode ✨
- **For**: Non-technical users, clients, content teams
- **Features**:
  - Clean scene cards (max 10 scenes)
  - Text-only view without technical details
  - One-click actions for common tasks
  - Mobile-friendly interface

#### Easy Mode 🔬
- **For**: Developers, AE engineers, technical users
- **Features**:
  - Full timeline view of all text elements
  - Complete metadata (layers, comps, timing)
  - Natural language chat for modifications
  - JSON export for pipelines

#### Advanced Mode 🎯
- **For**: Project managers, technical analysts
- **Features**:
  - Complete project statistics dashboard
  - Composition analysis tables
  - Effect usage reports
  - Issue detection and warnings
  - Professional dark theme

### 📝 Simple Story Viewer

**Best for**: Quick text extraction for non-technical users

```bash
./start-simple-viewer.sh
```

**Features**:
- Simplified scene cards
- No technical terminology
- Quick export to text
- Perfect for content review

**Workflow**:
1. Upload AEP file
2. Review scene cards
3. Export clean text
4. Share with team

### 🔬 Easy Mode Story Viewer

**Best for**: Technical users needing full control

```bash
./start-easy-mode.sh
```

**Features**:
- Timeline visualization
- Layer-level detail
- Chat interface for edits
- JSON export with metadata

**Workflow**:
1. Upload AEP file
2. Analyze timeline
3. Use chat to modify text
4. Export JSON for rendering

### 🔄 Unified Story Viewer

**Best for**: Teams with mixed technical levels

```bash
./start-unified-viewer.sh
```

**Features**:
- Switch between Simple and Easy modes
- Same features as standalone versions
- Mode selection on startup
- Requires re-upload to switch modes

## Search and Discovery

### Basic Search

```bash
# Search via web interface
Navigate to: http://localhost:8080/search

# Search via CLI
./mobot search "motion graphics"

# Search via API
curl "http://localhost:8080/api/templates/search?q=motion+graphics"
```

### Advanced Search

Use filters to refine results:

- **Automation Score**: Find templates ready for automation
  ```
  min_score:0.8
  ```

- **Tags**: Filter by categories
  ```
  tags:social,marketing
  ```

- **Date Range**: Recent templates
  ```
  created:2025-07-01..2025-07-31
  ```

### Search Types

1. **Keyword Search**: Basic text matching
2. **Semantic Search**: AI-powered concept matching
3. **Pattern Search**: Find similar structures

Example combined search:
```
"text animation" min_score:0.7 tags:promotional type:semantic
```

## Automation Workflows

### Creating a Workflow

1. **Via Web Interface**:
   - Go to Workflows → Create New
   - Select templates
   - Choose agents
   - Configure settings
   - Start workflow

2. **Via API**:
   ```json
   POST /api/workflows
   {
     "name": "Update Brand Colors",
     "template_ids": [1, 2, 3],
     "agents": ["planning", "implementation"],
     "config": {
       "parallel": true,
       "quality_threshold": 0.9
     }
   }
   ```

### Workflow Types

#### Text Replacement Workflow
Automatically update text across multiple templates:
```json
{
  "type": "text_replacement",
  "mappings": {
    "{{company}}": "MoBot Inc.",
    "{{year}}": "2025",
    "{{tagline}}": "Automate Everything"
  }
}
```

#### Asset Update Workflow
Replace images or videos:
```json
{
  "type": "asset_update",
  "assets": {
    "logo.png": "new-logo.png",
    "background.mp4": "new-bg.mp4"
  }
}
```

#### Batch Processing Workflow
Process multiple templates with same settings:
```json
{
  "type": "batch_process",
  "templates": ["template1", "template2", "template3"],
  "operations": ["optimize", "validate", "export"]
}
```

### Monitoring Workflows

Track workflow progress in real-time:

1. **Web Dashboard**: `http://localhost:8080/workflows`
2. **CLI**: `./mobot workflow status <id>`
3. **API**: `GET /api/workflows/<id>/status`

Status indicators:
- 🟡 **Planning**: Analyzing requirements
- 🔵 **Running**: Executing tasks
- 🟢 **Complete**: Successfully finished
- 🔴 **Failed**: Error occurred
- ⏸️ **Paused**: Waiting for approval

## Export Options

### Text Export

#### Simple Text
```bash
# Export clean text only
./mobot export <template-id> --format text

# Output:
Scene 1: Welcome to our presentation
Scene 2: Introducing new features
Scene 3: Thank you for watching
```

#### JSON Export
```bash
# Export with full metadata
./mobot export <template-id> --format json

# Output includes:
- Text content
- Layer names
- Timing information
- Composition structure
```

### Automation Export

#### NexRender Format
```bash
# Export for NexRender processing
./mobot export <template-id> --format nexrender

# Creates job.json with:
- Template reference
- Asset mappings
- Output settings
```

#### Custom Format
```bash
# Export with custom template
./mobot export <template-id> --format custom --template my-format.tmpl
```

### Batch Export

```bash
# Export multiple templates
./mobot export --batch template1,template2,template3 --format json

# Export search results
./mobot export --search "marketing" --format csv
```

## Best Practices

### 1. Template Organization

- **Use consistent naming**: `[Category]_[Name]_[Version]`
- **Tag thoroughly**: Add relevant tags during import
- **Add descriptions**: Help others understand template purpose

### 2. Workflow Design

- **Start simple**: Test with one template before batch processing
- **Set quality thresholds**: Ensure automated changes meet standards
- **Enable notifications**: Stay informed of workflow progress

### 3. Performance Tips

- **Batch similar templates**: Process related templates together
- **Use appropriate viewer**: Choose based on user needs
- **Cache search results**: Reuse common searches

### 4. Collaboration

- **Share viewer links**: Each viewer generates shareable URLs
- **Export for review**: Use appropriate format for audience
- **Document workflows**: Save successful workflow configs

## Common Use Cases

### 1. Marketing Campaign Update
```bash
# Import campaign templates
./mobot import ./campaign/*.aep --tags campaign2025

# Create text update workflow
./mobot workflow create \
  --type text_update \
  --templates "tags:campaign2025" \
  --data campaign-text.json

# Export results
./mobot export --search "tags:campaign2025" --format nexrender
```

### 2. Brand Refresh
```bash
# Find all templates with old branding
./mobot search "old logo" --min-score 0.7

# Create asset replacement workflow
./mobot workflow create \
  --type asset_replace \
  --search "old logo" \
  --assets brand-assets.json

# Verify changes
./start-ultimate-viewer.sh
```

### 3. Template Audit
```bash
# Analyze all templates
./mobot analyze --all --output audit-report.json

# View in Advanced Mode
./start-ultimate-viewer.sh
# Select Advanced Mode for detailed analysis
```

## Troubleshooting

### Viewer Issues

**Viewer won't start**:
```bash
# Check if port is in use
lsof -i :8080

# Use alternative port
./start-viewer.sh --port 8081
```

**Upload fails**:
- Check file size (max 500MB)
- Ensure valid AEP format
- Try command-line import

**Text not extracted**:
- Verify text layers aren't rasterized
- Check composition visibility
- Review layer naming

### Workflow Issues

**Workflow stuck**:
- Check for pending approvals
- Review agent status
- Check system resources

**Quality check failures**:
- Lower quality threshold
- Review error logs
- Manually verify one template

### Performance Issues

**Slow processing**:
- Reduce batch size
- Enable parallel processing
- Check system resources

**Search timeout**:
- Use more specific queries
- Add filters to narrow results
- Enable search caching

## Keyboard Shortcuts

### In Viewers

- `⌘/Ctrl + O`: Open file
- `⌘/Ctrl + E`: Export
- `⌘/Ctrl + F`: Search
- `⌘/Ctrl + S`: Save changes
- `ESC`: Close dialog

### In Chat (Easy Mode)

- `↑`: Previous command
- `↓`: Next command
- `Tab`: Autocomplete
- `⌘/Ctrl + Enter`: Submit

## Getting Help

### Resources

1. **Built-in Help**: Type `/help` in any chat interface
2. **API Docs**: `http://localhost:8080/docs`
3. **Examples**: See [Examples](EXAMPLES.md)
4. **Troubleshooting**: See [Troubleshooting Guide](TROUBLESHOOTING.md)

### Support

- **GitHub Issues**: Report bugs and request features
- **Discord**: Join our community for discussions
- **Email**: support@mobot2025.ai

---

For technical details, see [Developer Guide](DEVELOPER_GUIDE.md). For API usage, see [API Reference](API_REFERENCE.md).