package aep

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf8"
)

// TextDocument represents parsed text layer content
type TextDocument struct {
	Text           string
	FontName       string
	FontSize       float64
	FontStyle      string
	FillColor      [4]float32 // RGBA
	StrokeColor    [4]float32 // RGBA
	StrokeWidth    float32
	Justification  string
	Tracking       float32
	LineHeight     float32
	BaselineShift  float32
}

// TextKeyframe represents text at a specific time
type TextKeyframe struct {
	Time     float64
	Document TextDocument
}

// ExtractTextContent extracts text content from a text layer's property tree
func ExtractTextContent(textProperty *Property) (*TextDocument, error) {
	if textProperty == nil {
		return nil, fmt.Errorf("text property is nil")
	}
	
	doc := &TextDocument{
		// Default values
		FontName: "Arial",
		FontSize: 12.0,
		FillColor: [4]float32{1, 1, 1, 1}, // White
		StrokeColor: [4]float32{0, 0, 0, 1}, // Black
	}
	
	// Navigate through the property tree to find text data
	if err := extractFromPropertyTree(textProperty, doc); err != nil {
		return nil, err
	}
	
	// If we still don't have text, check the ADBE Text Document child directly
	if doc.Text == "" {
		for _, child := range textProperty.Properties {
			if child.MatchName == "ADBE Text Document" {
				// Try enhanced extraction on this specific property
				if len(child.RawData) > 0 {
					text := parseTextDocumentBinaryEnhanced(child.RawData)
					if text != "" {
						doc.Text = text
						break
					}
				}
				
				// If still no text, check deeper in the property tree
				if err := extractFromPropertyTree(child, doc); err == nil && doc.Text != "" {
					break
				}
				
				// Try enhanced extraction on the entire property tree
				if text := extractTextFromRawDataEnhanced(child); text != "" {
					doc.Text = text
					break
				}
				
				// Last resort - indicate where text should be
				if doc.Text == "" {
					doc.Text = "[Text content not extracted - check keyframes/expressions]"
				}
				break
			}
		}
	}
	
	return doc, nil
}

// extractFromPropertyTree recursively searches for text data in the property tree
func extractFromPropertyTree(prop *Property, doc *TextDocument) error {
	if prop == nil {
		return nil
	}
	
	// Check common text property match names
	switch prop.MatchName {
	case "ADBE Text Document":
		// This usually contains the actual text
		if err := parseTextDocument(prop, doc); err != nil {
			return err
		}
		
	case "ADBE Source Text":
		// Alternative location for text content
		if err := parseSourceText(prop, doc); err != nil {
			return err
		}
		
	case "ADBE Text Properties":
		// Root text properties, search children
		for _, child := range prop.Properties {
			if err := extractFromPropertyTree(child, doc); err != nil {
				return err
			}
		}
		
	case "ADBE Text Animators":
		// Text animators may contain text modifications
		// TODO: Handle animated text
		
	default:
		// Search in child properties
		for _, child := range prop.Properties {
			if err := extractFromPropertyTree(child, doc); err != nil {
				return err
			}
		}
	}
	
	// Also check property name for clues
	if strings.Contains(strings.ToLower(prop.Name), "text") ||
	   strings.Contains(strings.ToLower(prop.Name), "source") {
		// Try to extract text from this property
		parsePropertyData(prop, doc)
	}
	
	return nil
}

// parseTextDocument parses ADBE Text Document property
func parseTextDocument(prop *Property, doc *TextDocument) error {
	// Text documents in AEP files are typically stored as:
	// 1. String data in property labels or names
	// 2. Binary data in keyframe values
	// 3. Expression strings
	
	// Check if text is in the label (user-entered text)
	if prop.Label != "" && prop.Label != "-_0_/-" {
		doc.Text = prop.Label
		return nil
	}
	
	// Enhanced: Check raw data for binary text document
	if len(prop.RawData) > 0 {
		text := parseTextDocumentBinaryEnhanced(prop.RawData)
		if text != "" {
			doc.Text = text
			return nil
		}
	}
	
	// Look for keyframe data in child properties
	for _, child := range prop.Properties {
		if strings.Contains(child.MatchName, "Keyframe") || 
		   strings.Contains(child.MatchName, "Key") {
			// This might contain keyframe data
			if len(child.RawData) > 0 {
				// Try to parse keyframe text
				text := extractTextFromBytes(child.RawData)
				if text != "" && doc.Text == "" {
					doc.Text = text
					return nil
				}
			}
		}
	}
	
	return nil
}

