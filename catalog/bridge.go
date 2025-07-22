// Package catalog provides a bridge between the Go AEP parser and Python cataloging system
package catalog

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	aep "github.com/mojosolo/mobot2025"
)

// ProjectMetadata represents the enhanced metadata for an AE project
type ProjectMetadata struct {
	// Basic Info
	FilePath     string    `json:"file_path"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	ParsedAt     time.Time `json:"parsed_at"`
	
	// AE Project Info
	BitDepth         uint8   `json:"bit_depth"`
	ExpressionEngine string  `json:"expression_engine"`
	TotalItems       int     `json:"total_items"`
	
	// Composition Info
	Compositions []CompositionInfo `json:"compositions"`
	
	// Asset Discovery
	TextLayers    []TextLayerInfo  `json:"text_layers"`
	MediaAssets   []MediaAssetInfo `json:"media_assets"`
	Effects       []EffectInfo     `json:"effects"`
	
	// Categorization
	Categories    []string         `json:"categories"`
	Tags          []string         `json:"tags"`
	Capabilities  ProjectCapabilities `json:"capabilities"`
	
	// Opportunities
	Opportunities []Opportunity    `json:"opportunities"`
	
	// S3 Storage fields
	S3Bucket      string    `json:"s3_bucket,omitempty"`
	S3Key         string    `json:"s3_key,omitempty"`
	S3VersionID   string    `json:"s3_version_id,omitempty"`
}

// CompositionInfo represents a composition within the project
type CompositionInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	FrameRate    float64   `json:"frame_rate"`
	Duration     float64   `json:"duration"`
	LayerCount   int       `json:"layer_count"`
	Is3D         bool      `json:"is_3d"`
	HasEffects   bool      `json:"has_effects"`
}

// TextLayerInfo represents a text layer that can be customized
type TextLayerInfo struct {
	ID           string `json:"id"`
	CompID       string `json:"comp_id"`
	LayerName    string `json:"layer_name"`
	SourceText   string `json:"source_text"`
	FontUsed     string `json:"font_used,omitempty"`
	IsAnimated   bool   `json:"is_animated"`
	HasExpressions bool `json:"has_expressions"`
	Is3D         bool   `json:"is_3d"`
}

// MediaAssetInfo represents a media asset (image, video, audio)
type MediaAssetInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"` // image, video, audio
	Path         string `json:"path,omitempty"`
	IsPlaceholder bool  `json:"is_placeholder"`
	UsageCount   int    `json:"usage_count"`
}

// EffectInfo represents an effect applied in the project
type EffectInfo struct {
	Name         string `json:"name"`
	Category     string `json:"category"`
	UsageCount   int    `json:"usage_count"`
	IsCustomizable bool `json:"is_customizable"`
}

// ProjectCapabilities represents what can be customized in the project
type ProjectCapabilities struct {
	HasTextReplacement  bool `json:"has_text_replacement"`
	HasImageReplacement bool `json:"has_image_replacement"`
	HasColorControl     bool `json:"has_color_control"`
	HasAudioReplacement bool `json:"has_audio_replacement"`
	HasDataDriven       bool `json:"has_data_driven"`
	HasExpressions      bool `json:"has_expressions"`
	IsModular           bool `json:"is_modular"`
}

// Opportunity represents a customization opportunity
type Opportunity struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"` // easy, medium, hard
	Impact      string `json:"impact"`     // low, medium, high
	Components  []string `json:"components"`
}

// Parser provides the interface for parsing AEP files
type Parser struct {
	// Configuration
	ExtractText   bool
	ExtractMedia  bool
	DeepAnalysis  bool
}

// NewParser creates a new AEP parser with catalog capabilities
func NewParser() *Parser {
	return &Parser{
		ExtractText:  true,
		ExtractMedia: true,
		DeepAnalysis: true,
	}
}

