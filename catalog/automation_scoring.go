// Package catalog provides automation scoring and asset relationship mapping
package catalog

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

// AutomationScorer provides comprehensive scoring for templates
type AutomationScorer struct {
	weights ScoreWeights
}

// ScoreWeights defines the importance of different factors in scoring
type ScoreWeights struct {
	TextComplexity     float64 `json:"text_complexity"`      // Weight for text automation potential
	MediaComplexity    float64 `json:"media_complexity"`     // Weight for media replacement potential
	ModularScore       float64 `json:"modular_score"`        // Weight for modular system benefits
	EffectComplexity   float64 `json:"effect_complexity"`    // Weight for effect customization
	DataBindingScore   float64 `json:"data_binding_score"`   // Weight for data-driven potential
	APIReadiness       float64 `json:"api_readiness"`        // Weight for API integration readiness
	MaintenanceScore   float64 `json:"maintenance_score"`    // Weight for maintenance complexity
}

// DefaultScoreWeights provides balanced scoring weights
func DefaultScoreWeights() ScoreWeights {
	return ScoreWeights{
		TextComplexity:     0.25,
		MediaComplexity:    0.20,
		ModularScore:       0.15,
		EffectComplexity:   0.10,
		DataBindingScore:   0.15,
		APIReadiness:       0.10,
		MaintenanceScore:   0.05,
	}
}

// AssetRelationshipMapper analyzes relationships between assets
type AssetRelationshipMapper struct {
	scorer *AutomationScorer
}

// AssetRelationship represents a relationship between assets
type AssetRelationship struct {
	SourceAssetID   string                 `json:"source_asset_id"`
	TargetAssetID   string                 `json:"target_asset_id"`
	RelationshipType string                `json:"relationship_type"` // replacement, dependency, group
	Strength        float64                `json:"strength"`          // 0-1 scale
	Properties      map[string]interface{} `json:"properties"`
	Bidirectional   bool                   `json:"bidirectional"`
}

