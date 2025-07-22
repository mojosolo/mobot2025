// Package catalog provides meta-orchestrator for multi-agent coordination
package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// MetaOrchestrator coordinates all agents and manages workflow execution
type MetaOrchestrator struct {
	database         *Database
	planningAgent    *PlanningAgent
	implementAgent   *ImplementationAgent
	verificationAgent *VerificationAgent
	reviewAgent      *ReviewAgent
	
	// Workflow management
	workflows        map[string]*Workflow
	workflowMutex    sync.RWMutex
	
	// Loop detection and safety
	maxIterations    int
	currentIteration int
	loopDetection    map[string]int // task_id -> iteration count
	
	// Human approval gates
	approvalGates    map[string]bool // workflow_id -> requires_approval
	pendingApprovals map[string]*ApprovalRequest
	approvalMutex    sync.RWMutex
	
	// Monitoring
	metrics          *OrchestrationMetrics
	eventLog         []OrchestrationEvent
	eventMutex       sync.RWMutex
}

// Workflow represents a complete multi-agent workflow
type Workflow struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ProjectPath     string                 `json:"project_path"`
	Status          string                 `json:"status"` // pending, running, completed, failed, requires_approval
	CurrentStage    string                 `json:"current_stage"`
	Progress        float64                `json:"progress"` // 0.0-1.0
	
	// Execution plan
	PlanningResult  *PlanningResult        `json:"planning_result,omitempty"`
	Tasks           []*WorkflowTask        `json:"tasks"`
	Dependencies    map[string][]string    `json:"dependencies"`
	
	// Execution state
	CompletedTasks  []string               `json:"completed_tasks"`
	FailedTasks     []string               `json:"failed_tasks"`
	CurrentTasks    []string               `json:"current_tasks"`
	
	// Safety and approval
	RequiresApproval bool                  `json:"requires_approval"`
	ApprovalReason   string                `json:"approval_reason,omitempty"`
	IterationCount   int                   `json:"iteration_count"`
	
	// Timing
	CreatedAt       time.Time              `json:"created_at"`
	StartedAt       *time.Time             `json:"started_at,omitempty"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	EstimatedTime   time.Duration          `json:"estimated_time"`
	ActualTime      time.Duration          `json:"actual_time"`
	
	// Metadata
	Metadata        map[string]interface{} `json:"metadata"`
	
	// Agent Results
	Requirements        []string               `json:"requirements,omitempty"`
	GeneratedCode       map[string]string      `json:"generated_code,omitempty"`
	VerificationResults *VerificationResult    `json:"verification_results,omitempty"`
	ReviewResults       *ReviewResult          `json:"review_results,omitempty"`
}

// WorkflowTask represents a task within a workflow
type WorkflowTask struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"` // planning, implementation, verification, review
	Status          string                 `json:"status"` // pending, running, completed, failed
	Agent           string                 `json:"agent"`
	Input           map[string]interface{} `json:"input"`
	Output          map[string]interface{} `json:"output,omitempty"`
	Error           string                 `json:"error,omitempty"`
	Dependencies    []string               `json:"dependencies"`
	Priority        int                    `json:"priority"`
	RetryCount      int                    `json:"retry_count"`
	MaxRetries      int                    `json:"max_retries"`
	CreatedAt       time.Time              `json:"created_at"`
	StartedAt       *time.Time             `json:"started_at,omitempty"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	EstimatedTime   time.Duration          `json:"estimated_time"`
	ActualTime      time.Duration          `json:"actual_time"`
}

// ApprovalRequest represents a request for human approval
type ApprovalRequest struct {
	ID              string                 `json:"id"`
	WorkflowID      string                 `json:"workflow_id"`
	TaskID          string                 `json:"task_id,omitempty"`
	Type            string                 `json:"type"` // workflow_start, critical_decision, high_risk_task
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Context         map[string]interface{} `json:"context"`
	Options         []string               `json:"options"`
	Recommendation  string                 `json:"recommendation"`
	RiskLevel       string                 `json:"risk_level"` // low, medium, high, critical
	CreatedAt       time.Time              `json:"created_at"`
	Deadline        *time.Time             `json:"deadline,omitempty"`
	Status          string                 `json:"status"` // pending, approved, rejected, expired
	Response        string                 `json:"response,omitempty"`
	RespondedAt     *time.Time             `json:"responded_at,omitempty"`
	RespondedBy     string                 `json:"responded_by,omitempty"`
}

// OrchestrationMetrics tracks system performance
type OrchestrationMetrics struct {
	TotalWorkflows     int           `json:"total_workflows"`
	ActiveWorkflows    int           `json:"active_workflows"`
	CompletedWorkflows int           `json:"completed_workflows"`
	FailedWorkflows    int           `json:"failed_workflows"`
	
	TotalTasks         int           `json:"total_tasks"`
	CompletedTasks     int           `json:"completed_tasks"`
	FailedTasks        int           `json:"failed_tasks"`
	
	AvgWorkflowTime    time.Duration `json:"avg_workflow_time"`
	AvgTaskTime        time.Duration `json:"avg_task_time"`
	
	LoopDetections     int           `json:"loop_detections"`
	ApprovalRequests   int           `json:"approval_requests"`
	ApprovalRate       float64       `json:"approval_rate"`
	
	LastUpdated        time.Time     `json:"last_updated"`
}

// OrchestrationEvent represents a system event
type OrchestrationEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // workflow_created, task_completed, approval_requested, etc.
	Level       string                 `json:"level"` // debug, info, warning, error, critical
	Message     string                 `json:"message"`
	WorkflowID  string                 `json:"workflow_id,omitempty"`
	TaskID      string                 `json:"task_id,omitempty"`
	Agent       string                 `json:"agent,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewMetaOrchestrator creates a new meta-orchestrator
