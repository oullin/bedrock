package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/bedrock/packages/database"
)

// Connector creates SQLite database connections.
type Connector struct{}

// Connect opens a SQLite connection using the given configuration.
func (c *Connector) Connect(config database.ConnectionConfig) (*database.Connection, error) {
	dsn := config.Database

	if dsn == "" {
		dsn = ":memory:"
	}

	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, fmt.Errorf("sqlite: failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("sqlite: failed to ping: %w", err)
	}

	// Enable WAL mode and foreign keys by default.
	_, _ = db.Exec("PRAGMA journal_mode=WAL")
	_, _ = db.Exec("PRAGMA foreign_keys=ON")

	conn := database.NewConnection(db, "", config.Database, config.Prefix, map[string]any{
		"driver":   "sqlite",
		"database": config.Database,
	})
	conn.SetDriverName("sqlite")

	return conn, nil
}

// NewConnectorFactory returns a ConnectorFactory for the database Manager.
func NewConnectorFactory() database.ConnectorFactory {
	c := &Connector{}

	return func(config database.ConnectionConfig) (*database.Connection, error) {
		return c.Connect(config)
	}
}
