// Package catalog provides implementation agent for code generation and execution
package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// ImplementationAgent extends parser code and generates functional implementations
type ImplementationAgent struct {
	database        *Database
	planningAgent   *PlanningAgent
	modelCascade    []string // Model priority: Claude → GPT-4 → Gemini
	codeTemplates   map[string]string
	rifxPatterns    map[string]RIFXPattern
	maxRetries      int
}

// CodeGenRequest represents a code generation request
type CodeGenRequest struct {
	TaskID          string                 `json:"task_id"`
	BlockType       string                 `json:"block_type"`
	Description     string                 `json:"description"`
	ExistingCode    string                 `json:"existing_code,omitempty"`
	Requirements    []string               `json:"requirements"`
	Context         map[string]interface{} `json:"context"`
	Model           string                 `json:"model,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// CodeGenResult contains generated implementation
type CodeGenResult struct {
	TaskID       string            `json:"task_id"`
	BlockType    string            `json:"block_type"`
	GeneratedCode string           `json:"generated_code"`
	TestCode     string            `json:"test_code,omitempty"`
	ModelUsed    string            `json:"model_used"`
	Confidence   float64           `json:"confidence"`
	Integration  IntegrationInfo   `json:"integration"`
	Metrics      GenerationMetrics `json:"metrics"`
	Status       string            `json:"status"` // success, partial, failed
	Error        string            `json:"error,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
}

// IntegrationInfo describes how code integrates with existing system
type IntegrationInfo struct {
	PackageName     string   `json:"package_name"`
	Dependencies    []string `json:"dependencies"`
	ExportedFuncs   []string `json:"exported_functions"`
	InterfaceCompat bool     `json:"interface_compatible"`
	BackwardCompat  bool     `json:"backward_compatible"`
}

// GenerationMetrics tracks code generation performance
type GenerationMetrics struct {
	GenerationTime  time.Duration `json:"generation_time"`
	CodeLines       int           `json:"code_lines"`
	TestLines       int           `json:"test_lines"`
	CyclomaticComplexity int      `json:"cyclomatic_complexity"`
	TestCoverage    float64       `json:"test_coverage"`
	RetryCount      int           `json:"retry_count"`
}

// RIFXPattern defines patterns for RIFX block parsing
type RIFXPattern struct {
	BlockType       string            `json:"block_type"`
	ByteSignature   []byte            `json:"byte_signature"`
	StructTemplate  string            `json:"struct_template"`
	ParserTemplate  string            `json:"parser_template"`
	ValidationRules []ValidationRule  `json:"validation_rules"`
	Examples        []PatternExample  `json:"examples"`
}

// ValidationRule defines validation logic for parsed data
type ValidationRule struct {
	Field       string `json:"field"`
	Type        string `json:"type"`        // required, range, format, custom
	Rule        string `json:"rule"`        // validation expression
	ErrorMsg    string `json:"error_msg"`
}

// PatternExample provides working examples for pattern matching
type PatternExample struct {
	Description string `json:"description"`
	Input       []byte `json:"input"`
	Expected    string `json:"expected"`
	ParsedData  string `json:"parsed_data"`
}

// NewImplementationAgent creates a new implementation agent
func NewImplementationAgent(database *Database, planningAgent *PlanningAgent) *ImplementationAgent {
	agent := &ImplementationAgent{
		database:      database,
		planningAgent: planningAgent,
		modelCascade:  []string{"claude", "gpt-4", "gemini"}, // Priority order
		maxRetries:    3,
		codeTemplates: make(map[string]string),
		rifxPatterns:  make(map[string]RIFXPattern),
	}
	
	// Initialize code templates
	agent.initializeCodeTemplates()
	
	// Initialize RIFX patterns
	agent.initializeRIFXPatterns()
	
	return agent
}

