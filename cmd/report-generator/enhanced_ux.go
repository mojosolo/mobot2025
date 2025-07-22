package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	
	aep "github.com/yourusername/mobot2025"
)

// Core type definitions for report structure
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

type MediaAsset struct {
	ID         uint32
	Name       string
	Type       string
	Width      uint16
	Height     uint16
	Framerate  float64
	Duration   float64
	Frames     int
	UsedIn     []string
	UsageCount int
}

type TextLayer struct {
	LayerName string
	CompName  string
	Text      string
	Font      string
	Size      float64
}

type LayerAttribute struct {
	LayerName     string
	CompName      string
	Is3D          bool
	IsSolo        bool
	IsShy         bool
	IsLocked      bool
	IsGuide       bool
	IsAdjustment  bool
	HasMotionBlur bool
	HasEffects    bool
}

type FolderNode struct {
	Name      string
	Type      string
	ItemCount int
	Children  []FolderNode
}

// EnhancedReportData holds comprehensive data with enhanced UX features
type EnhancedReportData struct {
	Filename         string
	ParsedAt         string
	BitDepth         string
	ExpressionEngine string
	TotalItems       int
	FileSize         string
	FilePath         string
	
	// Enhanced Overview with more metrics
	Overview struct {
		FolderCount      int
		CompositionCount int
		FootageCount     int
		TotalLayers      int
		TotalEffects     int
		TotalKeyframes   int
		TotalExpressions int
		MainComps        []CompositionDetail
		Statistics       EnhancedProjectStats
		Warnings         []string
		Insights         []string
	}
	
	// All existing tabs plus enhanced features
	Compositions    []CompositionDetail
	AllLayers       []LayerDetail
	MediaAssets     []MediaAsset
	TextLayers      []TextLayer
	LayerAttributes []LayerAttribute
	FolderTree      FolderNode
	
	// New enhanced sections
	Effects         []EffectDetail
	Expressions     []ExpressionDetail
	Keyframes       []KeyframeDetail
	ProjectInsights ProjectInsights
	RawData         []RawDataSection
}

// Enhanced structures for comprehensive content
type EnhancedProjectStats struct {
	ColorDepth       string
	TotalFrames      int
	TotalDuration    float64
	AverageFramerate float64
	MemoryUsage      string
	ComplexityScore  int
	UsedEffects      map[string]int
	AssetTypes       map[string]int
	LayerTypes       map[string]int
}

type EffectDetail struct {
	LayerName    string
	CompName     string
	EffectName   string
	MatchName    string
	Parameters   []EffectParameter
	IsEnabled    bool
	RenderOrder  int
}

type EffectParameter struct {
	Name     string
	Value    string
	Type     string
	Animated bool
}

type ExpressionDetail struct {
	LayerName    string
	CompName     string
	PropertyPath string
	Expression   string
	Language     string
	HasErrors    bool
	LineCount    int
}

type KeyframeDetail struct {
	LayerName      string
	CompName       string
	PropertyPath   string
	Time           float64
	Value          string
	InterpolationType string
	EaseIn         float32
	EaseOut        float32
}

type ProjectInsights struct {
	PerformanceMetrics  []string
	OptimizationTips    []string
	CompatibilityNotes  []string
	WorkflowSuggestions []string
	AssetRecommendations []string
}

type RawDataSection struct {
	PropertyPath string
	DataType     string
	Size         int
	Preview      string
	HexDump      string
}

const enhancedHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Enhanced AEP Analysis - {{.Filename}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        :root {
            --primary-color: #667eea;
            --primary-dark: #5a67d8;
            --primary-light: #7f9cf5;
            --secondary-color: #ed64a6;
            --success-color: #48bb78;
            --warning-color: #f6ad55;
            --danger-color: #fc8181;
            --info-color: #63b3ed;
            --dark-bg: #1a202c;
            --light-bg: #f7fafc;
            --text-primary: #2d3748;
            --text-secondary: #718096;
            --border-color: #e2e8f0;
            --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
            --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
            --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
            --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: var(--text-primary);
            background-color: var(--light-bg);
            overflow-x: hidden;
        }
        
        /* Enhanced Header with Gradient and Animation */
        .header {
            background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
            color: white;
            padding: 40px 30px;
            text-align: center;
            box-shadow: var(--shadow-lg);
            position: relative;
            overflow: hidden;
        }
        
        .header::before {
            content: '';
            position: absolute;
            top: -50%;
            left: -50%;
            width: 200%;
            height: 200%;
            background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
            animation: pulse 4s ease-in-out infinite;
        }
        
        @keyframes pulse {
            0%, 100% { transform: scale(1); opacity: 0.5; }
            50% { transform: scale(1.1); opacity: 0.3; }
        }
        
        .header h1 {
            font-size: 3em;
            margin-bottom: 15px;
            position: relative;
            z-index: 1;
            font-weight: 700;
            letter-spacing: -1px;
        }
        
        .subtitle {
            opacity: 0.95;
            font-size: 1.2em;
            position: relative;
            z-index: 1;
            font-weight: 300;
        }
        
        /* File Info Bar */
        .file-info {
            background: white;
            padding: 20px 30px;
            box-shadow: var(--shadow-sm);
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 15px;
        }
        
        .file-info-item {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .file-info-icon {
            font-size: 1.2em;
            opacity: 0.7;
        }
        
        /* Enhanced Tabs with Icons and Badges */
        .tabs {
            background: white;
            padding: 0;
            box-shadow: var(--shadow-md);
            position: sticky;
            top: 0;
            z-index: 100;
            border-bottom: 1px solid var(--border-color);
        }
        
        .tab-list {
            display: flex;
            list-style: none;
            margin: 0;
            padding: 0 20px;
            overflow-x: auto;
            -webkit-overflow-scrolling: touch;
        }
        
        .tab-button {
            padding: 18px 25px;
            background: none;
            border: none;
            font-size: 1em;
            cursor: pointer;
            color: var(--text-secondary);
            transition: all 0.3s ease;
            border-bottom: 3px solid transparent;
            font-weight: 500;
            display: flex;
            align-items: center;
            gap: 8px;
            white-space: nowrap;
            position: relative;
        }
        
        .tab-icon {
            font-size: 1.1em;
            opacity: 0.8;
        }
        
        .tab-badge {
            background: var(--primary-light);
            color: white;
            font-size: 0.75em;
            padding: 2px 6px;
            border-radius: 10px;
            font-weight: 600;
        }
        
        .tab-button:hover {
            color: var(--primary-color);
            background: rgba(102, 126, 234, 0.05);
        }
        
        .tab-button.active {
            color: var(--primary-color);
            border-bottom-color: var(--primary-color);
            background: rgba(102, 126, 234, 0.08);
        }
        
        /* Tab Content with Smooth Transitions */
        .tab-content {
            display: none;
            padding: 30px;
            animation: fadeIn 0.4s ease;
            min-height: 500px;
        }
        
        .tab-content.active {
            display: block;
        }
        
        @keyframes fadeIn {
            from { 
                opacity: 0; 
                transform: translateY(10px);
            }
            to { 
                opacity: 1; 
                transform: translateY(0);
            }
        }
        
        /* Container and Layout */
        .container {
            max-width: 1600px;
            margin: 0 auto;
        }
        
        /* Enhanced Cards with Hover Effects */
        .card {
            background: white;
            border-radius: 12px;
            padding: 25px;
            margin-bottom: 25px;
            box-shadow: var(--shadow-sm);
            transition: all 0.3s ease;
            border: 1px solid var(--border-color);
        }
        
        .card:hover {
            box-shadow: var(--shadow-lg);
            transform: translateY(-2px);
        }
        
        .card h3 {
            color: var(--text-primary);
            margin-bottom: 20px;
            font-size: 1.5em;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .card-icon {
            font-size: 1.2em;
            color: var(--primary-color);
        }
        
        /* Enhanced Stats Grid */
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 25px;
            margin-bottom: 30px;
        }
        
        .stat-card {
            background: linear-gradient(135deg, #f6f9fc 0%, #ffffff 100%);
            border-radius: 12px;
            padding: 25px;
            box-shadow: var(--shadow-sm);
            transition: all 0.3s ease;
            border: 1px solid var(--border-color);
            position: relative;
            overflow: hidden;
        }
        
        .stat-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 4px;
            height: 100%;
            background: var(--primary-color);
            opacity: 0;
            transition: opacity 0.3s ease;
        }
        
        .stat-card:hover {
            transform: translateY(-3px);
            box-shadow: var(--shadow-lg);
        }
        
        .stat-card:hover::before {
            opacity: 1;
        }
        
        .stat-label {
            color: var(--text-secondary);
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 8px;
            font-weight: 600;
        }
        
        .stat-value {
            color: var(--text-primary);
            font-size: 2em;
            font-weight: 700;
            line-height: 1.2;
        }
        
        .stat-icon {
            position: absolute;
            top: 20px;
            right: 20px;
            font-size: 2.5em;
            opacity: 0.1;
        }
        
        /* Enhanced Tables with Sorting and Filtering */
        .table-controls {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            flex-wrap: wrap;
            gap: 15px;
        }
        
        .search-box {
            position: relative;
            flex: 1;
            max-width: 400px;
        }
        
        .search-icon {
            position: absolute;
            left: 15px;
            top: 50%;
            transform: translateY(-50%);
            color: var(--text-secondary);
            font-size: 1.1em;
        }
        
        .search-input {
            width: 100%;
            padding: 12px 15px 12px 45px;
            border: 2px solid var(--border-color);
            border-radius: 8px;
            font-size: 1em;
            transition: all 0.3s ease;
            background: white;
        }
        
        .search-input:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        .table-container {
            overflow-x: auto;
            border-radius: 8px;
            box-shadow: var(--shadow-sm);
            border: 1px solid var(--border-color);
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
            background: white;
        }
        
        th {
            background: var(--light-bg);
            padding: 15px;
            text-align: left;
            font-weight: 600;
            color: var(--text-primary);
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            border-bottom: 2px solid var(--border-color);
            cursor: pointer;
            user-select: none;
            position: sticky;
            top: 0;
            z-index: 10;
        }
        
        th:hover {
            background: #e2e8f0;
        }
        
        .sort-icon {
            float: right;
            opacity: 0.5;
            font-size: 0.9em;
        }
        
        td {
            padding: 15px;
            border-bottom: 1px solid var(--border-color);
            color: var(--text-primary);
        }
        
        tr:hover {
            background: rgba(102, 126, 234, 0.03);
        }
        
        tr:last-child td {
            border-bottom: none;
        }
        
        /* Enhanced Tree View */
        .tree-container {
            background: white;
            border-radius: 8px;
            padding: 25px;
            box-shadow: var(--shadow-sm);
            border: 1px solid var(--border-color);
        }
        
        .tree-node {
            padding: 10px 0 10px 25px;
            border-left: 2px solid var(--border-color);
            margin-left: 10px;
            position: relative;
            transition: all 0.3s ease;
        }
        
        .tree-node::before {
            content: '';
            position: absolute;
            left: -2px;
            top: 25px;
            width: 25px;
            height: 2px;
            background: var(--border-color);
        }
        
        .tree-node:hover {
            background: rgba(102, 126, 234, 0.05);
            border-left-color: var(--primary-color);
        }
        
        .tree-node-root {
            border-left: none;
            margin-left: 0;
            padding-left: 0;
            font-weight: 600;
            font-size: 1.1em;
        }
        
        .tree-icon {
            margin-right: 10px;
            font-size: 1.2em;
        }
        
        .tree-count {
            color: var(--text-secondary);
            font-size: 0.9em;
            margin-left: 10px;
        }
        
        /* Status Badges */
        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.85em;
            font-weight: 500;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .badge-success {
            background: rgba(72, 187, 120, 0.1);
            color: var(--success-color);
            border: 1px solid rgba(72, 187, 120, 0.3);
        }
        
        .badge-warning {
            background: rgba(246, 173, 85, 0.1);
            color: var(--warning-color);
            border: 1px solid rgba(246, 173, 85, 0.3);
        }
        
        .badge-danger {
            background: rgba(252, 129, 129, 0.1);
            color: var(--danger-color);
            border: 1px solid rgba(252, 129, 129, 0.3);
        }
        
        .badge-info {
            background: rgba(99, 179, 237, 0.1);
            color: var(--info-color);
            border: 1px solid rgba(99, 179, 237, 0.3);
        }
        
        /* Progress Bars */
        .progress-container {
            margin: 20px 0;
        }
        
        .progress-label {
            display: flex;
            justify-content: space-between;
            margin-bottom: 8px;
            color: var(--text-secondary);
            font-size: 0.9em;
        }
        
        .progress-bar {
            height: 8px;
            background: var(--border-color);
            border-radius: 4px;
            overflow: hidden;
        }
        
        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, var(--primary-color) 0%, var(--primary-light) 100%);
            border-radius: 4px;
            transition: width 0.6s ease;
            position: relative;
            overflow: hidden;
        }
        
        .progress-fill::after {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            bottom: 0;
            right: 0;
            background: linear-gradient(
                90deg,
                transparent,
                rgba(255, 255, 255, 0.3),
                transparent
            );
            animation: shimmer 2s infinite;
        }
        
        @keyframes shimmer {
            0% { transform: translateX(-100%); }
            100% { transform: translateX(100%); }
        }
        
        /* Alert Boxes */
        .alert {
            padding: 15px 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 15px;
            font-size: 0.95em;
        }
        
        .alert-icon {
            font-size: 1.3em;
            flex-shrink: 0;
        }
        
        .alert-info {
            background: rgba(99, 179, 237, 0.1);
            color: var(--info-color);
            border: 1px solid rgba(99, 179, 237, 0.3);
        }
        
        .alert-warning {
            background: rgba(246, 173, 85, 0.1);
            color: var(--warning-color);
            border: 1px solid rgba(246, 173, 85, 0.3);
        }
        
        .alert-success {
            background: rgba(72, 187, 120, 0.1);
            color: var(--success-color);
            border: 1px solid rgba(72, 187, 120, 0.3);
        }
        
        /* Code Display */
        .code-block {
            background: #2d3748;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.9em;
            line-height: 1.5;
            margin: 15px 0;
        }
        
        .code-line {
            display: block;
            padding: 2px 0;
        }
        
        .code-line:hover {
            background: rgba(255, 255, 255, 0.05);
        }
        
        /* Hex Dump Display */
        .hex-dump {
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            font-size: 0.85em;
            line-height: 1.6;
            background: #f7fafc;
            padding: 15px;
            border-radius: 8px;
            overflow-x: auto;
            border: 1px solid var(--border-color);
        }
        
        .hex-line {
            display: flex;
            gap: 20px;
            padding: 4px 0;
        }
        
        .hex-offset {
            color: var(--text-secondary);
            font-weight: 600;
            width: 80px;
        }
        
        .hex-bytes {
            color: var(--primary-color);
            flex: 1;
        }
        
        .hex-ascii {
            color: var(--text-primary);
            width: 150px;
        }
        
        /* Insights and Tips */
        .insight-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin: 20px 0;
        }
        
        .insight-card {
            background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(237, 100, 166, 0.05) 100%);
            border-radius: 12px;
            padding: 20px;
            border: 1px solid rgba(102, 126, 234, 0.2);
        }
        
        .insight-title {
            font-weight: 600;
            color: var(--primary-color);
            margin-bottom: 10px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .insight-content {
            color: var(--text-primary);
            font-size: 0.95em;
            line-height: 1.6;
        }
        
        /* Tooltips */
        .tooltip {
            position: relative;
            cursor: help;
            border-bottom: 1px dashed var(--text-secondary);
        }
        
        .tooltip-content {
            position: absolute;
            bottom: 100%;
            left: 50%;
            transform: translateX(-50%);
            background: var(--dark-bg);
            color: white;
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 0.85em;
            white-space: nowrap;
            opacity: 0;
            pointer-events: none;
            transition: opacity 0.3s ease;
            margin-bottom: 5px;
            z-index: 1000;
        }
        
        .tooltip:hover .tooltip-content {
            opacity: 1;
        }
        
        /* Footer */
        .footer {
            background: var(--dark-bg);
            color: white;
            text-align: center;
            padding: 30px;
            margin-top: 50px;
        }
        
        .footer p {
            margin: 5px 0;
            opacity: 0.8;
        }
        
        .footer a {
            color: var(--primary-light);
            text-decoration: none;
        }
        
        .footer a:hover {
            text-decoration: underline;
        }
        
        /* Responsive Design */
        @media (max-width: 768px) {
            .header h1 {
                font-size: 2em;
            }
            
            .tab-list {
                padding: 0 10px;
            }
            
            .tab-button {
                padding: 15px 20px;
                font-size: 0.9em;
            }
            
            .stats-grid {
                grid-template-columns: 1fr;
            }
            
            .table-container {
                font-size: 0.9em;
            }
            
            th, td {
                padding: 10px;
            }
        }
        
        /* Print Styles */
        @media print {
            .tabs, .search-box, .footer {
                display: none;
            }
            
            .tab-content {
                display: block !important;
                page-break-inside: avoid;
            }
            
            .card {
                box-shadow: none;
                border: 1px solid #ddd;
            }
        }
        
        /* Loading Animation */
        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid rgba(102, 126, 234, 0.3);
            border-radius: 50%;
            border-top-color: var(--primary-color);
            animation: spin 1s ease-in-out infinite;
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        
        /* Custom Scrollbar */
        ::-webkit-scrollbar {
            width: 10px;
            height: 10px;
        }
        
        ::-webkit-scrollbar-track {
            background: var(--light-bg);
        }
        
        ::-webkit-scrollbar-thumb {
            background: var(--primary-color);
            border-radius: 5px;
        }
        
        ::-webkit-scrollbar-thumb:hover {
            background: var(--primary-dark);
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Filename}}</h1>
        <div class="subtitle">Enhanced After Effects Project Analysis Report</div>
    </div>
    
    <div class="file-info">
        <div class="file-info-item">
            <span class="file-info-icon">📅</span>
            <span>{{.ParsedAt}}</span>
        </div>
        <div class="file-info-item">
            <span class="file-info-icon">🎨</span>
            <span>{{.BitDepth}} Color Depth</span>
        </div>
        <div class="file-info-item">
            <span class="file-info-icon">⚙️</span>
            <span>{{.ExpressionEngine}} Engine</span>
        </div>
        <div class="file-info-item">
            <span class="file-info-icon">📦</span>
            <span>{{.TotalItems}} Total Items</span>
        </div>
        {{if .FileSize}}
        <div class="file-info-item">
            <span class="file-info-icon">💾</span>
            <span>{{.FileSize}}</span>
        </div>
        {{end}}
    </div>
    
    <div class="tabs">
        <ul class="tab-list">
            <li>
                <button class="tab-button active" onclick="showTab('overview')">
                    <span class="tab-icon">📊</span>
                    Overview
                    <span class="tab-badge">{{.Overview.CompositionCount}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('compositions')">
                    <span class="tab-icon">🎬</span>
                    Compositions
                    <span class="tab-badge">{{len .Compositions}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('layers')">
                    <span class="tab-icon">📚</span>
                    Layers
                    <span class="tab-badge">{{len .AllLayers}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('media')">
                    <span class="tab-icon">🎞️</span>
                    Media
                    <span class="tab-badge">{{len .MediaAssets}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('text')">
                    <span class="tab-icon">🔤</span>
                    Text
                    <span class="tab-badge">{{len .TextLayers}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('effects')">
                    <span class="tab-icon">✨</span>
                    Effects
                    <span class="tab-badge">{{len .Effects}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('expressions')">
                    <span class="tab-icon">🧮</span>
                    Expressions
                    <span class="tab-badge">{{len .Expressions}}</span>
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('attributes')">
                    <span class="tab-icon">🎛️</span>
                    Attributes
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('hierarchy')">
                    <span class="tab-icon">🌳</span>
                    Hierarchy
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('insights')">
                    <span class="tab-icon">💡</span>
                    Insights
                </button>
            </li>
            <li>
                <button class="tab-button" onclick="showTab('raw')">
                    <span class="tab-icon">🔧</span>
                    Raw Data
                </button>
            </li>
        </ul>
    </div>
    
    <div class="container">
        <!-- Overview Tab -->
        <div id="overview" class="tab-content active">
            <div class="stats-grid">
                <div class="stat-card">
                    <div class="stat-icon">📁</div>
                    <div class="stat-label">Folders</div>
                    <div class="stat-value">{{.Overview.FolderCount}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">🎬</div>
                    <div class="stat-label">Compositions</div>
                    <div class="stat-value">{{.Overview.CompositionCount}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">🎞️</div>
                    <div class="stat-label">Footage Items</div>
                    <div class="stat-value">{{.Overview.FootageCount}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">📚</div>
                    <div class="stat-label">Total Layers</div>
                    <div class="stat-value">{{.Overview.TotalLayers}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">✨</div>
                    <div class="stat-label">Total Effects</div>
                    <div class="stat-value">{{.Overview.TotalEffects}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">🔑</div>
                    <div class="stat-label">Total Keyframes</div>
                    <div class="stat-value">{{.Overview.TotalKeyframes}}</div>
                </div>
            </div>
            
            <!-- Project Insights -->
            {{if .Overview.Insights}}
            <div class="card">
                <h3><span class="card-icon">💡</span>Project Insights</h3>
                <div class="insight-grid">
                    {{range .Overview.Insights}}
                    <div class="insight-card">
                        <div class="insight-content">{{.}}</div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
            
            <!-- Warnings -->
            {{if .Overview.Warnings}}
            <div class="card">
                <h3><span class="card-icon">⚠️</span>Warnings</h3>
                {{range .Overview.Warnings}}
                <div class="alert alert-warning">
                    <span class="alert-icon">⚠️</span>
                    <span>{{.}}</span>
                </div>
                {{end}}
            </div>
            {{end}}
            
            <!-- Main Compositions -->
            {{if .Overview.MainComps}}
            <div class="card">
                <h3><span class="card-icon">🎬</span>Main Compositions</h3>
                <div class="table-container">
                    <table>
                        <thead>
                            <tr>
                                <th onclick="sortTable(this, 0)">Name <span class="sort-icon">↕</span></th>
                                <th onclick="sortTable(this, 1)">Resolution <span class="sort-icon">↕</span></th>
                                <th onclick="sortTable(this, 2)">Framerate <span class="sort-icon">↕</span></th>
                                <th onclick="sortTable(this, 3)">Duration <span class="sort-icon">↕</span></th>
                                <th onclick="sortTable(this, 4)">Layers <span class="sort-icon">↕</span></th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Overview.MainComps}}
                            <tr>
                                <td><strong>{{.Name}}</strong></td>
                                <td>{{.Width}}×{{.Height}} <span class="badge badge-info">{{.Resolution}}</span></td>
                                <td>{{.Framerate}} fps</td>
                                <td>{{.Duration}}s ({{.Frames}} frames)</td>
                                <td>{{.LayerCount}} layers</td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
            </div>
            {{end}}
            
            <!-- Project Statistics -->
            <div class="card">
                <h3><span class="card-icon">📊</span>Project Statistics</h3>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-label">Color Space</div>
                        <div class="stat-value" style="font-size: 1.2em;">{{.Overview.Statistics.ColorDepth}}</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Total Frames</div>
                        <div class="stat-value" style="font-size: 1.5em;">{{.Overview.Statistics.TotalFrames}}</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Total Duration</div>
                        <div class="stat-value" style="font-size: 1.5em;">{{printf "%.1f" .Overview.Statistics.TotalDuration}}s</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Avg Framerate</div>
                        <div class="stat-value" style="font-size: 1.5em;">{{printf "%.1f" .Overview.Statistics.AverageFramerate}} fps</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Complexity Score</div>
                        <div class="stat-value" style="font-size: 1.5em;">{{.Overview.Statistics.ComplexityScore}}/100</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-label">Memory Usage</div>
                        <div class="stat-value" style="font-size: 1.5em;">{{.Overview.Statistics.MemoryUsage}}</div>
                    </div>
                </div>
                
                <!-- Asset Type Distribution -->
                {{if .Overview.Statistics.AssetTypes}}
                <h4 style="margin-top: 30px;">Asset Distribution</h4>
                <div class="progress-container">
                    {{range $type, $count := .Overview.Statistics.AssetTypes}}
                    <div class="progress-label">
                        <span>{{$type}}</span>
                        <span>{{$count}} items</span>
                    </div>
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: {{calcPercent $count $.TotalItems}}%"></div>
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>
        
        <!-- Compositions Tab -->
        <div id="compositions" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="compSearch" placeholder="Search compositions..." onkeyup="filterTable('compSearch', 'compTable')">
                </div>
                <div>
                    <span class="badge badge-info">{{len .Compositions}} Total</span>
                </div>
            </div>
            
            <div class="table-container">
                <table id="compTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Name <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Resolution <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)">Framerate <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Duration <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Layers <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 5)">Background <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Compositions}}
                        <tr>
                            <td><strong>{{.Name}}</strong></td>
                            <td>
                                {{.Width}}×{{.Height}}
                                <span class="badge badge-info">{{.Resolution}}</span>
                            </td>
                            <td>{{.Framerate}} fps</td>
                            <td>{{.Duration}}s ({{.Frames}} frames)</td>
                            <td>
                                {{.LayerCount}}
                                {{if .LayerCount}}
                                    {{if gt .LayerCount 20}}
                                    <span class="badge badge-warning">Complex</span>
                                    {{end}}
                                {{end}}
                            </td>
                            <td>
                                <div style="display: inline-block; width: 20px; height: 20px; background: {{.BGColor}}; border: 1px solid #ddd; border-radius: 4px; vertical-align: middle;"></div>
                                <code style="font-size: 0.85em;">{{.BGColor}}</code>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <!-- Layers Tab -->
        <div id="layers" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="layerSearch" placeholder="Search layers..." onkeyup="filterTable('layerSearch', 'layerTable')">
                </div>
                <div>
                    <span class="badge badge-info">{{len .AllLayers}} Total</span>
                </div>
            </div>
            
            <div class="table-container">
                <table id="layerTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Layer <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Composition <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)">Source <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Quality <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Properties <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 5)">Effects <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .AllLayers}}
                        <tr>
                            <td>
                                <strong>#{{.Index}} {{.Name}}</strong>
                                {{if .IsLocked}}<span class="badge badge-warning">🔒</span>{{end}}
                                {{if .IsShyred}}<span class="badge badge-info">👁️</span>{{end}}
                            </td>
                            <td>{{.CompName}}</td>
                            <td>
                                {{if .SourceName}}
                                    {{.SourceName}}
                                {{else}}
                                    <em style="color: #999;">No source</em>
                                {{end}}
                            </td>
                            <td>
                                {{if .Quality}}
                                    <span class="badge badge-{{if eq .Quality "Best"}}success{{else if eq .Quality "Draft"}}warning{{else}}info{{end}}">
                                        {{.Quality}}
                                    </span>
                                {{end}}
                            </td>
                            <td>
                                {{if .Is3D}}<span class="badge badge-info">3D</span>{{end}}
                                {{if .IsSolo}}<span class="badge badge-warning">Solo</span>{{end}}
                                {{if .IsGuide}}<span class="badge badge-info">Guide</span>{{end}}
                                {{if .IsAdjustment}}<span class="badge badge-success">Adjustment</span>{{end}}
                                {{if .HasMotionBlur}}<span class="badge badge-info">Motion Blur</span>{{end}}
                                {{if .IsCollapsed}}<span class="badge badge-info">Collapsed</span>{{end}}
                            </td>
                            <td>
                                {{if .HasEffects}}
                                    <span class="badge badge-success">✨ Effects</span>
                                {{else}}
                                    <span style="color: #999;">None</span>
                                {{end}}
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <!-- Media Tab -->
        <div id="media" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="mediaSearch" placeholder="Search media..." onkeyup="filterTable('mediaSearch', 'mediaTable')">
                </div>
                <div>
                    <span class="badge badge-info">{{len .MediaAssets}} Total</span>
                </div>
            </div>
            
            <div class="table-container">
                <table id="mediaTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Name <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Type <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)">Dimensions <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Duration <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Usage <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .MediaAssets}}
                        <tr>
                            <td><strong>{{.Name}}</strong></td>
                            <td>
                                <span class="badge badge-{{if eq .Type "Video"}}success{{else if eq .Type "Image"}}info{{else if eq .Type "Audio"}}warning{{else}}secondary{{end}}">
                                    {{.Type}}
                                </span>
                            </td>
                            <td>
                                {{if and .Width .Height}}
                                    {{.Width}}×{{.Height}}
                                {{else}}
                                    <em style="color: #999;">N/A</em>
                                {{end}}
                            </td>
                            <td>
                                {{if .Duration}}
                                    {{printf "%.2f" .Duration}}s
                                    {{if .Framerate}}
                                        @ {{printf "%.1f" .Framerate}}fps
                                    {{end}}
                                {{else}}
                                    <em style="color: #999;">Static</em>
                                {{end}}
                            </td>
                            <td>
                                <span class="badge badge-{{if gt .UsageCount 5}}danger{{else if gt .UsageCount 2}}warning{{else}}success{{end}}">
                                    {{.UsageCount}} uses
                                </span>
                                <div class="tooltip">
                                    {{if .UsedIn}}
                                        <span style="font-size: 0.85em; color: #666;">View</span>
                                        <div class="tooltip-content">
                                            {{range .UsedIn}}{{.}}<br>{{end}}
                                        </div>
                                    {{end}}
                                </div>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <!-- Text Tab -->
        <div id="text" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="textSearch" placeholder="Search text layers..." onkeyup="filterTable('textSearch', 'textTable')">
                </div>
                <div>
                    <span class="badge badge-info">{{len .TextLayers}} Text Layers</span>
                </div>
            </div>
            
            {{if .TextLayers}}
            <div class="table-container">
                <table id="textTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Layer <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Composition <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)" style="width: 40%;">Text Content <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Font <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Size <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 5)">Status <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .TextLayers}}
                        <tr>
                            <td><strong>{{.LayerName}}</strong></td>
                            <td>{{.CompName}}</td>
                            <td>
                                <div class="code-block" style="padding: 10px; margin: 0;">
                                    {{.Text}}
                                </div>
                            </td>
                            <td>{{.Font}}</td>
                            <td>{{if .Size}}{{.Size}}pt{{else}}-{{end}}</td>
                            <td>
                                {{if contains .Text "["}}
                                    <span class="badge badge-warning">Keyframed/Expression</span>
                                {{else if eq .Text ""}}
                                    <span class="badge badge-danger">Empty</span>
                                {{else}}
                                    <span class="badge badge-success">Extracted</span>
                                {{end}}
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{else}}
            <div class="alert alert-info">
                <span class="alert-icon">ℹ️</span>
                <span>No text layers found in this project</span>
            </div>
            {{end}}
        </div>
        
        <!-- Effects Tab -->
        <div id="effects" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="effectSearch" placeholder="Search effects..." onkeyup="filterTable('effectSearch', 'effectTable')">
                </div>
                <div>
                    <span class="badge badge-info">{{len .Effects}} Effects</span>
                </div>
            </div>
            
            {{if .Effects}}
            <div class="table-container">
                <table id="effectTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Layer <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Composition <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)">Effect <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Match Name <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Parameters <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 5)">Status <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Effects}}
                        <tr>
                            <td><strong>{{.LayerName}}</strong></td>
                            <td>{{.CompName}}</td>
                            <td><strong>{{.EffectName}}</strong></td>
                            <td><code style="font-size: 0.85em;">{{.MatchName}}</code></td>
                            <td>
                                {{if .Parameters}}
                                    <details>
                                        <summary>{{len .Parameters}} parameters</summary>
                                        <div style="margin-top: 10px;">
                                            {{range .Parameters}}
                                            <div style="margin: 5px 0;">
                                                <strong>{{.Name}}:</strong> {{.Value}}
                                                {{if .Animated}}<span class="badge badge-info">Animated</span>{{end}}
                                            </div>
                                            {{end}}
                                        </div>
                                    </details>
                                {{else}}
                                    <em style="color: #999;">No parameters</em>
                                {{end}}
                            </td>
                            <td>
                                {{if .IsEnabled}}
                                    <span class="badge badge-success">Enabled</span>
                                {{else}}
                                    <span class="badge badge-warning">Disabled</span>
                                {{end}}
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{else}}
            <div class="alert alert-info">
                <span class="alert-icon">ℹ️</span>
                <span>No effects found in this project</span>
            </div>
            {{end}}
        </div>
        
        <!-- Expressions Tab -->
        <div id="expressions" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="exprSearch" placeholder="Search expressions..." onkeyup="filterTable('exprSearch', 'exprTable')">
                </div>
                <div>
                    <span class="badge badge-info">{{len .Expressions}} Expressions</span>
                </div>
            </div>
            
            {{if .Expressions}}
            <div class="table-container">
                <table id="exprTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Layer <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Composition <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)">Property <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Lines <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Language <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 5)">Expression <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Expressions}}
                        <tr>
                            <td><strong>{{.LayerName}}</strong></td>
                            <td>{{.CompName}}</td>
                            <td><code style="font-size: 0.85em;">{{.PropertyPath}}</code></td>
                            <td>{{.LineCount}}</td>
                            <td>
                                <span class="badge badge-{{if eq .Language "JavaScript"}}success{{else}}info{{end}}">
                                    {{.Language}}
                                </span>
                            </td>
                            <td>
                                <details>
                                    <summary>View Expression</summary>
                                    <div class="code-block" style="margin-top: 10px;">
                                        {{.Expression}}
                                    </div>
                                </details>
                                {{if .HasErrors}}
                                    <span class="badge badge-danger">Has Errors</span>
                                {{end}}
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{else}}
            <div class="alert alert-info">
                <span class="alert-icon">ℹ️</span>
                <span>No expressions found in this project</span>
            </div>
            {{end}}
        </div>
        
        <!-- Attributes Tab -->
        <div id="attributes" class="tab-content">
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="attrSearch" placeholder="Search attributes..." onkeyup="filterTable('attrSearch', 'attrTable')">
                </div>
            </div>
            
            <div class="table-container">
                <table id="attrTable">
                    <thead>
                        <tr>
                            <th onclick="sortTable(this, 0)">Layer <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 1)">Composition <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 2)">3D <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 3)">Solo <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 4)">Shy <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 5)">Locked <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 6)">Guide <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 7)">Adjustment <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 8)">Motion Blur <span class="sort-icon">↕</span></th>
                            <th onclick="sortTable(this, 9)">Effects <span class="sort-icon">↕</span></th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .LayerAttributes}}
                        <tr>
                            <td><strong>{{.LayerName}}</strong></td>
                            <td>{{.CompName}}</td>
                            <td>{{if .Is3D}}✅{{else}}❌{{end}}</td>
                            <td>{{if .IsSolo}}✅{{else}}❌{{end}}</td>
                            <td>{{if .IsShy}}✅{{else}}❌{{end}}</td>
                            <td>{{if .IsLocked}}✅{{else}}❌{{end}}</td>
                            <td>{{if .IsGuide}}✅{{else}}❌{{end}}</td>
                            <td>{{if .IsAdjustment}}✅{{else}}❌{{end}}</td>
                            <td>{{if .HasMotionBlur}}✅{{else}}❌{{end}}</td>
                            <td>{{if .HasEffects}}✅{{else}}❌{{end}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <!-- Hierarchy Tab -->
        <div id="hierarchy" class="tab-content">
            <div class="tree-container">
                <h3><span class="card-icon">🌳</span>Project Hierarchy</h3>
                {{template "renderTree" .FolderTree}}
            </div>
        </div>
        
        <!-- Insights Tab -->
        <div id="insights" class="tab-content">
            {{if .ProjectInsights.PerformanceMetrics}}
            <div class="card">
                <h3><span class="card-icon">⚡</span>Performance Metrics</h3>
                <div class="insight-grid">
                    {{range .ProjectInsights.PerformanceMetrics}}
                    <div class="insight-card">
                        <div class="insight-title">
                            <span>📊</span>
                            Performance
                        </div>
                        <div class="insight-content">{{.}}</div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
            
            {{if .ProjectInsights.OptimizationTips}}
            <div class="card">
                <h3><span class="card-icon">🚀</span>Optimization Suggestions</h3>
                <div class="insight-grid">
                    {{range .ProjectInsights.OptimizationTips}}
                    <div class="insight-card">
                        <div class="insight-title">
                            <span>💡</span>
                            Tip
                        </div>
                        <div class="insight-content">{{.}}</div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
            
            {{if .ProjectInsights.CompatibilityNotes}}
            <div class="card">
                <h3><span class="card-icon">🔧</span>Compatibility Notes</h3>
                {{range .ProjectInsights.CompatibilityNotes}}
                <div class="alert alert-warning">
                    <span class="alert-icon">⚠️</span>
                    <span>{{.}}</span>
                </div>
                {{end}}
            </div>
            {{end}}
            
            {{if .ProjectInsights.WorkflowSuggestions}}
            <div class="card">
                <h3><span class="card-icon">🎯</span>Workflow Improvements</h3>
                <div class="insight-grid">
                    {{range .ProjectInsights.WorkflowSuggestions}}
                    <div class="insight-card">
                        <div class="insight-title">
                            <span>🔄</span>
                            Workflow
                        </div>
                        <div class="insight-content">{{.}}</div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>
        
        <!-- Raw Data Tab -->
        <div id="raw" class="tab-content">
            <div class="alert alert-info">
                <span class="alert-icon">ℹ️</span>
                <span>This tab shows raw binary data from properties that couldn't be fully parsed. This is useful for debugging and understanding the AEP file structure.</span>
            </div>
            
            {{if .RawData}}
            <div class="table-controls">
                <div class="search-box">
                    <span class="search-icon">🔍</span>
                    <input type="text" class="search-input" id="rawSearch" placeholder="Search raw data..." onkeyup="filterRawData('rawSearch')">
                </div>
            </div>
            
            {{range .RawData}}
            <div class="card raw-data-card" data-path="{{.PropertyPath}}">
                <h4>
                    <code>{{.PropertyPath}}</code>
                    <span class="badge badge-info">{{.DataType}}</span>
                    <span class="badge badge-secondary">{{.Size}} bytes</span>
                </h4>
                
                {{if .Preview}}
                <div style="margin: 15px 0;">
                    <strong>Preview:</strong>
                    <div class="code-block">{{.Preview}}</div>
                </div>
                {{end}}
                
                {{if .HexDump}}
                <details>
                    <summary>View Hex Dump</summary>
                    <div class="hex-dump">
                        {{range splitLines .HexDump}}
                        <div class="hex-line">
                            <span class="hex-offset">{{index . 0}}</span>
                            <span class="hex-bytes">{{index . 1}}</span>
                            <span class="hex-ascii">{{index . 2}}</span>
                        </div>
                        {{end}}
                    </div>
                </details>
                {{end}}
            </div>
            {{end}}
            {{else}}
            <div class="alert alert-success">
                <span class="alert-icon">✅</span>
                <span>All properties were successfully parsed. No raw data to display.</span>
            </div>
            {{end}}
        </div>
    </div>
    
    <div class="footer">
        <p>Enhanced AEP Parser by mobot2025</p>
        <p><a href="https://github.com/yourusername/mobot2025">github.com/yourusername/mobot2025</a></p>
        <p>Report generated on {{.ParsedAt}}</p>
    </div>
    
    <script>
        // Tab switching
        function showTab(tabName) {
            const tabs = document.querySelectorAll('.tab-content');
            tabs.forEach(tab => tab.classList.remove('active'));
            
            const buttons = document.querySelectorAll('.tab-button');
            buttons.forEach(btn => btn.classList.remove('active'));
            
            document.getElementById(tabName).classList.add('active');
            event.target.classList.add('active');
            
            // Save active tab
            localStorage.setItem('activeTab', tabName);
        }
        
        // Table filtering
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
        
        // Raw data filtering
        function filterRawData(inputId) {
            const input = document.getElementById(inputId);
            const filter = input.value.toUpperCase();
            const cards = document.querySelectorAll('.raw-data-card');
            
            cards.forEach(card => {
                const path = card.getAttribute('data-path');
                if (path.toUpperCase().indexOf(filter) > -1) {
                    card.style.display = '';
                } else {
                    card.style.display = 'none';
                }
            });
        }
        
        // Table sorting
        function sortTable(header, columnIndex) {
            const table = header.closest('table');
            const tbody = table.querySelector('tbody');
            const rows = Array.from(tbody.querySelectorAll('tr'));
            
            // Determine sort direction
            const isAsc = header.classList.contains('sort-asc');
            
            // Remove all sort indicators
            table.querySelectorAll('th').forEach(th => {
                th.classList.remove('sort-asc', 'sort-desc');
                th.querySelector('.sort-icon').textContent = '↕';
            });
            
            // Sort rows
            rows.sort((a, b) => {
                const aText = a.cells[columnIndex].textContent.trim();
                const bText = b.cells[columnIndex].textContent.trim();
                
                // Try to parse as numbers
                const aNum = parseFloat(aText);
                const bNum = parseFloat(bText);
                
                if (!isNaN(aNum) && !isNaN(bNum)) {
                    return isAsc ? bNum - aNum : aNum - bNum;
                }
                
                // Sort as strings
                return isAsc ? 
                    bText.localeCompare(aText) : 
                    aText.localeCompare(bText);
            });
            
            // Update sort indicator
            if (isAsc) {
                header.classList.remove('sort-asc');
                header.classList.add('sort-desc');
                header.querySelector('.sort-icon').textContent = '↓';
            } else {
                header.classList.add('sort-asc');
                header.querySelector('.sort-icon').textContent = '↑';
            }
            
            // Reorder rows
            rows.forEach(row => tbody.appendChild(row));
        }
        
        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            // Ctrl/Cmd + F to focus search
            if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
                e.preventDefault();
                const activeTab = document.querySelector('.tab-content.active');
                const searchInput = activeTab.querySelector('.search-input');
                if (searchInput) {
                    searchInput.focus();
                    searchInput.select();
                }
            }
            
            // Tab navigation with number keys
            if (e.key >= '1' && e.key <= '9' && !e.ctrlKey && !e.metaKey && !e.altKey) {
                const tabIndex = parseInt(e.key) - 1;
                const tabButtons = document.querySelectorAll('.tab-button');
                if (tabButtons[tabIndex]) {
                    tabButtons[tabIndex].click();
                }
            }
        });
        
        // Restore active tab
        const savedTab = localStorage.getItem('activeTab');
        if (savedTab) {
            const tabButton = document.querySelector('[onclick*="' + savedTab + '"]');
            if (tabButton) {
                tabButton.click();
            }
        }
        
        // Progress bar animations
        window.addEventListener('load', () => {
            const progressFills = document.querySelectorAll('.progress-fill');
            progressFills.forEach(fill => {
                const width = fill.style.width;
                fill.style.width = '0';
                setTimeout(() => {
                    fill.style.width = width;
                }, 100);
            });
        });
        
        // Copy to clipboard functionality
        function copyToClipboard(text) {
            const temp = document.createElement('textarea');
            temp.value = text;
            document.body.appendChild(temp);
            temp.select();
            document.execCommand('copy');
            document.body.removeChild(temp);
            
            // Show feedback
            const tooltip = document.createElement('div');
            tooltip.className = 'copy-tooltip';
            tooltip.textContent = 'Copied!';
            tooltip.style.cssText = 'position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%); background: #48bb78; color: white; padding: 10px 20px; border-radius: 6px; z-index: 10000;';
            document.body.appendChild(tooltip);
            
            setTimeout(() => {
                document.body.removeChild(tooltip);
            }, 1000);
        }
        
        // Add copy functionality to code blocks
        document.querySelectorAll('.code-block, code').forEach(block => {
            block.style.cursor = 'pointer';
            block.title = 'Click to copy';
            block.addEventListener('click', () => {
                copyToClipboard(block.textContent);
            });
        });
    </script>
