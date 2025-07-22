// Package catalog provides PostgreSQL/Neon database implementation
package catalog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// PostgresDatabase implements DatabaseInterface for PostgreSQL/Neon
type PostgresDatabase struct {
	db     *sql.DB
	config DatabaseConfig
}

// NewPostgresDatabase creates a new PostgreSQL/Neon database connection
func NewPostgresDatabase(config DatabaseConfig) (*PostgresDatabase, error) {
	// Build connection string if not provided
	if config.ConnectionString == "" {
		return nil, fmt.Errorf("PostgreSQL connection string is required")
	}

	// Parse and ensure SSL mode is set for Neon
	connStr := config.ConnectionString
	if !strings.Contains(connStr, "sslmode=") {
		if strings.Contains(connStr, "?") {
			connStr += "&sslmode=" + config.SSLMode
		} else {
			connStr += "?sslmode=" + config.SSLMode
		}
	}

	// Add application name
	if config.ApplicationName != "" && !strings.Contains(connStr, "application_name=") {
		if strings.Contains(connStr, "?") {
			connStr += "&application_name=" + config.ApplicationName
		} else {
			connStr += "?application_name=" + config.ApplicationName
		}
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	pgDB := &PostgresDatabase{
		db:     db,
		config: config,
	}

	// Run migrations
	if err := pgDB.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Printf("Connected to PostgreSQL database (Neon)")
	return pgDB, nil
}

// Close closes the database connection
func (d *PostgresDatabase) Close() error {
	return d.db.Close()
}

// GetDatabaseType returns the database type
func (d *PostgresDatabase) GetDatabaseType() string {
	return "postgres"
}

// HealthCheck verifies database connectivity
func (d *PostgresDatabase) HealthCheck() error {
	return d.db.Ping()
}

// Migrate runs database migrations
func (d *PostgresDatabase) Migrate() error {
	// Create migrations table
	if err := d.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := []struct {
		version int
		name    string
		sql     string
	}{
		{1, "create_projects_table", createProjectsTablePG},
		{2, "create_compositions_table", createCompositionsTablePG},
		{3, "create_text_layers_table", createTextLayersTablePG},
		{4, "create_media_assets_table", createMediaAssetsTablePG},
		{5, "create_effects_table", createEffectsTablePG},
		{6, "create_categories_table", createCategoriesTablePG},
		{7, "create_tags_table", createTagsTablePG},
		{8, "create_project_categories_table", createProjectCategoriesTablePG},
		{9, "create_project_tags_table", createProjectTagsTablePG},
		{10, "create_opportunities_table", createOpportunitiesTablePG},
		{11, "create_search_index_table", createSearchIndexTablePG},
		{12, "create_analysis_results_table", createAnalysisResultsTablePG},
		{13, "add_s3_storage_fields", addS3StorageFieldsPG},
	}

	for _, migration := range migrations {
		if err := d.runMigration(migration.version, migration.name, migration.sql); err != nil {
			return fmt.Errorf("migration %d (%s) failed: %w", migration.version, migration.name, err)
		}
	}

	return nil
}

func (d *PostgresDatabase) createMigrationsTable() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (d *PostgresDatabase) runMigration(version int, name, sql string) error {
	// Check if migration already applied
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if count > 0 {
		return nil // Already applied
	}

	// Apply migration
	_, err = d.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration
	_, err = d.db.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", version, name)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	log.Printf("Applied migration %d: %s", version, name)
	return nil
}

// StoreProject stores a project and all its metadata
func (d *PostgresDatabase) StoreProject(metadata *ProjectMetadata) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert or update project
	var projectID int64
	err = tx.QueryRow(`
		INSERT INTO projects (
			file_path, file_name, file_size, bit_depth, expression_engine,
			total_items, parsed_at, s3_bucket, s3_key, s3_version_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (file_path) DO UPDATE SET
			file_name = EXCLUDED.file_name,
			file_size = EXCLUDED.file_size,
			bit_depth = EXCLUDED.bit_depth,
			expression_engine = EXCLUDED.expression_engine,
			total_items = EXCLUDED.total_items,
			parsed_at = EXCLUDED.parsed_at,
			s3_bucket = EXCLUDED.s3_bucket,
			s3_key = EXCLUDED.s3_key,
			s3_version_id = EXCLUDED.s3_version_id,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, metadata.FilePath, metadata.FileName, metadata.FileSize,
		metadata.BitDepth, metadata.ExpressionEngine, metadata.TotalItems,
		metadata.ParsedAt, metadata.S3Bucket, metadata.S3Key, metadata.S3VersionID).Scan(&projectID)

	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	// Store related data (same as SQLite implementation)
	// ... (implement storing compositions, text layers, etc.)

	return tx.Commit()
}

// GetProject retrieves a project by ID
func (d *PostgresDatabase) GetProject(projectID int64) (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{}

	err := d.db.QueryRow(`
		SELECT file_path, file_name, file_size, bit_depth, expression_engine, 
			   total_items, parsed_at, s3_bucket, s3_key, s3_version_id
		FROM projects WHERE id = $1
	`, projectID).Scan(
		&metadata.FilePath, &metadata.FileName, &metadata.FileSize,
		&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
		&metadata.ParsedAt, &metadata.S3Bucket, &metadata.S3Key, &metadata.S3VersionID)

	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Load related data
	// ... (implement loading compositions, text layers, etc.)

	return metadata, nil
}

// SearchProjects performs full-text search using PostgreSQL's capabilities
func (d *PostgresDatabase) SearchProjects(query string, limit int) ([]*ProjectMetadata, error) {
	// Use PostgreSQL's full-text search
	rows, err := d.db.Query(`
		SELECT DISTINCT p.id, p.file_path, p.file_name, p.file_size,
			   p.bit_depth, p.expression_engine, p.total_items, p.parsed_at,
			   p.s3_bucket, p.s3_key, p.s3_version_id
		FROM projects p
		JOIN search_index si ON p.id = si.project_id
		WHERE to_tsvector('english', si.content) @@ plainto_tsquery('english', $1)
		ORDER BY p.parsed_at DESC
		LIMIT $2
	`, query, limit)

	if err != nil {
		// Fallback to LIKE search if full-text search fails
		rows, err = d.db.Query(`
			SELECT DISTINCT p.id, p.file_path, p.file_name, p.file_size,
				   p.bit_depth, p.expression_engine, p.total_items, p.parsed_at,
				   p.s3_bucket, p.s3_key, p.s3_version_id
			FROM projects p
			JOIN search_index si ON p.id = si.project_id
			WHERE si.content ILIKE $1
			ORDER BY p.parsed_at DESC
			LIMIT $2
		`, "%"+query+"%", limit)
		
		if err != nil {
			return nil, fmt.Errorf("failed to search projects: %w", err)
		}
	}
	defer rows.Close()

	var results []*ProjectMetadata
	for rows.Next() {
		metadata := &ProjectMetadata{}
		var projectID int64
		err := rows.Scan(
			&projectID, &metadata.FilePath, &metadata.FileName, &metadata.FileSize,
			&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
			&metadata.ParsedAt, &metadata.S3Bucket, &metadata.S3Key, &metadata.S3VersionID)
		if err != nil {
			continue
		}
		results = append(results, metadata)
	}

	return results, nil
}

// FilterProjects filters projects by criteria
func (d *PostgresDatabase) FilterProjects(filter ProjectFilter) ([]*ProjectMetadata, error) {
	// Similar to SQLite implementation but with PostgreSQL syntax
	// ... (implement filtering logic)
	return nil, nil
}

// StoreAnalysisResult stores deep analysis results
func (d *PostgresDatabase) StoreAnalysisResult(projectID int64, analysis *DeepAnalysisResult) error {
	analysisJSON, err := json.Marshal(analysis)
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	_, err = d.db.Exec(`
		INSERT INTO analysis_results (
			project_id, complexity_score, automation_score, 
			analysis_data, created_at
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (project_id) DO UPDATE SET
			complexity_score = EXCLUDED.complexity_score,
			automation_score = EXCLUDED.automation_score,
			analysis_data = EXCLUDED.analysis_data,
			created_at = EXCLUDED.created_at
	`, projectID, analysis.ComplexityScore, analysis.AutomationScore, 
		string(analysisJSON), time.Now())

	return err
}

// GetAnalysisResult retrieves analysis results
func (d *PostgresDatabase) GetAnalysisResult(projectID int64) (*DeepAnalysisResult, error) {
	var analysisJSON string
	var result DeepAnalysisResult

	err := d.db.QueryRow(`
		SELECT analysis_data FROM analysis_results WHERE project_id = $1
	`, projectID).Scan(&analysisJSON)

	if err != nil {
		return nil, fmt.Errorf("failed to get analysis result: %w", err)
	}

	err = json.Unmarshal([]byte(analysisJSON), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal analysis: %w", err)
	}

	return &result, nil
}

// PostgreSQL migration scripts
const createProjectsTablePG = `
CREATE TABLE IF NOT EXISTS projects (
	id SERIAL PRIMARY KEY,
	file_path TEXT UNIQUE NOT NULL,
	file_name TEXT NOT NULL,
	file_size BIGINT NOT NULL,
	bit_depth INTEGER NOT NULL,
	expression_engine TEXT NOT NULL,
	total_items INTEGER NOT NULL,
	parsed_at TIMESTAMP NOT NULL,
	s3_bucket TEXT,
	s3_key TEXT,
	s3_version_id TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_projects_file_path ON projects(file_path);
CREATE INDEX IF NOT EXISTS idx_projects_parsed_at ON projects(parsed_at);
CREATE INDEX IF NOT EXISTS idx_projects_s3_key ON projects(s3_bucket, s3_key);
`

const createCompositionsTablePG = `
CREATE TABLE IF NOT EXISTS compositions (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	comp_id TEXT NOT NULL,
	name TEXT NOT NULL,
	width INTEGER NOT NULL,
	height INTEGER NOT NULL,
	frame_rate REAL NOT NULL,
	duration REAL NOT NULL,
	layer_count INTEGER NOT NULL,
	is_3d BOOLEAN NOT NULL DEFAULT FALSE,
	has_effects BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_compositions_project_id ON compositions(project_id);
CREATE INDEX IF NOT EXISTS idx_compositions_resolution ON compositions(width, height);
`

const createTextLayersTablePG = `
CREATE TABLE IF NOT EXISTS text_layers (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	layer_id TEXT NOT NULL,
	comp_id TEXT,
	layer_name TEXT NOT NULL,
	source_text TEXT NOT NULL,
	font_used TEXT,
	is_animated BOOLEAN NOT NULL DEFAULT FALSE,
	has_expressions BOOLEAN NOT NULL DEFAULT FALSE,
	is_3d BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_text_layers_project_id ON text_layers(project_id);
CREATE INDEX IF NOT EXISTS idx_text_layers_content ON text_layers(source_text);
`

const createMediaAssetsTablePG = `
CREATE TABLE IF NOT EXISTS media_assets (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	asset_id TEXT NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	path TEXT,
	is_placeholder BOOLEAN NOT NULL DEFAULT FALSE,
	usage_count INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_media_assets_project_id ON media_assets(project_id);
CREATE INDEX IF NOT EXISTS idx_media_assets_type ON media_assets(type);
CREATE INDEX IF NOT EXISTS idx_media_assets_placeholder ON media_assets(is_placeholder);
`

const createEffectsTablePG = `
CREATE TABLE IF NOT EXISTS effects (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	category TEXT,
	usage_count INTEGER NOT NULL DEFAULT 0,
	is_customizable BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_effects_project_id ON effects(project_id);
CREATE INDEX IF NOT EXISTS idx_effects_name ON effects(name);
`

const createCategoriesTablePG = `
CREATE TABLE IF NOT EXISTS categories (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
`

const createTagsTablePG = `
CREATE TABLE IF NOT EXISTS tags (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
`

const createProjectCategoriesTablePG = `
CREATE TABLE IF NOT EXISTS project_categories (
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
	PRIMARY KEY (project_id, category_id)
);
`

const createProjectTagsTablePG = `
CREATE TABLE IF NOT EXISTS project_tags (
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
	PRIMARY KEY (project_id, tag_id)
);
`

const createOpportunitiesTablePG = `
CREATE TABLE IF NOT EXISTS opportunities (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	type TEXT NOT NULL,
	description TEXT NOT NULL,
	difficulty TEXT NOT NULL,
	impact TEXT NOT NULL,
	components JSONB,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_opportunities_project_id ON opportunities(project_id);
CREATE INDEX IF NOT EXISTS idx_opportunities_type ON opportunities(type);
CREATE INDEX IF NOT EXISTS idx_opportunities_impact ON opportunities(impact);
`

const createSearchIndexTablePG = `
CREATE TABLE IF NOT EXISTS search_index (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_search_index_project_id ON search_index(project_id);
CREATE INDEX IF NOT EXISTS idx_search_index_content_gin ON search_index USING gin(to_tsvector('english', content));
`

const createAnalysisResultsTablePG = `
CREATE TABLE IF NOT EXISTS analysis_results (
	id SERIAL PRIMARY KEY,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	complexity_score REAL NOT NULL,
	automation_score REAL NOT NULL,
	analysis_data JSONB,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(project_id)
);

CREATE INDEX IF NOT EXISTS idx_analysis_results_project_id ON analysis_results(project_id);
CREATE INDEX IF NOT EXISTS idx_analysis_results_complexity ON analysis_results(complexity_score);
CREATE INDEX IF NOT EXISTS idx_analysis_results_automation ON analysis_results(automation_score);
`

const addS3StorageFieldsPG = `
-- Already added in projects table creation, but keeping for clarity
-- ALTER TABLE projects ADD COLUMN IF NOT EXISTS s3_bucket TEXT;
-- ALTER TABLE projects ADD COLUMN IF NOT EXISTS s3_key TEXT;
-- ALTER TABLE projects ADD COLUMN IF NOT EXISTS s3_version_id TEXT;
`