// GenerateImplementation creates functional code for parsing tasks
func (ia *ImplementationAgent) GenerateImplementation(request *CodeGenRequest) (*CodeGenResult, error) {
	log.Printf("Implementation Agent: Generating code for %s (%s)", request.TaskID, request.BlockType)
	
	startTime := time.Now()
	
	result := &CodeGenResult{
		TaskID:    request.TaskID,
		BlockType: request.BlockType,
		CreatedAt: time.Now(),
	}
	
	// 1. Analyze existing patterns
	pattern, exists := ia.rifxPatterns[request.BlockType]
	if !exists {
		log.Printf("Warning: No RIFX pattern found for %s, using generic pattern", request.BlockType)
		pattern = ia.getGenericPattern(request.BlockType)
	}
	
	// 2. Generate code with model cascade
	var err error
	for attempt := 0; attempt < ia.maxRetries; attempt++ {
		modelIndex := attempt % len(ia.modelCascade)
		model := ia.modelCascade[modelIndex]
		
		log.Printf("Attempting code generation with %s (attempt %d)", model, attempt+1)
		
		err = ia.generateCodeWithModel(request, &pattern, model, result)
		if err == nil && result.Status == "success" {
			result.ModelUsed = model
			result.Metrics.RetryCount = attempt
			break
		}
		
		log.Printf("Generation failed with %s: %v", model, err)
	}
	
	if err != nil && result.Status != "success" {
		result.Status = "failed"
		result.Error = fmt.Sprintf("All models failed: %v", err)
		return result, err
	}
	
	// 3. Generate integration info
	ia.analyzeIntegration(result)
	
	// 4. Calculate metrics
	result.Metrics.GenerationTime = time.Since(startTime)
	ia.calculateMetrics(result)
	
	// 5. Store result
	if err := ia.storeImplementationResult(result); err != nil {
		log.Printf("Warning: Failed to store implementation result: %v", err)
	}
	
	log.Printf("Implementation Agent: Generated %d lines of code in %v", 
		result.Metrics.CodeLines, result.Metrics.GenerationTime)
	
	return result, nil
}

// generateCodeWithModel attempts code generation with specific model
func (ia *ImplementationAgent) generateCodeWithModel(request *CodeGenRequest, pattern *RIFXPattern, model string, result *CodeGenResult) error {
	// Simulate model-specific code generation
	switch model {
	case "claude":
		return ia.generateWithClaude(request, pattern, result)
	case "gpt-4":
		return ia.generateWithGPT4(request, pattern, result)
	case "gemini":
		return ia.generateWithGemini(request, pattern, result)
	default:
		return fmt.Errorf("unknown model: %s", model)
	}
}

// generateWithClaude simulates Claude-based code generation
func (ia *ImplementationAgent) generateWithClaude(request *CodeGenRequest, pattern *RIFXPattern, result *CodeGenResult) error {
	// Claude is excellent at structured code generation
	template := ia.getClaudeTemplate(request.BlockType)
	
	generatedCode := ia.expandTemplate(template, map[string]interface{}{
		"BlockType":    request.BlockType,
		"Description":  request.Description,
		"Pattern":      pattern,
		"Requirements": request.Requirements,
		"Context":      request.Context,
	})
	
	// Generate comprehensive test code
	testCode := ia.generateTestCode(request.BlockType, generatedCode)
	
	result.GeneratedCode = generatedCode
	result.TestCode = testCode
	result.Status = "success"
	result.Confidence = 0.90 // Claude typically produces high-quality code
	
	return nil
}

// generateWithGPT4 simulates GPT-4 based code generation
func (ia *ImplementationAgent) generateWithGPT4(request *CodeGenRequest, pattern *RIFXPattern, result *CodeGenResult) error {
	// GPT-4 is strong at complex logic
	template := ia.getGPT4Template(request.BlockType)
	
	generatedCode := ia.expandTemplate(template, map[string]interface{}{
		"BlockType":    request.BlockType,
		"Description":  request.Description,
		"Pattern":      pattern,
		"Requirements": request.Requirements,
	})
	
	testCode := ia.generateBasicTestCode(request.BlockType, generatedCode)
	
	result.GeneratedCode = generatedCode
	result.TestCode = testCode
	result.Status = "success"
	result.Confidence = 0.85 // Good quality, slightly lower than Claude
	
	return nil
}

// generateWithGemini simulates Gemini-based code generation
func (ia *ImplementationAgent) generateWithGemini(request *CodeGenRequest, pattern *RIFXPattern, result *CodeGenResult) error {
	// Gemini provides functional but basic implementations
	template := ia.getGeminiTemplate(request.BlockType)
	
	generatedCode := ia.expandTemplate(template, map[string]interface{}{
		"BlockType":   request.BlockType,
		"Description": request.Description,
		"Pattern":     pattern,
	})
	
	result.GeneratedCode = generatedCode
	result.TestCode = "" // Basic implementation without tests
	result.Status = "partial"
	result.Confidence = 0.70 // Functional but may need refinement
	
	return nil
}