// parseSourceText parses ADBE Source Text property
func parseSourceText(prop *Property, doc *TextDocument) error {
	// Source text might be in various formats
	
	// Check direct property values
	if prop.Name != "" && prop.Name != prop.MatchName && !strings.HasPrefix(prop.Name, "ADBE") {
		// Sometimes the text is stored directly in the name
		doc.Text = prop.Name
		return nil
	}
	
	// Check select options (for expression-based text)
	if len(prop.SelectOptions) > 0 {
		doc.Text = strings.Join(prop.SelectOptions, " ")
		return nil
	}
	
	// Enhanced: Check raw data
	if len(prop.RawData) > 0 {
		text := extractTextFromBytes(prop.RawData)
		if text != "" {
			doc.Text = text
			return nil
		}
	}
	
	return nil
}

// parsePropertyData attempts to extract text from generic property data
func parsePropertyData(prop *Property, doc *TextDocument) {
	// Try various heuristics to find text
	
	// Check if property name contains actual text
	if prop.Name != "" && prop.Name != prop.MatchName {
		// Filter out property names that are clearly not text content
		if !strings.HasPrefix(prop.Name, "ADBE") && 
		   !strings.Contains(prop.Name, "::") &&
		   len(prop.Name) > 3 {
			// This might be actual text content
			if doc.Text == "" {
				doc.Text = prop.Name
			}
		}
	}
	
	// Check label
	if prop.Label != "" && prop.Label != "-_0_/-" {
		if doc.Text == "" {
			doc.Text = prop.Label
		}
	}
	
	// Check raw data
	if len(prop.RawData) > 0 && doc.Text == "" {
		// Try to parse as UTF-8 string
		text := string(prop.RawData)
		// Clean up the text
		text = strings.TrimSpace(text)
		if text != "" && !strings.HasPrefix(text, "ADBE") {
			doc.Text = text
		}
	}
}

// ExtractEnhancedTextContent provides enhanced text extraction with binary parsing support
func ExtractEnhancedTextContent(textProperty *Property) (*TextDocument, error) {
	// This is an alias for ExtractTextContent which now includes enhanced extraction
	return ExtractTextContent(textProperty)
}

// ExtractAllTextLayers finds all text layers in a project and extracts their content
func ExtractAllTextLayers(project *Project) map[string][]*TextDocument {
	textByComp := make(map[string][]*TextDocument)
	
	for _, item := range project.Items {
		if item.ItemType == ItemTypeComposition {
			compTexts := make([]*TextDocument, 0)
			
			for _, layer := range item.CompositionLayers {
				if layer.Text != nil {
					// This is a text layer
					doc, err := ExtractTextContent(layer.Text)
					if err == nil && doc != nil {
						// Try to get better text if current is empty
						if doc.Text == "" {
							// Use layer name as fallback
							if strings.Contains(strings.ToLower(layer.Name), "text") ||
							   strings.Contains(strings.ToLower(layer.Name), "title") ||
							   strings.Contains(strings.ToLower(layer.Name), "placeholder") {
								doc.Text = fmt.Sprintf("[%s]", layer.Name)
							}
						}
						compTexts = append(compTexts, doc)
					}
				}
			}
			
			if len(compTexts) > 0 {
				textByComp[item.Name] = compTexts
			}
		}
	}
	
	return textByComp
}

// Helper function to parse binary font data (placeholder for future implementation)
func parseFontData(data []byte, doc *TextDocument) error {
	// TODO: Implement binary font data parsing
	// AEP files store font information in a binary format that needs reverse engineering
	
	// Basic structure (guessed):
	// - Font name as UTF-8 string
	// - Font size as float32
	// - Style flags as uint32
	// - Colors as float32 arrays
	
	if len(data) < 4 {
		return fmt.Errorf("insufficient data for font parsing")
	}
	
	// This is a placeholder - actual implementation would need to reverse engineer
	// the binary format through experimentation
	
	return nil
}

// ParseTextExpression parses text that comes from expressions
func ParseTextExpression(expression string) string {
	// Simple expression parsing
	// Real expressions would need JavaScript evaluation
	
	// Remove common expression syntax
	expression = strings.TrimSpace(expression)
	expression = strings.TrimPrefix(expression, "\"")
	expression = strings.TrimSuffix(expression, "\"")
	expression = strings.TrimPrefix(expression, "'")
	expression = strings.TrimSuffix(expression, "'")
	
	return expression
}

