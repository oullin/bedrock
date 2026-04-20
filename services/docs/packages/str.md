# str

<!-- upstream-docs: strings.md#strings -->
<!-- upstream-docs: strings.md#available-methods -->
<!-- upstream-docs: strings.md#fluent-strings -->

String helpers split out from the old `support` module.

## Overview

The `str` module contains Bedrock's Upstream-style `Str*` helpers, string
builder utilities, pluralization, UUID and ULID helpers, transliteration, and
Markdown rendering.

**Module:** `github.com/bedrock/packages/str`

```bash
go get github.com/bedrock/packages/str@latest
```

## Usage

```go
str.StrSlug("Hello World!")  // "hello-world"
str.StrCamel("hello_world")  // "helloWorld"
str.StrUuid()                // "550e8400-e29b-41d4-a716-..."
str.StrContains("abc", "b")  // true
str.StrLimit("long text", 5) // "long …"
```
