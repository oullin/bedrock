package migrations

import (
	"context"
	"fmt"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// MigrationResult represents the outcome of a single migration execution.
type MigrationResult struct {
	Migration string
	Direction string // "up" or "down"
	Error     error
}

// Migrator orchestrates running and rolling back migrations.
type Migrator struct {
	repository Repository
	resolver   dbcontract.ConnectionResolver
	connection string
}

// NewMigrator creates a new Migrator.
func NewMigrator(repository Repository, resolver dbcontract.ConnectionResolver, connection string) *Migrator {
	return &Migrator{
		repository: repository,
		resolver:   resolver,
		connection: connection,
	}
}

// Run executes all pending migrations.
func (m *Migrator) Run(ctx context.Context, migrations []Migration) ([]MigrationResult, error) {
	ran, err := m.repository.GetRan(ctx)

	if err != nil {
		return nil, err
	}

	ranSet := make(map[string]bool, len(ran))

	for _, name := range ran {
		ranSet[name] = true
	}

	var pending []Migration

	for _, migration := range migrations {
		if !ranSet[migration.Name()] {
			pending = append(pending, migration)
		}
	}

	if len(pending) == 0 {
		return nil, nil
	}

	batch, err := m.repository.GetNextBatchNumber(ctx)

	if err != nil {
		return nil, err
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return nil, err
	}

	var results []MigrationResult

	for _, migration := range pending {
		err := migration.Up(ctx, conn)
		result := MigrationResult{Migration: migration.Name(), Direction: "up", Error: err}
		results = append(results, result)

		if err != nil {
			return results, fmt.Errorf("%w: %s: %v", ErrMigrationFailed, migration.Name(), err)
		}

		if err := m.repository.Log(ctx, migration.Name(), batch); err != nil {
			return results, err
		}
	}

	return results, nil
}

// Rollback rolls back the last batch of migrations.
func (m *Migrator) Rollback(ctx context.Context, migrations []Migration) ([]MigrationResult, error) {
	last, err := m.repository.GetLast(ctx)

	if err != nil {
		return nil, err
	}

	if len(last) == 0 {
		return nil, nil
	}

	migrationMap := make(map[string]Migration, len(migrations))

	for _, migration := range migrations {
		migrationMap[migration.Name()] = migration
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return nil, err
	}

	var results []MigrationResult

	for i := len(last) - 1; i >= 0; i-- {
		record := last[i]
		migration, ok := migrationMap[record.Migration]

		if !ok {
			results = append(results, MigrationResult{
				Migration: record.Migration, Direction: "down",
				Error: fmt.Errorf("%w: %s", ErrMigrationNotFound, record.Migration),
			})

			continue
		}

		err := migration.Down(ctx, conn)
		result := MigrationResult{Migration: record.Migration, Direction: "down", Error: err}
		results = append(results, result)

		if err != nil {
			return results, fmt.Errorf("%w: %s: %v", ErrMigrationFailed, record.Migration, err)
		}

		if err := m.repository.Delete(ctx, record.Migration); err != nil {
			return results, err
		}
	}

	return results, nil
}

// Reset rolls back ALL migrations.
func (m *Migrator) Reset(ctx context.Context, migrations []Migration) ([]MigrationResult, error) {
	ran, err := m.repository.GetRan(ctx)

	if err != nil {
		return nil, err
	}

	if len(ran) == 0 {
		return nil, nil
	}

	migrationMap := make(map[string]Migration, len(migrations))

	for _, migration := range migrations {
		migrationMap[migration.Name()] = migration
	}

	conn, err := m.resolver.Connection(ctx, m.connection)

	if err != nil {
		return nil, err
	}

	var results []MigrationResult

	for i := len(ran) - 1; i >= 0; i-- {
		name := ran[i]
		migration, ok := migrationMap[name]

		if !ok {
			continue
		}

		err := migration.Down(ctx, conn)
		result := MigrationResult{Migration: name, Direction: "down", Error: err}
		results = append(results, result)

		if err == nil {
			_ = m.repository.Delete(ctx, name)
		}
	}

	return results, nil
}

// Status returns the migration status (ran or pending).
func (m *Migrator) Status(ctx context.Context, migrations []Migration) ([]map[string]string, error) {
	ran, err := m.repository.GetRan(ctx)

	if err != nil {
		return nil, err
	}

	ranSet := make(map[string]bool, len(ran))

	for _, name := range ran {
		ranSet[name] = true
	}

	var status []map[string]string

	for _, migration := range migrations {
		s := "Pending"

		if ranSet[migration.Name()] {
			s = "Ran"
		}

		status = append(status, map[string]string{
			"migration": migration.Name(),
			"status":    s,
		})
	}

	return status, nil
}

// GetRepository returns the migration repository.
func (m *Migrator) GetRepository() Repository { return m.repository }
