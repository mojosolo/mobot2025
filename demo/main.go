package main

import (
	"fmt"
	"log"
	"path/filepath"
	
	aep "github.com/mojosolo/mobot2025"
)

func main() {
	fmt.Println("🎬 Adobe After Effects AEP Parser Demo")
	fmt.Println("=====================================")
	
	// Parse a test AEP file with more content
	testFile := "data/Layer-01.aep"  // This file has 17 layers!
	fmt.Printf("Parsing: %s\n\n", filepath.Base(testFile))
	
	project, err := aep.Open(testFile)
	if err != nil {
		log.Fatal("Error parsing AEP file:", err)
	}
	
	// Display project information
	fmt.Println("📋 Project Details:")
	fmt.Printf("   - Bit Depth: ")
	switch project.Depth {
	case aep.BPC8:
		fmt.Println("8-bit")
	case aep.BPC16:
		fmt.Println("16-bit")
	case aep.BPC32:
		fmt.Println("32-bit")
	default:
		fmt.Printf("Unknown (%v)\n", project.Depth)
	}
	
	if project.ExpressionEngine != "" {
		fmt.Printf("   - Expression Engine: %s\n", project.ExpressionEngine)
	}
	fmt.Printf("   - Total Items: %d\n\n", len(project.Items))
	
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
	
	fmt.Println("📊 Item Summary:")
	fmt.Printf("   - Folders: %d\n", folders)
	fmt.Printf("   - Compositions: %d\n", comps)
	fmt.Printf("   - Footage Items: %d\n", footage)
	
	// Show composition details
	if comps > 0 {
		fmt.Println("\n🎥 Compositions:")
		for _, item := range project.Items {
			if item.ItemType == aep.ItemTypeComposition {
				fmt.Printf("\n   Composition: %s\n", item.Name)
				fmt.Printf("      - ID: %d\n", item.ID)
				fmt.Printf("      - Dimensions: %dx%d\n", 
					item.FootageDimensions[0], 
					item.FootageDimensions[1])
				fmt.Printf("      - Framerate: %.2f fps\n", item.FootageFramerate)
				fmt.Printf("      - Duration: %.2f seconds\n", item.FootageSeconds)
				fmt.Printf("      - Background Color: RGB(%d, %d, %d)\n",
					item.BackgroundColor[0],
					item.BackgroundColor[1], 
					item.BackgroundColor[2])
				fmt.Printf("      - Layers: %d\n", len(item.CompositionLayers))
				
				// Show layer details
				if len(item.CompositionLayers) > 0 {
					fmt.Println("      - Layer Stack:")
					for i, layer := range item.CompositionLayers {
						fmt.Printf("         %d. %s", layer.Index, layer.Name)
						if layer.SoloEnabled {
							fmt.Print(" [SOLO]")
						}
						if layer.ThreeDEnabled {
							fmt.Print(" [3D]")
						}
						fmt.Println()
						if i >= 2 && len(item.CompositionLayers) > 3 {
							fmt.Printf("         ... and %d more layers\n", 
								len(item.CompositionLayers)-3)
							break
						}
					}
				}
			}
		}
	}
	
	// Show footage details
	if footage > 0 {
		fmt.Println("\n📁 Footage Items:")
		count := 0
		for _, item := range project.Items {
			if item.ItemType == aep.ItemTypeFootage {
				fmt.Printf("\n   %s\n", item.Name)
				fmt.Printf("      - ID: %d\n", item.ID)
				fmt.Printf("      - Dimensions: %dx%d\n",
					item.FootageDimensions[0],
					item.FootageDimensions[1])
				fmt.Printf("      - Framerate: %.2f fps\n", item.FootageFramerate)
				fmt.Printf("      - Duration: %.2f seconds\n", item.FootageSeconds)
				
				switch item.FootageType {
				case aep.FootageTypeSolid:
					fmt.Println("      - Type: Solid Color")
				case aep.FootageTypePlaceholder:
					fmt.Println("      - Type: Placeholder")
				default:
					fmt.Printf("      - Type: Other (%v)\n", item.FootageType)
				}
				
				count++
				if count >= 3 && footage > 3 {
					fmt.Printf("\n   ... and %d more footage items\n", footage-3)
					break
				}
			}
		}
	}
	
	fmt.Println("\n✅ Parsing complete!")
}