// Package catalog provides nexrender integration for automated rendering
package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NexRenderIntegration handles nexrender job management
type NexRenderIntegration struct {
	apiURL      string
	apiKey      string
	workDir     string
	httpClient  *http.Client
}

// NexRenderJob represents a nexrender job configuration
type NexRenderJob struct {
	UID         string                   `json:"uid,omitempty"`
	Type        string                   `json:"type"`
	State       string                   `json:"state,omitempty"`
	Template    NexRenderTemplate        `json:"template"`
	Assets      []NexRenderAsset         `json:"assets"`
	Actions     NexRenderActions         `json:"actions,omitempty"`
	Settings    NexRenderSettings        `json:"settings,omitempty"`
	OnComplete  string                   `json:"oncomplete,omitempty"`
	OnError     string                   `json:"onerror,omitempty"`
}

// NexRenderTemplate defines the AE template
type NexRenderTemplate struct {
	Src         string `json:"src"`
	Composition string `json:"composition,omitempty"`
	FrameStart  int    `json:"frameStart,omitempty"`
	FrameEnd    int    `json:"frameEnd,omitempty"`
	Continual   bool   `json:"continual,omitempty"`
}

// NexRenderAsset represents a replaceable asset
type NexRenderAsset struct {
	Src         string      `json:"src,omitempty"`
	Type        string      `json:"type"`
	LayerName   string      `json:"layerName,omitempty"`
	LayerIndex  int         `json:"layerIndex,omitempty"`
	Property    string      `json:"property,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Expression  string      `json:"expression,omitempty"`
	Composition string      `json:"composition,omitempty"`
}

// NexRenderActions defines post-render actions
type NexRenderActions struct {
	Prerender   []NexRenderAction `json:"prerender,omitempty"`
	Postrender  []NexRenderAction `json:"postrender,omitempty"`
}

// NexRenderAction represents a single action
type NexRenderAction struct {
	Module      string                 `json:"module"`
	Output      string                 `json:"output,omitempty"`
	Input       string                 `json:"input,omitempty"`
	Preset      string                 `json:"preset,omitempty"`
	Params      map[string]interface{} `json:"params,omitempty"`
}

// NexRenderSettings contains render settings
type NexRenderSettings struct {
	OutputModule string `json:"outputModule,omitempty"`
	OutputExt    string `json:"outputExt,omitempty"`
	SettingsTemplate string `json:"settingsTemplate,omitempty"`
}

// JobStatus represents the status of a render job
type JobStatus struct {
	UID        string    `json:"uid"`
	State      string    `json:"state"`
	Progress   float64   `json:"progress"`
	RenderTime float64   `json:"renderTime"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
	Error      string    `json:"error,omitempty"`
	Output     string    `json:"output,omitempty"`
}

