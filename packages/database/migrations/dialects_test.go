package migrations

import (
	"context"
	"strings"
	"testing"

	dbcontract "github.com/bedrock/packages/contracts/database"
)

// fakeConn implements just enough of dbcontract.Connection to drive
// DatabaseRepository.CreateRepository and RepositoryExists. It records every
// statement/select call so tests can assert the dialect-keyed SQL.
type fakeConn struct {
	driver       string
	statements   []string
	selectQuery  string
	selectArgs   []any
	selectResult []map[string]any
}

func (c *fakeConn) Statement(_ context.Context, query string, _ ...any) (bool, error) {
	c.statements = append(c.statements, query)

	return true, nil
}

func (c *fakeConn) Select(_ context.Context, query string, bindings ...any) ([]map[string]any, error) {
	c.selectQuery = query
	c.selectArgs = bindings

	return c.selectResult, nil
}

func (c *fakeConn) GetDriverName() string { return c.driver }

// The remaining methods are never called by the code paths under test.
func (c *fakeConn) Table(context.Context, string, ...string) any { return nil }
func (c *fakeConn) Raw(string) dbcontract.Expression             { return nil }
func (c *fakeConn) SelectOne(context.Context, string, ...any) (map[string]any, error) {
	return nil, nil
}
func (c *fakeConn) Insert(context.Context, string, ...any) (bool, error)  { return false, nil }
func (c *fakeConn) Update(context.Context, string, ...any) (int64, error) { return 0, nil }
func (c *fakeConn) Delete(context.Context, string, ...any) (int64, error) { return 0, nil }
func (c *fakeConn) AffectingStatement(context.Context, string, ...any) (int64, error) {
	return 0, nil
}
func (c *fakeConn) Unprepared(context.Context, string) (bool, error) { return false, nil }
func (c *fakeConn) PrepareBindings(b []any) []any                    { return b }
func (c *fakeConn) Transaction(context.Context, func(dbcontract.Connection) error, ...int) error {
	return nil
}
func (c *fakeConn) BeginTransaction(context.Context) error { return nil }
func (c *fakeConn) Commit(context.Context) error           { return nil }
func (c *fakeConn) Rollback(context.Context, ...int) error { return nil }
func (c *fakeConn) TransactionLevel() int                  { return 0 }
func (c *fakeConn) AfterCommit(func())                     {}
func (c *fakeConn) GetTablePrefix() string                 { return "" }
func (c *fakeConn) GetDatabaseName() string                { return "" }
func (c *fakeConn) GetName() string                        { return "" }
func (c *fakeConn) GetConfig(string) any                   { return nil }

func TestCreateRepository_UsesDialectForRegisteredDrivers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		driver      string
		mustContain []string
	}{
		{"pgsql", []string{"BIGSERIAL", "migrations"}},
		{"mysql", []string{"AUTO_INCREMENT", "ENGINE=InnoDB"}},
		{"mariadb", []string{"AUTO_INCREMENT", "ENGINE=InnoDB"}},
		{"sqlite", []string{"AUTOINCREMENT"}},
		{"clickhouse", []string{"MergeTree", "ORDER BY"}},
	}

	for _, tc := range cases {
		conn := &fakeConn{driver: tc.driver}
		repo := NewDatabaseRepository(conn, "migrations")

		if err := repo.CreateRepository(context.Background()); err != nil {
			t.Fatalf("%s: CreateRepository: %v", tc.driver, err)
		}

		if len(conn.statements) != 1 {
			t.Fatalf("%s: expected 1 statement, got %d", tc.driver, len(conn.statements))
		}

		for _, frag := range tc.mustContain {
			if !strings.Contains(conn.statements[0], frag) {
				t.Errorf("%s: expected DDL to contain %q, got %q", tc.driver, frag, conn.statements[0])
			}
		}
	}
}

func TestRepositoryExists_UsesDialectForRegisteredDrivers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		driver      string
		mustContain string
	}{
		{"pgsql", "pg_catalog.pg_tables"},
		{"mysql", "information_schema.tables"},
		{"mariadb", "information_schema.tables"},
		{"sqlite", "sqlite_master"},
		{"clickhouse", "system.tables"},
	}

	for _, tc := range cases {
		conn := &fakeConn{driver: tc.driver}
		repo := NewDatabaseRepository(conn, "migrations")

		if _, err := repo.RepositoryExists(context.Background()); err != nil {
			t.Fatalf("%s: RepositoryExists: %v", tc.driver, err)
		}

		if !strings.Contains(conn.selectQuery, tc.mustContain) {
			t.Errorf("%s: expected query to contain %q, got %q", tc.driver, tc.mustContain, conn.selectQuery)
		}
	}
}

func TestCreateRepository_ManualOverrideWinsOverDialect(t *testing.T) {
	t.Parallel()

	conn := &fakeConn{driver: "pgsql"}
	repo := NewDatabaseRepository(conn, "migrations")
	repo.SetCreateDDL("CREATE TABLE %s (custom INT)")

	if err := repo.CreateRepository(context.Background()); err != nil {
		t.Fatalf("CreateRepository: %v", err)
	}

	if got := conn.statements[0]; got != "CREATE TABLE migrations (custom INT)" {
		t.Fatalf("expected manual override, got %q", got)
	}
}

func TestCreateRepository_UnknownDriverFallsBackToSQLiteShape(t *testing.T) {
	t.Parallel()

	conn := &fakeConn{driver: "unknown"}
	repo := NewDatabaseRepository(conn, "migrations")

	if err := repo.CreateRepository(context.Background()); err != nil {
		t.Fatalf("CreateRepository: %v", err)
	}

	if !strings.Contains(conn.statements[0], "AUTOINCREMENT") {
		t.Fatalf("expected SQLite-shaped fallback, got %q", conn.statements[0])
	}
}