// AssetGroup represents a logical grouping of assets
type AssetGroup struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`        // theme, brand, content
	AssetIDs    []string `json:"asset_ids"`
	Replaceable bool     `json:"replaceable"` // Can entire group be replaced as unit
	Priority    int      `json:"priority"`    // 1-10 importance scale
}

// AutomationOpportunity represents a scored automation opportunity
type AutomationOpportunity struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	AutomationScore  float64                `json:"automation_score"`  // 0-100
	ComplexityScore  float64                `json:"complexity_score"`  // 0-100  
	ROIScore         float64                `json:"roi_score"`         // 0-100
	ImplementationEffort string             `json:"implementation_effort"`
	BusinessImpact   string                 `json:"business_impact"`
	TechnicalDetails map[string]interface{} `json:"technical_details"`
	Prerequisites    []string               `json:"prerequisites"`
	Recommendations  []string               `json:"recommendations"`
}

// ScoringResult contains comprehensive scoring results
type ScoringResult struct {
	ProjectID           int64                    `json:"project_id"`
	OverallScore        float64                  `json:"overall_score"`
	CategoryScores      map[string]float64       `json:"category_scores"`
	AssetRelationships  []AssetRelationship      `json:"asset_relationships"`
	AssetGroups         []AssetGroup             `json:"asset_groups"`
	Opportunities       []AutomationOpportunity  `json:"opportunities"`
	ScoreBreakdown      ScoreBreakdown          `json:"score_breakdown"`
	Recommendations     []string                 `json:"recommendations"`
}

// ScoreBreakdown provides detailed score explanation
type ScoreBreakdown struct {
	TextScore        float64 `json:"text_score"`
	MediaScore       float64 `json:"media_score"`
	ModularScore     float64 `json:"modular_score"`
	EffectScore      float64 `json:"effect_score"`
	DataBindingScore float64 `json:"data_binding_score"`
	APIScore         float64 `json:"api_score"`
	MaintenanceScore float64 `json:"maintenance_score"`
	WeightedTotal    float64 `json:"weighted_total"`
	OverallScore     float64 `json:"overall_score"`
}

// NewAutomationScorer creates a new automation scorer
func NewAutomationScorer() *AutomationScorer {
	return &AutomationScorer{
		weights: DefaultScoreWeights(),
	}
}

// NewAssetRelationshipMapper creates a new asset relationship mapper
func NewAssetRelationshipMapper() *AssetRelationshipMapper {
	return &AssetRelationshipMapper{
		scorer: NewAutomationScorer(),
	}
}

// ScoreAutomationPotential calculates comprehensive automation scores
func (as *AutomationScorer) ScoreAutomationPotential(metadata *ProjectMetadata, analysis *DeepAnalysisResult) *ScoringResult {
	result := &ScoringResult{
		CategoryScores:     make(map[string]float64),
		AssetRelationships: []AssetRelationship{},
		AssetGroups:        []AssetGroup{},
		Opportunities:      []AutomationOpportunity{},
	}

	// Calculate individual scores
	breakdown := ScoreBreakdown{}
	breakdown.TextScore = as.scoreTextAutomation(metadata, analysis)
	breakdown.MediaScore = as.scoreMediaAutomation(metadata, analysis)
	breakdown.ModularScore = as.scoreModularSystem(metadata, analysis)
	breakdown.EffectScore = as.scoreEffectAutomation(metadata, analysis)
	breakdown.DataBindingScore = as.scoreDataBinding(metadata, analysis)
	breakdown.APIScore = as.scoreAPIReadiness(metadata, analysis)
	breakdown.MaintenanceScore = as.scoreMaintenanceComplexity(metadata, analysis)

	// Calculate weighted total
	weights := as.weights
	breakdown.WeightedTotal = 
		breakdown.TextScore*weights.TextComplexity +
		breakdown.MediaScore*weights.MediaComplexity +
		breakdown.ModularScore*weights.ModularScore +
		breakdown.EffectScore*weights.EffectComplexity +
		breakdown.DataBindingScore*weights.DataBindingScore +
		breakdown.APIScore*weights.APIReadiness +
		breakdown.MaintenanceScore*weights.MaintenanceScore

	result.OverallScore = breakdown.WeightedTotal
	result.ScoreBreakdown = breakdown

	// Create category scores
	result.CategoryScores["text_automation"] = breakdown.TextScore
	result.CategoryScores["media_automation"] = breakdown.MediaScore
	result.CategoryScores["modular_system"] = breakdown.ModularScore
	result.CategoryScores["effect_automation"] = breakdown.EffectScore
	result.CategoryScores["api_integration"] = breakdown.APIScore

	// Generate opportunities
	result.Opportunities = as.generateOpportunities(metadata, analysis, result.CategoryScores)

	// Generate recommendations
	result.Recommendations = as.generateRecommendations(result.CategoryScores, breakdown)

	return result
}

// scoreTextAutomation scores text automation potential
func (as *AutomationScorer) scoreTextAutomation(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	if len(metadata.TextLayers) == 0 {
		return 0.0
	}

	score := 0.0
	maxScore := 100.0

	// Base score from number of text layers
	layerScore := math.Min(float64(len(metadata.TextLayers))*5, 40) // Cap at 40 points
	score += layerScore

	// Intelligence bonus from deep analysis
	if analysis != nil && analysis.TextIntelligence != nil {
		intel := analysis.TextIntelligence
		
		// Bonus for detected field types
		typeBonus := math.Min(float64(len(intel.PatternGroups))*8, 30) // Cap at 30 points
		score += typeBonus

		// Bonus for data binding potential
		if len(intel.DataBindingOptions) > 0 {
			score += 15
		}

		// Bonus for localization readiness
		if intel.LocalizationReady {
			score += 15
		}
	}

	return math.Min(score, maxScore)
}

// scoreMediaAutomation scores media replacement automation potential
func (as *AutomationScorer) scoreMediaAutomation(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	if len(metadata.MediaAssets) == 0 {
		return 0.0
	}

	score := 0.0
	maxScore := 100.0
	replaceableCount := 0

	// Count replaceable assets
	for _, asset := range metadata.MediaAssets {
		if asset.IsPlaceholder {
			replaceableCount++
		}
	}

	if replaceableCount == 0 {
		return 20.0 // Some base score for having media assets
	}

	// Base score from replaceable assets
	replaceableScore := math.Min(float64(replaceableCount)*10, 50) // Cap at 50 points
	score += replaceableScore

	// Enhanced scoring from deep analysis
	if analysis != nil && analysis.MediaMapping != nil {
		mapping := analysis.MediaMapping
		
		// Bonus for smart suggestions
		suggestionBonus := math.Min(float64(len(mapping.SmartSuggestions))*5, 25)
		score += suggestionBonus

		// Bonus for asset groups
		if len(mapping.AssetGroups) > 1 {
			score += 15 // Multiple asset types
		}

		// Type diversity bonus
		typeCount := len(mapping.AssetGroups)
		if typeCount >= 3 {
			score += 10 // Has image, video, and audio
		}
	}

	return math.Min(score, maxScore)
}

// scoreModularSystem scores modular system benefits
func (as *AutomationScorer) scoreModularSystem(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	if len(metadata.Compositions) < 2 {
		return 0.0 // Need multiple compositions for modularity
	}

	score := 0.0
	maxScore := 100.0

	// Base score from composition count
	compScore := math.Min(float64(len(metadata.Compositions))*8, 40) // Cap at 40 points
	score += compScore

	// Enhanced scoring from deep analysis
	if analysis != nil && analysis.ModularSystem != nil {
		modular := analysis.ModularSystem
		
		// Bonus for discovered modular components
		componentBonus := math.Min(float64(modular.TotalModules)*10, 40)
		score += componentBonus

		// Bonus for mix-match options
		mixMatchBonus := math.Min(float64(len(modular.MixMatchOptions))*5, 20)
		score += mixMatchBonus

		// Huge bonus for high variant potential
		if modular.VariantPotential > 10 {
			score += 20
		} else if modular.VariantPotential > 5 {
			score += 10
		}
	}

	return math.Min(score, maxScore)
}

// scoreEffectAutomation scores effect customization potential
func (as *AutomationScorer) scoreEffectAutomation(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	if len(metadata.Effects) == 0 {
		return 0.0
	}

	score := 0.0
	maxScore := 100.0

	// Base score from effect count and customizability
	customizableCount := 0
	for _, effect := range metadata.Effects {
		if effect.IsCustomizable {
			customizableCount++
		}
	}

	effectScore := math.Min(float64(customizableCount)*8, 60) // Cap at 60 points
	score += effectScore

	// Bonus from effect chains in deep analysis
	if analysis != nil && len(analysis.EffectChains) > 0 {
		chainBonus := math.Min(float64(len(analysis.EffectChains))*15, 40)
		score += chainBonus
	}

	return math.Min(score, maxScore)
}

// scoreDataBinding scores data-driven content potential
func (as *AutomationScorer) scoreDataBinding(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	score := 0.0
	maxScore := 100.0

	// Base score from capabilities
	if metadata.Capabilities.HasDataDriven {
		score += 30
	}

	if metadata.Capabilities.HasExpressions {
		score += 25
	}

	// Enhanced scoring from deep analysis
	if analysis != nil && analysis.TextIntelligence != nil {
		intel := analysis.TextIntelligence
		
		// Bonus for data binding options
		if len(intel.DataBindingOptions) > 0 {
			bindingScore := math.Min(float64(len(intel.DataBindingOptions))*8, 45)
			score += bindingScore
		}
	}

	return math.Min(score, maxScore)
}

// scoreAPIReadiness scores how ready the template is for API integration
func (as *AutomationScorer) scoreAPIReadiness(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	score := 0.0
	maxScore := 100.0

	// Base score from automation opportunities
	if len(metadata.Opportunities) > 0 {
		oppScore := math.Min(float64(len(metadata.Opportunities))*15, 45)
		score += oppScore
	}

	// Strong bonus if API schema was generated
	if analysis != nil && analysis.APISchema != nil {
		score += 40 // Major bonus for having auto-generated API

		// Extra bonus for comprehensive API
		if len(analysis.APISchema.Endpoints) > 1 {
			score += 10
		}
		if len(analysis.APISchema.Examples) > 0 {
			score += 5
		}
	}

	return math.Min(score, maxScore)
}

// scoreMaintenanceComplexity scores maintenance burden (lower is better, but inverted for final score)
func (as *AutomationScorer) scoreMaintenanceComplexity(metadata *ProjectMetadata, analysis *DeepAnalysisResult) float64 {
	complexity := 0.0
	maxComplexity := 100.0

	// Complexity from layer count
	totalLayers := 0
	for _, comp := range metadata.Compositions {
		totalLayers += comp.LayerCount
	}
	layerComplexity := math.Min(float64(totalLayers)*2, 50)
	complexity += layerComplexity

	// Complexity from effects
	effectComplexity := math.Min(float64(len(metadata.Effects))*3, 30)
	complexity += effectComplexity

	// Complexity from analysis
	if analysis != nil {
		if analysis.ComplexityScore > 70 {
			complexity += 20
		} else if analysis.ComplexityScore > 40 {
			complexity += 10
		}
	}

	// Convert complexity to score (invert so higher is better)
	maintenanceScore := math.Max(0, maxComplexity-complexity)
	return maintenanceScore
}

// generateOpportunities creates specific automation opportunities
func (as *AutomationScorer) generateOpportunities(metadata *ProjectMetadata, analysis *DeepAnalysisResult, categoryScores map[string]float64) []AutomationOpportunity {
	opportunities := []AutomationOpportunity{}

	// Text automation opportunities
	if categoryScores["text_automation"] > 60 {
		opp := AutomationOpportunity{
			ID:               "text_automation_" + metadata.FileName,
			Type:             "text_automation",
			Title:            "Dynamic Text Content System",
			Description:      fmt.Sprintf("Implement automated text replacement for %d text layers with intelligent field detection", len(metadata.TextLayers)),
			AutomationScore:  categoryScores["text_automation"],
			ComplexityScore:  30.0, // Text automation is usually straightforward
			ROIScore:         as.calculateROI(categoryScores["text_automation"], 30.0),
			ImplementationEffort: "Low",
			BusinessImpact:   "High",
			TechnicalDetails: map[string]interface{}{
				"text_layers":     len(metadata.TextLayers),
				"field_types":     as.extractFieldTypes(analysis),
				"api_endpoints":   []string{"/render", "/preview"},
				"data_validation": true,
			},
			Prerequisites:    []string{"Define text field schema", "Create validation rules"},
			Recommendations:  []string{"Start with highest-impact text fields", "Implement progressive enhancement"},
		}
		opportunities = append(opportunities, opp)
	}

	// Media automation opportunities
	if categoryScores["media_automation"] > 50 {
		replaceableCount := as.countReplaceableAssets(metadata)
		opp := AutomationOpportunity{
			ID:               "media_automation_" + metadata.FileName,
			Type:             "media_automation",
			Title:            "Smart Media Replacement System",
			Description:      fmt.Sprintf("Automate replacement of %d media assets with intelligent sizing and validation", replaceableCount),
			AutomationScore:  categoryScores["media_automation"],
			ComplexityScore:  50.0, // Media requires more validation
			ROIScore:         as.calculateROI(categoryScores["media_automation"], 50.0),
			ImplementationEffort: "Medium",
			BusinessImpact:   "High",
			TechnicalDetails: map[string]interface{}{
				"replaceable_assets": replaceableCount,
				"asset_types":        as.extractAssetTypes(metadata),
				"validation_rules":   true,
				"smart_cropping":     true,
			},
			Prerequisites:    []string{"Define asset requirements", "Implement validation pipeline"},
			Recommendations:  []string{"Start with image replacement", "Add video support incrementally"},
		}
		opportunities = append(opportunities, opp)
	}

	// Modular system opportunities
	if categoryScores["modular_system"] > 40 {
		opp := AutomationOpportunity{
			ID:               "modular_system_" + metadata.FileName,
			Type:             "modular_system",
			Title:            "Mix-and-Match Template System",
			Description:      fmt.Sprintf("Create flexible template system with %d modular components", len(metadata.Compositions)),
			AutomationScore:  categoryScores["modular_system"],
			ComplexityScore:  70.0, // Modular systems are complex
			ROIScore:         as.calculateROI(categoryScores["modular_system"], 70.0),
			ImplementationEffort: "High",
			BusinessImpact:   "Very High",
			TechnicalDetails: map[string]interface{}{
				"module_count":    len(metadata.Compositions),
				"variant_potential": as.extractVariantPotential(analysis),
				"dependency_graph": true,
				"rule_engine":     true,
			},
			Prerequisites:    []string{"Map component dependencies", "Create combination rules", "Build preview system"},
			Recommendations:  []string{"Phase implementation by module type", "Start with most independent modules"},
		}
		opportunities = append(opportunities, opp)
	}

	// API integration opportunities
	if categoryScores["api_integration"] > 55 {
		opp := AutomationOpportunity{
			ID:               "api_integration_" + metadata.FileName,
			Type:             "api_integration",
			Title:            "Production API Development",
			Description:      "Build comprehensive API for programmatic video generation with full documentation",
			AutomationScore:  categoryScores["api_integration"],
			ComplexityScore:  60.0, // API development requires careful design
			ROIScore:         as.calculateROI(categoryScores["api_integration"], 60.0),
			ImplementationEffort: "High",
			BusinessImpact:   "Very High",
			TechnicalDetails: map[string]interface{}{
				"endpoints":      as.extractAPIEndpoints(analysis),
				"authentication": true,
				"rate_limiting":  true,
				"documentation":  true,
			},
			Prerequisites:    []string{"Design API schema", "Implement authentication", "Create documentation"},
			Recommendations:  []string{"Start with render endpoint", "Add webhook support", "Implement rate limiting"},
		}
		opportunities = append(opportunities, opp)
	}

	// Sort opportunities by ROI score
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].ROIScore > opportunities[j].ROIScore
	})

	return opportunities
}

// generateRecommendations creates actionable recommendations
func (as *AutomationScorer) generateRecommendations(categoryScores map[string]float64, breakdown ScoreBreakdown) []string {
	recommendations := []string{}

	// Text automation recommendations
	if categoryScores["text_automation"] > 70 {
		recommendations = append(recommendations, "🎯 High Priority: Implement text automation system - excellent ROI potential")
	} else if categoryScores["text_automation"] > 40 {
		recommendations = append(recommendations, "📝 Consider text automation for frequently updated content")
	}

	// Media automation recommendations  
	if categoryScores["media_automation"] > 60 {
		recommendations = append(recommendations, "🖼️ Implement smart media replacement system with validation")
	}

	// Modular system recommendations
	if categoryScores["modular_system"] > 50 {
		recommendations = append(recommendations, "🔧 Build modular template system for maximum flexibility")
	}

	// API integration recommendations
	if categoryScores["api_integration"] > 55 {
		recommendations = append(recommendations, "🚀 Develop production API for programmatic video generation")
	} else if breakdown.OverallScore > 60 {
		recommendations = append(recommendations, "📡 Consider API development for external integrations")
	}

	// Overall recommendations based on total score
	if breakdown.OverallScore > 80 {
		recommendations = append(recommendations, "⭐ Excellent automation candidate - invest in full pipeline")
	} else if breakdown.OverallScore > 60 {
		recommendations = append(recommendations, "✅ Good automation potential - focus on highest-ROI areas")
	} else if breakdown.OverallScore > 40 {
		recommendations = append(recommendations, "⚠️ Moderate potential - selective automation recommended")
	} else {
		recommendations = append(recommendations, "📋 Low automation priority - focus on manual workflows")
	}

	return recommendations
}

// Helper methods

func (as *AutomationScorer) calculateROI(automationScore, complexityScore float64) float64 {
	// ROI = (Automation Benefit - Implementation Cost) / Implementation Cost
	// Simplified as: Automation Score weighted against Complexity
	benefit := automationScore
	cost := complexityScore
	
	if cost == 0 {
		return benefit // Avoid division by zero
	}
	
	roi := ((benefit - cost*0.5) / cost) * 100
	return math.Max(0, math.Min(roi, 100)) // Clamp between 0-100
}

func (as *AutomationScorer) extractFieldTypes(analysis *DeepAnalysisResult) []string {
	if analysis == nil || analysis.TextIntelligence == nil {
		return []string{}
	}
	
	types := []string{}
	for fieldType := range analysis.TextIntelligence.PatternGroups {
		types = append(types, fieldType)
	}
	return types
}

func (as *AutomationScorer) countReplaceableAssets(metadata *ProjectMetadata) int {
	count := 0
	for _, asset := range metadata.MediaAssets {
		if asset.IsPlaceholder {
			count++
		}
	}
	return count
}

func (as *AutomationScorer) extractAssetTypes(metadata *ProjectMetadata) []string {
	typeSet := make(map[string]bool)
	for _, asset := range metadata.MediaAssets {
		if asset.IsPlaceholder {
			typeSet[asset.Type] = true
		}
	}
	
	types := []string{}
	for assetType := range typeSet {
		types = append(types, assetType)
	}
	return types
}

func (as *AutomationScorer) extractVariantPotential(analysis *DeepAnalysisResult) int {
	if analysis == nil || analysis.ModularSystem == nil {
		return 1
	}
	return analysis.ModularSystem.VariantPotential
}

func (as *AutomationScorer) extractAPIEndpoints(analysis *DeepAnalysisResult) []string {
	if analysis == nil || analysis.APISchema == nil {
		return []string{"/render"}
	}
	
	endpoints := []string{}
	for _, endpoint := range analysis.APISchema.Endpoints {
		endpoints = append(endpoints, endpoint.Path)
	}
	return endpoints
}

// MapAssetRelationships discovers relationships between assets
func (arm *AssetRelationshipMapper) MapAssetRelationships(metadata *ProjectMetadata, analysis *DeepAnalysisResult) ([]AssetRelationship, []AssetGroup) {
	relationships := []AssetRelationship{}
	groups := []AssetGroup{}

	// Create asset groups by type
	typeGroups := make(map[string][]string)
	for _, asset := range metadata.MediaAssets {
		typeGroups[asset.Type] = append(typeGroups[asset.Type], asset.ID)
	}

	// Create groups
	groupID := 0
	for assetType, assetIDs := range typeGroups {
		if len(assetIDs) > 1 {
			group := AssetGroup{
				ID:          fmt.Sprintf("group_%d", groupID),
				Name:        fmt.Sprintf("%s Assets", strings.Title(assetType)),
				Type:        "content",
				AssetIDs:    assetIDs,
				Replaceable: true,
				Priority:    arm.calculateGroupPriority(assetType),
			}
			groups = append(groups, group)
			groupID++
		}
	}

	// Create replacement relationships for placeholders
	for _, asset := range metadata.MediaAssets {
		if asset.IsPlaceholder {
			// Find similar assets that could be replacements
			for _, otherAsset := range metadata.MediaAssets {
				if otherAsset.ID != asset.ID && otherAsset.Type == asset.Type && !otherAsset.IsPlaceholder {
					rel := AssetRelationship{
						SourceAssetID:    asset.ID,
						TargetAssetID:    otherAsset.ID,
						RelationshipType: "replacement",
						Strength:         0.8,
						Properties: map[string]interface{}{
							"type_match": true,
							"size_compatible": true,
						},
						Bidirectional: false,
					}
					relationships = append(relationships, rel)
				}
			}
		}
	}

	return relationships, groups
}

func (arm *AssetRelationshipMapper) calculateGroupPriority(assetType string) int {
	priorities := map[string]int{
		"image": 8,
		"video": 9,
		"audio": 6,
		"text":  7,
	}
	
	if priority, ok := priorities[assetType]; ok {
		return priority
	}
	return 5 // Default priority
}

// ToJSON exports scoring results as JSON
func (sr *ScoringResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(sr, "", "  ")
}

// GenerateScoringReport creates a human-readable scoring report
func (sr *ScoringResult) GenerateScoringReport() string {
	report := fmt.Sprintf(`# Automation Scoring Report