// NewNexRenderIntegration creates a new nexrender integration
func NewNexRenderIntegration(apiURL, apiKey, workDir string) *NexRenderIntegration {
	return &NexRenderIntegration{
		apiURL:  apiURL,
		apiKey:  apiKey,
		workDir: workDir,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateJobFromMetadata creates a nexrender job from project metadata
func (n *NexRenderIntegration) CreateJobFromMetadata(metadata *ProjectMetadata, replacements map[string]interface{}) (*NexRenderJob, error) {
	if len(metadata.Compositions) == 0 {
		return nil, fmt.Errorf("no compositions found in project")
	}

	// Create base job
	job := &NexRenderJob{
		Type: "default",
		Template: NexRenderTemplate{
			Src:         metadata.FilePath,
			Composition: metadata.Compositions[0].Name, // Use first composition
		},
		Assets:   []NexRenderAsset{},
		Settings: NexRenderSettings{
			OutputModule: "h264",
			OutputExt:    "mp4",
		},
	}

	// Add text replacements
	for _, textLayer := range metadata.TextLayers {
		if replacement, ok := replacements[textLayer.LayerName]; ok {
			asset := NexRenderAsset{
				Type:      "data",
				LayerName: textLayer.LayerName,
				Property:  "Source Text",
				Value:     replacement,
			}
			
			// If specific composition
			if textLayer.CompID != "" {
				for _, comp := range metadata.Compositions {
					if comp.ID == textLayer.CompID {
						asset.Composition = comp.Name
						break
					}
				}
			}
			
			job.Assets = append(job.Assets, asset)
		}
	}

	// Add media replacements
	for _, mediaAsset := range metadata.MediaAssets {
		key := n.generateReplacementKey(mediaAsset.Name)
		if replacement, ok := replacements[key]; ok {
			var asset NexRenderAsset
			
			switch mediaAsset.Type {
			case "image":
				asset = NexRenderAsset{
					Type:      "image",
					Src:       replacement.(string),
					LayerName: mediaAsset.Name,
				}
			case "video":
				asset = NexRenderAsset{
					Type:      "video",
					Src:       replacement.(string),
					LayerName: mediaAsset.Name,
				}
			case "audio":
				asset = NexRenderAsset{
					Type:      "audio",
					Src:       replacement.(string),
					LayerName: mediaAsset.Name,
				}
			}
			
			job.Assets = append(job.Assets, asset)
		}
	}

	// Add output action
	outputPath := filepath.Join(n.workDir, fmt.Sprintf("output_%d.mp4", time.Now().Unix()))
	job.Actions.Postrender = []NexRenderAction{
		{
			Module: "@nexrender/action-encode",
			Preset: "mp4",
			Output: outputPath,
			Params: map[string]interface{}{
				"preset": "youtube-1080p",
			},
		},
	}

	return job, nil
}

// CreateJobFromDeepAnalysis creates an optimized job from deep analysis
func (n *NexRenderIntegration) CreateJobFromDeepAnalysis(analysis *DeepAnalysisResult, config RenderConfig) (*NexRenderJob, error) {
	job := &NexRenderJob{
		Type: "advanced",
		Template: NexRenderTemplate{
			Src: analysis.Metadata.FilePath,
		},
		Assets:   []NexRenderAsset{},
		Settings: NexRenderSettings{},
	}

	// Select composition based on config
	if config.CompositionName != "" {
		job.Template.Composition = config.CompositionName
	} else if len(analysis.Metadata.Compositions) > 0 {
		job.Template.Composition = analysis.Metadata.Compositions[0].Name
	}

	// Add intelligent text replacements
	if analysis.TextIntelligence != nil {
		for _, field := range analysis.TextIntelligence.DynamicFields {
			if value, ok := config.TextReplacements[field.FieldName]; ok {
				asset := NexRenderAsset{
					Type:      "data",
					LayerName: field.LayerID,
					Property:  "Source Text",
					Value:     n.formatTextValue(value, field),
				}
				job.Assets = append(job.Assets, asset)
			}
		}
	}

	// Add smart media replacements
	if analysis.MediaMapping != nil {
		for _, replaceableAsset := range analysis.MediaMapping.ReplaceableAssets {
			key := n.generateReplacementKey(replaceableAsset.Name)
			if mediaURL, ok := config.MediaReplacements[key]; ok {
				asset := n.createMediaAsset(replaceableAsset, mediaURL)
				job.Assets = append(job.Assets, asset)
			}
		}
	}

	// Configure quality settings
	job.Settings = n.configureQualitySettings(config.Quality)

	// Add render actions
	job.Actions = n.createRenderActions(config)

	// Set job metadata
	if config.JobMetadata != nil {
		job.OnComplete = config.JobMetadata["oncomplete"]
		job.OnError = config.JobMetadata["onerror"]
	}

	return job, nil
}

// SubmitJob submits a job to nexrender
func (n *NexRenderIntegration) SubmitJob(job *NexRenderJob) (*JobStatus, error) {
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job: %w", err)
	}

	req, err := http.NewRequest("POST", n.apiURL+"/api/jobs", strings.NewReader(string(jobJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if n.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+n.apiKey)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to submit job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("job submission failed: %s", string(body))
	}

	var status JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

// GetJobStatus retrieves the status of a job
func (n *NexRenderIntegration) GetJobStatus(jobID string) (*JobStatus, error) {
	req, err := http.NewRequest("GET", n.apiURL+"/api/jobs/"+jobID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if n.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+n.apiKey)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get job status: %s", string(body))
	}

	var status JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

// CancelJob cancels a running job
func (n *NexRenderIntegration) CancelJob(jobID string) error {
	req, err := http.NewRequest("DELETE", n.apiURL+"/api/jobs/"+jobID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if n.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+n.apiKey)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to cancel job: %s", string(body))
	}

	return nil
}

// WaitForCompletion waits for a job to complete
func (n *NexRenderIntegration) WaitForCompletion(jobID string, timeout time.Duration) (*JobStatus, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status, err := n.GetJobStatus(jobID)
			if err != nil {
				return nil, err
			}

			switch status.State {
			case "finished":
				return status, nil
			case "error", "failed":
				return status, fmt.Errorf("job failed: %s", status.Error)
			}

			if time.Now().After(deadline) {
				return status, fmt.Errorf("timeout waiting for job completion")
			}
		}
	}
}

