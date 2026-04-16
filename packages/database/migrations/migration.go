package migrations

import (
	"context"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// Migration defines a single database migration with up and down operations.
type Migration interface {
	// Name returns the unique migration name (e.g., "2024_01_01_create_users_table").
	Name() string
	// Up applies the migration.
	Up(ctx context.Context, conn dbcontract.Connection) error
	// Down reverts the migration.
	Down(ctx context.Context, conn dbcontract.Connection) error
}

// FuncMigration implements Migration using function values.
type FuncMigration struct {
	MigrationName string
	UpFunc        func(ctx context.Context, conn dbcontract.Connection) error
	DownFunc      func(ctx context.Context, conn dbcontract.Connection) error
}

// Name returns the migration name.
func (m *FuncMigration) Name() string { return m.MigrationName }

// Up runs the up migration.
func (m *FuncMigration) Up(ctx context.Context, conn dbcontract.Connection) error {
	if m.UpFunc == nil {
		return nil
	}
	return m.UpFunc(ctx, conn)
}

// Down runs the down migration.
func (m *FuncMigration) Down(ctx context.Context, conn dbcontract.Connection) error {
	if m.DownFunc == nil {
		return nil
	}
	return m.DownFunc(ctx, conn)
}
