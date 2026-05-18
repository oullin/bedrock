package migrations

// dialect holds per-driver overrides for the migrations-table CREATE and
// existence-check statements. Drivers whose default DDL is incompatible
// register an entry here so DatabaseRepository can pick it up automatically
// from the connection's driver name.
type dialect struct {
	// createDDL is a fmt.Sprintf format with a single %s placeholder for
	// the table name.
	createDDL string
	// existsSQL is a parameterized query receiving the table name as its
	// single positional binding.
	existsSQL string
}

var dialects = map[string]dialect{
	"clickhouse": {
		createDDL: `CREATE TABLE IF NOT EXISTS %s (
			id Int64,
			migration String,
			batch Int32
		) ENGINE = MergeTree
		ORDER BY (batch, migration)`,
		existsSQL: "SELECT name FROM system.tables WHERE database = currentDatabase() AND name = ?",
	},
}

// dialectFor returns the registered dialect for a given driver name, falling
// back to the zero value when none is registered. Callers should treat empty
// fields as "use the package default".
func dialectFor(driverName string) dialect {
	return dialects[driverName]
}