## Overall Score: %.1f/100

### Score Breakdown

- **Text Automation**: %.1f/100
- **Media Automation**: %.1f/100  
- **Modular System**: %.1f/100
- **Effect Automation**: %.1f/100
- **API Integration**: %.1f/100

### Detailed Scores

- Text Score: %.1f (Weight: %.1f%%)
- Media Score: %.1f (Weight: %.1f%%)
- Modular Score: %.1f (Weight: %.1f%%)
- Effect Score: %.1f (Weight: %.1f%%)
- Data Binding: %.1f (Weight: %.1f%%)
- API Readiness: %.1f (Weight: %.1f%%)
- Maintenance: %.1f (Weight: %.1f%%)

`, sr.OverallScore,
		sr.CategoryScores["text_automation"],
		sr.CategoryScores["media_automation"],
		sr.CategoryScores["modular_system"],
		sr.CategoryScores["effect_automation"],
		sr.CategoryScores["api_integration"],
		sr.ScoreBreakdown.TextScore, 25.0,
		sr.ScoreBreakdown.MediaScore, 20.0,
		sr.ScoreBreakdown.ModularScore, 15.0,
		sr.ScoreBreakdown.EffectScore, 10.0,
		sr.ScoreBreakdown.DataBindingScore, 15.0,
		sr.ScoreBreakdown.APIScore, 10.0,
		sr.ScoreBreakdown.MaintenanceScore, 5.0)

	// Add opportunities
	if len(sr.Opportunities) > 0 {
		report += "## Automation Opportunities\n\n"
		for i, opp := range sr.Opportunities {
			report += fmt.Sprintf("### %d. %s (ROI: %.1f)\n", i+1, opp.Title, opp.ROIScore)
			report += fmt.Sprintf("- **Type**: %s\n", opp.Type)
			report += fmt.Sprintf("- **Description**: %s\n", opp.Description)
			report += fmt.Sprintf("- **Effort**: %s\n", opp.ImplementationEffort)
			report += fmt.Sprintf("- **Impact**: %s\n", opp.BusinessImpact)
			report += fmt.Sprintf("- **Automation Score**: %.1f/100\n\n", opp.AutomationScore)
		}
	}

	// Add recommendations
	if len(sr.Recommendations) > 0 {
		report += "## Recommendations\n\n"
		for _, rec := range sr.Recommendations {
			report += fmt.Sprintf("- %s\n", rec)
		}
	}

	return report
}