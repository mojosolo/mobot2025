# 👨‍💻 MoBot 2025 Developer Guide

## Overview

This guide provides comprehensive information for developers who want to contribute to MoBot 2025, extend its functionality, or integrate it into their applications. 

## Table of Contents

1. [Development Setup](#development-setup)
2. [Project Structure](#project-structure)
3. [Architecture Overview](#architecture-overview)
4. [Coding Standards](#coding-standards)
5. [Adding Features](#adding-features)
6. [Testing](#testing)
7. [Building and Deployment](#building-and-deployment)
8. [Contributing](#contributing)

## Development Setup

### Prerequisites

- Go 1.19 or higher
- Python 3.8+ (for Python bridge)
- Git
- Make (optional but recommended)
- golangci-lint (for linting)

### Environment Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/mobot2025.git
cd mobot2025

# Install Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest

# Install pre-commit hooks
cp scripts/pre-commit .git/hooks/
chmod +x .git/hooks/pre-commit

# Set up Python environment (optional)
cd bridge/python
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows
pip install -e ".[dev]"
```

### IDE Configuration

#### VS Code
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "go.buildTags": "integration"
}
```

#### GoLand/IntelliJ
- Enable Go modules support
- Configure golangci-lint as external tool
- Set up file watchers for goimports

## Project Structure

```
mobot2025/
├── cmd/                    # Command-line entry points
│   └── parser/
│       └── main.go        # Main CLI application
├── catalog/               # Core business logic
│   ├── api_service.go     # REST API implementation
│   ├── database.go        # Database operations
│   ├── parser.go          # AEP parsing logic
│   ├── search_engine.go   # Search functionality
│   ├── automation_scoring.go  # Scoring algorithms
│   ├── planning_agent.go      # Agent implementations
│   ├── implementation_agent.go
│   ├── verification_agent.go
│   ├── review_agent.go
│   ├── meta_orchestrator.go
│   ├── agent_communication.go
│   ├── workflow_automation.go
│   ├── quality_assurance.go
│   └── system_integration_testing.go
├── parser/                # Binary parsing utilities
│   ├── rifx.go           # RIFX format handling
│   └── blocks.go         # Block type definitions
├── bridge/               # Language bridges
│   ├── python/           # Python integration
│   └── node/             # Node.js integration (future)
├── demo/                 # Demo applications
│   └── viewers/          # Story viewer implementations
├── docs/                 # Documentation
├── scripts/              # Build and utility scripts
├── tests/                # Test files
│   ├── unit/            # Unit tests
│   ├── integration/     # Integration tests
│   └── fixtures/        # Test data
└── vendor/              # Vendored dependencies (optional)
```

## Architecture Overview

### Core Components

#### 1. Parser Engine
The heart of MoBot 2025, responsible for reading and interpreting AEP files.

```go
// Key interfaces
type Parser interface {
    Parse(data []byte) (*Project, error)
    ParseFile(path string) (*Project, error)
}

type BlockHandler interface {
    CanHandle(blockType string) bool
    Parse(data []byte) (Block, error)
}
```

#### 2. Catalog System
Manages templates, metadata, and provides search capabilities.

```go
type Catalog interface {
    ImportTemplate(path string) (*Template, error)
    GetTemplate(id string) (*Template, error)
    Search(query string, opts SearchOptions) ([]*Template, error)
    CalculateAutomationScore(template *Template) float64
}
```

#### 3. Agent System
Implements the multi-agent orchestration for intelligent automation.

```go
type Agent interface {
    GetID() string
    GetType() string
    HandleMessage(msg *Message) error
    GetState() AgentState
}
```

### Data Flow

```
AEP File → Parser → Block Handlers → Project Model → Catalog → Database
                                                        ↓
                                          Search Index ← → API Service
                                                        ↓
                                                    Agent System
```

## Coding Standards

### Go Code Style

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

Key points:
- Use `gofmt` and `goimports` for formatting
- Keep functions focused and small
- Prefer composition over inheritance
- Handle errors explicitly
- Use meaningful variable names

#### Naming Conventions
```go
// Interfaces - noun with 'er' suffix
type Parser interface {}
type Handler interface {}

// Structs - noun
type Template struct {}
type Workflow struct {}

// Functions - verb or verb phrase
func ParseFile(path string) error {}
func CalculateScore(t *Template) float64 {}

// Constants - CamelCase or CAPS_WITH_UNDERSCORE
const MaxRetries = 3
const DEFAULT_TIMEOUT = 30 * time.Second
```

#### Error Handling
```go
// Always return errors
func ProcessTemplate(id string) (*Result, error) {
    template, err := db.GetTemplate(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get template: %w", err)
    }
    // ...
}

// Custom errors for specific cases
var (
    ErrTemplateNotFound = errors.New("template not found")
    ErrInvalidFormat    = errors.New("invalid file format")
)
```

#### Comments and Documentation
```go
// Package catalog provides template management and automation capabilities.
package catalog

// Template represents an After Effects project template with metadata.
// It includes parsing results, automation scoring, and search indices.
type Template struct {
    // ID is the unique identifier for the template
    ID string `json:"id"`
    
    // Name is the human-readable template name
    Name string `json:"name"`
    
    // AutomationScore indicates how suitable this template is for automation (0.0-1.0)
    AutomationScore float64 `json:"automation_score"`
}

// ImportTemplate imports an AEP file and creates a new template in the catalog.
// It performs parsing, analysis, and scoring before storing the template.
//
// Example:
//     template, err := catalog.ImportTemplate("marketing.aep")
//     if err != nil {
//         log.Fatal(err)
//     }
func (c *Catalog) ImportTemplate(path string) (*Template, error) {
    // Implementation
}
```

### Testing Standards

#### Test File Organization
```go
// parser_test.go
package catalog

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        want    *Project
        wantErr bool
    }{
        {
            name:  "valid project",
            input: validProjectData,
            want:  &Project{Name: "Test"},
        },
        {
            name:    "invalid format",
            input:   []byte("invalid"),
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := NewParser()
            got, err := p.Parse(tt.input)
            
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

#### Test Categories
- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete workflows
- **Benchmark Tests**: Measure performance

## Adding Features

### 1. Adding a New Block Type

```go
// 1. Define the block structure in parser/blocks.go
type TextLayerBlock struct {
    BaseBlock
    Text       string            `json:"text"`
    Font       string            `json:"font"`
    Size       float64           `json:"size"`
    Properties map[string]interface{} `json:"properties"`
}

// 2. Create handler in parser/handlers/text_layer.go
type TextLayerHandler struct{}

func (h *TextLayerHandler) CanHandle(blockType string) bool {
    return blockType == "TeLa" // Text Layer
}

func (h *TextLayerHandler) Parse(data []byte) (Block, error) {
    block := &TextLayerBlock{}
    // Parsing logic
    return block, nil
}

// 3. Register handler in parser/registry.go
func init() {
    RegisterHandler(&TextLayerHandler{})
}

// 4. Add tests in parser/handlers/text_layer_test.go
func TestTextLayerHandler_Parse(t *testing.T) {
    // Test implementation
}
```

### 2. Adding a New Agent

```go
// 1. Define agent in catalog/optimization_agent.go
type OptimizationAgent struct {
    id       string
    state    AgentState
    comm     *AgentCommunicationSystem
    analyzer *OptimizationAnalyzer
}

// 2. Implement Agent interface
func (oa *OptimizationAgent) GetID() string {
    return oa.id
}

func (oa *OptimizationAgent) HandleMessage(msg *Message) error {
    switch msg.Type {
    case "optimize_request":
        return oa.handleOptimizeRequest(msg)
    default:
        return fmt.Errorf("unknown message type: %s", msg.Type)
    }
}

// 3. Register with orchestrator
func (mo *MetaOrchestrator) RegisterOptimizationAgent() {
    agent := NewOptimizationAgent()
    mo.RegisterAgent(agent)
}

// 4. Add agent-specific logic
func (oa *OptimizationAgent) OptimizeTemplate(templateID string) (*OptimizationResult, error) {
    // Implementation
}
```

### 3. Adding API Endpoints

```go
// 1. Define handler in catalog/api_handlers.go
func (s *APIService) handleOptimize(w http.ResponseWriter, r *http.Request) {
    templateID := chi.URLParam(r, "id")
    
    result, err := s.optimizer.OptimizeTemplate(templateID)
    if err != nil {
        s.sendError(w, err)
        return
    }
    
    s.sendJSON(w, result)
}

// 2. Register route in catalog/api_routes.go
func (s *APIService) setupRoutes() {
    r := chi.NewRouter()
    
    r.Route("/api", func(r chi.Router) {
        r.Route("/templates/{id}", func(r chi.Router) {
            r.Post("/optimize", s.handleOptimize)
        })
    })
}

// 3. Add API documentation
// Update docs/API_REFERENCE.md with new endpoint details
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific tests
go test -run TestParser ./catalog

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./catalog

# Run integration tests
go test -tags=integration ./tests/integration
```

### Writing Tests

#### Unit Test Example
```go
func TestCalculateAutomationScore(t *testing.T) {
    template := &Template{
        TextElements: []TextElement{
            {Replaceable: true},
            {Replaceable: true},
            {Replaceable: false},
        },
        Effects: []Effect{
            {Complexity: "simple"},
            {Complexity: "complex"},
        },
    }
    
    score := CalculateAutomationScore(template)
    assert.InRange(t, score, 0.0, 1.0)
    assert.Greater(t, score, 0.5) // Should be automatable
}
```

#### Integration Test Example
```go
// +build integration

func TestWorkflowExecution(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer cleanupTestDB(db)
    
    orchestrator := NewMetaOrchestrator(db)
    
    // Create workflow
    workflow, err := orchestrator.CreateWorkflow(&WorkflowConfig{
        Name:   "test workflow",
        Agents: []string{"planning", "implementation"},
    })
    require.NoError(t, err)
    
    // Execute
    err = orchestrator.ExecuteWorkflow(workflow.ID)
    require.NoError(t, err)
    
    // Verify
    status, err := orchestrator.GetWorkflowStatus(workflow.ID)
    require.NoError(t, err)
    assert.Equal(t, "completed", status.Status)
}
```

### Test Data - REAL FILES ONLY

**CRITICAL: NO MOCK DATA ALLOWED**

We use ONLY real AEP files for testing:
```
data/                      # Real AEP test files
├── BPC-8.aep             # 8-bit color depth
├── BPC-16.aep            # 16-bit color depth
├── BPC-32.aep            # 32-bit color depth
├── ExEn-es.aep           # ExtendScript engine
├── ExEn-js.aep           # JavaScript engine
├── Item-01.aep           # Item structure test
├── Layer-01.aep          # Layer structure test
└── Property-01.aep       # Property structure test

sample-aep/
└── Ai Text Intro.aep     # Complex real project (3157 items!)
```

Key learnings from real data testing:
- Complex files parse successfully with 3000+ items
- Performance is excellent: ~10ms for complex files
- Real files reveal patterns mocks would miss
- The parser is more robust than expected

## Building and Deployment

### Building

```bash
# Development build
go build -o mobot cmd/parser/main.go

# Production build (with optimizations)
go build -ldflags="-s -w" -o mobot cmd/parser/main.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o mobot-linux
GOOS=darwin GOARCH=amd64 go build -o mobot-macos
GOOS=windows GOARCH=amd64 go build -o mobot.exe

# Using Makefile
make build
make build-all
```

### Docker

```dockerfile
# Dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o mobot cmd/parser/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mobot .
EXPOSE 8080
CMD ["./mobot", "api"]
```

### Deployment Options

#### 1. Systemd Service
```ini
[Unit]
Description=MoBot 2025 API Service
After=network.target

[Service]
Type=simple
User=mobot
WorkingDirectory=/opt/mobot
ExecStart=/opt/mobot/mobot api
Restart=always
Environment=MOBOT_ENV=production

[Install]
WantedBy=multi-user.target
```

#### 2. Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mobot-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mobot-api
  template:
    metadata:
      labels:
        app: mobot-api
    spec:
      containers:
      - name: mobot
        image: mobot2025:latest
        ports:
        - containerPort: 8080
        env:
        - name: MOBOT_ENV
          value: "production"
```

## Contributing

### Getting Started

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Pull Request Guidelines

#### PR Title Format
```
type(scope): description

Examples:
feat(parser): add support for shape layers
fix(api): handle nil pointer in search endpoint
docs(readme): update installation instructions
test(catalog): add integration tests for workflows
```

#### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new warnings
```

### Code Review Process

1. **Automated Checks**: CI runs tests, linting, and coverage
2. **Peer Review**: At least one maintainer review required
3. **Testing**: Reviewer should test locally if possible
4. **Documentation**: Ensure docs are updated

### Release Process

1. **Version Bump**: Update version in code
2. **Changelog**: Update CHANGELOG.md
3. **Tag Release**: `git tag -a v1.2.3 -m "Release v1.2.3"`
4. **Build Artifacts**: CI builds release binaries
5. **Publish**: Upload to GitHub releases

## Debugging

### Debug Mode

```bash
# Enable debug logging
export MOBOT_DEBUG=true
export MOBOT_LOG_LEVEL=debug

# Run with delve debugger
dlv debug cmd/parser/main.go -- api

# Attach to running process
dlv attach $(pgrep mobot)
```

### Performance Profiling

```go
// Add profiling endpoints
import _ "net/http/pprof"

// In main()
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Trace
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out
go tool trace trace.out
```

## Resources

### Internal Documentation
- [Architecture](ARCHITECTURE.md)
- [API Reference](API_REFERENCE.md)
- [Multi-Agent System](MULTI_AGENT_SYSTEM.md)

### External Resources
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [After Effects SDK](https://adobe.io/after-effects/)

### Community
- GitHub Discussions
- Discord Server
- Stack Overflow tag: `mobot2025`

---

Happy coding! If you have questions, please check existing issues or create a new one.