package clickhouse

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/bedrock/packages/database"
)

// Connector creates ClickHouse database connections.
type Connector struct{}

// Connect opens a ClickHouse connection using the given configuration.
// The clickhouse-go/v2 driver must be blank-imported by the application:
//
//	import _ "github.com/ClickHouse/clickhouse-go/v2"
func (c *Connector) Connect(config database.ConnectionConfig) (*database.Connection, error) {
	dsn := buildDSN(config)

	db, err := sql.Open(DriverName, dsn)

	if err != nil {
		return nil, fmt.Errorf("clickhouse: failed to open connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("clickhouse: failed to ping: %w", err)
	}

	conn := database.NewConnection(db, "", config.Database, config.Prefix, map[string]any{
		"driver":   DriverName,
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
		"username": config.Username,
	})
	conn.SetDriverName(DriverName)

	return conn, nil
}

// NewConnectorFactory returns a ConnectorFactory for the database Manager.
func NewConnectorFactory() database.ConnectorFactory {
	c := &Connector{}

	return func(config database.ConnectionConfig) (*database.Connection, error) {
		return c.Connect(config)
	}
}

// buildDSN renders a clickhouse-go/v2-compatible URL from the connection
// config. It defaults the port to 9000 (native TCP) and forwards extra
// options as query parameters.
func buildDSN(config database.ConnectionConfig) string {
	port := config.Port

	if port == 0 {
		port = defaultPort
	}

	scheme := DriverName

	if proto, ok := config.Options["protocol"].(string); ok && strings.EqualFold(proto, "http") {
		scheme = "http"
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", config.Host, port),
		Path:   "/" + config.Database,
	}

	if config.Username != "" {
		u.User = url.UserPassword(config.Username, config.Password)
	}

	q := u.Query()

	for k, v := range config.Options {
		if k == "protocol" {
			continue
		}

		q.Set(k, fmt.Sprintf("%v", v))
	}

	u.RawQuery = q.Encode()

	return u.String()
}
