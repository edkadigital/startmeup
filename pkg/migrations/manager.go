package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manager handles database migrations
type Manager struct {
	db              *sql.DB
	migrationsDir   string
	migrationsTable string
}

// NewManager creates a new migrations manager
func NewManager(db *sql.DB, migrationsDir, migrationsTable string) *Manager {
	return &Manager{
		db:              db,
		migrationsDir:   migrationsDir,
		migrationsTable: migrationsTable,
	}
}

// NewRiverManager creates a migrations manager specifically for River
func NewRiverManager(db *sql.DB) *Manager {
	return NewManager(db, "migrations/river", "river_migrations")
}

// EnsureMigrationsTable ensures the migrations table exists
func (m *Manager) EnsureMigrationsTable(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			migration_name TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`, m.migrationsTable)

	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// GetAppliedMigrations returns a list of already applied migrations
func (m *Manager) GetAppliedMigrations(ctx context.Context) ([]string, error) {
	// Make sure the migrations table exists
	if err := m.EnsureMigrationsTable(ctx); err != nil {
		return nil, err
	}

	// Query applied migrations
	query := fmt.Sprintf("SELECT migration_name FROM %s ORDER BY id", m.migrationsTable)
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	// Collect migration names
	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, name)
	}

	return migrations, nil
}

// GetPendingMigrations returns migrations that haven't been applied yet
func (m *Manager) GetPendingMigrations(ctx context.Context) ([]string, error) {
	// Get applied migrations
	applied, err := m.GetAppliedMigrations(ctx)
	if err != nil {
		return nil, err
	}

	// Build a map for faster lookup
	appliedMap := make(map[string]bool)
	for _, name := range applied {
		appliedMap[name] = true
	}

	// Get available migration files
	files, err := filepath.Glob(filepath.Join(m.migrationsDir, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to list migration files: %w", err)
	}

	// Collect pending migrations
	var pending []string
	for _, file := range files {
		name := filepath.Base(file)
		name = strings.TrimSuffix(name, ".sql")
		if !appliedMap[name] {
			pending = append(pending, name)
		}
	}

	// Sort migrations by name
	sort.Strings(pending)

	return pending, nil
}

// ApplyMigration applies a single migration
func (m *Manager) ApplyMigration(ctx context.Context, name string) error {
	// Read migration file
	filePath := filepath.Join(m.migrationsDir, name+".sql")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start a transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			if err == nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
			} else {
				log.Printf("rollback error suppressed: %v", rollbackErr)
			}
		}
	}()

	// Execute migration SQL
	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Make sure the migrations table wasn't dropped by the migration
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			migration_name TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`, m.migrationsTable)

	if _, err := tx.ExecContext(ctx, createTableQuery); err != nil {
		return fmt.Errorf("failed to ensure migrations table exists: %w", err)
	}

	// Record the migration
	insertQuery := fmt.Sprintf("INSERT INTO %s (migration_name) VALUES ($1)", m.migrationsTable)
	if _, err := tx.ExecContext(ctx, insertQuery, name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ApplyPendingMigrations applies all pending migrations
func (m *Manager) ApplyPendingMigrations(ctx context.Context) error {
	// Get pending migrations
	pending, err := m.GetPendingMigrations(ctx)
	if err != nil {
		return err
	}

	// Apply each migration
	for _, name := range pending {
		log.Printf("Applying migration: %s", name)
		if err := m.ApplyMigration(ctx, name); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", name, err)
		}
		log.Printf("Migration applied: %s", name)
	}

	return nil
}

// ApplyAllMigrations applies all migrations regardless of what's already applied
func (m *Manager) ApplyAllMigrations(ctx context.Context) error {
	// Get available migration files
	files, err := filepath.Glob(filepath.Join(m.migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	// Extract migration names
	var names []string
	for _, file := range files {
		name := filepath.Base(file)
		name = strings.TrimSuffix(name, ".sql")
		names = append(names, name)
	}

	// Sort migrations by name
	sort.Strings(names)

	// Apply each migration
	for _, name := range names {
		log.Printf("Applying migration: %s", name)
		if err := m.ApplyMigration(ctx, name); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", name, err)
		}
		log.Printf("Migration applied: %s", name)
	}

	return nil
}
