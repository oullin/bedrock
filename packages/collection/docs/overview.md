# Overview

A Go port of [Upstream Collections](https://upstream.com/docs/collections) — fluent, type-safe, and powered by Go generics and `iter.Seq` lazy evaluation.

---

## 🎯 Why This Library Exists

Go's built-in slice and map operations are intentionally minimal. Common tasks like filtering, mapping, partitioning, and reducing require writing boilerplate loops every time. This library provides:

- **💎 Readable Pipelines:** Chain `.Filter().Map().Take()` instead of nesting loops.
- **🛡️ Type-safe Generics:** No `interface{}` casting; the compiler catches type errors at build time.
- **♻️ Immutable Returns:** Methods return new collections, leaving the original untouched.
- **🐢 Lazy Evaluation:** Process large or infinite datasets without allocating the entire result up front.
- **🧩 Standalone Utils:** Use individual functions from `arr` or `kv` for quick, one-off operations.

---

## 🏗️ Package Structure

```text
github.com/bedrock/packages/collection
    ├── collection/      Core: Collection[T] (fluent slice wrapper)
    ├── lazy/            Lazy sequences backed by iter.Seq[T]
    ├── collectible/     collectible.Collection[K, V] (ordered map with fluent API)
    ├── arr/             Standalone generic slice utilities
    └── kv/              Map utilities with dot-notation support
```

### Core Packages

- **`collection`**: The primary workhorse for fluent slice manipulation.
- **`lazy`**: Provides deferred evaluation for large or infinite datasets.
- **`collectible`**: Handles key-value data while maintaining insertion order.

---

## 🚀 Quick Start

### Basic Collection

```go
import "github.com/bedrock/packages/collection/collection"

// Create a collection from a slice
numbers := collection.Collect([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

// Filter to even numbers, take the first 3
result := numbers.
    Filter(func(n int, _ int) bool { return n % 2 == 0 }).
    Take(3)

fmt.Println(result.All()) // [2, 4, 6]
```

### Map Collection

```go
import "github.com/bedrock/packages/collection/collectible"

m := collectible.New(map[string]int{
    "apples":  5,
    "bananas": 3,
})

m.IsNotEmpty() // true
```

### Lazy Collection

```go
import "github.com/bedrock/packages/collection/lazy"

// Computation is deferred until .All() or .Iter() is called
lc := lazy.Range(1, 1000000).
    Filter(func(n int, _ int) bool { return n % 2 == 0 }).
    Take(5)

fmt.Println(lc.All()) // [2, 4, 6, 8, 10]
```

---

## 🛠️ Key Design Decisions

### 1. Go Generics
Every collection type is parameterized (`Collection[T]`, `collectible.Collection[K, V]`).
- **No type assertions:** Zero runtime overhead for type checking.
- **Compile-time safety:** Errors are caught during development, not in production.

### 2. `iter.Seq` Integration
`lazy.Collection` is built on Go's standard iterator protocol (`iter.Seq[T]`).
- **Native Range Support:** Use `range` directly over `lc.Iter()`.
- **Deferred Execution:** Computation only happens when results are requested.

### 3. Immutable Returns
Methods like `.Filter()`, `.Take()`, or `.Flatten()` return **new** collections.
- **Side-effect free:** Original data remains untouched.
- **Concurrency safe:** Collections can be shared across goroutines for read-only access.

---

## 📊 Comparison: Which to Use?

| Type | Best for... |
|:---|:---|
| **`Collection[T]`** | Standard slice manipulation with fluent chaining. |
| **`collectible.Collection[K, V]`** | Key-value data requiring ordered iteration or set operations. |
| **`lazy.Collection[T]`** | Large datasets or streams where deferred execution is critical. |
| **`arr.FuncName`** | Single, one-off operations on raw slices. |
| **`kv.FuncName`** | Nested maps (JSON/Config) using dot-notation paths. |

---

## 📚 Further Reading

- [**Collection API**](collection.md)
- [**Collectible API**](collectible.md)
- [**Lazy API**](lazy.md)
- [**Array Utilities (`arr`)**](arr.md)
- [**Key-Value Utilities (`kv`)**](kv.md)
- [**Helper Functions**](helpers.md)