// BatchRender submits multiple render jobs
func (n *NexRenderIntegration) BatchRender(jobs []*NexRenderJob) ([]JobStatus, error) {
	statuses := make([]JobStatus, 0, len(jobs))
	
	for _, job := range jobs {
		status, err := n.SubmitJob(job)
		if err != nil {
			// Continue with other jobs even if one fails
			statuses = append(statuses, JobStatus{
				State: "error",
				Error: err.Error(),
			})
			continue
		}
		statuses = append(statuses, *status)
	}
	
	return statuses, nil
}

// Helper methods

func (n *NexRenderIntegration) generateReplacementKey(assetName string) string {
	// Clean asset name for use as key
	key := strings.ToLower(assetName)
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, ".", "_")
	return key + "_url"
}

func (n *NexRenderIntegration) formatTextValue(value interface{}, field DynamicTextField) interface{} {
	// Apply formatting based on field type
	switch field.DetectedType {
	case "date":
		// Could apply date formatting here
		return value
	case "price":
		// Could format as currency
		if numVal, ok := value.(float64); ok {
			return fmt.Sprintf("$%.2f", numVal)
		}
		return value
	default:
		return value
	}
}

func (n *NexRenderIntegration) createMediaAsset(replaceableAsset ReplaceableAsset, mediaURL string) NexRenderAsset {
	asset := NexRenderAsset{
		Type:      replaceableAsset.Type,
		Src:       mediaURL,
		LayerName: replaceableAsset.Name,
	}
	
	// Add smart replacement features
	if replaceableAsset.ReplacementType == "smart" {
		// Could add scaling, positioning parameters
		asset.Property = "Scale"
		asset.Value = []float64{100, 100} // Maintain aspect ratio
	}
	
	return asset
}

func (n *NexRenderIntegration) configureQualitySettings(quality string) NexRenderSettings {
	settings := NexRenderSettings{
		OutputExt: "mp4",
	}
	
	switch quality {
	case "ultra":
		settings.OutputModule = "h265"
		settings.SettingsTemplate = "Best Settings"
	case "high":
		settings.OutputModule = "h264"
		settings.SettingsTemplate = "High Quality"
	case "medium":
		settings.OutputModule = "h264"
		settings.SettingsTemplate = "Medium Quality"
	case "low":
		settings.OutputModule = "h264"
		settings.SettingsTemplate = "Draft Quality"
	default:
		settings.OutputModule = "h264"
		settings.SettingsTemplate = "High Quality"
	}
	
	return settings
}

func (n *NexRenderIntegration) createRenderActions(config RenderConfig) NexRenderActions {
	actions := NexRenderActions{
		Postrender: []NexRenderAction{},
	}
	
	// Add encoding action
	outputPath := config.OutputPath
	if outputPath == "" {
		outputPath = filepath.Join(n.workDir, fmt.Sprintf("render_%d.mp4", time.Now().Unix()))
	}
	
	encodeAction := NexRenderAction{
		Module: "@nexrender/action-encode",
		Output: outputPath,
		Preset: n.getPresetForQuality(config.Quality),
	}
	actions.Postrender = append(actions.Postrender, encodeAction)
	
	// Add upload action if specified
	if config.UploadDestination != "" {
		uploadAction := NexRenderAction{
			Module: "@nexrender/action-upload",
			Input:  outputPath,
			Params: map[string]interface{}{
				"provider": config.UploadProvider,
				"options":  config.UploadOptions,
			},
		}
		actions.Postrender = append(actions.Postrender, uploadAction)
	}
	
	// Add notification action
	if config.NotificationWebhook != "" {
		notifyAction := NexRenderAction{
			Module: "@nexrender/action-webhook",
			Params: map[string]interface{}{
				"url":    config.NotificationWebhook,
				"method": "POST",
			},
		}
		actions.Postrender = append(actions.Postrender, notifyAction)
	}
	
	return actions
}

