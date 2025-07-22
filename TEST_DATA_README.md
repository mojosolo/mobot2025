# Test Data for MoBot 2025

## Important: Test Files Not Included

AEP files and other binary test data are NOT included in this repository to keep it lightweight and avoid versioning binary files.

## How to Get Test Files

### Option 1: Download Sample AEP Files
The test suite expects AEP files in the `data/` directory. You can:
1. Create your own AEP files using Adobe After Effects
2. Use the AEP file paths referenced in the tests as a guide for naming

### Option 2: Use Minimal Test Data
For basic testing, create empty files with the expected names:
```bash
mkdir -p data sample-aep/Assets sample-aep/Help
touch data/BPC-8.aep data/BPC-16.aep data/BPC-32.aep
touch data/ExEn-js.aep data/ExEn-es.aep
touch data/Item-01.aep data/Layer-01.aep data/Property-01.aep
touch "sample-aep/Ai Text Intro.aep"
```

### Expected Test Files
The test suite references these files:
- `data/BPC-8.aep` - 8-bit color depth test
- `data/BPC-16.aep` - 16-bit color depth test  
- `data/BPC-32.aep` - 32-bit color depth test
- `data/ExEn-js.aep` - JavaScript expression engine test
- `data/ExEn-es.aep` - ExtendScript expression engine test
- `data/Item-01.aep` - Item parsing test
- `data/Layer-01.aep` - Layer parsing test
- `data/Property-01.aep` - Property parsing test
- `sample-aep/Ai Text Intro.aep` - Complex project test

## Why These Files Are Excluded

1. **Binary files don't version well** - Git isn't designed for binary files
2. **Repository size** - AEP files can be very large (MB to GB)
3. **Licensing** - AEP files may contain proprietary content
4. **Not needed for code review** - The parser code is what matters

## Running Tests Without Real AEP Files

Many tests will fail without real AEP files, but you can still:
- Review the parser implementation
- Run unit tests that don't require files
- Build the project
- Understand the architecture