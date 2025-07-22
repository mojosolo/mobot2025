package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	aep "github.com/yourusername/mobot2025"
)

// ReportData holds all data for the HTML report
type ReportData struct {
	Filename          string
	ParsedAt          string
	BitDepth          string
	ExpressionEngine  string
	TotalItems        int
	FolderCount       int
	CompositionCount  int
	FootageCount      int
	TotalLayers       int
	Compositions      []CompositionInfo
	MainCompositions  []CompositionInfo
	SceneCompositions []CompositionInfo
	FootageItems      []FootageInfo
	Statistics        ProjectStats
}

// CompositionInfo holds composition details
type CompositionInfo struct {
	Name       string
	Width      uint16
	Height     uint16
	Resolution string
	Framerate  float64
	Duration   float64
	Layers     int
	BGColor    string
}

// FootageInfo holds footage details
type FootageInfo struct {
	Name      string
	Width     uint16
	Height    uint16
	Framerate float64
	Duration  float64
	Type      string
}

// ProjectStats holds project statistics
type ProjectStats struct {
	NullCount       int
	AdjustmentCount int
	SolidCount      int
	OtherCount      int
	HasLogo         bool
	HasVideo        bool
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AEP Analysis Report - {{.Filename}}</title>
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
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px 0;
            text-align: center;
            margin-bottom: 30px;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .subtitle {
            font-size: 1.2em;
            opacity: 0.9;
        }
        
        .section {
            background: white;
            padding: 30px;
            margin-bottom: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        
        h2 {
            color: #667eea;
            margin-bottom: 20px;
            font-size: 1.8em;
            border-bottom: 2px solid #eee;
            padding-bottom: 10px;
        }
        
        h3 {
            color: #555;
            margin-bottom: 15px;
            font-size: 1.3em;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .stat-card {
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
            padding: 20px;
            border-radius: 8px;
            text-align: center;
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
        
        .composition-card {
            background: #f8f9fa;
            padding: 20px;
            margin-bottom: 15px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        
        .composition-title {
            font-size: 1.2em;
            font-weight: bold;
            color: #333;
            margin-bottom: 10px;
        }
        
        .composition-details {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 10px;
            font-size: 0.9em;
            color: #666;
        }
        
        .detail-item {
            display: flex;
            align-items: center;
        }
        
        .detail-label {
            font-weight: 600;
            margin-right: 5px;
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #eee;
        }
        
        th {
            background: #f8f9fa;
            font-weight: 600;
            color: #667eea;
        }
        
        tr:hover {
            background: #f8f9fa;
        }
        
        .badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.8em;
            font-weight: 600;
            margin-left: 10px;
        }
        
        .badge-4k {
            background: #10b981;
            color: white;
        }
        
        .badge-hd {
            background: #3b82f6;
            color: white;
        }
        
        .badge-3d {
            background: #f59e0b;
            color: white;
        }
        
        .footer {
            text-align: center;
            color: #666;
            margin-top: 40px;
            padding: 20px;
        }
        
        .icon {
            display: inline-block;
            width: 20px;
            margin-right: 5px;
            vertical-align: middle;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🎬 AEP Analysis Report</h1>
            <div class="subtitle">{{.Filename}}</div>
            <div style="margin-top: 10px; opacity: 0.8;">Generated on {{.ParsedAt}}</div>
        </header>
        
        <div class="section">
            <h2>📊 Project Overview</h2>
            <div class="stats-grid">
                <div class="stat-card">
                    <span class="stat-number">{{.TotalItems}}</span>
                    <span class="stat-label">Total Items</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.CompositionCount}}</span>
                    <span class="stat-label">Compositions</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.FootageCount}}</span>
                    <span class="stat-label">Footage Items</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.TotalLayers}}</span>
                    <span class="stat-label">Total Layers</span>
                </div>
            </div>
            
            <div style="margin-top: 20px;">
                <p><strong>Bit Depth:</strong> {{.BitDepth}}</p>
                <p><strong>Expression Engine:</strong> {{.ExpressionEngine}}</p>
                <p><strong>Folder Structure:</strong> {{.FolderCount}} folders</p>
            </div>
        </div>
        
        {{if .MainCompositions}}
        <div class="section">
            <h2>🎯 Main Compositions</h2>
            {{range .MainCompositions}}
            <div class="composition-card">
                <div class="composition-title">
                    {{.Name}}
                    {{if eq .Resolution "4K"}}<span class="badge badge-4k">4K</span>{{end}}
                    {{if eq .Resolution "HD"}}<span class="badge badge-hd">HD</span>{{end}}
                </div>
                <div class="composition-details">
                    <div class="detail-item">
                        <span class="detail-label">Resolution:</span> {{.Width}}×{{.Height}}
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Framerate:</span> {{printf "%.2f" .Framerate}} fps
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Duration:</span> {{printf "%.2f" .Duration}}s
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">Layers:</span> {{.Layers}}
                    </div>
                </div>
            </div>
            {{end}}
        </div>
        {{end}}
        
        {{if .SceneCompositions}}
        <div class="section">
            <h2>🎬 Scene Compositions</h2>
            <table>
                <thead>
                    <tr>
                        <th>Scene</th>
                        <th>Resolution</th>
                        <th>Framerate</th>
                        <th>Duration</th>
                        <th>Layers</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .SceneCompositions}}
                    <tr>
                        <td><strong>{{.Name}}</strong></td>
                        <td>{{.Width}}×{{.Height}}</td>
                        <td>{{printf "%.2f" .Framerate}} fps</td>
                        <td>{{printf "%.2f" .Duration}}s</td>
                        <td>{{.Layers}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}
        
        <div class="section">
            <h2>📦 Asset Breakdown</h2>
            <div class="stats-grid">
                <div class="stat-card">
                    <span class="stat-number">{{.Statistics.NullCount}}</span>
                    <span class="stat-label">Null Objects</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Statistics.AdjustmentCount}}</span>
                    <span class="stat-label">Adjustment Layers</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Statistics.SolidCount}}</span>
                    <span class="stat-label">Solid Layers</span>
                </div>
                <div class="stat-card">
                    <span class="stat-number">{{.Statistics.OtherCount}}</span>
                    <span class="stat-label">Other Assets</span>
                </div>
            </div>
            
            {{if or .Statistics.HasLogo .Statistics.HasVideo}}
            <div style="margin-top: 20px;">
                <h3>Detected Asset Types:</h3>
                <ul style="list-style: none; padding-left: 0;">
                    {{if .Statistics.HasLogo}}<li>✅ Logo/Image Assets</li>{{end}}
                    {{if .Statistics.HasVideo}}<li>✅ Video Footage</li>{{end}}
                </ul>
            </div>
            {{end}}
        </div>
        
        {{if gt (len .FootageItems) 0}}
        <div class="section">
            <h2>📹 Sample Footage Items</h2>
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Type</th>
                        <th>Dimensions</th>
                        <th>Framerate</th>
                        <th>Duration</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .FootageItems}}
                    <tr>
                        <td><strong>{{.Name}}</strong></td>
                        <td>{{.Type}}</td>
                        <td>{{if gt .Width 0}}{{.Width}}×{{.Height}}{{else}}-{{end}}</td>
                        <td>{{if gt .Framerate 0.0}}{{printf "%.2f" .Framerate}} fps{{else}}-{{end}}</td>
                        <td>{{if gt .Duration 0.0}}{{printf "%.2f" .Duration}}s{{else}}-{{end}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}
        
        <div class="footer">
            <p>Generated by mobot2025 AEP Parser</p>
            <p>github.com/yourusername/mobot2025</p>
        </div>
    </div>
