// Package fixtures provides test data generation utilities for MoBot 2025
package fixtures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// MockAEPGenerator creates mock AEP files for testing
type MockAEPGenerator struct {
	rand *rand.Rand
}

// NewMockAEPGenerator creates a new generator with optional seed
func NewMockAEPGenerator(seed ...int64) *MockAEPGenerator {
	var s int64
	if len(seed) > 0 {
		s = seed[0]
	} else {
		s = time.Now().UnixNano()
	}
	return &MockAEPGenerator{
		rand: rand.New(rand.NewSource(s)),
	}
}

// GenerateMinimalAEP creates a minimal valid AEP file structure
func (g *MockAEPGenerator) GenerateMinimalAEP() ([]byte, error) {
	buf := &bytes.Buffer{}
	
	// Write RIFX header
	if err := g.writeRIFXHeader(buf); err != nil {
		return nil, fmt.Errorf("failed to write RIFX header: %w", err)
	}
	
	// Write minimal chunks
	if err := g.writeProjectChunk(buf); err != nil {
		return nil, fmt.Errorf("failed to write project chunk: %w", err)
	}
	
	return buf.Bytes(), nil
}

// GenerateComplexAEP creates a complex AEP file with multiple compositions
func (g *MockAEPGenerator) GenerateComplexAEP(numComps, numLayers int) ([]byte, error) {
	buf := &bytes.Buffer{}
	
	// Write RIFX header
	if err := g.writeRIFXHeader(buf); err != nil {
		return nil, fmt.Errorf("failed to write RIFX header: %w", err)
	}
	
	// Write project info
	if err := g.writeProjectChunk(buf); err != nil {
		return nil, fmt.Errorf("failed to write project chunk: %w", err)
	}
	
	// Write compositions
	for i := 0; i < numComps; i++ {
		if err := g.writeCompositionChunk(buf, fmt.Sprintf("Comp %d", i+1), numLayers); err != nil {
			return nil, fmt.Errorf("failed to write composition %d: %w", i, err)
		}
	}
	
	return buf.Bytes(), nil
}

// writeRIFXHeader writes the RIFX header
func (g *MockAEPGenerator) writeRIFXHeader(w io.Writer) error {
	// RIFX magic number
	if _, err := w.Write([]byte("RIFX")); err != nil {
		return err
	}
	
	// File size (placeholder - would be updated in real implementation)
	if err := binary.Write(w, binary.BigEndian, uint32(1024)); err != nil {
		return err
	}
	
	// Form type
	if _, err := w.Write([]byte("Egg!")); err != nil {
		return err
	}
	
	return nil
}

// writeProjectChunk writes a project information chunk
func (g *MockAEPGenerator) writeProjectChunk(w io.Writer) error {
	// Chunk ID
	if _, err := w.Write([]byte("proj")); err != nil {
		return err
	}
	
	// Chunk size
	if err := binary.Write(w, binary.BigEndian, uint32(100)); err != nil {
		return err
	}
	
	// Project data (simplified)
	data := struct {
		Version   uint32
		FrameRate uint32
		Width     uint32
		Height    uint32
	}{
		Version:   2024,
		FrameRate: 30,
		Width:     1920,
		Height:    1080,
	}
	
	return binary.Write(w, binary.BigEndian, data)
}

// writeCompositionChunk writes a composition chunk
func (g *MockAEPGenerator) writeCompositionChunk(w io.Writer, name string, numLayers int) error {
	// Chunk ID
	if _, err := w.Write([]byte("comp")); err != nil {
		return err
	}
	
	// Calculate chunk size
	nameLen := len(name)
	chunkSize := 16 + nameLen + (numLayers * 20) // Simplified calculation
	
	if err := binary.Write(w, binary.BigEndian, uint32(chunkSize)); err != nil {
		return err
	}
	
	// Write composition name length and name
	if err := binary.Write(w, binary.BigEndian, uint32(nameLen)); err != nil {
		return err
	}
	if _, err := w.Write([]byte(name)); err != nil {
		return err
	}
	
	// Write number of layers
	if err := binary.Write(w, binary.BigEndian, uint32(numLayers)); err != nil {
		return err
	}
	
	// Write simplified layer data
	for i := 0; i < numLayers; i++ {
		layer := struct {
			Type     uint32
			StartTime uint32
			Duration uint32
			Index    uint32
			Flags    uint32
		}{
			Type:     1, // Text layer
			StartTime: 0,
			Duration: 300,
			Index:    uint32(i),
			Flags:    0,
		}
		if err := binary.Write(w, binary.BigEndian, layer); err != nil {
			return err
		}
	}
	
	return nil
}

