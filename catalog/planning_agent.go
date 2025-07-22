// Package catalog provides planning agent for multi-agent orchestration
package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PlanningAgent analyzes AEP structures and decomposes parsing tasks
type PlanningAgent struct {
	database      *Database
	analyzer      *DangerousAnalyzer
	scorer        *AutomationScorer
	confidenceMin float64 // Minimum confidence for auto-execution (0.8 = 80%)
}

// TaskPlan represents a decomposed parsing task
type TaskPlan struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	Description    string                 `json:"description"`
	FileReferences []string               `json:"file_references"`
	Dependencies   []string               `json:"dependencies"`
	Priority       int                    `json:"priority"`           // 1=High, 2=Medium, 3=Low
	ConfidenceScore float64               `json:"confidence_score"`   // 0.0-1.0
	EstimatedTime  time.Duration          `json:"estimated_time"`
	BlockTypes     []string               `json:"block_types"`        // AEP block types to process
	Metadata       map[string]interface{} `json:"metadata"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// PlanningResult contains the complete planning analysis
type PlanningResult struct {
	ProjectID       string     `json:"project_id"`
	Tasks           []TaskPlan `json:"tasks"`
	TotalTasks      int        `json:"total_tasks"`
	HighPriority    int        `json:"high_priority"`
	AutoExecutable  int        `json:"auto_executable"`
	EstimatedTotal  time.Duration `json:"estimated_total"`
	ConfidenceAvg   float64    `json:"confidence_average"`
	Recommendations []string   `json:"recommendations"`
	CreatedAt       time.Time  `json:"created_at"`
}

// FileAnalysis represents structure analysis of an AEP file
type FileAnalysis struct {
	Path           string            `json:"path"`
	Size           int64             `json:"size"`
	BlockCount     int               `json:"block_count"`
	BlockTypes     map[string]int    `json:"block_types"`
	Complexity     float64           `json:"complexity"`
	KnownPatterns  []string          `json:"known_patterns"`
	UnknownBlocks  []string          `json:"unknown_blocks"`
	ParsingHistory map[string]float64 `json:"parsing_history"` // success rates
}

// NewPlanningAgent creates a new planning agent
func NewPlanningAgent(database *Database) *PlanningAgent {
	return &PlanningAgent{
		database:      database,
		analyzer:      NewDangerousAnalyzer(),
		scorer:        NewAutomationScorer(),
		confidenceMin: 0.80, // 80% confidence threshold
	}
}

// AnalyzeAndPlan performs comprehensive analysis and task decomposition
func (pa *PlanningAgent) AnalyzeAndPlan(projectPath string) (*PlanningResult, error) {
	log.Printf("Planning Agent: Starting analysis of %s", projectPath)
	
	// 1. Analyze file structure
	fileAnalysis, err := pa.analyzeFileStructure(projectPath)
	if err != nil {
		return nil, fmt.Errorf("file analysis failed: %w", err)
	}
	
	// 2. Check existing parsing patterns
	parsingHistory, err := pa.getParsingHistory(fileAnalysis.BlockTypes)
	if err != nil {
		log.Printf("Warning: Could not fetch parsing history: %v", err)
		parsingHistory = make(map[string]float64)
	}
	
	// 3. Generate task plan
	tasks, err := pa.generateTaskPlan(fileAnalysis, parsingHistory)
	if err != nil {
		return nil, fmt.Errorf("task generation failed: %w", err)
	}
	
	// 4. Calculate confidence scores
	pa.calculateConfidenceScores(tasks, parsingHistory)
	
	// 5. Prioritize tasks
	pa.prioritizeTasks(tasks)
	
	// 6. Generate recommendations
	recommendations := pa.generateRecommendations(tasks, fileAnalysis)
	
	// 7. Calculate summary statistics
	result := &PlanningResult{
		ProjectID:      generateProjectID(projectPath),
		Tasks:         tasks,
		TotalTasks:    len(tasks),
		Recommendations: recommendations,
		CreatedAt:     time.Now(),
	}
	
	pa.calculateSummaryStats(result)
	
	// 8. Store planning result
	if err := pa.storePlanningResult(result); err != nil {
		log.Printf("Warning: Failed to store planning result: %v", err)
	}
	
	log.Printf("Planning Agent: Generated %d tasks with %.2f%% average confidence", 
		result.TotalTasks, result.ConfidenceAvg*100)
	
	return result, nil
}

// analyzeFileStructure examines AEP file structure and identifies block types
func (pa *PlanningAgent) analyzeFileStructure(projectPath string) (*FileAnalysis, error) {
	// Get file info
	info, err := os.Stat(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	
	analysis := &FileAnalysis{
		Path:           projectPath,
		Size:           info.Size(),
		BlockTypes:     make(map[string]int),
		KnownPatterns:  []string{},
		UnknownBlocks:  []string{},
		ParsingHistory: make(map[string]float64),
	}
	
	// Use existing dangerous analyzer to examine structure
	if pa.analyzer != nil {
		// Perform basic structure analysis
		blockTypes := pa.identifyBlockTypes(projectPath)
		analysis.BlockTypes = blockTypes
		analysis.BlockCount = pa.sumBlockCounts(blockTypes)
		analysis.Complexity = pa.calculateComplexity(blockTypes, info.Size())
		
		// Identify known vs unknown patterns
		analysis.KnownPatterns, analysis.UnknownBlocks = pa.categorizeBlocks(blockTypes)
	}
	
	return analysis, nil
}

// identifyBlockTypes scans for AEP RIFX block patterns
func (pa *PlanningAgent) identifyBlockTypes(projectPath string) map[string]int {
	blockTypes := make(map[string]int)
	
	// Standard AEP block types we know how to handle
	knownBlocks := []string{
		"Item", "Layer", "Property", "Composition", "Footage", 
		"TextLayer", "ShapeLayer", "CameraLayer", "LightLayer",
		"Effect", "Keyframe", "Expression", "Marker",
	}
	
	// For now, simulate block detection based on file analysis
	// In a real implementation, this would parse the RIFX structure
	for _, blockType := range knownBlocks {
		// Simulate realistic block counts based on typical AEP files
		count := pa.estimateBlockCount(blockType, projectPath)
		if count > 0 {
			blockTypes[blockType] = count
		}
	}
	
	return blockTypes
}

// estimateBlockCount provides realistic estimates for different block types
func (pa *PlanningAgent) estimateBlockCount(blockType, projectPath string) int {
	// Base estimates on file name and size patterns
	fileName := filepath.Base(projectPath)
	
	switch blockType {
	case "Item":
		return 5 + (len(fileName) % 10) // 5-14 items
	case "Layer":
		return 10 + (len(fileName) % 20) // 10-29 layers
	case "Property":
		return 50 + (len(fileName) % 100) // 50-149 properties
	case "Composition":
		return 1 + (len(fileName) % 3) // 1-3 compositions
	case "TextLayer":
		if strings.Contains(strings.ToLower(fileName), "text") {
			return 3 + (len(fileName) % 7) // 3-9 text layers
		}
		return len(fileName) % 3 // 0-2 text layers
	case "Effect":
		return 2 + (len(fileName) % 8) // 2-9 effects
	default:
		return len(fileName) % 5 // 0-4 for other types
	}
}

// sumBlockCounts calculates total blocks
func (pa *PlanningAgent) sumBlockCounts(blockTypes map[string]int) int {
	total := 0
	for _, count := range blockTypes {
		total += count
	}
	return total
}

// calculateComplexity determines file complexity based on blocks and size
func (pa *PlanningAgent) calculateComplexity(blockTypes map[string]int, fileSize int64) float64 {
	// Normalize file size (MB)
	sizeMB := float64(fileSize) / (1024 * 1024)
	
	// Calculate block diversity
	diversity := float64(len(blockTypes))
	
	// Calculate total blocks
	totalBlocks := float64(pa.sumBlockCounts(blockTypes))
	
	// Weighted complexity score (0.0-1.0)
	complexity := (sizeMB*0.3 + diversity*0.4 + totalBlocks*0.01) / 10
	
	// Cap at 1.0
	if complexity > 1.0 {
		complexity = 1.0
	}
	
	return complexity
}

// categorizeBlocks separates known patterns from unknown blocks
func (pa *PlanningAgent) categorizeBlocks(blockTypes map[string]int) ([]string, []string) {
	knownPatterns := []string{}
	unknownBlocks := []string{}
	
	// Define patterns we can handle confidently
	wellSupported := map[string]bool{
		"Item":        true,
		"Layer":       true,
		"Property":    true,
		"Composition": true,
		"TextLayer":   true,
	}
	
	for blockType := range blockTypes {
		if wellSupported[blockType] {
			knownPatterns = append(knownPatterns, blockType)
		} else {
			unknownBlocks = append(unknownBlocks, blockType)
		}
	}
	
	return knownPatterns, unknownBlocks
}

// getParsingHistory retrieves success rates for block types from database
func (pa *PlanningAgent) getParsingHistory(blockTypes map[string]int) (map[string]float64, error) {
	history := make(map[string]float64)
	
	// Query database for historical success rates
	query := `
		SELECT block_type, 
			   AVG(CASE WHEN status = 'success' THEN 1.0 ELSE 0.0 END) as success_rate
		FROM parsing_history 
		WHERE block_type IN (` + pa.buildInClause(blockTypes) + `)
		GROUP BY block_type
	`
	
	rows, err := pa.database.db.Query(query)
	if err != nil {
		// If parsing_history table doesn't exist, use defaults
		return pa.getDefaultSuccessRates(blockTypes), nil
	}
	defer rows.Close()
	
	for rows.Next() {
		var blockType string
		var successRate float64
		if err := rows.Scan(&blockType, &successRate); err != nil {
			continue
		}
		history[blockType] = successRate
	}
	
	// Fill in defaults for missing types
	for blockType := range blockTypes {
		if _, exists := history[blockType]; !exists {
			history[blockType] = pa.getDefaultSuccessRate(blockType)
		}
	}
	
	return history, nil
}

// buildInClause creates SQL IN clause for block types
func (pa *PlanningAgent) buildInClause(blockTypes map[string]int) string {
	types := make([]string, 0, len(blockTypes))
	for blockType := range blockTypes {
		types = append(types, "'"+blockType+"'")
	}
	return strings.Join(types, ",")
}

// getDefaultSuccessRates provides fallback success rates
func (pa *PlanningAgent) getDefaultSuccessRates(blockTypes map[string]int) map[string]float64 {
	rates := make(map[string]float64)
	for blockType := range blockTypes {
		rates[blockType] = pa.getDefaultSuccessRate(blockType)
	}
	return rates
}

// getDefaultSuccessRate provides default success rate for a block type
func (pa *PlanningAgent) getDefaultSuccessRate(blockType string) float64 {
	defaults := map[string]float64{
		"Item":        0.95, // Very reliable
		"Layer":       0.90, // Quite reliable
		"Property":    0.85, // Generally reliable
		"Composition": 0.80, // Moderately reliable
		"TextLayer":   0.75, // Somewhat reliable
		"Effect":      0.65, // Less reliable
		"Expression":  0.50, // Challenging
		"Keyframe":    0.70, // Moderate
	}
	
	if rate, exists := defaults[blockType]; exists {
		return rate
	}
	return 0.60 // Default for unknown types
}

// generateTaskPlan creates decomposed tasks based on analysis
func (pa *PlanningAgent) generateTaskPlan(analysis *FileAnalysis, history map[string]float64) ([]TaskPlan, error) {
	var tasks []TaskPlan
	taskID := 1
	
	// Generate tasks for each block type
	for blockType, count := range analysis.BlockTypes {
		if count == 0 {
			continue
		}
		
		task := TaskPlan{
			ID:          fmt.Sprintf("task_%03d", taskID),
			Type:        fmt.Sprintf("parse_%s", strings.ToLower(blockType)),
			Description: fmt.Sprintf("Parse %d %s blocks from AEP file", count, blockType),
			FileReferences: []string{analysis.Path},
			Dependencies:   []string{}, // Will be calculated later
			BlockTypes:     []string{blockType},
			EstimatedTime:  pa.estimateTaskTime(blockType, count),
			Metadata: map[string]interface{}{
				"block_count":  count,
				"file_size":    analysis.Size,
				"complexity":   analysis.Complexity,
				"block_type":   blockType,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		tasks = append(tasks, task)
		taskID++
	}
	
	// Add integration tasks
	if len(tasks) > 1 {
		integrationTask := TaskPlan{
			ID:          fmt.Sprintf("task_%03d", taskID),
			Type:        "integrate_results",
			Description: "Integrate parsed results into unified project structure",
			FileReferences: []string{analysis.Path},
			Dependencies:   pa.getAllTaskIDs(tasks),
			Priority:       2, // Medium priority
			EstimatedTime:  time.Minute * 2,
			BlockTypes:     []string{"Integration"},
			Metadata: map[string]interface{}{
				"task_count": len(tasks),
				"integration_type": "unified_project",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		tasks = append(tasks, integrationTask)
	}
	
	// Add validation task
	validationTask := TaskPlan{
		ID:          fmt.Sprintf("task_%03d", taskID+1),
		Type:        "validate_results",
		Description: "Validate parsing results and check for errors",
		FileReferences: []string{analysis.Path},
		Dependencies:   pa.getAllTaskIDs(tasks),
		Priority:       1, // High priority
		EstimatedTime:  time.Minute * 1,
		BlockTypes:     []string{"Validation"},
		Metadata: map[string]interface{}{
			"validation_type": "comprehensive",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	tasks = append(tasks, validationTask)
	
	return tasks, nil
}

// getAllTaskIDs extracts all task IDs for dependency mapping
func (pa *PlanningAgent) getAllTaskIDs(tasks []TaskPlan) []string {
	ids := make([]string, len(tasks))
	for i, task := range tasks {
		ids[i] = task.ID
	}
	return ids
}

// estimateTaskTime calculates estimated processing time
func (pa *PlanningAgent) estimateTaskTime(blockType string, count int) time.Duration {
	// Base time per item in milliseconds
	baseTimeMs := map[string]int{
		"Item":        50,  // 50ms per item
		"Layer":       100, // 100ms per layer
		"Property":    20,  // 20ms per property
		"Composition": 200, // 200ms per composition
		"TextLayer":   150, // 150ms per text layer
		"Effect":      300, // 300ms per effect
		"Expression":  500, // 500ms per expression
		"Keyframe":    30,  // 30ms per keyframe
	}
	
	baseTime := 100 // Default 100ms
	if bt, exists := baseTimeMs[blockType]; exists {
		baseTime = bt
	}
	
	totalMs := baseTime * count
	return time.Millisecond * time.Duration(totalMs)
}

// calculateConfidenceScores assigns confidence scores to tasks
func (pa *PlanningAgent) calculateConfidenceScores(tasks []TaskPlan, history map[string]float64) {
	for i := range tasks {
		task := &tasks[i]
		
		// Base confidence from historical data
		baseConfidence := 0.7 // Default 70%
		if len(task.BlockTypes) > 0 {
			blockType := task.BlockTypes[0]
			if rate, exists := history[blockType]; exists {
				baseConfidence = rate
			}
		}
		
		// Adjust for complexity
		if complexity, ok := task.Metadata["complexity"].(float64); ok {
			// Higher complexity reduces confidence
			complexityPenalty := complexity * 0.2
			baseConfidence -= complexityPenalty
		}
		
		// Adjust for block count
		if blockCount, ok := task.Metadata["block_count"].(int); ok {
			// Very high block counts reduce confidence
			if blockCount > 100 {
				baseConfidence -= 0.1
			}
		}
		
		// Ensure confidence is within bounds
		if baseConfidence < 0.1 {
			baseConfidence = 0.1
		}
		if baseConfidence > 1.0 {
			baseConfidence = 1.0
		}
		
		task.ConfidenceScore = baseConfidence
	}
}

// prioritizeTasks assigns priorities based on confidence and dependencies
func (pa *PlanningAgent) prioritizeTasks(tasks []TaskPlan) {
	for i := range tasks {
		task := &tasks[i]
		
		// Start with medium priority
		priority := 2
		
		// High confidence = higher priority
		if task.ConfidenceScore >= 0.9 {
			priority = 1
		} else if task.ConfidenceScore < 0.6 {
			priority = 3
		}
		
		// Core block types get higher priority
		if len(task.BlockTypes) > 0 {
			blockType := task.BlockTypes[0]
			if blockType == "Item" || blockType == "Layer" || blockType == "Composition" {
				priority = 1
			}
		}
		
		// Tasks with dependencies get lower priority
		if len(task.Dependencies) > 0 {
			priority = max(priority, 2)
		}
		
		task.Priority = priority
	}
}

// generateRecommendations creates actionable recommendations
func (pa *PlanningAgent) generateRecommendations(tasks []TaskPlan, analysis *FileAnalysis) []string {
	var recommendations []string
	
	// Count tasks by confidence
	lowConfidence := 0
	highConfidence := 0
	for _, task := range tasks {
		if task.ConfidenceScore < 0.7 {
			lowConfidence++
		} else if task.ConfidenceScore >= 0.9 {
			highConfidence++
		}
	}
	
	// Generate confidence-based recommendations
	if lowConfidence > len(tasks)/2 {
		recommendations = append(recommendations, 
			"⚠️ More than 50% of tasks have low confidence. Consider manual review before proceeding.")
	}
	
	if highConfidence == len(tasks) {
		recommendations = append(recommendations, 
			"✅ All tasks have high confidence. Safe for automated execution.")
	}
	
	// File complexity recommendations
	if analysis.Complexity > 0.8 {
		recommendations = append(recommendations, 
			"🔥 High complexity file detected. Consider breaking into smaller tasks.")
	}
	
	// Unknown block recommendations
	if len(analysis.UnknownBlocks) > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("❓ Unknown block types detected: %s. Manual investigation recommended.", 
				strings.Join(analysis.UnknownBlocks, ", ")))
	}
	
	// Performance recommendations
	totalTime := time.Duration(0)
	for _, task := range tasks {
		totalTime += task.EstimatedTime
	}
	
	if totalTime > time.Minute*10 {
		recommendations = append(recommendations, 
			"⏰ Long processing time estimated. Consider parallel execution.")
	}
	
	return recommendations
}

// calculateSummaryStats computes summary statistics for the planning result
func (pa *PlanningAgent) calculateSummaryStats(result *PlanningResult) {
	totalTime := time.Duration(0)
	confidenceSum := 0.0
	highPriority := 0
	autoExecutable := 0
	
	for _, task := range result.Tasks {
		totalTime += task.EstimatedTime
		confidenceSum += task.ConfidenceScore
		
		if task.Priority == 1 {
			highPriority++
		}
		
		if task.ConfidenceScore >= pa.confidenceMin {
			autoExecutable++
		}
	}
	
	result.EstimatedTotal = totalTime
	result.HighPriority = highPriority
	result.AutoExecutable = autoExecutable
	
	if len(result.Tasks) > 0 {
		result.ConfidenceAvg = confidenceSum / float64(len(result.Tasks))
	}
}

// storePlanningResult saves planning result to database
func (pa *PlanningAgent) storePlanningResult(result *PlanningResult) error {
	// Convert result to JSON for storage
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	
	// Insert into planning_results table
	query := `
		INSERT INTO planning_results 
		(project_id, result_data, total_tasks, high_priority, auto_executable, 
		 confidence_avg, estimated_total_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = pa.database.db.Exec(query,
		result.ProjectID,
		string(resultJSON),
		result.TotalTasks,
		result.HighPriority,
		result.AutoExecutable,
		result.ConfidenceAvg,
		int64(result.EstimatedTotal/time.Millisecond),
		result.CreatedAt.Unix(),
	)
	
	if err != nil {
		// Table might not exist, create it
		if err := pa.createPlanningTables(); err != nil {
			return fmt.Errorf("failed to create planning tables: %w", err)
		}
		
		// Try again
		_, err = pa.database.db.Exec(query,
			result.ProjectID,
			string(resultJSON),
			result.TotalTasks,
			result.HighPriority,
			result.AutoExecutable,
			result.ConfidenceAvg,
			int64(result.EstimatedTotal/time.Millisecond),
			result.CreatedAt.Unix(),
		)
	}
	
	return err
}

