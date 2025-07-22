package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	aep "github.com/yourusername/mobot2025"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_text_extraction_final.go <aep-file-or-directory>")
		os.Exit(1)
	}
	
	path := os.Args[1]
	
	// Check if path is a directory or file
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	if info.IsDir() {
		// Process all AEP files in directory
		processDirectory(path)
	} else {
		// Process single file
		processFile(path)
	}
}

func processDirectory(dir string) {
	fmt.Printf("📁 Processing directory: %s\n\n", dir)
	
	files, err := filepath.Glob(filepath.Join(dir, "*.aep"))
	if err != nil {
		fmt.Printf("❌ Error scanning directory: %v\n", err)
		return
	}
	
	if len(files) == 0 {
		fmt.Println("❌ No AEP files found in directory")
		return
	}
	
	totalFiles := len(files)
	filesWithText := 0
	totalTextLayers := 0
	extractedTexts := 0
	
	for _, file := range files {
		fmt.Printf("──────────────────────────────────────────\n")
		stats := processFile(file)
		if stats.hasText {
			filesWithText++
		}
		totalTextLayers += stats.textLayers
		extractedTexts += stats.extracted
	}
	
	// Summary
	fmt.Printf("\n══════════════════════════════════════════\n")
	fmt.Printf("📊 Directory Summary:\n")
	fmt.Printf("   Total AEP files: %d\n", totalFiles)
	fmt.Printf("   Files with text layers: %d\n", filesWithText)
	fmt.Printf("   Total text layers: %d\n", totalTextLayers)
	fmt.Printf("   Successfully extracted: %d\n", extractedTexts)
	if totalTextLayers > 0 {
		fmt.Printf("   Success rate: %.1f%%\n", float64(extractedTexts)/float64(totalTextLayers)*100)
	}
}

type FileStats struct {
	hasText    bool
	textLayers int
	extracted  int
}

func processFile(aepFile string) FileStats {
	stats := FileStats{}
	
	fmt.Printf("🎬 File: %s\n", filepath.Base(aepFile))
	
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n\n", err)
		return stats
	}
	
	// Extract all text layers
	textByComp := aep.ExtractAllTextLayers(project)
	
	if len(textByComp) == 0 {
		fmt.Printf("   ❌ No text layers found\n\n")
		return stats
	}
	
	stats.hasText = true
	
	for compName, textDocs := range textByComp {
		fmt.Printf("\n   📄 Composition: %s\n", compName)
		
		for i, doc := range textDocs {
			stats.textLayers++
			
			if doc.Text != "" && !isPlaceholderText(doc.Text) {
				stats.extracted++
				fmt.Printf("      ✅ Text %d: \"%s\"\n", i+1, truncateText(doc.Text, 50))
				
				// Show font info if not default
				if doc.FontName != "Arial" || doc.FontSize != 12.0 {
					fmt.Printf("         Font: %s @ %.1fpt\n", doc.FontName, doc.FontSize)
				}
			} else {
				fmt.Printf("      ❌ Text %d: (empty or not extracted)\n", i+1)
			}
		}
	}
	
	fmt.Printf("\n   Summary: %d/%d texts extracted (%.0f%%)\n\n", 
		stats.extracted, stats.textLayers, 
		float64(stats.extracted)/float64(stats.textLayers)*100)
	
	return stats
}

func isPlaceholderText(text string) bool {
	return text == "[Text content not extracted - check keyframes/expressions]" ||
	       text == "[Text content in keyframes/expressions]"
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}