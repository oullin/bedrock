# 💎 Collection[T] API Reference

`import "github.com/bedrock/packages/collection/collection"`

`Collection[T]` is a generic wrapper around a Go slice that provides a fluent, chainable API for filtering, sorting, transforming, and aggregating data. It is the core type of the `collection` package.

> **⚠️ Note on Generics:** Go generics do not allow methods to introduce new type parameters. Functions like `Map`, `Reduce`, `Pluck`, and `GroupBy` are implemented as package-level functions (e.g., `collection.Map(c, fn)`), not methods.

---

## 🛠 Available Functions

### Constructors & Query
| Function | Purpose |
|:---|:---|
| [**New**](#new) | Creates a new collection from variadic arguments. |
| [**Collect**](#collect) | Creates a new collection from an existing slice. |
| [**Empty**](#empty) | Creates an empty collection of the given type. |
| [**Wrap**](#wrap) | Wraps any value into a collection. |
| [**Unwrap**](#unwrap) | Extracts the underlying slice from a collection. |
| [**Times**](#times) | Creates a collection by invoking a callback N times. |
| [**Range**](#range) | Creates a collection of consecutive integers. |
| [**All**](#all) | Returns the underlying slice. |
| [**Count**](#count) | Returns the number of items in the collection. |
| [**IsEmpty**](#isempty) | Reports whether the collection is empty. |
| [**IsNotEmpty**](#isnotempty) | Reports whether the collection is not empty. |
| [**ContainsOneItem**](#containsoneitem) | Reports if there is exactly one item. |
| [**ContainsManyItems**](#containsmanyitems) | Reports if there are multiple items. |
| [**Has**](#has) | Reports whether a given index exists. |

### Retrieval & Search
| Function | Purpose |
|:---|:---|
| [**First**](#first) | Returns the first item matching a predicate. |
| [**FirstOrFail**](#firstorfail) | Returns the first item or an error if not found. |
| [**Last**](#last) | Returns the last item matching a predicate. |
| [**Sole**](#sole) | Returns the only item matching a predicate (or error). |
| [**Get**](#get) | Returns the item at a given index (with optional default). |
| [**Search**](#search) | Returns the index of the first matching item. |
| [**Contains**](#contains) | Reports if any item satisfies a predicate. |
| [**Before**](#before) | Returns the item before the first match. |
| [**After**](#after) | Returns the item after the first match. |

### Mutation & Iteration
| Function | Purpose |
|:---|:---|
| [**Push**](#push) / [**Add**](#add) | Appends items to the collection (mutates). |
| [**Prepend**](#prepend) | Adds an item to the beginning (mutates). |
| [**Pop**](#pop) | Removes and returns the last item. |
| [**Shift**](#shift) | Removes and returns the first item. |
| [**Put**](#put) | Sets the item at a given index. |
| [**Pull**](#pull) | Removes and returns an item by index. |
| [**Forget**](#forget) | Removes an item by index (mutates). |
| [**Each**](#each) | Iterates over items with a callback. |
| [**Tap**](#tap) | Passes the collection to a callback for side effects. |
| [**Iter**](#iter) | Returns a Go 1.23+ iterator (`iter.Seq[T]`). |

### Transformation & Sorting
| Function | Purpose |
|:---|:---|
| [**Map**](#map) | Transforms items into a new type (Package-level). |
| [**FlatMap**](#flatmap) | Maps and flattens results (Package-level). |
| [**Reduce**](#reduce) | Accumulates a single result (Package-level). |
| [**Filter**](#filter) / [**Where**](#where) | Returns items matching a predicate. |
| [**Unique**](#unique) | Returns distinct items (Package-level). |
| [**Sort**](#sort) / [**SortBy**](#sortby) | Returns a new sorted collection. |
| [**Reverse**](#reverse) | Returns a reversed collection. |
| [**Shuffle**](#shuffle) | Returns a randomized collection. |

### Partitioning & Slicing
| Function | Purpose |
|:---|:---|
| [**Chunk**](#chunk) | Breaks the collection into multiple slices of N size. |
| [**Split**](#split) | Splits items into a fixed number of groups. |
| [**Slice**](#slice) | Extracts a portion of the collection. |
| [**Take**](#take) | Returns N items from the front or back. |
| [**Skip**](#skip) | Skips N items and returns the rest. |
| [**ForPage**](#forpage) | Returns a subset for a specific page. |

### Combining & Aggregation
| Function | Purpose |
|:---|:---|
| [**Concat**](#concat) / [**Merge**](#merge) | Joins two datasets. |
| [**Zip**](#zip) | Merges multiple slices element-by-element. |
| [**Sum**](#sum) / [**Avg**](#avg) | Computes totals or averages (Package-level). |
| [**Min**](#min) / [**Max**](#max) | Finds minimum or maximum values (Package-level). |
| [**GroupBy**](#groupby) | Groups items into a map of collections. |

---

## 🚀 Constructors

### New
```go
func New[T any](items ...T) *Collection[T]
```
Creates a new collection from variadic arguments.

### Collect
```go
func Collect[T any](items []T) *Collection[T]
```
Creates a new collection from an existing slice.

### Empty
```go
func Empty[T any]() *Collection[T]
```
Creates an empty collection of the given type.

---

## 🔍 Query

### All
```go
func (c *Collection[T]) All() []T
```
Returns the underlying slice.

### Count
```go
func (c *Collection[T]) Count() int
```
Returns the number of items in the collection.

### IsEmpty / IsNotEmpty
```go
func (c *Collection[T]) IsEmpty() bool
func (c *Collection[T]) IsNotEmpty() bool
```

---

## 🎯 Retrieval

### First
```go
func (c *Collection[T]) First(predicates ...func(T, int) bool) (T, bool)
```
Returns the first matching element.

### Get
```go
func (c *Collection[T]) Get(index int, defaults ...T) (T, bool)
```
Returns the item at the given index with an optional default value.

---

## 🧪 Search

### Contains
```go
func (c *Collection[T]) Contains(predicate func(T, int) bool) bool
```
Reports whether any item satisfies the predicate.

---

## 🛠 Mutation

### Push / Add
```go
func (c *Collection[T]) Push(values ...T) *Collection[T]
```
Appends items to the collection. Mutates in place.

### Forget
```go
func (c *Collection[T]) Forget(index int) *Collection[T]
```
Removes an item by index. Mutates in place.

---

## ♻️ Iteration

### Each
```go
func (c *Collection[T]) Each(callback func(T, int) bool) *Collection[T]
```
Iterates over items. Return `false` to stop early.

### Iter
```go
func (c *Collection[T]) Iter() iter.Seq[T]
```
Returns a Go 1.23+ iterator compatible with `for range`.

---

## 💎 Transformation (Package-Level)

### Map
```go
func Map[T any, R any](c *Collection[T], callback func(T, int) R) *Collection[R]
```
Transforms each item into a new type.

### Reduce
```go
func Reduce[T any, R any](c *Collection[T], callback func(R, T, int) R, initial R) R
```
Reduces the collection to a single value.

---

## 🧹 Filtering

### Filter / Where
```go
func (c *Collection[T]) Filter(callback func(T, int) bool) *Collection[T]
```
Returns items matching the predicate.

### Unique
```go
func Unique[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T]
```
Returns a new collection with distinct keys.

---

## 📊 Sorting

### Sort
```go
func (c *Collection[T]) Sort(less func(a, b T) bool) *Collection[T]
```
Returns a new collection sorted by the provided function.

---

## 🍰 Partitioning

### Chunk
```go
func (c *Collection[T]) Chunk(size int) [][]T
```
Breaks the collection into multiple slices of the given size.

---

## ✂️ Slicing

### Take
```go
func (c *Collection[T]) Take(limit int) *Collection[T]
```
Returns N items from the front (positive) or back (negative).

---

## ➕ Combining

### Merge
```go
func (c *Collection[T]) Merge(items []T) *Collection[T]
```
Joins the collection with another slice.

---

## 📈 Aggregation (Package-Level)

### Sum / Avg
```go
func Sum[T Numeric](c *Collection[T]) T
func Avg[T Numeric](c *Collection[T]) float64
```

### Min / Max
```go
func Min[T cmp.Ordered](c *Collection[T]) (T, bool)
func Max[T cmp.Ordered](c *Collection[T]) (T, bool)
```

---

## 📦 Serialization

### ToJSON / ToPrettyJSON
```go
func (c *Collection[T]) ToJSON() ([]byte, error)
func (c *Collection[T]) ToPrettyJSON() ([]byte, error)
```

---

👉 [**Back to Overview**](overview.md)
