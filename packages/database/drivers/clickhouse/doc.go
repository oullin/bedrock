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
//
// For columnar bulk inserts and async inserts, use Raw to drop down to the
// underlying clickhouse-go native connection without losing the rest of the
// shared Connection wiring (events, query log, transactions, etc.).
package clickhouse
