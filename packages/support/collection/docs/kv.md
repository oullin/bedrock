# 🧩 Key-Value Utilities (`kv`)

`import "github.com/gollin/packages/support/collection/kv"`

The `kv` package provides standalone utility functions for working with Go maps. Its headline feature is **dot-notation access** for nested structures, along with general-purpose filtering, sorting, and transformation helpers.

---

## 🛠 Available Functions

### Dot-Notation Helpers

| Function              | Purpose                                                        |
| :-------------------- | :------------------------------------------------------------- |
| [**Get**](#get)       | Retrieves a value from a nested map using dot-separated paths. |
| [**Set**](#set)       | Sets a value in a nested map using a dot-separated path.       |
| [**Has**](#has)       | Reports whether a dot-notated key exists.                      |
| [**Fill**](#fill)     | Sets a value only if the key does not exist.                   |
| [**Add**](#add)       | Alias for `Fill`.                                              |
| [**Forget**](#forget) | Removes a value at a given dot-notated key.                    |
| [**Pull**](#pull)     | Retrieves and then removes a value.                            |
| [**Dot**](#dot)       | Flattens a nested map into a single-level map.                 |
| [**Undot**](#undot)   | Expands a flat map back into a nested structure.               |

### Map Selection & Transformation

| Function                | Purpose                                               |
| :---------------------- | :---------------------------------------------------- |
| [**Only**](#only)       | Returns a new map containing only the specified keys. |
| [**Except**](#except)   | Returns a new map with the specified keys removed.    |
| [**Map**](#map)         | Transforms map values while preserving keys.          |
| [**Where**](#where)     | Filters map entries by a custom predicate.            |
| [**Replace**](#replace) | Merges maps, with later maps taking precedence.       |
| [**Sort**](#sort)       | Returns a new map with keys in ascending order.       |

### Encoding

| Function                          | Purpose                                         |
| :-------------------------------- | :---------------------------------------------- |
| [**Query**](#query)               | Encodes a map as a URL query string.            |
| [**ToCssClasses**](#tocssclasses) | Generates a CSS class list from a boolean map.  |
| [**ToCssStyles**](#tocssstyles)   | Generates inline CSS styles from a boolean map. |

---

## 💎 Dot-Notation Helpers

These functions navigate nested `map[string]any` structures using paths such as `"user.address.city"`.

### Get

```go
func Get(target map[string]any, key string, defaults ...any) any
```

Retrieves a nested value. Returns `nil` or a default if not found.

### Set

```go
func Set(target map[string]any, key string, value any, overwrite ...bool) map[string]any
```

Sets a nested value, creating intermediate maps as needed.

### Has

```go
func Has(target map[string]any, key string) bool
```

Reports if a nested key exists.

### Forget

```go
func Forget(target map[string]any, key string) map[string]any
```

Removes a nested key.

---

## 💎 Transformation & Selection

### Only / Except

```go
func Only(items map[string]any, keys ...string) map[string]any
func Except(items map[string]any, keys ...string) map[string]any
```

Whitelists or blacklists specific keys.

### Map

```go
func Map[V any, R any](items map[string]V, callback func(V, string) R) map[string]R
```

Transforms values using generics.

### Sort

```go
func Sort(items map[string]any) map[string]any
```

Returns a new map with keys sorted alphabetically.

---

## 💎 Encoding

### Query

```go
func Query(items map[string]any) string
```

Converts a map to a URL-encoded query string.

### ToCssClasses

```go
func ToCssClasses(classes map[string]bool) string
```

Builds a CSS class string from active (true) boolean flags.

---

👉 [**Back to Overview**](overview.md)
