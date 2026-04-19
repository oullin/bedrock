package query

// Grammar compiles a Builder into SQL for a specific database driver.
type Grammar interface {
	// CompileSelect compiles a SELECT query.
	CompileSelect(b *Builder) string
	// CompileExists wraps a SELECT in an EXISTS check.
	CompileExists(b *Builder) string
	// CompileInsert compiles an INSERT statement.
	CompileInsert(b *Builder, values []map[string]any) string
	// CompileInsertOrIgnore compiles an INSERT IGNORE statement.
	CompileInsertOrIgnore(b *Builder, values []map[string]any) string
	// CompileInsertGetId compiles an INSERT returning the new ID.
	CompileInsertGetId(b *Builder, values map[string]any, sequence string) string
	// CompileInsertUsing compiles an INSERT ... SELECT statement.
	CompileInsertUsing(b *Builder, columns []string, sql string) string
	// CompileUpdate compiles an UPDATE statement.
	CompileUpdate(b *Builder, values map[string]any) string
	// CompileUpsert compiles an UPSERT (INSERT ON CONFLICT UPDATE) statement.
	CompileUpsert(b *Builder, values []map[string]any, uniqueBy []string, update []string) string
	// CompileDelete compiles a DELETE statement.
	CompileDelete(b *Builder) string
	// CompileTruncate compiles a TRUNCATE statement.
	CompileTruncate(b *Builder) map[string]string
	// CompileRandom compiles a random ordering expression.
	CompileRandom(seed string) string
	// Wrap wraps a value in keyword identifiers.
	Wrap(value string) string
	// WrapTable wraps a table name.
	WrapTable(table string) string
	// Columnize converts a slice of column names to a comma-separated string.
	Columnize(columns []string) string
	// Parameterize converts a slice of values to a comma-separated parameter string.
	Parameterize(values []any) string
	// Parameter converts a single value to a parameter placeholder.
	Parameter(value any) string
	// GetTablePrefix returns the table prefix.
	GetTablePrefix() string
	// SetTablePrefix sets the table prefix.
	SetTablePrefix(prefix string)
	// IsExpression reports whether the value is an Expression.
	IsExpression(value any) bool
	// GetValue extracts the value from an Expression.
	GetValue(expression any) string
}

// Processor processes query results returned by the database.
type Processor interface {
	// ProcessSelect processes the results of a select query.
	ProcessSelect(b *Builder, results []map[string]any) []map[string]any
	// ProcessInsertGetId processes an auto-incrementing ID after insert.
	ProcessInsertGetId(b *Builder, sql string, values []any, sequence string) (int64, error)
	// ProcessColumnListing processes the results of a column listing query.
	ProcessColumnListing(results []map[string]any) []string
}
