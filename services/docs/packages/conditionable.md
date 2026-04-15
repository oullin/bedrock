# conditionable

Conditional method execution with a fluent proxy.

## Overview

The `conditionable` package provides generic `When` and `Unless` functions for
conditional branching without breaking a fluent call chain. It mirrors
`Framework\Support\Traits\Conditionable`.

**Module:** `github.com/gocanto/bedrock/packages/conditionable`

```bash
go get github.com/gocanto/bedrock/packages/conditionable@latest
```

## When

Execute a callback only when a value is truthy:

```go
result := conditionable.When(
    builder,          // target
    isAdmin,          // condition value
    func(b *Builder, v bool) *Builder {
        return b.Where("role", "admin")
    },
)
```

With an optional else branch:

```go
result := conditionable.When(
    builder,
    isAdmin,
    func(b *Builder, v bool) *Builder { return b.Where("role", "admin") },
    func(b *Builder, v bool) *Builder { return b.Where("role", "user") },
)
```

## Unless

Execute a callback only when a value is falsy:

```go
result := conditionable.Unless(
    query,
    isCached,
    func(q *Query, v bool) *Query { return q.Load("relations") },
)
```

## WhenFunc / UnlessFunc

Resolve the condition lazily from the target:

```go
result := conditionable.WhenFunc(
    user,
    func(u *User) bool { return u.IsVerified() },
    func(u *User, verified bool) *User {
        u.SendWelcomeEmail()
        return u
    },
)
```

## Truthiness Rules

A value is considered truthy if it is:

- A `bool` set to `true`
- A non-zero number
- A non-empty string
- A non-nil pointer, slice, map, or channel
- Any non-zero comparable value

Zero values, empty collections, nil pointers, and empty strings are falsy.
