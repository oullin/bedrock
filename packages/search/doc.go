// Package search provides full-text search with pluggable engine backends.
// It is the Go port of Upstream's Search package, adapted to Go idioms while
// maintaining 100% function parity.
//
// Search integrates with Orm models through the Searchable interface and
// SearchableMixin, providing automatic index synchronisation on model
// create, update, and delete operations.
//
// Supported search engines:
//   - Database: full-text search via SQL (MATCH/AGAINST, tsvector, LIKE)
//   - Collection: in-memory filtering without external dependencies
//   - Meilisearch: HTTP-based search via the Meilisearch Go client
//   - Algolia: search via the Algolia Go client
//   - Typesense: search via the Typesense Go client
//   - Null: no-op engine for testing and development
//
// Key components:
//   - Engine: abstract interface for search backends
//   - Builder: fluent API for constructing search queries
//   - EngineManager: Manager/Driver pattern for pluggable engines
//   - ModelObserver: automatic model-to-index synchronisation
//   - SearchableMixin: embeddable struct providing default Searchable behaviour
package search
