package fixtures

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestEnvironment provides a complete test environment
type TestEnvironment struct {
	t       *testing.T
	TempDir string
	Files   map[string]string
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	tempDir, err := ioutil.TempDir("", "mobot-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	return &TestEnvironment{
		t:       t,
		TempDir: tempDir,
		Files:   make(map[string]string),
	}
}

// Cleanup removes all test files and directories
func (env *TestEnvironment) Cleanup() {
	if err := os.RemoveAll(env.TempDir); err != nil {
		env.t.Errorf("Failed to cleanup temp dir: %v", err)
	}
}

// CreateFile creates a file in the test environment
func (env *TestEnvironment) CreateFile(name string, content []byte) string {
	path := filepath.Join(env.TempDir, name)
	dir := filepath.Dir(path)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		env.t.Fatalf("Failed to create directory %s: %v", dir, err)
	}
	
	if err := ioutil.WriteFile(path, content, 0644); err != nil {
		env.t.Fatalf("Failed to write file %s: %v", path, err)
	}
	
	env.Files[name] = path
	return path
}

// CreateMockAEP creates a mock AEP file in the test environment
func (env *TestEnvironment) CreateMockAEP(name string) string {
	gen := NewMockAEPGenerator()
	content, err := gen.GenerateMinimalAEP()
	if err != nil {
		env.t.Fatalf("Failed to generate mock AEP: %v", err)
	}
	
	return env.CreateFile(name, content)
}

// GetFilePath returns the full path for a file in the test environment
func (env *TestEnvironment) GetFilePath(name string) string {
	if path, ok := env.Files[name]; ok {
		return path
	}
	return filepath.Join(env.TempDir, name)
}

// AssertFileExists checks if a file exists in the test environment
func (env *TestEnvironment) AssertFileExists(name string) {
	path := env.GetFilePath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		env.t.Errorf("Expected file %s to exist, but it doesn't", path)
	}
}

// AssertFileContains checks if a file contains expected content
func (env *TestEnvironment) AssertFileContains(name string, expected string) {
	path := env.GetFilePath(name)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		env.t.Errorf("Failed to read file %s: %v", path, err)
		return
	}
	
	if string(content) != expected {
		env.t.Errorf("File %s content mismatch:\nExpected: %s\nActual: %s", 
			path, expected, string(content))
	}
}

// TestDatabase provides a test database
type TestDatabase struct {
	t      *testing.T
	dbPath string
}

// NewTestDatabase creates a new test database
func NewTestDatabase(t *testing.T) *TestDatabase {
	tempDir, err := ioutil.TempDir("", "mobot-db-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	return &TestDatabase{
		t:      t,
		dbPath: filepath.Join(tempDir, "test.db"),
	}
}

// Path returns the database path
func (db *TestDatabase) Path() string {
	return db.dbPath
}

// Cleanup removes the test database
func (db *TestDatabase) Cleanup() {
	dir := filepath.Dir(db.dbPath)
	if err := os.RemoveAll(dir); err != nil {
		db.t.Errorf("Failed to cleanup database: %v", err)
	}
}

// Common test assertions

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotNil fails the test if value is nil
func AssertNotNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: expected non-nil value", msg)
	}
}

// AssertTrue fails the test if condition is false
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: expected true", msg)
	}
}

// AssertFalse fails the test if condition is true
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("%s: expected false", msg)
	}
}

// AssertContains fails the test if slice doesn't contain value
func AssertContains(t *testing.T, slice []string, value string, msg string) {
	t.Helper()
	for _, v := range slice {
		if v == value {
			return
		}
	}
	t.Errorf("%s: slice %v doesn't contain %s", msg, slice, value)
}

// Benchmark helpers

// BenchmarkEnvironment provides a benchmark environment
type BenchmarkEnvironment struct {
	b       *testing.B
	TempDir string
}

// NewBenchmarkEnvironment creates a new benchmark environment
func NewBenchmarkEnvironment(b *testing.B) *BenchmarkEnvironment {
	tempDir, err := ioutil.TempDir("", "mobot-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	
	return &BenchmarkEnvironment{
		b:       b,
		TempDir: tempDir,
	}
}

// Cleanup removes all benchmark files
func (env *BenchmarkEnvironment) Cleanup() {
	os.RemoveAll(env.TempDir)
}

// ResetTimer resets the benchmark timer after setup
func (env *BenchmarkEnvironment) ResetTimer() {
	env.b.ResetTimer()
}