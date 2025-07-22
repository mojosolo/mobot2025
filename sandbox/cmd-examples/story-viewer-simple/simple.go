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
	
	aep "github.com/mojosolo/mobot2025"
)

// Scene represents a grouped scene with its text
type Scene struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	HeroText  string   `json:"heroText"`
	OtherText []string `json:"otherText"`
	HasPlaceholders bool `json:"hasPlaceholders"`
}

// SimpleStory holds the simplified story structure
type SimpleStory struct {
	ProjectName string   `json:"projectName"`
	Scenes      []Scene  `json:"scenes"`
	TotalText   int      `json:"totalText"`
}

const simpleHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Story Viewer</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, system-ui, sans-serif;
            background: #f5f5f5;
            color: #333;
            line-height: 1.6;
            padding: 20px;
        }
        
        .container {
            max-width: 600px;
            margin: 0 auto;
        }
        
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        
        .header h1 {
            font-size: 24px;
            font-weight: 400;
            margin-bottom: 10px;
        }
        
        .upload-area {
            background: white;
            border: 2px dashed #ddd;
            border-radius: 8px;
            padding: 40px;
            text-align: center;
            cursor: pointer;
            transition: all 0.2s;
            margin-bottom: 30px;
        }
        
        .upload-area:hover {
            border-color: #007AFF;
            background: #f9f9f9;
        }
        
        .upload-area.dragover {
            border-color: #007AFF;
            background: #E3F2FD;
        }
        
        .upload-btn {
            background: #007AFF;
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 6px;
            font-size: 16px;
            cursor: pointer;
        }
        
        .upload-btn:hover {
            background: #0056D3;
        }
        
        input[type="file"] {
            display: none;
        }
        
        .loading {
            display: none;
            text-align: center;
            padding: 40px;
        }
        
        .loading-spinner {
            width: 40px;
            height: 40px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #007AFF;
            border-radius: 50%;
            margin: 0 auto 20px;
            animation: spin 1s linear infinite;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .story-view {
            display: none;
        }
        
        .project-title {
            font-size: 20px;
            margin-bottom: 20px;
            text-align: center;
        }
        
        .scene-card {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 16px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .scene-card:hover {
            box-shadow: 0 4px 8px rgba(0,0,0,0.15);
        }
        
        .scene-card.expanded {
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        
        .scene-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;
        }
        
        .scene-number {
            background: #007AFF;
            color: white;
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 500;
            font-size: 14px;
        }
        
        .scene-title {
            flex: 1;
            margin: 0 12px;
            font-weight: 500;
        }
        
        .placeholder-badge {
            background: #FFB800;
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            display: none;
        }
        
        .scene-card.has-placeholders .placeholder-badge {
            display: inline-block;
        }
        
        .hero-text {
            font-size: 18px;
            margin-bottom: 12px;
            color: #000;
        }
        
        .more-text {
            display: none;
            margin-top: 12px;
            padding-top: 12px;
            border-top: 1px solid #eee;
        }
        
        .scene-card.expanded .more-text {
            display: block;
        }
        
        .other-text-item {
            padding: 8px 0;
            color: #666;
        }
        
        .more-indicator {
            color: #007AFF;
            font-size: 14px;
            text-align: center;
        }
        
        .actions {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin-top: 24px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        
        .actions-title {
            font-size: 16px;
            margin-bottom: 16px;
            font-weight: 500;
        }
        
        .action-buttons {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 12px;
        }
        
        .action-btn {
            background: #f0f0f0;
            border: none;
            padding: 12px;
            border-radius: 6px;
            font-size: 14px;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .action-btn:hover {
            background: #e0e0e0;
        }
        
        .action-btn.primary {
            background: #007AFF;
            color: white;
            grid-column: 1 / -1;
        }
        
        .action-btn.primary:hover {
            background: #0056D3;
        }
        
        .success-message {
            position: fixed;
            top: 20px;
            right: 20px;
            background: #4CAF50;
            color: white;
            padding: 12px 20px;
            border-radius: 6px;
            display: none;
            animation: slideIn 0.3s ease;
        }
        
        @keyframes slideIn {
            from {
                transform: translateX(100%);
                opacity: 0;
            }
            to {
                transform: translateX(0);
                opacity: 1;
            }
        }
        
        @media (max-width: 600px) {
            .action-buttons {
                grid-template-columns: 1fr;
            }
            
            .action-btn.primary {
                grid-column: 1;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Story Viewer</h1>
            <p>View your project narrative in seconds</p>
        </div>
        
        <div class="upload-area" id="uploadArea">
            <input type="file" id="fileInput" accept=".aep">
            <button class="upload-btn" onclick="document.getElementById('fileInput').click()">
                Open Your Project
            </button>
            <p style="margin-top: 16px; color: #666; font-size: 14px;">
                Or drag and drop your file here
            </p>
        </div>
        
        <div class="loading" id="loading">
            <div class="loading-spinner"></div>
            <p>Reading your story...</p>
        </div>
        
        <div class="story-view" id="storyView">
            <h2 class="project-title" id="projectTitle"></h2>
            <div id="scenesContainer"></div>
            
            <div class="actions">
                <h3 class="actions-title">Quick Actions</h3>
                <div class="action-buttons">
                    <button class="action-btn" onclick="makeAction('professional')">Make Professional</button>
                    <button class="action-btn" onclick="makeAction('casual')">Make Casual</button>
                    <button class="action-btn" onclick="makeAction('shorter')">Make Shorter</button>
                    <button class="action-btn" onclick="makeAction('placeholders')">Fill Placeholders</button>
                    <button class="action-btn primary" onclick="exportText()">Export Text</button>
                </div>
            </div>
        </div>
    </div>
    
    <div class="success-message" id="successMessage"></div>
    
    <script>
        let storyData = null;
        
        // File upload handling
        const fileInput = document.getElementById('fileInput');
        const uploadArea = document.getElementById('uploadArea');
        
        fileInput.addEventListener('change', handleFileSelect);
        
        uploadArea.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadArea.classList.add('dragover');
        });
        
        uploadArea.addEventListener('dragleave', () => {
            uploadArea.classList.remove('dragover');
        });
        
        uploadArea.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadArea.classList.remove('dragover');
            
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                handleFile(files[0]);
            }
        });
        
        function handleFileSelect(e) {
            const file = e.target.files[0];
            if (file) {
                handleFile(file);
            }
        }
        
        function handleFile(file) {
            const formData = new FormData();
            formData.append('file', file);
            
            document.getElementById('uploadArea').style.display = 'none';
            document.getElementById('loading').style.display = 'block';
            
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
                console.error('Error:', error);
                showMessage('Error loading file. Please try again.', 'error');
                resetUpload();
            });
        }
        
        function displayStory() {
            document.getElementById('loading').style.display = 'none';
            document.getElementById('storyView').style.display = 'block';
            
            document.getElementById('projectTitle').textContent = storyData.projectName;
            
            const container = document.getElementById('scenesContainer');
            container.innerHTML = '';
            
            storyData.scenes.forEach((scene, index) => {
                const card = document.createElement('div');
                card.className = 'scene-card' + (scene.hasPlaceholders ? ' has-placeholders' : '');
                
                const otherCount = scene.otherText.length;
                const moreText = otherCount > 0 ? '<div class="more-indicator">+' + otherCount + ' more</div>' : '';
                
                card.innerHTML = 
                    '<div class="scene-header">' +
                        '<div class="scene-number">' + scene.number + '</div>' +
                        '<div class="scene-title">' + scene.title + '</div>' +
                        '<span class="placeholder-badge">Has placeholders</span>' +
                    '</div>' +
                    '<div class="hero-text">' + scene.heroText + '</div>' +
                    moreText +
                    '<div class="more-text">' +
                        scene.otherText.map(text => 
                            '<div class="other-text-item">' + text + '</div>'
                        ).join('') +
                    '</div>';
                
                card.onclick = () => {
                    card.classList.toggle('expanded');
                };
                
                container.appendChild(card);
            });
        }
        
        function makeAction(action) {
            fetch('/action', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    story: storyData,
                    action: action
                })
            })
            .then(response => response.json())
            .then(data => {
                storyData = data;
                displayStory();
                showMessage('Text updated successfully!');
            })
            .catch(error => {
                console.error('Error:', error);
                showMessage('Error updating text.', 'error');
            });
        }
        
        function exportText() {
            let text = storyData.projectName + '\\n\\n';
            
            storyData.scenes.forEach(scene => {
                text += 'Scene ' + scene.number + ': ' + scene.title + '\\n';
                text += scene.heroText + '\\n';
                scene.otherText.forEach(t => {
                    text += t + '\\n';
                });
                text += '\\n';
            });
            
            navigator.clipboard.writeText(text).then(() => {
                showMessage('Text copied to clipboard!');
            });
        }
        
        function showMessage(message, type = 'success') {
            const msg = document.getElementById('successMessage');
            msg.textContent = message;
            msg.style.display = 'block';
            msg.style.background = type === 'error' ? '#f44336' : '#4CAF50';
            
            setTimeout(() => {
                msg.style.display = 'none';
            }, 3000);
        }
        
        function resetUpload() {
            document.getElementById('uploadArea').style.display = 'block';
            document.getElementById('loading').style.display = 'none';
            document.getElementById('storyView').style.display = 'none';
        }
    </script>