func NewMetaOrchestrator(database *Database) *MetaOrchestrator {
	orchestrator := &MetaOrchestrator{
		database:         database,
		workflows:        make(map[string]*Workflow),
		maxIterations:    50, // Safety limit
		loopDetection:    make(map[string]int),
		approvalGates:    make(map[string]bool),
		pendingApprovals: make(map[string]*ApprovalRequest),
		metrics:          &OrchestrationMetrics{},
		eventLog:         make([]OrchestrationEvent, 0, 1000),
	}
	
	// Initialize agents
	orchestrator.planningAgent = NewPlanningAgent(database)
	orchestrator.implementAgent = NewImplementationAgent(database, orchestrator.planningAgent)
	// Note: VerificationAgent and ReviewAgent will be implemented in subsequent tasks
	
	// Create database tables
	if err := orchestrator.createOrchestratorTables(); err != nil {
		log.Printf("Warning: Failed to create orchestrator tables: %v", err)
	}
	
	return orchestrator
}

// CreateWorkflow creates and initializes a new workflow
func (mo *MetaOrchestrator) CreateWorkflow(projectPath, name, description string) (*Workflow, error) {
	mo.logEvent("info", "workflow_create_started", fmt.Sprintf("Creating workflow for %s", projectPath), "")
	
	workflow := &Workflow{
		ID:              generateWorkflowID(),
		Name:            name,
		Description:     description,
		ProjectPath:     projectPath,
		Status:          "pending",
		CurrentStage:    "initialization",
		Progress:        0.0,
		Tasks:           []*WorkflowTask{},
		Dependencies:    make(map[string][]string),
		CompletedTasks:  []string{},
		FailedTasks:     []string{},
		CurrentTasks:    []string{},
		RequiresApproval: mo.shouldRequireApproval(projectPath),
		IterationCount:  0,
		CreatedAt:       time.Now(),
		Metadata:        make(map[string]interface{}),
	}
	
	// Store workflow
	mo.workflowMutex.Lock()
	mo.workflows[workflow.ID] = workflow
	mo.workflowMutex.Unlock()
	
	// Request approval if needed
	if workflow.RequiresApproval {
		if err := mo.requestWorkflowApproval(workflow); err != nil {
			mo.logEvent("error", "approval_request_failed", err.Error(), workflow.ID)
			return nil, fmt.Errorf("failed to request approval: %w", err)
		}
		workflow.Status = "requires_approval"
		workflow.ApprovalReason = "High-risk automated AEP processing workflow"
	}
	
	// Save to database
	if err := mo.saveWorkflow(workflow); err != nil {
		mo.logEvent("error", "workflow_save_failed", err.Error(), workflow.ID)
		return nil, fmt.Errorf("failed to save workflow: %w", err)
	}
	
	mo.updateMetrics()
	mo.logEvent("info", "workflow_created", fmt.Sprintf("Workflow %s created successfully", workflow.ID), workflow.ID)
	
	return workflow, nil
}

