package main

import (
	"fmt"
	"path/filepath"
	
	aep "github.com/mojosolo/mobot2025"
)

func main() {
	fmt.Println("🎬 Parsing Your Sample After Effects Project")
	fmt.Println("==========================================")
	
	// Parse your sample AEP file
	sampleFile := "sample-aep/Ai Text Intro.aep"
	fmt.Printf("📁 Project: %s\n\n", filepath.Base(sampleFile))
	
	project, err := aep.Open(sampleFile)
	if err != nil {
		fmt.Printf("❌ Error parsing AEP file: %v\n", err)
		fmt.Println("\nTroubleshooting:")
		fmt.Println("1. Make sure you're running from the mobot2025 directory")
		fmt.Println("2. The file path should be: sample-aep/Ai Text Intro.aep")
		return
	}
	
	// Display project information
	fmt.Println("✅ Successfully parsed project!")
	
	fmt.Println("📋 Project Details:")
	fmt.Printf("   - Bit Depth: ")
	switch project.Depth {
	case aep.BPC8:
		fmt.Println("8-bit")
	case aep.BPC16:
		fmt.Println("16-bit")
	case aep.BPC32:
		fmt.Println("32-bit (float)")
	default:
		fmt.Printf("Unknown (%v)\n", project.Depth)
	}
	
	if project.ExpressionEngine != "" {
		fmt.Printf("   - Expression Engine: %s\n", project.ExpressionEngine)
	}
	fmt.Printf("   - Total Items: %d\n\n", len(project.Items))
	
	// Count and categorize items
	var folders, comps, footage int
	var totalLayers int
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeFolder:
			folders++
		case aep.ItemTypeComposition:
			comps++
			totalLayers += len(item.CompositionLayers)
		case aep.ItemTypeFootage:
			footage++
		}
	}
	
	fmt.Println("📊 Project Structure:")
	fmt.Printf("   - Folders: %d\n", folders)
	fmt.Printf("   - Compositions: %d\n", comps)
	fmt.Printf("   - Footage Items: %d (likely includes your BG videos and logo)\n", footage)
	fmt.Printf("   - Total Layers: %d\n\n", totalLayers)
	
	// Show folder structure
	fmt.Println("📂 Folder Hierarchy:")
	printFolderStructure(project.RootFolder, 1)
	
	// Show composition details
	if comps > 0 {
		fmt.Println("\n🎥 Compositions:")
		compNum := 1
		for _, item := range project.Items {
			if item.ItemType == aep.ItemTypeComposition {
				fmt.Printf("\n   %d. %s\n", compNum, item.Name)
				fmt.Printf("      - Dimensions: %dx%d\n", 
					item.FootageDimensions[0], 
					item.FootageDimensions[1])
				fmt.Printf("      - Framerate: %.2f fps\n", item.FootageFramerate)
				fmt.Printf("      - Duration: %.2f seconds (%.0f frames)\n", 
					item.FootageSeconds,
					item.FootageSeconds * item.FootageFramerate)
				fmt.Printf("      - Background: RGB(%d,%d,%d)\n",
					item.BackgroundColor[0],
					item.BackgroundColor[1], 
					item.BackgroundColor[2])
				fmt.Printf("      - Layers: %d\n", len(item.CompositionLayers))
				
				// Show first few layers
				if len(item.CompositionLayers) > 0 {
					fmt.Println("      - Layer Stack (top to bottom):")
					for i, layer := range item.CompositionLayers {
						if i < 5 { // Show first 5 layers
							fmt.Printf("         %d. %s", layer.Index, layer.Name)
							if layer.ThreeDEnabled {
								fmt.Print(" [3D]")
							}
							if layer.SoloEnabled {
								fmt.Print(" [SOLO]")
							}
							if layer.AdjustmentLayerEnabled {
								fmt.Print(" [Adjustment]")
							}
							fmt.Println()
						}
					}
					if len(item.CompositionLayers) > 5 {
						fmt.Printf("         ... and %d more layers\n", 
							len(item.CompositionLayers)-5)
					}
				}
				compNum++
			}
		}
	}
	
	// Show footage items (your assets)
	if footage > 0 {
		fmt.Println("\n📹 Footage/Assets:")
		footageNum := 1
		for _, item := range project.Items {
			if item.ItemType == aep.ItemTypeFootage {
				fmt.Printf("\n   %d. %s\n", footageNum, item.Name)
				if item.FootageDimensions[0] > 0 {
					fmt.Printf("      - Dimensions: %dx%d\n",
						item.FootageDimensions[0],
						item.FootageDimensions[1])
				}
				if item.FootageFramerate > 0 {
					fmt.Printf("      - Framerate: %.2f fps\n", item.FootageFramerate)
				}
				if item.FootageSeconds > 0 {
					fmt.Printf("      - Duration: %.2f seconds\n", item.FootageSeconds)
				}
				
				// Identify type
				switch item.FootageType {
				case aep.FootageTypeSolid:
					fmt.Println("      - Type: Solid Color")
				case aep.FootageTypePlaceholder:
					fmt.Println("      - Type: Placeholder")
				default:
					// Likely your video/image files
					if item.Name == "Reference Logo.png" {
						fmt.Println("      - Type: Image File (PNG)")
					} else if item.FootageFramerate > 0 {
						fmt.Println("      - Type: Video File")
					}
				}
				footageNum++
			}
		}
	}
	
	fmt.Println("\n✨ Analysis complete!")
	fmt.Println("\n💡 Note: The parser found the AEP structure. Your actual media files")
	fmt.Println("   (BG03.mp4, BG04.mp4, BG05.mp4, Reference Logo.png) are referenced")
	fmt.Println("   but not embedded in the AEP file itself.")
}

func printFolderStructure(folder *aep.Item, depth int) {
	if folder == nil {
		return
	}
	
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "   "
	}
	
	if folder.Name == "root" {
		fmt.Printf("%s📁 Project Root\n", indent)
	} else {
		fmt.Printf("%s📁 %s\n", indent, folder.Name)
	}
	
	// Print contents
	for _, item := range folder.FolderContents {
		if item.ItemType == aep.ItemTypeFolder {
			printFolderStructure(item, depth+1)
		} else {
			itemType := "📄"
			if item.ItemType == aep.ItemTypeComposition {
				itemType = "🎬"
			} else if item.ItemType == aep.ItemTypeFootage {
				itemType = "🎞️"
			}
			fmt.Printf("%s   %s %s\n", indent, itemType, item.Name)
		}
	}
}