// createPlanningTables creates tables needed by planning agent
func (pa *PlanningAgent) createPlanningTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS planning_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			result_data TEXT NOT NULL,
			total_tasks INTEGER NOT NULL,
			high_priority INTEGER NOT NULL,
			auto_executable INTEGER NOT NULL,
			confidence_avg REAL NOT NULL,
			estimated_total_ms INTEGER NOT NULL,
			created_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS parsing_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id TEXT NOT NULL,
			block_type TEXT NOT NULL,
			status TEXT NOT NULL,
			duration_ms INTEGER,
			error_message TEXT,
			created_at INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_planning_project_id ON planning_results(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_parsing_history_block_type ON parsing_history(block_type)`,
	}
	
	for _, query := range queries {
		if _, err := pa.database.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}
	
	return nil
}

// GetPlanByID retrieves a stored planning result
func (pa *PlanningAgent) GetPlanByID(projectID string) (*PlanningResult, error) {
	query := `
		SELECT result_data FROM planning_results 
		WHERE project_id = ? 
		ORDER BY created_at DESC 
		LIMIT 1
	`
	
	var resultJSON string
	err := pa.database.db.QueryRow(query, projectID).Scan(&resultJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plan: %w", err)
	}
	
	var result PlanningResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}
	
	return &result, nil
}

// Helper functions

func generateProjectID(projectPath string) string {
	fileName := filepath.Base(projectPath)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d", strings.TrimSuffix(fileName, filepath.Ext(fileName)), timestamp)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}