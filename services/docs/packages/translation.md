# translation

Localisation and i18n support.

## Overview

The `translation` package provides a `Translator` that resolves translation keys
against a `Loader`, applies placeholder substitution, supports namespace/fallback
locale chains, and delegates pluralisation to a CLDR-based `MessageSelector`.

**Module:** `github.com/gocanto/bedrock/packages/translation`

```bash
go get github.com/gocanto/bedrock/packages/translation@latest
```

## Loaders

| Loader        | Description                                         |
|---------------|-----------------------------------------------------|
| `FileLoader`  | Reads JSON files from a directory tree              |
| `ArrayLoader` | In-memory map — ideal for tests and small apps      |

## File Structure

```
lang/
  en/
    messages.json   {"welcome": "Welcome, :name!"}
    auth.json       {"failed": "These credentials do not match."}
  fr/
    messages.json   {"welcome": "Bienvenue, :name !"}
```

## Creating a Translator

```go
loader := translation.NewFileLoader("/path/to/lang")
t := translation.NewTranslator(loader, "en")
t.SetFallback("en") // fallback locale
```

## Translating Keys

```go
t.Get("messages.welcome", map[string]string{"name": "Alice"})
// "Welcome, Alice!"

t.Get("auth.failed")
// "These credentials do not match."

// With a specific locale (without changing the global locale)
t.GetWithLocale("fr", "messages.welcome", map[string]string{"name": "Alice"})
// "Bienvenue, Alice !"
```

## Pluralization

```go
// messages.json: {"apples": "one apple|many apples"}
t.Choice("messages.apples", 1, nil)  // "one apple"
t.Choice("messages.apples", 5, nil)  // "many apples"

// Extended range syntax: "[0] no apples|[1] one apple|[2,*] :count apples"
t.Choice("messages.apples", 0, map[string]string{"count": "0"}) // "no apples"
```

## Namespaced Keys

```go
loader.AddNamespace("vendor/auth", "/path/to/vendor/auth/lang")
t.Get("vendor/auth::messages.login")
```

## In-memory Loader (tests)

```go
loader := translation.NewArrayLoader()
loader.Load("en", "messages", map[string]any{
    "greeting": "Hello, :name!",
})
t := translation.NewTranslator(loader, "en")
t.Get("messages.greeting", map[string]string{"name": "Bob"}) // "Hello, Bob!"
```

## Locale Management

```go
t.SetLocale("fr")
t.GetLocale()     // "fr"
t.SetFallback("en")
```
