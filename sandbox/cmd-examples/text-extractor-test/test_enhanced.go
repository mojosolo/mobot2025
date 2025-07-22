package main

import (
	"fmt"
	"os"
	"strings"
	
	aep "github.com/mojosolo/mobot2025"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_enhanced_extraction.go <aep-file>")
		os.Exit(1)
	}
	
	aepFile := os.Args[1]
	fmt.Printf("🔍 Testing enhanced text extraction for: %s\n\n", aepFile)
	
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	textCount := 0
	extractedCount := 0
	
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			fmt.Printf("📄 Composition: %s\n", item.Name)
			
			for _, layer := range item.CompositionLayers {
				if layer.Text != nil {
					textCount++
					fmt.Printf("\n  🔤 Text Layer: %s (Index: %d)\n", layer.Name, layer.Index)
					
					// Extract text using enhanced method
					doc, err := aep.ExtractTextContent(layer.Text)
					if err != nil {
						fmt.Printf("     ❌ Error extracting text: %v\n", err)
						continue
					}
					
					if doc != nil {
						fmt.Printf("     📝 Extracted Text: \"%s\"\n", doc.Text)
						
						// Check if we got real text or placeholder
						if !strings.Contains(doc.Text, "[") && !strings.Contains(doc.Text, "keyframes") {
							extractedCount++
							fmt.Printf("     ✅ Successfully extracted real text content!\n")
						} else {
							fmt.Printf("     ⚠️  Text appears to be placeholder or in keyframes\n")
						}
						
						if doc.FontName != "Arial" || doc.FontSize != 12.0 {
							fmt.Printf("     🎨 Font: %s, Size: %.1f\n", doc.FontName, doc.FontSize)
						}
					}
				}
			}
		}
	}
	
	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("   - Total text layers found: %d\n", textCount)
	fmt.Printf("   - Successfully extracted text: %d\n", extractedCount)
	fmt.Printf("   - Success rate: %.1f%%\n", float64(extractedCount)/float64(textCount)*100)
	
	if extractedCount == 0 && textCount > 0 {
		fmt.Println("\n💡 Tip: Text content might be stored in keyframes, expressions, or use a different encoding.")
		fmt.Println("   The enhanced parser tried multiple strategies but couldn't extract the text.")
	}
}