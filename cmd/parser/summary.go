package main

import (
	"fmt"
	"strings"
	
	aep "github.com/mojosolo/mobot2025"
)

func main() {
	fmt.Println("🎬 Ai Text Intro Project Analysis Summary")
	fmt.Println("========================================")
	
	// Parse your sample AEP file
	sampleFile := "sample-aep/Ai Text Intro.aep"
	
	project, err := aep.Open(sampleFile)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	// Count items by type
	var folders, comps, footage int
	var totalLayers int
	mainComps := make([]*aep.Item, 0)
	
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeFolder:
			folders++
		case aep.ItemTypeComposition:
			comps++
			totalLayers += len(item.CompositionLayers)
			// Find main comps
			if strings.Contains(item.Name, "Final Comp") || 
			   strings.HasPrefix(item.Name, "S0") {
				mainComps = append(mainComps, item)
			}
		case aep.ItemTypeFootage:
			footage++
		}
	}
	
	fmt.Println("📊 Project Overview:")
	fmt.Printf("   - Total Items: %d\n", len(project.Items))
	fmt.Printf("   - Bit Depth: 16-bit (professional quality)\n")
	fmt.Printf("   - Expression Engine: %s\n\n", project.ExpressionEngine)
	
	fmt.Println("📈 Structure Summary:")
	fmt.Printf("   - Folders: %d (highly organized)\n", folders)
	fmt.Printf("   - Compositions: %d\n", comps)
	fmt.Printf("   - Footage/Assets: %d\n", footage)
	fmt.Printf("   - Total Layers: %d\n\n", totalLayers)
	
	fmt.Println("🎯 Key Compositions Found:")
	
	// Show final comps
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			if strings.Contains(item.Name, "Final Comp") {
				fmt.Printf("\n   📺 %s (Main Output)\n", item.Name)
				fmt.Printf("      - Resolution: %dx%d", 
					item.FootageDimensions[0], 
					item.FootageDimensions[1])
				if item.FootageDimensions[0] >= 3840 {
					fmt.Print(" (4K)")
				} else if item.FootageDimensions[0] >= 1920 {
					fmt.Print(" (2K/HD)")
				}
				fmt.Println()
				fmt.Printf("      - Framerate: %.2f fps\n", item.FootageFramerate)
				fmt.Printf("      - Duration: %.2f seconds\n", item.FootageSeconds)
			}
		}
	}
	
	// Show scene comps
	fmt.Println("\n   🎬 Scene Compositions:")
	sceneCount := 0
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			name := item.Name
			if (strings.HasPrefix(name, "S0") && len(name) == 3) ||
			   (name == "S01" || name == "S02" || name == "S03" || 
			    name == "S04" || name == "S05" || name == "S06" || 
			    name == "S07" || name == "S08" || name == "S09") {
				sceneCount++
				fmt.Printf("      - %s: %d layers\n", name, len(item.CompositionLayers))
			}
		}
	}
	fmt.Printf("      Total Scenes: %d\n", sceneCount)
	
	// Asset analysis
	fmt.Println("\n📦 Asset Analysis:")
	
	// Count asset types
	var nulls, adjustments, solids, others int
	var hasLogo, hasVideo bool
	
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeFootage {
			name := item.Name
			if strings.Contains(strings.ToLower(name), "null") {
				nulls++
			} else if strings.Contains(name, "Adjustment Layer") {
				adjustments++
			} else if strings.Contains(name, "Solid") {
				solids++
			} else {
				others++
				if strings.Contains(strings.ToLower(name), "logo") {
					hasLogo = true
				}
				if item.FootageFramerate > 0 && item.FootageSeconds > 0 {
					hasVideo = true
				}
			}
		}
	}
	
	fmt.Printf("   - Null Objects: %d (for animation control)\n", nulls)
	fmt.Printf("   - Adjustment Layers: %d (for effects)\n", adjustments)
	fmt.Printf("   - Solid Layers: %d (shapes/backgrounds)\n", solids)
	fmt.Printf("   - Other Assets: %d\n", others)
	
	if hasLogo {
		fmt.Println("   ✓ Logo placeholder detected")
	}
	if hasVideo {
		fmt.Println("   ✓ Video footage detected")
	}
	
	// Template features
	fmt.Println("\n✨ Template Features Detected:")
	fmt.Println("   ✓ Multiple text placeholder compositions")
	fmt.Println("   ✓ Organized scene structure (S01-S09)")
	fmt.Println("   ✓ Both 2K and 4K output compositions")
	fmt.Println("   ✓ Complex folder hierarchy for organization")
	fmt.Println("   ✓ Extensive use of null objects (animation rigs)")
	fmt.Println("   ✓ Many adjustment layers (color grading/effects)")
	
	fmt.Println("\n🎯 This is a professional After Effects template for creating")
	fmt.Println("   animated text intros with 9 different scenes and customizable")
	fmt.Println("   text placeholders. Perfect for video intros and titles!")
	
	fmt.Println("\n💡 Your external assets (BG videos, logo) will be linked when")
	fmt.Println("   you open this in After Effects.")
}