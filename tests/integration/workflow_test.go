package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/mobot2025/catalog"
	"github.com/yourusername/mobot2025/tests/helpers"
)

// TestCompleteAEPWorkflow tests the complete AEP processing workflow
func TestCompleteAEPWorkflow(t *testing.T) {
	// Setup test environment
	db := helpers.CreateTestDatabase(t)
	api := catalog.NewAPIService(db)
	server := httptest.NewServer(api.Handler())
	defer server.Close()

	// Step 1: Upload AEP file
	t.Run("upload_aep_file", func(t *testing.T) {
		// Create multipart form with AEP file
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		
		file, err := os.Open(helpers.TestFixtures.ValidAEP)
		helpers.AssertNoError(t, err, "Failed to open test file")
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(helpers.TestFixtures.ValidAEP))
		helpers.AssertNoError(t, err, "Failed to create form file")
		
		_, err = io.Copy(part, file)
		helpers.AssertNoError(t, err, "Failed to copy file content")
		
		writer.Close()

		// Send upload request
		resp, err := http.Post(
			server.URL+"/api/upload",
			writer.FormDataContentType(),
			&buf,
		)
		helpers.AssertNoError(t, err, "Upload request failed")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Upload response status")

		// Parse response
		var uploadResp struct {
			ProjectID string `json:"project_id"`
			Message   string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&uploadResp)
		helpers.AssertNoError(t, err, "Failed to decode upload response")

		// Store project ID for next steps
		t.Setenv("TEST_PROJECT_ID", uploadResp.ProjectID)
	})

	// Step 2: Query project metadata
	t.Run("query_project_metadata", func(t *testing.T) {
		projectID := os.Getenv("TEST_PROJECT_ID")
		if projectID == "" {
			t.Skip("No project ID from upload step")
		}

		resp, err := http.Get(server.URL + "/api/projects/" + projectID)
		helpers.AssertNoError(t, err, "Failed to query project")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Query response status")

		var project catalog.ProjectMetadata
		err = json.NewDecoder(resp.Body).Decode(&project)
		helpers.AssertNoError(t, err, "Failed to decode project metadata")

		// Verify project data
		if project.Name == "" {
			t.Error("Project name is empty")
		}
		if project.Version == "" {
			t.Error("Project version is empty")
		}
	})

	// Step 3: Search for text layers
	t.Run("search_text_layers", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/search?type=text&query=title")
		helpers.AssertNoError(t, err, "Failed to search text layers")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Search response status")

		var results struct {
			Layers []catalog.TextLayerInfo `json:"layers"`
			Total  int                     `json:"total"`
		}
		err = json.NewDecoder(resp.Body).Decode(&results)
		helpers.AssertNoError(t, err, "Failed to decode search results")

		if results.Total == 0 {
			t.Log("No text layers found matching 'title'")
		}
	})

	// Step 4: Generate analysis report
	t.Run("generate_analysis_report", func(t *testing.T) {
		projectID := os.Getenv("TEST_PROJECT_ID")
		if projectID == "" {
			t.Skip("No project ID from upload step")
		}

		resp, err := http.Post(
			server.URL+"/api/analyze/"+projectID,
			"application/json",
			nil,
		)
		helpers.AssertNoError(t, err, "Failed to request analysis")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Analysis response status")

		var analysis struct {
			ProjectID       string   `json:"project_id"`
			Opportunities   []string `json:"opportunities"`
			Recommendations []string `json:"recommendations"`
			Score           float64  `json:"automation_score"`
		}
		err = json.NewDecoder(resp.Body).Decode(&analysis)
		helpers.AssertNoError(t, err, "Failed to decode analysis")

		if analysis.Score < 0 || analysis.Score > 100 {
			t.Errorf("Invalid automation score: %f", analysis.Score)
		}
	})
}

