# validation

<!-- laravel-docs: validation.md#validation -->
<!-- laravel-docs: validation.md#validation-quickstart -->
<!-- laravel-docs: validation.md#manually-creating-validators -->
<!-- laravel-docs: validation.md#custom-validation-rules -->
<!-- laravel-docs: validation.md#working-with-error-messages -->
<!-- laravel-docs: validation.md#available-validation-rules -->

Rule-based input validation — a 1:1 Go port of Laravel's validator.

## Overview

The `validation` package provides a `Factory` that creates `Validator` instances.
It ports 80+ Laravel validation rules, pipe-delimited rule syntax, custom
messages, and the `MessageBag` error collection.

**Module:** `github.com/bedrock/packages/validation`

```bash
go get github.com/bedrock/packages/validation@latest
```

## Quick Start

```go
factory := validation.NewFactory()

v := factory.Make(
    // data to validate
    map[string]any{
        "name":  "Alice",
        "email": "alice@example.com",
        "age":   17,
    },
    // rules (pipe-delimited string, slice, or ValidationRule object)
    map[string]any{
        "name":  "required|string|max:255",
        "email": "required|email",
        "age":   "required|integer|min:18",
    },
    nil, // custom messages (optional)
    nil, // attribute display names (optional)
)

if v.Fails() {
    errs := v.Errors() // *validation.MessageBag
    fmt.Println(errs.All())
}
```

## Validate (shorthand)

`Validate` creates and runs the validator in one call, returning validated data
or a `*ValidationException` on failure:

```go
data, err := factory.Validate(
    map[string]any{"email": "bad"},
    map[string]any{"email": "required|email"},
    nil, nil,
)
```

## Rule Syntax

Rules are expressed as pipe-delimited strings, slices, or objects:

```go
// String
"required|string|max:255"

// Slice of strings
[]string{"required", "string", "max:255"}

// Mixed slice (strings + rule objects)
[]any{"required", rules.Min(18)}
```

## Custom Messages

```go
factory.Make(data, rules, map[string]string{
    "email.required": "An e-mail address is required.",
    "email.email":    "Please enter a valid e-mail.",
}, nil)
```

## Custom Attribute Names

```go
factory.Make(data, rules, nil, map[string]string{
    "email": "E-mail Address",
})
```

## Custom Rules

Register an extension on the factory:

```go
factory.Extend("uppercase", func(attr, value string, params []string, v *validation.Validator) bool {
    return value == strings.ToUpper(value)
})
```

Use `ExtendImplicit` for rules that should run even when the field is absent.

## MessageBag

```go
bag := v.Errors()

bag.Has("email")        // bool
bag.Get("email")        // []string{"The email field is required."}
bag.First("email")      // "The email field is required."
bag.All()               // map[string][]string
bag.ToJSON()            // JSON bytes
```

## Selected Built-in Rules

| Rule                  | Description                          |
| --------------------- | ------------------------------------ |
| `required`            | Field must be present and non-empty  |
| `string`              | Must be a string                     |
| `integer` / `int`     | Must be an integer                   |
| `numeric`             | Must be a number                     |
| `email`               | Must be a valid e-mail               |
| `url`                 | Must be a valid URL                  |
| `min:N`               | Minimum value / length / count       |
| `max:N`               | Maximum value / length / count       |
| `between:N,M`         | Value between N and M                |
| `in:a,b,c`            | Must be one of the listed values     |
| `unique:table,column` | Must not exist in the database       |
| `confirmed`           | Must match `{field}_confirmation`    |
| `date`                | Must be a parseable date             |
| `regex:pattern`       | Must match the regular expression    |
| `nullable`            | Allow null values to pass validation |
