package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/yourusername/mobot2025/catalog"
)

func main() {
	// Define subcommands
	parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	catalogCmd := flag.NewFlagSet("catalog", flag.ExitOnError)
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	renderCmd := flag.NewFlagSet("render", flag.ExitOnError)
	
	// Parse flags
	parseFile := parseCmd.String("file", "", "AEP file to parse")
	parseOutput := parseCmd.String("output", "", "Output JSON file")
	
	// Analyze flags
	analyzeFile := analyzeCmd.String("file", "", "AEP file to analyze")
	analyzeDeep := analyzeCmd.Bool("deep", true, "Perform deep dangerous analysis")
	analyzeOutput := analyzeCmd.String("output", "", "Output analysis file")
	analyzeReport := analyzeCmd.Bool("report", false, "Generate human-readable report")
	
	// Catalog flags
	catalogDir := catalogCmd.String("dir", "", "Directory to catalog")
	catalogPattern := catalogCmd.String("pattern", "*.aep", "File pattern to match")
	catalogOutput := catalogCmd.String("output", "catalog.json", "Output catalog file")
	catalogReport := catalogCmd.Bool("report", false, "Generate catalog report")
	
	// Import flags
	importDir := importCmd.String("dir", "", "MoBot directory to import from")
	importOutput := importCmd.String("output", "import_report.md", "Output report file")
	importExport := importCmd.Bool("export", false, "Export import data to JSON")
	
	// Serve flags
	servePort := serveCmd.Int("port", 8080, "API server port")
	
	// Render flags
	renderFile := renderCmd.String("file", "", "AEP file to render")
	renderConfig := renderCmd.String("config", "", "Render configuration JSON")
	renderAPI := renderCmd.String("api", "", "NexRender API URL")
	renderLocal := renderCmd.Bool("local", false, "Use local nexrender")
	
	// Check if subcommand provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	
	// Parse subcommand
	switch os.Args[1] {
	case "parse":
		parseCmd.Parse(os.Args[2:])
		handleParse(*parseFile, *parseOutput)
		
	case "analyze":
		analyzeCmd.Parse(os.Args[2:])
		handleAnalyze(*analyzeFile, *analyzeDeep, *analyzeOutput, *analyzeReport)
		
	case "catalog":
		catalogCmd.Parse(os.Args[2:])
		handleCatalog(*catalogDir, *catalogPattern, *catalogOutput, *catalogReport)
		
	case "import":
		importCmd.Parse(os.Args[2:])
		handleImport(*importDir, *importOutput, *importExport)
		
	case "serve":
		serveCmd.Parse(os.Args[2:])
		handleServe(*servePort)
		
	case "render":
		renderCmd.Parse(os.Args[2:])
		handleRender(*renderFile, *renderConfig, *renderAPI, *renderLocal)
		
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`MoBot 2025 - Advanced AEP Processing Pipeline

Usage: mobot2025 <command> [options]

Commands:
  parse    Parse an AEP file and extract metadata
  analyze  Perform deep analysis on an AEP file
  catalog  Catalog a directory of AEP files
  import   Import templates from legacy MoBot system
  serve    Start the API server
  render   Submit a render job

Examples:
  mobot2025 parse -file template.aep -output metadata.json
  mobot2025 analyze -file template.aep -deep -report
  mobot2025 catalog -dir ./templates -output catalog.json
  mobot2025 import -dir ../mobot -output import_report.md
  mobot2025 serve -port 8080
  mobot2025 render -file template.aep -config render.json

For more help on a command: mobot2025 <command> -h`)
}

