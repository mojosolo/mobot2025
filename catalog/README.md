# MoBot 2025 Catalog System Implementation

## Overview

This directory contains the core implementation of the MoBot 2025 catalog system. The catalog provides comprehensive functionality for parsing, analyzing, and automating Adobe After Effects templates through a sophisticated multi-agent architecture.

For complete documentation, please refer to:
- 📖 [Architecture Overview](/docs/ARCHITECTURE.md) - System design and components
- 📡 [API Reference](/docs/API_REFERENCE.md) - REST API documentation
- 🤖 [Multi-Agent System](/docs/MULTI_AGENT_SYSTEM.md) - Agent details
- 👨‍💻 [Developer Guide](/docs/DEVELOPER_GUIDE.md) - Implementation details

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AEP Processing Pipeline                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐ │
│  │   Go Parser │────▶│ Bridge Layer │────▶│ Python API  │ │
│  │  (Fast I/O) │     │  (Unified)   │     │ (Compatible)│ │
│  └─────────────┘     └──────────────┘     └─────────────┘ │
│         │                    │                      │       │
│         ▼                    ▼                      ▼       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Catalog Database (JSON/SQL)             │   │
│  └─────────────────────────────────────────────────────┘   │
│                             │                               │
│                             ▼                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    REST API Service                  │   │
│  │  /parse  /catalog  /analyze  /opportunities  /jobs  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Key Features

### 1. Enhanced AEP Parser (Go)
- **Binary RIFX format parsing** - Direct reading of AEP file structure
- **Text extraction** - Multiple strategies for extracting text from binary data
- **Media asset discovery** - Identifies all images, videos, and audio files
- **Effect tracking** - Catalogs all effects and their usage
- **Composition analysis** - Full composition tree with dependencies

### 2. Intelligent Cataloging
- **Automatic categorization** - Resolution, duration, content type
- **Tag generation** - Searchable tags based on capabilities
- **Opportunity identification** - Automation opportunities with difficulty/impact scoring
- **Capability detection** - What can be customized in each template

### 3. Python Bridge
- **Backward compatibility** - Works with existing mobot Python tools
- **Format conversion** - Translates between Go and Python formats
- **Batch processing** - Catalog entire directories
- **Report generation** - JSON and Markdown reports

### 4. REST API
- **Asynchronous processing** - Background jobs for large catalogs
- **Caching** - 15-minute TTL for parsed projects
- **NexRender integration** - Generate render configurations
- **Health monitoring** - Service status and metrics

## Usage

### Command Line

**Parse single AEP file:**
```bash
# Using Go parser directly
go run cmd/parser/main.go -file template.aep -output metadata.json

# Using Python bridge
python catalog/python_bridge.py parse template.aep --mobot-format
```

**Catalog directory:**
```bash
# Catalog all AEP files in templates directory
python catalog/python_bridge.py catalog ./templates -o catalog.json

# Generate human-readable report
python catalog/python_bridge.py report ./templates -o catalog_report.json
```

### API Usage

**Start the API server:**
```bash
go run catalog/cmd/api/main.go -port 8080
```

**Parse a project:**
```bash
curl -X POST http://localhost:8080/api/v1/parse \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "/path/to/template.aep",
    "options": {
      "extract_text": true,
      "extract_media": true,
      "deep_analysis": true
    }
  }'
```

**Find opportunities:**
```bash
curl -X POST http://localhost:8080/api/v1/opportunities \
  -H "Content-Type: application/json" \
  -d '{
    "file_paths": ["/path/to/template1.aep", "/path/to/template2.aep"],
    "criteria": {
      "min_impact": "medium",
      "max_difficulty": "medium",
      "types": ["text_automation", "media_automation"]
    }
  }'
```

## Data Structures

### Project Metadata
```json
{
  "file_path": "/path/to/project.aep",
  "file_name": "project.aep",
  "bit_depth": 8,
  "expression_engine": "javascript",
  "compositions": [...],
  "text_layers": [...],
  "media_assets": [...],
  "capabilities": {
    "has_text_replacement": true,
    "has_image_replacement": true,
    "has_color_control": false,
    "has_audio_replacement": false,
    "is_modular": true
  },
  "categories": ["HD", "Text Animation", "Short Form"],
  "tags": ["text", "animated-text", "modular"],
  "opportunities": [...]
}
```

### Opportunity Structure
```json
{
  "type": "text_automation",
  "description": "Automate 5 text layers for dynamic content",
  "difficulty": "easy",
  "impact": "high",
  "components": ["Title", "Subtitle", "CTA", "Date", "Location"]
}
```

## Integration Points

### With NexRender
The system generates NexRender-compatible configurations for automated rendering:
```json
{
  "template": {
    "src": "template.aep",
    "composition": "Main"
  },
  "assets": [
    {
      "type": "data",
      "layerName": "Title",
      "property": "Source Text",
      "value": "New Title"
    }
  ]
}
```

### With Multi-Agent System
Future integration will include:
- **Planning Agent** - Analyzes projects and creates processing plans
- **Implementation Agent** - Executes modifications and replacements
- **Verification Agent** - Validates changes and quality
- **Review Agent** - Optimizes performance and output
- **Meta-Orchestrator** - Coordinates the entire workflow

## Performance

- **Go Parser**: ~100ms for average AEP file
- **Python Bridge**: ~200ms overhead for compatibility
- **API Response**: <500ms for cached results
- **Batch Processing**: ~10 projects/second

## Next Steps

1. Complete Phase 1 tasks (in progress)
2. Implement database persistence
3. Add WebSocket support for real-time updates
4. Integrate with render farm
5. Deploy multi-agent orchestration

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for development guidelines.