// ExecuteWorkflow runs a complete multi-agent workflow
func (mo *MetaOrchestrator) ExecuteWorkflow(workflowID string) error {
	mo.workflowMutex.RLock()
	workflow, exists := mo.workflows[workflowID]
	mo.workflowMutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("workflow %s not found", workflowID)
	}
	
	// Check if approval is required and not yet granted
	if workflow.RequiresApproval && workflow.Status == "requires_approval" {
		return fmt.Errorf("workflow %s requires human approval", workflowID)
	}
	
	mo.logEvent("info", "workflow_execution_started", fmt.Sprintf("Starting execution of workflow %s", workflowID), workflowID)
	
	workflow.Status = "running"
	workflow.CurrentStage = "planning"
	startTime := time.Now()
	workflow.StartedAt = &startTime
	
	// Execute workflow stages
	if err := mo.executePlanningStage(workflow); err != nil {
		return mo.handleWorkflowError(workflow, "planning", err)
	}
	
	if err := mo.executeImplementationStage(workflow); err != nil {
		return mo.handleWorkflowError(workflow, "implementation", err)
	}
	
	if err := mo.executeVerificationStage(workflow); err != nil {
		return mo.handleWorkflowError(workflow, "verification", err)
	}
	
	if err := mo.executeReviewStage(workflow); err != nil {
		return mo.handleWorkflowError(workflow, "review", err)
	}
	
	// Complete workflow
	return mo.completeWorkflow(workflow)
}

// executePlanningStage runs the planning phase
func (mo *MetaOrchestrator) executePlanningStage(workflow *Workflow) error {
	mo.logEvent("info", "planning_stage_started", "Starting planning stage", workflow.ID)
	
	workflow.CurrentStage = "planning"
	workflow.Progress = 0.1
	
	// Check for infinite loops
	if err := mo.checkLoopDetection(workflow.ID, "planning"); err != nil {
		return err
	}
	
	// Execute planning
	planningResult, err := mo.planningAgent.AnalyzeAndPlan(workflow.ProjectPath)
	if err != nil {
		return fmt.Errorf("planning failed: %w", err)
	}
	
	workflow.PlanningResult = planningResult
	workflow.EstimatedTime = planningResult.EstimatedTotal
	
	// Convert planning tasks to workflow tasks
	for _, planTask := range planningResult.Tasks {
		workflowTask := &WorkflowTask{
			ID:            planTask.ID,
			Type:          "implementation",
			Status:        "pending",
			Agent:         "implementation_agent",
			Dependencies:  planTask.Dependencies,
			Priority:      planTask.Priority,
			MaxRetries:    3,
			CreatedAt:     time.Now(),
			EstimatedTime: planTask.EstimatedTime,
			Input: map[string]interface{}{
				"block_type":    planTask.BlockTypes,
				"description":   planTask.Description,
				"file_refs":     planTask.FileReferences,
				"confidence":    planTask.ConfidenceScore,
			},
		}
		workflow.Tasks = append(workflow.Tasks, workflowTask)
	}
	
	workflow.Progress = 0.2
	mo.logEvent("info", "planning_stage_completed", fmt.Sprintf("Planning completed: %d tasks generated", len(workflow.Tasks)), workflow.ID)
	
	return nil
}

