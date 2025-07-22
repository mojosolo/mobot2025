// Package catalog provides database interfaces for flexible backend support
package catalog

import (
	"fmt"
	"os"
	"time"
)

// DatabaseInterface defines the contract for database implementations
type DatabaseInterface interface {
	// Core operations
	Close() error
	Migrate() error
	
	// Project operations
	StoreProject(metadata *ProjectMetadata) error
	GetProject(projectID int64) (*ProjectMetadata, error)
	SearchProjects(query string, limit int) ([]*ProjectMetadata, error)
	FilterProjects(filter ProjectFilter) ([]*ProjectMetadata, error)
	
	// Analysis operations
	StoreAnalysisResult(projectID int64, analysis *DeepAnalysisResult) error
	GetAnalysisResult(projectID int64) (*DeepAnalysisResult, error)
	
	// Utility operations
	GetDatabaseType() string
	HealthCheck() error
}

// DatabaseConfig holds configuration for database connections
type DatabaseConfig struct {
	Type           string // "sqlite" or "postgres"
	ConnectionString string
	MaxOpenConns   int
	MaxIdleConns   int
	ConnMaxLifetime time.Duration
	
	// PostgreSQL/Neon specific
	SSLMode        string
	ApplicationName string
	
	// SQLite specific
	JournalMode    string
	BusyTimeout    time.Duration
}

// NewDatabaseFromConfig creates a database instance based on configuration
func NewDatabaseFromConfig(config DatabaseConfig) (DatabaseInterface, error) {
	switch config.Type {
	case "sqlite":
		return NewSQLiteDatabase(config)
	case "postgres", "neon":
		return NewPostgresDatabase(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// GetDefaultConfig returns default configuration based on environment
func GetDefaultDatabaseConfig() DatabaseConfig {
	dbType := os.Getenv("MOBOT_DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite for backward compatibility
	}
	
	config := DatabaseConfig{
		Type:            dbType,
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
	}
	
	switch dbType {
	case "postgres", "neon":
		config.ConnectionString = os.Getenv("DATABASE_URL")
		if config.ConnectionString == "" {
			config.ConnectionString = os.Getenv("NEON_DATABASE_URL")
		}
		config.SSLMode = "require"
		config.ApplicationName = "mobot2025"
	case "sqlite":
		config.ConnectionString = os.Getenv("SQLITE_DB_PATH")
		if config.ConnectionString == "" {
			config.ConnectionString = "./catalog.db"
		}
		config.JournalMode = "WAL"
		config.BusyTimeout = 5 * time.Second
	}
	
	return config
}