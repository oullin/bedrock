package migrations

import "errors"

var (
	// ErrMigrationNotFound is returned when a migration cannot be found.
	ErrMigrationNotFound = errors.New("migrations: migration not found")
	// ErrMigrationFailed is returned when a migration fails to execute.
	ErrMigrationFailed = errors.New("migrations: migration failed")
	// ErrNothingToMigrate is returned when there are no pending migrations.
	ErrNothingToMigrate = errors.New("migrations: nothing to migrate")
)