// executeImplementationStage runs the implementation phase
func (mo *MetaOrchestrator) executeImplementationStage(workflow *Workflow) error {
	mo.logEvent("info", "implementation_stage_started", "Starting implementation stage", workflow.ID)
	
	workflow.CurrentStage = "implementation"
	workflow.Progress = 0.3
	
	// Execute implementation tasks
	totalTasks := len(workflow.Tasks)
	completedTasks := 0
	
	for _, task := range workflow.Tasks {
		if task.Type != "implementation" {
			continue
		}
		
		// Check dependencies
		if !mo.checkTaskDependencies(task, workflow) {
			mo.logEvent("warning", "task_dependencies_not_met", fmt.Sprintf("Task %s dependencies not met, skipping", task.ID), workflow.ID)
			continue
		}
		
		// Check for loops
		if err := mo.checkLoopDetection(task.ID, "implementation"); err != nil {
			return err
		}
		
		// Execute implementation task
		if err := mo.executeImplementationTask(task, workflow); err != nil {
			task.Status = "failed"
			task.Error = err.Error()
			workflow.FailedTasks = append(workflow.FailedTasks, task.ID)
			mo.logEvent("error", "task_failed", fmt.Sprintf("Task %s failed: %v", task.ID, err), workflow.ID)
			
			// Decide whether to continue or stop
			if !mo.shouldContinueAfterFailure(workflow, err) {
				return fmt.Errorf("critical task failure: %w", err)
			}
		} else {
			task.Status = "completed"
			workflow.CompletedTasks = append(workflow.CompletedTasks, task.ID)
			completedTasks++
		}
		
		// Update progress
		workflow.Progress = 0.3 + (0.4 * float64(completedTasks) / float64(totalTasks))
	}
	
	workflow.Progress = 0.7
	mo.logEvent("info", "implementation_stage_completed", fmt.Sprintf("Implementation completed: %d/%d tasks successful", completedTasks, totalTasks), workflow.ID)
	
	return nil
}

// executeVerificationStage runs the verification phase
func (mo *MetaOrchestrator) executeVerificationStage(workflow *Workflow) error {
	mo.logEvent("info", "verification_stage_started", "Starting verification stage", workflow.ID)
	
	workflow.CurrentStage = "verification"
	workflow.Progress = 0.8
	
	// Execute verification using the verification agent
	verificationAgent := NewVerificationAgent(mo.database)
	
	// For now, verify the overall workflow
	// In a real implementation, this would verify each task individually
	verificationRequest := &VerificationRequest{
		TaskID:        workflow.ID,
		BlockType:     "workflow",
		GeneratedCode: "// Workflow implementation completed",
		TestCode:      "// Test implementation",
		Requirements:  workflow.Requirements,
		Context: map[string]interface{}{
			"workflow_id": workflow.ID,
			"project_path": workflow.ProjectPath,
		},
		CreatedAt:     time.Now(),
	}
	
	verificationResult, err := verificationAgent.VerifyImplementation(verificationRequest)
	if err != nil {
		mo.logEvent("error", "verification_failed", fmt.Sprintf("Verification failed: %v", err), workflow.ID)
		return fmt.Errorf("verification failed: %w", err)
	}
	
	// Store verification results
	workflow.VerificationResults = verificationResult
	
	workflow.Progress = 0.9
	mo.logEvent("info", "verification_stage_completed", "Verification stage completed", workflow.ID)
	
	return nil
}

