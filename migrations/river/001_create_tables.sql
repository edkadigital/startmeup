-- Migration: 001_create_tables.sql
-- River Queue Schema - Initial Migration

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

-- Create a migrations table to track applied migrations
CREATE TABLE IF NOT EXISTS river_migrations (
    id SERIAL PRIMARY KEY,
    migration_name TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Record this migration
INSERT INTO river_migrations (migration_name) VALUES ('001_create_tables');