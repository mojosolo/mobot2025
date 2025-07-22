// Package catalog provides template import capabilities from legacy mobot system
package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MobotTemplateImporter handles importing templates from mobot
type MobotTemplateImporter struct {
	parser   *Parser
	database *Database
	analyzer *DangerousAnalyzer
}

// LegacyTemplateConfig represents the old mobot template configuration
type LegacyTemplateConfig struct {
	TemplatePath     string   `json:"template_path"`
	TemplateName     string   `json:"template_name"`
	Category         string   `json:"category"`
	Tags             []string `json:"tags"`
	DifficultyRating int      `json:"difficulty_rating"`
	CustomizableElements struct {
		TextLayers []struct {
			Name         string `json:"name"`
			DefaultText  string `json:"default_text"`
			FontFamily   string `json:"font_family"`
			Customizable bool   `json:"customizable"`
		} `json:"text_layers"`
		ImagePlaceholders []struct {
			Name         string `json:"name"`
			Type         string `json:"type"`
			Dimensions   string `json:"dimensions"`
			Placeholder  bool   `json:"placeholder"`
		} `json:"image_placeholders"`
	} `json:"customizable_elements"`
}

// PatternConfig represents pattern configurations from mobot
type PatternConfig struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Structure    map[string]interface{} `json:"structure"`
	Templates    []string               `json:"templates"`
	Complexity   string                 `json:"complexity"`
	UseCases     []string               `json:"use_cases"`
	Requirements map[string]interface{} `json:"requirements"`
}

// ImportResult tracks the import process results
type ImportResult struct {
	TotalFound       int                 `json:"total_found"`
	SuccessfulImports int                `json:"successful_imports"`
	FailedImports    []string            `json:"failed_imports"`
	ImportedProjects []*ProjectMetadata  `json:"imported_projects"`
	PatternConfigs   []PatternConfig     `json:"pattern_configs"`
	StartTime        time.Time           `json:"start_time"`
	EndTime          time.Time           `json:"end_time"`
	Duration         time.Duration       `json:"duration"`
}

// NewTemplateImporter creates a new template importer
func NewTemplateImporter(database *Database) *MobotTemplateImporter {
	return &MobotTemplateImporter{
		parser:   NewParser(),
		database: database,
		analyzer: NewDangerousAnalyzer(),
	}
}

// ImportFromMobot imports all templates from the mobot directory
func (ti *MobotTemplateImporter) ImportFromMobot(mobotDir string) (*ImportResult, error) {
	result := &ImportResult{
		StartTime:        time.Now(),
		FailedImports:    []string{},
		ImportedProjects: []*ProjectMetadata{},
		PatternConfigs:   []PatternConfig{},
	}

	fmt.Printf("🔍 Scanning mobot directory: %s\n", mobotDir)

	// Import pattern configurations first
	if err := ti.importPatternConfigs(mobotDir, result); err != nil {
		fmt.Printf("Warning: Failed to import pattern configs: %v\n", err)
	}

	// Find and import AEP templates
	templatesDir := filepath.Join(mobotDir, "templates")
	if err := ti.discoverAndImportTemplates(templatesDir, result); err != nil {
		return nil, fmt.Errorf("failed to import templates: %w", err)
	}

	// Also check root directory for any additional AEP files
	if err := ti.discoverAndImportTemplates(mobotDir, result); err != nil {
		fmt.Printf("Warning: Failed to scan root directory: %v\n", err)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	fmt.Printf("✅ Import completed in %v\n", result.Duration)
	fmt.Printf("   - Found: %d templates\n", result.TotalFound)
	fmt.Printf("   - Imported: %d successfully\n", result.SuccessfulImports)
	fmt.Printf("   - Failed: %d\n", len(result.FailedImports))

	return result, nil
}

// discoverAndImportTemplates recursively finds and imports AEP files
func (ti *MobotTemplateImporter) discoverAndImportTemplates(dir string, result *ImportResult) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking
		}

		// Skip directories and non-AEP files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".aep") {
			return nil
		}

		result.TotalFound++
		fmt.Printf("📁 Processing: %s\n", filepath.Base(path))

		// Import the template
		if err := ti.importSingleTemplate(path, result); err != nil {
			fmt.Printf("❌ Failed to import %s: %v\n", filepath.Base(path), err)
			result.FailedImports = append(result.FailedImports, path)
			return nil // Continue with other files
		}

		result.SuccessfulImports++
		return nil
	})
}

