# ⚙️ MoBot 2025 Configuration Guide

## Overview

MoBot 2025 can be configured through environment variables, configuration files, and command-line flags. This guide covers all available configuration options and their effects.

## Configuration Hierarchy

Configuration sources are applied in the following order (later sources override earlier ones):

1. Default values (built-in)
2. Configuration file (`config.yaml` or `config.json`)
3. Environment variables
4. Command-line flags

## Configuration File

### Location

MoBot looks for configuration files in these locations (in order):

1. `./config.yaml` (current directory)
2. `$HOME/.mobot/config.yaml`
3. `/etc/mobot/config.yaml`
4. Path specified by `--config` flag

### Format

Configuration files can be in YAML or JSON format:

#### YAML Example (`config.yaml`)
```yaml
# Server Configuration
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  shutdown_timeout: 10s
  cors:
    enabled: true
    origins: ["*"]
    methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]

# Database Configuration
database:
  path: "./data/templates.db"
  max_connections: 25
  idle_connections: 5
  connection_timeout: 30s
  migration_path: "./migrations"
  backup:
    enabled: true
    interval: 24h
    retention: 7

# Parser Configuration
parser:
  max_file_size: 524288000  # 500MB
  timeout: 5m
  workers: 4
  cache:
    enabled: true
    size: 100
    ttl: 1h

# Agent Configuration
agents:
  enabled: true
  registry:
    planning:
      workers: 3
      timeout: 5m
      max_retries: 3
    implementation:
      workers: 2
      timeout: 10m
      models: ["claude", "gpt-4", "gemini"]
    verification:
      workers: 5
      timeout: 3m
      coverage_threshold: 0.8
    review:
      workers: 2
      timeout: 5m
    orchestrator:
      max_workflows: 50
      loop_limit: 50
      checkpoint_interval: 1m

# Search Configuration
search:
  engine: "sqlite_fts5"  # Options: sqlite_fts5, elasticsearch
  min_score: 0.1
  max_results: 100
  timeout: 10s
  cache:
    enabled: true
    ttl: 5m

# API Configuration
api:
  rate_limiting:
    enabled: true
    requests_per_minute: 1000
    burst: 50
  authentication:
    enabled: false
    jwt_secret: "your-secret-key"
    token_expiry: 24h
  pagination:
    default_limit: 20
    max_limit: 100

# Logging Configuration
logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, text
  output: "stdout"  # stdout, file
  file:
    path: "./logs/mobot.log"
    max_size: 100  # MB
    max_backups: 5
    max_age: 30  # days
    compress: true

# Monitoring Configuration
monitoring:
  metrics:
    enabled: true
    port: 9090
    path: "/metrics"
  health_check:
    enabled: true
    interval: 30s
    timeout: 5s
  tracing:
    enabled: false
    endpoint: "http://localhost:14268/api/traces"
    sample_rate: 0.1

# Storage Configuration
storage:
  templates:
    path: "./storage/templates"
    max_size: 10737418240  # 10GB
  temp:
    path: "./storage/temp"
    cleanup_interval: 1h
    max_age: 24h
  exports:
    path: "./storage/exports"
    retention: 168h  # 7 days

# Workflow Configuration
workflow:
  batch:
    max_size: 100
    concurrent_jobs: 10
    retry_attempts: 3
    retry_delay: 1m
  scheduling:
    enabled: true
    timezone: "UTC"
  notifications:
    enabled: false
    webhook_url: ""
    events: ["completed", "failed"]

# Security Configuration
security:
  file_upload:
    allowed_extensions: [".aep", ".aepx"]
    scan_enabled: true
    quarantine_path: "./quarantine"
  api:
    csrf_protection: true
    request_id_header: "X-Request-ID"
  encryption:
    enabled: false
    key_path: "./keys/master.key"

# Development Configuration
development:
  debug: false
  pretty_logs: true
  reload_templates: true
  mock_agents: false
  test_mode: false
```

#### JSON Example (`config.json`)
```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": "30s",
    "write_timeout": "30s"
  },
  "database": {
    "path": "./data/templates.db",
    "max_connections": 25
  },
  "agents": {
    "enabled": true,
    "registry": {
      "planning": {
        "workers": 3,
        "timeout": "5m"
      }
    }
  }
}
```

## Environment Variables

All configuration options can be set via environment variables using the prefix `MOBOT_`:

### Server Variables
```bash
MOBOT_SERVER_HOST=0.0.0.0
MOBOT_SERVER_PORT=8080
MOBOT_SERVER_READ_TIMEOUT=30s
MOBOT_SERVER_WRITE_TIMEOUT=30s
MOBOT_SERVER_SHUTDOWN_TIMEOUT=10s
MOBOT_SERVER_CORS_ENABLED=true
MOBOT_SERVER_CORS_ORIGINS=*
```

