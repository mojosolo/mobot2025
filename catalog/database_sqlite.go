// Package catalog provides SQLite database implementation
package catalog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDatabase implements DatabaseInterface for SQLite
type SQLiteDatabase struct {
	db     *sql.DB
	path   string
	config DatabaseConfig
}

// NewSQLiteDatabase creates a new SQLite database connection
func NewSQLiteDatabase(config DatabaseConfig) (*SQLiteDatabase, error) {
	dbPath := config.ConnectionString
	if dbPath == "" {
		dbPath = "./catalog.db"
	}

	// Ensure directory exists
	if dbPath != ":memory:" {
		dir := filepath.Dir(dbPath)
		if err := createDirIfNotExists(dir); err != nil {
			return nil, fmt.Errorf("failed to create db directory: %w", err)
		}
	}

	// Open database with connection parameters
	connStr := dbPath
	if config.JournalMode != "" {
		connStr = fmt.Sprintf("%s?_journal_mode=%s", dbPath, config.JournalMode)
	}
	if config.BusyTimeout > 0 {
		if strings.Contains(connStr, "?") {
			connStr += fmt.Sprintf("&_busy_timeout=%d", config.BusyTimeout.Milliseconds())
		} else {
			connStr += fmt.Sprintf("?_busy_timeout=%d", config.BusyTimeout.Milliseconds())
		}
	}

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	sqliteDB := &SQLiteDatabase{
		db:     db,
		path:   dbPath,
		config: config,
	}

	// Run migrations
	if err := sqliteDB.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Printf("Connected to SQLite database: %s", dbPath)
	return sqliteDB, nil
}

// Close closes the database connection
func (d *SQLiteDatabase) Close() error {
	return d.db.Close()
}

// GetDatabaseType returns the database type
func (d *SQLiteDatabase) GetDatabaseType() string {
	return "sqlite"
}

// HealthCheck verifies database connectivity
func (d *SQLiteDatabase) HealthCheck() error {
	var result int
	return d.db.QueryRow("SELECT 1").Scan(&result)
}

// Migrate runs database migrations
func (d *SQLiteDatabase) Migrate() error {
	// Enable foreign keys
	_, err := d.db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	migrations := []struct {
		version int
		name    string
		sql     string
	}{
		{1, "create_projects_table", createProjectsTableSQLite},
		{2, "create_compositions_table", createCompositionsTableSQLite},
		{3, "create_text_layers_table", createTextLayersTableSQLite},
		{4, "create_media_assets_table", createMediaAssetsTableSQLite},
		{5, "create_effects_table", createEffectsTableSQLite},
		{6, "create_categories_table", createCategoriesTableSQLite},
		{7, "create_tags_table", createTagsTableSQLite},
		{8, "create_project_categories_table", createProjectCategoriesTableSQLite},
		{9, "create_project_tags_table", createProjectTagsTableSQLite},
		{10, "create_opportunities_table", createOpportunitiesTableSQLite},
		{11, "create_search_index_table", createSearchIndexTableSQLite},
		{12, "create_analysis_results_table", createAnalysisResultsTableSQLite},
		{13, "add_s3_storage_fields", addS3StorageFieldsSQLite},
	}

	// Create migrations table
	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, migration := range migrations {
		if err := d.runMigration(migration.version, migration.name, migration.sql); err != nil {
			return fmt.Errorf("migration %d (%s) failed: %w", migration.version, migration.name, err)
		}
	}

	return nil
}

func (d *SQLiteDatabase) runMigration(version int, name, sql string) error {
	// Check if migration already applied
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
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
	_, err = d.db.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", version, name)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	log.Printf("Applied migration %d: %s", version, name)
	return nil
}