// importSingleTemplate imports a single AEP template
func (ti *MobotTemplateImporter) importSingleTemplate(aepPath string, result *ImportResult) error {
	// Parse with standard parser first
	metadata, err := ti.parser.ParseProject(aepPath)
	if err != nil {
		return fmt.Errorf("failed to parse AEP: %w", err)
	}

	// Check if legacy config exists
	configPath := ti.findLegacyConfig(aepPath)
	if configPath != "" {
		if err := ti.enhanceWithLegacyConfig(metadata, configPath); err != nil {
			fmt.Printf("Warning: Failed to load legacy config for %s: %v\n", filepath.Base(aepPath), err)
		}
	}

	// Perform deep analysis for additional insights
	analysis, err := ti.analyzer.AnalyzeProject(aepPath)
	if err != nil {
		fmt.Printf("Warning: Deep analysis failed for %s: %v\n", filepath.Base(aepPath), err)
	} else {
		// Enhance metadata with analysis insights
		ti.enhanceWithAnalysis(metadata, analysis)
	}

	// Store in database
	if err := ti.database.StoreProject(metadata); err != nil {
		return fmt.Errorf("failed to store in database: %w", err)
	}

	// Store analysis if available
	if analysis != nil {
		if err := ti.database.StoreAnalysisResult(1, analysis); err != nil { // Note: Using placeholder project ID
			fmt.Printf("Warning: Failed to store analysis: %v\n", err)
		}
	}

	result.ImportedProjects = append(result.ImportedProjects, metadata)
	return nil
}

// findLegacyConfig looks for template_config.json near the AEP file
func (ti *MobotTemplateImporter) findLegacyConfig(aepPath string) string {
	dir := filepath.Dir(aepPath)
	baseName := strings.TrimSuffix(filepath.Base(aepPath), ".aep")

	// Try various config file patterns
	patterns := []string{
		filepath.Join(dir, "template_config.json"),
		filepath.Join(dir, baseName+"_config.json"),
		filepath.Join(dir, "config.json"),
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern
		}
	}

	return ""
}

