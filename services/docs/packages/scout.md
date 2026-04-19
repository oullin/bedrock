# scout

Full-text search with pluggable engine backends.

## Overview

The `scout` package is a Go port of Laravel Scout. It provides automatic index
synchronisation on model create, update, and delete operations through the
`Searchable` interface and `SearchableMixin`.

**Module:** `github.com/bedrock/packages/scout`

```bash
go get github.com/bedrock/packages/scout@latest
```

## Supported Engines

| Engine        | Module                                                      |
| ------------- | ----------------------------------------------------------- |
| Database      | Built-in — SQL full-text search (MATCH/AGAINST, tsvector)  |
| Algolia       | `github.com/bedrock/packages/scout/engines/algolia`         |
| Meilisearch   | `github.com/bedrock/packages/scout/engines/meilisearch`     |
| Typesense     | `github.com/bedrock/packages/scout/engines/typesense`       |

## Sub-packages

| Package         | Description                               |
| --------------- | ----------------------------------------- |
| `scout/events`  | Scout-specific model indexing events      |
| `scout/jobs`    | Background indexing and removal jobs      |

## Coming Soon

Full documentation is in progress.
