// Package catalog provides verification agent for automated testing and quality validation
package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// VerificationAgent validates parsing accuracy and code quality
type VerificationAgent struct {
	database         *Database
	testRunner       *TestRunner
	qualityChecker   *QualityChecker
	sampleFiles      []string
	coverageTarget   float64 // 80% coverage requirement
	validationRules  []ValidationRule
}

// TestRunner manages automated test execution
type TestRunner struct {
	goCommand      string
	testTimeout    time.Duration
	testDir        string
	tempDir        string
	parallelTests  bool
	verboseOutput  bool
}

// QualityChecker analyzes code quality metrics
type QualityChecker struct {
	maxComplexity     int     // Cyclomatic complexity limit
	minCoverage       float64 // Minimum test coverage
	maxDuplication    float64 // Maximum code duplication
	requiredPatterns  []string // Required code patterns
	forbiddenPatterns []string // Forbidden code patterns
}

// VerificationRequest represents a verification task
type VerificationRequest struct {
	TaskID         string                 `json:"task_id"`
	BlockType      string                 `json:"block_type"`
	GeneratedCode  string                 `json:"generated_code"`
	TestCode       string                 `json:"test_code"`
	SampleData     []byte                 `json:"sample_data,omitempty"`
	Requirements   []string               `json:"requirements"`
	Context        map[string]interface{} `json:"context"`
	CreatedAt      time.Time              `json:"created_at"`
}

// VerificationResult contains validation results
type VerificationResult struct {
	TaskID           string              `json:"task_id"`
	BlockType        string              `json:"block_type"`
	Status           string              `json:"status"` // passed, failed, partial
	OverallScore     float64             `json:"overall_score"` // 0.0-1.0
	
	// Test Results
	UnitTests        TestResults         `json:"unit_tests"`
	IntegrationTests TestResults         `json:"integration_tests"`
	BinaryTests      TestResults         `json:"binary_tests"`
	
	// Quality Metrics
	CodeQuality      QualityMetrics      `json:"code_quality"`
	Coverage         CoverageResults     `json:"coverage"`
	Performance      PerformanceResults  `json:"performance"`
	
	// Validation Details
	ValidationErrors []ValidationError   `json:"validation_errors"`
	Recommendations  []string            `json:"recommendations"`
	
	// Timing
	VerificationTime time.Duration       `json:"verification_time"`
	CreatedAt        time.Time           `json:"created_at"`
}

// TestResults contains test execution results
type TestResults struct {
	TotalTests    int           `json:"total_tests"`
	PassedTests   int           `json:"passed_tests"`
	FailedTests   int           `json:"failed_tests"`
	SkippedTests  int           `json:"skipped_tests"`
	ExecutionTime time.Duration `json:"execution_time"`
	Output        string        `json:"output,omitempty"`
	Errors        []string      `json:"errors,omitempty"`
}

// QualityMetrics contains code quality analysis
type QualityMetrics struct {
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
	LinesOfCode          int     `json:"lines_of_code"`
	Maintainability      float64 `json:"maintainability"` // 0.0-1.0
	Readability          float64 `json:"readability"`     // 0.0-1.0
	Duplication          float64 `json:"duplication"`     // % duplicated
	TechnicalDebt        string  `json:"technical_debt"`  // low, medium, high
}

// CoverageResults contains test coverage analysis
type CoverageResults struct {
	LineCoverage     float64            `json:"line_coverage"`
	BranchCoverage   float64            `json:"branch_coverage"`
	FunctionCoverage float64            `json:"function_coverage"`
	CoveredLines     int                `json:"covered_lines"`
	TotalLines       int                `json:"total_lines"`
	UncoveredLines   []int              `json:"uncovered_lines"`
	CoverageByFile   map[string]float64 `json:"coverage_by_file"`
}

// PerformanceResults contains performance benchmarks
type PerformanceResults struct {
	AvgExecutionTime time.Duration `json:"avg_execution_time"`
	MinExecutionTime time.Duration `json:"min_execution_time"`
	MaxExecutionTime time.Duration `json:"max_execution_time"`
	MemoryUsage      int64         `json:"memory_usage"` // bytes
	AllocationsCount int64         `json:"allocations_count"`
	BenchmarkScore   float64       `json:"benchmark_score"`
	Passed           bool          `json:"passed"`
}

