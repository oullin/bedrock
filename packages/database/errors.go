package database

import "errors"

var (
	// ErrConnectionNotConfigured is returned when a connection name has no
	// matching configuration entry.
	ErrConnectionNotConfigured = errors.New("database: connection is not configured")
	// ErrDriverNotSupported is returned when the requested database driver
	// is not registered with the manager.
	ErrDriverNotSupported = errors.New("database: driver is not supported")
	// ErrDeadlockDetected is returned when a database deadlock is detected.
	ErrDeadlockDetected = errors.New("database: deadlock detected")
	// ErrLostConnection is returned when the database connection is lost.
	ErrLostConnection = errors.New("database: lost connection")
	// ErrQueryFailed is returned when a query execution fails.
	ErrQueryFailed = errors.New("database: query execution failed")
	// ErrTransactionFailed is returned when a transaction could not be completed.
	ErrTransactionFailed = errors.New("database: transaction failed")
	// ErrRecordNotFound is returned when a query expecting a result finds none.
	ErrRecordNotFound = errors.New("database: record not found")
	// ErrMultipleRecords is returned when a single-result query returns more
	// than one row.
	ErrMultipleRecords = errors.New("database: multiple records found")
	// ErrUniqueConstraint is returned when an INSERT or UPDATE violates a
	// unique constraint.
	ErrUniqueConstraint = errors.New("database: unique constraint violation")
	// ErrForeignKeysNotSupported is returned by schema grammars whose driver
	// does not enforce referential integrity (e.g. ClickHouse).
	ErrForeignKeysNotSupported = errors.New("database: foreign keys are not supported by this driver")
	// ErrInsertGetIdNotSupported is returned by query grammars whose driver
	// has no auto-increment / last-insert-id semantics (e.g. ClickHouse).
	ErrInsertGetIdNotSupported = errors.New("database: InsertGetId is not supported by this driver")
	// ErrUpsertNotSupported is returned by query grammars whose driver does
	// not implement INSERT ... ON CONFLICT-style upserts (e.g. ClickHouse;
	// use ReplacingMergeTree at the table level instead).
	ErrUpsertNotSupported = errors.New("database: upsert is not supported by this driver")
	// ErrInsertOrIgnoreNotSupported is returned by query grammars whose
	// driver has no INSERT IGNORE / ON CONFLICT DO NOTHING equivalent
	// (e.g. ClickHouse).
	ErrInsertOrIgnoreNotSupported = errors.New("database: InsertOrIgnore is not supported by this driver")
)