</body>
</html>

{{define "renderTree"}}
    <div class="tree-node {{if eq .Name "root"}}tree-node-root{{end}}">
        <span class="tree-icon">
            {{if eq .Type "folder"}}📁{{else if eq .Type "comp"}}🎬{{else}}🎞️{{end}}
        </span>
        {{.Name}}
        {{if .ItemCount}}
            <span class="tree-count">({{.ItemCount}} items)</span>
        {{end}}
    </div>
    {{range .Children}}
        {{template "renderTree" .}}
    {{end}}
{{end}}`

// Helper functions for template
var funcMap = template.FuncMap{
	"contains": strings.Contains,
	"calcPercent": func(part, total int) float64 {
		if total == 0 {
			return 0
		}
		return float64(part) / float64(total) * 100
	},
	"splitLines": func(s string) [][]string {
		lines := strings.Split(s, "\n")
		result := make([][]string, 0)
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				offset := parts[0]
				hex := strings.Join(parts[1:len(parts)-1], " ")
				ascii := parts[len(parts)-1]
				result = append(result, []string{offset, hex, ascii})
			}
		}
		return result
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_enhanced_ux_report.go <aep-file>")
		fmt.Println("Example: go run generate_enhanced_ux_report.go sample-aep/Ai\\ Text\\ Intro.aep")
		os.Exit(1)
	}
	
	aepFile := os.Args[1]
	fmt.Printf("📄 Parsing: %s\n", aepFile)
	
	// Get file info
	fileInfo, err := os.Stat(aepFile)
	if err != nil {
		fmt.Printf("❌ Error getting file info: %v\n", err)
		os.Exit(1)
	}
	
	// Parse the AEP file
	project, err := aep.Open(aepFile)
	if err != nil {
		fmt.Printf("❌ Error parsing AEP file: %v\n", err)
		os.Exit(1)
	}
	
	// Prepare enhanced report data
	report := prepareEnhancedReport(aepFile, fileInfo, project)
	
	// Generate HTML filename
	baseFilename := strings.TrimSuffix(filepath.Base(aepFile), filepath.Ext(aepFile))
	htmlFilename := fmt.Sprintf("%s-enhanced-ux-report-%s.html", 
		baseFilename, 
		time.Now().Format("2006-01-02-150405"))
	
	// Create HTML file
	file, err := os.Create(htmlFilename)
	if err != nil {
		fmt.Printf("❌ Error creating HTML file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	
	// Parse and execute template with functions
	tmpl, err := template.New("report").Funcs(funcMap).Parse(enhancedHTMLTemplate)
	if err != nil {
		fmt.Printf("❌ Error parsing template: %v\n", err)
		os.Exit(1)
	}
	
	err = tmpl.Execute(file, report)
	if err != nil {
		fmt.Printf("❌ Error generating HTML: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ Enhanced UX report generated: %s\n", htmlFilename)
	fmt.Printf("📂 Open in browser: file://%s/%s\n", 
		mustGetWorkingDir(), htmlFilename)
}

func prepareEnhancedReport(filename string, fileInfo os.FileInfo, project *aep.Project) EnhancedReportData {
	report := EnhancedReportData{
		Filename:         filepath.Base(filename),
		FilePath:         filename,
		ParsedAt:         time.Now().Format("January 2, 2006 at 3:04 PM"),
		TotalItems:       len(project.Items),
		ExpressionEngine: project.ExpressionEngine,
		FileSize:         formatFileSize(fileInfo.Size()),
	}
	
	// Convert bit depth
	switch project.Depth {
	case aep.BPC8:
		report.BitDepth = "8-bit"
		report.Overview.Statistics.ColorDepth = "16.7 million colors (8 bits/channel)"
	case aep.BPC16:
		report.BitDepth = "16-bit"
		report.Overview.Statistics.ColorDepth = "281 trillion colors (16 bits/channel)"
	case aep.BPC32:
		report.BitDepth = "32-bit (float)"
		report.Overview.Statistics.ColorDepth = "Floating point precision (32 bits/channel)"
	default:
		report.BitDepth = fmt.Sprintf("Unknown (%v)", project.Depth)
		report.Overview.Statistics.ColorDepth = "Unknown color depth"
	}
	
	// Initialize maps and counters
	mediaUsage := make(map[uint32][]string)
	allComps := make([]CompositionDetail, 0)
	allLayers := make([]LayerDetail, 0)
	allEffects := make([]EffectDetail, 0)
	allExpressions := make([]ExpressionDetail, 0)
	allTextLayers := make([]TextLayer, 0)
	allAttributes := make([]LayerAttribute, 0)
	effectCounts := make(map[string]int)
	assetTypes := make(map[string]int)
	layerTypes := make(map[string]int)
	
	var totalFrames int
	var totalDuration float64
	var totalFramerate float64
	var framerateCount int
	
	// Process all items
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeFolder:
			report.Overview.FolderCount++
			assetTypes["Folders"]++
			
		case aep.ItemTypeComposition:
			report.Overview.CompositionCount++
			report.Overview.TotalLayers += len(item.CompositionLayers)
			assetTypes["Compositions"]++
			
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
					int(item.BackgroundColor[0]*255),
					int(item.BackgroundColor[1]*255),
					int(item.BackgroundColor[2]*255)),
			}
			
			// Calculate totals
			totalFrames += comp.Frames
			totalDuration += comp.Duration
			if comp.Framerate > 0 {
				totalFramerate += comp.Framerate
				framerateCount++
			}
			
			// Determine resolution
			if comp.Width >= 3840 {
				comp.Resolution = "4K"
			} else if comp.Width >= 1920 {
				comp.Resolution = "HD"
			} else if comp.Width >= 1280 {
				comp.Resolution = "HD Ready"
			} else {
				comp.Resolution = "SD"
			}
			
			// Process layers in this comp
			for _, layer := range item.CompositionLayers {
				// Count layer types
				if layer.AdjustmentLayerEnabled {
					layerTypes["Adjustment"]++
				} else if layer.Text != nil {
					layerTypes["Text"]++
				} else if layer.ThreeDEnabled {
					layerTypes["3D"]++
				} else {
					layerTypes["2D"]++
				}
				
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
					if sourceItem.ItemType == aep.ItemTypeFootage {
						mediaUsage[sourceItem.ID] = append(mediaUsage[sourceItem.ID], item.Name)
					}
				}
				
				// Determine quality
				switch layer.Quality {
				case 0:
					layerDetail.Quality = "Best"
				case 1:
					layerDetail.Quality = "Draft"
				case 2:
					layerDetail.Quality = "Wireframe"
				}
				
				// Check for text layer
				if layer.Text != nil {
					textLayer := TextLayer{
						LayerName: layer.Name,
						CompName:  item.Name,
						Text:      "[Text content not extracted]",
						Font:      "Unknown",
						Size:      0,
					}
					
					// Try to extract text content
					if layer.Text != nil && layer.Text.TextDocument != nil {
						doc := layer.Text.TextDocument
						if doc.Text != "" {
							textLayer.Text = doc.Text
						}
						if doc.FontName != "" {
							textLayer.Font = doc.FontName
						}
						textLayer.Size = doc.FontSize
					}
					
					// Check for keyframe/expression text
					if strings.Contains(textLayer.Text, "not extracted") && layer.Name != "" {
						// Use layer name as hint
						textLayer.Text = fmt.Sprintf("[Layer: %s - Text in keyframes/expressions]", layer.Name)
					}
					
					allTextLayers = append(allTextLayers, textLayer)
				}
				
				// Add layer attributes
				attr := LayerAttribute{
					LayerName:     layer.Name,
					CompName:      item.Name,
					Is3D:          layer.ThreeDEnabled,
					IsSolo:        layer.SoloEnabled,
					IsShy:         layer.ShyEnabled,
					IsLocked:      layer.LockEnabled,
					IsGuide:       layer.GuideEnabled,
					IsAdjustment:  layer.AdjustmentLayerEnabled,
					HasMotionBlur: layer.MotionBlurEnabled,
					HasEffects:    layer.EffectsEnabled,
				}
				allAttributes = append(allAttributes, attr)
				
				// Process effects (placeholder - would need actual effect parsing)
				if layer.EffectsEnabled {
					report.Overview.TotalEffects++
				}
				
				allLayers = append(allLayers, layerDetail)
			}
			
			// Add to main comps if it's likely a main comp
			if comp.LayerCount > 5 || strings.Contains(strings.ToLower(comp.Name), "main") || 
			   strings.Contains(strings.ToLower(comp.Name), "final") {
				report.Overview.MainComps = append(report.Overview.MainComps, comp)
			}
			
			allComps = append(allComps, comp)
			
		case aep.ItemTypeFootage:
			report.Overview.FootageCount++
			
			asset := MediaAsset{
				ID:        item.ID,
				Name:      item.Name,
				Width:     item.FootageDimensions[0],
				Height:    item.FootageDimensions[1],
				Framerate: item.FootageFramerate,
				Duration:  item.FootageSeconds,
			}
			
			if item.FootageFramerate > 0 && item.FootageSeconds > 0 {
				asset.Frames = int(item.FootageSeconds * item.FootageFramerate)
			}
			
			// Determine asset type
			switch item.FootageType {
			case 1:
				asset.Type = "Image"
				assetTypes["Images"]++
			case 2:
				asset.Type = "Video"
				assetTypes["Videos"]++
			case 3:
				asset.Type = "Audio"
				assetTypes["Audio"]++
			default:
				asset.Type = "Other"
				assetTypes["Other"]++
			}
			
			// Add usage info
			if usages, exists := mediaUsage[item.ID]; exists {
				asset.UsedIn = usages
				asset.UsageCount = len(usages)
			}
			
			report.MediaAssets = append(report.MediaAssets, asset)
		}
	}
	
	// Sort compositions by layer count (descending)
	sort.Slice(allComps, func(i, j int) bool {
		return allComps[i].LayerCount > allComps[j].LayerCount
	})
	
	// Sort main comps
	sort.Slice(report.Overview.MainComps, func(i, j int) bool {
		return report.Overview.MainComps[i].LayerCount > report.Overview.MainComps[j].LayerCount
	})
	
	// Limit main comps to top 5
	if len(report.Overview.MainComps) > 5 {
		report.Overview.MainComps = report.Overview.MainComps[:5]
	}
	
	// Calculate statistics
	report.Overview.Statistics.TotalFrames = totalFrames
	report.Overview.Statistics.TotalDuration = totalDuration
	if framerateCount > 0 {
		report.Overview.Statistics.AverageFramerate = totalFramerate / float64(framerateCount)
	}
	report.Overview.Statistics.UsedEffects = effectCounts
	report.Overview.Statistics.AssetTypes = assetTypes
	report.Overview.Statistics.LayerTypes = layerTypes
	
	// Calculate complexity score (0-100)
	complexity := 0
	complexity += min(report.Overview.CompositionCount*2, 30)
	complexity += min(report.Overview.TotalLayers/10, 30)
	complexity += min(report.Overview.TotalEffects/5, 20)
	complexity += min(len(allExpressions)*2, 20)
	report.Overview.Statistics.ComplexityScore = complexity
	
	// Estimate memory usage
	estimatedMemory := int64(0)
	estimatedMemory += int64(report.Overview.TotalLayers) * 1024 * 100 // 100KB per layer estimate
	estimatedMemory += int64(report.Overview.FootageCount) * 1024 * 500 // 500KB per footage estimate
	report.Overview.Statistics.MemoryUsage = formatFileSize(estimatedMemory)
	
	// Generate insights
	report.Overview.Insights = generateProjectInsights(report)
	
	// Generate warnings
	report.Overview.Warnings = generateProjectWarnings(report)
	
	// Generate project insights
	report.ProjectInsights = generateDetailedInsights(report)
	
	// Build folder tree
	report.FolderTree = buildFolderTree(project)
	
	// Assign data
	report.Compositions = allComps
	report.AllLayers = allLayers
	report.TextLayers = allTextLayers
	report.LayerAttributes = allAttributes
	report.Effects = allEffects
	report.Expressions = allExpressions
	
	return report
}

func generateProjectInsights(report EnhancedReportData) []string {
	insights := []string{}
	
	// Composition insights
	if report.Overview.CompositionCount > 20 {
		insights = append(insights, fmt.Sprintf("Large project with %d compositions. Consider organizing into folders for better management.", report.Overview.CompositionCount))
	}
	
	// Layer insights
	avgLayersPerComp := 0
	if report.Overview.CompositionCount > 0 {
		avgLayersPerComp = report.Overview.TotalLayers / report.Overview.CompositionCount
	}
	if avgLayersPerComp > 50 {
		insights = append(insights, fmt.Sprintf("High layer density (%d layers/comp average). This may impact performance.", avgLayersPerComp))
	}
	
	// Performance insights
	if report.Overview.Statistics.ComplexityScore > 70 {
		insights = append(insights, "High complexity project. Consider pre-rendering heavy compositions for better playback.")
	}
	
	// Text layer insights
	if len(report.TextLayers) > 0 {
		extractedCount := 0
		for _, text := range report.TextLayers {
			if !strings.Contains(text.Text, "[") {
				extractedCount++
			}
		}
		if extractedCount == 0 {
			insights = append(insights, "All text appears to be animated or expression-based. Text content is stored in keyframes.")
		}
	}
	
	return insights
}

func generateProjectWarnings(report EnhancedReportData) []string {
	warnings := []string{}
	
	// Check for missing footage
	if report.Overview.FootageCount == 0 && report.Overview.TotalLayers > 0 {
		warnings = append(warnings, "No footage items found but layers exist. Check for missing media.")
	}
	
	// Check for extreme framerates
	if report.Overview.Statistics.AverageFramerate > 60 {
		warnings = append(warnings, fmt.Sprintf("High average framerate (%.1f fps) detected. Ensure playback compatibility.", report.Overview.Statistics.AverageFramerate))
	}
	
	// Check for memory usage
	memStr := report.Overview.Statistics.MemoryUsage
	if strings.Contains(memStr, "GB") {
		warnings = append(warnings, fmt.Sprintf("High estimated memory usage (%s). This project may require significant RAM.", memStr))
	}
	
	return warnings
}

func generateDetailedInsights(report EnhancedReportData) ProjectInsights {
	insights := ProjectInsights{}
	
	// Performance metrics
	insights.PerformanceMetrics = []string{
		fmt.Sprintf("Project contains %d total frames across all compositions", report.Overview.Statistics.TotalFrames),
		fmt.Sprintf("Average composition framerate: %.1f fps", report.Overview.Statistics.AverageFramerate),
		fmt.Sprintf("Estimated memory footprint: %s", report.Overview.Statistics.MemoryUsage),
		fmt.Sprintf("Complexity score: %d/100", report.Overview.Statistics.ComplexityScore),
	}
	
	// Optimization tips
	if report.Overview.Statistics.ComplexityScore > 50 {
		insights.OptimizationTips = append(insights.OptimizationTips, 
			"Consider pre-composing complex layer groups to improve timeline performance")
	}
	if report.Overview.TotalLayers > 100 {
		insights.OptimizationTips = append(insights.OptimizationTips, 
			"Use guide layers for reference footage to reduce render overhead")
	}
	if report.Overview.TotalEffects > 50 {
		insights.OptimizationTips = append(insights.OptimizationTips, 
			"Apply effects to adjustment layers instead of individual layers when possible")
	}
	
	// Compatibility notes
	if report.BitDepth == "32-bit (float)" {
		insights.CompatibilityNotes = append(insights.CompatibilityNotes,
			"32-bit float color depth provides maximum quality but requires more processing power")
	}
	if report.ExpressionEngine == "javascript-1.0" {
		insights.CompatibilityNotes = append(insights.CompatibilityNotes,
			"JavaScript expression engine is enabled. Ensure expressions are compatible with target After Effects version.")
	}
	
	// Workflow suggestions
	if len(report.TextLayers) > 10 {
		insights.WorkflowSuggestions = append(insights.WorkflowSuggestions,
			"Multiple text layers detected. Consider using Essential Graphics for easier text management.")
	}
	if report.Overview.FolderCount < 3 && report.Overview.CompositionCount > 10 {
		insights.WorkflowSuggestions = append(insights.WorkflowSuggestions,
			"Organize compositions into folders for better project structure")
	}
	
	return insights
}

func buildFolderTree(project *aep.Project) FolderNode {
	root := FolderNode{
		Name:     "Project Root",
		Type:     "folder",
		Children: []FolderNode{},
	}
	
	// This is a simplified version - would need proper folder hierarchy parsing
	folderCount := 0
	compCount := 0
	footageCount := 0
	
	for _, item := range project.Items {
		switch item.ItemType {
		case aep.ItemTypeFolder:
			folderCount++
		case aep.ItemTypeComposition:
			compCount++
		case aep.ItemTypeFootage:
			footageCount++
		}
	}
	
	if folderCount > 0 {
		root.Children = append(root.Children, FolderNode{
			Name:      fmt.Sprintf("Folders (%d)", folderCount),
			Type:      "folder",
			ItemCount: folderCount,
		})
	}
	
	if compCount > 0 {
		root.Children = append(root.Children, FolderNode{
			Name:      fmt.Sprintf("Compositions (%d)", compCount),
			Type:      "comp",
			ItemCount: compCount,
		})
	}
	
	if footageCount > 0 {
		root.Children = append(root.Children, FolderNode{
			Name:      fmt.Sprintf("Footage (%d)", footageCount),
			Type:      "footage",
			ItemCount: footageCount,
		})
	}
	
	root.ItemCount = len(project.Items)
	
	return root
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mustGetWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}