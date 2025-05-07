package riveradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/edkadigital/startmeup/pkg/migrations"
)

// ExampleTaskKind is the kind for example tasks
const ExampleTaskKind = "example_task"

// ExampleTask represents a simple example task
type ExampleTask struct {
	Message string
}

// Worker creates and manages River worker instances
type Worker struct {
	// We're simulating the implementation with a placeholder
	db *sql.DB
}

// NewWorker creates a new Worker for River
func NewWorker(db *sql.DB) (*Worker, error) {
	return &Worker{db: db}, nil
}

// Start starts the River worker
func (w *Worker) Start(ctx context.Context) {
	log.Println("Starting River worker")
}

// Stop stops the River worker
func (w *Worker) Stop(ctx context.Context) {
	log.Println("Stopping River worker")
}

// RegisterExampleTask registers the example task handler
func (w *Worker) RegisterExampleTask(handler func(ctx context.Context, task ExampleTask) error) error {
	log.Println("Registering example task handler")
	return nil
}

// InsertExampleTask inserts a new example task with the given message
func (w *Worker) InsertExampleTask(ctx context.Context, message string, delay time.Duration) error {
	log.Printf("Adding task with message: %s, delay: %v\n", message, delay)

	// Insert directly into the DB using SQL
	scheduledTime := time.Now().Add(delay)

	// Build the query with parameters
	query := `
	INSERT INTO river_jobs (
		id, 
		kind, 
		args, 
		scheduled_at, 
		created_at,
		updated_at,
		state,
		queue
	) VALUES (
		gen_random_uuid(), 
		$1, 
		$2, 
		$3, 
		NOW(),
		NOW(),
		'available',
		'default'
	)
	`

	// Convert task to JSON
	argsJSON := fmt.Sprintf(`{"Args":{"Message":"%s"}}`, message)

	// Execute the query
	_, err := w.db.ExecContext(ctx, query, ExampleTaskKind, argsJSON, scheduledTime)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	return nil
}

// MigrateDB runs the River migrations using our migration manager
func MigrateDB(ctx context.Context, db *sql.DB) error {
	// Import the migrations package
	migrationManager := migrations.NewRiverManager(db)
	return migrationManager.ApplyPendingMigrations(ctx)
}

// CreateMigrationSQL returns the SQL needed to create the River tables
func CreateMigrationSQL() string {
	return `
-- River Queue Schema

-- Create extension for UUID generation if not exists
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- river_jobs table stores all queued, running, and completed jobs
CREATE TABLE IF NOT EXISTS river_jobs (
    id UUID PRIMARY KEY,
    args JSONB NOT NULL,
    attempt INT NOT NULL DEFAULT 0,
    attempt_total INT NOT NULL DEFAULT 0,
    attempted_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    discard_reason TEXT,
    error JSONB,
    finalizer_priority INT,
    finalizer_retry_at TIMESTAMPTZ,
    finalizer_total_attempt_limit INT,
    finalized_at TIMESTAMPTZ,
    kind TEXT NOT NULL,
    max_attempts INT NOT NULL DEFAULT 25,
    metadata JSONB,
    priority INT NOT NULL DEFAULT 0,
    queue TEXT NOT NULL DEFAULT 'default',
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    state TEXT NOT NULL DEFAULT 'available',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    worker_id TEXT
);

-- Create indexes needed for efficient operations
CREATE INDEX IF NOT EXISTS river_jobs_kind_idx ON river_jobs (kind);
CREATE INDEX IF NOT EXISTS river_jobs_queue_idx ON river_jobs (queue);
CREATE INDEX IF NOT EXISTS river_jobs_scheduled_at_idx ON river_jobs (scheduled_at);
CREATE INDEX IF NOT EXISTS river_jobs_state_idx ON river_jobs (state);
CREATE INDEX IF NOT EXISTS river_jobs_updated_at_idx ON river_jobs (updated_at);

-- River client info table stores information about River clients
CREATE TABLE IF NOT EXISTS river_client_info (
    id TEXT PRIMARY KEY,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- River Leader clock keeps track of leadership
CREATE TABLE IF NOT EXISTS river_leader_clock (
    id INT PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL
);

-- River job period tracks periodic jobs and their last execution
CREATE TABLE IF NOT EXISTS river_job_periodic (
    id TEXT PRIMARY KEY,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    args JSONB NOT NULL,
    inserted_job_id UUID,
    interval_ms BIGINT NOT NULL,
    kind TEXT NOT NULL,
    last_inserted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    state TEXT NOT NULL DEFAULT 'enabled'
);
`
}

// RunMigrations runs the River migrations directly using SQL
func RunMigrations(db *sql.DB) error {
	sql := CreateMigrationSQL()
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	log.Println("River migrations executed successfully")
	return nil
}
