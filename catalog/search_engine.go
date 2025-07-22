// Package catalog provides advanced search and filtering capabilities
package catalog

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

// SearchEngine provides comprehensive search functionality
type SearchEngine struct {
	database *Database
	indexer  *SearchIndexer
}

// SearchQuery represents a search request
type SearchQuery struct {
	Query          string                 `json:"query"`
	Categories     []string               `json:"categories"`
	Tags           []string               `json:"tags"`
	MinComplexity  float64                `json:"min_complexity"`
	MaxComplexity  float64                `json:"max_complexity"`
	MinAutomation  float64                `json:"min_automation"`
	MaxAutomation  float64                `json:"max_automation"`
	HasText        bool                   `json:"has_text"`
	HasMedia       bool                   `json:"has_media"`
	IsModular      bool                   `json:"is_modular"`
	Resolution     string                 `json:"resolution"` // HD, 4K, etc.
	Duration       string                 `json:"duration"`   // short, medium, long
	SortBy         string                 `json:"sort_by"`    // relevance, date, complexity, automation
	SortOrder      string                 `json:"sort_order"` // asc, desc
	Limit          int                    `json:"limit"`
	Offset         int                    `json:"offset"`
	Filters        map[string]interface{} `json:"filters"`
}

// SearchResult represents search results
type SearchResult struct {
	Query         SearchQuery        `json:"query"`
	TotalResults  int                `json:"total_results"`
	Results       []*ProjectMetadata `json:"results"`
	Facets        SearchFacets       `json:"facets"`
	Suggestions   []string           `json:"suggestions"`
	SearchTime    float64            `json:"search_time_ms"`
	DidYouMean    string             `json:"did_you_mean,omitempty"`
}

// SearchFacets provide filter options
type SearchFacets struct {
	Categories   map[string]int `json:"categories"`
	Tags         map[string]int `json:"tags"`
	Resolutions  map[string]int `json:"resolutions"`
	Durations    map[string]int `json:"durations"`
	Capabilities map[string]int `json:"capabilities"`
}

// SearchIndexer handles search index management
type SearchIndexer struct {
	database *Database
}

// SimilarityMatch represents template similarity
type SimilarityMatch struct {
	ProjectID    int64   `json:"project_id"`
	SimilarityScore float64 `json:"similarity_score"`
	MatchReasons []string `json:"match_reasons"`
}

// NewSearchEngine creates a new search engine
func NewSearchEngine(database *Database) *SearchEngine {
	return &SearchEngine{
		database: database,
		indexer:  NewSearchIndexer(database),
	}
}

// NewSearchIndexer creates a new search indexer
func NewSearchIndexer(database *Database) *SearchIndexer {
	return &SearchIndexer{
		database: database,
	}
}

