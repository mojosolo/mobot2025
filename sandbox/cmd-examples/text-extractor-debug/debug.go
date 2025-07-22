package main

import (
	"fmt"
	"os"
	
	aep "github.com/mojosolo/mobot2025"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_text_extraction.go <aep-file>")
		os.Exit(1)
	}
	
	aepFile := os.Args[1]
	fmt.Printf("🔍 Debugging text extraction for: %s\n\n", aepFile)
	
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	textFound := false
	
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			fmt.Printf("📄 Composition: %s\n", item.Name)
			
			for _, layer := range item.CompositionLayers {
				if layer.Text != nil {
					textFound = true
					fmt.Printf("\n  🔤 Text Layer: %s (Index: %d)\n", layer.Name, layer.Index)
					fmt.Printf("     - Has Text Property: ✅\n")
					
					// Debug the property structure
					debugProperty(layer.Text, 2)
				}
			}
		}
	}
	
	if !textFound {
		fmt.Println("❌ No text layers found in this project")
	}
}

func debugProperty(prop *aep.Property, indent int) {
	spaces := ""
	for i := 0; i < indent*2; i++ {
		spaces += " "
	}
	
	fmt.Printf("%s📦 Property:\n", spaces)
	fmt.Printf("%s   - MatchName: %s\n", spaces, prop.MatchName)
	fmt.Printf("%s   - Name: %s\n", spaces, prop.Name)
	if prop.Label != "" && prop.Label != "-_0_/-" {
		fmt.Printf("%s   - Label: %s\n", spaces, prop.Label)
	}
	fmt.Printf("%s   - Type: %s\n", spaces, prop.PropertyType.String())
	
	if prop.TextDocument != nil {
		fmt.Printf("%s   📝 TextDocument:\n", spaces)
		fmt.Printf("%s      - Text: \"%s\"\n", spaces, prop.TextDocument.Text)
		fmt.Printf("%s      - Font: %s\n", spaces, prop.TextDocument.FontName)
		fmt.Printf("%s      - Size: %.1f\n", spaces, prop.TextDocument.FontSize)
	}
	
	if prop.RawData != nil {
		fmt.Printf("%s   - RawData: %d bytes\n", spaces, len(prop.RawData))
		// Show first 100 chars if it looks like text
		preview := string(prop.RawData)
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("%s   - RawData preview: %q\n", spaces, preview)
	}
	
	if len(prop.SelectOptions) > 0 {
		fmt.Printf("%s   - SelectOptions: %v\n", spaces, prop.SelectOptions)
	}
	
	// Recursively debug child properties
	if len(prop.Properties) > 0 {
		fmt.Printf("%s   - Child Properties: %d\n", spaces, len(prop.Properties))
		for i, child := range prop.Properties {
			fmt.Printf("%s   [%d] %s:\n", spaces, i, child.MatchName)
			debugProperty(child, indent+2)
		}
	}
}