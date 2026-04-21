# broadcasting

<!-- laravel-docs: broadcasting.md#broadcasting -->
<!-- laravel-docs: broadcasting.md#server-side-installation -->
<!-- laravel-docs: broadcasting.md#defining-broadcast-events -->
<!-- laravel-docs: broadcasting.md#broadcasting-events -->
<!-- laravel-docs: broadcasting.md#authorizing-channels -->

Laravel-style server-side broadcasting for Go applications.

## Overview

The `broadcasting` package owns the framework side of real-time events:
channel authorization, broadcast event payloads, named broadcaster
connections, and Pusher-compatible channel conventions.

**Module:** `github.com/bedrock/packages/broadcasting`

```bash
go get github.com/bedrock/packages/broadcasting@latest
```

Use `packages/echo` on the client side to subscribe to channels and receive
events. Use `packages/reverb` when you need a Pusher-compatible WebSocket
server.

## Features

| Feature          | Description                                                                     |
| ---------------- | ------------------------------------------------------------------------------- |
| Channel auth     | Register channel patterns and authorize private or presence subscriptions       |
| Broadcast events | Resolve event names, channels, payloads, sockets, middleware, and failure hooks |
| Backends         | Pusher/Reverb, Redis, and Ably broadcaster implementations                      |
| Conventions      | Pusher-style `private-`, `private-encrypted-`, and `presence-` channel names    |

## Channel Authorization

```go
broadcaster := broadcasting.NewPusherBroadcaster(client)

broadcaster.Channel("orders.{id}", broadcasting.ChannelHandler(
	func(user auth.Authenticatable, params ...any) (any, error) {
		orderID := params[0].(string)

		if user.GetAuthIdentifier() == orderID {
			return true, nil
		}

		return false, nil
	},
))
```

## Echo Link

`broadcasting` authorizes and publishes server-side events. `echo` is the
client package that subscribes to those channels and listens for the events.
