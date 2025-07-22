// Package catalog provides configuration management
package catalog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the catalog system
type Config struct {
	Database DatabaseConfig
	Storage  StorageConfig
	Server   ServerConfig
	Features FeatureFlags
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type      string // "s3" or "local"
	S3        S3Config
	LocalPath string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Environment  string
	LogLevel     string
	EnableCORS   bool
	TrustedHosts []string
}

// FeatureFlags holds feature toggles
type FeatureFlags struct {
	EnableS3Storage      bool
	EnableNeonDB         bool
	EnablePinecone       bool
	EnableDeepAnalysis   bool
	EnableAutoTagging    bool
	EnableWebhooks       bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Database: GetDefaultDatabaseConfig(),
		Storage:  loadStorageConfig(),
		Server:   loadServerConfig(),
		Features: loadFeatureFlags(),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// loadStorageConfig loads storage configuration from environment
func loadStorageConfig() StorageConfig {
	cfg := StorageConfig{
		Type:      getEnv("STORAGE_TYPE", "local"),
		LocalPath: getEnv("LOCAL_STORAGE_PATH", "./storage"),
	}

	// Load S3 config if enabled
	if cfg.Type == "s3" || getBoolEnv("AWS_S3_ENABLED", false) {
		cfg.Type = "s3"
		cfg.S3 = GetDefaultS3Config()
	}

	return cfg
}

// loadServerConfig loads server configuration from environment
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		EnableCORS:   getBoolEnv("ENABLE_CORS", true),
		TrustedHosts: strings.Split(getEnv("TRUSTED_HOSTS", "localhost"), ","),
	}
}

// loadFeatureFlags loads feature flags from environment
func loadFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableS3Storage:    getBoolEnv("AWS_S3_ENABLED", false),
		EnableNeonDB:       getEnv("MOBOT_DB_TYPE", "sqlite") == "postgres",
		EnablePinecone:     getEnv("PINECONE_HOST", "") != "",
		EnableDeepAnalysis: getBoolEnv("ENABLE_DEEP_ANALYSIS", true),
		EnableAutoTagging:  getBoolEnv("ENABLE_AUTO_TAGGING", true),
		EnableWebhooks:     getBoolEnv("ENABLE_WEBHOOKS", false),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate database config
	if c.Features.EnableNeonDB && c.Database.ConnectionString == "" {
		return fmt.Errorf("Neon database enabled but NEON_DATABASE_URL not set")
	}

	// Validate S3 config
	if c.Features.EnableS3Storage {
		if c.Storage.S3.AccessKeyID == "" || c.Storage.S3.SecretAccessKey == "" {
			return fmt.Errorf("S3 storage enabled but AWS credentials not set")
		}
		if c.Storage.S3.Bucket == "" {
			return fmt.Errorf("S3 storage enabled but AWS_BUCKET not set")
		}
	}

	// Validate server config
	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return fmt.Errorf("invalid port number: %s", c.Server.Port)
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetStorageInterface returns the appropriate storage implementation
func (c *Config) GetStorageInterface() (StorageInterface, error) {
	switch c.Storage.Type {
	case "s3":
		return NewS3Storage(c.Storage.S3)
	case "local":
		return NewLocalStorage(c.Storage.LocalPath)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", c.Storage.Type)
	}
}

// GetDatabaseInterface returns the appropriate database implementation
func (c *Config) GetDatabaseInterface() (DatabaseInterface, error) {
	return NewDatabaseFromConfig(c.Database)
}

// Helper functions

// getEnv returns an environment variable or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getBoolEnv returns a boolean environment variable or default value
func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	
	return parsed
}

// getDurationEnv returns a duration environment variable or default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	
	return parsed
}

// MustLoadConfig loads configuration or panics
func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}

// LogConfig logs the current configuration (with secrets masked)
func (c *Config) LogConfig() {
	fmt.Println("=== MoBot 2025 Configuration ===")
	fmt.Printf("Environment: %s\n", c.Server.Environment)
	fmt.Printf("Port: %s\n", c.Server.Port)
	fmt.Printf("Database Type: %s\n", c.Database.Type)
	fmt.Printf("Storage Type: %s\n", c.Storage.Type)
	
	fmt.Println("\nFeature Flags:")
	fmt.Printf("  S3 Storage: %v\n", c.Features.EnableS3Storage)
	fmt.Printf("  Neon DB: %v\n", c.Features.EnableNeonDB)
	fmt.Printf("  Pinecone: %v\n", c.Features.EnablePinecone)
	fmt.Printf("  Deep Analysis: %v\n", c.Features.EnableDeepAnalysis)
	fmt.Printf("  Auto Tagging: %v\n", c.Features.EnableAutoTagging)
	fmt.Printf("  Webhooks: %v\n", c.Features.EnableWebhooks)
	
	if c.Features.EnableS3Storage {
		fmt.Printf("\nS3 Configuration:\n")
		fmt.Printf("  Bucket: %s\n", c.Storage.S3.Bucket)
		fmt.Printf("  Region: %s\n", c.Storage.S3.Region)
		fmt.Printf("  Prefix: %s\n", c.Storage.S3.Prefix)
	}
	
	fmt.Println("================================")
}