package clickhouse

import (
	"github.com/bedrock/packages/database/migrations"
)

// ConfigureMigrationsRepository swaps the migrations-table DDL and existence
// query on the given repository for ones ClickHouse will accept. Apps wiring
// up migrations against a ClickHouse connection should call this before
// invoking CreateRepository or running the migrator:
//
//	repo := migrations.NewDatabaseRepository(conn, "migrations")
//	clickhouse.ConfigureMigrationsRepository(repo)
func ConfigureMigrationsRepository(repo *migrations.DatabaseRepository) {
	repo.SetCreateDDL(`
		CREATE TABLE IF NOT EXISTS %s (
			id Int64,
			migration String,
			batch Int32
		) ENGINE = MergeTree
		ORDER BY (batch, migration)
	`)
	repo.SetExistsSQL("SELECT name FROM system.tables WHERE database = currentDatabase() AND name = ?")
}
