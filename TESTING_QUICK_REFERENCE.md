# Testing Quick Reference for mobot2025

## Quick Start

```bash
# Run all tests
./run_all_tests.sh

# Run with performance tests
./run_all_tests.sh --with-performance

# Run with linting
./run_all_tests.sh --with-lint

# Run specific test suite
go test -v ./tests/agents
go test -v ./tests/database
go test -v ./tests/demo
```

## Test Organization

```
tests/
├── agents/           # Multi-agent system tests
├── database/         # Database operation tests
├── demo/            # Demo viewer tests
├── helpers/         # Test utilities and helpers
├── integration/     # End-to-end integration tests
├── python/          # Python bridge tests
└── fixtures/        # Test data files
```

## Writing New Tests

### 1. Unit Test Template
```go
func TestFeatureName(t *testing.T) {
    // Arrange
    agent := catalog.NewAgent()
    input := helpers.CreateMockData()
    
    // Act
    result, err := agent.Process(input)
    
    // Assert
    helpers.AssertNoError(t, err, "Processing failed")
    helpers.AssertEqual(t, expected, result, "Unexpected result")
}
```

### 2. Table-Driven Test Template
```go
func TestMultipleScenarios(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        want    interface{}
        wantErr bool
    }{
        {"valid_input", validData, expectedResult, false},
        {"invalid_input", invalidData, nil, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessData(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 3. Integration Test Template
```go
func TestEndToEndWorkflow(t *testing.T) {
    // Setup
    db := helpers.CreateTestDatabase(t)
    api := catalog.NewAPIService(db)
    
    // Execute workflow
    projectID := uploadTestFile(t, api)
    processProject(t, api, projectID)
    verifyResults(t, api, projectID)
}
```

## Common Test Commands

### Running Tests
```bash
# All tests with verbose output
go test -v ./...

# Specific package
go test -v ./catalog

# Specific test function
go test -v -run TestPlanningAgent ./catalog

# With race detection
go test -race ./...

# With coverage
go test -cover ./...
```

### Coverage Analysis
```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open in browser (macOS)
open coverage.html
```

### Benchmarking
```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkParser ./

# With memory allocation stats
go test -bench=. -benchmem ./...

# Run benchmarks multiple times
go test -bench=. -benchtime=10s ./...
```

## Test Data Management

### Using Test Fixtures
```go
// Load test AEP file
data := helpers.LoadTestFixture(t, "fixtures/valid.aep")

// Use predefined test files
project, err := aep.Open(helpers.TestFixtures.ValidAEP)
```

### Creating Mock Data
```go
// Mock project
project := helpers.CreateMockProject("Test Project")

// Mock compositions
comp := helpers.CreateMockComposition("Main Comp", 1920, 1080, 30)

// Mock text layers
layers := helpers.CreateMockTextLayers(10)
```

## Debugging Tests

### Verbose Output
```bash
# Run with verbose flag
go test -v ./...

# Add custom debug output
t.Logf("Debug: value = %+v", value)
```

### Run Single Test
```bash
# Run one test function
go test -run TestSpecificFunction ./package

# Run tests matching pattern
go test -run "TestAgent.*" ./catalog
```

### Skip Long Tests
```go
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping long test in short mode")
    }
    // Long test logic
}
```

Run with: `go test -short ./...`

## Testing Best Practices

1. **Test Independence**: Each test should run independently
2. **Clean State**: Always clean up resources (use `t.Cleanup()`)
3. **Descriptive Names**: Use clear test names that describe what's being tested
4. **Error Messages**: Provide helpful error messages
5. **Mock External Dependencies**: Don't rely on external services
6. **Test Edge Cases**: Include boundary conditions and error scenarios
7. **Benchmark Important Code**: Add benchmarks for performance-critical sections

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: ./run_all_tests.sh
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage/coverage.out
```

## Troubleshooting

### Common Issues

1. **Test files not found**
   ```bash
   # Ensure you're in the project root
   cd /path/to/mobot2025
   ```

2. **Missing dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Database tests failing**
   - Check SQLite is available
   - Ensure temp directory permissions

4. **Python tests failing**
   ```bash
   # Install Python dependencies
   pip3 install pytest
   ```

## Test Metrics Goals

- **Unit Test Coverage**: 90%+
- **Integration Test Coverage**: 80%+
- **Response Time**: <200ms for API endpoints
- **Benchmark Stability**: <5% variance between runs

## Getting Help

- Check test output for specific errors
- Use `-v` flag for verbose output
- Add `t.Logf()` statements for debugging
- Review test helpers in `tests/helpers/`