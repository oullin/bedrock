package events

// SchemaDumped is dispatched after a database schema is dumped.
type SchemaDumped struct {
	ConnectionName string
	Path           string
}

// SchemaLoaded is dispatched after a schema dump is loaded.
type SchemaLoaded struct {
	ConnectionName string
	Path           string
}

// StatementPrepared is dispatched when a SQL statement is prepared.
type StatementPrepared struct {
	ConnectionName string
	Statement      string
}
