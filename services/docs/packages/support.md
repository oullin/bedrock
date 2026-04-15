# support

General-purpose helpers and types — a Go port of `Framework\Support`.

## Overview

The `support` package provides the utility layer shared across Bedrock: global
helpers, `Fluent` dynamic bags, `Optional[T]`, `MessageBag`, string builders,
probabilistic execution, and sleeping.

**Module:** `github.com/gocanto/bedrock/packages/support`

```bash
go get github.com/gocanto/bedrock/packages/support@latest
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

## String Helpers

```go
support.StrSlug("Hello World!")     // "hello-world"
support.StrCamel("hello_world")     // "helloWorld"
support.StrStudly("hello_world")    // "HelloWorld"
support.StrSnake("HelloWorld")      // "hello_world"
support.StrUuid()                   // "550e8400-e29b-41d4-a716-..."
support.StrContains("abc", "b")     // true
support.StrLimit("long text", 5)    // "long …"
```

## Lottery

Execute code probabilistically:

```go
support.Lottery(1, 100).Winner(func() {
    // runs ~1% of the time
}).Loser(func() {
    // runs the other 99%
}).Choose()
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