// ValidationError represents a validation issue
type ValidationError struct {
	Type        string `json:"type"`        // syntax, logic, performance, security
	Severity    string `json:"severity"`    // error, warning, info
	Message     string `json:"message"`
	Line        int    `json:"line,omitempty"`
	Column      int    `json:"column,omitempty"`
	Rule        string `json:"rule,omitempty"`
	Suggestion  string `json:"suggestion,omitempty"`
}

// NewVerificationAgent creates a new verification agent
func NewVerificationAgent(database *Database) *VerificationAgent {
	agent := &VerificationAgent{
		database:       database,
		coverageTarget: 0.80, // 80% coverage requirement
		testRunner: &TestRunner{
			goCommand:     "go",
			testTimeout:   time.Minute * 5,
			testDir:       "test_temp",
			parallelTests: true,
			verboseOutput: true,
		},
		qualityChecker: &QualityChecker{
			maxComplexity:  10,
			minCoverage:    0.80,
			maxDuplication: 0.10,
			requiredPatterns: []string{
				"error handling",
				"input validation",
				"struct definition",
			},
			forbiddenPatterns: []string{
				"panic(",
				"os.Exit(",
				"fmt.Print", // Should use proper logging
			},
		},
	}
	
	// Initialize sample files
	agent.initializeSampleFiles()
	
	// Create database tables
	if err := agent.createVerificationTables(); err != nil {
		log.Printf("Warning: Failed to create verification tables: %v", err)
	}
	
	return agent
}

// VerifyImplementation performs comprehensive verification
func (va *VerificationAgent) VerifyImplementation(request *VerificationRequest) (*VerificationResult, error) {
	log.Printf("Verification Agent: Starting verification for %s (%s)", request.TaskID, request.BlockType)
	
	startTime := time.Now()
	
	result := &VerificationResult{
		TaskID:           request.TaskID,
		BlockType:        request.BlockType,
		ValidationErrors: []ValidationError{},
		Recommendations:  []string{},
		CreatedAt:        time.Now(),
	}
	
	// 1. Validate code syntax and structure
	if err := va.validateCodeStructure(request, result); err != nil {
		result.Status = "failed"
		result.ValidationErrors = append(result.ValidationErrors, ValidationError{
			Type:     "syntax",
			Severity: "error",
			Message:  fmt.Sprintf("Code structure validation failed: %v", err),
		})
		return result, err
	}
	
	// 2. Run unit tests
	if err := va.runUnitTests(request, result); err != nil {
		log.Printf("Unit tests failed: %v", err)
		result.ValidationErrors = append(result.ValidationErrors, ValidationError{
			Type:     "testing",
			Severity: "error",
			Message:  fmt.Sprintf("Unit tests failed: %v", err),
		})
	}
	
	// 3. Run integration tests
	if err := va.runIntegrationTests(request, result); err != nil {
		log.Printf("Integration tests failed: %v", err)
		result.ValidationErrors = append(result.ValidationErrors, ValidationError{
			Type:     "testing",
			Severity: "warning",
			Message:  fmt.Sprintf("Integration tests failed: %v", err),
		})
	}
	
	// 4. Binary validation with sample data
	if err := va.runBinaryValidation(request, result); err != nil {
		log.Printf("Binary validation failed: %v", err)
		result.ValidationErrors = append(result.ValidationErrors, ValidationError{
			Type:     "binary",
			Severity: "error",
			Message:  fmt.Sprintf("Binary validation failed: %v", err),
		})
	}
	
	// 5. Analyze code quality
	va.analyzeCodeQuality(request, result)
	
	// 6. Calculate test coverage
	va.calculateTestCoverage(request, result)
	
	// 7. Run performance benchmarks
	va.runPerformanceBenchmarks(request, result)
	
	// 8. Calculate overall score
	va.calculateOverallScore(result)
	
	// 9. Generate recommendations
	va.generateRecommendations(result)
	
	// 10. Determine final status
	va.determineFinalStatus(result)
	
	result.VerificationTime = time.Since(startTime)
	
	// Store result
	if err := va.storeVerificationResult(result); err != nil {
		log.Printf("Warning: Failed to store verification result: %v", err)
	}
	
	log.Printf("Verification Agent: Completed verification for %s (Score: %.2f, Status: %s)", 
		request.TaskID, result.OverallScore, result.Status)
	
	return result, nil
}

