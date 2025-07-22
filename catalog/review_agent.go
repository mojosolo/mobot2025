// Package catalog provides review agent for performance optimization and improvement recommendations
package catalog

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

// ReviewAgent analyzes code quality, performance, and provides optimization recommendations
type ReviewAgent struct {
	database         *Database
	optimizers       map[string]Optimizer
	benchmarks       *BenchmarkSuite
	patterns         *ReviewPatternAnalyzer
	recommendations  *RecommendationEngine
}

// Optimizer defines optimization strategies
type Optimizer interface {
	Analyze(code string, metrics map[string]interface{}) OptimizationResult
	GetCategory() string
	GetPriority() int
}

// ReviewRequest represents a code review request
type ReviewRequest struct {
	TaskID            string                 `json:"task_id"`
	BlockType         string                 `json:"block_type"`
	GeneratedCode     string                 `json:"generated_code"`
	TestCode          string                 `json:"test_code"`
	VerificationResult *VerificationResult   `json:"verification_result,omitempty"`
	PerformanceData   map[string]interface{} `json:"performance_data,omitempty"`
	Context           map[string]interface{} `json:"context"`
	CreatedAt         time.Time              `json:"created_at"`
}

// ReviewResult contains comprehensive review analysis
type ReviewResult struct {
	TaskID               string                  `json:"task_id"`
	BlockType            string                  `json:"block_type"`
	Status               string                  `json:"status"` // excellent, good, needs_improvement, poor
	OverallScore         float64                 `json:"overall_score"` // 0.0-1.0
	
	// Optimization Analysis
	PerformanceOptimization OptimizationAnalysis `json:"performance_optimization"`
	CodeOptimization        OptimizationAnalysis `json:"code_optimization"`
	MemoryOptimization      OptimizationAnalysis `json:"memory_optimization"`
	
	// Improvement Recommendations
	CriticalIssues      []ReviewIssue           `json:"critical_issues"`
	Improvements        []ImprovementSuggestion `json:"improvements"`
	BestPractices       []BestPractice          `json:"best_practices"`
	OptimizationPlan    OptimizationPlan        `json:"optimization_plan"`
	
	// Maintainability Analysis
	Maintainability     MaintainabilityAnalysis `json:"maintainability"`
	TechnicalDebt       TechnicalDebtAnalysis   `json:"technical_debt"`
	
	// Comparative Analysis
	BenchmarkComparison BenchmarkComparison     `json:"benchmark_comparison"`
	IndustryStandards   IndustryComparison      `json:"industry_standards"`
	
	// Review Metrics
	ReviewTime          time.Duration           `json:"review_time"`
	CreatedAt           time.Time               `json:"created_at"`
}

// OptimizationAnalysis contains performance optimization insights
type OptimizationAnalysis struct {
	Category            string                 `json:"category"`
	CurrentScore        float64                `json:"current_score"` // 0.0-1.0
	PotentialScore      float64                `json:"potential_score"` // 0.0-1.0
	ImprovementPotential float64               `json:"improvement_potential"` // 0.0-1.0
	BottleneckAreas     []BottleneckArea       `json:"bottleneck_areas"`
	OptimizationOptions []OptimizationOption   `json:"optimization_options"`
	EstimatedGains      PerformanceGains       `json:"estimated_gains"`
}

// ReviewIssue represents a code issue
type ReviewIssue struct {
	ID          string `json:"id"`
	Type        string `json:"type"`        // performance, security, maintainability, correctness
	Severity    string `json:"severity"`    // critical, high, medium, low
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    CodeLocation `json:"location"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`     // high, medium, low
	Solution    string `json:"solution"`
}

// ImprovementSuggestion provides actionable improvement advice
type ImprovementSuggestion struct {
	ID               string            `json:"id"`
	Category         string            `json:"category"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Benefits         []string          `json:"benefits"`
	Implementation   string            `json:"implementation"`
	CodeExample      string            `json:"code_example,omitempty"`
	Priority         int               `json:"priority"` // 1=highest, 5=lowest
	EstimatedEffort  string            `json:"estimated_effort"`
	ExpectedImpact   ImpactAssessment  `json:"expected_impact"`
}

// BestPractice represents coding best practices
type BestPractice struct {
	Category    string `json:"category"`
	Practice    string `json:"practice"`
	Rationale   string `json:"rationale"`
	Example     string `json:"example,omitempty"`
	IsFollowed  bool   `json:"is_followed"`
	Importance  string `json:"importance"` // critical, important, nice-to-have
}

// OptimizationPlan provides structured optimization roadmap
type OptimizationPlan struct {
	Phases          []OptimizationPhase `json:"phases"`
	TotalEstimate   time.Duration       `json:"total_estimate"`
	ExpectedGains   PerformanceGains    `json:"expected_gains"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements"`
}

// OptimizationPhase represents a phase in optimization plan
type OptimizationPhase struct {
	Phase       int                     `json:"phase"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Tasks       []OptimizationTask      `json:"tasks"`
	Duration    time.Duration           `json:"duration"`
	Priority    string                  `json:"priority"`
	Dependencies []int                  `json:"dependencies"`
}

// OptimizationTask represents a specific optimization task
type OptimizationTask struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Type         string           `json:"type"` // refactor, optimize, test, document
	Effort       string           `json:"effort"`
	Impact       ImpactAssessment `json:"impact"`
	Prerequisites []string        `json:"prerequisites"`
}

// MaintainabilityAnalysis assesses code maintainability
type MaintainabilityAnalysis struct {
	OverallScore    float64                `json:"overall_score"`
	Readability     ReadabilityMetrics     `json:"readability"`
	Modularity      ModularityMetrics      `json:"modularity"`
	Documentation   DocumentationMetrics   `json:"documentation"`
	TestCoverage    TestCoverageMetrics    `json:"test_coverage"`
	CodeStructure   CodeStructureMetrics   `json:"code_structure"`
}

// TechnicalDebtAnalysis quantifies technical debt
type TechnicalDebtAnalysis struct {
	DebtLevel       string              `json:"debt_level"` // low, medium, high, critical
	DebtScore       float64             `json:"debt_score"` // 0.0-1.0 (higher = more debt)
	DebtAreas       []DebtArea          `json:"debt_areas"`
	PaydownPlan     DebtPaydownPlan     `json:"paydown_plan"`
	MaintenanceCost MaintenanceCost     `json:"maintenance_cost"`
}

// BenchmarkComparison compares against similar implementations
type ReviewBenchmarkComparison struct {
	Category           string                    `json:"category"`
	CurrentPerformance PerformanceBenchmark     `json:"current_performance"`
	BestInCategory     PerformanceBenchmark     `json:"best_in_category"`
	AveragePerformance PerformanceBenchmark     `json:"average_performance"`
	PerformanceGap     PerformanceGap           `json:"performance_gap"`
	RankingPercentile  float64                  `json:"ranking_percentile"`
}

// IndustryComparison compares against industry standards
type IndustryComparison struct {
	Standards       map[string]StandardMetric `json:"standards"`
	Compliance      map[string]bool           `json:"compliance"`
	ComplianceScore float64                   `json:"compliance_score"`
	Gaps            []ComplianceGap           `json:"gaps"`
}

// Supporting structs
type BottleneckArea struct {
	Area        string  `json:"area"`
	Impact      float64 `json:"impact"` // 0.0-1.0
	Description string  `json:"description"`
	Solutions   []string `json:"solutions"`
}

type OptimizationOption struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Effort       string           `json:"effort"`
	Impact       ImpactAssessment `json:"impact"`
	Feasibility  float64          `json:"feasibility"` // 0.0-1.0
	ROI          float64          `json:"roi"`         // Return on investment
}

type PerformanceGains struct {
	SpeedImprovement   float64 `json:"speed_improvement"`   // % improvement
	MemoryReduction    float64 `json:"memory_reduction"`    // % reduction
	ThroughputIncrease float64 `json:"throughput_increase"` // % increase
	LatencyReduction   float64 `json:"latency_reduction"`   // % reduction
}

type CodeLocation struct {
	File   string `json:"file,omitempty"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
	Function string `json:"function,omitempty"`
}