// initializeCodeTemplates sets up code generation templates
func (ia *ImplementationAgent) initializeCodeTemplates() {
	// Claude templates (comprehensive, well-documented)
	ia.codeTemplates["claude_item"] = `
// Parse{{.BlockType}} parses {{.Description}}
func Parse{{.BlockType}}(data []byte, offset int) (*{{.BlockType}}, int, error) {
	if len(data) < offset+{{.Pattern.MinSize}} {
		return nil, 0, fmt.Errorf("insufficient data for {{.BlockType}}")
	}
	
	item := &{{.BlockType}}{
		Type:      "{{.BlockType}}",
		Timestamp: time.Now(),
	}
	
	// Parse header
	if err := item.parseHeader(data[offset:]); err != nil {
		return nil, 0, fmt.Errorf("failed to parse header: %w", err)
	}
	
	// Parse content with validation
	contentSize, err := item.parseContent(data[offset+HeaderSize:])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse content: %w", err)
	}
	
	// Validate parsed data
	if err := item.validate(); err != nil {
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}
	
	totalSize := HeaderSize + contentSize
	return item, totalSize, nil
}

// parseHeader extracts header information
func (item *{{.BlockType}}) parseHeader(data []byte) error {
	// Implementation based on RIFX specification
	item.Size = binary.BigEndian.Uint32(data[0:4])
	item.ID = binary.BigEndian.Uint32(data[4:8])
	return nil
}

// parseContent extracts main content
func (item *{{.BlockType}}) parseContent(data []byte) (int, error) {
	// Implementation specific to {{.BlockType}}
	// Returns size of parsed content
	return int(item.Size) - HeaderSize, nil
}

// validate checks parsed data integrity
func (item *{{.BlockType}}) validate() error {
	{{range .Pattern.ValidationRules}}
	// {{.ErrorMsg}}
	{{end}}
	return nil
}
`
	
	// Similar templates for GPT-4 and Gemini (shorter versions)
	ia.codeTemplates["gpt4_item"] = `
func Parse{{.BlockType}}(data []byte, offset int) (*{{.BlockType}}, int, error) {
	item := &{{.BlockType}}{}
	
	// Basic parsing logic
	if len(data) < offset+8 {
		return nil, 0, errors.New("insufficient data")
	}
	
	item.Size = binary.BigEndian.Uint32(data[offset:])
	item.ID = binary.BigEndian.Uint32(data[offset+4:])
	
	return item, int(item.Size), nil
}
`
	
	ia.codeTemplates["gemini_item"] = `
func Parse{{.BlockType}}(data []byte, offset int) (*{{.BlockType}}, int, error) {
	// Basic implementation
	item := &{{.BlockType}}{}
	item.Size = binary.BigEndian.Uint32(data[offset:])
	return item, int(item.Size), nil
}
`
}

// initializeRIFXPatterns sets up RIFX parsing patterns
func (ia *ImplementationAgent) initializeRIFXPatterns() {
	// Item block pattern
	ia.rifxPatterns["Item"] = RIFXPattern{
		BlockType:     "Item",
		ByteSignature: []byte{0x49, 0x74, 0x65, 0x6D}, // "Item"
		StructTemplate: `
type Item struct {
	Size      uint32    ` + "`json:\"size\"`" + `
	ID        uint32    ` + "`json:\"id\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	Type      string    ` + "`json:\"type\"`" + `
	Timestamp time.Time ` + "`json:\"timestamp\"`" + `
}`,
		ParserTemplate: "claude_item",
		ValidationRules: []ValidationRule{
			{Field: "Size", Type: "range", Rule: "> 0", ErrorMsg: "Item size must be positive"},
			{Field: "Name", Type: "required", Rule: "!= \"\"", ErrorMsg: "Item name is required"},
		},
	}
	
	// Layer block pattern  
	ia.rifxPatterns["Layer"] = RIFXPattern{
		BlockType:     "Layer",
		ByteSignature: []byte{0x4C, 0x61, 0x79, 0x72}, // "Layr"
		StructTemplate: `
type Layer struct {
	Size        uint32    ` + "`json:\"size\"`" + `
	ID          uint32    ` + "`json:\"id\"`" + `
	Name        string    ` + "`json:\"name\"`" + `
	Type        string    ` + "`json:\"type\"`" + `
	Visible     bool      ` + "`json:\"visible\"`" + `
	Opacity     float32   ` + "`json:\"opacity\"`" + `
	BlendMode   string    ` + "`json:\"blend_mode\"`" + `
	Transform   Transform ` + "`json:\"transform\"`" + `
	Timestamp   time.Time ` + "`json:\"timestamp\"`" + `
}`,
		ParserTemplate: "claude_item",
		ValidationRules: []ValidationRule{
			{Field: "Opacity", Type: "range", Rule: "0.0 <= x <= 1.0", ErrorMsg: "Opacity must be between 0.0 and 1.0"},
			{Field: "BlendMode", Type: "format", Rule: "valid_blend_mode", ErrorMsg: "Invalid blend mode"},
		},
	}
	
	// Add more patterns for Property, Composition, etc.
	ia.addAdditionalPatterns()
}