// Enhanced extraction methods below

// extractTextFromRawDataEnhanced attempts to extract text from raw data using multiple strategies
func extractTextFromRawDataEnhanced(prop *Property) string {
	// Check the property's own raw data
	if len(prop.RawData) > 0 {
		if text := parseTextDocumentBinaryEnhanced(prop.RawData); text != "" {
			return text
		}
	}
	
	// Check all child properties for raw data
	for _, child := range prop.Properties {
		if text := extractTextFromRawDataEnhanced(child); text != "" {
			return text
		}
	}
	
	return ""
}

// parseTextDocumentBinaryEnhanced attempts to parse the binary format with multiple strategies
func parseTextDocumentBinaryEnhanced(data []byte) string {
	if len(data) < 8 {
		return ""
	}
	
	// Try different parsing strategies
	
	// Strategy 1: Look for UTF-16 encoded text (common in Windows AEP files)
	if text := extractUTF16TextEnhanced(data); text != "" {
		return text
	}
	
	// Strategy 2: Look for UTF-8 text blocks
	if text := extractUTF8TextEnhanced(data); text != "" {
		return text
	}
	
	// Strategy 3: Look for null-terminated strings
	if text := extractNullTerminatedStringEnhanced(data); text != "" {
		return text
	}
	
	// Strategy 4: Try to find text after specific markers
	if text := extractTextAfterMarkersEnhanced(data); text != "" {
		return text
	}
	
	return ""
}

// extractUTF16TextEnhanced attempts to extract UTF-16 encoded text
func extractUTF16TextEnhanced(data []byte) string {
	// Look for UTF-16 BOM or patterns
	for i := 0; i < len(data)-4; i++ {
		// Check for possible UTF-16 LE text (every other byte is 0 for ASCII)
		if i+20 < len(data) && data[i+1] == 0 && data[i+3] == 0 && data[i+5] == 0 {
			// Might be UTF-16 LE
			end := i
			for j := i; j < len(data)-1; j += 2 {
				if data[j] == 0 && data[j+1] == 0 {
					end = j
					break
				}
			}
			
			if end > i+4 {
				// Try to decode UTF-16 LE
				textBytes := data[i:end]
				text := decodeUTF16LEEnhanced(textBytes)
				if isValidTextEnhanced(text) {
					return text
				}
			}
		}
	}
	
	return ""
}

// extractUTF8TextEnhanced attempts to extract UTF-8 encoded text
func extractUTF8TextEnhanced(data []byte) string {
	// Look for consecutive valid UTF-8 characters
	for i := 0; i < len(data); i++ {
		if utf8.Valid(data[i:]) {
			// Find the longest valid UTF-8 sequence
			for j := len(data); j > i+4; j-- {
				if utf8.Valid(data[i:j]) {
					text := string(data[i:j])
					text = cleanTextStringEnhanced(text)
					if isValidTextEnhanced(text) {
						return text
					}
				}
			}
		}
	}
	
	return ""
}

// extractNullTerminatedStringEnhanced extracts null-terminated strings
func extractNullTerminatedStringEnhanced(data []byte) string {
	for i := 0; i < len(data)-4; i++ {
		// Look for printable ASCII followed by null
		if data[i] >= 32 && data[i] <= 126 {
			end := bytes.IndexByte(data[i:], 0)
			if end > 4 && i+end <= len(data) {
				text := string(data[i : i+end])
				if isValidTextEnhanced(text) {
					return text
				}
			}
		}
	}
	
	return ""
}

// extractTextAfterMarkersEnhanced looks for text after known markers
func extractTextAfterMarkersEnhanced(data []byte) string {
	// Common markers that might precede text in AEP files
	markers := [][]byte{
		[]byte("TEXT"),
		[]byte("text"),
		[]byte("Utf8"),
		[]byte("utf8"),
		[]byte("tdbs"), // Text document binary structure
		[]byte{0x00, 0x00, 0x00}, // Common padding before text
	}
	
	for _, marker := range markers {
		idx := bytes.Index(data, marker)
		if idx >= 0 && idx+len(marker)+4 < len(data) {
			// Skip marker and some potential header bytes
			start := idx + len(marker)
			
			// Skip any leading nulls or small values
			for start < len(data) && data[start] < 32 {
				start++
			}
			
			if start < len(data)-4 {
				// Try to extract text from this position
				text := extractTextFromBytesEnhanced(data[start:])
				if text != "" {
					return text
				}
			}
		}
	}
	
	return ""
}

