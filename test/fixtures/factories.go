package fixtures

import (
	"fmt"
	"time"
	
	"github.com/mojosolo/mobot2025/catalog"
)

// Factory functions for creating test data

// CreateTestProjectMetadata creates a test ProjectMetadata
func CreateTestProjectMetadata(name string) *catalog.ProjectMetadata {
	return &catalog.ProjectMetadata{
		FilePath:         fmt.Sprintf("/test/%s.aep", name),
		FileName:         fmt.Sprintf("%s.aep", name),
		FileSize:         1024 * 1024, // 1MB
		BitDepth:         8,
		ExpressionEngine: "javascript",
		TotalItems:       10,
		Compositions:     CreateTestCompositions(3),
		TextLayers:       CreateTestTextLayers(5),
		MediaAssets:      CreateTestMediaAssets(3),
		Effects:          CreateTestEffects(2),
		Categories:       []string{"motion-graphics", "title-sequence"},
		Tags:             []string{"corporate", "modern", "animated"},
		Capabilities: catalog.ProjectCapabilities{
			HasTextReplacement:  true,
			HasImageReplacement: true,
			HasColorControl:     true,
			HasAudioReplacement: false,
			HasDataDriven:       false,
			HasExpressions:      true,
			IsModular:           true,
		},
		Opportunities: []catalog.Opportunity{
			{
				Type:        "text_replacement",
				Description: "5 text layers can be automated",
				Difficulty:  "easy",
				Impact:      "high",
				Components:  []string{"text-layer-1", "text-layer-2"},
			},
		},
		ParsedAt: time.Now(),
	}
}

// CreateTestCompositions creates test compositions
func CreateTestCompositions(count int) []catalog.CompositionInfo {
	comps := make([]catalog.CompositionInfo, count)
	for i := 0; i < count; i++ {
		comps[i] = catalog.CompositionInfo{
			ID:         fmt.Sprintf("comp-%d", i+1),
			Name:       fmt.Sprintf("Composition %d", i+1),
			Width:      1920,
			Height:     1080,
			FrameRate:  30,
			Duration:   float64(10 + i*5),
			LayerCount: 5 + i*2,
			Is3D:       i%2 == 0,
			HasEffects: true,
		}
	}
	return comps
}

// CreateTestTextLayers creates test text layers
func CreateTestTextLayers(count int) []catalog.TextLayerInfo {
	layers := make([]catalog.TextLayerInfo, count)
	templates := []struct {
		name string
		text string
	}{
		{"Title", "Main Title Text"},
		{"Subtitle", "Supporting subtitle here"},
		{"Body", "Lorem ipsum dolor sit amet"},
		{"CTA", "Click Here to Learn More"},
		{"Footer", "© 2024 Company Name"},
	}
	
	for i := 0; i < count; i++ {
		template := templates[i%len(templates)]
		layers[i] = catalog.TextLayerInfo{
			ID:             fmt.Sprintf("text-%d", i+1),
			CompID:         fmt.Sprintf("comp-%d", (i%3)+1),
			LayerName:      template.name,
			SourceText:     template.text,
			FontUsed:       "Arial",
			IsAnimated:     i%2 == 0,
			HasExpressions: i%3 == 0,
			Is3D:           i%4 == 0,
		}
	}
	return layers
}

// CreateTestMediaAssets creates test media assets
func CreateTestMediaAssets(count int) []catalog.MediaAssetInfo {
	assets := make([]catalog.MediaAssetInfo, count)
	types := []string{"image", "video", "audio"}
	
	for i := 0; i < count; i++ {
		assetType := types[i%len(types)]
		assets[i] = catalog.MediaAssetInfo{
			ID:            fmt.Sprintf("asset-%d", i+1),
			Name:          fmt.Sprintf("%s-asset-%d", assetType, i+1),
			Type:          assetType,
			Path:          fmt.Sprintf("/assets/%s-%d.%s", assetType, i+1, getExtension(assetType)),
			IsPlaceholder: i%3 == 0,
			UsageCount:    1 + i,
		}
	}
	return assets
}

// CreateTestEffects creates test effects
func CreateTestEffects(count int) []catalog.EffectInfo {
	effects := make([]catalog.EffectInfo, count)
	effectTypes := []string{"Gaussian Blur", "Drop Shadow", "Glow", "Color Correction"}
	
	for i := 0; i < count; i++ {
		effects[i] = catalog.EffectInfo{
			Name:           effectTypes[i%len(effectTypes)],
			Category:       "visual",
			UsageCount:     1 + i,
			IsCustomizable: i%2 == 0,
		}
	}
	return effects
}

// CreateTestPlanningResult creates a test planning result
func CreateTestPlanningResult(projectID string) *catalog.PlanningResult {
	return &catalog.PlanningResult{
		ProjectID:      projectID,
		TotalTasks:     10,
		HighPriority:   3,
		AutoExecutable: 7,
		EstimatedTotal: time.Hour * 2,
		ConfidenceAvg:  0.85,
		Recommendations: []string{
			"Automate text replacements",
			"Use templates for similar compositions",
			"Batch process media assets",
		},
		CreatedAt: time.Now(),
	}
}

