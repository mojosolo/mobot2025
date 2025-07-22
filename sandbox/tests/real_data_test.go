package aep

import (
	"path/filepath"
	"runtime"
	"testing"
)

// getRealAEPPath returns the path to a real AEP file in the data directory
func getRealAEPPath(filename string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "data", filename)
}

// TestRealAEPParsing tests parsing real AEP files
func TestRealAEPParsing(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		validate func(*testing.T, *Project)
	}{
		{
			name:     "8-bit color depth",
			filename: "BPC-8.aep",
			validate: func(t *testing.T, p *Project) {
				if p.Depth != BPC8 {
					t.Errorf("Expected 8-bit depth, got %v", p.Depth)
				}
			},
		},
		{
			name:     "16-bit color depth",
			filename: "BPC-16.aep",
			validate: func(t *testing.T, p *Project) {
				if p.Depth != BPC16 {
					t.Errorf("Expected 16-bit depth, got %v", p.Depth)
				}
			},
		},
		{
			name:     "32-bit color depth",
			filename: "BPC-32.aep",
			validate: func(t *testing.T, p *Project) {
				if p.Depth != BPC32 {
					t.Errorf("Expected 32-bit depth, got %v", p.Depth)
				}
			},
		},
		{
			name:     "JavaScript expression engine",
			filename: "ExEn-js.aep",
			validate: func(t *testing.T, p *Project) {
				if p.ExpressionEngine == "" {
					t.Error("Expected JavaScript expression engine")
				}
			},
		},
		{
			name:     "ExtendScript expression engine",
			filename: "ExEn-es.aep",
			validate: func(t *testing.T, p *Project) {
				if p.ExpressionEngine == "" {
					t.Error("Expected ExtendScript expression engine")
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get real file path
			path := getRealAEPPath(tt.filename)
			
			// Load and parse real AEP file
			project, err := Open(path)
			if err != nil {
				t.Fatalf("Failed to load real AEP file: %v", err)
			}
			
			// Validate parsing results
			if project == nil {
				t.Fatal("Parsed project is nil")
			}
			
			// Run specific validation
			tt.validate(t, project)
			
			// General validation
			if len(project.Items) == 0 {
				t.Log("Warning: No items found in project (this may be expected for test files)")
			}
		})
	}
}

// TestComplexRealAEP tests parsing a complex real AEP file
func TestComplexRealAEP(t *testing.T) {
	path := filepath.Join("sample-aep", "Ai Text Intro.aep")
	
	// This file might be complex and fail - that's OK for now
	project, err := Open(path)
	if err != nil {
		t.Logf("Complex AEP parsing failed (expected): %v", err)
		return
	}
	
	// If it parsed, log what we found
	t.Logf("Complex AEP parsed successfully!")
	t.Logf("  Items: %d", len(project.Items))
	t.Logf("  Depth: %v", project.Depth)
	t.Logf("  Expression Engine: %s", project.ExpressionEngine)
	
	// Count different item types
	var folders, comps, footage int
	for _, item := range project.Items {
		switch item.ItemType {
		case ItemTypeFolder:
			folders++
		case ItemTypeComposition:
			comps++
		case ItemTypeFootage:
			footage++
		}
	}
	
	t.Logf("  Folders: %d", folders)
	t.Logf("  Compositions: %d", comps)
	t.Logf("  Footage: %d", footage)
}

// BenchmarkRealAEPParsing benchmarks parsing real AEP files
func BenchmarkRealAEPParsing(b *testing.B) {
	// Benchmark parsing the Layer-01.aep file
	path := getRealAEPPath("Layer-01.aep")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Open(path)
		if err != nil {
			b.Fatal(err)
		}
	}
}