// addAdditionalPatterns adds patterns for other block types
func (ia *ImplementationAgent) addAdditionalPatterns() {
	// Property pattern
	ia.rifxPatterns["Property"] = RIFXPattern{
		BlockType:     "Property",
		ByteSignature: []byte{0x50, 0x72, 0x6F, 0x70}, // "Prop"
		StructTemplate: `
type Property struct {
	Size      uint32      ` + "`json:\"size\"`" + `
	ID        uint32      ` + "`json:\"id\"`" + `
	Name      string      ` + "`json:\"name\"`" + `
	Type      string      ` + "`json:\"type\"`" + `
	Value     interface{} ` + "`json:\"value\"`" + `
	Animated  bool        ` + "`json:\"animated\"`" + `
	Keyframes []Keyframe  ` + "`json:\"keyframes,omitempty\"`" + `
	Timestamp time.Time   ` + "`json:\"timestamp\"`" + `
}`,
	}
	
	// Composition pattern
	ia.rifxPatterns["Composition"] = RIFXPattern{
		BlockType:     "Composition",
		ByteSignature: []byte{0x43, 0x6F, 0x6D, 0x70}, // "Comp"
		StructTemplate: `
type Composition struct {
	Size       uint32    ` + "`json:\"size\"`" + `
	ID         uint32    ` + "`json:\"id\"`" + `
	Name       string    ` + "`json:\"name\"`" + `
	Width      int       ` + "`json:\"width\"`" + `
	Height     int       ` + "`json:\"height\"`" + `
	Duration   float64   ` + "`json:\"duration\"`" + `
	FrameRate  float64   ` + "`json:\"frame_rate\"`" + `
	Layers     []Layer   ` + "`json:\"layers\"`" + `
	Timestamp  time.Time ` + "`json:\"timestamp\"`" + `
}`,
		ValidationRules: []ValidationRule{
			{Field: "Width", Type: "range", Rule: "> 0", ErrorMsg: "Width must be positive"},
			{Field: "Height", Type: "range", Rule: "> 0", ErrorMsg: "Height must be positive"},
			{Field: "FrameRate", Type: "range", Rule: "> 0", ErrorMsg: "Frame rate must be positive"},
		},
	}
}

// getGenericPattern returns a generic pattern for unknown block types
func (ia *ImplementationAgent) getGenericPattern(blockType string) RIFXPattern {
	return RIFXPattern{
		BlockType:     blockType,
		ByteSignature: []byte{0x47, 0x65, 0x6E, 0x72}, // "Genr"
		StructTemplate: fmt.Sprintf(`
type %s struct {
	Size      uint32                 ` + "`json:\"size\"`" + `
	ID        uint32                 ` + "`json:\"id\"`" + `
	Type      string                 ` + "`json:\"type\"`" + `
	Data      []byte                 ` + "`json:\"data\"`" + `
	Metadata  map[string]interface{} ` + "`json:\"metadata\"`" + `
	Timestamp time.Time              ` + "`json:\"timestamp\"`" + `
}`, blockType),
		ParserTemplate: "claude_item",
	}
}

// getClaudeTemplate returns Claude-optimized template
func (ia *ImplementationAgent) getClaudeTemplate(blockType string) string {
	if template, exists := ia.codeTemplates["claude_"+strings.ToLower(blockType)]; exists {
		return template
	}
	return ia.codeTemplates["claude_item"] // Default
}

// getGPT4Template returns GPT-4 optimized template
func (ia *ImplementationAgent) getGPT4Template(blockType string) string {
	if template, exists := ia.codeTemplates["gpt4_"+strings.ToLower(blockType)]; exists {
		return template
	}
	return ia.codeTemplates["gpt4_item"] // Default
}