// executeReviewStage runs the review phase
func (mo *MetaOrchestrator) executeReviewStage(workflow *Workflow) error {
	mo.logEvent("info", "review_stage_started", "Starting review stage", workflow.ID)
	
	workflow.CurrentStage = "review"
	workflow.Progress = 0.95
	
	// Execute review using the review agent
	reviewAgent := NewReviewAgent(mo.database)
	
	// For now, review the overall workflow
	// In a real implementation, this would review each task individually
	reviewRequest := &ReviewRequest{
		TaskID:            workflow.ID,
		BlockType:         "workflow",
		GeneratedCode:     "// Workflow implementation completed",
		TestCode:          "// Test implementation",
		VerificationResult: workflow.VerificationResults,
		Context: map[string]interface{}{
			"workflow_id": workflow.ID,
			"project_path": workflow.ProjectPath,
		},
		CreatedAt:         time.Now(),
	}
	
	reviewResult, err := reviewAgent.ReviewCode(reviewRequest)
	if err != nil {
		mo.logEvent("error", "review_failed", fmt.Sprintf("Review failed: %v", err), workflow.ID)
		return fmt.Errorf("review failed: %w", err)
	}
	
	// Store review results
	workflow.ReviewResults = reviewResult
	
	workflow.Progress = 1.0
	mo.logEvent("info", "review_stage_completed", "Review stage completed", workflow.ID)
	
	return nil
}

// executeImplementationTask executes a single implementation task
func (mo *MetaOrchestrator) executeImplementationTask(task *WorkflowTask, workflow *Workflow) error {
	task.Status = "running"
	startTime := time.Now()
	task.StartedAt = &startTime
	
	// Prepare implementation request
	request := &CodeGenRequest{
		TaskID:      task.ID,
		BlockType:   mo.getBlockTypeFromTask(task),
		Description: mo.getDescriptionFromTask(task),
		Requirements: []string{"RIFX compatible", "Error handling", "Test coverage"},
		Context:     task.Input,
		CreatedAt:   time.Now(),
	}
	
	// Execute implementation
	result, err := mo.implementAgent.GenerateImplementation(request)
	if err != nil {
		task.RetryCount++
		if task.RetryCount < task.MaxRetries {
			mo.logEvent("warning", "task_retry", fmt.Sprintf("Task %s failed, retrying (%d/%d)", task.ID, task.RetryCount, task.MaxRetries), workflow.ID)
			return mo.executeImplementationTask(task, workflow) // Retry
		}
		return err
	}
	
	// Store result
	task.Output = map[string]interface{}{
		"generated_code": result.GeneratedCode,
		"test_code":      result.TestCode,
		"model_used":     result.ModelUsed,
		"confidence":     result.Confidence,
		"status":         result.Status,
	}
	
	completedTime := time.Now()
	task.CompletedAt = &completedTime
	task.ActualTime = completedTime.Sub(*task.StartedAt)
	
	return nil
}

