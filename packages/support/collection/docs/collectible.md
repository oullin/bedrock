# 💎 Collectible API Reference

The `collectible` package provides a generic `Collection[K, V]` type that wraps a Go map. It maintains the insertion order of keys and offers a fluent API for common map operations, transformations, and queries.

```go
import "github.com/gollin/packages/support/collection/collectible"
```

### 🚀 Constructors

| Function | Description |
| --- | --- |
| `New(items map[K]V)` | Creates a new Collection from a map. |
| `FromPairs(pairs ...Pair[K, V])` | Creates a Collection from key-value pairs. |

### 🔍 Query & Access

| Method | Description |
| --- | --- |
| `All()` | Returns a shallow copy of the underlying map. |
| `Keys()` | Returns all keys as a slice, preserving insertion order. |
| `Values()` | Returns all values as a slice, preserving insertion order. |
| `Count()` | Returns the number of items. |
| `IsEmpty()` | Reports whether the collection is empty. |
| `IsNotEmpty()` | Reports whether the collection is not empty. |
| `ContainsOneItem()` | Reports whether the collection has exactly one item. |
| `ContainsManyItems()` | Reports whether the collection has more than one item. |
| `Iter()` | Returns a standard Go iterator (`iter.Seq2[K, V]`). |
| `Get(key, defaults...)` | Returns the value for a key, with an optional default. |
| `Has(key)` | Reports whether a key exists. |
| `HasAny(keys...)` | Reports whether any of the given keys exist. |
| `Contains(predicate)` | Reports whether any item satisfies the predicate. |
| `Some(predicate)` | Alias for `Contains`. |
| `DoesntContain(predicate)` | Reports whether no item satisfies the predicate. |
| `HasSole(predicate)` | Reports whether exactly one item satisfies the predicate. |
| `Search(predicate)` | Returns the key of the first item satisfying the predicate. |
| `First(predicates...)` | Returns the first item matching the optional predicate. |
| `Last(predicates...)` | Returns the last item matching the optional predicate. |
| `Every(predicate)` | Reports whether all items satisfy the predicate. |

### 🛠️ Mutation & Transformation

| Method / Function | Description |
| --- | --- |
| `Put(key, value)` | Sets a key-value pair in the collection. |
| `GetOrPut(key, value)` | Returns the value for a key or stores the default. |
| `Pull(key)` | Removes a key and returns its value. |
| `Forget(keys...)` | Removes one or more items by key. |
| `Merge(map)` | Returns a new collection with merged items (overwrites). |
| `Union(map)` | Returns a new collection with merged items (no overwrite). |
| `Replace(map)` | Alias for `Merge`. |
| `DiffKeys(map)` | Returns items whose keys are not in the given map. |
| `DiffKeysUsing(map, equals)` | `DiffKeys` with a custom comparison function. |
| `IntersectByKeys(map)` | Returns items whose keys are in the given map. |
| `Only(keys...)` | Returns a collection containing only the specified keys. |
| `Except(keys...)` | Returns a collection excluding the specified keys. |
| `Filter(predicate)` | Returns items satisfying the predicate. |
| `Reject(predicate)` | Returns items NOT satisfying the predicate. |
| `Partition(predicate)` | Splits the collection into two based on the predicate. |
| `MapValues(col, callback)` | (Function) Transforms values while keeping keys. |
| `Flip(col)` | (Function) Swaps keys and values. |

### 🎯 Pipeline & Utilities

| Method / Function | Description |
| --- | --- |
| `Each(callback)` | Iterates over items. Return `false` to stop. |
| `Tap(callback)` | Executes a callback with the collection and returns it. |
| `When(cond, callback, default)` | Applies a callback if the condition is true. |
| `Unless(cond, callback, default)` | Applies a callback unless the condition is true. |
| `Implode(glue)` | Concatenates values into a string. |
| `Join(glue, finalGlue)` | Joins values with a separator and optional final separator. |
| `ToJSON()` | Serializes the collection to JSON. |
| `Copy()` | Returns a shallow copy of the collection. |
| `ToPairs()` | Converts the collection to a slice of `Pair` values. |
| `Dump()` | Prints the items for debugging. |
| `DD()` | Prints the items and terminates the program. |

### ⚖️ Sorting

| Method / Function | Description |
| --- | --- |
| `SortKeys(col)` | (Function) Sorts string keys in ascending order. |
| `SortKeysDesc(col)` | (Function) Sorts string keys in descending order. |
| `SortKeysUsing(less)` | Sorts keys using a custom comparison function. |

---

## Detailed Methods

### `Filter`
Filters the collection using the given callback. The callback receives the value and the key.

```go
filtered := col.Filter(func(v int, k string) bool {
    return v > 10
})
```

### `MapValues`
Transforms each value in the collection. Note that this is a package-level function to allow changing the value type.

```go
names := collectible.New(map[int]string{1: "juan", 2: "pedro"})
lengths := collectible.MapValues(names, func(v string, k int) int {
    return len(v)
})
```

### `Merge` & `Union`
`Merge` overwrites existing keys, while `Union` only adds keys that are not already present.

```go
col := collectible.New(map[string]int{"a": 1})
merged := col.Merge(map[string]int{"a": 2, "b": 3}) // a: 2, b: 3
union := col.Union(map[string]int{"a": 3, "c": 4}) // a: 1, c: 4
```

### `Partition`
Splits the collection into two collections: one containing items that pass the truth test, and another containing those that don't.

```go
underage, adults := users.Partition(func(u User, id string) bool {
    return u.Age < 18
})
```

👉 [**Back to Overview**](overview.md)
