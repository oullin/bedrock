// Package clickhouse provides a ClickHouse driver for the database package.
// It connects through the standard database/sql interface backed by
// github.com/ClickHouse/clickhouse-go/v2, which must be blank-imported by the
// application:
//
//	import _ "github.com/ClickHouse/clickhouse-go/v2"
//
// The driver supports MergeTree-family engines. Apps declare the engine and
// sort key on the schema Blueprint via SetEngine and SetOrderBy. Per-row
// UPDATE and DELETE are emitted as ALTER TABLE ... UPDATE / ... DELETE
// mutations, which the server processes asynchronously. Foreign keys, last
// insert IDs, and TRUNCATE with cascade are not supported and surface as
// explicit errors rather than silent no-ops.
package clickhouse

// DriverName is the ClickHouse driver identifier used by Manager and grammars.
const DriverName = "clickhouse"

// Default network port for the ClickHouse native TCP protocol.
const defaultPort = 9000
