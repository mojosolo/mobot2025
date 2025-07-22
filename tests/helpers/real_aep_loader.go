package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	
	aep "github.com/mojosolo/mobot2025"
)

// GetRealAEPPath returns the absolute path to a real AEP file in the data directory
func GetRealAEPPath(filename string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..", "..", "data", filename)
}

// LoadRealAEPFile loads an actual AEP file from the data directory
func LoadRealAEPFile(filename string) (*aep.Project, error) {
	path := GetRealAEPPath(filename)
	
	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("real AEP file not found: %s", path)
	}
	
	// Parse the real AEP file
	project, err := aep.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse real AEP file %s: %w", filename, err)
	}
	
	return project, nil
}

// GetAllRealAEPFiles returns a list of all real AEP files available for testing
func GetAllRealAEPFiles() []string {
	return []string{
		"BPC-8.aep",      // 8-bit color depth test
		"BPC-16.aep",     // 16-bit color depth test  
		"BPC-32.aep",     // 32-bit color depth test
		"ExEn-es.aep",    // ExtendScript expression engine
		"ExEn-js.aep",    // JavaScript expression engine
		"Item-01.aep",    // Item structure test
		"Layer-01.aep",   // Layer structure test
		"Property-01.aep", // Property structure test
	}
}

// LoadRealAEPBytes loads raw bytes from a real AEP file for low-level testing
func LoadRealAEPBytes(filename string) ([]byte, error) {
	path := GetRealAEPPath(filename)
	return ioutil.ReadFile(path)
}

// ValidateRealAEPFile performs basic validation on a real AEP file
func ValidateRealAEPFile(filename string) error {
	data, err := LoadRealAEPBytes(filename)
	if err != nil {
		return err
	}
	
	// Check RIFX header
	if len(data) < 4 {
		return fmt.Errorf("file too small to be valid AEP")
	}
	
	// AEP files start with RIFX
	if string(data[0:4]) != "RIFX" {
		return fmt.Errorf("invalid AEP file: missing RIFX header")
	}
	
	return nil
}

// GetComplexRealAEPPath returns path to the complex sample project
func GetComplexRealAEPPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..", "..", "sample-aep", "Ai Text Intro.aep")
}