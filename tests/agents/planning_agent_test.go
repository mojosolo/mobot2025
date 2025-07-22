package agents_test

import (
	"testing"
	"time"

	"github.com/mojosolo/mobot2025/catalog"
)

// TestPlanningAgentBasic tests basic planning agent functionality
func TestPlanningAgentBasic(t *testing.T) {
	// Create a temporary database for testing
	db, err := catalog.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()
	
	// Create planning agent
	agent := catalog.NewPlanningAgent(db)
	
	// Test with a non-existent file (will fail but that's expected)
	result, err := agent.AnalyzeAndPlan("/tmp/test_nonexistent.aep")
	
	// We expect an error since the file doesn't exist
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	// Result should be nil on error
	if result != nil {
		t.Error("Expected nil result on error")
	}
}

// TestPlanningResultStructure tests the planning result structure
func TestPlanningResultStructure(t *testing.T) {
	// Create a mock planning result to test the structure
	result := &catalog.PlanningResult{
		ProjectID:      "test-project-123",
		TotalTasks:     5,
		HighPriority:   2,
		AutoExecutable: 3,
		EstimatedTotal: time.Hour,
		ConfidenceAvg:  0.85,
		Recommendations: []string{
			"Consider automating text replacements",
			"Optimize render settings",
		},
		CreatedAt: time.Now(),
	}
	
	// Verify structure
	if result.ProjectID == "" {
		t.Error("ProjectID should not be empty")
	}
	
	if result.TotalTasks <= 0 {
		t.Error("TotalTasks should be positive")
	}
	
	if result.ConfidenceAvg < 0 || result.ConfidenceAvg > 1 {
		t.Error("ConfidenceAvg should be between 0 and 1")
	}
	
	if len(result.Recommendations) == 0 {
		t.Error("Should have recommendations")
	}
	
	if result.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

// TestTaskPlanStructure tests the task plan structure
func TestTaskPlanStructure(t *testing.T) {
	task := &catalog.TaskPlan{
		ID:              "task-123",
		Type:            "parse_composition",
		Description:     "Parse main composition",
		FileReferences:  []string{"/tmp/test.aep"},
		Dependencies:    []string{},
		Priority:        1,
		ConfidenceScore: 0.9,
		EstimatedTime:   time.Second * 30,
		BlockTypes:      []string{"CompItem", "Layer"},
		Metadata:        map[string]interface{}{"comp_name": "Main"},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Verify task structure
	if task.ID == "" {
		t.Error("Task ID should not be empty")
	}
	
	if task.Priority < 1 || task.Priority > 3 {
		t.Error("Priority should be between 1 and 3")
	}
	
	if task.ConfidenceScore < 0 || task.ConfidenceScore > 1 {
		t.Error("ConfidenceScore should be between 0 and 1")
	}
	
	if len(task.BlockTypes) == 0 {
		t.Error("Should have block types")
	}
}