func (n *NexRenderIntegration) getPresetForQuality(quality string) string {
	presetMap := map[string]string{
		"ultra":  "youtube-4k",
		"high":   "youtube-1080p",
		"medium": "youtube-720p",
		"low":    "youtube-480p",
	}
	
	if preset, ok := presetMap[quality]; ok {
		return preset
	}
	return "youtube-1080p"
}

// RenderConfig contains configuration for rendering
type RenderConfig struct {
	CompositionName     string                 `json:"composition_name,omitempty"`
	Quality             string                 `json:"quality"`
	TextReplacements    map[string]interface{} `json:"text_replacements"`
	MediaReplacements   map[string]string      `json:"media_replacements"`
	OutputPath          string                 `json:"output_path"`
	UploadDestination   string                 `json:"upload_destination,omitempty"`
	UploadProvider      string                 `json:"upload_provider,omitempty"`
	UploadOptions       map[string]interface{} `json:"upload_options,omitempty"`
	NotificationWebhook string                 `json:"notification_webhook,omitempty"`
	JobMetadata         map[string]string      `json:"job_metadata,omitempty"`
}

// LocalNexRender provides local nexrender execution
type LocalNexRender struct {
	binaryPath   string
	workDir      string
	aerenderPath string
}

// NewLocalNexRender creates a local nexrender executor
func NewLocalNexRender(binaryPath, workDir, aerenderPath string) *LocalNexRender {
	return &LocalNexRender{
		binaryPath:   binaryPath,
		workDir:      workDir,
		aerenderPath: aerenderPath,
	}
}

// RenderLocal executes nexrender locally
func (l *LocalNexRender) RenderLocal(job *NexRenderJob) error {
	// Save job to file
	jobPath := filepath.Join(l.workDir, fmt.Sprintf("job_%d.json", time.Now().Unix()))
	jobData, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}
	
	if err := os.WriteFile(jobPath, jobData, 0644); err != nil {
		return fmt.Errorf("failed to write job file: %w", err)
	}
	
	// Execute nexrender
	// This would typically use os/exec to run the nexrender binary
	// For now, we'll just return success
	
	return nil
}

// CreatePreviewJob creates a low-res preview render job
func CreatePreviewJob(metadata *ProjectMetadata) *NexRenderJob {
	job := &NexRenderJob{
		Type: "preview",
		Template: NexRenderTemplate{
			Src: metadata.FilePath,
		},
		Settings: NexRenderSettings{
			OutputModule: "h264",
			OutputExt:    "mp4",
			SettingsTemplate: "Draft Settings",
		},
	}
	
	// Limit to first 5 seconds for preview
	if len(metadata.Compositions) > 0 {
		job.Template.Composition = metadata.Compositions[0].Name
		job.Template.FrameStart = 0
		job.Template.FrameEnd = int(metadata.Compositions[0].FrameRate * 5) // 5 seconds
	}
	
	return job
}

// CreateThumbnailJob creates a job to extract thumbnail
func CreateThumbnailJob(metadata *ProjectMetadata, frame int) *NexRenderJob {
	job := &NexRenderJob{
		Type: "thumbnail",
		Template: NexRenderTemplate{
			Src:        metadata.FilePath,
			FrameStart: frame,
			FrameEnd:   frame,
		},
		Settings: NexRenderSettings{
			OutputModule: "JPEG",
			OutputExt:    "jpg",
		},
	}
	
	if len(metadata.Compositions) > 0 {
		job.Template.Composition = metadata.Compositions[0].Name
	}
	
	return job
}