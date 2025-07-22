# 🧪 MoBot 2025 Testing Guide

This guide covers testing best practices, strategies, and tools for the MoBot 2025 project.

## Table of Contents

1. [Testing Philosophy](#testing-philosophy)
2. [Test Structure](#test-structure)
3. [Writing Tests](#writing-tests)
4. [Test Data](#test-data)
5. [Running Tests](#running-tests)
6. [Code Coverage](#code-coverage)
7. [CI/CD Integration](#cicd-integration)
8. [Troubleshooting](#troubleshooting)

## Testing Philosophy

### Core Principles

1. **Test Behavior, Not Implementation**: Focus on what the code does, not how
2. **Isolation**: Each test should be independent and repeatable
3. **Clarity**: Test names should clearly describe what is being tested
4. **Speed**: Unit tests should run quickly (< 100ms per test)
5. **Coverage**: Aim for 80%+ coverage, 100% for critical paths

### Testing Pyramid

```
        /\
       /  \  E2E Tests (5%)
      /----\
     /      \ Integration Tests (20%)
    /--------\
   /          \ Unit Tests (75%)
  /____________\
```

## Test Structure

### Directory Layout

```
mobot2025/
├── *_test.go              # Unit tests (next to source files)
├── tests/
│   ├── integration/       # Integration tests
│   ├── e2e/              # End-to-end tests
│   └── helpers/          # Shared test utilities
└── test/
    └── fixtures/         # Test data generators
```

### Test File Naming

- Unit tests: `<file>_test.go` (same package)
- External tests: `<package>_test.go` (separate package)
- Integration tests: `tests/integration/<feature>_test.go`
- Benchmarks: `<file>_bench_test.go`

## Writing Tests

### Unit Test Template

```go
package mypackage_test

import (
    "testing"
    "github.com/mojosolo/mobot2025/test/fixtures"
)

func TestFeatureName_WhenCondition_ShouldExpectedBehavior(t *testing.T) {
    // Arrange
    env := fixtures.NewTestEnvironment(t)
    defer env.Cleanup()
    
    // Act
    result, err := myFunction()
    
    // Assert
    fixtures.AssertNoError(t, err, "Unexpected error")
    fixtures.AssertEqual(t, expected, result, "Result mismatch")
}
```

### Table-Driven Tests

```go
func TestParser_ParseMultipleFormats(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected Output
        wantErr  bool
    }{
        {
            name:     "valid RIFX format",
            input:    "RIFX...",
            expected: Output{...},
            wantErr:  false,
        },
        {
            name:     "invalid header",
            input:    "INVALID",
            expected: Output{},
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Parse(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("Parse() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Integration Test Example

```go
func TestWorkflow_CompleteAutomation(t *testing.T) {
    // Skip in short mode
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup
    db := fixtures.NewTestDatabase(t)
    defer db.Cleanup()
    
    catalog, err := catalog.NewDatabase(db.Path())
    fixtures.AssertNoError(t, err, "Failed to create catalog")
    defer catalog.Close()
    
    // Create test scenario
    project, planning, tasks := fixtures.GenerateAutomationScenario()
    
    // Execute workflow
    orchestrator := catalog.NewMetaOrchestrator(catalog)
    workflow, err := orchestrator.ExecuteWorkflow(project, planning, tasks)
    
    // Verify results
    fixtures.AssertNoError(t, err, "Workflow execution failed")
    fixtures.AssertEqual(t, "completed", workflow.Status, "Workflow not completed")
}
```

### Benchmark Example

```go
func BenchmarkParser_LargeFile(b *testing.B) {
    // Setup
    env := fixtures.NewBenchmarkEnvironment(b)
    defer env.Cleanup()
    
    // Generate large test file
    gen := fixtures.NewMockAEPGenerator()
    data, _ := gen.GenerateComplexAEP(100, 50)
    
    // Reset timer after setup
    b.ResetTimer()
    
    // Run benchmark
    for i := 0; i < b.N; i++ {
        _, err := Parse(data)
        if err != nil {
            b.Fatal(err)
        }
    }
    
    // Report allocations
    b.ReportAllocs()
}
```

## Test Data

### Using Test Fixtures

```go
// Generate mock AEP file
gen := fixtures.NewMockAEPGenerator()
aepData, err := gen.GenerateMinimalAEP()

// Create test project metadata
project := fixtures.CreateTestProjectMetadata("test-project")

// Generate complete scenario
project, planning, tasks := fixtures.GenerateAutomationScenario()
```

### Test Environment

```go
func TestFileOperations(t *testing.T) {
    env := fixtures.NewTestEnvironment(t)
    defer env.Cleanup()
    
    // Create test files
    aepPath := env.CreateMockAEP("test.aep")
    
    // Use test files
    result, err := ProcessFile(aepPath)
    
    // Assertions
    env.AssertFileExists("output.json")
}
```

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./catalog

# Run specific test
go test -run TestParser ./parser

# Run with race detector
go test -race ./...

# Run short tests only
go test -short ./...
```

### Coverage Commands

```bash
# Run with coverage
go test -cover ./...

# Generate coverage report
./scripts/coverage.sh

# Generate detailed reports
./scripts/coverage-report.sh

# Open HTML coverage
./scripts/coverage-report.sh --open
```

### Integration Tests

```bash
# Run integration tests
go test -tags=integration ./tests/integration/...

# Run with timeout
go test -timeout 30m ./tests/integration/...
```

## Code Coverage

### Coverage Requirements

- **Overall**: Minimum 80% coverage
- **Critical Paths**: 100% coverage required for:
  - Parser core functionality
  - Database operations
  - Agent coordination
  - Error handling

### Checking Coverage

```bash
# Check coverage meets threshold
./scripts/coverage.sh

# View uncovered code
go tool cover -html=coverage.out
```

### Improving Coverage

1. Focus on error paths
2. Test edge cases
3. Cover all branches
4. Test concurrent scenarios
5. Verify error messages

## CI/CD Integration

### GitHub Actions

Tests run automatically on:
- Push to main/develop branches
- Pull requests
- Tagged releases

### Local CI Simulation

```bash
# Run tests like CI
go test -v -race -coverprofile=coverage.out ./...

# Run linter
golangci-lint run

# Build all platforms
GOOS=linux go build ./...
GOOS=darwin go build ./...
GOOS=windows go build ./...
```

## Troubleshooting

### Common Issues

#### 1. Tests Hanging

**Symptom**: Tests don't complete
**Solution**: Add timeouts to tests

```go
func TestLongRunning(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Use ctx in your test
}
```

#### 2. Flaky Tests

**Symptom**: Tests pass/fail randomly
**Common Causes**:
- Race conditions
- Dependency on external state
- Time-based logic

**Solution**: Use deterministic testing

```go
// Bad: Time-dependent
if time.Now().Hour() > 12 {
    // ...
}

// Good: Inject time
func doWork(now time.Time) {
    if now.Hour() > 12 {
        // ...
    }
}
```

#### 3. Database Lock Errors

**Symptom**: "database is locked" errors
**Solution**: Use separate database per test

```go
func TestDatabase(t *testing.T) {
    db := fixtures.NewTestDatabase(t)
    defer db.Cleanup()
    // Each test gets its own database
}
```

### Debugging Tests

```bash
# Run single test with verbose output
go test -v -run TestSpecific ./package

# Use delve debugger
dlv test ./package -- -test.run TestSpecific

# Print coverage for specific function
go tool cover -func=coverage.out | grep FunctionName

# Check for race conditions
go test -race -run TestSpecific ./package
```

## Best Practices

### DO:
- ✅ Write tests first (TDD)
- ✅ Keep tests simple and focused
- ✅ Use descriptive test names
- ✅ Test edge cases
- ✅ Mock external dependencies
- ✅ Use test fixtures for consistency
- ✅ Clean up resources with defer
- ✅ Run tests in parallel when possible

### DON'T:
- ❌ Test private functions directly
- ❌ Depend on test execution order
- ❌ Use global state in tests
- ❌ Skip error checking in tests
- ❌ Write tests that take > 1 second
- ❌ Commit failing tests
- ❌ Use production credentials
- ❌ Test implementation details

## Testing Checklist

Before submitting PR:

- [ ] All tests pass locally
- [ ] Coverage meets 80% threshold
- [ ] No race conditions detected
- [ ] Integration tests pass
- [ ] New features have tests
- [ ] Error paths are tested
- [ ] Documentation updated
- [ ] CI build is green

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Assertions](https://github.com/stretchr/testify)
- [Go Test Patterns](https://github.com/golang/go/wiki/TestComments)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)