package websockets

import (
	"context"
	"time"
)

// PingInactiveConnections sends a pusher:ping to every connection that has
// been inactive for longer than its application's ActivityTimeout. It
// iterates all registered applications and all connections for each one.
func PingInactiveConnections(ctx context.Context, conns *ConnectionManager, apps *AppManager) {
	ping := []byte(`{"event":"pusher:ping","data":"{}"}`)

	for _, app := range apps.All() {
		timeout := time.Duration(app.ActivityTimeout()) * time.Second

		for _, conn := range conns.All(app.ID()) {
			if time.Since(conn.LastSeenAt()) > timeout {
				_ = conn.Send(ctx, ping)
			}
		}
	}
}

// PruneStaleConnections closes and removes connections that did not respond
// to a ping within their application's PingInterval. A connection is
// considered stale when its pongedAt time is before its lastSeenAt time AND
// the connection has been inactive for longer than PingInterval — meaning a
// ping was sent (lastSeenAt advanced) but no pong was ever received
// (pongedAt did not advance past lastSeenAt).
func PruneStaleConnections(ctx context.Context, conns *ConnectionManager, apps *AppManager) {
	for _, app := range apps.All() {
		pingInterval := time.Duration(app.PingInterval()) * time.Second

		for _, conn := range conns.All(app.ID()) {
			lastSeen := conn.LastSeenAt()

			// pongedAt is zero on a brand-new connection; treat it the same as
			// "never ponged" — only prune if the connection is also past the
			// ping interval with no pong recorded after the last seen time.
			if conn.pongedAt.Before(lastSeen) && time.Since(lastSeen) > pingInterval {
				_ = conn.Close(ctx, CodeConnectionStale, "connection stale")
				conns.Remove(conn.AppID(), conn.SocketID())
			}
		}
	}
}

// StartJobLoop runs PingInactiveConnections followed by PruneStaleConnections
// on every tick of interval. It blocks until ctx is cancelled.
func StartJobLoop(ctx context.Context, conns *ConnectionManager, apps *AppManager, interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			PingInactiveConnections(ctx, conns, apps)
			PruneStaleConnections(ctx, conns, apps)
		}
	}
}
