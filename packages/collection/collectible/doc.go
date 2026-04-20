// Package collectible provides an ordered map collection with a fluent,
// chainable API.
//
// [Collection] wraps a Go map while preserving insertion order. Unlike
// Go's built-in map, iterating over a Collection always visits entries
// in the order they were added.
//
// Methods that extract keys or values return plain slices, keeping the
// package independent of other collection packages.
package collectible
