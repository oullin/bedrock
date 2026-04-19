# helpers

Laravel-style global helper functions for arrays and strings.

## Overview

The `helpers` package is a Go port of `laravel/helpers`. It provides array
helpers (`Arr*`) for map and slice operations with dot-notation support, and
string helpers that delegate to `packages/support`.

**Module:** `github.com/bedrock/packages/helpers`

```bash
go get github.com/bedrock/packages/helpers@latest
```

## Array Helpers

Array helpers implement the functionality of Laravel's `Illuminate\Support\Arr`,
offering dot-notation access, flattening, grouping, and more.

## String Helpers

String helpers re-export the corresponding `Str*` functions from
`packages/support`.

## Coming Soon

Full documentation is in progress.