// TestMultiAgentCoordination tests multi-agent system coordination
func TestMultiAgentCoordination(t *testing.T) {
	// Setup agents
	orchestrator := catalog.NewMetaOrchestrator()
	
	// Create test project
	project := helpers.CreateMockProject("Multi-Agent Test")
	
	// Step 1: Planning phase
	t.Run("planning_phase", func(t *testing.T) {
		plan, err := orchestrator.CreateProjectPlan(project)
		helpers.AssertNoError(t, err, "Planning failed")
		
		if len(plan.Tasks) == 0 {
			t.Error("No tasks generated in plan")
		}
		
		// Verify task dependencies
		for _, task := range plan.Tasks {
			for _, dep := range task.Dependencies {
				found := false
				for _, t := range plan.Tasks {
					if t.ID == dep {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Task %s has invalid dependency %s", task.ID, dep)
				}
			}
		}
		
		t.Setenv("TEST_PLAN_ID", plan.ID)
	})
	
	// Step 2: Implementation phase
	t.Run("implementation_phase", func(t *testing.T) {
		planID := os.Getenv("TEST_PLAN_ID")
		if planID == "" {
			t.Skip("No plan ID from planning phase")
		}
		
		results, err := orchestrator.ExecutePlan(planID)
		helpers.AssertNoError(t, err, "Execution failed")
		
		// Check all tasks completed
		for taskID, result := range results {
			if result.Status != "completed" {
				t.Errorf("Task %s not completed: %s", taskID, result.Status)
			}
		}
	})
	
	// Step 3: Review phase
	t.Run("review_phase", func(t *testing.T) {
		planID := os.Getenv("TEST_PLAN_ID")
		if planID == "" {
			t.Skip("No plan ID from planning phase")
		}
		
		review, err := orchestrator.ReviewResults(planID)
		helpers.AssertNoError(t, err, "Review failed")
		
		if review.QualityScore < 0 || review.QualityScore > 1 {
			t.Errorf("Invalid quality score: %f", review.QualityScore)
		}
		
		if len(review.Issues) > 0 {
			t.Logf("Review found %d issues", len(review.Issues))
		}
	})
	
	// Step 4: Verification phase
	t.Run("verification_phase", func(t *testing.T) {
		planID := os.Getenv("TEST_PLAN_ID")
		if planID == "" {
			t.Skip("No plan ID from planning phase")
		}
		
		verified, err := orchestrator.VerifyCompletion(planID)
		helpers.AssertNoError(t, err, "Verification failed")
		
		if !verified {
			t.Error("Project verification failed")
		}
	})
}

// TestErrorRecovery tests system error recovery capabilities
func TestErrorRecovery(t *testing.T) {
	db := helpers.CreateTestDatabase(t)
	api := catalog.NewAPIService(db)
	
	// Test corrupted file handling
	t.Run("corrupted_file_recovery", func(t *testing.T) {
		// Create corrupted data
		corruptedData := []byte("This is not a valid AEP file")
		
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/upload", bytes.NewReader(corruptedData))
		
		api.ServeHTTP(resp, req)
		
		// Should return error but not crash
		if resp.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", resp.Code)
		}
		
		var errResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		
		if errResp.Error == "" {
			t.Error("Expected error message in response")
		}
	})
	
	// Test database connection failure recovery
	t.Run("database_failure_recovery", func(t *testing.T) {
		// Close database to simulate failure
		db.Close()
		
		// Try to perform operation
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/projects", nil)
		
		api.ServeHTTP(resp, req)
		
		// Should handle gracefully
		if resp.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected 503, got %d", resp.Code)
		}
	})
}

// TestPerformanceUnderLoad tests system performance under load
func TestPerformanceUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	db := helpers.CreateTestDatabase(t)
	api := catalog.NewAPIService(db)
	server := httptest.NewServer(api.Handler())
	defer server.Close()
	
	// Generate load with concurrent requests
	concurrency := 10
	requestsPerClient := 100
	
	helpers.RunConcurrentTest(t, concurrency, func(clientID int) {
		for i := 0; i < requestsPerClient; i++ {
			resp, err := http.Get(server.URL + "/api/health")
			if err != nil {
				t.Errorf("Client %d request %d failed: %v", clientID, i, err)
				continue
			}
			resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Client %d request %d got status %d", clientID, i, resp.StatusCode)
			}
		}
	})
}