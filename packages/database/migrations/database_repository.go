package migrations

import (
	"context"
	"fmt"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// DatabaseRepository stores migration records in a database table.
type DatabaseRepository struct {
	conn  dbcontract.Connection
	table string
	// createDDL overrides the migrations-table CREATE statement when set.
	// It must contain a single %s placeholder for the table name. Drivers
	// whose default DDL is incompatible (e.g. ClickHouse, which requires a
	// table engine) supply their own via SetCreateDDL.
	createDDL string
	// existsSQL overrides the table-existence query when set. The query is
	// executed with the table name as its single positional binding.
	existsSQL string
}

var _ Repository = (*DatabaseRepository)(nil)

// NewDatabaseRepository creates a new DatabaseRepository.
func NewDatabaseRepository(conn dbcontract.Connection, table string) *DatabaseRepository {
	if table == "" {
		table = "migrations"
	}

	return &DatabaseRepository{conn: conn, table: table}
}

// SetCreateDDL overrides the migrations-table CREATE statement. The format
// string must contain exactly one %s placeholder for the table name.
func (r *DatabaseRepository) SetCreateDDL(format string) {
	r.createDDL = format
}

// SetExistsSQL overrides the table-existence lookup query. The query receives
// the table name as a single positional binding (`?`).
func (r *DatabaseRepository) SetExistsSQL(query string) {
	r.existsSQL = query
}

func (r *DatabaseRepository) GetRan(ctx context.Context) ([]string, error) {
	rows, err := r.conn.Select(ctx,
		fmt.Sprintf("select migration from %s order by batch, migration", r.table))

	if err != nil {
		return nil, err
	}

	var names []string

	for _, row := range rows {
		if name, ok := row["migration"]; ok {
			names = append(names, fmt.Sprintf("%v", name))
		}
	}

	return names, nil
}

func (r *DatabaseRepository) GetMigrationBatches(ctx context.Context) ([]MigrationRecord, error) {
	rows, err := r.conn.Select(ctx,
		fmt.Sprintf("select migration, batch from %s order by batch, migration", r.table))

	if err != nil {
		return nil, err
	}

	var records []MigrationRecord

	for _, row := range rows {
		records = append(records, MigrationRecord{
			Migration: fmt.Sprintf("%v", row["migration"]),
			Batch:     toInt(row["batch"]),
		})
	}

	return records, nil
}

func (r *DatabaseRepository) GetLast(ctx context.Context) ([]MigrationRecord, error) {
	batch, err := r.GetLastBatchNumber(ctx)

	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Select(ctx,
		fmt.Sprintf("select migration, batch from %s where batch = ? order by migration desc", r.table),
		batch)

	if err != nil {
		return nil, err
	}

	var records []MigrationRecord

	for _, row := range rows {
		records = append(records, MigrationRecord{
			Migration: fmt.Sprintf("%v", row["migration"]),
			Batch:     toInt(row["batch"]),
		})
	}

	return records, nil
}

func (r *DatabaseRepository) GetLastBatchNumber(ctx context.Context) (int, error) {
	row, err := r.conn.SelectOne(ctx,
		fmt.Sprintf("select max(batch) as batch from %s", r.table))

	if err != nil {
		return 0, err
	}

	if row == nil {
		return 0, nil
	}

	return toInt(row["batch"]), nil
}

func (r *DatabaseRepository) GetNextBatchNumber(ctx context.Context) (int, error) {
	last, err := r.GetLastBatchNumber(ctx)

	if err != nil {
		return 1, err
	}

	return last + 1, nil
}

func (r *DatabaseRepository) Log(ctx context.Context, name string, batch int) error {
	_, err := r.conn.Insert(ctx,
		fmt.Sprintf("insert into %s (migration, batch) values (?, ?)", r.table),
		name, batch)

	return err
}

func (r *DatabaseRepository) Delete(ctx context.Context, name string) error {
	_, err := r.conn.Delete(ctx,
		fmt.Sprintf("delete from %s where migration = ?", r.table),
		name)

	return err
}

func (r *DatabaseRepository) CreateRepository(ctx context.Context) error {
	format := r.createDDL

	if format == "" {
		format = dialectFor(r.conn.GetDriverName()).createDDL
	}

	if format == "" {
		format = `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL
		)
	`
	}

	_, err := r.conn.Statement(ctx, fmt.Sprintf(format, r.table))

	return err
}

func (r *DatabaseRepository) RepositoryExists(ctx context.Context) (bool, error) {
	query := r.existsSQL

	if query == "" {
		query = dialectFor(r.conn.GetDriverName()).existsSQL
	}

	if query == "" {
		query = "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	}

	rows, err := r.conn.Select(ctx, query, r.table)

	if err != nil {
		return false, err
	}

	return len(rows) > 0, nil
}

func (r *DatabaseRepository) DeleteRepository(ctx context.Context) error {
	_, err := r.conn.Statement(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", r.table))

	return err
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}
