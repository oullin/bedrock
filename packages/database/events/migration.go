package events

// MigrationStarted is dispatched when a single migration starts running.
type MigrationStarted struct {
	Migration string
	Direction string // "up" or "down"
}

// MigrationEnded is dispatched when a single migration finishes.
type MigrationEnded struct {
	Migration string
	Direction string
}

// MigrationSkipped is dispatched when a migration is skipped.
type MigrationSkipped struct {
	Migration string
	Direction string
}

// MigrationsStarted is dispatched before a batch of migrations begins.
type MigrationsStarted struct {
	Direction string
}

// MigrationsEnded is dispatched after a batch of migrations completes.
type MigrationsEnded struct {
	Direction string
}

// NoPendingMigrations is dispatched when there are no migrations to run.
type NoPendingMigrations struct {
	Direction string
}

// MigrationsPruned is dispatched after old migrations are pruned.
type MigrationsPruned struct{}
