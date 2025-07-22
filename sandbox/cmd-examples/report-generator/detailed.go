package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	
	aep "github.com/mojosolo/mobot2025"
)

// DetailedReportData holds comprehensive data for the HTML report
type DetailedReportData struct {
	Filename         string
	ParsedAt         string
	BitDepth         string
	ExpressionEngine string
	TotalItems       int
	
	// Overview tab
	Overview struct {
		FolderCount      int
		CompositionCount int
		FootageCount     int
		TotalLayers      int
		MainComps        []CompositionDetail
		Statistics       ProjectStats
	}
	
	// Compositions tab
	Compositions []CompositionDetail
	
	// Layers tab
	AllLayers []LayerDetail
	
	// Media tab
	MediaAssets []MediaAsset
	
	// Text tab
	TextLayers []TextLayer
	
	// Attributes tab
	LayerAttributes []LayerAttribute
	
	// Hierarchy tab
	FolderTree FolderNode
}

// CompositionDetail holds detailed composition information
type CompositionDetail struct {
	ID         uint32
	Name       string
	Width      uint16
	Height     uint16
	Resolution string
	Framerate  float64
	Duration   float64
	Frames     int
	BGColor    string
	Layers     []LayerDetail
	LayerCount int
}

// LayerDetail holds detailed layer information
type LayerDetail struct {
	Index          uint32
	Name           string
	CompName       string
	SourceID       uint32
	SourceName     string
	Quality        string
	SamplingMode   string
	FrameBlending  string
	Is3D           bool
	IsSolo         bool
	IsGuide        bool
	IsAdjustment   bool
	IsShyred       bool
	IsLocked       bool
	HasMotionBlur  bool
	HasEffects     bool
	IsCollapsed    bool
}

// MediaAsset holds media/footage information
type MediaAsset struct {
	ID        uint32
	Name      string
	Type      string
	Width     uint16
	Height    uint16
	Framerate float64
	Duration  float64
	Frames    int
	UsedIn    []string // Compositions using this asset
	UsageCount int
}

// TextLayer holds text-specific layer information
type TextLayer struct {
	LayerName string
	CompName  string
	Text      string // Would need expression parsing to get actual text
	Font      string
	Size      float64
}

// LayerAttribute holds layer attributes and effects
type LayerAttribute struct {
	LayerName  string
	CompName   string
	Attributes map[string]interface{}
}

// FolderNode represents folder hierarchy
type FolderNode struct {
	Name     string
	Type     string
	Children []FolderNode
	ItemCount int
}

// ProjectStats holds project statistics
type ProjectStats struct {
	NullCount       int
	AdjustmentCount int
	SolidCount      int
	VideoCount      int
	ImageCount      int
	OtherCount      int
	TotalSize       string
	ColorDepth      string
}

const detailedHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AEP Detailed Analysis - {{.Filename}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .subtitle {
            opacity: 0.9;
            font-size: 1.1em;
        }
        
        .tabs {
            background: white;
            padding: 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
            position: sticky;
            top: 0;
            z-index: 100;
        }
        
        .tab-list {
            display: flex;
            list-style: none;
            margin: 0;
            padding: 0;
            border-bottom: 2px solid #eee;
        }
        
        .tab-button {
            padding: 15px 30px;
            background: none;
            border: none;
            font-size: 1em;
            cursor: pointer;
            color: #666;
            transition: all 0.3s ease;
            border-bottom: 3px solid transparent;
            font-weight: 500;
        }
        
        .tab-button:hover {
            color: #667eea;
            background: #f8f9fa;
        }
        
        .tab-button.active {
            color: #667eea;
            border-bottom-color: #667eea;
            background: #f8f9fa;
        }
        
        .tab-content {
            display: none;
            padding: 30px;
            animation: fadeIn 0.3s ease;
        }
        
        .tab-content.active {
            display: block;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .stat-card {
            background: white;
            padding: 25px;
            border-radius: 10px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
            transition: transform 0.2s ease;
        }
        
        .stat-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        }
        
        .stat-number {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
            display: block;
        }
        
        .stat-label {
            color: #666;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        table {
            width: 100%;
            background: white;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
            margin-bottom: 20px;
        }
        
        th, td {
            padding: 12px 15px;
            text-align: left;
        }
        
        th {
            background: #f8f9fa;
            font-weight: 600;
            color: #667eea;
            position: sticky;
            top: 60px;
            z-index: 10;
        }
        
        tr {
            border-bottom: 1px solid #eee;
        }
        
        tr:hover {
            background: #f8f9fa;
        }
        
        .badge {
            display: inline-block;
            padding: 4px 10px;
            border-radius: 20px;
            font-size: 0.8em;
            font-weight: 600;
            margin: 0 2px;
        }
        
        .badge-4k { background: #10b981; color: white; }
        .badge-hd { background: #3b82f6; color: white; }
        .badge-3d { background: #f59e0b; color: white; }
        .badge-solo { background: #ef4444; color: white; }
        .badge-adjustment { background: #8b5cf6; color: white; }
        .badge-effects { background: #ec4899; color: white; }
        .badge-video { background: #06b6d4; color: white; }
        .badge-image { background: #84cc16; color: white; }
        .badge-solid { background: #6b7280; color: white; }
        
        .search-box {
            margin-bottom: 20px;
        }
        
        .search-box input {
            width: 100%;
            padding: 12px 20px;
            border: 2px solid #eee;
            border-radius: 8px;
            font-size: 1em;
            transition: border-color 0.3s ease;
        }
        
        .search-box input:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .tree-view {
            background: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        
        .tree-node {
            margin-left: 20px;
            padding: 5px 0;
        }
        
        .tree-node-root {
            margin-left: 0;
            font-weight: bold;
        }
        
        .tree-icon {
            display: inline-block;
            width: 20px;
            margin-right: 5px;
        }
        
        .comp-detail {
            background: white;
            padding: 20px;
            margin-bottom: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        
        .comp-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
            padding-bottom: 15px;
            border-bottom: 2px solid #eee;
        }
        
        .comp-title {
            font-size: 1.4em;
            font-weight: bold;
            color: #333;
        }
        
        .comp-meta {
            display: flex;
            gap: 20px;
            color: #666;
            font-size: 0.9em;
        }
        
        .empty-state {
            text-align: center;
            padding: 60px;
            color: #999;
        }
        
        .empty-state-icon {
            font-size: 4em;
            margin-bottom: 20px;
            opacity: 0.3;
        }
        
        .footer {
            text-align: center;
            padding: 30px;
            color: #666;
            background: white;
            margin-top: 50px;
        }
        
        .tag {
            display: inline-block;
            padding: 2px 8px;
            background: #e5e7eb;
            border-radius: 4px;
            font-size: 0.85em;
            margin: 2px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎬 {{.Filename}}</h1>
        <div class="subtitle">Detailed AEP Analysis Report</div>
        <div style="margin-top: 10px; opacity: 0.8;">Generated on {{.ParsedAt}}</div>
    </div>
    
    <div class="tabs">
        <ul class="tab-list">
            <li><button class="tab-button active" onclick="showTab('overview')">📊 Overview</button></li>
            <li><button class="tab-button" onclick="showTab('compositions')">🎬 Compositions</button></li>
            <li><button class="tab-button" onclick="showTab('layers')">📑 All Layers</button></li>
            <li><button class="tab-button" onclick="showTab('media')">📹 Media Assets</button></li>
            <li><button class="tab-button" onclick="showTab('text')">📝 Text Layers</button></li>
            <li><button class="tab-button" onclick="showTab('attributes')">⚙️ Attributes</button></li>
            <li><button class="tab-button" onclick="showTab('hierarchy')">📁 Hierarchy</button></li>
        </ul>
    </div>
    
    <div class="container">
        <!-- Overview Tab -->
        <div id="overview" class="tab-content active">
            <h2>Project Overview</h2>
            
            <div class="stats-grid">
                <div class="stat-card">
                    <span class="stat-number">{{.TotalItems}}</span>
                    <span class="stat-label">Total Items</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.CompositionCount}}</span>
                    <span class="stat-label">Compositions</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.FootageCount}}</span>
                    <span class="stat-label">Media Assets</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.TotalLayers}}</span>
                    <span class="stat-label">Total Layers</span>
                </div>
            </div>
            
            <h3>Project Properties</h3>
            <table>
                <tr>
                    <td><strong>Bit Depth:</strong></td>
                    <td>{{.BitDepth}}</td>
                    <td><strong>Expression Engine:</strong></td>
                    <td>{{.ExpressionEngine}}</td>
                </tr>
                <tr>
                    <td><strong>Folder Count:</strong></td>
                    <td>{{.Overview.FolderCount}}</td>
                    <td><strong>Color Depth:</strong></td>
                    <td>{{.Overview.Statistics.ColorDepth}}</td>
                </tr>
            </table>
            
            {{if .Overview.MainComps}}
            <h3>Main Compositions</h3>
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Resolution</th>
                        <th>Framerate</th>
                        <th>Duration</th>
                        <th>Layers</th>
                        <th>Background</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Overview.MainComps}}
                    <tr>
                        <td>
                            <strong>{{.Name}}</strong>
                            {{if eq .Resolution "4K"}}<span class="badge badge-4k">4K</span>{{end}}
                            {{if eq .Resolution "HD"}}<span class="badge badge-hd">HD</span>{{end}}
                        </td>
                        <td>{{.Width}} × {{.Height}}</td>
                        <td>{{printf "%.2f" .Framerate}} fps</td>
                        <td>{{printf "%.2f" .Duration}}s ({{.Frames}} frames)</td>
                        <td>{{.LayerCount}}</td>
                        <td><span style="display:inline-block;width:20px;height:20px;background:{{.BGColor}};border:1px solid #ddd;vertical-align:middle;"></span> {{.BGColor}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{end}}
            
            <h3>Asset Breakdown</h3>
            <div class="stats-grid">
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.Statistics.VideoCount}}</span>
                    <span class="stat-label">Videos</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.Statistics.ImageCount}}</span>
                    <span class="stat-label">Images</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.Statistics.SolidCount}}</span>
                    <span class="stat-label">Solids</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.Statistics.AdjustmentCount}}</span>
                    <span class="stat-label">Adjustments</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Overview.Statistics.NullCount}}</span>
                    <span class="stat-label">Nulls</span>
                </div>
            </div>
        </div>
        
        <!-- Compositions Tab -->
        <div id="compositions" class="tab-content">
            <h2>All Compositions ({{len .Compositions}})</h2>
            
            <div class="search-box">
                <input type="text" id="comp-search" placeholder="🔍 Search compositions..." onkeyup="filterTable('comp-search', 'comp-table')">
            </div>
            
            {{if .Compositions}}
            {{range .Compositions}}
            <div class="comp-detail">
                <div class="comp-header">
                    <div class="comp-title">
                        {{.Name}}
                        {{if eq .Resolution "4K"}}<span class="badge badge-4k">4K</span>{{end}}
                        {{if eq .Resolution "HD"}}<span class="badge badge-hd">HD</span>{{end}}
                    </div>
                    <div class="comp-meta">
                        <span>📐 {{.Width}}×{{.Height}}</span>
                        <span>🎞️ {{printf "%.2f" .Framerate}} fps</span>
                        <span>⏱️ {{printf "%.2f" .Duration}}s</span>
                        <span>📑 {{.LayerCount}} layers</span>
                    </div>
                </div>
                
                {{if .Layers}}
                <table>
                    <thead>
                        <tr>
                            <th width="50">#</th>
                            <th>Layer Name</th>
                            <th>Source</th>
                            <th>Properties</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Layers}}
                        <tr>
                            <td>{{.Index}}</td>
                            <td>{{.Name}}</td>
                            <td>{{.SourceName}}</td>
                            <td>
                                {{if .Is3D}}<span class="badge badge-3d">3D</span>{{end}}
                                {{if .IsSolo}}<span class="badge badge-solo">Solo</span>{{end}}
                                {{if .IsAdjustment}}<span class="badge badge-adjustment">Adjustment</span>{{end}}
                                {{if .HasEffects}}<span class="badge badge-effects">Effects</span>{{end}}
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
                {{end}}
            </div>
            {{end}}
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">📭</div>
                <p>No compositions found in this project</p>
            </div>
            {{end}}
        </div>
        
        <!-- Layers Tab -->
        <div id="layers" class="tab-content">
            <h2>All Layers ({{len .AllLayers}})</h2>
            
            <div class="search-box">
                <input type="text" id="layer-search" placeholder="🔍 Search layers..." onkeyup="filterTable('layer-search', 'layer-table')">
            </div>
            
            {{if .AllLayers}}
            <table id="layer-table">
                <thead>
                    <tr>
                        <th>Layer Name</th>
                        <th>Composition</th>
                        <th>Index</th>
                        <th>Source</th>
                        <th>Quality</th>
                        <th>Properties</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .AllLayers}}
                    <tr>
                        <td><strong>{{.Name}}</strong></td>
                        <td>{{.CompName}}</td>
                        <td>{{.Index}}</td>
                        <td>{{.SourceName}}</td>
                        <td>{{.Quality}}</td>
                        <td>
                            {{if .Is3D}}<span class="badge badge-3d">3D</span>{{end}}
                            {{if .IsSolo}}<span class="badge badge-solo">Solo</span>{{end}}
                            {{if .IsAdjustment}}<span class="badge badge-adjustment">Adj</span>{{end}}
                            {{if .HasMotionBlur}}<span class="badge">Motion Blur</span>{{end}}
                            {{if .HasEffects}}<span class="badge badge-effects">Effects</span>{{end}}
                            {{if .IsCollapsed}}<span class="badge">Collapsed</span>{{end}}
                            {{if .IsGuide}}<span class="badge">Guide</span>{{end}}
                            {{if .IsLocked}}<span class="badge">Locked</span>{{end}}
                            {{if .IsShyred}}<span class="badge">Shy</span>{{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">📭</div>
                <p>No layers found in this project</p>
            </div>
            {{end}}
        </div>
        
        <!-- Media Tab -->
        <div id="media" class="tab-content">
            <h2>Media Assets ({{len .MediaAssets}})</h2>
            
            <div class="search-box">
                <input type="text" id="media-search" placeholder="🔍 Search media assets..." onkeyup="filterTable('media-search', 'media-table')">
            </div>
            
            {{if .MediaAssets}}
            <table id="media-table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Type</th>
                        <th>Dimensions</th>
                        <th>Framerate</th>
                        <th>Duration</th>
                        <th>Usage</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .MediaAssets}}
                    <tr>
                        <td><strong>{{.Name}}</strong></td>
                        <td>
                            {{if eq .Type "Video"}}<span class="badge badge-video">Video</span>{{end}}
                            {{if eq .Type "Image"}}<span class="badge badge-image">Image</span>{{end}}
                            {{if eq .Type "Solid"}}<span class="badge badge-solid">Solid</span>{{end}}
                            {{if eq .Type "Null Object"}}<span class="badge">Null</span>{{end}}
                            {{if eq .Type "Adjustment"}}<span class="badge badge-adjustment">Adjustment</span>{{end}}
                        </td>
                        <td>{{if gt .Width 0}}{{.Width}}×{{.Height}}{{else}}-{{end}}</td>
                        <td>{{if gt .Framerate 0.0}}{{printf "%.2f" .Framerate}} fps{{else}}-{{end}}</td>
                        <td>{{if gt .Duration 0.0}}{{printf "%.2f" .Duration}}s{{else}}-{{end}}</td>
                        <td>
                            {{.UsageCount}} comp(s)
                            {{if .UsedIn}}
                            <div style="font-size:0.85em; color:#666;">
                                {{range .UsedIn}}<span class="tag">{{.}}</span>{{end}}
                            </div>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">📭</div>
                <p>No media assets found in this project</p>
            </div>
            {{end}}
        </div>
        
        <!-- Text Tab -->
        <div id="text" class="tab-content">
            <h2>Text Layers</h2>
            
            <div class="search-box">
                <input type="text" id="text-search" placeholder="🔍 Search text layers..." onkeyup="filterTable('text-search', 'text-table')">
            </div>
            
            {{if .TextLayers}}
            <table id="text-table">
                <thead>
                    <tr>
                        <th>Layer Name</th>
                        <th>Composition</th>
                        <th>Text Content</th>
                        <th>Font</th>
                        <th>Size</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .TextLayers}}
                    <tr>
                        <td><strong>{{.LayerName}}</strong></td>
                        <td>{{.CompName}}</td>
                        <td>{{.Text}}</td>
                        <td>{{.Font}}</td>
                        <td>{{.Size}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">📝</div>
                <p>No text layers detected in this project</p>
                <p style="margin-top:10px; font-size:0.9em; color:#999;">
                    Note: Text content extraction requires parsing expressions and keyframes,<br>
                    which is not yet implemented in the parser.
                </p>
            </div>
            {{end}}
        </div>
        
        <!-- Attributes Tab -->
        <div id="attributes" class="tab-content">
            <h2>Layer Attributes & Effects</h2>
            
            <div class="search-box">
                <input type="text" id="attr-search" placeholder="🔍 Search attributes..." onkeyup="filterTable('attr-search', 'attr-table')">
            </div>
            
            <p style="margin-bottom:20px; color:#666;">
                Detailed layer attributes including transformations, effects, and expressions.
            </p>
            
            {{if .LayerAttributes}}
            <table id="attr-table">
                <thead>
                    <tr>
                        <th>Layer</th>
                        <th>Composition</th>
                        <th>Attributes</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .LayerAttributes}}
                    <tr>
                        <td><strong>{{.LayerName}}</strong></td>
                        <td>{{.CompName}}</td>
                        <td>
                            {{range $key, $value := .Attributes}}
                            <span class="tag">{{$key}}: {{$value}}</span>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">⚙️</div>
                <p>Layer attribute extraction coming soon</p>
                <p style="margin-top:10px; font-size:0.9em; color:#999;">
                    This will include transform properties, effects, masks, and expressions.
                </p>
            </div>
            {{end}}
        </div>
        
        <!-- Hierarchy Tab -->
        <div id="hierarchy" class="tab-content">
            <h2>Project Hierarchy</h2>
            
            <div class="tree-view">
                {{template "renderTree" .FolderTree}}
            </div>
        </div>
    </div>
    
    <div class="footer">
        <p>Generated by mobot2025 AEP Parser</p>
        <p>github.com/mojosolo/mobot2025</p>
    </div>
    
    <script>
        function showTab(tabName) {
            // Hide all tabs
            const tabs = document.querySelectorAll('.tab-content');
            tabs.forEach(tab => tab.classList.remove('active'));
            
            // Remove active from all buttons
            const buttons = document.querySelectorAll('.tab-button');
            buttons.forEach(btn => btn.classList.remove('active'));
            
            // Show selected tab
            document.getElementById(tabName).classList.add('active');
            
            // Mark button as active
            event.target.classList.add('active');
        }
        
        function filterTable(inputId, tableId) {
            const input = document.getElementById(inputId);
            const filter = input.value.toUpperCase();
            const table = document.getElementById(tableId);
            
            if (!table) return;
            
            const rows = table.getElementsByTagName('tr');
            
            for (let i = 1; i < rows.length; i++) {
                const row = rows[i];
                const text = row.textContent || row.innerText;
                
                if (text.toUpperCase().indexOf(filter) > -1) {
                    row.style.display = '';
                } else {
                    row.style.display = 'none';
                }
            }
        }
    </script>