// StoreProject stores a project and all its metadata
func (d *SQLiteDatabase) StoreProject(metadata *ProjectMetadata) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert project
	result, err := tx.Exec(`
		INSERT OR REPLACE INTO projects (
			file_path, file_name, file_size, bit_depth, expression_engine,
			total_items, parsed_at, s3_bucket, s3_key, s3_version_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, metadata.FilePath, metadata.FileName, metadata.FileSize,
		metadata.BitDepth, metadata.ExpressionEngine, metadata.TotalItems,
		metadata.ParsedAt, metadata.S3Bucket, metadata.S3Key, metadata.S3VersionID)

	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get project ID: %w", err)
	}

	// Store related data
	if err := d.storeProjectRelatedData(tx, projectID, metadata); err != nil {
		return err
	}

	return tx.Commit()
}

// GetProject retrieves a project by ID
func (d *SQLiteDatabase) GetProject(projectID int64) (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{}

	err := d.db.QueryRow(`
		SELECT file_path, file_name, file_size, bit_depth, expression_engine, 
			   total_items, parsed_at, s3_bucket, s3_key, s3_version_id
		FROM projects WHERE id = ?
	`, projectID).Scan(
		&metadata.FilePath, &metadata.FileName, &metadata.FileSize,
		&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
		&metadata.ParsedAt, &metadata.S3Bucket, &metadata.S3Key, &metadata.S3VersionID)

	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Load related data
	d.loadProjectRelatedData(projectID, metadata)

	return metadata, nil
}

// SearchProjects searches projects by query
func (d *SQLiteDatabase) SearchProjects(query string, limit int) ([]*ProjectMetadata, error) {
	rows, err := d.db.Query(`
		SELECT DISTINCT p.id, p.file_path, p.file_name, p.file_size,
			   p.bit_depth, p.expression_engine, p.total_items, p.parsed_at,
			   p.s3_bucket, p.s3_key, p.s3_version_id
		FROM projects p
		JOIN search_index si ON p.id = si.project_id
		WHERE si.content LIKE ?
		ORDER BY p.parsed_at DESC
		LIMIT ?
	`, "%"+query+"%", limit)

	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
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
func (d *SQLiteDatabase) FilterProjects(filter ProjectFilter) ([]*ProjectMetadata, error) {
	// Implementation similar to original Database.FilterProjects
	// ... (implement filtering logic)
	return nil, nil
}

// StoreAnalysisResult stores deep analysis results
func (d *SQLiteDatabase) StoreAnalysisResult(projectID int64, analysis *DeepAnalysisResult) error {
	analysisJSON, err := json.Marshal(analysis)
	if err != nil {
		return fmt.Errorf("failed to marshal analysis: %w", err)
	}

	_, err = d.db.Exec(`
		INSERT OR REPLACE INTO analysis_results (
			project_id, complexity_score, automation_score, 
			analysis_data, created_at
		) VALUES (?, ?, ?, ?, ?)
	`, projectID, analysis.ComplexityScore, analysis.AutomationScore, 
		string(analysisJSON), time.Now())

	return err
}

// GetAnalysisResult retrieves analysis results
func (d *SQLiteDatabase) GetAnalysisResult(projectID int64) (*DeepAnalysisResult, error) {
	var analysisJSON string
	var result DeepAnalysisResult

	err := d.db.QueryRow(`
		SELECT analysis_data FROM analysis_results WHERE project_id = ?
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

// Helper methods for storing and loading related data
func (d *SQLiteDatabase) storeProjectRelatedData(tx *sql.Tx, projectID int64, metadata *ProjectMetadata) error {
	// Store compositions
	for _, comp := range metadata.Compositions {
		if err := d.insertComposition(tx, projectID, comp); err != nil {
			return fmt.Errorf("failed to insert composition: %w", err)
		}
	}

	// Store text layers
	for _, text := range metadata.TextLayers {
		if err := d.insertTextLayer(tx, projectID, text); err != nil {
			return fmt.Errorf("failed to insert text layer: %w", err)
		}
	}

	// Store media assets
	for _, asset := range metadata.MediaAssets {
		if err := d.insertMediaAsset(tx, projectID, asset); err != nil {
			return fmt.Errorf("failed to insert media asset: %w", err)
		}
	}

	// Store effects
	for _, effect := range metadata.Effects {
		if err := d.insertEffect(tx, projectID, effect); err != nil {
			return fmt.Errorf("failed to insert effect: %w", err)
		}
	}

	// Store categories and tags
	if err := d.insertProjectCategories(tx, projectID, metadata.Categories); err != nil {
		return fmt.Errorf("failed to insert categories: %w", err)
	}

	if err := d.insertProjectTags(tx, projectID, metadata.Tags); err != nil {
		return fmt.Errorf("failed to insert tags: %w", err)
	}

	// Store opportunities
	for _, opp := range metadata.Opportunities {
		if err := d.insertOpportunity(tx, projectID, opp); err != nil {
			return fmt.Errorf("failed to insert opportunity: %w", err)
		}
	}

	// Create search index
	if err := d.createSearchIndex(tx, projectID, metadata); err != nil {
		return fmt.Errorf("failed to create search index: %w", err)
	}

	return nil
}

func (d *SQLiteDatabase) loadProjectRelatedData(projectID int64, metadata *ProjectMetadata) {
	metadata.Compositions, _ = d.getProjectCompositions(projectID)
	metadata.TextLayers, _ = d.getProjectTextLayers(projectID)
	metadata.MediaAssets, _ = d.getProjectMediaAssets(projectID)
	metadata.Effects, _ = d.getProjectEffects(projectID)
	metadata.Categories, _ = d.getProjectCategories(projectID)
	metadata.Tags, _ = d.getProjectTags(projectID)
	metadata.Opportunities, _ = d.getProjectOpportunities(projectID)
}

// SQLite-specific helper methods (similar to original implementation)
func (d *SQLiteDatabase) insertComposition(tx *sql.Tx, projectID int64, comp CompositionInfo) error {
	_, err := tx.Exec(`
		INSERT INTO compositions (
			project_id, comp_id, name, width, height, frame_rate, duration,
			layer_count, is_3d, has_effects
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, projectID, comp.ID, comp.Name, comp.Width, comp.Height,
		comp.FrameRate, comp.Duration, comp.LayerCount, comp.Is3D, comp.HasEffects)
	return err
}

func (d *SQLiteDatabase) insertTextLayer(tx *sql.Tx, projectID int64, text TextLayerInfo) error {
	_, err := tx.Exec(`
		INSERT INTO text_layers (
			project_id, layer_id, comp_id, layer_name, source_text,
			font_used, is_animated, has_expressions, is_3d
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, projectID, text.ID, text.CompID, text.LayerName, text.SourceText,
		text.FontUsed, text.IsAnimated, text.HasExpressions, text.Is3D)
	return err
}

func (d *SQLiteDatabase) insertMediaAsset(tx *sql.Tx, projectID int64, asset MediaAssetInfo) error {
	_, err := tx.Exec(`
		INSERT INTO media_assets (
			project_id, asset_id, name, type, path, is_placeholder, usage_count
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, projectID, asset.ID, asset.Name, asset.Type, asset.Path,
		asset.IsPlaceholder, asset.UsageCount)
	return err
}

func (d *SQLiteDatabase) insertEffect(tx *sql.Tx, projectID int64, effect EffectInfo) error {
	_, err := tx.Exec(`
		INSERT INTO effects (
			project_id, name, category, usage_count, is_customizable
		) VALUES (?, ?, ?, ?, ?)
	`, projectID, effect.Name, effect.Category, effect.UsageCount, effect.IsCustomizable)
	return err
}

func (d *SQLiteDatabase) insertProjectCategories(tx *sql.Tx, projectID int64, categories []string) error {
	for _, category := range categories {
		// Insert category if not exists
		_, err := tx.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
		if err != nil {
			return err
		}

		// Get category ID
		var categoryID int64
		err = tx.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryID)
		if err != nil {
			return err
		}

		// Link project to category
		_, err = tx.Exec("INSERT OR IGNORE INTO project_categories (project_id, category_id) VALUES (?, ?)", 
			projectID, categoryID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *SQLiteDatabase) insertProjectTags(tx *sql.Tx, projectID int64, tags []string) error {
	for _, tag := range tags {
		// Insert tag if not exists
		_, err := tx.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return err
		}

		// Get tag ID
		var tagID int64
		err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return err
		}

		// Link project to tag
		_, err = tx.Exec("INSERT OR IGNORE INTO project_tags (project_id, tag_id) VALUES (?, ?)", 
			projectID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *SQLiteDatabase) insertOpportunity(tx *sql.Tx, projectID int64, opp Opportunity) error {
	componentsJSON, _ := json.Marshal(opp.Components)
	_, err := tx.Exec(`
		INSERT INTO opportunities (
			project_id, type, description, difficulty, impact, components
		) VALUES (?, ?, ?, ?, ?, ?)
	`, projectID, opp.Type, opp.Description, opp.Difficulty, opp.Impact, string(componentsJSON))
	return err
}

func (d *SQLiteDatabase) createSearchIndex(tx *sql.Tx, projectID int64, metadata *ProjectMetadata) error {
	// Build searchable content
	content := []string{
		metadata.FileName,
		strings.Join(metadata.Categories, " "),
		strings.Join(metadata.Tags, " "),
	}

	// Add text layer content
	for _, text := range metadata.TextLayers {
		content = append(content, text.LayerName, text.SourceText)
	}

	// Add asset names
	for _, asset := range metadata.MediaAssets {
		content = append(content, asset.Name)
	}

	// Add effect names
	for _, effect := range metadata.Effects {
		content = append(content, effect.Name)
	}

	searchText := strings.Join(content, " ")

	_, err := tx.Exec(`
		INSERT OR REPLACE INTO search_index (project_id, content) VALUES (?, ?)
	`, projectID, searchText)

	return err
}

// Helper functions for retrieving related data
func (d *SQLiteDatabase) getProjectCompositions(projectID int64) ([]CompositionInfo, error) {
	rows, err := d.db.Query(`
		SELECT comp_id, name, width, height, frame_rate, duration,
			   layer_count, is_3d, has_effects
		FROM compositions WHERE project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compositions []CompositionInfo
	for rows.Next() {
		var comp CompositionInfo
		err := rows.Scan(&comp.ID, &comp.Name, &comp.Width, &comp.Height,
			&comp.FrameRate, &comp.Duration, &comp.LayerCount, &comp.Is3D, &comp.HasEffects)
		if err != nil {
			continue
		}
		compositions = append(compositions, comp)
	}

	return compositions, nil
}

func (d *SQLiteDatabase) getProjectTextLayers(projectID int64) ([]TextLayerInfo, error) {
	rows, err := d.db.Query(`
		SELECT layer_id, comp_id, layer_name, source_text, font_used,
			   is_animated, has_expressions, is_3d
		FROM text_layers WHERE project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var textLayers []TextLayerInfo
	for rows.Next() {
		var text TextLayerInfo
		err := rows.Scan(&text.ID, &text.CompID, &text.LayerName, &text.SourceText,
			&text.FontUsed, &text.IsAnimated, &text.HasExpressions, &text.Is3D)
		if err != nil {
			continue
		}
		textLayers = append(textLayers, text)
	}

	return textLayers, nil
}

func (d *SQLiteDatabase) getProjectMediaAssets(projectID int64) ([]MediaAssetInfo, error) {
	rows, err := d.db.Query(`
		SELECT asset_id, name, type, path, is_placeholder, usage_count
		FROM media_assets WHERE project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []MediaAssetInfo
	for rows.Next() {
		var asset MediaAssetInfo
		err := rows.Scan(&asset.ID, &asset.Name, &asset.Type, &asset.Path,
			&asset.IsPlaceholder, &asset.UsageCount)
		if err != nil {
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (d *SQLiteDatabase) getProjectEffects(projectID int64) ([]EffectInfo, error) {
	rows, err := d.db.Query(`
		SELECT name, category, usage_count, is_customizable
		FROM effects WHERE project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var effects []EffectInfo
	for rows.Next() {
		var effect EffectInfo
		err := rows.Scan(&effect.Name, &effect.Category, &effect.UsageCount, &effect.IsCustomizable)
		if err != nil {
			continue
		}
		effects = append(effects, effect)
	}

	return effects, nil
}

func (d *SQLiteDatabase) getProjectCategories(projectID int64) ([]string, error) {
	rows, err := d.db.Query(`
		SELECT c.name FROM categories c
		JOIN project_categories pc ON c.id = pc.category_id
		WHERE pc.project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			categories = append(categories, name)
		}
	}

	return categories, nil
}

func (d *SQLiteDatabase) getProjectTags(projectID int64) ([]string, error) {
	rows, err := d.db.Query(`
		SELECT t.name FROM tags t
		JOIN project_tags pt ON t.id = pt.tag_id
		WHERE pt.project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			tags = append(tags, name)
		}
	}

	return tags, nil
}

func (d *SQLiteDatabase) getProjectOpportunities(projectID int64) ([]Opportunity, error) {
	rows, err := d.db.Query(`
		SELECT type, description, difficulty, impact, components
		FROM opportunities WHERE project_id = ?
	`, projectID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opportunities []Opportunity
	for rows.Next() {
		var opp Opportunity
		var componentsJSON string
		err := rows.Scan(&opp.Type, &opp.Description, &opp.Difficulty, &opp.Impact, &componentsJSON)
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(componentsJSON), &opp.Components)
		opportunities = append(opportunities, opp)
	}

	return opportunities, nil
}

// SQLite migration scripts
const createProjectsTableSQLite = `
CREATE TABLE IF NOT EXISTS projects (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	file_path TEXT UNIQUE NOT NULL,
	file_name TEXT NOT NULL,
	file_size INTEGER NOT NULL,
	bit_depth INTEGER NOT NULL,
	expression_engine TEXT NOT NULL,
	total_items INTEGER NOT NULL,
	parsed_at DATETIME NOT NULL,
	s3_bucket TEXT,
	s3_key TEXT,
	s3_version_id TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_projects_file_path ON projects(file_path);
CREATE INDEX IF NOT EXISTS idx_projects_parsed_at ON projects(parsed_at);
CREATE INDEX IF NOT EXISTS idx_projects_s3_key ON projects(s3_bucket, s3_key);
`

const createCompositionsTableSQLite = `
CREATE TABLE IF NOT EXISTS compositions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	comp_id TEXT NOT NULL,
	name TEXT NOT NULL,
	width INTEGER NOT NULL,
	height INTEGER NOT NULL,
	frame_rate REAL NOT NULL,
	duration REAL NOT NULL,
	layer_count INTEGER NOT NULL,
	is_3d BOOLEAN NOT NULL DEFAULT 0,
	has_effects BOOLEAN NOT NULL DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_compositions_project_id ON compositions(project_id);
CREATE INDEX IF NOT EXISTS idx_compositions_resolution ON compositions(width, height);
`

const createTextLayersTableSQLite = `
CREATE TABLE IF NOT EXISTS text_layers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	layer_id TEXT NOT NULL,
	comp_id TEXT,
	layer_name TEXT NOT NULL,
	source_text TEXT NOT NULL,
	font_used TEXT,
	is_animated BOOLEAN NOT NULL DEFAULT 0,
	has_expressions BOOLEAN NOT NULL DEFAULT 0,
	is_3d BOOLEAN NOT NULL DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_text_layers_project_id ON text_layers(project_id);
CREATE INDEX IF NOT EXISTS idx_text_layers_content ON text_layers(source_text);
`

const createMediaAssetsTableSQLite = `
CREATE TABLE IF NOT EXISTS media_assets (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	asset_id TEXT NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	path TEXT,
	is_placeholder BOOLEAN NOT NULL DEFAULT 0,
	usage_count INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_media_assets_project_id ON media_assets(project_id);
CREATE INDEX IF NOT EXISTS idx_media_assets_type ON media_assets(type);
CREATE INDEX IF NOT EXISTS idx_media_assets_placeholder ON media_assets(is_placeholder);
`

const createEffectsTableSQLite = `
CREATE TABLE IF NOT EXISTS effects (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	category TEXT,
	usage_count INTEGER NOT NULL DEFAULT 0,
	is_customizable BOOLEAN NOT NULL DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_effects_project_id ON effects(project_id);
CREATE INDEX IF NOT EXISTS idx_effects_name ON effects(name);
`

const createCategoriesTableSQLite = `
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
`

const createTagsTableSQLite = `
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
`

const createProjectCategoriesTableSQLite = `
CREATE TABLE IF NOT EXISTS project_categories (
	project_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	PRIMARY KEY (project_id, category_id),
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);
`

const createProjectTagsTableSQLite = `
CREATE TABLE IF NOT EXISTS project_tags (
	project_id INTEGER NOT NULL,
	tag_id INTEGER NOT NULL,
	PRIMARY KEY (project_id, tag_id),
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);
`

const createOpportunitiesTableSQLite = `
CREATE TABLE IF NOT EXISTS opportunities (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	type TEXT NOT NULL,
	description TEXT NOT NULL,
	difficulty TEXT NOT NULL,
	impact TEXT NOT NULL,
	components TEXT, -- JSON array
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_opportunities_project_id ON opportunities(project_id);
CREATE INDEX IF NOT EXISTS idx_opportunities_type ON opportunities(type);
CREATE INDEX IF NOT EXISTS idx_opportunities_impact ON opportunities(impact);
`

const createSearchIndexTableSQLite = `
CREATE TABLE IF NOT EXISTS search_index (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_search_index_project_id ON search_index(project_id);
CREATE INDEX IF NOT EXISTS idx_search_index_content ON search_index(content);
`

const createAnalysisResultsTableSQLite = `
CREATE TABLE IF NOT EXISTS analysis_results (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	complexity_score REAL NOT NULL,
	automation_score REAL NOT NULL,
	analysis_data TEXT, -- JSON
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_analysis_results_project_id ON analysis_results(project_id);
CREATE INDEX IF NOT EXISTS idx_analysis_results_complexity ON analysis_results(complexity_score);
CREATE INDEX IF NOT EXISTS idx_analysis_results_automation ON analysis_results(automation_score);
`

const addS3StorageFieldsSQLite = `
-- Add S3 storage fields if they don't exist
-- SQLite doesn't support IF NOT EXISTS for columns, so we'll handle this in code
`