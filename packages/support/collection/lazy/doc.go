// Package lazy provides lazily evaluated generic sequences backed by [iter.Seq].
//
// [Collection] it wraps an iterator and provides a fluent API for filtering,
// transforming, and consuming data on demand. Nothing is computed until the
// collection is materialised, making it ideal for large datasets, pipelines
// with early termination, and infinite sequences.
//
// To create a lazy collection from a slice:
//
//	lc := lazy.From([]int{1, 2, 3})
//
// To convert back to a slice:
//
//	items := lc.Eager()
package lazy
