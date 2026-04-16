package mariadb

import (
	"database/sql"
	"fmt"

	"github.com/bedrock/packages/database"
)

// Connector creates MariaDB database connections.
type Connector struct{}

// Connect opens a MariaDB connection using the given configuration.
// MariaDB uses the same wire protocol as MySQL but the grammar diverges
// on JSON operations, RETURNING support, and UUID functions.
func (c *Connector) Connect(config database.ConnectionConfig) (*database.Connection, error) {
	dsn := buildDSN(config)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mariadb: failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("mariadb: failed to ping: %w", err)
	}

	conn := database.NewConnection(db, "", config.Database, config.Prefix, map[string]any{
		"driver":   "mariadb",
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
		"username": config.Username,
		"charset":  config.Charset,
	})
	conn.SetDriverName("mariadb")

	return conn, nil
}

func buildDSN(config database.ConnectionConfig) string {
	charset := config.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	port := config.Port
	if port == 0 {
		port = 3306
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		config.Username, config.Password,
		config.Host, port,
		config.Database, charset,
	)
}

// NewConnectorFactory returns a ConnectorFactory for the database Manager.
func NewConnectorFactory() database.ConnectorFactory {
	c := &Connector{}
	return func(config database.ConnectionConfig) (*database.Connection, error) {
		return c.Connect(config)
	}
}