### Database Variables
```bash
MOBOT_DATABASE_PATH=./data/templates.db
MOBOT_DATABASE_MAX_CONNECTIONS=25
MOBOT_DATABASE_IDLE_CONNECTIONS=5
MOBOT_DATABASE_CONNECTION_TIMEOUT=30s
MOBOT_DATABASE_BACKUP_ENABLED=true
MOBOT_DATABASE_BACKUP_INTERVAL=24h
```

### Agent Variables
```bash
MOBOT_AGENTS_ENABLED=true
MOBOT_AGENTS_PLANNING_WORKERS=3
MOBOT_AGENTS_PLANNING_TIMEOUT=5m
MOBOT_AGENTS_IMPLEMENTATION_MODELS=claude,gpt-4,gemini
MOBOT_AGENTS_VERIFICATION_COVERAGE_THRESHOLD=0.8
MOBOT_AGENTS_ORCHESTRATOR_MAX_WORKFLOWS=50
MOBOT_AGENTS_ORCHESTRATOR_LOOP_LIMIT=50
```

### API Variables
```bash
MOBOT_API_RATE_LIMITING_ENABLED=true
MOBOT_API_RATE_LIMITING_REQUESTS_PER_MINUTE=1000
MOBOT_API_AUTHENTICATION_ENABLED=false
MOBOT_API_AUTHENTICATION_JWT_SECRET=your-secret-key
MOBOT_API_PAGINATION_DEFAULT_LIMIT=20
MOBOT_API_PAGINATION_MAX_LIMIT=100
```

### Logging Variables
```bash
MOBOT_LOGGING_LEVEL=info
MOBOT_LOGGING_FORMAT=json
MOBOT_LOGGING_OUTPUT=stdout
MOBOT_LOGGING_FILE_PATH=./logs/mobot.log
MOBOT_LOGGING_FILE_MAX_SIZE=100
```

### Development Variables
```bash
MOBOT_DEVELOPMENT_DEBUG=false
MOBOT_DEVELOPMENT_PRETTY_LOGS=true
MOBOT_DEVELOPMENT_RELOAD_TEMPLATES=true
MOBOT_DEVELOPMENT_MOCK_AGENTS=false
MOBOT_DEVELOPMENT_TEST_MODE=false
```

## Command-Line Flags

Command-line flags override both configuration files and environment variables:

### Global Flags
```bash
mobot --config ./custom-config.yaml
mobot --log-level debug
mobot --env production
```

### Server Flags
```bash
mobot api --port 8090
mobot api --host 127.0.0.1
mobot api --cors-enabled=false
mobot api --rate-limit 500
```

### Parser Flags
```bash
mobot parse --timeout 10m
mobot parse --workers 8
mobot parse --max-size 1GB
```

### Agent Flags
```bash
mobot api --enable-agents
mobot api --agent-workers 10
mobot api --disable-agent planning
```

## Configuration Profiles

Use profiles for different environments:

### Development Profile
```yaml
# config.dev.yaml
extends: config.yaml
server:
  port: 8080
logging:
  level: debug
  format: text
development:
  debug: true
  pretty_logs: true
  mock_agents: true
```

### Production Profile
```yaml
# config.prod.yaml
extends: config.yaml
server:
  port: 80
  cors:
    origins: ["https://app.example.com"]
logging:
  level: warn
  format: json
api:
  authentication:
    enabled: true
  rate_limiting:
    requests_per_minute: 100
monitoring:
  metrics:
    enabled: true
  tracing:
    enabled: true
```

### Testing Profile
```yaml
# config.test.yaml
extends: config.yaml
database:
  path: ":memory:"
development:
  test_mode: true
  mock_agents: true
logging:
  level: error
```

## Advanced Configuration

### Multi-Model Configuration
```yaml
agents:
  implementation:
    models:
      claude:
        api_key: "${CLAUDE_API_KEY}"
        max_tokens: 4096
        temperature: 0.7
      gpt-4:
        api_key: "${OPENAI_API_KEY}"
        max_tokens: 8192
        temperature: 0.5
      gemini:
        api_key: "${GEMINI_API_KEY}"
        max_tokens: 2048
        temperature: 0.8
```

### External Service Integration
```yaml
integrations:
  nexrender:
    enabled: true
    api_url: "http://nexrender.local"
    api_key: "${NEXRENDER_API_KEY}"
  storage:
    provider: "s3"  # s3, gcs, azure
    bucket: "mobot-templates"
    region: "us-east-1"
    credentials:
      access_key: "${AWS_ACCESS_KEY_ID}"
      secret_key: "${AWS_SECRET_ACCESS_KEY}"
```

