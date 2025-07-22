package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	
	aep "github.com/yourusername/mobot2025"
)

// StoryElement represents a text element in the narrative order
type StoryElement struct {
	Order        int     `json:"order"`
	Text         string  `json:"text"`
	LayerName    string  `json:"layerName"`
	CompName     string  `json:"compName"`
	TimeStart    float64 `json:"timeStart"`
	TimeEnd      float64 `json:"timeEnd"`
	Duration     float64 `json:"duration"`
	IsPlaceholder bool   `json:"isPlaceholder"`
}

// StoryData holds the complete narrative structure
type StoryData struct {
	ProjectName  string         `json:"projectName"`
	TotalScenes  int            `json:"totalScenes"`
	TotalDuration float64       `json:"totalDuration"`
	Elements     []StoryElement `json:"elements"`
	ExtractedAt  string         `json:"extractedAt"`
}

const easyModeHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Easy Mode - Story Viewer</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
            min-height: 100vh;
            padding: 20px;
            line-height: 1.6;
        }
        
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 300;
            letter-spacing: -1px;
        }
        
        .header p {
            opacity: 0.9;
            font-size: 1.1em;
        }
        
        .upload-section {
            padding: 40px;
            text-align: center;
            border-bottom: 1px solid #eee;
            background: #f8f9fa;
        }
        
        .upload-box {
            border: 3px dashed #667eea;
            border-radius: 15px;
            padding: 60px 40px;
            cursor: pointer;
            transition: all 0.3s ease;
            background: white;
        }
        
        .upload-box:hover {
            background: #f0f4ff;
            border-color: #764ba2;
            transform: translateY(-2px);
        }
        
        .upload-icon {
            font-size: 4em;
            margin-bottom: 20px;
        }
        
        .upload-text {
            font-size: 1.3em;
            color: #333;
            margin-bottom: 10px;
        }
        
        .upload-subtext {
            color: #666;
            font-size: 0.95em;
        }
        
        input[type="file"] {
            display: none;
        }
        
        .story-section {
            padding: 40px;
            display: none;
        }
        
        .story-header {
            text-align: center;
            margin-bottom: 40px;
        }
        
        .story-title {
            font-size: 2em;
            color: #333;
            margin-bottom: 10px;
        }
        
        .story-meta {
            color: #666;
            font-size: 1.1em;
        }
        
        .story-timeline {
            margin: 40px 0;
            position: relative;
            padding-left: 30px;
        }
        
        .timeline-line {
            position: absolute;
            left: 10px;
            top: 0;
            bottom: 0;
            width: 3px;
            background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
        }
        
        .story-element {
            margin-bottom: 30px;
            padding: 25px;
            background: #f8f9fa;
            border-radius: 15px;
            position: relative;
            transition: all 0.3s ease;
            cursor: pointer;
        }
        
        .story-element:hover {
            background: #e9ecef;
            transform: translateX(5px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .story-element::before {
            content: '';
            position: absolute;
            left: -23px;
            top: 30px;
            width: 16px;
            height: 16px;
            background: #667eea;
            border-radius: 50%;
            border: 3px solid white;
            box-shadow: 0 2px 5px rgba(0,0,0,0.2);
        }
        
        .element-number {
            position: absolute;
            top: 10px;
            right: 20px;
            background: #667eea;
            color: white;
            width: 30px;
            height: 30px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            font-size: 0.9em;
        }
        
        .element-text {
            font-size: 1.4em;
            color: #333;
            margin-bottom: 15px;
            line-height: 1.5;
        }
        
        .element-meta {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            font-size: 0.9em;
            color: #666;
        }
        
        .meta-item {
            display: flex;
            align-items: center;
            gap: 5px;
        }
        
        .meta-icon {
            opacity: 0.7;
        }
        
        .placeholder-badge {
            background: #ffc107;
            color: #000;
            padding: 3px 10px;
            border-radius: 20px;
            font-size: 0.8em;
            font-weight: 500;
        }
        
        .chat-section {
            padding: 40px;
            background: #f8f9fa;
            border-top: 1px solid #eee;
        }
        
        .chat-header {
            text-align: center;
            margin-bottom: 30px;
        }
        
        .chat-title {
            font-size: 1.5em;
            color: #333;
            margin-bottom: 10px;
        }
        
        .chat-box {
            background: white;
            border-radius: 15px;
            padding: 20px;
            min-height: 150px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.05);
        }
        
        .chat-input {
            width: 100%;
            padding: 15px;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            font-size: 1em;
            margin-bottom: 15px;
            transition: border-color 0.3s ease;
        }
        
        .chat-input:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .chat-buttons {
            display: flex;
            gap: 10px;
            justify-content: center;
        }
        
        .btn {
            padding: 12px 25px;
            border: none;
            border-radius: 25px;
            font-size: 1em;
            cursor: pointer;
            transition: all 0.3s ease;
            font-weight: 500;
        }
        
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        
        .btn-secondary {
            background: #e9ecef;
            color: #333;
        }
        
        .btn-secondary:hover {
            background: #dee2e6;
        }
        
        .json-output {
            background: #2d3748;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 10px;
            margin-top: 20px;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 0.9em;
            line-height: 1.5;
            overflow-x: auto;
            display: none;
        }
        
        .loading {
            display: none;
            text-align: center;
            padding: 40px;
        }
        
        .loading-spinner {
            width: 50px;
            height: 50px;
            margin: 0 auto 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #666;
        }
        
        .empty-icon {
            font-size: 5em;
            opacity: 0.3;
            margin-bottom: 20px;
        }
        
        .empty-text {
            font-size: 1.2em;
        }
        
        @media (max-width: 768px) {
            .container {
                margin: 0;
                border-radius: 0;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .story-element {
                padding: 20px;
            }
            
            .element-text {
                font-size: 1.2em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📖 Easy Mode Story Viewer</h1>
            <p>See your After Effects project as a story</p>
        </div>
        
        <div class="upload-section">
            <label for="file-input" class="upload-box">
                <div class="upload-icon">📤</div>
                <div class="upload-text">Drop your AEP file here or click to upload</div>
                <div class="upload-subtext">We'll extract all the text and show it in story order</div>
            </label>
            <input type="file" id="file-input" accept=".aep">
        </div>
        
        <div class="loading">
            <div class="loading-spinner"></div>
            <p>Extracting your story...</p>
        </div>
        
        <div class="story-section">
            <div class="story-header">
                <h2 class="story-title" id="story-title">Your Story</h2>
                <p class="story-meta" id="story-meta"></p>
            </div>
            
            <div class="story-timeline" id="story-timeline">
                <div class="timeline-line"></div>
                <!-- Story elements will be inserted here -->
            </div>
            
            <div class="empty-state" id="empty-state" style="display: none;">
                <div class="empty-icon">📭</div>
                <p class="empty-text">No text found in this project</p>
            </div>
        </div>
        
        <div class="chat-section" style="display: none;" id="chat-section">
            <div class="chat-header">
                <h3 class="chat-title">🤖 Chat with AI about your story</h3>
                <p>Ask questions or request changes to the narrative</p>
            </div>
            
            <div class="chat-box">
                <textarea class="chat-input" id="chat-input" placeholder="E.g., 'Make the story more exciting' or 'Change all placeholder text to be about cats'" rows="4"></textarea>
                <div class="chat-buttons">
                    <button class="btn btn-primary" onclick="processChat()">Generate JSON</button>
                    <button class="btn btn-secondary" onclick="copyJSON()">Copy JSON</button>
                    <button class="btn btn-secondary" onclick="downloadJSON()">Download JSON</button>
                </div>
                <pre class="json-output" id="json-output"></pre>
            </div>
        </div>
    </div>
    
    <script>
        let storyData = null;
        
        document.getElementById('file-input').addEventListener('change', handleFileUpload);
        
        function handleFileUpload(e) {
            const file = e.target.files[0];
            if (!file) return;
            
            const formData = new FormData();
            formData.append('file', file);
            
            document.querySelector('.upload-section').style.display = 'none';
            document.querySelector('.loading').style.display = 'block';
            
            // Upload to server
            fetch('/upload', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                storyData = data;
                displayStory();
            })
            .catch(error => {
                console.error('Upload error:', error);
                alert('Failed to process file. Please try again.');
                location.reload();
            });
        }
        
        function displayStory() {
            document.querySelector('.loading').style.display = 'none';
            document.querySelector('.story-section').style.display = 'block';
            document.querySelector('#chat-section').style.display = 'block';
            
            document.getElementById('story-title').textContent = storyData.projectName;
            document.getElementById('story-meta').textContent = storyData.totalScenes + ' scenes • ' + storyData.totalDuration + 's total • ' + storyData.elements.length + ' text elements';
            
            const timeline = document.getElementById('story-timeline');
            const existingElements = timeline.querySelectorAll('.story-element');
            existingElements.forEach(el => el.remove());
            
            if (storyData.elements.length === 0) {
                document.getElementById('empty-state').style.display = 'block';
                return;
            }
            
            storyData.elements.forEach((element, index) => {
                const elementDiv = document.createElement('div');
                elementDiv.className = 'story-element';
                elementDiv.innerHTML = '<div class="element-number">' + element.order + '</div>' +
                    '<div class="element-text">' + element.text + '</div>' +
                    '<div class="element-meta">' +
                        '<div class="meta-item">' +
                            '<span class="meta-icon">📍</span>' +
                            '<span>' + element.compName + '</span>' +
                        '</div>' +
                        '<div class="meta-item">' +
                            '<span class="meta-icon">📝</span>' +
                            '<span>' + element.layerName + '</span>' +
                        '</div>' +
                        '<div class="meta-item">' +
                            '<span class="meta-icon">⏱️</span>' +
                            '<span>' + element.timeStart + 's - ' + element.timeEnd + 's</span>' +
                        '</div>' +
                        (element.isPlaceholder ? '<span class="placeholder-badge">Placeholder</span>' : '') +
                    '</div>';
                
                elementDiv.addEventListener('click', () => {
                    selectElement(index);
                });
                
                timeline.appendChild(elementDiv);
            });
        }
        
        function selectElement(index) {
            const elements = document.querySelectorAll('.story-element');
            elements.forEach((el, i) => {
                if (i === index) {
                    el.style.background = '#e7f3ff';
                } else {
                    el.style.background = '#f8f9fa';
                }
            });
        }
        
        function processChat() {
            const userInput = document.getElementById('chat-input').value;
            if (!userInput.trim() || !storyData) return;
            
            // Show loading state
            const jsonOutput = document.getElementById('json-output');
            jsonOutput.textContent = 'Processing your request...';
            jsonOutput.style.display = 'block';
            
            // Send to NLP endpoint
            fetch('/nlp', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    story: storyData,
                    userRequest: userInput
                })
            })
            .then(response => response.json())
            .then(data => {
                // Update story data with AI-processed version
                storyData = data.story;
                
                // Display the full response with changes
                jsonOutput.textContent = JSON.stringify(data, null, 2);
                
                // Optionally refresh the story display
                if (data.changes && data.changes.length > 0) {
                    displayStory();
                    
                    // Show what changed
                    const changesSummary = data.changes.join('\\n• ');
                    alert('AI Processing Complete!\\n\\nChanges made:\\n• ' + changesSummary);
                }
            })
            .catch(error => {
                console.error('NLP error:', error);
                jsonOutput.textContent = 'Error processing request. Please try again.';
            });
        }
        
        function copyJSON() {
            const jsonText = document.getElementById('json-output').textContent;
            if (!jsonText) return;
            
            navigator.clipboard.writeText(jsonText).then(() => {
                alert('JSON copied to clipboard!');
            });
        }
        
        function downloadJSON() {
            const jsonText = document.getElementById('json-output').textContent;
            if (!jsonText) return;
            
            const blob = new Blob([jsonText], { type: 'application/json' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = storyData.projectName + '-story.json';
            a.click();
            URL.revokeObjectURL(url);
        }
        
        // Drag and drop support
        const uploadBox = document.querySelector('.upload-box');
        
        uploadBox.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadBox.style.background = '#f0f4ff';
        });
        
        uploadBox.addEventListener('dragleave', () => {
            uploadBox.style.background = 'white';
        });
        
        uploadBox.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadBox.style.background = 'white';
            
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                document.getElementById('file-input').files = files;
                handleFileUpload({ target: { files: files } });
            }
        });
    </script>