// ParseProject parses an AEP file and extracts comprehensive metadata
func (p *Parser) ParseProject(aepPath string) (*ProjectMetadata, error) {
	// Parse the AEP file using the Go parser
	project, err := aep.Open(aepPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AEP file: %w", err)
	}

	// Initialize metadata
	metadata := &ProjectMetadata{
		FilePath:         aepPath,
		FileName:         filepath.Base(aepPath),
		ParsedAt:         time.Now(),
		BitDepth:         uint8(project.Depth),
		ExpressionEngine: project.ExpressionEngine,
		TotalItems:       len(project.Items),
		Compositions:     []CompositionInfo{},
		TextLayers:       []TextLayerInfo{},
		MediaAssets:      []MediaAssetInfo{},
		Effects:          []EffectInfo{},
		Categories:       []string{},
		Tags:             []string{},
		Opportunities:    []Opportunity{},
	}

	// Extract compositions and layers
	effectMap := make(map[string]int)
	
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeComposition:
			comp := CompositionInfo{
				ID:        fmt.Sprintf("%d", item.ID),
				Name:      item.Name,
				Width:     int(item.FootageDimensions[0]),
				Height:    int(item.FootageDimensions[1]),
				FrameRate: item.FootageFramerate,
				Duration:  item.FootageSeconds,
			}
			
			// Analyze layers in composition
			layerCount := 0
			for _, layer := range item.CompositionLayers {
				layerCount++
				
				// Extract text layers - note: aep package doesn't expose SourceText directly
				// For now, we'll track all layers as potential text layers
				if p.ExtractText && layer.Name != "" {
					textLayer := TextLayerInfo{
						ID:         fmt.Sprintf("%d_%d", item.ID, layer.Index),
						CompID:     fmt.Sprintf("%d", item.ID),
						LayerName:  layer.Name,
						SourceText: "", // Would need deeper RIFX parsing
						IsAnimated: false, // Would need property tracking
					}
					
					// Check for 3D properties
					if layer.ThreeDEnabled {
						textLayer.Is3D = true
						comp.Is3D = true
					}
					
					metadata.TextLayers = append(metadata.TextLayers, textLayer)
				}
			}
			
			comp.LayerCount = layerCount
			metadata.Compositions = append(metadata.Compositions, comp)
			
		case aep.ItemTypeFootage:
			// Extract media assets
			if p.ExtractMedia {
				assetType := "unknown"
				if strings.HasPrefix(item.Name, "audio") || strings.HasSuffix(item.Name, ".mp3") || strings.HasSuffix(item.Name, ".wav") {
					assetType = "audio"
				} else if strings.HasSuffix(item.Name, ".mp4") || strings.HasSuffix(item.Name, ".mov") {
					assetType = "video"
				} else if strings.HasSuffix(item.Name, ".png") || strings.HasSuffix(item.Name, ".jpg") || strings.HasSuffix(item.Name, ".jpeg") {
					assetType = "image"
				}
				
				media := MediaAssetInfo{
					ID:            fmt.Sprintf("%d", item.ID),
					Name:          item.Name,
					Type:          assetType,
					IsPlaceholder: strings.Contains(strings.ToLower(item.Name), "placeholder"),
				}
				metadata.MediaAssets = append(metadata.MediaAssets, media)
			}
		}
	}
	
	// Process effects
	for effectName, count := range effectMap {
		effect := EffectInfo{
			Name:           effectName,
			UsageCount:     count,
			IsCustomizable: true, // Assume all effects are customizable for now
		}
		metadata.Effects = append(metadata.Effects, effect)
	}
	
	// Analyze capabilities
	metadata.Capabilities = p.analyzeCapabilities(metadata)
	
	// Categorize project
	metadata.Categories = p.categorizeProject(metadata)
	metadata.Tags = p.generateTags(metadata)
	
	// Identify opportunities
	if p.DeepAnalysis {
		metadata.Opportunities = p.identifyOpportunities(metadata)
	}
	
	return metadata, nil
}

// analyzeCapabilities determines what the project can do
func (p *Parser) analyzeCapabilities(metadata *ProjectMetadata) ProjectCapabilities {
	caps := ProjectCapabilities{}
	
	caps.HasTextReplacement = len(metadata.TextLayers) > 0
	caps.HasImageReplacement = false
	caps.HasAudioReplacement = false
	
	for _, asset := range metadata.MediaAssets {
		if asset.IsPlaceholder {
			switch asset.Type {
			case "image":
				caps.HasImageReplacement = true
			case "audio":
				caps.HasAudioReplacement = true
			}
		}
	}
	
	// Check for expressions
	for _, text := range metadata.TextLayers {
		if text.HasExpressions {
			caps.HasExpressions = true
			break
		}
	}
	
	// Check for modular compositions
	caps.IsModular = len(metadata.Compositions) > 3
	
	return caps
}

// categorizeProject determines project categories
func (p *Parser) categorizeProject(metadata *ProjectMetadata) []string {
	categories := []string{}
	
	// Resolution-based categories
	for _, comp := range metadata.Compositions {
		if comp.Width == 1920 && comp.Height == 1080 {
			categories = append(categories, "HD")
		} else if comp.Width == 3840 && comp.Height == 2160 {
			categories = append(categories, "4K")
		} else if comp.Width == 1080 && comp.Height == 1920 {
			categories = append(categories, "Vertical")
		} else if comp.Width == 1080 && comp.Height == 1080 {
			categories = append(categories, "Square")
		}
	}
	
	// Content-based categories
	if metadata.Capabilities.HasTextReplacement {
		categories = append(categories, "Text Animation")
	}
	
	// Duration-based categories
	for _, comp := range metadata.Compositions {
		if comp.Duration <= 5 {
			categories = append(categories, "Bumper")
		} else if comp.Duration <= 15 {
			categories = append(categories, "Short Form")
		} else if comp.Duration <= 30 {
			categories = append(categories, "Standard")
		} else {
			categories = append(categories, "Long Form")
		}
	}
	
	return categories
}

