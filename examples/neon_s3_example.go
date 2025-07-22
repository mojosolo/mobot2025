// Example of using MoBot 2025 with Neon and S3
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mojosolo/mobot2025/catalog"
)

func main() {
	// Load configuration from environment
	cfg, err := catalog.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Log configuration (secrets are masked)
	cfg.LogConfig()

	// Create database connection (Neon or SQLite based on config)
	db, err := cfg.GetDatabaseInterface()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create storage client (S3 or local based on config)
	storage, err := cfg.GetStorageInterface()
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	// Example: Parse and store an AEP file
	if len(os.Args) > 1 {
		aepPath := os.Args[1]
		if err := processAEPFile(aepPath, db, storage); err != nil {
			log.Fatalf("Failed to process AEP file: %v", err)
		}
	}

	// Example: Search projects
	fmt.Println("\nSearching for motion graphics projects...")
	results, err := db.SearchProjects("motion graphics", 10)
	if err != nil {
		log.Printf("Search failed: %v", err)
	} else {
		for _, project := range results {
			fmt.Printf("- %s (stored in: %s/%s)\n", 
				project.FileName, project.S3Bucket, project.S3Key)
		}
	}

	// Example: Get presigned URL for download
	if cfg.Features.EnableS3Storage && len(results) > 0 {
		ctx := context.Background()
		url, err := storage.GetURL(ctx, results[0].S3Key, 15*time.Minute)
		if err != nil {
			log.Printf("Failed to get download URL: %v", err)
		} else {
			fmt.Printf("\nDownload URL (valid for 15 minutes):\n%s\n", url)
		}
	}
}

func processAEPFile(aepPath string, db catalog.DatabaseInterface, storage catalog.StorageInterface) error {
	// Parse the AEP file
	fmt.Printf("Parsing %s...\n", aepPath)
	parser := catalog.NewParser()
	metadata, err := parser.ParseProject(aepPath)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	// Upload to S3 if enabled
	if storage != nil {
		fmt.Println("Uploading to S3...")
		ctx := context.Background()
		
		// Open the file
		file, err := os.Open(aepPath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Generate S3 key
		s3Key := fmt.Sprintf("projects/%s/%s", 
			time.Now().Format("2006/01/02"), 
			metadata.FileName)

		// Upload with metadata
		storageInfo, err := storage.Upload(ctx, s3Key, file, map[string]string{
			"project-name": metadata.FileName,
			"parsed-at":    metadata.ParsedAt.Format(time.RFC3339),
			"bit-depth":    fmt.Sprintf("%d", metadata.BitDepth),
		})
		if err != nil {
			return fmt.Errorf("failed to upload: %w", err)
		}

		// Update metadata with S3 info
		metadata.S3Bucket = storageInfo.Bucket
		metadata.S3Key = storageInfo.Key
		metadata.S3VersionID = storageInfo.VersionID

		fmt.Printf("Uploaded to S3: s3://%s/%s\n", 
			storageInfo.Bucket, storageInfo.Key)
	}

	// Store in database
	fmt.Println("Storing in database...")
	if err := db.StoreProject(metadata); err != nil {
		return fmt.Errorf("failed to store in database: %w", err)
	}

	// Perform analysis
	fmt.Println("Running deep analysis...")
	analyzer := catalog.NewDangerousAnalyzer()
	analysis, err := analyzer.AnalyzeProject(aepPath)
	if err != nil {
		log.Printf("Analysis failed: %v", err)
	} else {
		fmt.Printf("Complexity Score: %.2f\n", analysis.ComplexityScore)
		fmt.Printf("Automation Score: %.2f\n", analysis.AutomationScore)
		fmt.Printf("Opportunities: %d\n", len(analysis.Opportunities))
	}

	fmt.Println("✅ Successfully processed!")
	return nil
}