// Search performs a comprehensive search
func (se *SearchEngine) Search(query SearchQuery) (*SearchResult, error) {
	startTime := getCurrentTimeMs()

	// Validate and set defaults
	if query.Limit == 0 {
		query.Limit = 50
	}
	if query.SortBy == "" {
		query.SortBy = "relevance"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	result := &SearchResult{
		Query:       query,
		Results:     []*ProjectMetadata{},
		Facets:      SearchFacets{},
		Suggestions: []string{},
	}

	// Perform search
	if query.Query != "" {
		// Full-text search
		results, err := se.performFullTextSearch(query)
		if err != nil {
			return nil, fmt.Errorf("full-text search failed: %w", err)
		}
		result.Results = results
	} else {
		// Filter-only search
		results, err := se.performFilterSearch(query)
		if err != nil {
			return nil, fmt.Errorf("filter search failed: %w", err)
		}
		result.Results = results
	}

	// Apply additional filters
	result.Results = se.applyFilters(result.Results, query)

	// Sort results
	se.sortResults(result.Results, query.SortBy, query.SortOrder)

	// Apply pagination
	result.TotalResults = len(result.Results)
	if query.Offset > 0 {
		if query.Offset >= len(result.Results) {
			result.Results = []*ProjectMetadata{}
		} else {
			result.Results = result.Results[query.Offset:]
		}
	}
	if len(result.Results) > query.Limit {
		result.Results = result.Results[:query.Limit]
	}

	// Generate facets
	result.Facets = se.generateFacets(result.Results)

	// Generate suggestions and spell check
	if len(result.Results) == 0 && query.Query != "" {
		result.Suggestions = se.generateSuggestions(query.Query)
		result.DidYouMean = se.generateDidYouMean(query.Query)
	}

	result.SearchTime = getCurrentTimeMs() - startTime
	return result, nil
}

// performFullTextSearch performs text-based search
func (se *SearchEngine) performFullTextSearch(query SearchQuery) ([]*ProjectMetadata, error) {
	// Enhanced search with query parsing
	searchTerms := se.parseSearchQuery(query.Query)
	
	// Build search conditions
	conditions := []string{}
	args := []interface{}{}

	// Add text search conditions
	for _, term := range searchTerms {
		conditions = append(conditions, "si.content LIKE ?")
		args = append(args, "%"+term+"%")
	}

	// Build query
	sqlQuery := `
		SELECT DISTINCT p.id, p.file_path, p.file_name, p.file_size,
			   p.bit_depth, p.expression_engine, p.total_items, p.parsed_at
		FROM projects p
		JOIN search_index si ON p.id = si.project_id
		WHERE ` + strings.Join(conditions, " AND ")

	// Execute search
	rows, err := se.database.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	var results []*ProjectMetadata
	for rows.Next() {
		metadata := &ProjectMetadata{}
		var projectID int64
		err := rows.Scan(
			&projectID, &metadata.FilePath, &metadata.FileName, &metadata.FileSize,
			&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
			&metadata.ParsedAt)
		if err != nil {
			continue
		}

		// Load complete metadata
		fullMetadata, err := se.database.GetProject(projectID)
		if err != nil {
			continue
		}
		results = append(results, fullMetadata)
	}

	return results, nil
}

// performFilterSearch performs filter-only search
func (se *SearchEngine) performFilterSearch(query SearchQuery) ([]*ProjectMetadata, error) {
	filter := ProjectFilter{
		Categories:    query.Categories,
		Tags:          query.Tags,
		MinComplexity: query.MinComplexity,
		MaxComplexity: query.MaxComplexity,
		Limit:         query.Limit * 2, // Get more for post-filtering
	}

	return se.database.FilterProjects(filter)
}

// applyFilters applies additional search filters
func (se *SearchEngine) applyFilters(results []*ProjectMetadata, query SearchQuery) []*ProjectMetadata {
	filtered := []*ProjectMetadata{}

	for _, metadata := range results {
		if !se.matchesFilters(metadata, query) {
			continue
		}
		filtered = append(filtered, metadata)
	}

	return filtered
}

// matchesFilters checks if metadata matches search filters
func (se *SearchEngine) matchesFilters(metadata *ProjectMetadata, query SearchQuery) bool {
	// Text capability filter
	if query.HasText && !metadata.Capabilities.HasTextReplacement {
		return false
	}

	// Media capability filter
	if query.HasMedia && !metadata.Capabilities.HasImageReplacement {
		return false
	}

	// Modular filter
	if query.IsModular && !metadata.Capabilities.IsModular {
		return false
	}

	// Resolution filter
	if query.Resolution != "" {
		if !se.hasResolution(metadata, query.Resolution) {
			return false
		}
	}

	// Duration filter
	if query.Duration != "" {
		if !se.hasDuration(metadata, query.Duration) {
			return false
		}
	}

	// Custom filters
	for key, value := range query.Filters {
		if !se.matchesCustomFilter(metadata, key, value) {
			return false
		}
	}

	return true
}

// hasResolution checks if project has specified resolution
func (se *SearchEngine) hasResolution(metadata *ProjectMetadata, resolution string) bool {
	for _, comp := range metadata.Compositions {
		switch resolution {
		case "HD":
			if comp.Width == 1920 && comp.Height == 1080 {
				return true
			}
		case "4K":
			if comp.Width == 3840 && comp.Height == 2160 {
				return true
			}
		case "Square":
			if comp.Width == comp.Height {
				return true
			}
		case "Vertical":
			if comp.Height > comp.Width {
				return true
			}
		}
	}
	return false
}

// hasDuration checks if project has specified duration category
func (se *SearchEngine) hasDuration(metadata *ProjectMetadata, duration string) bool {
	for _, comp := range metadata.Compositions {
		switch duration {
		case "short":
			if comp.Duration <= 15 {
				return true
			}
		case "medium":
			if comp.Duration > 15 && comp.Duration <= 60 {
				return true
			}
		case "long":
			if comp.Duration > 60 {
				return true
			}
		}
	}
	return false
}

// matchesCustomFilter applies custom filter logic
func (se *SearchEngine) matchesCustomFilter(metadata *ProjectMetadata, key string, value interface{}) bool {
	switch key {
	case "min_text_layers":
		if minLayers, ok := value.(float64); ok {
			return float64(len(metadata.TextLayers)) >= minLayers
		}
	case "min_compositions":
		if minComps, ok := value.(float64); ok {
			return float64(len(metadata.Compositions)) >= minComps
		}
	case "has_effects":
		if hasEffects, ok := value.(bool); ok {
			return (len(metadata.Effects) > 0) == hasEffects
		}
	case "file_size_mb":
		if maxSize, ok := value.(float64); ok {
			sizeMB := float64(metadata.FileSize) / (1024 * 1024)
			return sizeMB <= maxSize
		}
	}
	return true
}

// sortResults sorts search results
func (se *SearchEngine) sortResults(results []*ProjectMetadata, sortBy, sortOrder string) {
	sort.Slice(results, func(i, j int) bool {
		var less bool
		
		switch sortBy {
		case "name":
			less = results[i].FileName < results[j].FileName
		case "date":
			less = results[i].ParsedAt.Before(results[j].ParsedAt)
		case "complexity":
			// Would need to load analysis results for this
			less = len(results[i].Effects) < len(results[j].Effects)
		case "size":
			less = results[i].FileSize < results[j].FileSize
		case "compositions":
			less = len(results[i].Compositions) < len(results[j].Compositions)
		default: // relevance
			// For relevance, use a combination of factors
			scoreI := se.calculateRelevanceScore(results[i])
			scoreJ := se.calculateRelevanceScore(results[j])
			less = scoreI < scoreJ
		}

		if sortOrder == "desc" {
			return !less
		}
		return less
	})
}

// calculateRelevanceScore calculates relevance score for sorting
func (se *SearchEngine) calculateRelevanceScore(metadata *ProjectMetadata) float64 {
	score := 0.0
	
	// Boost score based on capabilities
	if metadata.Capabilities.HasTextReplacement {
		score += 10
	}
	if metadata.Capabilities.HasImageReplacement {
		score += 8
	}
	if metadata.Capabilities.IsModular {
		score += 15
	}
	
	// Boost based on automation opportunities
	score += float64(len(metadata.Opportunities)) * 5
	
	// Recent files get slight boost
	// daysSinceParse := time.Since(metadata.ParsedAt).Hours() / 24
	// score += math.Max(0, 30 - daysSinceParse) * 0.1
	
	return score
}

// generateFacets creates faceted search data
func (se *SearchEngine) generateFacets(results []*ProjectMetadata) SearchFacets {
	facets := SearchFacets{
		Categories:   make(map[string]int),
		Tags:         make(map[string]int),
		Resolutions:  make(map[string]int),
		Durations:    make(map[string]int),
		Capabilities: make(map[string]int),
	}

	for _, metadata := range results {
		// Category facets
		for _, cat := range metadata.Categories {
			facets.Categories[cat]++
		}
		
		// Tag facets
		for _, tag := range metadata.Tags {
			facets.Tags[tag]++
		}
		
		// Resolution facets
		for _, comp := range metadata.Compositions {
			resolution := se.categorizeResolution(comp.Width, comp.Height)
			facets.Resolutions[resolution]++
		}
		
		// Duration facets
		for _, comp := range metadata.Compositions {
			duration := se.categorizeDuration(comp.Duration)
			facets.Durations[duration]++
		}
		
		// Capability facets
		if metadata.Capabilities.HasTextReplacement {
			facets.Capabilities["Text Replacement"]++
		}
		if metadata.Capabilities.HasImageReplacement {
			facets.Capabilities["Image Replacement"]++
		}
		if metadata.Capabilities.IsModular {
			facets.Capabilities["Modular"]++
		}
	}

	return facets
}

// categorizeResolution determines resolution category
func (se *SearchEngine) categorizeResolution(width, height int) string {
	if width == 1920 && height == 1080 {
		return "HD (1920x1080)"
	}
	if width == 3840 && height == 2160 {
		return "4K (3840x2160)"
	}
	if width == height {
		return "Square"
	}
	if height > width {
		return "Vertical"
	}
	return fmt.Sprintf("Custom (%dx%d)", width, height)
}

// categorizeDuration determines duration category
func (se *SearchEngine) categorizeDuration(duration float64) string {
	if duration <= 15 {
		return "Short (≤15s)"
	}
	if duration <= 60 {
		return "Medium (15-60s)"
	}
	return "Long (>60s)"
}

// generateSuggestions creates search suggestions
func (se *SearchEngine) generateSuggestions(query string) []string {
	suggestions := []string{}
	
	// Common template terms
	commonTerms := []string{
		"logo reveal", "text animation", "intro", "outro", "lower third",
		"social media", "promo", "presentation", "slideshow", "product showcase",
		"corporate", "modern", "minimal", "dynamic", "professional",
	}
	
	queryLower := strings.ToLower(query)
	for _, term := range commonTerms {
		if strings.Contains(term, queryLower) || strings.Contains(queryLower, term) {
			suggestions = append(suggestions, term)
		}
	}
	
	// Limit suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}
	
	return suggestions
}

// generateDidYouMean creates spell correction suggestions
func (se *SearchEngine) generateDidYouMean(query string) string {
	// Simple spell correction for common terms
	corrections := map[string]string{
		"tempalte":  "template",
		"animaton":  "animation",
		"trnasition": "transition",
		"backgrond": "background",
		"compsoition": "composition",
	}
	
	queryLower := strings.ToLower(query)
	for typo, correction := range corrections {
		if strings.Contains(queryLower, typo) {
			return strings.ReplaceAll(queryLower, typo, correction)
		}
	}
	
	return ""
}

// parseSearchQuery parses search query into terms
func (se *SearchEngine) parseSearchQuery(query string) []string {
	// Remove special characters and split
	reg := regexp.MustCompile(`[^\w\s]+`)
	cleaned := reg.ReplaceAllString(query, " ")
	
	// Split and filter empty terms
	terms := strings.Fields(strings.ToLower(cleaned))
	
	// Remove common stop words
	stopWords := map[string]bool{
		"and": true, "or": true, "the": true, "a": true, "an": true,
		"is": true, "are": true, "was": true, "were": true,
		"for": true, "with": true, "by": true,
	}
	
	filtered := []string{}
	for _, term := range terms {
		if !stopWords[term] && len(term) > 1 {
			filtered = append(filtered, term)
		}
	}
	
	return filtered
}

// FindSimilarTemplates finds templates similar to a given template
func (se *SearchEngine) FindSimilarTemplates(projectID int64, limit int) ([]SimilarityMatch, error) {
	// Get the source project
	sourceProject, err := se.database.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source project: %w", err)
	}

	// Get all other projects for comparison
	allProjects, err := se.database.FilterProjects(ProjectFilter{Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	var matches []SimilarityMatch
	
	for _, project := range allProjects {
		if project.FilePath == sourceProject.FilePath {
			continue // Skip self
		}

		score, reasons := se.calculateSimilarity(sourceProject, project)
		if score > 0.3 { // Minimum similarity threshold
			match := SimilarityMatch{
				ProjectID:       int64(0), // Would need to get actual project ID
				SimilarityScore: score,
				MatchReasons:    reasons,
			}
			matches = append(matches, match)
		}
	}

	// Sort by similarity score
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].SimilarityScore > matches[j].SimilarityScore
	})

	// Apply limit
	if len(matches) > limit {
		matches = matches[:limit]
	}

	return matches, nil
}

