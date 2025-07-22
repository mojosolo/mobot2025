// Package catalog provides API service for AEP processing pipeline
package catalog

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "sync"
    "time"
)

// APIService provides HTTP endpoints for the AEP processing pipeline
type APIService struct {
    parser     *Parser
    database   *Database
    cache      *ProjectCache
    jobQueue   *JobQueue
    port       int
}

// ProjectCache provides thread-safe caching for parsed projects
type ProjectCache struct {
    mu       sync.RWMutex
    projects map[string]*ProjectMetadata
    ttl      time.Duration
}

// Job represents a processing job
type Job struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Status      string                 `json:"status"`
    Progress    float64                `json:"progress"`
    Input       map[string]interface{} `json:"input"`
    Output      map[string]interface{} `json:"output,omitempty"`
    Error       string                 `json:"error,omitempty"`
    CreatedAt   time.Time              `json:"created_at"`
    CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// JobQueue manages background processing jobs
type JobQueue struct {
    mu   sync.RWMutex
    jobs map[string]*Job
}

// NewAPIService creates a new API service
func NewAPIService(port int, dbPath string) (*APIService, error) {
    // Initialize database
    database, err := NewDatabase(dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }

    return &APIService{
        parser: NewParser(),
        database: database,
        cache: &ProjectCache{
            projects: make(map[string]*ProjectMetadata),
            ttl:      15 * time.Minute,
        },
        jobQueue: &JobQueue{
            jobs: make(map[string]*Job),
        },
        port: port,
    }, nil
}

// Start begins the HTTP server
func (s *APIService) Start() error {
    mux := http.NewServeMux()
    
    // API endpoints
    mux.HandleFunc("/api/v1/parse", s.handleParse)
    mux.HandleFunc("/api/v1/catalog", s.handleCatalog)
    mux.HandleFunc("/api/v1/analyze", s.handleAnalyze)
    mux.HandleFunc("/api/v1/opportunities", s.handleOpportunities)
    mux.HandleFunc("/api/v1/nexrender", s.handleNexRender)
    mux.HandleFunc("/api/v1/jobs/", s.handleJobs)
    mux.HandleFunc("/api/v1/health", s.handleHealth)
    mux.HandleFunc("/api/v1/search", s.handleSearch)
    mux.HandleFunc("/api/v1/filter", s.handleFilter)
    
    // Static file server for reports
    mux.Handle("/reports/", http.StripPrefix("/reports/", http.FileServer(http.Dir("./reports"))))
    
    addr := fmt.Sprintf(":%d", s.port)
    fmt.Printf("AEP Catalog API starting on %s\n", addr)
    
    return http.ListenAndServe(addr, mux)
}

// handleParse processes a single AEP file
func (s *APIService) handleParse(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        FilePath string `json:"file_path"`
        Options  struct {
            ExtractText  bool `json:"extract_text"`
            ExtractMedia bool `json:"extract_media"`
            DeepAnalysis bool `json:"deep_analysis"`
        } `json:"options"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Check cache first
    if cached := s.cache.Get(req.FilePath); cached != nil {
        s.writeJSON(w, cached)
        return
    }
    
    // Configure parser
    parser := NewParser()
    parser.ExtractText = req.Options.ExtractText
    parser.ExtractMedia = req.Options.ExtractMedia
    parser.DeepAnalysis = req.Options.DeepAnalysis
    
    // Parse the project
    metadata, err := parser.ParseProject(req.FilePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Store in database
    if err := s.database.StoreProject(metadata); err != nil {
        // Log error but continue (non-critical)
        fmt.Printf("Failed to store project in database: %v\n", err)
    }
    
    // Cache the result
    s.cache.Set(req.FilePath, metadata)
    
    s.writeJSON(w, metadata)
}

// handleCatalog processes a directory of AEP files
func (s *APIService) handleCatalog(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        Directory string `json:"directory"`
        Pattern   string `json:"pattern"`
        Async     bool   `json:"async"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if req.Pattern == "" {
        req.Pattern = "*.aep"
    }
    
    // Create a job for async processing
    if req.Async {
        job := s.createJob("catalog", map[string]interface{}{
            "directory": req.Directory,
            "pattern":   req.Pattern,
        })
        
        go s.processCatalogJob(job)
        
        s.writeJSON(w, map[string]interface{}{
            "job_id": job.ID,
            "status": "processing",
        })
        return
    }
    
    // Synchronous processing
    results, err := s.catalogDirectory(req.Directory, req.Pattern)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    s.writeJSON(w, map[string]interface{}{
        "total": len(results),
        "templates": results,
    })
}

