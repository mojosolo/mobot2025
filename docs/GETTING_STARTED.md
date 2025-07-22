# 🚀 Getting Started with MoBot 2025

This guide will help you get MoBot 2025 up and running on your system. Whether you're using it as a library, API service, or complete automation platform, we've got you covered.

## 📋 Prerequisites

### System Requirements
- **OS**: Linux, macOS, or Windows (with WSL)
- **RAM**: Minimum 4GB, recommended 8GB+
- **Disk**: 2GB free space for application and database
- **CPU**: Multi-core processor recommended for agent operations

### Software Requirements
- **Go**: Version 1.19 or higher
- **Python**: Version 3.8+ (for Python bridge)
- **Git**: For cloning the repository
- **SQLite**: Version 3.x (usually pre-installed)

### Optional Requirements
- **curl**: For testing REST API endpoints

## 🔧 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/mobot2025.git
cd mobot2025
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Build the Project

```bash
# Build the main binary
go build -o mobot ./cmd/mobot2025/main.go

# Note: The catalog package has compilation issues that need to be resolved
```

### 4. Verify Installation

```bash
# Check version
./mobot --version

# Run tests
go test ./...
```

## 🏃 Quick Start

### Option 1: Command Line Interface

```bash
# Parse a single AEP file
./mobot parse sample.aep

# Parse and analyze
./mobot analyze sample.aep

# Start interactive mode
./mobot interactive
```

### Option 2: REST API Server

```bash
# Start the API server (default port: 8080)
./mobot api

# Or specify a custom port
./mobot api --port 8090

# With verbose logging
./mobot api --verbose
```

### Option 3: Go Library

Create a new Go file:

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/mobot2025/catalog"
)

func main() {
    // Initialize catalog
    cat, err := catalog.NewCatalog("templates.db")
    if err != nil {
        log.Fatal(err)
    }
    defer cat.Close()

    // Import an AEP file
    template, err := cat.ImportTemplate("sample.aep")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Imported template: %s\n", template.Name)
    fmt.Printf("Automation score: %.2f\n", template.AutomationScore)
}
```

## 🔌 Python Integration

### 1. Install Python Bridge

```bash
cd bridge/python
pip install -e .
```

### 2. Use in Python

```python
from mobot2025 import MoBot

# Initialize MoBot
mobot = MoBot()

# Parse AEP file
result = mobot.parse_aep("sample.aep")
print(f"Found {len(result['blocks'])} blocks")

# Get automation score
score = mobot.calculate_automation_score("sample.aep")
print(f"Automation potential: {score:.2%}")
```

## 🌐 REST API Quick Start

### Starting the Server

```bash
# Start with default settings
./mobot api

# The server will start on http://localhost:8080
```

### Basic API Operations

```bash
# Import a template
curl -X POST http://localhost:8080/api/templates/import \
  -F "file=@sample.aep" \
  -H "Accept: application/json"

# List all templates
curl http://localhost:8080/api/templates

# Get specific template
curl http://localhost:8080/api/templates/1

# Search templates
curl "http://localhost:8080/api/templates/search?q=text+animation"

# Get automation score
curl http://localhost:8080/api/templates/1/automation-score
```

## 🤖 Multi-Agent System

### Enabling Agents

```bash
# Start with agents enabled
./mobot api --enable-agents

# Or set via environment variable
export MOBOT_ENABLE_AGENTS=true
./mobot api
```

### Creating a Workflow

```bash
# Create a new workflow
curl -X POST http://localhost:8080/api/workflows \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": 1,
    "agents": ["planning", "implementation", "verification"],
    "config": {
      "quality_threshold": 0.9,
      "max_iterations": 30
    }
  }'

# Check workflow status
curl http://localhost:8080/api/workflows/1/status
```

## ⚙️ Configuration

### Environment Variables

```bash
# API Configuration
export MOBOT_PORT=8080
export MOBOT_HOST=0.0.0.0
export MOBOT_ENV=production

# Database Configuration
export MOBOT_DB_PATH=./data/templates.db
export MOBOT_DB_TIMEOUT=30s

# Agent Configuration
export MOBOT_ENABLE_AGENTS=true
export MOBOT_AGENT_WORKERS=5
export MOBOT_AGENT_TIMEOUT=5m

# Logging
export MOBOT_LOG_LEVEL=info
export MOBOT_LOG_FORMAT=json
```

### Configuration File

Create `config.yaml`:

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 30s
  write_timeout: 30s

database:
  path: "./data/templates.db"
  max_connections: 10
  timeout: 30s

agents:
  enabled: true
  workers: 5
  timeout: 5m
  retry_attempts: 3

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

## 📊 Database Setup

### Initialize Database

```bash
# Create database directory
mkdir -p data

# Initialize with migrations
./mobot db init

# Run migrations
./mobot db migrate
```

### Import Sample Data

```bash
# Import sample templates
./mobot db seed

# Or import from directory
./mobot import ./samples/
```

## 🐳 Docker Deployment

### Build Docker Image

```bash
docker build -t mobot2025:latest .
```

### Run with Docker

```bash
# Run API server
docker run -p 8080:8080 mobot2025:latest

# With persistent storage
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  mobot2025:latest

# With environment variables
docker run -p 8080:8080 \
  -e MOBOT_ENABLE_AGENTS=true \
  -e MOBOT_LOG_LEVEL=debug \
  mobot2025:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  mobot:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - MOBOT_ENABLE_AGENTS=true
      - MOBOT_ENV=production
    restart: unless-stopped
```

## 🧪 Testing Your Installation

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "agents": "enabled",
  "database": "connected"
}
```

### 2. API Test

```bash
# Create a test file
echo '{"test": "data"}' > test.aep

# Try importing
curl -X POST http://localhost:8080/api/templates/import \
  -F "file=@test.aep"
```

### 3. Agent Test

```bash
# Check agent status
curl http://localhost:8080/api/agents/status
```

## 🔍 Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Change port
   ./mobot api --port 8090
   ```

2. **Database Lock Error**
   ```bash
   # Remove lock file
   rm data/templates.db-journal
   ```

3. **Permission Denied**
   ```bash
   # Make binary executable
   chmod +x mobot
   ```

4. **Missing Dependencies**
   ```bash
   # Reinstall dependencies
   go mod tidy
   go mod download
   ```

### Debug Mode

```bash
# Run with debug logging
./mobot api --debug

# Or set environment variable
export MOBOT_LOG_LEVEL=debug
./mobot api
```

## 📚 Next Steps

Now that you have MoBot 2025 running:

1. **Explore the API**: See [API Reference](API_REFERENCE.md)
2. **Learn about Agents**: Read [Multi-Agent System](MULTI_AGENT_SYSTEM.md)
3. **Import Templates**: Follow the [User Guide](USER_GUIDE.md)
4. **Contribute**: Check out [Developer Guide](DEVELOPER_GUIDE.md)

## 💡 Tips

- Use `--help` with any command for more options
- Check logs in `./logs/` directory for debugging
- The database is automatically created on first run
- API documentation is available at `http://localhost:8080/docs`

## 🆘 Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
2. Search existing [GitHub Issues](https://github.com/yourusername/mobot2025/issues)
3. Join our [Discord Community](https://discord.gg/mobot2025)
4. Create a new issue with:
   - System information
   - Error messages
   - Steps to reproduce

---

Happy automating with MoBot 2025! 🎬✨