// checkTaskDependencies verifies all task dependencies are met
func (mo *MetaOrchestrator) checkTaskDependencies(task *WorkflowTask, workflow *Workflow) bool {
	for _, depID := range task.Dependencies {
		found := false
		for _, completedID := range workflow.CompletedTasks {
			if depID == completedID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// checkLoopDetection prevents infinite loops
func (mo *MetaOrchestrator) checkLoopDetection(taskID, stage string) error {
	key := fmt.Sprintf("%s:%s", taskID, stage)
	mo.loopDetection[key]++
	
	if mo.loopDetection[key] > mo.maxIterations {
		mo.logEvent("critical", "loop_detected", fmt.Sprintf("Loop detected for %s after %d iterations", key, mo.maxIterations), "")
		return fmt.Errorf("loop detected: %s has been executed %d times", key, mo.maxIterations)
	}
	
	mo.currentIteration++
	if mo.currentIteration > mo.maxIterations {
		return fmt.Errorf("maximum iterations exceeded: %d", mo.maxIterations)
	}
	
	return nil
}

// shouldContinueAfterFailure determines if workflow should continue after task failure
func (mo *MetaOrchestrator) shouldContinueAfterFailure(workflow *Workflow, err error) bool {
	// Continue for non-critical failures
	failureRate := float64(len(workflow.FailedTasks)) / float64(len(workflow.Tasks))
	
	// Stop if failure rate is too high
	if failureRate > 0.5 {
		return false
	}
	
	// Continue for most errors
	return true
}

// completeWorkflow finalizes workflow execution
func (mo *MetaOrchestrator) completeWorkflow(workflow *Workflow) error {
	completedTime := time.Now()
	workflow.CompletedAt = &completedTime
	
	if workflow.StartedAt != nil {
		workflow.ActualTime = completedTime.Sub(*workflow.StartedAt)
	}
	
	// Determine final status
	if len(workflow.FailedTasks) == 0 {
		workflow.Status = "completed"
	} else if len(workflow.CompletedTasks) > len(workflow.FailedTasks) {
		workflow.Status = "completed_with_errors"
	} else {
		workflow.Status = "failed"
	}
	
	workflow.Progress = 1.0
	
	// Save final state
	if err := mo.saveWorkflow(workflow); err != nil {
		mo.logEvent("error", "workflow_save_failed", err.Error(), workflow.ID)
	}
	
	mo.updateMetrics()
	mo.logEvent("info", "workflow_completed", fmt.Sprintf("Workflow %s completed with status: %s", workflow.ID, workflow.Status), workflow.ID)
	
	return nil
}

// handleWorkflowError handles workflow execution errors
func (mo *MetaOrchestrator) handleWorkflowError(workflow *Workflow, stage string, err error) error {
	workflow.Status = "failed"
	completedTime := time.Now()
	workflow.CompletedAt = &completedTime
	
	if workflow.StartedAt != nil {
		workflow.ActualTime = completedTime.Sub(*workflow.StartedAt)
	}
	
	mo.logEvent("error", "workflow_failed", fmt.Sprintf("Workflow %s failed in %s stage: %v", workflow.ID, stage, err), workflow.ID)
	
	// Save error state
	mo.saveWorkflow(workflow)
	mo.updateMetrics()
	
	return fmt.Errorf("workflow failed in %s stage: %w", stage, err)
}

// requestWorkflowApproval creates an approval request for human oversight
func (mo *MetaOrchestrator) requestWorkflowApproval(workflow *Workflow) error {
	approval := &ApprovalRequest{
		ID:          generateApprovalID(),
		WorkflowID:  workflow.ID,
		Type:        "workflow_start",
		Title:       fmt.Sprintf("Approve AEP Processing Workflow: %s", workflow.Name),
		Description: fmt.Sprintf("Automated processing of AEP file: %s", workflow.ProjectPath),
		Context: map[string]interface{}{
			"project_path": workflow.ProjectPath,
			"description":  workflow.Description,
			"estimated_time": "10-30 minutes",
		},
		Options:        []string{"Approve", "Reject", "Request Changes"},
		Recommendation: "Approve - Standard AEP processing workflow",
		RiskLevel:      "medium",
		CreatedAt:      time.Now(),
		Status:         "pending",
	}
	
	// Set deadline (1 hour from now)
	deadline := time.Now().Add(time.Hour)
	approval.Deadline = &deadline
	
	mo.approvalMutex.Lock()
	mo.pendingApprovals[approval.ID] = approval
	mo.approvalMutex.Unlock()
	
	mo.logEvent("info", "approval_requested", fmt.Sprintf("Approval requested for workflow %s", workflow.ID), workflow.ID)
	
	return mo.saveApprovalRequest(approval)
}

// shouldRequireApproval determines if a workflow needs human approval
func (mo *MetaOrchestrator) shouldRequireApproval(projectPath string) bool {
	// Require approval for:
	// - Large files (>10MB)
	// - Files with specific patterns
	// - First time processing new file types
	
	if info, err := os.Stat(projectPath); err == nil {
		if info.Size() > 10*1024*1024 { // 10MB
			return true
		}
	}
	
	// Check for high-risk patterns
	riskPatterns := []string{"experimental", "beta", "test", "sample"}
	for _, pattern := range riskPatterns {
		if strings.Contains(strings.ToLower(projectPath), pattern) {
			return true
		}
	}
	
	return false
}

// Helper methods for task data extraction
func (mo *MetaOrchestrator) getBlockTypeFromTask(task *WorkflowTask) string {
	if blockTypes, ok := task.Input["block_type"].([]string); ok && len(blockTypes) > 0 {
		return blockTypes[0]
	}
	return "Item" // Default
}

func (mo *MetaOrchestrator) getDescriptionFromTask(task *WorkflowTask) string {
	if desc, ok := task.Input["description"].(string); ok {
		return desc
	}
	return "Parse AEP block data"
}

// logEvent adds an event to the orchestration log
func (mo *MetaOrchestrator) logEvent(level, eventType, message, workflowID string) {
	event := OrchestrationEvent{
		ID:        generateEventID(),
		Type:      eventType,
		Level:     level,
		Message:   message,
		WorkflowID: workflowID,
		Timestamp: time.Now(),
	}
	
	mo.eventMutex.Lock()
	mo.eventLog = append(mo.eventLog, event)
	
	// Keep only last 1000 events
	if len(mo.eventLog) > 1000 {
		mo.eventLog = mo.eventLog[1:]
	}
	mo.eventMutex.Unlock()
	
	// Also log to standard logger
	log.Printf("[%s] %s: %s", strings.ToUpper(level), eventType, message)
}

// updateMetrics recalculates orchestration metrics
func (mo *MetaOrchestrator) updateMetrics() {
	mo.workflowMutex.RLock()
	defer mo.workflowMutex.RUnlock()
	
	mo.metrics.TotalWorkflows = len(mo.workflows)
	mo.metrics.ActiveWorkflows = 0
	mo.metrics.CompletedWorkflows = 0
	mo.metrics.FailedWorkflows = 0
	
	totalTime := time.Duration(0)
	completedCount := 0
	
	for _, workflow := range mo.workflows {
		switch workflow.Status {
		case "running":
			mo.metrics.ActiveWorkflows++
		case "completed", "completed_with_errors":
			mo.metrics.CompletedWorkflows++
			if workflow.ActualTime > 0 {
				totalTime += workflow.ActualTime
				completedCount++
			}
		case "failed":
			mo.metrics.FailedWorkflows++
		}
	}
	
	if completedCount > 0 {
		mo.metrics.AvgWorkflowTime = totalTime / time.Duration(completedCount)
	}
	
	mo.metrics.LastUpdated = time.Now()
}

// Database operations
func (mo *MetaOrchestrator) createOrchestratorTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS workflows (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			project_path TEXT NOT NULL,
			status TEXT NOT NULL,
			current_stage TEXT,
			progress REAL NOT NULL,
			requires_approval BOOLEAN NOT NULL,
			iteration_count INTEGER NOT NULL,
			workflow_data TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			started_at INTEGER,
			completed_at INTEGER
		)`,
		`CREATE TABLE IF NOT EXISTS approval_requests (
			id TEXT PRIMARY KEY,
			workflow_id TEXT NOT NULL,
			task_id TEXT,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			risk_level TEXT NOT NULL,
			status TEXT NOT NULL,
			request_data TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			deadline INTEGER,
			responded_at INTEGER
		)`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_status ON workflows(status)`,
		`CREATE INDEX IF NOT EXISTS idx_approvals_status ON approval_requests(status)`,
	}
	
	for _, query := range queries {
		if _, err := mo.database.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}
	
	return nil
}

func (mo *MetaOrchestrator) saveWorkflow(workflow *Workflow) error {
	workflowJSON, err := json.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}
	
	query := `
		INSERT OR REPLACE INTO workflows
		(id, name, description, project_path, status, current_stage, progress,
		 requires_approval, iteration_count, workflow_data, created_at, started_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	var startedAtUnix, completedAtUnix *int64
	if workflow.StartedAt != nil {
		t := workflow.StartedAt.Unix()
		startedAtUnix = &t
	}
	if workflow.CompletedAt != nil {
		t := workflow.CompletedAt.Unix()
		completedAtUnix = &t
	}
	
	_, err = mo.database.db.Exec(query,
		workflow.ID, workflow.Name, workflow.Description, workflow.ProjectPath,
		workflow.Status, workflow.CurrentStage, workflow.Progress,
		workflow.RequiresApproval, workflow.IterationCount,
		string(workflowJSON), workflow.CreatedAt.Unix(),
		startedAtUnix, completedAtUnix)
	
	return err
}

func (mo *MetaOrchestrator) saveApprovalRequest(approval *ApprovalRequest) error {
	approvalJSON, err := json.Marshal(approval)
	if err != nil {
		return fmt.Errorf("failed to marshal approval: %w", err)
	}
	
	query := `
		INSERT INTO approval_requests
		(id, workflow_id, task_id, type, title, description, risk_level,
		 status, request_data, created_at, deadline, responded_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	var deadlineUnix, respondedAtUnix *int64
	if approval.Deadline != nil {
		t := approval.Deadline.Unix()
		deadlineUnix = &t
	}
	if approval.RespondedAt != nil {
		t := approval.RespondedAt.Unix()
		respondedAtUnix = &t
	}
	
	_, err = mo.database.db.Exec(query,
		approval.ID, approval.WorkflowID, approval.TaskID,
		approval.Type, approval.Title, approval.Description,
		approval.RiskLevel, approval.Status, string(approvalJSON),
		approval.CreatedAt.Unix(), deadlineUnix, respondedAtUnix)
	
	return err
}

// Public API methods

// GetWorkflowStatus returns current workflow status
func (mo *MetaOrchestrator) GetWorkflowStatus(workflowID string) (*Workflow, error) {
	mo.workflowMutex.RLock()
	workflow, exists := mo.workflows[workflowID]
	mo.workflowMutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}
	
	return workflow, nil
}

