package migrations

// dialect holds per-driver overrides for the migrations-table CREATE and
// existence-check statements. Each registered driver supplies the column
// types and metadata query that match its system catalogs, so apps get
// correct DDL out of the box without manual SetCreateDDL/SetExistsSQL calls.
type dialect struct {
	// createDDL is a fmt.Sprintf format with a single %s placeholder for
	// the table name.
	createDDL string
	// existsSQL is a parameterized query that receives the table name as
	// its single positional binding. Drivers using $N placeholders (postgres)
	// must spell them explicitly here.
	existsSQL string
}

// mysqlFamily is shared by mysql and mariadb: identical column types,
// storage engine, and information_schema lookup.
var mysqlFamily = dialect{
	createDDL: `CREATE TABLE IF NOT EXISTS %s (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		migration VARCHAR(255) NOT NULL,
		batch INT NOT NULL
	) ENGINE=InnoDB`,
	existsSQL: "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?",
}

var dialects = map[string]dialect{
	"pgsql": {
		createDDL: `CREATE TABLE IF NOT EXISTS %s (
			id BIGSERIAL PRIMARY KEY,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL
		)`,
		existsSQL: "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = current_schema() AND tablename = $1",
	},
	"mysql":   mysqlFamily,
	"mariadb": mysqlFamily,
	"sqlite": {
		createDDL: `CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration VARCHAR(255) NOT NULL,
			batch INTEGER NOT NULL
		)`,
		existsSQL: "SELECT name FROM sqlite_master WHERE type='table' AND name=?",
	},
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