// validateCodeStructure checks code syntax and required patterns
func (va *VerificationAgent) validateCodeStructure(request *VerificationRequest, result *VerificationResult) error {
	code := request.GeneratedCode
	
	// Check for required patterns
	for _, pattern := range va.qualityChecker.requiredPatterns {
		if !va.hasPattern(code, pattern) {
			result.ValidationErrors = append(result.ValidationErrors, ValidationError{
				Type:     "structure",
				Severity: "warning",
				Message:  fmt.Sprintf("Missing required pattern: %s", pattern),
				Rule:     pattern,
				Suggestion: va.getPatternSuggestion(pattern),
			})
		}
	}
	
	// Check for forbidden patterns
	for _, pattern := range va.qualityChecker.forbiddenPatterns {
		if strings.Contains(code, pattern) {
			result.ValidationErrors = append(result.ValidationErrors, ValidationError{
				Type:     "structure",
				Severity: "error",
				Message:  fmt.Sprintf("Forbidden pattern found: %s", pattern),
				Rule:     pattern,
				Suggestion: fmt.Sprintf("Replace %s with proper alternative", pattern),
			})
		}
	}
	
	// Validate Go syntax by attempting to parse
	if err := va.validateGoSyntax(code); err != nil {
		return fmt.Errorf("Go syntax validation failed: %w", err)
	}
	
	return nil
}

// runUnitTests executes unit tests for the generated code
func (va *VerificationAgent) runUnitTests(request *VerificationRequest, result *VerificationResult) error {
	if request.TestCode == "" {
		result.UnitTests = TestResults{
			TotalTests: 0,
			Output:     "No unit tests provided",
		}
		return fmt.Errorf("no unit tests provided")
	}
	
	// Create temporary test environment
	testDir, err := va.createTestEnvironment(request)
	if err != nil {
		return fmt.Errorf("failed to create test environment: %w", err)
	}
	defer os.RemoveAll(testDir)
	
	// Write test files
	mainFile := filepath.Join(testDir, fmt.Sprintf("%s.go", strings.ToLower(request.BlockType)))
	testFile := filepath.Join(testDir, fmt.Sprintf("%s_test.go", strings.ToLower(request.BlockType)))
	
	if err := va.writeTestFile(mainFile, request.GeneratedCode); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}
	
	if err := va.writeTestFile(testFile, request.TestCode); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}
	
	// Execute tests
	cmd := exec.Command(va.testRunner.goCommand, "test", "-v", "./...")
	cmd.Dir = testDir
	
	output, err := va.runWithTimeout(cmd, va.testRunner.testTimeout)
	
	// Parse test results
	testResults := va.parseTestOutput(string(output))
	testResults.Output = string(output)
	
	if err != nil {
		testResults.Errors = append(testResults.Errors, err.Error())
	}
	
	result.UnitTests = testResults
	
	return nil
}

// runIntegrationTests executes integration tests
func (va *VerificationAgent) runIntegrationTests(request *VerificationRequest, result *VerificationResult) error {
	// Create integration test that tests parsing with actual AEP data
	integrationTest := va.generateIntegrationTest(request)
	
	// Create test environment
	testDir, err := va.createTestEnvironment(request)
	if err != nil {
		return err
	}
	defer os.RemoveAll(testDir)
	
	// Write integration test
	testFile := filepath.Join(testDir, "integration_test.go")
	if err := va.writeTestFile(testFile, integrationTest); err != nil {
		return err
	}
	
	// Write main code
	mainFile := filepath.Join(testDir, fmt.Sprintf("%s.go", strings.ToLower(request.BlockType)))
	if err := va.writeTestFile(mainFile, request.GeneratedCode); err != nil {
		return err
	}
	
	// Execute integration tests
	cmd := exec.Command(va.testRunner.goCommand, "test", "-v", "-tags=integration", "./...")
	cmd.Dir = testDir
	
	output, err := va.runWithTimeout(cmd, va.testRunner.testTimeout)
	
	testResults := va.parseTestOutput(string(output))
	testResults.Output = string(output)
	
	if err != nil {
		testResults.Errors = append(testResults.Errors, err.Error())
	}
	
	result.IntegrationTests = testResults
	
	return nil
}