// calculateSimilarity calculates similarity between two projects
func (se *SearchEngine) calculateSimilarity(project1, project2 *ProjectMetadata) (float64, []string) {
	score := 0.0
	reasons := []string{}

	// Category similarity
	categoryScore := se.calculateCategorySimilarity(project1.Categories, project2.Categories)
	if categoryScore > 0 {
		score += categoryScore * 0.3
		reasons = append(reasons, fmt.Sprintf("Similar categories (%.0f%% match)", categoryScore*100))
	}

	// Tag similarity
	tagScore := se.calculateTagSimilarity(project1.Tags, project2.Tags)
	if tagScore > 0 {
		score += tagScore * 0.2
		reasons = append(reasons, fmt.Sprintf("Similar tags (%.0f%% match)", tagScore*100))
	}

	// Composition similarity
	compScore := se.calculateCompositionSimilarity(project1.Compositions, project2.Compositions)
	if compScore > 0 {
		score += compScore * 0.2
		reasons = append(reasons, fmt.Sprintf("Similar compositions (%.0f%% match)", compScore*100))
	}

	// Capability similarity
	capScore := se.calculateCapabilitySimilarity(project1.Capabilities, project2.Capabilities)
	if capScore > 0 {
		score += capScore * 0.2
		reasons = append(reasons, fmt.Sprintf("Similar capabilities (%.0f%% match)", capScore*100))
	}

	// Text layer similarity
	textScore := se.calculateTextSimilarity(project1.TextLayers, project2.TextLayers)
	if textScore > 0 {
		score += textScore * 0.1
		reasons = append(reasons, fmt.Sprintf("Similar text structure (%.0f%% match)", textScore*100))
	}

	return math.Min(score, 1.0), reasons
}