</body>
</html>`

// SimpleServer handles the simplified web interface
type SimpleServer struct {
	Port string
}

// ExtractSimpleStory extracts a simplified story structure
func ExtractSimpleStory(filename string, project *aep.Project) SimpleStory {
	story := SimpleStory{
		ProjectName: strings.TrimSuffix(filepath.Base(filename), ".aep"),
		Scenes:      []Scene{},
		TotalText:   0,
	}
	
	// Group compositions into scenes
	type CompText struct {
		CompName string
		Texts    []string
		Order    int
	}
	
	var compTexts []CompText
	
	// Extract text from all compositions
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			comp := CompText{
				CompName: item.Name,
				Texts:    []string{},
				Order:    extractSceneNumber(item.Name),
			}
			
			// Get all text layers
			for _, layer := range item.CompositionLayers {
				if layer.Text != nil {
					text := layer.Name
					if text != "" && !strings.HasPrefix(text, "Layer") {
						comp.Texts = append(comp.Texts, text)
						story.TotalText++
					}
				}
			}
			
			if len(comp.Texts) > 0 {
				compTexts = append(compTexts, comp)
			}
		}
	}
	
	// Sort by scene order
	sort.Slice(compTexts, func(i, j int) bool {
		return compTexts[i].Order < compTexts[j].Order
	})
	
	// Group into max 10 scenes
	maxScenes := 10
	if len(compTexts) <= maxScenes {
		// One comp per scene
		for i, comp := range compTexts {
			scene := Scene{
				Number:    i + 1,
				Title:     cleanSceneName(comp.CompName),
				HeroText:  comp.Texts[0],
				OtherText: []string{},
			}
			
			if len(comp.Texts) > 1 {
				scene.OtherText = comp.Texts[1:]
			}
			
			// Check for placeholders
			for _, text := range comp.Texts {
				if strings.Contains(text, "[") || strings.Contains(strings.ToLower(text), "placeholder") {
					scene.HasPlaceholders = true
					break
				}
			}
			
			story.Scenes = append(story.Scenes, scene)
		}
	} else {
		// Group multiple comps per scene
		compsPerScene := (len(compTexts) + maxScenes - 1) / maxScenes
		sceneNum := 1
		
		for i := 0; i < len(compTexts); i += compsPerScene {
			end := i + compsPerScene
			if end > len(compTexts) {
				end = len(compTexts)
			}
			
			// Combine texts from multiple comps
			var allTexts []string
			title := "Part " + fmt.Sprint(sceneNum)
			
			for j := i; j < end; j++ {
				allTexts = append(allTexts, compTexts[j].Texts...)
				if j == i {
					title = cleanSceneName(compTexts[j].CompName)
				}
			}
			
			scene := Scene{
				Number:    sceneNum,
				Title:     title,
				HeroText:  allTexts[0],
				OtherText: []string{},
			}
			
			if len(allTexts) > 1 {
				scene.OtherText = allTexts[1:]
			}
			
			// Check for placeholders
			for _, text := range allTexts {
				if strings.Contains(text, "[") || strings.Contains(strings.ToLower(text), "placeholder") {
					scene.HasPlaceholders = true
					break
				}
			}
			
			story.Scenes = append(story.Scenes, scene)
			sceneNum++
		}
	}
	
	return story
}

// cleanSceneName removes technical prefixes from scene names
func cleanSceneName(name string) string {
	// Remove common prefixes
	name = strings.TrimPrefix(name, "Scene ")
	name = strings.TrimPrefix(name, "S")
	name = strings.TrimPrefix(name, "Comp ")
	
	// Remove numbers at start
	parts := strings.Fields(name)
	if len(parts) > 1 {
		// Check if first part is just a number
		var num int
		if n, _ := fmt.Sscanf(parts[0], "%d", &num); n == 1 {
			name = strings.Join(parts[1:], " ")
		}
	}
	
	// Default names
	if name == "" || strings.HasPrefix(strings.ToLower(name), "comp") {
		name = "Scene"
	}
	
	return name
}

// extractSceneNumber tries to extract a scene number from comp name
func extractSceneNumber(name string) int {
	parts := strings.Fields(name)
	for i, part := range parts {
		if strings.ToLower(part) == "scene" && i+1 < len(parts) {
			var num int
			if _, err := fmt.Sscanf(parts[i+1], "%d", &num); err == nil {
				return num
			}
		}
		if strings.HasPrefix(strings.ToUpper(part), "S") {
			var num int
			if _, err := fmt.Sscanf(part[1:], "%d", &num); err == nil {
				return num
			}
		}
	}
	return 999
}

// HandleUpload processes uploaded AEP files
func (s *SimpleServer) HandleUpload(w http.ResponseWriter, r *http.Request) {
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
	
	// Copy file content
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
	
	// Extract simplified story
	story := ExtractSimpleStory(header.Filename, project)
	
	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(story)
}

// HandleAction processes quick actions
func (s *SimpleServer) HandleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Story  SimpleStory `json:"story"`
		Action string      `json:"action"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	story := req.Story
	
	// Apply transformations based on action
	switch req.Action {
	case "professional":
		for i := range story.Scenes {
			story.Scenes[i].HeroText = makeProfessional(story.Scenes[i].HeroText)
			for j := range story.Scenes[i].OtherText {
				story.Scenes[i].OtherText[j] = makeProfessional(story.Scenes[i].OtherText[j])
			}
		}
		
	case "casual":
		for i := range story.Scenes {
			story.Scenes[i].HeroText = makeCasual(story.Scenes[i].HeroText)
			for j := range story.Scenes[i].OtherText {
				story.Scenes[i].OtherText[j] = makeCasual(story.Scenes[i].OtherText[j])
			}
		}
		
	case "shorter":
		for i := range story.Scenes {
			story.Scenes[i].HeroText = makeShort(story.Scenes[i].HeroText)
			for j := range story.Scenes[i].OtherText {
				story.Scenes[i].OtherText[j] = makeShort(story.Scenes[i].OtherText[j])
			}
		}
		
	case "placeholders":
		fillCount := 1
		for i := range story.Scenes {
			if strings.Contains(story.Scenes[i].HeroText, "[") {
				story.Scenes[i].HeroText = fmt.Sprintf("Amazing Content %d", fillCount)
				fillCount++
			}
			for j := range story.Scenes[i].OtherText {
				if strings.Contains(story.Scenes[i].OtherText[j], "[") {
					story.Scenes[i].OtherText[j] = fmt.Sprintf("Brilliant Text %d", fillCount)
					fillCount++
				}
			}
			story.Scenes[i].HasPlaceholders = false
		}
	}
	
	// Return modified story
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(story)
}

