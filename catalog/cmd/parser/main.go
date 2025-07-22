package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    
    "github.com/mojosolo/mobot2025/catalog"
)

func main() {
    var (
        aepPath      = flag.String("file", "", "Path to AEP file")
        outputPath   = flag.String("output", "", "Output JSON path (default: stdout)")
        extractText  = flag.Bool("text", true, "Extract text layers")
        extractMedia = flag.Bool("media", true, "Extract media assets")
        deepAnalysis = flag.Bool("deep", true, "Perform deep analysis")
    )
    
    flag.Parse()
    
    if *aepPath == "" {
        log.Fatal("Please provide an AEP file path with -file")
    }
    
    // Create parser with options
    parser := catalog.NewParser()
    parser.ExtractText = *extractText
    parser.ExtractMedia = *extractMedia
    parser.DeepAnalysis = *deepAnalysis
    
    // Parse the project
    metadata, err := parser.ParseProject(*aepPath)
    if err != nil {
        log.Fatalf("Failed to parse project: %v", err)
    }
    
    // Convert to JSON
    jsonData, err := metadata.ToJSON()
    if err != nil {
        log.Fatalf("Failed to convert to JSON: %v", err)
    }
    
    // Output results
    if *outputPath != "" {
        err = os.WriteFile(*outputPath, jsonData, 0644)
        if err != nil {
            log.Fatalf("Failed to write output: %v", err)
        }
    } else {
        fmt.Println(string(jsonData))
    }
}