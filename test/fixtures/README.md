# Test Fixtures

This directory contains test data generation utilities and fixtures for MoBot 2025.

## Structure

```
test/fixtures/
├── generator.go    # Mock AEP file generator
├── factories.go    # Factory functions for test data
├── helpers.go      # Test environment and assertion helpers
├── aep/           # Generated AEP test files (gitignored)
├── compositions/  # Test composition data
├── layers/        # Test layer data
└── assets/        # Test media asset data
```

## Usage

### Creating Test Data

```go
import "github.com/mojosolo/mobot2025/test/fixtures"

// Create a mock AEP file
gen := fixtures.NewMockAEPGenerator()
aepData, err := gen.GenerateMinimalAEP()

// Create test project metadata
project := fixtures.CreateTestProjectMetadata("test-project")

// Generate a complete test scenario
project, planning, tasks := fixtures.GenerateAutomationScenario()
```

### Test Environment

```go
func TestMyFeature(t *testing.T) {
    // Create test environment
    env := fixtures.NewTestEnvironment(t)
    defer env.Cleanup()
    
    // Create mock AEP file
    aepPath := env.CreateMockAEP("test.aep")
    
    // Run your tests...
}
```

### Test Database

```go
func TestDatabaseFeature(t *testing.T) {
    // Create test database
    testDB := fixtures.NewTestDatabase(t)
    defer testDB.Cleanup()
    
    // Use database path
    db, err := catalog.NewDatabase(testDB.Path())
    // ...
}
```

### Assertions

```go
// Common assertions
fixtures.AssertNoError(t, err, "Failed to parse")
fixtures.AssertEqual(t, expected, actual, "Values don't match")
fixtures.AssertTrue(t, condition, "Condition failed")
```

## Test Scenarios

### 1. Basic Test Project
- Single composition
- Few text layers
- Minimal complexity

### 2. Complex Test Project
- Multiple compositions
- Various layer types
- 3D layers and effects

### 3. Automation Scenario
- Complete workflow
- Task dependencies
- Planning results

### 4. Search Scenario
- Multiple projects
- Different categories
- Varying complexity

### 5. Quality Assurance Scenario
- Projects with issues
- Missing assets
- Poor practices

## Mock AEP Format

The mock AEP generator creates simplified RIFX files with:
- Proper RIFX header
- Basic chunk structure
- Minimal valid data

These files are sufficient for testing parsing logic without requiring real AEP files.

## Best Practices

1. **Use factories for consistency**: Always use factory functions for creating test data
2. **Clean up resources**: Always defer cleanup of test environments and databases
3. **Isolate tests**: Each test should create its own environment
4. **Mock external dependencies**: Use the generators instead of real files
5. **Test edge cases**: Use the various scenarios to test different conditions

## Adding New Fixtures

To add new test fixtures:

1. Add factory functions to `factories.go`
2. Add generator methods to `generator.go`
3. Update test scenarios as needed
4. Document new fixtures in this README