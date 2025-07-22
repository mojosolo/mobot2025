# Demo Directory

The demo tools have been reorganized into the `cmd/` directory for better organization:

## Available Tools

### Story Viewers
- `cmd/story-viewer/simple.go` - Basic story extraction
- `cmd/story-viewer/easy.go` - Easy mode with simplified output
- `cmd/story-viewer/unified.go` - Unified story viewer
- `cmd/story-viewer/ultimate.go` - Ultimate story viewer with all features

### Text Extractors
- `cmd/text-extractor/debug.go` - Debug text extraction
- `cmd/text-extractor/final.go` - Final text extraction
- `cmd/text-extractor/enhanced.go` - Enhanced text extraction
- `cmd/text-extractor/test_enhanced.go` - Test enhanced extraction

### Report Generators
- `cmd/report-generator/enhanced_ux.go` - Enhanced UX report
- `cmd/report-generator/html.go` - HTML report generator
- `cmd/report-generator/detailed.go` - Detailed report generator

### Scanners and Parsers
- `cmd/scanner/main.go` - Scan all AEP files
- `cmd/parser/sample.go` - Parse sample files
- `cmd/parser/summary.go` - Parse with summary
- `cmd/parser/debug_raw.go` - Debug raw RIFX data

## Running the Tools

Each tool can be run independently:

```bash
# Run story viewer
go run cmd/story-viewer/simple.go data/sample.aep

# Run text extractor
go run cmd/text-extractor/debug.go data/sample.aep

# Run report generator
go run cmd/report-generator/html.go data/sample.aep

# Run scanner
go run cmd/scanner/main.go data/
```

## Main Demo

The main demo (`demo/main.go`) provides a simple example of parsing an AEP file and displaying project information.

```bash
go run demo/main.go
```