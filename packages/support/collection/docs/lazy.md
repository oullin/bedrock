# 🐢 Lazy API Reference

The `lazy` package provides a `Collection[T]` type that leverages Go's iterators (`iter.Seq[T]`) for lazy evaluation. Items are only computed when needed, making it highly efficient for processing large datasets or infinite sequences.

```go
import "github.com/gollin/packages/support/collection/lazy"
```

### 🚀 Constructors

| Function                  | Description                                              |
| ------------------------- | -------------------------------------------------------- |
| `New(source iter.Seq[T])` | Creates a new Collection from an iterator function.      |
| `From(slice)`             | Creates a Collection from a slice.                       |
| `Empty()`                 | Creates an empty Collection.                             |
| `Range(from, to)`         | Creates a collection of sequential integers.             |
| `Times(n, callback)`      | Creates a collection by invoking the callback `n` times. |

### 🔍 Query & Access

| Method / Function            | Description                                                    |
| ---------------------------- | -------------------------------------------------------------- |
| `Iter()`                     | Returns the underlying iterator (`iter.Seq[T]`).               |
| `All()`                      | Eagerly evaluates the collection and returns a slice.          |
| `Eager()`                    | Alias for `All`.                                               |
| `Collect()`                  | Alias for `All`.                                               |
| `Count()`                    | Returns the total number of items (requires eager evaluation). |
| `IsEmpty()`                  | Reports whether the collection is empty.                       |
| `IsNotEmpty()`               | Reports whether the collection is not empty.                   |
| `ContainsOneItem()`          | Reports whether the collection has exactly one item.           |
| `ContainsManyItems()`        | Reports whether the collection has more than one item.         |
| `First(predicates...)`       | Returns the first item matching the optional predicate.        |
| `FirstOrFail(predicates...)` | Returns the first item or an error if not found.               |
| `Last(predicates...)`        | Returns the last item matching the optional predicate.         |
| `Sole(predicates...)`        | Returns the only item matching the predicate (error if > 1).   |
| `Get(index)`                 | Returns the item at the given zero-based index.                |
| `Has(index)`                 | Reports whether the given index exists.                        |
| `HasAny(indices...)`         | Reports whether any of the given indices exist.                |
| `Contains(predicate)`        | Reports whether any item satisfies the predicate.              |
| `Some(predicate)`            | Alias for `Contains`.                                          |
| `DoesntContain(predicate)`   | Reports whether no item satisfies the predicate.               |
| `HasSole(predicate)`         | Reports whether exactly one item satisfies the predicate.      |
| `Search(predicate)`          | Returns the index of the first item satisfying the predicate.  |
| `Before(predicate)`          | Returns the item before the first match.                       |
| `After(predicate)`           | Returns the item after the first match.                        |
| `Every(predicate)`           | Reports whether all items satisfy the predicate.               |

### ♻️ Transformation

| Method / Function             | Description                                          |
| ----------------------------- | ---------------------------------------------------- |
| `Map(col, callback)`          | (Function) Transforms each item using the callback.  |
| `FlatMap(col, callback)`      | (Function) Transforms and flattens results.          |
| `Filter(predicate)`           | Returns items satisfying the predicate.              |
| `Reject(predicate)`           | Returns items NOT satisfying the predicate.          |
| `Unique(col, keyFunc)`        | (Function) Returns unique items by key.              |
| `Pluck(col, valueFunc)`       | (Function) Extracts a value from each item.          |
| `Reduce(col, callback, init)` | (Function) Reduces the collection to a single value. |
| `GroupBy(col, keyFunc)`       | (Function) Groups items by the key.                  |
| `KeyBy(col, keyFunc)`         | (Function) Indexes items by the key.                 |
| `CountBy(col, keyFunc)`       | (Function) Counts occurrences by key.                |

### 🎯 Pipeline & Utilities

| Method / Function               | Description                                                |
| ------------------------------- | ---------------------------------------------------------- |
| `Take(limit)`                   | Limits the number of items.                                |
| `TakeUntil(predicate)`          | Takes items until the predicate returns true.              |
| `TakeWhile(predicate)`          | Takes items while the predicate returns true.              |
| `Skip(count)`                   | Skips the first `count` items.                             |
| `SkipUntil(predicate)`          | Skips items until the predicate returns true.              |
| `SkipWhile(predicate)`          | Skips items while the predicate returns true.              |
| `Slice(offset, lengths...)`     | Returns a subset of the collection.                        |
| `Chunk(size)`                   | Eagerly splits items into groups of the given size.        |
| `ChunkWhile(predicate)`         | Splits items into groups based on consecutive matches.     |
| `Nth(step, offset)`             | Returns every step-th element.                             |
| `Concat(slice)`                 | Appends items from a slice.                                |
| `Pad(size, value)`              | Pads the collection to a given size.                       |
| `Each(callback)`                | Iterates over items.                                       |
| `Tap(callback)`                 | Passes the collection to a callback.                       |
| `TapEach(callback)`             | Calls the callback on each item as it passes through.      |
| `Throttle(delay)`               | Inserts a delay between each yielded item.                 |
| `Remember()`                    | Caches items after the first iteration.                    |
| `Implode(glue)`                 | Joins items into a string.                                 |
| `Join(glue, finalGlue)`         | Joins items with a separator and optional final separator. |
| `When(cond, callback, default)` | Conditional application.                                   |
| `WhenEmpty(callback)`           | Applies if the collection is empty.                        |
| `WhenNotEmpty(callback)`        | Applies if the collection is not empty.                    |
| `Unless(cond, callback)`        | Inverse of `When`.                                         |
| `Dump()`                        | Prints evaluated items for debugging.                      |

---

## Detailed Methods

### `Map`

Transforms each item in the collection. This is a package-level function to allow changing the type of the collection.

```go
numbers := lazy.From([]int{1, 2, 3})
doubled := lazy.Map(numbers, func(v int, i int) int {
    return v * 2
})
```

### `Take` & `Skip`

Lazy collections can be sliced without evaluating the entire source.

```go
// Only takes the first 5 even numbers
evens := lazy.Range(1, 100).
    Filter(func(v int, _ int) bool { return v%2 == 0 }).
    Take(5)
```

### `Remember`

By default, lazy collections re-run their logic every time they are iterated. `Remember` caches the results in memory after the first pass.

```go
cached := lazy.New(expensiveSource).Remember()
cached.All() // Computes
cached.All() // Uses cache
```

### `Throttle`

`Throttle` is useful for rate-limiting operations like API calls or log processing.

```go
// Yields one item every 500ms
lazy.From(items).Throttle(500 * time.Millisecond).Each(func(v T, _ int) bool {
    process(v)
    return true
})
```

👉 [**Back to Overview**](overview.md)
