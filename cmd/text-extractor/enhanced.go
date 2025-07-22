package main

import (
	"fmt"
	"os"
	"strings"
	
	aep "github.com/yourusername/mobot2025"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_enhanced_text_extraction.go <aep-file>")
		os.Exit(1)
	}
	
	aepFile := os.Args[1]
	fmt.Printf("🔍 Testing enhanced text extraction for: %s\n\n", aepFile)
	
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("❌ Error opening file: %v\n", err)
		os.Exit(1)
	}
	
	totalTextLayers := 0
	extractedTexts := 0
	
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			fmt.Printf("📄 Composition: %s\n", item.Name)
			fmt.Printf("   Resolution: %dx%d @ %.2f fps\n", item.Width, item.Height, item.FrameRate)
			fmt.Printf("   Duration: %.2f seconds\n\n", item.Duration)
			
			for _, layer := range item.CompositionLayers {
				if layer.Text != nil {
					totalTextLayers++
					fmt.Printf("  🔤 Text Layer #%d: %s\n", layer.Index, layer.Name)
					
					// Try standard extraction first
					standardDoc, err1 := aep.ExtractTextContent(layer.Text)
					if err1 == nil && standardDoc != nil && standardDoc.Text != "" {
						fmt.Printf("     ✅ Standard extraction: \"%s\"\n", standardDoc.Text)
						extractedTexts++
					} else {
						fmt.Printf("     ❌ Standard extraction failed\n")
					}
					
					// Try enhanced extraction
					enhancedDoc, err2 := aep.ExtractEnhancedTextContent(layer.Text)
					if err2 == nil && enhancedDoc != nil && enhancedDoc.Text != "" {
						fmt.Printf("     ✅ Enhanced extraction: \"%s\"\n", enhancedDoc.Text)
						if err1 != nil || standardDoc == nil || standardDoc.Text == "" {
							extractedTexts++
						}
						
						// Show additional info if available
						if enhancedDoc.FontName != "Arial" || enhancedDoc.FontSize != 12.0 {
							fmt.Printf("     📝 Font: %s @ %.1fpt\n", enhancedDoc.FontName, enhancedDoc.FontSize)
						}
						if enhancedDoc.IsAnimated {
							fmt.Printf("     🎬 Animated text detected\n")
						}
					} else {
						fmt.Printf("     ❌ Enhanced extraction failed\n")
						
						// Debug: Show property structure
						debugPropertyStructure(layer.Text, 5)
					}
					
					fmt.Println()
				}
			}
		}
	}
	
	// Summary
	fmt.Printf("\n📊 Summary:\n")
	fmt.Printf("   Total text layers found: %d\n", totalTextLayers)
	fmt.Printf("   Successfully extracted: %d\n", extractedTexts)
	fmt.Printf("   Success rate: %.1f%%\n", float64(extractedTexts)/float64(totalTextLayers)*100)
}

func debugPropertyStructure(prop *aep.Property, indent int) {
	spaces := strings.Repeat(" ", indent)
	
	fmt.Printf("%s🔍 Debug - Property Structure:\n", spaces)
	fmt.Printf("%s  MatchName: %s\n", spaces, prop.MatchName)
	fmt.Printf("%s  Name: %s\n", spaces, prop.Name)
	if prop.Label != "" && prop.Label != "-_0_/-" {
		fmt.Printf("%s  Label: %s\n", spaces, prop.Label)
	}
	
	if prop.RawData != nil && len(prop.RawData) > 0 {
		fmt.Printf("%s  RawData: %d bytes\n", spaces, len(prop.RawData))
		
		// Show hex dump of first 64 bytes
		fmt.Printf("%s  Hex dump (first 64 bytes):\n", spaces)
		for i := 0; i < len(prop.RawData) && i < 64; i += 16 {
			fmt.Printf("%s    %04x: ", spaces, i)
			
			// Hex values
			for j := 0; j < 16 && i+j < len(prop.RawData) && i+j < 64; j++ {
				fmt.Printf("%02x ", prop.RawData[i+j])
			}
			
			// ASCII representation
			fmt.Printf(" |")
			for j := 0; j < 16 && i+j < len(prop.RawData) && i+j < 64; j++ {
				b := prop.RawData[i+j]
				if b >= 32 && b <= 126 {
					fmt.Printf("%c", b)
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("|\n")
		}
	}
	
	if len(prop.Properties) > 0 {
		fmt.Printf("%s  Child properties: %d\n", spaces, len(prop.Properties))
		for _, child := range prop.Properties {
			fmt.Printf("%s    - %s (%s)\n", spaces, child.MatchName, child.Name)
		}
	}
}