func handleParse(aepFile, outputFile string) {
	if aepFile == "" {
		log.Fatal("Please provide an AEP file with -file")
	}
	
	// Create parser
	parser := catalog.NewParser()
	
	// Parse the project
	fmt.Printf("Parsing %s...\n", aepFile)
	metadata, err := parser.ParseProject(aepFile)
	if err != nil {
		log.Fatalf("Failed to parse project: %v", err)
	}
	
	// Convert to JSON
	jsonData, err := metadata.ToJSON()
	if err != nil {
		log.Fatalf("Failed to convert to JSON: %v", err)
	}
	
	// Output results
	if outputFile != "" {
		err = os.WriteFile(outputFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
		fmt.Printf("Metadata saved to %s\n", outputFile)
	} else {
		fmt.Println(string(jsonData))
	}
	
	// Print summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Compositions: %d\n", len(metadata.Compositions))
	fmt.Printf("- Text Layers: %d\n", len(metadata.TextLayers))
	fmt.Printf("- Media Assets: %d\n", len(metadata.MediaAssets))
	fmt.Printf("- Effects: %d\n", len(metadata.Effects))
	fmt.Printf("- Opportunities: %d\n", len(metadata.Opportunities))
}

func handleAnalyze(aepFile string, deep bool, outputFile string, generateReport bool) {
	if aepFile == "" {
		log.Fatal("Please provide an AEP file with -file")
	}
	
	if !deep {
		// Use standard parser
		handleParse(aepFile, outputFile)
		return
	}
	
	// Create dangerous analyzer
	analyzer := catalog.NewDangerousAnalyzer()
	
	// Perform deep analysis
	fmt.Printf("Performing deep analysis on %s...\n", aepFile)
	analysis, err := analyzer.AnalyzeProject(aepFile)
	if err != nil {
		log.Fatalf("Failed to analyze project: %v", err)
	}
	
	// Save JSON output
	if outputFile != "" || !generateReport {
		jsonData, err := analysis.ToJSON()
		if err != nil {
			log.Fatalf("Failed to convert to JSON: %v", err)
		}
		
		outFile := outputFile
		if outFile == "" {
			outFile = strings.TrimSuffix(aepFile, filepath.Ext(aepFile)) + "_analysis.json"
		}
		
		err = os.WriteFile(outFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
		fmt.Printf("Analysis saved to %s\n", outFile)
	}
	
	// Generate report
	if generateReport {
		report := analysis.GenerateReport()
		reportFile := strings.TrimSuffix(aepFile, filepath.Ext(aepFile)) + "_report.md"
		if outputFile != "" && strings.HasSuffix(outputFile, ".md") {
			reportFile = outputFile
		}
		
		err = os.WriteFile(reportFile, []byte(report), 0644)
		if err != nil {
			log.Fatalf("Failed to write report: %v", err)
		}
		fmt.Printf("Report saved to %s\n", reportFile)
		
		// Also generate API documentation
		apiDoc := analysis.GenerateAPIDocumentation()
		apiFile := strings.TrimSuffix(aepFile, filepath.Ext(aepFile)) + "_api.md"
		err = os.WriteFile(apiFile, []byte(apiDoc), 0644)
		if err != nil {
			log.Fatalf("Failed to write API doc: %v", err)
		}
		fmt.Printf("API documentation saved to %s\n", apiFile)
	}
	
	// Print summary
	fmt.Printf("\nAnalysis Summary:\n")
	fmt.Printf("- Complexity Score: %.1f/100\n", analysis.ComplexityScore)
	fmt.Printf("- Automation Score: %.1f/100\n", analysis.AutomationScore)
	if analysis.ModularSystem != nil {
		fmt.Printf("- Modular Components: %d\n", analysis.ModularSystem.TotalModules)
		fmt.Printf("- Variant Potential: %d\n", analysis.ModularSystem.VariantPotential)
	}
	if analysis.TextIntelligence != nil {
		fmt.Printf("- Dynamic Text Fields: %d\n", len(analysis.TextIntelligence.DynamicFields))
	}
	if analysis.MediaMapping != nil {
		fmt.Printf("- Replaceable Assets: %d\n", len(analysis.MediaMapping.ReplaceableAssets))
	}
	fmt.Printf("- Recommendations: %d\n", len(analysis.Recommendations))
}

func handleCatalog(directory, pattern, outputFile string, generateReport bool) {
	if directory == "" {
		log.Fatal("Please provide a directory with -dir")
	}
	
	// Create parser
	parser := catalog.NewParser()
	
	fmt.Printf("Cataloging %s with pattern %s...\n", directory, pattern)
	
	// Find all matching files
	var results []*catalog.ProjectMetadata
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if matched {
				fmt.Printf("Processing %s...\n", path)
				metadata, err := parser.ParseProject(path)
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", path, err)
					return nil // Continue with other files
				}
				results = append(results, metadata)
			}
		}
		return nil
	})
	
	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}
	
	fmt.Printf("Cataloged %d projects\n", len(results))
	
	// Create catalog structure
	catalogData := map[string]interface{}{
		"version":    "1.0.0",
		"total":      len(results),
		"generated":  filepath.Base(outputFile),
		"templates":  results,
	}
	
	// Save catalog
	jsonData, err := json.MarshalIndent(catalogData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to create catalog JSON: %v", err)
	}
	
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write catalog: %v", err)
	}
	fmt.Printf("Catalog saved to %s\n", outputFile)
	
	// Generate report if requested
	if generateReport {
		reportContent := generateCatalogReport(results)
		reportFile := strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + "_report.md"
		err = os.WriteFile(reportFile, []byte(reportContent), 0644)
		if err != nil {
			log.Fatalf("Failed to write report: %v", err)
		}
		fmt.Printf("Report saved to %s\n", reportFile)
	}
}

