# telescope

<!-- laravel-docs: telescope.md#laravel-telescope -->
<!-- laravel-docs: telescope.md#available-watchers -->
<!-- laravel-docs: telescope.md#filtering -->
<!-- laravel-docs: telescope.md#tagging -->

Debugging and introspection tool for Bedrock applications.

## Overview

The `telescope` package mirrors Laravel Telescope. It captures HTTP requests,
database queries, exceptions, log messages, events, queued jobs, cache
operations, mail, notifications, model changes, commands, scheduled tasks,
Redis commands, authorization gates, outbound HTTP calls, and debug dumps.
Every captured item is a typed, UUID-keyed entry grouped into batches and
persisted via a pluggable repository contract.

**Module:** `github.com/bedrock/packages/telescope`

```bash
go get github.com/bedrock/packages/telescope@latest
```

## Captured Entry Types

| Type          | Description                               |
| ------------- | ----------------------------------------- |
| Requests      | Incoming HTTP requests and responses      |
| Queries       | Database queries with bindings and timing |
| Exceptions    | Captured errors and stack traces          |
| Log           | Structured log messages                   |
| Events        | Dispatched application events             |
| Jobs          | Queued and processed jobs                 |
| Cache         | Cache hits, misses, and writes            |
| Mail          | Outgoing email messages                   |
| Notifications | Sent notifications across channels        |

## Coming Soon

Full documentation is in progress.
