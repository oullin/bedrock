# 🧩 Helper Functions

The root `collection` package includes a set of general-purpose helper functions. These standalone utilities complement the collection types and simplify common Go patterns.

---

## 🛠 Available Helpers

| Function | Purpose |
|:---|:---|
| [**Value**](#value) | Returns the given value unchanged (no-op). |
| [**ValueFunc**](#valuefunc) | Calls a callback and returns its result (lazy evaluation). |
| [**Head**](#head) | Safely returns the first element of a slice. |
| [**Last**](#last) | Safely returns the last element of a slice. |
| [**WhenValue**](#whenvalue) | Inline conditional for values. |
| [**WhenFunc**](#whenfunc) | Inline conditional for callbacks (lazy). |
| [**Error Types**](#error-types) | Specialized errors for collection operations. |
| [**Shared Types**](#shared-types) | Common type constraints and structures. |

---

## 💎 Value

```go
func Value[T any](value T) T
```

Returns the given value unchanged. Useful as a no-op callback or for consistent value resolution patterns.

```go
v := collection.Value(42) // 42
```

---

## 💎 ValueFunc

```go
func ValueFunc[T any](callback func() T) T
```

Calls the given callback and returns its result. Useful for deferred evaluation.

```go
v := collection.ValueFunc(func() string {
    return computeExpensiveDefault()
})
```

---

## 💎 Head

```go
func Head[T any](items []T) (T, bool)
```

Returns the first element of a slice. Returns `false` if the slice is empty.

```go
first, ok := collection.Head([]string{"alice", "bob"})
// first = "alice", ok = true
```

---

## 💎 Last

```go
func Last[T any](items []T) (T, bool)
```

Returns the last element of a slice. Returns `false` if the slice is empty.

---

## 💎 WhenValue

```go
func WhenValue[T any](condition bool, value T, defaults ...T) T
```

Returns `value` if `condition` is true, otherwise returns the first default or the zero value.

```go
label := collection.WhenValue(isAdmin, "Admin Panel", "Dashboard")
```

---

## 💎 WhenFunc

```go
func WhenFunc[T any](condition bool, callback func() T, defaults ...func() T) T
```

Like `WhenValue` but accepts callbacks for deferred evaluation. Only the required branch is executed.

---

## 🛡 Error Types

### ItemNotFoundError
Returned when a requested item is missing (e.g., `FirstOrFail`, `Sole`).

### MultipleItemsFoundError
Returned by `Sole` when more than one item matches the predicate.

---

## 🧩 Shared Types

### Numeric
A type constraint for numeric types, used by aggregation functions like `Sum` and `Avg`.

### Pair[K, V]
A generic key-value pair used by `Combine` and `MapCollection`.

```go
type Pair[K any, V any] struct {
    Key   K
    Value V
}
```

---

👉 [**Back to Overview**](overview.md)