// handleAnalyze performs deep analysis on a project
func (s *APIService) handleAnalyze(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        FilePath string `json:"file_path"`
        Analysis struct {
            ModularSystem bool `json:"modular_system"`
            AssetMapping  bool `json:"asset_mapping"`
            Optimization  bool `json:"optimization"`
        } `json:"analysis"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Parse the project first
    metadata, err := s.parser.ParseProject(req.FilePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Perform additional analysis
    analysis := map[string]interface{}{
        "metadata": metadata,
    }
    
    if req.Analysis.ModularSystem && metadata.Capabilities.IsModular {
        analysis["modular_analysis"] = s.analyzeModularSystem(metadata)
    }
    
    if req.Analysis.AssetMapping {
        analysis["asset_mapping"] = s.analyzeAssetMapping(metadata)
    }
    
    if req.Analysis.Optimization {
        analysis["optimization_suggestions"] = s.analyzeOptimization(metadata)
    }
    
    s.writeJSON(w, analysis)
}

// handleOpportunities identifies automation opportunities
func (s *APIService) handleOpportunities(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        FilePaths []string `json:"file_paths"`
        Criteria  struct {
            MinTextLayers int     `json:"min_text_layers"`
            MinImpact     string  `json:"min_impact"`
            MaxDifficulty string  `json:"max_difficulty"`
            Types         []string `json:"types"`
        } `json:"criteria"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    allOpportunities := []map[string]interface{}{}
    
    for _, filePath := range req.FilePaths {
        metadata, err := s.parser.ParseProject(filePath)
        if err != nil {
            continue
        }
        
        for _, opp := range metadata.Opportunities {
            // Apply filters
            if req.Criteria.MinImpact != "" && !s.meetsImpactCriteria(opp.Impact, req.Criteria.MinImpact) {
                continue
            }
            
            if req.Criteria.MaxDifficulty != "" && !s.meetsDifficultyCriteria(opp.Difficulty, req.Criteria.MaxDifficulty) {
                continue
            }
            
            if len(req.Criteria.Types) > 0 && !s.containsType(opp.Type, req.Criteria.Types) {
                continue
            }
            
            allOpportunities = append(allOpportunities, map[string]interface{}{
                "file": filePath,
                "opportunity": opp,
                "metadata": map[string]interface{}{
                    "file_name": metadata.FileName,
                    "categories": metadata.Categories,
                    "tags": metadata.Tags,
                },
            })
        }
    }
    
    s.writeJSON(w, map[string]interface{}{
        "total": len(allOpportunities),
        "opportunities": allOpportunities,
    })
}