type ImpactAssessment struct {
	Performance   string `json:"performance"`   // high, medium, low, none
	Maintainability string `json:"maintainability"` // high, medium, low, none
	Reliability   string `json:"reliability"`   // high, medium, low, none
	Security      string `json:"security"`      // high, medium, low, none
}

type ResourceRequirements struct {
	DeveloperHours int      `json:"developer_hours"`
	TestingHours   int      `json:"testing_hours"`
	ReviewHours    int      `json:"review_hours"`
	Skills         []string `json:"skills"`
	Tools          []string `json:"tools"`
}

type ReadabilityMetrics struct {
	Score          float64 `json:"score"`
	CommentRatio   float64 `json:"comment_ratio"`
	NamingQuality  float64 `json:"naming_quality"`
	CodeClarity    float64 `json:"code_clarity"`
}

type ModularityMetrics struct {
	Score           float64 `json:"score"`
	CouplingLevel   string  `json:"coupling_level"`
	CohesionLevel   string  `json:"cohesion_level"`
	FunctionLength  float64 `json:"avg_function_length"`
}

type DocumentationMetrics struct {
	Score           float64 `json:"score"`
	APIDocCoverage  float64 `json:"api_doc_coverage"`
	InlineComments  float64 `json:"inline_comments"`
	ExamplesCoverage float64 `json:"examples_coverage"`
}

type TestCoverageMetrics struct {
	LineCoverage     float64 `json:"line_coverage"`
	BranchCoverage   float64 `json:"branch_coverage"`
	TestQuality      float64 `json:"test_quality"`
	TestMaintainability float64 `json:"test_maintainability"`
}

type CodeStructureMetrics struct {
	Complexity      int     `json:"complexity"`
	Duplication     float64 `json:"duplication"`
	CodeSmells      int     `json:"code_smells"`
	DesignPatterns  int     `json:"design_patterns"`
}

type DebtArea struct {
	Category    string  `json:"category"`
	Description string  `json:"description"`
	DebtLevel   float64 `json:"debt_level"`
	Impact      string  `json:"impact"`
	PaydownCost string  `json:"paydown_cost"`
}

type DebtPaydownPlan struct {
	TotalCost     time.Duration    `json:"total_cost"`
	Phases        []DebtPhase      `json:"phases"`
	ROI           float64          `json:"roi"`
	Priority      string           `json:"priority"`
}

type DebtPhase struct {
	Name        string        `json:"name"`
	Cost        time.Duration `json:"cost"`
	Benefits    []string      `json:"benefits"`
	Areas       []string      `json:"areas"`
}

type MaintenanceCost struct {
	CurrentCost   float64 `json:"current_cost"`   // relative cost
	OptimizedCost float64 `json:"optimized_cost"` // after debt paydown
	Savings       float64 `json:"savings"`        // % savings
}

type PerformanceBenchmark struct {
	ExecutionTime time.Duration `json:"execution_time"`
	MemoryUsage   int64         `json:"memory_usage"`
	Throughput    float64       `json:"throughput"`
	Allocations   int64         `json:"allocations"`
}

type PerformanceGap struct {
	SpeedGap    float64 `json:"speed_gap"`    // % slower than best
	MemoryGap   float64 `json:"memory_gap"`   // % more memory than best
	RankPosition int    `json:"rank_position"` // position in benchmark ranking
}

type StandardMetric struct {
	Name      string  `json:"name"`
	Standard  float64 `json:"standard"`
	Current   float64 `json:"current"`
	Unit      string  `json:"unit"`
	Meets     bool    `json:"meets"`
}

type ComplianceGap struct {
	Standard    string  `json:"standard"`
	Current     float64 `json:"current"`
	Required    float64 `json:"required"`
	Gap         float64 `json:"gap"`
	Severity    string  `json:"severity"`
	Remediation string  `json:"remediation"`
}

// BenchmarkSuite manages performance benchmarks
type BenchmarkSuite struct {
	Benchmarks map[string]PerformanceBenchmark
}

// PatternAnalyzer identifies code patterns and anti-patterns
type ReviewPatternAnalyzer struct {
	GoodPatterns []CodePattern
	AntiPatterns []CodePattern
}