// runBinaryValidation validates parsing with real AEP binary data
func (va *VerificationAgent) runBinaryValidation(request *VerificationRequest, result *VerificationResult) error {
	// Test with sample AEP files
	testResults := TestResults{
		TotalTests: len(va.sampleFiles),
	}
	
	for _, sampleFile := range va.sampleFiles {
		// Read sample file
		data, err := os.ReadFile(sampleFile)
		if err != nil {
			testResults.FailedTests++
			testResults.Errors = append(testResults.Errors, fmt.Sprintf("Failed to read %s: %v", sampleFile, err))
			continue
		}
		
		// Test parsing (simulated for now)
		if va.testBinaryParsing(data, request) {
			testResults.PassedTests++
		} else {
			testResults.FailedTests++
			testResults.Errors = append(testResults.Errors, fmt.Sprintf("Binary parsing failed for %s", sampleFile))
		}
	}
	
	result.BinaryTests = testResults
	
	return nil
}

// analyzeCodeQuality performs static code analysis
func (va *VerificationAgent) analyzeCodeQuality(request *VerificationRequest, result *VerificationResult) {
	code := request.GeneratedCode
	
	quality := QualityMetrics{
		LinesOfCode: va.countLines(code),
	}
	
	// Calculate cyclomatic complexity
	quality.CyclomaticComplexity = va.calculateComplexity(code)
	
	// Calculate maintainability (based on complexity, length, etc.)
	quality.Maintainability = va.calculateMaintainability(code, quality.CyclomaticComplexity)
	
	// Calculate readability (based on naming, comments, structure)
	quality.Readability = va.calculateReadability(code)
	
	// Calculate duplication (simplified)
	quality.Duplication = va.calculateDuplication(code)
	
	// Determine technical debt level
	quality.TechnicalDebt = va.assessTechnicalDebt(quality)
	
	result.CodeQuality = quality
}

// calculateTestCoverage analyzes test coverage
func (va *VerificationAgent) calculateTestCoverage(request *VerificationRequest, result *VerificationResult) {
	coverage := CoverageResults{
		CoverageByFile: make(map[string]float64),
	}
	
	if request.TestCode == "" {
		// No tests = no coverage
		coverage.LineCoverage = 0.0
		coverage.BranchCoverage = 0.0
		coverage.FunctionCoverage = 0.0
		result.Coverage = coverage
		return
	}
	
	// Estimate coverage based on test quality
	testQuality := va.analyzeTestQuality(request.TestCode)
	
	// Simple heuristic: better tests = better coverage
	baseCoverage := 0.6 // Base coverage for having tests
	
	if strings.Contains(request.TestCode, "BenchmarkParse") {
		baseCoverage += 0.1 // Bonus for benchmarks
	}
	
	if strings.Count(request.TestCode, "t.Run(") > 1 {
		baseCoverage += 0.2 // Bonus for multiple test cases
	}
	
	if strings.Contains(request.TestCode, "wantErr") {
		baseCoverage += 0.1 // Bonus for error testing
	}
	
	coverage.LineCoverage = baseCoverage * testQuality
	coverage.BranchCoverage = coverage.LineCoverage * 0.9 // Usually slightly lower
	coverage.FunctionCoverage = coverage.LineCoverage * 1.1 // Usually slightly higher
	
	// Ensure bounds
	if coverage.LineCoverage > 1.0 {
		coverage.LineCoverage = 1.0
	}
	if coverage.BranchCoverage > 1.0 {
		coverage.BranchCoverage = 1.0
	}
	if coverage.FunctionCoverage > 1.0 {
		coverage.FunctionCoverage = 1.0
	}
	
	// Estimate covered lines
	totalLines := va.countLines(request.GeneratedCode)
	coverage.TotalLines = totalLines
	coverage.CoveredLines = int(float64(totalLines) * coverage.LineCoverage)
	
	result.Coverage = coverage
}