// enhanceWithLegacyConfig enhances metadata with legacy configuration
func (ti *MobotTemplateImporter) enhanceWithLegacyConfig(metadata *ProjectMetadata, configPath string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var legacyConfig LegacyTemplateConfig
	if err := json.Unmarshal(configData, &legacyConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Enhance categories
	if legacyConfig.Category != "" {
		metadata.Categories = append(metadata.Categories, legacyConfig.Category)
	}

	// Add tags
	metadata.Tags = append(metadata.Tags, legacyConfig.Tags...)

	// Enhance text layers with legacy info
	for i, textLayer := range metadata.TextLayers {
		for _, legacyText := range legacyConfig.CustomizableElements.TextLayers {
			if strings.Contains(textLayer.LayerName, legacyText.Name) {
				if legacyText.FontFamily != "" {
					metadata.TextLayers[i].FontUsed = legacyText.FontFamily
				}
				break
			}
		}
	}

	// Add difficulty-based tags
	switch {
	case legacyConfig.DifficultyRating <= 2:
		metadata.Tags = append(metadata.Tags, "beginner", "easy")
	case legacyConfig.DifficultyRating <= 4:
		metadata.Tags = append(metadata.Tags, "intermediate", "medium")
	default:
		metadata.Tags = append(metadata.Tags, "advanced", "hard")
	}

	return nil
}

// enhanceWithAnalysis enhances metadata with deep analysis results
func (ti *MobotTemplateImporter) enhanceWithAnalysis(metadata *ProjectMetadata, analysis *DeepAnalysisResult) {
	// Add analysis-based categories
	if analysis.ModularSystem != nil && analysis.ModularSystem.TotalModules > 2 {
		metadata.Categories = append(metadata.Categories, "Modular System")
	}

	if analysis.ComplexityScore > 70 {
		metadata.Categories = append(metadata.Categories, "Advanced")
	} else if analysis.ComplexityScore > 40 {
		metadata.Categories = append(metadata.Categories, "Intermediate")
	} else {
		metadata.Categories = append(metadata.Categories, "Beginner")
	}

	// Add automation-based tags
	if analysis.AutomationScore > 80 {
		metadata.Tags = append(metadata.Tags, "highly-automatable", "api-ready")
	} else if analysis.AutomationScore > 60 {
		metadata.Tags = append(metadata.Tags, "automatable")
	}

	// Add intelligence-based tags
	if analysis.TextIntelligence != nil {
		for fieldType := range analysis.TextIntelligence.PatternGroups {
			metadata.Tags = append(metadata.Tags, "has-"+fieldType)
		}

		if analysis.TextIntelligence.LocalizationReady {
			metadata.Tags = append(metadata.Tags, "localization-ready")
		}
	}

	// Add media-based tags
	if analysis.MediaMapping != nil {
		for _, asset := range analysis.MediaMapping.ReplaceableAssets {
			metadata.Tags = append(metadata.Tags, "replaceable-"+asset.Type)
		}
	}

	// Remove duplicates
	metadata.Categories = ti.removeDuplicates(metadata.Categories)
	metadata.Tags = ti.removeDuplicates(metadata.Tags)
}

// importPatternConfigs imports pattern configurations
func (ti *MobotTemplateImporter) importPatternConfigs(mobotDir string, result *ImportResult) error {
	patternsFile := filepath.Join(mobotDir, "pattern_selector.json")
	if _, err := os.Stat(patternsFile); err != nil {
		return nil // File doesn't exist, skip
	}

	fmt.Printf("📋 Importing pattern configurations...\n")

	data, err := os.ReadFile(patternsFile)
	if err != nil {
		return fmt.Errorf("failed to read patterns file: %w", err)
	}

	var patterns map[string]PatternConfig
	if err := json.Unmarshal(data, &patterns); err != nil {
		return fmt.Errorf("failed to parse patterns: %w", err)
	}

	for name, pattern := range patterns {
		pattern.Name = name
		result.PatternConfigs = append(result.PatternConfigs, pattern)
		fmt.Printf("   ✓ Pattern: %s\n", name)
	}

	return nil
}

// removeDuplicates removes duplicate strings from slice
func (ti *MobotTemplateImporter) removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// GenerateImportReport creates a comprehensive report of the import
func (ti *MobotTemplateImporter) GenerateImportReport(result *ImportResult) string {
	report := fmt.Sprintf(`# MoBot Template Import Report

**Generated**: %s
**Duration**: %v

## Summary

- **Templates Found**: %d
- **Successfully Imported**: %d
- **Failed Imports**: %d
- **Pattern Configurations**: %d

## Import Statistics

### Success Rate
- Success Rate: %.1f%%
- Processing Speed: %.2f templates/second

### Categories Discovered
`, result.EndTime.Format("2006-01-02 15:04:05"),
		result.Duration,
		result.TotalFound,
		result.SuccessfulImports,
		len(result.FailedImports),
		len(result.PatternConfigs),
		float64(result.SuccessfulImports)/float64(result.TotalFound)*100,
		float64(result.TotalFound)/result.Duration.Seconds())

	// Count categories
	categoryCount := make(map[string]int)
	tagCount := make(map[string]int)

	for _, project := range result.ImportedProjects {
		for _, cat := range project.Categories {
			categoryCount[cat]++
		}
		for _, tag := range project.Tags {
			tagCount[tag]++
		}
	}

	for cat, count := range categoryCount {
		report += fmt.Sprintf("- **%s**: %d templates\n", cat, count)
	}

	report += "\n### Top Tags\n"
	topTags := ti.getTopTags(tagCount, 10)
	for _, tag := range topTags {
		report += fmt.Sprintf("- `%s`: %d templates\n", tag.Name, tag.Count)
	}

	// Failed imports section
	if len(result.FailedImports) > 0 {
		report += "\n## Failed Imports\n\n"
		for _, failed := range result.FailedImports {
			report += fmt.Sprintf("- `%s`\n", failed)
		}
	}

	// Pattern configurations
	if len(result.PatternConfigs) > 0 {
		report += "\n## Pattern Configurations\n\n"
		for _, pattern := range result.PatternConfigs {
			report += fmt.Sprintf("### %s\n", pattern.Name)
			report += fmt.Sprintf("- **Description**: %s\n", pattern.Description)
			report += fmt.Sprintf("- **Complexity**: %s\n", pattern.Complexity)
			if len(pattern.UseCases) > 0 {
				report += fmt.Sprintf("- **Use Cases**: %s\n", strings.Join(pattern.UseCases, ", "))
			}
			report += "\n"
		}
	}

	// Detailed project listing
	report += "\n## Imported Projects\n\n"
	for _, project := range result.ImportedProjects {
		report += fmt.Sprintf("### %s\n", project.FileName)
		report += fmt.Sprintf("- **Path**: `%s`\n", project.FilePath)
		report += fmt.Sprintf("- **Compositions**: %d\n", len(project.Compositions))
		report += fmt.Sprintf("- **Text Layers**: %d\n", len(project.TextLayers))
		report += fmt.Sprintf("- **Media Assets**: %d\n", len(project.MediaAssets))
		if len(project.Categories) > 0 {
			report += fmt.Sprintf("- **Categories**: %s\n", strings.Join(project.Categories, ", "))
		}
		if len(project.Tags) > 0 {
			report += fmt.Sprintf("- **Tags**: %s\n", strings.Join(project.Tags, ", "))
		}
		report += fmt.Sprintf("- **Opportunities**: %d\n", len(project.Opportunities))
		report += "\n"
	}

	return report
}

// TagCount helper struct
type TagCount struct {
	Name  string
	Count int
}

// getTopTags returns the most frequent tags
func (ti *MobotTemplateImporter) getTopTags(tagCount map[string]int, limit int) []TagCount {
	var tags []TagCount
	for name, count := range tagCount {
		tags = append(tags, TagCount{Name: name, Count: count})
	}

	// Simple bubble sort by count
	for i := 0; i < len(tags)-1; i++ {
		for j := 0; j < len(tags)-i-1; j++ {
			if tags[j].Count < tags[j+1].Count {
				tags[j], tags[j+1] = tags[j+1], tags[j]
			}
		}
	}

	if len(tags) > limit {
		tags = tags[:limit]
	}

	return tags
}

// Import implements the TemplateImporter interface
func (ti *MobotTemplateImporter) Import(path string) (*Template, error) {
	// Create a basic template from the imported file
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	
	template := &Template{
		ID:        generateMobotTemplateID(),
		Name:      name,
		Path:      path,
		Type:      "mobot",
		Category:  "imported",
		Tags:      []string{"imported"},
		Metadata:  map[string]interface{}{"source": path},
		Status:    "imported",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return template, nil
}

// GetSupportedFormats implements the TemplateImporter interface
func (ti *MobotTemplateImporter) GetSupportedFormats() []string {
	return []string{".aep", ".aepx", ".aet"}
}

// ExportImportedData exports imported data to JSON for external use
func (ti *MobotTemplateImporter) ExportImportedData(result *ImportResult, outputPath string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	fmt.Printf("📄 Import data exported to: %s\n", outputPath)
	return nil
}

// generateMobotTemplateID generates a unique template ID
func generateMobotTemplateID() string {
	return fmt.Sprintf("tmpl_%d", time.Now().UnixNano())
}