// extractTextFromBytesEnhanced is a general-purpose text extractor
func extractTextFromBytesEnhanced(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	
	// First, try direct UTF-8 conversion
	if utf8.Valid(data) {
		text := cleanTextStringEnhanced(string(data))
		if isValidTextEnhanced(text) {
			return text
		}
	}
	
	// Try to find sequences of printable characters
	var result []byte
	inText := false
	
	for i := 0; i < len(data); i++ {
		if isPrintableEnhanced(data[i]) {
			if !inText {
				inText = true
			}
			result = append(result, data[i])
		} else if inText && len(result) > 4 {
			// End of text sequence
			text := cleanTextStringEnhanced(string(result))
			if isValidTextEnhanced(text) {
				return text
			}
			result = nil
			inText = false
		} else if inText {
			result = nil
			inText = false
		}
	}
	
	// Check final result
	if len(result) > 4 {
		text := cleanTextStringEnhanced(string(result))
		if isValidTextEnhanced(text) {
			return text
		}
	}
	
	return ""
}

// Helper functions for enhanced extraction

func decodeUTF16LEEnhanced(data []byte) string {
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}
	
	runes := make([]rune, 0, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		r := rune(binary.LittleEndian.Uint16(data[i : i+2]))
		if r > 0 && r < 0xFFFD { // Valid Unicode range
			runes = append(runes, r)
		}
	}
	
	return string(runes)
}

func isPrintableEnhanced(b byte) bool {
	return (b >= 32 && b <= 126) || b == '\n' || b == '\r' || b == '\t'
}

func cleanTextStringEnhanced(s string) string {
	// Remove null characters and trim
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.TrimSpace(s)
	
	// Remove any non-printable characters except newlines and tabs
	var result []rune
	for _, r := range s {
		if r >= 32 || r == '\n' || r == '\r' || r == '\t' {
			result = append(result, r)
		}
	}
	
	return string(result)
}

func isValidTextEnhanced(s string) bool {
	if len(s) < 1 {
		return false
	}
	
	// Check if it's likely to be actual text content
	// Reject strings that look like property names or technical identifiers
	if strings.HasPrefix(s, "ADBE") ||
		strings.HasPrefix(s, "tdbs") ||
		strings.HasPrefix(s, "pard") ||
		strings.Contains(s, "::") ||
		strings.HasPrefix(s, "-_") ||
		strings.HasPrefix(s, "@@") {
		return false
	}
	
	// Count printable characters
	printableCount := 0
	for _, r := range s {
		if r >= 32 && r <= 126 {
			printableCount++
		}
	}
	
	// At least 50% should be printable
	return float64(printableCount)/float64(len(s)) > 0.5
}

// parseTextDocumentBinary attempts to parse the binary format of text documents
func parseTextDocumentBinary(data []byte) string {
	if len(data) < 8 {
		return ""
	}
	
	// Try different parsing strategies
	
	// Strategy 1: Look for UTF-16 encoded text (common in Windows AEP files)
	if text := extractUTF16Text(data); text != "" {
		return text
	}
	
	// Strategy 2: Look for UTF-8 text blocks
	if text := extractUTF8Text(data); text != "" {
		return text
	}
	
	// Strategy 3: Look for null-terminated strings
	if text := extractNullTerminatedString(data); text != "" {
		return text
	}
	
	// Strategy 4: Try to find text after specific markers
	if text := extractTextAfterMarkers(data); text != "" {
		return text
	}
	
	return ""
}

// extractTextFromBytes is a general-purpose text extractor
func extractTextFromBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	
	// First, try direct UTF-8 conversion
	if utf8.Valid(data) {
		text := cleanTextString(string(data))
		if isValidText(text) {
			return text
		}
	}
	
	// Try to find sequences of printable characters
	var result []byte
	inText := false
	
	for i := 0; i < len(data); i++ {
		if isPrintable(data[i]) {
			if !inText {
				inText = true
			}
			result = append(result, data[i])
		} else if inText && len(result) > 4 {
			// End of text sequence
			text := cleanTextString(string(result))
			if isValidText(text) {
				return text
			}
			result = nil
			inText = false
		} else if inText {
			result = nil
			inText = false
		}
	}
	
	// Check final result
	if len(result) > 4 {
		text := cleanTextString(string(result))
		if isValidText(text) {
			return text
		}
	}
	
	return ""
}

