package helpers

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/yourusername/mobot2025/catalog"
	_ "github.com/mattn/go-sqlite3"
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
	
	parser := catalog.NewParser()
	
	// Parse and store all real AEP files
	for _, filename := range GetAllRealAEPFiles() {
		path := GetRealAEPPath(filename)
		
		metadata, err := parser.ParseAEPFile(path)
		if err != nil {
			t.Logf("Warning: Failed to parse %s: %v", filename, err)
			continue
		}
		
		// Store in real database
		projectID, err := db.StoreProject(metadata)
		if err != nil {
			t.Logf("Warning: Failed to store %s: %v", filename, err)
			continue
		}
		
		t.Logf("Stored real project %s with ID %d", filename, projectID)
	}
}

// QueryRealData performs real queries on populated database
func QueryRealData(t *testing.T, db *catalog.Database) {
	t.Helper()
	
	// Query real project count
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query project count: %v", err)
	}
	t.Logf("Real projects in database: %d", count)
	
	// Query real text layers
	var textCount int
	err = db.QueryRow("SELECT COUNT(*) FROM text_layers").Scan(&textCount)
	if err != nil {
		t.Fatalf("Failed to query text layer count: %v", err)
	}
	t.Logf("Real text layers in database: %d", textCount)
	
	// Query real media assets
	var mediaCount int
	err = db.QueryRow("SELECT COUNT(*) FROM media_assets").Scan(&mediaCount)
	if err != nil {
		t.Fatalf("Failed to query media asset count: %v", err)
	}
	t.Logf("Real media assets in database: %d", mediaCount)
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
			
			// Perform real queries
			var projects []string
			rows, err := db.Query("SELECT name FROM projects")
			if err != nil {
				t.Errorf("Goroutine %d: query failed: %v", id, err)
				return
			}
			defer rows.Close()
			
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					t.Errorf("Goroutine %d: scan failed: %v", id, err)
					return
				}
				projects = append(projects, name)
			}
			
			t.Logf("Goroutine %d: found %d real projects", id, len(projects))
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// BenchmarkRealDatabaseOperations benchmarks real database operations
func BenchmarkRealDatabaseOperations(b *testing.B, db *catalog.Database) {
	// Populate once
	parser := catalog.NewParser()
	path := GetRealAEPPath("Layer-01.aep")
	metadata, _ := parser.ParseAEPFile(path)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Benchmark real store operation
		_, err := db.StoreProject(metadata)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// VerifyRealDatabaseIntegrity verifies database integrity after real operations
func VerifyRealDatabaseIntegrity(t *testing.T, dbPath string) {
	t.Helper()
	
	// Open database directly
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Run integrity check
	var result string
	err = db.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to run integrity check: %v", err)
	}
	
	if result != "ok" {
		t.Fatalf("Database integrity check failed: %s", result)
	}
	
	t.Log("Database integrity check passed")
}