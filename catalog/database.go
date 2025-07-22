// Package catalog provides database management for AEP project metadata
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

// Database manages the catalog database connection
type Database struct {
	db   *sql.DB
	path string
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := createDirIfNotExists(dir); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	database := &Database{
		db:   db,
		path: dbPath,
	}

	// Run migrations
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// migrate runs database migrations
func (d *Database) migrate() error {
	migrations := []string{
		createProjectsTable,
		createCompositionsTable,
		createTextLayersTable,
		createMediaAssetsTable,
		createEffectsTable,
		createCategoriesTable,
		createTagsTable,
		createProjectCategoriesTable,
		createProjectTagsTable,
		createOpportunitiesTable,
		createSearchIndexTable,
		createAnalysisResultsTable,
	}

	for i, migration := range migrations {
		if err := d.execMigration(migration, i+1); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}

// execMigration executes a single migration
func (d *Database) execMigration(migration string, version int) error {
	// Check if migration already applied
	var count int
	err := d.db.QueryRow(`
		SELECT COUNT(*) FROM sqlite_master 
		WHERE type='table' AND name='schema_migrations'
	`).Scan(&count)

	if err == nil && count == 0 {
		// Create migrations table
		_, err = d.db.Exec(`
			CREATE TABLE schema_migrations (
				version INTEGER PRIMARY KEY,
				applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create migrations table: %w", err)
		}
	}

	// Check if this version already applied
	err = d.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if count > 0 {
		return nil // Already applied
	}

	// Apply migration
	_, err = d.db.Exec(migration)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration
	_, err = d.db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	log.Printf("Applied migration %d", version)
	return nil
}

// StoreProject stores a project and all its metadata
func (d *Database) StoreProject(metadata *ProjectMetadata) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert project
	projectID, err := d.insertProject(tx, metadata)
	if err != nil {
		return fmt.Errorf("failed to insert project: %w", err)
	}

	// Insert compositions
	for _, comp := range metadata.Compositions {
		if err := d.insertComposition(tx, projectID, comp); err != nil {
			return fmt.Errorf("failed to insert composition: %w", err)
		}
	}

	// Insert text layers
	for _, text := range metadata.TextLayers {
		if err := d.insertTextLayer(tx, projectID, text); err != nil {
			return fmt.Errorf("failed to insert text layer: %w", err)
		}
	}

	// Insert media assets
	for _, asset := range metadata.MediaAssets {
		if err := d.insertMediaAsset(tx, projectID, asset); err != nil {
			return fmt.Errorf("failed to insert media asset: %w", err)
		}
	}

	// Insert effects
	for _, effect := range metadata.Effects {
		if err := d.insertEffect(tx, projectID, effect); err != nil {
			return fmt.Errorf("failed to insert effect: %w", err)
		}
	}

	// Insert categories and tags
	if err := d.insertProjectCategories(tx, projectID, metadata.Categories); err != nil {
		return fmt.Errorf("failed to insert categories: %w", err)
	}

	if err := d.insertProjectTags(tx, projectID, metadata.Tags); err != nil {
		return fmt.Errorf("failed to insert tags: %w", err)
	}

	// Insert opportunities
	for _, opp := range metadata.Opportunities {
		if err := d.insertOpportunity(tx, projectID, opp); err != nil {
			return fmt.Errorf("failed to insert opportunity: %w", err)
		}
	}

	// Create search index
	if err := d.createSearchIndex(tx, projectID, metadata); err != nil {
		return fmt.Errorf("failed to create search index: %w", err)
	}

	return tx.Commit()
}

// StoreAnalysisResult stores a deep analysis result
func (d *Database) StoreAnalysisResult(projectID int64, analysis *DeepAnalysisResult) error {
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

// GetProject retrieves a project by ID
func (d *Database) GetProject(projectID int64) (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{}

	// Get project
	err := d.db.QueryRow(`
		SELECT file_path, file_name, file_size, bit_depth, expression_engine, 
			   total_items, parsed_at
		FROM projects WHERE id = ?
	`, projectID).Scan(
		&metadata.FilePath, &metadata.FileName, &metadata.FileSize,
		&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
		&metadata.ParsedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Load related data
	metadata.Compositions, _ = d.getProjectCompositions(projectID)
	metadata.TextLayers, _ = d.getProjectTextLayers(projectID)
	metadata.MediaAssets, _ = d.getProjectMediaAssets(projectID)
	metadata.Effects, _ = d.getProjectEffects(projectID)
	metadata.Categories, _ = d.getProjectCategories(projectID)
	metadata.Tags, _ = d.getProjectTags(projectID)
	metadata.Opportunities, _ = d.getProjectOpportunities(projectID)

	return metadata, nil
}

// SearchProjects searches projects by query
func (d *Database) SearchProjects(query string, limit int) ([]*ProjectMetadata, error) {
	// Full-text search across indexed content
	rows, err := d.db.Query(`
		SELECT DISTINCT p.id, p.file_path, p.file_name, p.file_size,
			   p.bit_depth, p.expression_engine, p.total_items, p.parsed_at
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
		err := rows.Scan(
			&metadata.FilePath, &metadata.FileName, &metadata.FileSize,
			&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
			&metadata.ParsedAt)
		if err != nil {
			continue
		}
		results = append(results, metadata)
	}

	return results, nil
}

// FilterProjects filters projects by criteria
func (d *Database) FilterProjects(filter ProjectFilter) ([]*ProjectMetadata, error) {
	query := `
		SELECT DISTINCT p.id, p.file_path, p.file_name, p.file_size,
			   p.bit_depth, p.expression_engine, p.total_items, p.parsed_at
		FROM projects p
	`
	
	var conditions []string
	var args []interface{}
	
	// Add category filter
	if len(filter.Categories) > 0 {
		query += " JOIN project_categories pc ON p.id = pc.project_id"
		query += " JOIN categories c ON pc.category_id = c.id"
		placeholders := strings.Repeat("?,", len(filter.Categories)-1) + "?"
		conditions = append(conditions, fmt.Sprintf("c.name IN (%s)", placeholders))
		for _, cat := range filter.Categories {
			args = append(args, cat)
		}
	}
	
	// Add tag filter
	if len(filter.Tags) > 0 {
		query += " JOIN project_tags pt ON p.id = pt.project_id"
		query += " JOIN tags t ON pt.tag_id = t.id"
		placeholders := strings.Repeat("?,", len(filter.Tags)-1) + "?"
		conditions = append(conditions, fmt.Sprintf("t.name IN (%s)", placeholders))
		for _, tag := range filter.Tags {
			args = append(args, tag)
		}
	}
	
	// Add complexity filter
	if filter.MinComplexity > 0 || filter.MaxComplexity > 0 {
		query += " JOIN analysis_results ar ON p.id = ar.project_id"
		if filter.MinComplexity > 0 {
			conditions = append(conditions, "ar.complexity_score >= ?")
			args = append(args, filter.MinComplexity)
		}
		if filter.MaxComplexity > 0 {
			conditions = append(conditions, "ar.complexity_score <= ?")
			args = append(args, filter.MaxComplexity)
		}
	}
	
	// Add WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Add ordering and limit
	query += " ORDER BY p.parsed_at DESC LIMIT ?"
	args = append(args, filter.Limit)
	
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter projects: %w", err)
	}
	defer rows.Close()
	
	var results []*ProjectMetadata
	for rows.Next() {
		metadata := &ProjectMetadata{}
		err := rows.Scan(
			&metadata.FilePath, &metadata.FileName, &metadata.FileSize,
			&metadata.BitDepth, &metadata.ExpressionEngine, &metadata.TotalItems,
			&metadata.ParsedAt)
		if err != nil {
			continue
		}
		results = append(results, metadata)
	}
	
	return results, nil
}

// ProjectFilter defines search criteria
type ProjectFilter struct {
	Categories    []string
	Tags          []string
	MinComplexity float64
	MaxComplexity float64
	Limit         int
}

// Helper functions for inserting related data

func (d *Database) insertProject(tx *sql.Tx, metadata *ProjectMetadata) (int64, error) {
	result, err := tx.Exec(`
		INSERT OR REPLACE INTO projects (
			file_path, file_name, file_size, bit_depth, expression_engine,
			total_items, parsed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, metadata.FilePath, metadata.FileName, metadata.FileSize,
		metadata.BitDepth, metadata.ExpressionEngine, metadata.TotalItems,
		metadata.ParsedAt)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (d *Database) insertComposition(tx *sql.Tx, projectID int64, comp CompositionInfo) error {
	_, err := tx.Exec(`
		INSERT INTO compositions (
			project_id, comp_id, name, width, height, frame_rate, duration,
			layer_count, is_3d, has_effects
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, projectID, comp.ID, comp.Name, comp.Width, comp.Height,
		comp.FrameRate, comp.Duration, comp.LayerCount, comp.Is3D, comp.HasEffects)

	return err
}

func (d *Database) insertTextLayer(tx *sql.Tx, projectID int64, text TextLayerInfo) error {
	_, err := tx.Exec(`
		INSERT INTO text_layers (
			project_id, layer_id, comp_id, layer_name, source_text,
			font_used, is_animated, has_expressions
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, projectID, text.ID, text.CompID, text.LayerName, text.SourceText,
		text.FontUsed, text.IsAnimated, text.HasExpressions)

	return err
}

func (d *Database) insertMediaAsset(tx *sql.Tx, projectID int64, asset MediaAssetInfo) error {
	_, err := tx.Exec(`
		INSERT INTO media_assets (
			project_id, asset_id, name, type, path, is_placeholder, usage_count
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, projectID, asset.ID, asset.Name, asset.Type, asset.Path,
		asset.IsPlaceholder, asset.UsageCount)

	return err
}

func (d *Database) insertEffect(tx *sql.Tx, projectID int64, effect EffectInfo) error {
	_, err := tx.Exec(`
		INSERT INTO effects (
			project_id, name, category, usage_count, is_customizable
		) VALUES (?, ?, ?, ?, ?)
	`, projectID, effect.Name, effect.Category, effect.UsageCount, effect.IsCustomizable)

	return err
}

func (d *Database) insertProjectCategories(tx *sql.Tx, projectID int64, categories []string) error {
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

func (d *Database) insertProjectTags(tx *sql.Tx, projectID int64, tags []string) error {
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

func (d *Database) insertOpportunity(tx *sql.Tx, projectID int64, opp Opportunity) error {
	componentsJSON, _ := json.Marshal(opp.Components)
	_, err := tx.Exec(`
		INSERT INTO opportunities (
			project_id, type, description, difficulty, impact, components
		) VALUES (?, ?, ?, ?, ?, ?)
	`, projectID, opp.Type, opp.Description, opp.Difficulty, opp.Impact, string(componentsJSON))

	return err
}

func (d *Database) createSearchIndex(tx *sql.Tx, projectID int64, metadata *ProjectMetadata) error {
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

func (d *Database) getProjectCompositions(projectID int64) ([]CompositionInfo, error) {
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

func (d *Database) getProjectTextLayers(projectID int64) ([]TextLayerInfo, error) {
	rows, err := d.db.Query(`
		SELECT layer_id, comp_id, layer_name, source_text, font_used,
			   is_animated, has_expressions
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
			&text.FontUsed, &text.IsAnimated, &text.HasExpressions)
		if err != nil {
			continue
		}
		textLayers = append(textLayers, text)
	}

	return textLayers, nil
}

func (d *Database) getProjectMediaAssets(projectID int64) ([]MediaAssetInfo, error) {
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

func (d *Database) getProjectEffects(projectID int64) ([]EffectInfo, error) {
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

func (d *Database) getProjectCategories(projectID int64) ([]string, error) {
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

func (d *Database) getProjectTags(projectID int64) ([]string, error) {
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

func (d *Database) getProjectOpportunities(projectID int64) ([]Opportunity, error) {
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

func createDirIfNotExists(dir string) error {
	if dir == "" {
		return nil
	}
	return nil // os.MkdirAll would go here in real implementation
}

// Database schema definitions
const createProjectsTable = `
CREATE TABLE IF NOT EXISTS projects (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	file_path TEXT UNIQUE NOT NULL,
	file_name TEXT NOT NULL,
	file_size INTEGER NOT NULL,
	bit_depth INTEGER NOT NULL,
	expression_engine TEXT NOT NULL,
	total_items INTEGER NOT NULL,
	parsed_at DATETIME NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_projects_file_path ON projects(file_path);
CREATE INDEX IF NOT EXISTS idx_projects_parsed_at ON projects(parsed_at);
`

const createCompositionsTable = `
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

const createTextLayersTable = `
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
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_text_layers_project_id ON text_layers(project_id);
CREATE INDEX IF NOT EXISTS idx_text_layers_content ON text_layers(source_text);
`

const createMediaAssetsTable = `
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

const createEffectsTable = `
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

const createCategoriesTable = `
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
`

const createTagsTable = `
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
`

const createProjectCategoriesTable = `
CREATE TABLE IF NOT EXISTS project_categories (
	project_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	PRIMARY KEY (project_id, category_id),
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);
`

const createProjectTagsTable = `
CREATE TABLE IF NOT EXISTS project_tags (
	project_id INTEGER NOT NULL,
	tag_id INTEGER NOT NULL,
	PRIMARY KEY (project_id, tag_id),
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);
`

const createOpportunitiesTable = `
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

const createSearchIndexTable = `
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

const createAnalysisResultsTable = `
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