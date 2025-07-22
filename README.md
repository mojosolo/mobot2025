# 🎬 MoBot 2025 - AI-Powered Video Production Pipeline

<div align="center">
  <img src="https://img.shields.io/badge/Phase%201-Complete-green?style=for-the-badge" alt="Phase 1 Complete">
  <img src="https://img.shields.io/badge/Phase%202-Complete-green?style=for-the-badge" alt="Phase 2 Complete">
  <img src="https://img.shields.io/badge/Go-1.19+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
</div>

## 🚨 Current Status

**✅ All compilation errors have been resolved!** The project now builds successfully. Both the core parser and the advanced catalog features are fully functional.

## 🚀 Overview

**MoBot 2025** is an advanced AI-powered system for analyzing, cataloging, and automating Adobe After Effects (AEP) project files. It combines sophisticated parsing capabilities with a multi-agent AI orchestration system to enable intelligent video production workflows at scale.

### 🎯 Key Features

- **🔍 Deep AEP Analysis**: Parse and analyze After Effects project files with RIFX format support
- **🤖 Multi-Agent AI System**: 5-agent orchestration for intelligent automation
- **📊 Automation Scoring**: Evaluate templates for automation potential (0-100 score)
- **⚡ Advanced Search**: Vector-based semantic search with quality filtering
- **🔄 Workflow Automation**: End-to-end template processing pipelines
- **📈 Quality Assurance**: Pattern matching and anti-pattern detection
- **🏗️ Batch Processing**: Concurrent processing of multiple templates
- **📡 REST API**: Production-ready API for all features
- **💾 Database Support**: SQLite (local) or PostgreSQL with Neon (production)
- **☁️ S3 Storage**: Optional AWS S3 integration for AEP file storage
- **🔐 GitHub Secrets**: Secure credential management for production deployments

## 📋 Quick Start

```bash
# Clone the repository
git clone https://github.com/mojosolo/mobot2025.git
cd mobot2025

# Build the project
go build -o mobot ./cmd/mobot2025/main.go

# Start the REST API server
./mobot serve

# Or use as a Go library
go get github.com/mojosolo/mobot2025
```

### Verify Your Setup

```bash
# Run the verification script
./verify-setup.sh
```

### ⚠️ Important: Test Data Not Included

**AEP files and other binary test data are NOT included in this repository.** See [TEST_DATA_README.md](TEST_DATA_README.md) for instructions on obtaining test files.

### Working Example

```bash
# Parse an AEP file (you must provide your own AEP file)
./mobot parse -file your-project.aep

# Analyze an AEP file (you must provide your own AEP file)
./mobot analyze -file your-project.aep

# Import template to catalog (you must provide your own AEP file)
./mobot import your-project.aep
```

For detailed setup instructions, see [Getting Started Guide](docs/GETTING_STARTED.md).

## 🏗️ Architecture

MoBot 2025 features a sophisticated multi-layer architecture:

```
┌─────────────────────────────────────────────────────────┐
│                    REST API Layer                        │
├─────────────────────────────────────────────────────────┤
│              Multi-Agent Orchestration                   │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐    │
│  │Planning │ │Implement │ │Verify    │ │Review   │    │
│  │Agent    │ │Agent     │ │Agent     │ │Agent    │    │
│  └────┬────┘ └────┬─────┘ └────┬─────┘ └────┬────┘    │
│       └───────────┴──────────┬──┴────────────┘         │
│                     ┌────────┴────────┐                 │
│                     │Meta-Orchestrator│                 │
│                     └────────┬────────┘                 │
├──────────────────────────────┴──────────────────────────┤
│                    Catalog System                        │
│  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐ │
│  │Database │ │Search    │ │Scoring   │ │Quality     │ │
│  │         │ │Engine    │ │Engine    │ │Assurance   │ │
│  └─────────┘ └──────────┘ └──────────┘ └────────────┘ │
├─────────────────────────────────────────────────────────┤
│                  Core Parser Engine                      │
│            (RIFX Format Support, Block Types)           │
└─────────────────────────────────────────────────────────┘
```

For complete architecture details, see [Architecture Documentation](docs/ARCHITECTURE.md).

## 🔧 Core Components

### Phase 1: Foundation (✅ Complete)
- **AEP Parser**: Advanced RIFX format parsing with 16+ block types
- **Database System**: SQLite with migration support
- **Search Engine**: Semantic search with multiple index types
- **Automation Scoring**: 7-factor weighted scoring algorithm
- **Import System**: Legacy template migration
- **REST API**: Comprehensive endpoint coverage

### Phase 2: Multi-Agent System (✅ Complete)
- **Planning Agent**: Task decomposition and analysis
- **Implementation Agent**: Code generation with model cascading
- **Verification Agent**: Automated testing and validation
- **Review Agent**: Performance optimization recommendations
- **Meta-Orchestrator**: Workflow coordination and state management
- **Communication Protocol**: Inter-agent messaging system
- **Workflow Automation**: Pipeline-based batch processing
- **Quality Assurance**: Pattern matching and validation
- **System Integration**: Comprehensive testing framework

## 📚 Documentation

