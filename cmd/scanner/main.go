package main

import (
	"fmt"
	"path/filepath"
	
	aep "github.com/yourusername/mobot2025"
)

func main() {
	fmt.Println("🔍 Scanning all AEP files in data directory...")
	fmt.Println("============================================\n")
	
	testFiles := []string{
		"data/BPC-8.aep",
		"data/BPC-16.aep", 
		"data/BPC-32.aep",
		"data/ExEn-es.aep",
		"data/ExEn-js.aep",
		"data/Item-01.aep",
		"data/Layer-01.aep",
		"data/Property-01.aep",
	}
	
	for _, testFile := range testFiles {
		fmt.Printf("📄 %s:\n", filepath.Base(testFile))
		
		project, err := aep.Open(testFile)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n\n", err)
			continue
		}
		
		// Count items by type
		var folders, comps, footage int
		for _, item := range project.Items {
			switch item.ItemType {
			case aep.ItemTypeFolder:
				folders++
			case aep.ItemTypeComposition:
				comps++
			case aep.ItemTypeFootage:
				footage++
			}
		}
		
		fmt.Printf("   - Bit Depth: ")
		switch project.Depth {
		case aep.BPC8:
			fmt.Print("8-bit")
		case aep.BPC16:
			fmt.Print("16-bit")
		case aep.BPC32:
			fmt.Print("32-bit")
		default:
			fmt.Printf("Unknown (%v)", project.Depth)
		}
		fmt.Println()
		
		if project.ExpressionEngine != "" {
			fmt.Printf("   - Expression Engine: %s\n", project.ExpressionEngine)
		}
		fmt.Printf("   - Items: %d (Folders: %d, Comps: %d, Footage: %d)\n", 
			len(project.Items), folders, comps, footage)
		
		// Show composition details if any
		for _, item := range project.Items {
			if item.ItemType == aep.ItemTypeComposition {
				fmt.Printf("   - Comp '%s': %dx%d @ %.2ffps, %d layers\n",
					item.Name,
					item.FootageDimensions[0],
					item.FootageDimensions[1],
					item.FootageFramerate,
					len(item.CompositionLayers))
			}
		}
		
		fmt.Println()
	}
}