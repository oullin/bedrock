# Collection

[![Go Reference](https://pkg.go.dev/badge/github.com/bedrock/packages/collection.svg)](https://pkg.go.dev/github.com/bedrock/packages/collection)
[![Go Version](https://img.shields.io/badge/go-1.26-blue.svg)](https://golang.org/doc/devel/release.html#go1.26)
[![CI](https://github.com/gocanto/bedrock/actions/workflows/ci.yml/badge.svg)](https://github.com/gocanto/bedrock/actions)
[![codecov](https://codecov.io/gh/gocanto/bedrock/branch/main/graph/badge.svg)](https://codecov.io/gh/gocanto/bedrock)
[![Go Report Card](https://goreportcard.com/badge/github.com/bedrock/packages/collection)](https://goreportcard.com/report/github.com/bedrock/packages/collection)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful Go port of [Upstream Collections](https://upstream.com/docs/collections) — fluent, type-safe, and powered by Go generics and `iter.Seq` lazy evaluation.

```bash
go get github.com/bedrock/packages/collection
```

---

## ⚡️ Why use this?

Go's slices and maps are powerful but often lead to repetitive boilerplate loops for common operations. `collection` provides a **fluent, expressive API** to handle data transformations with ease:

- **💎 Fluent Pipelines:** Chain `.Filter().Take().Each()` instead of nesting `for` loops.
- **🛡️ Type-Safe:** Built with Go generics—no `interface{}` casting or runtime type errors.
- **♻️ Immutable by Default:** Transformation methods return new collections, preserving your original data.
- **🐢 Lazy Evaluation:** Process massive datasets efficiently using `iter.Seq` sequences.
- **🧩 Modular:** Use the full fluent API or lightweight standalone utilities (`arr` and `kv`).

---

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/bedrock/packages/collection/collection"
)

func main() {
    // 1. Create a collection from a slice
    numbers := collection.Collect([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

    // 2. Chain operations fluently
    even := numbers.
        Filter(func(n int, _ int) bool { return n % 2 == 0 }).
        Take(3)

    // 3. Transform types using top-level generics
    // (Go doesn't allow methods to introduce new type parameters)
    result := collection.Map(even, func(n int, _ int) string {
        return fmt.Sprintf("Number: %d", n)
    })

    fmt.Println(result.All()) // ["Number: 2", "Number: 4", "Number: 6"]
}
```

---

## 📦 Package Ecosystem

| Package | Purpose | Use when... |
|:---|:---|:---|
| **`collection`** | **The Core** | You want the full fluent `Collection[T]` experience for slices. |
| **`lazy`** | **Lazy Sequences** | You're handling large/infinite streams and want deferred execution. |
| **`collectible`** | **Key-Value Maps** | You need an ordered `collectible.Collection[K, V]` with a fluent API. |
| **`arr`** | **Slice Utils** | You need a quick one-off helper (e.g., `Sort`, `Flatten`) on a raw `[]T`. |
| **`kv`** | **Map Utils** | You need dot-notation access (`"user.profile.name"`) for `map[string]any`. |

---

## 📖 Deep Dive

Explore our comprehensive documentation for each component:

- [**Architecture Overview**](docs/overview.md) — Design philosophy and core concepts.
- [**Collection API**](docs/collection.md) — Reference for the standard fluent collection.
- [**Collectible API**](docs/collectible.md) — Reference for key-value collection operations.
- [**Lazy API**](docs/lazy.md) — Reference for lazy sequences and `iter.Seq`.
- [**Array Utilities (`arr`)**](docs/arr.md) — Standalone functions for raw slices.
- [**Key-Value Utilities (`kv`)**](docs/kv.md) — Standalone functions for maps and dot-notation.

---

## ⚖️ License

Released under the [MIT License](LICENSE).