// handleNexRender generates nexrender configurations
func (s *APIService) handleNexRender(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        FilePath     string                 `json:"file_path"`
        Replacements map[string]interface{} `json:"replacements"`
        Output       struct {
            Module      string `json:"module"`
            Codec       string `json:"codec"`
            Preset      string `json:"preset"`
            Destination string `json:"destination"`
        } `json:"output"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Parse the project
    metadata, err := s.parser.ParseProject(req.FilePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Generate base nexrender config
    config := metadata.ToNexRenderConfig()
    
    // Apply custom replacements
    if assets, ok := config["assets"].([]map[string]interface{}); ok {
        for _, asset := range assets {
            if layerName, ok := asset["layerName"].(string); ok {
                if replacement, exists := req.Replacements[layerName]; exists {
                    asset["value"] = replacement
                }
            }
        }
    }
    
    // Configure output
    if req.Output.Module != "" {
        config["actions"] = map[string]interface{}{
            "postrender": []map[string]interface{}{
                {
                    "module": req.Output.Module,
                    "codec": req.Output.Codec,
                    "preset": req.Output.Preset,
                    "output": req.Output.Destination,
                },
            },
        }
    }
    
    s.writeJSON(w, config)
}

// handleJobs manages background job status
func (s *APIService) handleJobs(w http.ResponseWriter, r *http.Request) {
    jobID := filepath.Base(r.URL.Path)
    
    if jobID == "" || jobID == "jobs" {
        // List all jobs
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        
        jobs := s.jobQueue.GetAll()
        s.writeJSON(w, jobs)
        return
    }
    
    // Get specific job
    job := s.jobQueue.Get(jobID)
    if job == nil {
        http.Error(w, "Job not found", http.StatusNotFound)
        return
    }
    
    s.writeJSON(w, job)
}

// handleHealth provides health check endpoint
func (s *APIService) handleHealth(w http.ResponseWriter, r *http.Request) {
    s.writeJSON(w, map[string]interface{}{
        "status": "healthy",
        "service": "aep-catalog-api",
        "version": "1.0.0",
        "uptime": time.Since(startTime).String(),
    })
}

// Helper methods

func (s *APIService) writeJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func (s *APIService) catalogDirectory(directory, pattern string) ([]*ProjectMetadata, error) {
    // Implementation would scan directory and parse all matching files
    // For now, returning empty slice
    return []*ProjectMetadata{}, nil
}

func (s *APIService) createJob(jobType string, input map[string]interface{}) *Job {
    job := &Job{
        ID:        fmt.Sprintf("%s_%d", jobType, time.Now().Unix()),
        Type:      jobType,
        Status:    "pending",
        Progress:  0,
        Input:     input,
        CreatedAt: time.Now(),
    }
    
    s.jobQueue.Set(job)
    return job
}

func (s *APIService) processCatalogJob(job *Job) {
    // Implementation would process the catalog job asynchronously
    job.Status = "processing"
    s.jobQueue.Set(job)
    
    // Simulate processing
    time.Sleep(2 * time.Second)
    
    job.Status = "completed"
    job.Progress = 100
    now := time.Now()
    job.CompletedAt = &now
    s.jobQueue.Set(job)
}

func (s *APIService) analyzeModularSystem(metadata *ProjectMetadata) map[string]interface{} {
    // Analyze modular composition relationships
    return map[string]interface{}{
        "total_modules": len(metadata.Compositions),
        "can_mix_match": true,
        "suggested_combinations": []string{
            "Intro + Main + Outro",
            "Logo + Text Animation",
        },
    }
}

func (s *APIService) analyzeAssetMapping(metadata *ProjectMetadata) map[string]interface{} {
    // Map assets to usage locations
    mapping := map[string][]string{}
    
    for _, asset := range metadata.MediaAssets {
        mapping[asset.Name] = []string{
            "Main Composition",
            "Background Layer",
        }
    }
    
    return map[string]interface{}{
        "asset_usage": mapping,
        "replaceable_count": len(metadata.MediaAssets),
    }
}

func (s *APIService) analyzeOptimization(metadata *ProjectMetadata) []string {
    suggestions := []string{}
    
    if len(metadata.Effects) > 10 {
        suggestions = append(suggestions, "Consider pre-rendering heavy effects")
    }
    
    if len(metadata.TextLayers) > 20 {
        suggestions = append(suggestions, "Use text templating system for bulk updates")
    }
    
    return suggestions
}

func (s *APIService) meetsImpactCriteria(impact, minImpact string) bool {
    impactLevels := map[string]int{"low": 1, "medium": 2, "high": 3}
    return impactLevels[impact] >= impactLevels[minImpact]
}

func (s *APIService) meetsDifficultyCriteria(difficulty, maxDifficulty string) bool {
    difficultyLevels := map[string]int{"easy": 1, "medium": 2, "hard": 3}
    return difficultyLevels[difficulty] <= difficultyLevels[maxDifficulty]
}

func (s *APIService) containsType(oppType string, types []string) bool {
    for _, t := range types {
        if t == oppType {
            return true
        }
    }
    return false
}

// Cache methods

func (c *ProjectCache) Get(key string) *ProjectMetadata {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.projects[key]
}

func (c *ProjectCache) Set(key string, value *ProjectMetadata) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.projects[key] = value
}

// JobQueue methods

func (j *JobQueue) Get(id string) *Job {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return j.jobs[id]
}

func (j *JobQueue) Set(job *Job) {
    j.mu.Lock()
    defer j.mu.Unlock()
    j.jobs[job.ID] = job
}

func (j *JobQueue) GetAll() []*Job {
    j.mu.RLock()
    defer j.mu.RUnlock()
    
    jobs := make([]*Job, 0, len(j.jobs))
    for _, job := range j.jobs {
        jobs = append(jobs, job)
    }
    return jobs
}

// handleSearch handles full-text search requests
func (s *APIService) handleSearch(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
        return
    }
    
    limit := 50 // Default limit
    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        if parsedLimit := parseInt(limitStr, 50); parsedLimit > 0 && parsedLimit <= 100 {
            limit = parsedLimit
        }
    }
    
    results, err := s.database.SearchProjects(query, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    s.writeJSON(w, map[string]interface{}{
        "query": query,
        "total": len(results),
        "results": results,
    })
}

// handleFilter handles advanced filtering requests
func (s *APIService) handleFilter(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req struct {
        Categories    []string `json:"categories"`
        Tags          []string `json:"tags"`
        MinComplexity float64  `json:"min_complexity"`
        MaxComplexity float64  `json:"max_complexity"`
        Limit         int      `json:"limit"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if req.Limit == 0 {
        req.Limit = 50
    }
    
    filter := ProjectFilter{
        Categories:    req.Categories,
        Tags:          req.Tags,
        MinComplexity: req.MinComplexity,
        MaxComplexity: req.MaxComplexity,
        Limit:         req.Limit,
    }
    
    results, err := s.database.FilterProjects(filter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    s.writeJSON(w, map[string]interface{}{
        "filter": req,
        "total": len(results),
        "results": results,
    })
}

// Helper function to parse integers safely
func parseInt(s string, defaultVal int) int {
    var i int
    if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
        return i
    }
    return defaultVal
}

var startTime = time.Now()