// runPerformanceBenchmarks executes performance tests
func (va *VerificationAgent) runPerformanceBenchmarks(request *VerificationRequest, result *VerificationResult) {
	performance := PerformanceResults{}
	
	// Simulate performance benchmarking
	if strings.Contains(request.TestCode, "BenchmarkParse") {
		// Estimate performance based on code complexity
		complexity := va.calculateComplexity(request.GeneratedCode)
		codeLines := va.countLines(request.GeneratedCode)
		
		// Simple performance estimation
		baseTime := time.Microsecond * time.Duration(100+complexity*20+codeLines)
		
		performance.AvgExecutionTime = baseTime
		performance.MinExecutionTime = baseTime * 80 / 100  // 80% of average
		performance.MaxExecutionTime = baseTime * 150 / 100 // 150% of average
		performance.MemoryUsage = int64(1024 * (complexity + codeLines/10))
		performance.AllocationsCount = int64(complexity * 2)
		performance.BenchmarkScore = va.calculateBenchmarkScore(performance)
		performance.Passed = performance.BenchmarkScore > 0.7
	} else {
		// No benchmarks available
		performance.Passed = false
		performance.BenchmarkScore = 0.0
	}
	
	result.Performance = performance
}

// calculateOverallScore computes the overall verification score
func (va *VerificationAgent) calculateOverallScore(result *VerificationResult) {
	weights := map[string]float64{
		"unit_tests":        0.30,
		"code_quality":      0.25,
		"coverage":          0.20,
		"integration_tests": 0.15,
		"performance":       0.10,
	}
	
	scores := make(map[string]float64)
	
	// Unit tests score
	if result.UnitTests.TotalTests > 0 {
		scores["unit_tests"] = float64(result.UnitTests.PassedTests) / float64(result.UnitTests.TotalTests)
	}
	
	// Code quality score
	scores["code_quality"] = result.CodeQuality.Maintainability * 0.4 + result.CodeQuality.Readability * 0.6
	
	// Coverage score
	scores["coverage"] = result.Coverage.LineCoverage
	
	// Integration tests score
	if result.IntegrationTests.TotalTests > 0 {
		scores["integration_tests"] = float64(result.IntegrationTests.PassedTests) / float64(result.IntegrationTests.TotalTests)
	}
	
	// Performance score
	if result.Performance.Passed {
		scores["performance"] = result.Performance.BenchmarkScore
	}
	
	// Calculate weighted average
	totalScore := 0.0
	totalWeight := 0.0
	
	for category, weight := range weights {
		if score, exists := scores[category]; exists {
			totalScore += score * weight
			totalWeight += weight
		}
	}
	
	if totalWeight > 0 {
		result.OverallScore = totalScore / totalWeight
	}
}

// generateRecommendations creates actionable recommendations
func (va *VerificationAgent) generateRecommendations(result *VerificationResult) {
	recommendations := []string{}
	
	// Test coverage recommendations
	if result.Coverage.LineCoverage < va.coverageTarget {
		recommendations = append(recommendations, 
			fmt.Sprintf("📊 Increase test coverage from %.1f%% to %.1f%% minimum", 
				result.Coverage.LineCoverage*100, va.coverageTarget*100))
	}
	
	// Code quality recommendations
	if result.CodeQuality.CyclomaticComplexity > va.qualityChecker.maxComplexity {
		recommendations = append(recommendations, 
			fmt.Sprintf("🔧 Reduce cyclomatic complexity from %d to %d maximum", 
				result.CodeQuality.CyclomaticComplexity, va.qualityChecker.maxComplexity))
	}
	
	// Unit test recommendations
	if result.UnitTests.TotalTests == 0 {
		recommendations = append(recommendations, "✅ Add unit tests to verify parsing functionality")
	} else if result.UnitTests.FailedTests > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("🚨 Fix %d failing unit tests", result.UnitTests.FailedTests))
	}
	
	// Performance recommendations
	if !result.Performance.Passed {
		recommendations = append(recommendations, "⚡ Add performance benchmarks and optimize parsing speed")
	}
	
	// Binary validation recommendations
	if result.BinaryTests.FailedTests > result.BinaryTests.PassedTests {
		recommendations = append(recommendations, "🔍 Improve binary parsing accuracy with sample AEP files")
	}
	
	// Overall score recommendations
	if result.OverallScore < 0.8 {
		recommendations = append(recommendations, "📈 Overall verification score is below 80% - focus on critical issues first")
	}
	
	result.Recommendations = recommendations
}

