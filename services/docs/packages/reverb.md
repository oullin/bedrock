# websockets

WebSocket server implementing the Pusher protocol.

## Overview

The `websockets` package is a Go port of Upstream WebSockets — a self-hosted WebSocket
server that speaks the Pusher protocol. Clients connect using standard Pusher
client libraries (or `packages/broadcastclient`), authenticate via HMAC-SHA256, and receive
real-time events broadcast from the server.

**Module:** `github.com/bedrock/packages/websockets`

```bash
go get github.com/bedrock/packages/websockets@latest
```

## Features

- Drop-in Pusher replacement — no client-side changes needed
- HMAC-SHA256 channel authentication
- Public, private, and presence channel support
- Compatible with `packages/broadcastclient` for Go-side broadcasting

## Coming Soon

Full documentation is in progress.
