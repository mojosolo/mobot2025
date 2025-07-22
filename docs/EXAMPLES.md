# 🎯 MoBot 2025 Examples

This document provides practical examples for common use cases with MoBot 2025.

## Table of Contents

1. [Basic Usage](#basic-usage)
2. [Template Management](#template-management)
3. [Search Operations](#search-operations)
4. [Automation Workflows](#automation-workflows)
5. [Agent Integration](#agent-integration)
6. [API Integration](#api-integration)
7. [Advanced Scenarios](#advanced-scenarios)

## ⚠️ Important: AEP Files Required

**This repository does NOT include AEP test files.** All examples below assume you have your own Adobe After Effects project files. See [TEST_DATA_README.md](../TEST_DATA_README.md) for more information.

## Basic Usage

### Command Line Examples

```bash
# Parse a single AEP file (you must provide your own)
./mobot parse your-project.aep

# Parse complex project with verbose output
./mobot parse your-complex-project.aep --verbose

# Parse and export to JSON
./mobot parse your-project.aep --output json > template.json

# Start API server
./mobot serve

# Start API with custom port
./mobot serve --port 8090
```

### Go Library Examples

```go
package main

import (
    "fmt"
    "log"
    "github.com/mojosolo/mobot2025/catalog"
)

func main() {
    // Initialize catalog
    cat, err := catalog.NewCatalog("templates.db")
    if err != nil {
        log.Fatal(err)
    }
    defer cat.Close()

    // Import template (you must provide your own AEP file)
    template, err := cat.ImportTemplate("your-project.aep")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Template: %s (Score: %.2f)\n", 
        template.Name, template.AutomationScore)
}
```

## Expected AEP File Types

When testing MoBot 2025, you'll want to use various types of AEP files:

### Recommended Test Cases
```bash
# Color depth variations
your-8bit-project.aep    # 8-bit color depth projects
your-16bit-project.aep   # 16-bit color depth projects
your-32bit-project.aep   # 32-bit color depth projects

# Expression variations
your-extendscript-project.aep  # Projects with ExtendScript expressions
your-javascript-project.aep    # Projects with JavaScript expressions

# Structure variations
your-simple-project.aep      # Simple projects with few layers
your-complex-project.aep     # Complex projects with many compositions
your-nested-project.aep      # Projects with nested compositions
```

### Performance Expectations
```go
// Based on testing with real AEP files:
// Simple files (8-50 items): ~1-2ms parsing time
// Medium files (100-500 items): ~3-5ms parsing time
// Complex files (1000+ items): ~10-20ms parsing time
// Memory usage: Typically under 50MB

func BenchmarkYourFiles(b *testing.B) {
    files := []string{
        "your-simple-project.aep",      // Simple
        "your-medium-project.aep",      // Medium
        "your-complex-project.aep",     // Complex
    }
    
    for _, file := range files {
        b.Run(filepath.Base(file), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                _, err := parser.ParseFile(file)
                if err != nil {
                    b.Fatal(err)
                }
            }
        })
    }
}
```

## Template Management

### Importing Templates

```go
// Import with metadata (you must provide your own AEP file)
template, err := cat.ImportTemplate("your-promo.aep", catalog.ImportOptions{
    Name: "Summer Promo 2025",
    Tags: []string{"promotional", "summer", "social"},
    Metadata: map[string]interface{}{
        "campaign": "summer-2025",
        "duration": "15s",
        "format": "16:9",
    },
})

// Batch import (you must provide your own AEP files)
templates := []string{
    "your-template1.aep",
    "your-template2.aep",
    "your-template3.aep",
}

for _, path := range templates {
    go func(p string) {
        _, err := cat.ImportTemplate(p)
        if err != nil {
            log.Printf("Failed to import %s: %v", p, err)
        }
    }(path)
}
```

### Working with Templates

```go
// Get template by ID
template, err := cat.GetTemplate("template_123")

// Update template metadata
err = cat.UpdateTemplate("template_123", catalog.UpdateOptions{
    Metadata: map[string]interface{}{
        "reviewed": true,
        "approved_by": "john.doe",
        "approved_at": time.Now(),
    },
})

// Delete template
err = cat.DeleteTemplate("template_123")

// List all templates
templates, err := cat.ListTemplates(catalog.ListOptions{
    Limit: 50,
    Offset: 0,
    Sort: "automation_score DESC",
})
```

## Search Operations

### Basic Search

```go
// Keyword search
results, err := cat.Search("text animation", catalog.SearchOptions{
    Type: catalog.SearchTypeKeyword,
    Limit: 20,
})

// Semantic search
results, err := cat.Search("energetic intro", catalog.SearchOptions{
    Type: catalog.SearchTypeSemantic,
    MinScore: 0.7,
})

// Pattern search
results, err := cat.Search("", catalog.SearchOptions{
    Type: catalog.SearchTypePattern,
    Pattern: catalog.PatternTextHeavy,
})
```

### Advanced Search

```go
// Multi-criteria search
query := catalog.AdvancedQuery{
    Conditions: []catalog.Condition{
        {
            Field: "automation_score",
            Operator: ">=",
            Value: 0.8,
        },
        {
            Field: "metadata.duration",
            Operator: "<",
            Value: "30s",
        },
    },
    Logic: "AND",
    Sort: "created_at DESC",
}

results, err := cat.AdvancedSearch(query)

// Search with aggregation
results, err := cat.SearchWithAggregation("motion", catalog.AggregationOptions{
    GroupBy: "metadata.category",
    Metrics: []string{"count", "avg_score"},
})
```

## Automation Workflows

### Creating Simple Workflows

```go
// Text replacement workflow
workflow := &catalog.Workflow{
    Name: "Update Company Name",
    Type: catalog.WorkflowTypeTextReplace,
    Config: catalog.TextReplaceConfig{
        Replacements: map[string]string{
            "{{company}}": "MoBot Inc.",
            "{{year}}": "2025",
            "{{website}}": "www.mobot2025.ai",
        },
    },
    TemplateIDs: []string{"t1", "t2", "t3"},
}

err = orchestrator.CreateWorkflow(workflow)
```

### Multi-Agent Workflows

```go
// Complex workflow with multiple agents
workflow := &catalog.Workflow{
    Name: "Complete Template Optimization",
    Stages: []catalog.Stage{
        {
            Name: "Analysis",
            Agent: "planning",
            Config: map[string]interface{}{
                "depth": "comprehensive",
                "include_dependencies": true,
            },
        },
        {
            Name: "Optimization",
            Agent: "implementation",
            DependsOn: []string{"Analysis"},
            Config: map[string]interface{}{
                "optimize_expressions": true,
                "simplify_effects": true,
            },
        },
        {
            Name: "Validation",
            Agent: "verification",
            DependsOn: []string{"Optimization"},
            Config: map[string]interface{}{
                "run_tests": true,
                "coverage_threshold": 0.8,
            },
        },
    },
}

// Execute workflow
result, err := orchestrator.ExecuteWorkflow(workflow)

// Monitor progress
for {
    status, err := orchestrator.GetWorkflowStatus(workflow.ID)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Progress: %.0f%%\n", status.Progress * 100)
    
    if status.Status == "completed" {
        break
    }
    
    time.Sleep(5 * time.Second)
}
```

### Batch Processing

```go
// Batch process templates
batch := &catalog.BatchJob{
    Name: "Q4 Campaign Update",
    TemplateFilter: catalog.Filter{
        Tags: []string{"q4-campaign"},
        MinScore: 0.7,
    },
    Operations: []catalog.Operation{
        {
            Type: "text_replace",
            Config: map[string]string{
                "{{quarter}}": "Q4 2025",
            },
        },
        {
            Type: "asset_update",
            Config: map[string]string{
                "logo.png": "assets/new-logo-q4.png",
            },
        },
    },
    Options: catalog.BatchOptions{
        Parallel: true,
        Workers: 5,
        ContinueOnError: true,
    },
}

results, err := automation.ProcessBatch(batch)
```

## Agent Integration

### Direct Agent Communication

```go
// Send task to planning agent
message := &catalog.Message{
    Type: "request",
    To: "planning_agent",
    Subject: "analyze_template",
    Payload: map[string]interface{}{
        "template_id": "template_123",
        "analysis_type": "comprehensive",
        "include_opportunities": true,
    },
}

response, err := communication.SendMessage(message)

// Wait for response
result := <-response
fmt.Printf("Analysis complete: %v\n", result.Payload)
```

### Custom Agent Implementation

```go
// Implement custom agent
type CustomAgent struct {
    id    string
    comm  *catalog.AgentCommunicationSystem
}

func (ca *CustomAgent) GetID() string {
    return ca.id
}

func (ca *CustomAgent) HandleMessage(msg *catalog.Message) error {
    switch msg.Subject {
    case "custom_task":
        // Handle custom task
        result := ca.performCustomTask(msg.Payload)
        
        // Send response
        return ca.comm.SendMessage(&catalog.Message{
            Type: "response",
            To: msg.From,
            Subject: "task_complete",
            Payload: result,
        })
    }
    return nil
}

// Register agent
orchestrator.RegisterAgent(&CustomAgent{
    id: "custom_agent_1",
    comm: communication,
})
```

## API Integration

### REST API Examples

```bash
# Import template
curl -X POST http://localhost:8080/api/templates/import \
  -F "file=@template.aep" \
  -F "metadata={\"category\":\"social\"}"

# Search templates
curl "http://localhost:8080/api/templates/search?q=animation&type=semantic"

# Create workflow
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Batch Update",
    "template_ids": [1, 2, 3],
    "config": {
      "text_replacements": {
        "2024": "2025"
      }
    }
  }'

# Get workflow status
curl http://localhost:8080/api/workflows/workflow_123/status

# Export results
curl http://localhost:8080/api/templates/1/export?format=nexrender
```

### Python Client Examples

```python
from mobot2025 import MoBotClient

# Initialize client
client = MoBotClient("http://localhost:8080")

# Import template
template = client.import_template("promo.aep", {
    "name": "Summer Promo",
    "tags": ["summer", "promotional"]
})

# Search templates
results = client.search("text animation", 
    search_type="semantic",
    min_score=0.7
)

# Create workflow
workflow = client.create_workflow({
    "name": "Update Campaign",
    "template_ids": [t.id for t in results],
    "operations": [{
        "type": "text_replace",
        "replacements": {
            "Summer 2024": "Summer 2025"
        }
    }]
})

# Monitor workflow
while True:
    status = client.get_workflow_status(workflow.id)
    print(f"Progress: {status.progress * 100:.0f}%")
    
    if status.completed:
        break
    
    time.sleep(5)
```

### JavaScript/Node.js Examples

```javascript
const { MoBotClient } = require('mobot2025-client');

// Initialize client
const client = new MoBotClient({
    baseURL: 'http://localhost:8080',
    timeout: 30000
});

// Import template
async function importTemplate() {
    const template = await client.templates.import('promo.aep', {
        name: 'Summer Promo',
        metadata: {
            campaign: 'summer-2025'
        }
    });
    
    console.log(`Imported: ${template.name} (Score: ${template.automationScore})`);
}

// Search and process
async function searchAndProcess() {
    // Search for templates
    const results = await client.templates.search('motion graphics', {
        type: 'semantic',
        minScore: 0.8
    });
    
    // Create batch workflow
    const workflow = await client.workflows.create({
        name: 'Batch Motion Update',
        templateIds: results.map(r => r.id),
        agents: ['planning', 'implementation'],
        config: {
            parallel: true,
            qualityThreshold: 0.9
        }
    });
    
    // Monitor progress
    const interval = setInterval(async () => {
        const status = await client.workflows.getStatus(workflow.id);
        console.log(`Progress: ${(status.progress * 100).toFixed(0)}%`);
        
        if (status.status === 'completed') {
            clearInterval(interval);
            console.log('Workflow completed!');
        }
    }, 5000);
}
```

## Advanced Scenarios

### Custom Analysis Pipeline

```go
// Create custom analysis pipeline
pipeline := &catalog.Pipeline{
    Name: "Advanced Template Analysis",
    Stages: []catalog.PipelineStage{
        {
            Name: "Parse",
            Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
                path := data.(string)
                return parser.ParseFile(path)
            },
        },
        {
            Name: "Analyze",
            Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
                project := data.(*catalog.Project)
                return analyzer.AnalyzeProject(project)
            },
        },
        {
            Name: "Score",
            Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
                analysis := data.(*catalog.Analysis)
                return scorer.CalculateScore(analysis)
            },
        },
        {
            Name: "Optimize",
            Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
                scored := data.(*catalog.ScoredAnalysis)
                return optimizer.GenerateOptimizations(scored)
            },
        },
    },
}

// Execute pipeline
result, err := pipeline.Execute(context.Background(), "template.aep")
```

### Real-time Monitoring

```go
// Subscribe to system events
events := orchestrator.Subscribe([]string{
    "workflow.*",
    "agent.*",
    "template.imported",
})

// Handle events
go func() {
    for event := range events {
        switch event.Type {
        case "workflow.started":
            fmt.Printf("Workflow %s started\n", event.WorkflowID)
            
        case "workflow.progress":
            fmt.Printf("Progress: %.0f%%\n", event.Progress * 100)
            
        case "agent.error":
            fmt.Printf("Agent error: %s\n", event.Error)
            
        case "template.imported":
            fmt.Printf("New template: %s\n", event.TemplateName)
        }
    }
}()
```

### Performance Optimization

```go
// Configure for high performance
config := &catalog.Config{
    Parser: catalog.ParserConfig{
        Workers: runtime.NumCPU(),
        BufferSize: 1024 * 1024, // 1MB
        CacheEnabled: true,
    },
    Database: catalog.DatabaseConfig{
        MaxConnections: 50,
        CacheSize: 2000,
        JournalMode: "WAL",
    },
    Agents: catalog.AgentConfig{
        Planning: catalog.AgentWorkerConfig{
            Workers: 5,
            QueueSize: 100,
        },
        Implementation: catalog.AgentWorkerConfig{
            Workers: 3,
            QueueSize: 50,
        },
    },
}

cat, err := catalog.NewCatalogWithConfig("templates.db", config)
```

## Integration Examples

### NexRender Integration

```go
// Export for NexRender
export, err := exporter.ExportForNexRender(template, catalog.NexRenderOptions{
    OutputModule: "h264",
    OutputPath: "s3://bucket/output/",
    Actions: []catalog.NexRenderAction{
        {
            Type: "postrender",
            Module: "@nexrender/action-upload",
            Options: map[string]interface{}{
                "provider": "s3",
                "region": "us-east-1",
            },
        },
    },
})

// Create NexRender job
job := export.ToNexRenderJob()
```

### CI/CD Integration

```yaml
# GitHub Actions example
name: Template Validation

on:
  push:
    paths:
      - 'templates/**/*.aep'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup MoBot
        run: |
          wget https://github.com/mobot2025/releases/latest/mobot-linux
          chmod +x mobot-linux
      
      - name: Validate Templates
        run: |
          for template in templates/*.aep; do
            ./mobot-linux analyze "$template" --min-score 0.7
          done
```

---

For more examples, check the `examples/` directory in the repository or visit our [GitHub Examples](https://github.com/mojosolo/mobot2025/tree/main/examples).