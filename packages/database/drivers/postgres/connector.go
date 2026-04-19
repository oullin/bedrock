package postgres

import (
	"database/sql"
	"fmt"

	"github.com/bedrock/packages/database"
)

// Connector creates PostgreSQL database connections.
type Connector struct{}

// Connect opens a PostgreSQL connection using the given configuration.
func (c *Connector) Connect(config database.ConnectionConfig) (*database.Connection, error) {
	dsn := buildDSN(config)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("postgres: failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("postgres: failed to ping: %w", err)
	}

	conn := database.NewConnection(db, "", config.Database, config.Prefix, map[string]any{
		"driver":   "pgsql",
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
		"username": config.Username,
		"schema":   config.Schema,
	})
	conn.SetDriverName("pgsql")

	return conn, nil
}

func buildDSN(config database.ConnectionConfig) string {
	port := config.Port

	if port == 0 {
		port = 5432
	}

	sslmode := config.SSLMode

	if sslmode == "" {
		sslmode = "disable"
	}

	schema := config.Schema

	if schema == "" {
		schema = "public"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		config.Host, port,
		config.Username, config.Password,
		config.Database, sslmode, schema,
	)
}

// NewConnectorFactory returns a ConnectorFactory for the database Manager.
func NewConnectorFactory() database.ConnectorFactory {
	c := &Connector{}

	return func(config database.ConnectionConfig) (*database.Connection, error) {
		return c.Connect(config)
	}
}