// CreateTestTaskPlan creates a test task plan
func CreateTestTaskPlan(taskType string, priority int) *catalog.TaskPlan {
	return &catalog.TaskPlan{
		ID:              fmt.Sprintf("task-%s-%d", taskType, time.Now().Unix()),
		Type:            taskType,
		Description:     fmt.Sprintf("Test %s task", taskType),
		FileReferences:  []string{"/test/file1.aep", "/test/file2.aep"},
		Dependencies:    []string{},
		Priority:        priority,
		ConfidenceScore: 0.9,
		EstimatedTime:   time.Minute * 30,
		BlockTypes:      []string{"CompItem", "AVLayer"},
		Metadata: map[string]interface{}{
			"test": true,
			"generated": time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestWorkflow creates a test workflow
func CreateTestWorkflow(name string) *catalog.Workflow {
	return &catalog.Workflow{
		ID:              fmt.Sprintf("workflow-%d", time.Now().Unix()),
		Name:            name,
		Description:     fmt.Sprintf("Test workflow: %s", name),
		ProjectPath:     "/test/project",
		Status:          "pending",
		CurrentStage:    "planning",
		Progress:        0.0,
		Tasks:           []*catalog.WorkflowTask{},
		Dependencies:    make(map[string][]string),
		CompletedTasks:  []string{},
		FailedTasks:     []string{},
		CurrentTasks:    []string{},
		RequiresApproval: false,
		IterationCount:  0,
		CreatedAt:       time.Now(),
		EstimatedTime:   time.Hour,
		ActualTime:      0,
		Metadata:        make(map[string]interface{}),
	}
}

// Helper functions

func getExtension(assetType string) string {
	switch assetType {
	case "image":
		return "png"
	case "video":
		return "mp4"
	case "audio":
		return "wav"
	default:
		return "dat"
	}
}

// Test scenario generators

// GenerateAutomationScenario creates a complete automation test scenario
func GenerateAutomationScenario() (*catalog.ProjectMetadata, *catalog.PlanningResult, []*catalog.TaskPlan) {
	// Create project
	project := CreateTestProjectMetadata("automation-test")
	
	// Create planning result
	planning := CreateTestPlanningResult("test-project-123")
	
	// Create tasks
	tasks := []*catalog.TaskPlan{
		CreateTestTaskPlan("text_replacement", 1),
		CreateTestTaskPlan("composition_setup", 2),
		CreateTestTaskPlan("effect_application", 3),
		CreateTestTaskPlan("render_output", 1),
	}
	
	// Set up dependencies
	tasks[1].Dependencies = []string{tasks[0].ID}
	tasks[2].Dependencies = []string{tasks[1].ID}
	tasks[3].Dependencies = []string{tasks[0].ID, tasks[1].ID, tasks[2].ID}
	
	return project, planning, tasks
}

// GenerateSearchScenario creates test data for search functionality
func GenerateSearchScenario() []*catalog.ProjectMetadata {
	projects := []*catalog.ProjectMetadata{
		CreateTestProjectMetadata("corporate-intro"),
		CreateTestProjectMetadata("product-showcase"),
		CreateTestProjectMetadata("event-highlights"),
		CreateTestProjectMetadata("social-media-pack"),
		CreateTestProjectMetadata("logo-animation"),
	}
	
	// Customize each project
	projects[0].Categories = []string{"corporate", "intro"}
	projects[1].Categories = []string{"product", "marketing"}
	projects[2].Categories = []string{"event", "documentation"}
	projects[3].Categories = []string{"social-media", "template-pack"}
	projects[4].Categories = []string{"branding", "animation"}
	
	// Add varying complexity
	for i, p := range projects {
		p.TotalItems = 5 + i*3
		p.Compositions = CreateTestCompositions(1 + i)
		p.TextLayers = CreateTestTextLayers(3 + i*2)
	}
	
	return projects
}

// GenerateQualityAssuranceScenario creates test data for QA testing
func GenerateQualityAssuranceScenario() *catalog.ProjectMetadata {
	project := CreateTestProjectMetadata("qa-test-project")
	
	// Add various quality issues
	project.TextLayers = append(project.TextLayers, catalog.TextLayerInfo{
		ID:         "problematic-text",
		CompID:     "comp-1",
		LayerName:  "Lorem Ipsum Text",
		SourceText: "Lorem ipsum dolor sit amet", // Placeholder text
		FontUsed:   "Comic Sans MS", // Poor font choice
		IsAnimated: false,
	})
	
	// Add missing fonts
	project.MediaAssets = append(project.MediaAssets, catalog.MediaAssetInfo{
		ID:            "missing-asset",
		Name:          "background.jpg",
		Type:          "image",
		Path:          "", // Missing path
		IsPlaceholder: true,
		UsageCount:    5,
	})
	
	return project
}