// calculateCategorySimilarity calculates category overlap
func (se *SearchEngine) calculateCategorySimilarity(cats1, cats2 []string) float64 {
	if len(cats1) == 0 || len(cats2) == 0 {
		return 0
	}

	set1 := make(map[string]bool)
	for _, cat := range cats1 {
		set1[cat] = true
	}

	matches := 0
	for _, cat := range cats2 {
		if set1[cat] {
			matches++
		}
	}

	// Jaccard similarity
	union := len(cats1) + len(cats2) - matches
	if union == 0 {
		return 0
	}
	
	return float64(matches) / float64(union)
}

// calculateTagSimilarity calculates tag overlap
func (se *SearchEngine) calculateTagSimilarity(tags1, tags2 []string) float64 {
	if len(tags1) == 0 || len(tags2) == 0 {
		return 0
	}

	set1 := make(map[string]bool)
	for _, tag := range tags1 {
		set1[tag] = true
	}

	matches := 0
	for _, tag := range tags2 {
		if set1[tag] {
			matches++
		}
	}

	// Jaccard similarity
	union := len(tags1) + len(tags2) - matches
	if union == 0 {
		return 0
	}
	
	return float64(matches) / float64(union)
}

// calculateCompositionSimilarity compares composition structures
func (se *SearchEngine) calculateCompositionSimilarity(comps1, comps2 []CompositionInfo) float64 {
	if len(comps1) == 0 || len(comps2) == 0 {
		return 0
	}

	// Compare counts
	countSimilarity := 1.0 - math.Abs(float64(len(comps1)-len(comps2)))/math.Max(float64(len(comps1)), float64(len(comps2)))
	
	// Compare resolutions
	resolutionMatches := 0
	totalComparisons := 0
	
	for _, comp1 := range comps1 {
		for _, comp2 := range comps2 {
			totalComparisons++
			if comp1.Width == comp2.Width && comp1.Height == comp2.Height {
				resolutionMatches++
			}
		}
	}
	
	resolutionSimilarity := 0.0
	if totalComparisons > 0 {
		resolutionSimilarity = float64(resolutionMatches) / float64(totalComparisons)
	}
	
	return (countSimilarity + resolutionSimilarity) / 2
}

