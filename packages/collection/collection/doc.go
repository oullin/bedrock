// Package collection provides a fluent, generic wrapper for working with slices of data.
//
// The core type is:
//
//   - [Collection] — a generic wrapper around a slice with a rich set of chainable methods
//     for filtering, sorting, transforming, and aggregating data.
//
// # Related packages
//
//   - [github.com/bedrock/packages/collection/lazy] — lazily evaluated sequences backed by iter.Seq
//   - [github.com/bedrock/packages/collection/collectible] — ordered map with fluent key-value API
//   - [github.com/bedrock/packages/collection/support] — shared types (Pair, Numeric) and errors
//   - [github.com/bedrock/packages/collection/arr] — generic slice helpers
//   - [github.com/bedrock/packages/collection/kv] — map helpers with dot-notation support
package collection