// GetPendingApprovals returns all pending approval requests
func (mo *MetaOrchestrator) GetPendingApprovals() []*ApprovalRequest {
	mo.approvalMutex.RLock()
	defer mo.approvalMutex.RUnlock()
	
	var approvals []*ApprovalRequest
	for _, approval := range mo.pendingApprovals {
		if approval.Status == "pending" {
			approvals = append(approvals, approval)
		}
	}
	
	return approvals
}

// ApproveWorkflow approves a pending workflow
func (mo *MetaOrchestrator) ApproveWorkflow(approvalID, response, approvedBy string) error {
	mo.approvalMutex.Lock()
	approval, exists := mo.pendingApprovals[approvalID]
	mo.approvalMutex.Unlock()
	
	if !exists {
		return fmt.Errorf("approval request %s not found", approvalID)
	}
	
	approval.Status = "approved"
	approval.Response = response
	approval.RespondedBy = approvedBy
	respondedTime := time.Now()
	approval.RespondedAt = &respondedTime
	
	// Update workflow status
	if workflow, exists := mo.workflows[approval.WorkflowID]; exists {
		workflow.RequiresApproval = false
		workflow.Status = "pending" // Ready for execution
	}
	
	mo.logEvent("info", "workflow_approved", fmt.Sprintf("Workflow %s approved by %s", approval.WorkflowID, approvedBy), approval.WorkflowID)
	
	return mo.saveApprovalRequest(approval)
}

// GetMetrics returns current orchestration metrics
func (mo *MetaOrchestrator) GetMetrics() *OrchestrationMetrics {
	mo.updateMetrics()
	return mo.metrics
}

// GetRecentEvents returns recent orchestration events
func (mo *MetaOrchestrator) GetRecentEvents(limit int) []OrchestrationEvent {
	mo.eventMutex.RLock()
	defer mo.eventMutex.RUnlock()
	
	if limit <= 0 || limit > len(mo.eventLog) {
		limit = len(mo.eventLog)
	}
	
	start := len(mo.eventLog) - limit
	events := make([]OrchestrationEvent, limit)
	copy(events, mo.eventLog[start:])
	
	return events
}

// Helper functions for ID generation
func generateWorkflowID() string {
	return fmt.Sprintf("wf_%d", time.Now().UnixNano())
}

func generateApprovalID() string {
	return fmt.Sprintf("ap_%d", time.Now().UnixNano())
}

// generateEventID moved to agent_communication.go to avoid duplication