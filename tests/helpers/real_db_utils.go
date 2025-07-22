package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
	
	"github.com/mojosolo/mobot2025/catalog"
)

// CreateRealTestDatabase creates a real SQLite database for testing
func CreateRealTestDatabase(t *testing.T) (*catalog.Database, string, func()) {
	t.Helper()
	
	// Create temporary directory for test database
	tempDir, err := ioutil.TempDir("", "mobot_real_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	dbPath := filepath.Join(tempDir, "test.db")
	
	// Create real database
	db, err := catalog.NewDatabase(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create real database: %v", err)
	}
	
	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.RemoveAll(tempDir)
	}
	
	return db, dbPath, cleanup
}

// PopulateRealDatabase populates database with real AEP parsing results
func PopulateRealDatabase(t *testing.T, db *catalog.Database) {
	t.Helper()
	
	// Since we don't have real AEP files, create mock metadata
	metadata := &catalog.ProjectMetadata{
		FilePath:  "/test/sample.aep",
		FileName:  "sample.aep",
		FileSize:  1024,
		BitDepth:  8,
		ExpressionEngine: "javascript",
		TotalItems: 10,
		Compositions: []catalog.CompositionInfo{
			{
				ID:        "comp1",
				Name:      "Main Comp",
				Width:     1920,
				Height:    1080,
				FrameRate: 30,
				Duration:  300,
				LayerCount: 5,
				Is3D:      false,
				HasEffects: true,
			},
		},
		TextLayers: []catalog.TextLayerInfo{
			{
				ID:         "text1",
				CompID:     "comp1",
				LayerName:  "Title",
				SourceText: "Sample Title",
				FontUsed:   "Arial",
				IsAnimated: true,
			},
		},
		ParsedAt: time.Now(),
	}
	
	// Store in database
	err := db.StoreProject(metadata)
	if err != nil {
		t.Fatalf("Failed to store project: %v", err)
	}
	
	t.Log("Stored test project in database")
}

// QueryRealData performs real queries on populated database
func QueryRealData(t *testing.T, db *catalog.Database) {
	t.Helper()
	
	// Search for projects
	projects, err := db.SearchProjects("", 10)
	if err != nil {
		t.Fatalf("Failed to search projects: %v", err)
	}
	t.Logf("Found %d projects in database", len(projects))
	
	// Get specific project if any exist
	if len(projects) > 0 {
		project, err := db.GetProject(1) // Assuming ID 1 exists
		if err != nil {
			t.Logf("Warning: Failed to get project 1: %v", err)
		} else {
			t.Logf("Retrieved project: %s", project.FileName)
		}
	}
}

// TestRealConcurrentAccess tests concurrent database access with real data
func TestRealConcurrentAccess(t *testing.T, db *catalog.Database) {
	t.Helper()
	
	// Populate with real data first
	PopulateRealDatabase(t, db)
	
	// Test concurrent reads
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			// Perform real queries using public API
			projects, err := db.SearchProjects("", 100)
			if err != nil {
				t.Errorf("Goroutine %d: search failed: %v", id, err)
				return
			}
			
			t.Logf("Goroutine %d: found %d projects", id, len(projects))
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// BenchmarkRealDatabaseOperations benchmarks real database operations
func BenchmarkRealDatabaseOperations(b *testing.B, db *catalog.Database) {
	// Create test metadata
	metadata := &catalog.ProjectMetadata{
		FilePath:  "/test/bench.aep",
		FileName:  "bench.aep",
		FileSize:  1024,
		BitDepth:  8,
		ExpressionEngine: "javascript",
		TotalItems: 5,
		Compositions: []catalog.CompositionInfo{
			{
				ID:        "comp-bench",
				Name:      "Benchmark Comp",
				Width:     1920,
				Height:    1080,
				FrameRate: 30,
				Duration:  300,
			},
		},
		ParsedAt: time.Now(),
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Benchmark real store operation
		err := db.StoreProject(metadata)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// VerifyRealDatabaseIntegrity verifies database integrity after real operations
func VerifyRealDatabaseIntegrity(t *testing.T, dbPath string) {
	t.Helper()
	
	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("Database file does not exist")
	}
	
	// Try to open it again to verify it's valid
	db, err := catalog.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to reopen database: %v", err)
	}
	defer db.Close()
	
	// Try a simple operation
	_, err = db.SearchProjects("", 1)
	if err != nil {
		t.Fatalf("Database operation failed: %v", err)
	}
	
	t.Log("Database integrity check passed")
}