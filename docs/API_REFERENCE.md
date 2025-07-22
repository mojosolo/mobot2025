# 📡 MoBot 2025 API Reference

## Overview

The MoBot 2025 REST API provides comprehensive access to all features including template management, automation scoring, search capabilities, workflow orchestration, and agent control. All endpoints return JSON responses.

**Base URL**: `http://localhost:8080/api`  
**Authentication**: Bearer token (optional)  
**Content-Type**: `application/json` (unless specified)

## Table of Contents

1. [Templates](#templates)
2. [Search](#search)
3. [Workflows](#workflows)
4. [Agents](#agents)
5. [Analysis](#analysis)
6. [System](#system)
7. [Error Handling](#error-handling)
8. [WebSocket Events](#websocket-events)

## Templates

### List Templates

```http
GET /api/templates
```

Query parameters:
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 20, max: 100)
- `sort` (string): Sort field (name, created_at, automation_score)
- `order` (string): Sort order (asc, desc)

Response:
```json
{
  "templates": [
    {
      "id": 1,
      "name": "Text Animation Template",
      "file_path": "/storage/templates/text-anim.aep",
      "file_size": 2048576,
      "automation_score": 0.85,
      "metadata": {
        "version": "2023",
        "compositions": 3,
        "duration": 10.0
      },
      "created_at": "2025-07-21T10:00:00Z",
      "updated_at": "2025-07-21T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "pages": 3
  }
}
```

### Get Template

```http
GET /api/templates/:id
```

Response:
```json
{
  "id": 1,
  "name": "Text Animation Template",
  "file_path": "/storage/templates/text-anim.aep",
  "file_size": 2048576,
  "automation_score": 0.85,
  "metadata": {
    "version": "2023",
    "compositions": 3,
    "duration": 10.0,
    "frame_rate": 30,
    "resolution": "1920x1080"
  },
  "blocks": [
    {
      "type": "CompItem",
      "name": "Main Comp",
      "properties": {}
    }
  ],
  "opportunities": [
    {
      "type": "text_replacement",
      "confidence": 0.95,
      "effort": "low"
    }
  ],
  "created_at": "2025-07-21T10:00:00Z",
  "updated_at": "2025-07-21T10:00:00Z"
}
```

### Import Template

```http
POST /api/templates/import
```

Headers:
- `Content-Type`: `multipart/form-data`

Body:
- `file`: AEP file (required)
- `name`: Template name (optional)
- `metadata`: JSON metadata (optional)

Performance Notes (from real testing):
- Simple AEP files (8 items): ~1ms processing
- Complex AEP files (3,157 items): ~10ms processing
- Memory usage remains under 50MB
- Parser handles unknown blocks gracefully

Response:
```json
{
  "id": 2,
  "name": "Imported Template",
  "status": "processing",
  "import_id": "import_abc123",
  "message": "Template imported successfully and queued for processing"
}
```

### Update Template

```http
PUT /api/templates/:id
```

Body:
```json
{
  "name": "Updated Template Name",
  "metadata": {
    "tags": ["motion", "text"],
    "category": "promotional"
  }
}
```

### Delete Template

```http
DELETE /api/templates/:id
```

Response:
```json
{
  "message": "Template deleted successfully"
}
```

### Get Automation Score

```http
GET /api/templates/:id/automation-score
```

Response:
```json
{
  "template_id": 1,
  "overall_score": 0.85,
  "factors": {
    "text_replaceability": 0.95,
    "asset_modularity": 0.80,
    "effect_simplicity": 0.75,
    "composition_structure": 0.90,
    "expression_complexity": 0.70,
    "plugin_compatibility": 0.85,
    "render_optimization": 0.90
  },
  "recommendations": [
    {
      "type": "optimization",
      "priority": "high",
      "description": "Simplify expressions in layer 'Title'",
      "estimated_improvement": 0.05
    }
  ],
  "automation_ready": true,
  "calculated_at": "2025-07-21T10:30:00Z"
}
```

## Search

### Search Templates

```http
GET /api/templates/search
```

Query parameters:
- `q` (string): Search query (required)
- `type` (string): Search type (keyword, semantic, pattern)
- `min_score` (float): Minimum automation score (0.0-1.0)
- `max_score` (float): Maximum automation score (0.0-1.0)
- `tags` (array): Filter by tags
- `category` (string): Filter by category
- `limit` (int): Results limit (default: 20)

Response:
```json
{
  "results": [
    {
      "id": 1,
      "name": "Text Animation Template",
      "relevance_score": 0.92,
      "automation_score": 0.85,
      "match_context": "Matched in composition names and effects",
      "highlights": [
        "Text <mark>Animation</mark> Template"
      ]
    }
  ],
  "total": 15,
  "query": "animation",
  "search_time_ms": 24
}
```

### Advanced Search

```http
POST /api/templates/search/advanced
```

Body:
```json
{
  "queries": [
    {
      "field": "name",
      "operator": "contains",
      "value": "text"
    },
    {
      "field": "automation_score",
      "operator": "gte",
      "value": 0.8
    }
  ],
  "logic": "AND",
  "sort": {
    "field": "automation_score",
    "order": "desc"
  },
  "limit": 50
}
```

## Workflows

### Create Workflow

```http
POST /api/workflows
```

Body:
```json
{
  "name": "Batch Text Update",
  "template_ids": [1, 2, 3],
  "agents": ["planning", "implementation", "verification"],
  "config": {
    "parallel": true,
    "quality_threshold": 0.9,
    "max_iterations": 30,
    "require_approval": true,
    "timeout_minutes": 60
  },
  "schedule": {
    "type": "once",
    "start_at": "2025-07-22T10:00:00Z"
  }
}
```

Response:
```json
{
  "id": "workflow_xyz789",
  "name": "Batch Text Update",
  "status": "created",
  "created_at": "2025-07-21T11:00:00Z"
}
```

### Get Workflow

```http
GET /api/workflows/:id
```

Response:
```json
{
  "id": "workflow_xyz789",
  "name": "Batch Text Update",
  "status": "running",
  "progress": 0.45,
  "current_stage": "implementation",
  "stages": {
    "planning": {
      "status": "completed",
      "duration_ms": 1200,
      "agent": "planning_agent_1"
    },
    "implementation": {
      "status": "running",
      "started_at": "2025-07-21T11:01:00Z",
      "progress": 0.6
    },
    "verification": {
      "status": "pending"
    }
  },
  "results": {
    "processed": 2,
    "successful": 2,
    "failed": 0,
    "pending": 1
  },
  "created_at": "2025-07-21T11:00:00Z",
  "started_at": "2025-07-21T11:00:30Z"
}
```

### List Workflows

```http
GET /api/workflows
```

Query parameters:
- `status` (string): Filter by status (created, running, completed, failed)
- `agent` (string): Filter by agent involvement
- `created_after` (datetime): Filter by creation date
- `created_before` (datetime): Filter by creation date

### Cancel Workflow

```http
POST /api/workflows/:id/cancel
```

Response:
```json
{
  "id": "workflow_xyz789",
  "status": "cancelling",
  "message": "Workflow cancellation initiated"
}
```

### Retry Workflow

```http
POST /api/workflows/:id/retry
```

Body:
```json
{
  "from_stage": "verification",
  "config_overrides": {
    "quality_threshold": 0.85
  }
}
```

## Agents

### List Agents

```http
GET /api/agents
```

Response:
```json
{
  "agents": [
    {
      "id": "planning_agent_1",
      "type": "planning",
      "status": "idle",
      "health": "healthy",
      "capabilities": ["task_decomposition", "confidence_scoring"],
      "metrics": {
        "tasks_completed": 156,
        "avg_response_time_ms": 1200,
        "success_rate": 0.995
      },
      "last_activity": "2025-07-21T10:55:00Z"
    }
  ]
}
```

### Get Agent Status

```http
GET /api/agents/:id/status
```

Response:
```json
{
  "id": "planning_agent_1",
  "type": "planning",
  "status": "busy",
  "current_task": {
    "id": "task_abc123",
    "type": "analyze_template",
    "progress": 0.75,
    "started_at": "2025-07-21T11:10:00Z"
  },
  "health": {
    "status": "healthy",
    "memory_usage_mb": 256,
    "cpu_usage_percent": 15,
    "uptime_seconds": 3600
  },
  "queue_length": 3
}
```

### Send Agent Command

```http
POST /api/agents/:id/command
```

Body:
```json
{
  "command": "analyze",
  "parameters": {
    "template_id": 1,
    "depth": "comprehensive"
  },
  "priority": "high",
  "timeout_seconds": 300
}
```

Response:
```json
{
  "command_id": "cmd_def456",
  "status": "queued",
  "estimated_completion": "2025-07-21T11:15:00Z"
}
```

## Analysis

### Analyze Template

```http
POST /api/analysis/template/:id
```

Body:
```json
{
  "analysis_types": ["complexity", "opportunities", "dependencies"],
  "include_recommendations": true
}
```

Response:
```json
{
  "template_id": 1,
  "analysis": {
    "complexity": {
      "overall": "medium",
      "factors": {
        "expressions": 12,
        "effects": 24,
        "layers": 36,
        "compositions": 3
      },
      "score": 0.65
    },
    "opportunities": [
      {
        "type": "text_automation",
        "locations": ["Comp 1/Layer 2", "Comp 2/Layer 5"],
        "potential_time_saved": "2 hours",
        "confidence": 0.9
      }
    ],
    "dependencies": {
      "plugins": ["Particular", "Optical Flares"],
      "fonts": ["Arial", "Helvetica Neue"],
      "missing_assets": []
    }
  },
  "recommendations": [
    {
      "priority": "high",
      "category": "optimization",
      "action": "Convert expressions to keyframes where possible",
      "impact": "15% render time improvement"
    }
  ]
}
```

### Batch Analysis

```http
POST /api/analysis/batch
```

Body:
```json
{
  "template_ids": [1, 2, 3, 4, 5],
  "analysis_types": ["automation_score", "complexity"],
  "parallel": true
}
```

## System

### Health Check

```http
GET /api/health
```

Response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime_seconds": 3600,
  "components": {
    "database": "healthy",
    "agents": "healthy",
    "search": "healthy",
    "api": "healthy"
  },
  "timestamp": "2025-07-21T11:20:00Z"
}
```

### System Metrics

```http
GET /api/metrics
```

Response:
```json
{
  "system": {
    "cpu_usage_percent": 25,
    "memory_usage_mb": 512,
    "disk_usage_gb": 2.5,
    "network_io_mbps": 10
  },
  "application": {
    "active_workflows": 3,
    "templates_count": 150,
    "agents_active": 5,
    "api_requests_per_minute": 120
  },
  "performance": {
    "avg_response_time_ms": 45,
    "p95_response_time_ms": 120,
    "error_rate": 0.001,
    "parser_performance": {
      "simple_files_ms": 1,
      "complex_files_ms": 10,
      "max_items_tested": 3157
    }
  }
}
```

### Configuration

```http
GET /api/config
```

Response:
```json
{
  "features": {
    "agents_enabled": true,
    "batch_processing": true,
    "advanced_search": true
  },
  "limits": {
    "max_file_size_mb": 500,
    "max_batch_size": 100,
    "rate_limit_per_minute": 1000
  },
  "version": "1.0.0",
  "environment": "production"
}
```

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Template with ID 999 not found",
    "details": {
      "resource_type": "template",
      "resource_id": 999
    },
    "request_id": "req_ghi789",
    "timestamp": "2025-07-21T11:25:00Z"
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Malformed request or invalid parameters |
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `RESOURCE_NOT_FOUND` | 404 | Requested resource doesn't exist |
| `CONFLICT` | 409 | Resource conflict (e.g., duplicate) |
| `VALIDATION_ERROR` | 422 | Request validation failed |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

## WebSocket Events

Connect to receive real-time updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws');

ws.on('message', (data) => {
  const event = JSON.parse(data);
  console.log(event);
});
```

### Event Types

#### Workflow Updates
```json
{
  "type": "workflow.progress",
  "workflow_id": "workflow_xyz789",
  "data": {
    "progress": 0.75,
    "stage": "verification",
    "message": "Running quality checks"
  },
  "timestamp": "2025-07-21T11:30:00Z"
}
```

#### Agent Status
```json
{
  "type": "agent.status_change",
  "agent_id": "planning_agent_1",
  "data": {
    "old_status": "idle",
    "new_status": "busy",
    "task_id": "task_jkl012"
  },
  "timestamp": "2025-07-21T11:31:00Z"
}
```

#### Template Processing
```json
{
  "type": "template.imported",
  "template_id": 10,
  "data": {
    "name": "New Template",
    "automation_score": 0.88,
    "processing_time_ms": 2500
  },
  "timestamp": "2025-07-21T11:32:00Z"
}
```

## Rate Limiting

API requests are rate limited based on endpoint:

| Endpoint Category | Limit | Window |
|------------------|-------|---------|
| Read operations | 1000 | 1 minute |
| Write operations | 100 | 1 minute |
| Analysis operations | 50 | 1 minute |
| Workflow operations | 20 | 1 minute |

Rate limit headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1627836000
```

## Pagination

All list endpoints support pagination:

```
GET /api/templates?page=2&limit=50
```

Pagination response format:
```json
{
  "data": [...],
  "pagination": {
    "page": 2,
    "limit": 50,
    "total": 245,
    "pages": 5,
    "has_next": true,
    "has_prev": true
  }
}
```

## Authentication (Optional)

If authentication is enabled:

```http
Authorization: Bearer YOUR_API_TOKEN
```

Token endpoint:
```http
POST /api/auth/token
```

Body:
```json
{
  "username": "user",
  "password": "pass"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2025-07-22T11:00:00Z"
}
```

---

For more examples and use cases, see the [Examples](EXAMPLES.md) documentation.