// getGeminiTemplate returns Gemini optimized template
func (ia *ImplementationAgent) getGeminiTemplate(blockType string) string {
	if template, exists := ia.codeTemplates["gemini_"+strings.ToLower(blockType)]; exists {
		return template
	}
	return ia.codeTemplates["gemini_item"] // Default
}

// expandTemplate replaces template variables with actual values
func (ia *ImplementationAgent) expandTemplate(template string, vars map[string]interface{}) string {
	result := template
	
	// Simple template expansion (in production, use text/template)
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}
	
	// Handle pattern-specific replacements
	if pattern, ok := vars["Pattern"].(*RIFXPattern); ok {
		result = strings.ReplaceAll(result, "{{.Pattern.MinSize}}", "8")
		
		// Add validation rules
		if len(pattern.ValidationRules) > 0 {
			validations := ""
			for _, rule := range pattern.ValidationRules {
				validations += fmt.Sprintf("if !(%s) { return fmt.Errorf(\"%s\") }\n\t", 
					rule.Rule, rule.ErrorMsg)
			}
			result = strings.ReplaceAll(result, "{{range .Pattern.ValidationRules}}", validations)
			result = strings.ReplaceAll(result, "{{end}}", "")
		}
	}
	
	// Clean up any remaining template syntax
	result = strings.ReplaceAll(result, "{{range", "// Range:")
	result = strings.ReplaceAll(result, "{{end}}", "")
	
	return result
}

// generateTestCode creates comprehensive test code
func (ia *ImplementationAgent) generateTestCode(blockType, generatedCode string) string {
	return fmt.Sprintf(`
package catalog

import (
	"testing"
	"time"
)

func TestParse%s(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		offset   int
		wantErr  bool
	}{
		{
			name:    "valid %s",
			data:    []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x01}, // size=16, id=1
			offset:  0,
			wantErr: false,
		},
		{
			name:    "insufficient data",
			data:    []byte{0x00, 0x00},
			offset:  0,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, size, err := Parse%s(tt.data, tt.offset)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse%s() error = %%v, wantErr %%v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if item == nil {
					t.Error("Parse%s() returned nil item")
				}
				if size <= 0 {
					t.Error("Parse%s() returned invalid size")
				}
			}
		})
	}
}

func BenchmarkParse%s(b *testing.B) {
	data := make([]byte, 1024)
	// Initialize with valid %s data
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Parse%s(data, 0)
	}
}
`, blockType, blockType, blockType, blockType, blockType, blockType, blockType, blockType, blockType)
}

// generateBasicTestCode creates basic test code
func (ia *ImplementationAgent) generateBasicTestCode(blockType, generatedCode string) string {
	return fmt.Sprintf(`
func TestParse%s(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x01}
	item, size, err := Parse%s(data, 0)
	
	if err != nil {
		t.Errorf("Parse%s() failed: %%v", err)
	}
	if item == nil {
		t.Error("Parse%s() returned nil")
	}
	if size <= 0 {
		t.Error("Parse%s() returned invalid size")
	}
}
`, blockType, blockType, blockType, blockType, blockType)
}

// analyzeIntegration determines how generated code integrates
func (ia *ImplementationAgent) analyzeIntegration(result *CodeGenResult) {
	result.Integration = IntegrationInfo{
		PackageName:     "catalog",
		Dependencies:    ia.extractDependencies(result.GeneratedCode),
		ExportedFuncs:   ia.extractExportedFunctions(result.GeneratedCode),
		InterfaceCompat: ia.checkInterfaceCompatibility(result.GeneratedCode),
		BackwardCompat:  true, // Assume backward compatible for now
	}
}

// extractDependencies finds import dependencies in generated code
func (ia *ImplementationAgent) extractDependencies(code string) []string {
	deps := []string{}
	
	// Look for common imports in Go code
	if strings.Contains(code, "binary.BigEndian") {
		deps = append(deps, "encoding/binary")
	}
	if strings.Contains(code, "fmt.Errorf") {
		deps = append(deps, "fmt")
	}
	if strings.Contains(code, "time.Now") {
		deps = append(deps, "time")
	}
	if strings.Contains(code, "errors.New") {
		deps = append(deps, "errors")
	}
	
	return deps
}