type CodePattern struct {
	Name        string   `json:"name"`
	Pattern     string   `json:"pattern"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	Examples    []string `json:"examples"`
}

// RecommendationEngine generates improvement recommendations
type RecommendationEngine struct {
	Rules        []RecommendationRule
	Templates    map[string]string
	PriorityWeights map[string]float64
}

type RecommendationRule struct {
	Condition string
	Action    string
	Priority  int
	Category  string
}

// Performance Optimizers
type PerformanceOptimizer struct{}
type MemoryOptimizer struct{}
type CodeQualityOptimizer struct{}

type OptimizationResult struct {
	Category     string                 `json:"category"`
	Score        float64                `json:"score"`
	Issues       []ReviewIssue          `json:"issues"`
	Suggestions  []ImprovementSuggestion `json:"suggestions"`
	Potential    float64                `json:"potential"`
}

// NewReviewAgent creates a new review agent
func NewReviewAgent(database *Database) *ReviewAgent {
	agent := &ReviewAgent{
		database:   database,
		optimizers: make(map[string]Optimizer),
		benchmarks: &BenchmarkSuite{
			Benchmarks: make(map[string]PerformanceBenchmark),
		},
		patterns: &ReviewPatternAnalyzer{},
		recommendations: &RecommendationEngine{
			Templates: make(map[string]string),
			PriorityWeights: map[string]float64{
				"performance":     0.30,
				"maintainability": 0.25,
				"reliability":     0.20,
				"security":        0.15,
				"documentation":   0.10,
			},
		},
	}
	
	// Initialize optimizers
	agent.initializeOptimizers()
	
	// Initialize benchmarks
	agent.initializeBenchmarks()
	
	// Initialize patterns
	agent.initializePatterns()
	
	// Initialize recommendation engine
	agent.initializeRecommendationEngine()
	
	// Create database tables
	if err := agent.createReviewTables(); err != nil {
		log.Printf("Warning: Failed to create review tables: %v", err)
	}
	
	return agent
}

// ReviewCode performs comprehensive code review and optimization analysis
func (ra *ReviewAgent) ReviewCode(request *ReviewRequest) (*ReviewResult, error) {
	log.Printf("Review Agent: Starting review for %s (%s)", request.TaskID, request.BlockType)
	
	startTime := time.Now()
	
	result := &ReviewResult{
		TaskID:          request.TaskID,
		BlockType:       request.BlockType,
		CriticalIssues:  []ReviewIssue{},
		Improvements:    []ImprovementSuggestion{},
		BestPractices:   []BestPractice{},
		CreatedAt:       time.Now(),
	}
	
	// 1. Performance optimization analysis
	result.PerformanceOptimization = ra.analyzePerformanceOptimization(request)
	
	// 2. Code optimization analysis
	result.CodeOptimization = ra.analyzeCodeOptimization(request)
	
	// 3. Memory optimization analysis
	result.MemoryOptimization = ra.analyzeMemoryOptimization(request)
	
	// 4. Maintainability analysis
	result.Maintainability = ra.analyzeMaintainability(request)
	
	// 5. Technical debt analysis
	result.TechnicalDebt = ra.analyzeTechnicalDebt(request)
	
	// 6. Benchmark comparison
	result.BenchmarkComparison = ra.compareToBenchmarks(request)
	
	// 7. Industry standards comparison
	result.IndustryStandards = ra.compareToIndustryStandards(request)
	
	// 8. Identify critical issues
	result.CriticalIssues = ra.identifyCriticalIssues(request, result)
	
	// 9. Generate improvement suggestions
	result.Improvements = ra.generateImprovements(request, result)
	
	// 10. Evaluate best practices
	result.BestPractices = ra.evaluateBestPractices(request)
	
	// 11. Create optimization plan
	result.OptimizationPlan = ra.createOptimizationPlan(result)
	
	// 12. Calculate overall score
	result.OverallScore = ra.calculateOverallScore(result)
	
	// 13. Determine status
	result.Status = ra.determineReviewStatus(result)
	
	result.ReviewTime = time.Since(startTime)
	
	// Store result
	if err := ra.storeReviewResult(result); err != nil {
		log.Printf("Warning: Failed to store review result: %v", err)
	}
	
	log.Printf("Review Agent: Completed review for %s (Score: %.2f, Status: %s)", 
		request.TaskID, result.OverallScore, result.Status)
	
	return result, nil
}

// analyzePerformanceOptimization analyzes performance optimization opportunities
func (ra *ReviewAgent) analyzePerformanceOptimization(request *ReviewRequest) OptimizationAnalysis {
	optimizer := ra.optimizers["performance"]
	metrics := ra.extractPerformanceMetrics(request)
	
	optimizationResult := optimizer.Analyze(request.GeneratedCode, metrics)
	
	return OptimizationAnalysis{
		Category:            "performance",
		CurrentScore:        optimizationResult.Score,
		PotentialScore:      ra.calculatePotentialScore(optimizationResult),
		ImprovementPotential: optimizationResult.Potential,
		BottleneckAreas:     ra.identifyBottlenecks(request),
		OptimizationOptions: ra.generateOptimizationOptions(optimizationResult),
		EstimatedGains:      ra.estimatePerformanceGains(optimizationResult),
	}
}

// analyzeCodeOptimization analyzes code quality optimization
func (ra *ReviewAgent) analyzeCodeOptimization(request *ReviewRequest) OptimizationAnalysis {
	optimizer := ra.optimizers["code_quality"]
	metrics := ra.extractCodeMetrics(request)
	
	optimizationResult := optimizer.Analyze(request.GeneratedCode, metrics)
	
	return OptimizationAnalysis{
		Category:            "code_quality",
		CurrentScore:        optimizationResult.Score,
		PotentialScore:      ra.calculatePotentialScore(optimizationResult),
		ImprovementPotential: optimizationResult.Potential,
		BottleneckAreas:     ra.identifyCodeBottlenecks(request),
		OptimizationOptions: ra.generateCodeOptimizations(optimizationResult),
		EstimatedGains:      ra.estimateCodeGains(optimizationResult),
	}
}

// analyzeMemoryOptimization analyzes memory usage optimization
func (ra *ReviewAgent) analyzeMemoryOptimization(request *ReviewRequest) OptimizationAnalysis {
	optimizer := ra.optimizers["memory"]
	metrics := ra.extractMemoryMetrics(request)
	
	optimizationResult := optimizer.Analyze(request.GeneratedCode, metrics)
	
	return OptimizationAnalysis{
		Category:            "memory",
		CurrentScore:        optimizationResult.Score,
		PotentialScore:      ra.calculatePotentialScore(optimizationResult),
		ImprovementPotential: optimizationResult.Potential,
		BottleneckAreas:     ra.identifyMemoryBottlenecks(request),
		OptimizationOptions: ra.generateMemoryOptimizations(optimizationResult),
		EstimatedGains:      ra.estimateMemoryGains(optimizationResult),
	}
}

// analyzeMaintainability evaluates code maintainability
func (ra *ReviewAgent) analyzeMaintainability(request *ReviewRequest) MaintainabilityAnalysis {
	code := request.GeneratedCode
	
	return MaintainabilityAnalysis{
		OverallScore: ra.calculateMaintainabilityScore(code),
		Readability: ReadabilityMetrics{
			Score:         ra.calculateReadabilityScore(code),
			CommentRatio:  ra.calculateCommentRatio(code),
			NamingQuality: ra.evaluateNamingQuality(code),
			CodeClarity:   ra.evaluateCodeClarity(code),
		},
		Modularity: ModularityMetrics{
			Score:          ra.calculateModularityScore(code),
			CouplingLevel:  ra.assessCouplingLevel(code),
			CohesionLevel:  ra.assessCohesionLevel(code),
			FunctionLength: ra.calculateAvgFunctionLength(code),
		},
		Documentation: DocumentationMetrics{
			Score:            ra.calculateDocumentationScore(code),
			APIDocCoverage:   ra.calculateAPIDocCoverage(code),
			InlineComments:   ra.calculateInlineComments(code),
			ExamplesCoverage: ra.calculateExamplesCoverage(request),
		},
		TestCoverage: ra.extractTestCoverageMetrics(request),
		CodeStructure: CodeStructureMetrics{
			Complexity:     ra.calculateComplexity(code),
			Duplication:    ra.calculateDuplication(code),
			CodeSmells:     ra.countCodeSmells(code),
			DesignPatterns: ra.identifyDesignPatterns(code),
		},
	}
}

// analyzeTechnicalDebt quantifies technical debt
func (ra *ReviewAgent) analyzeTechnicalDebt(request *ReviewRequest) TechnicalDebtAnalysis {
	code := request.GeneratedCode
	
	debtAreas := ra.identifyDebtAreas(code)
	debtScore := ra.calculateDebtScore(debtAreas)
	
	return TechnicalDebtAnalysis{
		DebtLevel:   ra.determineDebtLevel(debtScore),
		DebtScore:   debtScore,
		DebtAreas:   debtAreas,
		PaydownPlan: ra.createDebtPaydownPlan(debtAreas),
		MaintenanceCost: ra.calculateMaintenanceCost(debtScore),
	}
}

// compareToBenchmarks compares performance to similar implementations
func (ra *ReviewAgent) compareToBenchmarks(request *ReviewRequest) BenchmarkComparison {
	category := request.BlockType
	currentPerf := ra.estimateCurrentPerformance(request)
	
	bestInCategory, exists := ra.benchmarks.Benchmarks[category+"_best"]
	if !exists {
		bestInCategory = ra.getDefaultBenchmark(category, "best")
	}
	
	avgPerf, exists := ra.benchmarks.Benchmarks[category+"_average"]
	if !exists {
		avgPerf = ra.getDefaultBenchmark(category, "average")
	}
	
	// Convert to BenchmarkComparison format
	categories := make(map[string]float64)
	categories["execution_time"] = float64(currentPerf.ExecutionTime.Milliseconds())
	categories["memory_usage"] = float64(currentPerf.MemoryUsage)
	categories["throughput"] = currentPerf.Throughput
	
	projectScore := (float64(currentPerf.ExecutionTime.Milliseconds()) + float64(currentPerf.MemoryUsage)/1000000.0) / 2
	benchmarkScore := (float64(bestInCategory.ExecutionTime.Milliseconds()) + float64(bestInCategory.MemoryUsage)/1000000.0) / 2
	industryAvg := (float64(avgPerf.ExecutionTime.Milliseconds()) + float64(avgPerf.MemoryUsage)/1000000.0) / 2
	
	percentile := ra.calculateRankingPercentile(currentPerf, avgPerf, bestInCategory)
	
	ranking := "Below Average"
	if percentile > 75 {
		ranking = "Excellent"
	} else if percentile > 50 {
		ranking = "Good"
	} else if percentile > 25 {
		ranking = "Average"
	}
	
	return BenchmarkComparison{
		ProjectScore:   projectScore,
		BenchmarkScore: benchmarkScore,
		IndustryAvg:    industryAvg,
		Percentile:     percentile,
		Categories:     categories,
		Ranking:        ranking,
	}
}

// compareToIndustryStandards compares against industry standards
func (ra *ReviewAgent) compareToIndustryStandards(request *ReviewRequest) IndustryComparison {
	standards := ra.getIndustryStandards(request.BlockType)
	compliance := make(map[string]bool)
	gaps := []ComplianceGap{}
	
	complianceCount := 0
	totalStandards := len(standards)
	
	for name, standard := range standards {
		current := ra.measureStandardMetric(request, name)
		meets := ra.meetsStandard(current, standard)
		compliance[name] = meets
		
		if meets {
			complianceCount++
		} else {
			gaps = append(gaps, ComplianceGap{
				Standard:    name,
				Current:     current,
				Required:    standard.Standard,
				Gap:         standard.Standard - current,
				Severity:    ra.determineGapSeverity(standard.Standard - current),
				Remediation: ra.getRemediationAdvice(name),
			})
		}
	}
	
	complianceScore := float64(complianceCount) / float64(totalStandards)
	
	return IndustryComparison{
		Standards:       standards,
		Compliance:      compliance,
		ComplianceScore: complianceScore,
		Gaps:           gaps,
	}
}

// identifyCriticalIssues finds critical issues requiring immediate attention
func (ra *ReviewAgent) identifyCriticalIssues(request *ReviewRequest, result *ReviewResult) []ReviewIssue {
	issues := []ReviewIssue{}
	
	// Performance critical issues
	if result.PerformanceOptimization.CurrentScore < 0.4 {
		issues = append(issues, ReviewIssue{
			ID:          "PERF001",
			Type:        "performance",
			Severity:    "critical",
			Title:       "Poor Performance Detected",
			Description: "Code shows significant performance issues that may impact user experience",
			Impact:      "High latency and resource consumption",
			Effort:      "medium",
			Solution:    "Implement performance optimizations identified in analysis",
		})
	}
	
	// Code quality critical issues
	if result.Maintainability.OverallScore < 0.3 {
		issues = append(issues, ReviewIssue{
			ID:          "QUAL001",
			Type:        "maintainability",
			Severity:    "high",
			Title:       "Low Code Maintainability",
			Description: "Code structure makes maintenance and future changes difficult",
			Impact:      "Increased development time and higher bug risk",
			Effort:      "high",
			Solution:    "Refactor code following maintainability recommendations",
		})
	}
	
	// Technical debt critical issues
	if result.TechnicalDebt.DebtLevel == "critical" || result.TechnicalDebt.DebtLevel == "high" {
		issues = append(issues, ReviewIssue{
			ID:          "DEBT001",
			Type:        "technical_debt",
			Severity:    "high",
			Title:       "High Technical Debt",
			Description: "Accumulated technical debt may impact long-term development velocity",
			Impact:      "Slower feature development and increased bug risk",
			Effort:      "high",
			Solution:    "Follow debt paydown plan to reduce technical debt",
		})
	}
	
	return issues
}

// generateImprovements creates actionable improvement suggestions
func (ra *ReviewAgent) generateImprovements(request *ReviewRequest, result *ReviewResult) []ImprovementSuggestion {
	suggestions := []ImprovementSuggestion{}
	
	// Performance improvements
	if result.PerformanceOptimization.ImprovementPotential > 0.3 {
		suggestions = append(suggestions, ImprovementSuggestion{
			ID:          "IMP001",
			Category:    "performance",
			Title:       "Optimize Parsing Algorithm",
			Description: "Current parsing implementation has optimization opportunities",
			Benefits:    []string{"Faster processing", "Lower CPU usage", "Better scalability"},
			Implementation: "Replace linear search with binary search for block type identification",
			Priority:    1,
			EstimatedEffort: "medium",
			ExpectedImpact: ImpactAssessment{
				Performance:     "high",
				Maintainability: "medium",
				Reliability:     "low",
				Security:        "none",
			},
		})
	}
	
	// Code quality improvements
	if result.Maintainability.Readability.Score < 0.6 {
		suggestions = append(suggestions, ImprovementSuggestion{
			ID:          "IMP002",
			Category:    "code_quality",
			Title:       "Improve Code Readability",
			Description: "Code readability can be enhanced with better naming and comments",
			Benefits:    []string{"Easier maintenance", "Faster onboarding", "Reduced bugs"},
			Implementation: "Add meaningful variable names and comprehensive comments",
			CodeExample: ra.getReadabilityExample(),
			Priority:    2,
			EstimatedEffort: "low",
			ExpectedImpact: ImpactAssessment{
				Performance:     "none",
				Maintainability: "high",
				Reliability:     "medium",
				Security:        "none",
			},
		})
	}
	
	// Memory improvements
	if result.MemoryOptimization.ImprovementPotential > 0.2 {
		suggestions = append(suggestions, ImprovementSuggestion{
			ID:          "IMP003",
			Category:    "memory",
			Title:       "Optimize Memory Usage",
			Description: "Memory allocations can be reduced through pooling and reuse",
			Benefits:    []string{"Lower memory footprint", "Reduced GC pressure", "Better performance"},
			Implementation: "Implement object pooling for frequently allocated structures",
			Priority:    2,
			EstimatedEffort: "medium",
			ExpectedImpact: ImpactAssessment{
				Performance:     "medium",
				Maintainability: "low",
				Reliability:     "medium",
				Security:        "none",
			},
		})
	}
	
	return ra.prioritizeSuggestions(suggestions)
}

// evaluateBestPractices checks adherence to coding best practices
func (ra *ReviewAgent) evaluateBestPractices(request *ReviewRequest) []BestPractice {
	practices := []BestPractice{}
	code := request.GeneratedCode
	
	practices = append(practices, BestPractice{
		Category:   "error_handling",
		Practice:   "Always check and handle errors appropriately",
		Rationale:  "Proper error handling prevents crashes and provides better user experience",
		IsFollowed: strings.Contains(code, "error") && strings.Contains(code, "return"),
		Importance: "critical",
	})
	
	practices = append(practices, BestPractice{
		Category:   "input_validation",
		Practice:   "Validate all input parameters",
		Rationale:  "Input validation prevents security vulnerabilities and runtime errors",
		IsFollowed: strings.Contains(code, "len(") || strings.Contains(code, "nil"),
		Importance: "critical",
	})
	
	practices = append(practices, BestPractice{
		Category:   "documentation",
		Practice:   "Document public functions and complex logic",
		Rationale:  "Good documentation improves code maintainability and team productivity",
		IsFollowed: strings.Contains(code, "//"),
		Importance: "important",
	})
	
	practices = append(practices, BestPractice{
		Category:   "testing",
		Practice:   "Provide comprehensive test coverage",
		Rationale:  "Tests ensure code correctness and prevent regressions",
		IsFollowed: request.TestCode != "",
		Importance: "critical",
	})
	
	return practices
}

// createOptimizationPlan generates structured optimization roadmap
func (ra *ReviewAgent) createOptimizationPlan(result *ReviewResult) OptimizationPlan {
	phases := []OptimizationPhase{}
	
	// Phase 1: Critical Issues
	if len(result.CriticalIssues) > 0 {
		tasks := ra.createTasksForIssues(result.CriticalIssues)
		phases = append(phases, OptimizationPhase{
			Phase:       1,
			Name:        "Critical Issues Resolution",
			Description: "Address critical issues that impact functionality or performance",
			Tasks:       tasks,
			Duration:    ra.estimatePhaseDuration(tasks),
			Priority:    "critical",
			Dependencies: []int{},
		})
	}
	
	// Phase 2: Performance Optimization
	if result.PerformanceOptimization.ImprovementPotential > 0.3 {
		tasks := ra.createPerformanceTasks(result.PerformanceOptimization)
		phases = append(phases, OptimizationPhase{
			Phase:       2,
			Name:        "Performance Optimization",
			Description: "Implement performance improvements for better speed and efficiency",
			Tasks:       tasks,
			Duration:    ra.estimatePhaseDuration(tasks),
			Priority:    "high",
			Dependencies: []int{1},
		})
	}
	
	// Phase 3: Code Quality Improvements
	if result.Maintainability.OverallScore < 0.8 {
		tasks := ra.createQualityTasks(result.Maintainability)
		phases = append(phases, OptimizationPhase{
			Phase:       3,
			Name:        "Code Quality Enhancement",
			Description: "Improve code maintainability and readability",
			Tasks:       tasks,
			Duration:    ra.estimatePhaseDuration(tasks),
			Priority:    "medium",
			Dependencies: []int{1},
		})
	}
	
	// Calculate total estimates
	totalDuration := time.Duration(0)
	for _, phase := range phases {
		totalDuration += phase.Duration
	}
	
	return OptimizationPlan{
		Phases:        phases,
		TotalEstimate: totalDuration,
		ExpectedGains: ra.calculateTotalExpectedGains(result),
		ResourceRequirements: ResourceRequirements{
			DeveloperHours: int(totalDuration / time.Hour),
			TestingHours:   int(totalDuration / time.Hour / 3), // 1/3 of dev time for testing
			ReviewHours:    int(totalDuration / time.Hour / 6), // 1/6 of dev time for review
			Skills:         []string{"Go programming", "Performance optimization", "Code review"},
			Tools:          []string{"Go profiler", "Benchmarking tools", "Code analysis tools"},
		},
	}
}

// calculateOverallScore computes overall review score
func (ra *ReviewAgent) calculateOverallScore(result *ReviewResult) float64 {
	weights := ra.recommendations.PriorityWeights
	
	scores := map[string]float64{
		"performance":     result.PerformanceOptimization.CurrentScore,
		"maintainability": result.Maintainability.OverallScore,
		"reliability":     ra.calculateReliabilityScore(result),
		"security":        ra.calculateSecurityScore(result),
		"documentation":   result.Maintainability.Documentation.Score,
	}
	
	weightedSum := 0.0
	totalWeight := 0.0
	
	for category, weight := range weights {
		if score, exists := scores[category]; exists {
			weightedSum += score * weight
			totalWeight += weight
		}
	}
	
	if totalWeight == 0 {
		return 0.0
	}
	
	return weightedSum / totalWeight
}

// determineReviewStatus determines final review status
func (ra *ReviewAgent) determineReviewStatus(result *ReviewResult) string {
	if result.OverallScore >= 0.9 {
		return "excellent"
	} else if result.OverallScore >= 0.7 {
		return "good"
	} else if result.OverallScore >= 0.5 {
		return "needs_improvement"
	} else {
		return "poor"
	}
}

// Initialize methods

func (ra *ReviewAgent) initializeOptimizers() {
	ra.optimizers["performance"] = &PerformanceOptimizer{}
	ra.optimizers["memory"] = &MemoryOptimizer{}
	ra.optimizers["code_quality"] = &CodeQualityOptimizer{}
}

func (ra *ReviewAgent) initializeBenchmarks() {
	// Initialize with realistic benchmarks for AEP parsing
	ra.benchmarks.Benchmarks["Item_best"] = PerformanceBenchmark{
		ExecutionTime: time.Microsecond * 50,
		MemoryUsage:   1024,
		Throughput:    1000.0,
		Allocations:   5,
	}
	
	ra.benchmarks.Benchmarks["Item_average"] = PerformanceBenchmark{
		ExecutionTime: time.Microsecond * 100,
		MemoryUsage:   2048,
		Throughput:    500.0,
		Allocations:   10,
	}
	
	// Add more benchmarks for other block types
	ra.addBlockTypeBenchmarks("Layer", time.Microsecond*75, 1536, 750.0, 7)
	ra.addBlockTypeBenchmarks("Property", time.Microsecond*25, 512, 2000.0, 3)
}

func (ra *ReviewAgent) initializePatterns() {
	ra.patterns.GoodPatterns = []CodePattern{
		{
			Name:        "Error Handling",
			Pattern:     "if.*error.*return",
			Category:    "reliability",
			Description: "Proper error checking and propagation",
			Impact:      "high",
		},
		{
			Name:        "Input Validation",
			Pattern:     "len\\(.*\\).*<",
			Category:    "security",
			Description: "Input bounds checking",
			Impact:      "high",
		},
	}
	
	ra.patterns.AntiPatterns = []CodePattern{
		{
			Name:        "Panic Usage",
			Pattern:     "panic\\(",
			Category:    "reliability",
			Description: "Using panic instead of proper error handling",
			Impact:      "critical",
		},
		{
			Name:        "Magic Numbers",
			Pattern:     "[0-9]{2,}",
			Category:    "maintainability",
			Description: "Hard-coded numeric values without explanation",
			Impact:      "medium",
		},
	}
}

func (ra *ReviewAgent) initializeRecommendationEngine() {
	ra.recommendations.Rules = []RecommendationRule{
		{
			Condition: "performance_score < 0.5",
			Action:    "recommend_performance_optimization",
			Priority:  1,
			Category:  "performance",
		},
		{
			Condition: "maintainability_score < 0.6",
			Action:    "recommend_refactoring",
			Priority:  2,
			Category:  "maintainability",
		},
	}
	
	ra.recommendations.Templates["readability_example"] = `
// Before: Poor readability
func ParseItem(d []byte, o int) (*Item, int, error) {
	if len(d) < o+8 {
		return nil, 0, errors.New("bad")
	}
	// ... more code
}

// After: Better readability
func ParseItem(data []byte, offset int) (*Item, int, error) {
	if len(data) < offset+MinimumItemSize {
		return nil, 0, ErrInsufficientData
	}
	// ... more code with clear variable names
}
`
}

// Helper methods implementation (abbreviated for brevity)

func (ra *ReviewAgent) extractPerformanceMetrics(request *ReviewRequest) map[string]interface{} {
	metrics := make(map[string]interface{})
	if request.VerificationResult != nil {
		metrics["execution_time"] = request.VerificationResult.Performance.AvgExecutionTime
		metrics["memory_usage"] = request.VerificationResult.Performance.MemoryUsage
	}
	return metrics
}

func (ra *ReviewAgent) calculateComplexity(code string) int {
	// Simplified complexity calculation
	complexity := 1
	keywords := []string{"if", "for", "switch", "case", "else"}
	for _, keyword := range keywords {
		complexity += strings.Count(code, keyword)
	}
	return complexity
}

func (ra *ReviewAgent) getReadabilityExample() string {
	if example, exists := ra.recommendations.Templates["readability_example"]; exists {
		return example
	}
	return "// Example code improvements here"
}

func (ra *ReviewAgent) prioritizeSuggestions(suggestions []ImprovementSuggestion) []ImprovementSuggestion {
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Priority < suggestions[j].Priority
	})
	return suggestions
}

// Database operations (simplified)
func (ra *ReviewAgent) createReviewTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS review_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		block_type TEXT NOT NULL,
		status TEXT NOT NULL,
		overall_score REAL NOT NULL,
		result_data TEXT NOT NULL,
		review_time_ms INTEGER NOT NULL,
		created_at INTEGER NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_review_task_id ON review_results(task_id);
	`
	
	_, err := ra.database.db.Exec(query)
	return err
}

