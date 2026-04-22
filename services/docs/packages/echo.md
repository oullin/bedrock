# echo

<!-- laravel-docs: broadcasting.md#receiving-broadcasts -->
<!-- laravel-docs: broadcasting.md#presence-channels -->
<!-- laravel-docs: broadcasting.md#client-events -->

Real-time event broadcasting client (Go port of Laravel Echo).

## Overview

The `echo` package provides real-time event broadcasting abstractions over
multiple transport backends (Pusher, Socket.IO, Null/stub) with a uniform
`Channel` and `Connector` interface. It is a Go port of the Laravel Echo
JavaScript library.

For server-side channel authorization and event publishing, use
`github.com/bedrock/packages/broadcasting`.

**Module:** `github.com/bedrock/packages/echo`

```bash
go get github.com/bedrock/packages/echo@latest
```

## Supported Backends

| Backend   | Description                               |
| --------- | ----------------------------------------- |
| Pusher    | Connects to Pusher or the `reverb` server |
| Socket.IO | Socket.IO-compatible transports           |
| Null      | No-op stub for testing                    |

## Channel Types

| Type     | Description                                     |
| -------- | ----------------------------------------------- |
| Public   | Unauthenticated channel open to all connections |
| Private  | Authenticated, user-scoped channel              |
| Presence | Authenticated with member tracking              |

## Coming Soon

Full documentation is in progress.
