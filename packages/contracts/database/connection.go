package database

import "context"

// Connection is the primary contract for a database connection.
// Ref: @bedrock/code-0203
type Connection interface {
	// Table begins a fluent query builder for the given table.
	Table(ctx context.Context, table string, as ...string) any
	// Raw creates a raw SQL Expression that will not be quoted.
	Raw(value string) Expression
	// SelectOne runs a SELECT query and returns a single row.
	SelectOne(ctx context.Context, query string, bindings ...any) (map[string]any, error)
	// Select runs a SELECT query and returns all matching rows.
	Select(ctx context.Context, query string, bindings ...any) ([]map[string]any, error)
	// Insert executes an INSERT statement.
	Insert(ctx context.Context, query string, bindings ...any) (bool, error)
	// Update executes an UPDATE statement and returns affected rows.
	Update(ctx context.Context, query string, bindings ...any) (int64, error)
	// Delete executes a DELETE statement and returns affected rows.
	Delete(ctx context.Context, query string, bindings ...any) (int64, error)
	// Statement executes a raw SQL statement.
	Statement(ctx context.Context, query string, bindings ...any) (bool, error)
	// AffectingStatement executes a statement and returns affected rows.
	AffectingStatement(ctx context.Context, query string, bindings ...any) (int64, error)
	// Unprepared runs a raw query without parameter binding.
	Unprepared(ctx context.Context, query string) (bool, error)
	// PrepareBindings prepares the query bindings for execution.
	PrepareBindings(bindings []any) []any
	// Transaction runs a closure inside a database transaction with retry support.
	Transaction(ctx context.Context, fn func(Connection) error, attempts ...int) error
	// BeginTransaction starts a new database transaction.
	BeginTransaction(ctx context.Context) error
	// Commit commits the active database transaction.
	Commit(ctx context.Context) error
	// Rollback rolls back the active database transaction.
	Rollback(ctx context.Context, toLevel ...int) error
	// TransactionLevel returns the current transaction nesting depth.
	TransactionLevel() int
	// AfterCommit registers a callback to run after the current transaction commits.
	AfterCommit(fn func())
	// GetTablePrefix returns the table prefix for this connection.
	GetTablePrefix() string
	// GetDatabaseName returns the name of the connected database.
	GetDatabaseName() string
	// GetDriverName returns the driver name (mysql, pgsql, sqlite).
	GetDriverName() string
	// GetName returns the connection name.
	GetName() string
	// GetConfig returns a configuration value.
	GetConfig(key string) any
}

// ConnectionResolver resolves named database connections.
// Ref: @bedrock/code-0204
type ConnectionResolver interface {
	// Connection returns the connection with the given name, or the default
	// connection if no name is provided.
	Connection(ctx context.Context, name ...string) (Connection, error)
	// GetDefaultConnection returns the default connection name.
	GetDefaultConnection() string
	// SetDefaultConnection sets the default connection name.
	SetDefaultConnection(name string)
}

// Expression wraps a raw SQL string so the query builder will not
// parameterise or quote it.
type Expression interface {
	// GetValue returns the raw SQL string.
	GetValue() string
}
