package migrations

import "context"

// MigrationRecord represents a row in the migrations tracking table.
type MigrationRecord struct {
	Migration string
	Batch     int
}

// Repository tracks which migrations have been run.
type Repository interface {
	// GetRan returns the names of all migrations that have been run.
	GetRan(ctx context.Context) ([]string, error)
	// GetMigrationBatches returns migrations grouped by batch number.
	GetMigrationBatches(ctx context.Context) ([]MigrationRecord, error)
	// GetLast returns the migrations from the last batch.
	GetLast(ctx context.Context) ([]MigrationRecord, error)
	// GetLastBatchNumber returns the last batch number.
	GetLastBatchNumber(ctx context.Context) (int, error)
	// GetNextBatchNumber returns the next batch number.
	GetNextBatchNumber(ctx context.Context) (int, error)
	// Log records that a migration was run.
	Log(ctx context.Context, name string, batch int) error
	// Delete records that a migration was rolled back.
	Delete(ctx context.Context, name string) error
	// CreateRepository creates the migrations table if it doesn't exist.
	CreateRepository(ctx context.Context) error
	// RepositoryExists checks if the migrations table exists.
	RepositoryExists(ctx context.Context) (bool, error)
	// DeleteRepository drops the migrations table.
	DeleteRepository(ctx context.Context) error
}
