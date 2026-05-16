package clickhouse

import (
	"context"
	"errors"

	"github.com/bedrock/packages/database"
)

// Raw exposes the underlying clickhouse-go driver connection so callers can
// invoke native APIs (columnar bulk inserts, async inserts, prepared batches)
// that are not reachable through database/sql. The fn argument receives the
// driver's clickhouse.Conn — assert it inside fn to the concrete type from
// github.com/ClickHouse/clickhouse-go/v2.
func Raw(ctx context.Context, conn *database.Connection, fn func(driverConn any) error) error {
	if conn == nil {
		return errors.New("clickhouse: nil connection")
	}

	db := conn.DB()

	if db == nil {
		return errors.New("clickhouse: connection has no underlying *sql.DB")
	}

	sqlConn, err := db.Conn(ctx)

	if err != nil {
		return err
	}

	defer sqlConn.Close()

	return sqlConn.Raw(fn)
}
