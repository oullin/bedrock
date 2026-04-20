# pagination

<!-- upstream-docs: pagination.md#database-pagination -->
<!-- upstream-docs: pagination.md#basic-usage -->
<!-- upstream-docs: pagination.md#cursor-pagination -->
<!-- upstream-docs: pagination.md#displaying-pagination-results -->

Offset-based and cursor-based pagination.

## Overview

The `pagination` package provides three generic paginators that mirror Upstream's
pagination layer. Each paginates a Go slice and generates URL windows for
rendering page-link ranges.

**Module:** `github.com/gocanto/bedrock/packages/pagination`

```bash
go get github.com/gocanto/bedrock/packages/pagination@latest
```

## Paginator Types

| Type                      | Knows total count | Use when                            |
| ------------------------- | ----------------- | ----------------------------------- |
| `Paginator[T]`            | No                | Simple "next/previous" paging       |
| `LengthAwarePaginator[T]` | Yes               | Full page-number navigation         |
| `CursorPaginator[T]`      | No                | Stable cursor-based infinite scroll |

## Simple Paginator

```go
users := []User{ /* ... */ }
page  := 2
perPage := 15

p := pagination.NewPaginator(users, perPage, page, map[string]any{
    "path": "/users",
})

p.Items()           // []User (current page slice)
p.CurrentPage()     // 2
p.HasMorePages()    // bool
p.NextPageUrl()     // "/users?page=3"
p.PreviousPageUrl() // "/users?page=1"
```

## Length-Aware Paginator

```go
total := 142

lp := pagination.NewLengthAwarePaginator(users, total, perPage, page, map[string]any{
    "path": "/users",
})

lp.Total()          // 142
lp.LastPage()       // 10
lp.HasPages()       // true

// URL window for rendering page links
window := lp.LinkCollection() // []pagination.PaginatorLink
```

## Cursor Paginator

```go
cursors := /* decoded from query string */
cp := pagination.NewCursorPaginator(users, perPage, cursor, map[string]any{
    "path": "/feed",
})

cp.NextCursor()     // *pagination.Cursor
cp.PreviousCursor() // *pagination.Cursor
cp.NextPageUrl()    // "/feed?cursor=eyJ..."
```

## Appending Query Parameters

```go
p.AppendQueryString(map[string]string{"sort": "name"})
p.WithFragment("results")
```

## JSON Serialization

All paginators implement `json.Marshaler`:

```go
out, _ := json.Marshal(lp)
// {"current_page":2,"data":[...],"total":142,"last_page":10,...}
```
