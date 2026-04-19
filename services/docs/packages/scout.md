# search

Full-text search with pluggable engine backends.

## Overview

The `search` package is a Go port of Upstream Search. It provides automatic index
synchronisation on model create, update, and delete operations through the
`Searchable` interface and `SearchableMixin`.

**Module:** `github.com/bedrock/packages/search`

```bash
go get github.com/bedrock/packages/search@latest
```

## Supported Engines

| Engine        | Module                                                      |
| ------------- | ----------------------------------------------------------- |
| Database      | Built-in — SQL full-text search (MATCH/AGAINST, tsvector)  |
| Algolia       | `github.com/bedrock/packages/search/engines/algolia`         |
| Meilisearch   | `github.com/bedrock/packages/search/engines/meilisearch`     |
| Typesense     | `github.com/bedrock/packages/search/engines/typesense`       |

## Sub-packages

| Package         | Description                               |
| --------------- | ----------------------------------------- |
| `search/events`  | Search-specific model indexing events      |
| `search/jobs`    | Background indexing and removal jobs      |

## Coming Soon

Full documentation is in progress.
