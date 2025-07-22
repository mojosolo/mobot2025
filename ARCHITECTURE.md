# mobot2025 Architecture

## Overview

The mobot2025 system is an advanced Adobe After Effects Project (AEP) file processing pipeline that combines high-performance binary parsing with intelligent multi-agent orchestration. The architecture leverages the strengths of multiple technologies:

- **Go** for fast, efficient binary parsing of AEP files
- **Python** for compatibility with existing mobot tools and workflows
- **REST API** for service integration and scalability
- **Multi-Agent LLM System** for intelligent analysis and automation

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         mobot2025 AEP Processing System                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Multi-Agent Orchestration Layer               │   │
│  │  ┌─────────┐  ┌──────────────┐  ┌──────────────┐  ┌─────────┐  │   │
│  │  │Planning │  │Implementation│  │ Verification │  │  Review  │  │   │
│  │  │ Agent   │  │    Agent     │  │    Agent     │  │  Agent   │  │   │
│  │  └────┬────┘  └──────┬───────┘  └──────┬───────┘  └────┬────┘  │   │
│  │       │              │                  │                │       │   │
│  │       └──────────────┴──────────────────┴───────────────┘       │   │
│  │                            │                                     │   │
│  │                    ┌───────▼────────┐                          │   │
│  │                    │Meta-Orchestrator│                          │   │
│  │                    └───────┬────────┘                          │   │
│  └────────────────────────────┼────────────────────────────────────┘   │
│                              │                                          │
│  ┌───────────────────────────▼─────────────────────────────────────┐   │
│  │                    Core Processing Pipeline                      │   │
│  │  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐     │   │
│  │  │  Go Parser  │────▶│ Bridge Layer │────▶│ Python API  │     │   │
│  │  │ (Fast I/O)  │     │  (Unified)   │     │(Compatible) │     │   │
│  │  └─────────────┘     └──────────────┘     └─────────────┘     │   │
│  │         │                    │                      │           │   │
│  │         ▼                    ▼                      ▼           │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │              Catalog Database (JSON/SQL)                 │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  │                             │                                   │   │
│  │                             ▼                                   │   │
│  │  ┌─────────────────────────────────────────────────────────┐   │   │
│  │  │                    REST API Service                      │   │   │
│  │  │  /parse  /catalog  /analyze  /opportunities  /jobs      │   │   │
│  │  └─────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

## Multi-Agent Architecture

### Agent Roles and Responsibilities

#### 1. Planning Agent
- **Role**: Analyzes AEP file structures and decomposes parsing tasks
- **Guardrails**: "Systematically evaluate existing parsing logic and extend capabilities by referencing specific block types"
- **Tools**: 
  - File analysis for AEP structure patterns
  - Documentation parsing for Adobe specifications
- **Output**: JSON plan with parsing steps and confidence scores (>80%)

#### 2. Implementation Agent
- **Role**: Extends parser code for new block types and features
- **Guardrails**: "Implement functional parsing code that integrates seamlessly with existing RIFX framework"
- **Approach**:
  - Reuse existing parseItem, parseLayer patterns
  - Add new block type handlers incrementally
  - Maintain backward compatibility

#### 3. Verification Agent
- **Role**: Validates parsing accuracy against test AEP files
- **Guardrails**: "Empirically validate all parsing logic with real AEP samples"
- **Tests**:
  - Unit tests for each block type
  - Integration tests with sample AEP files
  - Binary accuracy validation

#### 4. Final Review Agent
- **Role**: Ensures code quality and performance benchmarks
- **Guardrails**: "Review for efficient, maintainable parsing logic that advances project capabilities"
- **Checks**:
  - Memory efficiency for large files
  - Error handling completeness
  - API consistency

#### 5. Meta-Orchestrator Agent
- **Role**: Coordinates parsing workflow and manages iterations
- **Guardrails**: "Facilitate verifiable parsing improvements with structured testing"
- **Controls**:
  - Max 50 iteration limit
  - Human approval gates
  - Progress tracking

## Core Components

### 1. Go Parser (High Performance Layer)
- **Binary RIFX format parsing**: Direct reading of AEP file structure
- **Text extraction**: Multiple strategies for extracting text from binary data
- **Media asset discovery**: Identifies all images, videos, and audio files
- **Effect tracking**: Catalogs all effects and their usage
- **Composition analysis**: Full composition tree with dependencies
- **Performance**: ~100ms for average AEP file

### 2. Python Bridge (Compatibility Layer)
- **Backward compatibility**: Works with existing mobot Python tools
- **Format conversion**: Translates between Go and Python formats
- **Batch processing**: Catalog entire directories
- **Report generation**: JSON and Markdown reports
- **Performance**: ~200ms overhead for compatibility

### 3. REST API Service
- **Asynchronous processing**: Background jobs for large catalogs
- **Caching**: 15-minute TTL for parsed projects
- **NexRender integration**: Generate render configurations
- **Health monitoring**: Service status and metrics
- **Performance**: <500ms for cached results

### 4. Catalog Database
- **Storage formats**: JSON for flexibility, SQL for queries
- **Project metadata**: Complete project information
- **Opportunity tracking**: Automation opportunities with scoring
- **Capability detection**: What can be customized
- **Performance**: ~10 projects/second batch processing

## Data Models

### Project Metadata Structure
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

## Implementation Strategy

### Phase 1: Foundation (Current Sprint)
1. Analyze existing parser capabilities
2. Identify missing AEP block types
3. Create test suite with sample files
4. Document current limitations

### Phase 2: Enhancement
1. Implement new block parsers
2. Add streaming support
3. Enhance error handling
4. Create validation tools

### Phase 3: Integration
1. Performance optimization
2. API refinement
3. Documentation updates
4. Release preparation

## Code Reuse Principles
- Extend existing types rather than create new ones
- Leverage RIFX library patterns
- Maintain consistent error handling
- Reuse test utilities

## Integration Points

### NexRender Integration
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

### External Services
- **Adobe After Effects**: Direct integration via ExtendScript
- **Render Farms**: Batch job submission
- **Storage Services**: Asset management
- **Analytics**: Usage tracking and optimization

## Performance Targets
- Parse 90%+ of AEP block types
- 80%+ test coverage
- <100ms parse time for typical projects
- Zero breaking changes to existing API
- <500ms API response time
- 10+ projects/second batch processing

## Security Considerations
- Input validation for all file paths
- Sandboxed parsing environment
- Rate limiting on API endpoints
- Authentication for sensitive operations
- Audit logging for all actions

## Deployment Architecture
- **Containerized services**: Docker for each component
- **Load balancing**: Horizontal scaling for API layer
- **Queue management**: Background job processing
- **Monitoring**: Health checks and metrics
- **Logging**: Centralized log aggregation

## Future Enhancements
1. WebSocket support for real-time updates
2. Machine learning for template categorization
3. Advanced effect chain analysis
4. GPU-accelerated rendering integration
5. Distributed parsing for large projects

## Development Guidelines
- Follow Go best practices for parser code
- Use Python type hints for bridge layer
- Implement comprehensive error handling
- Write tests for all new functionality
- Document all API endpoints
- Maintain backward compatibility