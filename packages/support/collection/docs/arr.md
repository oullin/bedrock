# 🧩 Array Utilities (`arr`)

`import "github.com/gollin/packages/support/collection/arr"`

The `arr` package provides standalone, generic helper functions for working with Go slices. Every function operates on plain slices, returns new slices (the original is never mutated), and can be called without constructing a collection object.

Use `arr` when you need a single operation on a slice and do not need the fluent chaining that `collection.Collection` provides.

---

## 🛠 Available Functions

| Function                            | Purpose                                                                        |
| :---------------------------------- | :----------------------------------------------------------------------------- |
| [**Accessible**](#accessible)       | Reports whether the given value is non-nil.                                    |
| [**IsList**](#islist)               | Reports whether the given slice is a sequential list.                          |
| [**First**](#first)                 | Returns the first element or the first matching a predicate.                   |
| [**Last**](#last)                   | Returns the last element or the last matching a predicate.                     |
| [**Take**](#take)                   | Returns a new slice containing up to N elements.                               |
| [**Only**](#only)                   | Returns a new slice containing only the elements at the given indices.         |
| [**Except**](#except)               | Returns a new slice containing all elements except those at the given indices. |
| [**Flatten**](#flatten)             | Flattens a slice of slices into a single, flat slice.                          |
| [**Collapse**](#collapse)           | Alias for `Flatten`.                                                           |
| [**Wrap**](#wrap)                   | Wraps a single value in a one-element slice.                                   |
| [**WrapSlice**](#wrapslice)         | Returns the given slice unchanged.                                             |
| [**Prepend**](#prepend)             | Inserts a value at the beginning of the slice.                                 |
| [**Push**](#push)                   | Appends one or more values to the end of the slice.                            |
| [**Shuffle**](#shuffle)             | Returns a new slice with the elements in random order.                         |
| [**Random**](#random)               | Returns a new slice containing N randomly selected elements.                   |
| [**Sort**](#sort)                   | Returns a new slice sorted using a comparison function.                        |
| [**SortDesc**](#sortdesc)           | Returns a new slice sorted in descending order.                                |
| [**SortRecursive**](#sortrecursive) | Sorts a slice using a comparison function.                                     |
| [**Where**](#where)                 | Returns a new slice containing only elements matching a callback.              |
| [**WhereNotNull**](#wherenotnull)   | Returns a new slice with all zero-value elements removed.                      |
| [**Reject**](#reject)               | Returns a new slice excluding elements matching a callback.                    |
| [**Partition**](#partition)         | Splits the slice into two based on a callback.                                 |
| [**Every**](#every)                 | Reports whether every element satisfies a callback.                            |
| [**Some**](#some)                   | Reports whether at least one element satisfies a callback.                     |
| [**Exists**](#exists)               | Reports whether the given index is valid.                                      |
| [**Has**](#has)                     | Reports whether all of the given indices are valid.                            |
| [**HasAny**](#hasany)               | Reports whether at least one of the given indices is valid.                    |
| [**Join**](#join)                   | Concatenates string slice elements with a glue string.                         |
| [**CrossJoin**](#crossjoin)         | Returns the Cartesian product of the given slices.                             |
| [**Divide**](#divide)               | Splits a slice into indices and values.                                        |
| [**Map**](#map)                     | Transforms each element into a new slice of results.                           |
| [**MapWithKeys**](#mapwithkeys)     | Transforms each element into a key-value pair map.                             |
| [**MapSpread**](#mapspread)         | Alias for `Map`.                                                               |
| [**KeyBy**](#keyby)                 | Indexes the slice elements by a key returned from a function.                  |
| [**Pluck**](#pluck)                 | Extracts a single field from each element.                                     |

---

## 💎 Accessible

```go
func Accessible(value any) bool
```

Reports whether the given value is non-nil. This is useful as a guard before performing operations on interface values that may be nil.

```go
var s []int
arr.Accessible(s)   // false -- nil slice
arr.Accessible(42)  // true

m := map[string]int{}
arr.Accessible(m)   // true
```

---

## 💎 IsList

```go
func IsList[T any](items []T) bool
```

Reports whether the given slice is a sequential list. In Go every slice is inherently a sequential, integer-indexed list, so this always returns `true`.

```go
arr.IsList([]string{"a", "b", "c"}) // true
arr.IsList([]int{})                 // true
```

---

## 💎 First

```go
func First[T any](items []T, callbacks ...func(T, int) bool) (T, bool)
```

Returns the first element of the slice, or the first element that matches the provided callback predicate.

```go
// Without a callback -- returns the first element.
val, ok := arr.First([]int{10, 20, 30})
// val == 10, ok == true

// With a callback -- returns the first even number.
val, ok = arr.First([]int{1, 3, 4, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// val == 4, ok == true
```

---

## 💎 Last

```go
func Last[T any](items []T, callbacks ...func(T, int) bool) (T, bool)
```

Returns the last element of the slice, or the last element that matches the provided callback predicate.

```go
val, ok := arr.Last([]string{"a", "b", "c"})
// val == "c", ok == true

val, ok = arr.Last([]int{1, 2, 3, 4, 5}, func(v int, _ int) bool {
    return v < 4
})
// val == 3, ok == true
```

---

## 💎 Take

```go
func Take[T any](items []T, limit int) []T
```

Returns a new slice containing up to `limit` elements. A positive limit takes from the front; a negative limit takes from the end.

```go
arr.Take([]int{1, 2, 3, 4, 5}, 3)   // [1, 2, 3]
arr.Take([]int{1, 2, 3, 4, 5}, -2)  // [4, 5]
```

---

## 💎 Only

```go
func Only[T any](items []T, indices []int) []T
```

Returns a new slice containing only the elements at the given indices. Out-of-range indices are silently skipped.

```go
arr.Only([]string{"a", "b", "c", "d"}, []int{0, 2})
// ["a", "c"]
```

---

## 💎 Except

```go
func Except[T any](items []T, indices []int) []T
```

Returns a new slice containing all elements except those at the given indices.

```go
arr.Except([]string{"a", "b", "c", "d"}, []int{1, 3})
// ["a", "c"]
```

---

## 💎 Flatten

```go
func Flatten[T any](items [][]T) []T
```

Flattens a slice of slices into a single, flat slice.

```go
nested := [][]int{{1, 2}, {3, 4}, {5}}
arr.Flatten(nested)
// [1, 2, 3, 4, 5]
```

---

## 💎 Collapse

```go
func Collapse[T any](items [][]T) []T
```

Merges a slice of slices into a single slice. This is an alias for `Flatten`.

---

## 💎 Wrap

```go
func Wrap[T any](value T) []T
```

Wraps a single value in a one-element slice.

```go
arr.Wrap(42)       // [42]
arr.Wrap("hello")  // ["hello"]
```

---

## 💎 WrapSlice

```go
func WrapSlice[T any](value []T) []T
```

Returns the given slice unchanged. Use this when the value is already a slice and you want to avoid double-wrapping.

---

## 💎 Prepend

```go
func Prepend[T any](items []T, value T) []T
```

Inserts a value at the beginning of the slice and returns the new slice.

```go
arr.Prepend([]int{2, 3, 4}, 1)
// [1, 2, 3, 4]
```

---

## 💎 Push

```go
func Push[T any](items []T, values ...T) []T
```

Appends one or more values to the end of the slice and returns the new slice.

---

## 💎 Shuffle

```go
func Shuffle[T any](items []T) []T
```

Returns a new slice with the elements in random order. The original slice is not modified.

---

## 💎 Random

```go
func Random[T any](items []T, counts ...int) []T
```

Returns a new slice containing `count` randomly selected elements. Defaults to 1.

```go
arr.Random([]string{"a", "b", "c", "d"}, 2)
// e.g. ["c", "a"]
```

---

## 💎 Sort

```go
func Sort[T any](items []T, less func(a, b T) bool) []T
```

Returns a new slice sorted using the provided comparison function. The sort is stable.

```go
arr.Sort([]int{3, 1, 4, 1, 5}, func(a, b int) bool {
    return a < b
})
// [1, 1, 3, 4, 5]
```

---

## 💎 SortDesc

```go
func SortDesc[T any](items []T, less func(a, b T) bool) []T
```

Returns a new slice sorted in descending order.

---

## 💎 SortRecursive

```go
func SortRecursive[T any](items []T, less func(a, b T) bool) []T
```

Sorts a slice using the provided comparison function. Exists for API parity with Upstream collections.

---

## 💎 Where

```go
func Where[T any](items []T, callback func(T, int) bool) []T
```

Returns a new slice containing only the elements for which the callback returns `true`.

```go
evens := arr.Where([]int{1, 2, 3, 4, 5, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
// [2, 4, 6]
```

---

## 💎 WhereNotNull

```go
func WhereNotNull[T comparable](items []T) []T
```

Returns a new slice with all zero-value elements removed. The type must satisfy `comparable`.

```go
arr.WhereNotNull([]string{"a", "", "b", "", "c"})
// ["a", "b", "c"]
```

---

## 💎 Reject

```go
func Reject[T any](items []T, callback func(T, int) bool) []T
```

Returns a new slice containing only the elements for which the callback returns `false`. It is the inverse of `Where`.

---

## 💎 Partition

```go
func Partition[T any](items []T, callback func(T, int) bool) ([]T, []T)
```

Splits the slice into two: the first contains elements where the callback returns `true`, the second contains the rest.

```go
pass, fail := arr.Partition([]int{1, 2, 3, 4, 5, 6}, func(v int, _ int) bool {
    return v%2 == 0
})
```

---

## 💎 Every

```go
func Every[T any](items []T, callback func(T, int) bool) bool
```

Reports whether every element in the slice satisfies the callback.

---

## 💎 Some

```go
func Some[T any](items []T, callback func(T, int) bool) bool
```

Reports whether at least one element in the slice satisfies the callback.

---

## 💎 Exists

```go
func Exists[T any](items []T, index int) bool
```

Reports whether the given index is valid (in bounds) for the slice.

---

## 💎 Has

```go
func Has[T any](items []T, indices ...int) bool
```

Reports whether **all** of the given indices are valid for the slice.

---

## 💎 HasAny

```go
func HasAny[T any](items []T, indices ...int) bool
```

Reports whether **at least one** of the given indices is valid for the slice.

---

## 💎 Join

```go
func Join(items []string, glue string, finalGlues ...string) string
```

Concatenates string slice elements with a glue string. An optional final glue is placed between the last two elements.

```go
arr.Join([]string{"a", "b", "c"}, ", ", " and ")
// "a, b and c"
```

---

## 💎 CrossJoin

```go
func CrossJoin[T any](lists ...[]T) [][]T
```

Returns the Cartesian product of the given slices.

```go
arr.CrossJoin([]int{1, 2}, []int{10, 20})
// [[1, 10], [1, 20], [2, 10], [2, 20]]
```

---

## 💎 Divide

```go
func Divide[T any](items []T) ([]int, []T)
```

Splits a slice into two: a slice of indices and a slice of values.

---

## 💎 Map

```go
func Map[T any, R any](items []T, callback func(T, int) R) []R
```

Applies the callback to each element and returns a new slice of the transformed results.

```go
arr.Map([]int{1, 2, 3}, func(v int, _ int) int {
    return v * 2
})
// [2, 4, 6]
```

---

## 💎 MapWithKeys

```go
func MapWithKeys[T any, K comparable, V any](items []T, callback func(T) (K, V)) map[K]V
```

Applies the callback to each element, which produces a key-value pair, and collects the results into a map.

---

## 💎 MapSpread

```go
func MapSpread[T any, R any](items []T, callback func(T, int) R) []R
```

Alias for `Map`. Exists for API parity with Upstream collections.

---

## 💎 KeyBy

```go
func KeyBy[T any, K comparable](items []T, keyFunc func(T) K) map[K]T
```

Indexes the slice elements by the key returned from `keyFunc`, producing a map.

---

## 💎 Pluck

```go
func Pluck[T any, V any](items []T, valueFunc func(T) V) []V
```

Extracts a single field or derived value from each element and returns the collected values as a slice.

```go
type Employee struct {
    Name   string
    Salary float64
}

employees := []Employee{
    {Name: "Alice", Salary: 90000},
    {Name: "Bob", Salary: 75000},
}

names := arr.Pluck(employees, func(e Employee) string {
    return e.Name
})
// ["Alice", "Bob"]
```

---

👉 [**Back to Overview**](overview.md)
