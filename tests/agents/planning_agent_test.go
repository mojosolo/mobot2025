package agents_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/mobot2025/catalog"
)

// TestPlanningAgentAnalysis tests the planning agent's ability to analyze AEP projects
func TestPlanningAgentAnalysis(t *testing.T) {
	agent := catalog.NewPlanningAgent()

	tests := []struct {
		name        string
		metadata    *catalog.ProjectMetadata
		wantTasks   int
		wantError   bool
		description string
	}{
		{
			name: "simple_project",
			metadata: &catalog.ProjectMetadata{
				Name:         "Simple Animation",
				Version:      "2021",
				Compositions: 1,
				TextLayers:   []catalog.TextLayerInfo{
					{Text: "Hello World", LayerName: "Title"},
				},
			},
			wantTasks:   3, // Basic analysis, text processing, optimization
			wantError:   false,
			description: "Should generate basic tasks for simple project",
		},
		{
			name: "complex_project",
			metadata: &catalog.ProjectMetadata{
				Name:         "Corporate Presentation",
				Version:      "2021",
				Compositions: 10,
				TextLayers:   []catalog.TextLayerInfo{
					{Text: "Title", LayerName: "Main Title"},
					{Text: "Subtitle", LayerName: "Subtitle"},
					{Text: "Body Text", LayerName: "Description"},
				},
				MediaAssets: []catalog.MediaAsset{
					{Name: "background.mp4", Type: "video"},
					{Name: "logo.png", Type: "image"},
				},
			},
			wantTasks:   8, // More tasks for complex project
			wantError:   false,
			description: "Should generate comprehensive tasks for complex project",
		},
		{
			name:        "empty_project",
			metadata:    &catalog.ProjectMetadata{},
			wantTasks:   0,
			wantError:   true,
			description: "Should error on empty project metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := agent.AnalyzeProject(tt.metadata)
			
			if (err != nil) != tt.wantError {
				t.Errorf("AnalyzeProject() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if len(plan.Tasks) != tt.wantTasks {
					t.Errorf("Expected %d tasks, got %d", tt.wantTasks, len(plan.Tasks))
				}

				// Verify plan has essential components
				if plan.ProjectID == "" {
					t.Error("Plan missing project ID")
				}
				if plan.CreatedAt.IsZero() {
					t.Error("Plan missing creation timestamp")
				}
			}
		})
	}
}

// TestPlanningAgentPrioritization tests task prioritization logic
func TestPlanningAgentPrioritization(t *testing.T) {
	agent := catalog.NewPlanningAgent()

	metadata := &catalog.ProjectMetadata{
		Name:         "Priority Test Project",
		Version:      "2021",
		Compositions: 5,
		HasErrors:    true, // This should trigger high-priority error resolution tasks
		TextLayers: []catalog.TextLayerInfo{
			{Text: "Missing Font Text", HasMissingFonts: true},
		},
	}

	plan, err := agent.AnalyzeProject(metadata)
	if err != nil {
		t.Fatalf("Failed to analyze project: %v", err)
	}

	// Find error resolution task
	var errorTask *catalog.Task
	for _, task := range plan.Tasks {
		if task.Type == catalog.TaskTypeErrorResolution {
			errorTask = &task
			break
		}
	}

	if errorTask == nil {
		t.Fatal("Expected error resolution task for project with errors")
	}

	if errorTask.Priority != catalog.PriorityHigh {
		t.Errorf("Error resolution task should have high priority, got %s", errorTask.Priority)
	}
}

// TestPlanningAgentConcurrency tests concurrent project analysis
func TestPlanningAgentConcurrency(t *testing.T) {
	agent := catalog.NewPlanningAgent()
	
	// Create multiple projects to analyze concurrently
	projects := make([]*catalog.ProjectMetadata, 10)
	for i := range projects {
		projects[i] = &catalog.ProjectMetadata{
			Name:         fmt.Sprintf("Concurrent Project %d", i),
			Version:      "2021",
			Compositions: i + 1,
		}
	}

	// Analyze all projects concurrently
	results := make(chan error, len(projects))
	
	for _, project := range projects {
		go func(p *catalog.ProjectMetadata) {
			_, err := agent.AnalyzeProject(p)
			results <- err
		}(project)
	}

	// Collect results
	for i := 0; i < len(projects); i++ {
		select {
		case err := <-results:
			if err != nil {
				t.Errorf("Concurrent analysis failed: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent analysis timed out")
		}
	}
}

// TestPlanningAgentValidation tests input validation
func TestPlanningAgentValidation(t *testing.T) {
	agent := catalog.NewPlanningAgent()

	invalidInputs := []struct {
		name     string
		metadata *catalog.ProjectMetadata
		wantErr  string
	}{
		{
			name:     "nil_metadata",
			metadata: nil,
			wantErr:  "project metadata is nil",
		},
		{
			name: "invalid_version",
			metadata: &catalog.ProjectMetadata{
				Name:    "Test",
				Version: "1999", // Too old
			},
			wantErr: "unsupported After Effects version",
		},
		{
			name: "corrupted_data",
			metadata: &catalog.ProjectMetadata{
				Name:         "Corrupted",
				Compositions: -1, // Invalid composition count
			},
			wantErr: "invalid composition count",
		},
	}

	for _, tt := range invalidInputs {
		t.Run(tt.name, func(t *testing.T) {
			_, err := agent.AnalyzeProject(tt.metadata)
			if err == nil {
				t.Fatal("Expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// BenchmarkPlanningAgentAnalysis benchmarks the planning agent
func BenchmarkPlanningAgentAnalysis(b *testing.B) {
	agent := catalog.NewPlanningAgent()
	
	// Create a complex project for benchmarking
	metadata := &catalog.ProjectMetadata{
		Name:         "Benchmark Project",
		Version:      "2021",
		Compositions: 50,
		TextLayers:   make([]catalog.TextLayerInfo, 100),
		MediaAssets:  make([]catalog.MediaAsset, 200),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.AnalyzeProject(metadata)
		if err != nil {
			b.Fatal(err)
		}
	}
}