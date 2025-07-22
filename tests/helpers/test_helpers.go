package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/mobot2025/catalog"
)

// RealTestFixtures provides paths to REAL test fixture files
var RealTestFixtures = struct {
	BitDepth8    string
	BitDepth16   string
	BitDepth32   string
	ExtendScript string
	JavaScript   string
	ItemTest     string
	LayerTest    string
	PropertyTest string
	ComplexAEP   string
}{
	BitDepth8:    "../../data/BPC-8.aep",
	BitDepth16:   "../../data/BPC-16.aep",
	BitDepth32:   "../../data/BPC-32.aep",
	ExtendScript: "../../data/ExEn-es.aep",
	JavaScript:   "../../data/ExEn-js.aep",
	ItemTest:     "../../data/Item-01.aep",
	LayerTest:    "../../data/Layer-01.aep",
	PropertyTest: "../../data/Property-01.aep",
	ComplexAEP:   "../../sample-aep/Ai Text Intro.aep",
}

// CreateRealTempDir creates a temporary directory for testing with real data
func CreateRealTempDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "mobot_real_test_*")
	if err != nil {
		t.Fatal(err)
	}
	
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	
	return tempDir
}

// LoadRealTestFixture loads a REAL test fixture file
func LoadRealTestFixture(t *testing.T, filename string) []byte {
	t.Helper()
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to load real test fixture %s: %v", filename, err)
	}
	
	// Verify it's a real AEP file
	if len(data) < 4 || string(data[0:4]) != "RIFX" {
		t.Fatalf("Invalid AEP file (missing RIFX header): %s", filename)
	}
	
	return data
}

// ParseRealAEPFile parses a real AEP file for testing
func ParseRealAEPFile(t *testing.T, fixturePath string) *catalog.ProjectMetadata {
	t.Helper()
	
	parser := catalog.NewParser()
	metadata, err := parser.ParseAEPFile(fixturePath)
	if err != nil {
		t.Fatalf("Failed to parse real AEP file %s: %v", fixturePath, err)
	}
	
	return metadata
}

// CompareRealProjects compares two real parsed projects for testing
func CompareRealProjects(t *testing.T, expected, actual *catalog.ProjectMetadata) {
	t.Helper()
	
	if expected.FileName != actual.FileName {
		t.Errorf("Project filename mismatch: expected %s, got %s", expected.FileName, actual.FileName)
	}
	
	if expected.BitDepth != actual.BitDepth {
		t.Errorf("BitDepth mismatch: expected %d, got %d", expected.BitDepth, actual.BitDepth)
	}
	
	if expected.ExpressionEngine != actual.ExpressionEngine {
		t.Errorf("ExpressionEngine mismatch: expected %s, got %s", expected.ExpressionEngine, actual.ExpressionEngine)
	}
	
	if expected.Compositions != actual.Compositions {
		t.Errorf("Composition count mismatch: expected %d, got %d", expected.Compositions, actual.Compositions)
	}
}

// TestAllRealFiles runs a test function on all real AEP files
func TestAllRealFiles(t *testing.T, testFunc func(t *testing.T, filepath string)) {
	t.Helper()
	
	realFiles := []struct {
		name string
		path string
	}{
		{"8-bit depth", RealTestFixtures.BitDepth8},
		{"16-bit depth", RealTestFixtures.BitDepth16},
		{"32-bit depth", RealTestFixtures.BitDepth32},
		{"ExtendScript", RealTestFixtures.ExtendScript},
		{"JavaScript", RealTestFixtures.JavaScript},
		{"Item test", RealTestFixtures.ItemTest},
		{"Layer test", RealTestFixtures.LayerTest},
		{"Property test", RealTestFixtures.PropertyTest},
	}
	
	for _, file := range realFiles {
		t.Run(file.name, func(t *testing.T) {
			testFunc(t, file.path)
		})
	}
}

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", message)
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertContains fails the test if substr is not in str
func AssertContains(t *testing.T, str, substr, message string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Fatalf("%s: %q does not contain %q", message, str, substr)
	}
}

// WaitForRealCondition waits for a real condition to be true or times out
func WaitForRealCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()
	
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	t.Fatalf("Timeout waiting for condition: %s", message)
}

// RunRealConcurrentTest runs a test function concurrently with real data
func RunRealConcurrentTest(t *testing.T, count int, testFunc func(id int)) {
	t.Helper()
	
	var wg sync.WaitGroup
	wg.Add(count)
	
	for i := 0; i < count; i++ {
		go func(id int) {
			defer wg.Done()
			testFunc(id)
		}(i)
	}
	
	wg.Wait()
}

// CaptureRealOutput captures stdout/stderr during real test execution
func CaptureRealOutput(t *testing.T, f func()) string {
	t.Helper()
	
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	f()
	
	w.Close()
	os.Stdout = old
	
	out, _ := ioutil.ReadAll(r)
	return string(out)
}

// BenchmarkRealFile runs benchmarks on real AEP files
func BenchmarkRealFile(b *testing.B, filepath string) {
	data := LoadRealTestFixture(&testing.T{}, filepath)
	
	b.ResetTimer()
	b.SetBytes(int64(len(data)))
	
	for i := 0; i < b.N; i++ {
		// Benchmark operation here
		_ = data
	}
}

// ValidateRealParsingResult validates that parsing produced real results
func ValidateRealParsingResult(t *testing.T, metadata *catalog.ProjectMetadata) {
	t.Helper()
	
	if metadata == nil {
		t.Fatal("Parsing returned nil metadata")
	}
	
	if metadata.FilePath == "" {
		t.Error("Missing file path in metadata")
	}
	
	if metadata.ParsedAt.IsZero() {
		t.Error("Missing parsed timestamp")
	}
	
	// Real AEP files should have some content
	if metadata.TotalItems == 0 && metadata.Compositions == 0 {
		t.Error("No items or compositions found - suspicious for real AEP file")
	}
}

// DumpRealMetadata dumps real metadata for debugging
func DumpRealMetadata(t *testing.T, metadata *catalog.ProjectMetadata) {
	t.Helper()
	
	t.Logf("=== Real AEP Metadata ===")
	t.Logf("File: %s", metadata.FileName)
	t.Logf("Path: %s", metadata.FilePath)
	t.Logf("BitDepth: %d", metadata.BitDepth)
	t.Logf("ExpressionEngine: %s", metadata.ExpressionEngine)
	t.Logf("Total Items: %d", metadata.TotalItems)
	t.Logf("Compositions: %d", metadata.Compositions)
	t.Logf("Text Layers: %d", len(metadata.TextLayers))
	t.Logf("Media Assets: %d", len(metadata.MediaAssets))
	t.Logf("Effects: %d", len(metadata.Effects))
	t.Logf("Parsed At: %v", metadata.ParsedAt)
}