</body>
</html>`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_html_report.go <aep-file>")
		fmt.Println("Example: go run generate_html_report.go sample-aep/Ai\\ Text\\ Intro.aep")
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
	
	// Prepare report data
	report := prepareReportData(aepFile, project)
	
	// Generate HTML filename
	baseFilename := strings.TrimSuffix(filepath.Base(aepFile), filepath.Ext(aepFile))
	htmlFilename := fmt.Sprintf("%s-report-%s.html", 
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
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("❌ Error parsing template: %v\n", err)
		os.Exit(1)
	}
	
	err = tmpl.Execute(file, report)
	if err != nil {
		fmt.Printf("❌ Error generating HTML: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ Report generated: %s\n", htmlFilename)
	fmt.Printf("📂 Open in browser: file://%s/%s\n", 
		mustGetWorkingDir(), htmlFilename)
}

func prepareReportData(filename string, project *aep.Project) ReportData {
	report := ReportData{
		Filename:         filepath.Base(filename),
		ParsedAt:         time.Now().Format("January 2, 2006 at 3:04 PM"),
		TotalItems:       len(project.Items),
		ExpressionEngine: project.ExpressionEngine,
	}
	
	// Convert bit depth
	switch project.Depth {
	case aep.BPC8:
		report.BitDepth = "8-bit"
	case aep.BPC16:
		report.BitDepth = "16-bit"
	case aep.BPC32:
		report.BitDepth = "32-bit (float)"
	default:
		report.BitDepth = fmt.Sprintf("Unknown (%v)", project.Depth)
	}
	
	// Count items and collect compositions
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeFolder:
			report.FolderCount++
		case aep.ItemTypeComposition:
			report.CompositionCount++
			report.TotalLayers += len(item.CompositionLayers)
			
			comp := CompositionInfo{
				Name:      item.Name,
				Width:     item.FootageDimensions[0],
				Height:    item.FootageDimensions[1],
				Framerate: item.FootageFramerate,
				Duration:  item.FootageSeconds,
				Layers:    len(item.CompositionLayers),
				BGColor:   fmt.Sprintf("rgb(%d,%d,%d)", 
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
			
			// Categorize compositions
			if strings.Contains(item.Name, "Final Comp") {
				report.MainCompositions = append(report.MainCompositions, comp)
			} else if strings.HasPrefix(item.Name, "S0") || 
					  (len(item.Name) == 3 && item.Name[0] == 'S' && 
					   item.Name[1] >= '0' && item.Name[1] <= '9') {
				report.SceneCompositions = append(report.SceneCompositions, comp)
			}
			
			report.Compositions = append(report.Compositions, comp)
			
		case aep.ItemTypeFootage:
			report.FootageCount++
			
			// Analyze footage type
			name := item.Name
			if strings.Contains(strings.ToLower(name), "null") {
				report.Statistics.NullCount++
			} else if strings.Contains(name, "Adjustment Layer") {
				report.Statistics.AdjustmentCount++
			} else if strings.Contains(name, "Solid") {
				report.Statistics.SolidCount++
			} else {
				report.Statistics.OtherCount++
				if strings.Contains(strings.ToLower(name), "logo") {
					report.Statistics.HasLogo = true
				}
				if item.FootageFramerate > 0 && item.FootageSeconds > 0 {
					report.Statistics.HasVideo = true
				}
			}
			
			// Add to footage list (limit to first 20)
			if len(report.FootageItems) < 20 {
				footage := FootageInfo{
					Name:      item.Name,
					Width:     item.FootageDimensions[0],
					Height:    item.FootageDimensions[1],
					Framerate: item.FootageFramerate,
					Duration:  item.FootageSeconds,
				}
				
				// Determine type
				switch item.FootageType {
				case aep.FootageTypeSolid:
					footage.Type = "Solid"
				case aep.FootageTypePlaceholder:
					footage.Type = "Placeholder"
				default:
					if strings.Contains(strings.ToLower(name), "null") {
						footage.Type = "Null Object"
					} else if strings.Contains(name, "Adjustment") {
						footage.Type = "Adjustment"
					} else if footage.Framerate > 0 {
						footage.Type = "Video"
					} else if footage.Width > 0 {
						footage.Type = "Image"
					} else {
						footage.Type = "Other"
					}
				}
				
				report.FootageItems = append(report.FootageItems, footage)
			}
		}
	}
	
	return report
}

func mustGetWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}