- 📖 [Getting Started](docs/GETTING_STARTED.md) - Installation and setup
- 🏗️ [Architecture](docs/ARCHITECTURE.md) - System design and components
- 🤖 [Multi-Agent System](docs/MULTI_AGENT_SYSTEM.md) - AI orchestration details
- 📡 [API Reference](docs/API_REFERENCE.md) - Complete API documentation
- 👤 [User Guide](docs/USER_GUIDE.md) - How to use MoBot 2025
- 👨‍💻 [Developer Guide](docs/DEVELOPER_GUIDE.md) - Contributing and extending
- ⚙️ [Configuration](docs/CONFIGURATION.md) - Settings and options
- 🎯 [Examples](docs/EXAMPLES.md) - Code examples and use cases
- 🔧 [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues

### Sprint Reports
- 📊 [Phase 1 Completion](docs/PHASE_1_COMPLETION.md) - Foundation implementation
- 📊 [Phase 2 Completion](docs/PHASE_2_COMPLETION.md) - Multi-agent system

## 🚦 Project Status

| Component | Status | Description |
|-----------|--------|-------------|
| Core Parser | ✅ Complete | RIFX format parsing with 16+ block types |
| Database | ✅ Complete | SQLite with migrations and indexing |
| Search Engine | ✅ Complete | Semantic search with quality filtering |
| REST API | ✅ Complete | Full CRUD operations and advanced queries |
| Multi-Agent System | ✅ Complete | 5-agent orchestration platform |
| Workflow Automation | ✅ Complete | Pipeline-based batch processing |
| Quality Assurance | ✅ Complete | Pattern matching and validation |
| Documentation | ✅ Complete | Comprehensive documentation |

## 💻 Usage Examples

### As a Go Library

```go
import "github.com/mojosolo/mobot2025/catalog"

// Create a new catalog
cat := catalog.NewCatalog("templates.db")

// Import an AEP file (you must provide your own AEP file)
template, err := cat.ImportTemplate("your-project.aep")

// Search for templates
results := cat.Search("motion graphics", catalog.SearchOptions{
    Type: "semantic",
    MinScore: 0.8,
})

// Get automation score
score := cat.CalculateAutomationScore(template)
```

### REST API

```bash
# Import a template (you must provide your own AEP file)
curl -X POST http://localhost:8080/api/templates/import \
  -F "file=@your-project.aep"

# Search templates
curl "http://localhost:8080/api/templates/search?q=motion+graphics&type=semantic"

# Get automation recommendations
curl "http://localhost:8080/api/templates/123/automation-score"
```

### Python Integration

```python
from mobot2025 import Catalog

# Initialize catalog
catalog = Catalog("templates.db")

# Analyze template (you must provide your own AEP file)
template = catalog.import_template("your-project.aep")
score = catalog.get_automation_score(template.id)

# Batch process templates (you must provide your own AEP files)
results = catalog.batch_process([
    "your-template1.aep",
    "your-template2.aep",
    "your-template3.aep"
])
```

## ☁️ Cloud Deployment

MoBot 2025 supports production deployment with cloud services:

### Neon PostgreSQL Database
- Serverless PostgreSQL with auto-scaling
- Database branching for development/staging
- Point-in-time recovery

### AWS S3 Storage
- Secure file storage for AEP projects
- Presigned URLs for direct downloads
- Automatic file organization by date

### Configuration
```bash
# For production with Neon + S3
export MOBOT_DB_TYPE=postgres
export NEON_DATABASE_URL=postgresql://user:pass@host/db?sslmode=require
export AWS_S3_ENABLED=true
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
export AWS_BUCKET=mobot2025-storage
```

See [Deployment Guide](docs/DEPLOYMENT.md) for detailed instructions.

## 🛠️ Development

### Prerequisites
- Go 1.19 or higher
- Python 3.8+ (for Python bridge)
- SQLite 3 (for local development)
- PostgreSQL client (for Neon deployment)
- AWS CLI (for S3 deployment)
- Git

## 📁 Project Structure

```
mobot2025/
├── catalog/          # Core business logic and agents
├── cmd/              # Command-line entry points
├── data/             # Directory for test data (AEP files NOT included)
├── demo/             # Demo applications and viewers
├── docs/             # Documentation
├── enhancements/     # Enhancement modules
├── sample-aep/       # Directory for sample projects (AEP files NOT included)
├── sandbox/          # Temporary files and test scripts (not for production)
│   ├── scripts/      # Shell/Python scripts used during development
│   ├── tests/        # One-off test files
│   ├── temp/         # Temporary and backup files
│   └── reports/      # Generated HTML reports
└── tests/            # Organized test suite
```

**Note**: The `data/` and `sample-aep/` directories are placeholders. You must provide your own AEP files for testing. Binary files (AEP, video, PDF, etc.) are excluded from version control to keep the repository lightweight.

### Building from Source

```bash
# Clone repository
git clone https://github.com/mojosolo/mobot2025.git
cd mobot2025

# Install dependencies
go mod download

# Build
go build -o mobot ./cmd/mobot2025/main.go

# Run tests
go test ./...

# Run with race detector
go test -race ./...
```

For detailed development instructions, see [Developer Guide](docs/DEVELOPER_GUIDE.md).

## 🤝 Contributing

We welcome contributions! Please see our [Developer Guide](docs/DEVELOPER_GUIDE.md) for:
- Code style guidelines
- Testing requirements
- Pull request process
- Architecture decisions

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Adobe After Effects for the AEP format
- The Go community for excellent libraries
- Contributors to the multi-agent system design

## 📞 Support

- 📧 Email: support@mobot2025.ai
- 💬 Discord: [Join our community](https://discord.gg/mobot2025)
- 🐛 Issues: [GitHub Issues](https://github.com/mojosolo/mobot2025/issues)

---

<div align="center">
  <strong>Built with ❤️ for the video production community</strong>
</div>