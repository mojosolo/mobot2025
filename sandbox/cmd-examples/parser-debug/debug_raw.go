package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	
	aep "github.com/mojosolo/mobot2025"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_raw_data.go <aep-file>")
		os.Exit(1)
	}
	
	aepFile := os.Args[1]
	fmt.Printf("🔍 Debugging raw data for: %s\n\n", aepFile)
	
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	// Look for first text layer
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			for _, layer := range item.CompositionLayers {
				if layer.Text != nil {
					fmt.Printf("📄 Found text layer: %s in composition: %s\n", layer.Name, item.Name)
					fmt.Printf("🔍 Analyzing text property tree...\n\n")
					
					// Analyze the property tree
					analyzePropertyTree(layer.Text, 0)
					
					// Exit after first text layer for focused analysis
					return
				}
			}
		}
	}
	
	fmt.Println("❌ No text layers found")
}

func analyzePropertyTree(prop *aep.Property, indent int) {
	spaces := strings.Repeat("  ", indent)
	
	fmt.Printf("%s📦 Property: %s\n", spaces, prop.MatchName)
	fmt.Printf("%s   Name: %s\n", spaces, prop.Name)
	if prop.Label != "" && prop.Label != "-_0_/-" {
		fmt.Printf("%s   Label: %s\n", spaces, prop.Label)
	}
	
	// Check for raw data
	if len(prop.RawData) > 0 {
		fmt.Printf("%s   RawData: %d bytes\n", spaces, len(prop.RawData))
		
		// Show hex dump of first 256 bytes
		showBytes := 256
		if len(prop.RawData) < showBytes {
			showBytes = len(prop.RawData)
		}
		
		fmt.Printf("%s   Hex dump (first %d bytes):\n", spaces, showBytes)
		hexDump := hex.Dump(prop.RawData[:showBytes])
		for _, line := range strings.Split(hexDump, "\n") {
			if line != "" {
				fmt.Printf("%s   %s\n", spaces, line)
			}
		}
		
		// Try to find text patterns
		fmt.Printf("%s   Text patterns found:\n", spaces)
		findTextPatterns(prop.RawData, spaces+"   ")
	}
	
	// Check select options
	if len(prop.SelectOptions) > 0 {
		fmt.Printf("%s   SelectOptions: %v\n", spaces, prop.SelectOptions)
	}
	
	// Recursively analyze children
	if len(prop.Properties) > 0 {
		fmt.Printf("%s   Children: %d\n", spaces, len(prop.Properties))
		for _, child := range prop.Properties {
			analyzePropertyTree(child, indent+1)
		}
	}
	
	fmt.Println()
}

func findTextPatterns(data []byte, prefix string) {
	// Look for various text patterns
	
	// Pattern 1: ASCII text sequences
	var asciiText []byte
	for i := 0; i < len(data); i++ {
		if data[i] >= 32 && data[i] <= 126 {
			asciiText = append(asciiText, data[i])
		} else if len(asciiText) > 4 {
			fmt.Printf("%s   ASCII at %d: %s\n", prefix, i-len(asciiText), string(asciiText))
			asciiText = nil
		} else {
			asciiText = nil
		}
	}
	if len(asciiText) > 4 {
		fmt.Printf("%s   ASCII at end: %s\n", prefix, string(asciiText))
	}
	
	// Pattern 2: UTF-16 LE (common in Windows)
	for i := 0; i < len(data)-20; i++ {
		// Check for ASCII characters with null bytes between (UTF-16 LE pattern)
		if data[i] >= 32 && data[i] <= 126 && data[i+1] == 0 &&
		   data[i+2] >= 32 && data[i+2] <= 126 && data[i+3] == 0 &&
		   data[i+4] >= 32 && data[i+4] <= 126 && data[i+5] == 0 {
			// Found potential UTF-16 LE text
			end := i
			for j := i; j < len(data)-1; j += 2 {
				if data[j] == 0 && data[j+1] == 0 {
					end = j
					break
				}
			}
			if end > i+6 {
				// Extract UTF-16 LE text
				var utf16Text []byte
				for k := i; k < end; k += 2 {
					if data[k] != 0 {
						utf16Text = append(utf16Text, data[k])
					}
				}
				fmt.Printf("%s   UTF-16 LE at %d: %s\n", prefix, i, string(utf16Text))
				i = end // Skip past this text
			}
		}
	}
	
	// Pattern 3: Look for specific markers
	markers := []string{"TEXT", "text", "tdbs", "Utf8", "utf8"}
	for _, marker := range markers {
		if idx := strings.Index(string(data), marker); idx >= 0 {
			fmt.Printf("%s   Found marker '%s' at offset %d\n", prefix, marker, idx)
			// Show some context after the marker
			contextStart := idx + len(marker)
			contextEnd := contextStart + 50
			if contextEnd > len(data) {
				contextEnd = len(data)
			}
			if contextStart < contextEnd {
				context := data[contextStart:contextEnd]
				fmt.Printf("%s     Context: %s\n", prefix, hex.EncodeToString(context))
			}
		}
	}
}