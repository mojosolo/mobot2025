# Real Data Testing Guide for MoBot 2025

## Overview

This guide explains how to test MoBot 2025 using ONLY real AEP files and actual data. No mocks, no fakes - just real Adobe After Effects project files.

## Available Real Test Data

### Core Test Files (in `data/` directory)

1. **Bit Depth Tests**
   - `BPC-8.aep` - 8-bit color depth
   - `BPC-16.aep` - 16-bit color depth  
   - `BPC-32.aep` - 32-bit color depth

2. **Expression Engine Tests**
   - `ExEn-es.aep` - ExtendScript expression engine
   - `ExEn-js.aep` - JavaScript expression engine

3. **Structure Tests**
   - `Item-01.aep` - Item structure testing
   - `Layer-01.aep` - Layer structure testing
   - `Property-01.aep` - Property structure testing

4. **Complex Project Test**
   - `sample-aep/Ai Text Intro.aep` - Complex real-world project (3157 items!)

## Test Infrastructure

### Helper Files

- `tests/helpers/real_aep_loader.go` - Loads real AEP files
- `tests/helpers/dangerous_test_utils.go` - Deep analysis using DangerousAnalyzer
- `tests/helpers/real_db_utils.go` - Real SQLite database testing

### Running Tests

```bash
# Run all real data tests
./run_real_tests.sh

# Run specific test categories
go test -v . -run TestRealAEP           # Core parser tests
go test -v . -run TestComplexRealAEP    # Complex file test
go test -bench BenchmarkRealAEP -run ^$ # Performance benchmarks
```

## Test Categories

### 1. Parser Tests
Tests the core RIFX parser with real AEP files:
- Bit depth parsing (8/16/32-bit)
- Expression engine detection
- Item/Layer/Property extraction
- Complex project handling

### 2. DangerousAnalyzer Tests
Deep inspection of real AEP files:
- Hidden layer detection
- Modular component discovery
- Text intelligence analysis
- API schema generation
- Complexity scoring

### 3. Database Tests
Real SQLite operations:
- Storing parsed project data
- Concurrent access testing
- Query performance
- Data integrity checks

### 4. Integration Tests
Multi-component workflows:
- Parse → Store → Query pipeline
- Multi-agent coordination
- Error recovery with real failures

### 5. Performance Tests
Benchmarks with real data:
- Parsing speed (currently ~10ms/file)
- Memory usage
- Concurrent parsing
- Large file handling

## Writing New Real-Data Tests

### Example Test Structure

```go
func TestRealFeature(t *testing.T) {
    // Load real AEP file
    path := getRealAEPPath("Layer-01.aep")
    project, err := Open(path)
    if err != nil {
        t.Fatalf("Failed to parse real AEP: %v", err)
    }
    
    // Test with real data
    for _, item := range project.Items {
        // Verify real properties
        if item.ItemType == ItemTypeComposition {
            // Test real composition data
        }
    }
}
```

### Using DangerousAnalyzer

```go
func TestDangerousAnalysis(t *testing.T) {
    result := helpers.RunDangerousAnalysis(t, "complex.aep")
    
    // Check real findings
    t.Logf("Hidden layers found: %d", len(result.HiddenLayers))
    t.Logf("Modular components: %d", result.ModularSystem.TotalModules)
    t.Logf("Automation score: %.2f", result.AutomationScore)
}
```

## Edge Cases with Real Files

### Currently Tested
- Empty compositions
- Minimal AEP files
- Complex nested structures (3000+ items)
- Different color depths
- Multiple expression engines

### Future Real-Data Tests
- Corrupted AEP files (create by hex-editing)
- Very large files (>100MB)
- Unicode/international text
- Missing asset references
- Circular dependencies

## Performance Metrics

Current benchmarks with real files:
- Simple AEP: ~1ms parsing time
- Complex AEP (3157 items): ~10ms parsing time
- Memory usage: <50MB for complex files

## Troubleshooting

### Common Issues

1. **File Not Found**
   - Ensure you're running from project root
   - Check file paths are relative to test location

2. **RIFX Header Missing**
   - File may be corrupted
   - Not a valid AEP file

3. **Parsing Failures**
   - Check Go version (1.17+)
   - Verify rifx dependency installed

## NO MOCKS Policy

This project strictly uses real data:
- ❌ No mock objects
- ❌ No fake data generators
- ❌ No synthetic test cases
- ✅ Only real AEP files
- ✅ Only actual parsing results
- ✅ Only genuine error conditions

## Contributing Tests

When adding new tests:
1. Use real AEP files only
2. Document what real aspect you're testing
3. Add files to `data/` or `tests/fixtures/`
4. Update this guide with new test scenarios

Remember: Real bugs hide behind fake data!