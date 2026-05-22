// Package database provides a database abstraction layer with query builder,
// schema management, Orm-style ORM, and migration support. It is the Go
// Ref: @bedrock/code-0201
// maintaining 100% function parity.
//
// The package supports MySQL, PostgreSQL, and SQLite through driver-specific
// grammars, processors, and connectors. Connections are managed by a
// DatabaseManager that follows the Manager/Driver pattern used throughout
// the Bedrock framework.
//
// Key components:
//   - Connection: executes raw SQL, manages transactions, dispatches events
//   - Query Builder: fluent interface for building SELECT, INSERT, UPDATE, DELETE
//   - Schema Builder: Blueprint-based DDL for creating and modifying tables
//   - Orm ORM: Model[T] with relationships, scopes, events, and casting
//   - Migrations: version-controlled database schema changes
package database
