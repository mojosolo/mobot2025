# 🏗️ MoBot 2025 Architecture

## Overview

MoBot 2025 is a sophisticated AI-powered system that combines deep binary parsing capabilities with multi-agent orchestration to enable intelligent automation of Adobe After Effects templates. The architecture seamlessly integrates high-performance Go components with Python compatibility layers and an intelligent multi-agent system.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         REST API Service                         │
│                    (HTTP endpoints, JSON responses)              │
├─────────────────────────────────────────────────────────────────┤
│                    Multi-Agent Orchestration Layer              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  Planning   │  │Implementation│  │Verification│            │
│  │   Agent     │  │    Agent     │  │   Agent    │            │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘            │
│         │                 │                 │                    │
│  ┌──────┴─────────────────┴─────────────────┴──────┐           │
│  │              Meta-Orchestrator Agent             │           │
│  └──────┬───────────────────────────────────┬──────┘           │
│         │                                   │                    │
│  ┌──────┴──────┐                    ┌──────┴──────┐            │
│  │   Review    │                    │Communication│            │
│  │   Agent     │                    │   System    │            │
│  └─────────────┘                    └─────────────┘            │
├─────────────────────────────────────────────────────────────────┤
│                    Core Processing Pipeline                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │   Parser    │  │  Analyzer   │  │   Catalog   │            │
│  │   Engine    │  │   Engine    │  │  Database   │            │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘            │
│         │                 │                 │                    │
│  ┌──────┴─────────────────┴─────────────────┴──────┐           │
│  │              Binary Processing Core              │           │
│  │            (RIFX Format, Block Types)           │           │
│  └──────────────────────────────────────────────────┘           │
├─────────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │   SQLite    │  │   Search    │  │  Quality    │            │
│  │  Database   │  │   Engine    │  │ Assurance   │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
```

## Multi-Agent Architecture

### Planning Agent
- **Purpose**: Analyzes AEP file structures and creates comprehensive parsing plans
- **Responsibilities**:
  - Break down complex parsing tasks into manageable subtasks
  - Create file reference mappings for accurate parsing
  - Generate confidence scores for each parsing approach
  - Identify dependencies between different block types
- **Key Features**:
  - Task decomposition algorithms
  - Confidence scoring system (0.0-1.0)
  - File structure analysis
  - Dependency graph generation

### Implementation Agent
- **Purpose**: Extends parser code to handle newly discovered block types
- **Responsibilities**:
  - Generate parsing code for new block types
  - Implement model cascading (Claude → GPT-4 → Gemini)
  - Integrate RIFX pattern support
  - Ensure backward compatibility
- **Key Features**:
  - Multi-model code generation
  - Pattern recognition for block types
  - Code quality validation
  - Automatic fallback mechanisms

### Verification Agent
- **Purpose**: Validates parsing accuracy and code quality
- **Responsibilities**:
  - Execute automated tests for new parsers
  - Validate parsing results against sample files
  - Ensure 80% minimum test coverage
  - Generate quality metrics reports
- **Key Features**:
  - Automated test generation
  - Coverage analysis tools
  - Performance benchmarking
  - Regression testing

### Review Agent
- **Purpose**: Ensures code quality, performance, and maintainability
- **Responsibilities**:
  - Analyze code for performance bottlenecks
  - Generate optimization recommendations
  - Review maintainability metrics
  - Provide improvement suggestions
- **Key Features**:
  - Performance profiling
  - Code complexity analysis
  - Best practices enforcement
  - Technical debt tracking

### Meta-Orchestrator Agent
- **Purpose**: Coordinates all agents and manages execution workflows
- **Responsibilities**:
  - Maintain workflow state across all agents
  - Detect and prevent infinite loops (50-iteration limit)
  - Coordinate inter-agent communication
  - Manage human approval gates
- **Key Features**:
  - Workflow state management
  - Loop detection algorithms
  - Event-driven coordination
  - Progress tracking and reporting

## Core Components

### 1. Go Parser (High Performance Layer)
```go
// Core parsing engine
type Parser struct {
    decoder     *RIFXDecoder
    blockTypes  map[string]BlockHandler
    cache       *BlockCache
    metrics     *ParsingMetrics
}
```

**Features:**
- RIFX binary format parsing
- 16+ block type handlers
- Streaming support for large files
- Concurrent block processing
- Memory-efficient operation

### 2. Python Bridge (Compatibility Layer)
```python
# Python interface for Go parser
class AEPParser:
    def __init__(self, parser_path):
        self.bridge = GoBridge(parser_path)
    
    def parse(self, file_path):
        return self.bridge.call("Parse", file_path)