func (ra *ReviewAgent) storeReviewResult(result *ReviewResult) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	
	query := `
		INSERT INTO review_results
		(task_id, block_type, status, overall_score, result_data, review_time_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = ra.database.db.Exec(query,
		result.TaskID, result.BlockType, result.Status, result.OverallScore,
		string(resultJSON), int64(result.ReviewTime/time.Millisecond), result.CreatedAt.Unix())
	
	return err
}

// Optimizer implementations
func (po *PerformanceOptimizer) Analyze(code string, metrics map[string]interface{}) OptimizationResult {
	// Simplified performance analysis
	score := 0.7 // Default score
	
	if strings.Contains(code, "binary.BigEndian") {
		score += 0.1 // Good: using efficient binary operations
	}
	
	if strings.Count(code, "append(") > 3 {
		score -= 0.2 // Potential: multiple appends could be optimized
	}
	
	return OptimizationResult{
		Category:  "performance",
		Score:     math.Min(score, 1.0),
		Potential: 0.8 - score, // Room for improvement
	}
}

func (po *PerformanceOptimizer) GetCategory() string { return "performance" }
func (po *PerformanceOptimizer) GetPriority() int    { return 1 }

func (mo *MemoryOptimizer) Analyze(code string, metrics map[string]interface{}) OptimizationResult {
	score := 0.6
	
	if strings.Contains(code, "make([]") {
		score += 0.2 // Good: pre-allocating slices
	}
	
	return OptimizationResult{
		Category:  "memory",
		Score:     math.Min(score, 1.0),
		Potential: 0.9 - score,
	}
}

func (mo *MemoryOptimizer) GetCategory() string { return "memory" }
func (mo *MemoryOptimizer) GetPriority() int    { return 2 }

func (cqo *CodeQualityOptimizer) Analyze(code string, metrics map[string]interface{}) OptimizationResult {
	score := 0.5
	
	if strings.Contains(code, "//") {
		score += 0.2 // Good: has comments
	}
	
	if len(code) < 1000 {
		score += 0.2 // Good: reasonable function size
	}
	
	return OptimizationResult{
		Category:  "code_quality",
		Score:     math.Min(score, 1.0),
		Potential: 0.95 - score,
	}
}

func (cqo *CodeQualityOptimizer) GetCategory() string { return "code_quality" }
func (cqo *CodeQualityOptimizer) GetPriority() int    { return 3 }

// Additional helper methods (implementation abbreviated for brevity)
func (ra *ReviewAgent) calculateMaintainabilityScore(code string) float64 { return 0.7 }
func (ra *ReviewAgent) calculateReadabilityScore(code string) float64     { return 0.75 }
func (ra *ReviewAgent) calculateCommentRatio(code string) float64         { return 0.15 }
func (ra *ReviewAgent) evaluateNamingQuality(code string) float64         { return 0.8 }
func (ra *ReviewAgent) evaluateCodeClarity(code string) float64           { return 0.7 }
func (ra *ReviewAgent) calculateModularityScore(code string) float64      { return 0.8 }
func (ra *ReviewAgent) assessCouplingLevel(code string) string            { return "low" }
func (ra *ReviewAgent) assessCohesionLevel(code string) string            { return "high" }
func (ra *ReviewAgent) calculateAvgFunctionLength(code string) float64    { return 25.0 }
func (ra *ReviewAgent) calculateDocumentationScore(code string) float64   { return 0.6 }
func (ra *ReviewAgent) calculateAPIDocCoverage(code string) float64       { return 0.7 }
func (ra *ReviewAgent) calculateInlineComments(code string) float64       { return 0.2 }
func (ra *ReviewAgent) calculateExamplesCoverage(request *ReviewRequest) float64 { return 0.3 }
func (ra *ReviewAgent) calculateDuplication(code string) float64          { return 0.05 }
func (ra *ReviewAgent) countCodeSmells(code string) int                   { return 2 }
func (ra *ReviewAgent) identifyDesignPatterns(code string) int            { return 1 }

// More helper method stubs...
func (ra *ReviewAgent) extractCodeMetrics(request *ReviewRequest) map[string]interface{} {
	return make(map[string]interface{})
}
func (ra *ReviewAgent) extractMemoryMetrics(request *ReviewRequest) map[string]interface{} {
	return make(map[string]interface{})
}
func (ra *ReviewAgent) extractTestCoverageMetrics(request *ReviewRequest) TestCoverageMetrics {
	return TestCoverageMetrics{LineCoverage: 0.75, BranchCoverage: 0.7, TestQuality: 0.8, TestMaintainability: 0.75}
}
func (ra *ReviewAgent) calculatePotentialScore(result OptimizationResult) float64 {
	return math.Min(result.Score+result.Potential, 1.0)
}
func (ra *ReviewAgent) identifyBottlenecks(request *ReviewRequest) []BottleneckArea {
	return []BottleneckArea{{Area: "Binary parsing", Impact: 0.6, Description: "Main processing bottleneck"}}
}
func (ra *ReviewAgent) generateOptimizationOptions(result OptimizationResult) []OptimizationOption {
	return []OptimizationOption{{Name: "Optimize binary parsing", Effort: "medium", Feasibility: 0.8, ROI: 2.5}}
}
func (ra *ReviewAgent) estimatePerformanceGains(result OptimizationResult) PerformanceGains {
	return PerformanceGains{SpeedImprovement: 25.0, MemoryReduction: 15.0, ThroughputIncrease: 30.0, LatencyReduction: 20.0}
}
func (ra *ReviewAgent) identifyCodeBottlenecks(request *ReviewRequest) []BottleneckArea {
	return []BottleneckArea{{Area: "Code structure", Impact: 0.4, Description: "Complex function organization"}}
}
func (ra *ReviewAgent) generateCodeOptimizations(result OptimizationResult) []OptimizationOption {
	return []OptimizationOption{{Name: "Refactor functions", Effort: "low", Feasibility: 0.9, ROI: 1.8}}
}
func (ra *ReviewAgent) estimateCodeGains(result OptimizationResult) PerformanceGains {
	return PerformanceGains{SpeedImprovement: 5.0, MemoryReduction: 0.0, ThroughputIncrease: 10.0, LatencyReduction: 0.0}
}
func (ra *ReviewAgent) identifyMemoryBottlenecks(request *ReviewRequest) []BottleneckArea {
	return []BottleneckArea{{Area: "Memory allocation", Impact: 0.3, Description: "Frequent allocations"}}
}
func (ra *ReviewAgent) generateMemoryOptimizations(result OptimizationResult) []OptimizationOption {
	return []OptimizationOption{{Name: "Implement object pooling", Effort: "medium", Feasibility: 0.7, ROI: 2.0}}
}
func (ra *ReviewAgent) estimateMemoryGains(result OptimizationResult) PerformanceGains {
	return PerformanceGains{SpeedImprovement: 10.0, MemoryReduction: 40.0, ThroughputIncrease: 15.0, LatencyReduction: 5.0}
}

// More implementation stubs for database and analysis methods...
func (ra *ReviewAgent) identifyDebtAreas(code string) []DebtArea {
	return []DebtArea{{Category: "Structure", Description: "Function organization", DebtLevel: 0.3, Impact: "medium"}}
}
func (ra *ReviewAgent) calculateDebtScore(areas []DebtArea) float64 { return 0.3 }
func (ra *ReviewAgent) determineDebtLevel(score float64) string {
	if score > 0.7 {
		return "critical"
	} else if score > 0.5 {
		return "high"
	} else if score > 0.3 {
		return "medium"
	}
	return "low"
}
func (ra *ReviewAgent) createDebtPaydownPlan(areas []DebtArea) DebtPaydownPlan {
	return DebtPaydownPlan{TotalCost: time.Hour * 4, ROI: 1.8, Priority: "medium"}
}
func (ra *ReviewAgent) calculateMaintenanceCost(debtScore float64) MaintenanceCost {
	return MaintenanceCost{CurrentCost: 1.0, OptimizedCost: 0.7, Savings: 30.0}
}
func (ra *ReviewAgent) estimateCurrentPerformance(request *ReviewRequest) PerformanceBenchmark {
	return PerformanceBenchmark{ExecutionTime: time.Microsecond * 80, MemoryUsage: 1800, Throughput: 625.0, Allocations: 8}
}
func (ra *ReviewAgent) getDefaultBenchmark(category, level string) PerformanceBenchmark {
	if level == "best" {
		return PerformanceBenchmark{ExecutionTime: time.Microsecond * 50, MemoryUsage: 1000, Throughput: 1000.0, Allocations: 5}
	}
	return PerformanceBenchmark{ExecutionTime: time.Microsecond * 100, MemoryUsage: 2000, Throughput: 500.0, Allocations: 10}
}
func (ra *ReviewAgent) calculatePerformanceGap(current, best PerformanceBenchmark) PerformanceGap {
	speedGap := float64(current.ExecutionTime-best.ExecutionTime) / float64(best.ExecutionTime) * 100
	memoryGap := float64(current.MemoryUsage-best.MemoryUsage) / float64(best.MemoryUsage) * 100
	return PerformanceGap{SpeedGap: speedGap, MemoryGap: memoryGap, RankPosition: 3}
}
func (ra *ReviewAgent) calculateRankingPercentile(current, avg, best PerformanceBenchmark) float64 { return 65.0 }
func (ra *ReviewAgent) getIndustryStandards(blockType string) map[string]StandardMetric {
	return map[string]StandardMetric{
		"max_complexity": {Name: "Cyclomatic Complexity", Standard: 10, Unit: "points", Current: 8, Meets: true},
		"min_coverage":   {Name: "Test Coverage", Standard: 80, Unit: "percent", Current: 75, Meets: false},
	}
}
func (ra *ReviewAgent) measureStandardMetric(request *ReviewRequest, name string) float64 {
	switch name {
	case "max_complexity":
		return float64(ra.calculateComplexity(request.GeneratedCode))
	case "min_coverage":
		if request.VerificationResult != nil {
			return request.VerificationResult.Coverage.LineCoverage * 100
		}
		return 75.0
	}
	return 0.0
}
func (ra *ReviewAgent) meetsStandard(current float64, standard StandardMetric) bool {
	return current >= standard.Standard
}
func (ra *ReviewAgent) determineGapSeverity(gap float64) string {
	if gap > 20 {
		return "critical"
	} else if gap > 10 {
		return "high"
	} else if gap > 5 {
		return "medium"
	}
	return "low"
}
func (ra *ReviewAgent) getRemediationAdvice(standard string) string {
	advice := map[string]string{
		"max_complexity": "Refactor complex functions into smaller, focused functions",
		"min_coverage":   "Add unit tests to increase test coverage above 80%",
	}
	if advice, exists := advice[standard]; exists {
		return advice
	}
	return "Follow industry best practices"
}
func (ra *ReviewAgent) calculateReliabilityScore(result *ReviewResult) float64 { return 0.8 }
func (ra *ReviewAgent) calculateSecurityScore(result *ReviewResult) float64   { return 0.85 }
func (ra *ReviewAgent) addBlockTypeBenchmarks(blockType string, execTime time.Duration, memory int64, throughput float64, allocs int64) {
	ra.benchmarks.Benchmarks[blockType+"_best"] = PerformanceBenchmark{
		ExecutionTime: execTime,
		MemoryUsage:   memory,
		Throughput:    throughput,
		Allocations:   allocs,
	}
	ra.benchmarks.Benchmarks[blockType+"_average"] = PerformanceBenchmark{
		ExecutionTime: execTime * 2,
		MemoryUsage:   memory * 2,
		Throughput:    throughput / 2,
		Allocations:   allocs * 2,
	}
}
func (ra *ReviewAgent) createTasksForIssues(issues []ReviewIssue) []OptimizationTask { return []OptimizationTask{} }
func (ra *ReviewAgent) createPerformanceTasks(optimization OptimizationAnalysis) []OptimizationTask {
	return []OptimizationTask{}
}
func (ra *ReviewAgent) createQualityTasks(maintainability MaintainabilityAnalysis) []OptimizationTask {
	return []OptimizationTask{}
}
func (ra *ReviewAgent) estimatePhaseDuration(tasks []OptimizationTask) time.Duration { return time.Hour * 2 }
func (ra *ReviewAgent) calculateTotalExpectedGains(result *ReviewResult) PerformanceGains {
	return PerformanceGains{SpeedImprovement: 30.0, MemoryReduction: 25.0, ThroughputIncrease: 35.0, LatencyReduction: 20.0}
}

// Public API methods

// GetReviewByTaskID retrieves review result by task ID
func (ra *ReviewAgent) GetReviewByTaskID(taskID string) (*ReviewResult, error) {
	query := `SELECT result_data FROM review_results WHERE task_id = ? ORDER BY created_at DESC LIMIT 1`
	
	var resultJSON string
	err := ra.database.db.QueryRow(query, taskID).Scan(&resultJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review result: %w", err)
	}
	
	var result ReviewResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}
	
	return &result, nil
}