### Custom Agent Configuration
```yaml
agents:
  custom:
    my_agent:
      type: "custom"
      executable: "./agents/my_agent"
      args: ["--mode", "production"]
      env:
        CUSTOM_VAR: "value"
      health_check:
        endpoint: "http://localhost:9000/health"
        interval: 30s
```

## Configuration Validation

MoBot validates configuration on startup:

```bash
# Validate configuration without starting
mobot validate --config config.yaml

# Test configuration
mobot test-config --config config.yaml
```

### Validation Rules

1. **Required Fields**: Certain fields must be present
2. **Type Checking**: Values must match expected types
3. **Range Validation**: Numeric values must be within valid ranges
4. **Path Validation**: File paths must be accessible
5. **Dependency Checking**: Related settings must be compatible

## Dynamic Configuration

Some settings can be changed at runtime:

### Via API
```bash
# Update log level
curl -X PUT http://localhost:8080/api/config \
  -H "Content-Type: application/json" \
  -d '{"logging": {"level": "debug"}}'

# Update rate limiting
curl -X PUT http://localhost:8080/api/config \
  -H "Content-Type: application/json" \
  -d '{"api": {"rate_limiting": {"requests_per_minute": 500}}}'
```

### Via Admin Interface
Access `http://localhost:8080/admin/config` for web-based configuration.

## Performance Tuning

### Database Optimization
```yaml
database:
  max_connections: 50  # Increase for high load
  cache_size: 2000     # Pages in cache
  journal_mode: "WAL"  # Write-ahead logging
  synchronous: "NORMAL"
  temp_store: "MEMORY"
```

### Agent Performance
```yaml
agents:
  planning:
    workers: 10        # More workers for parallel processing
    buffer_size: 100   # Message buffer
    timeout: 2m        # Shorter timeout for responsiveness
```

### API Performance
```yaml
api:
  max_request_size: 104857600  # 100MB
  timeout: 5m
  keep_alive: true
  compression:
    enabled: true
    level: 5
```

## Security Configuration

### TLS/SSL
```yaml
server:
  tls:
    enabled: true
    cert_path: "./certs/server.crt"
    key_path: "./certs/server.key"
    client_auth: false
    ca_path: "./certs/ca.crt"
```

### Authentication
```yaml
api:
  authentication:
    enabled: true
    providers:
      jwt:
        secret: "${JWT_SECRET}"
        issuer: "mobot2025"
        audience: ["api.mobot2025.com"]
      oauth:
        enabled: true
        providers:
          - name: "google"
            client_id: "${GOOGLE_CLIENT_ID}"
            client_secret: "${GOOGLE_CLIENT_SECRET}"
```

## Monitoring and Alerts

### Prometheus Metrics
```yaml
monitoring:
  metrics:
    enabled: true
    port: 9090
    path: "/metrics"
    namespace: "mobot"
    subsystem: "api"
    buckets: [0.1, 0.5, 1, 2.5, 5, 10]
```

### Alert Configuration
```yaml
monitoring:
  alerts:
    enabled: true
    rules:
      - name: "high_error_rate"
        condition: "error_rate > 0.05"
        duration: "5m"
        severity: "critical"
      - name: "slow_response"
        condition: "p95_latency > 1s"
        duration: "10m"
        severity: "warning"
```

## Troubleshooting Configuration

### Debug Configuration Issues
```bash
# Show effective configuration
mobot config show

# Show configuration sources
mobot config sources

# Test specific configuration
mobot config test database
```

### Common Issues

1. **Port Already in Use**
   ```yaml
   server:
     port: 8081  # Change to available port
   ```

2. **Database Lock**
   ```yaml
   database:
     journal_mode: "DELETE"  # Less concurrent but avoids locks
   ```

3. **Memory Issues**
   ```yaml
   parser:
     workers: 2  # Reduce workers
     cache:
       size: 50  # Smaller cache
   ```

## Best Practices

1. **Use Environment Variables for Secrets**
   ```yaml
   api:
     authentication:
       jwt_secret: "${JWT_SECRET}"  # Not hardcoded
   ```

2. **Profile-Based Configuration**
   ```bash
   mobot api --profile production
   ```

3. **Version Control**
   - Include `config.example.yaml`
   - Exclude actual config files with secrets
   - Document all custom settings

4. **Regular Backups**
   ```yaml
   database:
     backup:
       enabled: true
       schedule: "0 2 * * *"  # 2 AM daily
   ```

---

For more details on specific configurations, see the relevant guides:
- [API Reference](API_REFERENCE.md) for API-specific settings
- [Multi-Agent System](MULTI_AGENT_SYSTEM.md) for agent configuration
- [Developer Guide](DEVELOPER_GUIDE.md) for development settings