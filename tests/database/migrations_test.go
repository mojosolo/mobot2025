package database_test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/mobot2025/catalog"
	_ "github.com/mattn/go-sqlite3"
)

// TestDatabaseMigrations tests the database migration system
func TestDatabaseMigrations(t *testing.T) {
	// Create temporary database
	tempDir, err := ioutil.TempDir("", "mobot_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	
	// Test initial migration
	t.Run("initial_migration", func(t *testing.T) {
		db, err := catalog.NewDatabase(dbPath)
		if err != nil {
			t.Fatalf("Failed to create database: %v", err)
		}
		defer db.Close()

		// Verify all tables were created
		tables := []string{
			"projects",
			"compositions",
			"text_layers",
			"media_assets",
			"effects",
			"categories",
			"tags",
			"project_categories",
			"project_tags",
			"opportunities",
			"search_index",
			"analysis_results",
			"schema_migrations",
		}

		for _, table := range tables {
			if !tableExists(t, db, table) {
				t.Errorf("Table %s was not created", table)
			}
		}
	})

	// Test migration idempotency
	t.Run("migration_idempotency", func(t *testing.T) {
		// Open database again - migrations should not run twice
		db, err := catalog.NewDatabase(dbPath)
		if err != nil {
			t.Fatalf("Failed to reopen database: %v", err)
		}
		defer db.Close()

		// Count migration records
		var count int
		err = db.QueryRow("SELECT COUNT(DISTINCT version) FROM schema_migrations").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count migrations: %v", err)
		}

		// Should have exactly the number of migrations defined
		expectedMigrations := 12 // Based on the migration list in database.go
		if count != expectedMigrations {
			t.Errorf("Expected %d migrations, got %d", expectedMigrations, count)
		}
	})
}

// TestDatabaseSchema tests the schema structure
func TestDatabaseSchema(t *testing.T) {
	db := createTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	// Test projects table schema
	t.Run("projects_schema", func(t *testing.T) {
		columns := getTableColumns(t, db, "projects")
		
		expectedColumns := map[string]string{
			"id":                "INTEGER",
			"file_path":         "TEXT",
			"name":              "TEXT",
			"version":           "TEXT",
			"created_date":      "DATETIME",
			"modified_date":     "DATETIME",
			"expression_engine": "TEXT",
			"composition_count": "INTEGER",
			"size_bytes":        "INTEGER",
			"analyzed_at":       "DATETIME",
		}

		for col, expectedType := range expectedColumns {
			if colType, exists := columns[col]; !exists {
				t.Errorf("Missing column %s in projects table", col)
			} else if colType != expectedType {
				t.Errorf("Column %s has type %s, expected %s", col, colType, expectedType)
			}
		}
	})

	// Test text_layers table schema with foreign key
	t.Run("text_layers_foreign_key", func(t *testing.T) {
		// Insert test project
		_, err := db.Exec(`
			INSERT INTO projects (file_path, name, version) 
			VALUES (?, ?, ?)
		`, "/test/project.aep", "Test Project", "2021")
		if err != nil {
			t.Fatal(err)
		}

		// Try to insert text layer with invalid project_id
		_, err = db.Exec(`
			INSERT INTO text_layers (project_id, composition_name, layer_name, text)
			VALUES (?, ?, ?, ?)
		`, 9999, "Comp", "Layer", "Text")
		
		if err == nil {
			t.Error("Expected foreign key constraint error")
		}
	})
}

// TestDatabaseIndexes tests that proper indexes are created
func TestDatabaseIndexes(t *testing.T) {
	db := createTestDatabase(t)
	defer cleanupTestDatabase(t, db)

	indexes := []struct {
		table string
		index string
	}{
		{"projects", "idx_projects_file_path"},
		{"text_layers", "idx_text_layers_project"},
		{"media_assets", "idx_media_assets_project"},
		{"search_index", "idx_search_content"},
	}

	for _, idx := range indexes {
		t.Run(idx.index, func(t *testing.T) {
			var count int
			err := db.QueryRow(`
				SELECT COUNT(*) FROM sqlite_master 
				WHERE type='index' AND name=? AND tbl_name=?
			`, idx.index, idx.table).Scan(&count)
			
			if err != nil {
				t.Fatal(err)
			}
			if count == 0 {
				t.Errorf("Index %s not found on table %s", idx.index, idx.table)
			}
		})
	}
}

// TestMigrationRollback tests migration rollback capability
func TestMigrationRollback(t *testing.T) {
	// This would test rollback functionality if implemented
	t.Skip("Rollback functionality not yet implemented")
}

// Helper functions

func createTestDatabase(t *testing.T) *catalog.Database {
	tempDir, err := ioutil.TempDir("", "mobot_test_*")
	if err != nil {
		t.Fatal(err)
	}
	
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := catalog.NewDatabase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	
	// Store temp dir for cleanup
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	
	return db
}

func cleanupTestDatabase(t *testing.T, db *catalog.Database) {
	if err := db.Close(); err != nil {
		t.Errorf("Failed to close database: %v", err)
	}
}

func tableExists(t *testing.T, db *catalog.Database, tableName string) bool {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master 
		WHERE type='table' AND name=?
	`, tableName).Scan(&count)
	
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	
	return count > 0
}

func getTableColumns(t *testing.T, db *catalog.Database, tableName string) map[string]string {
	rows, err := db.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	columns := make(map[string]string)
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		
		err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		if err != nil {
			t.Fatal(err)
		}
		
		columns[name] = ctype
	}
	
	return columns
}