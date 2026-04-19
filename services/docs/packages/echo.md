# broadcastclient

Real-time event broadcasting client (Go port of Upstream BroadcastClient).

## Overview

The `broadcastclient` package provides real-time event broadcasting abstractions over
multiple transport backends (Pusher, Socket.IO, Null/stub) with a uniform
`Channel` and `Connector` interface. It is a Go port of the Upstream BroadcastClient
JavaScript library.

**Module:** `github.com/bedrock/packages/broadcastclient`

```bash
go get github.com/bedrock/packages/broadcastclient@latest
```

## Supported Backends

| Backend   | Description                                  |
| --------- | -------------------------------------------- |
| Pusher    | Connects to Pusher or the `websockets` server    |
| Socket.IO | Socket.IO-compatible transports              |
| Null      | No-op stub for testing                       |

## Channel Types

| Type      | Description                                        |
| --------- | -------------------------------------------------- |
| Public    | Unauthenticated channel open to all connections    |
| Private   | Authenticated, user-scoped channel                 |
| Presence  | Authenticated with member tracking                 |

## Coming Soon

Full documentation is in progress.
