package demo_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/mobot2025/demo"
	"github.com/yourusername/mobot2025/tests/helpers"
)

// TestSimpleStoryViewer tests the simple story viewer functionality
func TestSimpleStoryViewer(t *testing.T) {
	// Create viewer instance
	viewer := demo.NewSimpleStoryViewer()
	server := httptest.NewServer(viewer)
	defer server.Close()

	// Test homepage
	t.Run("homepage", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		helpers.AssertNoError(t, err, "Failed to get homepage")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Homepage status")

		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "Simple Story Viewer") {
			t.Error("Homepage doesn't contain expected title")
		}
	})

	// Test file upload
	t.Run("file_upload", func(t *testing.T) {
		// Create multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		file, err := os.Open(helpers.TestFixtures.ValidAEP)
		helpers.AssertNoError(t, err, "Failed to open test file")
		defer file.Close()

		part, err := writer.CreateFormFile("file", "test.aep")
		helpers.AssertNoError(t, err, "Failed to create form file")

		_, err = io.Copy(part, file)
		helpers.AssertNoError(t, err, "Failed to copy file")

		writer.Close()

		// Upload file
		resp, err := http.Post(
			server.URL+"/upload",
			writer.FormDataContentType(),
			&buf,
		)
		helpers.AssertNoError(t, err, "Upload failed")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Upload status")

		// Check response
		var result struct {
			Success  bool   `json:"success"`
			Message  string `json:"message"`
			Redirect string `json:"redirect"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		helpers.AssertNoError(t, err, "Failed to decode response")

		if !result.Success {
			t.Errorf("Upload not successful: %s", result.Message)
		}
	})

	// Test story display
	t.Run("story_display", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/story/test-project")
		helpers.AssertNoError(t, err, "Failed to get story")
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("No test project available")
		}

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Story page status")

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Check for expected elements
		expectedElements := []string{
			"story-container",
			"text-layer",
			"composition-info",
		}

		for _, elem := range expectedElements {
			if !strings.Contains(bodyStr, elem) {
				t.Errorf("Story page missing element: %s", elem)
			}
		}
	})
}

// TestEasyModeViewer tests the easy mode viewer
func TestEasyModeViewer(t *testing.T) {
	viewer := demo.NewEasyModeViewer()
	server := httptest.NewServer(viewer)
	defer server.Close()

	// Test drag and drop upload
	t.Run("drag_drop_upload", func(t *testing.T) {
		// Simulate drag and drop by sending file data
		fileData, err := os.ReadFile(helpers.TestFixtures.ValidAEP)
		helpers.AssertNoError(t, err, "Failed to read test file")

		req, err := http.NewRequest("POST", server.URL+"/upload/easy", bytes.NewReader(fileData))
		helpers.AssertNoError(t, err, "Failed to create request")

		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-Filename", "dragged-file.aep")

		resp, err := http.DefaultClient.Do(req)
		helpers.AssertNoError(t, err, "Upload failed")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Upload status")
	})

	// Test automatic story generation
	t.Run("auto_story_generation", func(t *testing.T) {
		// Upload and wait for processing
		projectID := uploadTestFile(t, server.URL+"/upload/easy")

		// Poll for story completion
		var storyReady bool
		for i := 0; i < 10; i++ {
			resp, err := http.Get(server.URL + "/status/" + projectID)
			if err != nil {
				continue
			}

			var status struct {
				Ready bool `json:"ready"`
			}
			json.NewDecoder(resp.Body).Decode(&status)
			resp.Body.Close()

			if status.Ready {
				storyReady = true
				break
			}

			time.Sleep(500 * time.Millisecond)
		}

		if !storyReady {
			t.Error("Story generation timed out")
		}
	})
}

// TestUltimateViewer tests the ultimate viewer with advanced features
func TestUltimateViewer(t *testing.T) {
	viewer := demo.NewUltimateViewer()
	server := httptest.NewServer(viewer)
	defer server.Close()

	// Test advanced upload modes
	t.Run("advanced_upload_modes", func(t *testing.T) {
		modes := []string{"simple", "advanced", "batch"}

		for _, mode := range modes {
			resp, err := http.Get(server.URL + "/upload/" + mode)
			helpers.AssertNoError(t, err, "Failed to access upload mode: "+mode)
			defer resp.Body.Close()

			helpers.AssertEqual(t, http.StatusOK, resp.StatusCode, "Upload mode status: "+mode)
		}
	})

	// Test real-time updates via WebSocket
	t.Run("websocket_updates", func(t *testing.T) {
		// This would test WebSocket functionality
		// Simplified for this example
		t.Skip("WebSocket testing requires additional setup")
	})

	// Test export functionality
	t.Run("export_functionality", func(t *testing.T) {
		// Test various export formats
		formats := []string{"json", "xml", "pdf", "html"}

		for _, format := range formats {
			resp, err := http.Get(server.URL + "/export/test-project?format=" + format)
			if err != nil {
				t.Errorf("Export failed for format %s: %v", format, err)
				continue
			}
			resp.Body.Close()

			// Check content type
			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				t.Errorf("No content type for format %s", format)
			}
		}
	})
}

// TestViewerPerformance tests viewer performance metrics
func TestViewerPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	viewer := demo.NewSimpleStoryViewer()
	server := httptest.NewServer(viewer)
	defer server.Close()

	// Test response times
	t.Run("response_times", func(t *testing.T) {
		endpoints := []string{
			"/",
			"/upload",
			"/story/test",
			"/api/health",
		}

		for _, endpoint := range endpoints {
			start := time.Now()
			resp, err := http.Get(server.URL + endpoint)
			duration := time.Since(start)

			if err != nil {
				t.Errorf("Failed to test endpoint %s: %v", endpoint, err)
				continue
			}
			resp.Body.Close()

			// Check response time
			if duration > 200*time.Millisecond {
				t.Errorf("Endpoint %s took too long: %v", endpoint, duration)
			}
		}
	})

	// Test concurrent users
	t.Run("concurrent_users", func(t *testing.T) {
		concurrentUsers := 50
		requestsPerUser := 10

		helpers.RunConcurrentTest(t, concurrentUsers, func(userID int) {
			for i := 0; i < requestsPerUser; i++ {
				resp, err := http.Get(server.URL + "/")
				if err != nil {
					t.Errorf("User %d request %d failed: %v", userID, i, err)
					continue
				}
				resp.Body.Close()
			}
		})
	})
}

// TestViewerErrorHandling tests error handling in viewers
func TestViewerErrorHandling(t *testing.T) {
	viewer := demo.NewSimpleStoryViewer()
	server := httptest.NewServer(viewer)
	defer server.Close()

	// Test invalid file upload
	t.Run("invalid_file_upload", func(t *testing.T) {
		// Upload non-AEP file
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		part, err := writer.CreateFormFile("file", "test.txt")
		helpers.AssertNoError(t, err, "Failed to create form")

		part.Write([]byte("This is not an AEP file"))
		writer.Close()

		resp, err := http.Post(
			server.URL+"/upload",
			writer.FormDataContentType(),
			&buf,
		)
		helpers.AssertNoError(t, err, "Upload request failed")
		defer resp.Body.Close()

		// Should return error
		if resp.StatusCode == http.StatusOK {
			t.Error("Expected error for invalid file")
		}
	})

	// Test missing file handling
	t.Run("missing_file", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/story/non-existent-project")
		helpers.AssertNoError(t, err, "Request failed")
		defer resp.Body.Close()

		helpers.AssertEqual(t, http.StatusNotFound, resp.StatusCode, "Missing file status")
	})

	// Test large file handling
	t.Run("large_file_handling", func(t *testing.T) {
		// Create large dummy data (10MB)
		largeData := make([]byte, 10*1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		part, err := writer.CreateFormFile("file", "large.aep")
		helpers.AssertNoError(t, err, "Failed to create form")

		part.Write(largeData)
		writer.Close()

		resp, err := http.Post(
			server.URL+"/upload",
			writer.FormDataContentType(),
			&buf,
		)
		
		if err == nil {
			defer resp.Body.Close()
			// Should handle large files gracefully
			if resp.StatusCode == http.StatusInternalServerError {
				t.Error("Server crashed on large file")
			}
		}
	})
}

// Helper function to upload a test file
func uploadTestFile(t *testing.T, uploadURL string) string {
	t.Helper()

	fileData, err := os.ReadFile(helpers.TestFixtures.ValidAEP)
	helpers.AssertNoError(t, err, "Failed to read test file")

	resp, err := http.Post(uploadURL, "application/octet-stream", bytes.NewReader(fileData))
	helpers.AssertNoError(t, err, "Upload failed")
	defer resp.Body.Close()

	var result struct {
		ProjectID string `json:"project_id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.ProjectID
}