```

**Features:**
- JSON-based communication
- Automatic type conversion
- Error propagation
- Async operation support

### 3. REST API Service
```
POST   /api/templates/import      - Import AEP file
GET    /api/templates             - List all templates
GET    /api/templates/:id         - Get template details
GET    /api/templates/search      - Search templates
GET    /api/templates/:id/score   - Get automation score
POST   /api/workflows             - Create workflow
GET    /api/agents/status         - Agent health status
```

**Features:**
- RESTful design principles
- JSON request/response format
- Comprehensive error handling
- Rate limiting and authentication
- WebSocket support for real-time updates

### 4. Catalog Database
```sql
-- Core schema
CREATE TABLE templates (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    file_path TEXT,
    metadata JSON,
    automation_score REAL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE workflows (
    id INTEGER PRIMARY KEY,
    template_id INTEGER,
    status TEXT,
    config JSON,
    results JSON
);
```

**Features:**
- SQLite for portability
- JSON storage for flexible metadata
- Full-text search indexes
- Migration support
- Backup and restore capabilities

## Data Flow

### 1. Template Import Flow
```
AEP File → Parser → Analyzer → Scoring → Database → Search Index
```

### 2. Automation Workflow
```
Request → Planning → Implementation → Verification → Review → Execution
```

### 3. Search and Discovery
```
Query → Search Engine → Scoring Filter → Quality Check → Results
```

## Agent Communication Protocol

### Message Format
```json
{
    "id": "msg_123",
    "type": "task_assignment",
    "from": "meta_orchestrator",
    "to": "planning_agent",
    "payload": {
        "task_id": "task_456",
        "action": "analyze_template",
        "parameters": {
            "template_id": "789",
            "depth": "comprehensive"
        }
    },
    "metadata": {
        "priority": "high",
        "timeout": 300,
        "correlation_id": "workflow_001"
    }
}
```

### State Synchronization
- Event-driven updates
- Distributed state management
- Conflict resolution protocols
- Message queuing with priorities

## Performance Architecture

### Optimization Strategies
1. **Caching**
   - Block-level caching for repeated structures
   - Query result caching with TTL
   - Metadata caching for quick access

2. **Concurrency**
   - Parallel block processing
   - Concurrent agent execution
   - Async API endpoints

3. **Resource Management**
   - Connection pooling
   - Memory limits per operation
   - Graceful degradation

### Performance Targets
- Parser: <5 seconds for 100MB AEP file
- API Response: <200ms for queries
- Agent Response: <1 second for simple tasks
- Workflow Completion: <5 minutes for standard templates

## Security Architecture

### Security Layers
1. **API Security**
   - JWT authentication
   - Rate limiting per endpoint
   - Input validation and sanitization

2. **File Security**
   - Virus scanning on upload
   - File type validation
   - Sandboxed parsing environment

3. **Data Security**
   - Encryption at rest
   - Secure communication channels
   - Audit logging

## Deployment Architecture

### Container Structure
```yaml
services:
  api:
    image: mobot2025/api
    ports: ["8080:8080"]
    
  agents:
    image: mobot2025/agents
    scale: 3
    
  database:
    image: sqlite:latest
    volumes: ["./data:/data"]
```

### Scaling Strategy
- Horizontal scaling for API servers
- Agent pool scaling based on load
- Database replication for read scaling
- CDN for static assets

## Integration Points

### 1. NexRender Integration
```javascript
// NexRender job creation
const job = {
    template: {
        src: "http://mobot2025/api/templates/123/download",
        composition: "Main"
    },
    assets: automationData.assets,
    actions: automationData.actions
}
```

### 2. External Services
- Cloud storage (S3, GCS) for template files
- Video encoding services
- Asset management systems
- Project management tools

## Quality Assurance Architecture

### Testing Strategy
1. **Unit Tests**
   - Parser components
   - Individual agents
   - API endpoints

2. **Integration Tests**
   - Multi-agent workflows
   - End-to-end scenarios
   - Performance benchmarks

3. **Quality Gates**
   - Code coverage >80%
   - Performance regression checks
   - Security vulnerability scanning

## Monitoring and Observability

### Metrics Collection
```
- API request/response times
- Agent task completion rates
- Parser performance metrics:
  * Simple AEP files: ~1ms parsing time
  * Complex AEP files (3000+ items): ~10ms parsing time
  * Memory usage: <50MB even for complex files
- Error rates and types
- Resource utilization
```

### Real-World Performance Insights
Based on real AEP file testing:
- The parser handles extremely complex projects (3157 items)
- RIFX binary parsing is highly efficient
- No performance degradation with nested structures
- Memory usage remains low even with large files

### Logging Strategy
- Structured JSON logging
- Distributed tracing with correlation IDs
- Log aggregation and analysis
- Real-time alerting

## Future Architecture Enhancements

### Planned Improvements
1. **Machine Learning Integration**
   - Template similarity detection
   - Automated optimization suggestions
   - Predictive performance analysis

2. **Advanced Workflows**
   - Visual workflow designer
   - Conditional branching
   - Custom agent development

3. **Cloud-Native Features**
   - Kubernetes operators
   - Service mesh integration
   - Multi-region deployment

## Development Guidelines

### Code Organization
```
mobot2025/
├── cmd/           # Entry points
├── catalog/       # Core business logic
├── agents/        # Agent implementations
├── api/           # REST API handlers
├── parser/        # Binary parsing
├── bridge/        # Language bridges
└── docs/          # Documentation
```

### Design Principles
1. **Modularity**: Each component should be independently deployable
2. **Extensibility**: Easy to add new block types and agents
3. **Resilience**: Graceful failure handling at every level
4. **Performance**: Optimization without sacrificing readability
5. **Security**: Defense in depth approach

### Contributing Architecture
- Follow established patterns
- Document architectural decisions
- Performance test new components
- Maintain backward compatibility
- Update architecture diagrams