// extractUTF16Text attempts to extract UTF-16 encoded text
func extractUTF16Text(data []byte) string {
	// Look for UTF-16 BOM or patterns
	for i := 0; i < len(data)-4; i++ {
		// Check for possible UTF-16 LE text (every other byte is 0 for ASCII)
		if i+20 < len(data) && data[i+1] == 0 && data[i+3] == 0 && data[i+5] == 0 {
			// Might be UTF-16 LE
			end := i
			for j := i; j < len(data)-1; j += 2 {
				if data[j] == 0 && data[j+1] == 0 {
					end = j
					break
				}
			}
			
			if end > i+4 {
				// Try to decode UTF-16 LE
				textBytes := data[i:end]
				text := decodeUTF16LE(textBytes)
				if isValidText(text) {
					return text
				}
			}
		}
	}
	
	return ""
}

// extractUTF8Text attempts to extract UTF-8 encoded text
func extractUTF8Text(data []byte) string {
	// Look for consecutive valid UTF-8 characters
	for i := 0; i < len(data); i++ {
		if utf8.Valid(data[i:]) {
			// Find the longest valid UTF-8 sequence
			for j := len(data); j > i+4; j-- {
				if utf8.Valid(data[i:j]) {
					text := string(data[i:j])
					text = cleanTextString(text)
					if isValidText(text) {
						return text
					}
				}
			}
		}
	}
	
	return ""
}

// extractNullTerminatedString extracts null-terminated strings
func extractNullTerminatedString(data []byte) string {
	for i := 0; i < len(data)-4; i++ {
		// Look for printable ASCII followed by null
		if data[i] >= 32 && data[i] <= 126 {
			end := bytes.IndexByte(data[i:], 0)
			if end > 4 {
				text := string(data[i : i+end])
				if isValidText(text) {
					return text
				}
			}
		}
	}
	
	return ""
}

// extractTextAfterMarkers looks for text after known markers
func extractTextAfterMarkers(data []byte) string {
	// Common markers that might precede text in AEP files
	markers := [][]byte{
		[]byte("TEXT"),
		[]byte("text"),
		[]byte("Utf8"),
		[]byte("utf8"),
		[]byte{0x00, 0x00, 0x00}, // Common padding before text
	}
	
	for _, marker := range markers {
		idx := bytes.Index(data, marker)
		if idx >= 0 && idx+len(marker)+4 < len(data) {
			// Skip marker and some potential header bytes
			start := idx + len(marker)
			
			// Skip any leading nulls or small values
			for start < len(data) && data[start] < 32 {
				start++
			}
			
			if start < len(data)-4 {
				text := extractTextFromBytes(data[start:])
				if text != "" {
					return text
				}
			}
		}
	}
	
	return ""
}

// Helper functions

func decodeUTF16LE(data []byte) string {
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}
	
	runes := make([]rune, 0, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		r := rune(binary.LittleEndian.Uint16(data[i : i+2]))
		if r > 0 {
			runes = append(runes, r)
		}
	}
	
	return string(runes)
}

func isPrintable(b byte) bool {
	return (b >= 32 && b <= 126) || b == '\n' || b == '\r' || b == '\t'
}

func cleanTextString(s string) string {
	// Remove null characters and trim
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.TrimSpace(s)
	
	// Remove any non-printable characters except newlines and tabs
	var result []rune
	for _, r := range s {
		if r >= 32 || r == '\n' || r == '\r' || r == '\t' {
			result = append(result, r)
		}
	}
	
	return string(result)
}

func isValidText(s string) bool {
	if len(s) < 1 {
		return false
	}
	
	// Check if it's likely to be actual text content
	// Reject strings that look like property names or technical identifiers
	if strings.HasPrefix(s, "ADBE") ||
		strings.HasPrefix(s, "tdbs") ||
		strings.HasPrefix(s, "pard") ||
		strings.Contains(s, "::") ||
		strings.HasPrefix(s, "-_") {
		return false
	}
	
	// Count printable characters
	printableCount := 0
	for _, r := range s {
		if r >= 32 && r <= 126 {
			printableCount++
		}
	}
	
	// At least 50% should be printable
	return float64(printableCount)/float64(len(s)) > 0.5
}