func generateCatalogReport(results []*catalog.ProjectMetadata) string {
	report := "# AEP Template Catalog Report\n\n"
	report += fmt.Sprintf("Total Templates: %d\n\n", len(results))
	
	// Category summary
	categories := make(map[string]int)
	capabilities := make(map[string]int)
	
	for _, metadata := range results {
		for _, cat := range metadata.Categories {
			categories[cat]++
		}
		
		if metadata.Capabilities.HasTextReplacement {
			capabilities["Text Replacement"]++
		}
		if metadata.Capabilities.HasImageReplacement {
			capabilities["Image Replacement"]++
		}
		if metadata.Capabilities.IsModular {
			capabilities["Modular"]++
		}
	}
	
	report += "## Categories\n\n"
	for cat, count := range categories {
		report += fmt.Sprintf("- %s: %d templates\n", cat, count)
	}
	
	report += "\n## Capabilities\n\n"
	for cap, count := range capabilities {
		report += fmt.Sprintf("- %s: %d templates\n", cap, count)
	}
	
	report += "\n## Template Details\n\n"
	for _, metadata := range results {
		report += fmt.Sprintf("### %s\n\n", metadata.FileName)
		report += fmt.Sprintf("- Path: `%s`\n", metadata.FilePath)
		report += fmt.Sprintf("- Compositions: %d\n", len(metadata.Compositions))
		report += fmt.Sprintf("- Text Layers: %d\n", len(metadata.TextLayers))
		report += fmt.Sprintf("- Categories: %s\n", strings.Join(metadata.Categories, ", "))
		report += fmt.Sprintf("- Opportunities: %d\n\n", len(metadata.Opportunities))
	}
	
	return report
}

func handleServe(port int) {
	fmt.Printf("Starting AEP Catalog API on port %d...\n", port)
	
	// Create API service with database
	dbPath := "./catalog.db"
	service, err := catalog.NewAPIService(port, dbPath)
	if err != nil {
		log.Fatalf("Failed to create API service: %v", err)
	}
	
	fmt.Printf("Database: %s\n", dbPath)
	
	// Start server
	if err := service.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRender(aepFile, configFile, apiURL string, localMode bool) {
	if aepFile == "" {
		log.Fatal("Please provide an AEP file with -file")
	}
	
	if configFile == "" {
		log.Fatal("Please provide a render configuration with -config")
	}
	
	// Read render configuration
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	
	var config catalog.RenderConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	
	// Parse the AEP file first
	parser := catalog.NewParser()
	metadata, err := parser.ParseProject(aepFile)
	if err != nil {
		log.Fatalf("Failed to parse project: %v", err)
	}
	
	if localMode {
		// Local rendering
		fmt.Println("Local rendering not yet implemented")
		return
	}
	
	// API rendering
	if apiURL == "" {
		apiURL = "http://localhost:3000" // Default nexrender API
	}
	
	// Create nexrender integration
	nexrender := catalog.NewNexRenderIntegration(apiURL, "", "./renders")
	
	// Create render job
	replacements := make(map[string]interface{})
	for k, v := range config.TextReplacements {
		replacements[k] = v
	}
	for k, v := range config.MediaReplacements {
		replacements[k] = v
	}
	
	job, err := nexrender.CreateJobFromMetadata(metadata, replacements)
	if err != nil {
		log.Fatalf("Failed to create job: %v", err)
	}
	
	// Submit job
	fmt.Println("Submitting render job...")
	status, err := nexrender.SubmitJob(job)
	if err != nil {
		log.Fatalf("Failed to submit job: %v", err)
	}
	
	fmt.Printf("Job submitted: %s\n", status.UID)
	fmt.Printf("Status: %s\n", status.State)
	
	// Wait for completion
	fmt.Println("Waiting for render to complete...")
	finalStatus, err := nexrender.WaitForCompletion(status.UID, 30*60) // 30 minute timeout
	if err != nil {
		log.Fatalf("Render failed: %v", err)
	}
	
	fmt.Printf("Render completed!\n")
	fmt.Printf("Output: %s\n", finalStatus.Output)
}

func handleImport(directory, outputFile string, exportJSON bool) {
	if directory == "" {
		log.Fatal("Please provide a MoBot directory with -dir")
	}
	
	fmt.Printf("🚀 Starting import from MoBot directory: %s\n", directory)
	
	// Create database connection
	dbPath := "./catalog.db"
	database, err := catalog.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer database.Close()
	
	// Create importer
	importer := catalog.NewTemplateImporter(database)
	
	// Import templates
	result, err := importer.ImportFromMobot(directory)
	if err != nil {
		log.Fatalf("Import failed: %v", err)
	}
	
	// Generate report
	report := importer.GenerateImportReport(result)
	
	// Save report
	if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}
	
	fmt.Printf("📋 Import report saved to: %s\n", outputFile)
	
	// Export JSON data if requested
	if exportJSON {
		jsonFile := strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + ".json"
		if err := importer.ExportImportedData(result, jsonFile); err != nil {
			log.Fatalf("Failed to export JSON: %v", err)
		}
	}
	
	// Print summary
	fmt.Printf("\n✨ Import Summary:\n")
	fmt.Printf("   - Templates processed: %d\n", result.TotalFound)
	fmt.Printf("   - Successfully imported: %d\n", result.SuccessfulImports)
	fmt.Printf("   - Duration: %v\n", result.Duration)
	fmt.Printf("   - Success rate: %.1f%%\n", float64(result.SuccessfulImports)/float64(result.TotalFound)*100)
}