// ProjectConfig holds configuration for generating a project
type ProjectConfig struct {
	Name          string
	Compositions  []CompositionConfig
	BitDepth      uint8
	FrameRate     float64
	ExpressionEngine string
}

// CompositionConfig holds configuration for generating a composition
type CompositionConfig struct {
	Name      string
	Width     int
	Height    int
	Duration  float64
	Layers    []LayerConfig
}

// LayerConfig holds configuration for generating a layer
type LayerConfig struct {
	Name       string
	Type       string // "text", "solid", "null", "shape"
	SourceText string
	StartTime  float64
	Duration   float64
	Is3D       bool
}

// GenerateFromConfig creates an AEP file from a configuration
func (g *MockAEPGenerator) GenerateFromConfig(config ProjectConfig) ([]byte, error) {
	// This would generate a more realistic AEP structure based on the config
	// For now, we'll use the simple generator
	return g.GenerateComplexAEP(len(config.Compositions), 5)
}

// Common test configurations

// GetBasicTestProject returns a basic test project configuration
func GetBasicTestProject() ProjectConfig {
	return ProjectConfig{
		Name:          "Basic Test Project",
		BitDepth:      8,
		FrameRate:     30,
		ExpressionEngine: "javascript",
		Compositions: []CompositionConfig{
			{
				Name:     "Main Comp",
				Width:    1920,
				Height:   1080,
				Duration: 10,
				Layers: []LayerConfig{
					{
						Name:       "Title",
						Type:       "text",
						SourceText: "Test Title",
						StartTime:  0,
						Duration:   10,
						Is3D:       false,
					},
				},
			},
		},
	}
}

// GetComplexTestProject returns a complex test project configuration
func GetComplexTestProject() ProjectConfig {
	return ProjectConfig{
		Name:          "Complex Test Project",
		BitDepth:      16,
		FrameRate:     60,
		ExpressionEngine: "javascript",
		Compositions: []CompositionConfig{
			{
				Name:     "Intro",
				Width:    3840,
				Height:   2160,
				Duration: 5,
				Layers: []LayerConfig{
					{Name: "Background", Type: "solid", StartTime: 0, Duration: 5},
					{Name: "Logo", Type: "null", StartTime: 0, Duration: 5, Is3D: true},
					{Name: "Title", Type: "text", SourceText: "Welcome", StartTime: 1, Duration: 4},
				},
			},
			{
				Name:     "Main Content",
				Width:    1920,
				Height:   1080,
				Duration: 30,
				Layers: []LayerConfig{
					{Name: "Content BG", Type: "solid", StartTime: 0, Duration: 30},
					{Name: "Text 1", Type: "text", SourceText: "First Point", StartTime: 5, Duration: 10},
					{Name: "Text 2", Type: "text", SourceText: "Second Point", StartTime: 15, Duration: 10},
					{Name: "Transition", Type: "shape", StartTime: 25, Duration: 5},
				},
			},
			{
				Name:     "Outro",
				Width:    1920,
				Height:   1080,
				Duration: 3,
				Layers: []LayerConfig{
					{Name: "End Card", Type: "solid", StartTime: 0, Duration: 3},
					{Name: "CTA", Type: "text", SourceText: "Subscribe", StartTime: 0.5, Duration: 2.5},
				},
			},
		},
	}
}