// generateTags creates searchable tags
func (p *Parser) generateTags(metadata *ProjectMetadata) []string {
	tags := []string{}
	
	// Add capability tags
	if metadata.Capabilities.HasTextReplacement {
		tags = append(tags, "text", "typography", "animated-text")
	}
	if metadata.Capabilities.HasImageReplacement {
		tags = append(tags, "image", "photo", "placeholder")
	}
	if metadata.Capabilities.IsModular {
		tags = append(tags, "modular", "scenes", "multi-part")
	}
	
	// Add effect tags
	for _, effect := range metadata.Effects {
		if strings.Contains(strings.ToLower(effect.Name), "glow") {
			tags = append(tags, "glow", "lighting")
		}
		if strings.Contains(strings.ToLower(effect.Name), "blur") {
			tags = append(tags, "blur", "depth")
		}
	}
	
	return tags
}

// identifyOpportunities finds automation opportunities
func (p *Parser) identifyOpportunities(metadata *ProjectMetadata) []Opportunity {
	opportunities := []Opportunity{}
	
	// Text replacement opportunity
	if metadata.Capabilities.HasTextReplacement {
		opp := Opportunity{
			Type:        "text_automation",
			Description: fmt.Sprintf("Automate %d text layers for dynamic content", len(metadata.TextLayers)),
			Difficulty:  "easy",
			Impact:      "high",
			Components:  []string{},
		}
		for _, text := range metadata.TextLayers {
			opp.Components = append(opp.Components, text.LayerName)
		}
		opportunities = append(opportunities, opp)
	}
	
	// Media replacement opportunity
	if metadata.Capabilities.HasImageReplacement {
		count := 0
		components := []string{}
		for _, asset := range metadata.MediaAssets {
			if asset.IsPlaceholder && asset.Type == "image" {
				count++
				components = append(components, asset.Name)
			}
		}
		
		if count > 0 {
			opp := Opportunity{
				Type:        "media_automation",
				Description: fmt.Sprintf("Replace %d placeholder images dynamically", count),
				Difficulty:  "medium",
				Impact:      "high",
				Components:  components,
			}
			opportunities = append(opportunities, opp)
		}
	}
	
	// Batch rendering opportunity
	if len(metadata.Compositions) > 1 {
		opp := Opportunity{
			Type:        "batch_rendering",
			Description: fmt.Sprintf("Batch render %d compositions with variations", len(metadata.Compositions)),
			Difficulty:  "medium",
			Impact:      "medium",
			Components:  []string{},
		}
		for _, comp := range metadata.Compositions {
			opp.Components = append(opp.Components, comp.Name)
		}
		opportunities = append(opportunities, opp)
	}
	
	// Modular composition opportunity
	if metadata.Capabilities.IsModular {
		opp := Opportunity{
			Type:        "modular_system",
			Description: "Create mix-and-match system from modular compositions",
			Difficulty:  "hard",
			Impact:      "high",
			Components:  []string{"All compositions"},
		}
		opportunities = append(opportunities, opp)
	}
	
	return opportunities
}

// ToJSON exports the metadata as JSON
func (m *ProjectMetadata) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// ToNexRenderConfig generates a nexrender configuration
func (m *ProjectMetadata) ToNexRenderConfig() map[string]interface{} {
	config := map[string]interface{}{
		"template": map[string]interface{}{
			"src": m.FilePath,
			"composition": func() string {
				if len(m.Compositions) > 0 {
					return m.Compositions[0].Name
				}
				return ""
			}(),
		},
		"assets": []map[string]interface{}{},
	}
	
	// Add text replacements
	for _, text := range m.TextLayers {
		asset := map[string]interface{}{
			"type": "data",
			"layerName": text.LayerName,
			"property": "Source Text",
			"value": text.SourceText,
		}
		config["assets"] = append(config["assets"].([]map[string]interface{}), asset)
	}
	
	// Add media replacements
	for _, media := range m.MediaAssets {
		if media.IsPlaceholder {
			asset := map[string]interface{}{
				"type": "image",
				"layerName": media.Name,
				"src": "placeholder.jpg", // To be replaced
			}
			config["assets"] = append(config["assets"].([]map[string]interface{}), asset)
		}
	}
	
	return config
}