</body>
</html>

{{define "renderTree"}}
    <div class="tree-node {{if eq .Name "root"}}tree-node-root{{end}}">
        <span class="tree-icon">
            {{if eq .Type "folder"}}📁{{else if eq .Type "comp"}}🎬{{else}}🎞️{{end}}
        </span>
        {{.Name}}
        {{if gt .ItemCount 0}}
            <span style="color:#999; font-size:0.9em;">({{.ItemCount}} items)</span>
        {{end}}
    </div>
    {{range .Children}}
        {{template "renderTree" .}}
    {{end}}
{{end}}`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_detailed_report.go <aep-file>")
		fmt.Println("Example: go run generate_detailed_report.go sample-aep/Ai\\ Text\\ Intro.aep")
		os.Exit(1)
	}
	
	aepFile := os.Args[1]
	fmt.Printf("📄 Parsing: %s\n", aepFile)
	
	// Parse the AEP file
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("❌ Error parsing AEP file: %v\n", err)
		os.Exit(1)
	}
	
	// Prepare detailed report data
	report := prepareDetailedReport(aepFile, project)
	
	// Generate HTML filename
	baseFilename := strings.TrimSuffix(filepath.Base(aepFile), filepath.Ext(aepFile))
	htmlFilename := fmt.Sprintf("%s-detailed-report-%s.html", 
		baseFilename, 
		time.Now().Format("2006-01-02-150405"))
	
	// Create HTML file
	file, err := os.Create(htmlFilename)
	if err != nil {
		fmt.Printf("❌ Error creating HTML file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	
	// Parse and execute template
	tmpl, err := template.New("report").Parse(detailedHTMLTemplate)
	if err != nil {
		fmt.Printf("❌ Error parsing template: %v\n", err)
		os.Exit(1)
	}
	
	err = tmpl.Execute(file, report)
	if err != nil {
		fmt.Printf("❌ Error generating HTML: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ Detailed report generated: %s\n", htmlFilename)
	fmt.Printf("📂 Open in browser: file://%s/%s\n", 
		mustGetWorkingDir(), htmlFilename)
}

func prepareDetailedReport(filename string, project *aep.Project) DetailedReportData {
	report := DetailedReportData{
		Filename:         filepath.Base(filename),
		ParsedAt:         time.Now().Format("January 2, 2006 at 3:04 PM"),
		TotalItems:       len(project.Items),
		ExpressionEngine: project.ExpressionEngine,
	}
	
	// Convert bit depth
	switch project.Depth {
	case aep.BPC8:
		report.BitDepth = "8-bit"
		report.Overview.Statistics.ColorDepth = "16.7 million colors"
	case aep.BPC16:
		report.BitDepth = "16-bit"
		report.Overview.Statistics.ColorDepth = "281 trillion colors"
	case aep.BPC32:
		report.BitDepth = "32-bit (float)"
		report.Overview.Statistics.ColorDepth = "Floating point precision"
	default:
		report.BitDepth = fmt.Sprintf("Unknown (%v)", project.Depth)
	}
	
	// Maps for tracking usage
	mediaUsage := make(map[uint32][]string)
	allComps := make([]CompositionDetail, 0)
	allLayers := make([]LayerDetail, 0)
	
	// Process all items
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeFolder:
			report.Overview.FolderCount++
			
		case aep.ItemTypeComposition:
			report.Overview.CompositionCount++
			report.Overview.TotalLayers += len(item.CompositionLayers)
			
			comp := CompositionDetail{
				ID:         item.ID,
				Name:       item.Name,
				Width:      item.FootageDimensions[0],
				Height:     item.FootageDimensions[1],
				Framerate:  item.FootageFramerate,
				Duration:   item.FootageSeconds,
				Frames:     int(item.FootageSeconds * item.FootageFramerate),
				LayerCount: len(item.CompositionLayers),
				BGColor:    fmt.Sprintf("rgb(%d,%d,%d)", 
					item.BackgroundColor[0],
					item.BackgroundColor[1],
					item.BackgroundColor[2]),
			}
			
			// Determine resolution
			if comp.Width >= 3840 {
				comp.Resolution = "4K"
			} else if comp.Width >= 1920 {
				comp.Resolution = "HD"
			} else {
				comp.Resolution = "SD"
			}
			
			// Process layers in this comp
			for _, layer := range item.CompositionLayers {
				layerDetail := LayerDetail{
					Index:          layer.Index,
					Name:           layer.Name,
					CompName:       item.Name,
					SourceID:       layer.SourceID,
					Is3D:           layer.ThreeDEnabled,
					IsSolo:         layer.SoloEnabled,
					IsGuide:        layer.GuideEnabled,
					IsAdjustment:   layer.AdjustmentLayerEnabled,
					IsShyred:       layer.ShyEnabled,
					IsLocked:       layer.LockEnabled,
					HasMotionBlur:  layer.MotionBlurEnabled,
					HasEffects:     layer.EffectsEnabled,
					IsCollapsed:    layer.CollapseTransformEnabled,
				}
				
				// Get source name
				if sourceItem, exists := project.Items[layer.SourceID]; exists {
					layerDetail.SourceName = sourceItem.Name
					// Track media usage
					mediaUsage[layer.SourceID] = append(mediaUsage[layer.SourceID], item.Name)
				}
				
				// Quality
				switch layer.Quality {
				case aep.LayerQualityBest:
					layerDetail.Quality = "Best"
				case aep.LayerQualityDraft:
					layerDetail.Quality = "Draft"
				case aep.LayerQualityWireframe:
					layerDetail.Quality = "Wireframe"
				}
				
				// Sampling
				switch layer.SamplingMode {
				case aep.LayerSamplingModeBilinear:
					layerDetail.SamplingMode = "Bilinear"
				case aep.LayerSamplingModeBicubic:
					layerDetail.SamplingMode = "Bicubic"
				}
				
				comp.Layers = append(comp.Layers, layerDetail)
				allLayers = append(allLayers, layerDetail)
				
				// Check for text layers
				if layer.Text != nil {
					textLayer := TextLayer{
						LayerName: layer.Name,
						CompName:  item.Name,
						Text:      "(No text content found)",
						Font:      "Unknown",
						Size:      0,
					}
					
					// Extract text content if available
					if layer.Text.TextDocument != nil {
						doc := layer.Text.TextDocument
						if doc.Text != "" {
							textLayer.Text = doc.Text
						}
						if doc.FontName != "" {
							textLayer.Font = doc.FontName
						}
						if doc.FontSize > 0 {
							textLayer.Size = doc.FontSize
						}
					} else if layer.Text.RawData != nil {
						// Try to extract from raw data
						text := strings.TrimSpace(string(layer.Text.RawData))
						if text != "" && !strings.HasPrefix(text, "ADBE") {
							textLayer.Text = text
						}
					} else if layer.Text.Name != "" && layer.Text.Name != layer.Text.MatchName {
						// Use property name as fallback
						textLayer.Text = layer.Text.Name
					}
					
					// Smart fallback - use layer name as text content hint
					if textLayer.Text == "(No text content found)" || 
					   textLayer.Text == "[Text content in keyframes/expressions]" {
						// Many After Effects templates put the actual text in the layer name
						layerNameLower := strings.ToLower(layer.Name)
						
						// Skip generic names
						if !strings.Contains(layerNameLower, "text") &&
						   !strings.Contains(layerNameLower, "title") &&
						   !strings.Contains(layerNameLower, "placeholder") &&
						   !strings.Contains(layerNameLower, "main") &&
						   !strings.Contains(layerNameLower, "colortxt") &&
						   len(layer.Name) > 3 {
							// This is likely the actual text content
							textLayer.Text = layer.Name
						} else {
							textLayer.Text = fmt.Sprintf("[%s]", layer.Name)
						}
					}
					
					report.TextLayers = append(report.TextLayers, textLayer)
				}
			}
			
			allComps = append(allComps, comp)
			
			// Check if main comp
			if strings.Contains(item.Name, "Final Comp") || 
			   strings.Contains(item.Name, "Main Comp") {
				report.Overview.MainComps = append(report.Overview.MainComps, comp)
			}
			
		case aep.ItemTypeFootage:
			report.Overview.FootageCount++
			
			media := MediaAsset{
				ID:        item.ID,
				Name:      item.Name,
				Width:     item.FootageDimensions[0],
				Height:    item.FootageDimensions[1],
				Framerate: item.FootageFramerate,
				Duration:  item.FootageSeconds,
			}
			
			if media.Framerate > 0 && media.Duration > 0 {
				media.Frames = int(media.Duration * media.Framerate)
			}
			
			// Determine type and update stats
			name := item.Name
			if strings.Contains(strings.ToLower(name), "null") {
				report.Overview.Statistics.NullCount++
				media.Type = "Null Object"
			} else if strings.Contains(name, "Adjustment Layer") {
				report.Overview.Statistics.AdjustmentCount++
				media.Type = "Adjustment"
			} else if strings.Contains(name, "Solid") {
				report.Overview.Statistics.SolidCount++
				media.Type = "Solid"
			} else if item.FootageType == aep.FootageTypeSolid {
				report.Overview.Statistics.SolidCount++
				media.Type = "Solid"
			} else if item.FootageType == aep.FootageTypePlaceholder {
				media.Type = "Placeholder"
			} else if media.Framerate > 0 && media.Duration > 0 {
				report.Overview.Statistics.VideoCount++
				media.Type = "Video"
			} else if media.Width > 0 && media.Height > 0 {
				report.Overview.Statistics.ImageCount++
				media.Type = "Image"
			} else {
				report.Overview.Statistics.OtherCount++
				media.Type = "Other"
			}
			
			report.MediaAssets = append(report.MediaAssets, media)
		}
	}
	
	// Update media usage information
	for i := range report.MediaAssets {
		media := &report.MediaAssets[i]
		if comps, exists := mediaUsage[media.ID]; exists {
			media.UsedIn = unique(comps)
			media.UsageCount = len(media.UsedIn)
		}
	}
	
	// Sort compositions by name
	sort.Slice(allComps, func(i, j int) bool {
		return allComps[i].Name < allComps[j].Name
	})
	report.Compositions = allComps
	
	// Sort layers by composition then index
	sort.Slice(allLayers, func(i, j int) bool {
		if allLayers[i].CompName == allLayers[j].CompName {
			return allLayers[i].Index < allLayers[j].Index
		}
		return allLayers[i].CompName < allLayers[j].CompName
	})
	report.AllLayers = allLayers
	
	// Build folder hierarchy
	report.FolderTree = buildFolderTree(project.RootFolder)
	
	return report
}

func buildFolderTree(folder *aep.Item) FolderNode {
	if folder == nil {
		return FolderNode{}
	}
	
	node := FolderNode{
		Name:      folder.Name,
		Type:      "folder",
		ItemCount: len(folder.FolderContents),
	}
	
	if node.Name == "root" {
		node.Name = "📁 Project Root"
	}
	
	for _, item := range folder.FolderContents {
		if item.ItemType == aep.ItemTypeFolder {
			childNode := buildFolderTree(item)
			node.Children = append(node.Children, childNode)
		} else {
			childNode := FolderNode{
				Name: item.Name,
			}
			if item.ItemType == aep.ItemTypeComposition {
				childNode.Type = "comp"
			} else {
				childNode.Type = "footage"
			}
			node.Children = append(node.Children, childNode)
		}
	}
	
	return node
}

func unique(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

func mustGetWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}