// extractExportedFunctions finds exported function names
func (ia *ImplementationAgent) extractExportedFunctions(code string) []string {
	funcs := []string{}
	
	// Simple regex-like detection for Go functions
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func ") && !strings.Contains(line, "func (") {
			// Extract function name
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				funcName := strings.Split(parts[1], "(")[0]
				// Check if exported (starts with capital letter)
				if len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z' {
					funcs = append(funcs, funcName)
				}
			}
		}
	}
	
	return funcs
}

// checkInterfaceCompatibility verifies generated code follows expected interfaces
func (ia *ImplementationAgent) checkInterfaceCompatibility(code string) bool {
	// Check for expected function signatures
	expectedSignatures := []string{
		"func Parse",
		"error",
		"[]byte",
	}
	
	for _, sig := range expectedSignatures {
		if !strings.Contains(code, sig) {
			return false
		}
	}
	
	return true
}

// calculateMetrics computes code quality metrics
func (ia *ImplementationAgent) calculateMetrics(result *CodeGenResult) {
	// Count lines of code
	result.Metrics.CodeLines = len(strings.Split(result.GeneratedCode, "\n"))
	result.Metrics.TestLines = len(strings.Split(result.TestCode, "\n"))
	
	// Simple cyclomatic complexity (count decision points)
	complexity := 1 // Base complexity
	decisionKeywords := []string{"if", "for", "switch", "case"}
	for _, keyword := range decisionKeywords {
		complexity += strings.Count(result.GeneratedCode, keyword)
	}
	result.Metrics.CyclomaticComplexity = complexity
	
	// Estimate test coverage based on test code quality
	if result.TestCode != "" && len(result.TestCode) > 100 {
		if strings.Contains(result.TestCode, "BenchmarkParse") {
			result.Metrics.TestCoverage = 0.90 // Comprehensive tests
		} else {
			result.Metrics.TestCoverage = 0.70 // Basic tests
		}
	} else {
		result.Metrics.TestCoverage = 0.30 // Minimal or no tests
	}
}

// storeImplementationResult saves result to database
func (ia *ImplementationAgent) storeImplementationResult(result *CodeGenResult) error {
	// Create table if it doesn't exist
	if err := ia.createImplementationTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	
	// Convert result to JSON (not used for now, but could be used for debugging)
	_, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	
	// Insert into database
	query := `
		INSERT INTO implementation_results
		(task_id, block_type, generated_code, test_code, model_used, confidence,
		 status, error_message, code_lines, test_coverage, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = ia.database.db.Exec(query,
		result.TaskID,
		result.BlockType,
		result.GeneratedCode,
		result.TestCode,
		result.ModelUsed,
		result.Confidence,
		result.Status,
		result.Error,
		result.Metrics.CodeLines,
		result.Metrics.TestCoverage,
		result.CreatedAt.Unix(),
	)
	
	return err
}

// createImplementationTables creates necessary database tables
func (ia *ImplementationAgent) createImplementationTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS implementation_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		block_type TEXT NOT NULL,
		generated_code TEXT NOT NULL,
		test_code TEXT,
		model_used TEXT NOT NULL,
		confidence REAL NOT NULL,
		status TEXT NOT NULL,
		error_message TEXT,
		code_lines INTEGER NOT NULL,
		test_coverage REAL NOT NULL,
		created_at INTEGER NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_impl_task_id ON implementation_results(task_id);
	CREATE INDEX IF NOT EXISTS idx_impl_block_type ON implementation_results(block_type);
	`
	
	_, err := ia.database.db.Exec(query)
	return err
}

// GetImplementationByTaskID retrieves implementation result by task ID
func (ia *ImplementationAgent) GetImplementationByTaskID(taskID string) (*CodeGenResult, error) {
	query := `
		SELECT task_id, block_type, generated_code, test_code, model_used, 
		       confidence, status, error_message, code_lines, test_coverage, created_at
		FROM implementation_results
		WHERE task_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	var result CodeGenResult
	var createdAtUnix int64
	
	err := ia.database.db.QueryRow(query, taskID).Scan(
		&result.TaskID,
		&result.BlockType,
		&result.GeneratedCode,
		&result.TestCode,
		&result.ModelUsed,
		&result.Confidence,
		&result.Status,
		&result.Error,
		&result.Metrics.CodeLines,
		&result.Metrics.TestCoverage,
		&createdAtUnix,
	)
	
	if err != nil {
		return nil, err
	}
	
	result.CreatedAt = time.Unix(createdAtUnix, 0)
	
	return &result, nil
}

// Helper functions moved to database.go to avoid duplication