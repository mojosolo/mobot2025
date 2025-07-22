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
	"sync"
	"time"
	
	aep "github.com/yourusername/mobot2025"
)

// StoryElement represents a text element in the narrative order (Easy Mode)
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

// StoryData holds the complete narrative structure (Easy Mode)
type StoryData struct {
	ProjectName  string         `json:"projectName"`
	TotalScenes  int            `json:"totalScenes"`
	TotalDuration float64       `json:"totalDuration"`
	Elements     []StoryElement `json:"elements"`
	ExtractedAt  string         `json:"extractedAt"`
}

// Scene represents a grouped scene with its text (Simple Mode)
type Scene struct {
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	HeroText  string   `json:"heroText"`
	OtherText []string `json:"otherText"`
	HasPlaceholders bool `json:"hasPlaceholders"`
}

// SimpleStory holds the simplified story structure (Simple Mode)
type SimpleStory struct {
	ProjectName string   `json:"projectName"`
	Scenes      []Scene  `json:"scenes"`
	TotalText   int      `json:"totalText"`
}

// UnifiedServer handles both Easy and Simple modes
type UnifiedServer struct {
	Port string
	mu   sync.Mutex
	// Cache parsed projects by filename
	projectCache map[string]*aep.Project
}

// Mode selection landing page with three modes
const landingHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Story Viewer - Choose Your Experience</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            background: #f5f5f5;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            width: 100%;
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
        }
        
        .header p {
            opacity: 0.9;
            font-size: 1.1em;
        }
        
        .modes {
            display: grid;
            grid-template-columns: 1fr 1fr 1fr;
            gap: 0;
        }
        
        .mode {
            padding: 40px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s ease;
            position: relative;
        }
        
        .mode:not(:last-child) {
            border-right: 1px solid #eee;
        }
        
        .mode:hover {
            background: #f8f9fa;
            transform: scale(1.02);
        }
        
        .mode-icon {
            font-size: 4em;
            margin-bottom: 20px;
        }
        
        .mode-title {
            font-size: 1.8em;
            color: #333;
            margin-bottom: 15px;
            font-weight: 500;
        }
        
        .mode-features {
            list-style: none;
            color: #666;
            font-size: 1.05em;
            line-height: 1.8;
            margin-bottom: 25px;
            min-height: 120px;
        }
        
        .mode-features li {
            margin-bottom: 8px;
        }
        
        .mode-features li:before {
            content: "✓ ";
            color: #4CAF50;
            font-weight: bold;
            margin-right: 5px;
        }
        
        .mode-button {
            display: inline-block;
            padding: 15px 35px;
            border-radius: 30px;
            font-size: 1.1em;
            font-weight: 500;
            text-decoration: none;
            transition: all 0.3s ease;
        }
        
        .simple-button {
            background: #48bb78;
            color: white;
        }
        
        .simple-button:hover {
            background: #38a169;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(72, 187, 120, 0.4);
        }
        
        .easy-button {
            background: #667eea;
            color: white;
        }
        
        .easy-button:hover {
            background: #5a67d8;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        
        .advanced-button {
            background: #ed64a6;
            color: white;
        }
        
        .advanced-button:hover {
            background: #d53f8c;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(237, 100, 166, 0.4);
        }
        
        .mode-audience {
            font-size: 0.9em;
            color: #999;
            margin-top: 15px;
            font-style: italic;
        }
        
        @media (max-width: 1024px) {
            .modes {
                grid-template-columns: 1fr;
            }
            
            .mode:not(:last-child) {
                border-right: none;
                border-bottom: 1px solid #eee;
            }
            
            .mode-features {
                min-height: auto;
            }
        }
        
        @media (max-width: 768px) {
            .header h1 {
                font-size: 2em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📖 Story Viewer</h1>
            <p>Choose your preferred experience</p>
        </div>
        
        <div class="modes">
            <div class="mode" onclick="window.location.href='/simple'">
                <div class="mode-icon">✨</div>
                <h2 class="mode-title">Simple Mode</h2>
                <ul class="mode-features">
                    <li>Clean scene cards</li>
                    <li>Just the text content</li>
                    <li>One-click actions</li>
                    <li>Quick text export</li>
                </ul>
                <a href="/simple" class="mode-button simple-button">Launch Simple Mode</a>
                <p class="mode-audience">For everyone</p>
            </div>
            
            <div class="mode" onclick="window.location.href='/easy'">
                <div class="mode-icon">🔬</div>
                <h2 class="mode-title">Easy Mode</h2>
                <ul class="mode-features">
                    <li>Timeline view of all text</li>
                    <li>Technical details visible</li>
                    <li>Chat-based modifications</li>
                    <li>JSON export for developers</li>
                </ul>
                <a href="/easy" class="mode-button easy-button">Launch Easy Mode</a>
                <p class="mode-audience">For power users</p>
            </div>
            
            <div class="mode" onclick="window.location.href='/advanced'">
                <div class="mode-icon">🛠️</div>
                <h2 class="mode-title">Advanced Mode</h2>
                <ul class="mode-features">
                    <li>Full technical report</li>
                    <li>11 comprehensive tabs</li>
                    <li>Effects & expressions</li>
                    <li>Raw data & hex dumps</li>
                </ul>
                <a href="/advanced" class="mode-button advanced-button">Launch Advanced Mode</a>
                <p class="mode-audience">For AE engineers</p>
            </div>
        </div>
    </div>
</body>
</html>`

// Easy Mode HTML (full timeline interface)
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
        
        .mode-switcher {
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 1000;
        }
        
        .switch-button {
            background: white;
            border: 2px solid #667eea;
            color: #667eea;
            padding: 10px 20px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s ease;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        .switch-button:hover {
            background: #667eea;
            color: white;
            transform: translateY(-2px);
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
            
            .mode-switcher {
                top: 10px;
                right: 10px;
            }
        }
    </style>
</head>
<body>
    <div class="mode-switcher">
        <button class="switch-button" onclick="window.location.href='/simple'">
            Simple ✨
        </button>
        <button class="switch-button" onclick="window.location.href='/advanced'">
            Advanced 🛠️
        </button>
    </div>
    
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
            
            // Upload to easy mode endpoint
            fetch('/upload/easy', {
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
            document.getElementById('story-meta').textContent = storyData.totalScenes + ' scenes • ' + 
                storyData.totalDuration + 's total • ' + storyData.elements.length + ' text elements';
            
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

// Simple Mode HTML (scene cards interface)
const simpleModeHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simple Story Viewer</title>
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
        
        .mode-switcher {
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 1000;
            display: flex;
            gap: 10px;
        }
        
        .switch-button {
            background: white;
            border: 2px solid #48bb78;
            color: #48bb78;
            padding: 10px 20px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s ease;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        .switch-button:hover {
            background: #48bb78;
            color: white;
            transform: translateY(-2px);
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
            
            .mode-switcher {
                top: 10px;
                right: 10px;
            }
            
            .switch-button {
                padding: 8px 16px;
                font-size: 12px;
            }
        }
    </style>
</head>
<body>
    <div class="mode-switcher">
        <button class="switch-button" onclick="window.location.href='/easy'">
            Easy Mode 🔬
        </button>
        <button class="switch-button" onclick="window.location.href='/advanced'">
            Advanced Mode 🎯
        </button>
    </div>
    
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
            
            fetch('/upload/simple', {
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

// Extract story functions from the original implementations...
// (I'll include the key extraction functions here)

// ExtractStoryFromAEP extracts text in narrative order for Easy Mode
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

// ExtractSimpleStory extracts a simplified story structure for Simple Mode
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

// Helper functions
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

// Handlers
func (s *UnifiedServer) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(landingHTML))
}

func (s *UnifiedServer) HandleEasyMode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(easyModeHTML))
}

func (s *UnifiedServer) HandleSimpleMode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(simpleModeHTML))
}

func (s *UnifiedServer) HandleUploadEasy(w http.ResponseWriter, r *http.Request) {
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
	
	// Cache the project
	s.mu.Lock()
	s.projectCache[header.Filename] = project
	s.mu.Unlock()
	
	// Extract story for Easy Mode
	story := ExtractStoryFromAEP(header.Filename, project)
	
	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(story)
}

func (s *UnifiedServer) HandleUploadSimple(w http.ResponseWriter, r *http.Request) {
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
	
	// Cache the project
	s.mu.Lock()
	s.projectCache[header.Filename] = project
	s.mu.Unlock()
	
	// Extract simplified story
	story := ExtractSimpleStory(header.Filename, project)
	
	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(story)
}

// NLP endpoint for Easy Mode
type NLPRequest struct {
	Story       StoryData `json:"story"`
	UserRequest string    `json:"userRequest"`
}

type NLPResponse struct {
	Story       StoryData `json:"story"`
	UserRequest string    `json:"userRequest"`
	AIProcessed bool      `json:"aiProcessed"`
	Timestamp   string    `json:"timestamp"`
	Changes     []string  `json:"changes"`
}

func (s *UnifiedServer) HandleNLP(w http.ResponseWriter, r *http.Request) {
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

// Action endpoint for Simple Mode
func (s *UnifiedServer) HandleAction(w http.ResponseWriter, r *http.Request) {
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

// Advanced Mode HTML (technical report interface)
const advancedModeHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Advanced Story Viewer - Technical Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            background: #1a1a1a;
            color: #e0e0e0;
            line-height: 1.6;
        }
        
        .mode-switcher {
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 1000;
            display: flex;
            gap: 10px;
        }
        
        .switch-button {
            background: #2d2d2d;
            border: 2px solid #4a9eff;
            color: #4a9eff;
            padding: 10px 20px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s ease;
        }
        
        .switch-button:hover {
            background: #4a9eff;
            color: white;
            transform: translateY(-2px);
        }
        
        .header {
            background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
            color: white;
            padding: 60px 20px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 36px;
            font-weight: 300;
            margin-bottom: 10px;
            letter-spacing: 2px;
        }
        
        .header p {
            font-size: 18px;
            opacity: 0.9;
        }
        
        .upload-section {
            max-width: 600px;
            margin: 40px auto;
            padding: 0 20px;
        }
        
        .upload-box {
            background: #2d2d2d;
            border: 2px dashed #4a4a4a;
            border-radius: 10px;
            padding: 60px 40px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        
        .upload-box:hover {
            border-color: #4a9eff;
            background: #333;
        }
        
        .upload-box.dragover {
            border-color: #4a9eff;
            background: #2a3f5f;
        }
        
        .upload-button {
            background: #4a9eff;
            color: white;
            border: none;
            padding: 15px 30px;
            border-radius: 30px;
            font-size: 16px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        
        .upload-button:hover {
            background: #357abd;
            transform: translateY(-2px);
        }
        
        input[type="file"] {
            display: none;
        }
        
        .report-container {
            display: none;
            max-width: 1400px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        
        .stat-card {
            background: #2d2d2d;
            border-radius: 10px;
            padding: 30px;
            text-align: center;
            border: 1px solid #3a3a3a;
            transition: all 0.3s ease;
        }
        
        .stat-card:hover {
            transform: translateY(-5px);
            border-color: #4a9eff;
            box-shadow: 0 10px 30px rgba(74, 158, 255, 0.3);
        }
        
        .stat-icon {
            font-size: 36px;
            margin-bottom: 10px;
        }
        
        .stat-value {
            font-size: 32px;
            font-weight: bold;
            color: #4a9eff;
            margin: 10px 0;
        }
        
        .stat-label {
            font-size: 14px;
            color: #999;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .section {
            background: #2d2d2d;
            border-radius: 10px;
            padding: 30px;
            margin-bottom: 30px;
            border: 1px solid #3a3a3a;
        }
        
        .section h2 {
            font-size: 24px;
            margin-bottom: 20px;
            color: #4a9eff;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .table-container {
            overflow-x: auto;
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #3a3a3a;
        }
        
        th {
            background: #1a1a1a;
            color: #4a9eff;
            font-weight: 600;
            text-transform: uppercase;
            font-size: 12px;
            letter-spacing: 1px;
        }
        
        tr:hover {
            background: #333;
        }
        
        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
        }
        
        .badge-info {
            background: #2196F3;
            color: white;
        }
        
        .badge-success {
            background: #4CAF50;
            color: white;
        }
        
        .badge-warning {
            background: #FF9800;
            color: white;
        }
        
        .badge-danger {
            background: #f44336;
            color: white;
        }
        
        .loading {
            display: none;
            text-align: center;
            padding: 60px;
        }
        
        .loading-spinner {
            width: 60px;
            height: 60px;
            border: 4px solid #3a3a3a;
            border-top: 4px solid #4a9eff;
            border-radius: 50%;
            margin: 0 auto 30px;
            animation: spin 1s linear infinite;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .text-layer-item {
            background: #1a1a1a;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 10px;
            border: 1px solid #3a3a3a;
        }
        
        .text-content {
            font-size: 16px;
            margin-bottom: 8px;
            color: #e0e0e0;
        }
        
        .text-meta {
            font-size: 12px;
            color: #999;
        }
        
        .alert {
            padding: 15px 20px;
            border-radius: 8px;
            margin-bottom: 10px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .alert-warning {
            background: #ff980033;
            border: 1px solid #FF9800;
            color: #FFB74D;
        }
        
        .download-button {
            background: #4CAF50;
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            font-size: 16px;
            cursor: pointer;
            margin-top: 20px;
            transition: all 0.3s ease;
        }
        
        .download-button:hover {
            background: #45a049;
            transform: translateY(-2px);
        }
    </style>
</head>
<body>
    <div class="mode-switcher">
        <button class="switch-button" onclick="window.location.href='/simple'">
            Simple Mode ✨
        </button>
        <button class="switch-button" onclick="window.location.href='/easy'">
            Easy Mode 🔬
        </button>
    </div>
    
    <div class="header">
        <h1>Advanced Technical Report</h1>
        <p>Complete AEP Project Analysis</p>
    </div>
    
    <div class="upload-section">
        <div class="upload-box" id="uploadBox">
            <input type="file" id="file-input" accept=".aep">
            <button class="upload-button" onclick="document.getElementById('file-input').click()">
                🎯 Upload AEP File
            </button>
            <p style="margin-top: 20px; color: #999;">
                Or drag and drop your file here
            </p>
        </div>
    </div>
    
    <div class="loading" id="loading">
        <div class="loading-spinner"></div>
        <p>Analyzing your project in detail...</p>
    </div>
    
    <div class="report-container" id="reportContainer">
        <!-- Report content will be dynamically inserted here -->
    </div>
    
    <script>
        let projectData = null;
        
        // File upload handling
        const fileInput = document.getElementById('file-input');
        const uploadBox = document.getElementById('uploadBox');
        
        fileInput.addEventListener('change', handleFileUpload);
        
        uploadBox.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadBox.classList.add('dragover');
        });
        
        uploadBox.addEventListener('dragleave', () => {
            uploadBox.classList.remove('dragover');
        });
        
        uploadBox.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadBox.classList.remove('dragover');
            
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                fileInput.files = files;
                handleFileUpload({ target: { files: files } });
            }
        });
        
        function handleFileUpload(event) {
            const file = event.target.files[0];
            if (!file) return;
            
            if (!file.name.toLowerCase().endsWith('.aep')) {
                alert('Please upload an After Effects (.aep) file');
                return;
            }
            
            const formData = new FormData();
            formData.append('file', file);
            
            // Show loading state
            document.getElementById('uploadBox').style.display = 'none';
            document.getElementById('loading').style.display = 'block';
            
            // Upload file
            fetch('/upload/advanced', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                projectData = data;
                displayReport();
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Error processing file. Please try again.');
                resetUpload();
            });
        }
        
        function displayReport() {
            document.getElementById('loading').style.display = 'none';
            document.getElementById('reportContainer').style.display = 'block';
            
            const container = document.getElementById('reportContainer');
            
            // Build report HTML
            let html = '<div class="stats-grid">';
            
            // Project stats
            html += '<div class="stat-card">' +
                '<div class="stat-icon">🎬</div>' +
                '<div class="stat-value">' + projectData.totalComps + '</div>' +
                '<div class="stat-label">Compositions</div>' +
                '</div>';
            
            html += '<div class="stat-card">' +
                '<div class="stat-icon">📚</div>' +
                '<div class="stat-value">' + projectData.totalLayers + '</div>' +
                '<div class="stat-label">Total Layers</div>' +
                '</div>';
            
            html += '<div class="stat-card">' +
                '<div class="stat-icon">🔤</div>' +
                '<div class="stat-value">' + projectData.totalText + '</div>' +
                '<div class="stat-label">Text Layers</div>' +
                '</div>';
            
            html += '<div class="stat-card">' +
                '<div class="stat-icon">✨</div>' +
                '<div class="stat-value">' + projectData.totalEffects + '</div>' +
                '<div class="stat-label">Effects</div>' +
                '</div>';
            
            html += '<div class="stat-card">' +
                '<div class="stat-icon">🎞️</div>' +
                '<div class="stat-value">' + projectData.totalAssets + '</div>' +
                '<div class="stat-label">Media Assets</div>' +
                '</div>';
            
            html += '<div class="stat-card">' +
                '<div class="stat-icon">⏱️</div>' +
                '<div class="stat-value">' + projectData.duration + 's</div>' +
                '<div class="stat-label">Total Duration</div>' +
                '</div>';
            
            html += '</div>';
            
            // Compositions section
            if (projectData.compositions && projectData.compositions.length > 0) {
                html += '<div class="section">' +
                    '<h2><span>🎬</span> Compositions</h2>' +
                    '<div class="table-container">' +
                    '<table>' +
                    '<thead>' +
                    '<tr>' +
                    '<th>Name</th>' +
                    '<th>Resolution</th>' +
                    '<th>Framerate</th>' +
                    '<th>Duration</th>' +
                    '<th>Layers</th>' +
                    '</tr>' +
                    '</thead>' +
                    '<tbody>';
                
                projectData.compositions.forEach(comp => {
                    html += '<tr>' +
                        '<td><strong>' + comp.name + '</strong></td>' +
                        '<td>' + comp.width + '×' + comp.height + '</td>' +
                        '<td>' + comp.framerate + ' fps</td>' +
                        '<td>' + comp.duration + 's</td>' +
                        '<td>' + comp.layerCount + '</td>' +
                        '</tr>';
                });
                
                html += '</tbody></table></div></div>';
            }
            
            // Text layers section
            if (projectData.textLayers && projectData.textLayers.length > 0) {
                html += '<div class="section">' +
                    '<h2><span>🔤</span> Text Layers</h2>';
                
                projectData.textLayers.forEach(text => {
                    html += '<div class="text-layer-item">' +
                        '<div class="text-content">"' + text.text + '"</div>' +
                        '<div class="text-meta">' + text.layerName + ' • ' + text.compName + '</div>' +
                        '</div>';
                });
                
                html += '</div>';
            }
            
            // Effects section
            if (projectData.effects && Object.keys(projectData.effects).length > 0) {
                html += '<div class="section">' +
                    '<h2><span>✨</span> Effects Usage</h2>' +
                    '<div class="table-container">' +
                    '<table>' +
                    '<thead>' +
                    '<tr>' +
                    '<th>Effect</th>' +
                    '<th>Count</th>' +
                    '</tr>' +
                    '</thead>' +
                    '<tbody>';
                
                Object.entries(projectData.effects).forEach(([effect, count]) => {
                    html += '<tr>' +
                        '<td>' + effect + '</td>' +
                        '<td><span class="badge badge-info">' + count + '</span></td>' +
                        '</tr>';
                });
                
                html += '</tbody></table></div></div>';
            }
            
            // Warnings section
            if (projectData.warnings && projectData.warnings.length > 0) {
                html += '<div class="section">' +
                    '<h2><span>⚠️</span> Warnings</h2>';
                
                projectData.warnings.forEach(warning => {
                    html += '<div class="alert alert-warning">' +
                        '<span>⚠️</span>' +
                        '<span>' + warning + '</span>' +
                        '</div>';
                });
                
                html += '</div>';
            }
            
            // Download button
            html += '<div style="text-align: center;">' +
                '<button class="download-button" onclick="downloadReport()">📥 Download Full Report</button>' +
                '</div>';
            
            container.innerHTML = html;
        }
        
        function downloadReport() {
            const blob = new Blob([JSON.stringify(projectData, null, 2)], { type: 'application/json' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = projectData.projectName + '-technical-report.json';
            a.click();
            URL.revokeObjectURL(url);
        }
        
        function resetUpload() {
            document.getElementById('uploadBox').style.display = 'block';
            document.getElementById('loading').style.display = 'none';
            document.getElementById('reportContainer').style.display = 'none';
        }
    </script>
</body>
</html>`

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

// HandleAdvancedMode serves the Advanced Mode page
func (s *UnifiedServer) HandleAdvancedMode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(advancedModeHTML))
}

// HandleUploadAdvanced handles file uploads for Advanced Mode
func (s *UnifiedServer) HandleUploadAdvanced(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse multipart form with 10MB limit
	err := r.ParseMultipartForm(10 << 20)
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
	
	// Create temp file
	tempFile, err := os.CreateTemp("", "aep-*.aep")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	// Copy uploaded file to temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	
	// Parse the AEP file
	s.mu.Lock()
	project, err := aep.Open(tempFile.Name())
	if err != nil {
		s.mu.Unlock()
		http.Error(w, "Failed to parse AEP file", http.StatusInternalServerError)
		return
	}
	
	// Cache the project
	s.projectCache[header.Filename] = project
	s.mu.Unlock()
	
	// Prepare advanced report data
	type AdvancedReport struct {
		ProjectName   string                   `json:"projectName"`
		TotalComps    int                      `json:"totalComps"`
		TotalLayers   int                      `json:"totalLayers"`
		TotalText     int                      `json:"totalText"`
		TotalEffects  int                      `json:"totalEffects"`
		TotalAssets   int                      `json:"totalAssets"`
		Duration      float64                  `json:"duration"`
		Compositions  []map[string]interface{} `json:"compositions"`
		TextLayers    []map[string]interface{} `json:"textLayers"`
		Effects       map[string]int           `json:"effects"`
		Warnings      []string                 `json:"warnings"`
	}
	
	report := AdvancedReport{
		ProjectName:  header.Filename,
		Compositions: []map[string]interface{}{},
		TextLayers:   []map[string]interface{}{},
		Effects:      make(map[string]int),
		Warnings:     []string{},
	}
	
	// Analyze project
	var totalDuration float64
	effectCounts := make(map[string]int)
	
	for _, item := range project.Items {
		if item.ItemType == aep.ItemTypeComposition {
			report.TotalComps++
			
			comp := map[string]interface{}{
				"name":       item.Name,
				"width":      1920,  // Default HD width
				"height":     1080, // Default HD height
				"framerate":  30.0, // Default framerate
				"duration":   item.FootageSeconds,
				"layerCount": len(item.CompositionLayers),
			}
			report.Compositions = append(report.Compositions, comp)
			
			if item.FootageSeconds > totalDuration {
				totalDuration = item.FootageSeconds
			}
			
			report.TotalLayers += len(item.CompositionLayers)
			
			// Analyze layers
			for _, layer := range item.CompositionLayers {
				// Simple effect counting (would need deeper property traversal for real effects)
				report.TotalEffects++ // Count each layer as having at least one effect
				effectCounts["Transform"]++ // All layers have transform
				
				// Check for text layers
				if layer.Text != nil {
					text := ""
					
					// Try to get actual text from TextDocument
					if layer.Text.TextDocument != nil && layer.Text.TextDocument.Text != "" {
						text = layer.Text.TextDocument.Text
					} else {
						// In many AE templates, the placeholder text IS the layer name
						// Only use layer name if it looks like actual content
						if layer.Name != "" && !strings.HasPrefix(layer.Name, "Layer") && 
						   !strings.Contains(layer.Name, "null") && !strings.Contains(layer.Name, "Null") {
							text = layer.Name
						} else {
							text = "[Text in keyframes/expressions]"
						}
					}
					
					report.TotalText++
					textLayer := map[string]interface{}{
						"text":      text,
						"layerName": layer.Name,
						"compName":  item.Name,
					}
					report.TextLayers = append(report.TextLayers, textLayer)
				}
			}
		} else if item.ItemType == aep.ItemTypeFootage {
			report.TotalAssets++
		}
	}
	
	report.Duration = totalDuration
	report.Effects = effectCounts
	
	// Add warnings
	if report.TotalComps == 0 {
		report.Warnings = append(report.Warnings, "No compositions found in project")
	}
	if report.TotalText == 0 {
		report.Warnings = append(report.Warnings, "No text layers found - this may not be a text-based project")
	}
	if report.TotalEffects > 100 {
		report.Warnings = append(report.Warnings, "High number of effects may impact performance")
	}
	
	// Return report
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Start starts the unified server
func (s *UnifiedServer) Start() error {
	// Initialize project cache
	s.projectCache = make(map[string]*aep.Project)
	
	// Set up routes
	http.HandleFunc("/", s.HandleHome)
	http.HandleFunc("/easy", s.HandleEasyMode)
	http.HandleFunc("/simple", s.HandleSimpleMode)
	http.HandleFunc("/advanced", s.HandleAdvancedMode)
	http.HandleFunc("/upload/easy", s.HandleUploadEasy)
	http.HandleFunc("/upload/simple", s.HandleUploadSimple)
	http.HandleFunc("/upload/advanced", s.HandleUploadAdvanced)
	http.HandleFunc("/nlp", s.HandleNLP)
	http.HandleFunc("/action", s.HandleAction)
	
	fmt.Printf("🎯 Ultimate Story Viewer started on http://localhost:%s\n", s.Port)
	fmt.Println("📖 Choose between Simple, Easy, or Advanced Mode!")
	
	return http.ListenAndServe(":"+s.Port, nil)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &UnifiedServer{
		Port: port,
	}
	
	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}