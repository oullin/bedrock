# collection

<!-- upstream-docs: collections.md#collections -->
<!-- upstream-docs: collections.md#lazy-collections -->
<!-- upstream-docs: collections.md#available-methods -->

Fluent, type-safe collection helpers for slices, ordered maps, lazy sequences,
and one-off array or key-value operations.

**Module:** `github.com/bedrock/packages/collection`

```bash
go get github.com/bedrock/packages/collection@latest
```

## Packages

| Package       | Purpose                                                     |
| ------------- | ----------------------------------------------------------- |
| `collection`  | Fluent `Collection[T]` wrapper for slices                   |
| `lazy`        | Deferred `iter.Seq[T]` pipelines for large or streamed data |
| `collectible` | Ordered key-value collection with a fluent API              |
| `arr`         | Standalone generic slice helpers                            |
| `kv`          | `map[string]any` helpers with dot-notation access           |

## Collection

```go
import "github.com/bedrock/packages/collection/collection"

numbers := collection.Collect([]int{1, 2, 3, 4, 5, 6})

even := numbers.
	Filter(func(n int, _ int) bool { return n%2 == 0 }).
	Take(2)

even.All() // []int{2, 4}
```

Use top-level functions when a transformation changes the element type:

```go
labels := collection.Map(even, func(n int, _ int) string {
	return fmt.Sprintf("number:%d", n)
})
```

## Lazy Sequences

```go
import "github.com/bedrock/packages/collection/lazy"

values := lazy.Range(1, 1_000_000).
	Filter(func(n int, _ int) bool { return n%2 == 0 }).
	Take(3).
	All()
```

## Maps and Utilities

```go
import (
	"github.com/bedrock/packages/collection/arr"
	"github.com/bedrock/packages/collection/collectible"
	"github.com/bedrock/packages/collection/kv"
)

arr.Filter([]int{1, 2, 3}, func(n int, _ int) bool { return n > 1 })

data := map[string]any{"user": map[string]any{"name": "Ada"}}
kv.Get(data, "user.name")

scores := collectible.New(map[string]int{"a": 1, "b": 2})
scores.Keys()
```

## Notes

The existing `support` package still exposes its current `Arr*`, `Fluent`,
`Optional[T]`, `MessageBag`, sleep, retry, and timebox helpers. This package is
added as a standalone collection surface; cleanup of older overlapping helpers
is intentionally deferred.