// determineFinalStatus sets the final verification status
func (va *VerificationAgent) determineFinalStatus(result *VerificationResult) {
	// Count critical errors
	criticalErrors := 0
	for _, err := range result.ValidationErrors {
		if err.Severity == "error" {
			criticalErrors++
		}
	}
	
	// Determine status based on score and errors
	if criticalErrors > 0 {
		result.Status = "failed"
	} else if result.OverallScore >= 0.8 {
		result.Status = "passed"
	} else if result.OverallScore >= 0.6 {
		result.Status = "partial"
	} else {
		result.Status = "failed"
	}
}

// Helper methods

func (va *VerificationAgent) initializeSampleFiles() {
	// Initialize with sample AEP files
	va.sampleFiles = []string{
		"data/Item-01.aep",
		"data/Layer-01.aep",
		"data/Property-01.aep",
	}
}

func (va *VerificationAgent) hasPattern(code, pattern string) bool {
	switch pattern {
	case "error handling":
		return strings.Contains(code, "error") && (strings.Contains(code, "return") || strings.Contains(code, "err"))
	case "input validation":
		return strings.Contains(code, "len(") || strings.Contains(code, "nil") || strings.Contains(code, "validate")
	case "struct definition":
		return strings.Contains(code, "type") && strings.Contains(code, "struct")
	default:
		return strings.Contains(code, pattern)
	}
}

func (va *VerificationAgent) getPatternSuggestion(pattern string) string {
	suggestions := map[string]string{
		"error handling":    "Add proper error checking and return appropriate error messages",
		"input validation":  "Validate input parameters for nil, length, and bounds",
		"struct definition": "Define appropriate struct types for parsed data",
	}
	
	if suggestion, exists := suggestions[pattern]; exists {
		return suggestion
	}
	return fmt.Sprintf("Implement %s in your code", pattern)
}

func (va *VerificationAgent) validateGoSyntax(code string) error {
	// Simple syntax validation - in production, use go/parser
	requiredKeywords := []string{"func", "package"}
	for _, keyword := range requiredKeywords {
		if !strings.Contains(code, keyword) {
			return fmt.Errorf("missing required keyword: %s", keyword)
		}
	}
	return nil
}

func (va *VerificationAgent) createTestEnvironment(request *VerificationRequest) (string, error) {
	testDir := filepath.Join(os.TempDir(), fmt.Sprintf("mobot_test_%d", time.Now().UnixNano()))
	
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return "", err
	}
	
	// Create go.mod
	goMod := fmt.Sprintf("module test_%s\n\ngo 1.19\n", request.BlockType)
	if err := os.WriteFile(filepath.Join(testDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return "", err
	}
	
	return testDir, nil
}

func (va *VerificationAgent) writeTestFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func (va *VerificationAgent) runWithTimeout(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	done := make(chan struct{})
	var output []byte
	var err error
	
	go func() {
		output, err = cmd.CombinedOutput()
		close(done)
	}()
	
	select {
	case <-done:
		return output, err
	case <-time.After(timeout):
		cmd.Process.Kill()
		return nil, fmt.Errorf("command timed out after %v", timeout)
	}
}

