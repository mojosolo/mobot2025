package aep_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	
	aep "github.com/yourusername/mobot2025"
)

// Enhanced test suite following Verification Agent principles
// Builds on existing test patterns with 80%+ coverage target

// TestParseAllSampleFiles validates parser against all test AEP files
// Extends the pattern from existing test files
func TestParseAllSampleFiles(t *testing.T) {
	// Reuse the data directory pattern from existing tests
	testFiles, err := filepath.Glob("data/*.aep")
	if err != nil {
		t.Fatalf("Failed to find test files: %v", err)
	}
	
	for _, testFile := range testFiles {
		t.Run(filepath.Base(testFile), func(t *testing.T) {
			// Use existing Open function
			project, err := aep.Open(testFile)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", testFile, err)
				return
			}
			
			// Validate basic structure
			validateProjectStructure(t, project, testFile)
		})
	}
}

// validateProjectStructure performs comprehensive validation
// Following the verification agent's empirical validation approach
func validateProjectStructure(t *testing.T, project *aep.Project, filename string) {
	// Basic validations
	if project == nil {
		t.Error("Project is nil")
		return
	}
	
	if project.RootFolder == nil {
		t.Error("Root folder is missing")
		return
	}
	
	// Validate based on filename patterns (reusing existing test files)
	switch filepath.Base(filename) {
	case "BPC-8.aep":
		if project.Depth != aep.BPC8 {
			t.Errorf("Expected 8-bit depth, got %v", project.Depth)
		}
	case "BPC-16.aep":
		if project.Depth != aep.BPC16 {
			t.Errorf("Expected 16-bit depth, got %v", project.Depth)
		}
	case "BPC-32.aep":
		if project.Depth != aep.BPC32 {
			t.Errorf("Expected 32-bit depth, got %v", project.Depth)
		}
	case "ExEn-js.aep":
		if project.ExpressionEngine != "javascript-1.0" {
			t.Errorf("Expected JavaScript engine, got %s", project.ExpressionEngine)
		}
	}
	
	// Validate all items are properly indexed
	for id, item := range project.Items {
		if item.ID != id {
			t.Errorf("Item ID mismatch: map key %d != item.ID %d", id, item.ID)
		}
		
		// Validate item-specific properties
		validateItemProperties(t, item)
	}
}

// validateItemProperties checks item-specific constraints
func validateItemProperties(t *testing.T, item *aep.Item) {
	switch item.ItemType {
	case aep.ItemTypeFolder:
		// Folders should have contents or be empty
		if item.FolderContents == nil {
			t.Errorf("Folder %s has nil contents array", item.Name)
		}
		
	case aep.ItemTypeComposition:
		// Compositions should have valid dimensions
		if item.FootageDimensions[0] == 0 || item.FootageDimensions[1] == 0 {
			t.Errorf("Composition %s has invalid dimensions: %v", item.Name, item.FootageDimensions)
		}
		if item.FootageFramerate <= 0 {
			t.Errorf("Composition %s has invalid framerate: %f", item.Name, item.FootageFramerate)
		}
		
	case aep.ItemTypeFootage:
		// Footage should have valid properties
		if item.FootageType == 0 && item.FootageDimensions[0] == 0 {
			t.Errorf("Footage %s has no type or dimensions", item.Name)
		}
	}
}

// BenchmarkParseProject measures parsing performance
// Following performance optimization guidelines
func BenchmarkParseProject(b *testing.B) {
	// Find a test file
	testFile := "data/Item-01.aep"
	if _, err := os.Stat(testFile); err != nil {
		b.Skip("Test file not found")
	}
	
	// Read file once to avoid I/O in benchmark
	data, err := os.ReadFile(testFile)
	if err != nil {
		b.Fatalf("Failed to read test file: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Benchmark parsing from memory
		_, err := aep.FromReader(bytes.NewReader(data))
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}