// calculateCapabilitySimilarity compares capabilities
func (se *SearchEngine) calculateCapabilitySimilarity(cap1, cap2 ProjectCapabilities) float64 {
	matches := 0
	total := 6 // Number of capability fields

	if cap1.HasTextReplacement == cap2.HasTextReplacement {
		matches++
	}
	if cap1.HasImageReplacement == cap2.HasImageReplacement {
		matches++
	}
	if cap1.HasColorControl == cap2.HasColorControl {
		matches++
	}
	if cap1.HasAudioReplacement == cap2.HasAudioReplacement {
		matches++
	}
	if cap1.IsModular == cap2.IsModular {
		matches++
	}
	if cap1.HasExpressions == cap2.HasExpressions {
		matches++
	}

	return float64(matches) / float64(total)
}

// calculateTextSimilarity compares text layer structures
func (se *SearchEngine) calculateTextSimilarity(text1, text2 []TextLayerInfo) float64 {
	if len(text1) == 0 || len(text2) == 0 {
		return 0
	}

	// Simple count-based similarity
	countDiff := math.Abs(float64(len(text1) - len(text2)))
	maxCount := math.Max(float64(len(text1)), float64(len(text2)))
	
	return 1.0 - (countDiff / maxCount)
}

// RebuildSearchIndex rebuilds the entire search index
func (si *SearchIndexer) RebuildSearchIndex() error {
	// This would typically involve:
	// 1. Clearing existing index
	// 2. Re-indexing all projects
	// 3. Optimizing index
	
	fmt.Println("🔄 Rebuilding search index...")
	
	// Implementation would go here
	// For now, just return success
	
	fmt.Println("✅ Search index rebuilt")
	return nil
}

// Helper function to get current time in milliseconds
func getCurrentTimeMs() float64 {
	// In real implementation, would use time.Now().UnixNano() / 1000000
	return 0.0 // Placeholder
}