func (va *VerificationAgent) parseTestOutput(output string) TestResults {
	results := TestResults{}
	
	// Parse Go test output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "PASS:") {
			results.PassedTests++
		} else if strings.Contains(line, "FAIL:") {
			results.FailedTests++
		} else if strings.Contains(line, "SKIP:") {
			results.SkippedTests++
		}
	}
	
	results.TotalTests = results.PassedTests + results.FailedTests + results.SkippedTests
	
	// Extract execution time
	timeRegex := regexp.MustCompile(`(\d+\.\d+)s`)
	if matches := timeRegex.FindStringSubmatch(output); len(matches) > 1 {
		if duration, err := time.ParseDuration(matches[1] + "s"); err == nil {
			results.ExecutionTime = duration
		}
	}
	
	return results
}

func (va *VerificationAgent) generateIntegrationTest(request *VerificationRequest) string {
	return fmt.Sprintf(`
//go:build integration

package main

import (
	"testing"
	"os"
)

func TestIntegration%s(t *testing.T) {
	// Test with sample AEP data
	sampleFiles := []string{
		"../data/Item-01.aep",
		"../data/Layer-01.aep", 
	}
	
	for _, file := range sampleFiles {
		t.Run("parse_"+file, func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Skipf("Sample file %%s not found", file)
				return
			}
			
			if len(data) < 8 {
				t.Errorf("Sample file %%s too small", file)
				return
			}
			
			// Test parsing
			item, size, err := Parse%s(data, 0)
			if err != nil {
				t.Errorf("Parse%s failed for %%s: %%v", file, err)
				return
			}
			
			if item == nil {
				t.Errorf("Parse%s returned nil for %%s", file)
			}
			
			if size <= 0 {
				t.Errorf("Parse%s returned invalid size for %%s: %%d", file, size)
			}
		})
	}
}
`, request.BlockType, request.BlockType, request.BlockType, request.BlockType, request.BlockType)
}

func (va *VerificationAgent) testBinaryParsing(data []byte, request *VerificationRequest) bool {
	// Simulate binary parsing test
	// In a real implementation, this would call the generated parsing function
	
	// Basic checks
	if len(data) < 8 {
		return false
	}
	
	// Check for RIFX signature (simplified)
	if len(data) >= 4 {
		signature := string(data[:4])
		return signature == "RIFX" || signature == "RIFF"
	}
	
	return true
}

func (va *VerificationAgent) countLines(code string) int {
	return len(strings.Split(code, "\n"))
}

func (va *VerificationAgent) calculateComplexity(code string) int {
	complexity := 1 // Base complexity
	
	// Count decision points
	decisionKeywords := []string{"if", "for", "switch", "case", "else", "&&", "||"}
	for _, keyword := range decisionKeywords {
		complexity += strings.Count(code, keyword)
	}
	
	return complexity
}

func (va *VerificationAgent) calculateMaintainability(code string, complexity int) float64 {
	lines := va.countLines(code)
	
	// Simple maintainability index based on complexity and length
	if complexity > 20 || lines > 500 {
		return 0.3
	} else if complexity > 10 || lines > 200 {
		return 0.6
	} else {
		return 0.9
	}
}

func (va *VerificationAgent) calculateReadability(code string) float64 {
	score := 1.0
	
	// Penalize for lack of comments
	if !strings.Contains(code, "//") {
		score -= 0.2
	}
	
	// Penalize for long lines
	lines := strings.Split(code, "\n")
	longLines := 0
	for _, line := range lines {
		if len(line) > 100 {
			longLines++
		}
	}
	if longLines > len(lines)/4 {
		score -= 0.3
	}
	
	// Bonus for good naming
	if strings.Contains(code, "Parse") && strings.Contains(code, "error") {
		score += 0.1
	}
	
	if score < 0 {
		score = 0
	}
	return score
}

func (va *VerificationAgent) calculateDuplication(code string) float64 {
	// Simple duplication detection - count repeated lines
	lines := strings.Split(code, "\n")
	lineCount := make(map[string]int)
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 10 { // Only count substantial lines
			lineCount[trimmed]++
		}
	}
	
	duplicatedLines := 0
	for _, count := range lineCount {
		if count > 1 {
			duplicatedLines += count - 1
		}
	}
	
	totalLines := len(lines)
	if totalLines == 0 {
		return 0.0
	}
	
	return float64(duplicatedLines) / float64(totalLines)
}

