package migrations

import (
	"database/sql"
	"fmt"
)

// Run creates the tables from laravel/laravel's default skeleton migrations.
func Run(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("demo migrations: nil database")
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			email_verified_at DATETIME,
			password TEXT NOT NULL,
			remember_token TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS password_reset_tokens (
			email TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			created_at DATETIME
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER,
			ip_address TEXT,
			user_agent TEXT,
			payload TEXT NOT NULL,
			last_activity INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS sessions_user_id_index ON sessions (user_id);
		CREATE INDEX IF NOT EXISTS sessions_last_activity_index ON sessions (last_activity);

		CREATE TABLE IF NOT EXISTS cache (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			expiration INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS cache_expiration_index ON cache (expiration);

		CREATE TABLE IF NOT EXISTS cache_locks (
			key TEXT PRIMARY KEY,
			owner TEXT NOT NULL,
			expiration INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS cache_locks_expiration_index ON cache_locks (expiration);

		CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			queue TEXT NOT NULL,
			payload TEXT NOT NULL,
			attempts INTEGER NOT NULL,
			reserved_at INTEGER,
			available_at INTEGER NOT NULL,
			created_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS jobs_queue_index ON jobs (queue);

		CREATE TABLE IF NOT EXISTS job_batches (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			total_jobs INTEGER NOT NULL,
			pending_jobs INTEGER NOT NULL,
			failed_jobs INTEGER NOT NULL,
			failed_job_ids TEXT NOT NULL,
			options TEXT,
			cancelled_at INTEGER,
			created_at INTEGER NOT NULL,
			finished_at INTEGER
		);

		CREATE TABLE IF NOT EXISTS failed_jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL UNIQUE,
			connection TEXT NOT NULL,
			queue TEXT NOT NULL,
			payload TEXT NOT NULL,
			exception TEXT NOT NULL,
			failed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return fmt.Errorf("demo migrations: %w", err)
	}

	return nil
}
