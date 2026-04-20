# support

<!-- upstream-docs: helpers.md#helpers -->
<!-- upstream-docs: helpers.md#arrays-objects -->
<!-- upstream-docs: helpers.md#other-utilities -->

General-purpose helpers and types from Bedrock's Upstream support port.

## Overview

The `support` package provides the utility layer shared across Bedrock: global
helpers, array and dot-notation helpers, `Fluent` dynamic bags, `Optional[T]`,
`MessageBag`, sleeping helpers, and `Timebox`.

**Module:** `github.com/bedrock/packages/support`

```bash
go get github.com/bedrock/packages/support@latest
```

## Global Helpers

```go
support.Blank("")          // true  — nil, empty string, empty slice/map
support.Filled("hello")    // true  — not blank

support.Tap(user, func(u *User) { u.Name = "Alice" }) // returns user after callback

// Resolve a value or call it if it's a function
support.Value[string]("hello")            // "hello"
support.Value[string](func() string { return "hi" }) // "hi"

// Pass through an optional transform
support.With(user, func(u *User) *User { u.Active = true; return u })

// Transform to a different type
name := support.Transform[*User, string](user, func(u *User) string { return u.Name })

support.Env("APP_KEY", "default") // os.Getenv with fallback

// Retry with exponential backoff
result, err := support.Retry(3, func() (string, error) {
    return callExternalAPI()
}, 100*time.Millisecond)
```

## Fluent

`Fluent` is a dynamic key-value bag inspired by Upstream's `Fluent` class:

```go
f := support.NewFluent(map[string]any{
    "name":  "Alice",
    "email": "alice@example.com",
})

f.Get("name")           // "Alice"
f.Set("role", "admin")
f.Has("email")          // true
f.All()                 // map[string]any{...}
f.ToJSON()              // []byte
```

## Optional

`Optional[T]` wraps a nullable value — a Go equivalent of Upstream's `optional()`:

```go
opt := support.NewOptional[*User](user)

opt.IsPresent()          // true
opt.Get()                // *User
opt.OrElse(guestUser)    // returns user if present, otherwise guestUser
opt.Map(func(u *User) string { return u.Name }) // Optional[string]
```

## MessageBag

```go
bag := support.NewMessageBag(map[string][]string{
    "email": {"The email field is required."},
})

bag.Add("name", "The name must be at least 3 characters.")
bag.Has("email")    // true
bag.Get("email")    // []string
bag.First("email")  // "The email field is required."
bag.All()           // map[string][]string
bag.IsEmpty()       // false
```

## Array Helpers

```go
values := map[string]any{"user": map[string]any{"name": "Ada"}}

support.ArrGet(values, "user.name")            // "Ada"
support.ArrSet(values, "user.email", "a@b.c")  // mutates nested map
support.ArrHas(values, "user.name", "user.email")
support.ArrDot(values)                         // flatten nested map to dot keys
```

## Split Packages

String utilities and probabilistic execution now live in dedicated modules:

```go
import (
    "github.com/bedrock/packages/lottery"
    "github.com/bedrock/packages/str"
)

str.StrSlug("Hello World!")
lottery.NewLottery(1, 100).Choose()
```

## Sleep

Fluent sleeping with pause/fake support:

```go
support.Sleep(5 * time.Second)
support.SleepUntil(target time.Time)

// In tests — fake and assert
support.FakeSleep()
support.AssertSlept(5 * time.Second, 1)
```

## Timebox

Cap execution to a fixed duration:

```go
support.Timebox(100*time.Millisecond, func() {
    expensiveOperation()
})
```
