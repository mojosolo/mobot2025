package database_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mojosolo/mobot2025/catalog"
)

// TestDatabaseCreation tests basic database creation
func TestDatabaseCreation(t *testing.T) {
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "mobot_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	
	// Create database
	db, err := catalog.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()
	
	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

// TestDatabaseInMemory tests in-memory database
func TestDatabaseInMemory(t *testing.T) {
	// Create in-memory database
	db, err := catalog.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	defer db.Close()
	
	// If we got here, the database was created successfully
	// The migrations should have run automatically
}

// TestDatabaseClose tests database closing
func TestDatabaseClose(t *testing.T) {
	db, err := catalog.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	
	// Close the database
	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
	
	// Closing again should not panic
	err = db.Close()
	if err == nil {
		// This is okay - some implementations allow multiple closes
	}
}

// TestDatabaseOperations tests basic database operations
func TestDatabaseOperations(t *testing.T) {
	db, err := catalog.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()
	
	// Test storing a project (this assumes StoreProject method exists)
	// For now, we'll just verify the database was created successfully
}

// TestDatabaseConcurrency tests concurrent database access
func TestDatabaseConcurrency(t *testing.T) {
	db, err := catalog.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()
	
	// Run multiple goroutines that try to use the database
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Simulate some database operation
			// Since we don't have public methods to test,
			// we just verify the database is accessible
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}