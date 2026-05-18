package clickhouse

import "errors"

// ErrForeignKeysNotSupported is returned when a schema operation requests a
// foreign key. ClickHouse does not enforce referential integrity; surfacing
// an error prevents silent migration failures.
var ErrForeignKeysNotSupported = errors.New("clickhouse: foreign keys are not supported")

// ErrInsertGetIdNotSupported is returned when InsertGetId is invoked against
// a ClickHouse connection. ClickHouse has no auto-increment / last-insert-id
// semantics; callers should assign IDs explicitly.
var ErrInsertGetIdNotSupported = errors.New("clickhouse: InsertGetId is not supported (no auto-increment)")

// ErrUpsertNotSupported is returned when an UPSERT-style statement is
// requested. ClickHouse offers ReplacingMergeTree for deduplication instead.
var ErrUpsertNotSupported = errors.New("clickhouse: upsert is not supported (use ReplacingMergeTree)")