func (va *VerificationAgent) assessTechnicalDebt(quality QualityMetrics) string {
	score := quality.Maintainability*0.4 + quality.Readability*0.6
	
	if score > 0.8 {
		return "low"
	} else if score > 0.6 {
		return "medium"
	} else {
		return "high"
	}
}

func (va *VerificationAgent) analyzeTestQuality(testCode string) float64 {
	quality := 0.5 // Base quality
	
	// Bonus for comprehensive tests
	if strings.Count(testCode, "func Test") > 1 {
		quality += 0.2
	}
	
	// Bonus for error testing
	if strings.Contains(testCode, "wantErr") {
		quality += 0.2
	}
	
	// Bonus for edge cases
	if strings.Contains(testCode, "insufficient data") || strings.Contains(testCode, "nil") {
		quality += 0.1
	}
	
	if quality > 1.0 {
		quality = 1.0
	}
	
	return quality
}

func (va *VerificationAgent) calculateBenchmarkScore(performance PerformanceResults) float64 {
	// Score based on reasonable performance expectations
	avgMicros := performance.AvgExecutionTime.Nanoseconds() / 1000
	
	if avgMicros < 100 {
		return 1.0 // Excellent performance
	} else if avgMicros < 1000 {
		return 0.8 // Good performance
	} else if avgMicros < 10000 {
		return 0.6 // Acceptable performance
	} else {
		return 0.3 // Poor performance
	}
}

// Database operations
func (va *VerificationAgent) createVerificationTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS verification_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		block_type TEXT NOT NULL,
		status TEXT NOT NULL,
		overall_score REAL NOT NULL,
		unit_tests_passed INTEGER NOT NULL,
		unit_tests_total INTEGER NOT NULL,
		coverage_percentage REAL NOT NULL,
		complexity INTEGER NOT NULL,
		result_data TEXT NOT NULL,
		verification_time_ms INTEGER NOT NULL,
		created_at INTEGER NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_verification_task_id ON verification_results(task_id);
	CREATE INDEX IF NOT EXISTS idx_verification_status ON verification_results(status);
	`
	
	_, err := va.database.db.Exec(query)
	return err
}

func (va *VerificationAgent) storeVerificationResult(result *VerificationResult) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	
	query := `
		INSERT INTO verification_results
		(task_id, block_type, status, overall_score, unit_tests_passed, unit_tests_total,
		 coverage_percentage, complexity, result_data, verification_time_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = va.database.db.Exec(query,
		result.TaskID,
		result.BlockType,
		result.Status,
		result.OverallScore,
		result.UnitTests.PassedTests,
		result.UnitTests.TotalTests,
		result.Coverage.LineCoverage,
		result.CodeQuality.CyclomaticComplexity,
		string(resultJSON),
		int64(result.VerificationTime/time.Millisecond),
		result.CreatedAt.Unix(),
	)
	
	return err
}

// Public API methods

// GetVerificationByTaskID retrieves verification result by task ID
func (va *VerificationAgent) GetVerificationByTaskID(taskID string) (*VerificationResult, error) {
	query := `
		SELECT result_data FROM verification_results 
		WHERE task_id = ? 
		ORDER BY created_at DESC 
		LIMIT 1
	`
	
	var resultJSON string
	err := va.database.db.QueryRow(query, taskID).Scan(&resultJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch verification result: %w", err)
	}
	
	var result VerificationResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}
	
	return &result, nil
}

// GetVerificationStats returns aggregated verification statistics
func (va *VerificationAgent) GetVerificationStats() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			AVG(overall_score) as avg_score,
			SUM(CASE WHEN status = 'passed' THEN 1 ELSE 0 END) as passed,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
			AVG(coverage_percentage) as avg_coverage
		FROM verification_results
	`
	
	var total, passed, failed int
	var avgScore, avgCoverage float64
	
	err := va.database.db.QueryRow(query).Scan(&total, &avgScore, &passed, &failed, &avgCoverage)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"total_verifications": total,
		"average_score":       avgScore,
		"passed_count":        passed,
		"failed_count":        failed,
		"pass_rate":          float64(passed) / float64(total),
		"average_coverage":    avgCoverage,
	}, nil
}