</body>
</html>`

// StoryServer handles the web interface
type StoryServer struct {
	Port string
}

// ExtractStoryFromAEP extracts text in narrative order
func ExtractStoryFromAEP(filename string, project *aep.Project) StoryData {
	story := StoryData{
		ProjectName:  filepath.Base(filename),
		ExtractedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Elements:     []StoryElement{},
	}
	
	// Track compositions by their start time
	type CompInfo struct {
		Name      string
		StartTime float64
		Duration  float64
		Layers    []struct {
			Name      string
			Text      string
			InPoint   float64
			OutPoint  float64
		}
	}
	
	compositions := []CompInfo{}
	
	// Process all compositions
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			comp := CompInfo{
				Name:      item.Name,
				Duration:  item.FootageSeconds,
				StartTime: 0, // Would need timeline analysis
			}
			
			// Process layers in this comp
			for _, layer := range item.CompositionLayers {
				// Look for text layers
				if layer.Text != nil {
					text := layer.Name // In many templates, the placeholder text IS the layer name
					
					// Try to extract actual text if available
					if layer.Text.TextDocument != nil && layer.Text.TextDocument.Text != "" && 
					   !strings.Contains(layer.Text.TextDocument.Text, "[") {
						text = layer.Text.TextDocument.Text
					}
					
					comp.Layers = append(comp.Layers, struct {
						Name      string
						Text      string
						InPoint   float64
						OutPoint  float64
					}{
						Name:     layer.Name,
						Text:     text,
						InPoint:  0.0, // Layer timing would need to be extracted from timeline data
						OutPoint: comp.Duration,
					})
				}
			}
			
			if len(comp.Layers) > 0 {
				compositions = append(compositions, comp)
				story.TotalScenes++
				story.TotalDuration += comp.Duration
			}
		}
	}
	
	// Sort compositions by name (approximation of timeline order)
	sort.Slice(compositions, func(i, j int) bool {
		// Try to extract scene numbers
		iNum := extractSceneNumber(compositions[i].Name)
		jNum := extractSceneNumber(compositions[j].Name)
		if iNum != jNum {
			return iNum < jNum
		}
		return compositions[i].Name < compositions[j].Name
	})
	
	// Build story elements in order
	order := 1
	currentTime := 0.0
	
	for _, comp := range compositions {
		// Sort layers by in-point
		sort.Slice(comp.Layers, func(i, j int) bool {
			return comp.Layers[i].InPoint < comp.Layers[j].InPoint
		})
		
		for _, layer := range comp.Layers {
			element := StoryElement{
				Order:        order,
				Text:         layer.Text,
				LayerName:    layer.Name,
				CompName:     comp.Name,
				TimeStart:    currentTime + layer.InPoint,
				TimeEnd:      currentTime + layer.OutPoint,
				Duration:     layer.OutPoint - layer.InPoint,
				IsPlaceholder: strings.Contains(layer.Text, "[") || strings.Contains(strings.ToLower(layer.Text), "placeholder"),
			}
			
			story.Elements = append(story.Elements, element)
			order++
		}
		
		currentTime += comp.Duration
	}
	
	return story
}

// extractSceneNumber tries to extract a scene number from comp name
func extractSceneNumber(name string) int {
	// Look for patterns like "Scene 1", "S01", etc.
	parts := strings.Fields(name)
	for i, part := range parts {
		if strings.ToLower(part) == "scene" && i+1 < len(parts) {
			var num int
			fmt.Sscanf(parts[i+1], "%d", &num)
			return num
		}
		if strings.HasPrefix(strings.ToUpper(part), "S") {
			var num int
			fmt.Sscanf(part[1:], "%d", &num)
			return num
		}
	}
	return 999 // Default for comps without scene numbers
}

// HandleUpload processes uploaded AEP files
func (s *StoryServer) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	
	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// Save temporary file
	tempFile, err := os.CreateTemp("", "upload-*.aep")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	
	// Copy uploaded file
	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	
	// Copy file content to temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	tempFile.Close()
	
	// Parse AEP file
	project, err := aep.Open(tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to parse AEP file", http.StatusBadRequest)
		return
	}
	
	// Extract story
	story := ExtractStoryFromAEP(header.Filename, project)
	
	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(story)
}

// HandleHome serves the main HTML interface
func (s *StoryServer) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(easyModeHTML))
}

// NLPRequest represents a chat request
type NLPRequest struct {
	Story       StoryData `json:"story"`
	UserRequest string    `json:"userRequest"`
}

// NLPResponse represents the processed story
type NLPResponse struct {
	Story       StoryData `json:"story"`
	UserRequest string    `json:"userRequest"`
	AIProcessed bool      `json:"aiProcessed"`
	Timestamp   string    `json:"timestamp"`
	Changes     []string  `json:"changes"`
}

// HandleNLP processes natural language requests
func (s *StoryServer) HandleNLP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req NLPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	// Simple NLP processing - in production this would use AI
	response := NLPResponse{
		Story:       req.Story,
		UserRequest: req.UserRequest,
		AIProcessed: true,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		Changes:     []string{},
	}
	
	// Process common requests
	lowerRequest := strings.ToLower(req.UserRequest)
	
	if strings.Contains(lowerRequest, "exciting") || strings.Contains(lowerRequest, "dramatic") {
		// Make text more exciting
		for i := range response.Story.Elements {
			if !response.Story.Elements[i].IsPlaceholder {
				response.Story.Elements[i].Text = strings.ToUpper(response.Story.Elements[i].Text) + "!"
			}
		}
		response.Changes = append(response.Changes, "Made text more dramatic with caps and exclamation marks")
	}
	
	if strings.Contains(lowerRequest, "cat") || strings.Contains(lowerRequest, "kitten") {
		// Cat-ify the placeholders
		catPhrases := []string{
			"[Meow meow meow]",
			"[Purr-fect content here]",
			"[Insert cat wisdom]",
			"[Feline features go here]",
			"[Paws for effect]",
		}
		catIndex := 0
		for i := range response.Story.Elements {
			if response.Story.Elements[i].IsPlaceholder {
				response.Story.Elements[i].Text = catPhrases[catIndex%len(catPhrases)]
				catIndex++
			}
		}
		response.Changes = append(response.Changes, "Replaced placeholders with cat-themed content")
	}
	
	if strings.Contains(lowerRequest, "short") || strings.Contains(lowerRequest, "brief") {
		// Shorten text
		for i := range response.Story.Elements {
			words := strings.Fields(response.Story.Elements[i].Text)
			if len(words) > 3 {
				response.Story.Elements[i].Text = strings.Join(words[:3], " ") + "..."
			}
		}
		response.Changes = append(response.Changes, "Shortened all text to 3 words maximum")
	}
	
	if strings.Contains(lowerRequest, "professional") || strings.Contains(lowerRequest, "business") {
		// Make more professional
		replacements := map[string]string{
			"amazing": "innovative",
			"awesome": "excellent",
			"cool": "professional",
			"great": "exceptional",
		}
		for i := range response.Story.Elements {
			text := response.Story.Elements[i].Text
			for old, new := range replacements {
				text = strings.ReplaceAll(strings.ToLower(text), old, new)
			}
			response.Story.Elements[i].Text = strings.Title(text)
		}
		response.Changes = append(response.Changes, "Applied professional tone and vocabulary")
	}
	
	// Return processed story
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Start starts the web server
func (s *StoryServer) Start() error {
	http.HandleFunc("/", s.HandleHome)
	http.HandleFunc("/upload", s.HandleUpload)
	http.HandleFunc("/nlp", s.HandleNLP)
	
	fmt.Printf("🚀 Easy Mode Story Viewer started on http://localhost:%s\n", s.Port)
	fmt.Println("📤 Upload an AEP file to see your story!")
	
	return http.ListenAndServe(":"+s.Port, nil)
}

func main() {
	server := &StoryServer{
		Port: "8080",
	}
	
	if len(os.Args) > 1 && os.Args[1] == "test" {
		// Test mode - process a file directly
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run easy_mode_story_viewer.go test <aep-file>")
			os.Exit(1)
		}
		
		aepFile := os.Args[2]
		project, err := aep.Open(aepFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		
		story := ExtractStoryFromAEP(aepFile, project)
		
		// Print story
		fmt.Printf("\n📖 Story from: %s\n", story.ProjectName)
		fmt.Printf("📊 %d scenes, %.1fs total\n\n", story.TotalScenes, story.TotalDuration)
		
		for _, element := range story.Elements {
			marker := "📝"
			if element.IsPlaceholder {
				marker = "🔤"
			}
			fmt.Printf("%s %d. %s\n", marker, element.Order, element.Text)
			fmt.Printf("   📍 %s > %s (%.1fs - %.1fs)\n\n", element.CompName, element.LayerName, element.TimeStart, element.TimeEnd)
		}
		
		// Output JSON
		jsonData, _ := json.MarshalIndent(story, "", "  ")
		fmt.Printf("\n📋 JSON Output:\n%s\n", string(jsonData))
		
	} else {
		// Web server mode
		if err := server.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	}
}