// Text transformation helpers
func makeProfessional(text string) string {
	if strings.Contains(text, "[") {
		return text // Don't change placeholders
	}
	
	replacements := map[string]string{
		"awesome": "excellent",
		"cool": "professional", 
		"great": "exceptional",
		"amazing": "innovative",
		"fun": "engaging",
	}
	
	result := text
	for old, new := range replacements {
		result = strings.ReplaceAll(strings.ToLower(result), old, new)
	}
	
	// Capitalize first letter
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	
	return result
}

func makeCasual(text string) string {
	if strings.Contains(text, "[") {
		return text
	}
	
	replacements := map[string]string{
		"utilize": "use",
		"implement": "do",
		"exceptional": "great",
		"innovative": "cool",
		"professional": "nice",
	}
	
	result := text
	for old, new := range replacements {
		result = strings.ReplaceAll(strings.ToLower(result), old, new)
	}
	
	return result
}

func makeShort(text string) string {
	if strings.Contains(text, "[") {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) > 5 {
		return strings.Join(words[:5], " ") + "..."
	}
	
	return text
}

// HandleHome serves the main HTML interface
func (s *SimpleServer) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(simpleHTML))
}

// Start starts the web server
func (s *SimpleServer) Start() error {
	http.HandleFunc("/", s.HandleHome)
	http.HandleFunc("/upload", s.HandleUpload)
	http.HandleFunc("/action", s.HandleAction)
	
	fmt.Printf("🚀 Simple Story Viewer started on http://localhost:%s\n", s.Port)
	fmt.Println("📖 Open your browser to view stories!")
	
	return http.ListenAndServe(":"+s.Port, nil)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &SimpleServer{
		Port: port,
	}
	
	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}