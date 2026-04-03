// Package support provides shared types and error definitions used across
// the collection family of packages.
//
// Exported types:
//
//   - [Pair] — a generic key-value pair.
//   - [Numeric] — a type constraint covering all built-in numeric types.
//   - [ItemNotFoundError] — returned when a lookup finds no matching item.
//   - [MultipleItemsFoundError] — returned when a single item is expected but several matches.
package support
